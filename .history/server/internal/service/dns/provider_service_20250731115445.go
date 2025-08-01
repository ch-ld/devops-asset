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
