package cmdb

import (
	"fmt"
	"time"

	"api-server/internal/model"
	"api-server/internal/service/cmdb/adapter"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RegionService 区域服务
type RegionService struct {
	db *gorm.DB
}

// NewRegionService 创建区域服务
func NewRegionService(db *gorm.DB) *RegionService {
	return &RegionService{db: db}
}

// GetRegionsByProvider 获取指定云厂商的区域列表
func (s *RegionService) GetRegionsByProvider(providerType string) ([]string, error) {
	var regions []model.Region
	err := s.db.Where("provider_type = ? AND status = ?", providerType, "active").
		Order("is_default DESC, region_id ASC").
		Find(&regions).Error
	if err != nil {
		return nil, fmt.Errorf("查询区域失败: %w", err)
	}

	// 如果数据库中没有区域数据，说明系统未正确初始化，返回默认区域作为备选
	if len(regions) == 0 {
		zap.L().Warn("数据库中未找到区域数据，返回默认区域", zap.String("provider", providerType))
		return s.getDefaultRegions(providerType), nil
	}

	// 转换为字符串数组
	var regionList []string
	for _, region := range regions {
		regionList = append(regionList, region.RegionID)
	}

	return regionList, nil
}

// SyncRegionsFromProvider 从云厂商同步区域数据
func (s *RegionService) SyncRegionsFromProvider(provider *model.Provider) error {
	zap.L().Info("开始同步区域数据", zap.String("provider", provider.Name), zap.String("type", provider.Type))

	// 获取云适配器
	cloudAdapter, err := adapter.GetCloudAdapter(provider)
	if err != nil {
		zap.L().Error("创建云适配器失败", zap.Error(err))
		// 如果适配器创建失败，使用默认区域
		return s.initDefaultRegions(provider.Type)
	}

	// 尝试从云厂商API获取区域列表
	regions, err := cloudAdapter.GetRegions()
	if err != nil {
		zap.L().Warn("从云厂商API获取区域失败，使用默认区域",
			zap.String("provider", provider.Type),
			zap.Error(err))
		return s.initDefaultRegions(provider.Type)
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 清除旧的区域数据
	if err := tx.Where("provider_type = ?", provider.Type).Delete(&model.Region{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("清除旧区域数据失败: %w", err)
	}

	// 插入新的区域数据
	now := time.Now()
	for _, region := range regions {
		regionModel := &model.Region{
			ProviderType: provider.Type,
			RegionID:     region.ID,
			RegionName:   region.Name,
			Status:       "active",
			IsDefault:    false, // 从API获取的区域默认不设为默认
			SyncTime:     now,
		}

		if err := tx.Create(regionModel).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("插入区域数据失败: %w", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	zap.L().Info("区域数据同步完成",
		zap.String("provider", provider.Type),
		zap.Int("count", len(regions)))

	return nil
}

// initDefaultRegions 初始化默认区域数据
func (s *RegionService) initDefaultRegions(providerType string) error {
	// 检查是否已有区域数据
	var count int64
	s.db.Model(&model.Region{}).Where("provider_type = ?", providerType).Count(&count)
	if count > 0 {
		return nil // 已有数据，不需要初始化
	}

	defaultRegions := s.getDefaultRegions(providerType)
	if len(defaultRegions) == 0 {
		return fmt.Errorf("不支持的云厂商类型: %s", providerType)
	}

	// 插入默认区域数据
	now := time.Now()
	for i, regionID := range defaultRegions {
		regionModel := &model.Region{
			ProviderType: providerType,
			RegionID:     regionID,
			RegionName:   regionID, // 默认区域名称与ID相同
			Status:       "active",
			IsDefault:    i == 0, // 第一个区域设为默认
			SyncTime:     now,
		}

		if err := s.db.Create(regionModel).Error; err != nil {
			return fmt.Errorf("插入默认区域数据失败: %w", err)
		}
	}

	zap.L().Info("默认区域数据初始化完成",
		zap.String("provider", providerType),
		zap.Int("count", len(defaultRegions)))

	return nil
}

// getDefaultRegions 获取默认区域列表
func (s *RegionService) getDefaultRegions(providerType string) []string {
	switch providerType {
	case "alicloud":
		return []string{
			"cn-hangzhou",
			"cn-shanghai",
			"cn-beijing",
			"cn-shenzhen",
			"cn-guangzhou",
			"cn-chengdu",
			"cn-qingdao",
			"cn-zhangjiakou",
		}
	case "tencentcloud":
		return []string{
			"ap-beijing",
			"ap-shanghai",
			"ap-guangzhou",
			"ap-chengdu",
			"ap-chongqing",
			"ap-nanjing",
			"ap-shenzhen-fsi",
		}
	case "aws":
		return []string{
			"us-east-1",
			"us-east-2",
			"us-west-1",
			"us-west-2",
			"ca-central-1",
			"eu-central-1",
			"eu-west-1",
			"eu-west-2",
			"eu-west-3",
			"eu-north-1",
			"eu-south-1",
			"ap-northeast-1",
			"ap-northeast-2",
			"ap-northeast-3",
			"ap-southeast-1",
			"ap-southeast-2",
			"ap-southeast-3",
			"ap-south-1",
			"ap-east-1",
			"me-south-1",
			"af-south-1",
			"sa-east-1",
		}
	case "huaweicloud":
		return []string{
			"cn-north-1",
			"cn-north-4",
			"cn-east-2",
			"cn-east-3",
			"cn-south-1",
		}
	case "baiduyun":
		return []string{
			"bj",
			"gz",
			"su",
			"hkg",
		}
	default:
		return []string{}
	}
}

// SyncAllProviderRegions 同步所有云账号的区域数据
func (s *RegionService) SyncAllProviderRegions() error {
	// 获取所有云账号
	var providers []model.Provider
	if err := s.db.Find(&providers).Error; err != nil {
		return fmt.Errorf("获取云账号列表失败: %w", err)
	}

	// 按云厂商类型分组，避免重复同步
	providerTypes := make(map[string]*model.Provider)
	for _, provider := range providers {
		if _, exists := providerTypes[provider.Type]; !exists {
			providerTypes[provider.Type] = &provider
		}
	}

	// 同步每种云厂商的区域数据
	for providerType, provider := range providerTypes {
		if err := s.SyncRegionsFromProvider(provider); err != nil {
			zap.L().Error("同步区域数据失败",
				zap.String("provider_type", providerType),
				zap.Error(err))
			// 继续同步其他云厂商，不中断整个过程
		}
	}

	return nil
}

// GetRegionSyncStatus 获取区域同步状态
func (s *RegionService) GetRegionSyncStatus() (map[string]time.Time, error) {
	var regions []model.Region
	err := s.db.Select("provider_type, MAX(sync_time) as sync_time").
		Group("provider_type").
		Find(&regions).Error
	if err != nil {
		return nil, fmt.Errorf("查询同步状态失败: %w", err)
	}

	status := make(map[string]time.Time)
	for _, region := range regions {
		status[region.ProviderType] = region.SyncTime
	}

	return status, nil
}
