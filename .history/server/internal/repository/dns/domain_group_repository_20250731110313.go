package dns

import (
	"api-server/internal/model/dns"
	"fmt"

	"gorm.io/gorm"
)

// DomainGroupRepository 域名分组数据访问层
type DomainGroupRepository struct {
	db *gorm.DB
}

// NewDomainGroupRepository 创建域名分组仓库实例
func NewDomainGroupRepository(db *gorm.DB) *DomainGroupRepository {
	return &DomainGroupRepository{db: db}
}

// Create 创建域名分组
func (r *DomainGroupRepository) Create(group *dns.DomainGroup) error {
	return r.db.Create(group).Error
}

// Update 更新域名分组
func (r *DomainGroupRepository) Update(group *dns.DomainGroup) error {
	return r.db.Model(group).Updates(group).Error
}

// Delete 删除域名分组
func (r *DomainGroupRepository) Delete(id uint) error {
	return r.db.Delete(&dns.DomainGroup{}, id).Error
}

// FindByID 根据ID查找域名分组
func (r *DomainGroupRepository) FindByID(id uint) (*dns.DomainGroup, error) {
	var group dns.DomainGroup
	err := r.db.Preload("Parent").Preload("Children").First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// FindByName 根据名称查找域名分组
func (r *DomainGroupRepository) FindByName(name string) (*dns.DomainGroup, error) {
	var group dns.DomainGroup
	err := r.db.Where("name = ?", name).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// FindByParentID 根据父分组ID查找子分组
func (r *DomainGroupRepository) FindByParentID(parentID *uint) ([]*dns.DomainGroup, error) {
	var groups []*dns.DomainGroup
	query := r.db.Model(&dns.DomainGroup{})
	
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	
	err := query.Order("sort ASC, name ASC").Find(&groups).Error
	return groups, err
}

// FindRootGroups 查找根分组
func (r *DomainGroupRepository) FindRootGroups() ([]*dns.DomainGroup, error) {
	var groups []*dns.DomainGroup
	err := r.db.Where("parent_id IS NULL").Order("sort ASC, name ASC").Find(&groups).Error
	return groups, err
}

// FindByTenantID 根据租户ID查找域名分组
func (r *DomainGroupRepository) FindByTenantID(tenantID uint) ([]*dns.DomainGroup, error) {
	var groups []*dns.DomainGroup
	err := r.db.Where("tenant_id = ?", tenantID).Order("sort ASC, name ASC").Find(&groups).Error
	return groups, err
}

// FindAll 查找所有域名分组
func (r *DomainGroupRepository) FindAll() ([]*dns.DomainGroup, error) {
	var groups []*dns.DomainGroup
	err := r.db.Order("sort ASC, name ASC").Find(&groups).Error
	return groups, err
}

// GetTree 获取分组树结构
func (r *DomainGroupRepository) GetTree(tenantID uint) ([]*dns.DomainGroup, error) {
	var groups []*dns.DomainGroup
	err := r.db.Where("tenant_id = ?", tenantID).Order("sort ASC, name ASC").Find(&groups).Error
	if err != nil {
		return nil, err
	}

	// 构建树结构
	groupMap := make(map[uint]*dns.DomainGroup)
	var rootGroups []*dns.DomainGroup

	// 第一遍：创建映射
	for i := range groups {
		groupMap[groups[i].ID] = &groups[i]
		groups[i].Children = []*dns.DomainGroup{}
	}

	// 第二遍：构建父子关系
	for i := range groups {
		if groups[i].ParentID == nil {
			rootGroups = append(rootGroups, &groups[i])
		} else {
			if parent, exists := groupMap[*groups[i].ParentID]; exists {
				parent.Children = append(parent.Children, &groups[i])
			}
		}
	}

	return rootGroups, nil
}

// Search 搜索域名分组
func (r *DomainGroupRepository) Search(keyword string, limit, offset int) ([]*dns.DomainGroup, int64, error) {
	var groups []*dns.DomainGroup
	var total int64

	query := r.db.Model(&dns.DomainGroup{})

	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
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

	err := query.Order("sort ASC, name ASC").Find(&groups).Error
	return groups, total, err
}

// SearchWithFilters 带过滤条件的搜索
func (r *DomainGroupRepository) SearchWithFilters(filters map[string]interface{}, limit, offset int) ([]*dns.DomainGroup, int64, error) {
	var groups []*dns.DomainGroup
	var total int64

	query := r.db.Model(&dns.DomainGroup{})

	// 应用过滤条件
	for key, value := range filters {
		if value != nil && value != "" {
			switch key {
			case "keyword":
				searchPattern := "%" + fmt.Sprintf("%v", value) + "%"
				query = query.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
			case "parent_id":
				if value == 0 {
					query = query.Where("parent_id IS NULL")
				} else {
					query = query.Where("parent_id = ?", value)
				}
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

	err := query.Order("sort ASC, name ASC").Find(&groups).Error
	return groups, total, err
}

// BatchCreate 批量创建域名分组
func (r *DomainGroupRepository) BatchCreate(groups []*dns.DomainGroup) error {
	return r.db.Create(&groups).Error
}

// BatchUpdate 批量更新域名分组
func (r *DomainGroupRepository) BatchUpdate(groups []*dns.DomainGroup) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, group := range groups {
			if err := tx.Model(group).Updates(group).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchDelete 批量删除域名分组
func (r *DomainGroupRepository) BatchDelete(ids []uint) error {
	return r.db.Delete(&dns.DomainGroup{}, ids).Error
}

// ExistsByName 检查分组名称是否存在
func (r *DomainGroupRepository) ExistsByName(name string, parentID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&dns.DomainGroup{}).Where("name = ?", name)
	
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	
	err := query.Count(&count).Error
	return count > 0, err
}

// CountByTenantID 统计租户分组数量
func (r *DomainGroupRepository) CountByTenantID(tenantID uint) (int64, error) {
	var count int64
	err := r.db.Model(&dns.DomainGroup{}).Where("tenant_id = ?", tenantID).Count(&count).Error
	return count, err
}

// HasChildren 检查分组是否有子分组
func (r *DomainGroupRepository) HasChildren(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&dns.DomainGroup{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

// HasDomains 检查分组是否有域名
func (r *DomainGroupRepository) HasDomains(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&dns.Domain{}).Where("group_id = ?", id).Count(&count).Error
	return count > 0, err
}

// UpdateSort 更新排序
func (r *DomainGroupRepository) UpdateSort(id uint, sort int) error {
	return r.db.Model(&dns.DomainGroup{}).Where("id = ?", id).Update("sort", sort).Error
}

// GetMaxSort 获取最大排序值
func (r *DomainGroupRepository) GetMaxSort(parentID *uint) (int, error) {
	var maxSort int
	query := r.db.Model(&dns.DomainGroup{}).Select("COALESCE(MAX(sort), 0)")
	
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	
	err := query.Scan(&maxSort).Error
	return maxSort, err
}
