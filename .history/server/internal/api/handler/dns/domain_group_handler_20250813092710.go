package dns

import (
	"strconv"

	"api-server/internal/model/dns"
	"api-server/internal/response/response"
	svc "api-server/internal/service/dns"

	"github.com/gin-gonic/gin"
)

// DomainGroupHandler 域名分组管理接口处理器
type DomainGroupHandler struct {
	domainGroupService *svc.DomainGroupService
}

// NewDomainGroupHandler 创建域名分组处理器实例
func NewDomainGroupHandler(domainGroupService *svc.DomainGroupService) *DomainGroupHandler {
	return &DomainGroupHandler{
		domainGroupService: domainGroupService,
	}
}

// CreateDomainGroupRequest 创建域名分组请求
type CreateDomainGroupRequest struct {
	Name        string `json:"name" binding:"required,max=100" example:"生产环境" comment:"分组名称"`
	ParentID    *uint  `json:"parent_id" example:"1" comment:"父分组ID"`
	Description string `json:"description" max:"500" example:"生产环境域名分组" comment:"描述"`
	Sort        int    `json:"sort" example:"1" comment:"排序"`
	Color       string `json:"color" example:"#3b82f6" comment:"分组颜色"`
	Status      string `json:"status" example:"active" comment:"状态"`
}

// UpdateDomainGroupRequest 更新域名分组请求
type UpdateDomainGroupRequest struct {
	Name        string `json:"name" max:"100" example:"生产环境" comment:"分组名称"`
	ParentID    *uint  `json:"parent_id" example:"1" comment:"父分组ID"`
	Description string `json:"description" max:"500" example:"生产环境域名分组" comment:"描述"`
	Sort        int    `json:"sort" example:"1" comment:"排序"`
	Color       string `json:"color" example:"#3b82f6" comment:"分组颜色"`
	Status      string `json:"status" example:"active" comment:"状态"`
}

// DomainGroupListResponse 域名分组列表响应
type DomainGroupListResponse struct {
	Total int                    `json:"total" example:"10" comment:"总数"`
	Items []*DomainGroupResponse `json:"items" comment:"分组列表"`
}

// CreateDomainGroup 创建域名分组
// @Summary 创建域名分组
// @Description 创建新的域名分组
// @Tags DNS域名分组
// @Accept json
// @Produce json
// @Param request body CreateDomainGroupRequest true "域名分组创建请求"
// @Success 200 {object} response.Response{data=DomainGroupResponse} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/dns/domain-groups [post]
func (h *DomainGroupHandler) CreateDomainGroup(c *gin.Context) {
	var req CreateDomainGroupRequest
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

	// 构建域名分组对象
	group := &dns.DomainGroup{
		Name:        req.Name,
		ParentID:    req.ParentID,
		Description: req.Description,
		Sort:        req.Sort,
		Color:       req.Color,
		Status:      req.Status,
		TenantID:    tenantID.(uint),
		CreatedBy:   userID.(uint),
		UpdatedBy:   userID.(uint),
	}

	// 设置默认值
	if group.Color == "" {
		group.Color = "#3b82f6"
	}
	if group.Status == "" {
		group.Status = "active"
	}

	// 获取客户端IP
	clientIP := c.ClientIP()

	// 创建域名分组
	if err := h.domainGroupService.CreateDomainGroup(group, userID.(uint), clientIP); err != nil {
		response.ReturnError(c, response.INTERNAL, "创建域名分组失败: "+err.Error())
		return
	}

	// 返回创建的域名分组信息
	groupResp := h.convertToResponse(group)
	response.ReturnData(c, groupResp)
}

// UpdateDomainGroup 更新域名分组
// @Summary 更新域名分组
// @Description 更新指定ID的域名分组信息
// @Tags DNS域名分组
// @Accept json
// @Produce json
// @Param id path int true "域名分组ID"
// @Param request body UpdateDomainGroupRequest true "域名分组更新请求"
// @Success 200 {object} response.Response{data=DomainGroupResponse} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "域名分组不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/dns/domain-groups/{id} [put]
func (h *DomainGroupHandler) UpdateDomainGroup(c *gin.Context) {
	// 获取分组ID
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的分组ID")
		return
	}

	var req UpdateDomainGroupRequest
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

	// 获取现有分组
	group, err := h.domainGroupService.GetDomainGroup(uint(groupID))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "域名分组不存在")
		return
	}

	// 更新分组信息
	if req.Name != "" {
		group.Name = req.Name
	}
	if req.ParentID != nil {
		group.ParentID = req.ParentID
	}
	if req.Description != "" {
		group.Description = req.Description
	}
	if req.Sort != 0 {
		group.Sort = req.Sort
	}
	if req.Color != "" {
		group.Color = req.Color
	}
	if req.Status != "" {
		group.Status = req.Status
	}
	group.UpdatedBy = userID.(uint)

	// 获取客户端IP
	clientIP := c.ClientIP()

	// 更新域名分组
	if err := h.domainGroupService.UpdateDomainGroup(group, userID.(uint), clientIP); err != nil {
		response.ReturnError(c, response.INTERNAL, "更新域名分组失败: "+err.Error())
		return
	}

	// 返回更新的域名分组信息
	groupResp := h.convertToResponse(group)
	response.ReturnData(c, groupResp)
}

// DeleteDomainGroup 删除域名分组
// @Summary 删除域名分组
// @Description 删除指定ID的域名分组
// @Tags DNS域名分组
// @Produce json
// @Param id path int true "域名分组ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "域名分组不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/dns/domain-groups/{id} [delete]
func (h *DomainGroupHandler) DeleteDomainGroup(c *gin.Context) {
	// 获取分组ID
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的分组ID")
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

	// 删除域名分组
	if err := h.domainGroupService.DeleteDomainGroup(uint(groupID), userID.(uint), clientIP); err != nil {
		response.ReturnError(c, response.INTERNAL, "删除域名分组失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// GetDomainGroup 获取域名分组详情
// @Summary 获取域名分组详情
// @Description 获取指定ID的域名分组详细信息
// @Tags DNS域名分组
// @Produce json
// @Param id path int true "域名分组ID"
// @Success 200 {object} response.Response{data=DomainGroupResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "域名分组不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/dns/domain-groups/{id} [get]
func (h *DomainGroupHandler) GetDomainGroup(c *gin.Context) {
	// 获取分组ID
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的分组ID")
		return
	}

	// 获取租户信息
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 获取域名分组
	group, err := h.domainGroupService.GetDomainGroupByID(uint(groupID), tenantID.(uint))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "域名分组不存在")
		return
	}

	// 返回域名分组信息
	groupResp := h.convertToResponse(group)
	response.ReturnData(c, groupResp)
}

// ListDomainGroups 获取域名分组列表
// @Summary 获取域名分组列表
// @Description 获取域名分组列表，支持分页和筛选
// @Tags DNS域名分组
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "关键词搜索"
// @Param parent_id query int false "父分组ID筛选"
// @Param status query string false "状态筛选"
// @Success 200 {object} response.Response{data=DomainGroupListResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/dns/domain-groups [get]
func (h *DomainGroupHandler) ListDomainGroups(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
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
	if parentIDStr := c.Query("parent_id"); parentIDStr != "" {
		if parentID, err := strconv.ParseUint(parentIDStr, 10, 32); err == nil {
			filters["parent_id"] = uint(parentID)
		}
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	// 获取域名分组列表
	groups, total, err := h.domainGroupService.ListDomainGroups(tenantID.(uint), filters, pageSize, (page-1)*pageSize)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取域名分组列表失败: "+err.Error())
		return
	}

	// 转换响应格式
	items := make([]*DomainGroupResponse, len(groups))
	for i, group := range groups {
		items[i] = h.convertToResponse(group)
	}

	resp := &DomainGroupListResponse{
		Total: int(total),
		Items: items,
	}

	response.ReturnData(c, resp)
}

// GetDomainGroupTree 获取域名分组树
// @Summary 获取域名分组树
// @Description 获取完整的域名分组树结构
// @Tags DNS域名分组
// @Produce json
// @Success 200 {object} response.Response{data=[]*DomainGroupResponse} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/dns/domain-groups/tree [get]
func (h *DomainGroupHandler) GetDomainGroupTree(c *gin.Context) {
	// 获取租户信息
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 获取域名分组树
	groups, err := h.domainGroupService.GetDomainGroupTree(tenantID.(uint))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取域名分组树失败: "+err.Error())
		return
	}

	// 转换响应格式
	resp := h.convertToTreeResponse(groups)
	response.ReturnData(c, resp)
}

// MoveDomainGroup 移动域名分组
// @Summary 移动域名分组
// @Description 移动域名分组到新的父分组
// @Tags DNS域名分组
// @Accept json
// @Produce json
// @Param id path int true "域名分组ID"
// @Param request body struct{ParentID *uint `json:"parent_id"`;Sort int `json:"sort"`} true "移动请求"
// @Success 200 {object} response.Response "移动成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "域名分组不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/dns/domain-groups/{id}/move [put]
func (h *DomainGroupHandler) MoveDomainGroup(c *gin.Context) {
	// 获取分组ID
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的分组ID")
		return
	}

	var req struct {
		ParentID *uint `json:"parent_id"`
		Sort     int   `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请求参数错误: "+err.Error())
		return
	}

	// 获取租户信息
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 移动域名分组
	if err := h.domainGroupService.MoveDomainGroup(uint(groupID), req.ParentID, req.Sort, tenantID.(uint)); err != nil {
		response.ReturnError(c, response.INTERNAL, "移动域名分组失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c, "移动成功")
}

// convertToResponse 转换为响应格式
func (h *DomainGroupHandler) convertToResponse(group *dns.DomainGroup) *DomainGroupResponse {
	resp := &DomainGroupResponse{
		ID:               group.ID,
		Name:             group.Name,
		ParentID:         group.ParentID,
		Description:      group.Description,
		Sort:             group.Sort,
		Color:            group.Color,
		Status:           group.Status,
		DomainCount:      group.DomainCount,
		TotalDomainCount: group.TotalDomainCount,
		TenantID:         group.TenantID,
		CreatedBy:        group.CreatedBy,
		UpdatedBy:        group.UpdatedBy,
		CreatedAt:        group.CreatedAt,
		UpdatedAt:        group.UpdatedAt,
	}

	// 转换父分组信息
	if group.Parent != nil {
		resp.Parent = &DomainGroupResponse{
			ID:   group.Parent.ID,
			Name: group.Parent.Name,
		}
	}

	// 转换子分组信息
	if len(group.Children) > 0 {
		resp.Children = make([]*DomainGroupResponse, len(group.Children))
		for i, child := range group.Children {
			resp.Children[i] = h.convertToResponse(child)
		}
	}

	return resp
}

// convertToTreeResponse 转换为树形响应格式
func (h *DomainGroupHandler) convertToTreeResponse(groups []*dns.DomainGroup) []*DomainGroupResponse {
	result := make([]*DomainGroupResponse, len(groups))
	for i, group := range groups {
		result[i] = h.convertToResponse(group)
	}
	return result
}
