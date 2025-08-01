package dns

import (
	"strconv"
	"time"

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

	// TODO: 实现证书列表查询
	// certificates, total, err := h.certificateService.ListCertificates(...)

	// 临时响应
	resp := map[string]interface{}{
		"list":       []CertificateResponse{},
		"total":      0,
		"page":       page,
		"page_size":  pageSize,
		"total_page": 0,
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
	_, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// TODO: 实现获取证书详情
	// cert, err := h.certificateService.GetCertificate(uint(id))

	// 临时响应
	now := time.Now()
	cert := &CertificateResponse{
		ID:         uint(id),
		DomainID:   1,
		CommonName: "example.com",
		CAType:     "letsencrypt",
		Status:     "issued",
		IssuedAt:   &now,
		ExpiresAt:  &now,
		AutoRenew:  true,
		Remark:     "生产环境SSL证书",
		TenantID:   1,
		CreatedBy:  1,
	}

	response.ReturnData(c, cert)
}

// RenewCertificate 续期证书
func (h *CertificateHandler) RenewCertificate(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的证书ID")
		return
	}

	// 获取用户信息
	_, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// TODO: 实现证书续期
	// cert, err := h.certificateService.RenewCertificate(uint(id), userID.(uint), c.ClientIP())

	response.ReturnSuccess(c, "证书续期成功")
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

	// TODO: 实现证书吊销
	// err := h.certificateService.RevokeCertificate(uint(id), userID.(uint), c.ClientIP())

	response.ReturnSuccess(c, "证书吊销成功")
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

	// TODO: 实现证书部署
	// err := h.certificateService.DeployCertificate(uint(id), req.HostIDs, userID.(uint), c.ClientIP())

	response.ReturnSuccess(c, "证书部署成功")
}

// DownloadCertificate 下载证书
func (h *CertificateHandler) DownloadCertificate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的证书ID")
		return
	}

	format := c.DefaultQuery("format", "pem") // pem, pfx, jks

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未登录")
		return
	}

	// TODO: 实现证书下载
	// certData, filename, err := h.certificateService.DownloadCertificate(uint(id), format)

	// 设置响应头
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=certificate.pem")

	// 返回证书文件内容
	c.String(200, "# Certificate file content would be here")
}
