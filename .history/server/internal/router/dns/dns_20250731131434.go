package dns

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitDnsRoutes initializes all routes for the DNS module.
func InitDnsRoutes(r *gin.RouterGroup, db *gorm.DB) {
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
}
