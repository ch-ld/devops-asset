package cmdb

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 资源类型常量
const (
	ResourceTypeECS = "ecs" // 阿里云ECS
	ResourceTypeEC2 = "ec2" // AWS EC2
	ResourceTypeCVM = "cvm" // 腾讯云CVM
)

// Host 主机模型
type Host struct {
	gorm.Model
	InstanceID        string         `gorm:"size:100;uniqueIndex;comment:实例ID" json:"instance_id"`
	Name              string         `gorm:"size:100;not null;comment:主机名称" json:"name"`
	Status            string         `gorm:"size:20;not null;default:unknown;comment:状态" json:"status"`
	PublicIP          datatypes.JSON `gorm:"type:json;comment:公网IP" json:"public_ip"`
	PrivateIP         datatypes.JSON `gorm:"type:json;comment:私网IP" json:"private_ip"`
	Configuration     datatypes.JSON `gorm:"type:json;comment:配置信息" json:"configuration"`
	OS                string         `gorm:"size:100;comment:操作系统" json:"os"`
	Region            string         `gorm:"size:50;comment:区域" json:"region"`
	Username          string         `gorm:"size:100;comment:SSH用户名" json:"username"`
	Password          string         `gorm:"size:1000;comment:SSH密码" json:"password"`
	PrivateKey        string         `gorm:"type:text;comment:SSH私钥" json:"private_key"`
	Port              int            `gorm:"default:22;comment:SSH端口" json:"port"`
	AuthType          string         `gorm:"size:20;default:password;comment:SSH认证类型(password/privatekey/both)" json:"auth_type"`
	SSHStatus         string         `gorm:"size:20;default:unknown;comment:SSH连接状态(unknown/online/offline/error)" json:"ssh_status"`
	LastConnectedAt   *time.Time     `gorm:"comment:最后SSH连接时间" json:"last_connected_at"`
	ConnectionTimeout int            `gorm:"default:30;comment:SSH连接超时时间(秒)" json:"connection_timeout"`
	ExpiredAt         *time.Time     `gorm:"comment:过期时间" json:"expired_at"`
	ProviderType      string         `gorm:"size:20;not null;default:manual;comment:提供商类型(aliyun/aws/tencent/manual)" json:"provider_type"`
	ResourceType      string         `gorm:"size:20;not null;default:unknown;comment:资源类型(ecs/ec2/cvm)" json:"resource_type"`
	ProviderID        *uint          `gorm:"comment:云账号ID" json:"provider_id"`
	Provider          *Provider      `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	GroupID           *uint          `gorm:"comment:主机组ID" json:"group_id"`
	Group             *HostGroup     `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Tags              datatypes.JSON `gorm:"type:json;comment:标签" json:"tags"`
	ExtraFields       datatypes.JSON `gorm:"type:json;comment:自定义字段" json:"extra_fields"`
	Remark            string         `gorm:"size:1000;comment:备注" json:"remark"`
}

// TableName 获取表名
func (h *Host) TableName() string {
	return "cmdb_hosts"
}
