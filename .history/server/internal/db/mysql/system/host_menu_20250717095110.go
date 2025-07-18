package system

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AddHostMenusIfNotExist 如果主机相关菜单不存在则添加
func AddHostMenusIfNotExist(db *gorm.DB) error {
	// 1. 检查CMDB父菜单是否存在
	var parentMenu SystemMenu
	err := db.Where("name = 'CMDB'").First(&parentMenu).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("failed to query CMDB menu", zap.Error(err))
		return err
	}

	// 如果父菜单不存在，创建它
	if err == gorm.ErrRecordNotFound {
		parentMenu = SystemMenu{
			Path:      "/cmdb",
			Name:      "CMDB",
			Component: "/index/index",
			Title:     "CMDB资产管理",
			Icon:      "&#xe6a0;", // 服务器图标
			KeepAlive: 2,
			Status:    1,
			Level:     1,
			ParentID:  0,
			Sort:      90, // 排在仪表盘后面
		}
		if err := db.Create(&parentMenu).Error; err != nil {
			zap.L().Error("failed to create CMDB parent menu", zap.Error(err))
			return err
		}
		zap.L().Info("successfully created CMDB parent menu")

		// 为管理员角色关联此菜单
		var adminRole SystemRole
		if err := db.First(&adminRole, 1).Error; err == nil {
			db.Model(&adminRole).Association("SystemMenus").Append(&parentMenu)
		}
	}

	// 2. 添加主机管理子菜单
	err = addHostManagementMenuIfNotExists(db, parentMenu.ID)
	if err != nil {
		return err
	}

	// 3. 添加主机详情子菜单
	err = addHostDetailMenuIfNotExists(db, parentMenu.ID)
	if err != nil {
		return err
	}

	// 4. 添加主机编辑子菜单
	err = addHostEditMenuIfNotExists(db, parentMenu.ID)
	if err != nil {
		return err
	}

	// 5. 添加主机创建子菜单
	err = addHostCreateMenuIfNotExists(db, parentMenu.ID)
	if err != nil {
		return err
	}

	// 6. 添加主机仪表盘子菜单
	err = addHostDashboardMenuIfNotExists(db, parentMenu.ID)
	if err != nil {
		return err
	}

	// 7. 添加云账号管理子菜单
	err = addProviderManagementMenuIfNotExists(db, parentMenu.ID)
	if err != nil {
		return err
	}

	zap.L().Info("host menus check/update complete")
	return nil
}

// 添加主机管理菜单
func addHostManagementMenuIfNotExists(db *gorm.DB, parentID uint) error {
	var menu SystemMenu
	err := db.Where("name = 'HostManagement' AND parent_id = ?", parentID).First(&menu).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("failed to query host management menu", zap.Error(err))
		return err
	}

	if err == gorm.ErrRecordNotFound {
		menu = SystemMenu{
			Path:      "hosts",
			Name:      "HostManagement",
			Component: "/cmdb/host/simple-index",
			Title:     "主机管理",
			Icon:      "&#xe73e;",
			KeepAlive: 1, // 保持缓存
			Status:    1,
			Level:     2,
			ParentID:  parentID,
			Sort:      99,
		}
		if err := db.Create(&menu).Error; err != nil {
			zap.L().Error("failed to create host management menu", zap.Error(err))
			return err
		}
		zap.L().Info("successfully created host management menu")

		// 为管理员角色关联此菜单
		var adminRole SystemRole
		if err := db.First(&adminRole, 1).Error; err == nil {
			db.Model(&adminRole).Association("SystemMenus").Append(&menu)
		}
	}
	return nil
}

// 添加主机详情菜单
func addHostDetailMenuIfNotExists(db *gorm.DB, parentID uint) error {
	var menu SystemMenu
	err := db.Where("name = 'HostDetail' AND parent_id = ?", parentID).First(&menu).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("failed to query host detail menu", zap.Error(err))
		return err
	}

	if err == gorm.ErrRecordNotFound {
		menu = SystemMenu{
			Path:      "host/:id",
			Name:      "HostDetail",
			Component: "/cmdb/host/detail",
			Title:     "主机详情",
			Icon:      "&#xe74c;",
			IsHide:    1, // 隐藏菜单
			KeepAlive: 2,
			Status:    1,
			Level:     2,
			ParentID:  parentID,
			Sort:      0,
		}
		if err := db.Create(&menu).Error; err != nil {
			zap.L().Error("failed to create host detail menu", zap.Error(err))
			return err
		}
		zap.L().Info("successfully created host detail menu")

		// 为管理员角色关联此菜单
		var adminRole SystemRole
		if err := db.First(&adminRole, 1).Error; err == nil {
			db.Model(&adminRole).Association("SystemMenus").Append(&menu)
		}
	}
	return nil
}

// 添加主机编辑菜单
func addHostEditMenuIfNotExists(db *gorm.DB, parentID uint) error {
	var menu SystemMenu
	err := db.Where("name = 'HostEdit' AND parent_id = ?", parentID).First(&menu).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("failed to query host edit menu", zap.Error(err))
		return err
	}

	if err == gorm.ErrRecordNotFound {
		menu = SystemMenu{
			Path:      "host-edit/:id",
			Name:      "HostEdit",
			Component: "/cmdb/host/edit",
			Title:     "编辑主机",
			Icon:      "&#xe66e;",
			IsHide:    1, // 隐藏菜单
			KeepAlive: 2,
			Status:    1,
			Level:     2,
			ParentID:  parentID,
			Sort:      0,
		}
		if err := db.Create(&menu).Error; err != nil {
			zap.L().Error("failed to create host edit menu", zap.Error(err))
			return err
		}
		zap.L().Info("successfully created host edit menu")

		// 为管理员角色关联此菜单
		var adminRole SystemRole
		if err := db.First(&adminRole, 1).Error; err == nil {
			db.Model(&adminRole).Association("SystemMenus").Append(&menu)
		}
	}
	return nil
}

// 添加主机创建菜单
func addHostCreateMenuIfNotExists(db *gorm.DB, parentID uint) error {
	var menu SystemMenu
	err := db.Where("name = 'HostCreate' AND parent_id = ?", parentID).First(&menu).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("failed to query host create menu", zap.Error(err))
		return err
	}

	if err == gorm.ErrRecordNotFound {
		menu = SystemMenu{
			Path:      "host-create",
			Name:      "HostCreate",
			Component: "/cmdb/host/edit",
			Title:     "添加主机",
			Icon:      "&#xe727;",
			IsHide:    1, // 隐藏菜单
			KeepAlive: 2,
			Status:    1,
			Level:     2,
			ParentID:  parentID,
			Sort:      0,
		}
		if err := db.Create(&menu).Error; err != nil {
			zap.L().Error("failed to create host create menu", zap.Error(err))
			return err
		}
		zap.L().Info("successfully created host create menu")

		// 为管理员角色关联此菜单
		var adminRole SystemRole
		if err := db.First(&adminRole, 1).Error; err == nil {
			db.Model(&adminRole).Association("SystemMenus").Append(&menu)
		}
	}
	return nil
}

// 添加主机仪表盘菜单
func addHostDashboardMenuIfNotExists(db *gorm.DB, parentID uint) error {
	var menu SystemMenu
	err := db.Where("name = 'HostDashboard' AND parent_id = ?", parentID).First(&menu).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("failed to query host dashboard menu", zap.Error(err))
		return err
	}

	if err == gorm.ErrRecordNotFound {
		menu = SystemMenu{
			Path:      "host-dashboard",
			Name:      "HostDashboard",
			Component: "/cmdb/host/dashboard",
			Title:     "主机概览",
			Icon:      "&#xe64c;",
			KeepAlive: 1, // 保持缓存
			Status:    1,
			Level:     2,
			ParentID:  parentID,
			Sort:      88,
		}
		if err := db.Create(&menu).Error; err != nil {
			zap.L().Error("failed to create host dashboard menu", zap.Error(err))
			return err
		}
		zap.L().Info("successfully created host dashboard menu")

		// 为管理员角色关联此菜单
		var adminRole SystemRole
		if err := db.First(&adminRole, 1).Error; err == nil {
			db.Model(&adminRole).Association("SystemMenus").Append(&menu)
		}
	}
	return nil
}

// 添加云账号管理菜单
func addProviderManagementMenuIfNotExists(db *gorm.DB, parentID uint) error {
	var menu SystemMenu
	err := db.Where("name = 'ProviderManagement' AND parent_id = ?", parentID).First(&menu).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("failed to query provider management menu", zap.Error(err))
		return err
	}

	if err == gorm.ErrRecordNotFound {
		menu = SystemMenu{
			Path:      "providers",
			Name:      "ProviderManagement",
			Component: "/cmdb/provider/index",
			Title:     "云账号管理",
			Icon:      "&#xe67c;",
			KeepAlive: 1, // 保持缓存
			Status:    1,
			Level:     2,
			ParentID:  parentID,
			Sort:      77,
		}
		if err := db.Create(&menu).Error; err != nil {
			zap.L().Error("failed to create provider management menu", zap.Error(err))
			return err
		}
		zap.L().Info("successfully created provider management menu")

		// 为管理员角色关联此菜单
		var adminRole SystemRole
		if err := db.First(&adminRole, 1).Error; err == nil {
			db.Model(&adminRole).Association("SystemMenus").Append(&menu)
		}
	}
	return nil
}
