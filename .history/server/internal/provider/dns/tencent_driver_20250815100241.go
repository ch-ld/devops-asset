package dns

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"go.uber.org/zap"
)

// TencentDriver 腾讯云DNS驱动
type TencentDriver struct {
	client *dnspod.Client
	config *TencentConfig
	logger *zap.Logger
	info   *ProviderInfo
}

// TencentConfig 腾讯云配置
type TencentConfig struct {
	*BaseConfig
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

// NewTencentDriver 创建腾讯云DNS驱动
func NewTencentDriver(config Config) (Driver, error) {
	tencentConfig, ok := config.(*TencentConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type for Tencent driver")
	}

	if err := tencentConfig.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// 创建认证信息
	credential := common.NewCredential(
		tencentConfig.SecretID,
		tencentConfig.SecretKey,
	)

	// 创建客户端配置
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	cpf.HttpProfile.ReqTimeout = int(tencentConfig.GetTimeout().Seconds())

	// 创建客户端
	client, err := dnspod.NewClient(credential, tencentConfig.GetRegion(), cpf)
	if err != nil {
		return nil, fmt.Errorf("failed to create Tencent DNS client: %w", err)
	}

	driver := &TencentDriver{
		client: client,
		config: tencentConfig,
		logger: zap.L().Named("tencent-dns"),
		info: &ProviderInfo{
			Name:    "Tencent DNS",
			Type:    "tencent",
			Version: "1.0.0",
			Features: []string{
				"dns_management",
				"batch_operations",
				"txt_challenges",
				"zone_sync",
				"statistics",
			},
			Limits: map[string]int{
				"max_domains":            500,
				"max_records_per_domain": 50000,
				"max_batch_size":         100,
				"rate_limit_per_second":  20,
			},
			Regions: []string{
				"ap-beijing",
				"ap-shanghai",
				"ap-guangzhou",
				"ap-chengdu",
			},
			RecordTypes: []string{
				"A", "AAAA", "CNAME", "MX", "TXT", "NS", "SRV", "CAA", "PTR",
			},
			Metadata: map[string]string{
				"api_version": "2021-03-23",
				"endpoint":    "dnspod.tencentcloudapi.com",
			},
		},
	}

	return driver, nil
}

// GetInfo 获取提供商信息
func (d *TencentDriver) GetInfo() *ProviderInfo {
	return d.info
}

// GetCapabilities 获取功能列表
func (d *TencentDriver) GetCapabilities() []string {
	return d.info.Features
}

// GetSupportedRecordTypes 获取支持的记录类型
func (d *TencentDriver) GetSupportedRecordTypes() []string {
	return d.info.RecordTypes
}

// Test 测试连接
func (d *TencentDriver) Test(ctx context.Context) *TestResult {
	start := time.Now()

	request := dnspod.NewDescribeDomainListRequest()
	request.Limit = common.Int64Ptr(1)

	_, err := d.client.DescribeDomainList(request)
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
func (d *TencentDriver) ValidateCredentials(ctx context.Context, credentials map[string]string) *ValidationResult {
	result := &ValidationResult{
		Details: make(map[string]string),
	}

	secretID, ok := credentials["secret_id"]
	if !ok || secretID == "" {
		result.Valid = false
		result.ErrorMsg = "secret_id is required"
		result.Suggestions = append(result.Suggestions, "Please provide a valid secret_id")
		return result
	}

	secretKey, ok := credentials["secret_key"]
	if !ok || secretKey == "" {
		result.Valid = false
		result.ErrorMsg = "secret_key is required"
		result.Suggestions = append(result.Suggestions, "Please provide a valid secret_key")
		return result
	}

	// 尝试创建客户端并测试
	credential := common.NewCredential(secretID, secretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"

	testClient, err := dnspod.NewClient(credential, "ap-beijing", cpf)
	if err != nil {
		result.Valid = false
		result.ErrorMsg = fmt.Sprintf("Failed to create client: %v", err)
		result.Suggestions = append(result.Suggestions, "Check your secret ID and key")
		return result
	}

	// 测试API调用
	request := dnspod.NewDescribeDomainListRequest()
	request.Limit = common.Int64Ptr(1)

	_, err = testClient.DescribeDomainList(request)
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
func (d *TencentDriver) ListZones(ctx context.Context, options *ListOptions) ([]Zone, error) {
	request := dnspod.NewDescribeDomainListRequest()

	if options != nil {
		if options.PageSize > 0 {
			request.Limit = common.Int64Ptr(int64(options.PageSize))
		} else {
			request.Limit = common.Int64Ptr(100)
		}

		if options.Page > 0 {
			offset := int64((options.Page - 1) * options.PageSize)
			request.Offset = common.Int64Ptr(offset)
		}

		// 处理过滤条件
		if keyword, ok := options.Filter["keyword"]; ok {
			request.Keyword = common.StringPtr(keyword)
		}
	} else {
		request.Limit = common.Int64Ptr(100)
	}

	response, err := d.client.DescribeDomainList(request)
	if err != nil {
		d.logger.Error("Failed to list zones", zap.Error(err))
		return nil, fmt.Errorf("failed to list zones: %w", err)
	}

	var zones []Zone
	for _, domain := range response.Response.DomainList {
		status := "active"
		if *domain.Status != "ENABLE" {
			status = "inactive"
		}

		zones = append(zones, Zone{
			ID:     strconv.FormatUint(*domain.DomainId, 10),
			Name:   *domain.Name,
			Status: status,
		})
	}

	d.logger.Info("Listed zones successfully", zap.Int("count", len(zones)))
	return zones, nil
}

// GetZone 获取指定DNS区域
func (d *TencentDriver) GetZone(ctx context.Context, zoneName string) (*Zone, error) {
	request := dnspod.NewDescribeDomainRequest()
	request.Domain = common.StringPtr(zoneName)

	response, err := d.client.DescribeDomain(request)
	if err != nil {
		d.logger.Error("Failed to get zone", zap.String("zone", zoneName), zap.Error(err))
		return nil, fmt.Errorf("failed to get zone %s: %w", zoneName, err)
	}

	domain := response.Response.DomainInfo
	status := "active"
	if *domain.Status != "ENABLE" {
		status = "inactive"
	}

	zone := &Zone{
		ID:     strconv.FormatUint(*domain.DomainId, 10),
		Name:   *domain.Domain,
		Status: status,
	}

	d.logger.Info("Got zone successfully", zap.String("zone", zoneName))
	return zone, nil
}

// CreateZone 创建DNS区域
func (d *TencentDriver) CreateZone(ctx context.Context, zoneName string) (*Zone, error) {
	request := dnspod.NewCreateDomainRequest()
	request.Domain = common.StringPtr(zoneName)

	response, err := d.client.CreateDomain(request)
	if err != nil {
		d.logger.Error("Failed to create zone", zap.String("zone", zoneName), zap.Error(err))
		return nil, fmt.Errorf("failed to create zone %s: %w", zoneName, err)
	}

	zone := &Zone{
		ID:     strconv.FormatUint(*response.Response.DomainInfo.Id, 10),
		Name:   zoneName,
		Status: "active",
	}

	d.logger.Info("Created zone successfully", zap.String("zone", zoneName))
	return zone, nil
}

// UpdateZone 更新DNS区域
func (d *TencentDriver) UpdateZone(ctx context.Context, zone *Zone) (*Zone, error) {
	// 腾讯云DNS支持有限的域名更新操作
	return zone, nil
}

// DeleteZone 删除DNS区域
func (d *TencentDriver) DeleteZone(ctx context.Context, zoneName string) error {
	request := dnspod.NewDeleteDomainRequest()
	request.Domain = common.StringPtr(zoneName)

	_, err := d.client.DeleteDomain(request)
	if err != nil {
		d.logger.Error("Failed to delete zone", zap.String("zone", zoneName), zap.Error(err))
		return fmt.Errorf("failed to delete zone %s: %w", zoneName, err)
	}

	d.logger.Info("Deleted zone successfully", zap.String("zone", zoneName))
	return nil
}

// ListRecords 列出指定区域的所有记录
func (d *TencentDriver) ListRecords(ctx context.Context, zoneName string, options *ListOptions) ([]Record, error) {
	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = common.StringPtr(zoneName)

	if options != nil {
		if options.PageSize > 0 {
			request.Limit = common.Uint64Ptr(uint64(options.PageSize))
		} else {
			request.Limit = common.Uint64Ptr(3000)
		}

		if options.Page > 0 {
			offset := uint64((options.Page - 1) * options.PageSize)
			request.Offset = common.Uint64Ptr(offset)
		}

		// 处理过滤条件
		if recordType, ok := options.Filter["type"]; ok {
			request.RecordType = common.StringPtr(recordType)
		}
		if keyword, ok := options.Filter["keyword"]; ok {
			request.Keyword = common.StringPtr(keyword)
		}
	} else {
		request.Limit = common.Uint64Ptr(3000)
	}

	response, err := d.client.DescribeRecordList(request)
	if err != nil {
		d.logger.Error("Failed to list records", zap.String("zone", zoneName), zap.Error(err))
		return nil, fmt.Errorf("failed to list records for zone %s: %w", zoneName, err)
	}

	var records []Record
	for _, record := range response.Response.RecordList {
		r := Record{
			ID:    strconv.FormatUint(*record.RecordId, 10),
			Name:  *record.Name,
			Type:  *record.Type,
			Value: *record.Value,
			TTL:   int(*record.TTL),
		}

		// 处理优先级
		if record.MX != nil {
			priority := int(*record.MX)
			r.Priority = &priority
		}

		// 处理权重
		if record.Weight != nil {
			weight := int(*record.Weight)
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
func (d *TencentDriver) GetRecord(ctx context.Context, zoneName, recordID string) (*Record, error) {
	recordIDUint, err := strconv.ParseUint(recordID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid record ID: %w", err)
	}

	request := dnspod.NewDescribeRecordRequest()
	request.Domain = common.StringPtr(zoneName)
	request.RecordId = common.Uint64Ptr(recordIDUint)

	response, err := d.client.DescribeRecord(request)
	if err != nil {
		d.logger.Error("Failed to get record",
			zap.String("zone", zoneName),
			zap.String("record_id", recordID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get record %s: %w", recordID, err)
	}

	recordInfo := response.Response.RecordInfo
	record := &Record{
		ID:    recordID,
		Name:  *recordInfo.SubDomain,
		Type:  *recordInfo.RecordType,
		Value: *recordInfo.Value,
		TTL:   int(*recordInfo.TTL),
	}

	// 处理优先级
	if recordInfo.MX != nil {
		priority := int(*recordInfo.MX)
		record.Priority = &priority
	}

	d.logger.Info("Got record successfully",
		zap.String("zone", zoneName),
		zap.String("record_id", recordID))

	return record, nil
}

// CreateRecord 创建DNS记录
func (d *TencentDriver) CreateRecord(ctx context.Context, zoneName string, record *Record) (*Record, error) {
	request := dnspod.NewCreateRecordRequest()
	request.Domain = common.StringPtr(zoneName)
	request.SubDomain = common.StringPtr(record.Name)
	request.RecordType = common.StringPtr(record.Type)
	request.Value = common.StringPtr(record.Value)

	if record.TTL > 0 {
		request.TTL = common.Uint64Ptr(uint64(record.TTL))
	}

	if record.Priority != nil {
		request.MX = common.Uint64Ptr(uint64(*record.Priority))
	}

	response, err := d.client.CreateRecord(request)
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
		ID:       strconv.FormatUint(*response.Response.RecordId, 10),
		Name:     record.Name,
		Type:     record.Type,
		Value:    record.Value,
		TTL:      record.TTL,
		Priority: record.Priority,
	}

	d.logger.Info("Created record successfully",
		zap.String("zone", zoneName),
		zap.String("record_id", createdRecord.ID),
		zap.String("name", record.Name),
		zap.String("type", record.Type))

	return createdRecord, nil
}

// UpdateRecord 更新DNS记录
func (d *TencentDriver) UpdateRecord(ctx context.Context, zoneName string, record *Record) (*Record, error) {
	recordIDUint, err := strconv.ParseUint(record.ID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid record ID: %w", err)
	}

	request := dnspod.NewModifyRecordRequest()
	request.Domain = common.StringPtr(zoneName)
	request.RecordId = common.Uint64Ptr(recordIDUint)
	request.SubDomain = common.StringPtr(record.Name)
	request.RecordType = common.StringPtr(record.Type)
	request.Value = common.StringPtr(record.Value)

	if record.TTL > 0 {
		request.TTL = common.Uint64Ptr(uint64(record.TTL))
	}

	if record.Priority != nil {
		request.MX = common.Uint64Ptr(uint64(*record.Priority))
	}

	_, err = d.client.ModifyRecord(request)
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
func (d *TencentDriver) DeleteRecord(ctx context.Context, zoneName, recordID string) error {
	recordIDUint, err := strconv.ParseUint(recordID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid record ID: %w", err)
	}

	request := dnspod.NewDeleteRecordRequest()
	request.Domain = common.StringPtr(zoneName)
	request.RecordId = common.Uint64Ptr(recordIDUint)

	_, err = d.client.DeleteRecord(request)
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
