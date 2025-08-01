package dns

import (
	"strconv"

	"api-server/internal/response/response"
	svc "api-server/internal/service/dns"

	"github.com/gin-gonic/gin"
)

// CertificateHandler 证书管理处理器
type CertificateHandler struct {
	certificateService *svc.CertificateService
}

// NewCertificateHandler 创建证书管理处理器
func NewCertificateHandler(certificateService *svc.CertificateService) *CertificateHandler {
	return &CertificateHandler{
		certificateService: certificateService,
	}
}

// IssueCertificateRequest 申请证书请求
type IssueCertificateRequest struct {
	DomainID    uint     `json:"domain_id" binding:"required" example:"1" comment:"域名ID"`
	Domains     []string `json:"domains" binding:"required" example:"[\"example.com\", \"www.example.com\"]" comment:"域名列表"`
	Email       string   `json:"email" binding:"required,email" example:"admin@example.com" comment:"邮箱地址"`
	ProviderID  uint     `json:"provider_id" binding:"required" example:"1" comment:"DNS提供商ID"`
	KeyType     string   `json:"key_type" example:"RSA2048" comment:"密钥类型"`
	ValidDays   int      `json:"valid_days" example:"90" comment:"有效天数"`
	AutoRenew   bool     `json:"auto_renew" example:"true" comment:"自动续期"`
	DeployHosts []uint   `json:"deploy_hosts" example:"[1,2,3]" comment:"部署主机ID列表"`
	Remark      string   `json:"remark" example:"生产环境SSL证书" comment:"备注"`
}

// IssueCertificate 申请证书
func (h *CertificateHandler) IssueCertificate(c *gin.Context) {
	var req IssueCertificateRequest
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

	// 转换请求
	certReq := &svc.CertificateRequest{
		DomainID:    req.DomainID,
		Domains:     req.Domains,
		Email:       req.Email,
		ProviderID:  req.ProviderID,
		KeyType:     req.KeyType,
		ValidDays:   req.ValidDays,
		AutoRenew:   req.AutoRenew,
		DeployHosts: req.DeployHosts,
		Remark:      req.Remark,
	}

	// 申请证书
	cert, err := h.certificateService.IssueCertificate(c.Request.Context(), certReq, userID.(uint), c.ClientIP())
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "申请证书失败: "+err.Error())
		return
	}

	// 转换响应
	certResp := &CertificateResponse{
		ID:         cert.ID,
		DomainID:   cert.DomainID,
		CommonName: cert.CommonName,
		CAType:     cert.CAType,
		Status:     cert.Status,
		IssuedAt:   cert.IssuedAt,
		ExpiresAt:  cert.ExpiresAt,
		AutoRenew:  cert.AutoRenew,
		Remark:     cert.Remark,
		TenantID:   cert.ID, // 临时使用ID，实际应该从cert获取
		CreatedBy:  cert.ID, // 临时使用ID，实际应该从cert获取
	}

	response.ReturnData(c, certResp)
}

// ListCertificates 获取证书列表
func (h *CertificateHandler) ListCertificates(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	caType := c.Query("ca_type")
	domainID, _ := strconv.ParseUint(c.Query("domain_id"), 10, 32)
	autoRenewStr := c.Query("auto_renew")

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

	// 构建查询参数
	params := &svc.ListCertificateParams{
		Page:     page,
		PageSize: pageSize,
		Keyword:  keyword,
		Status:   status,
		CAType:   caType,
		TenantID: tenantID.(uint),
	}

	if domainID > 0 {
		params.DomainID = uint(domainID)
	}

	if autoRenewStr != "" {
		if autoRenew, err := strconv.ParseBool(autoRenewStr); err == nil {
			params.AutoRenew = &autoRenew
		}
	}

	// 查询证书列表
	result, err := h.certificateService.ListCertificates(c.Request.Context(), params)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "查询证书列表失败: "+err.Error())
		return
	}

	// 转换为Handler响应格式
	certificateResponses := make([]*CertificateResponse, len(result.List))
	for i, cert := range result.List {
		certificateResponses[i] = &CertificateResponse{
			ID:         cert.ID,
			DomainID:   cert.DomainID,
			CommonName: cert.CommonName,
			CAType:     cert.CAType,
			Status:     cert.Status,
			IssuedAt:   cert.IssuedAt,
			ExpiresAt:  cert.ExpiresAt,
			AutoRenew:  cert.AutoRenew,
			Remark:     cert.Remark,
			CreatedBy:  userID.(uint), // 简化处理，实际应该从cert获取
			CreatedAt:  cert.CreatedAt,
			UpdatedAt:  cert.UpdatedAt,
		}
	}

	resp := map[string]interface{}{
		"list":       certificateResponses,
		"total":      result.Total,
		"page":       result.Page,
		"page_size":  result.PageSize,
		"total_page": result.TotalPage,
	}

	response.ReturnData(c, resp)
}

// GetCertificate 获取证书详情
func (h *CertificateHandler) GetCertificate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的证书ID")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 获取证书详情
	cert, err := h.certificateService.GetCertificate(c.Request.Context(), uint(id))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取证书详情失败: "+err.Error())
		return
	}

	// 转换为响应格式
	certResp := &CertificateResponse{
		ID:         cert.ID,
		DomainID:   cert.DomainID,
		CommonName: cert.CommonName,
		CAType:     cert.CAType,
		Status:     cert.Status,
		IssuedAt:   cert.IssuedAt,
		ExpiresAt:  cert.ExpiresAt,
		AutoRenew:  cert.AutoRenew,
		Remark:     cert.Remark,
		CreatedBy:  userID.(uint), // 简化处理
		CreatedAt:  cert.CreatedAt,
		UpdatedAt:  cert.UpdatedAt,
	}

	response.ReturnData(c, certResp)
}

// RenewCertificate 续期证书
func (h *CertificateHandler) RenewCertificate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的证书ID")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 续期证书
	newCert, err := h.certificateService.RenewCertificate(c.Request.Context(), uint(id), userID.(uint), c.ClientIP())
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "证书续期失败: "+err.Error())
		return
	}

	// 转换响应
	certResp := &CertificateResponse{
		ID:         newCert.ID,
		DomainID:   newCert.DomainID,
		CommonName: newCert.CommonName,
		CAType:     newCert.CAType,
		Status:     newCert.Status,
		IssuedAt:   newCert.IssuedAt,
		ExpiresAt:  newCert.ExpiresAt,
		AutoRenew:  newCert.AutoRenew,
		Remark:     newCert.Remark,
		CreatedBy:  userID.(uint),
		CreatedAt:  newCert.CreatedAt,
		UpdatedAt:  newCert.UpdatedAt,
	}

	response.ReturnData(c, certResp)
}

// RevokeCertificate 吊销证书
func (h *CertificateHandler) RevokeCertificate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的证书ID")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 吊销证书
	err = h.certificateService.RevokeCertificate(c.Request.Context(), uint(id), userID.(uint), c.ClientIP())
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "证书吊销失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// DeployCertificate 部署证书
func (h *CertificateHandler) DeployCertificate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的证书ID")
		return
	}

	var req struct {
		HostIDs []uint `json:"host_ids" binding:"required" comment:"主机ID列表"`
	}
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

	// 部署证书
	err = h.certificateService.DeployCertificate(c.Request.Context(), uint(id), req.HostIDs, userID.(uint), c.ClientIP())
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "证书部署失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// DownloadCertificate 下载证书
func (h *CertificateHandler) DownloadCertificate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的证书ID")
		return
	}

	format := c.DefaultQuery("format", "pem") // pem, crt, key, chain

	// 获取用户信息
	_, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// 下载证书
	downloadInfo, err := h.certificateService.DownloadCertificate(c.Request.Context(), uint(id), format)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "证书下载失败: "+err.Error())
		return
	}

	// 设置响应头
	c.Header("Content-Type", downloadInfo.MimeType)
	c.Header("Content-Disposition", "attachment; filename="+downloadInfo.FileName)
	c.Header("Content-Length", strconv.Itoa(len(downloadInfo.Content)))

	// 返回证书文件内容
	c.Data(200, downloadInfo.MimeType, downloadInfo.Content)
}

// GetCertificateStats 获取证书统计信息
func (h *CertificateHandler) GetCertificateStats(c *gin.Context) {
	// 获取用户信息
	_, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 获取证书统计信息
	stats, err := h.certificateService.GetCertificateStats(c.Request.Context(), tenantID.(uint))
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取证书统计信息失败: "+err.Error())
		return
	}

	response.ReturnData(c, stats)
}
