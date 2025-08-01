package dns

import (
	"api-server/internal/model/dns"
	repo "api-server/internal/repository/dns"
	"encoding/json"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RecordService DNS记录业务逻辑服务
type RecordService struct {
	recordRepo    *repo.RecordRepository
	domainRepo    *repo.DomainRepository
	providerRepo  *repo.ProviderRepository
	changeLogRepo *repo.ChangeLogRepository
	db            *gorm.DB
}

// NewRecordService 创建DNS记录业务服务实例
func NewRecordService(
	recordRepo *repo.RecordRepository,
	domainRepo *repo.DomainRepository,
	providerRepo *repo.ProviderRepository,
	changeLogRepo *repo.ChangeLogRepository,
	db *gorm.DB,
) *RecordService {
	return &RecordService{
		recordRepo:    recordRepo,
		domainRepo:    domainRepo,
		providerRepo:  providerRepo,
		changeLogRepo: changeLogRepo,
		db:            db,
	}
}

// CreateRecord 创建DNS记录
func (s *RecordService) CreateRecord(record *dns.Record, actorID uint, clientIP string) error {
	// 验证域名是否存在
	domain, err := s.domainRepo.FindByID(record.DomainID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("指定的域名不存在")
		}
		return fmt.Errorf("验证域名失败: %w", err)
	}

	// 验证DNS提供商是否存在
	provider, err := s.providerRepo.FindByID(record.ProviderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("指定的DNS提供商不存在")
		}
		return fmt.Errorf("验证DNS提供商失败: %w", err)
	}

	// 验证记录格式
	if err := s.validateRecord(record); err != nil {
		return fmt.Errorf("DNS记录验证失败: %w", err)
	}

	// 检查记录是否已存在
	exists, err := s.recordRepo.ExistsByDomainAndNameAndType(record.DomainID, record.Name, record.Type)
	if err != nil {
		return fmt.Errorf("检查DNS记录是否存在失败: %w", err)
	}
	if exists {
		return errors.New("相同名称和类型的DNS记录已存在")
	}

	// 设置创建信息
	record.TenantID = domain.TenantID
	record.CreatedBy = actorID
	record.UpdatedBy = actorID
	record.SyncStatus = "pending"

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 创建记录
		if err := s.recordRepo.Create(record); err != nil {
			return fmt.Errorf("创建DNS记录失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "record",
			ResourceID:   record.ID,
			Action:       "create",
			Description:  fmt.Sprintf("创建DNS记录: %s.%s (%s)", record.Name, domain.Name, record.Type),
			NewData:      s.recordToJSON(record),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     record.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录DNS记录创建日志失败", zap.Error(err))
		}

		// TODO: 异步同步到DNS提供商
		go s.syncRecordToProvider(record, provider, "create")

		return nil
	})
}

// UpdateRecord 更新DNS记录
func (s *RecordService) UpdateRecord(record *dns.Record, actorID uint, clientIP string) error {
	// 获取原始数据
	oldRecord, err := s.recordRepo.FindByID(record.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("DNS记录不存在")
		}
		return fmt.Errorf("获取DNS记录信息失败: %w", err)
	}

	// 验证域名是否存在
	domain, err := s.domainRepo.FindByID(record.DomainID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("指定的域名不存在")
		}
		return fmt.Errorf("验证域名失败: %w", err)
	}

	// 验证DNS提供商是否存在
	provider, err := s.providerRepo.FindByID(record.ProviderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("指定的DNS提供商不存在")
		}
		return fmt.Errorf("验证DNS提供商失败: %w", err)
	}

	// 验证记录格式
	if err := s.validateRecord(record); err != nil {
		return fmt.Errorf("DNS记录验证失败: %w", err)
	}

	// 如果名称或类型发生变化，检查是否冲突
	if oldRecord.Name != record.Name || oldRecord.Type != record.Type {
		exists, err := s.recordRepo.ExistsByDomainAndNameAndType(record.DomainID, record.Name, record.Type)
		if err != nil {
			return fmt.Errorf("检查DNS记录是否存在失败: %w", err)
		}
		if exists {
			return errors.New("相同名称和类型的DNS记录已存在")
		}
	}

	// 设置更新信息
	record.UpdatedBy = actorID
	record.SyncStatus = "pending"

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 更新记录
		if err := s.recordRepo.Update(record); err != nil {
			return fmt.Errorf("更新DNS记录失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "record",
			ResourceID:   record.ID,
			Action:       "update",
			Description:  fmt.Sprintf("更新DNS记录: %s.%s (%s)", record.Name, domain.Name, record.Type),
			OldData:      s.recordToJSON(oldRecord),
			NewData:      s.recordToJSON(record),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     record.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录DNS记录更新日志失败", zap.Error(err))
		}

		// TODO: 异步同步到DNS提供商
		go s.syncRecordToProvider(record, provider, "update")

		return nil
	})
}

// DeleteRecord 删除DNS记录
func (s *RecordService) DeleteRecord(id uint, actorID uint, clientIP string) error {
	// 获取记录信息
	record, err := s.recordRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("DNS记录不存在")
		}
		return fmt.Errorf("获取DNS记录信息失败: %w", err)
	}

	// 获取域名信息
	domain, err := s.domainRepo.FindByID(record.DomainID)
	if err != nil {
		return fmt.Errorf("获取域名信息失败: %w", err)
	}

	// 获取提供商信息
	provider, err := s.providerRepo.FindByID(record.ProviderID)
	if err != nil {
		return fmt.Errorf("获取DNS提供商信息失败: %w", err)
	}

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除记录
		if err := s.recordRepo.Delete(id); err != nil {
			return fmt.Errorf("删除DNS记录失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "record",
			ResourceID:   id,
			Action:       "delete",
			Description:  fmt.Sprintf("删除DNS记录: %s.%s (%s)", record.Name, domain.Name, record.Type),
			OldData:      s.recordToJSON(record),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     record.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录DNS记录删除日志失败", zap.Error(err))
		}

		// TODO: 异步从DNS提供商删除
		go s.syncRecordToProvider(record, provider, "delete")

		return nil
	})
}

// GetRecord 获取DNS记录详情
func (s *RecordService) GetRecord(id uint) (*dns.Record, error) {
	record, err := s.recordRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("DNS记录不存在")
		}
		return nil, fmt.Errorf("获取DNS记录信息失败: %w", err)
	}
	return record, nil
}

// ListRecords 获取DNS记录列表
func (s *RecordService) ListRecords(tenantID uint, filters map[string]interface{}, limit, offset int) ([]*dns.Record, int64, error) {
	// 添加租户过滤
	if filters == nil {
		filters = make(map[string]interface{})
	}
	filters["tenant_id"] = tenantID

	records, total, err := s.recordRepo.SearchWithFilters(filters, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("获取DNS记录列表失败: %w", err)
	}

	return records, total, nil
}

// GetRecordsByDomain 获取指定域名的DNS记录
func (s *RecordService) GetRecordsByDomain(domainID uint) ([]*dns.Record, error) {
	records, err := s.recordRepo.FindByDomainID(domainID)
	if err != nil {
		return nil, fmt.Errorf("获取域名DNS记录失败: %w", err)
	}
	return records, nil
}

// GetRecordStatistics 获取DNS记录统计信息
func (s *RecordService) GetRecordStatistics(tenantID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总数统计
	total, err := s.recordRepo.CountByTenantID(tenantID)
	if err != nil {
		return nil, fmt.Errorf("统计DNS记录总数失败: %w", err)
	}
	stats["total"] = total

	// 类型统计
	typeCounts, err := s.recordRepo.CountByType()
	if err != nil {
		return nil, fmt.Errorf("统计DNS记录类型失败: %w", err)
	}
	stats["by_type"] = typeCounts

	// 状态统计
	statusCounts, err := s.recordRepo.CountByStatus()
	if err != nil {
		return nil, fmt.Errorf("统计DNS记录状态失败: %w", err)
	}
	stats["by_status"] = statusCounts

	return stats, nil
}

// validateRecord 验证DNS记录
func (s *RecordService) validateRecord(record *dns.Record) error {
	if record.Name == "" {
		return errors.New("记录名称不能为空")
	}

	if record.Type == "" {
		return errors.New("记录类型不能为空")
	}

	if record.Value == "" {
		return errors.New("记录值不能为空")
	}

	// 验证记录类型
	validTypes := []string{"A", "AAAA", "CNAME", "MX", "TXT", "SRV", "NS", "PTR", "CAA"}
	isValidType := false
	for _, validType := range validTypes {
		if record.Type == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("不支持的记录类型: %s", record.Type)
	}

	// 验证TTL
	if record.TTL < 60 || record.TTL > 86400 {
		return errors.New("TTL值必须在60-86400秒之间")
	}

	// TODO: 根据记录类型验证记录值格式

	return nil
}

// syncRecordToProvider 同步记录到DNS提供商
func (s *RecordService) syncRecordToProvider(record *dns.Record, provider *dns.Provider, action string) {
	// TODO: 实现DNS提供商同步逻辑
	// 这里应该调用provider的Driver来同步记录
	zap.L().Info("DNS记录同步到提供商",
		zap.Uint("record_id", record.ID),
		zap.String("provider", provider.Name),
		zap.String("action", action),
	)

	// 更新同步状态
	var syncStatus string
	// 这里应该根据实际同步结果设置状态
	syncStatus = "synced" // 假设同步成功

	if err := s.recordRepo.UpdateSyncStatus(record.ID, syncStatus); err != nil {
		zap.L().Error("更新DNS记录同步状态失败", zap.Error(err))
	}
}

// recordToJSON 将DNS记录对象转换为JSON
func (s *RecordService) recordToJSON(record *dns.Record) []byte {
	data, err := encryption.ToJSON(record)
	if err != nil {
		zap.L().Error("DNS记录对象转JSON失败", zap.Error(err))
		return []byte("{}")
	}
	return data
}
