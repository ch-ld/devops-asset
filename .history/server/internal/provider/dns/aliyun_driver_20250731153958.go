package dns

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"go.uber.org/zap"
)

// AliyunDriver 阿里云DNS驱动
type AliyunDriver struct {
	client *alidns.Client
	config *AliyunConfig
	logger *zap.Logger
	info   *ProviderInfo
}

// AliyunConfig 阿里云配置
type AliyunConfig struct {
	*BaseConfig
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

// NewAliyunDriver 创建阿里云DNS驱动
func NewAliyunDriver(config Config) (Driver, error) {
	aliyunConfig, ok := config.(*AliyunConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type for Aliyun driver")
	}

	if err := aliyunConfig.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	client, err := alidns.NewClientWithAccessKey(
		aliyunConfig.GetRegion(),
		aliyunConfig.AccessKeyID,
		aliyunConfig.AccessKeySecret,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Aliyun DNS client: %w", err)
	}

	// 设置超时
	client.SetConnectTimeout(aliyunConfig.GetTimeout())
	client.SetReadTimeout(aliyunConfig.GetTimeout())

	driver := &AliyunDriver{
		client: client,
		config: aliyunConfig,
		logger: zap.L().Named("aliyun-dns"),
		info: &ProviderInfo{
			Name:    "Aliyun DNS",
			Type:    "aliyun",
			Version: "1.0.0",
			Features: []string{
				"dns_management",
				"batch_operations",
				"txt_challenges",
				"zone_sync",
				"dnssec",
				"statistics",
			},
			Limits: map[string]int{
				"max_domains":            100,
				"max_records_per_domain": 10000,
				"max_batch_size":         100,
				"rate_limit_per_second":  20,
			},
			Regions: []string{
				"cn-hangzhou",
				"cn-beijing",
				"cn-shanghai",
				"cn-shenzhen",
			},
			RecordTypes: []string{
				"A", "AAAA", "CNAME", "MX", "TXT", "NS", "SRV", "CAA", "PTR",
			},
			Metadata: map[string]string{
				"api_version": "2015-01-09",
				"endpoint":    "https://alidns.aliyuncs.com",
			},
		},
	}

	return driver, nil
}

// GetInfo 获取提供商信息
func (d *AliyunDriver) GetInfo() *ProviderInfo {
	return d.info
}

// GetCapabilities 获取功能列表
func (d *AliyunDriver) GetCapabilities() []string {
	return d.info.Features
}

// GetSupportedRecordTypes 获取支持的记录类型
func (d *AliyunDriver) GetSupportedRecordTypes() []string {
	return d.info.RecordTypes
}

// Test 测试连接
func (d *AliyunDriver) Test(ctx context.Context) *TestResult {
	start := time.Now()

	request := alidns.CreateDescribeDomainsRequest()
	request.Scheme = "https"
	request.PageSize = "1"

	_, err := d.client.DescribeDomains(request)
	latency := time.Since(start)

	result := &TestResult{
		TestedAt: time.Now(),
		TestType: "api_call",
		Endpoint: d.info.Metadata["endpoint"],
		Latency:  latency,
		Details:  make(map[string]string),
	}

	if err != nil {
		result.Success = false
		result.ErrorMsg = err.Error()
		result.StatusCode = 500
		d.logger.Error("Connection test failed", zap.Error(err))
	} else {
		result.Success = true
		result.StatusCode = 200
		result.Details["message"] = "Connection successful"
		d.logger.Info("Connection test successful", zap.Duration("latency", latency))
	}

	return result
}

// ValidateCredentials 验证凭证
func (d *AliyunDriver) ValidateCredentials(ctx context.Context, credentials map[string]string) *ValidationResult {
	result := &ValidationResult{
		Details: make(map[string]string),
	}

	accessKeyID, ok := credentials["access_key_id"]
	if !ok || accessKeyID == "" {
		result.Valid = false
		result.ErrorMsg = "access_key_id is required"
		result.Suggestions = append(result.Suggestions, "Please provide a valid access_key_id")
		return result
	}

	accessKeySecret, ok := credentials["access_key_secret"]
	if !ok || accessKeySecret == "" {
		result.Valid = false
		result.ErrorMsg = "access_key_secret is required"
		result.Suggestions = append(result.Suggestions, "Please provide a valid access_key_secret")
		return result
	}

	// 尝试创建客户端并测试
	testClient, err := alidns.NewClientWithAccessKey("cn-hangzhou", accessKeyID, accessKeySecret)
	if err != nil {
		result.Valid = false
		result.ErrorMsg = fmt.Sprintf("Failed to create client: %v", err)
		result.Suggestions = append(result.Suggestions, "Check your access key and secret")
		return result
	}

	// 测试API调用
	request := alidns.CreateDescribeDomainsRequest()
	request.Scheme = "https"
	request.PageSize = "1"

	_, err = testClient.DescribeDomains(request)
	if err != nil {
		result.Valid = false
		result.ErrorMsg = fmt.Sprintf("API test failed: %v", err)
		result.Suggestions = append(result.Suggestions, "Verify your credentials have DNS permissions")
		return result
	}

	result.Valid = true
	result.Details["message"] = "Credentials are valid"
	return result
}

// ListZones 列出所有DNS区域
func (d *AliyunDriver) ListZones(ctx context.Context, options *ListOptions) ([]Zone, error) {
	request := alidns.CreateDescribeDomainsRequest()
	request.Scheme = "https"

	if options != nil {
		if options.PageSize > 0 {
			request.PageSize = requests.NewInteger(options.PageSize)
		} else {
			request.PageSize = requests.NewInteger(100)
		}

		if options.Page > 0 {
			request.PageNumber = requests.NewInteger(options.Page)
		} else {
			request.PageNumber = requests.NewInteger(1)
		}

		// 处理过滤条件
		if keyword, ok := options.Filter["keyword"]; ok {
			request.KeyWord = keyword
		}
	} else {
		request.PageSize = "100"
		request.PageNumber = "1"
	}

	response, err := d.client.DescribeDomains(request)
	if err != nil {
		d.logger.Error("Failed to list zones", zap.Error(err))
		return nil, fmt.Errorf("failed to list zones: %w", err)
	}

	var zones []Zone
	for _, domain := range response.Domains.Domain {
		status := "active"
		if domain.InstanceId == "" {
			status = "inactive"
		}

		zones = append(zones, Zone{
			ID:     domain.DomainId,
			Name:   domain.DomainName,
			Status: status,
		})
	}

	d.logger.Info("Listed zones successfully", zap.Int("count", len(zones)))
	return zones, nil
}

// GetZone 获取指定DNS区域
func (d *AliyunDriver) GetZone(ctx context.Context, zoneName string) (*Zone, error) {
	request := alidns.CreateDescribeDomainInfoRequest()
	request.Scheme = "https"
	request.DomainName = zoneName

	response, err := d.client.DescribeDomainInfo(request)
	if err != nil {
		d.logger.Error("Failed to get zone", zap.String("zone", zoneName), zap.Error(err))
		return nil, fmt.Errorf("failed to get zone %s: %w", zoneName, err)
	}

	status := "active"
	if response.InstanceId == "" {
		status = "inactive"
	}

	zone := &Zone{
		ID:     response.DomainId,
		Name:   response.DomainName,
		Status: status,
	}

	d.logger.Info("Got zone successfully", zap.String("zone", zoneName))
	return zone, nil
}

// CreateZone 创建DNS区域
func (d *AliyunDriver) CreateZone(ctx context.Context, zoneName string) (*Zone, error) {
	// 阿里云DNS不支持通过API创建域名，域名需要在控制台添加
	return nil, fmt.Errorf("unsupported operation: create zone")
}

// UpdateZone 更新DNS区域
func (d *AliyunDriver) UpdateZone(ctx context.Context, zone *Zone) (*Zone, error) {
	// 阿里云DNS不支持通过API更新域名信息
	return nil, fmt.Errorf("unsupported operation: update zone")
}

// DeleteZone 删除DNS区域
func (d *AliyunDriver) DeleteZone(ctx context.Context, zoneName string) error {
	// 阿里云DNS不支持通过API删除域名
	return fmt.Errorf("unsupported operation: delete zone")
}

// ListRecords 列出指定区域的所有记录
func (d *AliyunDriver) ListRecords(ctx context.Context, zoneName string, options *ListOptions) ([]Record, error) {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = zoneName

	if options != nil {
		if options.PageSize > 0 {
			request.PageSize = strconv.Itoa(options.PageSize)
		} else {
			request.PageSize = "500"
		}

		if options.Page > 0 {
			request.PageNumber = strconv.Itoa(options.Page)
		} else {
			request.PageNumber = "1"
		}

		// 处理过滤条件
		if recordType, ok := options.Filter["type"]; ok {
			request.Type = recordType
		}
		if keyword, ok := options.Filter["keyword"]; ok {
			request.KeyWord = keyword
		}
	} else {
		request.PageSize = "500"
		request.PageNumber = "1"
	}

	response, err := d.client.DescribeDomainRecords(request)
	if err != nil {
		d.logger.Error("Failed to list records", zap.String("zone", zoneName), zap.Error(err))
		return nil, fmt.Errorf("failed to list records for zone %s: %w", zoneName, err)
	}

	var records []Record
	for _, record := range response.DomainRecords.Record {
		r := Record{
			ID:    record.RecordId,
			Name:  record.RR,
			Type:  record.Type,
			Value: record.Value,
			TTL:   int(record.TTL),
		}

		// 处理优先级
		if record.Priority != "" {
			if priority, err := strconv.Atoi(record.Priority); err == nil {
				r.Priority = &priority
			}
		}

		// 处理权重
		if record.Weight != 0 {
			weight := int(record.Weight)
			r.Weight = &weight
		}

		records = append(records, r)
	}

	d.logger.Info("Listed records successfully",
		zap.String("zone", zoneName),
		zap.Int("count", len(records)))

	return records, nil
}

// GetRecord 获取指定记录
func (d *AliyunDriver) GetRecord(ctx context.Context, zoneName, recordID string) (*Record, error) {
	request := alidns.CreateDescribeDomainRecordInfoRequest()
	request.Scheme = "https"
	request.RecordId = recordID

	response, err := d.client.DescribeDomainRecordInfo(request)
	if err != nil {
		d.logger.Error("Failed to get record",
			zap.String("zone", zoneName),
			zap.String("record_id", recordID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get record %s: %w", recordID, err)
	}

	record := &Record{
		ID:    response.RecordId,
		Name:  response.RR,
		Type:  response.Type,
		Value: response.Value,
		TTL:   int(response.TTL),
	}

	// 处理优先级
	if response.Priority != "" {
		if priority, err := strconv.Atoi(response.Priority); err == nil {
			record.Priority = &priority
		}
	}

	d.logger.Info("Got record successfully",
		zap.String("zone", zoneName),
		zap.String("record_id", recordID))

	return record, nil
}

// CreateRecord 创建DNS记录
func (d *AliyunDriver) CreateRecord(ctx context.Context, zoneName string, record *Record) (*Record, error) {
	request := alidns.CreateAddDomainRecordRequest()
	request.Scheme = "https"
	request.DomainName = zoneName
	request.RR = record.Name
	request.Type = record.Type
	request.Value = record.Value

	if record.TTL > 0 {
		request.TTL = strconv.Itoa(record.TTL)
	}

	if record.Priority != nil {
		request.Priority = strconv.Itoa(*record.Priority)
	}

	response, err := d.client.AddDomainRecord(request)
	if err != nil {
		d.logger.Error("Failed to create record",
			zap.String("zone", zoneName),
			zap.String("name", record.Name),
			zap.String("type", record.Type),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create record: %w", err)
	}

	// 返回创建的记录
	createdRecord := &Record{
		ID:       response.RecordId,
		Name:     record.Name,
		Type:     record.Type,
		Value:    record.Value,
		TTL:      record.TTL,
		Priority: record.Priority,
	}

	d.logger.Info("Created record successfully",
		zap.String("zone", zoneName),
		zap.String("record_id", response.RecordId),
		zap.String("name", record.Name),
		zap.String("type", record.Type))

	return createdRecord, nil
}

// UpdateRecord 更新DNS记录
func (d *AliyunDriver) UpdateRecord(ctx context.Context, zoneName string, record *Record) (*Record, error) {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = record.ID
	request.RR = record.Name
	request.Type = record.Type
	request.Value = record.Value

	if record.TTL > 0 {
		request.TTL = strconv.Itoa(record.TTL)
	}

	if record.Priority != nil {
		request.Priority = strconv.Itoa(*record.Priority)
	}

	_, err := d.client.UpdateDomainRecord(request)
	if err != nil {
		d.logger.Error("Failed to update record",
			zap.String("zone", zoneName),
			zap.String("record_id", record.ID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to update record %s: %w", record.ID, err)
	}

	d.logger.Info("Updated record successfully",
		zap.String("zone", zoneName),
		zap.String("record_id", record.ID))

	return record, nil
}

// DeleteRecord 删除DNS记录
func (d *AliyunDriver) DeleteRecord(ctx context.Context, zoneName, recordID string) error {
	request := alidns.CreateDeleteDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = recordID

	_, err := d.client.DeleteDomainRecord(request)
	if err != nil {
		d.logger.Error("Failed to delete record",
			zap.String("zone", zoneName),
			zap.String("record_id", recordID),
			zap.Error(err))
		return fmt.Errorf("failed to delete record %s: %w", recordID, err)
	}

	d.logger.Info("Deleted record successfully",
		zap.String("zone", zoneName),
		zap.String("record_id", recordID))

	return nil
}

// BatchCreateRecords 批量创建DNS记录
func (d *AliyunDriver) BatchCreateRecords(ctx context.Context, zoneName string, records []*Record) (*BatchResult, error) {
	result := &BatchResult{
		Total:   len(records),
		Results: make([]*OperationResult, 0, len(records)),
	}

	start := time.Now()

	for i, record := range records {
		opResult := &OperationResult{
			ID: fmt.Sprintf("batch_%d", i),
		}

		createdRecord, err := d.CreateRecord(ctx, zoneName, record)
		if err != nil {
			opResult.Success = false
			opResult.ErrorMsg = err.Error()
			result.Failed++
		} else {
			opResult.Success = true
			opResult.Data = createdRecord
			result.Success++
		}

		result.Results = append(result.Results, opResult)
	}

	result.Duration = time.Since(start)

	d.logger.Info("Batch create records completed",
		zap.String("zone", zoneName),
		zap.Int("total", result.Total),
		zap.Int("success", result.Success),
		zap.Int("failed", result.Failed))

	return result, nil
}

// BatchUpdateRecords 批量更新DNS记录
func (d *AliyunDriver) BatchUpdateRecords(ctx context.Context, zoneName string, records []*Record) (*BatchResult, error) {
	result := &BatchResult{
		Total:   len(records),
		Results: make([]*OperationResult, 0, len(records)),
	}

	start := time.Now()

	for i, record := range records {
		opResult := &OperationResult{
			ID: record.ID,
		}

		updatedRecord, err := d.UpdateRecord(ctx, zoneName, record)
		if err != nil {
			opResult.Success = false
			opResult.ErrorMsg = err.Error()
			result.Failed++
		} else {
			opResult.Success = true
			opResult.Data = updatedRecord
			result.Success++
		}

		result.Results = append(result.Results, opResult)
	}

	result.Duration = time.Since(start)

	d.logger.Info("Batch update records completed",
		zap.String("zone", zoneName),
		zap.Int("total", result.Total),
		zap.Int("success", result.Success),
		zap.Int("failed", result.Failed))

	return result, nil
}

// BatchDeleteRecords 批量删除DNS记录
func (d *AliyunDriver) BatchDeleteRecords(ctx context.Context, zoneName string, recordIDs []string) (*BatchResult, error) {
	result := &BatchResult{
		Total:   len(recordIDs),
		Results: make([]*OperationResult, 0, len(recordIDs)),
	}

	start := time.Now()

	for _, recordID := range recordIDs {
		opResult := &OperationResult{
			ID: recordID,
		}

		err := d.DeleteRecord(ctx, zoneName, recordID)
		if err != nil {
			opResult.Success = false
			opResult.ErrorMsg = err.Error()
			result.Failed++
		} else {
			opResult.Success = true
			result.Success++
		}

		result.Results = append(result.Results, opResult)
	}

	result.Duration = time.Since(start)

	d.logger.Info("Batch delete records completed",
		zap.String("zone", zoneName),
		zap.Int("total", result.Total),
		zap.Int("success", result.Success),
		zap.Int("failed", result.Failed))

	return result, nil
}
