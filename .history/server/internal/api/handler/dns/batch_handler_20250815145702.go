package dns

import (
	"strconv"

	"api-server/internal/model/dns"
	"api-server/internal/response/response"
	svc "api-server/internal/service/dns"

	"github.com/gin-gonic/gin"
)

// BatchRecordRequest 批量DNS记录请求结构体
type BatchRecordRequest struct {
	DomainID   uint                   `json:"domain_id" binding:"required" example:"1" comment:"域名ID"`
	Records    []*CreateRecordRequest `json:"records" binding:"required" comment:"DNS记录列表"`
	Sync       bool                   `json:"sync" example:"false" comment:"创建后是否立即同步云厂商"`
	ProviderID uint                   `json:"provider_id" example:"1" comment:"DNS提供商ID（用于同步）"`
}

// BatchRecordResponse 批量DNS记录响应结构体
type BatchRecordResponse struct {
	Total   int               `json:"total" example:"10" comment:"总数"`
	Success int               `json:"success" example:"8" comment:"成功数"`
	Failed  int               `json:"failed" example:"2" comment:"失败数"`
	Records []*RecordResponse `json:"records" comment:"成功创建的记录"`
	Errors  []BatchError      `json:"errors" comment:"错误信息"`
}

// BatchError 批量操作错误信息
type BatchError struct {
	Index  int    `json:"index" example:"1" comment:"记录索引"`
	Record string `json:"record" example:"www.example.com" comment:"记录名称"`
	Error  string `json:"error" example:"记录已存在" comment:"错误信息"`
}

// BatchDeleteRequest 批量删除请求结构体
type BatchDeleteRequest struct {
	RecordIDs []uint `json:"record_ids" binding:"required" comment:"记录ID列表"`
}

// BatchUpdateRequest 批量更新请求结构体
type BatchUpdateRequest struct {
	Records []*BatchUpdateRecord `json:"records" binding:"required" comment:"更新记录列表"`
}

// BatchUpdateRecord 批量更新记录结构体
type BatchUpdateRecord struct {
	ID       uint   `json:"id" binding:"required" example:"1" comment:"记录ID"`
	Name     string `json:"name" binding:"required" example:"www" comment:"记录名称"`
	Type     string `json:"type" binding:"required" example:"A" comment:"记录类型"`
	Value    string `json:"value" binding:"required" example:"192.168.1.1" comment:"记录值"`
	TTL      int    `json:"ttl" example:"300" comment:"TTL值"`
	Priority int    `json:"priority" example:"10" comment:"优先级"`
	Weight   int    `json:"weight" example:"5" comment:"权重"`
	Port     int    `json:"port" example:"80" comment:"端口"`
	Remark   string `json:"remark" example:"批量更新记录" comment:"备注"`
}

// BatchSyncRequest DNS记录同步请求结构体
type BatchSyncRequest struct {
	DomainID   uint `json:"domain_id" binding:"required" example:"1" comment:"域名ID"`
	ProviderID uint `json:"provider_id" binding:"required" example:"1" comment:"DNS提供商ID"`
	DryRun     bool `json:"dry_run" example:"false" comment:"是否为试运行"`
}

// BatchSyncResponse DNS记录同步响应结构体
type BatchSyncResponse struct {
	DomainName     string            `json:"domain_name" example:"example.com" comment:"域名"`
	Provider       string            `json:"provider" example:"aliyun" comment:"DNS提供商"`
	TotalLocal     int               `json:"total_local" example:"10" comment:"本地记录总数"`
	TotalRemote    int               `json:"total_remote" example:"12" comment:"远程记录总数"`
	ToAdd          int               `json:"to_add" example:"3" comment:"需要添加的记录数"`
	ToUpdate       int               `json:"to_update" example:"2" comment:"需要更新的记录数"`
	ToDelete       int               `json:"to_delete" example:"1" comment:"需要删除的记录数"`
	AddedRecords   []*RecordResponse `json:"added_records" comment:"已添加的记录"`
	UpdatedRecords []*RecordResponse `json:"updated_records" comment:"已更新的记录"`
	DeletedRecords []string          `json:"deleted_records" comment:"已删除的记录ID"`
	Errors         []string          `json:"errors" comment:"错误信息"`
	DryRun         bool              `json:"dry_run" comment:"是否为试运行"`
}

// BatchHandler DNS批量操作处理器
type BatchHandler struct {
	recordService *svc.RecordService
	domainService *svc.DomainService
}

// NewBatchHandler 创建DNS批量操作处理器
func NewBatchHandler(recordService *svc.RecordService, domainService *svc.DomainService) *BatchHandler {
	return &BatchHandler{
		recordService: recordService,
		domainService: domainService,
	}
}

// BatchCreateRecords 批量创建DNS记录
func (h *BatchHandler) BatchCreateRecords(c *gin.Context) {
	var req BatchRecordRequest
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

	var successRecords []*RecordResponse
	var errors []BatchError
	successCount := 0

	for i, recordReq := range req.Records {
		// 构建DNS记录对象
		record := &dns.Record{
			DomainID:  req.DomainID,
			Name:      recordReq.Name,
			Type:      recordReq.Type,
			Value:     recordReq.Value,
			TTL:       recordReq.TTL,
			Priority:  &recordReq.Priority,
			Weight:    &recordReq.Weight,
			Port:      &recordReq.Port,
			Status:    "active",
			Remark:    recordReq.Remark,
			TenantID:  tenantID.(uint),
			CreatedBy: userID.(uint),
			UpdatedBy: userID.(uint),
		}

		// 创建DNS记录
		if err := h.recordService.CreateRecord(record, userID.(uint), c.ClientIP()); err != nil {
			errors = append(errors, BatchError{
				Index:  i,
				Record: recordReq.Name,
				Error:  err.Error(),
			})
			continue
		}

		// 转换为响应格式
		recordResp := &RecordResponse{
			ID:        record.ID,
			DomainID:  record.DomainID,
			Name:      record.Name,
			Type:      record.Type,
			Value:     record.Value,
			TTL:       record.TTL,
			Priority:  *record.Priority,
			Weight:    *record.Weight,
			Port:      *record.Port,
			Status:    record.Status,
			Remark:    record.Remark,
			CreatedAt: record.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: record.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}

		successRecords = append(successRecords, recordResp)
		successCount++
	}
	// 如果要求创建后同步
	if req.Sync && req.ProviderID != 0 {
		params := &svc.SyncDomainRecordsParams{
			DomainID:   req.DomainID,
			ProviderID: req.ProviderID,
			DryRun:     false,
			ActorID:    userID.(uint),
			ClientIP:   c.ClientIP(),
		}
		if _, err := h.recordService.SyncDomainRecords(params); err != nil {
			errors = append(errors, BatchError{
				Index:  -1,
				Record: "sync",
				Error:  "创建完成，但同步失败: " + err.Error(),
			})
		}
	}

	// 返回批量操作结果
	resp := &BatchRecordResponse{
		Total:   len(req.Records),
		Success: successCount,
		Failed:  len(errors),
		Records: successRecords,
		Errors:  errors,
	}

	response.ReturnData(c, resp)
}

// BatchDeleteRecords 批量删除DNS记录
func (h *BatchHandler) BatchDeleteRecords(c *gin.Context) {
	var req BatchDeleteRequest
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

	var errors []BatchError
	successCount := 0

	for i, recordID := range req.RecordIDs {
		// 删除DNS记录
		if err := h.recordService.DeleteRecord(recordID, userID.(uint), c.ClientIP()); err != nil {
			errors = append(errors, BatchError{
				Index:  i,
				Record: strconv.FormatUint(uint64(recordID), 10),
				Error:  err.Error(),
			})
			continue
		}
		successCount++
	}

	// 返回批量操作结果
	resp := &BatchRecordResponse{
		Total:   len(req.RecordIDs),
		Success: successCount,
		Failed:  len(errors),
		Records: nil,
		Errors:  errors,
	}

	response.ReturnData(c, resp)
}

// BatchUpdateRecords 批量更新DNS记录
func (h *BatchHandler) BatchUpdateRecords(c *gin.Context) {
	var req BatchUpdateRequest
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

	var successRecords []*RecordResponse
	var errors []BatchError
	successCount := 0

	for i, recordReq := range req.Records {
		// 构建DNS记录对象
		record := &dns.Record{
			Name:      recordReq.Name,
			Type:      recordReq.Type,
			Value:     recordReq.Value,
			TTL:       recordReq.TTL,
			Priority:  &recordReq.Priority,
			Weight:    &recordReq.Weight,
			Port:      &recordReq.Port,
			Remark:    recordReq.Remark,
			TenantID:  tenantID.(uint),
			UpdatedBy: userID.(uint),
		}
		record.ID = recordReq.ID

		// 更新DNS记录
		if err := h.recordService.UpdateRecord(record, userID.(uint), c.ClientIP()); err != nil {
			errors = append(errors, BatchError{
				Index:  i,
				Record: recordReq.Name,
				Error:  err.Error(),
			})
			continue
		}

		// 获取更新后的记录信息
		updatedRecord, err := h.recordService.GetRecord(recordReq.ID)
		if err != nil {
			errors = append(errors, BatchError{
				Index:  i,
				Record: recordReq.Name,
				Error:  "获取更新后的记录失败: " + err.Error(),
			})
			continue
		}

		// 转换为响应格式
		recordResp := &RecordResponse{
			ID:        updatedRecord.ID,
			DomainID:  updatedRecord.DomainID,
			Name:      updatedRecord.Name,
			Type:      updatedRecord.Type,
			Value:     updatedRecord.Value,
			TTL:       updatedRecord.TTL,
			Priority:  *updatedRecord.Priority,
			Weight:    *updatedRecord.Weight,
			Port:      *updatedRecord.Port,
			Status:    updatedRecord.Status,
			Remark:    updatedRecord.Remark,
			CreatedAt: updatedRecord.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: updatedRecord.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}

		successRecords = append(successRecords, recordResp)
		successCount++
	}

	// 返回批量操作结果
	resp := &BatchRecordResponse{
		Total:   len(req.Records),
		Success: successCount,
		Failed:  len(errors),
		Records: successRecords,
		Errors:  errors,
	}

	response.ReturnData(c, resp)
}

// SyncDomainRecords 同步域名DNS记录
func (h *BatchHandler) SyncDomainRecords(c *gin.Context) {
	var req BatchSyncRequest
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

	// 调用Service层的同步方法
	params := &svc.SyncDomainRecordsParams{
		DomainID:   req.DomainID,
		ProviderID: req.ProviderID,
		DryRun:     req.DryRun,
		ActorID:    userID.(uint),
		ClientIP:   c.ClientIP(),
	}

	syncResult, err := h.recordService.SyncDomainRecords(params)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "DNS记录同步失败: "+err.Error())
		return
	}

	// 转换为API响应格式
	addedRecords := make([]*RecordResponse, len(syncResult.AddedRecords))
	for i, record := range syncResult.AddedRecords {
		addedRecords[i] = &RecordResponse{
			ID:        record.ID,
			DomainID:  record.DomainID,
			Name:      record.Name,
			Type:      record.Type,
			Value:     record.Value,
			TTL:       record.TTL,
			Priority:  *record.Priority,
			Weight:    *record.Weight,
			Port:      *record.Port,
			Status:    record.Status,
			Remark:    record.Remark,
			CreatedAt: record.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: record.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	updatedRecords := make([]*RecordResponse, len(syncResult.UpdatedRecords))
	for i, record := range syncResult.UpdatedRecords {
		updatedRecords[i] = &RecordResponse{
			ID:        record.ID,
			DomainID:  record.DomainID,
			Name:      record.Name,
			Type:      record.Type,
			Value:     record.Value,
			TTL:       record.TTL,
			Priority:  *record.Priority,
			Weight:    *record.Weight,
			Port:      *record.Port,
			Status:    record.Status,
			Remark:    record.Remark,
			CreatedAt: record.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: record.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	resp := &BatchSyncResponse{
		DomainName:     syncResult.DomainName,
		Provider:       syncResult.Provider,
		TotalLocal:     syncResult.TotalLocal,
		TotalRemote:    syncResult.TotalRemote,
		ToAdd:          syncResult.ToAdd,
		ToUpdate:       syncResult.ToUpdate,
		ToDelete:       syncResult.ToDelete,
		AddedRecords:   addedRecords,
		UpdatedRecords: updatedRecords,
		DeletedRecords: syncResult.DeletedRecords,
		Errors:         syncResult.Errors,
		DryRun:         syncResult.DryRun,
	}

	response.ReturnData(c, resp)
}
