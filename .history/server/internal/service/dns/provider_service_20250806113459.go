package dns

import (
	"api-server/internal/config"
	"api-server/internal/model/dns"
	dnsprovider "api-server/internal/provider/dns"
	repo "api-server/internal/repository/dns"
	"api-server/pkg/crypto/encryption"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProviderService DNS提供商业务逻辑服务
type ProviderService struct {
	providerRepo  *repo.ProviderRepository
	changeLogRepo *repo.ChangeLogRepository
	domainRepo    *repo.DomainRepository
	db            *gorm.DB
}

// NewProviderService 创建DNS提供商业务服务实例
func NewProviderService(
	providerRepo *repo.ProviderRepository,
	changeLogRepo *repo.ChangeLogRepository,
	domainRepo *repo.DomainRepository,
	db *gorm.DB,
) *ProviderService {
	return &ProviderService{
		providerRepo:  providerRepo,
		changeLogRepo: changeLogRepo,
		domainRepo:    domainRepo,
		db:            db,
	}
}

// CreateProvider 创建DNS提供商
func (s *ProviderService) CreateProvider(provider *dns.Provider, actorID uint, clientIP string) error {
	// 验证提供商配置
	if err := s.validateProvider(provider); err != nil {
		return fmt.Errorf("DNS提供商配置验证失败: %w", err)
	}

	// 检查提供商名称是否已存在
	exists, err := s.providerRepo.ExistsByName(provider.Name)
	if err != nil {
		return fmt.Errorf("检查提供商名称是否存在失败: %w", err)
	}
	if exists {
		return errors.New("提供商名称已存在")
	}

	// 加密凭证信息
	if err := s.encryptCredentials(provider); err != nil {
		return fmt.Errorf("加密凭证信息失败: %w", err)
	}

	// 设置创建信息
	provider.CreatedBy = actorID
	provider.UpdatedBy = actorID

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 创建提供商
		if err := s.providerRepo.Create(provider); err != nil {
			return fmt.Errorf("创建DNS提供商失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "provider",
			ResourceID:   provider.ID,
			Action:       "create",
			Description:  fmt.Sprintf("创建DNS提供商: %s (%s)", provider.Name, provider.Type),
			NewData:      s.providerToJSON(provider),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     provider.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录DNS提供商创建日志失败", zap.Error(err))
		}

		return nil
	})
}

// UpdateProvider 更新DNS提供商
func (s *ProviderService) UpdateProvider(provider *dns.Provider, actorID uint, clientIP string) error {
	// 获取原始数据
	oldProvider, err := s.providerRepo.FindByID(provider.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("DNS提供商不存在")
		}
		return fmt.Errorf("获取DNS提供商信息失败: %w", err)
	}

	// 验证提供商配置
	if err := s.validateProvider(provider); err != nil {
		return fmt.Errorf("DNS提供商配置验证失败: %w", err)
	}

	// 如果名称发生变化，检查新名称是否已存在
	if oldProvider.Name != provider.Name {
		exists, err := s.providerRepo.ExistsByName(provider.Name)
		if err != nil {
			return fmt.Errorf("检查提供商名称是否存在失败: %w", err)
		}
		if exists {
			return errors.New("提供商名称已存在")
		}
	}

	// 如果凭证信息发生变化，重新加密
	if provider.CredentialsEnc != oldProvider.CredentialsEnc {
		if err := s.encryptCredentials(provider); err != nil {
			return fmt.Errorf("加密凭证信息失败: %w", err)
		}
	}

	// 设置更新信息
	provider.UpdatedBy = actorID

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 更新提供商
		if err := s.providerRepo.Update(provider); err != nil {
			return fmt.Errorf("更新DNS提供商失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "provider",
			ResourceID:   provider.ID,
			Action:       "update",
			Description:  fmt.Sprintf("更新DNS提供商: %s (%s)", provider.Name, provider.Type),
			OldData:      s.providerToJSON(oldProvider),
			NewData:      s.providerToJSON(provider),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     provider.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录DNS提供商更新日志失败", zap.Error(err))
		}

		return nil
	})
}

// DeleteProvider 删除DNS提供商
func (s *ProviderService) DeleteProvider(id uint, actorID uint, clientIP string) error {
	// 获取提供商信息
	provider, err := s.providerRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("DNS提供商不存在")
		}
		return fmt.Errorf("获取DNS提供商信息失败: %w", err)
	}

	// TODO: 检查是否有关联的DNS记录
	// 这里可以添加业务规则，比如有关联记录时不允许删除

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除提供商
		if err := s.providerRepo.Delete(id); err != nil {
			return fmt.Errorf("删除DNS提供商失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "provider",
			ResourceID:   id,
			Action:       "delete",
			Description:  fmt.Sprintf("删除DNS提供商: %s (%s)", provider.Name, provider.Type),
			OldData:      s.providerToJSON(provider),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     provider.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录DNS提供商删除日志失败", zap.Error(err))
		}

		return nil
	})
}

// GetProvider 获取DNS提供商详情
func (s *ProviderService) GetProvider(id uint) (*dns.Provider, error) {
	provider, err := s.providerRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("DNS提供商不存在")
		}
		return nil, fmt.Errorf("获取DNS提供商信息失败: %w", err)
	}
	return provider, nil
}

// ListProviders 获取DNS提供商列表
func (s *ProviderService) ListProviders(tenantID uint, filters map[string]interface{}, limit, offset int) ([]*dns.Provider, int64, error) {
	// 添加租户过滤
	if filters == nil {
		filters = make(map[string]interface{})
	}
	filters["tenant_id"] = tenantID

	providers, total, err := s.providerRepo.SearchWithFilters(filters, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("获取DNS提供商列表失败: %w", err)
	}

	return providers, total, nil
}

// TestProvider 测试DNS提供商连接
func (s *ProviderService) TestProvider(id uint) (*dnsprovider.TestResult, error) {
	provider, err := s.providerRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("DNS提供商不存在")
		}
		return nil, fmt.Errorf("获取DNS提供商信息失败: %w", err)
	}

	// 解密凭证信息
	credentials, err := s.decryptCredentials(provider)
	if err != nil {
		return nil, fmt.Errorf("解密凭证信息失败: %w", err)
	}

	// 测试连接
	result := dnsprovider.TestDriver(provider.Type, provider.Name, credentials)

	// 更新测试结果
	var testResult, errorMessage string
	if result.Success {
		testResult = "success"
	} else {
		testResult = "failed"
		errorMessage = result.ErrorMsg
	}

	if err := s.providerRepo.UpdateTestResult(id, testResult, errorMessage); err != nil {
		zap.L().Error("更新DNS提供商测试结果失败", zap.Error(err))
	}

	return result, nil
}

// TestConnectionTemporary 临时测试DNS提供商连接（不保存提供商）
func (s *ProviderService) TestConnectionTemporary(providerType, providerName string, credentials map[string]string) (*dnsprovider.TestResult, error) {
	// 直接测试连接，不需要从数据库获取提供商信息
	result := dnsprovider.TestDriver(providerType, providerName, credentials)

	return result, nil
}

// SetDefaultProvider 设置默认DNS提供商
func (s *ProviderService) SetDefaultProvider(id uint, actorID uint, clientIP string) error {
	// 验证提供商是否存在
	provider, err := s.providerRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("DNS提供商不存在")
		}
		return fmt.Errorf("获取DNS提供商信息失败: %w", err)
	}

	// 设置默认提供商
	if err := s.providerRepo.SetDefault(id); err != nil {
		return fmt.Errorf("设置默认DNS提供商失败: %w", err)
	}

	// 记录变更日志
	changeLog := &dns.ChangeLog{
		ResourceType: "provider",
		ResourceID:   id,
		Action:       "set_default",
		Description:  fmt.Sprintf("设置默认DNS提供商: %s", provider.Name),
		Status:       "success",
		ClientIP:     clientIP,
		TenantID:     provider.TenantID,
		ActorID:      actorID,
	}
	if err := s.changeLogRepo.Create(changeLog); err != nil {
		zap.L().Error("记录设置默认提供商日志失败", zap.Error(err))
	}

	return nil
}

// validateProvider 验证DNS提供商配置
func (s *ProviderService) validateProvider(provider *dns.Provider) error {
	if provider.Name == "" {
		return errors.New("提供商名称不能为空")
	}

	if provider.Type == "" {
		return errors.New("提供商类型不能为空")
	}

	// 验证支持的提供商类型
	supportedTypes := []string{"route53", "aliyun", "godaddy", "cloudflare", "dnspod"}
	isSupported := false
	for _, supportedType := range supportedTypes {
		if provider.Type == supportedType {
			isSupported = true
			break
		}
	}
	if !isSupported {
		return fmt.Errorf("不支持的提供商类型: %s", provider.Type)
	}

	return nil
}

// encryptCredentials 加密凭证信息
func (s *ProviderService) encryptCredentials(provider *dns.Provider) error {
	if provider.CredentialsEnc == "" {
		return nil // 没有凭证信息需要加密
	}

	encrypted, err := s.encryptPassword(provider.CredentialsEnc)
	if err != nil {
		return fmt.Errorf("加密凭证信息失败: %w", err)
	}
	provider.CredentialsEnc = encrypted

	zap.L().Info("成功加密DNS提供商凭证", zap.String("provider", provider.Name))
	return nil
}

// encryptPassword 加密密码/凭证信息
func (s *ProviderService) encryptPassword(plainText string) (string, error) {
	if plainText == "" {
		return "", nil
	}

	// 获取AES密钥
	var keys [][]byte
	if len(config.GlobalConfig.App.AesKeys) > 0 {
		// 使用多密钥加密（只用第一个密钥加密）
		for _, keyStr := range config.GlobalConfig.App.AesKeys {
			keys = append(keys, []byte(keyStr))
		}
		return encryption.EncryptAESWithKeys([]byte(plainText), keys)
	} else if config.GlobalConfig.App.AesKey != "" {
		// 使用单密钥加密
		return encryption.EncryptAES([]byte(plainText), []byte(config.GlobalConfig.App.AesKey))
	}

	return "", fmt.Errorf("no encryption key configured")
}

// decryptCredentials 解密凭证信息
func (s *ProviderService) decryptCredentials(provider *dns.Provider) (map[string]string, error) {
	if provider.CredentialsEnc == "" {
		return make(map[string]string), nil
	}

	var decrypted string
	var err error

	// 首先尝试解密（新格式的加密数据）
	decrypted, err = s.decryptPassword(provider.CredentialsEnc)
	if err != nil {
		// 如果解密失败，可能是旧格式的明文JSON数据，直接使用
		zap.L().Warn("解密失败，尝试直接解析为JSON（可能是旧格式数据）",
			zap.String("provider", provider.Name),
			zap.Error(err))
		decrypted = provider.CredentialsEnc

		// 异步升级旧格式数据到新格式（如果可能的话）
		go s.upgradeProviderCredentials(provider)
	}

	// 解析JSON到map
	var credentials map[string]string
	if err := json.Unmarshal([]byte(decrypted), &credentials); err != nil {
		return nil, fmt.Errorf("解析凭证信息失败: %w", err)
	}

	zap.L().Info("成功解密/解析DNS提供商凭证",
		zap.String("provider", provider.Name),
		zap.Int("credentials_count", len(credentials)))

	return credentials, nil
}

// decryptPassword 解密密码/凭证信息
func (s *ProviderService) decryptPassword(encryptedText string) (string, error) {
	if encryptedText == "" {
		return "", nil
	}

	// 获取AES密钥
	var keys [][]byte
	if len(config.GlobalConfig.App.AesKeys) > 0 {
		// 使用多密钥解密
		for _, keyStr := range config.GlobalConfig.App.AesKeys {
			keys = append(keys, []byte(keyStr))
		}
		return encryption.DecryptAESWithKeys(encryptedText, keys)
	} else if config.GlobalConfig.App.AesKey != "" {
		// 使用单密钥解密
		return encryption.DecryptAES(encryptedText, []byte(config.GlobalConfig.App.AesKey))
	}

	return "", fmt.Errorf("no encryption key configured")
}

// upgradeProviderCredentials 异步升级旧格式凭证到新格式
func (s *ProviderService) upgradeProviderCredentials(provider *dns.Provider) {
	// 检查是否是明文JSON格式（可以正确解析但不是base64）
	var testCredentials map[string]string
	if err := json.Unmarshal([]byte(provider.CredentialsEnc), &testCredentials); err != nil {
		// 如果不能解析为JSON，说明不是我们预期的格式，跳过升级
		return
	}

	// 尝试解密，如果能解密说明已经是新格式了
	if _, err := s.decryptPassword(provider.CredentialsEnc); err == nil {
		// 已经是加密格式，无需升级
		return
	}

	// 确认是旧格式，进行升级
	encrypted, err := s.encryptPassword(provider.CredentialsEnc)
	if err != nil {
		zap.L().Error("升级提供商凭证格式失败",
			zap.String("provider", provider.Name),
			zap.Error(err))
		return
	}

	// 更新数据库
	if err := s.providerRepo.UpdateCredentials(provider.ID, encrypted); err != nil {
		zap.L().Error("更新提供商加密凭证失败",
			zap.String("provider", provider.Name),
			zap.Error(err))
		return
	}

	zap.L().Info("成功升级提供商凭证格式",
		zap.String("provider", provider.Name))
}

// SyncProviderDomains 同步单个提供商的域名
func (s *ProviderService) SyncProviderDomains(providerID uint, actorID uint, clientIP string) (map[string]interface{}, error) {
	zap.L().Info("开始同步DNS提供商域名", zap.Uint("provider_id", providerID))
	// 获取提供商信息
	provider, err := s.providerRepo.FindByID(providerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("DNS提供商不存在")
		}
		return nil, fmt.Errorf("获取DNS提供商信息失败: %w", err)
	}

	// 解密凭证信息
	credentials, err := s.decryptCredentials(provider)
	if err != nil {
		return nil, fmt.Errorf("解密凭证信息失败: %w", err)
	}

	// 创建DNS驱动
	zap.L().Info("开始创建DNS驱动",
		zap.String("provider", provider.Name),
		zap.String("type", provider.Type),
		zap.Int("credentials_count", len(credentials)))

	driver, err := dnsprovider.CreateDriver(provider.Type, provider.Name, credentials)
	if err != nil {
		zap.L().Error("创建DNS驱动失败",
			zap.String("provider", provider.Name),
			zap.String("type", provider.Type),
			zap.Error(err))
		return nil, fmt.Errorf("创建DNS驱动失败: %w", err)
	}

	zap.L().Info("DNS驱动创建成功", zap.String("provider", provider.Name))

	// 创建带超时的context（10秒超时）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 先测试连接以验证凭证
	zap.L().Info("开始测试DNS提供商连接", zap.String("provider", provider.Name))
	testResult := driver.Test(ctx)
	if !testResult.Success {
		zap.L().Error("DNS提供商连接测试失败",
			zap.String("provider", provider.Name),
			zap.String("error", testResult.ErrorMsg))
		return nil, fmt.Errorf("DNS提供商连接测试失败: %s", testResult.ErrorMsg)
	}
	zap.L().Info("DNS提供商连接测试成功", zap.String("provider", provider.Name))

	// 获取所有域名区域
	zap.L().Info("开始获取DNS提供商域名列表", zap.String("provider", provider.Name))
	zones, err := driver.ListZones(ctx, nil)
	if err != nil {
		// 检查是否属于未实现的功能
		if errors.Is(err, dnsprovider.ErrNotImplemented) {
			return nil, fmt.Errorf("当前提供商类型暂未支持域名同步 (Not Implemented)")
		}
		// 记录详细错误日志
		zap.L().Error("获取DNS提供商域名列表失败",
			zap.String("provider", provider.Name),
			zap.String("type", provider.Type),
			zap.Error(err))

		// 增强错误处理，识别常见的认证错误
		errMsg := err.Error()
		if strings.Contains(errMsg, "status 401") || strings.Contains(errMsg, "Unauthorized") {
			return nil, fmt.Errorf("DNS提供商认证失败，请检查API密钥是否正确: %w", err)
		}
		if strings.Contains(errMsg, "status 403") || strings.Contains(errMsg, "Forbidden") {
			return nil, fmt.Errorf("DNS提供商权限不足，请检查API密钥权限: %w", err)
		}
		if strings.Contains(errMsg, "status 404") {
			return nil, fmt.Errorf("DNS提供商API端点不存在，请检查提供商配置: %w", err)
		}
		if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "context deadline exceeded") {
			return nil, fmt.Errorf("DNS提供商API请求超时，请稍后重试: %w", err)
		}
		if strings.Contains(errMsg, "connectex") || strings.Contains(errMsg, "connection attempt failed") {
			return nil, fmt.Errorf("无法连接到DNS提供商API服务器，请检查：1) 网络连接是否正常 2) 防火墙是否阻止访问 3) 是否需要配置代理。详细错误: %v", err)
		}
		if strings.Contains(errMsg, "context deadline exceeded") {
			return nil, fmt.Errorf("DNS提供商API请求超时（10秒），请检查网络连接或稍后重试")
		}
		return nil, fmt.Errorf("获取域名列表失败: %w", err)
	}

	zap.L().Info("成功获取DNS提供商域名列表",
		zap.String("provider", provider.Name),
		zap.Int("zones_count", len(zones)))

	// 同步域名到数据库
	syncedCount := 0
	errorCount := 0
	var errors []string

	for _, zone := range zones {
		// 检查域名是否已存在
		exists, err := s.domainRepo.ExistsByName(zone.Name)
		if err != nil {
			zap.L().Error("检查域名是否存在失败",
				zap.String("domain", zone.Name),
				zap.Error(err))
			errorCount++
			errors = append(errors, fmt.Sprintf("检查域名 %s 失败: %v", zone.Name, err))
			continue
		}

		if !exists {
			// 创建新域名
			domain := &dns.Domain{
				Name:          zone.Name,
				RegistrarID:   &providerID,
				RegistrarType: provider.Type,
				Status:        "active",
				ExpiresAt:     nil, // 过期时间需要通过其他API获取
				AutoRenew:     false,
				TenantID:      provider.TenantID,
				CreatedBy:     actorID,
				UpdatedBy:     actorID,
			}

			if err := s.domainRepo.Create(domain); err != nil {
				zap.L().Error("创建域名失败",
					zap.String("domain", zone.Name),
					zap.Error(err))
				errorCount++
				errors = append(errors, fmt.Sprintf("创建域名 %s 失败: %v", zone.Name, err))
				continue
			}

			zap.L().Info("成功同步域名",
				zap.String("provider", provider.Name),
				zap.String("zone", zone.Name))
			syncedCount++
		} else {
			// 域名已存在，更新提供商关联
			domain, err := s.domainRepo.FindByName(zone.Name)
			if err != nil {
				zap.L().Error("获取域名信息失败",
					zap.String("domain", zone.Name),
					zap.Error(err))
				errorCount++
				errors = append(errors, fmt.Sprintf("获取域名 %s 信息失败: %v", zone.Name, err))
				continue
			}

			// 如果域名还没有关联提供商，则关联当前提供商
			if domain.RegistrarID == nil || *domain.RegistrarID != providerID {
				domain.RegistrarID = &providerID
				domain.RegistrarType = provider.Type
				domain.UpdatedBy = actorID
				if err := s.domainRepo.Update(domain); err != nil {
					zap.L().Error("更新域名提供商关联失败",
						zap.String("domain", zone.Name),
						zap.Error(err))
					errorCount++
					errors = append(errors, fmt.Sprintf("更新域名 %s 提供商关联失败: %v", zone.Name, err))
					continue
				}
			}

			zap.L().Info("域名已存在，跳过同步",
				zap.String("provider", provider.Name),
				zap.String("zone", zone.Name))
		}
	}

	// 记录同步日志
	changeLog := &dns.ChangeLog{
		ResourceType: "provider",
		ResourceID:   providerID,
		Action:       "sync_domains",
		Description:  fmt.Sprintf("同步提供商 %s 的域名，成功: %d, 失败: %d", provider.Name, syncedCount, errorCount),
		Status:       "success",
		ClientIP:     clientIP,
		TenantID:     provider.TenantID,
		ActorID:      actorID,
	}
	if err := s.changeLogRepo.Create(changeLog); err != nil {
		zap.L().Error("记录域名同步日志失败", zap.Error(err))
	}

	result := map[string]interface{}{
		"provider_id":   providerID,
		"provider_name": provider.Name,
		"synced_count":  syncedCount,
		"error_count":   errorCount,
		"errors":        errors,
	}

	return result, nil
}

// SyncAllProviderDomains 同步所有提供商的域名
func (s *ProviderService) SyncAllProviderDomains(tenantID uint, actorID uint, clientIP string) (map[string]interface{}, error) {
	// 获取所有活跃的提供商
	filters := map[string]interface{}{
		"tenant_id": tenantID,
		"status":    "active",
	}

	providers, _, err := s.providerRepo.SearchWithFilters(filters, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("获取提供商列表失败: %w", err)
	}

	totalSynced := 0
	totalErrors := 0
	providerResults := make([]map[string]interface{}, 0)

	// 逐个同步每个提供商的域名
	for _, provider := range providers {
		result, err := s.SyncProviderDomains(provider.ID, actorID, clientIP)
		if err != nil {
			zap.L().Error("同步提供商域名失败",
				zap.String("provider", provider.Name),
				zap.Error(err))
			totalErrors++
			providerResults = append(providerResults, map[string]interface{}{
				"provider_id":   provider.ID,
				"provider_name": provider.Name,
				"synced_count":  0,
				"error_count":   1,
				"error":         err.Error(),
			})
		} else {
			syncedCount := result["synced_count"].(int)
			errorCount := result["error_count"].(int)
			totalSynced += syncedCount
			totalErrors += errorCount
			providerResults = append(providerResults, result)
		}
	}

	// 记录总体同步日志
	changeLog := &dns.ChangeLog{
		ResourceType: "provider",
		ResourceID:   0, // 全局操作
		Action:       "sync_all_domains",
		Description:  fmt.Sprintf("同步所有提供商域名，总计成功: %d, 失败: %d", totalSynced, totalErrors),
		Status:       "success",
		ClientIP:     clientIP,
		TenantID:     tenantID,
		ActorID:      actorID,
	}
	if err := s.changeLogRepo.Create(changeLog); err != nil {
		zap.L().Error("记录全局域名同步日志失败", zap.Error(err))
	}

	result := map[string]interface{}{
		"total_providers":  len(providers),
		"total_synced":     totalSynced,
		"total_errors":     totalErrors,
		"provider_results": providerResults,
	}

	return result, nil
}

// providerToJSON 将DNS提供商对象转换为JSON
func (s *ProviderService) providerToJSON(provider *dns.Provider) []byte {
	// 创建副本，隐藏敏感信息
	providerCopy := *provider
	providerCopy.CredentialsEnc = "[ENCRYPTED]"

	data, err := json.Marshal(&providerCopy)
	if err != nil {
		zap.L().Error("DNS提供商对象转JSON失败", zap.Error(err))
		return []byte("{}")
	}
	return data
}
