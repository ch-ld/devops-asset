package dns

import (
	"api-server/internal/model/dns"

	"gorm.io/gorm"
)

// NotificationRepository 通知仓库
type NotificationRepository struct {
	DB *gorm.DB
}

// NewNotificationRepository 创建通知仓库
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		DB: db,
	}
}

// FindByType 根据类型查找通知
func (r *NotificationRepository) FindByType(notificationType string) ([]*dns.Notification, error) {
	var notifications []*dns.Notification
	err := r.DB.Where("type = ?", notificationType).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// FindByStatus 根据状态查找通知
func (r *NotificationRepository) FindByStatus(status string) ([]*dns.Notification, error) {
	var notifications []*dns.Notification
	err := r.DB.Where("status = ?", status).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// FindBySeverity 根据严重程度查找通知
func (r *NotificationRepository) FindBySeverity(severity string) ([]*dns.Notification, error) {
	var notifications []*dns.Notification
	err := r.DB.Where("severity = ?", severity).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// FindRecentNotifications 查找最近的通知
func (r *NotificationRepository) FindRecentNotifications(limit int) ([]*dns.Notification, error) {
	var notifications []*dns.Notification
	err := r.DB.Order("created_at DESC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

// FindFailedNotifications 查找发送失败的通知
func (r *NotificationRepository) FindFailedNotifications() ([]*dns.Notification, error) {
	var notifications []*dns.Notification
	err := r.DB.Where("status = ?", "failed").
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// CountByType 统计各类型的通知数量
func (r *NotificationRepository) CountByType() (map[string]int64, error) {
	var results []struct {
		Type  string
		Count int64
	}

	err := r.DB.Model(&dns.Notification{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	typeCount := make(map[string]int64)
	for _, result := range results {
		typeCount[result.Type] = result.Count
	}

	return typeCount, nil
}

// CountByStatus 统计各状态的通知数量
func (r *NotificationRepository) CountByStatus() (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	err := r.DB.Model(&dns.Notification{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	statusCount := make(map[string]int64)
	for _, result := range results {
		statusCount[result.Status] = result.Count
	}

	return statusCount, nil
}

// CountBySeverity 统计各严重程度的通知数量
func (r *NotificationRepository) CountBySeverity() (map[string]int64, error) {
	var results []struct {
		Severity string
		Count    int64
	}

	err := r.DB.Model(&dns.Notification{}).
		Select("severity, COUNT(*) as count").
		Group("severity").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	severityCount := make(map[string]int64)
	for _, result := range results {
		severityCount[result.Severity] = result.Count
	}

	return severityCount, nil
}

// ListWithPagination 分页查询通知
func (r *NotificationRepository) ListWithPagination(page, pageSize int, filters map[string]interface{}) ([]*dns.Notification, int64, error) {
	var notifications []*dns.Notification
	var total int64

	query := r.DB.Model(&dns.Notification{})

	// 应用过滤条件
	if notificationType, ok := filters["type"]; ok && notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if severity, ok := filters["severity"]; ok && severity != "" {
		query = query.Where("severity = ?", severity)
	}

	if keyword, ok := filters["keyword"]; ok && keyword != "" {
		query = query.Where("title LIKE ? OR message LIKE ?", 
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
		Find(&notifications).Error

	return notifications, total, err
}

// GetNotificationStatistics 获取通知统计信息
func (r *NotificationRepository) GetNotificationStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总通知数
	var totalCount int64
	err := r.DB.Model(&dns.Notification{}).Count(&totalCount).Error
	if err != nil {
		return nil, err
	}
	stats["total_notifications"] = totalCount

	// 今日通知数
	var todayCount int64
	err = r.DB.Model(&dns.Notification{}).
		Where("DATE(created_at) = CURDATE()").
		Count(&todayCount).Error
	if err != nil {
		return nil, err
	}
	stats["today_notifications"] = todayCount

	// 失败通知数
	var failedCount int64
	err = r.DB.Model(&dns.Notification{}).
		Where("status = ?", "failed").
		Count(&failedCount).Error
	if err != nil {
		return nil, err
	}
	stats["failed_notifications"] = failedCount

	// 按类型统计
	typeStats, err := r.CountByType()
	if err != nil {
		return nil, err
	}
	stats["type_distribution"] = typeStats

	// 按状态统计
	statusStats, err := r.CountByStatus()
	if err != nil {
		return nil, err
	}
	stats["status_distribution"] = statusStats

	// 按严重程度统计
	severityStats, err := r.CountBySeverity()
	if err != nil {
		return nil, err
	}
	stats["severity_distribution"] = severityStats

	// 最近7天的通知趋势
	var trendData []struct {
		Date  string
		Count int64
	}
	err = r.DB.Model(&dns.Notification{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&trendData).Error
	if err != nil {
		return nil, err
	}
	stats["weekly_trend"] = trendData

	return stats, nil
}

// MarkAsRead 标记通知为已读
func (r *NotificationRepository) MarkAsRead(id uint) error {
	return r.DB.Model(&dns.Notification{}).
		Where("id = ?", id).
		Update("read_at", "NOW()").Error
}

// MarkAllAsRead 标记所有通知为已读
func (r *NotificationRepository) MarkAllAsRead(tenantID uint) error {
	return r.DB.Model(&dns.Notification{}).
		Where("tenant_id = ?", tenantID).
		Where("read_at IS NULL").
		Update("read_at", "NOW()").Error
}

// DeleteOldNotifications 删除旧通知
func (r *NotificationRepository) DeleteOldNotifications(days int) error {
	return r.DB.Where("created_at < DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Delete(&dns.Notification{}).Error
}

// FindUnreadNotifications 查找未读通知
func (r *NotificationRepository) FindUnreadNotifications(tenantID uint) ([]*dns.Notification, error) {
	var notifications []*dns.Notification
	err := r.DB.Where("tenant_id = ?", tenantID).
		Where("read_at IS NULL").
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// CountUnreadNotifications 统计未读通知数量
func (r *NotificationRepository) CountUnreadNotifications(tenantID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&dns.Notification{}).
		Where("tenant_id = ?", tenantID).
		Where("read_at IS NULL").
		Count(&count).Error
	return count, err
}

// FindByTenantID 根据租户ID查找通知
func (r *NotificationRepository) FindByTenantID(tenantID uint) ([]*dns.Notification, error) {
	var notifications []*dns.Notification
	err := r.DB.Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// DeleteByTenantID 根据租户ID删除通知
func (r *NotificationRepository) DeleteByTenantID(tenantID uint) error {
	return r.DB.Where("tenant_id = ?", tenantID).Delete(&dns.Notification{}).Error
}

// UpdateNotificationStatus 更新通知状态
func (r *NotificationRepository) UpdateNotificationStatus(id uint, status string, errorMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if errorMsg != "" {
		updates["error_msg"] = errorMsg
	}

	return r.DB.Model(&dns.Notification{}).
		Where("id = ?", id).
		Updates(updates).Error
}
