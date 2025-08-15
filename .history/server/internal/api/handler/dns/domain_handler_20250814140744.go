package dns

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"api-server/internal/model/dns"
	"api-server/internal/response/response"
	svc "api-server/internal/service/dns"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CreateDomainRequest 域名创建请求体结构体
// @Description 域名创建时的请求参数
type CreateDomainRequest struct {
	Name          string         `json:"name" binding:"required" example:"example.com" comment:"域名"`
	Status        string         `json:"status" example:"active" comment:"状态(active/inactive/expired)"`
	RegistrarType string         `json:"registrar_type" example:"godaddy" comment:"注册商类型"`
	RegistrarID   *uint          `json:"registrar_id" example:"1" comment:"注册商配置ID"`
	ExpiresAt     *time.Time     `json:"expires_at" example:"2024-12-31T23:59:59Z" comment:"过期时间"`
	AutoRenew     bool           `json:"auto_renew" example:"false" comment:"是否自动续费"`
	GroupID       *uint          `json:"group_id" example:"1" comment:"分组ID"`
	Configuration datatypes.JSON `json:"configuration" example:"{\"ns\":[\"ns1.example.com\"]}" comment:"配置信息"`
	Remark        string         `json:"remark" example:"主域名" comment:"备注"`
}

// UpdateDomainRequest 域名更新请求体结构体
// @Description 域名更新时的请求参数
type UpdateDomainRequest struct {
	Name          string         `json:"name" binding:"required" example:"example.com" comment:"域名"`
	Status        string         `json:"status" example:"active" comment:"状态(active/inactive/expired)"`
	RegistrarType string         `json:"registrar_type" example:"godaddy" comment:"注册商类型"`
	RegistrarID   *uint          `json:"registrar_id" example:"1" comment:"注册商配置ID"`
	ExpiresAt     *time.Time     `json:"expires_at" example:"2024-12-31T23:59:59Z" comment:"过期时间"`
	AutoRenew     bool           `json:"auto_renew" example:"false" comment:"是否自动续费"`
	GroupID       *uint          `json:"group_id" example:"1" comment:"分组ID"`
	Configuration datatypes.JSON `json:"configuration" example:"{\"ns\":[\"ns1.example.com\"]}" comment:"配置信息"`
	Remark        string         `json:"remark" example:"主域名" comment:"备注"`
}

// DomainResponse 域名响应结构体
// @Description 域名信息响应
type DomainResponse struct {
	ID            uint                 `json:"id" example:"1" comment:"域名ID"`
	Name          string               `json:"name" example:"example.com" comment:"域名"`
	Status        string               `json:"status" example:"active" comment:"状态"`
	RegistrarType string               `json:"registrar_type" example:"godaddy" comment:"注册商类型"`
	RegistrarID   *uint                `json:"registrar_id" example:"1" comment:"注册商配置ID"`
	ProviderID    *uint                `json:"provider_id" example:"1" comment:"DNS提供商ID"`
	Provider      *ProviderInfo        `json:"provider,omitempty" comment:"DNS提供商信息"`
	ExpiresAt     *time.Time           `json:"expires_at" example:"2024-12-31T23:59:59Z" comment:"过期时间"`
	AutoRenew     bool                 `json:"auto_renew" example:"false" comment:"是否自动续费"`
	GroupID       *uint                `json:"group_id" example:"1" comment:"分组ID"`
	Group         *DomainGroupResponse `json:"group,omitempty" comment:"分组信息"`
	Tags          []TagResponse        `json:"tags,omitempty" comment:"标签列表"`
	Configuration datatypes.JSON       `json:"configuration" example:"{\"ns\":[\"ns1.example.com\"]}" comment:"配置信息"`
	Remark        string               `json:"remark" example:"主域名" comment:"备注"`
	TenantID      uint                 `json:"tenant_id" example:"1" comment:"租户ID"`
	CreatedBy     uint                 `json:"created_by" example:"1" comment:"创建人ID"`
	UpdatedBy     uint                 `json:"updated_by" example:"1" comment:"更新人ID"`
	CreatedAt     time.Time            `json:"created_at" example:"2024-01-01T00:00:00Z" comment:"创建时间"`
	UpdatedAt     time.Time            `json:"updated_at" example:"2024-01-01T00:00:00Z" comment:"更新时间"`
}

// DomainListResponse 域名列表响应结构体
// @Description 域名列表响应
type DomainListResponse struct {
	Total int               `json:"total" example:"100" comment:"总数量"`
	Items []*DomainResponse `json:"items" comment:"域名列表"`
}

// DomainHandler 域名管理接口处理器
// @Description 域名管理相关接口处理器，负责域名的增删改查等HTTP请求处理
type DomainHandler struct {
	domainService *svc.DomainService
}

// NewDomainHandler 创建域名处理器实例
func NewDomainHandler(domainService *svc.DomainService) *DomainHandler {
	return &DomainHandler{
		domainService: domainService,
	}
}

// CreateDomain 创建域名
// @Summary 创建域名
// @Description 创建新的域名记录
// @Tags DNS域名管理
// @Accept json
// @Produce json
// @Param request body CreateDomainRequest true "域名创建请求"
// @Success 200 {object} response.Response{data=DomainResponse} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/domains [post]
func (h *DomainHandler) CreateDomain(c *gin.Context) {
	var req CreateDomainRequest
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

	// 构建域名对象
	domain := &dns.Domain{
		Name:          req.Name,
		Status:        req.Status,
		RegistrarType: req.RegistrarType,
		RegistrarID:   req.RegistrarID,
		ExpiresAt:     req.ExpiresAt,
		AutoRenew:     req.AutoRenew,
		GroupID:       req.GroupID,
		Configuration: req.Configuration,
		Remark:        req.Remark,
		TenantID:      tenantID.(uint),
	}

	// 创建域名
	if err := h.domainService.CreateDomain(domain, userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "创建域名失败: "+err.Error())
		return
	}

	// 返回创建的域名信息
	domainResp := h.convertToResponse(domain)
	response.ReturnData(c, domainResp)
}

// UpdateDomain 更新域名
// @Summary 更新域名
// @Description 更新指定ID的域名信息
// @Tags DNS域名管理
// @Accept json
// @Produce json
// @Param id path int true "域名ID"
// @Param request body UpdateDomainRequest true "域名更新请求"
// @Success 200 {object} response.Response{data=DomainResponse} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "域名不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/domains/{id} [put]
func (h *DomainHandler) UpdateDomain(c *gin.Context) {
	// 获取域名ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的域名ID")
		return
	}

	var req UpdateDomainRequest
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

	// 构建域名对象
	domain := &dns.Domain{
		Model:         gorm.Model{ID: uint(id)},
		Name:          req.Name,
		Status:        req.Status,
		RegistrarType: req.RegistrarType,
		RegistrarID:   req.RegistrarID,
		ExpiresAt:     req.ExpiresAt,
		AutoRenew:     req.AutoRenew,
		GroupID:       req.GroupID,
		Configuration: req.Configuration,
		Remark:        req.Remark,
		TenantID:      tenantID.(uint),
	}

	// 更新域名
	if err := h.domainService.UpdateDomain(domain, userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "更新域名失败: "+err.Error())
		return
	}

	// 获取更新后的域名信息
	updatedDomain, err := h.domainService.GetDomain(uint(id))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取更新后的域名信息失败: "+err.Error())
		return
	}

	// 返回更新后的域名信息
	domainResp := h.convertToResponse(updatedDomain)
	response.ReturnData(c, domainResp)
}

// DeleteDomain 删除域名
// @Summary 删除域名
// @Description 删除指定ID的域名
// @Tags DNS域名管理
// @Produce json
// @Param id path int true "域名ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "域名不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/domains/{id} [delete]
func (h *DomainHandler) DeleteDomain(c *gin.Context) {
	// 获取域名ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的域名ID")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 删除域名
	if err := h.domainService.DeleteDomain(uint(id), userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "删除域名失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// GetDomain 获取域名详情
// @Summary 获取域名详情
// @Description 根据ID获取域名详细信息
// @Tags DNS域名管理
// @Produce json
// @Param id path int true "域名ID"
// @Success 200 {object} response.Response{data=DomainResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "域名不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/domains/{id} [get]
func (h *DomainHandler) GetDomain(c *gin.Context) {
	// 获取域名ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的域名ID")
		return
	}

	// 获取域名信息
	domain, err := h.domainService.GetDomain(uint(id))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "域名不存在: "+err.Error())
		return
	}

	// 返回域名信息
	domainResp := h.convertToResponse(domain)
	response.ReturnData(c, domainResp)
}

// ListDomains 获取域名列表
// @Summary 获取域名列表
// @Description 获取域名列表，支持分页和筛选
// @Tags DNS域名管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param keyword query string false "关键词搜索"
// @Param status query string false "状态筛选"
// @Param group_id query int false "分组ID筛选"
// @Success 200 {object} response.Response{data=DomainListResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/domains [get]
func (h *DomainHandler) ListDomains(c *gin.Context) {
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

	// 构建筛选条件
	filters := make(map[string]interface{})
	if keyword := c.Query("keyword"); keyword != "" {
		filters["keyword"] = keyword
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if groupIDStr := c.Query("group_id"); groupIDStr != "" {
		if groupID, err := strconv.ParseUint(groupIDStr, 10, 32); err == nil {
			filters["group_id"] = uint(groupID)
		}
	}

	// 获取域名列表
	domains, total, err := h.domainService.ListDomains(tenantID.(uint), filters, size, (page-1)*size)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取域名列表失败: "+err.Error())
		return
	}

	// 转换响应格式
	items := make([]*DomainResponse, len(domains))
	for i, domain := range domains {
		items[i] = h.convertToResponse(domain)
	}

	resp := &DomainListResponse{
		Total: int(total),
		Items: items,
	}

	response.ReturnData(c, resp)
}

// convertToResponse 转换为响应格式
func (h *DomainHandler) convertToResponse(domain *dns.Domain) *DomainResponse {
	resp := &DomainResponse{
		ID:            domain.ID,
		Name:          domain.Name,
		Status:        domain.Status,
		RegistrarType: domain.RegistrarType,
		RegistrarID:   domain.RegistrarID,
		ProviderID:    domain.ProviderID,
		ExpiresAt:     domain.ExpiresAt,
		AutoRenew:     domain.AutoRenew,
		GroupID:       domain.GroupID,
		Configuration: domain.Configuration,
		Remark:        domain.Remark,
		TenantID:      domain.TenantID,
		CreatedBy:     domain.CreatedBy,
		UpdatedBy:     domain.UpdatedBy,
		CreatedAt:     domain.CreatedAt,
		UpdatedAt:     domain.UpdatedAt,
	}

	// 转换分组信息
	if domain.Group != nil {
		resp.Group = &DomainGroupResponse{
			ID:          domain.Group.ID,
			Name:        domain.Group.Name,
			Description: domain.Group.Description,
		}
	}

	// 转换标签信息
	if len(domain.Tags) > 0 {
		resp.Tags = make([]TagResponse, len(domain.Tags))
		for i, tag := range domain.Tags {
			resp.Tags[i] = TagResponse{
				ID:    tag.ID,
				Name:  tag.Name,
				Color: tag.Color,
			}
		}
	}

	return resp
}

// ImportDomainsFromCSV 从CSV文件导入域名
// @Summary 从CSV文件导入域名
// @Description 支持批量导入域名，CSV格式需包含：域名、状态、注册商类型、过期时间、自动续费、分组、备注等列
// @Tags 域名管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV文件"
// @Success 200 {object} response.Response{data=svc.DomainImportResult} "导入结果"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dns/domains/import/csv [post]
func (h *DomainHandler) ImportDomainsFromCSV(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息获取失败")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "文件上传失败: "+err.Error())
		return
	}

	// 验证文件类型
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".csv") {
		response.ReturnError(c, response.INVALID_ARGUMENT, "文件格式不正确，请上传CSV文件")
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "打开文件失败: "+err.Error())
		return
	}
	defer src.Close()

	// 调用Service层导入域名
	result, err := h.domainService.ImportDomainsFromCSV(
		src,
		tenantID.(uint),
		userID.(uint),
		c.ClientIP(),
	)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "导入域名失败: "+err.Error())
		return
	}

	response.ReturnData(c, result)
}

// ImportDomainsFromExcel 从Excel文件导入域名
// @Summary 从Excel文件导入域名
// @Description 支持批量导入域名，Excel格式需包含：域名、状态、注册商类型、过期时间、自动续费、分组、备注等列
// @Tags 域名管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel文件"
// @Success 200 {object} response.Response{data=svc.DomainImportResult} "导入结果"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dns/domains/import/excel [post]
func (h *DomainHandler) ImportDomainsFromExcel(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息获取失败")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "文件上传失败: "+err.Error())
		return
	}

	// 验证文件类型
	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	if fileExt != ".xlsx" && fileExt != ".xls" {
		response.ReturnError(c, response.INVALID_ARGUMENT, "文件格式不正确，请上传Excel文件")
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "打开文件失败: "+err.Error())
		return
	}
	defer src.Close()

	// 调用Service层导入域名
	result, err := h.domainService.ImportDomainsFromExcel(
		src,
		tenantID.(uint),
		userID.(uint),
		c.ClientIP(),
	)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "导入域名失败: "+err.Error())
		return
	}

	response.ReturnData(c, result)
}

// ExportDomainsToCSV 导出域名到CSV文件
// @Summary 导出域名到CSV文件
// @Description 导出域名列表到CSV格式文件
// @Tags 域名管理
// @Produce application/octet-stream
// @Param group_id query int false "分组ID"
// @Param status query string false "状态"
// @Param registrar_type query string false "注册商类型"
// @Success 200 {file} file "CSV文件"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dns/domains/export/csv [get]
func (h *DomainHandler) ExportDomainsToCSV(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息获取失败")
		return
	}

	// 解析查询参数
	filters := make(map[string]interface{})
	filters["tenant_id"] = tenantID.(uint)

	if groupID := c.Query("group_id"); groupID != "" {
		if id, err := strconv.ParseUint(groupID, 10, 32); err == nil {
			filters["group_id"] = uint(id)
		}
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	if registrarType := c.Query("registrar_type"); registrarType != "" {
		filters["registrar_type"] = registrarType
	}

	// 设置响应头
	filename := fmt.Sprintf("domains_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")

	// 导出数据
	if err := h.domainService.ExportDomainsToCSV(c.Writer, tenantID.(uint), filters); err != nil {
		response.ReturnError(c, response.INTERNAL, "导出域名失败: "+err.Error())
		return
	}
}

// ExportDomainsToExcel 导出域名到Excel文件
// @Summary 导出域名到Excel文件
// @Description 导出域名列表到Excel格式文件
// @Tags 域名管理
// @Produce application/octet-stream
// @Param group_id query int false "分组ID"
// @Param status query string false "状态"
// @Param registrar_type query string false "注册商类型"
// @Success 200 {file} file "Excel文件"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dns/domains/export/excel [get]
func (h *DomainHandler) ExportDomainsToExcel(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息获取失败")
		return
	}

	// 解析查询参数
	filters := make(map[string]interface{})
	filters["tenant_id"] = tenantID.(uint)

	if groupID := c.Query("group_id"); groupID != "" {
		if id, err := strconv.ParseUint(groupID, 10, 32); err == nil {
			filters["group_id"] = uint(id)
		}
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	if registrarType := c.Query("registrar_type"); registrarType != "" {
		filters["registrar_type"] = registrarType
	}

	// 设置响应头
	filename := fmt.Sprintf("domains_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")

	// 导出数据
	if err := h.domainService.ExportDomainsToExcel(c.Writer, tenantID.(uint), filters); err != nil {
		response.ReturnError(c, response.INTERNAL, "导出域名失败: "+err.Error())
		return
	}
}

// GetDomainStatistics 获取域名统计信息
// @Summary 获取域名统计信息
// @Description 获取域名统计信息，包括总数、状态分布、即将过期数量等
// @Tags DNS域名管理
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}} "统计信息"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/stats/domains [get]
func (h *DomainHandler) GetDomainStatistics(c *gin.Context) {
	// 获取租户信息
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 获取统计信息
	stats, err := h.domainService.GetDomainStatistics(tenantID.(uint))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取域名统计信息失败: "+err.Error())
		return
	}

	response.ReturnData(c, stats)
}

// GetDomainWhois 获取域名WHOIS信息
// @Summary 获取域名WHOIS信息
// @Description 获取指定域名的WHOIS信息
// @Tags DNS域名管理
// @Produce json
// @Param id path int true "域名ID"
// @Success 200 {object} response.Response{data=map[string]interface{}} "WHOIS信息"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "域名不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/domains/{id}/whois [get]
func (h *DomainHandler) GetDomainWhois(c *gin.Context) {
	// 获取域名ID
	domainID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的域名ID")
		return
	}

	// 获取域名信息
	domain, err := h.domainService.GetDomain(uint(domainID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.ReturnError(c, response.NOT_FOUND, "域名不存在")
		} else {
			response.ReturnError(c, response.INTERNAL, "获取域名信息失败: "+err.Error())
		}
		return
	}

	// 获取WHOIS信息
	whoisInfo, err := h.domainService.GetDomainWhois(domain.Name)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取WHOIS信息失败: "+err.Error())
		return
	}

	response.ReturnData(c, whoisInfo)
}

// ExportWhoisReport 导出WHOIS报告
// @Summary 导出域名WHOIS报告
// @Description 导出指定域名的WHOIS报告
// @Tags DNS域名管理
// @Produce application/pdf
// @Param id path int true "域名ID"
// @Success 200 {file} application/pdf "WHOIS报告PDF文件"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "域名不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/domains/{id}/whois/report [get]
func (h *DomainHandler) ExportWhoisReport(c *gin.Context) {
	// 获取域名ID
	domainID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的域名ID")
		return
	}

	// 获取域名信息
	domain, err := h.domainService.GetDomain(uint(domainID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.ReturnError(c, response.NOT_FOUND, "域名不存在")
		} else {
			response.ReturnError(c, response.INTERNAL, "获取域名信息失败: "+err.Error())
		}
		return
	}

	// 生成WHOIS报告
	reportData, err := h.domainService.GenerateWhoisReport(domain.Name)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "生成WHOIS报告失败: "+err.Error())
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s_whois_report.pdf\"", domain.Name))
	c.Data(200, "application/pdf", reportData)
}

// BatchDeleteDomains 批量删除域名
// @Summary 批量删除域名
// @Description 批量删除多个域名
// @Tags DNS域名管理
// @Accept json
// @Produce json
// @Param request body struct{Action string `json:"action"`;IDs []uint `json:"ids"`} true "批量删除请求"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/domains/batch [post]
func (h *DomainHandler) BatchDeleteDomains(c *gin.Context) {
	var req struct {
		Action string `json:"action" binding:"required"`
		IDs    []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请求参数错误: "+err.Error())
		return
	}

	if req.Action != "delete" {
		response.ReturnError(c, response.INVALID_ARGUMENT, "不支持的操作类型")
		return
	}

	if len(req.IDs) == 0 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请选择要删除的域名")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 获取客户端IP
	clientIP := c.ClientIP()

	// 批量删除域名
	err := h.domainService.BatchDeleteDomains(req.IDs, userID.(uint), clientIP)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "批量删除域名失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// DownloadCSVTemplate 下载CSV模板
// @Summary 下载CSV导入模板
// @Description 下载域名导入的CSV模板文件
// @Tags DNS域名管理
// @Produce application/octet-stream
// @Success 200 {file} application/octet-stream "CSV模板文件"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/domains/template/csv [get]
func (h *DomainHandler) DownloadCSVTemplate(c *gin.Context) {
	// CSV模板内容
	csvContent := `域名,状态,注册商类型,过期时间,自动续费,分组ID,备注
example.com,active,godaddy,2025-12-31,true,1,示例域名
test.com,active,dnspod,2025-06-30,false,2,测试域名`

	// 设置响应头
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=\"domain_import_template.csv\"")
	c.Data(200, "text/csv", []byte(csvContent))
}

// DownloadExcelTemplate 下载Excel模板
// @Summary 下载Excel导入模板
// @Description 下载域名导入的Excel模板文件
// @Tags DNS域名管理
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Success 200 {file} application/vnd.openxmlformats-officedocument.spreadsheetml.sheet "Excel模板文件"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/domains/template/excel [get]
func (h *DomainHandler) DownloadExcelTemplate(c *gin.Context) {
	// 生成Excel模板
	templateData, err := h.domainService.GenerateExcelTemplate()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "生成Excel模板失败: "+err.Error())
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\"domain_import_template.xlsx\"")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", templateData)
}
