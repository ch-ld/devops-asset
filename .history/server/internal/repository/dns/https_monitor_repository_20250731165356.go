package dns

import (
	"api-server/internal/model/dns"

	"gorm.io/gorm"
)

// HTTPSMonitorRepository HTTPS监控仓库
type HTTPSMonitorRepository struct {
	DB *gorm.DB
}

// NewHTTPSMonitorRepository 创建HTTPS监控仓库
func NewHTTPSMonitorRepository(db *gorm.DB) *HTTPSMonitorRepository {
	return &HTTPSMonitorRepository{
		DB: db,
	}
}

// Create 创建监控
func (r *HTTPSMonitorRepository) Create(monitor *dns.HTTPSMonitor) error {
	return r.DB.Create(monitor).Error
}

// Update 更新监控
func (r *HTTPSMonitorRepository) Update(monitor *dns.HTTPSMonitor) error {
	return r.DB.Model(monitor).Updates(monitor).Error
}

// Delete 删除监控
func (r *HTTPSMonitorRepository) Delete(id uint) error {
	return r.DB.Delete(&dns.HTTPSMonitor{}, id).Error
}

// FindByID 根据ID查找监控
func (r *HTTPSMonitorRepository) FindByID(id uint) (*dns.HTTPSMonitor, error) {
	var monitor dns.HTTPSMonitor
	err := r.DB.First(&monitor, id).Error
	if err != nil {
		return nil, err
	}
	return &monitor, nil
}

// FindActiveMonitors 查找活跃的监控
func (r *HTTPSMonitorRepository) FindActiveMonitors() ([]*dns.HTTPSMonitor, error) {
	var monitors []*dns.HTTPSMonitor
	err := r.DB.Where("enabled = ?", true).
		Where("status = ?", "active").
		Find(&monitors).Error
	return monitors, err
}

// FindByStatus 根据状态查找监控
func (r *HTTPSMonitorRepository) FindByStatus(status string) ([]*dns.HTTPSMonitor, error) {
	var monitors []*dns.HTTPSMonitor
	err := r.DB.Where("status = ?", status).Find(&monitors).Error
	return monitors, err
}

// FindByURL 根据URL查找监控
func (r *HTTPSMonitorRepository) FindByURL(url string) (*dns.HTTPSMonitor, error) {
	var monitor dns.HTTPSMonitor
	err := r.DB.Where("url = ?", url).First(&monitor).Error
	if err != nil {
		return nil, err
	}
	return &monitor, nil
}

// CountAll 统计所有监控数量
func (r *HTTPSMonitorRepository) CountAll() (int64, error) {
	var count int64
	err := r.DB.Model(&dns.HTTPSMonitor{}).Count(&count).Error
	return count, err
}

// CountActive 统计活跃监控数量
func (r *HTTPSMonitorRepository) CountActive() (int64, error) {
	var count int64
	err := r.DB.Model(&dns.HTTPSMonitor{}).
		Where("enabled = ?", true).
		Where("status = ?", "active").
		Count(&count).Error
	return count, err
}

// GetStatusStatistics 获取状态统计
func (r *HTTPSMonitorRepository) GetStatusStatistics() (map[string]int64, error) {
	var results []struct {
		LastStatus string
		Count      int64
	}

	err := r.DB.Model(&dns.HTTPSMonitor{}).
		Select("last_status, COUNT(*) as count").
		Where("enabled = ?", true).
		Group("last_status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	statusCount := make(map[string]int64)
	for _, result := range results {
		statusCount[result.LastStatus] = result.Count
	}

	return statusCount, nil
}

// CountExpiringCertificates 统计即将过期的证书数量
func (r *HTTPSMonitorRepository) CountExpiringCertificates(days int) (int64, error) {
	var count int64
	err := r.DB.Model(&dns.HTTPSMonitor{}).
		Where("enabled = ?", true).
		Where("cert_expiry IS NOT NULL").
		Where("cert_expiry <= DATE_ADD(NOW(), INTERVAL ? DAY)", days).
		Count(&count).Error
	return count, err
}

// FindExpiringCertificates 查找即将过期的证书
func (r *HTTPSMonitorRepository) FindExpiringCertificates(days int) ([]*dns.HTTPSMonitor, error) {
	var monitors []*dns.HTTPSMonitor
	err := r.DB.Where("enabled = ?", true).
		Where("cert_expiry IS NOT NULL").
		Where("cert_expiry <= DATE_ADD(NOW(), INTERVAL ? DAY)", days).
		Order("cert_expiry ASC").
		Find(&monitors).Error
	return monitors, err
}

// FindFailingMonitors 查找失败的监控
func (r *HTTPSMonitorRepository) FindFailingMonitors(minFailures int) ([]*dns.HTTPSMonitor, error) {
	var monitors []*dns.HTTPSMonitor
	err := r.DB.Where("enabled = ?", true).
		Where("failure_count >= ?", minFailures).
		Where("last_status = ?", "error").
		Find(&monitors).Error
	return monitors, err
}

// ListWithPagination 分页查询监控
func (r *HTTPSMonitorRepository) ListWithPagination(page, pageSize int, filters map[string]interface{}) ([]*dns.HTTPSMonitor, int64, error) {
	var monitors []*dns.HTTPSMonitor
	var total int64

	query := r.DB.Model(&dns.HTTPSMonitor{})

	// 应用过滤条件
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if lastStatus, ok := filters["last_status"]; ok && lastStatus != "" {
		query = query.Where("last_status = ?", lastStatus)
	}

	if enabled, ok := filters["enabled"]; ok {
		query = query.Where("enabled = ?", enabled)
	}

	if keyword, ok := filters["keyword"]; ok && keyword != "" {
		query = query.Where("name LIKE ? OR url LIKE ?",
			"%"+keyword.(string)+"%",
			"%"+keyword.(string)+"%")
	}

	if tenantID, ok := filters["tenant_id"]; ok && tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&monitors).Error

	return monitors, total, err
}

// UpdateMonitorStatus 更新监控状态
func (r *HTTPSMonitorRepository) UpdateMonitorStatus(id uint, status string, responseTime int, errorMsg string) error {
	updates := map[string]interface{}{
		"last_status":        status,
		"last_response_time": responseTime,
		"last_checked":       "NOW()",
	}

	if errorMsg != "" {
		updates["last_error"] = errorMsg
	} else {
		updates["last_error"] = ""
	}

	return r.DB.Model(&dns.HTTPSMonitor{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// UpdateCertificateInfo 更新证书信息
func (r *HTTPSMonitorRepository) UpdateCertificateInfo(id uint, certExpiry string, daysLeft int) error {
	return r.DB.Model(&dns.HTTPSMonitor{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"cert_expiry":    certExpiry,
			"cert_days_left": daysLeft,
		}).Error
}

// IncrementFailureCount 增加失败次数
func (r *HTTPSMonitorRepository) IncrementFailureCount(id uint) error {
	return r.DB.Model(&dns.HTTPSMonitor{}).
		Where("id = ?", id).
		Update("failure_count", gorm.Expr("failure_count + 1")).Error
}

// ResetFailureCount 重置失败次数
func (r *HTTPSMonitorRepository) ResetFailureCount(id uint) error {
	return r.DB.Model(&dns.HTTPSMonitor{}).
		Where("id = ?", id).
		Update("failure_count", 0).Error
}

// GetMonitorStatistics 获取监控统计信息
func (r *HTTPSMonitorRepository) GetMonitorStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总监控数
	totalCount, err := r.CountAll()
	if err != nil {
		return nil, err
	}
	stats["total_monitors"] = totalCount

	// 活跃监控数
	activeCount, err := r.CountActive()
	if err != nil {
		return nil, err
	}
	stats["active_monitors"] = activeCount

	// 状态分布
	statusStats, err := r.GetStatusStatistics()
	if err != nil {
		return nil, err
	}
	stats["status_distribution"] = statusStats

	// 即将过期的证书数量（7天内）
	expiringCount, err := r.CountExpiringCertificates(7)
	if err != nil {
		return nil, err
	}
	stats["expiring_certificates"] = expiringCount

	// 失败的监控数量
	var failingCount int64
	err = r.DB.Model(&dns.HTTPSMonitor{}).
		Where("enabled = ?", true).
		Where("last_status = ?", "error").
		Count(&failingCount).Error
	if err != nil {
		return nil, err
	}
	stats["failing_monitors"] = failingCount

	// 平均响应时间
	var avgResponseTime float64
	err = r.DB.Model(&dns.HTTPSMonitor{}).
		Where("enabled = ?", true).
		Where("last_response_time > 0").
		Select("AVG(last_response_time)").
		Scan(&avgResponseTime).Error
	if err != nil {
		return nil, err
	}
	stats["avg_response_time"] = avgResponseTime

	return stats, nil
}

// FindMonitorsNeedingCheck 查找需要检查的监控
func (r *HTTPSMonitorRepository) FindMonitorsNeedingCheck() ([]*dns.HTTPSMonitor, error) {
	var monitors []*dns.HTTPSMonitor
	err := r.DB.Where("enabled = ?", true).
		Where("status = ?", "active").
		Where("last_checked < DATE_SUB(NOW(), INTERVAL check_interval SECOND)").
		Find(&monitors).Error
	return monitors, err
}

// UpdateMonitorEnabled 更新监控启用状态
func (r *HTTPSMonitorRepository) UpdateMonitorEnabled(id uint, enabled bool) error {
	return r.DB.Model(&dns.HTTPSMonitor{}).
		Where("id = ?", id).
		Update("enabled", enabled).Error
}

// FindByTenantID 根据租户ID查找监控
func (r *HTTPSMonitorRepository) FindByTenantID(tenantID uint) ([]*dns.HTTPSMonitor, error) {
	var monitors []*dns.HTTPSMonitor
	err := r.DB.Where("tenant_id = ?", tenantID).Find(&monitors).Error
	return monitors, err
}

// DeleteByTenantID 根据租户ID删除监控
func (r *HTTPSMonitorRepository) DeleteByTenantID(tenantID uint) error {
	return r.DB.Where("tenant_id = ?", tenantID).Delete(&dns.HTTPSMonitor{}).Error
}
