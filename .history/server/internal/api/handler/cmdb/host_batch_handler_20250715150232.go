package cmdb

import (
	"api-server/internal/model/cmdb"
	"api-server/internal/response/response"
	"api-server/internal/service/cmdb"
	"time"

	"github.com/gin-gonic/gin"
)

// HostBatchHandler 主机批量操作处理器
type HostBatchHandler struct {
	hostService *cmdb.HostService
}

// NewHostBatchHandler 创建主机批量操作处理器
func NewHostBatchHandler(hostService *cmdb.HostService) *HostBatchHandler {
	return &HostBatchHandler{
		hostService: hostService,
	}
}

// BatchIds 批量操作ID请求
type BatchIds struct {
	Ids []uint `json:"ids" binding:"required"`
}

// BatchStatusRequest 批量状态变更请求
type BatchStatusRequest struct {
	Ids    []uint `json:"ids" binding:"required"`
	Status string `json:"status" binding:"required"`
}

// BatchMoveRequest 批量移动请求
type BatchMoveRequest struct {
	Ids     []uint `json:"ids" binding:"required"`
	GroupID *uint  `json:"group_id"`
}

// BatchLifecycleRequest 批量生命周期请求
type BatchLifecycleRequest struct {
	Ids       []uint  `json:"ids" binding:"required"`
	ExpiredAt *string `json:"expired_at"`
	Status    string  `json:"status"`
	Recycle   bool    `json:"recycle"`
}

// BatchCustomFieldRequest 批量自定义字段请求
type BatchCustomFieldRequest struct {
	Ids         []uint                 `json:"ids" binding:"required"`
	ExtraFields map[string]interface{} `json:"extra_fields" binding:"required"`
}

// BatchSSHRequest 批量SSH请求
type BatchSSHRequest struct {
	Ids     []uint `json:"ids" binding:"required"`
	Cmd     string `json:"cmd" binding:"required"`
	Timeout int    `json:"timeout"`
}

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	Total     int      `json:"total"`
	Success   int      `json:"success"`
	Failed    int      `json:"failed"`
	FailedMsg []string `json:"failed_msg"`
}

// HandleBatchDelete 批量删除主机
func (h *HostBatchHandler) HandleBatchDelete(c *gin.Context) {
	var req BatchIds
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	if err := h.hostService.BatchDeleteHosts(req.Ids); err != nil {
		response.SystemError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "批量删除成功",
	})
}

// HandleBatchChangeStatus 批量变更主机状态
func (h *HostBatchHandler) HandleBatchChangeStatus(c *gin.Context) {
	var req BatchStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 验证状态值是否合法
	validStatus := map[string]bool{
		"running":   true,
		"stopped":   true,
		"rebooting": true,
	}

	if !validStatus[req.Status] {
		response.ParamError(c, "无效的状态值")
		return
	}

	batchReq := cmdb.BatchChangeStatusRequest{
		IDs:    req.Ids,
		Status: req.Status,
	}

	if err := h.hostService.BatchChangeStatus(batchReq); err != nil {
		response.SystemError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "批量状态变更成功",
	})
}

// HandleBatchMove 批量移动主机
func (h *HostBatchHandler) HandleBatchMove(c *gin.Context) {
	var req BatchMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	var results BatchOperationResult
	results.Total = len(req.Ids)

	for _, id := range req.Ids {
		if err := h.hostService.MoveHost(id, req.GroupID); err != nil {
			results.Failed++
			results.FailedMsg = append(results.FailedMsg, err.Error())
		} else {
			results.Success++
		}
	}

	response.Success(c, results)
}

// HandleBatchLifecycle 批量设置主机生命周期
func (h *HostBatchHandler) HandleBatchLifecycle(c *gin.Context) {
	var req BatchLifecycleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	var expiredAt *time.Time
	if req.ExpiredAt != nil && *req.ExpiredAt != "" {
		t, err := time.Parse("2006-01-02", *req.ExpiredAt)
		if err != nil {
			response.ParamError(c, "过期时间格式错误，应为YYYY-MM-DD")
			return
		}
		expiredAt = &t
	}

	batchReq := cmdb.BatchLifecycleRequest{
		IDs:       req.Ids,
		ExpiredAt: expiredAt,
		Status:    req.Status,
		Recycle:   req.Recycle,
	}

	if err := h.hostService.BatchLifecycleHosts(batchReq); err != nil {
		response.SystemError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "批量设置生命周期成功",
	})
}

// HandleBatchSetCustomFields 批量设置自定义字段
func (h *HostBatchHandler) HandleBatchSetCustomFields(c *gin.Context) {
	var req BatchCustomFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	batchReq := cmdb.BatchSetCustomFieldsRequest{
		IDs:         req.Ids,
		ExtraFields: req.ExtraFields,
	}

	if err := h.hostService.BatchSetCustomFields(batchReq); err != nil {
		response.SystemError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "批量设置自定义字段成功",
	})
}

// HandleBatchSSH 批量执行SSH命令
func (h *HostBatchHandler) HandleBatchSSH(c *gin.Context) {
	var req BatchSSHRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 设置默认超时时间
	if req.Timeout <= 0 {
		req.Timeout = 30 // 默认30秒
	}

	batchReq := cmdb.BatchSSHRequest{
		IDs:     req.Ids,
		Cmd:     req.Cmd,
		Timeout: req.Timeout,
	}

	results := h.hostService.BatchSSH(batchReq)

	response.Success(c, results)
}

// HandleBatchStart 批量启动主机
func (h *HostBatchHandler) HandleBatchStart(c *gin.Context) {
	var req BatchIds
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	batchReq := cmdb.BatchChangeStatusRequest{
		IDs:    req.Ids,
		Status: "running",
	}

	if err := h.hostService.BatchChangeStatus(batchReq); err != nil {
		response.SystemError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "批量启动主机成功",
	})
}

// HandleBatchStop 批量停止主机
func (h *HostBatchHandler) HandleBatchStop(c *gin.Context) {
	var req BatchIds
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	batchReq := cmdb.BatchChangeStatusRequest{
		IDs:    req.Ids,
		Status: "stopped",
	}

	if err := h.hostService.BatchChangeStatus(batchReq); err != nil {
		response.SystemError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "批量停止主机成功",
	})
}

// HandleBatchReboot 批量重启主机
func (h *HostBatchHandler) HandleBatchReboot(c *gin.Context) {
	var req BatchIds
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	batchReq := cmdb.BatchChangeStatusRequest{
		IDs:    req.Ids,
		Status: "rebooting",
	}

	if err := h.hostService.BatchChangeStatus(batchReq); err != nil {
		response.SystemError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "批量重启主机成功",
	})
}

// RegisterRoutes 注册路由
func (h *HostBatchHandler) RegisterRoutes(router *gin.RouterGroup) {
	batchRouter := router.Group("/hosts")
	{
		batchRouter.POST("/batch_delete", h.HandleBatchDelete)
		batchRouter.PUT("/batch_status", h.HandleBatchChangeStatus)
		batchRouter.POST("/batch_move", h.HandleBatchMove)
		batchRouter.PUT("/lifecycle", h.HandleBatchLifecycle)
		batchRouter.PUT("/custom_fields", h.HandleBatchSetCustomFields)
		batchRouter.POST("/batch_ssh", h.HandleBatchSSH)

		// 操作快捷方法
		batchRouter.POST("/batch_start", h.HandleBatchStart)
		batchRouter.POST("/batch_stop", h.HandleBatchStop)
		batchRouter.POST("/batch_reboot", h.HandleBatchReboot)
	}
}
