package dns

import (
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Domain 域名模型
type Domain struct {
	gorm.Model
	Name          string         `gorm:"size:255;not null;uniqueIndex;comment:域名" json:"name"`
	Status        string         `gorm:"size:20;not null;default:active;comment:状态(active/inactive/expired)" json:"status"`
	IsWildcard    bool           `gorm:"default:false;comment:是否为通配符域名" json:"is_wildcard"`
	ParentDomain  string         `gorm:"size:255;comment:父域名(通配符域名的验证域名)" json:"parent_domain"`
	Level         int            `gorm:"default:1;comment:域名级别(1=顶级域名,2=二级域名,3=三级域名等)" json:"level"`
	RegistrarType string         `gorm:"size:50;comment:注册商类型(godaddy/aliyun/tencent)" json:"registrar_type"`
	RegistrarID   *uint          `gorm:"comment:注册商配置ID" json:"registrar_id"`
	ProviderID    *uint          `gorm:"comment:DNS提供商ID" json:"provider_id"`
	Provider      *Provider      `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	ExpiresAt     *time.Time     `gorm:"comment:过期时间" json:"expires_at"`
	AutoRenew     bool           `gorm:"default:false;comment:是否自动续费" json:"auto_renew"`
	GroupID       *uint          `gorm:"comment:分组ID" json:"group_id"`
	Group         *DomainGroup   `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Tags          []Tag          `gorm:"many2many:dns_domain_tags;" json:"tags,omitempty"`
	Configuration datatypes.JSON `gorm:"type:json;comment:配置信息" json:"configuration"`
	Remark        string         `gorm:"size:1000;comment:备注" json:"remark"`
	TenantID      uint           `gorm:"not null;index;comment:租户ID" json:"tenant_id"`
	CreatedBy     uint           `gorm:"comment:创建人ID" json:"created_by"`
	UpdatedBy     uint           `gorm:"comment:更新人ID" json:"updated_by"`
}

// TableName 获取表名
func (d *Domain) TableName() string {
	return "dns_domains"
}

// DomainGroup 域名分组模型
type DomainGroup struct {
	gorm.Model
	Name             string         `gorm:"size:100;not null;comment:分组名称" json:"name"`
	ParentID         *uint          `gorm:"comment:父分组ID" json:"parent_id"`
	Parent           *DomainGroup   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children         []*DomainGroup `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Description      string         `gorm:"size:500;comment:描述" json:"description"`
	Sort             int            `gorm:"default:0;comment:排序" json:"sort"`
	Color            string         `gorm:"size:20;default:#3b82f6;comment:分组颜色" json:"color"`
	Status           string         `gorm:"size:20;not null;default:active;comment:状态(active/inactive)" json:"status"`
	DomainCount      int            `gorm:"-" json:"domain_count,omitempty"`       // 直接归属域名数量
	TotalDomainCount int            `gorm:"-" json:"total_domain_count,omitempty"` // 包含子分组的总域名数量
	TenantID         uint           `gorm:"not null;index;comment:租户ID" json:"tenant_id"`
	CreatedBy        uint           `gorm:"comment:创建人ID" json:"created_by"`
	UpdatedBy        uint           `gorm:"comment:更新人ID" json:"updated_by"`
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

// TableName 获取表名
func (d *Domain) TableName() string {
	return "dns_domains"
}

// IsWildcardDomain 检查是否为通配符域名
func (d *Domain) IsWildcardDomain() bool {
	return d.IsWildcard || strings.HasPrefix(d.Name, "*.")
}

// GetParentDomain 获取父域名（用于DNS验证）
func (d *Domain) GetParentDomain() string {
	if d.ParentDomain != "" {
		return d.ParentDomain
	}

	// 如果是通配符域名，返回去掉*. 的部分
	if strings.HasPrefix(d.Name, "*.") {
		return d.Name[2:]
	}

	return d.Name
}

// GetDomainLevel 计算域名级别
func GetDomainLevel(domain string) int {
	// 去掉通配符前缀
	if strings.HasPrefix(domain, "*.") {
		domain = domain[2:]
	}

	// 计算点的数量 + 1
	return strings.Count(domain, ".") + 1
}

// ParseWildcardDomain 解析通配符域名信息
func ParseWildcardDomain(domain string) (isWildcard bool, parentDomain string, level int) {
	isWildcard = strings.HasPrefix(domain, "*.")
	if isWildcard {
		parentDomain = domain[2:]
	} else {
		parentDomain = domain
	}
	level = GetDomainLevel(domain)
	return
}
