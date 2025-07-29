package cmdb

import (
	"api-server/internal/api/handler/cmdb"
	"api-server/internal/middleware/middleware"
	repo "api-server/internal/repository/cmdb"
	svc "api-server/internal/service/cmdb"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 定义导出变量以便在main中访问服务
var (
	CronSvc    *svc.CronService
	MetricsSvc *svc.MetricsService
	HostSvc    *svc.HostService
	AlertSvc   *svc.AlertService
)

// InitCmdbRoutes initializes all routes for the CMDB module.
func InitCmdbRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// 1. Dependency Injection
	// Repositories
	providerRepo := repo.NewProviderRepository(db)
	hostRepo := repo.NewHostRepository(db)
	hostGroupRepo := repo.NewHostGroupRepository(db)
	customFieldRepo := repo.NewCustomFieldRepository(db)

	// Services
	providerSvc := svc.NewProviderService(providerRepo, hostRepo)
	HostSvc = svc.NewHostService(hostRepo, providerRepo, hostGroupRepo)
	hostGroupSvc := svc.NewHostGroupService(hostGroupRepo, hostRepo)
	customFieldSvc := svc.NewCustomFieldService(customFieldRepo)
	MetricsSvc = svc.NewMetricsService(hostRepo)                            // 指标监控服务
	AlertSvc = svc.NewAlertService(hostRepo, providerRepo)                  // 告警服务
	CronSvc = svc.NewCronService(hostRepo, providerRepo, AlertSvc, HostSvc) // 定时任务服务

	// Handlers
	providerHandler := cmdb.NewProviderHandler(providerSvc)
	hostHandler := cmdb.NewHostHandler(HostSvc, hostGroupSvc)
	hostGroupHandler := cmdb.NewHostGroupHandler(hostGroupSvc)
	webSshHandler := cmdb.NewWebSshHandler(hostRepo)
	sftpHandler := cmdb.NewSftpHandler(hostRepo)
	customFieldHandler := cmdb.NewCustomFieldHandler(customFieldSvc)
	hostBatchHandler := cmdb.NewHostBatchHandler(HostSvc)                        // 批量操作处理器
	dashboardHandler := cmdb.NewDashboardHandler(hostRepo, MetricsSvc, AlertSvc) // 仪表板处理器
	hostMetricsHandler := cmdb.NewHostMetricsHandler(MetricsSvc, HostSvc)        // 主机监控指标处理器
	alertHandler := cmdb.NewAlertHandler(AlertSvc)                               // 告警处理器

	// 创建新的Repository实例 (暂时注释)
	// loadBalancerRepo := repo.NewLoadBalancerRepository(db)
	// databaseRepo := repo.NewDatabaseRepository(db)
	// redisRepo := repo.NewRedisRepository(db)

	// 创建资源管理服务和处理器 (暂时注释)
	// resourceService := svc.NewResourceService(hostRepo, loadBalancerRepo, databaseRepo, redisRepo, providerRepo)
	// resourceHandler := cmdb.NewResourceHandler(resourceService)

	// 2. Route Grouping & Middleware
	// All CMDB routes will be authenticated
	cmdbGroup := r.Group("/cmdb")
	cmdbGroup.Use(middleware.JWTAuth())
	{
		// Provider routes
		providerGroup := cmdbGroup.Group("/providers")
		{
			providerGroup.POST("", providerHandler.CreateProvider)
			providerGroup.GET("", providerHandler.ListProviders)
			providerGroup.GET("/types", providerHandler.GetProviderTypes)
			providerGroup.GET("/types/:type/regions", providerHandler.GetProviderRegions)
			providerGroup.GET("/:id", providerHandler.GetProvider)
			providerGroup.PUT("/:id", providerHandler.UpdateProvider)
			providerGroup.DELETE("/:id", providerHandler.DeleteProvider)
			providerGroup.POST("/:id/sync", providerHandler.SyncResources)
			providerGroup.POST("/validate", providerHandler.ValidateCredentials)
		}

		// Host routes
		hostGroup := cmdbGroup.Group("/hosts")
		{
			hostGroup.POST("", hostHandler.CreateHost)
			hostGroup.GET("", hostHandler.ListHosts)
			hostGroup.POST("/batch_import", hostHandler.BatchImportHosts)
			hostGroup.GET("/batch_export", hostHandler.BatchExportHosts)
			hostGroup.GET("/filter_options", hostHandler.GetHostFilterOptions)
			hostGroup.GET("/:id", hostHandler.GetHost)
			hostGroup.PUT("/:id", hostHandler.UpdateHost)
			hostGroup.DELETE("/:id", hostHandler.DeleteHost)

			// 主机同步相关路由
			hostGroup.POST("/sync", hostHandler.SyncHosts)
			hostGroup.POST("/:id/sync_status", hostHandler.SyncHostStatus)
			// 自建主机路由
			hostGroup.POST("/manual", hostHandler.CreateManualHost)
			// 主机移动路由
			hostGroup.POST("/:id/move", hostHandler.MoveHost)
		}

		// 批量操作路由
		hostBatchHandler.RegisterRoutes(cmdbGroup)

		// 主机监控指标路由
		hostMetricsHandler.RegisterRoutes(cmdbGroup)

		// Host Group routes
		hostGroupGroup := cmdbGroup.Group("/host_groups")
		{
			hostGroupGroup.POST("", hostGroupHandler.CreateHostGroup)
			hostGroupGroup.GET("/tree", hostGroupHandler.GetHostGroupTree)
			hostGroupGroup.GET("/:id", hostGroupHandler.GetHostGroup)
			hostGroupGroup.PUT("/:id", hostGroupHandler.UpdateHostGroup)
			hostGroupGroup.DELETE("/:id", hostGroupHandler.DeleteHostGroup)
			hostGroupGroup.PUT("/:id/move", hostGroupHandler.MoveHostGroup)
		}

		// 主机组下的主机列表
		cmdbGroup.GET("/groups/:group_id/hosts", hostHandler.ListGroupHosts)

		// 主机自定义字段定义管理接口
		customFieldGroup := cmdbGroup.Group("/host/custom_fields")
		{
			customFieldGroup.POST("", customFieldHandler.CreateCustomField)
			customFieldGroup.PUT("", customFieldHandler.UpdateCustomField)
			customFieldGroup.DELETE("", customFieldHandler.DeleteCustomField)
			customFieldGroup.GET("", customFieldHandler.ListCustomFields)
		}

		// Alert routes - 告警管理路由
		alertGroup := cmdbGroup.Group("/alerts")
		{
			alertGroup.GET("", alertHandler.GetAlerts)
			alertGroup.GET("/statistics", alertHandler.GetAlertStatistics)
			alertGroup.GET("/summary", alertHandler.GetAlertSummary)
			alertGroup.GET("/expiring", alertHandler.CheckExpiringHosts)
			alertGroup.GET("/status", alertHandler.CheckHostStatus)
			alertGroup.POST("/check", alertHandler.TriggerAlertCheck)
		}

		// 主机告警路由
		cmdbGroup.GET("/hosts/:id/alerts", alertHandler.GetHostAlerts)

		// Resource routes - 资源管理路由 (暂时注释)
		// resourceGroup := cmdbGroup.Group("/resources")
		// {
		// 	resourceGroup.GET("/statistics", resourceHandler.GetOverallStatistics)
		// 	resourceGroup.GET("/summary", resourceHandler.GetResourceSummary)
		// 	resourceGroup.GET("/providers", resourceHandler.GetProviderSummary)
		// 	resourceGroup.GET("/regions", resourceHandler.GetRegionSummary)
		// 	resourceGroup.GET("/status", resourceHandler.GetStatusSummary)
		// 	resourceGroup.GET("/expiring", resourceHandler.GetExpiringResources)
		// 	resourceGroup.GET("/dashboard", resourceHandler.GetDashboardData)
		// 	resourceGroup.GET("/health", resourceHandler.GetResourceHealth)
		// }

		// Dashboard路由
		dashboardHandler.RegisterRoutes(cmdbGroup)

		// SFTP路由
		sftpHandler.RegisterRoutes(cmdbGroup)

		// WebTerminal相关路由
		webTerminalHandler.RegisterRoutes(r)
	}

	// WebSocket routes - WebSSH路由，需要在API组下注册
	wsGroup := cmdbGroup.Group("/ws")
	// WebSocket连接也需要JWT认证，但认证方式可能需要特殊处理
	{
		wsGroup.GET("/ssh", webSshHandler.HandleSSH)
	}

	// 注释掉在路由初始化时启动定时任务服务，避免重复启动
	// if err := cronSvc.Start(); err != nil {
	// 	panic("Failed to start cron service: " + err.Error())
	// }
}

// RegisterHostRoutes 注册主机管理相关路由
func RegisterHostRoutes(r *gin.RouterGroup, hostHandler *cmdb.HostHandler) {
	hosts := r.Group("/hosts")
	{
		// 主机同步接口
		hosts.POST("/sync", hostHandler.SyncHosts)
		hosts.POST("/:id/sync_status", hostHandler.SyncHostStatus)

		// 自建主机接口
		hosts.POST("/manual", hostHandler.CreateManualHost)

		// 主机移动接口
		hosts.POST("/:id/move", hostHandler.MoveHost)
	}

	groups := r.Group("/groups")
	{
		// 主机组下的主机列表
		groups.GET("/:group_id/hosts", hostHandler.ListGroupHosts)
	}
}
