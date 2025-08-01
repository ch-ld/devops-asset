package dns

import (
	"strconv"

	"api-server/internal/model/dns"
	"api-server/internal/response/response"
	svc "api-server/internal/service/dns"

	"github.com/gin-gonic/gin"
)

// HTTPSMonitorHandler HTTPS监控处理器
type HTTPSMonitorHandler struct {
	httpsMonitorService *svc.HTTPSMonitorService
}

// NewHTTPSMonitorHandler 创建HTTPS监控处理器
func NewHTTPSMonitorHandler(httpsMonitorService *svc.HTTPSMonitorService) *HTTPSMonitorHandler {
	return &HTTPSMonitorHandler{
		httpsMonitorService: httpsMonitorService,
	}
}

// CreateMonitorRequest 创建监控请求
type CreateMonitorRequest struct {
	URL            string `json:"url" binding:"required,url" example:"https://example.com" comment:"监控URL"`
	Name           string `json:"name" binding:"required" example:"主站监控" comment:"监控名称"`
	CheckInterval  int    `json:"check_interval" example:"600" comment:"检查间隔(秒)"`
	Timeout        int    `json:"timeout" example:"30" comment:"超时时间(秒)"`
	AlertThreshold int    `json:"alert_threshold" example:"7" comment:"证书过期提醒天数"`
	Enabled        bool   `json:"enabled" example:"true" comment:"是否启用"`
	Remark         string `json:"remark" example:"生产环境主站" comment:"备注"`
}

// MonitorResponse 监控响应
type MonitorResponse struct {
	ID               uint   `json:"id" example:"1" comment:"监控ID"`
	URL              string `json:"url" example:"https://example.com" comment:"监控URL"`
	Name             string `json:"name" example:"主站监控" comment:"监控名称"`
	Status           string `json:"status" example:"active" comment:"监控状态"`
	CheckInterval    int    `json:"check_interval" example:"600" comment:"检查间隔(秒)"`
	Timeout          int    `json:"timeout" example:"30" comment:"超时时间(秒)"`
	AlertThreshold   int    `json:"alert_threshold" example:"7" comment:"证书过期提醒天数"`
	Enabled          bool   `json:"enabled" example:"true" comment:"是否启用"`
	LastChecked      string `json:"last_checked" example:"2024-01-01T00:00:00Z" comment:"最后检查时间"`
	LastStatus       string `json:"last_status" example:"ok" comment:"最后检查状态"`
	LastResponseTime int    `json:"last_response_time" example:"200" comment:"最后响应时间(毫秒)"`
	FailureCount     int    `json:"failure_count" example:"0" comment:"连续失败次数"`
	CertExpiry       string `json:"cert_expiry" example:"2024-04-01T00:00:00Z" comment:"证书过期时间"`
	CertDaysLeft     int    `json:"cert_days_left" example:"90" comment:"证书剩余天数"`
	LastError        string `json:"last_error" example:"" comment:"最后错误信息"`
	Remark           string `json:"remark" example:"生产环境主站" comment:"备注"`
	CreatedAt        string `json:"created_at" example:"2024-01-01T00:00:00Z" comment:"创建时间"`
	UpdatedAt        string `json:"updated_at" example:"2024-01-01T00:00:00Z" comment:"更新时间"`
}

// MonitorResultResponse 监控结果响应
type MonitorResultResponse struct {
	URL               string `json:"url" example:"https://example.com" comment:"监控URL"`
	Status            string `json:"status" example:"ok" comment:"状态"`
	StatusCode        int    `json:"status_code" example:"200" comment:"HTTP状态码"`
	ResponseTime      int    `json:"response_time" example:"200" comment:"响应时间(毫秒)"`
	CertValid         bool   `json:"cert_valid" example:"true" comment:"证书是否有效"`
	CertExpiry        string `json:"cert_expiry" example:"2024-04-01T00:00:00Z" comment:"证书过期时间"`
	CertDaysLeft      int    `json:"cert_days_left" example:"90" comment:"证书剩余天数"`
	CertIssuer        string `json:"cert_issuer" example:"Let's Encrypt" comment:"证书颁发者"`
	CertSubject       string `json:"cert_subject" example:"CN=example.com" comment:"证书主题"`
	ErrorMessage      string `json:"error_message" example:"" comment:"错误信息"`
	CheckedAt         string `json:"checked_at" example:"2024-01-01T00:00:00Z" comment:"检查时间"`
	DNSResolutionTime int    `json:"dns_resolution_time" example:"50" comment:"DNS解析时间(毫秒)"`
	TLSHandshakeTime  int    `json:"tls_handshake_time" example:"100" comment:"TLS握手时间(毫秒)"`
}

// CreateMonitor 创建HTTPS监控
func (h *HTTPSMonitorHandler) CreateMonitor(c *gin.Context) {
	var req CreateMonitorRequest
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

	// 构建监控对象
	monitor := &dns.HTTPSMonitor{
		URL:            req.URL,
		Name:           req.Name,
		CheckInterval:  req.CheckInterval,
		Timeout:        req.Timeout,
		AlertThreshold: req.AlertThreshold,
		Enabled:        req.Enabled,
		Remark:         req.Remark,
		TenantID:       tenantID.(uint),
		CreatedBy:      userID.(uint),
		UpdatedBy:      userID.(uint),
	}

	// 创建监控
	err := h.httpsMonitorService.CreateMonitor(monitor)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "创建监控失败: "+err.Error())
		return
	}

	// 转换响应
	monitorResp := h.convertToResponse(monitor)
	response.ReturnData(c, monitorResp)
}

// ListMonitors 获取监控列表
func (h *HTTPSMonitorHandler) ListMonitors(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 获取用户信息
	_, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	_, exists = c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// TODO: 实现监控列表查询
	// monitors, total, err := h.httpsMonitorService.ListMonitors(...)

	// 临时响应
	resp := map[string]interface{}{
		"list":       []MonitorResponse{},
		"total":      0,
		"page":       page,
		"page_size":  pageSize,
		"total_page": 0,
	}

	response.ReturnData(c, resp)
}

// GetMonitor 获取监控详情
func (h *HTTPSMonitorHandler) GetMonitor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的监控ID")
		return
	}

	// 获取用户信息
	_, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// TODO: 实现获取监控详情
	// monitor, err := h.httpsMonitorService.GetMonitor(uint(id))

	// 临时响应
	monitor := &MonitorResponse{
		ID:               uint(id),
		URL:              "https://example.com",
		Name:             "主站监控",
		Status:           "active",
		CheckInterval:    600,
		Timeout:          30,
		AlertThreshold:   7,
		Enabled:          true,
		LastChecked:      "2024-01-01T00:00:00Z",
		LastStatus:       "ok",
		LastResponseTime: 200,
		FailureCount:     0,
		CertExpiry:       "2024-04-01T00:00:00Z",
		CertDaysLeft:     90,
		LastError:        "",
		Remark:           "生产环境主站",
		CreatedAt:        "2024-01-01T00:00:00Z",
		UpdatedAt:        "2024-01-01T00:00:00Z",
	}

	response.ReturnData(c, monitor)
}

// UpdateMonitor 更新监控
func (h *HTTPSMonitorHandler) UpdateMonitor(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的监控ID")
		return
	}

	var req CreateMonitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请求参数错误: "+err.Error())
		return
	}

	// 获取用户信息
	_, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// TODO: 实现更新监控
	// err := h.httpsMonitorService.UpdateMonitor(uint(id), req, userID.(uint))

	response.ReturnSuccess(c)
}

// DeleteMonitor 删除监控
func (h *HTTPSMonitorHandler) DeleteMonitor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的监控ID")
		return
	}

	// 获取用户信息
	_, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 删除监控
	err = h.httpsMonitorService.DeleteMonitor(uint(id))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "删除监控失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// CheckMonitor 手动检查监控
func (h *HTTPSMonitorHandler) CheckMonitor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的监控ID")
		return
	}

	// 获取用户信息
	_, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 执行检查
	result, err := h.httpsMonitorService.CheckMonitor(c.Request.Context(), uint(id))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "检查监控失败: "+err.Error())
		return
	}

	// 转换响应
	resultResp := &MonitorResultResponse{
		URL:               result.URL,
		Status:            result.Status,
		StatusCode:        result.StatusCode,
		ResponseTime:      int(result.ResponseTime.Milliseconds()),
		CertValid:         result.CertValid,
		CertExpiry:        result.CertExpiry.Format("2006-01-02T15:04:05Z"),
		CertDaysLeft:      result.CertDaysLeft,
		CertIssuer:        result.CertIssuer,
		CertSubject:       result.CertSubject,
		ErrorMessage:      result.ErrorMessage,
		CheckedAt:         result.CheckedAt.Format("2006-01-02T15:04:05Z"),
		DNSResolutionTime: int(result.DNSResolutionTime.Milliseconds()),
		TLSHandshakeTime:  int(result.TLSHandshakeTime.Milliseconds()),
	}

	response.ReturnData(c, resultResp)
}

// GetMonitorStatistics 获取监控统计信息
func (h *HTTPSMonitorHandler) GetMonitorStatistics(c *gin.Context) {
	// 获取用户信息
	_, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 获取统计信息
	stats, err := h.httpsMonitorService.GetMonitorStatistics(c.Request.Context())
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取统计信息失败: "+err.Error())
		return
	}

	response.ReturnData(c, stats)
}

// GetExpiringCertificates 获取即将过期的证书
func (h *HTTPSMonitorHandler) GetExpiringCertificates(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 获取即将过期的证书
	monitors, err := h.httpsMonitorService.GetExpiringCertificates(c.Request.Context(), days)
	if err != nil {
		response.ReturnError(c, response.INTERNAL_ERROR, "获取即将过期证书失败: "+err.Error())
		return
	}

	// 转换响应
	var responses []MonitorResponse
	for _, monitor := range monitors {
		responses = append(responses, *h.convertToResponse(monitor))
	}

	response.ReturnData(c, responses)
}

// convertToResponse 转换为响应格式
func (h *HTTPSMonitorHandler) convertToResponse(monitor *dns.HTTPSMonitor) *MonitorResponse {
	resp := &MonitorResponse{
		ID:               monitor.ID,
		URL:              monitor.URL,
		Name:             monitor.Name,
		Status:           monitor.Status,
		CheckInterval:    monitor.CheckInterval,
		Timeout:          monitor.Timeout,
		AlertThreshold:   monitor.AlertThreshold,
		Enabled:          monitor.Enabled,
		LastChecked:      monitor.LastChecked.Format("2006-01-02T15:04:05Z"),
		LastStatus:       monitor.LastStatus,
		LastResponseTime: monitor.LastResponseTime,
		FailureCount:     monitor.FailureCount,
		LastError:        monitor.LastError,
		Remark:           monitor.Remark,
		CreatedAt:        monitor.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        monitor.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if monitor.CertExpiry != nil {
		resp.CertExpiry = monitor.CertExpiry.Format("2006-01-02T15:04:05Z")
	}

	if monitor.CertDaysLeft != nil {
		resp.CertDaysLeft = *monitor.CertDaysLeft
	}

	return resp
}
