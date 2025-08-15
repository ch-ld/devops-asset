package dns

import (
	"api-server/internal/model/dns"
	"fmt"

	"gorm.io/gorm"
)

// RecordRepository DNS记录数据访问层
type RecordRepository struct {
	db *gorm.DB
}

// NewRecordRepository 创建DNS记录仓库实例
func NewRecordRepository(db *gorm.DB) *RecordRepository {
	return &RecordRepository{db: db}
}

// Create 创建DNS记录
func (r *RecordRepository) Create(record *dns.Record) error {
	return r.db.Create(record).Error
}

// Update 更新DNS记录
func (r *RecordRepository) Update(record *dns.Record) error {
	return r.db.Model(record).Updates(record).Error
}

// Delete 删除DNS记录
func (r *RecordRepository) Delete(id uint) error {
	return r.db.Delete(&dns.Record{}, id).Error
}

// FindByID 根据ID查找DNS记录
func (r *RecordRepository) FindByID(id uint) (*dns.Record, error) {
	var record dns.Record
	err := r.db.Preload("Domain").Preload("Provider").First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// FindByDomainID 根据域名ID查找DNS记录
func (r *RecordRepository) FindByDomainID(domainID uint) ([]*dns.Record, error) {
	var records []*dns.Record
	err := r.db.Preload("Domain").Preload("Provider").Where("domain_id = ?", domainID).Find(&records).Error
	return records, err
}

// FindByProviderID 根据提供商ID查找DNS记录
func (r *RecordRepository) FindByProviderID(providerID uint) ([]*dns.Record, error) {
	var records []*dns.Record
	err := r.db.Preload("Domain").Preload("Provider").Where("provider_id = ?", providerID).Find(&records).Error
	return records, err
}

// FindByType 根据记录类型查找DNS记录
func (r *RecordRepository) FindByType(recordType string) ([]*dns.Record, error) {
	var records []*dns.Record
	err := r.db.Preload("Domain").Preload("Provider").Where("type = ?", recordType).Find(&records).Error
	return records, err
}

// FindByStatus 根据状态查找DNS记录
func (r *RecordRepository) FindByStatus(status string) ([]*dns.Record, error) {
	var records []*dns.Record
	err := r.db.Preload("Domain").Preload("Provider").Where("status = ?", status).Find(&records).Error
	return records, err
}

// FindBySyncStatus 根据同步状态查找DNS记录
func (r *RecordRepository) FindBySyncStatus(syncStatus string) ([]*dns.Record, error) {
	var records []*dns.Record
	err := r.db.Preload("Domain").Preload("Provider").Where("sync_status = ?", syncStatus).Find(&records).Error
	return records, err
}

// FindByDomainAndName 根据域名和记录名查找DNS记录
func (r *RecordRepository) FindByDomainAndName(domainID uint, name string) ([]*dns.Record, error) {
	var records []*dns.Record
	err := r.db.Preload("Domain").Preload("Provider").
		Where("domain_id = ? AND name = ?", domainID, name).Find(&records).Error
	return records, err
}

// FindByDomainAndType 根据域名和记录类型查找DNS记录
func (r *RecordRepository) FindByDomainAndType(domainID uint, recordType string) ([]*dns.Record, error) {
	var records []*dns.Record
	err := r.db.Preload("Domain").Preload("Provider").
		Where("domain_id = ? AND type = ?", domainID, recordType).Find(&records).Error
	return records, err
}

// FindByCloudRecordID 根据云端记录ID查找DNS记录
func (r *RecordRepository) FindByCloudRecordID(cloudRecordID string) (*dns.Record, error) {
	var record dns.Record
	err := r.db.Preload("Domain").Preload("Provider").
		Where("cloud_record_id = ?", cloudRecordID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// FindByTenantID 根据租户ID查找DNS记录
func (r *RecordRepository) FindByTenantID(tenantID uint) ([]*dns.Record, error) {
	var records []*dns.Record
	err := r.db.Preload("Domain").Preload("Provider").Where("tenant_id = ?", tenantID).Find(&records).Error
	return records, err
}

// Search 搜索DNS记录
func (r *RecordRepository) Search(keyword string, limit, offset int) ([]*dns.Record, int64, error) {
	var records []*dns.Record
	var total int64

	query := r.db.Preload("Domain").Preload("Provider").Model(&dns.Record{})

	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR value LIKE ? OR remark LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&records).Error
	return records, total, err
}

// SearchWithFilters 带过滤条件的搜索
func (r *RecordRepository) SearchWithFilters(filters map[string]interface{}, limit, offset int) ([]*dns.Record, int64, error) {
	var records []*dns.Record
	var total int64

	query := r.db.Preload("Domain").Preload("Provider").Model(&dns.Record{})

	// 应用过滤条件
	for key, value := range filters {
		if value != nil && value != "" {
			switch key {
			case "keyword":
				searchPattern := "%" + fmt.Sprintf("%v", value) + "%"
				query = query.Where("name LIKE ? OR value LIKE ? OR remark LIKE ?",
					searchPattern, searchPattern, searchPattern)
			case "domain_id":
				query = query.Where("domain_id = ?", value)
			case "provider_id":
				query = query.Where("provider_id = ?", value)
			case "type":
				query = query.Where("type = ?", value)
			case "status":
				query = query.Where("status = ?", value)
			case "sync_status":
				query = query.Where("sync_status = ?", value)
			case "tenant_id":
				query = query.Where("tenant_id = ?", value)
			}
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&records).Error
	return records, total, err
}

// CountByType 统计各类型记录数量
func (r *RecordRepository) CountByType() (map[string]int64, error) {
	var results []struct {
		Type  string
		Count int64
	}

	err := r.db.Model(&dns.Record{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.Type] = result.Count
	}

	return counts, nil
}

// CountByStatus 统计各状态记录数量
func (r *RecordRepository) CountByStatus() (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	err := r.db.Model(&dns.Record{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.Status] = result.Count
	}

	return counts, nil
}

// CountByDomainIDs 批量统计各域名的记录数量（按租户过滤）
func (r *RecordRepository) CountByDomainIDs(domainIDs []uint, tenantID uint) (map[uint]int64, error) {
	result := make(map[uint]int64)
	if len(domainIDs) == 0 {
		return result, nil
	}

	type row struct {
		DomainID uint  `json:"domain_id"`
		Count    int64 `json:"count"`
	}
	var rows []row

	err := r.db.Model(&dns.Record{}).
		Select("domain_id, COUNT(*) as count").
		Where("domain_id IN ?", domainIDs).
		Where("tenant_id = ?", tenantID).
		Group("domain_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		result[r.DomainID] = r.Count
	}
	return result, nil
}


// BatchCreate 批量创建DNS记录
func (r *RecordRepository) BatchCreate(records []*dns.Record) error {
	return r.db.Create(&records).Error
}

// BatchUpdate 批量更新DNS记录
func (r *RecordRepository) BatchUpdate(records []*dns.Record) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, record := range records {
			if err := tx.Model(record).Updates(record).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchDelete 批量删除DNS记录
func (r *RecordRepository) BatchDelete(ids []uint) error {
	return r.db.Delete(&dns.Record{}, ids).Error
}

// BatchUpdateSyncStatus 批量更新同步状态
func (r *RecordRepository) BatchUpdateSyncStatus(ids []uint, syncStatus string) error {
	return r.db.Model(&dns.Record{}).Where("id IN ?", ids).Update("sync_status", syncStatus).Error
}

// ExistsByDomainAndNameAndType 检查记录是否存在
func (r *RecordRepository) ExistsByDomainAndNameAndType(domainID uint, name, recordType string) (bool, error) {
	var count int64
	err := r.db.Model(&dns.Record{}).
		Where("domain_id = ? AND name = ? AND type = ?", domainID, name, recordType).
		Count(&count).Error
	return count > 0, err
}

// CountByTenantID 统计租户记录数量
func (r *RecordRepository) CountByTenantID(tenantID uint) (int64, error) {
	var count int64
	err := r.db.Model(&dns.Record{}).Where("tenant_id = ?", tenantID).Count(&count).Error
	return count, err
}

// UpdateSyncStatus 更新同步状态
func (r *RecordRepository) UpdateSyncStatus(id uint, syncStatus string) error {
	return r.db.Model(&dns.Record{}).Where("id = ?", id).Update("sync_status", syncStatus).Error
}

// UpdateCloudRecordID 更新云端记录ID
func (r *RecordRepository) UpdateCloudRecordID(id uint, cloudRecordID string) error {
	return r.db.Model(&dns.Record{}).Where("id = ?", id).Update("cloud_record_id", cloudRecordID).Error
}
