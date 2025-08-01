package dns

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

// SyncZone 同步整个区域
func (d *AliyunDriver) SyncZone(ctx context.Context, zoneName string, options *SyncOptions) (*SyncResult, error) {
	start := time.Now()

	result := &SyncResult{
		Success:      true,
		TotalRecords: 0,
		Duration:     0,
		Errors:       []string{},
		Details:      make(map[string]string),
	}

	// 获取远程记录
	remoteRecords, err := d.ListRecords(ctx, zoneName, nil)
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to list remote records: %v", err))
		result.Duration = time.Since(start)
		return result, err
	}

	result.TotalRecords = len(remoteRecords)
	result.Duration = time.Since(start)
	result.Details["zone"] = zoneName
	result.Details["remote_records"] = fmt.Sprintf("%d", len(remoteRecords))

	d.logger.Info("Zone sync completed",
		zap.String("zone", zoneName),
		zap.Int("records", len(remoteRecords)),
		zap.Duration("duration", result.Duration))

	return result, nil
}

// CompareZone 比较区域
func (d *AliyunDriver) CompareZone(ctx context.Context, zoneName string, localRecords []*Record) (*ZoneComparison, error) {
	remoteRecords, err := d.ListRecords(ctx, zoneName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list remote records: %w", err)
	}

	comparison := &ZoneComparison{
		Domain:        zoneName,
		LocalRecords:  len(localRecords),
		RemoteRecords: len(remoteRecords),
		ToAdd:         []*Record{},
		ToUpdate:      []*Record{},
		ToDelete:      []*Record{},
		Conflicts:     []*Record{},
		ComparedAt:    time.Now(),
	}

	// 创建远程记录映射
	remoteMap := make(map[string]*Record)
	for i := range remoteRecords {
		key := fmt.Sprintf("%s:%s", remoteRecords[i].Name, remoteRecords[i].Type)
		remoteMap[key] = &remoteRecords[i]
	}

	// 比较本地记录
	for _, localRecord := range localRecords {
		key := fmt.Sprintf("%s:%s", localRecord.Name, localRecord.Type)
		if remoteRecord, exists := remoteMap[key]; exists {
			if localRecord.Value != remoteRecord.Value || localRecord.TTL != remoteRecord.TTL {
				comparison.ToUpdate = append(comparison.ToUpdate, localRecord)
			}
			delete(remoteMap, key)
		} else {
			comparison.ToAdd = append(comparison.ToAdd, localRecord)
		}
	}

	// 剩余的远程记录需要删除
	for _, remoteRecord := range remoteMap {
		comparison.ToDelete = append(comparison.ToDelete, remoteRecord)
	}

	return comparison, nil
}

// CreateTXTChallenge 创建TXT验证记录
func (d *AliyunDriver) CreateTXTChallenge(ctx context.Context, domain, token string, ttl int) (*ChallengeResult, error) {
	// 构造TXT记录
	record := &Record{
		Name:  "_acme-challenge",
		Type:  "TXT",
		Value: token,
		TTL:   ttl,
	}

	if ttl <= 0 {
		record.TTL = 600 // 默认10分钟
	}

	createdRecord, err := d.CreateRecord(ctx, domain, record)
	if err != nil {
		return nil, fmt.Errorf("failed to create TXT challenge: %w", err)
	}

	result := &ChallengeResult{
		Domain:    domain,
		Token:     token,
		RecordID:  createdRecord.ID,
		TTL:       record.TTL,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(record.TTL) * time.Second),
	}

	d.logger.Info("Created TXT challenge",
		zap.String("domain", domain),
		zap.String("record_id", createdRecord.ID))

	return result, nil
}

// DeleteTXTChallenge 删除TXT验证记录
func (d *AliyunDriver) DeleteTXTChallenge(ctx context.Context, domain, token string) error {
	// 查找TXT记录
	records, err := d.GetRecordsByName(ctx, domain, "_acme-challenge")
	if err != nil {
		return fmt.Errorf("failed to find TXT challenge records: %w", err)
	}

	for _, record := range records {
		if record.Type == "TXT" && record.Value == token {
			err := d.DeleteRecord(ctx, domain, record.ID)
			if err != nil {
				d.logger.Error("Failed to delete TXT challenge",
					zap.String("domain", domain),
					zap.String("record_id", record.ID),
					zap.Error(err))
				return fmt.Errorf("failed to delete TXT challenge: %w", err)
			}

			d.logger.Info("Deleted TXT challenge",
				zap.String("domain", domain),
				zap.String("record_id", record.ID))
			return nil
		}
	}

	return fmt.Errorf("TXT challenge record not found")
}

// ValidateChallenge 验证DNS-01挑战
func (d *AliyunDriver) ValidateChallenge(ctx context.Context, domain, token string) (*ChallengeValidation, error) {
	validation := &ChallengeValidation{
		Valid:         false,
		Propagated:    false,
		ExpectedValue: token,
		CheckedAt:     time.Now(),
		Servers:       []string{},
	}

	// 查找TXT记录
	records, err := d.GetRecordsByName(ctx, domain, "_acme-challenge")
	if err != nil {
		validation.Value = ""
		return validation, nil
	}

	for _, record := range records {
		if record.Type == "TXT" {
			validation.Value = record.Value
			validation.Valid = (record.Value == token)
			validation.Propagated = true
			break
		}
	}

	return validation, nil
}

// WaitForPropagation 等待DNS传播
func (d *AliyunDriver) WaitForPropagation(ctx context.Context, domain, recordType, expectedValue string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 检查记录是否已传播
		records, err := d.GetRecordsByType(ctx, domain, recordType)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		for _, record := range records {
			if record.Value == expectedValue {
				d.logger.Info("DNS propagation completed",
					zap.String("domain", domain),
					zap.String("type", recordType),
					zap.String("value", expectedValue))
				return nil
			}
		}

		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("DNS propagation timeout after %v", timeout)
}

// EnableDNSSEC 启用DNSSEC
func (d *AliyunDriver) EnableDNSSEC(ctx context.Context, zoneName string) (*DNSSECResult, error) {
	// 阿里云DNS不支持DNSSEC
	return nil, fmt.Errorf("DNSSEC not supported")
}

// DisableDNSSEC 禁用DNSSEC
func (d *AliyunDriver) DisableDNSSEC(ctx context.Context, zoneName string) error {
	// 阿里云DNS不支持DNSSEC
	return fmt.Errorf("DNSSEC not supported")
}

// GetDNSSECKeys 获取DNSSEC密钥
func (d *AliyunDriver) GetDNSSECKeys(ctx context.Context, zoneName string) ([]*DNSSECKey, error) {
	// 阿里云DNS不支持DNSSEC
	return nil, fmt.Errorf("DNSSEC not supported")
}

// RotateDNSSECKeys 轮换DNSSEC密钥
func (d *AliyunDriver) RotateDNSSECKeys(ctx context.Context, zoneName string) (*DNSSECResult, error) {
	// 阿里云DNS不支持DNSSEC
	return nil, fmt.Errorf("DNSSEC not supported")
}

// GetZoneFile 获取区域文件
func (d *AliyunDriver) GetZoneFile(ctx context.Context, zoneName string) (string, error) {
	records, err := d.ListRecords(ctx, zoneName, nil)
	if err != nil {
		return "", fmt.Errorf("failed to list records: %w", err)
	}

	var zoneFile strings.Builder
	zoneFile.WriteString(fmt.Sprintf("; Zone file for %s\n", zoneName))
	zoneFile.WriteString(fmt.Sprintf("; Generated at %s\n", time.Now().Format(time.RFC3339)))
	zoneFile.WriteString("\n")

	for _, record := range records {
		line := fmt.Sprintf("%s\t%d\tIN\t%s\t%s\n",
			record.Name, record.TTL, record.Type, record.Value)
		zoneFile.WriteString(line)
	}

	return zoneFile.String(), nil
}

// ImportZoneFile 导入区域文件
func (d *AliyunDriver) ImportZoneFile(ctx context.Context, zoneName string, zoneFile string) (*ImportResult, error) {
	// 简化实现，实际应该解析zone file格式
	return &ImportResult{
		Total:    0,
		Success:  0,
		Failed:   0,
		Skipped:  0,
		Results:  []*OperationResult{},
		Duration: 0,
	}, fmt.Errorf("unsupported operation: import zone file")
}

// ExportZoneFile 导出区域文件
func (d *AliyunDriver) ExportZoneFile(ctx context.Context, zoneName string, format string) (string, error) {
	if format != "bind" && format != "rfc1035" {
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	return d.GetZoneFile(ctx, zoneName)
}

// GetRecordsByType 根据类型获取记录
func (d *AliyunDriver) GetRecordsByType(ctx context.Context, zoneName, recordType string) ([]Record, error) {
	options := &ListOptions{
		Filter: map[string]string{
			"type": recordType,
		},
	}

	return d.ListRecords(ctx, zoneName, options)
}

// GetRecordsByName 根据名称获取记录
func (d *AliyunDriver) GetRecordsByName(ctx context.Context, zoneName, recordName string) ([]Record, error) {
	options := &ListOptions{
		Filter: map[string]string{
			"keyword": recordName,
		},
	}

	return d.ListRecords(ctx, zoneName, options)
}

// SearchRecords 搜索记录
func (d *AliyunDriver) SearchRecords(ctx context.Context, zoneName string, query map[string]string) ([]Record, error) {
	options := &ListOptions{
		Filter: query,
	}

	return d.ListRecords(ctx, zoneName, options)
}

// ValidateRecord 验证记录格式
func (d *AliyunDriver) ValidateRecord(record *Record) error {
	if record.Name == "" {
		return fmt.Errorf("invalid record: name is required")
	}

	if record.Type == "" {
		return fmt.Errorf("invalid record: type is required")
	}

	if record.Value == "" {
		return fmt.Errorf("invalid record: value is required")
	}

	// 检查记录类型是否支持
	supported := false
	for _, recordType := range d.GetSupportedRecordTypes() {
		if record.Type == recordType {
			supported = true
			break
		}
	}

	if !supported {
		return fmt.Errorf("unsupported record type: %s", record.Type)
	}

	return nil
}

// ValidateZone 验证区域
func (d *AliyunDriver) ValidateZone(ctx context.Context, zoneName string) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:   true,
		Details: make(map[string]string),
	}

	// 检查区域是否存在
	_, err := d.GetZone(ctx, zoneName)
	if err != nil {
		result.Valid = false
		result.ErrorMsg = fmt.Sprintf("Zone validation failed: %v", err)
		return result, nil
	}

	result.Details["message"] = "Zone is valid"
	return result, nil
}

// CheckRecordConflicts 检查记录冲突
func (d *AliyunDriver) CheckRecordConflicts(ctx context.Context, zoneName string, record *Record) ([]string, error) {
	conflicts := []string{}

	// 获取同名同类型的记录
	existingRecords, err := d.GetRecordsByName(ctx, zoneName, record.Name)
	if err != nil {
		return conflicts, err
	}

	for _, existing := range existingRecords {
		if existing.Type == record.Type && existing.ID != record.ID {
			if record.Type == "CNAME" || existing.Type == "CNAME" {
				conflicts = append(conflicts, fmt.Sprintf("CNAME conflict with %s record", existing.Type))
			}
		}
	}

	return conflicts, nil
}

// GetStatistics 获取统计信息
func (d *AliyunDriver) GetStatistics(ctx context.Context, zoneName string) (*Statistics, error) {
	records, err := d.ListRecords(ctx, zoneName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	recordsByType := make(map[string]int)
	for _, record := range records {
		recordsByType[record.Type]++
	}

	stats := &Statistics{
		Domain:        zoneName,
		TotalRecords:  len(records),
		RecordsByType: recordsByType,
		LastSync:      time.Now(),
		SyncCount:     1,
		ErrorCount:    0,
		Metrics:       make(map[string]interface{}),
	}

	return stats, nil
}

// GetQuota 获取配额信息
func (d *AliyunDriver) GetQuota(ctx context.Context) (*Quota, error) {
	quota := &Quota{
		Provider:      "aliyun",
		TotalDomains:  d.info.Limits["max_domains"],
		UsedDomains:   0, // 需要实际查询
		TotalRecords:  d.info.Limits["max_records_per_domain"],
		UsedRecords:   0, // 需要实际查询
		RateLimit:     d.info.Limits["rate_limit_per_second"],
		RateRemaining: d.info.Limits["rate_limit_per_second"],
		ResetTime:     time.Now().Add(time.Hour),
		Features:      make(map[string]bool),
		Limits:        d.info.Limits,
	}

	// 设置功能支持
	for _, feature := range d.info.Features {
		quota.Features[feature] = true
	}

	return quota, nil
}

// GetUsage 获取使用情况
func (d *AliyunDriver) GetUsage(ctx context.Context) (map[string]interface{}, error) {
	usage := make(map[string]interface{})

	// 获取域名列表来计算使用量
	zones, err := d.ListZones(ctx, nil)
	if err != nil {
		return usage, err
	}

	usage["domains_used"] = len(zones)
	usage["domains_limit"] = d.info.Limits["max_domains"]

	totalRecords := 0
	for _, zone := range zones {
		records, err := d.ListRecords(ctx, zone.Name, nil)
		if err == nil {
			totalRecords += len(records)
		}
	}

	usage["records_used"] = totalRecords
	usage["records_limit"] = d.info.Limits["max_records_per_domain"] * len(zones)

	return usage, nil
}

// HealthCheck 健康检查
func (d *AliyunDriver) HealthCheck(ctx context.Context) *TestResult {
	return d.Test(ctx)
}

// GetMetrics 获取指标
func (d *AliyunDriver) GetMetrics(ctx context.Context) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	// 基础指标
	metrics["provider"] = "aliyun"
	metrics["api_version"] = d.info.Metadata["api_version"]
	metrics["last_check"] = time.Now()

	// 获取使用情况作为指标
	usage, err := d.GetUsage(ctx)
	if err == nil {
		for k, v := range usage {
			metrics[k] = v
		}
	}

	return metrics, nil
}

// SetRecordComment 设置记录注释
func (d *AliyunDriver) SetRecordComment(ctx context.Context, zoneName, recordID, comment string) error {
	// 阿里云DNS不支持记录注释
	return fmt.Errorf("unsupported operation: set record comment")
}

// GetRecordHistory 获取记录历史
func (d *AliyunDriver) GetRecordHistory(ctx context.Context, zoneName, recordID string) ([]interface{}, error) {
	// 阿里云DNS不支持记录历史
	return nil, fmt.Errorf("unsupported operation: get record history")
}

// CreateRecordSet 创建记录集
func (d *AliyunDriver) CreateRecordSet(ctx context.Context, zoneName string, records []*Record) ([]*Record, error) {
	result := make([]*Record, 0, len(records))

	for _, record := range records {
		createdRecord, err := d.CreateRecord(ctx, zoneName, record)
		if err != nil {
			return result, err
		}
		result = append(result, createdRecord)
	}

	return result, nil
}

// UpdateRecordSet 更新记录集
func (d *AliyunDriver) UpdateRecordSet(ctx context.Context, zoneName string, records []*Record) ([]*Record, error) {
	result := make([]*Record, 0, len(records))

	for _, record := range records {
		updatedRecord, err := d.UpdateRecord(ctx, zoneName, record)
		if err != nil {
			return result, err
		}
		result = append(result, updatedRecord)
	}

	return result, nil
}
