package cmdb

import (
	"fmt"
	"strconv"

	model "api-server/internal/model/cmdb"
	"api-server/internal/response/response"
	svc "api-server/internal/service/cmdb"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ProviderRequest 云账号请求结构体
// @Description 云账号创建或更新时的请求参数
type ProviderRequest struct {
	Name         string `json:"name" binding:"required" example:"阿里云账号" comment:"云账号名称"`
	Type         string `json:"type" binding:"required" example:"alicloud" comment:"云服务商类型：alicloud,aws,tencentcloud"`
	AccessKey    string `json:"access_key" binding:"required" example:"LTAI4..." comment:"访问密钥ID"`
	SecretKey    string `json:"secret_key" binding:"required" example:"abc123..." comment:"访问密钥Secret"`
	Region       string `json:"region" example:"cn-hangzhou" comment:"默认区域"`
	Status       string `json:"status" example:"active" comment:"状态：active,inactive"`
	Description  string `json:"description" example:"生产环境阿里云账号" comment:"描述信息"`
}

// ProviderResponse 云账号响应结构体
// @Description 云账号信息响应
type ProviderResponse struct {
	ID          uint   `json:"id" example:"1" comment:"云账号ID"`
	Name        string `json:"name" example:"阿里云账号" comment:"云账号名称"`
	Type        string `json:"type" example:"alicloud" comment:"云服务商类型"`
	AccessKey   string `json:"access_key" example:"LTAI4..." comment:"访问密钥ID"`
	Region      string `json:"region" example:"cn-hangzhou" comment:"默认区域"`
	Status      string `json:"status" example:"active" comment:"状态"`
	Description string `json:"description" example:"生产环境阿里云账号" comment:"描述信息"`
	CreatedAt   string `json:"created_at" example:"2024-01-01T00:00:00Z" comment:"创建时间"`
	UpdatedAt   string `json:"updated_at" example:"2024-01-01T00:00:00Z" comment:"更新时间"`
}

// ProviderListResponse 云账号列表响应结构体
// @Description 云账号列表响应
type ProviderListResponse struct {
	Total int                 `json:"total" example:"10" comment:"总数量"`
	Items []*ProviderResponse `json:"items" comment:"云账号列表"`
}

// SyncResourcesRequest 同步资源请求结构体
// @Description 同步云资源请求参数
type SyncResourcesRequest struct {
	GroupID *uint `json:"group_id" example:"1" comment:"目标主机组ID，可选"`
}

// SyncResourcesResponse 同步资源响应结构体
// @Description 同步云资源响应
type SyncResourcesResponse struct {
	Total    int `json:"total" example:"100" comment:"总同步数量"`
	Success  int `json:"success" example:"95" comment:"成功同步数量"`
	Failed   int `json:"failed" example:"5" comment:"失败同步数量"`
	Skipped  int `json:"skipped" example:"10" comment:"跳过数量"`
	Message  string `json:"message" example:"同步完成" comment:"同步结果消息"`
}

// ProviderHandler 云账号管理相关接口处理器
// @Description 负责云账号的增删改查、密钥加密等HTTP请求处理
type ProviderHandler struct {
	providerSvc *svc.ProviderService
}

// NewProviderHandler 创建云账号处理器实例
// @Description 创建云账号管理处理器的新实例
// @Param providerSvc 云账号服务实例
// @Return *ProviderHandler 云账号处理器实例
func NewProviderHandler(providerSvc *svc.ProviderService) *ProviderHandler {
	return &ProviderHandler{providerSvc: providerSvc}
}

// CreateProvider 创建云账号
// @Summary 创建新云账号
// @Description 创建一个新的云账号配置，支持阿里云、腾讯云、AWS等
// @Tags 云账号管理
// @Accept json
// @Produce json
// @Param provider body ProviderRequest true "云账号信息"
// @Success 200 {object} response.Response{data=ProviderResponse} "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/cmdb/providers [post]
// @Security BearerAuth
func (h *ProviderHandler) CreateProvider(c *gin.Context) {
	var provider model.Provider
	if err := c.ShouldBindJSON(&provider); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "参数格式错误: "+err.Error())
		return
	}

	if err := h.providerSvc.CreateProvider(&provider); err != nil {
		response.ReturnError(c, response.INTERNAL, "创建云账号失败: "+err.Error())
		return
	}
	
	response.ReturnData(c, provider)
}

// UpdateProvider 更新云账号
// @Summary 更新云账号信息
// @Description 根据云账号ID更新云账号的配置信息
// @Tags 云账号管理
// @Accept json
// @Produce json
// @Param id path int true "云账号ID"
// @Param provider body ProviderRequest true "更新的云账号信息"
// @Success 200 {object} response.Response{data=ProviderResponse} "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "云账号不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/cmdb/providers/{id} [put]
// @Security BearerAuth
func (h *ProviderHandler) UpdateProvider(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的云账号ID")
		return
	}

	var provider model.Provider
	if err := c.ShouldBindJSON(&provider); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "参数格式错误: "+err.Error())
		return
	}

	provider.ID = uint(id)
	if err := h.providerSvc.UpdateProvider(&provider); err != nil {
		response.ReturnError(c, response.INTERNAL, "更新云账号失败: "+err.Error())
		return
	}
	
	response.ReturnData(c, provider)
}

// DeleteProvider 删除云账号
// @Summary 删除云账号
// @Description 根据云账号ID删除云账号配置
// @Tags 云账号管理
// @Accept json
// @Produce json
// @Param id path int true "云账号ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "云账号不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/cmdb/providers/{id} [delete]
// @Security BearerAuth
func (h *ProviderHandler) DeleteProvider(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的云账号ID")
		return
	}

	if err := h.providerSvc.DeleteProvider(uint(id)); err != nil {
		response.ReturnError(c, response.INTERNAL, "删除云账号失败: "+err.Error())
		return
	}
	
	response.ReturnSuccess(c, "云账号删除成功")
}

// GetProvider 获取云账号详情
// @Summary 获取云账号详情
// @Description 根据云账号ID获取云账号的详细信息
// @Tags 云账号管理
// @Accept json
// @Produce json
// @Param id path int true "云账号ID"
// @Success 200 {object} response.Response{data=ProviderResponse} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "云账号不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/cmdb/providers/{id} [get]
// @Security BearerAuth
func (h *ProviderHandler) GetProvider(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的云账号ID")
		return
	}

	provider, err := h.providerSvc.GetProvider(uint(id))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "云账号不存在: "+err.Error())
		return
	}
	
	response.ReturnData(c, provider)
}

// ListProviders 获取云账号列表
// @Summary 获取云账号列表
// @Description 获取云账号列表，支持分页和过滤
// @Tags 云账号管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Param name query string false "云账号名称过滤"
// @Param type query string false "云服务商类型过滤"
// @Param status query string false "状态过滤"
// @Success 200 {object} response.Response{data=ProviderListResponse} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/cmdb/providers [get]
// @Security BearerAuth
func (h *ProviderHandler) ListProviders(c *gin.Context) {
	providers, err := h.providerSvc.ListProviders()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取云账号列表失败: "+err.Error())
		return
	}
	
	response.ReturnData(c, providers)
}

// SyncResources 同步云资源
// @Summary 同步云账号资源
// @Description 同步指定云账号的资源到CMDB系统
// @Tags 云账号管理
// @Accept json
// @Produce json
// @Param id path int true "云账号ID"
// @Param request body SyncResourcesRequest false "同步参数"
// @Success 200 {object} response.Response{data=SyncResourcesResponse} "同步成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "云账号不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/cmdb/providers/{id}/sync [post]
// @Security BearerAuth
func (h *ProviderHandler) SyncResources(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的云账号ID")
		return
	}

	var req SyncResourcesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果解析失败，使用默认值
		req.GroupID = nil
	}

	// 执行云账号主机同步
	result, err := h.providerSvc.SyncResources(uint(id), req.GroupID)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "同步资源失败: "+err.Error())
		return
	}

	response.ReturnData(c, result)
}

// ValidateCredentials 验证云账号凭证
// @Summary 验证云账号凭证
// @Description 验证云账号的访问密钥是否有效
// @Tags 云账号管理
// @Accept json
// @Produce json
// @Param id path int true "云账号ID"
// @Success 200 {object} response.Response{data=object} "验证成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "云账号不存在"
// @Failure 500 {object} response.Response "验证失败"
// @Router /api/v1/cmdb/providers/{id}/validate [post]
// @Security BearerAuth
func (h *ProviderHandler) ValidateCredentials(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的云账号ID")
		return
	}

	result, err := h.providerSvc.ValidateCredentials(uint(id))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "凭证验证失败: "+err.Error())
		return
	}

	response.ReturnData(c, result)
}

// GetRegions 获取云账号可用区域
// @Summary 获取云账号可用区域
// @Description 获取指定云账号支持的所有区域列表
// @Tags 云账号管理
// @Accept json
// @Produce json
// @Param id path int true "云账号ID"
// @Success 200 {object} response.Response{data=[]string} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "云账号不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/cmdb/providers/{id}/regions [get]
// @Security BearerAuth
func (h *ProviderHandler) GetRegions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的云账号ID")
		return
	}

	regions, err := h.providerSvc.GetRegions(uint(id))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取区域列表失败: "+err.Error())
		return
	}

	response.ReturnData(c, regions)
}

// GetResourceStatistics 获取云账号资源统计
// @Summary 获取云账号资源统计
// @Description 获取指定云账号的资源统计信息
// @Tags 云账号管理
// @Accept json
// @Produce json
// @Param id path int true "云账号ID"
// @Success 200 {object} response.Response{data=object} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "云账号不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/cmdb/providers/{id}/statistics [get]
// @Security BearerAuth
func (h *ProviderHandler) GetResourceStatistics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的云账号ID")
		return
	}

	stats, err := h.providerSvc.GetResourceStatistics(uint(id))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取资源统计失败: "+err.Error())
		return
	}

	response.ReturnData(c, stats)
}
