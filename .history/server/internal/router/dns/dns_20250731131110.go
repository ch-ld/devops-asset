package dns

import (
	"api-server/internal/api/handler/dns"
	"api-server/internal/middleware/middleware"
	repo "api-server/internal/repository/dns"
	svc "api-server/internal/service/dns"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 定义导出变量以便在main中访问服务
var (
	DomainSvc      *svc.DomainService
	RecordSvc      *svc.RecordService
	ProviderSvc    *svc.ProviderService
	DomainGroupSvc *svc.DomainGroupService
)

// InitDnsRoutes initializes all routes for the DNS module.
func InitDnsRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// TODO: 暂时注释掉DNS路由注册，因为相关依赖还未完全实现
	// 等DNS模块的repository、service、handler都实现后再启用

	// 添加一个简单的测试路由来验证DNS模块是否被加载
	dnsGroup := r.Group("/dns")
	{
		dnsGroup.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "DNS module is loaded",
				"status":  "ok",
			})
		})
	}

	return

	// 以下代码暂时注释，等依赖完善后启用
	/*
	// 1. Dependency Injection
	// Repositories
	domainRepo := repo.NewDomainRepository(db)
	domainGroupRepo := repo.NewDomainGroupRepository(db)
	certificateRepo := repo.NewCertificateRepository(db)
	recordRepo := repo.NewRecordRepository(db)
	providerRepo := repo.NewProviderRepository(db)
	changeLogRepo := repo.NewChangeLogRepository(db)

	// Services
	DomainSvc = svc.NewDomainService(domainRepo, domainGroupRepo, changeLogRepo, db)
	RecordSvc = svc.NewRecordService(recordRepo, domainRepo, providerRepo, changeLogRepo, db)
	ProviderSvc = svc.NewProviderService(providerRepo, changeLogRepo, db)
	DomainGroupSvc = svc.NewDomainGroupService(domainGroupRepo, changeLogRepo, db)

	// Handlers
	domainHandler := dns.NewDomainHandler(DomainSvc)
	recordHandler := dns.NewRecordHandler(RecordSvc)
	providerHandler := dns.NewProviderHandler(ProviderSvc)

	// 2. Route Registration
	dnsGroup := r.Group("/dns")
	{
		// 域名管理路由
		domainGroup := dnsGroup.Group("/domains")
		domainGroup.Use(middleware.JWTAuth())
		{
			domainGroup.POST("", domainHandler.CreateDomain)           // 创建域名
			domainGroup.GET("", domainHandler.ListDomains)             // 获取域名列表
			domainGroup.GET("/:id", domainHandler.GetDomain)           // 获取域名详情
			domainGroup.PUT("/:id", domainHandler.UpdateDomain)        // 更新域名
			domainGroup.DELETE("/:id", domainHandler.DeleteDomain)     // 删除域名
		}

		// DNS记录管理路由
		recordGroup := dnsGroup.Group("/records")
		recordGroup.Use(middleware.JWTAuth())
		{
			recordGroup.POST("", recordHandler.CreateRecord)           // 创建DNS记录
			recordGroup.GET("", recordHandler.ListRecords)             // 获取DNS记录列表
			recordGroup.GET("/:id", recordHandler.GetRecord)           // 获取DNS记录详情
			recordGroup.PUT("/:id", recordHandler.UpdateRecord)        // 更新DNS记录
			recordGroup.DELETE("/:id", recordHandler.DeleteRecord)     // 删除DNS记录
		}

		// DNS提供商管理路由
		providerGroup := dnsGroup.Group("/providers")
		providerGroup.Use(middleware.JWTAuth())
		{
			providerGroup.POST("", providerHandler.CreateProvider)           // 创建DNS提供商
			providerGroup.GET("", providerHandler.ListProviders)             // 获取DNS提供商列表
			providerGroup.GET("/:id", providerHandler.GetProvider)           // 获取DNS提供商详情
			providerGroup.PUT("/:id", providerHandler.UpdateProvider)        // 更新DNS提供商
			providerGroup.DELETE("/:id", providerHandler.DeleteProvider)     // 删除DNS提供商
			providerGroup.POST("/:id/test", providerHandler.TestProvider)    // 测试DNS提供商连接
			providerGroup.POST("/:id/default", providerHandler.SetDefaultProvider) // 设置默认DNS提供商
		}
	}
	*/
}
		// groupGroup := dnsGroup.Group("/groups")
		// groupGroup.Use(middleware.JWTAuth())
		// {
		//     groupGroup.POST("", domainGroupHandler.CreateDomainGroup)
		//     groupGroup.GET("", domainGroupHandler.ListDomainGroups)
		//     groupGroup.GET("/tree", domainGroupHandler.GetDomainGroupTree)
		//     groupGroup.GET("/:id", domainGroupHandler.GetDomainGroup)
		//     groupGroup.PUT("/:id", domainGroupHandler.UpdateDomainGroup)
		//     groupGroup.DELETE("/:id", domainGroupHandler.DeleteDomainGroup)
		//     groupGroup.PUT("/:id/sort", domainGroupHandler.UpdateDomainGroupSort)
		// }

		// TODO: 证书管理路由
		// certGroup := dnsGroup.Group("/certificates")
		// certGroup.Use(middleware.JWTAuth())
		// {
		//     certGroup.POST("", certificateHandler.CreateCertificate)
		//     certGroup.GET("", certificateHandler.ListCertificates)
		//     certGroup.GET("/:id", certificateHandler.GetCertificate)
		//     certGroup.PUT("/:id", certificateHandler.UpdateCertificate)
		//     certGroup.DELETE("/:id", certificateHandler.DeleteCertificate)
		//     certGroup.POST("/:id/renew", certificateHandler.RenewCertificate)
		//     certGroup.GET("/:id/download", certificateHandler.DownloadCertificate)
		//     certGroup.POST("/:id/deploy", certificateHandler.DeployCertificate)
		// }

		// TODO: 统计和监控路由
		// statsGroup := dnsGroup.Group("/stats")
		// statsGroup.Use(middleware.JWTAuth())
		// {
		//     statsGroup.GET("/domains", domainHandler.GetDomainStatistics)
		//     statsGroup.GET("/records", recordHandler.GetRecordStatistics)
		//     statsGroup.GET("/certificates", certificateHandler.GetCertificateStatistics)
		//     statsGroup.GET("/providers", providerHandler.GetProviderStatistics)
		// }

		// TODO: 批量操作路由
		// batchGroup := dnsGroup.Group("/batch")
		// batchGroup.Use(middleware.JWTAuth())
		// {
		//     batchGroup.POST("/domains", domainHandler.BatchOperateDomains)
		//     batchGroup.POST("/records", recordHandler.BatchOperateRecords)
		//     batchGroup.POST("/sync", recordHandler.BatchSyncRecords)
		// }

		// TODO: 导入导出路由
		// importExportGroup := dnsGroup.Group("/import-export")
		// importExportGroup.Use(middleware.JWTAuth())
		// {
		//     importExportGroup.POST("/domains/import", domainHandler.ImportDomains)
		//     importExportGroup.POST("/domains/export", domainHandler.ExportDomains)
		//     importExportGroup.POST("/records/import", recordHandler.ImportRecords)
		//     importExportGroup.POST("/records/export", recordHandler.ExportRecords)
		// }

		// TODO: 同步和验证路由
		// syncGroup := dnsGroup.Group("/sync")
		// syncGroup.Use(middleware.JWTAuth())
		// {
		//     syncGroup.POST("/records", recordHandler.SyncRecords)
		//     syncGroup.POST("/zones", recordHandler.SyncZones)
		//     syncGroup.GET("/status", recordHandler.GetSyncStatus)
		// }

		// TODO: 测试和验证路由
		// testGroup := dnsGroup.Group("/test")
		// testGroup.Use(middleware.JWTAuth())
		// {
		//     testGroup.POST("/connection", providerHandler.TestConnection)
		//     testGroup.POST("/dns-resolution", recordHandler.TestDNSResolution)
		//     testGroup.POST("/certificate-validation", certificateHandler.ValidateCertificate)
		// }
	}
}
