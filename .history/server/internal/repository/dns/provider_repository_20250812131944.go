package dns

import (
	"api-server/internal/model/dns"
	"fmt"

	"gorm.io/gorm"
)

// ProviderRepository DNS提供商数据访问层
type ProviderRepository struct {
	db *gorm.DB
}

// NewProviderRepository 创建DNS提供商仓库实例
func NewProviderRepository(db *gorm.DB) *ProviderRepository {
	return &ProviderRepository{db: db}
}

// Create 创建DNS提供商
func (r *ProviderRepository) Create(provider *dns.Provider) error {
	return r.db.Create(provider).Error
}

// Update 更新DNS提供商
func (r *ProviderRepository) Update(provider *dns.Provider) error {
	return r.db.Model(provider).Updates(provider).Error
}

// Delete 删除DNS提供商
func (r *ProviderRepository) Delete(id uint) error {
	return r.db.Delete(&dns.Provider{}, id).Error
}

// FindByID 根据ID查找DNS提供商
func (r *ProviderRepository) FindByID(id uint) (*dns.Provider, error) {
	var provider dns.Provider
	err := r.db.First(&provider, id).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// FindByName 根据名称查找DNS提供商
func (r *ProviderRepository) FindByName(name string) (*dns.Provider, error) {
	var provider dns.Provider
	err := r.db.Where("name = ?", name).First(&provider).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// FindByType 根据类型查找DNS提供商
func (r *ProviderRepository) FindByType(providerType string) ([]*dns.Provider, error) {
	var providers []*dns.Provider
	err := r.db.Where("type = ?", providerType).Find(&providers).Error
	return providers, err
}

// FindByStatus 根据状态查找DNS提供商
func (r *ProviderRepository) FindByStatus(status string) ([]*dns.Provider, error) {
	var providers []*dns.Provider
	err := r.db.Where("status = ?", status).Find(&providers).Error
	return providers, err
}

// FindDefault 查找默认DNS提供商
func (r *ProviderRepository) FindDefault() (*dns.Provider, error) {
	var provider dns.Provider
	err := r.db.Where("is_default = ? AND status = ?", true, "active").First(&provider).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// FindByTenantID 根据租户ID查找DNS提供商
func (r *ProviderRepository) FindByTenantID(tenantID uint) ([]*dns.Provider, error) {
	var providers []*dns.Provider
	err := r.db.Where("tenant_id = ?", tenantID).Find(&providers).Error
	return providers, err
}

// FindAll 查找所有DNS提供商
func (r *ProviderRepository) FindAll() ([]*dns.Provider, error) {
	var providers []*dns.Provider
	err := r.db.Find(&providers).Error
	return providers, err
}

// Search 搜索DNS提供商
func (r *ProviderRepository) Search(keyword string, limit, offset int) ([]*dns.Provider, int64, error) {
	var providers []*dns.Provider
	var total int64

	query := r.db.Model(&dns.Provider{})

	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR type LIKE ? OR remark LIKE ?", 
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

	err := query.Find(&providers).Error
	return providers, total, err
}

// SearchWithFilters 带过滤条件的搜索
func (r *ProviderRepository) SearchWithFilters(filters map[string]interface{}, limit, offset int) ([]*dns.Provider, int64, error) {
	var providers []*dns.Provider
	var total int64

	query := r.db.Model(&dns.Provider{})

	// 应用过滤条件
	for key, value := range filters {
		if value != nil && value != "" {
			switch key {
			case "keyword":
				searchPattern := "%" + fmt.Sprintf("%v", value) + "%"
				query = query.Where("name LIKE ? OR type LIKE ? OR remark LIKE ?", 
					searchPattern, searchPattern, searchPattern)
			case "type":
				query = query.Where("type = ?", value)
			case "status":
				query = query.Where("status = ?", value)
			case "is_default":
				query = query.Where("is_default = ?", value)
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

	err := query.Find(&providers).Error
	return providers, total, err
}

// CountByType 统计各类型提供商数量
func (r *ProviderRepository) CountByType() (map[string]int64, error) {
	var results []struct {
		Type  string
		Count int64
	}

	err := r.db.Model(&dns.Provider{}).
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

// CountByStatus 统计各状态提供商数量
func (r *ProviderRepository) CountByStatus() (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	err := r.db.Model(&dns.Provider{}).
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

// CountByStatusAndTenant 统计各状态提供商数量（按租户过滤）
func (r *ProviderRepository) CountByStatusAndTenant(tenantID uint) (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	err := r.db.Model(&dns.Provider{}).
		Select("status, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
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

// CountByTypeAndTenant 统计各类型提供商数量（按租户过滤）
func (r *ProviderRepository) CountByTypeAndTenant(tenantID uint) (map[string]int64, error) {
	var results []struct {
		Type  string
		Count int64
	}

	err := r.db.Model(&dns.Provider{}).
		Select("type, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
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

// BatchCreate 批量创建DNS提供商
func (r *ProviderRepository) BatchCreate(providers []*dns.Provider) error {
	return r.db.Create(&providers).Error
}

// BatchUpdate 批量更新DNS提供商
func (r *ProviderRepository) BatchUpdate(providers []*dns.Provider) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, provider := range providers {
			if err := tx.Model(provider).Updates(provider).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchDelete 批量删除DNS提供商
func (r *ProviderRepository) BatchDelete(ids []uint) error {
	return r.db.Delete(&dns.Provider{}, ids).Error
}

// ExistsByName 检查提供商名称是否存在
func (r *ProviderRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&dns.Provider{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// CountByTenantID 统计租户提供商数量
func (r *ProviderRepository) CountByTenantID(tenantID uint) (int64, error) {
	var count int64
	err := r.db.Model(&dns.Provider{}).Where("tenant_id = ?", tenantID).Count(&count).Error
	return count, err
}

// UpdateStatus 更新提供商状态
func (r *ProviderRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&dns.Provider{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateTestResult 更新测试结果
func (r *ProviderRepository) UpdateTestResult(id uint, testResult, errorMessage string) error {
	return r.db.Model(&dns.Provider{}).Where("id = ?", id).Updates(map[string]interface{}{
		"test_result":   testResult,
		"error_message": errorMessage,
		"last_test_at":  gorm.Expr("NOW()"),
	}).Error
}

// SetDefault 设置默认提供商
func (r *ProviderRepository) SetDefault(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先取消所有默认设置
		if err := tx.Model(&dns.Provider{}).Update("is_default", false).Error; err != nil {
			return err
		}
		// 设置指定提供商为默认
		return tx.Model(&dns.Provider{}).Where("id = ?", id).Update("is_default", true).Error
	})
}

// UpdateCredentials 更新凭证信息
func (r *ProviderRepository) UpdateCredentials(id uint, credentialsEnc string) error {
	return r.db.Model(&dns.Provider{}).Where("id = ?", id).Update("credentials_enc", credentialsEnc).Error
}
