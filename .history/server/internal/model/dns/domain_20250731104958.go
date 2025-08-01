package dns

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Domain 域名模型
type Domain struct {
	gorm.Model
	Name           string         `gorm:"size:255;not null;uniqueIndex;comment:域名" json:"name"`
	Status         string         `gorm:"size:20;not null;default:active;comment:状态(active/inactive/expired)" json:"status"`
	RegistrarType  string         `gorm:"size:50;comment:注册商类型(godaddy/aliyun/tencent)" json:"registrar_type"`
	RegistrarID    *uint          `gorm:"comment:注册商配置ID" json:"registrar_id"`
	ExpiresAt      *time.Time     `gorm:"comment:过期时间" json:"expires_at"`
	AutoRenew      bool           `gorm:"default:false;comment:是否自动续费" json:"auto_renew"`
	GroupID        *uint          `gorm:"comment:分组ID" json:"group_id"`
	Group          *DomainGroup   `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Tags           []Tag          `gorm:"many2many:dns_domain_tags;" json:"tags,omitempty"`
	Configuration  datatypes.JSON `gorm:"type:json;comment:配置信息" json:"configuration"`
	Remark         string         `gorm:"size:1000;comment:备注" json:"remark"`
	TenantID       uint           `gorm:"not null;index;comment:租户ID" json:"tenant_id"`
	CreatedBy      uint           `gorm:"comment:创建人ID" json:"created_by"`
	UpdatedBy      uint           `gorm:"comment:更新人ID" json:"updated_by"`
}

// TableName 获取表名
func (d *Domain) TableName() string {
	return "dns_domains"
}

// DomainGroup 域名分组模型
type DomainGroup struct {
	gorm.Model
	Name        string   `gorm:"size:100;not null;comment:分组名称" json:"name"`
	ParentID    *uint    `gorm:"comment:父分组ID" json:"parent_id"`
	Parent      *DomainGroup `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []DomainGroup `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Description string   `gorm:"size:500;comment:描述" json:"description"`
	Sort        int      `gorm:"default:0;comment:排序" json:"sort"`
	TenantID    uint     `gorm:"not null;index;comment:租户ID" json:"tenant_id"`
	CreatedBy   uint     `gorm:"comment:创建人ID" json:"created_by"`
	UpdatedBy   uint     `gorm:"comment:更新人ID" json:"updated_by"`
}

// TableName 获取表名
func (dg *DomainGroup) TableName() string {
	return "dns_domain_groups"
}

// Tag 标签模型
type Tag struct {
	gorm.Model
	Name        string `gorm:"size:50;not null;comment:标签名称" json:"name"`
	Color       string `gorm:"size:20;default:#1890ff;comment:标签颜色" json:"color"`
	Description string `gorm:"size:200;comment:描述" json:"description"`
	TenantID    uint   `gorm:"not null;index;comment:租户ID" json:"tenant_id"`
	CreatedBy   uint   `gorm:"comment:创建人ID" json:"created_by"`
	UpdatedBy   uint   `gorm:"comment:更新人ID" json:"updated_by"`
}

// TableName 获取表名
func (t *Tag) TableName() string {
	return "dns_tags"
}
