package router

import (
	"net/http"

	"api-server/internal/config"
	"api-server/internal/db/mysql"
	"api-server/internal/middleware/middleware"
	"api-server/internal/router/cmdb"
	"api-server/internal/router/dns"
	"api-server/internal/router/system"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	// 设置运行模式
	if !config.IsDev() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(middleware.ErrorRecovery())
	r.Use(middleware.CorssDomainHandler())

	// 静态文件服务 - 头像上传文件访问
	r.Static("/uploads", "./uploads")

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
			"version": "1.0.0",
		})
	})

	// API信息
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "DevOps Asset Management System",
			"version":     "1.0.0",
			"description": "企业级DevOps资产管理系统",
			"docs":        "/swagger/index.html",
			"health":      "/health",
		})
	})

	// Swagger登录页面
	r.GET("/swagger-login", middleware.SwaggerLogin)
	r.POST("/swagger-login", middleware.SwaggerLogin)

	// Swagger文档（需要认证）
	r.GET("/swagger/*any", middleware.SwaggerAuth(), ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API路由组
	apiV1 := r.Group("/api/v1")
	{
		// 系统模块路由
		system.InitSystemRoutes(apiV1)
		// 主机管理模块路由
		cmdb.InitCmdbRoutes(apiV1, mysql.GetClient())
		// DNS管理模块路由
		dns.InitDnsRoutes(apiV1, mysql.GetClient())
	}

	return r
}
