package dns

import (
	"api-server/internal/model/dns"
	dnsprovider "api-server/internal/provider/dns"
	repo "api-server/internal/repository/dns"
	"encoding/json"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProviderService DNS提供商业务逻辑服务
type ProviderService struct {
	providerRepo  *repo.ProviderRepository
	changeLogRepo *repo.ChangeLogRepository
	db            *gorm.DB
}

// NewProviderService 创建DNS提供商业务服务实例
func NewProviderService(
	providerRepo *repo.ProviderRepository,
	changeLogRepo *repo.ChangeLogRepository,
	db *gorm.DB,
) *ProviderService {
	return &ProviderService{
		providerRepo:  providerRepo,
		changeLogRepo: changeLogRepo,
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

	// TODO: 实现凭证加密逻辑
	// 这里应该使用AES加密凭证信息
	zap.L().Info("加密DNS提供商凭证", zap.String("provider", provider.Name))

	return nil
}

// decryptCredentials 解密凭证信息
func (s *ProviderService) decryptCredentials(provider *dns.Provider) (map[string]string, error) {
	if provider.CredentialsEnc == "" {
		return make(map[string]string), nil
	}

	// TODO: 实现凭证解密逻辑
	// 这里应该使用AES解密凭证信息
	zap.L().Info("解密DNS提供商凭证", zap.String("provider", provider.Name))

	// 临时返回空的凭证信息
	return make(map[string]string), nil
}

// SyncProviderDomains 同步单个提供商的域名
func (s *ProviderService) SyncProviderDomains(providerID uint, actorID uint, clientIP string) (map[string]interface{}, error) {
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
	driver, err := dnsprovider.CreateDriver(provider.Type, provider.Name, credentials)
	if err != nil {
		return nil, fmt.Errorf("创建DNS驱动失败: %w", err)
	}

	// 获取所有域名区域
	zones, err := driver.ListZones(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("获取域名列表失败: %w", err)
	}

	// 同步域名到数据库
	syncedCount := 0
	errorCount := 0
	var errors []string

	for _, zone := range zones {
		// TODO: 将域名保存到数据库
		// 这里需要调用DomainService来创建或更新域名
		zap.L().Info("同步域名",
			zap.String("provider", provider.Name),
			zap.String("zone", zone.Name))
		syncedCount++
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
