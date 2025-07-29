package cmdb

import (
	"gorm.io/gorm"
)

// HostGroup 主机组模型
type HostGroup struct {
	gorm.Model
	Name           string      `gorm:"size:100;not null;comment:主机组名称" json:"name"`
	Description    string      `gorm:"size:500;comment:描述信息" json:"description"`
	ParentID       *uint       `gorm:"comment:父组ID" json:"parent_id"`
	Path           string      `gorm:"size:500;comment:组路径" json:"path"`
	Sort           int         `gorm:"default:0;comment:排序" json:"sort"`
	HostCount      int64       `gorm:"-" json:"host_count"`         // 直接归属的主机数量，不存储到数据库
	TotalHostCount int64       `gorm:"-" json:"total_host_count"`   // 包含子分组的总主机数量，不存储到数据库
	Children       []HostGroup `gorm:"-" json:"children,omitempty"` // 子组，不存储到数据库
}

// TableName 获取表名
func (g *HostGroup) TableName() string {
	return "cmdb_host_groups"
}
