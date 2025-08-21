package dns

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Certificate 证书模型
type Certificate struct {
	gorm.Model
	DomainID          *uint          `gorm:"index;comment:域名ID（可选）" json:"domain_id"`
	Domain            *Domain        `gorm:"foreignKey:DomainID" json:"domain,omitempty"`
	CommonName        string         `gorm:"size:255;not null;comment:主域名" json:"common_name"`
	SubjectAltNames   datatypes.JSON `gorm:"type:json;comment:SAN域名列表" json:"subject_alt_names"`
	CAType            string         `gorm:"size:50;not null;default:letsencrypt;comment:CA类型(letsencrypt/zerossl/buypass)" json:"ca_type"`
	Status            string         `gorm:"size:20;not null;default:pending;comment:状态(pending/issued/expired/revoked)" json:"status"`
	CertificateEnc    string         `gorm:"type:text;comment:证书内容(加密)" json:"-"`
	PrivateKeyEnc     string         `gorm:"type:text;comment:私钥内容(加密)" json:"-"`
	ChainEnc          string         `gorm:"type:text;comment:证书链(加密)" json:"-"`
	SerialNumber      string         `gorm:"size:100;comment:证书序列号" json:"serial_number"`
	Fingerprint       string         `gorm:"size:100;comment:证书指纹" json:"fingerprint"`
	KeyType           string         `gorm:"size:20;default:RSA2048;comment:密钥类型" json:"key_type"`
	ChallengeType     string         `gorm:"size:20;default:dns;comment:验证方式" json:"challenge_type"`
	IssuedAt          *time.Time     `gorm:"comment:签发时间" json:"issued_at"`
	ExpiresAt         *time.Time     `gorm:"comment:过期时间" json:"expires_at"`
	AutoRenew         bool           `gorm:"default:true;comment:是否自动续期" json:"auto_renew"`
	RenewDays         int            `gorm:"default:30;comment:提前续期天数" json:"renew_days"`
	ScheduleTime      string         `gorm:"size:10;comment:定时触发时间" json:"schedule_time"`
	NotificationType  string         `gorm:"size:20;default:default;comment:通知类型" json:"notification_type"`
	NotificationEmail string         `gorm:"size:255;comment:自定义通知邮箱" json:"notification_email"`
	LastRenewAt       *time.Time     `gorm:"comment:最后续期时间" json:"last_renew_at"`
	Configuration     datatypes.JSON `gorm:"type:json;comment:配置信息" json:"configuration"`
	Remark            string         `gorm:"size:1000;comment:备注" json:"remark"`
	TenantID          uint           `gorm:"not null;index;comment:租户ID" json:"tenant_id"`
	CreatedBy         uint           `gorm:"comment:创建人ID" json:"created_by"`
	UpdatedBy         uint           `gorm:"comment:更新人ID" json:"updated_by"`
}

// TableName 获取表名
func (c *Certificate) TableName() string {
	return "dns_certificates"
}

// CertificateDeployment 证书部署记录模型
type CertificateDeployment struct {
	gorm.Model
	CertificateID uint           `gorm:"not null;index;comment:证书ID" json:"certificate_id"`
	Certificate   *Certificate   `gorm:"foreignKey:CertificateID" json:"certificate,omitempty"`
	HostID        uint           `gorm:"not null;comment:主机ID" json:"host_id"`
	DeployPath    string         `gorm:"size:500;comment:部署路径" json:"deploy_path"`
	ServiceName   string         `gorm:"size:100;comment:服务名称(nginx/apache/tomcat)" json:"service_name"`
	Status        string         `gorm:"size:20;not null;default:pending;comment:部署状态(pending/success/failed)" json:"status"`
	DeployedAt    *time.Time     `gorm:"comment:部署时间" json:"deployed_at"`
	ErrorMessage  string         `gorm:"size:1000;comment:错误信息" json:"error_message"`
	Configuration datatypes.JSON `gorm:"type:json;comment:部署配置" json:"configuration"`
	TenantID      uint           `gorm:"not null;index;comment:租户ID" json:"tenant_id"`
	CreatedBy     uint           `gorm:"comment:创建人ID" json:"created_by"`
	UpdatedBy     uint           `gorm:"comment:更新人ID" json:"updated_by"`
}

// TableName 获取表名
func (cd *CertificateDeployment) TableName() string {
	return "dns_certificate_deployments"
}
