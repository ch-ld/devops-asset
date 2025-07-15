package main

import (
	"api-server/internal/config"
	"api-server/internal/db/mysql"
	"api-server/internal/db/redis"
	"api-server/internal/middleware/middleware"
	repo "api-server/internal/repository/cmdb"
	"api-server/internal/router/cmdb"
	"api-server/internal/router/system"
	cmdbService "api-server/internal/service/cmdb"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	if err := config.LoadConfig(""); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. 设置运行模式
	if config.GlobalConfig.App.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 3. 初始化数据库连接
	if err := mysql.InitClient(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer mysql.Close()

	// 4. 初始化Redis连接
	if err := redis.Init(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redis.Close()

	// 5. 创建Gin引擎
	r := gin.Default()

	// 6. 注册中间件
	r.Use(gin.Recovery())
	r.Use(middleware.CorssDomainHandler())
	// 移除未定义的中间件

	// 添加健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
			"version": config.GlobalConfig.App.Version,
		})
	})

	// 7. 注册路由
	apiGroup := r.Group("/api/v1")
	system.InitSystemRoutes(apiGroup)                // 注册系统路由
	cmdb.InitCmdbRoutes(apiGroup, mysql.GetClient()) // 注册CMDB路由

	// 8. 启动后台服务
	// 获取数据库实例
	db := mysql.GetClient()

	// 8.1 启动CMDB定时任务服务
	hostRepo := repo.NewHostRepository(db)
	providerRepo := repo.NewProviderRepository(db)
	alertService := cmdbService.NewAlertService(hostRepo, providerRepo)
	hostService := cmdbService.NewHostService(hostRepo, providerRepo, nil)
	cronService := cmdbService.NewCronService(hostRepo, providerRepo, alertService, hostService)

	go func() {
		log.Println("Starting CMDB cron service...")
		if err := cronService.Start(); err != nil {
			log.Printf("Failed to start CMDB cron service: %v", err)
		}
	}()

	// 8.2 启动指标监控服务
	metricsService := cmdbService.NewMetricsService(hostRepo)
	go func() {
		log.Println("Starting CMDB metrics service...")
		if err := metricsService.Start(); err != nil {
			log.Printf("Failed to start CMDB metrics service: %v", err)
		}
	}()

	// 9. 启动HTTP服务器
	port := fmt.Sprintf(":%d", config.GlobalConfig.App.Port)
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	// 10. 优雅关闭
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 11. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 12. 关闭后台服务
	cronService.Stop()
	metricsService.Stop()

	// 13. 给服务器一些时间完成未完成的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
