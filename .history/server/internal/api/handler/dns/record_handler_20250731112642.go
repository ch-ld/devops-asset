package dns

import (
	"net/http"
	"strconv"
	"time"

	"api-server/internal/model/dns"
	"api-server/internal/response/response"
	svc "api-server/internal/service/dns"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// CreateRecordRequest DNS记录创建请求体结构体
// @Description DNS记录创建时的请求参数
type CreateRecordRequest struct {
	DomainID      uint           `json:"domain_id" binding:"required" example:"1" comment:"域名ID"`
	ProviderID    uint           `json:"provider_id" binding:"required" example:"1" comment:"DNS提供商ID"`
	Name          string         `json:"name" binding:"required" example:"www" comment:"记录名称"`
	Type          string         `json:"type" binding:"required" example:"A" comment:"记录类型"`
	Value         string         `json:"value" binding:"required" example:"1.2.3.4" comment:"记录值"`
	TTL           int            `json:"ttl" example:"600" comment:"TTL值(秒)"`
	Priority      *int           `json:"priority" example:"10" comment:"优先级(MX/SRV记录)"`
	Weight        *int           `json:"weight" example:"5" comment:"权重(SRV记录)"`
	Port          *int           `json:"port" example:"80" comment:"端口(SRV记录)"`
	Status        string         `json:"status" example:"active" comment:"状态(active/inactive)"`
	Configuration datatypes.JSON `json:"configuration" example:"{\"proxy\":true}" comment:"配置信息"`
	Remark        string         `json:"remark" example:"主站点" comment:"备注"`
}

// UpdateRecordRequest DNS记录更新请求体结构体
// @Description DNS记录更新时的请求参数
type UpdateRecordRequest struct {
	DomainID      uint           `json:"domain_id" binding:"required" example:"1" comment:"域名ID"`
	ProviderID    uint           `json:"provider_id" binding:"required" example:"1" comment:"DNS提供商ID"`
	Name          string         `json:"name" binding:"required" example:"www" comment:"记录名称"`
	Type          string         `json:"type" binding:"required" example:"A" comment:"记录类型"`
	Value         string         `json:"value" binding:"required" example:"1.2.3.4" comment:"记录值"`
	TTL           int            `json:"ttl" example:"600" comment:"TTL值(秒)"`
	Priority      *int           `json:"priority" example:"10" comment:"优先级(MX/SRV记录)"`
	Weight        *int           `json:"weight" example:"5" comment:"权重(SRV记录)"`
	Port          *int           `json:"port" example:"80" comment:"端口(SRV记录)"`
	Status        string         `json:"status" example:"active" comment:"状态(active/inactive)"`
	Configuration datatypes.JSON `json:"configuration" example:"{\"proxy\":true}" comment:"配置信息"`
	Remark        string         `json:"remark" example:"主站点" comment:"备注"`
}

// RecordResponse DNS记录响应结构体
// @Description DNS记录信息响应
type RecordResponse struct {
	ID            uint           `json:"id" example:"1" comment:"记录ID"`
	DomainID      uint           `json:"domain_id" example:"1" comment:"域名ID"`
	Domain        *DomainResponse `json:"domain,omitempty" comment:"域名信息"`
	ProviderID    uint           `json:"provider_id" example:"1" comment:"DNS提供商ID"`
	Provider      *ProviderResponse `json:"provider,omitempty" comment:"提供商信息"`
	Name          string         `json:"name" example:"www" comment:"记录名称"`
	Type          string         `json:"type" example:"A" comment:"记录类型"`
	Value         string         `json:"value" example:"1.2.3.4" comment:"记录值"`
	TTL           int            `json:"ttl" example:"600" comment:"TTL值(秒)"`
	Priority      *int           `json:"priority" example:"10" comment:"优先级"`
	Weight        *int           `json:"weight" example:"5" comment:"权重"`
	Port          *int           `json:"port" example:"80" comment:"端口"`
	Status        string         `json:"status" example:"active" comment:"状态"`
	CloudRecordID string         `json:"cloud_record_id" example:"rec_123456" comment:"云端记录ID"`
	SyncStatus    string         `json:"sync_status" example:"synced" comment:"同步状态"`
	LastSyncAt    *time.Time     `json:"last_sync_at" example:"2024-01-01T00:00:00Z" comment:"最后同步时间"`
	Configuration datatypes.JSON `json:"configuration" example:"{\"proxy\":true}" comment:"配置信息"`
	Remark        string         `json:"remark" example:"主站点" comment:"备注"`
	TenantID      uint           `json:"tenant_id" example:"1" comment:"租户ID"`
	CreatedBy     uint           `json:"created_by" example:"1" comment:"创建人ID"`
	UpdatedBy     uint           `json:"updated_by" example:"1" comment:"更新人ID"`
	CreatedAt     time.Time      `json:"created_at" example:"2024-01-01T00:00:00Z" comment:"创建时间"`
	UpdatedAt     time.Time      `json:"updated_at" example:"2024-01-01T00:00:00Z" comment:"更新时间"`
}

// RecordListResponse DNS记录列表响应结构体
// @Description DNS记录列表响应
type RecordListResponse struct {
	Total int              `json:"total" example:"100" comment:"总数量"`
	Items []*RecordResponse `json:"items" comment:"记录列表"`
}

// RecordHandler DNS记录管理接口处理器
// @Description DNS记录管理相关接口处理器，负责DNS记录的增删改查等HTTP请求处理
type RecordHandler struct {
	recordService *svc.RecordService
}

// NewRecordHandler 创建DNS记录处理器实例
func NewRecordHandler(recordService *svc.RecordService) *RecordHandler {
	return &RecordHandler{
		recordService: recordService,
	}
}

// CreateRecord 创建DNS记录
// @Summary 创建DNS记录
// @Description 创建新的DNS解析记录
// @Tags DNS记录管理
// @Accept json
// @Produce json
// @Param request body CreateRecordRequest true "DNS记录创建请求"
// @Success 200 {object} response.Response{data=RecordResponse} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/records [post]
func (h *RecordHandler) CreateRecord(c *gin.Context) {
	var req CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	// 构建DNS记录对象
	record := &dns.Record{
		DomainID:      req.DomainID,
		ProviderID:    req.ProviderID,
		Name:          req.Name,
		Type:          req.Type,
		Value:         req.Value,
		TTL:           req.TTL,
		Priority:      req.Priority,
		Weight:        req.Weight,
		Port:          req.Port,
		Status:        req.Status,
		Configuration: req.Configuration,
		Remark:        req.Remark,
	}

	// 设置默认值
	if record.TTL == 0 {
		record.TTL = 600
	}
	if record.Status == "" {
		record.Status = "active"
	}

	// 创建DNS记录
	if err := h.recordService.CreateRecord(record, userID.(uint), c.ClientIP()); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建DNS记录失败: "+err.Error())
		return
	}

	// 返回创建的记录信息
	recordResp := h.convertToResponse(record)
	response.Success(c, recordResp)
}

// UpdateRecord 更新DNS记录
// @Summary 更新DNS记录
// @Description 更新指定ID的DNS记录信息
// @Tags DNS记录管理
// @Accept json
// @Produce json
// @Param id path int true "记录ID"
// @Param request body UpdateRecordRequest true "DNS记录更新请求"
// @Success 200 {object} response.Response{data=RecordResponse} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "记录不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/records/{id} [put]
func (h *RecordHandler) UpdateRecord(c *gin.Context) {
	// 获取记录ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的记录ID")
		return
	}

	var req UpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	// 构建DNS记录对象
	record := &dns.Record{
		Model:         dns.Model{ID: uint(id)},
		DomainID:      req.DomainID,
		ProviderID:    req.ProviderID,
		Name:          req.Name,
		Type:          req.Type,
		Value:         req.Value,
		TTL:           req.TTL,
		Priority:      req.Priority,
		Weight:        req.Weight,
		Port:          req.Port,
		Status:        req.Status,
		Configuration: req.Configuration,
		Remark:        req.Remark,
	}

	// 更新DNS记录
	if err := h.recordService.UpdateRecord(record, userID.(uint), c.ClientIP()); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新DNS记录失败: "+err.Error())
		return
	}

	// 获取更新后的记录信息
	updatedRecord, err := h.recordService.GetRecord(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取更新后的记录信息失败: "+err.Error())
		return
	}

	// 返回更新后的记录信息
	recordResp := h.convertToResponse(updatedRecord)
	response.Success(c, recordResp)
}

// DeleteRecord 删除DNS记录
// @Summary 删除DNS记录
// @Description 删除指定ID的DNS记录
// @Tags DNS记录管理
// @Produce json
// @Param id path int true "记录ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "记录不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/records/{id} [delete]
func (h *RecordHandler) DeleteRecord(c *gin.Context) {
	// 获取记录ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的记录ID")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	// 删除DNS记录
	if err := h.recordService.DeleteRecord(uint(id), userID.(uint), c.ClientIP()); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除DNS记录失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// GetRecord 获取DNS记录详情
// @Summary 获取DNS记录详情
// @Description 根据ID获取DNS记录详细信息
// @Tags DNS记录管理
// @Produce json
// @Param id path int true "记录ID"
// @Success 200 {object} response.Response{data=RecordResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "记录不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/records/{id} [get]
func (h *RecordHandler) GetRecord(c *gin.Context) {
	// 获取记录ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的记录ID")
		return
	}

	// 获取DNS记录信息
	record, err := h.recordService.GetRecord(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "DNS记录不存在: "+err.Error())
		return
	}

	// 返回记录信息
	recordResp := h.convertToResponse(record)
	response.Success(c, recordResp)
}

// ListRecords 获取DNS记录列表
// @Summary 获取DNS记录列表
// @Description 获取DNS记录列表，支持分页和筛选
// @Tags DNS记录管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param keyword query string false "关键词搜索"
// @Param domain_id query int false "域名ID筛选"
// @Param type query string false "记录类型筛选"
// @Param status query string false "状态筛选"
// @Success 200 {object} response.Response{data=RecordListResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/dns/records [get]
func (h *RecordHandler) ListRecords(c *gin.Context) {
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
		response.Error(c, http.StatusUnauthorized, "租户信息缺失")
		return
	}

	// 构建筛选条件
	filters := make(map[string]interface{})
	if keyword := c.Query("keyword"); keyword != "" {
		filters["keyword"] = keyword
	}
	if domainIDStr := c.Query("domain_id"); domainIDStr != "" {
		if domainID, err := strconv.ParseUint(domainIDStr, 10, 32); err == nil {
			filters["domain_id"] = uint(domainID)
		}
	}
	if recordType := c.Query("type"); recordType != "" {
		filters["type"] = recordType
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	// 获取DNS记录列表
	records, total, err := h.recordService.ListRecords(tenantID.(uint), filters, size, (page-1)*size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取DNS记录列表失败: "+err.Error())
		return
	}

	// 转换响应格式
	items := make([]*RecordResponse, len(records))
	for i, record := range records {
		items[i] = h.convertToResponse(record)
	}

	resp := &RecordListResponse{
		Total: int(total),
		Items: items,
	}

	response.Success(c, resp)
}

// convertToResponse 转换为响应格式
func (h *RecordHandler) convertToResponse(record *dns.Record) *RecordResponse {
	resp := &RecordResponse{
		ID:            record.ID,
		DomainID:      record.DomainID,
		ProviderID:    record.ProviderID,
		Name:          record.Name,
		Type:          record.Type,
		Value:         record.Value,
		TTL:           record.TTL,
		Priority:      record.Priority,
		Weight:        record.Weight,
		Port:          record.Port,
		Status:        record.Status,
		CloudRecordID: record.CloudRecordID,
		SyncStatus:    record.SyncStatus,
		LastSyncAt:    record.LastSyncAt,
		Configuration: record.Configuration,
		Remark:        record.Remark,
		TenantID:      record.TenantID,
		CreatedBy:     record.CreatedBy,
		UpdatedBy:     record.UpdatedBy,
		CreatedAt:     record.CreatedAt,
		UpdatedAt:     record.UpdatedAt,
	}

	// 转换域名信息
	if record.Domain != nil {
		resp.Domain = &DomainResponse{
			ID:   record.Domain.ID,
			Name: record.Domain.Name,
		}
	}

	// 转换提供商信息
	if record.Provider != nil {
		resp.Provider = &ProviderResponse{
			ID:   record.Provider.ID,
			Name: record.Provider.Name,
			Type: record.Provider.Type,
		}
	}

	return resp
}
