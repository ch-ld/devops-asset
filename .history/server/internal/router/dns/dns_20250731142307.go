package dns

import (
	"api-server/internal/api/handler/dns"
	"api-server/internal/middleware/middleware"
	repo "api-server/internal/repository/dns"
	svc "api-server/internal/service/dns"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	// certificateRepo := repo.NewCertificateRepository(db) // TODO: 暂时注释，等证书功能实现后启用
	recordRepo := repo.NewRecordRepository(db)
	providerRepo := repo.NewProviderRepository(db)
	changeLogRepo := repo.NewChangeLogRepository(db)
	fmt.Println("[DNS路由] Repository创建完成")

	// Services
	fmt.Println("[DNS路由] 正在创建Service...")
	DomainSvc = svc.NewDomainService(domainRepo, domainGroupRepo, changeLogRepo, db)
	RecordSvc = svc.NewRecordService(recordRepo, domainRepo, providerRepo, changeLogRepo, db)
	ProviderSvc = svc.NewProviderService(providerRepo, changeLogRepo, db)
	DomainGroupSvc = svc.NewDomainGroupService(domainGroupRepo, changeLogRepo, db)
	fmt.Println("[DNS路由] Service创建完成")

	// Handlers
	fmt.Println("[DNS路由] 正在创建Handler...")
	domainHandler := dns.NewDomainHandler(DomainSvc)
	recordHandler := dns.NewRecordHandler(RecordSvc)
	providerHandler := dns.NewProviderHandler(ProviderSvc)
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
			providerGroup.POST("", providerHandler.CreateProvider)                 // 创建DNS提供商
			providerGroup.GET("", providerHandler.ListProviders)                   // 获取DNS提供商列表
			providerGroup.GET("/:id", providerHandler.GetProvider)                 // 获取DNS提供商详情
			providerGroup.PUT("/:id", providerHandler.UpdateProvider)              // 更新DNS提供商
			providerGroup.DELETE("/:id", providerHandler.DeleteProvider)           // 删除DNS提供商
			providerGroup.POST("/:id/test", providerHandler.TestProvider)          // 测试DNS提供商连接
			providerGroup.POST("/:id/default", providerHandler.SetDefaultProvider) // 设置默认DNS提供商
		}
	}
	fmt.Println("[DNS路由] DNS路由注册完成！")
}
