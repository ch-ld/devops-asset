package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"api-server/internal/config"
	"api-server/internal/constants"
	"api-server/internal/response/response"
	"api-server/pkg/auth/authentication"
)

const (
	JWTDataKey = "jwtData"
)

// JWTAuth JWT认证中间件 - 统一的JWT验证逻辑
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开发模式下跳过JWT认证
		if config.IsDev() {
			// 设置一个默认的用户ID到上下文，模拟已登录用户
			c.Set(JWTDataKey, "1") // 使用管理员用户ID
			c.Next()
			return
		}

		// 防止文件未发送完成就返回错误, 导致前端504而不是正确响应
		c.FormFile("file")

		// 获取token - 支持多种方式
		token := extractToken(c)
		if token == "" {
			response.ReturnError(c, response.UNAUTHENTICATED, "未携带token")
			return
		}

		// 验证token
		data, err := authentication.JWTDecrypt(token)
		if err != nil {
			response.ReturnError(c, response.UNAUTHENTICATED, "token验证失败")
			return
		}

		// 设置用户数据到上下文
		c.Set(JWTDataKey, data)
		c.Next()
	}
}

// extractToken 从多个来源提取token
func extractToken(c *gin.Context) string {
	// 1. 优先从Authorization头部获取 (标准方式)
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// 支持 "Bearer <token>" 格式
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer ")
		}
		// 直接返回token
		return authHeader
	}

	// 2. 从Access-Token头部获取 (兼容旧版本)
	accessToken := c.GetHeader("Access-Token")
	if accessToken != "" {
		return accessToken
	}

	// 3. 从自定义头部获取
	jwtToken := c.GetHeader(constants.JWTHeaderKey)
	if jwtToken != "" {
		return jwtToken
	}

	// 4. 从查询参数获取 (用于特殊场景，如WebSocket)
	queryToken := c.Query("token")
	if queryToken != "" {
		return queryToken
	}

	return ""
}

// TokenVerify 兼容旧版本的函数名 - 已弃用，请使用JWTAuth()
// Deprecated: 请使用JWTAuth()替代
func TokenVerify(c *gin.Context) {
	JWTAuth()(c)
}
