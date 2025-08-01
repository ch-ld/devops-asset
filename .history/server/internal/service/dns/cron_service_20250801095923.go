package dns

import (
	"context"
	"fmt"
	"time"

	"api-server/internal/model/dns"
	dnsprovider "api-server/internal/provider/dns"
	dnsrepo "api-server/internal/repository/dns"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CertificateServiceInterface 证书服务接口（占位符）
type CertificateServiceInterface interface {
	// RenewCertificate 续期证书
	RenewCertificate(ctx context.Context, certificateID uint, actorID uint, clientIP string) (*dns.Certificate, error)
}

// CronService DNS定时任务服务
type CronService struct {
	db                 *gorm.DB
	cron               *cron.Cron
	domainRepo         *dnsrepo.DomainRepository
	certificateRepo    *dnsrepo.CertificateRepository
	recordRepo         *dnsrepo.RecordRepository
	providerRepo       *dnsrepo.ProviderRepository
	domainService      *DomainService
	certificateService CertificateServiceInterface
	recordService      *RecordService
	providerService    *ProviderService
}

// NewCronService 创建DNS定时任务服务
func NewCronService(
	db *gorm.DB,
	domainRepo *dnsrepo.DomainRepository,
	certificateRepo *dnsrepo.CertificateRepository,
	recordRepo *dnsrepo.RecordRepository,
	providerRepo *dnsrepo.ProviderRepository,
	domainService *DomainService,
	certificateService CertificateServiceInterface,
	recordService *RecordService,
	providerService *ProviderService,
) *CronService {
	return &CronService{
		db:                 db,
		cron:               cron.New(cron.WithSeconds()),
		domainRepo:         domainRepo,
		certificateRepo:    certificateRepo, // 可以为nil
		recordRepo:         recordRepo,
		providerRepo:       providerRepo,
		domainService:      domainService,
		certificateService: certificateService, // 可以为nil
		recordService:      recordService,
		providerService:    providerService,
	}
}

// Start 启动定时任务
func (s *CronService) Start() error {
	zap.L().Info("Starting DNS cron service...")

	// 注册定时任务
	if err := s.registerJobs(); err != nil {
		return fmt.Errorf("failed to register cron jobs: %w", err)
	}

	// 启动定时任务调度器
	s.cron.Start()

	zap.L().Info("DNS cron service started successfully")
	return nil
}

// Stop 停止定时任务
func (s *CronService) Stop() {
	zap.L().Info("Stopping DNS cron service...")
	s.cron.Stop()
	zap.L().Info("DNS cron service stopped")
}

// registerJobs 注册定时任务
func (s *CronService) registerJobs() error {
	// 域名到期提醒任务 (每天上午9点执行)
	if _, err := s.cron.AddFunc("0 0 9 * * *", s.domainExpiryNotifier); err != nil {
		return fmt.Errorf("failed to add domain expiry notifier job: %w", err)
	}

	// 证书到期提醒任务 (每天上午9点执行) - 仅在证书服务可用时注册
	if s.certificateRepo != nil && s.certificateService != nil {
		if _, err := s.cron.AddFunc("0 0 9 * * *", s.certificateExpiryNotifier); err != nil {
			return fmt.Errorf("failed to add certificate expiry notifier job: %w", err)
		}
	}

	// DNS区域同步任务 (每5分钟执行一次)
	if _, err := s.cron.AddFunc("0 */5 * * * *", s.dnsZoneSync); err != nil {
		return fmt.Errorf("failed to add DNS zone sync job: %w", err)
	}

	// 证书自动续期任务 (每天凌晨2点执行) - 仅在证书服务可用时注册
	if s.certificateRepo != nil && s.certificateService != nil {
		if _, err := s.cron.AddFunc("0 0 2 * * *", s.certificateAutoRenew); err != nil {
			return fmt.Errorf("failed to add certificate auto renew job: %w", err)
		}
	}

	// DNS提供商健康检查任务 (每小时执行一次)
	if _, err := s.cron.AddFunc("0 0 * * * *", s.providerHealthCheck); err != nil {
		return fmt.Errorf("failed to add provider health check job: %w", err)
	}

	zap.L().Info("DNS cron jobs registered successfully")
	return nil
}

// domainExpiryNotifier 域名到期提醒任务
func (s *CronService) domainExpiryNotifier() {
	zap.L().Info("Running DNS cron job: domain_expiry_notifier")

	ctx := context.Background()

	// 检查30天内到期的域名
	domains30, err := s.domainRepo.FindExpiring(30)
	if err != nil {
		zap.L().Error("Failed to find domains expiring in 30 days", zap.Error(err))
		return
	}

	// 检查7天内到期的域名
	domains7, err := s.domainRepo.FindExpiring(7)
	if err != nil {
		zap.L().Error("Failed to find domains expiring in 7 days", zap.Error(err))
		return
	}

	// 检查1天内到期的域名
	domains1, err := s.domainRepo.FindExpiring(1)
	if err != nil {
		zap.L().Error("Failed to find domains expiring in 1 day", zap.Error(err))
		return
	}

	// 发送通知
	s.sendDomainExpiryNotifications(ctx, domains30, 30)
	s.sendDomainExpiryNotifications(ctx, domains7, 7)
	s.sendDomainExpiryNotifications(ctx, domains1, 1)

	zap.L().Info("Domain expiry notifier job completed",
		zap.Int("domains_30d", len(domains30)),
		zap.Int("domains_7d", len(domains7)),
		zap.Int("domains_1d", len(domains1)))
}

// certificateExpiryNotifier 证书到期提醒任务
func (s *CronService) certificateExpiryNotifier() {
	zap.L().Info("Running DNS cron job: certificate_expiry_notifier")

	ctx := context.Background()

	// 检查30天内到期的证书
	certs30, err := s.certificateRepo.FindExpiring(30)
	if err != nil {
		zap.L().Error("Failed to find certificates expiring in 30 days", zap.Error(err))
		return
	}

	// 检查7天内到期的证书
	certs7, err := s.certificateRepo.FindExpiring(7)
	if err != nil {
		zap.L().Error("Failed to find certificates expiring in 7 days", zap.Error(err))
		return
	}

	// 检查1天内到期的证书
	certs1, err := s.certificateRepo.FindExpiring(1)
	if err != nil {
		zap.L().Error("Failed to find certificates expiring in 1 day", zap.Error(err))
		return
	}

	// 发送通知
	s.sendCertificateExpiryNotifications(ctx, certs30, 30)
	s.sendCertificateExpiryNotifications(ctx, certs7, 7)
	s.sendCertificateExpiryNotifications(ctx, certs1, 1)

	zap.L().Info("Certificate expiry notifier job completed",
		zap.Int("certificates_30d", len(certs30)),
		zap.Int("certificates_7d", len(certs7)),
		zap.Int("certificates_1d", len(certs1)))
}

// dnsZoneSync DNS区域同步任务
func (s *CronService) dnsZoneSync() {
	zap.L().Info("Running DNS cron job: dns_zone_sync")

	ctx := context.Background()

	// 获取所有活跃的DNS提供商
	providers, err := s.providerRepo.FindByStatus("active")
	if err != nil {
		zap.L().Error("Failed to find active DNS providers", zap.Error(err))
		return
	}

	syncCount := 0
	errorCount := 0

	for _, provider := range providers {
		// 获取该提供商下的所有域名
		// TODO: 实现FindByProviderID方法
		domains := []*dns.Domain{} // 临时占位符
		_ = provider.ID            // 避免未使用变量警告

		// 同步每个域名的DNS记录
		for _, domain := range domains {
			if err := s.syncDomainRecords(ctx, domain, provider); err != nil {
				zap.L().Error("Failed to sync domain records",
					zap.String("domain", domain.Name),
					zap.String("provider", provider.Name),
					zap.Error(err))
				errorCount++
			} else {
				syncCount++
			}
		}
	}

	zap.L().Info("DNS zone sync job completed",
		zap.Int("synced_domains", syncCount),
		zap.Int("errors", errorCount))
}

// certificateAutoRenew 证书自动续期任务
func (s *CronService) certificateAutoRenew() {
	zap.L().Info("Running DNS cron job: certificate_auto_renew")

	ctx := context.Background()

	// 查找需要自动续期的证书
	certificates, err := s.certificateRepo.FindAutoRenew()
	if err != nil {
		zap.L().Error("Failed to find certificates for auto renewal", zap.Error(err))
		return
	}

	renewCount := 0
	errorCount := 0

	for _, cert := range certificates {
		if err := s.renewCertificate(ctx, cert); err != nil {
			zap.L().Error("Failed to auto renew certificate",
				zap.String("domain", cert.CommonName),
				zap.Error(err))
			errorCount++
		} else {
			renewCount++
		}
	}

	zap.L().Info("Certificate auto renew job completed",
		zap.Int("renewed_certificates", renewCount),
		zap.Int("errors", errorCount))
}

// providerHealthCheck DNS提供商健康检查任务
func (s *CronService) providerHealthCheck() {
	zap.L().Info("Running DNS cron job: provider_health_check")

	ctx := context.Background()

	// 获取所有DNS提供商
	providers, err := s.providerRepo.FindAll()
	if err != nil {
		zap.L().Error("Failed to find DNS providers", zap.Error(err))
		return
	}

	healthyCount := 0
	unhealthyCount := 0

	for _, provider := range providers {
		// 测试提供商连接
		result, err := s.providerService.TestProvider(provider.ID)
		if err != nil {
			zap.L().Error("Failed to test provider",
				zap.String("provider", provider.Name),
				zap.Error(err))
			unhealthyCount++
			continue
		}

		// 更新提供商健康状态
		if err := s.updateProviderHealthStatus(ctx, provider, result); err != nil {
			zap.L().Error("Failed to update provider health status",
				zap.String("provider", provider.Name),
				zap.Error(err))
		}

		if result.Success {
			healthyCount++
		} else {
			unhealthyCount++
		}
	}

	zap.L().Info("Provider health check job completed",
		zap.Int("healthy_providers", healthyCount),
		zap.Int("unhealthy_providers", unhealthyCount))
}

// syncDomainRecords 同步域名DNS记录
func (s *CronService) syncDomainRecords(ctx context.Context, domain *dns.Domain, provider *dns.Provider) error {
	zap.L().Info("Syncing domain records",
		zap.String("domain", domain.Name),
		zap.String("provider", provider.Name))

	// 调用RecordService的同步方法
	if s.recordService != nil {
		params := &SyncDomainRecordsParams{
			DomainID:   domain.ID,
			ProviderID: provider.ID,
			DryRun:     false,
			ActorID:    1, // 系统用户ID
			ClientIP:   "127.0.0.1",
		}

		result, err := s.recordService.SyncDomainRecords(params)
		if err != nil {
			zap.L().Error("Failed to sync domain records",
				zap.String("domain", domain.Name),
				zap.Error(err))
			return err
		}

		zap.L().Info("Domain records synced successfully",
			zap.String("domain", domain.Name),
			zap.Int("added", result.ToAdd),
			zap.Int("updated", result.ToUpdate),
			zap.Int("deleted", result.ToDelete),
			zap.Int("errors", len(result.Errors)))

		// 如果有错误，记录日志但不中断
		if len(result.Errors) > 0 {
			for _, errMsg := range result.Errors {
				zap.L().Warn("Sync error", zap.String("error", errMsg))
			}
		}
	}

	return nil
}

// renewCertificate 续期证书
func (s *CronService) renewCertificate(ctx context.Context, certificate *dns.Certificate) error {
	zap.L().Info("Renewing certificate",
		zap.String("common_name", certificate.CommonName),
		zap.Uint("certificate_id", certificate.ID))

	// 这里应该调用CertificateService的续期方法
	// 由于CronService初始化时证书服务可能为nil，需要检查
	if s.certificateService == nil {
		zap.L().Warn("Certificate service not available, skipping renewal")
		return nil
	}

	// 调用证书续期服务
	_, err := s.certificateService.RenewCertificate(ctx, certificate.ID, 1, "127.0.0.1")
	if err != nil {
		zap.L().Error("Failed to renew certificate",
			zap.String("common_name", certificate.CommonName),
			zap.Uint("certificate_id", certificate.ID),
			zap.Error(err))

		// 更新证书状态为续期失败
		certificate.Status = "renew_failed"
		if updateErr := s.certificateRepo.Update(certificate); updateErr != nil {
			zap.L().Error("Failed to update certificate status", zap.Error(updateErr))
		}

		return err
	}

	zap.L().Info("Certificate renewed successfully",
		zap.String("common_name", certificate.CommonName),
		zap.Uint("certificate_id", certificate.ID))

	return nil
}

// updateProviderHealthStatus 更新提供商健康状态
func (s *CronService) updateProviderHealthStatus(ctx context.Context, provider *dns.Provider, result *dnsprovider.TestResult) error {
	zap.L().Info("Updating provider health status",
		zap.String("provider", provider.Name),
		zap.Bool("healthy", result.Success))

	// 更新提供商的健康状态
	var healthStatus string
	var lastError string

	if result.Success {
		healthStatus = "healthy"
		lastError = ""
	} else {
		healthStatus = "unhealthy"
		lastError = result.ErrorMsg
	}

	// 构造更新数据
	updateData := map[string]interface{}{
		"health_status":   healthStatus,
		"last_checked_at": result.TestedAt,
		"response_time":   result.Latency.Milliseconds(),
		"last_error":      lastError,
	}

	// 更新到数据库
	err := s.db.Model(&dns.Provider{}).
		Where("id = ?", provider.ID).
		Updates(updateData).Error

	if err != nil {
		zap.L().Error("Failed to update provider health status",
			zap.String("provider", provider.Name),
			zap.Error(err))
		return err
	}

	zap.L().Info("Provider health status updated",
		zap.String("provider", provider.Name),
		zap.String("status", healthStatus),
		zap.Duration("response_time", result.Latency))

	return nil
}

// sendDomainExpiryNotifications 发送域名到期通知
func (s *CronService) sendDomainExpiryNotifications(ctx context.Context, domains []*dns.Domain, days int) {
	for _, domain := range domains {
		zap.L().Info("Domain expiry notification",
			zap.String("domain", domain.Name),
			zap.Int("days_until_expiry", days),
			zap.Time("expires_at", *domain.ExpiresAt))

		// TODO: 集成通知服务发送邮件/短信/webhook等
		// 这里可以调用统一的通知服务
		s.sendNotification(ctx, &NotificationRequest{
			Type:    "domain_expiry",
			Title:   fmt.Sprintf("域名 %s 即将过期", domain.Name),
			Content: fmt.Sprintf("域名 %s 将在 %d 天后过期，到期时间：%s", domain.Name, days, domain.ExpiresAt.Format("2006-01-02")),
			Level:   s.getNotificationLevel(days),
			Target:  domain.TenantID,
			Metadata: map[string]interface{}{
				"domain_id":   domain.ID,
				"domain_name": domain.Name,
				"expires_at":  domain.ExpiresAt,
				"days_left":   days,
			},
		})
	}
}

// sendCertificateExpiryNotifications 发送证书到期通知
func (s *CronService) sendCertificateExpiryNotifications(ctx context.Context, certificates []*dns.Certificate, days int) {
	for _, cert := range certificates {
		zap.L().Info("Certificate expiry notification",
			zap.String("common_name", cert.CommonName),
			zap.Int("days_until_expiry", days),
			zap.Time("expires_at", *cert.ExpiresAt))

		// 发送证书到期通知
		s.sendNotification(ctx, &NotificationRequest{
			Type:    "certificate_expiry",
			Title:   fmt.Sprintf("SSL证书 %s 即将过期", cert.CommonName),
			Content: fmt.Sprintf("SSL证书 %s 将在 %d 天后过期，到期时间：%s", cert.CommonName, days, cert.ExpiresAt.Format("2006-01-02")),
			Level:   s.getNotificationLevel(days),
			Target:  cert.TenantID,
			Metadata: map[string]interface{}{
				"certificate_id": cert.ID,
				"common_name":    cert.CommonName,
				"expires_at":     cert.ExpiresAt,
				"days_left":      days,
				"auto_renew":     cert.AutoRenew,
			},
		})
	}
}

// NotificationRequest 通知请求
type NotificationRequest struct {
	Type     string                 `json:"type"`
	Title    string                 `json:"title"`
	Content  string                 `json:"content"`
	Level    string                 `json:"level"`
	Target   uint                   `json:"target"` // 租户ID
	Metadata map[string]interface{} `json:"metadata"`
}

// sendNotification 发送通知
func (s *CronService) sendNotification(ctx context.Context, req *NotificationRequest) {
	// TODO: 实现通知发送逻辑
	// 这里可以集成邮件服务、短信服务、webhook等
	zap.L().Info("Sending notification",
		zap.String("type", req.Type),
		zap.String("title", req.Title),
		zap.String("level", req.Level),
		zap.Uint("target", req.Target))

	// 简化实现：记录到日志
	// 实际应该调用通知服务API或发送到消息队列
}

// getNotificationLevel 根据剩余天数获取通知级别
func (s *CronService) getNotificationLevel(days int) string {
	if days <= 1 {
		return "critical"
	} else if days <= 7 {
		return "warning"
	} else if days <= 30 {
		return "info"
	}
	return "low"
}

// cleanupExpiredRecords 清理过期记录
func (s *CronService) cleanupExpiredRecords(ctx context.Context) error {
	zap.L().Info("Starting cleanup of expired records")

	// 清理过期的变更日志（保留90天）
	cutoff := time.Now().AddDate(0, 0, -90)
	result := s.db.Where("created_at < ?", cutoff).Delete(&dns.ChangeLog{})
	if result.Error != nil {
		zap.L().Error("Failed to cleanup expired change logs", zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected > 0 {
		zap.L().Info("Cleaned up expired change logs",
			zap.Int64("rows_deleted", result.RowsAffected))
	}

	// 清理过期的DNS记录（软删除30天后的记录）
	cutoff = time.Now().AddDate(0, 0, -30)
	result = s.db.Unscoped().Where("deleted_at IS NOT NULL AND deleted_at < ?", cutoff).Delete(&dns.Record{})
	if result.Error != nil {
		zap.L().Error("Failed to cleanup expired DNS records", zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected > 0 {
		zap.L().Info("Cleaned up expired DNS records",
			zap.Int64("rows_deleted", result.RowsAffected))
	}

	zap.L().Info("Cleanup of expired records completed")
	return nil
}

// generateStatistics 生成统计报告
func (s *CronService) generateStatistics(ctx context.Context) error {
	zap.L().Info("Generating DNS module statistics")

	// 统计域名数量
	var domainCount int64
	s.db.Model(&dns.Domain{}).Count(&domainCount)

	// 统计证书数量
	var certCount int64
	s.db.Model(&dns.Certificate{}).Count(&certCount)

	// 统计DNS记录数量
	var recordCount int64
	s.db.Model(&dns.Record{}).Count(&recordCount)

	// 统计提供商数量
	var providerCount int64
	s.db.Model(&dns.Provider{}).Count(&providerCount)

	zap.L().Info("DNS module statistics",
		zap.Int64("domains", domainCount),
		zap.Int64("certificates", certCount),
		zap.Int64("records", recordCount),
		zap.Int64("providers", providerCount))

	// TODO: 将统计数据存储到metrics表或发送到监控系统

	return nil
}
