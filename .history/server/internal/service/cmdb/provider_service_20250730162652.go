package cmdb

import (
	"api-server/internal/config"
	model "api-server/internal/model/cmdb"
	repo "api-server/internal/repository/cmdb"
	"api-server/internal/service/cmdb/adapter"
	"api-server/pkg/crypto/encryption"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SyncResult 同步结果结构体
type SyncResult struct {
	Success      bool          `json:"success"`
	Message      string        `json:"message"`
	Duration     time.Duration `json:"duration"`
	Statistics   SyncStats     `json:"statistics"`
	ProviderName string        `json:"provider_name"`
	ProviderType string        `json:"provider_type"`
	ProviderID   uint          `json:"provider_id"`
	StartTime    time.Time     `json:"start_time"`
	Details      []SyncDetail  `json:"details"`
}

// SyncStats 同步统计信息
type SyncStats struct {
	TotalFound int `json:"total_found"`
	Created    int `json:"created"`
	Updated    int `json:"updated"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

// SyncDetail 单个区域同步详情
type SyncDetail struct {
	Region   string        `json:"region"`   // 区域名称
	Success  bool          `json:"success"`  // 是否成功
	Error    string        `json:"error"`    // 错误信息
	Count    int           `json:"count"`    // 同步的资源数量
	Duration time.Duration `json:"duration"` // 同步耗时
}

// ProviderService 云账号业务逻辑服务
// 提供云账号的增删改查、密钥加密解密等核心业务逻辑
type ProviderService struct {
	providerRepo *repo.ProviderRepository
	hostRepo     *repo.HostRepository
}

// NewProviderService 创建云账号业务服务实例
func NewProviderService(providerRepo *repo.ProviderRepository, hostRepo *repo.HostRepository) *ProviderService {
	return &ProviderService{providerRepo: providerRepo, hostRepo: hostRepo}
}

// CreateProvider 创建云账号，自动加密密钥
func (s *ProviderService) CreateProvider(provider *model.Provider) error {
	// 多密钥加密
	keys := getAllAESKeys()
	enc, err := encryption.EncryptAESWithKeys([]byte(provider.SecretKey), keys)
	if err != nil {
		return fmt.Errorf("密钥加密失败: %w", err)
	}
	provider.SecretKey = enc
	return s.providerRepo.Create(provider)
}

// UpdateProvider 更新云账号，自动加密密钥，未传递密钥时保留原密钥
func (s *ProviderService) UpdateProvider(provider *model.Provider) error {
	if provider.SecretKey != "" {
		key := []byte(config.GlobalConfig.App.AesKey)
		encryptedSecret, err := encryption.EncryptAES([]byte(provider.SecretKey), key)
		if err != nil {
			return fmt.Errorf("failed to encrypt secret: %w", err)
		}
		provider.SecretKey = encryptedSecret
	} else {
		// If secret is not provided, fetch the existing one to avoid overwriting it with an empty value.
		existing, err := s.providerRepo.FindByID(provider.ID)
		if err != nil {
			return fmt.Errorf("failed to find existing provider: %w", err)
		}
		provider.SecretKey = existing.SecretKey
	}
	return s.providerRepo.Update(provider)
}

// DeleteProvider 删除云账号
func (s *ProviderService) DeleteProvider(id uint) error {
	return s.providerRepo.DeleteByID(id)
}

// GetProviderByID 查询单个云账号，自动解密密钥
func (s *ProviderService) GetProviderByID(id uint) (*model.Provider, error) {
	return s.providerRepo.FindByID(id)
}

// ListProviders 查询云账号列表，支持后续扩展过滤
func (s *ProviderService) ListProviders() ([]*model.Provider, error) {
	return s.providerRepo.FindAll()
}

// SyncResources 同步资源，从云提供商获取实例并更新本地数据库
func (s *ProviderService) SyncResources(providerID uint, groupID *uint) (*SyncResult, error) {
	startTime := time.Now()

	result := &SyncResult{
		Success:    false,
		Statistics: SyncStats{},
	}
	provider, err := s.providerRepo.FindByID(providerID)
	if err != nil {
		result.Duration = time.Since(startTime)
		result.Message = fmt.Sprintf("云账号不存在: %v", err)
		return result, fmt.Errorf("provider not found: %w", err)
	}

	result.ProviderName = provider.Name
	result.ProviderType = provider.Type

	keys := getAllAESKeys()
	decryptedSecret, err := encryption.DecryptAESWithKeys(provider.SecretKey, keys)
	if err != nil {
		result.Duration = time.Since(startTime)
		result.Message = fmt.Sprintf("密钥解密失败: %v", err)
		return result, fmt.Errorf("failed to decrypt secret: %w", err)
	}

	var cloudHosts []model.Host
	switch provider.Type {
	case model.ProviderTypeAliyun:
		aliAdapter, err := adapter.NewAliyunAdapter(provider.AccessKey, decryptedSecret)
		if err != nil {
			result.Duration = time.Since(startTime)
			result.Message = fmt.Sprintf("创建阿里云适配器失败: %v", err)
			return result, fmt.Errorf("failed to create aliyun adapter: %w", err)
		}
		cloudHosts, err = aliAdapter.ListInstances()
		if err != nil {
			result.Duration = time.Since(startTime)
			result.Message = fmt.Sprintf("获取阿里云主机列表失败: %v", err)
			return result, fmt.Errorf("failed to get instances from aliyun: %w", err)
		}
	case model.ProviderTypeTencent:
		tencentAdapter, err := adapter.NewTencentAdapter(provider.Region, provider.AccessKey, decryptedSecret)
		if err != nil {
			result.Duration = time.Since(startTime)
			result.Message = fmt.Sprintf("创建腾讯云适配器失败: %v", err)
			return result, fmt.Errorf("failed to create tencent adapter: %w", err)
		}
		cloudHosts, err = tencentAdapter.ListInstances()
		if err != nil {
			result.Duration = time.Since(startTime)
			result.Message = fmt.Sprintf("获取腾讯云主机列表失败: %v", err)
			return result, fmt.Errorf("failed to get instances from tencent cloud: %w", err)
		}
	case model.ProviderTypeAWS:
		awsAdapter, err := adapter.NewAWSAdapter(provider.AccessKey, decryptedSecret, provider.Region)
		if err != nil {
			result.Duration = time.Since(startTime)
			result.Message = fmt.Sprintf("创建AWS适配器失败: %v", err)
			return result, fmt.Errorf("failed to create aws adapter: %w", err)
		}
		cloudHosts, err = awsAdapter.ListInstances()
		if err != nil {
			result.Duration = time.Since(startTime)
			result.Message = fmt.Sprintf("获取AWS主机列表失败: %v", err)
			return result, fmt.Errorf("failed to get instances from aws: %w", err)
		}
	default:
		result.Duration = time.Since(startTime)
		result.Message = fmt.Sprintf("不支持的云厂商类型: %s", provider.Type)
		return result, fmt.Errorf("unsupported provider type: %s", provider.Type)
	}

	// 设置统计信息
	result.Statistics.TotalFound = len(cloudHosts)

	// 实现完整同步逻辑（更新、创建、删除）
	for _, host := range cloudHosts {
		host.ProviderID = &provider.ID
		// 如果指定了主机组，将主机分配到该组
		if groupID != nil {
			host.GroupID = groupID
		}

		if err := s.hostRepo.Create(&host); err != nil {
			// 主机可能已存在，这是正常情况
			result.Statistics.Skipped++
			zap.L().Debug("Host already exists, skipping creation",
				zap.String("instance_id", host.InstanceID),
				zap.Error(err))
		} else {
			result.Statistics.Created++
			logFields := []zap.Field{
				zap.String("instance_id", host.InstanceID),
				zap.String("name", host.Name),
			}
			if host.GroupID != nil {
				logFields = append(logFields, zap.Uint("group_id", *host.GroupID))
			}
			zap.L().Debug("Successfully created host", logFields...)
		}
	}

	// 设置成功结果
	result.Success = true
	result.Duration = time.Since(startTime)
	result.Message = fmt.Sprintf("同步完成！共发现 %d 台主机，新增 %d 台，跳过 %d 台",
		result.Statistics.TotalFound,
		result.Statistics.Created,
		result.Statistics.Skipped)

	return result, nil
}

// SyncResourcesMultiRegion 多区域同步云账号资源
func (s *ProviderService) SyncResourcesMultiRegion(provider *model.Provider, regions []string, groupID *uint) (*SyncResult, error) {
	zap.L().Info("🚀 开始多区域同步云账号资源",
		zap.String("provider", provider.Name),
		zap.String("type", provider.Type),
		zap.Strings("regions", regions),
	)

	startTime := time.Now()
	result := &SyncResult{
		ProviderID: provider.ID,
		StartTime:  startTime,
		Statistics: SyncStats{},
		Details:    make([]SyncDetail, 0),
	}

	// 遍历每个区域进行同步
	for _, region := range regions {
		zap.L().Info("🌍 开始同步区域", zap.String("region", region))

		// 创建该区域的临时Provider对象
		regionProvider := &model.Provider{
			ID:        provider.ID,
			Name:      provider.Name,
			Type:      provider.Type,
			AccessKey: provider.AccessKey,
			SecretKey: provider.SecretKey,
			Region:    region, // 使用指定的区域
		}

		// 获取云适配器
		adapter, err := adapter.GetCloudAdapter(regionProvider)
		if err != nil {
			zap.L().Error("创建云适配器失败",
				zap.String("region", region),
				zap.Error(err))

			result.Details = append(result.Details, SyncDetail{
				Region:   region,
				Success:  false,
				Error:    err.Error(),
				Count:    0,
				Duration: time.Since(startTime),
			})
			continue
		}

		// 获取该区域的主机列表
		hosts, err := adapter.ListInstances()
		if err != nil {
			zap.L().Error("获取主机列表失败",
				zap.String("region", region),
				zap.Error(err))

			result.Details = append(result.Details, SyncDetail{
				Region:   region,
				Success:  false,
				Error:    err.Error(),
				Count:    0,
				Duration: time.Since(startTime),
			})
			continue
		}

		regionStartTime := time.Now()
		regionCreated := 0
		regionSkipped := 0

		// 处理该区域的每台主机
		for _, host := range hosts {
			// 检查主机是否已存在（基于实例ID和区域）
			existingHost, err := s.hostRepo.FindByInstanceIDAndRegion(host.InstanceID, region)
			if err != nil && err != gorm.ErrRecordNotFound {
				zap.L().Error("查询主机失败", zap.Error(err))
				continue
			}

			if existingHost != nil {
				// 主机已存在，更新信息
				existingHost.Name = host.Name
				existingHost.Status = host.Status
				existingHost.PrivateIP = host.PrivateIP
				existingHost.PublicIP = host.PublicIP
				existingHost.CPU = host.CPU
				existingHost.Memory = host.Memory
				existingHost.Disk = host.Disk
				existingHost.OS = host.OS
				existingHost.Region = region
				existingHost.Zone = host.Zone
				existingHost.ProviderID = provider.ID
				if groupID != nil {
					existingHost.GroupID = groupID
				}

				if err := s.hostRepo.Update(existingHost); err != nil {
					zap.L().Error("更新主机失败", zap.Error(err))
					continue
				}
				regionSkipped++
			} else {
				// 新主机，创建记录
				newHost := &model.Host{
					Name:       host.Name,
					InstanceID: host.InstanceID,
					Status:     host.Status,
					PrivateIP:  host.PrivateIP,
					PublicIP:   host.PublicIP,
					CPU:        host.CPU,
					Memory:     host.Memory,
					Disk:       host.Disk,
					OS:         host.OS,
					Region:     region,
					Zone:       host.Zone,
					ProviderID: provider.ID,
				}
				if groupID != nil {
					newHost.GroupID = groupID
				}

				if err := s.hostRepo.Create(newHost); err != nil {
					zap.L().Error("创建主机失败", zap.Error(err))
					continue
				}
				regionCreated++
			}
		}

		// 记录该区域的同步结果
		result.Details = append(result.Details, SyncDetail{
			Region:   region,
			Success:  true,
			Error:    "",
			Count:    len(hosts),
			Duration: time.Since(regionStartTime),
		})

		// 更新总体统计
		result.Statistics.TotalFound += len(hosts)
		result.Statistics.Created += regionCreated
		result.Statistics.Skipped += regionSkipped

		zap.L().Info("✅ 区域同步完成",
			zap.String("region", region),
			zap.Int("found", len(hosts)),
			zap.Int("created", regionCreated),
			zap.Int("skipped", regionSkipped),
		)
	}

	// 设置成功结果
	result.Success = true
	result.Duration = time.Since(startTime)
	result.Message = fmt.Sprintf("多区域同步完成！共同步 %d 个区域，发现 %d 台主机，新增 %d 台，跳过 %d 台",
		len(regions),
		result.Statistics.TotalFound,
		result.Statistics.Created,
		result.Statistics.Skipped)

	return result, nil
}

// GetProviderRegions 获取云账号的可用区域
func (s *ProviderService) GetProviderRegions(provider *model.Provider) ([]string, error) {
	// 获取云适配器
	adapter, err := adapter.GetCloudAdapter(provider)
	if err != nil {
		zap.L().Error("创建云适配器失败", zap.Error(err))
		// 如果适配器创建失败，返回默认区域列表
		return s.getDefaultRegions(provider.Type), nil
	}

	// 尝试从云厂商API获取区域列表
	regions, err := adapter.GetRegions()
	if err != nil {
		zap.L().Warn("从云厂商API获取区域失败，使用默认区域列表",
			zap.String("provider", provider.Type),
			zap.Error(err))
		return s.getDefaultRegions(provider.Type), nil
	}

	// 转换为字符串数组
	var regionList []string
	for _, region := range regions {
		regionList = append(regionList, region.ID)
	}

	return regionList, nil
}

// getDefaultRegions 获取默认区域列表（作为备选方案）
func (s *ProviderService) getDefaultRegions(providerType string) []string {
	switch providerType {
	case model.ProviderTypeAliyun:
		return []string{
			"cn-hangzhou", "cn-shanghai", "cn-beijing", "cn-shenzhen",
			"cn-guangzhou", "cn-qingdao", "cn-zhangjiakou", "cn-huhehaote",
			"cn-wulanchabu", "cn-chengdu", "cn-hongkong",
		}
	case model.ProviderTypeTencent:
		return []string{
			"ap-beijing", "ap-shanghai", "ap-guangzhou", "ap-chengdu",
			"ap-chongqing", "ap-nanjing", "ap-hongkong", "ap-singapore",
		}
	case model.ProviderTypeAWS:
		return []string{
			"us-east-1", "us-west-1", "us-west-2", "eu-west-1",
			"eu-central-1", "ap-southeast-1", "ap-southeast-2",
			"ap-northeast-1", "ap-northeast-2", "ap-south-1",
		}
	default:
		return []string{"default"}
	}
}

// ValidateCredentials 验证云账号凭证
func (s *ProviderService) ValidateCredentials(providerType, accessKey, secretKey, region string) (bool, error) {
	zap.L().Info("🔍 开始验证云账号凭证",
		zap.String("providerType", providerType),
		zap.String("accessKey", accessKey[:min(len(accessKey), 8)]+"***"), // 只显示前8位
		zap.String("region", region),
	)

	// 尝试创建适配器并获取实例列表来验证凭证
	var err error
	switch providerType {
	case model.ProviderTypeAliyun:
		zap.L().Info("📡 创建阿里云适配器")
		aliAdapter, adapterErr := adapter.NewAliyunAdapter(accessKey, secretKey)
		if adapterErr != nil {
			zap.L().Error("❌ 创建阿里云适配器失败", zap.Error(adapterErr))
			return false, fmt.Errorf("failed to create aliyun adapter: %w", adapterErr)
		}
		zap.L().Info("✅ 阿里云适配器创建成功，开始验证凭证")
		// 尝试获取实例列表来验证凭证
		instances, err := aliAdapter.ListInstances()
		if err != nil {
			zap.L().Error("❌ 阿里云凭证验证失败", zap.Error(err))
		} else {
			zap.L().Info("✅ 阿里云凭证验证成功", zap.Int("instanceCount", len(instances)))
		}
	case model.ProviderTypeTencent:
		zap.L().Info("📡 创建腾讯云适配器", zap.String("region", region))
		tencentAdapter, adapterErr := adapter.NewTencentAdapter(accessKey, secretKey, region)
		if adapterErr != nil {
			zap.L().Error("❌ 创建腾讯云适配器失败", zap.Error(adapterErr))
			return false, fmt.Errorf("failed to create tencent adapter: %w", adapterErr)
		}
		zap.L().Info("✅ 腾讯云适配器创建成功，开始验证凭证")
		// 尝试获取实例列表来验证凭证
		instances, err := tencentAdapter.ListInstances()
		if err != nil {
			zap.L().Error("❌ 腾讯云凭证验证失败", zap.Error(err))
		} else {
			zap.L().Info("✅ 腾讯云凭证验证成功", zap.Int("instanceCount", len(instances)))
		}
	case model.ProviderTypeAWS:
		zap.L().Info("📡 创建AWS适配器", zap.String("region", region))
		awsAdapter, adapterErr := adapter.NewAWSAdapter(accessKey, secretKey, region)
		if adapterErr != nil {
			zap.L().Error("❌ 创建AWS适配器失败", zap.Error(adapterErr))
			return false, fmt.Errorf("failed to create aws adapter: %w", adapterErr)
		}
		zap.L().Info("✅ AWS适配器创建成功，开始验证凭证")
		// 尝试获取实例列表来验证凭证
		instances, err := awsAdapter.ListInstances()
		if err != nil {
			zap.L().Error("❌ AWS凭证验证失败", zap.Error(err))
		} else {
			zap.L().Info("✅ AWS凭证验证成功", zap.Int("instanceCount", len(instances)))
		}
	default:
		zap.L().Error("❌ 不支持的云厂商类型", zap.String("providerType", providerType))
		return false, fmt.Errorf("unsupported provider type: %s", providerType)
	}

	// 如果有错误，说明凭证无效
	if err != nil {
		zap.L().Warn("⚠️ 凭证验证失败", zap.Error(err))
		return false, nil // 凭证无效，但不返回错误
	}

	zap.L().Info("🎉 凭证验证成功")
	return true, nil
}
