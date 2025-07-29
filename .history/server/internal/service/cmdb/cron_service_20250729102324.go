package cmdb

import (
	repo "api-server/internal/repository/cmdb"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// CronService 定时任务服务
// 负责调度和管理主机相关的定时任务
type CronService struct {
	cron         *cron.Cron
	hostRepo     *repo.HostRepository
	providerRepo *repo.ProviderRepository
	alertService *AlertService
	hostService  *HostService
	jobs         map[string]cron.EntryID
	mu           sync.Mutex
}

// NewCronService 创建定时任务服务实例
func NewCronService(
	hostRepo *repo.HostRepository,
	providerRepo *repo.ProviderRepository,
	alertService *AlertService,
	hostService *HostService,
) *CronService {
	cronSvc := &CronService{
		cron:         cron.New(cron.WithSeconds()),
		hostRepo:     hostRepo,
		providerRepo: providerRepo,
		alertService: alertService,
		hostService:  hostService,
		jobs:         make(map[string]cron.EntryID),
	}

	return cronSvc
}

// Start 启动定时任务服务
func (s *CronService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 注册默认任务
	if err := s.registerDefaultJobs(); err != nil {
		return err
	}

	// 启动cron服务
	s.cron.Start()
	log.Println("CronService started successfully")

	return nil
}

// Stop 停止定时任务服务
func (s *CronService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Println("🔍 [CronService] Starting graceful shutdown...")

	if s.cron != nil {
		log.Println("🔍 [CronService] Stopping cron scheduler...")
		ctx := s.cron.Stop()
		log.Println("🔍 [CronService] Waiting for running jobs to complete...")
		<-ctx.Done() // 等待任务完成
		log.Println("🔍 [CronService] All jobs completed")
	}

	log.Println("✅ [CronService] Graceful shutdown completed")
}

// registerDefaultJobs 注册默认任务
func (s *CronService) registerDefaultJobs() error {
	// 同步所有云资源的任务 (每天凌晨2点执行)
	if err := s.AddJob("sync_all_cloud_resources", "0 0 2 * * ?", func() {
		log.Println("Running scheduled job: sync_all_cloud_resources")
		if err := s.hostService.SyncAllProviderHosts(); err != nil {
			log.Printf("Error syncing all provider hosts: %v", err)
		} else {
			log.Println("Successfully synced all provider hosts")
		}
	}); err != nil {
		return err
	}

	// 检查主机状态的任务 (每5分钟执行一次)
	if err := s.AddJob("check_host_status", "0 */5 * * * ?", func() {
		log.Println("Running scheduled job: check_host_status")
		if err := s.hostService.AlertHostStatus(); err != nil {
			log.Printf("Error checking host status: %v", err)
		} else {
			log.Println("Successfully checked host status")
		}
	}); err != nil {
		return err
	}

	// 检查即将过期主机的任务 (每天早上9点执行)
	if err := s.AddJob("check_expiring_hosts", "0 0 9 * * ?", func() {
		log.Println("Running scheduled job: check_expiring_hosts")
		alerts, err := s.alertService.CheckExpiringHosts(30) // 检查30天内过期的主机
		if err != nil {
			log.Printf("Error checking expiring hosts: %v", err)
		} else {
			log.Printf("Successfully checked expiring hosts, found %d alerts", len(alerts))
		}
	}); err != nil {
		return err
	}

	// 运行所有告警检查的任务 (每小时执行一次)
	if err := s.AddJob("run_all_alert_checks", "0 0 * * * ?", func() {
		log.Println("Running scheduled job: run_all_alert_checks")
		alerts, err := s.alertService.RunAllChecks()
		if err != nil {
			log.Printf("Error running all alert checks: %v", err)
		} else {
			log.Printf("Successfully ran all alert checks, found %d alerts", len(alerts))
		}
	}); err != nil {
		return err
	}

	return nil
}

// AddJob 添加定时任务
func (s *CronService) AddJob(name string, spec string, job func()) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果任务已存在，先移除
	if entryID, exists := s.jobs[name]; exists {
		s.cron.Remove(entryID)
		delete(s.jobs, name)
	}

	// 添加新任务
	entryID, err := s.cron.AddFunc(spec, job)
	if err != nil {
		return fmt.Errorf("添加定时任务 '%s' 失败: %w", name, err)
	}

	// 记录任务
	s.jobs[name] = entryID
	return nil
}

// RemoveJob 移除定时任务
func (s *CronService) RemoveJob(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryID, exists := s.jobs[name]; exists {
		s.cron.Remove(entryID)
		delete(s.jobs, name)
	}
}

// ListJobs 列出所有定时任务
func (s *CronService) ListJobs() map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make(map[string]string)
	entries := s.cron.Entries()

	for name, entryID := range s.jobs {
		for _, entry := range entries {
			if entry.ID == entryID {
				result[name] = entry.Schedule.Next(time.Now()).Format(time.RFC3339)
				break
			}
		}
	}

	return result
}
