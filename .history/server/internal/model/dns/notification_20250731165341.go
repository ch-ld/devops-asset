package dns

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Notification 通知模型
type Notification struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Type      string         `gorm:"type:varchar(50);not null;index" json:"type" comment:"通知类型"`
	Title     string         `gorm:"type:varchar(200);not null" json:"title" comment:"通知标题"`
	Message   string         `gorm:"type:text;not null" json:"message" comment:"通知内容"`
	Severity  string         `gorm:"type:varchar(20);not null;default:'info'" json:"severity" comment:"严重程度"`
	Channels  string         `gorm:"type:varchar(100)" json:"channels" comment:"通知渠道"`
	Status    string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status" comment:"发送状态"`
	SentAt    *time.Time     `gorm:"type:datetime" json:"sent_at" comment:"发送时间"`
	ReadAt    *time.Time     `gorm:"type:datetime" json:"read_at" comment:"阅读时间"`
	ErrorMsg  string         `gorm:"type:text" json:"error_msg" comment:"错误信息"`
	Metadata  string         `gorm:"type:json" json:"metadata" comment:"元数据"`
	TenantID  uint           `gorm:"not null;index" json:"tenant_id" comment:"租户ID"`
	CreatedBy uint           `gorm:"not null" json:"created_by" comment:"创建者ID"`
	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at" comment:"创建时间"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at" comment:"更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" comment:"删除时间"`
}

// TableName 指定表名
func (Notification) TableName() string {
	return "dns_notifications"
}

// BeforeCreate 创建前钩子
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	n.CreatedAt = now
	n.UpdatedAt = now
	return nil
}

// BeforeUpdate 更新前钩子
func (n *Notification) BeforeUpdate(tx *gorm.DB) error {
	n.UpdatedAt = time.Now()
	return nil
}

// IsPending 检查是否待发送
func (n *Notification) IsPending() bool {
	return n.Status == "pending"
}

// IsSent 检查是否已发送
func (n *Notification) IsSent() bool {
	return n.Status == "sent"
}

// IsFailed 检查是否发送失败
func (n *Notification) IsFailed() bool {
	return n.Status == "failed"
}

// IsRead 检查是否已读
func (n *Notification) IsRead() bool {
	return n.ReadAt != nil
}

// MarkAsSent 标记为已发送
func (n *Notification) MarkAsSent() {
	n.Status = "sent"
	now := time.Now()
	n.SentAt = &now
}

// MarkAsFailed 标记为发送失败
func (n *Notification) MarkAsFailed(errorMsg string) {
	n.Status = "failed"
	n.ErrorMsg = errorMsg
}

// MarkAsRead 标记为已读
func (n *Notification) MarkAsRead() {
	now := time.Now()
	n.ReadAt = &now
}

// GetSeverityColor 获取严重程度颜色
func (n *Notification) GetSeverityColor() string {
	switch n.Severity {
	case "info":
		return "blue"
	case "warning":
		return "orange"
	case "error":
		return "red"
	case "critical":
		return "purple"
	default:
		return "gray"
	}
}

// GetSeverityIcon 获取严重程度图标
func (n *Notification) GetSeverityIcon() string {
	switch n.Severity {
	case "info":
		return "info-circle"
	case "warning":
		return "exclamation-triangle"
	case "error":
		return "times-circle"
	case "critical":
		return "exclamation-circle"
	default:
		return "bell"
	}
}

// GetTypeDisplayName 获取类型显示名称
func (n *Notification) GetTypeDisplayName() string {
	switch n.Type {
	case "domain_expiry":
		return "域名过期提醒"
	case "cert_expiry":
		return "证书过期提醒"
	case "https_down":
		return "HTTPS服务异常"
	case "dns_sync_failed":
		return "DNS同步失败"
	case "system_alert":
		return "系统告警"
	default:
		return "通知"
	}
}

// GetChannelList 获取通知渠道列表
func (n *Notification) GetChannelList() []string {
	if n.Channels == "" {
		return []string{}
	}
	
	channels := []string{}
	for _, channel := range strings.Split(n.Channels, ",") {
		channel = strings.TrimSpace(channel)
		if channel != "" {
			channels = append(channels, channel)
		}
	}
	return channels
}

// SetChannels 设置通知渠道
func (n *Notification) SetChannels(channels []string) {
	n.Channels = strings.Join(channels, ",")
}

// IsUrgent 检查是否紧急
func (n *Notification) IsUrgent() bool {
	return n.Severity == "error" || n.Severity == "critical"
}

// GetAge 获取通知年龄
func (n *Notification) GetAge() time.Duration {
	return time.Since(n.CreatedAt)
}

// IsExpired 检查是否过期（超过7天）
func (n *Notification) IsExpired() bool {
	return n.GetAge() > 7*24*time.Hour
}

// ShouldRetry 检查是否应该重试发送
func (n *Notification) ShouldRetry() bool {
	if !n.IsFailed() {
		return false
	}
	
	// 失败后1小时内可以重试
	if n.SentAt != nil && time.Since(*n.SentAt) < time.Hour {
		return true
	}
	
	return false
}

// GetRetryCount 获取重试次数（从错误信息中解析）
func (n *Notification) GetRetryCount() int {
	// 简化实现，实际应该从元数据中获取
	if n.IsFailed() {
		return 1
	}
	return 0
}

// CanDelete 检查是否可以删除
func (n *Notification) CanDelete() bool {
	// 已读且超过30天的通知可以删除
	return n.IsRead() && n.GetAge() > 30*24*time.Hour
}

// GetPriority 获取优先级
func (n *Notification) GetPriority() int {
	switch n.Severity {
	case "critical":
		return 4
	case "error":
		return 3
	case "warning":
		return 2
	case "info":
		return 1
	default:
		return 0
	}
}

// Validate 验证模型数据
func (n *Notification) Validate() error {
	if n.Type == "" {
		return ErrInvalidNotificationType
	}
	
	if n.Title == "" {
		return ErrInvalidNotificationTitle
	}
	
	if n.Message == "" {
		return ErrInvalidNotificationMessage
	}
	
	validSeverities := []string{"info", "warning", "error", "critical"}
	validSeverity := false
	for _, severity := range validSeverities {
		if n.Severity == severity {
			validSeverity = true
			break
		}
	}
	if !validSeverity {
		return ErrInvalidNotificationSeverity
	}
	
	validStatuses := []string{"pending", "sent", "failed"}
	validStatus := false
	for _, status := range validStatuses {
		if n.Status == status {
			validStatus = true
			break
		}
	}
	if !validStatus {
		return ErrInvalidNotificationStatus
	}
	
	return nil
}

// 错误定义
var (
	ErrInvalidNotificationType     = errors.New("invalid notification type")
	ErrInvalidNotificationTitle    = errors.New("invalid notification title")
	ErrInvalidNotificationMessage  = errors.New("invalid notification message")
	ErrInvalidNotificationSeverity = errors.New("invalid notification severity")
	ErrInvalidNotificationStatus   = errors.New("invalid notification status")
)
