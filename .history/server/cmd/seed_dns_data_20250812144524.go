package main

import (
	"api-server/internal/config"
	"api-server/internal/db/mysql"
	"api-server/internal/model/dns"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("./configs/config.dev.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	db, err := mysql.NewClient(cfg.Database.MySQL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 添加测试数据
	if err := seedDNSTestData(db); err != nil {
		log.Fatalf("Failed to seed DNS test data: %v", err)
	}

	fmt.Println("DNS test data seeded successfully!")
}

func seedDNSTestData(db *gorm.DB) error {
	// 创建测试分组
	groups := []*dns.DomainGroup{
		{
			Model:       gorm.Model{ID: 1},
			Name:        "生产环境",
			Description: "生产环境域名",
			Sort:        1,
			TenantID:    1,
			CreatedBy:   1,
			UpdatedBy:   1,
		},
		{
			Model:       gorm.Model{ID: 2},
			Name:        "测试环境",
			Description: "测试环境域名",
			Sort:        2,
			TenantID:    1,
			CreatedBy:   1,
			UpdatedBy:   1,
		},
	}

	// 插入分组数据
	for _, group := range groups {
		if err := db.FirstOrCreate(group, "id = ?", group.ID).Error; err != nil {
			return fmt.Errorf("failed to create domain group: %w", err)
		}
	}

	// 创建测试标签
	tags := []*dns.Tag{
		{
			Model:       gorm.Model{ID: 1},
			Name:        "重要",
			Color:       "#ff4d4f",
			Description: "重要域名",
			TenantID:    1,
			CreatedBy:   1,
			UpdatedBy:   1,
		},
		{
			Model:       gorm.Model{ID: 2},
			Name:        "API",
			Color:       "#1890ff",
			Description: "API服务域名",
			TenantID:    1,
			CreatedBy:   1,
			UpdatedBy:   1,
		},
		{
			Model:       gorm.Model{ID: 3},
			Name:        "CDN",
			Color:       "#52c41a",
			Description: "CDN域名",
			TenantID:    1,
			CreatedBy:   1,
			UpdatedBy:   1,
		},
	}

	// 插入标签数据
	for _, tag := range tags {
		if err := db.FirstOrCreate(tag, "id = ?", tag.ID).Error; err != nil {
			return fmt.Errorf("failed to create tag: %w", err)
		}
	}

	// 创建测试域名
	expireDate1 := time.Now().AddDate(0, 6, 0) // 6个月后过期
	expireDate2 := time.Now().AddDate(1, 0, 0) // 1年后过期
	expireDate3 := time.Now().AddDate(0, 1, 15) // 1个月15天后过期
	expireDate4 := time.Now().AddDate(0, 0, 20) // 20天后过期

	groupID1 := uint(1)
	groupID2 := uint(2)

	domains := []*dns.Domain{
		{
			Model:         gorm.Model{ID: 1},
			Name:          "example.com",
			Status:        "active",
			RegistrarType: "aliyun",
			ExpiresAt:     &expireDate1,
			AutoRenew:     true,
			GroupID:       &groupID1,
			Remark:        "主域名",
			TenantID:      1,
			CreatedBy:     1,
			UpdatedBy:     1,
		},
		{
			Model:         gorm.Model{ID: 2},
			Name:          "api.example.com",
			Status:        "active",
			RegistrarType: "aliyun",
			ExpiresAt:     &expireDate2,
			AutoRenew:     false,
			GroupID:       &groupID1,
			Remark:        "API服务域名",
			TenantID:      1,
			CreatedBy:     1,
			UpdatedBy:     1,
		},
		{
			Model:         gorm.Model{ID: 3},
			Name:          "test.com",
			Status:        "active",
			RegistrarType: "godaddy",
			ExpiresAt:     &expireDate3,
			AutoRenew:     true,
			GroupID:       &groupID2,
			Remark:        "测试域名",
			TenantID:      1,
			CreatedBy:     1,
			UpdatedBy:     1,
		},
		{
			Model:         gorm.Model{ID: 4},
			Name:          "cdn.example.com",
			Status:        "active",
			RegistrarType: "route53",
			ExpiresAt:     &expireDate4,
			AutoRenew:     false,
			GroupID:       &groupID1,
			Remark:        "CDN域名",
			TenantID:      1,
			CreatedBy:     1,
			UpdatedBy:     1,
		},
		{
			Model:         gorm.Model{ID: 5},
			Name:          "staging.example.com",
			Status:        "inactive",
			RegistrarType: "tencent",
			ExpiresAt:     &expireDate2,
			AutoRenew:     false,
			GroupID:       &groupID2,
			Remark:        "预发布环境域名",
			TenantID:      1,
			CreatedBy:     1,
			UpdatedBy:     1,
		},
	}

	// 插入域名数据
	for _, domain := range domains {
		if err := db.FirstOrCreate(domain, "id = ?", domain.ID).Error; err != nil {
			return fmt.Errorf("failed to create domain: %w", err)
		}
	}

	// 建立域名和标签的关联关系
	domainTags := []struct {
		DomainID uint
		TagID    uint
	}{
		{1, 1}, // example.com - 重要
		{2, 2}, // api.example.com - API
		{4, 3}, // cdn.example.com - CDN
		{4, 1}, // cdn.example.com - 重要
	}

	// 插入域名标签关联
	for _, dt := range domainTags {
		var count int64
		db.Table("dns_domain_tags").Where("domain_id = ? AND tag_id = ?", dt.DomainID, dt.TagID).Count(&count)
		if count == 0 {
			if err := db.Table("dns_domain_tags").Create(map[string]interface{}{
				"domain_id": dt.DomainID,
				"tag_id":    dt.TagID,
			}).Error; err != nil {
				return fmt.Errorf("failed to create domain-tag relation: %w", err)
			}
		}
	}

	fmt.Println("Created test data:")
	fmt.Printf("- %d domain groups\n", len(groups))
	fmt.Printf("- %d tags\n", len(tags))
	fmt.Printf("- %d domains\n", len(domains))
	fmt.Printf("- %d domain-tag relations\n", len(domainTags))

	return nil
}
