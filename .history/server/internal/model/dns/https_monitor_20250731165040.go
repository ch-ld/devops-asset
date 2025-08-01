package dns

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// HTTPSMonitor HTTPS监控模型
type HTTPSMonitor struct {
	ID               uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	URL              string         `gorm:"type:varchar(500);not null;index" json:"url" comment:"监控URL"`
	Name             string         `gorm:"type:varchar(100);not null" json:"name" comment:"监控名称"`
	Status           string         `gorm:"type:varchar(20);not null;default:'active'" json:"status" comment:"监控状态"`
	CheckInterval    int            `gorm:"not null;default:600" json:"check_interval" comment:"检查间隔(秒)"`
	Timeout          int            `gorm:"not null;default:30" json:"timeout" comment:"超时时间(秒)"`
	AlertThreshold   int            `gorm:"not null;default:7" json:"alert_threshold" comment:"证书过期提醒天数"`
	Enabled          bool           `gorm:"not null;default:true" json:"enabled" comment:"是否启用"`
	LastChecked      time.Time      `gorm:"type:datetime" json:"last_checked" comment:"最后检查时间"`
	LastStatus       string         `gorm:"type:varchar(20)" json:"last_status" comment:"最后检查状态"`
	LastResponseTime int            `gorm:"default:0" json:"last_response_time" comment:"最后响应时间(毫秒)"`
	FailureCount     int            `gorm:"default:0" json:"failure_count" comment:"连续失败次数"`
	CertExpiry       *time.Time     `gorm:"type:datetime" json:"cert_expiry" comment:"证书过期时间"`
	CertDaysLeft     *int           `json:"cert_days_left" comment:"证书剩余天数"`
	LastError        string         `gorm:"type:text" json:"last_error" comment:"最后错误信息"`
	Remark           string         `gorm:"type:varchar(500)" json:"remark" comment:"备注"`
	TenantID         uint           `gorm:"not null;index" json:"tenant_id" comment:"租户ID"`
	CreatedBy        uint           `gorm:"not null" json:"created_by" comment:"创建者ID"`
	UpdatedBy        uint           `gorm:"not null" json:"updated_by" comment:"更新者ID"`
	CreatedAt        time.Time      `gorm:"type:datetime;not null" json:"created_at" comment:"创建时间"`
	UpdatedAt        time.Time      `gorm:"type:datetime;not null" json:"updated_at" comment:"更新时间"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at" comment:"删除时间"`
}

// TableName 指定表名
func (HTTPSMonitor) TableName() string {
	return "dns_https_monitors"
}

// BeforeCreate 创建前钩子
func (m *HTTPSMonitor) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	m.LastChecked = now
	return nil
}

// BeforeUpdate 更新前钩子
func (m *HTTPSMonitor) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}

// IsActive 检查监控是否活跃
func (m *HTTPSMonitor) IsActive() bool {
	return m.Enabled && m.Status == "active"
}

// IsHealthy 检查监控是否健康
func (m *HTTPSMonitor) IsHealthy() bool {
	return m.LastStatus == "ok" || m.LastStatus == "warning"
}

// IsCertExpiring 检查证书是否即将过期
func (m *HTTPSMonitor) IsCertExpiring() bool {
	if m.CertDaysLeft == nil {
		return false
	}
	return *m.CertDaysLeft <= m.AlertThreshold
}

// GetStatusColor 获取状态颜色
func (m *HTTPSMonitor) GetStatusColor() string {
	switch m.LastStatus {
	case "ok":
		return "green"
	case "warning":
		return "orange"
	case "error":
		return "red"
	default:
		return "gray"
	}
}

// GetFailureLevel 获取失败级别
func (m *HTTPSMonitor) GetFailureLevel() string {
	if m.FailureCount == 0 {
		return "none"
	} else if m.FailureCount <= 3 {
		return "low"
	} else if m.FailureCount <= 10 {
		return "medium"
	} else {
		return "high"
	}
}

// UpdateCertInfo 更新证书信息
func (m *HTTPSMonitor) UpdateCertInfo(expiry time.Time, daysLeft int) {
	m.CertExpiry = &expiry
	m.CertDaysLeft = &daysLeft
}

// UpdateStatus 更新状态
func (m *HTTPSMonitor) UpdateStatus(status string, responseTime int, errorMsg string) {
	m.LastStatus = status
	m.LastResponseTime = responseTime
	m.LastError = errorMsg
	m.LastChecked = time.Now()
	
	if status == "error" {
		m.FailureCount++
	} else {
		m.FailureCount = 0
	}
}

// ResetFailures 重置失败计数
func (m *HTTPSMonitor) ResetFailures() {
	m.FailureCount = 0
}

// Disable 禁用监控
func (m *HTTPSMonitor) Disable() {
	m.Enabled = false
	m.Status = "disabled"
}

// Enable 启用监控
func (m *HTTPSMonitor) Enable() {
	m.Enabled = true
	m.Status = "active"
}

// GetDisplayName 获取显示名称
func (m *HTTPSMonitor) GetDisplayName() string {
	if m.Name != "" {
		return m.Name
	}
	return m.URL
}

// NeedsCheck 检查是否需要检查
func (m *HTTPSMonitor) NeedsCheck() bool {
	if !m.IsActive() {
		return false
	}
	
	interval := time.Duration(m.CheckInterval) * time.Second
	return time.Since(m.LastChecked) >= interval
}

// GetNextCheckTime 获取下次检查时间
func (m *HTTPSMonitor) GetNextCheckTime() time.Time {
	interval := time.Duration(m.CheckInterval) * time.Second
	return m.LastChecked.Add(interval)
}

// IsOverdue 检查是否过期未检查
func (m *HTTPSMonitor) IsOverdue() bool {
	if !m.IsActive() {
		return false
	}
	
	interval := time.Duration(m.CheckInterval) * time.Second
	overdueThreshold := interval * 2 // 超过2倍间隔时间认为过期
	return time.Since(m.LastChecked) >= overdueThreshold
}

// GetUptimePercentage 获取可用性百分比（简化计算）
func (m *HTTPSMonitor) GetUptimePercentage() float64 {
	// 这里是简化实现，实际应该基于历史数据计算
	if m.FailureCount == 0 {
		return 100.0
	} else if m.FailureCount <= 5 {
		return 95.0
	} else if m.FailureCount <= 10 {
		return 90.0
	} else {
		return 80.0
	}
}

// Validate 验证模型数据
func (m *HTTPSMonitor) Validate() error {
	if m.URL == "" {
		return ErrInvalidURL
	}
	
	if m.Name == "" {
		return ErrInvalidName
	}
	
	if m.CheckInterval < 60 {
		return ErrInvalidCheckInterval
	}
	
	if m.Timeout < 5 || m.Timeout > 300 {
		return ErrInvalidTimeout
	}
	
	if m.AlertThreshold < 1 || m.AlertThreshold > 90 {
		return ErrInvalidAlertThreshold
	}
	
	return nil
}

// 错误定义
var (
	ErrInvalidURL              = errors.New("invalid URL")
	ErrInvalidName             = errors.New("invalid name")
	ErrInvalidCheckInterval    = errors.New("check interval must be at least 60 seconds")
	ErrInvalidTimeout          = errors.New("timeout must be between 5 and 300 seconds")
	ErrInvalidAlertThreshold   = errors.New("alert threshold must be between 1 and 90 days")
)
