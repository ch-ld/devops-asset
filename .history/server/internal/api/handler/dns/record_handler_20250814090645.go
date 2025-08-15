package dns

import (
	"strconv"

	"api-server/internal/model/dns"
	"api-server/internal/response/response"
	svc "api-server/internal/service/dns"

	"github.com/gin-gonic/gin"
)

// CreateRecordRequest DNS记录创建请求体结构体
type CreateRecordRequest struct {
	DomainID uint   `json:"domain_id" binding:"required" example:"1" comment:"域名ID"`
	Name     string `json:"name" binding:"required" example:"www" comment:"记录名称"`
	Type     string `json:"type" binding:"required" example:"A" comment:"记录类型"`
	Value    string `json:"value" binding:"required" example:"192.168.1.1" comment:"记录值"`
	TTL      int    `json:"ttl" example:"300" comment:"TTL值"`
	Priority int    `json:"priority" example:"10" comment:"优先级(MX记录使用)"`
	Weight   int    `json:"weight" example:"5" comment:"权重(SRV记录使用)"`
	Port     int    `json:"port" example:"80" comment:"端口(SRV记录使用)"`
	Remark   string `json:"remark" example:"主要A记录" comment:"备注"`
}

// UpdateRecordRequest DNS记录更新请求体结构体
type UpdateRecordRequest struct {
	Name     string `json:"name" binding:"required" example:"www" comment:"记录名称"`
	Type     string `json:"type" binding:"required" example:"A" comment:"记录类型"`
	Value    string `json:"value" binding:"required" example:"192.168.1.1" comment:"记录值"`
	TTL      int    `json:"ttl" example:"300" comment:"TTL值"`
	Priority int    `json:"priority" example:"10" comment:"优先级(MX记录使用)"`
	Weight   int    `json:"weight" example:"5" comment:"权重(SRV记录使用)"`
	Port     int    `json:"port" example:"80" comment:"端口(SRV记录使用)"`
	Remark   string `json:"remark" example:"主要A记录" comment:"备注"`
}

// DomainInfo 域名信息结构体
type DomainInfo struct {
	ID   uint   `json:"id" example:"1" comment:"域名ID"`
	Name string `json:"name" example:"example.com" comment:"域名名称"`
}

// ProviderInfo 提供商信息结构体
type ProviderInfo struct {
	ID   uint   `json:"id" example:"1" comment:"提供商ID"`
	Name string `json:"name" example:"阿里云DNS" comment:"提供商名称"`
	Type string `json:"type" example:"aliyun" comment:"提供商类型"`
}

// RecordResponse DNS记录响应结构体
type RecordResponse struct {
	ID           uint          `json:"id" example:"1" comment:"记录ID"`
	DomainID     uint          `json:"domain_id" example:"1" comment:"域名ID"`
	Domain       *DomainInfo   `json:"domain,omitempty" comment:"域名信息"`
	ProviderID   uint          `json:"provider_id" example:"1" comment:"提供商ID"`
	Provider     *ProviderInfo `json:"provider,omitempty" comment:"提供商信息"`
	Name         string        `json:"name" example:"www" comment:"记录名称"`
	Type         string        `json:"type" example:"A" comment:"记录类型"`
	Value        string        `json:"value" example:"192.168.1.1" comment:"记录值"`
	TTL          int           `json:"ttl" example:"300" comment:"TTL值"`
	Priority     int           `json:"priority" example:"10" comment:"优先级"`
	Weight       int           `json:"weight" example:"5" comment:"权重"`
	Port         int           `json:"port" example:"80" comment:"端口"`
	Status       string        `json:"status" example:"active" comment:"状态"`
	SyncStatus   string        `json:"sync_status" example:"synced" comment:"同步状态"`
	LastSyncAt   string        `json:"last_sync_at,omitempty" comment:"最后同步时间"`
	Remark       string        `json:"remark" example:"主要A记录" comment:"备注"`
	CreatedAt    string        `json:"created_at" example:"2024-01-01T00:00:00Z" comment:"创建时间"`
	UpdatedAt    string        `json:"updated_at" example:"2024-01-01T00:00:00Z" comment:"更新时间"`
}

// RecordListResponse DNS记录列表响应结构体
type RecordListResponse struct {
	Total int               `json:"total" example:"100" comment:"总数"`
	Items []*RecordResponse `json:"items" comment:"记录列表"`
}

// RecordHandler DNS记录处理器
type RecordHandler struct {
	recordService *svc.RecordService
}

// NewRecordHandler 创建DNS记录处理器
func NewRecordHandler(recordService *svc.RecordService) *RecordHandler {
	return &RecordHandler{
		recordService: recordService,
	}
}

// CreateRecord 创建DNS记录
func (h *RecordHandler) CreateRecord(c *gin.Context) {
	var req CreateRecordRequest
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

	// 构建DNS记录对象
	record := &dns.Record{
		DomainID:  req.DomainID,
		Name:      req.Name,
		Type:      req.Type,
		Value:     req.Value,
		TTL:       req.TTL,
		Priority:  &req.Priority,
		Weight:    &req.Weight,
		Port:      &req.Port,
		Status:    "active",
		Remark:    req.Remark,
		TenantID:  tenantID.(uint),
		CreatedBy: userID.(uint),
		UpdatedBy: userID.(uint),
	}

	// 创建DNS记录
	if err := h.recordService.CreateRecord(record, userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "创建DNS记录失败: "+err.Error())
		return
	}

	// 返回创建的记录信息
	recordResp := h.convertToResponse(record)
	response.ReturnData(c, recordResp)
}

// UpdateRecord 更新DNS记录
func (h *RecordHandler) UpdateRecord(c *gin.Context) {
	// 获取记录ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的记录ID")
		return
	}

	var req UpdateRecordRequest
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

	// 构建DNS记录对象
	record := &dns.Record{
		Name:      req.Name,
		Type:      req.Type,
		Value:     req.Value,
		TTL:       req.TTL,
		Priority:  &req.Priority,
		Weight:    &req.Weight,
		Port:      &req.Port,
		Remark:    req.Remark,
		TenantID:  tenantID.(uint),
		UpdatedBy: userID.(uint),
	}
	record.ID = uint(id)

	// 更新DNS记录
	if err := h.recordService.UpdateRecord(record, userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "更新DNS记录失败: "+err.Error())
		return
	}

	// 获取更新后的记录信息
	updatedRecord, err := h.recordService.GetRecord(uint(id))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取更新后的记录信息失败: "+err.Error())
		return
	}

	// 返回更新后的记录信息
	recordResp := h.convertToResponse(updatedRecord)
	response.ReturnData(c, recordResp)
}

// DeleteRecord 删除DNS记录
func (h *RecordHandler) DeleteRecord(c *gin.Context) {
	// 获取记录ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的记录ID")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 删除DNS记录
	if err := h.recordService.DeleteRecord(uint(id), userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "删除DNS记录失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// GetRecord 获取DNS记录详情
func (h *RecordHandler) GetRecord(c *gin.Context) {
	// 获取记录ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的记录ID")
		return
	}

	// 获取DNS记录信息
	record, err := h.recordService.GetRecord(uint(id))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "DNS记录不存在: "+err.Error())
		return
	}

	// 返回记录信息
	recordResp := h.convertToResponse(record)
	response.ReturnData(c, recordResp)
}

// ListRecords 获取DNS记录列表
func (h *RecordHandler) ListRecords(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// 获取筛选参数
	domainIDStr := c.Query("domain_id")
	name := c.Query("name")
	recordType := c.Query("type")
	status := c.Query("status")

	// 获取租户信息
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 构建筛选条件
	filters := make(map[string]interface{})
	if domainIDStr != "" {
		if domainID, err := strconv.ParseUint(domainIDStr, 10, 32); err == nil {
			filters["domain_id"] = uint(domainID)
		}
	}
	if name != "" {
		filters["name"] = name
	}
	if recordType != "" {
		filters["type"] = recordType
	}
	if status != "" {
		filters["status"] = status
	}

	// 获取DNS记录列表
	records, total, err := h.recordService.ListRecords(tenantID.(uint), filters, size, (page-1)*size)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取DNS记录列表失败: "+err.Error())
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

	response.ReturnData(c, resp)
}

// SyncRecord 同步DNS记录到云厂商
// @Summary 同步DNS记录到云厂商
// @Description 将本地DNS记录同步到对应的云厂商
// @Tags DNS记录
// @Accept json
// @Produce json
// @Param id path int true "记录ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/dns/records/{id}/sync [post]
func (h *RecordHandler) SyncRecord(c *gin.Context) {
	// 获取记录ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的记录ID")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 同步DNS记录
	if err := h.recordService.SyncRecord(uint(id), userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "同步DNS记录失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// BatchSyncRecords 批量同步DNS记录
// @Summary 批量同步DNS记录
// @Description 批量同步多个DNS记录到云厂商
// @Tags DNS记录
// @Accept json
// @Produce json
// @Param ids body []uint true "记录ID列表"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/dns/records/batch/sync [post]
func (h *RecordHandler) BatchSyncRecords(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required" example:"[1,2,3]" comment:"记录ID列表"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请求参数错误: "+err.Error())
		return
	}

	if len(req.IDs) == 0 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "记录ID列表不能为空")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 批量同步DNS记录
	if err := h.recordService.BatchSyncRecords(req.IDs, userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "批量同步DNS记录失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// BatchDeleteRecords 批量删除DNS记录
// @Summary 批量删除DNS记录
// @Description 批量删除多个DNS记录
// @Tags DNS记录
// @Accept json
// @Produce json
// @Param ids body []uint true "记录ID列表"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/dns/records/batch/delete [post]
func (h *RecordHandler) BatchDeleteRecords(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required" example:"[1,2,3]" comment:"记录ID列表"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请求参数错误: "+err.Error())
		return
	}

	if len(req.IDs) == 0 {
		response.ReturnError(c, response.INVALID_ARGUMENT, "记录ID列表不能为空")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 批量删除DNS记录
	if err := h.recordService.BatchDeleteRecords(req.IDs, userID.(uint), c.ClientIP()); err != nil {
		response.ReturnError(c, response.INTERNAL, "批量删除DNS记录失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// convertToResponse 转换为响应格式
func (h *RecordHandler) convertToResponse(record *dns.Record) *RecordResponse {
	var priority, weight, port int
	if record.Priority != nil {
		priority = *record.Priority
	}
	if record.Weight != nil {
		weight = *record.Weight
	}
	if record.Port != nil {
		port = *record.Port
	}

	response := &RecordResponse{
		ID:         record.ID,
		DomainID:   record.DomainID,
		ProviderID: record.ProviderID,
		Name:       record.Name,
		Type:       record.Type,
		Value:      record.Value,
		TTL:        record.TTL,
		Priority:   priority,
		Weight:     weight,
		Port:       port,
		Status:     record.Status,
		SyncStatus: record.SyncStatus,
		Remark:     record.Remark,
		CreatedAt:  record.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  record.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// 添加域名信息
	if record.Domain != nil {
		response.Domain = &DomainInfo{
			ID:   record.Domain.ID,
			Name: record.Domain.Name,
		}
	}

	// 添加提供商信息
	if record.Provider != nil {
		response.Provider = &ProviderInfo{
			ID:   record.Provider.ID,
			Name: record.Provider.Name,
			Type: record.Provider.Type,
		}
	}

	// 格式化最后同步时间
	if record.LastSyncAt != nil && !record.LastSyncAt.Time.IsZero() {
		response.LastSyncAt = record.LastSyncAt.Time.Format("2006-01-02T15:04:05Z")
	}

	return response
}
