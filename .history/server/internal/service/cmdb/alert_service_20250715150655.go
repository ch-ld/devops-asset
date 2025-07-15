package cmdb

import (
	"api-server/internal/model/cmdb"
	repo "api-server/internal/repository/cmdb"
	"fmt"
	"time"
)

// AlertService 主机告警服务
// 负责检查主机过期、异常状态等情况并产生告警
type AlertService struct {
	hostRepo       *repo.HostRepository
	providerRepo   *repo.ProviderRepository
	notificationCh chan HostAlert
}

// HostAlert 主机告警结构
type HostAlert struct {
	Host      *cmdb.Host `json:"host"`
	AlertType string     `json:"alert_type"` // expired, expiring, error, abnormal
	Message   string     `json:"message"`
	Time      time.Time  `json:"time"`
}

// NewAlertService 创建告警服务实例
func NewAlertService(hostRepo *repo.HostRepository, providerRepo *repo.ProviderRepository) *AlertService {
	return &AlertService{
		hostRepo:       hostRepo,
		providerRepo:   providerRepo,
		notificationCh: make(chan HostAlert, 100), // 缓冲通道，避免阻塞
	}
}

// CheckExpiringHosts 检查即将过期的主机
func (s *AlertService) CheckExpiringHosts(days int) ([]HostAlert, error) {
	var alerts []HostAlert

	// 获取所有有过期时间的主机
	hosts, err := s.hostRepo.FindByExpiring(days)
	if err != nil {
		return nil, fmt.Errorf("查询即将过期主机失败: %w", err)
	}

	now := time.Now()
	for _, host := range hosts {
		if host.ExpiredAt == nil {
			continue
		}

		// 计算剩余天数
		remainingDays := int(host.ExpiredAt.Sub(now).Hours() / 24)

		// 已过期
		if now.After(*host.ExpiredAt) {
			alert := HostAlert{
				Host:      &host,
				AlertType: "expired",
				Message:   fmt.Sprintf("主机 %s 已于 %s 过期", host.Name, host.ExpiredAt.Format("2006-01-02")),
				Time:      now,
			}
			alerts = append(alerts, alert)
			s.notificationCh <- alert
		} else if remainingDays <= days {
			// 即将过期
			alert := HostAlert{
				Host:      &host,
				AlertType: "expiring",
				Message:   fmt.Sprintf("主机 %s 将于 %s 过期，剩余 %d 天", host.Name, host.ExpiredAt.Format("2006-01-02"), remainingDays),
				Time:      now,
			}
			alerts = append(alerts, alert)
			s.notificationCh <- alert
		}
	}

	return alerts, nil
}

// CheckAbnormalHosts 检查状态异常的主机
func (s *AlertService) CheckAbnormalHosts() ([]HostAlert, error) {
	var alerts []HostAlert

	// 获取所有状态为异常或错误的主机
	hosts, err := s.hostRepo.FindByAbnormalStatus()
	if err != nil {
		return nil, fmt.Errorf("查询状态异常主机失败: %w", err)
	}

	now := time.Now()
	for _, host := range hosts {
		alert := HostAlert{
			Host:      &host,
			AlertType: "abnormal",
			Message:   fmt.Sprintf("主机 %s 状态异常: %s", host.Name, host.Status),
			Time:      now,
		}
		alerts = append(alerts, alert)
		s.notificationCh <- alert
	}

	return alerts, nil
}

// VerifyHostConnectivity 检查主机连接性
func (s *AlertService) VerifyHostConnectivity() ([]HostAlert, error) {
	var alerts []HostAlert

	// 获取所有状态为运行中的主机
	hosts, err := s.hostRepo.FindByStatus("running")
	if err != nil {
		return nil, fmt.Errorf("查询运行中主机失败: %w", err)
	}

	now := time.Now()
	for _, host := range hosts {
		// 尝试通过SSH连接验证主机是否可达
		// 简化实现，实际可能需要使用SSH库进行连接测试
		isReachable := testHostReachable(&host)

		if !isReachable {
			alert := HostAlert{
				Host:      &host,
				AlertType: "unreachable",
				Message:   fmt.Sprintf("主机 %s 无法连接，但状态显示为运行中", host.Name),
				Time:      now,
			}
			alerts = append(alerts, alert)
			s.notificationCh <- alert
		}
	}

	return alerts, nil
}

// RunAllChecks 运行所有检查
func (s *AlertService) RunAllChecks() ([]HostAlert, error) {
	var allAlerts []HostAlert

	// 检查即将过期的主机（30天内）
	expiringAlerts, err := s.CheckExpiringHosts(30)
	if err != nil {
		return nil, err
	}
	allAlerts = append(allAlerts, expiringAlerts...)

	// 检查状态异常的主机
	abnormalAlerts, err := s.CheckAbnormalHosts()
	if err != nil {
		return nil, err
	}
	allAlerts = append(allAlerts, abnormalAlerts...)

	// 检查主机连接性
	// 在实际生产环境中可能会消耗较多资源，可选择性启用
	/*
		connectivityAlerts, err := s.VerifyHostConnectivity()
		if err != nil {
			return nil, err
		}
		allAlerts = append(allAlerts, connectivityAlerts...)
	*/

	return allAlerts, nil
}

// GetNotificationChannel 获取告警通知通道
func (s *AlertService) GetNotificationChannel() <-chan HostAlert {
	return s.notificationCh
}

// testHostReachable 测试主机是否可达
// 实际实现可能需要使用SSH或PING等方式验证连接性
func testHostReachable(host *cmdb.Host) bool {
	// 简化实现，实际应使用SSH库进行连接测试
	// 这里默认返回true，表示可连接
	return true
}
