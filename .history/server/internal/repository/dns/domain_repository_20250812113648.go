package dns

import (
	"api-server/internal/model/dns"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DomainRepository 域名数据访问层
type DomainRepository struct {
	db *gorm.DB
}

// NewDomainRepository 创建域名仓库实例
func NewDomainRepository(db *gorm.DB) *DomainRepository {
	return &DomainRepository{db: db}
}

// Create 创建域名
func (r *DomainRepository) Create(domain *dns.Domain) error {
	return r.db.Create(domain).Error
}

// Update 更新域名
func (r *DomainRepository) Update(domain *dns.Domain) error {
	return r.db.Model(domain).Updates(domain).Error
}

// Delete 删除域名
func (r *DomainRepository) Delete(id uint) error {
	return r.db.Delete(&dns.Domain{}, id).Error
}

// FindByID 根据ID查找域名
func (r *DomainRepository) FindByID(id uint) (*dns.Domain, error) {
	var domain dns.Domain
	err := r.db.Preload("Group").Preload("Tags").First(&domain, id).Error
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

// FindByName 根据域名查找
func (r *DomainRepository) FindByName(name string) (*dns.Domain, error) {
	var domain dns.Domain
	err := r.db.Preload("Group").Preload("Tags").Where("name = ?", name).First(&domain).Error
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

// FindAll 查找所有域名
func (r *DomainRepository) FindAll() ([]*dns.Domain, error) {
	var domains []*dns.Domain
	err := r.db.Preload("Group").Preload("Tags").Find(&domains).Error
	return domains, err
}

// FindByTenantID 根据租户ID查找域名
func (r *DomainRepository) FindByTenantID(tenantID uint) ([]*dns.Domain, error) {
	var domains []*dns.Domain
	err := r.db.Preload("Group").Preload("Tags").Where("tenant_id = ?", tenantID).Find(&domains).Error
	return domains, err
}

// FindByGroupID 根据分组ID查找域名
func (r *DomainRepository) FindByGroupID(groupID uint) ([]*dns.Domain, error) {
	var domains []*dns.Domain
	err := r.db.Preload("Group").Preload("Tags").Where("group_id = ?", groupID).Find(&domains).Error
	return domains, err
}

// FindByStatus 根据状态查找域名
func (r *DomainRepository) FindByStatus(status string) ([]*dns.Domain, error) {
	var domains []*dns.Domain
	err := r.db.Preload("Group").Preload("Tags").Where("status = ?", status).Find(&domains).Error
	return domains, err
}

// FindExpiring 查找即将过期的域名
func (r *DomainRepository) FindExpiring(days int) ([]*dns.Domain, error) {
	var domains []*dns.Domain
	err := r.db.Preload("Group").Preload("Tags").
		Where("expires_at IS NOT NULL AND expires_at <= DATE_ADD(NOW(), INTERVAL ? DAY)", days).
		Find(&domains).Error
	return domains, err
}

// Search 搜索域名
func (r *DomainRepository) Search(keyword string, limit, offset int) ([]*dns.Domain, int64, error) {
	var domains []*dns.Domain
	var total int64

	query := r.db.Preload("Group").Preload("Tags").Model(&dns.Domain{})

	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR remark LIKE ?", searchPattern, searchPattern)
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

	err := query.Find(&domains).Error
	return domains, total, err
}

// SearchWithFilters 带过滤条件的搜索
func (r *DomainRepository) SearchWithFilters(filters map[string]interface{}, limit, offset int) ([]*dns.Domain, int64, error) {
	var domains []*dns.Domain
	var total int64

	query := r.db.Preload("Group").Preload("Tags").Model(&dns.Domain{})

	// 应用过滤条件
	for key, value := range filters {
		if value != nil && value != "" {
			switch key {
			case "keyword":
				searchPattern := "%" + fmt.Sprintf("%v", value) + "%"
				query = query.Where("name LIKE ? OR remark LIKE ?", searchPattern, searchPattern)
			case "status":
				query = query.Where("status = ?", value)
			case "registrar_type":
				query = query.Where("registrar_type = ?", value)
			case "group_id":
				query = query.Where("group_id = ?", value)
			case "auto_renew":
				query = query.Where("auto_renew = ?", value)
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

	err := query.Find(&domains).Error
	return domains, total, err
}

// CountByStatus 统计各状态域名数量
func (r *DomainRepository) CountByStatus() (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	err := r.db.Model(&dns.Domain{}).
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

// BatchCreate 批量创建域名
func (r *DomainRepository) BatchCreate(domains []*dns.Domain) error {
	return r.db.Create(&domains).Error
}

// BatchUpdate 批量更新域名
func (r *DomainRepository) BatchUpdate(domains []*dns.Domain) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, domain := range domains {
			if err := tx.Model(domain).Updates(domain).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchDelete 批量删除域名
func (r *DomainRepository) BatchDelete(ids []uint) error {
	return r.db.Delete(&dns.Domain{}, ids).Error
}

// ExistsByName 检查域名是否存在
func (r *DomainRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&dns.Domain{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// CountByTenantID 统计租户域名数量
func (r *DomainRepository) CountByTenantID(tenantID uint) (int64, error) {
	var count int64
	err := r.db.Model(&dns.Domain{}).Where("tenant_id = ?", tenantID).Count(&count).Error
	return count, err
}

// UpsertByName 根据域名名称插入或更新，不修改已存在的 tenant_id
func (r *DomainRepository) UpsertByName(domain *dns.Domain) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"registrar_id", "registrar_type", "updated_by", "updated_at"}),
	}).Create(domain).Error
}

// FindByProviderID 根据提供商ID(RegistrarID)查找域名
func (r *DomainRepository) FindByProviderID(providerID uint) ([]*dns.Domain, error) {
	var domains []*dns.Domain
	err := r.db.Preload("Group").Preload("Tags").Where("registrar_id = ?", providerID).Find(&domains).Error
	return domains, err
}

// CountByProviderID 根据提供商ID统计域名数量
func (r *DomainRepository) CountByProviderID(providerID uint) (int64, error) {
	var count int64
	err := r.db.Model(&dns.Domain{}).Where("registrar_id = ?", providerID).Count(&count).Error
	return count, err
}
