package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GoDaddyConfig GoDaddy DNS配置
type GoDaddyConfig struct {
	*BaseConfig
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

// GoDaddyDriver GoDaddy DNS驱动
type GoDaddyDriver struct {
	config     *GoDaddyConfig
	httpClient *http.Client
	baseURL    string
}

// GoDaddy API响应结构
type GoDaddyDomain struct {
	Domain      string   `json:"domain"`
	Status      string   `json:"status"`
	Expires     string   `json:"expires"`
	CreatedAt   string   `json:"createdAt"`
	DomainID    int64    `json:"domainId"`
	Renewable   bool     `json:"renewable"`
	Privacy     bool     `json:"privacy"`
	Locked      bool     `json:"locked"`
	NameServers []string `json:"nameServers"`
}

type GoDaddyRecord struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Data     string `json:"data"`
	Priority int    `json:"priority,omitempty"`
	Port     int    `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Service  string `json:"service,omitempty"`
	TTL      int    `json:"ttl"`
	Weight   int    `json:"weight,omitempty"`
}

type GoDaddyError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Fields  []struct {
		Code        string `json:"code"`
		Message     string `json:"message"`
		Path        string `json:"path"`
		PathRelated string `json:"pathRelated,omitempty"`
	} `json:"fields,omitempty"`
}

// Validate 验证配置
func (c *GoDaddyConfig) Validate() error {
	if err := c.BaseConfig.Validate(); err != nil {
		return err
	}

	if c.APIKey == "" {
		return fmt.Errorf("API Key is required")
	}

	if c.APISecret == "" {
		return fmt.Errorf("API Secret is required")
	}

	return nil
}

// NewGoDaddyDriver 创建GoDaddy DNS驱动
func NewGoDaddyDriver(config Config) (Driver, error) {
	goDaddyConfig, ok := config.(*GoDaddyConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type for GoDaddy driver")
	}

	if err := goDaddyConfig.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	driver := &GoDaddyDriver{
		config:  goDaddyConfig,
		baseURL: "https://api.godaddy.com/v1",
		httpClient: &http.Client{
			Timeout: goDaddyConfig.GetTimeout(),
		},
	}

	return driver, nil
}

// Test 测试连接
func (d *GoDaddyDriver) Test(ctx context.Context) *TestResult {
	result := &TestResult{
		TestedAt: time.Now(),
		TestType: "connection",
	}

	// 测试API连接 - 获取域名列表
	req, err := http.NewRequestWithContext(ctx, "GET", d.baseURL+"/domains", nil)
	if err != nil {
		result.Success = false
		result.ErrorMsg = fmt.Sprintf("Failed to create request: %v", err)
		return result
	}

	// 设置认证头
	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", d.config.APIKey, d.config.APISecret))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		result.Success = false
		result.ErrorMsg = fmt.Sprintf("Failed to make request: %v", err)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		result.Success = false
		result.ErrorMsg = "Authentication failed: Invalid API Key or Secret"
		return result
	}

	if resp.StatusCode == 403 {
		result.Success = false
		result.ErrorMsg = "Access denied: Check API Key permissions"
		return result
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		result.Success = false
		result.ErrorMsg = fmt.Sprintf("API request failed with status %d: %s", resp.StatusCode, string(body))
		return result
	}

	result.Success = true
	result.Details = map[string]string{
		"message": "GoDaddy DNS connection successful",
	}
	return result
}

// ListZones 获取域名列表
func (d *GoDaddyDriver) ListZones(ctx context.Context) ([]*Zone, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", d.baseURL+"/domains", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置认证头
	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", d.config.APIKey, d.config.APISecret))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var domains []GoDaddyDomain
	if err := json.NewDecoder(resp.Body).Decode(&domains); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	zones := make([]*Zone, 0, len(domains))
	for _, domain := range domains {
		zone := &Zone{
			ID:   strconv.FormatInt(domain.DomainID, 10),
			Name: domain.Domain,
		}

		// 解析过期时间
		if domain.Expires != "" {
			if expireTime, err := time.Parse("2006-01-02T15:04:05.000Z", domain.Expires); err == nil {
				zone.ExpiresAt = &expireTime
			}
		}

		// 设置状态
		switch strings.ToLower(domain.Status) {
		case "active":
			zone.Status = "active"
		case "expired":
			zone.Status = "expired"
		case "pending":
			zone.Status = "pending"
		default:
			zone.Status = "unknown"
		}

		zones = append(zones, zone)
	}

	return zones, nil
}

// GetZone 获取单个域名信息
func (d *GoDaddyDriver) GetZone(ctx context.Context, zoneID string) (*Zone, error) {
	// GoDaddy API使用域名而不是ID
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/domains/%s", d.baseURL, zoneID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", d.config.APIKey, d.config.APISecret))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var domain GoDaddyDomain
	if err := json.NewDecoder(resp.Body).Decode(&domain); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	zone := &Zone{
		ID:   strconv.FormatInt(domain.DomainID, 10),
		Name: domain.Domain,
	}

	// 解析过期时间
	if domain.Expires != "" {
		if expireTime, err := time.Parse("2006-01-02T15:04:05.000Z", domain.Expires); err == nil {
			zone.ExpiresAt = &expireTime
		}
	}

	return zone, nil
}

// ListRecords 获取DNS记录列表
func (d *GoDaddyDriver) ListRecords(ctx context.Context, zoneID string) ([]*Record, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/domains/%s/records", d.baseURL, zoneID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", d.config.APIKey, d.config.APISecret))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var goDaddyRecords []GoDaddyRecord
	if err := json.NewDecoder(resp.Body).Decode(&goDaddyRecords); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	records := make([]*Record, 0, len(goDaddyRecords))
	for i, gdRecord := range goDaddyRecords {
		record := &Record{
			ID:       fmt.Sprintf("%s-%d", zoneID, i), // GoDaddy没有记录ID，生成一个
			ZoneID:   zoneID,
			Name:     gdRecord.Name,
			Type:     gdRecord.Type,
			Value:    gdRecord.Data,
			TTL:      gdRecord.TTL,
			Priority: gdRecord.Priority,
		}
		records = append(records, record)
	}

	return records, nil
}

// CreateRecord 创建DNS记录
func (d *GoDaddyDriver) CreateRecord(ctx context.Context, zoneID string, record *Record) (*Record, error) {
	goDaddyRecord := GoDaddyRecord{
		Type:     record.Type,
		Name:     record.Name,
		Data:     record.Value,
		TTL:      record.TTL,
		Priority: record.Priority,
	}

	recordData, err := json.Marshal([]GoDaddyRecord{goDaddyRecord})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal record: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/records", d.baseURL, zoneID)
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(recordData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", d.config.APIKey, d.config.APISecret))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// GoDaddy PATCH返回成功但不返回记录，我们返回原记录
	record.ID = fmt.Sprintf("%s-%s-%s", zoneID, record.Type, record.Name)
	return record, nil
}

// UpdateRecord 更新DNS记录
func (d *GoDaddyDriver) UpdateRecord(ctx context.Context, zoneID string, record *Record) (*Record, error) {
	// GoDaddy使用PATCH方法更新记录，与创建类似
	return d.CreateRecord(ctx, zoneID, record)
}

// DeleteRecord 删除DNS记录
func (d *GoDaddyDriver) DeleteRecord(ctx context.Context, zoneID, recordID string) error {
	// GoDaddy删除记录需要指定类型和名称
	parts := strings.Split(recordID, "-")
	if len(parts) < 3 {
		return fmt.Errorf("invalid record ID format")
	}

	recordType := parts[1]
	recordName := parts[2]

	url := fmt.Sprintf("%s/domains/%s/records/%s/%s", d.baseURL, zoneID, recordType, recordName)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", d.config.APIKey, d.config.APISecret))
	req.Header.Set("Accept", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetInfo 获取驱动信息
func (d *GoDaddyDriver) GetInfo() *DriverInfo {
	return &DriverInfo{
		Name:        "GoDaddy",
		Type:        "godaddy",
		Version:     "1.0.0",
		Description: "GoDaddy DNS Provider",
		Author:      "DevOps Asset System",
		SupportedRecordTypes: []string{
			"A", "AAAA", "CNAME", "MX", "TXT", "NS", "SRV", "SOA",
		},
		Features: []string{
			"zone_management",
			"record_management",
			"bulk_operations",
		},
	}
}
