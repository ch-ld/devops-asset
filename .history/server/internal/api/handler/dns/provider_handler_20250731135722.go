package dns

import (
	"net/http"
	"strconv"

	"api-server/internal/model/dns"
	"api-server/internal/response/response"
	svc "api-server/internal/service/dns"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// CreateProviderRequest DNS提供商创建请求体结构体
// @Description DNS提供商创建时的请求参数
type CreateProviderRequest struct {
	Name          string            `json:"name" binding:"required" example:"阿里云DNS" comment:"提供商名称"`
	Type          string            `json:"type" binding:"required" example:"aliyun" comment:"提供商类型"`
	Status        string            `json:"status" example:"active" comment:"状态(active/inactive)"`
	Credentials   map[string]string `json:"credentials" binding:"required" comment:"凭证信息"`
	Configuration datatypes.JSON    `json:"configuration" example:"{\"endpoint\":\"https://alidns.aliyuncs.com\"}" comment:"配置信息"`
	IsDefault     bool              `json:"is_default" example:"false" comment:"是否默认提供�?`
	Priority      int               `json:"priority" example:"0" comment:"优先�?`
	RateLimit     int               `json:"rate_limit" example:"10" comment:"速率限制(req/s)"`
	Concurrent    int               `json:"concurrent" example:"5" comment:"并发�?`
	Timeout       int               `json:"timeout" example:"30" comment:"超时时间(�?"`
	Remark        string            `json:"remark" example:"主要DNS提供�? comment:"备注"`
}

// UpdateProviderRequest DNS提供商更新请求体结构�?// @Description DNS提供商更新时的请求参�?type UpdateProviderRequest struct {
	Name          string            `json:"name" binding:"required" example:"阿里云DNS" comment:"提供商名�?`
	Type          string            `json:"type" binding:"required" example:"aliyun" comment:"提供商类�?`
	Status        string            `json:"status" example:"active" comment:"状�?active/inactive)"`
	Credentials   map[string]string `json:"credentials,omitempty" comment:"凭证信息"`
	Configuration datatypes.JSON    `json:"configuration" example:"{\"endpoint\":\"https://alidns.aliyuncs.com\"}" comment:"配置信息"`
	IsDefault     bool              `json:"is_default" example:"false" comment:"是否默认提供�?`
	Priority      int               `json:"priority" example:"0" comment:"优先�?`
	RateLimit     int               `json:"rate_limit" example:"10" comment:"速率限制(req/s)"`
	Concurrent    int               `json:"concurrent" example:"5" comment:"并发�?`
	Timeout       int               `json:"timeout" example:"30" comment:"超时时间(�?"`
	Remark        string            `json:"remark" example:"主要DNS提供�? comment:"备注"`
}

// ProviderListResponse DNS提供商列表响应结构体
// @Description DNS提供商列表响�?type ProviderListResponse struct {
	Total int                 `json:"total" example:"10" comment:"总数�?`
	Items []*ProviderResponse `json:"items" comment:"提供商列�?`
}

// ProviderHandler DNS提供商管理接口处理器
// @Description DNS提供商管理相关接口处理器，负责DNS提供商的增删改查等HTTP请求处理
type ProviderHandler struct {
	providerService *svc.ProviderService
}

// NewProviderHandler 创建DNS提供商处理器实例
func NewProviderHandler(providerService *svc.ProviderService) *ProviderHandler {
	return &ProviderHandler{
		providerService: providerService,
	}
}

// CreateProvider 创建DNS提供�?// @Summary 创建DNS提供�?// @Description 创建新的DNS提供商配�?// @Tags DNS提供商管�?// @Accept json
// @Produce json
// @Param request body CreateProviderRequest true "DNS提供商创建请�?
// @Success 200 {object} response.Response{data=ProviderResponse} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错�?
// @Router /api/v1/dns/providers [post]
func (h *ProviderHandler) CreateProvider(c *gin.Context) {
	var req CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请求参数错误: "+err.Error())
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登�?)
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 构建DNS提供商对�?	provider := &dns.Provider{
		Name:          req.Name,
		Type:          req.Type,
		Status:        req.Status,
		Configuration: req.Configuration,
		IsDefault:     req.IsDefault,
		Priority:      req.Priority,
		RateLimit:     req.RateLimit,
		Concurrent:    req.Concurrent,
		Timeout:       req.Timeout,
		Remark:        req.Remark,
		TenantID:      tenantID.(uint),
	}

	// 设置默认�?	if provider.Status == "" {
		provider.Status = "active"
	}
	if provider.RateLimit == 0 {
		provider.RateLimit = 10
	}
	if provider.Concurrent == 0 {
		provider.Concurrent = 5
	}
	if provider.Timeout == 0 {
		provider.Timeout = 30
	}

	// TODO: 处理凭证信息加密
	// provider.CredentialsEnc = encryptCredentials(req.Credentials)

	// 创建DNS提供�?	if err := h.providerService.CreateProvider(provider, userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "创建DNS提供商失�? "+err.Error())
		return
	}

	// 返回创建的提供商信息
	providerResp := h.convertToResponse(provider)
	response.ReturnData(c, providerResp)
}

// UpdateProvider 更新DNS提供�?// @Summary 更新DNS提供�?// @Description 更新指定ID的DNS提供商信�?// @Tags DNS提供商管�?// @Accept json
// @Produce json
// @Param id path int true "提供商ID"
// @Param request body UpdateProviderRequest true "DNS提供商更新请�?
// @Success 200 {object} response.Response{data=ProviderResponse} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "提供商不存在"
// @Failure 500 {object} response.Response "服务器内部错�?
// @Router /api/v1/dns/providers/{id} [put]
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
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登�?)
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 构建DNS提供商对�?	provider := &dns.Provider{
		Model:         dns.Model{ID: uint(id)},
		Name:          req.Name,
		Type:          req.Type,
		Status:        req.Status,
		Configuration: req.Configuration,
		IsDefault:     req.IsDefault,
		Priority:      req.Priority,
		RateLimit:     req.RateLimit,
		Concurrent:    req.Concurrent,
		Timeout:       req.Timeout,
		Remark:        req.Remark,
		TenantID:      tenantID.(uint),
	}

	// TODO: 处理凭证信息加密
	// if len(req.Credentials) > 0 {
	//     provider.CredentialsEnc = encryptCredentials(req.Credentials)
	// }

	// 更新DNS提供�?	if err := h.providerService.UpdateProvider(provider, userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "更新DNS提供商失�? "+err.Error())
		return
	}

	// 获取更新后的提供商信�?	updatedProvider, err := h.providerService.GetProvider(uint(id))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取更新后的提供商信息失�? "+err.Error())
		return
	}

	// 返回更新后的提供商信�?	providerResp := h.convertToResponse(updatedProvider)
	response.ReturnData(c, providerResp)
}

// DeleteProvider 删除DNS提供�?// @Summary 删除DNS提供�?// @Description 删除指定ID的DNS提供�?// @Tags DNS提供商管�?// @Produce json
// @Param id path int true "提供商ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "提供商不存在"
// @Failure 500 {object} response.Response "服务器内部错�?
// @Router /api/v1/dns/providers/{id} [delete]
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
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登�?)
		return
	}

	// 删除DNS提供�?	if err := h.providerService.DeleteProvider(uint(id), userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "删除DNS提供商失�? "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// GetProvider 获取DNS提供商详�?// @Summary 获取DNS提供商详�?// @Description 根据ID获取DNS提供商详细信�?// @Tags DNS提供商管�?// @Produce json
// @Param id path int true "提供商ID"
// @Success 200 {object} response.Response{data=ProviderResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "提供商不存在"
// @Failure 500 {object} response.Response "服务器内部错�?
// @Router /api/v1/dns/providers/{id} [get]
func (h *ProviderHandler) GetProvider(c *gin.Context) {
	// 获取提供商ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的提供商ID")
		return
	}

	// 获取DNS提供商信�?	provider, err := h.providerService.GetProvider(uint(id))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "DNS提供商不存在: "+err.Error())
		return
	}

	// 返回提供商信�?	providerResp := h.convertToResponse(provider)
	response.ReturnData(c, providerResp)
}

// ListProviders 获取DNS提供商列�?// @Summary 获取DNS提供商列�?// @Description 获取DNS提供商列表，支持分页和筛�?// @Tags DNS提供商管�?// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param keyword query string false "关键词搜�?
// @Param type query string false "提供商类型筛�?
// @Param status query string false "状态筛�?
// @Success 200 {object} response.Response{data=ProviderListResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错�?
// @Router /api/v1/dns/providers [get]
func (h *ProviderHandler) ListProviders(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	// 获取租户信息
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 构建筛选条�?	filters := make(map[string]interface{})
	if keyword := c.Query("keyword"); keyword != "" {
		filters["keyword"] = keyword
	}
	if providerType := c.Query("type"); providerType != "" {
		filters["type"] = providerType
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	// 获取DNS提供商列�?	providers, total, err := h.providerService.ListProviders(tenantID.(uint), filters, size, (page-1)*size)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取DNS提供商列表失�? "+err.Error())
		return
	}

	// 转换响应格式
	items := make([]*ProviderResponse, len(providers))
	for i, provider := range providers {
		items[i] = h.convertToResponse(provider)
	}

	resp := &ProviderListResponse{
		Total: int(total),
		Items: items,
	}

	response.ReturnData(c, resp)
}

// TestProvider 测试DNS提供商连�?// @Summary 测试DNS提供商连�?// @Description 测试指定DNS提供商的连接状�?// @Tags DNS提供商管�?// @Produce json
// @Param id path int true "提供商ID"
// @Success 200 {object} response.Response{data=TestConnectionResponse} "测试成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "提供商不存在"
// @Failure 500 {object} response.Response "服务器内部错�?
// @Router /api/v1/dns/providers/{id}/test [post]
func (h *ProviderHandler) TestProvider(c *gin.Context) {
	// 获取提供商ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的提供商ID")
		return
	}

	// 测试DNS提供商连�?	result, err := h.providerService.TestProvider(uint(id))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "测试DNS提供商连接失�? "+err.Error())
		return
	}

	// 转换响应格式
	testResp := &TestConnectionResponse{
		Success:    result.Success,
		Latency:    result.Latency.String(),
		ErrorMsg:   result.ErrorMsg,
		Details:    result.Details,
		TestedAt:   result.TestedAt,
		TestType:   result.TestType,
		Endpoint:   result.Endpoint,
		StatusCode: result.StatusCode,
	}

	response.ReturnData(c, testResp)
}

// SetDefaultProvider 设置默认DNS提供�?// @Summary 设置默认DNS提供�?// @Description 设置指定的DNS提供商为默认提供�?// @Tags DNS提供商管�?// @Produce json
// @Param id path int true "提供商ID"
// @Success 200 {object} response.Response "设置成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "提供商不存在"
// @Failure 500 {object} response.Response "服务器内部错�?
// @Router /api/v1/dns/providers/{id}/default [post]
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
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登�?)
		return
	}

	// 设置默认DNS提供�?	if err := h.providerService.SetDefaultProvider(uint(id), userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "设置默认DNS提供商失�? "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// convertToResponse 转换为响应格�?func (h *ProviderHandler) convertToResponse(provider *dns.Provider) *ProviderResponse {
	return &ProviderResponse{
		ID:           provider.ID,
		Name:         provider.Name,
		Type:         provider.Type,
		Status:       provider.Status,
		IsDefault:    provider.IsDefault,
		Priority:     provider.Priority,
		RateLimit:    provider.RateLimit,
		Concurrent:   provider.Concurrent,
		Timeout:      provider.Timeout,
		LastTestAt:   provider.LastTestAt,
		TestResult:   provider.TestResult,
		ErrorMessage: provider.ErrorMessage,
		Remark:       provider.Remark,
		TenantID:     provider.TenantID,
		CreatedBy:    provider.CreatedBy,
		UpdatedBy:    provider.UpdatedBy,
		CreatedAt:    provider.CreatedAt,
		UpdatedAt:    provider.UpdatedAt,
	}
}
