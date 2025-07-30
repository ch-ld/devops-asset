package db

import (
	"fmt"
	"time"

	"api-server/internal/db/mysql/system"
	cmdb_model "api-server/internal/model/cmdb"
	"api-server/internal/model/nav"

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

// TableConfig 表配置结构
type TableConfig struct {
	SystemModels []interface{}
	CMDBModels   []interface{}
	TableNames   map[string]string // 模型名 -> 表名映射
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
			&MigrationVersion{}, // 添加版本控制表
		},
		CMDBModels: []interface{}{
			&cmdb_model.Provider{},
			&cmdb_model.Host{},
			&cmdb_model.HostGroup{},
			&cmdb_model.LoadBalancer{},
			&cmdb_model.Database{},
			&cmdb_model.Redis{},
		},
		TableNames: map[string]string{
			"system_departments":    "system_departments",
			"system_roles":          "system_roles", 
			"system_menus":          "system_menus",
			"system_menu_auths":     "system_menu_auths",
			"system_users":          "system_users",
			"system_user_login_logs": "system_user_login_logs",
			"system_announcements":  "system_announcements",
			"navs":                  "navs",
			"migration_versions":    "migration_versions",
			"cmdb_providers":        "cmdb_providers",
			"cmdb_hosts":            "cmdb_hosts", 
			"cmdb_host_groups":      "cmdb_host_groups",
			"cmdb_load_balancers":   "cmdb_load_balancers",
			"cmdb_databases":        "cmdb_databases",
			"cmdb_redis":            "cmdb_redis",
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
	if err := initializeBaseData(db); err != nil {
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
	if err := cmdb_model.Migrate(db); err != nil {
		zap.L().Error("Failed to migrate CMDB tables", zap.Error(err))
		return err
	}
	zap.L().Info("CMDB tables migrated successfully")

	return nil
}

// initializeBaseData 初始化基础数据
func initializeBaseData(db *gorm.DB) error {
	// 1. 初始化系统数据
	if err := system.MigrateData(db); err != nil {
		zap.L().Error("Failed to initialize system data", zap.Error(err))
		return err
	}
	zap.L().Info("System data initialized successfully")

	// 2. 初始化导航栏测试数据
	if err := initializeNavData(db); err != nil {
		zap.L().Error("Failed to initialize navigation data", zap.Error(err))
		return err
	}
	zap.L().Info("Navigation data initialized successfully")

	return nil
}

// initializeNavData 初始化导航栏测试数据
func initializeNavData(db *gorm.DB) error {
	// 检查是否已有导航数据
	var count int64
	db.Model(&nav.Nav{}).Count(&count)
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
			IconUrl:      "&#xe7a8;", // GitHub图标
			Links:        "https://github.com",
			OrderNum:     1,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "开发工具", 
			Name:         "GitLab",
			Description:  "开源的Git仓库管理平台",
			IconUrl:      "&#xe7a9;", // GitLab图标
			Links:        "https://gitlab.com",
			OrderNum:     2,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "开发工具",
			Name:         "Docker Hub",
			Description:  "容器镜像仓库",
			IconUrl:      "&#xe7aa;", // Docker图标
			Links:        "https://hub.docker.com",
			OrderNum:     3,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "云服务",
			Name:         "阿里云",
			Description:  "阿里巴巴云计算服务",
			IconUrl:      "&#xe7ab;", // 云图标
			Links:        "https://www.aliyun.com",
			OrderNum:     1,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "云服务",
			Name:         "腾讯云", 
			Description:  "腾讯云计算服务",
			IconUrl:      "&#xe7ac;", // 云图标
			Links:        "https://cloud.tencent.com",
			OrderNum:     2,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "文档工具",
			Name:         "语雀",
			Description:  "专业的云端知识库",
			IconUrl:      "&#xe7ad;", // 文档图标
			Links:        "https://www.yuque.com",
			OrderNum:     1,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "文档工具",
			Name:         "Notion",
			Description:  "一体化工作空间",
			IconUrl:      "&#xe7ae;", // Notion图标
			Links:        "https://www.notion.so",
			OrderNum:     2,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "监控工具",
			Name:         "Grafana",
			Description:  "开源数据可视化和监控平台",
			IconUrl:      "&#xe7af;", // 图表图标
			Links:        "https://grafana.com",
			OrderNum:     1,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "监控工具",
			Name:         "Prometheus",
			Description:  "开源监控和告警工具",
			IconUrl:      "&#xe7b0;", // 监控图标
			Links:        "https://prometheus.io",
			OrderNum:     2,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "学习资源",
			Name:         "MDN Web Docs",
			Description:  "Web开发技术文档",
			IconUrl:      "&#xe7b1;", // 书籍图标
			Links:        "https://developer.mozilla.org",
			OrderNum:     1,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "学习资源",
			Name:         "Stack Overflow",
			Description:  "程序员问答社区",
			IconUrl:      "&#xe7b2;", // 问答图标
			Links:        "https://stackoverflow.com",
			OrderNum:     2,
			Status:       1,
			OpenInNewTab: true,
		},
		{
			GroupName:    "设计工具",
			Name:         "Figma",
			Description:  "在线协作设计工具",
			IconUrl:      "&#xe7b3;", // 设计图标
			Links:        "https://www.figma.com",
			OrderNum:     1,
			Status:       1,
			OpenInNewTab: true,
		},
	}

	// 批量创建导航数据
	if err := db.Create(&navItems).Error; err != nil {
		zap.L().Error("Failed to create navigation data", zap.Error(err))
		return err
	}

	zap.L().Info("Navigation test data created successfully", zap.Int("count", len(navItems)))
	return nil
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
	return initializeNavData(db)
}

// ensureMenuUpdates 确保菜单更新
func ensureMenuUpdates(db *gorm.DB) error {
	// 这里可以添加菜单更新逻辑
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

// GetMigrationStatus 获取迁移状态（优化版）
func GetMigrationStatus(db *gorm.DB) (map[string]bool, error) {
	status := make(map[string]bool)
	config := getTableConfig()
	
	// 检查所有表
	for tableName := range config.TableNames {
		status[tableName] = db.Migrator().HasTable(tableName)
	}
	
	return status, nil
}

// ResetDatabase 重置数据库（危险操作，仅用于开发环境）
func ResetDatabase(db *gorm.DB) error {
	zap.L().Warn("Resetting database - this will drop all tables!")
	config := getTableConfig()
	
	// 删除所有表
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
		db.Model(&system.SystemUser{}).Count(&dataStats["users"])
	}
	if db.Migrator().HasTable("navs") {
		db.Model(&nav.Nav{}).Count(&dataStats["navigations"]) 
	}
	info["data_statistics"] = dataStats
	
	return info
} 
