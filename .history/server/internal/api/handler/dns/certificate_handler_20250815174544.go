package dns

import (
	"path/filepath"
	"strconv"
	"strings"

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

// ValidateCSRRequest CSR验证请求
type ValidateCSRRequest struct {
	CSRContent string `json:"csr_content" binding:"required" example:"-----BEGIN CERTIFICATE REQUEST-----..." comment:"CSR内容(PEM格式)"`
}

// ValidateCSR 验证CSR内容
// @Summary 验证CSR内容
// @Description 验证CSR格式和内容的有效性，返回域名、组织等信息
// @Tags 证书管理
// @Accept json
// @Produce json
// @Param request body ValidateCSRRequest true "CSR验证请求"
// @Success 200 {object} response.Response{data=svc.CSRValidationResult} "验证结果"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dns/certificates/validate-csr [post]
func (h *CertificateHandler) ValidateCSR(c *gin.Context) {
	var req ValidateCSRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请求参数错误: "+err.Error())
		return
	}

	// 验证CSR
	result, err := h.certificateService.ValidateCSR(req.CSRContent)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "CSR验证失败: "+err.Error())
		return
	}

	response.ReturnData(c, result)
}

// UploadCSRFile 上传CSR文件
// @Summary 上传CSR文件
// @Description 上传CSR文件并验证其内容
// @Tags 证书管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSR文件"
// @Success 200 {object} response.Response{data=map[string]interface{}} "上传结果包含CSR内容和验证结果"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dns/certificates/upload-csr [post]
func (h *CertificateHandler) UploadCSRFile(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "文件上传失败: "+err.Error())
		return
	}

	// 验证文件类型
	allowedExts := []string{".csr", ".pem", ".txt"}
	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	isValidExt := false
	for _, ext := range allowedExts {
		if fileExt == ext {
			isValidExt = true
			break
		}
	}

	if !isValidExt {
		response.ReturnError(c, response.INVALID_ARGUMENT, "文件格式不正确，请上传.csr、.pem或.txt格式的CSR文件")
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "打开文件失败: "+err.Error())
		return
	}
	defer src.Close()

	// 上传并验证CSR
	csrContent, validation, err := h.certificateService.UploadCSRFile(src)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "处理CSR文件失败: "+err.Error())
		return
	}

	result := map[string]interface{}{
		"csr_content": csrContent,
		"validation":  validation,
		"filename":    file.Filename,
		"size":        file.Size,
	}

	response.ReturnData(c, result)
}

// IssueCertificateWithCSR 使用自定义CSR申请证书
// @Summary 使用自定义CSR申请证书
// @Description 使用用户提供的CSR申请SSL证书
// @Tags 证书管理
// @Accept json
// @Produce json
// @Param request body svc.CSRUploadRequest true "CSR证书申请请求"
// @Success 200 {object} response.Response{data=dns.Certificate} "申请成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dns/certificates/issue-with-csr [post]
func (h *CertificateHandler) IssueCertificateWithCSR(c *gin.Context) {
	var req svc.CSRUploadRequest
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
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息获取失败")
		return
	}

	// 使用CSR申请证书
	certificate, err := h.certificateService.IssueCertificateWithCSR(
		c.Request.Context(),
		&req,
		tenantID.(uint),
		userID.(uint),
	)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "证书申请失败: "+err.Error())
		return
	}

	response.ReturnData(c, certificate)
}

// UploadCertificateRequest 上传证书请求
type UploadCertificateRequest struct {
	CertContent       string `json:"cert_content" binding:"required" comment:"证书内容"`
	KeyContent        string `json:"key_content" binding:"required" comment:"私钥内容"`
	ChainContent      string `json:"chain_content" comment:"证书链内容"`
	AutoDeploy        bool   `json:"auto_deploy" comment:"自动部署"`
	DeployHosts       []uint `json:"deploy_hosts" comment:"部署主机ID列表"`
	DeployPath        string `json:"deploy_path" comment:"部署路径"`
	RestartCommand    string `json:"restart_command" comment:"重启命令"`
	EmailNotification bool   `json:"email_notification" comment:"邮件通知"`
	NotificationEmail string `json:"notification_email" comment:"通知邮箱"`
	Remark            string `json:"remark" comment:"备注"`
}

// UploadCertificate 上传证书
func (h *CertificateHandler) UploadCertificate(c *gin.Context) {
	var req UploadCertificateRequest
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

	// 上传证书
	certificate, err := h.certificateService.UploadCertificate(
		c.Request.Context(),
		&svc.CertificateUploadRequest{
			CertContent:       req.CertContent,
			KeyContent:        req.KeyContent,
			ChainContent:      req.ChainContent,
			AutoDeploy:        req.AutoDeploy,
			DeployHosts:       req.DeployHosts,
			DeployPath:        req.DeployPath,
			RestartCommand:    req.RestartCommand,
			EmailNotification: req.EmailNotification,
			NotificationEmail: req.NotificationEmail,
			Remark:            req.Remark,
		},
		tenantID.(uint),
		userID.(uint),
	)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "证书上传失败: "+err.Error())
		return
	}

	response.ReturnData(c, certificate)
}

// ValidateCertificateRequest 验证证书请求
type ValidateCertificateRequest struct {
	CertContent  string `json:"certContent" binding:"required" comment:"证书内容"`
	KeyContent   string `json:"keyContent" binding:"required" comment:"私钥内容"`
	ChainContent string `json:"chainContent" comment:"证书链内容"`
}

// ValidateCertificate 验证证书
func (h *CertificateHandler) ValidateCertificate(c *gin.Context) {
	var req ValidateCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "请求参数错误: "+err.Error())
		return
	}

	// 验证证书
	result, err := h.certificateService.ValidateCertificate(
		c.Request.Context(),
		req.CertContent,
		req.KeyContent,
		req.ChainContent,
	)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "证书验证失败: "+err.Error())
		return
	}

	response.ReturnData(c, result)
}

// BatchRenewRequest 批量续期请求
type BatchRenewRequest struct {
	IDs []uint `json:"ids" binding:"required" comment:"证书ID列表"`
}

// BatchRenewCertificates 批量续期证书
func (h *CertificateHandler) BatchRenewCertificates(c *gin.Context) {
	var req BatchRenewRequest
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

	// 批量续期证书
	err := h.certificateService.BatchRenewCertificates(
		c.Request.Context(),
		req.IDs,
		tenantID.(uint),
		userID.(uint),
	)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "批量续期失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required" comment:"证书ID列表"`
}

// BatchDeleteCertificates 批量删除证书
func (h *CertificateHandler) BatchDeleteCertificates(c *gin.Context) {
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

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 批量删除证书
	err := h.certificateService.BatchDeleteCertificates(
		c.Request.Context(),
		req.IDs,
		tenantID.(uint),
		userID.(uint),
	)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "批量删除失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}

// BatchExportCertificates 批量导出证书
func (h *CertificateHandler) BatchExportCertificates(c *gin.Context) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		response.ReturnError(c, response.INVALID_ARGUMENT, "证书ID列表不能为空")
		return
	}

	// 解析ID列表
	idStrs := strings.Split(idsStr, ",")
	var ids []uint
	for _, idStr := range idStrs {
		id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32)
		if err != nil {
			response.ReturnError(c, response.INVALID_ARGUMENT, "无效的证书ID: "+idStr)
			return
		}
		ids = append(ids, uint(id))
	}

	// 获取用户信息
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 导出证书报告
	exportData, err := h.certificateService.BatchExportCertificates(
		c.Request.Context(),
		ids,
		tenantID.(uint),
	)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "导出失败: "+err.Error())
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=certificates_report.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", exportData)
}

// BatchDeployRequest 批量部署请求
type BatchDeployRequest struct {
	CertIDs []uint `json:"cert_ids" binding:"required" comment:"证书ID列表"`
	HostIDs []uint `json:"host_ids" binding:"required" comment:"主机ID列表"`
}

// BatchDeployCertificates 批量部署证书
func (h *CertificateHandler) BatchDeployCertificates(c *gin.Context) {
	var req BatchDeployRequest
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

	// 批量部署证书
	err := h.certificateService.BatchDeployCertificates(
		c.Request.Context(),
		req.CertIDs,
		req.HostIDs,
		tenantID.(uint),
		userID.(uint),
	)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "批量部署失败: "+err.Error())
		return
	}

	response.ReturnSuccess(c)
}
