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
	// 创建根级分组
	rootGroups := []cmdb.HostGroup{
		{Name: "生产环境", Description: "生产环境主机组", ParentID: nil, Sort: 1},
		{Name: "测试环境", Description: "测试环境主机组", ParentID: nil, Sort: 2},
		{Name: "开发环境", Description: "开发环境主机组", ParentID: nil, Sort: 3},
	}

	// 先创建根级分组
	var createdGroups []cmdb.HostGroup
	for _, group := range rootGroups {
		if err := s.groupRepo.Create(&group); err != nil {
			return fmt.Errorf("创建默认主机组失败: %w", err)
		}
		createdGroups = append(createdGroups, group)
	}

	// 创建子分组
	subGroups := []struct {
		name        string
		description string
		parentName  string
		sort        int
	}{
		{"Web服务器", "Web应用服务器", "生产环境", 1},
		{"数据库服务器", "数据库服务器", "生产环境", 2},
		{"缓存服务器", "Redis/Memcached服务器", "生产环境", 3},
		{"测试服务器", "功能测试服务器", "测试环境", 1},
		{"性能测试", "性能压测服务器", "测试环境", 2},
		{"开发服务器", "开发人员使用的服务器", "开发环境", 1},
		{"CI/CD", "持续集成部署服务器", "开发环境", 2},
	}

	// 为子分组找到父分组ID并创建
	for _, subGroup := range subGroups {
		var parentID *uint
		for _, parent := range createdGroups {
			if parent.Name == subGroup.parentName {
				parentID = &parent.ID
				break
			}
		}

		if parentID != nil {
			childGroup := cmdb.HostGroup{
				Name:        subGroup.name,
				Description: subGroup.description,
				ParentID:    parentID,
				Sort:        subGroup.sort,
			}
			if err := s.groupRepo.Create(&childGroup); err != nil {
				return fmt.Errorf("创建子主机组失败: %w", err)
			}
		}
	}

	return nil
}

// buildTree 构建树状结构
func (s *HostGroupService) buildTree(groups []cmdb.HostGroup) []cmdb.HostGroup {
	// 创建一个map来快速查找节点
	groupMap := make(map[uint]*cmdb.HostGroup)

	// 初始化所有节点的children
	for i := range groups {
		groups[i].Children = []cmdb.HostGroup{}
		groupMap[groups[i].ID] = &groups[i]
	}

	// 构建树状结构
	var tree []cmdb.HostGroup
	for i := range groups {
		group := &groups[i]
		if group.ParentID == nil {
			// 根节点，直接添加到树中
			tree = append(tree, *group)
		} else {
			// 子节点，添加到父节点的children中
			if parent, exists := groupMap[*group.ParentID]; exists {
				parent.Children = append(parent.Children, *group)
			}
		}
	}

	// 更新树中根节点的children（因为上面的操作修改了原始数据）
	for i := range tree {
		if rootNode, exists := groupMap[tree[i].ID]; exists {
			tree[i].Children = rootNode.Children
		}
	}

	// 递归计算每个节点的总主机数量（包含子分组）
	s.calculateTotalHostCount(&tree)

	return tree
}

// calculateTotalHostCount 递归计算节点的总主机数量（包含子分组）
func (s *HostGroupService) calculateTotalHostCount(groups *[]cmdb.HostGroup) {
	for i := range *groups {
		group := &(*groups)[i]

		// 递归计算子分组的主机数量
		if len(group.Children) > 0 {
			s.calculateTotalHostCount(&group.Children)

			// 计算子分组的主机总数
			var childrenTotalHostCount int64
			for _, child := range group.Children {
				childrenTotalHostCount += child.TotalHostCount
			}

			// 当前分组的总主机数量 = 直接归属的主机数量 + 子分组的总主机数量
			group.TotalHostCount = group.HostCount + childrenTotalHostCount
		} else {
			// 没有子分组，总主机数量就是直接归属的主机数量
			group.TotalHostCount = group.HostCount
		}
	}
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
