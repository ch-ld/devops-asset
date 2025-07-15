package system

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/internal/config"
	"api-server/internal/model/nav"
)

func migrateTable(db *gorm.DB) error {
	err := db.AutoMigrate(&SystemDepartment{}, &SystemRole{}, &SystemMenu{}, &SystemMenuAuth{}, &SystemUser{}, &SystemUserLoginLog{}, &SystemAnnouncement{}, &nav.Nav{})
	if err != nil {
		zap.L().Error("failed to migrate system model", zap.Error(err))
		return err
	}
	return nil
}

func migrateData(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 检查是否已有数据，如果有则跳过初始化
		var count int64
		tx.Model(&SystemMenu{}).Count(&count)
		if count > 0 {
			zap.L().Info("menu data already exists, skipping initial data creation")
			// 检查是否存在公告管理菜单，如果不存在则添加
			if err := addAnnouncementMenuIfNotExists(tx); err != nil {
				return err
			}
			// 检查是否存在用户中心菜单，如果不存在则添加
			if err := addUserCenterMenuIfNotExists(tx); err != nil {
				return err
			}
			// 检查是否存在CMDB菜单，如果不存在则添加
			if err := addCmdbMenuIfNotExists(tx); err != nil {
				return err
			}
			return nil
		}

		// 创建菜单
		menus := []SystemMenu{
			{Model: gorm.Model{ID: 1}, Path: "/dashboard", Name: "Dashboard", Component: "/index/index", Title: "仪表盘", Icon: "&#xe6b2;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 99},
			{Model: gorm.Model{ID: 2}, Path: "/system", Name: "System", Component: "/index/index", Title: "系统管理", Icon: "&#xe72b;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 20},
			{Model: gorm.Model{ID: 3}, Path: "menu", Name: "SystemMenu", Component: "/system/menu/index", Title: "菜单管理", Icon: "&#xe662;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 99},
			{Model: gorm.Model{ID: 4}, Path: "role", Name: "SystemRole", Component: "/system/role/index", Title: "角色管理", Icon: "&#xe734;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 88},
			{Model: gorm.Model{ID: 5}, Path: "department", Name: "SystemDepartment", Component: "/system/department/index", Title: "部门管理", Icon: "&#xe753;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 77},
			{Model: gorm.Model{ID: 6}, Path: "user", Name: "SystemUser", Component: "/system/user/index", Title: "用户管理", Icon: "&#xe608;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 66},
			{Model: gorm.Model{ID: 7}, Path: "console", Name: "DashboardConsole", Component: "/dashboard/console/index", Title: "工作台", Icon: "&#xe651;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 1, Sort: 99},
			{Model: gorm.Model{ID: 8}, Path: "/private", Name: "Private", Component: "/index/index", Title: "隐藏页面", Icon: "&#xe636;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 99, IsHide: 1},
			{Model: gorm.Model{ID: 9}, Path: "announcement", Name: "SystemAnnouncement", Component: "/system/announcement/index", Title: "公告管理", Icon: "&#xe747;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 55},
			{Model: gorm.Model{ID: 10}, Path: "user-center", Name: "UserCenter", Component: "/system/user-center/index", Title: "个人中心", Icon: "&#xe734;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 50},
		}
		err := db.Create(&menus).Error
		if err != nil {
			zap.L().Error("failed to create menu", zap.Error(err))
			return err
		}

		// 检查是否已有角色数据
		tx.Model(&SystemRole{}).Count(&count)
		if count > 0 {
			zap.L().Info("role data already exists, skipping role creation")
			return nil
		}

		// 创建角色
		roles := []SystemRole{
			{Model: gorm.Model{ID: 1}, Name: "超级管理员", Desc: "拥有所有权限", Status: 1},
			{Model: gorm.Model{ID: 2}, Name: "普通用户", Desc: "普通用户", Status: 1},
		}
		err = db.Create(&roles).Error
		if err != nil {
			zap.L().Error("failed to create role", zap.Error(err))
			return err
		}

		// 为角色分配菜单权限
		// 超级管理员拥有所有菜单权限
		adminRole := SystemRole{}
		err = db.First(&adminRole, 1).Error
		if err != nil {
			zap.L().Error("failed to find admin role", zap.Error(err))
			return err
		}
		// 为超级管理员分配所有菜单
		var allMenus []SystemMenu
		err = db.Find(&allMenus).Error
		if err != nil {
			zap.L().Error("failed to find menus", zap.Error(err))
			return err
		}
		err = db.Model(&adminRole).Association("SystemMenus").Append(&allMenus)
		if err != nil {
			zap.L().Error("failed to associate menus with admin role", zap.Error(err))
			return err
		}
		// 为普通用户分配首页菜单
		normalRole := SystemRole{}
		err = db.First(&normalRole, 2).Error
		if err != nil {
			zap.L().Error("failed to find normal role", zap.Error(err))
			return err
		}
		// 为普通用户分配工作台和分析页菜单
		var consoleMenu, dashboardMenu SystemMenu
		err = db.First(&dashboardMenu, 1).Error
		if err != nil {
			zap.L().Error("failed to find dashboard menu", zap.Error(err))
			return err
		}
		err = db.First(&consoleMenu, 7).Error
		if err != nil {
			zap.L().Error("failed to find console menu", zap.Error(err))
			return err
		}
		err = db.Model(&normalRole).Association("SystemMenus").Append([]SystemMenu{dashboardMenu, consoleMenu})
		if err != nil {
			zap.L().Error("failed to associate console menus with normal role", zap.Error(err))
			return err
		}

		// 检查是否已有部门数据
		tx.Model(&SystemDepartment{}).Count(&count)
		if count > 0 {
			zap.L().Info("department data already exists, skipping department creation")
			return nil
		}

		// 创建部门
		departments := []SystemDepartment{
			{Model: gorm.Model{ID: 1}, Name: "管理中心", Sort: 1, Status: 1},
		}
		err = db.Create(&departments).Error
		if err != nil {
			zap.L().Error("failed to create department", zap.Error(err))
			return err
		}

		// 检查是否已有用户数据
		tx.Model(&SystemUser{}).Count(&count)
		if count > 0 {
			zap.L().Info("user data already exists, skipping user creation")
			return nil
		}

		// 创建用户
		pwd := encryptionPWD(config.AdminPassword)
		users := []SystemUser{
			{Model: gorm.Model{ID: 1}, DepartmentID: 1, RoleID: 1, Name: "超级管理员", Username: "admin", Password: pwd, Status: 1, Gender: 1},
		}
		err = db.Create(&users).Error
		if err != nil {
			zap.L().Error("failed to create user", zap.Error(err))
			return err
		}
		return nil
	})
	return err
}

// addUserCenterMenuIfNotExists 如果用户中心菜单不存在则添加
func addUserCenterMenuIfNotExists(db *gorm.DB) error {
	var count int64
	db.Model(&SystemMenu{}).Where("name = ? AND title = ?", "UserCenter", "个人中心").Count(&count)
	if count > 0 {
		zap.L().Info("user center menu already exists")
		return nil
	}

	// 添加用户中心菜单
	userCenterMenu := SystemMenu{
		Path:      "user-center",
		Name:      "UserCenter",
		Component: "/system/user-center/index",
		Title:     "个人中心",
		Icon:      "&#xe734;",
		KeepAlive: 2,
		Status:    1,
		Level:     2,
		ParentID:  2, // 系统管理的子菜单
		Sort:      50,
	}

	err := db.Create(&userCenterMenu).Error
	if err != nil {
		zap.L().Error("failed to create user center menu", zap.Error(err))
		return err
	}

	// 为超级管理员分配新菜单权限
	var adminRole SystemRole
	err = db.First(&adminRole, 1).Error
	if err != nil {
		zap.L().Error("failed to find admin role", zap.Error(err))
		return err
	}

	err = db.Model(&adminRole).Association("SystemMenus").Append(&userCenterMenu)
	if err != nil {
		zap.L().Error("failed to associate user center menu with admin role", zap.Error(err))
		return err
	}

	// 为普通用户也分配用户中心菜单权限
	var normalRole SystemRole
	err = db.First(&normalRole, 2).Error
	if err != nil {
		zap.L().Error("failed to find normal role", zap.Error(err))
		return err
	}

	err = db.Model(&normalRole).Association("SystemMenus").Append(&userCenterMenu)
	if err != nil {
		zap.L().Error("failed to associate user center menu with normal role", zap.Error(err))
		return err
	}

	zap.L().Info("successfully added user center menu")
	return nil
}

// addCmdbMenuIfNotExists adds the CMDB menu items if they don't exist.
func addCmdbMenuIfNotExists(db *gorm.DB) error {
	// 1. Check for the parent CMDB menu
	var parentMenu SystemMenu
	err := db.Where("name = 'CMDB'").First(&parentMenu).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("failed to query CMDB menu", zap.Error(err))
		return err
	}

	// If parent menu doesn't exist, create it.
	if err == gorm.ErrRecordNotFound {
		parentMenu = SystemMenu{
			Path:      "/cmdb",
			Name:      "CMDB",
			Component: "Layout",
			Title:     "主机管理",
			Icon:      "&#xe6a0;", // Server icon
			KeepAlive: 2,
			Status:    1,
			Level:     1,
			ParentID:  0,
			Sort:      90, // Placed after Dashboard
		}
		if err := db.Create(&parentMenu).Error; err != nil {
			zap.L().Error("failed to create CMDB parent menu", zap.Error(err))
			return err
		}
		zap.L().Info("successfully created CMDB parent menu")

		// Associate with admin role
		var adminRole SystemRole
		if err := db.First(&adminRole, 1).Error; err == nil {
			db.Model(&adminRole).Association("SystemMenus").Append(&parentMenu)
		}
	}

	// 2. Check for Provider Management sub-menu
	var providerMenu SystemMenu
	err = db.Where("name = 'Provider' AND parent_id = ?", parentMenu.ID).First(&providerMenu).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("failed to query Provider menu", zap.Error(err))
		return err
	}
	if err == gorm.ErrRecordNotFound {
		providerMenu = SystemMenu{
			Path:      "provider",
			Name:      "Provider",
			Component: "/cmdb/provider/index",
			Title:     "厂商管理",
			Icon:      "&#xe64c;", // Cloud icon
			KeepAlive: 2,
			Status:    1,
			Level:     2,
			ParentID:  parentMenu.ID,
			Sort:      1,
		}
		if err := db.Create(&providerMenu).Error; err != nil {
			zap.L().Error("failed to create Provider menu", zap.Error(err))
			return err
		}
		zap.L().Info("successfully created Provider menu")
		var adminRole SystemRole
		if err := db.First(&adminRole, 1).Error; err == nil {
			db.Model(&adminRole).Association("SystemMenus").Append(&providerMenu)
		}
	}

	// 3. Check for Host Management sub-menu
	var hostMenu SystemMenu
	err = db.Where("name = 'Host' AND parent_id = ?", parentMenu.ID).First(&hostMenu).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("failed to query Host menu", zap.Error(err))
		return err
	}
	if err == gorm.ErrRecordNotFound {
		hostMenu = SystemMenu{
			Path:      "host",
			Name:      "Host",
			Component: "/cmdb/host/index",
			Title:     "主机列表",
			Icon:      "&#xe6b3;", // Desktop icon
			KeepAlive: 2,
			Status:    1,
			Level:     2,
			ParentID:  parentMenu.ID,
			Sort:      2,
		}
		if err := db.Create(&hostMenu).Error; err != nil {
			zap.L().Error("failed to create Host menu", zap.Error(err))
			return err
		}
		zap.L().Info("successfully created Host menu")
		var adminRole SystemRole
		if err := db.First(&adminRole, 1).Error; err == nil {
			db.Model(&adminRole).Association("SystemMenus").Append(&hostMenu)
		}
	}

	zap.L().Info("CMDB menu check/update complete.")
	return nil
}

// addAnnouncementMenuIfNotExists 如果公告管理菜单不存在则添加
func addAnnouncementMenuIfNotExists(db *gorm.DB) error {
	var count int64
	db.Model(&SystemMenu{}).Where("name = ? AND title = ?", "SystemAnnouncement", "公告管理").Count(&count)
	if count > 0 {
		zap.L().Info("announcement menu already exists")
		return nil
	}

	// 添加公告管理菜单
	announcementMenu := SystemMenu{
		Path:      "announcement",
		Name:      "SystemAnnouncement",
		Component: "/system/announcement/index",
		Title:     "公告管理",
		Icon:      "&#xe747;",
		KeepAlive: 2,
		Status:    1,
		Level:     2,
		ParentID:  2, // 系统管理的子菜单
		Sort:      55,
	}

	err := db.Create(&announcementMenu).Error
	if err != nil {
		zap.L().Error("failed to create announcement menu", zap.Error(err))
		return err
	}

	// 为超级管理员分配新菜单权限
	var adminRole SystemRole
	err = db.First(&adminRole, 1).Error
	if err != nil {
		zap.L().Error("failed to find admin role", zap.Error(err))
		return err
	}

	err = db.Model(&adminRole).Association("SystemMenus").Append(&announcementMenu)
	if err != nil {
		zap.L().Error("failed to associate announcement menu with admin role", zap.Error(err))
		return err
	}

	zap.L().Info("successfully added announcement menu")
	return nil
}

// updateMenuIcons 更新菜单图标
func updateMenuIcons(db *gorm.DB) error {
	// 定义菜单图标映射
	menuIcons := map[string]string{
		"Dashboard":          "&#xe6b2;", // 仪表盘 - 显示器
		"System":             "&#xe72b;", // 系统管理 - 设置
		"SystemMenu":         "&#xe662;", // 菜单管理 - 菜单
		"SystemRole":         "&#xe734;", // 角色管理 - 个人
		"SystemDepartment":   "&#xe753;", // 部门管理 - 团队
		"SystemUser":         "&#xe608;", // 用户管理 - 用户
		"DashboardConsole":   "&#xe651;", // 工作台 - 计划工作
		"Private":            "&#xe636;", // 隐藏页面 - 帮助
		"SystemAnnouncement": "&#xe747;", // 公告管理 - 公告
		"UserCenter":         "&#xe734;", // 个人中心 - 个人
	}

	// 更新每个菜单的图标
	for name, icon := range menuIcons {
		err := db.Model(&SystemMenu{}).Where("name = ?", name).Update("icon", icon).Error
		if err != nil {
			zap.L().Error("failed to update menu icon", zap.String("menu", name), zap.Error(err))
			return err
		}
	}

	zap.L().Info("menu icons updated successfully")
	return nil
}

// resetSequences 重置序列（MySQL不需要）
func resetSequences(db *gorm.DB) error {
	zap.L().Info("skipping sequence reset for MySQL database")
	return nil
}

func Migrate(db *gorm.DB) error {
	err := migrateTable(db)
	if err != nil {
		return err
	}
	err = migrateData(db)
	if err != nil {
		return err
	}
	// 更新菜单图标
	err = updateMenuIcons(db)
	if err != nil {
		return err
	}
	// 添加序列重置操作
	err = resetSequences(db)
	if err != nil {
		return err
	}
	return nil
}
