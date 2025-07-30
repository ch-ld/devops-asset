package main

import (
	"encoding/json"
	"fmt"
	"log"

	"api-server/internal/config"
	"api-server/internal/db"
	"api-server/internal/db/mysql"
	"api-server/internal/model/nav"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(""); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := mysql.InitClient(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer mysql.Close()

	dbClient := mysql.GetClient()

	fmt.Println("=== 数据库迁移状态检查 ===")

	// 1. 检查迁移状态
	status, err := db.GetMigrationStatus(dbClient)
	if err != nil {
		log.Fatalf("Failed to get migration status: %v", err)
	}

	fmt.Println("\n📊 表状态:")
	for table, exists := range status {
		statusIcon := "❌"
		if exists {
			statusIcon = "✅"
		}
		fmt.Printf("  %s %s\n", statusIcon, table)
	}

	// 2. 检查迁移版本历史
	fmt.Println("\n📋 迁移历史:")
	var migrations []db.MigrationVersion
	if err := dbClient.Order("applied_at DESC").Find(&migrations).Error; err != nil {
		fmt.Printf("  ⚠️  无法获取迁移历史: %v\n", err)
	} else if len(migrations) == 0 {
		fmt.Println("  📝 没有迁移记录（可能是首次运行）")
	} else {
		for _, migration := range migrations {
			statusIcon := "❌"
			if migration.Applied {
				statusIcon = "✅"
			}
			fmt.Printf("  %s v%s - %s (应用于: %s)\n", 
				statusIcon, migration.Version, migration.Description, 
				migration.AppliedAt.Format("2006-01-02 15:04:05"))
		}
	}

	// 3. 检查导航栏数据
	fmt.Println("\n🧭 导航栏数据:")
	var navCount int64
	dbClient.Model(&nav.Nav{}).Count(&navCount)
	fmt.Printf("  📊 总导航项数量: %d\n", navCount)

	if navCount > 0 {
		var navItems []nav.Nav
		dbClient.Order("group_name, order_num").Find(&navItems)
		
		// 按组分类显示
		groupMap := make(map[string][]nav.Nav)
		for _, item := range navItems {
			groupMap[item.GroupName] = append(groupMap[item.GroupName], item)
		}

		for groupName, items := range groupMap {
			fmt.Printf("\n  📂 %s 组:\n", groupName)
			for _, item := range items {
				statusIcon := "✅"
				if item.Status != 1 {
					statusIcon = "❌"
				}
				fmt.Printf("    %s %s - %s\n", statusIcon, item.Name, item.Description)
			}
		}
	}

	// 4. 获取数据库详细信息
	fmt.Println("\n📈 数据库统计信息:")
	dbInfo := db.GetDatabaseInfo(dbClient)
	
	if totalTables, ok := dbInfo["total_tables"].(int); ok {
		fmt.Printf("  📊 总表数量: %d\n", totalTables)
	}
	if systemTables, ok := dbInfo["system_tables"].(int); ok {
		fmt.Printf("  🏢 系统表数量: %d\n", systemTables)
	}
	if cmdbTables, ok := dbInfo["cmdb_tables"].(int); ok {
		fmt.Printf("  💾 CMDB表数量: %d\n", cmdbTables)
	}

	if dataStats, ok := dbInfo["data_statistics"].(map[string]int64); ok {
		fmt.Println("  📋 数据统计:")
		for key, value := range dataStats {
			fmt.Printf("    - %s: %d 条记录\n", key, value)
		}
	}

	fmt.Println("\n✨ 迁移状态检查完成！")
} 
