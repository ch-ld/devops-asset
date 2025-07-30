package system

import (
	"api-server/internal/model/nav"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MigrateData 执行系统数据迁移（兼容性接口，实际逻辑已移至统一的数据库初始化）
// 该函数保留是为了兼容现有调用，但实际数据初始化已整合到 internal/db/migration.go 中
func MigrateData(db *gorm.DB) error {
	zap.L().Info("System MigrateData called - redirecting to unified migration system")
	// 这个函数现在只是一个占位符，实际的数据初始化逻辑已经移到了
	// internal/db/migration.go 中的 initializeSystemData 函数
	return nil
}

// Migrate 执行完整的系统模块迁移（兼容性接口）
func Migrate(db *gorm.DB) error {
	zap.L().Info("System Migrate called - performing table migration only")
	
	// 只执行表结构迁移，数据初始化由统一的迁移系统处理
	if err := migrateTable(db); err != nil {
		return err
	}
	
	zap.L().Info("System module table migration completed")
	return nil
}

// migrateTable 迁移表结构
func migrateTable(db *gorm.DB) error {
	err := db.AutoMigrate(&SystemDepartment{}, &SystemRole{}, &SystemMenu{}, &SystemMenuAuth{}, &SystemUser{}, &SystemUserLoginLog{}, &SystemAnnouncement{}, &nav.Nav{})
	if err != nil {
		zap.L().Error("failed to migrate system model", zap.Error(err))
		return err
	}
	return nil
}
