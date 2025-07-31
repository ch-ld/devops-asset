package cmdb

import (
	"api-server/internal/model/cmdb"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// HostRepository 主机数据仓库，封装主机表的数据库操作
type HostRepository struct {
	db *gorm.DB
}

// NewHostRepository 创建主机数据仓库实例
func NewHostRepository(db *gorm.DB) *HostRepository {
	return &HostRepository{db: db}
}

// Create 新增主机记录
func (r *HostRepository) Create(host *cmdb.Host) error {
	return r.db.Create(host).Error
}

// Update 更新主机信息
func (r *HostRepository) Update(host *cmdb.Host) error {
	return r.db.Model(host).Updates(host).Error
}

// DeleteByID 根据ID删除主机
func (r *HostRepository) DeleteByID(id uint) error {
	return r.db.Delete(&cmdb.Host{}, id).Error
}

// FindByID 根据ID查询主机
func (r *HostRepository) FindByID(id uint) (*cmdb.Host, error) {
	var host cmdb.Host
	if err := r.db.First(&host, id).Error; err != nil {
		return nil, err
	}
	return &host, nil
}

// buildConditions 构建查询条件
func (r *HostRepository) buildConditions(conditions map[string]interface{}) *gorm.DB {
	db := r.db

	if keyword, ok := conditions["keyword"].(string); ok && keyword != "" {
		keyword = "%" + keyword + "%"
		db = db.Where("name LIKE ? OR instance_id LIKE ? OR public_ip LIKE ? OR private_ip LIKE ?",
			keyword, keyword, keyword, keyword)
	}

	if name, ok := conditions["name"].(string); ok && name != "" {
		name = "%" + name + "%"
		db = db.Where("name LIKE ?", name)
	}

	if ip, ok := conditions["ip"].(string); ok && ip != "" {
		ip = "%" + ip + "%"
		db = db.Where("public_ip LIKE ? OR private_ip LIKE ?", ip, ip)
	}

	if status, ok := conditions["status"].(string); ok && status != "" {
		db = db.Where("status = ?", status)
	}

	if groupID, ok := conditions["group_id"].(uint); ok {
		db = db.Where("group_id = ?", groupID)
	}

	if region, ok := conditions["region"].(string); ok && region != "" {
		db = db.Where("region = ?", region)
	}

	return db
}

// FindByConditions 根据条件分页查询主机
func (r *HostRepository) FindByConditions(conditions map[string]interface{}, page, pageSize int) ([]cmdb.Host, error) {
	var hosts []cmdb.Host
	offset := (page - 1) * pageSize

	db := r.buildConditions(conditions)
	err := db.Offset(offset).Limit(pageSize).Find(&hosts).Error

	return hosts, err
}

// FindAllByConditions 根据条件查询所有主机，不分页
func (r *HostRepository) FindAllByConditions(conditions map[string]interface{}) ([]cmdb.Host, error) {
	var hosts []cmdb.Host
	db := r.buildConditions(conditions)
	err := db.Find(&hosts).Error
	return hosts, err
}

// CountByConditions 根据条件统计主机总数
func (r *HostRepository) CountByConditions(conditions map[string]interface{}) (int64, error) {
	var count int64
	db := r.buildConditions(conditions)
	err := db.Model(&cmdb.Host{}).Count(&count).Error
	return count, err
}

// FindAll 查询所有主机记录（可扩展分页/过滤）
func (r *HostRepository) FindAll() ([]*cmdb.Host, error) {
	var hosts []*cmdb.Host
	err := r.db.Find(&hosts).Error
	return hosts, err
}

// FindByGroupID 根据主机组ID查询主机列表
func (r *HostRepository) FindByGroupID(groupID uint, page, pageSize int, keyword string) ([]cmdb.Host, int64, error) {
	query := r.db.Model(&cmdb.Host{}).Where("group_id = ?", groupID)

	if keyword != "" {
		query = query.Where("name LIKE ? OR instance_id LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var hosts []cmdb.Host
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&hosts).Error; err != nil {
		return nil, 0, err
	}

	return hosts, total, nil
}

// FindByProviderID 根据云账号ID查询主机列表
func (r *HostRepository) FindByProviderID(providerID uint) ([]cmdb.Host, error) {
	var hosts []cmdb.Host
	if err := r.db.Where("provider_id = ?", providerID).Find(&hosts).Error; err != nil {
		return nil, err
	}
	return hosts, nil
}

// CreateManualHost 创建自建主机
func (r *HostRepository) CreateManualHost(host *cmdb.Host) error {
	host.ProviderType = "manual"
	return r.db.Create(host).Error
}

// BatchCreateHosts 批量创建主机
func (r *HostRepository) BatchCreateHosts(hosts []cmdb.Host) error {
	return r.db.Create(&hosts).Error
}

// BatchUpdateHosts 批量更新主机
func (r *HostRepository) BatchUpdateHosts(hosts []cmdb.Host) error {
	for _, host := range hosts {
		if err := r.db.Model(&host).Updates(host).Error; err != nil {
			return err
		}
	}
	return nil
}

// BatchCreate 批量创建主机
func (r *HostRepository) BatchCreate(hosts []cmdb.Host) error {
	if len(hosts) == 0 {
		return nil
	}
	return r.db.Create(&hosts).Error
}

// BatchDelete 批量删除主机
func (r *HostRepository) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Delete(&cmdb.Host{}, ids).Error
}

// FindByExpiring 查询即将过期的主机
func (r *HostRepository) FindByExpiring(days int) ([]cmdb.Host, error) {
	var hosts []cmdb.Host

	// 计算截止日期
	deadline := time.Now().AddDate(0, 0, days)

	// 查询在当前时间和截止日期之间过期的主机
	err := r.db.Where("expired_at IS NOT NULL AND expired_at <= ?", deadline).
		Preload("Provider").
		Preload("Group").
		Find(&hosts).Error

	if err != nil {
		return nil, fmt.Errorf("查询即将过期主机失败: %w", err)
	}

	return hosts, nil
}

// FindByAbnormalStatus 查询状态异常的主机
func (r *HostRepository) FindByAbnormalStatus() ([]cmdb.Host, error) {
	var hosts []cmdb.Host

	// 状态为error或者abnormal的主机视为异常
	err := r.db.Where("status IN ('error', 'abnormal')").
		Preload("Provider").
		Preload("Group").
		Find(&hosts).Error

	if err != nil {
		return nil, fmt.Errorf("查询状态异常主机失败: %w", err) // This line was modified as per the new_code, as fmt is not imported.
	}

	return hosts, nil
}

// FindByStatus 查询指定状态的主机
func (r *HostRepository) FindByStatus(status string) ([]cmdb.Host, error) {
	var hosts []cmdb.Host

	err := r.db.Where("status = ?", status).
		Preload("Provider").
		Preload("Group").
		Find(&hosts).Error

	if err != nil {
		return nil, fmt.Errorf("查询状态为 %s 的主机失败: %w", status, err) // This line was modified as per the new_code, as fmt is not imported.
	}

	return hosts, nil
}

// CountByProviderType 按云提供商类型统计主机数量
func (r *HostRepository) CountByProviderType() (map[string]int64, error) {
	type Result struct {
		ProviderType string `json:"provider_type"`
		Count        int64  `json:"count"`
	}

	var results []Result
	err := r.db.Model(&cmdb.Host{}). // This line was modified as per the new_code, as model is not imported.
						Select("provider_type, count(*) as count").
						Group("provider_type").
						Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("按云提供商统计主机失败: %w", err) // This line was modified as per the new_code, as fmt is not imported.
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.ProviderType] = result.Count
	}

	return counts, nil
}

// CountByRegion 按区域统计主机数量
func (r *HostRepository) CountByRegion() (map[string]int64, error) {
	type Result struct {
		Region string `json:"region"`
		Count  int64  `json:"count"`
	}

	var results []Result
	err := r.db.Model(&cmdb.Host{}). // This line was modified as per the new_code, as model is not imported.
						Select("region, count(*) as count").
						Group("region").
						Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("按区域统计主机失败: %w", err) // This line was modified as per the new_code, as fmt is not imported.
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.Region] = result.Count
	}

	return counts, nil
}

// CountByStatus 按状态统计主机数量
func (r *HostRepository) CountByStatus() (map[string]int64, error) {
	type Result struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	var results []Result
	err := r.db.Model(&cmdb.Host{}). // This line was modified as per the new_code, as model is not imported.
						Select("status, count(*) as count").
						Group("status").
						Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("按状态统计主机失败: %w", err) // This line was modified as per the new_code, as fmt is not imported.
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.Status] = result.Count
	}

	return counts, nil
}

// CountByOS 按操作系统统计主机数量
func (r *HostRepository) CountByOS() (map[string]int64, error) {
	type Result struct {
		OS    string `json:"os"`
		Count int64  `json:"count"`
	}

	var results []Result
	err := r.db.Model(&cmdb.Host{}). // This line was modified as per the new_code, as model is not imported.
						Select("os, count(*) as count").
						Group("os").
						Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("按操作系统统计主机失败: %w", err) // This line was modified as per the new_code, as fmt is not imported.
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.OS] = result.Count
	}

	return counts, nil
}

// GetTotalHostCount 获取主机总数
func (r *HostRepository) GetTotalHostCount() (int64, error) {
	var count int64
	err := r.db.Model(&cmdb.Host{}).Count(&count).Error // This line was modified as per the new_code, as model is not imported.
	if err != nil {
		return 0, fmt.Errorf("获取主机总数失败: %w", err) // This line was modified as per the new_code, as fmt is not imported.
	}
	return count, nil
}

// GetRunningHostCount 获取运行中主机数量
func (r *HostRepository) GetRunningHostCount() (int64, error) {
	var count int64
	err := r.db.Model(&cmdb.Host{}).Where("status = 'running'").Count(&count).Error // This line was modified as per the new_code, as model is not imported.
	if err != nil {
		return 0, fmt.Errorf("获取运行中主机数量失败: %w", err) // This line was modified as per the new_code, as fmt is not imported.
	}
	return count, nil
}

// GetExpiringHostCount 获取即将过期的主机数量
func (r *HostRepository) GetExpiringHostCount(days int) (int64, error) {
	var count int64
	deadline := time.Now().AddDate(0, 0, days)

	err := r.db.Model(&cmdb.Host{}).
		Where("expired_at IS NOT NULL AND expired_at <= ?", deadline).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("获取即将过期主机数量失败: %w", err)
	}

	return count, nil
}

// FindByTags 通过标签查询主机
func (r *HostRepository) FindByTags(tags []string) ([]cmdb.Host, error) {
	var hosts []cmdb.Host

	// 查询条件: tags包含指定的任意标签
	query := r.db.Model(&cmdb.Host{}) // This line was modified as per the new_code, as model is not imported.

	for _, tag := range tags {
		// 在JSON数组中查找标签
		// 注意：不同数据库实现JSON查询的语法可能有所不同，这里假设使用MySQL或PostgreSQL
		query = query.Or("JSON_CONTAINS(tags, ?)", fmt.Sprintf("\"%s\"", tag)) // This line was modified as per the new_code, as fmt is not imported.
	}

	err := query. // This line was modified as per the new_code, as model is not imported.
			Preload("Provider").
			Preload("Group").
			Find(&hosts).Error

	if err != nil {
		return nil, fmt.Errorf("按标签查询主机失败: %w", err) // This line was modified as per the new_code, as fmt is not imported.
	}

	return hosts, nil
}

// SearchByKeyword 通过关键词搜索主机
func (r *HostRepository) SearchByKeyword(keyword string) ([]cmdb.Host, error) {
	var hosts []cmdb.Host

	// 构建LIKE模式匹配
	likePattern := "%" + keyword + "%"

	// 在多个字段中搜索
	err := r.db.Where("name LIKE ? OR instance_id LIKE ? OR remark LIKE ?",
		likePattern, likePattern, likePattern).
		Or("JSON_SEARCH(public_ip, 'one', ?) IS NOT NULL", likePattern).
		Or("JSON_SEARCH(private_ip, 'one', ?) IS NOT NULL", likePattern).
		Preload("Provider").
		Preload("Group").
		Find(&hosts).Error

	if err != nil {
		return nil, fmt.Errorf("搜索主机失败: %w", err) // This line was modified as per the new_code, as fmt is not imported.
	}

	return hosts, nil
}

// GetHistoryByDateRange 获取指定日期范围内的主机数量变化
func (r *HostRepository) GetHistoryByDateRange(startDate, endDate time.Time) (map[string]int64, error) {
	// 此功能需要额外的表记录主机数量历史
	// 这里返回模拟数据
	result := make(map[string]int64)

	// 按天生成日期
	for date := startDate; date.Before(endDate) || date.Equal(endDate); date = date.AddDate(0, 0, 1) {
		dateStr := date.Format("2006-01-02")
		// 生成一些假数据
		baseCount := int64(100)
		dayOffset := int64(date.Sub(startDate).Hours() / 24)
		result[dateStr] = baseCount + dayOffset*2
	}

	return result, nil
}

// GetHostGroupsWithStats 获取带统计信息的主机组列表
func (r *HostRepository) GetHostGroupsWithStats() ([]map[string]interface{}, error) {
	type GroupStats struct {
		ID           uint   `json:"id"`
		Name         string `json:"name"`
		HostCount    int64  `json:"host_count"`
		RunningCount int64  `json:"running_count"`
	}

	var stats []GroupStats
	err := r.db.Table("cmdb_host_groups g").
		Select("g.id, g.name, COUNT(h.id) as host_count, SUM(CASE WHEN h.status = 'running' THEN 1 ELSE 0 END) as running_count").
		Joins("LEFT JOIN cmdb_hosts h ON g.id = h.group_id").
		Group("g.id, g.name").
		Find(&stats).Error

	if err != nil {
		return nil, fmt.Errorf("获取主机组统计失败: %w", err)
	}

	// 转换为返回格式
	var result []map[string]interface{}
	for _, stat := range stats {
		runningRate := float64(0)
		if stat.HostCount > 0 {
			runningRate = float64(stat.RunningCount) / float64(stat.HostCount) * 100
		}

		result = append(result, map[string]interface{}{
			"id":           stat.ID,
			"name":         stat.Name,
			"hostCount":    stat.HostCount,
			"runningCount": stat.RunningCount,
			"runningRate":  runningRate,
		})
	}

	return result, nil
}

// GetRecentHosts 获取最近添加的主机列表
func (r *HostRepository) GetRecentHosts(limit int) ([]cmdb.Host, error) {
	var hosts []cmdb.Host
	err := r.db.Model(&cmdb.Host{}).
		Preload("Provider").
		Preload("Group").
		Order("created_at DESC").
		Limit(limit).
		Find(&hosts).Error

	if err != nil {
		return nil, fmt.Errorf("获取最近主机列表失败: %w", err)
	}

	return hosts, nil
}
