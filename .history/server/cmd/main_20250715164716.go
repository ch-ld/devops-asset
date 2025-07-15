// @title DevOps Asset Management System API
// @version 1.0.0
// @description 企业级DevOps资产管理系统API文档
// @termsOfService https://example.com/terms/

// @contact.name API Support
// @contact.url https://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token authentication

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-server/internal/common/common/cron"
	"api-server/internal/config"
	"api-server/internal/db/mysql"
	"api-server/internal/db/mysql/system"
	"api-server/internal/db/redis"
	"api-server/internal/middleware/middleware"
	cmdb_model "api-server/internal/model/cmdb"
	repo "api-server/internal/repository/cmdb"
	"api-server/internal/router/cmdb"
	"api-server/internal/router/system"
	cmdbService "api-server/internal/service/cmdb"
	"api-server/pkg/logger/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// migrate 数据库迁移函数
func migrate() error {
	// 初始化数据库连接
	if err := mysql.InitClient(); err != nil {
		zap.L().Error("Failed to initialize database", zap.Error(err))
		return err
	}

	err := system.Migrate(mysql.GetClient())
	if err != nil {
		zap.L().Error("migrate failed", zap.Error(err))
		return err
	}
	err = cmdb_model.Migrate(mysql.GetClient())
	if err != nil {
		zap.L().Error("cmdb migrate failed", zap.Error(err))
		return err
	}
	zap.L().Info("migration completed successfully")
	return nil
}

func main() {
	// 解析命令行参数
	var configPath string
	for i, arg := range os.Args[1:] {
		if arg == "--migrate" {
			// 设置开发模式环境变量
			os.Setenv("APP_MODE", "dev")
			config.RunModel = config.RunModelDevValue
			// 加载配置
			if err := config.LoadConfig(""); err != nil {
				fmt.Printf("Failed to load config: %v\n", err)
				os.Exit(1)
			}
			log.SetLogger()
			migrate()
			return
		}
		if arg == "--dev" {
			os.Setenv("APP_MODE", "dev")
		}
		if arg == "--config" && i+1 < len(os.Args[1:]) {
			configPath = os.Args[i+2]
		}
	}

	// 1. 加载配置文件
	if err := config.LoadConfig(configPath); err != nil {
		log.Printf("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	// 2. 设置日志
	log.SetLogger()

	// 3. 设置运行模式
	if config.GlobalConfig.App.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 4. 初始化数据库连接
	if err := mysql.InitClient(); err != nil {
		zap.L().Error("Failed to initialize database", zap.Error(err))
		os.Exit(1)
	}
	defer mysql.Close()

	// 5. 初始化Redis连接
	if err := redis.Init(); err != nil {
		zap.L().Warn("Redis initialization failed, continuing without cache", zap.Error(err))
	} else {
		defer redis.Close()
	}

	// 6. 初始化定时任务
	cron.InitCronJobs()

	// 7. 创建Gin引擎
	r := gin.Default()

	// 8. 注册中间件
	r.Use(gin.Recovery())
	r.Use(middleware.CorssDomainHandler())

	// 9. 添加健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
			"version": config.GlobalConfig.App.Version,
		})
	})

	// 10. 注册API路由
	apiGroup := r.Group("/api/v1")
	system.InitSystemRoutes(apiGroup)                // 注册系统路由
	cmdb.InitCmdbRoutes(apiGroup, mysql.GetClient()) // 注册CMDB路由

	// 11. 启动后台服务
	// 获取数据库实例
	db := mysql.GetClient()

	// 11.1 启动CMDB定时任务服务
	hostRepo := repo.NewHostRepository(db)
	providerRepo := repo.NewProviderRepository(db)
	alertService := cmdbService.NewAlertService(hostRepo, providerRepo)
	hostService := cmdbService.NewHostService(hostRepo, providerRepo, nil)
	cronService := cmdbService.NewCronService(hostRepo, providerRepo, alertService, hostService)

	go func() {
		zap.L().Info("Starting CMDB cron service...")
		if err := cronService.Start(); err != nil {
			zap.L().Error("Failed to start CMDB cron service", zap.Error(err))
		}
	}()

	// 11.2 启动指标监控服务
	metricsService := cmdbService.NewMetricsService(hostRepo)
	go func() {
		zap.L().Info("Starting CMDB metrics service...")
		if err := metricsService.Start(); err != nil {
			zap.L().Error("Failed to start CMDB metrics service", zap.Error(err))
		}
	}()

	// 12. 启动HTTP服务器
	port := fmt.Sprintf(":%d", config.GlobalConfig.App.Port)
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	// 13. 优雅关闭设置
	go func() {
		zap.L().Info("Starting server", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 14. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutting down server...")

	// 15. 关闭后台服务
	cronService.Stop()
	metricsService.Stop()

	// 16. 给服务器一些时间完成未完成的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server forced to shutdown", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
