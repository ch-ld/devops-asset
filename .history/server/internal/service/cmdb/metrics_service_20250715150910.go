package cmdb

import (
	"api-server/internal/model/cmdb"
	repo "api-server/internal/repository/cmdb"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/datatypes"
)

// MetricsService 主机指标监控服务
// 负责收集和统计主机的CPU、内存、磁盘等指标
type MetricsService struct {
	hostRepo    *repo.HostRepository
	metricsChan chan HostMetrics
	collectors  map[uint]*hostCollector
	mu          sync.RWMutex
}

// HostMetrics 主机指标数据
type HostMetrics struct {
	HostID      uint       `json:"host_id"`
	Timestamp   time.Time  `json:"timestamp"`
	CPUUsage    float64    `json:"cpu_usage"`    // CPU使用率 (百分比)
	MemoryUsage float64    `json:"memory_usage"` // 内存使用率 (百分比)
	DiskUsage   float64    `json:"disk_usage"`   // 磁盘使用率 (百分比)
	LoadAvg     [3]float64 `json:"load_avg"`     // 1分钟、5分钟、15分钟系统负载
	NetIO       struct {
		RxBytes int64 `json:"rx_bytes"` // 网络接收字节数/秒
		TxBytes int64 `json:"tx_bytes"` // 网络发送字节数/秒
	} `json:"net_io"`
	ProcessCount int `json:"process_count"` // 进程数量
}

// hostCollector 单主机指标收集器
type hostCollector struct {
	hostID       uint
	host         *cmdb.Host
	metricsChan  chan<- HostMetrics
	stopChan     chan struct{}
	interval     time.Duration
	lastNetStats struct {
		timestamp time.Time
		rxBytes   int64
		txBytes   int64
	}
}

// NewMetricsService 创建主机指标监控服务实例
func NewMetricsService(hostRepo *repo.HostRepository) *MetricsService {
	return &MetricsService{
		hostRepo:    hostRepo,
		metricsChan: make(chan HostMetrics, 100),
		collectors:  make(map[uint]*hostCollector),
	}
}

// Start 启动指标监控服务
func (s *MetricsService) Start() error {
	// 启动指标处理协程
	go s.processMetrics()

	// 初始化已存在主机的指标收集器
	hosts, err := s.hostRepo.FindByStatus("running")
	if err != nil {
		return fmt.Errorf("获取运行中主机失败: %w", err)
	}

	for _, host := range hosts {
		s.StartHostMonitoring(host.ID)
	}

	return nil
}

// Stop 停止指标监控服务
func (s *MetricsService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 停止所有收集器
	for _, collector := range s.collectors {
		close(collector.stopChan)
	}
	s.collectors = make(map[uint]*hostCollector)

	// 关闭指标通道
	close(s.metricsChan)
}

// StartHostMonitoring 开始监控指定主机
func (s *MetricsService) StartHostMonitoring(hostID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果已经在监控中，先停止
	if collector, exists := s.collectors[hostID]; exists {
		close(collector.stopChan)
		delete(s.collectors, hostID)
	}

	// 获取主机信息
	host, err := s.hostRepo.FindByID(hostID)
	if err != nil {
		return fmt.Errorf("获取主机信息失败: %w", err)
	}

	// 创建收集器
	collector := &hostCollector{
		hostID:      hostID,
		host:        host,
		metricsChan: s.metricsChan,
		stopChan:    make(chan struct{}),
		interval:    time.Minute, // 默认收集间隔1分钟
	}

	// 存储收集器
	s.collectors[hostID] = collector

	// 启动收集协程
	go s.runCollector(collector)

	return nil
}

// StopHostMonitoring 停止监控指定主机
func (s *MetricsService) StopHostMonitoring(hostID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if collector, exists := s.collectors[hostID]; exists {
		close(collector.stopChan)
		delete(s.collectors, hostID)
	}
}

// runCollector 运行主机指标收集器
func (s *MetricsService) runCollector(collector *hostCollector) {
	ticker := time.NewTicker(collector.interval)
	defer ticker.Stop()

	// 立即收集一次指标
	s.collectHostMetrics(collector)

	for {
		select {
		case <-collector.stopChan:
			return
		case <-ticker.C:
			s.collectHostMetrics(collector)
		}
	}
}

// collectHostMetrics 收集主机指标
func (s *MetricsService) collectHostMetrics(collector *hostCollector) {
	host := collector.host

	// 跳过没有SSH访问权限的主机
	var ipList []string
	if err := json.Unmarshal(host.PublicIP, &ipList); err != nil || len(ipList) == 0 {
		return
	}

	if host.Username == "" || host.Password == "" {
		return
	}

	// 构建SSH命令
	cmd := `
		echo "===CPU==="
		top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk '{print 100 - $1}'
		echo "===MEMORY==="
		free | grep Mem | awk '{print $3/$2 * 100.0}'
		echo "===DISK==="
		df -h / | grep -v Filesystem | awk '{print $5}' | tr -d '%'
		echo "===LOAD==="
		cat /proc/loadavg
		echo "===NET==="
		cat /proc/net/dev | grep -E "eth0|ens|enp" | awk '{print $2, $10}'
		echo "===PROCESS==="
		ps -ef | wc -l
	`

	// 执行SSH命令
	output, err := s.execSSHCommand(host, cmd)
	if err != nil {
		log.Printf("收集主机 %d 指标失败: %v", host.ID, err)
		return
	}

	// 解析输出
	metrics, err := s.parseMetricsOutput(output, collector)
	if err != nil {
		log.Printf("解析主机 %d 指标失败: %v", host.ID, err)
		return
	}

	// 发送指标数据
	metrics.HostID = host.ID
	metrics.Timestamp = time.Now()
	collector.metricsChan <- metrics
}

// execSSHCommand 执行SSH命令
func (s *MetricsService) execSSHCommand(host *cmdb.Host, cmd string) (string, error) {
	// 示例代码，实际应使用SSH库执行命令
	// 这里简化实现，假设已经有execSSHCommand函数
	return "===CPU===\n25.5\n===MEMORY===\n65.2\n===DISK===\n78\n===LOAD===\n0.52 0.58 0.65 2/352 24685\n===NET===\n1024000 512000\n===PROCESS===\n243", nil
}

// parseMetricsOutput 解析指标输出
func (s *MetricsService) parseMetricsOutput(output string, collector *hostCollector) (HostMetrics, error) {
	var metrics HostMetrics
	lines := strings.Split(output, "\n")
	var section string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "===") && strings.HasSuffix(line, "===") {
			section = line
			continue
		}

		switch section {
		case "===CPU===":
			cpu, err := strconv.ParseFloat(line, 64)
			if err != nil {
				return metrics, fmt.Errorf("解析CPU使用率失败: %w", err)
			}
			metrics.CPUUsage = cpu

		case "===MEMORY===":
			mem, err := strconv.ParseFloat(line, 64)
			if err != nil {
				return metrics, fmt.Errorf("解析内存使用率失败: %w", err)
			}
			metrics.MemoryUsage = mem

		case "===DISK===":
			disk, err := strconv.ParseFloat(line, 64)
			if err != nil {
				return metrics, fmt.Errorf("解析磁盘使用率失败: %w", err)
			}
			metrics.DiskUsage = disk

		case "===LOAD===":
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				for i := 0; i < 3; i++ {
					load, err := strconv.ParseFloat(parts[i], 64)
					if err != nil {
						continue
					}
					metrics.LoadAvg[i] = load
				}
			}

		case "===NET===":
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				rx, _ := strconv.ParseInt(parts[0], 10, 64)
				tx, _ := strconv.ParseInt(parts[1], 10, 64)

				// 计算网络速率
				now := time.Now()
				if !collector.lastNetStats.timestamp.IsZero() {
					duration := now.Sub(collector.lastNetStats.timestamp).Seconds()
					if duration > 0 {
						metrics.NetIO.RxBytes = int64(float64(rx-collector.lastNetStats.rxBytes) / duration)
						metrics.NetIO.TxBytes = int64(float64(tx-collector.lastNetStats.txBytes) / duration)
					}
				}

				// 更新上次统计
				collector.lastNetStats.timestamp = now
				collector.lastNetStats.rxBytes = rx
				collector.lastNetStats.txBytes = tx
			}

		case "===PROCESS===":
			procs, err := strconv.Atoi(line)
			if err != nil {
				continue
			}
			metrics.ProcessCount = procs
		}
	}

	return metrics, nil
}

// 完善processMetrics方法，添加对指标数据的存储逻辑
func (s *MetricsService) processMetrics() {
	for metrics := range s.metricsChan {
		// 保存指标数据到数据库或时序数据库
		if err := s.saveMetricsToStorage(metrics); err != nil {
			log.Printf("保存主机 %d 指标数据失败: %v", metrics.HostID, err)
		}

		// 检查异常值并触发告警
		s.checkMetricsThresholds(metrics)
	}
}

// saveMetricsToStorage 保存指标数据到存储
func (s *MetricsService) saveMetricsToStorage(metrics HostMetrics) error {
	// 实现保存逻辑
	// 可以保存到关系型数据库或时序数据库
	// 这里简化实现，仅打印日志
	metricsJSON, _ := json.Marshal(metrics)
	log.Printf("保存主机 %d 的指标数据: %s", metrics.HostID, string(metricsJSON))

	// TODO: 实际生产环境应该将数据保存到数据库
	// 例如使用influxdb、prometheus等时序数据库

	return nil
}

// checkMetricsThresholds 检查指标阈值并触发告警
func (s *MetricsService) checkMetricsThresholds(metrics HostMetrics) {
	// CPU 使用率告警阈值
	if metrics.CPUUsage > 90 {
		log.Printf("警告: 主机 %d 的CPU使用率过高: %.2f%%", metrics.HostID, metrics.CPUUsage)
		// TODO: 触发告警通知
	}

	// 内存使用率告警阈值
	if metrics.MemoryUsage > 90 {
		log.Printf("警告: 主机 %d 的内存使用率过高: %.2f%%", metrics.HostID, metrics.MemoryUsage)
		// TODO: 触发告警通知
	}

	// 磁盘使用率告警阈值
	if metrics.DiskUsage > 90 {
		log.Printf("警告: 主机 %d 的磁盘使用率过高: %.2f%%", metrics.HostID, metrics.DiskUsage)
		// TODO: 触发告警通知
	}

	// 系统负载告警阈值
	if metrics.LoadAvg[0] > float64(5) { // 假设单核CPU
		log.Printf("警告: 主机 %d 的系统负载过高: %.2f", metrics.HostID, metrics.LoadAvg[0])
		// TODO: 触发告警通知
	}
}

// GetLatestMetrics 获取指定主机的最新指标
func (s *MetricsService) GetLatestMetrics(hostID uint) (*HostMetrics, error) {
	// 实际实现中，应该从存储中读取最新指标
	// 这里简化实现，直接触发一次收集并返回
	s.mu.RLock()
	collector, exists := s.collectors[hostID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("主机 %d 未在监控中", hostID)
	}

	metricsChan := make(chan HostMetrics, 1)
	collector.metricsChan = metricsChan

	s.collectHostMetrics(collector)

	select {
	case metrics := <-metricsChan:
		return &metrics, nil
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("获取主机 %d 指标超时", hostID)
	}
}

// GetHostMetricsHistory 获取主机指标历史数据
func (s *MetricsService) GetHostMetricsHistory(hostID uint, startTime, endTime time.Time, metricType string) ([]HostMetrics, error) {
	// 这里应该从存储中查询历史数据
	// 简化实现，生成模拟数据
	history := []HostMetrics{}

	// 计算时间点数量，每小时一个点
	hours := int(endTime.Sub(startTime).Hours()) + 1
	if hours <= 0 {
		return history, nil
	}

	// 限制最大返回点数
	if hours > 1000 {
		hours = 1000
	}

	// 生成模拟数据
	for i := 0; i < hours; i++ {
		timestamp := startTime.Add(time.Duration(i) * time.Hour)
		metrics := HostMetrics{
			HostID:      hostID,
			Timestamp:   timestamp,
			CPUUsage:    10 + rand.Float64()*50, // 10-60%
			MemoryUsage: 20 + rand.Float64()*60, // 20-80%
			DiskUsage:   30 + rand.Float64()*50, // 30-80%
			LoadAvg:     [3]float64{rand.Float64() * 3, rand.Float64() * 2, rand.Float64() * 1.5},
		}

		// 网络IO
		metrics.NetIO.RxBytes = rand.Int63n(1024 * 1024 * 10) // 0-10MB/s
		metrics.NetIO.TxBytes = rand.Int63n(1024 * 1024 * 5)  // 0-5MB/s

		// 进程数
		metrics.ProcessCount = 100 + rand.Intn(200) // 100-300个进程

		// 根据metricType过滤
		if metricType != "all" {
			// 保持所有字段，但实际可以根据metricType返回不同的数据结构
			history = append(history, metrics)
		} else {
			history = append(history, metrics)
		}
	}

	return history, nil
}

// GetOverallMetrics 获取所有主机的汇总指标
func (s *MetricsService) GetOverallMetrics() (map[string]interface{}, error) {
	// 获取所有主机
	hosts, err := s.hostRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("获取主机列表失败: %w", err)
	}

	// 准备结果
	result := map[string]interface{}{
		"host_count": len(hosts),
		"metrics": map[string]interface{}{
			"cpu_usage_avg":    0.0,
			"memory_usage_avg": 0.0,
			"disk_usage_avg":   0.0,
			"top_cpu_hosts":    []map[string]interface{}{},
			"top_memory_hosts": []map[string]interface{}{},
		},
	}

	// 统计汇总信息
	var totalCPUUsage, totalMemoryUsage, totalDiskUsage float64
	var hostMetricsCount int

	type hostMetricItem struct {
		HostID uint
		Name   string
		Value  float64
	}

	var cpuUsageList []hostMetricItem
	var memoryUsageList []hostMetricItem

	// 遍历所有主机，收集最新指标
	for _, host := range hosts {
		// 简化实现，随机生成数据
		// 实际应从存储中查询最新指标
		cpuUsage := 10 + rand.Float64()*60
		memoryUsage := 20 + rand.Float64()*70
		diskUsage := 30 + rand.Float64()*60

		totalCPUUsage += cpuUsage
		totalMemoryUsage += memoryUsage
		totalDiskUsage += diskUsage
		hostMetricsCount++

		cpuUsageList = append(cpuUsageList, hostMetricItem{
			HostID: host.ID,
			Name:   host.Name,
			Value:  cpuUsage,
		})

		memoryUsageList = append(memoryUsageList, hostMetricItem{
			HostID: host.ID,
			Name:   host.Name,
			Value:  memoryUsage,
		})
	}

	// 计算平均值
	if hostMetricsCount > 0 {
		result["metrics"].(map[string]interface{})["cpu_usage_avg"] = totalCPUUsage / float64(hostMetricsCount)
		result["metrics"].(map[string]interface{})["memory_usage_avg"] = totalMemoryUsage / float64(hostMetricsCount)
		result["metrics"].(map[string]interface{})["disk_usage_avg"] = totalDiskUsage / float64(hostMetricsCount)
	}

	// 排序CPU使用率，获取前5名
	sort.Slice(cpuUsageList, func(i, j int) bool {
		return cpuUsageList[i].Value > cpuUsageList[j].Value
	})

	topCPUHosts := []map[string]interface{}{}
	for i := 0; i < 5 && i < len(cpuUsageList); i++ {
		topCPUHosts = append(topCPUHosts, map[string]interface{}{
			"host_id": cpuUsageList[i].HostID,
			"name":    cpuUsageList[i].Name,
			"value":   cpuUsageList[i].Value,
		})
	}
	result["metrics"].(map[string]interface{})["top_cpu_hosts"] = topCPUHosts

	// 排序内存使用率，获取前5名
	sort.Slice(memoryUsageList, func(i, j int) bool {
		return memoryUsageList[i].Value > memoryUsageList[j].Value
	})

	topMemoryHosts := []map[string]interface{}{}
	for i := 0; i < 5 && i < len(memoryUsageList); i++ {
		topMemoryHosts = append(topMemoryHosts, map[string]interface{}{
			"host_id": memoryUsageList[i].HostID,
			"name":    memoryUsageList[i].Name,
			"value":   memoryUsageList[i].Value,
		})
	}
	result["metrics"].(map[string]interface{})["top_memory_hosts"] = topMemoryHosts

	return result, nil
}

// extractFirstIP 从JSON字符串中提取第一个IP地址
func extractFirstIP(jsonStr datatypes.JSON) string {
	var ips []string
	if err := json.Unmarshal([]byte(jsonStr), &ips); err != nil || len(ips) == 0 {
		return ""
	}
	return ips[0]
}
