package dns

import (
	"api-server/internal/api/handler/dns"
	"api-server/internal/middleware/middleware"
	repo "api-server/internal/repository/dns"
	svc "api-server/internal/service/dns"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	dnsprovider "api-server/internal/provider/dns"
)

// 定义导出变量以便在main中访问服务
var (
	DomainSvc      *svc.DomainService
	RecordSvc      *svc.RecordService
	ProviderSvc    *svc.ProviderService
	DomainGroupSvc *svc.DomainGroupService
	CronSvc        *svc.CronService
)

// InitDnsRoutes initializes all routes for the DNS module.
func InitDnsRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// 添加调试日志
	fmt.Println("[DNS路由] 开始初始化DNS路由...")

	// 1. Dependency Injection
	// Repositories
	fmt.Println("[DNS路由] 正在创建Repository...")
	domainRepo := repo.NewDomainRepository(db)
	domainGroupRepo := repo.NewDomainGroupRepository(db)
	certificateRepo := repo.NewCertificateRepository(db)
	recordRepo := repo.NewRecordRepository(db)
	providerRepo := repo.NewProviderRepository(db)
	changeLogRepo := repo.NewChangeLogRepository(db)
	httpsMonitorRepo := repo.NewHTTPSMonitorRepository(db)
	_ = repo.NewNotificationRepository(db) // TODO: 后续使用
	fmt.Println("[DNS路由] Repository创建完成")

	// Services
	fmt.Println("[DNS路由] 正在创建Service...")
	DomainSvc = svc.NewDomainService(domainRepo, domainGroupRepo, changeLogRepo, db)
	RecordSvc = svc.NewRecordService(recordRepo, domainRepo, providerRepo, changeLogRepo, db)
	ProviderSvc = svc.NewProviderService(providerRepo, changeLogRepo, domainRepo, db)
	DomainGroupSvc = svc.NewDomainGroupService(domainGroupRepo, changeLogRepo, db)

	// 创建DNS定时任务服务
	CronSvc = svc.NewCronService(
		db,
		domainRepo,
		certificateRepo,
		recordRepo,
		providerRepo,
		DomainSvc,
		nil, // certificateService - TODO: 等证书功能实现后替换
		RecordSvc,
		ProviderSvc,
	)
	fmt.Println("[DNS路由] Service创建完成")

	// Handlers
	fmt.Println("[DNS路由] 正在创建Handler...")
	domainHandler := dns.NewDomainHandler(DomainSvc)
	domainGroupHandler := dns.NewDomainGroupHandler(DomainGroupSvc)
	recordHandler := dns.NewRecordHandler(RecordSvc)
	providerHandler := dns.NewProviderHandler(ProviderSvc, DomainSvc)
	batchHandler := dns.NewBatchHandler(RecordSvc, DomainSvc)

	// 创建证书服务和Handler
	dnsFactory := dnsprovider.GetFactory() // 获取DNS Provider工厂实例
	certificateService := svc.NewCertificateService(
		certificateRepo, domainRepo, providerRepo, recordRepo, dnsFactory,
	)
	certificateHandler := dns.NewCertificateHandler(certificateService)

	// 创建HTTPS监控服务和Handler
	httpsMonitorService := svc.NewHTTPSMonitorService(
		httpsMonitorRepo, certificateRepo, domainRepo,
	)
	httpsMonitorHandler := dns.NewHTTPSMonitorHandler(httpsMonitorService)

	fmt.Println("[DNS路由] Handler创建完成")

	// 2. Route Registration
	fmt.Println("[DNS路由] 正在注册路由...")
	dnsGroup := r.Group("/dns")
	{
		// 域名管理路由
		fmt.Println("[DNS路由] 注册域名管理路由...")
		domainGroup := dnsGroup.Group("/domains")
		domainGroup.Use(middleware.JWTAuth())
		{
			domainGroup.POST("", domainHandler.CreateDomain)       // 创建域名
			domainGroup.GET("", domainHandler.ListDomains)         // 获取域名列表
			domainGroup.GET("/:id", domainHandler.GetDomain)       // 获取域名详情
			domainGroup.PUT("/:id", domainHandler.UpdateDomain)    // 更新域名
			domainGroup.DELETE("/:id", domainHandler.DeleteDomain) // 删除域名

			// 导入导出功能
			domainGroup.POST("/import/csv", domainHandler.ImportDomainsFromCSV)     // 从CSV导入域名
			domainGroup.POST("/import/excel", domainHandler.ImportDomainsFromExcel) // 从Excel导入域名
			domainGroup.GET("/export/csv", domainHandler.ExportDomainsToCSV)        // 导出域名到CSV
			domainGroup.GET("/export/excel", domainHandler.ExportDomainsToExcel)    // 导出域名到Excel
		}

		// DNS统计路由
		fmt.Println("[DNS路由] 注册DNS统计路由...")
		statsGroup := dnsGroup.Group("/stats")
		statsGroup.Use(middleware.JWTAuth())
		{
			statsGroup.GET("/domains", domainHandler.GetDomainStatistics)       // 获取域名统计信息
			statsGroup.GET("/providers", providerHandler.GetProviderStatistics) // 获取提供商统计信息
		}

		// DNS记录管理路由
		fmt.Println("[DNS路由] 注册DNS记录管理路由...")
		recordGroup := dnsGroup.Group("/records")
		recordGroup.Use(middleware.JWTAuth())
		{
			recordGroup.POST("", recordHandler.CreateRecord)       // 创建DNS记录
			recordGroup.GET("", recordHandler.ListRecords)         // 获取DNS记录列表
			recordGroup.GET("/:id", recordHandler.GetRecord)       // 获取DNS记录详情
			recordGroup.PUT("/:id", recordHandler.UpdateRecord)    // 更新DNS记录
			recordGroup.DELETE("/:id", recordHandler.DeleteRecord) // 删除DNS记录
		}

		// DNS提供商管理路由
		fmt.Println("[DNS路由] 注册DNS提供商管理路由...")
		providerGroup := dnsGroup.Group("/providers")
		providerGroup.Use(middleware.JWTAuth())
		{
			providerGroup.POST("", providerHandler.CreateProvider)                          // 创建DNS提供商
			providerGroup.POST("/test-connection", providerHandler.TestConnection)          // 临时测试连接
			providerGroup.GET("", providerHandler.ListProviders)                            // 获取DNS提供商列表
			providerGroup.GET("/:id", providerHandler.GetProvider)                          // 获取DNS提供商详情
			providerGroup.PUT("/:id", providerHandler.UpdateProvider)                       // 更新DNS提供商
			providerGroup.DELETE("/:id", providerHandler.DeleteProvider)                    // 删除DNS提供商
			providerGroup.POST("/:id/test", providerHandler.TestProvider)                   // 测试DNS提供商连接
			providerGroup.POST("/:id/default", providerHandler.SetDefaultProvider)          // 设置默认DNS提供商
			providerGroup.POST("/:id/sync-domains", providerHandler.SyncProviderDomains)    // 同步单个提供商的域名
			providerGroup.POST("/sync-all-domains", providerHandler.SyncAllProviderDomains) // 同步所有提供商的域名
			providerGroup.DELETE("/batch", providerHandler.BatchDeleteProviders)            // 批量删除DNS提供商
			providerGroup.POST("/batch-test", providerHandler.BatchTestProviders)           // 批量测试DNS提供商
		}

		// DNS批量操作路由
		fmt.Println("[DNS路由] 注册DNS批量操作路由...")
		batchGroup := dnsGroup.Group("/batch")
		batchGroup.Use(middleware.JWTAuth())
		{
			batchGroup.POST("/records", batchHandler.BatchCreateRecords)   // 批量创建DNS记录
			batchGroup.PUT("/records", batchHandler.BatchUpdateRecords)    // 批量更新DNS记录
			batchGroup.DELETE("/records", batchHandler.BatchDeleteRecords) // 批量删除DNS记录
			batchGroup.POST("/sync", batchHandler.SyncDomainRecords)       // 同步域名DNS记录
		}

		// SSL证书管理路由
		fmt.Println("[DNS路由] 注册SSL证书管理路由...")
		certGroup := dnsGroup.Group("/certificates")
		certGroup.Use(middleware.JWTAuth())
		{
			certGroup.POST("", certificateHandler.IssueCertificate)                // 申请证书
			certGroup.GET("", certificateHandler.ListCertificates)                 // 获取证书列表
			certGroup.GET("/stats", certificateHandler.GetCertificateStats)        // 获取证书统计信息
			certGroup.GET("/:id", certificateHandler.GetCertificate)               // 获取证书详情
			certGroup.POST("/:id/renew", certificateHandler.RenewCertificate)      // 续期证书
			certGroup.POST("/:id/revoke", certificateHandler.RevokeCertificate)    // 吊销证书
			certGroup.POST("/:id/deploy", certificateHandler.DeployCertificate)    // 部署证书
			certGroup.GET("/:id/download", certificateHandler.DownloadCertificate) // 下载证书

			// CSR相关功能
			certGroup.POST("/validate-csr", certificateHandler.ValidateCSR)               // 验证CSR内容
			certGroup.POST("/upload-csr", certificateHandler.UploadCSRFile)               // 上传CSR文件
			certGroup.POST("/issue-with-csr", certificateHandler.IssueCertificateWithCSR) // 使用CSR申请证书
		}

		// HTTPS监控路由
		fmt.Println("[DNS路由] 注册HTTPS监控路由...")
		monitorGroup := dnsGroup.Group("/monitors")
		monitorGroup.Use(middleware.JWTAuth())
		{
			monitorGroup.POST("", httpsMonitorHandler.CreateMonitor)                   // 创建监控
			monitorGroup.GET("", httpsMonitorHandler.ListMonitors)                     // 获取监控列表
			monitorGroup.GET("/:id", httpsMonitorHandler.GetMonitor)                   // 获取监控详情
			monitorGroup.PUT("/:id", httpsMonitorHandler.UpdateMonitor)                // 更新监控
			monitorGroup.DELETE("/:id", httpsMonitorHandler.DeleteMonitor)             // 删除监控
			monitorGroup.POST("/:id/check", httpsMonitorHandler.CheckMonitor)          // 手动检查监控
			monitorGroup.GET("/statistics", httpsMonitorHandler.GetMonitorStatistics)  // 获取监控统计
			monitorGroup.GET("/expiring", httpsMonitorHandler.GetExpiringCertificates) // 获取即将过期的证书
		}
	}
	fmt.Println("[DNS路由] DNS路由注册完成！")
}
