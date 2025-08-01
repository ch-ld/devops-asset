package dns

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"strings"
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
	dnsFactory      *dnsprovider.Factory
	logger          *zap.Logger
}

// NewCertificateService 创建证书管理服务
func NewCertificateService(
	certificateRepo *dnsrepo.CertificateRepository,
	domainRepo *dnsrepo.DomainRepository,
	providerRepo *dnsrepo.ProviderRepository,
	recordRepo *dnsrepo.RecordRepository,
	dnsFactory *dnsprovider.Factory,
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
	ID           uint      `json:"id"`
	DomainID     uint      `json:"domain_id"`
	DomainName   string    `json:"domain_name"`
	Domains      []string  `json:"domains"`
	Email        string    `json:"email"`
	Status       string    `json:"status"`
	KeyType      string    `json:"key_type"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidTo      time.Time `json:"valid_to"`
	AutoRenew    bool      `json:"auto_renew"`
	LastRenewed  time.Time `json:"last_renewed"`
	RenewCount   int       `json:"renew_count"`
	DeployStatus string    `json:"deploy_status"`
	Remark       string    `json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
		ID:           cert.ID,
		DomainID:     cert.DomainID,
		DomainName:   domainName,
		Domains:      strings.Split(cert.Domains, ","),
		Email:        cert.Email,
		Status:       cert.Status,
		KeyType:      cert.KeyType,
		ValidFrom:    cert.ValidFrom,
		ValidTo:      cert.ValidTo,
		AutoRenew:    cert.AutoRenew,
		LastRenewed:  cert.LastRenewed,
		RenewCount:   cert.RenewCount,
		DeployStatus: cert.DeployStatus,
		Remark:       cert.Remark,
		CreatedAt:    cert.CreatedAt,
		UpdatedAt:    cert.UpdatedAt,
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
		zap.Uints("host_ids", hostIDs))
}
