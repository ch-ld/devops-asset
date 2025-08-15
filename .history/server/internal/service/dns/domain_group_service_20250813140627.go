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

// DomainGroupService 域名分组业务逻辑服务
type DomainGroupService struct {
	domainGroupRepo *repo.DomainGroupRepository
	domainRepo      *repo.DomainRepository
	changeLogRepo   *repo.ChangeLogRepository
	db              *gorm.DB
}

// NewDomainGroupService 创建域名分组业务服务实例
func NewDomainGroupService(
	domainGroupRepo *repo.DomainGroupRepository,
	domainRepo *repo.DomainRepository,
	changeLogRepo *repo.ChangeLogRepository,
	db *gorm.DB,
) *DomainGroupService {
	return &DomainGroupService{
		domainGroupRepo: domainGroupRepo,
		domainRepo:      domainRepo,
		changeLogRepo:   changeLogRepo,
		db:              db,
	}
}

// CreateDomainGroup 创建域名分组
func (s *DomainGroupService) CreateDomainGroup(group *dns.DomainGroup, actorID uint, clientIP string) error {
	// 验证分组名称
	if group.Name == "" {
		return errors.New("分组名称不能为空")
	}

	// 验证父分组是否存在
	if group.ParentID != nil {
		_, err := s.domainGroupRepo.FindByID(*group.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("指定的父分组不存在")
			}
			return fmt.Errorf("验证父分组失败: %w", err)
		}
	}

	// 检查同级分组名称是否已存在
	exists, err := s.domainGroupRepo.ExistsByName(group.Name, group.ParentID)
	if err != nil {
		return fmt.Errorf("检查分组名称是否存在失败: %w", err)
	}
	if exists {
		return errors.New("同级分组中已存在相同名称的分组")
	}

	// 设置排序值
	if group.Sort == 0 {
		maxSort, err := s.domainGroupRepo.GetMaxSort(group.ParentID)
		if err != nil {
			return fmt.Errorf("获取最大排序值失败: %w", err)
		}
		group.Sort = maxSort + 1
	}

	// 设置创建信息
	group.CreatedBy = actorID
	group.UpdatedBy = actorID

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 创建分组
		if err := s.domainGroupRepo.Create(group); err != nil {
			return fmt.Errorf("创建域名分组失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "domain_group",
			ResourceID:   group.ID,
			Action:       "create",
			Description:  fmt.Sprintf("创建域名分组: %s", group.Name),
			NewData:      s.domainGroupToJSON(group),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     group.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录域名分组创建日志失败", zap.Error(err))
		}

		return nil
	})
}

// UpdateDomainGroup 更新域名分组
func (s *DomainGroupService) UpdateDomainGroup(group *dns.DomainGroup, actorID uint, clientIP string) error {
	// 获取原始数据
	oldGroup, err := s.domainGroupRepo.FindByID(group.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("域名分组不存在")
		}
		return fmt.Errorf("获取域名分组信息失败: %w", err)
	}

	// 验证分组名称
	if group.Name == "" {
		return errors.New("分组名称不能为空")
	}

	// 验证父分组是否存在
	if group.ParentID != nil {
		_, err := s.domainGroupRepo.FindByID(*group.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("指定的父分组不存在")
			}
			return fmt.Errorf("验证父分组失败: %w", err)
		}

		// 防止循环引用
		if *group.ParentID == group.ID {
			return errors.New("不能将分组设置为自己的子分组")
		}
	}

	// 如果名称或父分组发生变化，检查是否冲突
	if oldGroup.Name != group.Name ||
		(oldGroup.ParentID == nil && group.ParentID != nil) ||
		(oldGroup.ParentID != nil && group.ParentID == nil) ||
		(oldGroup.ParentID != nil && group.ParentID != nil && *oldGroup.ParentID != *group.ParentID) {

		exists, err := s.domainGroupRepo.ExistsByName(group.Name, group.ParentID)
		if err != nil {
			return fmt.Errorf("检查分组名称是否存在失败: %w", err)
		}
		if exists {
			return errors.New("同级分组中已存在相同名称的分组")
		}
	}

	// 设置更新信息
	group.UpdatedBy = actorID

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 更新分组
		if err := s.domainGroupRepo.Update(group); err != nil {
			return fmt.Errorf("更新域名分组失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "domain_group",
			ResourceID:   group.ID,
			Action:       "update",
			Description:  fmt.Sprintf("更新域名分组: %s", group.Name),
			OldData:      s.domainGroupToJSON(oldGroup),
			NewData:      s.domainGroupToJSON(group),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     group.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录域名分组更新日志失败", zap.Error(err))
		}

		return nil
	})
}

// DeleteDomainGroup 删除域名分组
func (s *DomainGroupService) DeleteDomainGroup(id uint, actorID uint, clientIP string) error {
	// 获取分组信息
	group, err := s.domainGroupRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("域名分组不存在")
		}
		return fmt.Errorf("获取域名分组信息失败: %w", err)
	}

	// 检查是否有子分组
	hasChildren, err := s.domainGroupRepo.HasChildren(id)
	if err != nil {
		return fmt.Errorf("检查子分组失败: %w", err)
	}
	if hasChildren {
		return errors.New("存在子分组，无法删除")
	}

	// 检查是否有关联的域名
	hasDomains, err := s.domainGroupRepo.HasDomains(id)
	if err != nil {
		return fmt.Errorf("检查关联域名失败: %w", err)
	}
	if hasDomains {
		return errors.New("存在关联的域名，无法删除")
	}

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除分组
		if err := s.domainGroupRepo.Delete(id); err != nil {
			return fmt.Errorf("删除域名分组失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "domain_group",
			ResourceID:   id,
			Action:       "delete",
			Description:  fmt.Sprintf("删除域名分组: %s", group.Name),
			OldData:      s.domainGroupToJSON(group),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     group.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录域名分组删除日志失败", zap.Error(err))
		}

		return nil
	})
}

// GetDomainGroup 获取域名分组详情
func (s *DomainGroupService) GetDomainGroup(id uint) (*dns.DomainGroup, error) {
	group, err := s.domainGroupRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("域名分组不存在")
		}
		return nil, fmt.Errorf("获取域名分组信息失败: %w", err)
	}
	return group, nil
}

// GetDomainGroupTree 获取域名分组树
func (s *DomainGroupService) GetDomainGroupTree(tenantID uint) ([]*dns.DomainGroup, error) {
	tree, err := s.domainGroupRepo.GetTree(tenantID)
	if err != nil {
		return nil, fmt.Errorf("获取域名分组树失败: %w", err)
	}

	// 获取每个分组的域名数量统计
	domainCounts, err := s.domainRepo.GetDomainCountByGroups(tenantID)
	if err != nil {
		return nil, fmt.Errorf("获取分组域名统计失败: %w", err)
	}

	// 为每个分组添加域名数量信息
	s.calculateDomainCounts(tree, domainCounts)

	return tree, nil
}

// ListDomainGroups 获取域名分组列表
func (s *DomainGroupService) ListDomainGroups(tenantID uint, filters map[string]interface{}, limit, offset int) ([]*dns.DomainGroup, int64, error) {
	// 添加租户过滤
	if filters == nil {
		filters = make(map[string]interface{})
	}
	filters["tenant_id"] = tenantID

	groups, total, err := s.domainGroupRepo.SearchWithFilters(filters, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("获取域名分组列表失败: %w", err)
	}

	return groups, total, nil
}

// UpdateDomainGroupSort 更新域名分组排序
func (s *DomainGroupService) UpdateDomainGroupSort(id uint, sort int, actorID uint, clientIP string) error {
	// 验证分组是否存在
	group, err := s.domainGroupRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("域名分组不存在")
		}
		return fmt.Errorf("获取域名分组信息失败: %w", err)
	}

	// 更新排序
	if err := s.domainGroupRepo.UpdateSort(id, sort); err != nil {
		return fmt.Errorf("更新域名分组排序失败: %w", err)
	}

	// 记录变更日志
	changeLog := &dns.ChangeLog{
		ResourceType: "domain_group",
		ResourceID:   id,
		Action:       "update_sort",
		Description:  fmt.Sprintf("更新域名分组排序: %s (排序: %d)", group.Name, sort),
		Status:       "success",
		ClientIP:     clientIP,
		TenantID:     group.TenantID,
		ActorID:      actorID,
	}
	if err := s.changeLogRepo.Create(changeLog); err != nil {
		zap.L().Error("记录域名分组排序更新日志失败", zap.Error(err))
	}

	return nil
}

// GetDomainGroupStatistics 获取域名分组统计信息
func (s *DomainGroupService) GetDomainGroupStatistics(tenantID uint) (map[string]interface{}, error) {
	// 获取总分组数
	totalGroups, err := s.domainGroupRepo.CountByTenant(tenantID)
	if err != nil {
		return nil, fmt.Errorf("获取分组总数失败: %w", err)
	}

	// 获取活跃分组数
	activeGroups, err := s.domainGroupRepo.CountByStatus(tenantID, "active")
	if err != nil {
		return nil, fmt.Errorf("获取活跃分组数失败: %w", err)
	}

	// 获取非活跃分组数
	inactiveGroups := totalGroups - activeGroups

	// 构建统计信息
	stats := map[string]interface{}{
		"total_groups":    totalGroups,
		"active_groups":   activeGroups,
		"inactive_groups": inactiveGroups,
	}

	return stats, nil
}

// domainGroupToJSON 将域名分组对象转换为JSON
func (s *DomainGroupService) domainGroupToJSON(group *dns.DomainGroup) []byte {
	data, err := json.Marshal(group)
	if err != nil {
		zap.L().Error("域名分组对象转JSON失败", zap.Error(err))
		return []byte("{}")
	}
	return data
}
