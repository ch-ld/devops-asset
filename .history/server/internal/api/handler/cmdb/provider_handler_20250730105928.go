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

// 云账号管理相关接口处理器
// 负责云账号的增删改查、密钥加密等HTTP请求处理
type ProviderHandler struct {
	providerSvc *svc.ProviderService
}

// NewProviderHandler 创建云账号处理器实例
func NewProviderHandler(providerSvc *svc.ProviderService) *ProviderHandler {
	return &ProviderHandler{providerSvc: providerSvc}
}

// CreateProvider 创建云账号接口
func (h *ProviderHandler) CreateProvider(c *gin.Context) {
	var provider model.Provider
	if err := c.ShouldBindJSON(&provider); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, err.Error())
		return
	}

	if err := h.providerSvc.CreateProvider(&provider); err != nil {
		response.ReturnError(c, response.INTERNAL, "Failed to create provider: "+err.Error())
		return
	}
	response.ReturnData(c, provider)
}

// UpdateProvider 更新云账号接口
func (h *ProviderHandler) UpdateProvider(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "Invalid provider ID")
		return
	}

	var provider model.Provider
	if err := c.ShouldBindJSON(&provider); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, err.Error())
		return
	}
	provider.ID = uint(id)

	if err := h.providerSvc.UpdateProvider(&provider); err != nil {
		response.ReturnError(c, response.INTERNAL, "Failed to update provider: "+err.Error())
		return
	}
	response.ReturnData(c, provider)
}

// DeleteProvider 删除云账号接口
func (h *ProviderHandler) DeleteProvider(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "Invalid provider ID")
		return
	}

	if err := h.providerSvc.DeleteProvider(uint(id)); err != nil {
		response.ReturnError(c, response.INTERNAL, "Failed to delete provider: "+err.Error())
		return
	}
	response.ReturnSuccess(c)
}

// ListProviders 查询云账号列表接口
func (h *ProviderHandler) ListProviders(c *gin.Context) {
	providers, err := h.providerSvc.ListProviders()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "Failed to list providers: "+err.Error())
		return
	}
	response.ReturnData(c, providers)
}

// GetProvider 查询单个云账号接口
func (h *ProviderHandler) GetProvider(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "Invalid provider ID")
		return
	}

	provider, err := h.providerSvc.GetProviderByID(uint(id))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "Provider not found: "+err.Error())
		return
	}
	response.ReturnData(c, provider)
}

// SyncResourcesRequest 同步资源请求
type SyncResourcesRequest struct {
	GroupID *uint `json:"group_id" form:"group_id"` // 可选的主机组ID
}

// SyncResources handles triggering a resource sync for a provider.
func (h *ProviderHandler) SyncResources(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "Invalid provider ID")
		return
	}

	// 解析请求参数 - 同时支持查询参数和请求体参数
	var req SyncResourcesRequest

	// 首先尝试从查询参数获取
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "Invalid query parameters: "+err.Error())
		return
	}

	// 如果查询参数中没有group_id，尝试从请求体获取
	if req.GroupID == nil {
		var bodyReq SyncResourcesRequest
		if err := c.ShouldBindJSON(&bodyReq); err == nil && bodyReq.GroupID != nil {
			req.GroupID = bodyReq.GroupID
		}
	}

	// 添加调试日志
	// 执行云账号主机同步

	result, err := h.providerSvc.SyncResources(uint(id), req.GroupID)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "Failed to sync resources: "+err.Error())
		return
	}

	response.ReturnData(c, result)
}

// ValidateCredentials 验证云账号凭证接口
func (h *ProviderHandler) ValidateCredentials(c *gin.Context) {
	var req struct {
		Type      string `json:"type" binding:"required"`
		AccessKey string `json:"access_key" binding:"required"`
		SecretKey string `json:"secret_key" binding:"required"`
		Region    string `json:"region"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("❌ 验证凭证请求参数错误", zap.Error(err))
		response.ReturnError(c, response.INVALID_ARGUMENT, err.Error())
		return
	}

	zap.L().Info("🚀 收到验证凭证请求",
		zap.String("type", req.Type),
		zap.String("accessKey", req.AccessKey[:min(len(req.AccessKey), 8)]+"***"),
		zap.String("region", req.Region),
		zap.String("clientIP", c.ClientIP()),
	)

	isValid, err := h.providerSvc.ValidateCredentials(req.Type, req.AccessKey, req.SecretKey, req.Region)
	if err != nil {
		zap.L().Error("❌ 验证凭证服务调用失败", zap.Error(err))
		response.ReturnError(c, response.INTERNAL, "Failed to validate credentials: "+err.Error())
		return
	}

	resultMessage := func() string {
		if isValid {
			return "凭证验证成功"
		}
		return "凭证验证失败，请检查AccessKey和SecretKey是否正确"
	}()

	zap.L().Info("📋 验证凭证结果",
		zap.Bool("valid", isValid),
		zap.String("message", resultMessage),
	)

	response.ReturnData(c, map[string]interface{}{
		"valid":   isValid,
		"message": resultMessage,
	})
}

// GetProviderTypes 获取支持的云厂商类型接口
func (h *ProviderHandler) GetProviderTypes(c *gin.Context) {
	types := []map[string]interface{}{
		{
			"value":       "aliyun",
			"label":       "阿里云",
			"icon":        "aliyun",
			"needsRegion": true,
		},
		{
			"value":       "tencent",
			"label":       "腾讯云",
			"icon":        "tencent",
			"needsRegion": true,
		},
		{
			"value":       "aws",
			"label":       "AWS",
			"icon":        "aws",
			"needsRegion": true,
		},
	}

	response.ReturnData(c, types)
}

// GetProviderRegions 获取云厂商支持的地域列表
func (h *ProviderHandler) GetProviderRegions(c *gin.Context) {
	providerType := c.Param("type")
	if providerType == "" {
		response.ReturnError(c, response.INVALID_ARGUMENT, "云厂商类型不能为空")
		return
	}

	var regions []map[string]interface{}

	switch providerType {
	case "aliyun":
		regions = []map[string]interface{}{
			{"value": "cn-hangzhou", "label": "华东1（杭州）"},
			{"value": "cn-shanghai", "label": "华东2（上海）"},
			{"value": "cn-qingdao", "label": "华北1（青岛）"},
			{"value": "cn-beijing", "label": "华北2（北京）"},
			{"value": "cn-zhangjiakou", "label": "华北3（张家口）"},
			{"value": "cn-huhehaote", "label": "华北5（呼和浩特）"},
			{"value": "cn-wulanchabu", "label": "华北6（乌兰察布）"},
			{"value": "cn-shenzhen", "label": "华南1（深圳）"},
			{"value": "cn-heyuan", "label": "华南2（河源）"},
			{"value": "cn-guangzhou", "label": "华南3（广州）"},
			{"value": "cn-chengdu", "label": "西南1（成都）"},
			{"value": "cn-hongkong", "label": "中国香港"},
		}
	case "tencent":
		regions = []map[string]interface{}{
			{"value": "ap-beijing", "label": "北京"},
			{"value": "ap-chengdu", "label": "成都"},
			{"value": "ap-chongqing", "label": "重庆"},
			{"value": "ap-guangzhou", "label": "广州"},
			{"value": "ap-hongkong", "label": "中国香港"},
			{"value": "ap-mumbai", "label": "孟买"},
			{"value": "ap-seoul", "label": "首尔"},
			{"value": "ap-shanghai", "label": "上海"},
			{"value": "ap-singapore", "label": "新加坡"},
			{"value": "ap-tokyo", "label": "东京"},
			{"value": "eu-frankfurt", "label": "法兰克福"},
			{"value": "na-ashburn", "label": "弗吉尼亚"},
			{"value": "na-siliconvalley", "label": "硅谷"},
		}
	case "aws":
		regions = []map[string]interface{}{
			{"value": "us-east-1", "label": "美国东部（弗吉尼亚北部）"},
			{"value": "us-east-2", "label": "美国东部（俄亥俄）"},
			{"value": "us-west-1", "label": "美国西部（加利福尼亚北部）"},
			{"value": "us-west-2", "label": "美国西部（俄勒冈）"},
			{"value": "ap-south-1", "label": "亚太地区（孟买）"},
			{"value": "ap-northeast-1", "label": "亚太地区（东京）"},
			{"value": "ap-northeast-2", "label": "亚太地区（首尔）"},
			{"value": "ap-southeast-1", "label": "亚太地区（新加坡）"},
			{"value": "ap-southeast-2", "label": "亚太地区（悉尼）"},
			{"value": "eu-central-1", "label": "欧洲（法兰克福）"},
			{"value": "eu-west-1", "label": "欧洲（爱尔兰）"},
			{"value": "eu-west-2", "label": "欧洲（伦敦）"},
			{"value": "cn-north-1", "label": "中国（北京）"},
			{"value": "cn-northwest-1", "label": "中国（宁夏）"},
		}
	default:
		response.ReturnError(c, response.INVALID_ARGUMENT, "不支持的云厂商类型")
		return
	}

	response.ReturnData(c, regions)
}
