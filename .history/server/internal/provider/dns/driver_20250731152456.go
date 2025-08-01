package dns

import (
	"context"
	"errors"
	"time"
)

// 错误定义
var (
	ErrInvalidProviderType = errors.New("invalid provider type")
	ErrInvalidProviderName = errors.New("invalid provider name")
	ErrMissingCredentials  = errors.New("missing credentials")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrProviderNotFound    = errors.New("provider not found")
	ErrZoneNotFound        = errors.New("zone not found")
	ErrRecordNotFound      = errors.New("record not found")
	ErrRecordExists        = errors.New("record already exists")
	ErrInvalidRecord       = errors.New("invalid record")
	ErrQuotaExceeded       = errors.New("quota exceeded")
	ErrRateLimitExceeded   = errors.New("rate limit exceeded")
	ErrOperationTimeout    = errors.New("operation timeout")
	ErrUnsupportedOperation = errors.New("unsupported operation")
	ErrDNSSECNotSupported  = errors.New("DNSSEC not supported")
	ErrInvalidZoneFile     = errors.New("invalid zone file")
)

// Record DNS记录结构
type Record struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	TTL      int    `json:"ttl"`
	Priority *int   `json:"priority,omitempty"`
	Weight   *int   `json:"weight,omitempty"`
	Port     *int   `json:"port,omitempty"`
}

// Zone DNS区域结构
type Zone struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// ProviderInfo 提供商信息
type ProviderInfo struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Version     string            `json:"version"`
	Features    []string          `json:"features"`
	Limits      map[string]int    `json:"limits"`
	Regions     []string          `json:"regions"`
	RecordTypes []string          `json:"record_types"`
	Metadata    map[string]string `json:"metadata"`
}

// SyncResult 同步结果
type SyncResult struct {
	Success        bool              `json:"success"`
	TotalRecords   int               `json:"total_records"`
	AddedRecords   int               `json:"added_records"`
	UpdatedRecords int               `json:"updated_records"`
	DeletedRecords int               `json:"deleted_records"`
	FailedRecords  int               `json:"failed_records"`
	Duration       time.Duration     `json:"duration"`
	Errors         []string          `json:"errors"`
	Details        map[string]string `json:"details"`
}

// TestResult 测试结果
type TestResult struct {
	Success    bool              `json:"success"`
	Latency    time.Duration     `json:"latency"`
	ErrorMsg   string            `json:"error_msg"`
	Details    map[string]string `json:"details"`
	TestedAt   time.Time         `json:"tested_at"`
	TestType   string            `json:"test_type"`
	Endpoint   string            `json:"endpoint"`
	StatusCode int               `json:"status_code"`
}

// ValidationResult 验证结果
type ValidationResult struct {
	Valid       bool              `json:"valid"`
	ErrorMsg    string            `json:"error_msg"`
	Suggestions []string          `json:"suggestions"`
	Details     map[string]string `json:"details"`
}

// ListOptions 列表查询选项
type ListOptions struct {
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Filter   map[string]string `json:"filter"`
	Sort     string            `json:"sort"`
	Order    string            `json:"order"`
}

// BatchResult 批量操作结果
type BatchResult struct {
	Total    int                `json:"total"`
	Success  int                `json:"success"`
	Failed   int                `json:"failed"`
	Results  []*OperationResult `json:"results"`
	Duration time.Duration      `json:"duration"`
}

// OperationResult 单个操作结果
type OperationResult struct {
	ID       string      `json:"id"`
	Success  bool        `json:"success"`
	ErrorMsg string      `json:"error_msg"`
	Data     interface{} `json:"data"`
}

// ChallengeResult DNS-01验证结果
type ChallengeResult struct {
	Domain    string    `json:"domain"`
	Token     string    `json:"token"`
	RecordID  string    `json:"record_id"`
	TTL       int       `json:"ttl"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ChallengeValidation DNS-01验证状态
type ChallengeValidation struct {
	Valid         bool      `json:"valid"`
	Propagated    bool      `json:"propagated"`
	Value         string    `json:"value"`
	ExpectedValue string    `json:"expected_value"`
	CheckedAt     time.Time `json:"checked_at"`
	Servers       []string  `json:"servers"`
}

// SyncOptions 同步选项
type SyncOptions struct {
	DryRun       bool     `json:"dry_run"`
	Force        bool     `json:"force"`
	RecordTypes  []string `json:"record_types"`
	ExcludeNames []string `json:"exclude_names"`
}

// ZoneComparison 区域比较结果
type ZoneComparison struct {
	Domain        string    `json:"domain"`
	LocalRecords  int       `json:"local_records"`
	RemoteRecords int       `json:"remote_records"`
	ToAdd         []*Record `json:"to_add"`
	ToUpdate      []*Record `json:"to_update"`
	ToDelete      []*Record `json:"to_delete"`
	Conflicts     []*Record `json:"conflicts"`
	ComparedAt    time.Time `json:"compared_at"`
}

// DNSSECResult DNSSEC操作结果
type DNSSECResult struct {
	Enabled   bool         `json:"enabled"`
	Keys      []*DNSSECKey `json:"keys"`
	Status    string       `json:"status"`
	Message   string       `json:"message"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// DNSSECKey DNSSEC密钥
type DNSSECKey struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // KSK, ZSK
	Algorithm int       `json:"algorithm"`
	KeyTag    int       `json:"key_tag"`
	PublicKey string    `json:"public_key"`
	DS        string    `json:"ds"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ImportResult 导入结果
type ImportResult struct {
	Total    int                `json:"total"`
	Success  int                `json:"success"`
	Failed   int                `json:"failed"`
	Skipped  int                `json:"skipped"`
	Results  []*OperationResult `json:"results"`
	Duration time.Duration      `json:"duration"`
}

// Statistics 统计信息
type Statistics struct {
	Domain        string                 `json:"domain"`
	TotalRecords  int                    `json:"total_records"`
	RecordsByType map[string]int         `json:"records_by_type"`
	LastSync      time.Time              `json:"last_sync"`
	SyncCount     int                    `json:"sync_count"`
	ErrorCount    int                    `json:"error_count"`
	Metrics       map[string]interface{} `json:"metrics"`
}

// Quota 配额信息
type Quota struct {
	Provider      string          `json:"provider"`
	TotalDomains  int             `json:"total_domains"`
	UsedDomains   int             `json:"used_domains"`
	TotalRecords  int             `json:"total_records"`
	UsedRecords   int             `json:"used_records"`
	RateLimit     int             `json:"rate_limit"`
	RateRemaining int             `json:"rate_remaining"`
	ResetTime     time.Time       `json:"reset_time"`
	Features      map[string]bool `json:"features"`
	Limits        map[string]int  `json:"limits"`
}

// Driver DNS提供商驱动接口
type Driver interface {
	// 基础信息
	GetInfo() *ProviderInfo
	GetCapabilities() []string
	GetSupportedRecordTypes() []string

	// 连接与验证
	Test(ctx context.Context) *TestResult
	ValidateCredentials(ctx context.Context, credentials map[string]string) *ValidationResult

	// 区域管理
	ListZones(ctx context.Context, options *ListOptions) ([]Zone, error)
	GetZone(ctx context.Context, zoneName string) (*Zone, error)
	CreateZone(ctx context.Context, zoneName string) (*Zone, error)
	UpdateZone(ctx context.Context, zone *Zone) (*Zone, error)
	DeleteZone(ctx context.Context, zoneName string) error

	// DNS记录管理
	ListRecords(ctx context.Context, zoneName string, options *ListOptions) ([]Record, error)
	GetRecord(ctx context.Context, zoneName, recordID string) (*Record, error)
	CreateRecord(ctx context.Context, zoneName string, record *Record) (*Record, error)
	UpdateRecord(ctx context.Context, zoneName string, record *Record) (*Record, error)
	DeleteRecord(ctx context.Context, zoneName, recordID string) error

	// 批量操作
	BatchCreateRecords(ctx context.Context, zoneName string, records []*Record) (*BatchResult, error)
	BatchUpdateRecords(ctx context.Context, zoneName string, records []*Record) (*BatchResult, error)
	BatchDeleteRecords(ctx context.Context, zoneName string, recordIDs []string) (*BatchResult, error)

	// 区域同步与比较
	SyncZone(ctx context.Context, zoneName string, options *SyncOptions) (*SyncResult, error)
	CompareZone(ctx context.Context, zoneName string, localRecords []*Record) (*ZoneComparison, error)

	// 证书DNS-01验证支持
	CreateTXTChallenge(ctx context.Context, domain, token string, ttl int) (*ChallengeResult, error)
	DeleteTXTChallenge(ctx context.Context, domain, token string) error
	ValidateChallenge(ctx context.Context, domain, token string) (*ChallengeValidation, error)
	WaitForPropagation(ctx context.Context, domain, recordType, expectedValue string, timeout time.Duration) error

	// DNSSEC支持
	EnableDNSSEC(ctx context.Context, zoneName string) (*DNSSECResult, error)
	DisableDNSSEC(ctx context.Context, zoneName string) error
	GetDNSSECKeys(ctx context.Context, zoneName string) ([]*DNSSECKey, error)
	RotateDNSSECKeys(ctx context.Context, zoneName string) (*DNSSECResult, error)

	// 区域文件操作
	GetZoneFile(ctx context.Context, zoneName string) (string, error)
	ImportZoneFile(ctx context.Context, zoneName string, zoneFile string) (*ImportResult, error)
	ExportZoneFile(ctx context.Context, zoneName string, format string) (string, error)

	// 查询与过滤
	GetRecordsByType(ctx context.Context, zoneName, recordType string) ([]Record, error)
	GetRecordsByName(ctx context.Context, zoneName, recordName string) ([]Record, error)
	SearchRecords(ctx context.Context, zoneName string, query map[string]string) ([]Record, error)

	// 验证与检查
	ValidateRecord(record *Record) error
	ValidateZone(ctx context.Context, zoneName string) (*ValidationResult, error)
	CheckRecordConflicts(ctx context.Context, zoneName string, record *Record) ([]string, error)

	// 统计与配额
	GetStatistics(ctx context.Context, zoneName string) (*Statistics, error)
	GetQuota(ctx context.Context) (*Quota, error)
	GetUsage(ctx context.Context) (map[string]interface{}, error)

	// 监控与健康检查
	HealthCheck(ctx context.Context) *TestResult
	GetMetrics(ctx context.Context) (map[string]interface{}, error)

	// 高级功能
	SetRecordComment(ctx context.Context, zoneName, recordID, comment string) error
	GetRecordHistory(ctx context.Context, zoneName, recordID string) ([]interface{}, error)
	CreateRecordSet(ctx context.Context, zoneName string, records []*Record) ([]*Record, error)
	UpdateRecordSet(ctx context.Context, zoneName string, records []*Record) ([]*Record, error)
}

// Config 提供商配置接口
type Config interface {
	// GetType 获取提供商类型
	GetType() string

	// GetName 获取提供商名称
	GetName() string

	// GetCredentials 获取凭证信息
	GetCredentials() map[string]string

	// GetEndpoint 获取API端点
	GetEndpoint() string

	// GetRegion 获取区域
	GetRegion() string

	// GetTimeout 获取超时时间
	GetTimeout() time.Duration

	// GetRateLimit 获取速率限制
	GetRateLimit() int

	// GetConcurrent 获取并发数
	GetConcurrent() int

	// Validate 验证配置
	Validate() error
}

// BaseConfig 基础配置结构
type BaseConfig struct {
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	Credentials map[string]string `json:"credentials"`
	Endpoint    string            `json:"endpoint"`
	Region      string            `json:"region"`
	Timeout     time.Duration     `json:"timeout"`
	RateLimit   int               `json:"rate_limit"`
	Concurrent  int               `json:"concurrent"`
}

// GetType 获取提供商类型
func (c *BaseConfig) GetType() string {
	return c.Type
}

// GetName 获取提供商名称
func (c *BaseConfig) GetName() string {
	return c.Name
}

// GetCredentials 获取凭证信息
func (c *BaseConfig) GetCredentials() map[string]string {
	return c.Credentials
}

// GetEndpoint 获取API端点
func (c *BaseConfig) GetEndpoint() string {
	return c.Endpoint
}

// GetRegion 获取区域
func (c *BaseConfig) GetRegion() string {
	return c.Region
}

// GetTimeout 获取超时时间
func (c *BaseConfig) GetTimeout() time.Duration {
	if c.Timeout == 0 {
		return 30 * time.Second
	}
	return c.Timeout
}

// GetRateLimit 获取速率限制
func (c *BaseConfig) GetRateLimit() int {
	if c.RateLimit == 0 {
		return 10
	}
	return c.RateLimit
}

// GetConcurrent 获取并发数
func (c *BaseConfig) GetConcurrent() int {
	if c.Concurrent == 0 {
		return 5
	}
	return c.Concurrent
}

// Validate 验证配置
func (c *BaseConfig) Validate() error {
	if c.Type == "" {
		return ErrInvalidProviderType
	}
	if c.Name == "" {
		return ErrInvalidProviderName
	}
	if len(c.Credentials) == 0 {
		return ErrMissingCredentials
	}
	return nil
}
