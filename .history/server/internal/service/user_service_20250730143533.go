package service

import (
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/internal/config"
	"api-server/internal/db/mysql"
	"api-server/internal/db/mysql/system"
	"api-server/internal/response/response"
	"api-server/pkg/auth/authentication"
	"api-server/pkg/crypto/encryption"
)

// UserService 用户服务接口
type UserService interface {
	// 认证相关
	Login(username, password string) (string, error)
	VerifyUser(username, password string) (*system.SystemUser, error)
	Register(req *RegisterRequest) error

	// 用户管理
	CreateUser(req *CreateUserRequest) error
	UpdateUser(id uint, req *UpdateUserRequest) error
	DeleteUser(id uint) error
	GetUserByID(id uint) (*system.UserWithRelations, error)
	GetUserList(req *GetUserListRequest) ([]system.UserWithRelations, int64, error)

	// 密码管理
	ChangePassword(userID uint, oldPassword, newPassword string) error

	// 日志管理
	CreateLoginLog(username, ip string, status uint) error
	GetLoginLogList(req *GetLoginLogListRequest) ([]system.SystemUserLoginLog, int64, error)
}

// userService 用户服务实现
type userService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService() UserService {
	return &userService{
		db: mysql.GetClient(),
	}
}

// 请求参数结构体
type CreateUserRequest struct {
	Username     string `json:"username" validate:"required,min=3,max=50" example:"john_doe"`
	Password     string `json:"password" validate:"required,min=6,max=50" example:"password123"`
	Name         string `json:"name" validate:"required,min=2,max=50" example:"John Doe"`
	Phone        string `json:"phone" validate:"required,len=11" example:"13800138000"`
	Gender       uint   `json:"gender" validate:"required,oneof=1 2" example:"1"`
	Status       uint   `json:"status" validate:"required,oneof=1 2" example:"1"`
	RoleID       uint   `json:"role_id" validate:"required,gt=0" example:"1"`
	DepartmentID uint   `json:"department_id" validate:"required,gt=0" example:"1"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"john_doe"`
	Password string `json:"password" validate:"required,min=6,max=50" example:"password123"`
	Name     string `json:"name" validate:"required,min=2,max=50" example:"John Doe"`
	Phone    string `json:"phone" validate:"required,len=11" example:"13800138000"`
	Gender   uint   `json:"gender" validate:"required,oneof=1 2" example:"1"`
}

type UpdateUserRequest struct {
	Username     string `json:"username" validate:"required,min=3,max=50" example:"john_doe"`
	Password     string `json:"password" validate:"omitempty,min=6,max=50" example:"password123"`
	Name         string `json:"name" validate:"required,min=2,max=50" example:"John Doe"`
	Phone        string `json:"phone" validate:"required,len=11" example:"13800138000"`
	Gender       uint   `json:"gender" validate:"required,oneof=1 2" example:"1"`
	Status       uint   `json:"status" validate:"required,oneof=1 2" example:"1"`
	RoleID       uint   `json:"role_id" validate:"required,gt=0" example:"1"`
	DepartmentID uint   `json:"department_id" validate:"required,gt=0" example:"1"`
}

type GetUserListRequest struct {
	Username     string `json:"username" validate:"omitempty,max=50"`
	Name         string `json:"name" validate:"omitempty,max=50"`
	Phone        string `json:"phone" validate:"omitempty,len=11"`
	DepartmentID uint   `json:"department_id" validate:"omitempty,gt=0"`
	RoleID       uint   `json:"role_id" validate:"omitempty,gt=0"`
	Page         int    `json:"page" validate:"omitempty,min=1"`
	PageSize     int    `json:"page_size" validate:"omitempty,min=1,max=100"`
}

type GetLoginLogListRequest struct {
	IP       string `json:"ip" validate:"omitempty,max=50"`
	Username string `json:"username" validate:"omitempty,max=50"`
	Page     int    `json:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size" validate:"omitempty,min=1,max=100"`
}

// Login 用户登录
func (s *userService) Login(username, password string) (string, error) {
	// 验证用户
	user, err := s.VerifyUser(username, password)
	if err != nil {
		zap.L().Error("验证用户失败", zap.Error(err), zap.String("username", username))
		return "", response.NewServiceError(response.DATA_LOSS, "用户验证失败")
	}

	if user.ID == 0 {
		return "", response.NewServiceError(response.INVALID_ARGUMENT, "用户名或密码错误")
	}

	if user.Status != 1 {
		return "", response.NewServiceError(response.PERMISSION_DENIED, "用户已被禁用")
	}

	// 生成token
	token, err := authentication.JWTIssue(fmt.Sprintf("%d", user.ID))
	if err != nil {
		zap.L().Error("生成token失败", zap.Error(err), zap.Uint("user_id", user.ID))
		return "", response.NewServiceError(response.INTERNAL, "生成令牌失败")
	}

	return token, nil
}

// VerifyUser 验证用户
func (s *userService) VerifyUser(username, password string) (*system.SystemUser, error) {
	user := &system.SystemUser{}
	encryptedPassword := encryption.MD5WithSalt(config.PWDSalt + password)

	err := s.db.Where(&system.SystemUser{
		Username: username,
		Password: encryptedPassword,
	}).First(user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, nil
		}
		zap.L().Error("查询用户失败", zap.Error(err), zap.String("username", username))
		return nil, err
	}

	return user, nil
}

// CreateUser 创建用户
func (s *userService) CreateUser(req *CreateUserRequest) error {
	// 检查用户名是否已存在
	var count int64
	if err := s.db.Model(&system.SystemUser{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		zap.L().Error("检查用户名失败", zap.Error(err), zap.String("username", req.Username))
		return response.NewServiceError(response.DATA_LOSS, "检查用户名失败")
	}

	if count > 0 {
		return response.NewServiceError(response.ALREADY_EXISTS, "用户名已存在")
	}

	// 检查手机号是否已存在
	if err := s.db.Model(&system.SystemUser{}).Where("phone = ?", req.Phone).Count(&count).Error; err != nil {
		zap.L().Error("检查手机号失败", zap.Error(err), zap.String("phone", req.Phone))
		return response.NewServiceError(response.DATA_LOSS, "检查手机号失败")
	}

	if count > 0 {
		return response.NewServiceError(response.ALREADY_EXISTS, "手机号已存在")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建用户
	user := &system.SystemUser{
		Username:     req.Username,
		Password:     encryption.MD5WithSalt(config.PWDSalt + req.Password),
		Name:         req.Name,
		Phone:        req.Phone,
		Gender:       req.Gender,
		Status:       req.Status,
		RoleID:       req.RoleID,
		DepartmentID: req.DepartmentID,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		zap.L().Error("创建用户失败", zap.Error(err), zap.String("username", req.Username))
		return response.NewServiceError(response.DATA_LOSS, "创建用户失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		zap.L().Error("提交事务失败", zap.Error(err))
		return response.NewServiceError(response.INTERNAL, "提交事务失败")
	}

	// 缓存用户信息
	if err := CacheSvc.CacheUserInfo(user.ID); err != nil {
		zap.L().Warn("缓存用户信息失败", zap.Error(err), zap.Uint("user_id", user.ID))
	}

	zap.L().Info("用户创建成功", zap.String("username", req.Username), zap.Uint("user_id", user.ID))
	return nil
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(id uint, req *UpdateUserRequest) error {
	if id == 0 {
		return response.NewServiceError(response.INVALID_ARGUMENT, "无效的用户ID")
	}

	// 检查用户是否存在
	var user system.SystemUser
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewServiceError(response.NOT_FOUND, "用户不存在")
		}
		zap.L().Error("查询用户失败", zap.Error(err), zap.Uint("id", id))
		return response.NewServiceError(response.DATA_LOSS, "查询用户失败")
	}

	// 检查用户名是否被其他用户使用
	var count int64
	if err := s.db.Model(&system.SystemUser{}).Where("username = ? AND id != ?", req.Username, id).Count(&count).Error; err != nil {
		zap.L().Error("检查用户名失败", zap.Error(err), zap.String("username", req.Username))
		return response.NewServiceError(response.DATA_LOSS, "检查用户名失败")
	}

	if count > 0 {
		return response.NewServiceError(response.ALREADY_EXISTS, "用户名已被其他用户使用")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新用户信息
	updateData := map[string]interface{}{
		"username":      req.Username,
		"name":          req.Name,
		"phone":         req.Phone,
		"gender":        req.Gender,
		"status":        req.Status,
		"role_id":       req.RoleID,
		"department_id": req.DepartmentID,
	}

	if req.Password != "" {
		updateData["password"] = encryption.MD5WithSalt(config.PWDSalt + req.Password)
	}

	if err := tx.Model(&user).Updates(updateData).Error; err != nil {
		tx.Rollback()
		zap.L().Error("更新用户失败", zap.Error(err), zap.Uint("id", id))
		return response.NewServiceError(response.DATA_LOSS, "更新用户失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		zap.L().Error("提交事务失败", zap.Error(err))
		return response.NewServiceError(response.INTERNAL, "提交事务失败")
	}

	// 更新缓存
	if err := CacheSvc.CacheUserInfo(id); err != nil {
		zap.L().Warn("更新用户缓存失败", zap.Error(err), zap.Uint("id", id))
	}

	zap.L().Info("用户更新成功", zap.Uint("id", id), zap.String("username", req.Username))
	return nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id uint) error {
	if id == 0 {
		return response.NewServiceError(response.INVALID_ARGUMENT, "无效的用户ID")
	}

	if id == 1 {
		return response.NewServiceError(response.PERMISSION_DENIED, "不能删除超级管理员")
	}

	// 检查用户是否存在
	var user system.SystemUser
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewServiceError(response.NOT_FOUND, "用户不存在")
		}
		zap.L().Error("查询用户失败", zap.Error(err), zap.Uint("id", id))
		return response.NewServiceError(response.DATA_LOSS, "查询用户失败")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 软删除用户
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		zap.L().Error("删除用户失败", zap.Error(err), zap.Uint("id", id))
		return response.NewServiceError(response.DATA_LOSS, "删除用户失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		zap.L().Error("提交事务失败", zap.Error(err))
		return response.NewServiceError(response.INTERNAL, "提交事务失败")
	}

	// 删除缓存
	if err := CacheSvc.DeleteUserFromCache(id); err != nil {
		zap.L().Warn("删除用户缓存失败", zap.Error(err), zap.Uint("id", id))
	}

	zap.L().Info("用户删除成功", zap.Uint("id", id), zap.String("username", user.Username))
	return nil
}

// GetUserByID 根据ID获取用户详情
func (s *userService) GetUserByID(id uint) (*system.UserWithRelations, error) {
	if id == 0 {
		return nil, response.NewServiceError(response.INVALID_ARGUMENT, "无效的用户ID")
	}

	userInfo, err := CacheSvc.GetUserFromCache(id)
	if err != nil {
		zap.L().Error("获取用户信息失败", zap.Error(err), zap.Uint("id", id))
		return nil, response.NewServiceError(response.DATA_LOSS, "获取用户信息失败")
	}

	return userInfo, nil
}

// GetUserList 获取用户列表
func (s *userService) GetUserList(req *GetUserListRequest) ([]system.UserWithRelations, int64, error) {
	user := &system.SystemUser{
		Username:     req.Username,
		Name:         req.Name,
		Phone:        req.Phone,
		RoleID:       req.RoleID,
		DepartmentID: req.DepartmentID,
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	users, total, err := system.FindUserList(user, page, pageSize)
	if err != nil {
		zap.L().Error("查询用户列表失败", zap.Error(err))
		return nil, 0, response.NewServiceError(response.DATA_LOSS, "查询用户列表失败")
	}

	// 清空密码字段
	for i := range users {
		users[i].SystemUser.Password = ""
	}

	return users, total, nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	if userID == 0 {
		return response.NewServiceError(response.INVALID_ARGUMENT, "无效的用户ID")
	}

	// 查询用户
	var user system.SystemUser
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewServiceError(response.NOT_FOUND, "用户不存在")
		}
		zap.L().Error("查询用户失败", zap.Error(err), zap.Uint("user_id", userID))
		return response.NewServiceError(response.DATA_LOSS, "查询用户失败")
	}

	// 验证旧密码
	encryptedOldPassword := encryption.MD5WithSalt(config.PWDSalt + oldPassword)
	if user.Password != encryptedOldPassword {
		return response.NewServiceError(response.PERMISSION_DENIED, "原密码错误")
	}

	// 更新密码
	encryptedNewPassword := encryption.MD5WithSalt(config.PWDSalt + newPassword)
	if err := s.db.Model(&user).Update("password", encryptedNewPassword).Error; err != nil {
		zap.L().Error("更新密码失败", zap.Error(err), zap.Uint("user_id", userID))
		return response.NewServiceError(response.DATA_LOSS, "更新密码失败")
	}

	// 更新缓存
	if err := CacheSvc.CacheUserInfo(userID); err != nil {
		zap.L().Warn("更新用户缓存失败", zap.Error(err), zap.Uint("user_id", userID))
	}

	zap.L().Info("用户密码修改成功", zap.Uint("user_id", userID))
	return nil
}

// CreateLoginLog 创建登录日志
func (s *userService) CreateLoginLog(username, ip string, status uint) error {
	log := &system.SystemUserLoginLog{
		UserName: username,
		IP:       ip,
		Status:   status,
	}

	if err := s.db.Create(log).Error; err != nil {
		zap.L().Error("创建登录日志失败", zap.Error(err), zap.String("username", username))
		return response.NewServiceError(response.DATA_LOSS, "创建登录日志失败")
	}

	return nil
}

// GetLoginLogList 获取登录日志列表
func (s *userService) GetLoginLogList(req *GetLoginLogListRequest) ([]system.SystemUserLoginLog, int64, error) {
	log := &system.SystemUserLoginLog{
		IP:       req.IP,
		UserName: req.Username,
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	logs, total, err := system.FindLoginLogList(log, page, pageSize)
	if err != nil {
		zap.L().Error("查询登录日志失败", zap.Error(err))
		return nil, 0, response.NewServiceError(response.DATA_LOSS, "查询登录日志失败")
	}

	return logs, total, nil
}

// 全局服务实例（延迟初始化）
var (
	UserSvc  UserService
	userOnce sync.Once
)

// GetUserService 获取用户服务实例（延迟初始化）
func GetUserService() UserService {
	userOnce.Do(func() {
		UserSvc = NewUserService()
	})
	return UserSvc
}
