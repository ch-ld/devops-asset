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
	stdlog "log" // 标准库log使用别名避免冲突
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-server/internal/common/common/cron"
	"api-server/internal/config"
	"api-server/internal/db"
	"api-server/internal/db/mysql"
	"api-server/internal/db/redis"
	"api-server/internal/middleware/middleware"
	"api-server/internal/router/cmdb"
	systemRouter "api-server/internal/router/system" // 使用别名避免冲突
	"api-server/internal/service/cmdb"
	applog "api-server/pkg/logger/log" // 使用别名避免冲突

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// migrate 数据库迁移函数
func migrate() error {
	zap.L().Info("Starting database migration")

	// 初始化数据库连接
	if err := mysql.InitClient(); err != nil {
		zap.L().Error("Failed to initialize database", zap.Error(err))
		return err
	}
	zap.L().Info("Database connection initialized successfully")

	// 执行统一的数据库迁移
	if err := db.MigrateAll(mysql.GetClient()); err != nil {
		zap.L().Error("Database migration failed", zap.Error(err))
		return err
	}

	zap.L().Info("Database migration completed successfully")
	return nil
}

func main() {
	stdlog.Println("=== 服务启动开始 ===")
	stdlog.Println("[启动日志] 解析命令行参数...")

	// 解析命令行参数
	var configPath string
	for i, arg := range os.Args[1:] {
		if arg == "--migrate" {
			stdlog.Println("[启动日志] 检测到--migrate参数，将执行数据库迁移")
			// 设置开发模式环境变量
			os.Setenv("APP_MODE", "dev")
			config.RunModel = config.RunModelDevValue
			// 加载配置
			stdlog.Println("[启动日志] 加载配置...")
			if err := config.LoadConfig(""); err != nil {
				stdlog.Printf("[启动日志] 配置加载失败: %v", err)
				fmt.Printf("Failed to load config: %v\n", err)
				os.Exit(1)
			}
			stdlog.Println("[启动日志] 配置加载成功")

			stdlog.Println("[启动日志] 设置日志...")
			applog.SetLogger() // 使用别名调用
			stdlog.Println("[启动日志] 日志设置完成")

			migrate()
			stdlog.Println("[启动日志] 数据库迁移完成，程序退出")
			return
		}
		if arg == "--dev" {
			stdlog.Println("[启动日志] 检测到--dev参数，将使用开发模式")
			os.Setenv("APP_MODE", "dev")
		}
		if arg == "--config" && i+1 < len(os.Args[1:]) {
			configPath = os.Args[i+2]
			stdlog.Printf("[启动日志] 指定配置文件: %s", configPath)
		}
	}

	// 1. 加载配置文件
	stdlog.Println("[启动日志] 步骤1: 加载配置文件")
	if err := config.LoadConfig(configPath); err != nil {
		stdlog.Printf("[启动日志] 配置加载失败: %v", err)
		os.Exit(1)
	}
	stdlog.Println("[启动日志] 配置加载成功")

	// 2. 设置日志
	stdlog.Println("[启动日志] 步骤2: 设置日志系统")
	applog.SetLogger() // 使用别名调用
	stdlog.Println("[启动日志] 日志系统设置成功")

	// 3. 设置运行模式
	stdlog.Println("[启动日志] 步骤3: 设置运行模式")
	if config.GlobalConfig.App.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
		stdlog.Println("[启动日志] 设置为生产模式")
	} else {
		gin.SetMode(gin.DebugMode)
		stdlog.Println("[启动日志] 设置为开发模式")
	}

	// 4. 初始化数据库连接
	stdlog.Println("[启动日志] 步骤4: 初始化数据库连接")
	if err := mysql.InitClient(); err != nil {
		stdlog.Printf("[启动日志] 数据库连接失败: %v", err)
		zap.L().Error("Failed to initialize database", zap.Error(err))
		os.Exit(1)
	}
	stdlog.Println("[启动日志] 数据库连接成功")
	defer mysql.Close()

	// 5. 初始化Redis连接
	stdlog.Println("[启动日志] 步骤5: 初始化Redis连接")
	if err := redis.Init(); err != nil {
		stdlog.Printf("[启动日志] Redis连接失败，但将继续运行: %v", err)
		zap.L().Warn("Redis initialization failed, continuing without cache", zap.Error(err))
	} else {
		stdlog.Println("[启动日志] Redis连接成功")
		defer redis.Close()
	}

	// 6. 初始化定时任务
	stdlog.Println("[启动日志] 步骤6: 初始化定时任务")
	cron.InitCronJobs()
	stdlog.Println("[启动日志] 定时任务初始化成功")

	// 7. 创建Gin引擎
	stdlog.Println("[启动日志] 步骤7: 创建Gin引擎")
	r := gin.Default()
	stdlog.Println("[启动日志] Gin引擎创建成功")

	// 8. 注册中间件
	stdlog.Println("[启动日志] 步骤8: 注册中间件")
	r.Use(gin.Recovery())
	r.Use(middleware.CorssDomainHandler())
	stdlog.Println("[启动日志] 中间件注册成功")

	// 9. 添加健康检查路由
	stdlog.Println("[启动日志] 步骤9: 添加健康检查路由")
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
			"version": config.GlobalConfig.App.Version,
		})
	})

	// 根路由信息
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "DevOps Asset Management System",
			"version":     config.GlobalConfig.App.Version,
			"description": "企业级DevOps资产管理系统",
			"docs":        "/swagger/index.html",
			"health":      "/health",
		})
	})
	stdlog.Println("[启动日志] 基础路由添加成功")

	// 10. 注册API路由
	stdlog.Println("[启动日志] 步骤10: 注册API路由")
	apiGroup := r.Group("/api/v1")
	stdlog.Println("[启动日志] 步骤10.1: 注册系统路由")
	systemRouter.InitSystemRoutes(apiGroup) // 使用别名调用
	stdlog.Println("[启动日志] 步骤10.2: 注册CMDB路由")
	cmdb.InitCmdbRoutes(apiGroup, mysql.GetClient()) // 注册CMDB路由
	stdlog.Println("[启动日志] API路由注册成功")

	// 11. 启动后台服务
	stdlog.Println("[启动日志] 步骤11: 启动后台服务")

	// 11.1 启动CMDB定时任务服务
	stdlog.Println("[启动日志] 步骤11.1: 使用已初始化的CMDB服务组件")

	// 这些服务已经在路由初始化时创建
	cronService := cmdb.CronSvc
	metricsService := cmdb.MetricsSvc

	stdlog.Println("[启动日志] 步骤11.2: 启动CMDB定时任务服务")
	go func() {
		stdlog.Println("[启动日志] 开始启动CMDB定时任务服务...")
		zap.L().Info("Starting CMDB cron service...")
		if err := cronService.Start(); err != nil {
			stdlog.Printf("[启动日志] CMDB定时任务服务启动失败: %v", err)
			zap.L().Error("Failed to start CMDB cron service", zap.Error(err))
		} else {
			stdlog.Println("[启动日志] CMDB定时任务服务启动成功")
		}
	}()

	// 11.3 启动指标监控服务
	stdlog.Println("[启动日志] 步骤11.3: 启动指标监控服务")
	go func() {
		stdlog.Println("[启动日志] 开始启动CMDB指标监控服务...")
		zap.L().Info("Starting CMDB metrics service...")
		if err := metricsService.Start(); err != nil {
			stdlog.Printf("[启动日志] CMDB指标监控服务启动失败: %v", err)
			zap.L().Error("Failed to start CMDB metrics service", zap.Error(err))
		} else {
			stdlog.Println("[启动日志] CMDB指标监控服务启动成功")
		}
	}()

	// 12. 启动HTTP服务器
	stdlog.Println("[启动日志] 步骤12: 准备启动HTTP服务器")
	port := fmt.Sprintf(":%d", config.GlobalConfig.App.Port)
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}
	stdlog.Printf("[启动日志] HTTP服务器配置在端口: %s", port)

	// 13. 优雅关闭设置
	stdlog.Println("[启动日志] 步骤13: 启动HTTP服务器")
	stdlog.Printf("[启动日志] 开始监听端口 %s", port)
	zap.L().Info("Starting server", zap.String("port", port))
	stdlog.Println("[启动日志] HTTP服务器启动成功")
	stdlog.Println("=== 服务启动完成 ===")
	stdlog.Printf("=== 可通过 http://localhost%s/health 访问健康检查 ===", port)

	// 创建一个通道来捕获服务器错误
	errChan := make(chan error)

	// 启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			stdlog.Printf("[启动日志] HTTP服务启动失败: %v", err)
			zap.L().Error("Failed to start server", zap.Error(err))
			errChan <- err
		}
	}()

	// 检查是否有立即的错误
	select {
	case err := <-errChan:
		stdlog.Printf("[启动日志] 服务器启动出错: %v", err)
		os.Exit(1)
	case <-time.After(1 * time.Second):
		stdlog.Println("[启动日志] 服务器成功启动并监听请求")
	}

	// 14. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	stdlog.Println("[启动日志] 接收到停止信号，开始关闭服务...")
	zap.L().Info("Shutting down server...")

	// 15. 关闭后台服务
	zap.L().Info("Shutting down background services")

	// 停止定时任务服务
	if cronService != nil {
		zap.L().Info("Stopping cron service")
		cronService.Stop()
		zap.L().Info("Cron service stopped")
	}

	// 停止指标监控服务
	if metricsService != nil {
		zap.L().Info("Stopping metrics service")
		metricsService.Stop()
		zap.L().Info("Metrics service stopped")
	}

	// 16. 优雅关闭HTTP服务器
	zap.L().Info("Gracefully shutting down HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 创建一个通道来接收关闭结果
	shutdownComplete := make(chan error, 1)
	go func() {
		shutdownComplete <- srv.Shutdown(ctx)
	}()

	// 等待关闭完成或超时
	select {
	case err := <-shutdownComplete:
		if err != nil {
			zap.L().Error("HTTP server shutdown error", zap.Error(err))
			os.Exit(1)
		} else {
			zap.L().Info("HTTP server shutdown completed")
		}
	case <-ctx.Done():
		zap.L().Error("HTTP server shutdown timeout, forcing exit")
		os.Exit(1)
	}

	// 17. 关闭数据库连接
	zap.L().Info("Closing database connections")
	mysql.Close()
	redis.Close()

	zap.L().Info("Server graceful shutdown completed")
	os.Exit(0)
}
