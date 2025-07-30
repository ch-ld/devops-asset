package system

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/internal/config"
	"api-server/internal/model/nav"
)

// SystemDataConfig 系统数据配置
type SystemDataConfig struct {
	Menus       []MenuConfig       `json:"menus"`
	Roles       []RoleConfig       `json:"roles"`
	Departments []DepartmentConfig `json:"departments"`
	Users       []UserConfig       `json:"users"`
}

// MenuConfig 菜单配置
type MenuConfig struct {
	ID        uint   `json:"id"`
	Path      string `json:"path"`
	Name      string `json:"name"`
	Component string `json:"component"`
	Title     string `json:"title"`
	Icon      string `json:"icon"`
	KeepAlive int    `json:"keep_alive"`
	Status    int    `json:"status"`
	Level     int    `json:"level"`
	ParentID  uint   `json:"parent_id"`
	Sort      int    `json:"sort"`
	IsHide    int    `json:"is_hide"`
}

// RoleConfig 角色配置
type RoleConfig struct {
	ID      uint     `json:"id"`
	Name    string   `json:"name"`
	Desc    string   `json:"desc"`
	Status  int      `json:"status"`
	MenuIDs []uint   `json:"menu_ids"` // 分配的菜单ID列表
}

// DepartmentConfig 部门配置
type DepartmentConfig struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
}

// UserConfig 用户配置
type UserConfig struct {
	ID           uint   `json:"id"`
	DepartmentID uint   `json:"department_id"`
	RoleID       uint   `json:"role_id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	Password     string `json:"password"` // 将在创建时加密
	Status       int    `json:"status"`
	Gender       int    `json:"gender"`
}

// getSystemDataConfig 获取系统数据配置
func getSystemDataConfig() *SystemDataConfig {
	return &SystemDataConfig{
		Menus: []MenuConfig{
			{ID: 1, Path: "/dashboard", Name: "Dashboard", Component: "/index/index", Title: "仪表盘", Icon: "&#xe6b2;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 99},
			{ID: 2, Path: "/system", Name: "System", Component: "/index/index", Title: "系统管理", Icon: "&#xe72b;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 20},
			{ID: 3, Path: "menu", Name: "SystemMenu", Component: "/system/menu/index", Title: "菜单管理", Icon: "&#xe662;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 99},
			{ID: 4, Path: "role", Name: "SystemRole", Component: "/system/role/index", Title: "角色管理", Icon: "&#xe734;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 88},
			{ID: 5, Path: "department", Name: "SystemDepartment", Component: "/system/department/index", Title: "部门管理", Icon: "&#xe753;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 77},
			{ID: 6, Path: "user", Name: "SystemUser", Component: "/system/user/index", Title: "用户管理", Icon: "&#xe608;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 66},
			{ID: 7, Path: "console", Name: "DashboardConsole", Component: "/dashboard/console/index", Title: "工作台", Icon: "&#xe651;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 1, Sort: 99},
			{ID: 8, Path: "/private", Name: "Private", Component: "/index/index", Title: "隐藏页面", Icon: "&#xe636;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 99, IsHide: 1},
			{ID: 9, Path: "announcement", Name: "SystemAnnouncement", Component: "/system/announcement/index", Title: "公告管理", Icon: "&#xe747;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 55},
			{ID: 10, Path: "user-center", Name: "UserCenter", Component: "/system/user-center/index", Title: "个人中心", Icon: "&#xe734;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 50},
			// CMDB菜单
			{ID: 11, Path: "/cmdb", Name: "CMDB", Component: "/index/index", Title: "CMDB资产管理", Icon: "&#xe6a0;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 90},
		},
		Roles: []RoleConfig{
			{ID: 1, Name: "超级管理员", Desc: "拥有所有权限", Status: 1, MenuIDs: []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}}, // 所有菜单
			{ID: 2, Name: "普通用户", Desc: "普通用户", Status: 1, MenuIDs: []uint{1, 7, 10}}, // 仪表盘、工作台、个人中心
		},
		Departments: []DepartmentConfig{
			{ID: 1, Name: "管理中心", Sort: 1, Status: 1},
		},
		Users: []UserConfig{
			{ID: 1, DepartmentID: 1, RoleID: 1, Name: "超级管理员", Username: "admin", Password: config.AdminPassword, Status: 1, Gender: 1},
		},
	}
}

// MigrateData 执行系统数据迁移（优化版）
func MigrateData(db *gorm.DB) error {
	zap.L().Info("Starting optimized system data migration")

	return db.Transaction(func(tx *gorm.DB) error {
		config := getSystemDataConfig()

		// 1. 迁移菜单数据
		if err := migrateMenuData(tx, config.Menus); err != nil {
			return fmt.Errorf("failed to migrate menu data: %w", err)
		}

		// 2. 迁移角色数据
		if err := migrateRoleData(tx, config.Roles); err != nil {
			return fmt.Errorf("failed to migrate role data: %w", err)
		}

		// 3. 迁移部门数据
		if err := migrateDepartmentData(tx, config.Departments); err != nil {
			return fmt.Errorf("failed to migrate department data: %w", err)
		}

		// 4. 迁移用户数据
		if err := migrateUserData(tx, config.Users); err != nil {
			return fmt.Errorf("failed to migrate user data: %w", err)
		}

		// 5. 分配角色菜单权限
		if err := assignRoleMenuPermissions(tx, config.Roles); err != nil {
			return fmt.Errorf("failed to assign role menu permissions: %w", err)
		}

		zap.L().Info("System data migration completed successfully")
		return nil
	})
}

// migrateMenuData 迁移菜单数据
func migrateMenuData(tx *gorm.DB, menuConfigs []MenuConfig) error {
	var existingCount int64
	tx.Model(&SystemMenu{}).Count(&existingCount)
	
	if existingCount > 0 {
		zap.L().Info("Menu data already exists, performing incremental update")
		return updateMenusIfNeeded(tx, menuConfigs)
	}

	// 批量创建菜单
	var menus []SystemMenu
	for _, config := range menuConfigs {
		menu := SystemMenu{
			Model:     gorm.Model{ID: config.ID},
			Path:      config.Path,
			Name:      config.Name,
			Component: config.Component,
			Title:     config.Title,
			Icon:      config.Icon,
			KeepAlive: config.KeepAlive,
			Status:    config.Status,
			Level:     config.Level,
			ParentID:  config.ParentID,
			Sort:      config.Sort,
			IsHide:    config.IsHide,
		}
		menus = append(menus, menu)
	}

	if err := tx.Create(&menus).Error; err != nil {
		return fmt.Errorf("failed to create menus: %w", err)
	}

	zap.L().Info("Menu data created successfully", zap.Int("count", len(menus)))
	return nil
}

// updateMenusIfNeeded 根据需要更新菜单
func updateMenusIfNeeded(tx *gorm.DB, menuConfigs []MenuConfig) error {
	for _, config := range menuConfigs {
		var existingMenu SystemMenu
		err := tx.Where("id = ?", config.ID).First(&existingMenu).Error
		
		if err == gorm.ErrRecordNotFound {
			// 菜单不存在，创建新菜单
			newMenu := SystemMenu{
				Model:     gorm.Model{ID: config.ID},
				Path:      config.Path,
				Name:      config.Name,
				Component: config.Component,
				Title:     config.Title,
				Icon:      config.Icon,
				KeepAlive: config.KeepAlive,
				Status:    config.Status,
				Level:     config.Level,
				ParentID:  config.ParentID,
				Sort:      config.Sort,
				IsHide:    config.IsHide,
			}
			
			if err := tx.Create(&newMenu).Error; err != nil {
				return fmt.Errorf("failed to create new menu %s: %w", config.Name, err)
			}
			zap.L().Info("Created new menu", zap.String("name", config.Name))
		} else if err != nil {
			return fmt.Errorf("failed to query menu %d: %w", config.ID, err)
		}
		// 如果菜单存在，可以在这里添加更新逻辑
	}
	return nil
}

// migrateRoleData 迁移角色数据
func migrateRoleData(tx *gorm.DB, roleConfigs []RoleConfig) error {
	var existingCount int64
	tx.Model(&SystemRole{}).Count(&existingCount)
	
	if existingCount > 0 {
		zap.L().Info("Role data already exists, skipping role creation")
		return nil
	}

	var roles []SystemRole
	for _, config := range roleConfigs {
		role := SystemRole{
			Model:  gorm.Model{ID: config.ID},
			Name:   config.Name,
			Desc:   config.Desc,
			Status: config.Status,
		}
		roles = append(roles, role)
	}

	if err := tx.Create(&roles).Error; err != nil {
		return fmt.Errorf("failed to create roles: %w", err)
	}

	zap.L().Info("Role data created successfully", zap.Int("count", len(roles)))
	return nil
}

// migrateDepartmentData 迁移部门数据
func migrateDepartmentData(tx *gorm.DB, deptConfigs []DepartmentConfig) error {
	var existingCount int64
	tx.Model(&SystemDepartment{}).Count(&existingCount)
	
	if existingCount > 0 {
		zap.L().Info("Department data already exists, skipping department creation")
		return nil
	}

	var departments []SystemDepartment
	for _, config := range deptConfigs {
		dept := SystemDepartment{
			Model:  gorm.Model{ID: config.ID},
			Name:   config.Name,
			Sort:   config.Sort,
			Status: config.Status,
		}
		departments = append(departments, dept)
	}

	if err := tx.Create(&departments).Error; err != nil {
		return fmt.Errorf("failed to create departments: %w", err)
	}

	zap.L().Info("Department data created successfully", zap.Int("count", len(departments)))
	return nil
}

// migrateUserData 迁移用户数据
func migrateUserData(tx *gorm.DB, userConfigs []UserConfig) error {
	var existingCount int64
	tx.Model(&SystemUser{}).Count(&existingCount)
	
	if existingCount > 0 {
		zap.L().Info("User data already exists, skipping user creation")
		return nil
	}

	var users []SystemUser
	for _, config := range userConfigs {
		user := SystemUser{
			Model:        gorm.Model{ID: config.ID},
			DepartmentID: config.DepartmentID,
			RoleID:       config.RoleID,
			Name:         config.Name,
			Username:     config.Username,
			Password:     encryptionPWD(config.Password), // 加密密码
			Status:       config.Status,
			Gender:       config.Gender,
		}
		users = append(users, user)
	}

	if err := tx.Create(&users).Error; err != nil {
		return fmt.Errorf("failed to create users: %w", err)
	}

	zap.L().Info("User data created successfully", zap.Int("count", len(users)))
	return nil
}

// assignRoleMenuPermissions 分配角色菜单权限
func assignRoleMenuPermissions(tx *gorm.DB, roleConfigs []RoleConfig) error {
	for _, roleConfig := range roleConfigs {
		var role SystemRole
		if err := tx.First(&role, roleConfig.ID).Error; err != nil {
			zap.L().Warn("Role not found, skipping permission assignment", 
				zap.Uint("role_id", roleConfig.ID))
			continue
		}

		// 检查是否已有权限分配
		var existingCount int64
		tx.Model(&role).Association("SystemMenus").Count(&existingCount)
		if existingCount > 0 {
			zap.L().Info("Role permissions already exist, skipping", 
				zap.String("role", role.Name))
			continue
		}

		// 获取菜单并分配权限
		var menus []SystemMenu
		if err := tx.Where("id IN ?", roleConfig.MenuIDs).Find(&menus).Error; err != nil {
			return fmt.Errorf("failed to find menus for role %s: %w", role.Name, err)
		}

		if err := tx.Model(&role).Association("SystemMenus").Append(&menus); err != nil {
			return fmt.Errorf("failed to assign menus to role %s: %w", role.Name, err)
		}

		zap.L().Info("Role permissions assigned successfully", 
			zap.String("role", role.Name), 
			zap.Int("menu_count", len(menus)))
	}

	return nil
}

// Migrate 执行完整的系统模块迁移（保持兼容性）
func Migrate(db *gorm.DB) error {
	zap.L().Info("Starting system module migration")
	
	// 1. 迁移表结构
	if err := migrateTable(db); err != nil {
		return fmt.Errorf("failed to migrate system tables: %w", err)
	}

	// 2. 迁移数据
	if err := MigrateData(db); err != nil {
		return fmt.Errorf("failed to migrate system data: %w", err)
	}

	zap.L().Info("System module migration completed successfully")
	return nil
}

// migrateTable 迁移表结构（保持原有逻辑）
func migrateTable(db *gorm.DB) error {
	err := db.AutoMigrate(&SystemDepartment{}, &SystemRole{}, &SystemMenu{}, &SystemMenuAuth{}, &SystemUser{}, &SystemUserLoginLog{}, &SystemAnnouncement{}, &nav.Nav{})
	if err != nil {
		zap.L().Error("failed to migrate system model", zap.Error(err))
		return err
	}
	return nil
}
