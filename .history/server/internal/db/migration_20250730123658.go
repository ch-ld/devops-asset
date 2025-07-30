package db

import (
	"api-server/internal/db/mysql/system"
	cmdb_model "api-server/internal/model/cmdb"
	"api-server/internal/model/nav"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MigrateAll 执行所有数据库迁移
// 这是统一的数据库迁移入口点，整合了系统模块和CMDB模块的所有迁移
func MigrateAll(db *gorm.DB) error {
	zap.L().Info("Starting database migration process")

	// 1. 执行系统模块表结构迁移
	if err := migrateSystemTables(db); err != nil {
		zap.L().Error("Failed to migrate system tables", zap.Error(err))
		return err
	}
	zap.L().Info("System tables migrated successfully")

	// 2. 执行CMDB模块表结构迁移
	if err := migrateCMDBTables(db); err != nil {
		zap.L().Error("Failed to migrate CMDB tables", zap.Error(err))
		return err
	}
	zap.L().Info("CMDB tables migrated successfully")

	// 3. 执行系统基础数据迁移
	if err := migrateSystemData(db); err != nil {
		zap.L().Error("Failed to migrate system data", zap.Error(err))
		return err
	}
	zap.L().Info("System data migrated successfully")

	zap.L().Info("All database migrations completed successfully")
	return nil
}

// migrateSystemTables 迁移系统模块表结构
func migrateSystemTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&system.SystemDepartment{},
		&system.SystemRole{},
		&system.SystemMenu{},
		&system.SystemMenuAuth{},
		&system.SystemUser{},
		&system.SystemUserLoginLog{},
		&system.SystemAnnouncement{},
		&nav.Nav{},
	)
}

// migrateCMDBTables 迁移CMDB模块表结构
func migrateCMDBTables(db *gorm.DB) error {
	return cmdb_model.Migrate(db)
}

// migrateSystemData 迁移系统基础数据
func migrateSystemData(db *gorm.DB) error {
	return system.MigrateData(db)
}

// GetMigrationStatus 获取迁移状态
func GetMigrationStatus(db *gorm.DB) (map[string]bool, error) {
	status := make(map[string]bool)
	
	// 检查系统表是否存在
	systemTables := []string{
		"system_departments", "system_roles", "system_menus", 
		"system_menu_auths", "system_users", "system_user_login_logs", 
		"system_announcements", "navs",
	}
	
	for _, table := range systemTables {
		status["system_"+table] = db.Migrator().HasTable(table)
	}
	
	// 检查CMDB表是否存在
	cmdbTables := []string{
		"providers", "hosts", "host_groups", "load_balancers", 
		"databases", "redis", "operation_logs",
	}
	
	for _, table := range cmdbTables {
		status["cmdb_"+table] = db.Migrator().HasTable(table)
	}
	
	return status, nil
}

// ResetDatabase 重置数据库（危险操作，仅用于开发环境）
func ResetDatabase(db *gorm.DB) error {
	zap.L().Warn("Resetting database - this will drop all tables!")
	
	// 删除CMDB表
	cmdbTables := []interface{}{
		&cmdb_model.Provider{},
		&cmdb_model.Host{},
		&cmdb_model.HostGroup{},
		&cmdb_model.LoadBalancer{},
		&cmdb_model.Database{},
		&cmdb_model.Redis{},
	}
	
	for _, model := range cmdbTables {
		if err := db.Migrator().DropTable(model); err != nil {
			zap.L().Error("Failed to drop CMDB table", zap.Error(err))
		}
	}
	
	// 删除系统表
	systemTables := []interface{}{
		&system.SystemDepartment{},
		&system.SystemRole{},
		&system.SystemMenu{},
		&system.SystemMenuAuth{},
		&system.SystemUser{},
		&system.SystemUserLoginLog{},
		&system.SystemAnnouncement{},
		&nav.Nav{},
	}
	
	for _, model := range systemTables {
		if err := db.Migrator().DropTable(model); err != nil {
			zap.L().Error("Failed to drop system table", zap.Error(err))
		}
	}
	
	zap.L().Info("Database reset completed")
	return nil
} 
