package dns

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Record DNS解析记录模型
type Record struct {
	gorm.Model
	DomainID     uint           `gorm:"not null;index;comment:域名ID" json:"domain_id"`
	Domain       *Domain        `gorm:"foreignKey:DomainID" json:"domain,omitempty"`
	ProviderID   uint           `gorm:"not null;index;comment:DNS提供商ID" json:"provider_id"`
	Provider     *Provider      `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	Name         string         `gorm:"size:255;not null;comment:记录名称" json:"name"`
	Type         string         `gorm:"size:10;not null;comment:记录类型(A/AAAA/CNAME/TXT/MX/SRV/NS)" json:"type"`
	Value        string         `gorm:"size:1000;not null;comment:记录值" json:"value"`
	TTL          int            `gorm:"default:600;comment:TTL值(秒)" json:"ttl"`
	Priority     *int           `gorm:"comment:优先级(MX/SRV记录)" json:"priority"`
	Weight       *int           `gorm:"comment:权重(SRV记录)" json:"weight"`
	Port         *int           `gorm:"comment:端口(SRV记录)" json:"port"`
	Status       string         `gorm:"size:20;not null;default:active;comment:状态(active/inactive)" json:"status"`
	CloudRecordID string        `gorm:"size:100;comment:云端记录ID" json:"cloud_record_id"`
	SyncStatus   string         `gorm:"size:20;default:synced;comment:同步状态(synced/pending/failed)" json:"sync_status"`
	LastSyncAt   *gorm.DeletedAt `gorm:"comment:最后同步时间" json:"last_sync_at"`
	Configuration datatypes.JSON `gorm:"type:json;comment:配置信息" json:"configuration"`
	Remark       string         `gorm:"size:1000;comment:备注" json:"remark"`
	TenantID     uint           `gorm:"not null;index;comment:租户ID" json:"tenant_id"`
	CreatedBy    uint           `gorm:"comment:创建人ID" json:"created_by"`
	UpdatedBy    uint           `gorm:"comment:更新人ID" json:"updated_by"`
}

// TableName 获取表名
func (r *Record) TableName() string {
	return "dns_records"
}

// RecordTemplate DNS记录模板模型
type RecordTemplate struct {
	gorm.Model
	Name         string         `gorm:"size:100;not null;comment:模板名称" json:"name"`
	Description  string         `gorm:"size:500;comment:模板描述" json:"description"`
	Records      datatypes.JSON `gorm:"type:json;comment:记录模板数据" json:"records"`
	Category     string         `gorm:"size:50;comment:分类(web/mail/cdn)" json:"category"`
	IsPublic     bool           `gorm:"default:false;comment:是否公开模板" json:"is_public"`
	UsageCount   int            `gorm:"default:0;comment:使用次数" json:"usage_count"`
	TenantID     uint           `gorm:"not null;index;comment:租户ID" json:"tenant_id"`
	CreatedBy    uint           `gorm:"comment:创建人ID" json:"created_by"`
	UpdatedBy    uint           `gorm:"comment:更新人ID" json:"updated_by"`
}

// TableName 获取表名
func (rt *RecordTemplate) TableName() string {
	return "dns_record_templates"
}
