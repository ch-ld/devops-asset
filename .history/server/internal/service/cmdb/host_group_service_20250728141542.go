package cmdb

import (
	"api-server/internal/model/cmdb"
	repo "api-server/internal/repository/cmdb"
	"fmt"

	"gorm.io/gorm"
)

// HostGroupService 主机组业务逻辑服务
type HostGroupService struct {
	groupRepo *repo.HostGroupRepository
	hostRepo  *repo.HostRepository
}

// NewHostGroupService 创建主机组业务服务实例
func NewHostGroupService(groupRepo *repo.HostGroupRepository, hostRepo *repo.HostRepository) *HostGroupService {
	return &HostGroupService{
		groupRepo: groupRepo,
		hostRepo:  hostRepo,
	}
}

// CreateHostGroup 创建主机组（别名方法）
func (s *HostGroupService) CreateHostGroup(group *cmdb.HostGroup) error {
	return s.CreateGroup(group)
}

// CreateGroup 创建主机组
func (s *HostGroupService) CreateGroup(group *cmdb.HostGroup) error {
	// 检查父组是否存在
	if group.ParentID != nil {
		if _, err := s.groupRepo.FindByID(*group.ParentID); err != nil {
			return fmt.Errorf("父组不存在: %w", err)
		}
	}

	return s.groupRepo.Create(group)
}

// UpdateHostGroup 更新主机组（别名方法）
func (s *HostGroupService) UpdateHostGroup(group *cmdb.HostGroup) error {
	return s.UpdateGroup(group)
}

// UpdateGroup 更新主机组
func (s *HostGroupService) UpdateGroup(group *cmdb.HostGroup) error {
	// 检查组是否存在
	oldGroup, err := s.groupRepo.FindByID(group.ID)
	if err != nil {
		return fmt.Errorf("主机组不存在: %w", err)
	}

	// 检查父组是否存在
	if group.ParentID != nil {
		if _, err := s.groupRepo.FindByID(*group.ParentID); err != nil {
			return fmt.Errorf("父组不存在: %w", err)
		}
		// 检查是否形成循环依赖
		if *group.ParentID == group.ID {
			return fmt.Errorf("不能将组的父组设置为自身")
		}
	}

	// 如果修改了名称，需要更新所有子组的路径
	if oldGroup.Name != group.Name {
		descendants, err := s.groupRepo.FindDescendants(oldGroup.Path)
		if err != nil {
			return fmt.Errorf("获取子组失败: %w", err)
		}
		for _, desc := range descendants {
			desc.Path = group.Path + desc.Path[len(oldGroup.Path):]
			if err := s.groupRepo.Update(&desc); err != nil {
				return fmt.Errorf("更新子组路径失败: %w", err)
			}
		}
	}

	return s.groupRepo.Update(group)
}

// DeleteHostGroup 删除主机组（别名方法）
func (s *HostGroupService) DeleteHostGroup(id uint) error {
	return s.DeleteGroup(id)
}

// DeleteGroup 删除主机组
func (s *HostGroupService) DeleteGroup(id uint) error {
	return s.groupRepo.Delete(id)
}

// GetHostGroup 获取主机组信息（别名方法）
func (s *HostGroupService) GetHostGroup(id uint) (*cmdb.HostGroup, error) {
	return s.GetGroup(id)
}

// GetGroup 获取主机组信息
func (s *HostGroupService) GetGroup(id uint) (*cmdb.HostGroup, error) {
	return s.groupRepo.FindByID(id)
}

// GetHostGroupTree 获取主机组树结构（别名方法）
func (s *HostGroupService) GetHostGroupTree() ([]cmdb.HostGroup, error) {
	groups, err := s.ListGroups("")
	if err != nil {
		return nil, err
	}

	// 如果没有数据，创建一些默认的主机组
	if len(groups) == 0 {
		err = s.createDefaultHostGroups()
		if err != nil {
			return nil, err
		}
		// 重新获取数据
		groups, err = s.ListGroups("")
		if err != nil {
			return nil, err
		}
	}

	// 获取每个主机组的主机数量统计
	hostCounts, err := s.groupRepo.GetHostCountByGroups()
	if err != nil {
		return nil, fmt.Errorf("获取主机组统计失败: %w", err)
	}

	// 为每个主机组添加主机数量信息
	for i := range groups {
		if count, exists := hostCounts[groups[i].ID]; exists {
			groups[i].HostCount = count
		} else {
			groups[i].HostCount = 0
		}
	}

	// 构建树状结构
	return s.buildTree(groups), nil
}

// createDefaultHostGroups 创建默认主机组
func (s *HostGroupService) createDefaultHostGroups() error {
	defaultGroups := []cmdb.HostGroup{
		{Name: "生产环境", Description: "生产环境主机组", ParentID: nil, Sort: 1},
		{Name: "测试环境", Description: "测试环境主机组", ParentID: nil, Sort: 2},
		{Name: "开发环境", Description: "开发环境主机组", ParentID: nil, Sort: 3},
	}

	for _, group := range defaultGroups {
		if err := s.groupRepo.Create(&group); err != nil {
			return fmt.Errorf("创建默认主机组失败: %w", err)
		}
	}

	return nil
}

// ListGroups 查询主机组列表
func (s *HostGroupService) ListGroups(keyword string) ([]cmdb.HostGroup, error) {
	return s.groupRepo.List(keyword)
}

// GetGroupHosts 获取主机组下的主机列表
func (s *HostGroupService) GetGroupHosts(groupID uint, page, pageSize int, keyword string) ([]cmdb.Host, int64, error) {
	return s.hostRepo.FindByGroupID(groupID, page, pageSize, keyword)
}

// MoveHostGroup 移动主机组（别名方法）
func (s *HostGroupService) MoveHostGroup(groupID uint, parentID *uint) error {
	// 检查组是否存在
	if _, err := s.groupRepo.FindByID(groupID); err != nil {
		return fmt.Errorf("主机组不存在: %w", err)
	}

	// 如果指定了父组，检查父组是否存在
	if parentID != nil {
		if _, err := s.groupRepo.FindByID(*parentID); err != nil {
			return fmt.Errorf("父组不存在: %w", err)
		}
	}

	// 更新父组ID
	group := &cmdb.HostGroup{
		Model:    gorm.Model{ID: groupID},
		ParentID: parentID,
	}
	return s.groupRepo.Update(group)
}

// MoveHost 移动主机到指定组
func (s *HostGroupService) MoveHost(hostID uint, groupID *uint) error {
	// 检查主机是否存在
	if _, err := s.hostRepo.FindByID(hostID); err != nil {
		return fmt.Errorf("主机不存在: %w", err)
	}

	// 如果指定了组，检查组是否存在
	if groupID != nil {
		if _, err := s.groupRepo.FindByID(*groupID); err != nil {
			return fmt.Errorf("主机组不存在: %w", err)
		}
	}

	return s.groupRepo.UpdateHostGroup(hostID, groupID)
}
