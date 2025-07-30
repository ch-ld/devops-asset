package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	// "api-server/internal/db/redis" // 暂时注释掉
	"api-server/internal/middleware/middleware"
	"api-server/internal/response/response"
	"api-server/internal/service"
)

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username  string `json:"username" form:"username" binding:"required" validate:"required,min=3,max=50" example:"admin"`
	Password  string `json:"password" form:"password" binding:"required" validate:"required,min=6,max=50" example:"admin123"`
	Captcha   string `json:"captcha" form:"captcha" binding:"required" validate:"required,min=4,max=6" example:"1234"`
	CaptchaID string `json:"captcha_id" form:"captcha_id" binding:"required" validate:"required" example:"captcha_123"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，需要提供用户名、密码和验证码
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param body body LoginRequest true "登录参数"
// @Success 200 {object} response.ResponseData{data=LoginResponse} "登录成功"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 401 {object} response.ResponseError "认证失败"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if !middleware.BindAndValidate(c, &req) {
		return
	}

	// 暂时跳过验证码验证以便测试
	// if !redis.GetCaptchaStore().Verify(req.CaptchaID, req.Captcha, true) {
	// 	zap.L().Warn("验证码验证失败", zap.String("username", req.Username), zap.String("ip", c.ClientIP()))
	// 	response.ReturnError(c, response.INVALID_ARGUMENT, "验证码错误")
	// 	return
	// }

	// 获取客户端IP
	clientIP := c.ClientIP()

	// 调用服务层进行登录
	token, err := service.GetUserService().Login(req.Username, req.Password)

	// 记录登录日志
	var status uint = 1 // 成功
	if err != nil {
		status = 0 // 失败
	}

	// 异步记录登录日志，不影响主流程
	go func() {
		if logErr := service.GetUserService().CreateLoginLog(req.Username, clientIP, status); logErr != nil {
			zap.L().Error("记录登录日志失败", zap.Error(logErr))
		}
	}()

	// 处理登录错误
	if err != nil {
		middleware.HandleServiceError(c, err)
		return
	}

	// 登录成功
	zap.L().Info("用户登录成功", zap.String("username", req.Username), zap.String("ip", clientIP))
	response.ReturnData(c, LoginResponse{
		AccessToken: token,
	})
}

// RegisterAPI 用户注册
// @Summary 用户注册
// @Description 用户注册接口，需要提供用户名、密码、姓名、手机号和性别
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param body body service.RegisterRequest true "注册参数"
// @Success 200 {object} response.ResponseData "注册成功"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 403 {object} response.ResponseError "注册功能未开启"
// @Failure 409 {object} response.ResponseError "用户已存在"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /auth/register [post]
func RegisterAPI(c *gin.Context) {
	var req service.RegisterRequest
	if !middleware.BindAndValidate(c, &req) {
		return
	}

	// 调用服务层
	err := service.GetUserService().Register(&req)
	if err != nil {
		middleware.HandleServiceError(c, err)
		return
	}

	zap.L().Info("用户注册成功", zap.String("username", req.Username))
	response.ReturnData(c, gin.H{"message": "注册成功"})
}

// GetUserListAPI 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表，支持按用户名、姓名、手机号、部门ID、角色ID过滤
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username query string false "用户名"
// @Param name query string false "姓名"
// @Param phone query string false "手机号"
// @Param department_id query int false "部门ID"
// @Param role_id query int false "角色ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "页面大小" default(10)
// @Success 200 {object} response.ResponseDataWithCount{data=[]docs.SystemUserWithRelations} "查询成功"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 401 {object} response.ResponseError "未授权"
// @Router /users [get]
func GetUserListAPI(c *gin.Context) {
	var req service.GetUserListRequest
	if !middleware.BindAndValidate(c, &req) {
		return
	}

	// 处理分页参数
	if req.Page == 0 {
		req.Page = middleware.GetPage(c)
	}
	if req.PageSize == 0 {
		req.PageSize = middleware.GetPageSize(c)
	}

	// 调用服务层
	users, total, err := service.GetUserService().GetUserList(&req)
	if err != nil {
		middleware.HandleServiceError(c, err)
		return
	}

	response.ReturnDataWithCount(c, int(total), users)
}

// GetUserDetail 获取用户详情
// @Summary 获取用户详情
// @Description 根据用户ID获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.ResponseData{data=docs.SystemUserWithRelations} "查询成功"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 401 {object} response.ResponseError "未授权"
// @Failure 404 {object} response.ResponseError "用户不存在"
// @Router /users/{id} [get]
func GetUserDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的用户ID")
		return
	}

	// 调用服务层
	userInfo, err := service.GetUserService().GetUserByID(uint(id))
	if err != nil {
		middleware.HandleServiceError(c, err)
		return
	}

	response.ReturnData(c, userInfo)
}

// CreateUserAPI 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body service.CreateUserRequest true "用户信息"
// @Success 200 {object} response.ResponseData "创建成功"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 401 {object} response.ResponseError "未授权"
// @Failure 409 {object} response.ResponseError "用户已存在"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /users [post]
func CreateUserAPI(c *gin.Context) {
	var req service.CreateUserRequest
	if !middleware.BindAndValidate(c, &req) {
		return
	}

	// 调用服务层
	err := service.GetUserService().CreateUser(&req)
	if err != nil {
		middleware.HandleServiceError(c, err)
		return
	}

	zap.L().Info("用户创建成功", zap.String("username", req.Username))
	response.ReturnData(c, gin.H{"message": "用户创建成功"})
}

// UpdateUserAPI 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param body body service.UpdateUserRequest true "用户信息"
// @Success 200 {object} response.ResponseData "更新成功"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 401 {object} response.ResponseError "未授权"
// @Failure 404 {object} response.ResponseError "用户不存在"
// @Failure 409 {object} response.ResponseError "用户名已存在"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /users/{id} [put]
func UpdateUserAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的用户ID")
		return
	}

	var req service.UpdateUserRequest
	if !middleware.BindAndValidate(c, &req) {
		return
	}

	// 调用服务层
	err = service.GetUserService().UpdateUser(uint(id), &req)
	if err != nil {
		middleware.HandleServiceError(c, err)
		return
	}

	zap.L().Info("用户更新成功", zap.Uint("id", uint(id)), zap.String("username", req.Username))
	response.ReturnData(c, gin.H{"message": "用户更新成功"})
}

// DeleteUserAPI 删除用户
// @Summary 删除用户
// @Description 删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.ResponseData "删除成功"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 401 {object} response.ResponseError "未授权"
// @Failure 403 {object} response.ResponseError "禁止删除"
// @Failure 404 {object} response.ResponseError "用户不存在"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /users/{id} [delete]
func DeleteUserAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的用户ID")
		return
	}

	// 调用服务层
	err = service.GetUserService().DeleteUser(uint(id))
	if err != nil {
		middleware.HandleServiceError(c, err)
		return
	}

	zap.L().Info("用户删除成功", zap.Uint("id", uint(id)))
	response.ReturnData(c, gin.H{"message": "用户删除成功"})
}

// ChangePasswordRequest 修改密码请求参数
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" form:"old_password" validate:"required,min=6,max=50" example:"oldpassword123"`
	NewPassword string `json:"new_password" form:"new_password" validate:"required,min=6,max=50" example:"newpassword123"`
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 用户修改自己的密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.ResponseData "修改成功"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 401 {object} response.ResponseError "未授权"
// @Failure 403 {object} response.ResponseError "原密码错误"
// @Failure 500 {object} response.ResponseError "服务器错误"
// @Router /users/password [post]
func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if !middleware.BindAndValidate(c, &req) {
		return
	}

	// 从JWT中获取用户ID
	userIDStr := c.GetString(middleware.JWTDataKey)
	if userIDStr == "" {
		response.ReturnError(c, response.UNAUTHENTICATED, "未找到用户信息")
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "用户ID格式错误")
		return
	}

	// 调用服务层
	err = service.GetUserService().ChangePassword(uint(userID), req.OldPassword, req.NewPassword)
	if err != nil {
		middleware.HandleServiceError(c, err)
		return
	}

	zap.L().Info("用户密码修改成功", zap.Uint("user_id", uint(userID)))
	response.ReturnData(c, gin.H{"message": "密码修改成功"})
}

// FindLoginLogList 查询登录日志列表
// @Summary 查询登录日志列表
// @Description 查询用户登录日志列表，支持按IP和用户名过滤
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ip query string false "IP地址"
// @Param username query string false "用户名"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "页面大小" default(10)
// @Success 200 {object} response.ResponseDataWithCount{data=[]docs.SystemUserLoginLog} "查询成功"
// @Failure 400 {object} response.ResponseError "参数错误"
// @Failure 401 {object} response.ResponseError "未授权"
// @Router /users/login-logs [get]
func FindLoginLogList(c *gin.Context) {
	var req service.GetLoginLogListRequest
	if !middleware.BindAndValidate(c, &req) {
		return
	}

	// 处理分页参数
	if req.Page == 0 {
		req.Page = middleware.GetPage(c)
	}
	if req.PageSize == 0 {
		req.PageSize = middleware.GetPageSize(c)
	}

	// 调用服务层
	logs, total, err := service.GetUserService().GetLoginLogList(&req)
	if err != nil {
		middleware.HandleServiceError(c, err)
		return
	}

	response.ReturnDataWithCount(c, int(total), logs)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出，清除token
// @Tags 用户认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.ResponseData "登出成功"
// @Failure 401 {object} response.ResponseError "未授权"
// @Router /auth/logout [post]
func Logout(c *gin.Context) {
	// 这里可以实现token黑名单逻辑
	// 暂时返回成功
	zap.L().Info("用户登出", zap.String("user_agent", c.GetHeader("User-Agent")))
	response.ReturnData(c, gin.H{"message": "登出成功"})
}
