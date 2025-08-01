package dns

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"time"

	"api-server/internal/model/dns"
	dnsprovider "api-server/internal/provider/dns"
	dnsrepo "api-server/internal/repository/dns"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"go.uber.org/zap"
)

// CertificateService 证书管理服务
type CertificateService struct {
	certificateRepo *dnsrepo.CertificateRepository
	domainRepo      *dnsrepo.DomainRepository
	providerRepo    *dnsrepo.ProviderRepository
	recordRepo      *dnsrepo.RecordRepository
	dnsFactory      dnsprovider.DriverFactory
	logger          *zap.Logger
}

// NewCertificateService 创建证书管理服务
func NewCertificateService(
	certificateRepo *dnsrepo.CertificateRepository,
	domainRepo *dnsrepo.DomainRepository,
	providerRepo *dnsrepo.ProviderRepository,
	recordRepo *dnsrepo.RecordRepository,
	dnsFactory dnsprovider.DriverFactory,
) *CertificateService {
	return &CertificateService{
		certificateRepo: certificateRepo,
		domainRepo:      domainRepo,
		providerRepo:    providerRepo,
		recordRepo:      recordRepo,
		dnsFactory:      dnsFactory,
		logger:          zap.L().Named("certificate-service"),
	}
}

// ACMEUser ACME用户实现
type ACMEUser struct {
	Email        string
	Registration *registration.Resource
	key          *rsa.PrivateKey
}

func (u *ACMEUser) GetEmail() string {
	return u.Email
}

func (u *ACMEUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *ACMEUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

// CertificateRequest 证书申请请求
type CertificateRequest struct {
	DomainID    uint     `json:"domain_id" binding:"required"`
	Domains     []string `json:"domains" binding:"required"`
	Email       string   `json:"email" binding:"required,email"`
	ProviderID  uint     `json:"provider_id" binding:"required"`
	KeyType     string   `json:"key_type"` // RSA2048, RSA4096, ECDSA256, ECDSA384
	ValidDays   int      `json:"valid_days"`
	AutoRenew   bool     `json:"auto_renew"`
	DeployHosts []uint   `json:"deploy_hosts"`
	Remark      string   `json:"remark"`
}

// CertificateInfo 证书信息
type CertificateInfo struct {
	ID         uint       `json:"id"`
	DomainID   uint       `json:"domain_id"`
	DomainName string     `json:"domain_name"`
	CommonName string     `json:"common_name"`
	CAType     string     `json:"ca_type"`
	Status     string     `json:"status"`
	IssuedAt   *time.Time `json:"issued_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	AutoRenew  bool       `json:"auto_renew"`
	Remark     string     `json:"remark"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// IssueCertificate 申请证书
func (s *CertificateService) IssueCertificate(ctx context.Context, req *CertificateRequest, userID uint, clientIP string) (*CertificateInfo, error) {
	s.logger.Info("Starting certificate issuance",
		zap.Uint("domain_id", req.DomainID),
		zap.Strings("domains", req.Domains))

	// 验证域名
	domain, err := s.domainRepo.FindByID(req.DomainID)
	if err != nil {
		return nil, fmt.Errorf("failed to find domain: %w", err)
	}

	// 验证DNS提供商
	provider, err := s.providerRepo.FindByID(req.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find DNS provider: %w", err)
	}

	// 创建DNS驱动
	dnsDriver, err := s.createDNSDriver(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create DNS driver: %w", err)
	}

	// 创建ACME用户
	user, err := s.createACMEUser(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to create ACME user: %w", err)
	}

	// 创建ACME客户端
	config := lego.NewConfig(user)
	config.CADirURL = lego.LEDirectoryProduction // 生产环境
	// config.CADirURL = lego.LEDirectoryStaging // 测试环境

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create ACME client: %w", err)
	}

	// 注册用户
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, fmt.Errorf("failed to register ACME user: %w", err)
	}
	user.Registration = reg

	// 设置DNS-01挑战
	dnsProvider := &DNSProvider{
		driver: dnsDriver,
		logger: s.logger,
	}

	err = client.Challenge.SetDNS01Provider(dnsProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to set DNS provider: %w", err)
	}

	// 申请证书
	request := certificate.ObtainRequest{
		Domains: req.Domains,
		Bundle:  true,
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain certificate: %w", err)
	}

	// 保存证书到数据库
	now := time.Now()
	expiresAt := now.AddDate(0, 0, 90) // Let's Encrypt证书有效期90天
	cert := &dns.Certificate{
		DomainID:    req.DomainID,
		CommonName:  req.Domains[0],
		CAType:      "letsencrypt",
		Status:      "issued",
		IssuedAt:    &now,
		ExpiresAt:   &expiresAt,
		AutoRenew:   req.AutoRenew,
		LastRenewAt: &now,
		Remark:      req.Remark,
		TenantID:    domain.TenantID,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	// 加密存储证书和私钥
	certData, err := s.encryptData(certificates.Certificate)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt certificate: %w", err)
	}
	cert.CertificateEnc = string(certData)

	keyData, err := s.encryptData(certificates.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt private key: %w", err)
	}
	cert.PrivateKeyEnc = string(keyData)

	if certificates.IssuerCertificate != nil {
		chainData, err := s.encryptData(certificates.IssuerCertificate)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt CA certificate: %w", err)
		}
		cert.ChainEnc = string(chainData)
	}

	err = s.certificateRepo.Create(cert)
	if err != nil {
		return nil, fmt.Errorf("failed to save certificate: %w", err)
	}

	// 记录操作日志
	s.logCertificateOperation("issue", cert.ID, userID, clientIP, "Certificate issued successfully")

	// 如果需要部署到主机，启动部署任务
	if len(req.DeployHosts) > 0 {
		go s.deployCertificateToHosts(cert.ID, req.DeployHosts)
	}

	s.logger.Info("Certificate issued successfully",
		zap.Uint("certificate_id", cert.ID),
		zap.Strings("domains", req.Domains))

	return s.convertToInfo(cert, domain.Name), nil
}

// createACMEUser 创建ACME用户
func (s *CertificateService) createACMEUser(email string) (*ACMEUser, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	user := &ACMEUser{
		Email: email,
		key:   privateKey,
	}

	return user, nil
}

// createDNSDriver 创建DNS驱动
func (s *CertificateService) createDNSDriver(provider *dns.Provider) (dnsprovider.Driver, error) {
	config := &dnsprovider.BaseConfig{
		Type:        provider.Type,
		Name:        provider.Name,
		Credentials: make(map[string]string),
	}

	// TODO: 解密凭证信息
	// credentials := s.decryptCredentials(provider.CredentialsEnc)
	// config.Credentials = credentials

	return s.dnsFactory.CreateDriver(config)
}

// encryptData 加密数据
func (s *CertificateService) encryptData(data []byte) ([]byte, error) {
	// TODO: 实现AES-256-GCM加密
	// 这里暂时返回原数据，实际应该使用KMS或Vault进行加密
	return data, nil
}

// decryptData 解密数据
func (s *CertificateService) decryptData(encryptedData []byte) ([]byte, error) {
	// TODO: 实现AES-256-GCM解密
	// 这里暂时返回原数据，实际应该使用KMS或Vault进行解密
	return encryptedData, nil
}

// convertToInfo 转换为证书信息
func (s *CertificateService) convertToInfo(cert *dns.Certificate, domainName string) *CertificateInfo {
	return &CertificateInfo{
		ID:         cert.ID,
		DomainID:   cert.DomainID,
		DomainName: domainName,
		CommonName: cert.CommonName,
		CAType:     cert.CAType,
		Status:     cert.Status,
		IssuedAt:   cert.IssuedAt,
		ExpiresAt:  cert.ExpiresAt,
		AutoRenew:  cert.AutoRenew,
		Remark:     cert.Remark,
		CreatedAt:  cert.CreatedAt,
		UpdatedAt:  cert.UpdatedAt,
	}
}

// logCertificateOperation 记录证书操作日志
func (s *CertificateService) logCertificateOperation(operation string, certificateID uint, userID uint, clientIP string, details string) {
	// TODO: 实现操作日志记录
	s.logger.Info("Certificate operation",
		zap.String("operation", operation),
		zap.Uint("certificate_id", certificateID),
		zap.Uint("user_id", userID),
		zap.String("client_ip", clientIP),
		zap.String("details", details))
}

// deployCertificateToHosts 部署证书到主机
func (s *CertificateService) deployCertificateToHosts(certificateID uint, hostIDs []uint) {
	// TODO: 实现证书部署到主机的逻辑
	s.logger.Info("Deploying certificate to hosts",
		zap.Uint("certificate_id", certificateID),
		zap.Any("host_ids", hostIDs))
}

// ListCertificateParams 证书列表查询参数
type ListCertificateParams struct {
	Page      int                    `json:"page"`
	PageSize  int                    `json:"page_size"`
	Keyword   string                 `json:"keyword"`
	Status    string                 `json:"status"`
	CAType    string                 `json:"ca_type"`
	DomainID  uint                   `json:"domain_id"`
	AutoRenew *bool                  `json:"auto_renew"`
	TenantID  uint                   `json:"tenant_id"`
	Filters   map[string]interface{} `json:"filters"`
}

// CertificateListResponse 证书列表响应
type CertificateListResponse struct {
	List      []*CertificateInfo `json:"list"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
	TotalPage int                `json:"total_page"`
}

// ListCertificates 获取证书列表
func (s *CertificateService) ListCertificates(ctx context.Context, params *ListCertificateParams) (*CertificateListResponse, error) {
	s.logger.Info("Listing certificates", zap.Any("params", params))

	// 构建过滤条件
	filters := make(map[string]interface{})
	if params.Keyword != "" {
		filters["keyword"] = params.Keyword
	}
	if params.Status != "" {
		filters["status"] = params.Status
	}
	if params.CAType != "" {
		filters["ca_type"] = params.CAType
	}
	if params.DomainID > 0 {
		filters["domain_id"] = params.DomainID
	}
	if params.AutoRenew != nil {
		filters["auto_renew"] = *params.AutoRenew
	}
	if params.TenantID > 0 {
		filters["tenant_id"] = params.TenantID
	}

	// 合并自定义过滤条件
	for k, v := range params.Filters {
		filters[k] = v
	}

	// 计算分页参数
	offset := (params.Page - 1) * params.PageSize
	limit := params.PageSize

	// 查询数据
	certificates, total, err := s.certificateRepo.SearchWithFilters(filters, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search certificates: %w", err)
	}

	// 转换为响应格式
	certInfos := make([]*CertificateInfo, len(certificates))
	for i, cert := range certificates {
		domainName := ""
		if cert.Domain != nil {
			domainName = cert.Domain.Name
		}
		certInfos[i] = s.convertToInfo(cert, domainName)
	}

	totalPage := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &CertificateListResponse{
		List:      certInfos,
		Total:     total,
		Page:      params.Page,
		PageSize:  params.PageSize,
		TotalPage: totalPage,
	}, nil
}

// GetCertificate 获取证书详情
func (s *CertificateService) GetCertificate(ctx context.Context, id uint) (*CertificateInfo, error) {
	s.logger.Info("Getting certificate", zap.Uint("id", id))

	cert, err := s.certificateRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find certificate: %w", err)
	}

	domainName := ""
	if cert.Domain != nil {
		domainName = cert.Domain.Name
	}

	return s.convertToInfo(cert, domainName), nil
}

// RenewCertificate 续期证书
func (s *CertificateService) RenewCertificate(ctx context.Context, id uint, userID uint, clientIP string) (*CertificateInfo, error) {
	s.logger.Info("Renewing certificate", zap.Uint("id", id))

	// 获取现有证书
	cert, err := s.certificateRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find certificate: %w", err)
	}

	if cert.Status != "issued" {
		return nil, fmt.Errorf("only issued certificates can be renewed")
	}

	// 获取域名信息
	domain, err := s.domainRepo.FindByID(cert.DomainID)
	if err != nil {
		return nil, fmt.Errorf("failed to find domain: %w", err)
	}

	// 获取DNS提供商 - 证书申请时直接使用request中的ProviderID
	provider, err := s.providerRepo.FindByID(req.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find DNS provider: %w", err)
	}

	// 从证书中解析域名列表（简化处理，实际应该从证书中解析SAN）
	domains := []string{cert.CommonName}

	// 创建续期请求
	req := &CertificateRequest{
		DomainID:    cert.DomainID,
		Domains:     domains,
		Email:       "admin@" + domain.Name, // 使用域名生成邮箱，实际应该从配置获取
		ProviderID:  domain.ProviderID,
		KeyType:     "RSA2048",
		ValidDays:   90,
		AutoRenew:   cert.AutoRenew,
		DeployHosts: []uint{}, // 续期时不自动部署，需要用户手动部署
		Remark:      "Certificate renewed from ID " + fmt.Sprintf("%d", id),
	}

	// 申请新证书
	newCert, err := s.IssueCertificate(ctx, req, userID, clientIP)
	if err != nil {
		return nil, fmt.Errorf("failed to issue renewed certificate: %w", err)
	}

	// 更新原证书状态为已续期
	cert.Status = "renewed"
	cert.UpdatedBy = userID
	err = s.certificateRepo.Update(cert)
	if err != nil {
		s.logger.Warn("Failed to update old certificate status", zap.Error(err))
	}

	// 记录续期操作
	s.logCertificateOperation("renew", cert.ID, userID, clientIP,
		fmt.Sprintf("Certificate renewed, new certificate ID: %d", newCert.ID))

	return newCert, nil
}

// RevokeCertificate 吊销证书
func (s *CertificateService) RevokeCertificate(ctx context.Context, id uint, userID uint, clientIP string) error {
	s.logger.Info("Revoking certificate", zap.Uint("id", id))

	cert, err := s.certificateRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find certificate: %w", err)
	}

	if cert.Status != "issued" {
		return fmt.Errorf("only issued certificates can be revoked")
	}

	// TODO: 实现ACME证书吊销逻辑
	// 这里应该调用ACME客户端的证书吊销接口

	// 更新证书状态
	cert.Status = "revoked"
	cert.UpdatedBy = userID
	err = s.certificateRepo.Update(cert)
	if err != nil {
		return fmt.Errorf("failed to update certificate status: %w", err)
	}

	// 记录操作日志
	s.logCertificateOperation("revoke", id, userID, clientIP, "Certificate revoked")

	s.logger.Info("Certificate revoked successfully", zap.Uint("id", id))
	return nil
}

// DeployCertificate 部署证书到主机
func (s *CertificateService) DeployCertificate(ctx context.Context, id uint, hostIDs []uint, userID uint, clientIP string) error {
	s.logger.Info("Deploying certificate", zap.Uint("id", id), zap.Any("host_ids", hostIDs))

	cert, err := s.certificateRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find certificate: %w", err)
	}

	if cert.Status != "issued" {
		return fmt.Errorf("only issued certificates can be deployed")
	}

	// 异步部署证书到各个主机
	go s.deployCertificateToHosts(id, hostIDs)

	// 记录操作日志
	s.logCertificateOperation("deploy", id, userID, clientIP,
		fmt.Sprintf("Certificate deployment started for %d hosts", len(hostIDs)))

	return nil
}

// CertificateDownloadInfo 证书下载信息
type CertificateDownloadInfo struct {
	FileName string `json:"filename"`
	Content  []byte `json:"content"`
	MimeType string `json:"mime_type"`
}

// DownloadCertificate 下载证书
func (s *CertificateService) DownloadCertificate(ctx context.Context, id uint, format string) (*CertificateDownloadInfo, error) {
	s.logger.Info("Downloading certificate", zap.Uint("id", id), zap.String("format", format))

	cert, err := s.certificateRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find certificate: %w", err)
	}

	if cert.Status != "issued" {
		return nil, fmt.Errorf("only issued certificates can be downloaded")
	}

	// 解密证书数据
	certData, err := s.decryptData([]byte(cert.CertificateEnc))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt certificate: %w", err)
	}

	keyData, err := s.decryptData([]byte(cert.PrivateKeyEnc))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}

	var chainData []byte
	if cert.ChainEnc != "" {
		chainData, err = s.decryptData([]byte(cert.ChainEnc))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt certificate chain: %w", err)
		}
	}

	// 根据格式生成文件内容
	var content []byte
	var filename string
	var mimeType string

	switch format {
	case "pem":
		// PEM格式：证书 + 私钥 + 证书链
		content = append(content, certData...)
		content = append(content, []byte("\n")...)
		content = append(content, keyData...)
		if len(chainData) > 0 {
			content = append(content, []byte("\n")...)
			content = append(content, chainData...)
		}
		filename = fmt.Sprintf("%s_%d.pem", cert.CommonName, cert.ID)
		mimeType = "application/x-pem-file"

	case "crt":
		// 仅证书文件
		content = certData
		filename = fmt.Sprintf("%s_%d.crt", cert.CommonName, cert.ID)
		mimeType = "application/x-x509-ca-cert"

	case "key":
		// 仅私钥文件
		content = keyData
		filename = fmt.Sprintf("%s_%d.key", cert.CommonName, cert.ID)
		mimeType = "application/pkcs8"

	case "chain":
		// 仅证书链
		if len(chainData) == 0 {
			return nil, fmt.Errorf("certificate chain not available")
		}
		content = chainData
		filename = fmt.Sprintf("%s_%d_chain.pem", cert.CommonName, cert.ID)
		mimeType = "application/x-pem-file"

	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	return &CertificateDownloadInfo{
		FileName: filename,
		Content:  content,
		MimeType: mimeType,
	}, nil
}

// GetCertificateStats 获取证书统计信息
func (s *CertificateService) GetCertificateStats(ctx context.Context, tenantID uint) (map[string]interface{}, error) {
	s.logger.Info("Getting certificate stats", zap.Uint("tenant_id", tenantID))

	// 统计各状态证书数量
	statusCounts, err := s.certificateRepo.CountByStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to count by status: %w", err)
	}

	// 统计各CA类型证书数量
	caCounts, err := s.certificateRepo.CountByCAType()
	if err != nil {
		return nil, fmt.Errorf("failed to count by CA type: %w", err)
	}

	// 统计即将过期的证书
	expiring30, err := s.certificateRepo.FindExpiring(30)
	if err != nil {
		return nil, fmt.Errorf("failed to find expiring certificates: %w", err)
	}

	expiring7, err := s.certificateRepo.FindExpiring(7)
	if err != nil {
		return nil, fmt.Errorf("failed to find expiring certificates: %w", err)
	}

	// 统计需要自动续期的证书
	autoRenew, err := s.certificateRepo.FindAutoRenew()
	if err != nil {
		return nil, fmt.Errorf("failed to find auto-renew certificates: %w", err)
	}

	stats := map[string]interface{}{
		"status_counts":      statusCounts,
		"ca_type_counts":     caCounts,
		"expiring_30_days":   len(expiring30),
		"expiring_7_days":    len(expiring7),
		"auto_renew_pending": len(autoRenew),
		"total_certificates": 0,
	}

	// 计算总数
	total := int64(0)
	for _, count := range statusCounts {
		total += count
	}
	stats["total_certificates"] = total

	return stats, nil
}
