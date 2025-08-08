package dns

import (
    "context"
    "time"
)

// StubDriver 为尚未完成的 Provider 提供一个占位实现，使编译通过且基础接口（Test/ListZones）可正常工作
// 其余方法全部返回 ErrNotImplemented。
// 这样在 Provider 同步流程中，未实现的云厂商至少可以返回空列表而不是 500 错误。

type StubDriver struct {
    info *ProviderInfo
}

// NewStubDriver 创建一个简单的 StubDriver
func NewStubDriver(providerType, providerName string) *StubDriver {
    return &StubDriver{
        info: &ProviderInfo{
            Name:    providerName,
            Type:    providerType,
            Version: "0.1.0",
            Features: []string{"dns_management"},
            RecordTypes: []string{"A", "AAAA", "CNAME", "TXT"},
            Metadata: map[string]string{},
        },
    }
}

// --- 基础信息 ---
func (d *StubDriver) GetInfo() *ProviderInfo                     { return d.info }
func (d *StubDriver) GetCapabilities() []string                  { return d.info.Features }
func (d *StubDriver) GetSupportedRecordTypes() []string          { return d.info.RecordTypes }

// --- 连接与验证 ---
func (d *StubDriver) Test(ctx context.Context) *TestResult {
    return &TestResult{Success: true, Latency: 0, TestedAt: time.Now(), TestType: "noop"}
}
func (d *StubDriver) ValidateCredentials(ctx context.Context, credentials map[string]string) *ValidationResult {
    return &ValidationResult{Valid: true}
}

// --- 区域管理 ---
func (d *StubDriver) ListZones(ctx context.Context, options *ListOptions) ([]Zone, error) {
    return []Zone{}, nil
}
func (d *StubDriver) GetZone(ctx context.Context, zoneName string) (*Zone, error)             { return nil, ErrNotImplemented }
func (d *StubDriver) CreateZone(ctx context.Context, zoneName string) (*Zone, error)          { return nil, ErrNotImplemented }
func (d *StubDriver) UpdateZone(ctx context.Context, zone *Zone) (*Zone, error)               { return nil, ErrNotImplemented }
func (d *StubDriver) DeleteZone(ctx context.Context, zoneName string) error                   { return ErrNotImplemented }

// --- DNS记录管理 ---
func (d *StubDriver) ListRecords(ctx context.Context, zoneName string, options *ListOptions) ([]Record, error) {
    return []Record{}, ErrNotImplemented
}
func (d *StubDriver) GetRecord(ctx context.Context, zoneName, recordID string) (*Record, error) {
    return nil, ErrNotImplemented
}
func (d *StubDriver) CreateRecord(ctx context.Context, zoneName string, record *Record) (*Record, error) {
    return nil, ErrNotImplemented
}
func (d *StubDriver) UpdateRecord(ctx context.Context, zoneName string, record *Record) (*Record, error) {
    return nil, ErrNotImplemented
}
func (d *StubDriver) DeleteRecord(ctx context.Context, zoneName, recordID string) error { return ErrNotImplemented }

// --- 批量操作 ---
func (d *StubDriver) BatchCreateRecords(ctx context.Context, zoneName string, records []*Record) (*BatchResult, error) {
    return nil, ErrNotImplemented
}
func (d *StubDriver) BatchUpdateRecords(ctx context.Context, zoneName string, records []*Record) (*BatchResult, error) {
    return nil, ErrNotImplemented
}
func (d *StubDriver) BatchDeleteRecords(ctx context.Context, zoneName string, recordIDs []string) (*BatchResult, error) {
    return nil, ErrNotImplemented
}

// --- 区域同步与比较 ---
func (d *StubDriver) SyncZone(ctx context.Context, zoneName string, options *SyncOptions) (*SyncResult, error) {
    return nil, ErrNotImplemented
}
func (d *StubDriver) CompareZone(ctx context.Context, zoneName string, localRecords []*Record) (*ZoneComparison, error) {
    return nil, ErrNotImplemented
}

// --- 证书 DNS-01 ---
func (d *StubDriver) CreateTXTChallenge(ctx context.Context, domain, token string, ttl int) (*ChallengeResult, error) {
    return nil, ErrNotImplemented
}
func (d *StubDriver) DeleteTXTChallenge(ctx context.Context, domain, token string) error { return ErrNotImplemented }
func (d *StubDriver) ValidateChallenge(ctx context.Context, domain, token string) (*ChallengeValidation, error) {
    return nil, ErrNotImplemented
}
func (d *StubDriver) WaitForPropagation(ctx context.Context, domain, recordType, expectedValue string, timeout time.Duration) error {
    return ErrNotImplemented
}

// --- DNSSEC ---
func (d *StubDriver) EnableDNSSEC(ctx context.Context, zoneName string) (*DNSSECResult, error) { return nil, ErrNotImplemented }
func (d *StubDriver) DisableDNSSEC(ctx context.Context, zoneName string) error                 { return ErrNotImplemented }
func (d *StubDriver) GetDNSSECKeys(ctx context.Context, zoneName string) ([]*DNSSECKey, error) { return nil, ErrNotImplemented }
func (d *StubDriver) RotateDNSSECKeys(ctx context.Context, zoneName string) (*DNSSECResult, error) {
    return nil, ErrNotImplemented
}

// --- 区域文件 ---
func (d *StubDriver) GetZoneFile(ctx context.Context, zoneName string) (string, error)             { return "", ErrNotImplemented }
func (d *StubDriver) ImportZoneFile(ctx context.Context, zoneName string, zoneFile string) (*ImportResult, error) {
    return nil, ErrNotImplemented
}
func (d *StubDriver) ExportZoneFile(ctx context.Context, zoneName string, format string) (string, error) {
    return "", ErrNotImplemented
}

// --- 查询与过滤 ---
func (d *StubDriver) GetRecordsByType(ctx context.Context, zoneName, recordType string) ([]Record, error) {
    return []Record{}, ErrNotImplemented
}
func (d *StubDriver) GetRecordsByName(ctx context.Context, zoneName, recordName string) ([]Record, error) {
    return []Record{}, ErrNotImplemented
}
func (d *StubDriver) SearchRecords(ctx context.Context, zoneName string, query map[string]string) ([]Record, error) {
    return []Record{}, ErrNotImplemented
}

// --- 验证与检查 ---
func (d *StubDriver) ValidateRecord(record *Record) error                                                 { return ErrNotImplemented }
func (d *StubDriver) ValidateZone(ctx context.Context, zoneName string) (*ValidationResult, error)        { return nil, ErrNotImplemented }
func (d *StubDriver) CheckRecordConflicts(ctx context.Context, zoneName string, record *Record) ([]string, error) {
    return []string{}, ErrNotImplemented
}

// --- 统计与配额 ---
func (d *StubDriver) GetStatistics(ctx context.Context, zoneName string) (*Statistics, error) { return nil, ErrNotImplemented }
func (d *StubDriver) GetQuota(ctx context.Context) (*Quota, error)                            { return nil, ErrNotImplemented }
func (d *StubDriver) GetUsage(ctx context.Context) (map[string]interface{}, error)            { return nil, ErrNotImplemented }

// --- 监控与健康检查 ---
func (d *StubDriver) HealthCheck(ctx context.Context) *TestResult                    { return d.Test(ctx) }
func (d *StubDriver) GetMetrics(ctx context.Context) (map[string]interface{}, error) { return nil, ErrNotImplemented }

// --- 高级功能 ---
func (d *StubDriver) SetRecordComment(ctx context.Context, zoneName, recordID, comment string) error { return ErrNotImplemented }
func (d *StubDriver) GetRecordHistory(ctx context.Context, zoneName, recordID string) ([]interface{}, error) {
    return nil, ErrNotImplemented
}
func (d *StubDriver) CreateRecordSet(ctx context.Context, zoneName string, records []*Record) ([]*Record, error) {
    return nil, ErrNotImplemented
}
func (d *StubDriver) UpdateRecordSet(ctx context.Context, zoneName string, records []*Record) ([]*Record, error) {
    return nil, ErrNotImplemented
} 
