package cmdb

import (
	"strconv"

	model "api-server/internal/model/cmdb"
	"api-server/internal/response/response"
	svc "api-server/internal/service/cmdb"
	"api-server/internal/service/cmdb/adapter"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ProviderRequest 云账号请求结构体
// @Description 云账号创建或更新时的请求参数
type ProviderRequest struct {
	Name        string `json:"name" binding:"required" example:"阿里云账号" comment:"云账号名称"`
	Type        string `json:"type" binding:"required" example:"alicloud" comment:"云服务商类型：alicloud,aws,tencentcloud"`
	AccessKey   string `json:"access_key" binding:"required" example:"LTAI4..." comment:"访问密钥ID"`
	SecretKey   string `json:"secret_key" binding:"required" example:"abc123..." comment:"访问密钥Secret"`
	Region      string `json:"region" example:"cn-hangzhou" comment:"默认区域"`
	Status      string `json:"status" example:"active" comment:"状态：active,inactive"`
	Description string `json:"description" example:"生产环境阿里云账号" comment:"描述信息"`
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
	Total   int    `json:"total" example:"100" comment:"总同步数量"`
	Success int    `json:"success" example:"95" comment:"成功同步数量"`
	Failed  int    `json:"failed" example:"5" comment:"失败同步数量"`
	Skipped int    `json:"skipped" example:"10" comment:"跳过数量"`
	Message string `json:"message" example:"同步完成" comment:"同步结果消息"`
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

	response.ReturnSuccess(c)
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

	provider, err := h.providerSvc.GetProviderByID(uint(id))
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
	_, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的云账号ID")
		return
	}

	// TODO: 实现凭证验证功能
	response.ReturnData(c, map[string]interface{}{
		"valid":   true,
		"message": "凭证验证功能待实现",
	})
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
	_, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的云账号ID")
		return
	}

	// TODO: 实现获取区域列表功能
	response.ReturnData(c, []string{"cn-hangzhou", "cn-shanghai", "cn-beijing"})
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
	_, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的云账号ID")
		return
	}

	// TODO: 实现资源统计功能
	response.ReturnData(c, map[string]interface{}{
		"hosts":          0,
		"load_balancers": 0,
		"databases":      0,
		"total":          0,
	})
}

// GetProviderTypes 获取支持的云厂商类型
// @Summary 获取支持的云厂商类型
// @Description 获取系统支持的所有云厂商类型列表
// @Tags 云账号管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]string} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/cmdb/providers/types [get]
// @Security BearerAuth
func (h *ProviderHandler) GetProviderTypes(c *gin.Context) {
	types := []string{
		"alicloud",     // 阿里云
		"tencentcloud", // 腾讯云
		"aws",          // 亚马逊云
		"huaweicloud",  // 华为云
		"baiduyun",     // 百度云
	}

	response.ReturnData(c, types)
}

// GetProviderRegions 获取云厂商支持的地域列表
// @Summary 获取云厂商支持的地域列表
// @Description 根据云厂商类型获取支持的地域列表，优先返回静态列表确保用户体验
// @Tags 云账号管理
// @Accept json
// @Produce json
// @Param type path string true "云厂商类型" Enums(alicloud,tencentcloud,aws,huaweicloud,baiduyun)
// @Success 200 {object} response.Response{data=[]string} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/cmdb/providers/types/{type}/regions [get]
// @Security BearerAuth
func (h *ProviderHandler) GetProviderRegions(c *gin.Context) {
	providerType := c.Param("type")

	// 直接返回静态区域列表，确保用户体验
	// 用户可以在验证凭证后通过 GetProviderRegionsWithCredentials 获取动态列表
	h.getFallbackRegions(c, providerType)
}

// GetProviderRegionsWithCredentials 使用凭证获取云厂商支持的地域列表
// @Summary 使用凭证获取云厂商支持的地域列表
// @Description 使用提供的凭证动态获取云厂商真实的地域列表
// @Tags 云账号管理
// @Accept json
// @Produce json
// @Param request body object{type=string,access_key=string,secret_key=string} true "请求参数"
// @Success 200 {object} response.Response{data=[]string} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/cmdb/providers/regions-with-credentials [post]
// @Security BearerAuth
func (h *ProviderHandler) GetProviderRegionsWithCredentials(c *gin.Context) {
	var req struct {
		Type      string `json:"type" binding:"required"`
		AccessKey string `json:"access_key" binding:"required"`
		SecretKey string `json:"secret_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "参数错误: "+err.Error())
		return
	}

	// 对于支持动态获取区域的云厂商，使用CloudAdapter
	switch req.Type {
	case "alicloud", "tencentcloud", "aws":
		tempProvider := &model.Provider{
			Type:      req.Type,
			AccessKey: req.AccessKey,
			SecretKey: req.SecretKey,
			Region:    "",
		}

		adapter, err := adapter.GetCloudAdapter(tempProvider)
		if err != nil {
			response.ReturnError(c, response.INTERNAL, "创建云适配器失败: "+err.Error())
			return
		}

		regions, err := adapter.GetRegions()
		if err != nil {
			zap.L().Warn("Failed to get regions from cloud provider API with credentials",
				zap.String("provider", req.Type),
				zap.Error(err))
			response.ReturnError(c, response.INTERNAL, "获取区域列表失败: "+err.Error())
			return
		}

		// 转换为字符串数组格式（保持前端兼容性）
		var regionList []string
		for _, region := range regions {
			regionList = append(regionList, region.ID)
		}

		response.ReturnData(c, regionList)
		return

	default:
		// 对于不支持动态获取的云厂商，返回静态列表
		h.getFallbackRegions(c, req.Type)
	}
}

// getFallbackRegions 获取静态区域列表（作为备选方案）
func (h *ProviderHandler) getFallbackRegions(c *gin.Context, providerType string) {
	var regions []string

	switch providerType {
	case "alicloud":
		regions = []string{
			"cn-hangzhou",    // 华东1（杭州）
			"cn-shanghai",    // 华东2（上海）
			"cn-qingdao",     // 华北1（青岛）
			"cn-beijing",     // 华北2（北京）
			"cn-zhangjiakou", // 华北3（张家口）
			"cn-huhehaote",   // 华北5（呼和浩特）
			"cn-wulanchabu",  // 华北6（乌兰察布）
			"cn-shenzhen",    // 华南1（深圳）
			"cn-heyuan",      // 华南2（河源）
			"cn-guangzhou",   // 华南3（广州）
			"cn-chengdu",     // 西南1（成都）
			"cn-hongkong",    // 中国香港
			"ap-northeast-1", // 亚太东北1（东京）
			"ap-southeast-1", // 亚太东南1（新加坡）
			"ap-southeast-2", // 亚太东南2（悉尼）
			"ap-southeast-3", // 亚太东南3（吉隆坡）
			"ap-southeast-5", // 亚太东南5（雅加达）
			"ap-south-1",     // 亚太南部1（孟买）
			"us-east-1",      // 美国东部1（弗吉尼亚）
			"us-west-1",      // 美国西部1（硅谷）
			"eu-west-1",      // 欧洲西部1（伦敦）
			"eu-central-1",   // 欧洲中部1（法兰克福）
			"me-east-1",      // 中东东部1（迪拜）
		}
	case "tencentcloud":
		regions = []string{
			"ap-beijing",       // 华北地区（北京）
			"ap-beijing-fsi",   // 华北地区（北京金融）
			"ap-tianjin",       // 华北地区（天津）
			"ap-shijiazhuang",  // 华北地区（石家庄）
			"ap-shanghai",      // 华东地区（上海）
			"ap-shanghai-fsi",  // 华东地区（上海金融）
			"ap-nanjing",       // 华东地区（南京）
			"ap-guangzhou",     // 华南地区（广州）
			"ap-shenzhen-fsi",  // 华南地区（深圳金融）
			"ap-chengdu",       // 西南地区（成都）
			"ap-chongqing",     // 西南地区（重庆）
			"ap-hongkong",      // 港澳台地区（中国香港）
			"ap-taipei",        // 港澳台地区（中国台北）
			"ap-singapore",     // 亚太东南（新加坡）
			"ap-bangkok",       // 亚太东南（曼谷）
			"ap-jakarta",       // 亚太东南（雅加达）
			"ap-seoul",         // 亚太东北（首尔）
			"ap-tokyo",         // 亚太东北（东京）
			"ap-mumbai",        // 亚太南部（孟买）
			"na-siliconvalley", // 美国西部（硅谷）
			"na-ashburn",       // 美国东部（弗吉尼亚）
			"na-toronto",       // 北美地区（多伦多）
			"sa-saopaulo",      // 南美地区（圣保罗）
			"eu-frankfurt",     // 欧洲地区（法兰克福）
			"eu-moscow",        // 欧洲地区（莫斯科）
		}
	case "aws":
		regions = []string{
			"us-east-1",      // 美国东部（弗吉尼亚北部）
			"us-west-2",      // 美国西部（俄勒冈）
			"eu-west-1",      // 欧洲（爱尔兰）
			"ap-southeast-1", // 亚太（新加坡）
			"ap-northeast-1", // 亚太（东京）
		}
	case "huaweicloud":
		regions = []string{
			"cn-north-1", // 华北-北京一
			"cn-north-4", // 华北-北京四
			"cn-east-2",  // 华东-上海二
			"cn-east-3",  // 华东-上海一
			"cn-south-1", // 华南-广州
		}
	case "baiduyun":
		regions = []string{
			"bj",  // 北京
			"gz",  // 广州
			"su",  // 苏州
			"hkg", // 香港
		}
	default:
		response.ReturnError(c, response.INVALID_ARGUMENT, "不支持的云厂商类型: "+providerType)
		return
	}

	response.ReturnData(c, regions)
}
