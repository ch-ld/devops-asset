package dns

import (
	"strconv"
	"time"

	"api-server/internal/model/dns"
	"api-server/internal/response/response"
	svc "api-server/internal/service/dns"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// CreateProviderRequest DNS提供商创建请求体结构体
type CreateProviderRequest struct {
	Name          string            `json:"name" binding:"required" example:"阿里云DNS" comment:"提供商名称"`
	Type          string            `json:"type" binding:"required" example:"aliyun" comment:"提供商类型"`
	Status        string            `json:"status" example:"active" comment:"状态(active/inactive)"`
	Credentials   map[string]string `json:"credentials" binding:"required" comment:"凭证信息"`
	Configuration datatypes.JSON    `json:"configuration" example:"{\"endpoint\":\"https://alidns.aliyuncs.com\"}" comment:"配置信息"`
	IsDefault     bool              `json:"is_default" example:"false" comment:"是否默认提供商"`
	Priority      int               `json:"priority" example:"0" comment:"优先级"`
	RateLimit     int               `json:"rate_limit" example:"10" comment:"速率限制(req/s)"`
	Remark        string            `json:"remark" example:"主要DNS提供商" comment:"备注"`
}

// UpdateProviderRequest DNS提供商更新请求体结构体
type UpdateProviderRequest struct {
	Name          string            `json:"name" binding:"required" example:"阿里云DNS" comment:"提供商名称"`
	Type          string            `json:"type" binding:"required" example:"aliyun" comment:"提供商类型"`
	Status        string            `json:"status" example:"active" comment:"状态(active/inactive)"`
	Credentials   map[string]string `json:"credentials,omitempty" comment:"凭证信息"`
	Configuration datatypes.JSON    `json:"configuration" example:"{\"endpoint\":\"https://alidns.aliyuncs.com\"}" comment:"配置信息"`
	IsDefault     bool              `json:"is_default" example:"false" comment:"是否默认提供商"`
	Priority      int               `json:"priority" example:"0" comment:"优先级"`
	RateLimit     int               `json:"rate_limit" example:"10" comment:"速率限制(req/s)"`
	Remark        string            `json:"remark" example:"主要DNS提供商" comment:"备注"`
}

// ProviderHandler DNS提供商处理器
type ProviderHandler struct {
	providerService *svc.ProviderService
}

// NewProviderHandler 创建DNS提供商处理器
func NewProviderHandler(providerService *svc.ProviderService) *ProviderHandler {
	return &ProviderHandler{
		providerService: providerService,
	}
}

// CreateProvider 创建DNS提供商
func (h *ProviderHandler) CreateProvider(c *gin.Context) {
	var req CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请求参数错误: "+err.Error())
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 构建DNS提供商对象
	provider := &dns.Provider{
		Name:          req.Name,
		Type:          req.Type,
		Status:        req.Status,
		Configuration: req.Configuration,
		IsDefault:     req.IsDefault,
		Priority:      req.Priority,
		RateLimit:     req.RateLimit,
		Remark:        req.Remark,
		TenantID:      tenantID.(uint),
		CreatedBy:     userID.(uint),
		UpdatedBy:     userID.(uint),
	}

	// TODO: 处理凭证信息加密
	// provider.CredentialsEnc = encryptCredentials(req.Credentials)

	// 创建DNS提供商
	if err := h.providerService.CreateProvider(provider, userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "创建DNS提供商失败: "+err.Error())
		return
	}

	// 返回创建的提供商信息
	providerResp := h.convertToResponse(provider)
	response.ReturnData(c, providerResp)
}

// UpdateProvider 更新DNS提供商
func (h *ProviderHandler) UpdateProvider(c *gin.Context) {
	// 获取提供商ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的提供商ID")
		return
	}

	var req UpdateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请求参数错误: "+err.Error())
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 构建DNS提供商对象
	provider := &dns.Provider{
		Name:          req.Name,
		Type:          req.Type,
		Status:        req.Status,
		Configuration: req.Configuration,
		IsDefault:     req.IsDefault,
		Priority:      req.Priority,
		RateLimit:     req.RateLimit,
		Remark:        req.Remark,
		TenantID:      tenantID.(uint),
		UpdatedBy:     userID.(uint),
	}
	provider.ID = uint(id)

	// 更新DNS提供商
	if err := h.providerService.UpdateProvider(provider, userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "更新DNS提供商失败: "+err.Error())
		return
	}

	// 获取更新后的提供商信息
	updatedProvider, err := h.providerService.GetProvider(uint(id))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取更新后的提供商信息失败: "+err.Error())
		return
	}

	// 返回更新后的提供商信息
	providerResp := h.convertToResponse(updatedProvider)
	response.ReturnData(c, providerResp)
}

// DeleteProvider 删除DNS提供商
func (h *ProviderHandler) DeleteProvider(c *gin.Context) {
	// 获取提供商ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的提供商ID")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 删除DNS提供商
	if err := h.providerService.DeleteProvider(uint(id), userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "删除DNS提供商失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// GetProvider 获取DNS提供商详情
func (h *ProviderHandler) GetProvider(c *gin.Context) {
	// 获取提供商ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的提供商ID")
		return
	}

	// 获取DNS提供商信息
	provider, err := h.providerService.GetProvider(uint(id))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "DNS提供商不存在: "+err.Error())
		return
	}

	// 返回提供商信息
	providerResp := h.convertToResponse(provider)
	response.ReturnData(c, providerResp)
}

// ListProviders 获取DNS提供商列表
func (h *ProviderHandler) ListProviders(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// 获取筛选参数
	name := c.Query("name")
	providerType := c.Query("type")
	status := c.Query("status")

	// 获取租户信息
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 构建筛选条件
	filters := make(map[string]interface{})
	if name != "" {
		filters["name"] = name
	}
	if providerType != "" {
		filters["type"] = providerType
	}
	if status != "" {
		filters["status"] = status
	}

	// 获取DNS提供商列表
	providers, total, err := h.providerService.ListProviders(tenantID.(uint), filters, size, (page-1)*size)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取DNS提供商列表失败: "+err.Error())
		return
	}

	// 转换响应格式
	items := make([]*ProviderResponse, len(providers))
	for i, provider := range providers {
		items[i] = h.convertToResponse(provider)
	}

	resp := map[string]interface{}{
		"total": int(total),
		"items": items,
	}

	response.ReturnData(c, resp)
}

// TestProvider 测试DNS提供商连接
func (h *ProviderHandler) TestProvider(c *gin.Context) {
	// 获取提供商ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的提供商ID")
		return
	}

	// 测试DNS提供商连接
	result, err := h.providerService.TestProvider(uint(id))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "测试DNS提供商连接失败: "+err.Error())
		return
	}

	// 转换响应格式
	testResp := &TestConnectionResponse{
		Success:  result.Success,
		Latency:  result.Latency.String(),
		ErrorMsg: "",
		Details:  make(map[string]string),
		TestedAt: time.Now(),
	}

	response.ReturnData(c, testResp)
}

// SetDefaultProvider 设置默认DNS提供商
func (h *ProviderHandler) SetDefaultProvider(c *gin.Context) {
	// 获取提供商ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的提供商ID")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 设置默认DNS提供商
	if err := h.providerService.SetDefaultProvider(uint(id), userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "设置默认DNS提供商失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// SyncProviderDomains 同步单个提供商的域名
func (h *ProviderHandler) SyncProviderDomains(c *gin.Context) {
	// 获取提供商ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的提供商ID")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 同步提供商域名
	result, err := h.providerService.SyncProviderDomains(uint(id), userID.(uint), c.ClientIP())
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "同步域名失败: "+err.Error())
		return
	}

	response.ReturnData(c, result)
}

// SyncAllProviderDomains 同步所有提供商的域名
func (h *ProviderHandler) SyncAllProviderDomains(c *gin.Context) {
	// 获取租户信息
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 同步所有提供商域名
	result, err := h.providerService.SyncAllProviderDomains(tenantID.(uint), userID.(uint), c.ClientIP())
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "同步所有域名失败: "+err.Error())
		return
	}

	response.ReturnData(c, result)
}

// convertToResponse 转换为响应格式
func (h *ProviderHandler) convertToResponse(provider *dns.Provider) *ProviderResponse {
	var lastTestAt *time.Time
	if provider.LastTestAt != nil && !provider.LastTestAt.Time.IsZero() {
		lastTestAt = &provider.LastTestAt.Time
	}

	return &ProviderResponse{
		ID:         provider.ID,
		Name:       provider.Name,
		Type:       provider.Type,
		Status:     provider.Status,
		IsDefault:  provider.IsDefault,
		Priority:   provider.Priority,
		RateLimit:  provider.RateLimit,
		Concurrent: provider.Concurrent,
		Timeout:    provider.Timeout,
		LastTestAt: lastTestAt,
		Remark:     provider.Remark,
		TenantID:   provider.TenantID,
		CreatedBy:  provider.CreatedBy,
		UpdatedBy:  provider.UpdatedBy,
		CreatedAt:  provider.CreatedAt,
		UpdatedAt:  provider.UpdatedAt,
	}
}
