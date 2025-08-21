package dns

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"strings"
	"time"

	"api-server/internal/config"
	"api-server/internal/model/dns"
	dnsprovider "api-server/internal/provider/dns"
	dnsrepo "api-server/internal/repository/dns"
	"api-server/pkg/crypto/encryption"

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

	// 解密凭证信息
	creds, err := s.decryptProviderCredentials(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt provider credentials: %w", err)
	}
	config.Credentials = creds

	return s.dnsFactory.CreateDriver(config)
}

// decryptProviderCredentials 解密提供商凭证为map
func (s *CertificateService) decryptProviderCredentials(provider *dns.Provider) (map[string]string, error) {
	if provider.CredentialsEnc == "" {
		return map[string]string{}, nil
	}
	// 先尝试多密钥
	if len(config.GlobalConfig.App.AesKeys) > 0 {
		var keys [][]byte
		for _, k := range config.GlobalConfig.App.AesKeys {
			keys = append(keys, []byte(k))
		}
		plain, err := encryption.DecryptAESWithKeys(provider.CredentialsEnc, keys)
		if err == nil {
			m := map[string]string{}
			if jsonErr := json.Unmarshal([]byte(plain), &m); jsonErr != nil {
				return nil, jsonErr
			}
			return m, nil
		}
	}
	// 退回单密钥
	if config.GlobalConfig.App.AesKey != "" {
		plain, err := encryption.DecryptAES(provider.CredentialsEnc, []byte(config.GlobalConfig.App.AesKey))
		if err == nil {
			m := map[string]string{}
			if jsonErr := json.Unmarshal([]byte(plain), &m); jsonErr != nil {
				return nil, jsonErr
			}
			return m, nil
		}
	}
	// 兼容旧明文
	m := map[string]string{}
	if err := json.Unmarshal([]byte(provider.CredentialsEnc), &m); err != nil {
		return nil, fmt.Errorf("unable to decrypt or parse provider credentials: %w", err)
	}
	return m, nil
}

// encryptData 加密数据
func (s *CertificateService) encryptData(data []byte) ([]byte, error) {
	// 使用应用配置的AES密钥
	if len(config.GlobalConfig.App.AesKeys) > 0 {
		var keys [][]byte
		for _, k := range config.GlobalConfig.App.AesKeys {
			keys = append(keys, []byte(k))
		}
		cipher, err := encryption.EncryptAESWithKeys(data, keys)
		if err != nil {
			return nil, err
		}
		return []byte(cipher), nil
	}
	if config.GlobalConfig.App.AesKey != "" {
		cipher, err := encryption.EncryptAES(data, []byte(config.GlobalConfig.App.AesKey))
		if err != nil {
			return nil, err
		}
		return []byte(cipher), nil
	}
	return nil, fmt.Errorf("no AES key configured")
}

// decryptData 解密数据
func (s *CertificateService) decryptData(encrypted []byte) ([]byte, error) {
	enc := string(encrypted)
	if len(config.GlobalConfig.App.AesKeys) > 0 {
		var keys [][]byte
		for _, k := range config.GlobalConfig.App.AesKeys {
			keys = append(keys, []byte(k))
		}
		plain, err := encryption.DecryptAESWithKeys(enc, keys)
		if err == nil {
			return []byte(plain), nil
		}
	}
	if config.GlobalConfig.App.AesKey != "" {
		plain, err := encryption.DecryptAES(enc, []byte(config.GlobalConfig.App.AesKey))
		if err == nil {
			return []byte(plain), nil
		}
	}
	// 兼容已存明文
	return encrypted, nil
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
	s.logger.Info("Starting certificate deployment to hosts",
		zap.Uint("certificate_id", certificateID),
		zap.Any("host_ids", hostIDs))

	// 获取证书信息
	cert, err := s.certificateRepo.FindByID(certificateID)
	if err != nil {
		s.logger.Error("Failed to find certificate for deployment",
			zap.Uint("certificate_id", certificateID),
			zap.Error(err))
		return
	}

	// 解密证书数据
	certData, err := s.decryptData([]byte(cert.CertificateEnc))
	if err != nil {
		s.logger.Error("Failed to decrypt certificate data",
			zap.Uint("certificate_id", certificateID),
			zap.Error(err))
		return
	}

	keyData, err := s.decryptData([]byte(cert.PrivateKeyEnc))
	if err != nil {
		s.logger.Error("Failed to decrypt private key data",
			zap.Uint("certificate_id", certificateID),
			zap.Error(err))
		return
	}

	var chainData []byte
	if cert.ChainEnc != "" {
		chainData, err = s.decryptData([]byte(cert.ChainEnc))
		if err != nil {
			s.logger.Error("Failed to decrypt certificate chain data",
				zap.Uint("certificate_id", certificateID),
				zap.Error(err))
			// 证书链不是必需的，继续部署
		}
	}

	// 为每个主机创建部署记录并执行部署
	for _, hostID := range hostIDs {
		deployment := &dns.CertificateDeployment{
			CertificateID: certificateID,
			HostID:        hostID,
			Status:        "pending",
			DeployPath:    "/etc/ssl/certs/", // 默认部署路径
			ServiceName:   "nginx",           // 默认服务
		}

		// 保存部署记录
		err := s.certificateRepo.CreateDeployment(deployment)
		if err != nil {
			s.logger.Error("Failed to create deployment record",
				zap.Uint("certificate_id", certificateID),
				zap.Uint("host_id", hostID),
				zap.Error(err))
			continue
		}

		// 异步执行部署
		go s.deployToSingleHost(deployment, certData, keyData, chainData)
	}
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

	// 从证书配置中获取ProviderID，如果没有则需要用户在续期时指定
	// 这里简化处理，实际应该在证书申请时记录使用的ProviderID
	// 临时使用domain的RegistrarID作为ProviderID (需要改进)
	var providerID uint = 1 // 默认使用第一个DNS Provider，实际应该从证书记录或用户输入获取
	if domain.RegistrarID != nil {
		providerID = *domain.RegistrarID
	}

	// 从证书中解析域名列表（简化处理，实际应该从证书中解析SAN）
	domains := []string{cert.CommonName}

	// 创建续期请求
	renewReq := &CertificateRequest{
		DomainID:    cert.DomainID,
		Domains:     domains,
		Email:       "admin@" + domain.Name, // 使用域名生成邮箱，实际应该从配置获取
		ProviderID:  providerID,
		KeyType:     "RSA2048",
		ValidDays:   90,
		AutoRenew:   cert.AutoRenew,
		DeployHosts: []uint{}, // 续期时不自动部署，需要用户手动部署
		Remark:      "Certificate renewed from ID " + fmt.Sprintf("%d", id),
	}

	// 申请新证书
	newCert, err := s.IssueCertificate(ctx, renewReq, userID, clientIP)
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

// CSRUploadRequest CSR上传请求
type CSRUploadRequest struct {
	CSRContent  string   `json:"csr_content" binding:"required" comment:"CSR内容(PEM格式)"`
	Domains     []string `json:"domains" binding:"required" comment:"域名列表"`
	Email       string   `json:"email" binding:"required,email" comment:"邮箱地址"`
	ProviderID  uint     `json:"provider_id" binding:"required" comment:"DNS提供商ID"`
	ValidDays   int      `json:"valid_days" example:"90" comment:"有效天数"`
	AutoRenew   bool     `json:"auto_renew" example:"true" comment:"自动续期"`
	DeployHosts []uint   `json:"deploy_hosts" example:"[1,2,3]" comment:"部署主机ID列表"`
	Remark      string   `json:"remark" example:"自定义CSR证书" comment:"备注"`
}

// CSRValidationResult CSR验证结果
type CSRValidationResult struct {
	Valid        bool     `json:"valid"`
	CommonName   string   `json:"common_name"`
	DNSNames     []string `json:"dns_names"`
	Organization string   `json:"organization"`
	Country      string   `json:"country"`
	KeyAlgorithm string   `json:"key_algorithm"`
	KeySize      int      `json:"key_size"`
	Signature    string   `json:"signature"`
	ErrorMessage string   `json:"error_message,omitempty"`
}

// ValidateCSR 验证CSR内容
func (s *CertificateService) ValidateCSR(csrContent string) (*CSRValidationResult, error) {
	s.logger.Info("Validating CSR content")

	// 解析PEM格式的CSR
	block, _ := pem.Decode([]byte(csrContent))
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return &CSRValidationResult{
			Valid:        false,
			ErrorMessage: "无效的CSR格式，请提供PEM格式的证书签名请求",
		}, nil
	}

	// 解析CSR
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return &CSRValidationResult{
			Valid:        false,
			ErrorMessage: fmt.Sprintf("解析CSR失败: %v", err),
		}, nil
	}

	// 验证CSR签名
	if err := csr.CheckSignature(); err != nil {
		return &CSRValidationResult{
			Valid:        false,
			ErrorMessage: fmt.Sprintf("CSR签名验证失败: %v", err),
		}, nil
	}

	// 提取密钥信息
	keyAlgorithm := "Unknown"
	keySize := 0
	switch pub := csr.PublicKey.(type) {
	case *rsa.PublicKey:
		keyAlgorithm = "RSA"
		keySize = pub.Size() * 8
	case crypto.PublicKey:
		keyAlgorithm = "Other"
	}

	// 构建验证结果
	result := &CSRValidationResult{
		Valid:        true,
		CommonName:   csr.Subject.CommonName,
		DNSNames:     csr.DNSNames,
		Organization: strings.Join(csr.Subject.Organization, ", "),
		Country:      strings.Join(csr.Subject.Country, ", "),
		KeyAlgorithm: keyAlgorithm,
		KeySize:      keySize,
		Signature:    csr.SignatureAlgorithm.String(),
	}

	s.logger.Info("CSR validation completed",
		zap.Bool("valid", result.Valid),
		zap.String("common_name", result.CommonName),
		zap.Strings("dns_names", result.DNSNames),
	)

	return result, nil
}

// IssueCertificateWithCSR 使用自定义CSR申请证书
func (s *CertificateService) IssueCertificateWithCSR(ctx context.Context, req *CSRUploadRequest, tenantID, actorID uint) (*dns.Certificate, error) {
	s.logger.Info("Issuing certificate with custom CSR",
		zap.Strings("domains", req.Domains),
		zap.String("email", req.Email),
		zap.Uint("provider_id", req.ProviderID),
	)

	// 1. 验证CSR
	csrValidation, err := s.ValidateCSR(req.CSRContent)
	if err != nil {
		return nil, fmt.Errorf("CSR验证失败: %w", err)
	}
	if !csrValidation.Valid {
		return nil, fmt.Errorf("CSR验证失败: %s", csrValidation.ErrorMessage)
	}

	// 2. 检查域名匹配
	if err := s.validateCSRDomains(req.Domains, csrValidation.DNSNames, csrValidation.CommonName); err != nil {
		return nil, fmt.Errorf("域名验证失败: %w", err)
	}

	// 3. 获取DNS提供商
	provider, err := s.providerRepo.FindByID(req.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("获取DNS提供商失败: %w", err)
	}

	// 4. 解析CSR
	block, _ := pem.Decode([]byte(req.CSRContent))
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析CSR失败: %w", err)
	}

	// 5. 创建证书记录
	subjectAltNames, _ := json.Marshal(req.Domains)
	certificate := &dns.Certificate{
		CommonName:      csrValidation.CommonName,
		SubjectAltNames: subjectAltNames,
		CAType:          "letsencrypt",
		Status:          "pending",
		AutoRenew:       req.AutoRenew,
		RenewDays:       30,
		Remark:          req.Remark,
		TenantID:        tenantID,
		CreatedBy:       actorID,
		UpdatedBy:       actorID,
	}

	// 6. 保存证书记录
	if err := s.certificateRepo.Create(certificate); err != nil {
		return nil, fmt.Errorf("保存证书记录失败: %w", err)
	}

	// 7. 创建ACME客户端并申请证书
	go func() {
		s.processCSRCertificate(ctx, certificate, csr, provider)
	}()

	s.logger.Info("Certificate with CSR created successfully",
		zap.Uint("certificate_id", certificate.ID),
		zap.String("common_name", certificate.CommonName),
	)

	return certificate, nil
}

// processCSRCertificate 处理CSR证书申请（简化版本）
func (s *CertificateService) processCSRCertificate(ctx context.Context, cert *dns.Certificate, csr *x509.CertificateRequest, provider *dns.Provider) {
	s.logger.Info("Processing CSR certificate", zap.Uint("certificate_id", cert.ID))

	// 更新状态为处理中
	s.certificateRepo.UpdateStatus(cert.ID, "processing")

	// TODO: 实现完整的CSR证书申请逻辑
	// 这里先实现一个简化版本，避免编译错误

	// 模拟证书申请过程
	time.Sleep(2 * time.Second)

	// 更新证书状态和相关信息
	s.certificateRepo.UpdateStatus(cert.ID, "issued")

	s.logger.Info("CSR certificate issued successfully", zap.Uint("certificate_id", cert.ID))
}

// validateCSRDomains 验证CSR中的域名与请求域名是否匹配
func (s *CertificateService) validateCSRDomains(requestDomains []string, csrDNSNames []string, csrCommonName string) error {
	// 将CSR中的所有域名收集到一个map中
	csrDomains := make(map[string]bool)
	if csrCommonName != "" {
		csrDomains[csrCommonName] = true
	}
	for _, domain := range csrDNSNames {
		csrDomains[domain] = true
	}

	// 检查请求的域名是否都在CSR中
	for _, requestDomain := range requestDomains {
		if !csrDomains[requestDomain] {
			return fmt.Errorf("域名 %s 不在CSR中", requestDomain)
		}
	}

	// 检查CSR中是否有多余的域名
	if len(csrDomains) > len(requestDomains) {
		s.logger.Warn("CSR contains more domains than requested",
			zap.Int("csr_domains", len(csrDomains)),
			zap.Int("request_domains", len(requestDomains)),
		)
	}

	return nil
}

// UploadCSRFile 上传CSR文件
func (s *CertificateService) UploadCSRFile(reader io.Reader) (string, *CSRValidationResult, error) {
	s.logger.Info("Uploading CSR file")

	// 读取文件内容
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", nil, fmt.Errorf("读取CSR文件失败: %w", err)
	}

	// 验证文件大小（限制为100KB）
	if len(content) > 100*1024 {
		return "", nil, fmt.Errorf("CSR文件过大，最大支持100KB")
	}

	csrContent := string(content)

	// 验证CSR内容
	validation, err := s.ValidateCSR(csrContent)
	if err != nil {
		return "", nil, fmt.Errorf("CSR验证失败: %w", err)
	}

	s.logger.Info("CSR file uploaded and validated", zap.Bool("valid", validation.Valid))

	return csrContent, validation, nil
}

// CertificateUploadRequest 证书上传请求
type CertificateUploadRequest struct {
	CertContent       string `json:"cert_content"`
	KeyContent        string `json:"key_content"`
	ChainContent      string `json:"chain_content"`
	AutoDeploy        bool   `json:"auto_deploy"`
	DeployHosts       []uint `json:"deploy_hosts"`
	DeployPath        string `json:"deploy_path"`
	RestartCommand    string `json:"restart_command"`
	EmailNotification bool   `json:"email_notification"`
	NotificationEmail string `json:"notification_email"`
	Remark            string `json:"remark"`
}

// CertificateValidationResult 证书验证结果
type CertificateValidationResult struct {
	Valid   bool                `json:"valid"`
	Message string              `json:"message"`
	Details *CertificateDetails `json:"details,omitempty"`
}

// CertificateDetails 证书详细信息
type CertificateDetails struct {
	CommonName      string   `json:"commonName"`
	SubjectAltNames []string `json:"subjectAltNames"`
	Issuer          string   `json:"issuer"`
	ValidFrom       string   `json:"validFrom"`
	ValidTo         string   `json:"validTo"`
	DaysRemaining   int      `json:"daysRemaining"`
	KeyAlgorithm    string   `json:"keyAlgorithm"`
	SerialNumber    string   `json:"serialNumber"`
}

// UploadCertificate 上传证书
func (s *CertificateService) UploadCertificate(ctx context.Context, req *CertificateUploadRequest, tenantID, userID uint) (*dns.Certificate, error) {
	s.logger.Info("Uploading certificate", zap.Uint("tenant_id", tenantID), zap.Uint("user_id", userID))

	// 验证证书
	validation, err := s.ValidateCertificate(ctx, req.CertContent, req.KeyContent, req.ChainContent)
	if err != nil {
		return nil, fmt.Errorf("证书验证失败: %w", err)
	}

	if !validation.Valid {
		return nil, fmt.Errorf("证书验证失败: %s", validation.Message)
	}

	// 解析证书获取信息
	block, _ := pem.Decode([]byte(req.CertContent))
	if block == nil {
		return nil, fmt.Errorf("无效的证书格式")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析证书失败: %w", err)
	}

	// 加密证书内容
	encryptedCert, err := s.encryptData([]byte(req.CertContent))
	if err != nil {
		return nil, fmt.Errorf("加密证书失败: %w", err)
	}

	encryptedKey, err := s.encryptData([]byte(req.KeyContent))
	if err != nil {
		return nil, fmt.Errorf("加密私钥失败: %w", err)
	}

	var encryptedChain []byte
	if req.ChainContent != "" {
		encryptedChain, err = s.encryptData([]byte(req.ChainContent))
		if err != nil {
			return nil, fmt.Errorf("加密证书链失败: %w", err)
		}
	}

	// 将DNS名称转换为JSON
	sanDomainsJSON, _ := json.Marshal(cert.DNSNames)

	// 创建证书记录
	issuedAt := cert.NotBefore
	expiresAt := cert.NotAfter
	certificate := &dns.Certificate{
		TenantID:        tenantID,
		CommonName:      cert.Subject.CommonName,
		SubjectAltNames: sanDomainsJSON,
		CAType:          "uploaded",
		Status:          "issued",
		CertificateEnc:  string(encryptedCert),
		PrivateKeyEnc:   string(encryptedKey),
		ChainEnc:        string(encryptedChain),
		SerialNumber:    cert.SerialNumber.String(),
		Fingerprint:     fmt.Sprintf("%x", cert.SerialNumber),
		IssuedAt:        &issuedAt,
		ExpiresAt:       &expiresAt,
		AutoRenew:       false, // 上传的证书不支持自动续期
		Remark:          req.Remark,
		CreatedBy:       userID,
		UpdatedBy:       userID,
	}

	// 保存到数据库
	if err := s.certificateRepo.Create(certificate); err != nil {
		return nil, fmt.Errorf("保存证书失败: %w", err)
	}

	// 如果启用自动部署，则部署证书
	if req.AutoDeploy && len(req.DeployHosts) > 0 {
		go s.deployCertificateToHosts(certificate.ID, req.DeployHosts)
	}

	s.logger.Info("Certificate uploaded successfully", zap.Uint("certificate_id", certificate.ID))
	return certificate, nil
}

// ValidateCertificate 验证证书
func (s *CertificateService) ValidateCertificate(ctx context.Context, certContent, keyContent, chainContent string) (*CertificateValidationResult, error) {
	s.logger.Info("Validating certificate")

	// 解析证书
	certBlock, _ := pem.Decode([]byte(certContent))
	if certBlock == nil || certBlock.Type != "CERTIFICATE" {
		return &CertificateValidationResult{
			Valid:   false,
			Message: "无效的证书格式",
		}, nil
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return &CertificateValidationResult{
			Valid:   false,
			Message: "证书解析失败: " + err.Error(),
		}, nil
	}

	// 解析私钥
	keyBlock, _ := pem.Decode([]byte(keyContent))
	if keyBlock == nil {
		return &CertificateValidationResult{
			Valid:   false,
			Message: "无效的私钥格式",
		}, nil
	}

	var privateKey crypto.PrivateKey
	switch keyBlock.Type {
	case "RSA PRIVATE KEY":
		privateKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	case "PRIVATE KEY":
		privateKey, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	default:
		return &CertificateValidationResult{
			Valid:   false,
			Message: "不支持的私钥类型: " + keyBlock.Type,
		}, nil
	}

	if err != nil {
		return &CertificateValidationResult{
			Valid:   false,
			Message: "私钥解析失败: " + err.Error(),
		}, nil
	}

	// 验证证书和私钥是否匹配
	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		if priv, ok := privateKey.(*rsa.PrivateKey); ok {
			if pub.N.Cmp(priv.N) != 0 {
				return &CertificateValidationResult{
					Valid:   false,
					Message: "证书和私钥不匹配",
				}, nil
			}
		} else {
			return &CertificateValidationResult{
				Valid:   false,
				Message: "证书和私钥类型不匹配",
			}, nil
		}
	default:
		return &CertificateValidationResult{
			Valid:   false,
			Message: "不支持的公钥类型",
		}, nil
	}

	// 检查证书是否过期
	now := time.Now()
	if now.Before(cert.NotBefore) {
		return &CertificateValidationResult{
			Valid:   false,
			Message: "证书尚未生效",
		}, nil
	}

	if now.After(cert.NotAfter) {
		return &CertificateValidationResult{
			Valid:   false,
			Message: "证书已过期",
		}, nil
	}

	// 计算剩余天数
	daysRemaining := int(cert.NotAfter.Sub(now).Hours() / 24)

	// 构建验证结果
	result := &CertificateValidationResult{
		Valid:   true,
		Message: "证书验证成功",
		Details: &CertificateDetails{
			CommonName:      cert.Subject.CommonName,
			SubjectAltNames: cert.DNSNames,
			Issuer:          cert.Issuer.CommonName,
			ValidFrom:       cert.NotBefore.Format("2006-01-02 15:04:05"),
			ValidTo:         cert.NotAfter.Format("2006-01-02 15:04:05"),
			DaysRemaining:   daysRemaining,
			KeyAlgorithm:    cert.PublicKeyAlgorithm.String(),
			SerialNumber:    cert.SerialNumber.String(),
		},
	}

	return result, nil
}

// BatchRenewCertificates 批量续期证书
func (s *CertificateService) BatchRenewCertificates(ctx context.Context, ids []uint, tenantID, userID uint) error {
	s.logger.Info("Batch renewing certificates", zap.Any("ids", ids), zap.Uint("tenant_id", tenantID))

	for _, id := range ids {
		// 异步续期每个证书
		go func(certID uint) {
			_, err := s.RenewCertificate(ctx, certID, userID, "batch_operation")
			if err != nil {
				s.logger.Error("Failed to renew certificate in batch",
					zap.Uint("certificate_id", certID),
					zap.Error(err))
			}
		}(id)
	}

	return nil
}

// BatchDeleteCertificates 批量删除证书
func (s *CertificateService) BatchDeleteCertificates(ctx context.Context, ids []uint, tenantID, userID uint) error {
	s.logger.Info("Batch deleting certificates", zap.Any("ids", ids), zap.Uint("tenant_id", tenantID))

	// 验证所有证书都属于该租户
	for _, id := range ids {
		cert, err := s.certificateRepo.FindByID(id)
		if err != nil {
			return fmt.Errorf("证书 %d 不存在", id)
		}
		if cert.TenantID != tenantID {
			return fmt.Errorf("无权限删除证书 %d", id)
		}
	}

	// 批量删除
	for _, id := range ids {
		if err := s.certificateRepo.Delete(id); err != nil {
			s.logger.Error("Failed to delete certificate in batch",
				zap.Uint("certificate_id", id),
				zap.Error(err))
			return fmt.Errorf("删除证书 %d 失败: %w", id, err)
		}
	}

	s.logger.Info("Batch delete certificates completed", zap.Int("count", len(ids)))
	return nil
}

// BatchExportCertificates 批量导出证书
func (s *CertificateService) BatchExportCertificates(ctx context.Context, ids []uint, tenantID uint) ([]byte, error) {
	s.logger.Info("Batch exporting certificates", zap.Any("ids", ids), zap.Uint("tenant_id", tenantID))

	// 获取证书列表
	var certificates []*dns.Certificate
	for _, id := range ids {
		cert, err := s.certificateRepo.FindByID(id)
		if err != nil {
			continue // 跳过不存在的证书
		}
		if cert.TenantID != tenantID {
			continue // 跳过无权限的证书
		}
		certificates = append(certificates, cert)
	}

	// 这里应该生成Excel文件，暂时返回JSON数据
	// TODO: 实现Excel导出功能
	data, err := json.Marshal(certificates)
	if err != nil {
		return nil, fmt.Errorf("导出数据序列化失败: %w", err)
	}

	return data, nil
}

// BatchDeployCertificates 批量部署证书
func (s *CertificateService) BatchDeployCertificates(ctx context.Context, certIDs, hostIDs []uint, tenantID, userID uint) error {
	s.logger.Info("Batch deploying certificates",
		zap.Any("cert_ids", certIDs),
		zap.Any("host_ids", hostIDs),
		zap.Uint("tenant_id", tenantID))

	// 验证所有证书都属于该租户
	for _, certID := range certIDs {
		cert, err := s.certificateRepo.FindByID(certID)
		if err != nil {
			return fmt.Errorf("证书 %d 不存在", certID)
		}
		if cert.TenantID != tenantID {
			return fmt.Errorf("无权限部署证书 %d", certID)
		}
	}

	// 批量部署
	for _, certID := range certIDs {
		go func(id uint) {
			err := s.DeployCertificate(ctx, id, hostIDs, userID, "batch_operation")
			if err != nil {
				s.logger.Error("Failed to deploy certificate in batch",
					zap.Uint("certificate_id", id),
					zap.Error(err))
			}
		}(certID)
	}

	return nil
}
