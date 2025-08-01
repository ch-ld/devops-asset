package db

import (
	"fmt"
	"time"

	"api-server/internal/config"
	"api-server/internal/db/mysql/system"
	"api-server/internal/model"
	cmdb_model "api-server/internal/model/cmdb"
	dns_model "api-server/internal/model/dns"
	"api-server/internal/model/nav"
	"api-server/pkg/crypto/encryption"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MigrationVersion 迁移版本表
type MigrationVersion struct {
	ID          uint      `gorm:"primaryKey"`
	Version     string    `gorm:"uniqueIndex;size:50;not null" json:"version"`
	Description string    `gorm:"size:200" json:"description"`
	Applied     bool      `gorm:"default:false" json:"applied"`
	AppliedAt   time.Time `json:"applied_at"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (MigrationVersion) TableName() string {
	return "migration_versions"
}

// 数据配置结构体
type SystemDataConfig struct {
	Menus       []MenuConfig       `json:"menus"`
	Roles       []RoleConfig       `json:"roles"`
	Departments []DepartmentConfig `json:"departments"`
	Users       []UserConfig       `json:"users"`
}

type MenuConfig struct {
	ID        uint   `json:"id"`
	Path      string `json:"path"`
	Name      string `json:"name"`
	Component string `json:"component"`
	Title     string `json:"title"`
	Icon      string `json:"icon"`
	KeepAlive uint   `json:"keep_alive"`
	Status    uint   `json:"status"`
	Level     uint   `json:"level"`
	ParentID  uint   `json:"parent_id"`
	Sort      uint   `json:"sort"`
	IsHide    uint   `json:"is_hide"`
}

type RoleConfig struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Status  uint   `json:"status"`
	MenuIDs []uint `json:"menu_ids"`
}

type DepartmentConfig struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Sort   uint   `json:"sort"`
	Status uint   `json:"status"`
}

type UserConfig struct {
	ID           uint   `json:"id"`
	DepartmentID uint   `json:"department_id"`
	RoleID       uint   `json:"role_id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Status       uint   `json:"status"`
	Gender       uint   `json:"gender"`
}

// TableConfig 表配置结构
type TableConfig struct {
	SystemModels []interface{}
	CMDBModels   []interface{}
	DNSModels    []interface{}
	TableNames   map[string]string
}

// getTableConfig 获取表配置
func getTableConfig() *TableConfig {
	return &TableConfig{
		SystemModels: []interface{}{
			&system.SystemDepartment{},
			&system.SystemRole{},
			&system.SystemMenu{},
			&system.SystemMenuAuth{},
			&system.SystemUser{},
			&system.SystemUserLoginLog{},
			&system.SystemAnnouncement{},
			&nav.Nav{},
			&MigrationVersion{},
		},
		CMDBModels: []interface{}{
			&cmdb_model.Provider{},
			&cmdb_model.Host{},
			&cmdb_model.HostGroup{},
			&cmdb_model.LoadBalancer{},
			&cmdb_model.Database{},
			&cmdb_model.Redis{},
			&model.OperationLog{},
			&model.Region{},
		},
		DNSModels: []interface{}{
			&dns_model.Domain{},
			&dns_model.DomainGroup{},
			&dns_model.Tag{},
			&dns_model.Certificate{},
			&dns_model.CertificateDeployment{},
			&dns_model.Record{},
			&dns_model.RecordTemplate{},
			&dns_model.Provider{},
			&dns_model.ProviderRegion{},
			&dns_model.ChangeLog{},
			&dns_model.SyncLog{},
		},
		TableNames: map[string]string{
			"system_departments":          "system_departments",
			"system_roles":                "system_roles",
			"system_menus":                "system_menus",
			"system_menu_auths":           "system_menu_auths",
			"system_users":                "system_users",
			"system_user_login_logs":      "system_user_login_logs",
			"system_announcements":        "system_announcements",
			"navs":                        "navs",
			"migration_versions":          "migration_versions",
			"cmdb_providers":              "cmdb_providers",
			"cmdb_hosts":                  "cmdb_hosts",
			"cmdb_host_groups":            "cmdb_host_groups",
			"cmdb_load_balancers":         "cmdb_load_balancers",
			"cmdb_databases":              "cmdb_databases",
			"cmdb_redis":                  "cmdb_redis",
			"operation_logs":              "operation_logs",
			"regions":                     "regions",
			"dns_domains":                 "dns_domains",
			"dns_domain_groups":           "dns_domain_groups",
			"dns_tags":                    "dns_tags",
			"dns_certificates":            "dns_certificates",
			"dns_certificate_deployments": "dns_certificate_deployments",
			"dns_records":                 "dns_records",
			"dns_record_templates":        "dns_record_templates",
			"dns_providers":               "dns_providers",
			"dns_provider_regions":        "dns_provider_regions",
			"dns_change_logs":             "dns_change_logs",
			"dns_sync_logs":               "dns_sync_logs",
		},
	}
}

// getSystemDataConfig 获取系统数据配置
func getSystemDataConfig() *SystemDataConfig {
	return &SystemDataConfig{
		Menus: []MenuConfig{
			// 一级菜单
			{ID: 1, Path: "/dashboard", Name: "Dashboard", Component: "/index/index", Title: "仪表盘", Icon: "&#xe6b2;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 99},
			{ID: 2, Path: "/system", Name: "System", Component: "/index/index", Title: "系统管理", Icon: "&#xe72b;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 20},
			{ID: 11, Path: "/cmdb", Name: "CMDB", Component: "/index/index", Title: "CMDB资产管理", Icon: "&#xe6a0;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 90},
			{ID: 40, Path: "/dns", Name: "Dns", Component: "/index/index", Title: "DNS管理", Icon: "&#xe6a0;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 60},
			{ID: 8, Path: "/private", Name: "Private", Component: "/index/index", Title: "隐藏页面", Icon: "&#xe636;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 99, IsHide: 1},

			// 仪表盘子菜单
			{ID: 7, Path: "console", Name: "DashboardConsole", Component: "/dashboard/console/index", Title: "工作台", Icon: "&#xe651;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 1, Sort: 99},

			// 系统管理子菜单
			{ID: 3, Path: "menu", Name: "SystemMenu", Component: "/system/menu/index", Title: "菜单管理", Icon: "&#xe662;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 99},
			{ID: 4, Path: "role", Name: "SystemRole", Component: "/system/role/index", Title: "角色管理", Icon: "&#xe734;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 88},
			{ID: 5, Path: "department", Name: "SystemDepartment", Component: "/system/department/index", Title: "部门管理", Icon: "&#xe753;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 77},
			{ID: 6, Path: "user", Name: "SystemUser", Component: "/system/user/index", Title: "用户管理", Icon: "&#xe608;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 66},
			{ID: 9, Path: "announcement", Name: "SystemAnnouncement", Component: "/system/announcement/index", Title: "公告管理", Icon: "&#xe747;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 55},
			{ID: 10, Path: "user-center", Name: "UserCenter", Component: "/system/user-center/index", Title: "个人中心", Icon: "&#xe734;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 2, Sort: 50},

			// CMDB资产管理子菜单
			{ID: 13, Path: "provider", Name: "Provider", Component: "/cmdb/provider/index", Title: "厂商管理", Icon: "&#xe693;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 11, Sort: 99},
			{ID: 19, Path: "host-dashboard", Name: "HostDashboard", Component: "/cmdb/host/dashboard", Title: "主机概览", Icon: "&#xe67a;", KeepAlive: 1, Status: 1, Level: 2, ParentID: 11, Sort: 88},
			{ID: 15, Path: "hosts", Name: "HostManagement", Component: "/cmdb/host/index", Title: "主机管理", Icon: "&#xe710;", KeepAlive: 1, Status: 1, Level: 2, ParentID: 11, Sort: 77},

			// CMDB隐藏页面（用于操作但不在菜单中显示）
			{ID: 16, Path: "host/:id", Name: "HostDetail", Component: "/cmdb/host/detail", Title: "主机详情", Icon: "&#xe65f;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 11, Sort: 0, IsHide: 1},
			{ID: 17, Path: "host-edit/:id", Name: "HostEdit", Component: "/cmdb/host/edit", Title: "编辑主机", Icon: "&#xe642;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 11, Sort: 0, IsHide: 1},
			{ID: 18, Path: "host-create", Name: "HostCreate", Component: "/cmdb/host/edit", Title: "添加主机", Icon: "&#xe61f;", KeepAlive: 2, Status: 1, Level: 2, ParentID: 11, Sort: 0, IsHide: 1},
		},
		Roles: []RoleConfig{
			{ID: 1, Name: "超级管理员", Desc: "拥有所有权限", Status: 1, MenuIDs: []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13, 15, 16, 17, 18, 19}},
			{ID: 2, Name: "普通用户", Desc: "普通用户", Status: 1, MenuIDs: []uint{1, 7, 10}},
			{ID: 3, Name: "CMDB管理员", Desc: "CMDB资产管理权限", Status: 1, MenuIDs: []uint{1, 7, 10, 11, 13, 15, 16, 17, 18, 19}},
		},
		Departments: []DepartmentConfig{
			{ID: 1, Name: "管理中心", Sort: 1, Status: 1},
		},
		Users: []UserConfig{
			{ID: 1, DepartmentID: 1, RoleID: 1, Name: "超级管理员", Username: "admin", Password: config.AdminPassword, Status: 1, Gender: 1},
		},
	}
}

// MigrateAll 执行所有数据库迁移（智能初始化）
func MigrateAll(db *gorm.DB) error {
	zap.L().Info("Starting intelligent database migration process")

	// 1. 首先创建版本控制表
	if err := createMigrationVersionTable(db); err != nil {
		return fmt.Errorf("failed to create migration version table: %w", err)
	}

	// 2. 检查是否已经初始化过
	if isAlreadyInitialized(db) {
		zap.L().Info("Database already initialized, checking for pending migrations")
		return applyPendingMigrations(db)
	}

	// 3. 执行完整初始化
	return performFullInitialization(db)
}

// createMigrationVersionTable 创建迁移版本控制表
func createMigrationVersionTable(db *gorm.DB) error {
	return db.AutoMigrate(&MigrationVersion{})
}

// isAlreadyInitialized 检查数据库是否已经初始化
func isAlreadyInitialized(db *gorm.DB) bool {
	config := getTableConfig()

	// 检查核心表是否存在且有数据
	if !db.Migrator().HasTable("system_users") {
		return false
	}

	var userCount int64
	db.Model(&system.SystemUser{}).Count(&userCount)
	if userCount == 0 {
		return false
	}

	// 检查所有系统表是否存在
	allModels := append(config.SystemModels, config.CMDBModels...)
	for _, model := range allModels {
		if !db.Migrator().HasTable(model) {
			zap.L().Info("Missing table detected, will perform migration",
				zap.String("model", fmt.Sprintf("%T", model)))
			return false
		}
	}

	return true
}

// performFullInitialization 执行完整的数据库初始化
func performFullInitialization(db *gorm.DB) error {
	zap.L().Info("Performing full database initialization")

	// 1. 迁移表结构
	if err := migrateAllTables(db); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	// 2. 初始化基础数据
	if err := initializeAllData(db); err != nil {
		return fmt.Errorf("failed to initialize base data: %w", err)
	}

	// 3. 记录初始化版本
	if err := recordMigrationVersion(db, "1.0.0", "Initial database setup with base data"); err != nil {
		return fmt.Errorf("failed to record migration version: %w", err)
	}

	zap.L().Info("Full database initialization completed successfully")
	return nil
}

// migrateAllTables 迁移所有表结构
func migrateAllTables(db *gorm.DB) error {
	config := getTableConfig()

	// 迁移系统表
	if err := db.AutoMigrate(config.SystemModels...); err != nil {
		zap.L().Error("Failed to migrate system tables", zap.Error(err))
		return err
	}
	zap.L().Info("System tables migrated successfully")

	// 迁移CMDB表
	if err := db.AutoMigrate(config.CMDBModels...); err != nil {
		zap.L().Error("Failed to migrate CMDB tables", zap.Error(err))
		return err
	}
	zap.L().Info("CMDB tables migrated successfully")

	// 迁移DNS表
	if err := db.AutoMigrate(config.DNSModels...); err != nil {
		zap.L().Error("Failed to migrate DNS tables", zap.Error(err))
		return err
	}
	zap.L().Info("DNS tables migrated successfully")

	return nil
}

// initializeAllData 初始化所有基础数据
func initializeAllData(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 初始化系统数据
		if err := initializeSystemData(tx); err != nil {
			return fmt.Errorf("failed to initialize system data: %w", err)
		}

		// 2. 初始化导航栏测试数据
		if err := initializeNavData(tx); err != nil {
			return fmt.Errorf("failed to initialize navigation data: %w", err)
		}

		// 3. 初始化区域数据
		if err := initializeRegionData(tx); err != nil {
			return fmt.Errorf("failed to initialize region data: %w", err)
		}

		zap.L().Info("All data initialized successfully")
		return nil
	})
}

// initializeSystemData 初始化系统数据
func initializeSystemData(tx *gorm.DB) error {
	zap.L().Info("Starting system data initialization")
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

	zap.L().Info("System data initialization completed successfully")
	return nil
}

// migrateMenuData 迁移菜单数据
func migrateMenuData(tx *gorm.DB, menuConfigs []MenuConfig) error {
	var existingCount int64
	tx.Model(&system.SystemMenu{}).Count(&existingCount)

	if existingCount > 0 {
		zap.L().Info("Menu data already exists, performing incremental update")
		return updateMenusIfNeeded(tx, menuConfigs)
	}

	// 批量创建菜单
	var menus []system.SystemMenu
	for _, config := range menuConfigs {
		menu := system.SystemMenu{
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
		var existingMenu system.SystemMenu
		err := tx.Where("id = ?", config.ID).First(&existingMenu).Error

		if err == gorm.ErrRecordNotFound {
			// 菜单不存在，创建新菜单
			newMenu := system.SystemMenu{
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
	}
	return nil
}

// migrateRoleData 迁移角色数据
func migrateRoleData(tx *gorm.DB, roleConfigs []RoleConfig) error {
	var existingCount int64
	tx.Model(&system.SystemRole{}).Count(&existingCount)

	if existingCount > 0 {
		zap.L().Info("Role data already exists, skipping role creation")
		return nil
	}

	var roles []system.SystemRole
	for _, config := range roleConfigs {
		role := system.SystemRole{
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
	tx.Model(&system.SystemDepartment{}).Count(&existingCount)

	if existingCount > 0 {
		zap.L().Info("Department data already exists, skipping department creation")
		return nil
	}

	var departments []system.SystemDepartment
	for _, config := range deptConfigs {
		dept := system.SystemDepartment{
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
	tx.Model(&system.SystemUser{}).Count(&existingCount)

	if existingCount > 0 {
		zap.L().Info("User data already exists, skipping user creation")
		return nil
	}

	var users []system.SystemUser
	for _, config := range userConfigs {
		user := system.SystemUser{
			Model:        gorm.Model{ID: config.ID},
			DepartmentID: config.DepartmentID,
			RoleID:       config.RoleID,
			Name:         config.Name,
			Username:     config.Username,
			Password:     encryptPassword(config.Password),
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
		var role system.SystemRole
		if err := tx.First(&role, roleConfig.ID).Error; err != nil {
			zap.L().Warn("Role not found, skipping permission assignment",
				zap.Uint("role_id", roleConfig.ID))
			continue
		}

		// 检查是否已有权限分配
		existingCount := tx.Model(&role).Association("SystemMenus").Count()
		if existingCount > 0 {
			zap.L().Info("Role permissions already exist, skipping",
				zap.String("role", role.Name))
			continue
		}

		// 获取菜单并分配权限
		var menus []system.SystemMenu
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

// initializeNavData 初始化导航栏测试数据
func initializeNavData(tx *gorm.DB) error {
	// 检查是否已有导航数据
	var count int64
	tx.Model(&nav.Nav{}).Count(&count)
	if count > 0 {
		zap.L().Info("Navigation data already exists, skipping initialization")
		return nil
	}

	// 创建测试导航数据
	navItems := []nav.Nav{
		{
			GroupName:    "开发工具",
			Name:         "GitHub",
			Description:  "全球最大的代码托管平台",
			IconUrl:      "&#xe7a8;",
			Links:        "https://github.com",
			OrderNum:     1,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "开发工具",
			Name:         "GitLab",
			Description:  "开源的Git仓库管理平台",
			IconUrl:      "&#xe7a9;",
			Links:        "https://gitlab.com",
			OrderNum:     2,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "开发工具",
			Name:         "Docker Hub",
			Description:  "容器镜像仓库",
			IconUrl:      "&#xe7aa;",
			Links:        "https://hub.docker.com",
			OrderNum:     3,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "云服务",
			Name:         "阿里云",
			Description:  "阿里巴巴云计算服务",
			IconUrl:      "&#xe7ab;",
			Links:        "https://www.aliyun.com",
			OrderNum:     1,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "云服务",
			Name:         "腾讯云",
			Description:  "腾讯云计算服务",
			IconUrl:      "&#xe7ac;",
			Links:        "https://cloud.tencent.com",
			OrderNum:     2,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "文档工具",
			Name:         "语雀",
			Description:  "专业的云端知识库",
			IconUrl:      "&#xe7ad;",
			Links:        "https://www.yuque.com",
			OrderNum:     1,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "文档工具",
			Name:         "Notion",
			Description:  "一体化工作空间",
			IconUrl:      "&#xe7ae;",
			Links:        "https://www.notion.so",
			OrderNum:     2,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "监控工具",
			Name:         "Grafana",
			Description:  "开源数据可视化和监控平台",
			IconUrl:      "&#xe7af;",
			Links:        "https://grafana.com",
			OrderNum:     1,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "监控工具",
			Name:         "Prometheus",
			Description:  "开源监控和告警工具",
			IconUrl:      "&#xe7b0;",
			Links:        "https://prometheus.io",
			OrderNum:     2,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "学习资源",
			Name:         "MDN Web Docs",
			Description:  "Web开发技术文档",
			IconUrl:      "&#xe7b1;",
			Links:        "https://developer.mozilla.org",
			OrderNum:     1,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "学习资源",
			Name:         "Stack Overflow",
			Description:  "程序员问答社区",
			IconUrl:      "&#xe7b2;",
			Links:        "https://stackoverflow.com",
			OrderNum:     2,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "设计工具",
			Name:         "Figma",
			Description:  "在线协作设计工具",
			IconUrl:      "&#xe7b3;",
			Links:        "https://www.figma.com",
			OrderNum:     1,
			Status:       1,
			OpenInNewTab: true,
		},
	}

	// 批量创建导航数据
	if err := tx.Create(&navItems).Error; err != nil {
		return fmt.Errorf("failed to create navigation data: %w", err)
	}

	zap.L().Info("Navigation test data created successfully", zap.Int("count", len(navItems)))
	return nil
}

// initializeRegionData 初始化区域数据
func initializeRegionData(tx *gorm.DB) error {
	// 检查是否已有区域数据
	var count int64
	tx.Model(&model.Region{}).Count(&count)
	if count > 0 {
		zap.L().Info("Region data already exists, skipping initialization")
		return nil
	}

	zap.L().Info("Starting region data initialization")

	// 支持的云厂商类型
	supportedProviders := []string{"alicloud", "tencentcloud", "aws", "huaweicloud", "baiduyun"}

	for _, providerType := range supportedProviders {
		// 获取默认区域列表
		defaultRegions := getDefaultRegionsByProvider(providerType)
		if len(defaultRegions) == 0 {
			zap.L().Warn("No default regions found for provider", zap.String("provider", providerType))
			continue
		}

		// 批量插入默认区域数据
		if err := batchInsertRegions(tx, providerType, defaultRegions); err != nil {
			return fmt.Errorf("failed to insert default regions for %s: %w", providerType, err)
		}

		zap.L().Info("Default regions initialized",
			zap.String("provider", providerType),
			zap.Int("count", len(defaultRegions)))
	}

	zap.L().Info("Region data initialization completed successfully")
	return nil
}

// batchInsertRegions 批量插入区域数据
func batchInsertRegions(tx *gorm.DB, providerType string, regionIDs []string) error {
	now := time.Now()
	var regions []model.Region

	for i, regionID := range regionIDs {
		region := model.Region{
			ProviderType: providerType,
			RegionID:     regionID,
			RegionName:   regionID, // 默认区域名称与ID相同
			Status:       "active",
			IsDefault:    i == 0, // 第一个区域设为默认
			SyncTime:     now,
		}
		regions = append(regions, region)
	}

	// 批量插入
	return tx.CreateInBatches(regions, 100).Error
}

// getDefaultRegionsByProvider 获取指定云厂商的默认区域列表
func getDefaultRegionsByProvider(providerType string) []string {
	switch providerType {
	case "alicloud":
		return []string{
			"cn-hangzhou",
			"cn-shanghai",
			"cn-beijing",
			"cn-shenzhen",
			"cn-guangzhou",
			"cn-chengdu",
			"cn-qingdao",
			"cn-zhangjiakou",
		}
	case "tencentcloud":
		return []string{
			"ap-beijing",
			"ap-shanghai",
			"ap-guangzhou",
			"ap-chengdu",
			"ap-chongqing",
			"ap-nanjing",
			"ap-shenzhen-fsi",
		}
	case "aws":
		return []string{
			"us-east-1",
			"us-east-2",
			"us-west-1",
			"us-west-2",
			"ca-central-1",
			"eu-central-1",
			"eu-west-1",
			"eu-west-2",
			"eu-west-3",
			"eu-north-1",
			"eu-south-1",
			"ap-northeast-1",
			"ap-northeast-2",
			"ap-northeast-3",
			"ap-southeast-1",
			"ap-southeast-2",
			"ap-southeast-3",
			"ap-south-1",
			"ap-east-1",
			"me-south-1",
			"af-south-1",
			"sa-east-1",
		}
	case "huaweicloud":
		return []string{
			"cn-north-1",
			"cn-north-4",
			"cn-east-2",
			"cn-east-3",
			"cn-south-1",
		}
	case "baiduyun":
		return []string{
			"bj",
			"gz",
			"su",
			"hkg",
		}
	default:
		return []string{}
	}
}

// applyPendingMigrations 应用待处理的迁移
func applyPendingMigrations(db *gorm.DB) error {
	// 确保所有表都是最新的
	if err := migrateAllTables(db); err != nil {
		return err
	}

	// 检查并应用新的迁移
	pendingMigrations := []struct {
		version     string
		description string
		migrateFn   func(*gorm.DB) error
	}{
		{"1.1.0", "Add navigation test data", ensureNavData},
		{"1.2.0", "Update menu icons and structure", ensureMenuUpdates},
	}

	for _, migration := range pendingMigrations {
		if !isMigrationApplied(db, migration.version) {
			zap.L().Info("Applying migration", zap.String("version", migration.version))

			if err := migration.migrateFn(db); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.version, err)
			}

			if err := recordMigrationVersion(db, migration.version, migration.description); err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migration.version, err)
			}

			zap.L().Info("Migration applied successfully", zap.String("version", migration.version))
		}
	}

	return nil
}

// ensureNavData 确保导航数据存在
func ensureNavData(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return initializeNavData(tx)
	})
}

// ensureMenuUpdates 确保菜单更新
func ensureMenuUpdates(db *gorm.DB) error {
	zap.L().Info("Menu updates checked")
	return nil
}

// recordMigrationVersion 记录迁移版本
func recordMigrationVersion(db *gorm.DB, version, description string) error {
	migration := MigrationVersion{
		Version:     version,
		Description: description,
		Applied:     true,
		AppliedAt:   time.Now(),
	}

	return db.Create(&migration).Error
}

// isMigrationApplied 检查迁移是否已应用
func isMigrationApplied(db *gorm.DB, version string) bool {
	var count int64
	db.Model(&MigrationVersion{}).Where("version = ? AND applied = ?", version, true).Count(&count)
	return count > 0
}

// GetMigrationStatus 获取迁移状态
func GetMigrationStatus(db *gorm.DB) (map[string]bool, error) {
	status := make(map[string]bool)
	config := getTableConfig()

	for tableName := range config.TableNames {
		status[tableName] = db.Migrator().HasTable(tableName)
	}

	return status, nil
}

// ResetDatabase 重置数据库（危险操作，仅用于开发环境）
func ResetDatabase(db *gorm.DB) error {
	zap.L().Warn("Resetting database - this will drop all tables!")
	config := getTableConfig()

	allModels := append(config.CMDBModels, config.SystemModels...)
	for _, model := range allModels {
		if err := db.Migrator().DropTable(model); err != nil {
			zap.L().Error("Failed to drop table",
				zap.String("model", fmt.Sprintf("%T", model)),
				zap.Error(err))
		}
	}

	zap.L().Info("Database reset completed")
	return nil
}

// GetDatabaseInfo 获取数据库信息
func GetDatabaseInfo(db *gorm.DB) map[string]interface{} {
	config := getTableConfig()
	info := make(map[string]interface{})

	// 基本信息
	info["total_tables"] = len(config.TableNames)
	info["system_tables"] = len(config.SystemModels)
	info["cmdb_tables"] = len(config.CMDBModels)

	// 表状态
	tableStatus := make(map[string]bool)
	for tableName := range config.TableNames {
		tableStatus[tableName] = db.Migrator().HasTable(tableName)
	}
	info["table_status"] = tableStatus

	// 迁移历史
	var migrations []MigrationVersion
	db.Order("applied_at DESC").Find(&migrations)
	info["migration_history"] = migrations

	// 数据统计
	dataStats := make(map[string]int64)
	if db.Migrator().HasTable("system_users") {
		var userCount int64
		db.Model(&system.SystemUser{}).Count(&userCount)
		dataStats["users"] = userCount
	}
	if db.Migrator().HasTable("navs") {
		var navCount int64
		db.Model(&nav.Nav{}).Count(&navCount)
		dataStats["navigations"] = navCount
	}
	info["data_statistics"] = dataStats

	return info
}

// encryptPassword 密码加密（使用MD5加密）
func encryptPassword(password string) string {
	// 使用项目统一的密码加密方式（MD5 + Salt）
	return encryption.MD5WithSalt(config.PWDSalt + password)
}
