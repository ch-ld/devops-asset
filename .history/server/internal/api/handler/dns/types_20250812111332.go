package dns

import (
	"time"
)

// DomainGroupResponse 域名分组响应结构体
// @Description 域名分组信息响应
type DomainGroupResponse struct {
	ID          uint                   `json:"id" example:"1" comment:"分组ID"`
	Name        string                 `json:"name" example:"生产环境" comment:"分组名称"`
	ParentID    *uint                  `json:"parent_id" example:"1" comment:"父分组ID"`
	Parent      *DomainGroupResponse   `json:"parent,omitempty" comment:"父分组信息"`
	Children    []*DomainGroupResponse `json:"children,omitempty" comment:"子分组列表"`
	Description string                 `json:"description" example:"生产环境域名" comment:"描述"`
	Sort        int                    `json:"sort" example:"1" comment:"排序"`
	TenantID    uint                   `json:"tenant_id" example:"1" comment:"租户ID"`
	CreatedBy   uint                   `json:"created_by" example:"1" comment:"创建人ID"`
	UpdatedBy   uint                   `json:"updated_by" example:"1" comment:"更新人ID"`
	CreatedAt   time.Time              `json:"created_at" example:"2024-01-01T00:00:00Z" comment:"创建时间"`
	UpdatedAt   time.Time              `json:"updated_at" example:"2024-01-01T00:00:00Z" comment:"更新时间"`
}

// TagResponse 标签响应结构体
// @Description 标签信息响应
type TagResponse struct {
	ID          uint      `json:"id" example:"1" comment:"标签ID"`
	Name        string    `json:"name" example:"重要" comment:"标签名称"`
	Color       string    `json:"color" example:"#1890ff" comment:"标签颜色"`
	Description string    `json:"description" example:"重要域名" comment:"描述"`
	TenantID    uint      `json:"tenant_id" example:"1" comment:"租户ID"`
	CreatedBy   uint      `json:"created_by" example:"1" comment:"创建人ID"`
	UpdatedBy   uint      `json:"updated_by" example:"1" comment:"更新人ID"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z" comment:"创建时间"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z" comment:"更新时间"`
}

// ProviderResponse DNS提供商响应结构体
// @Description DNS提供商信息响应
type ProviderResponse struct {
	ID          uint                   `json:"id" example:"1" comment:"提供商ID"`
	Name        string                 `json:"name" example:"阿里云DNS" comment:"提供商名称"`
	Type        string                 `json:"type" example:"aliyun" comment:"提供商类型"`
	Status      string                 `json:"status" example:"active" comment:"状态"`
	IsDefault   bool                   `json:"is_default" example:"false" comment:"是否默认提供商"`
	Priority    int                    `json:"priority" example:"0" comment:"优先级"`
	RateLimit   int                    `json:"rate_limit" example:"10" comment:"速率限制"`
	Concurrent  int                    `json:"concurrent" example:"5" comment:"并发数"`
	Timeout     int                    `json:"timeout" example:"30" comment:"超时时间"`
	LastTestAt  *time.Time             `json:"last_test_at" example:"2024-01-01T00:00:00Z" comment:"最后测试时间"`
	TestResult  string                 `json:"test_result" example:"success" comment:"测试结果"`
	ErrorMessage string                `json:"error_message" example:"" comment:"错误信息"`
	Remark      string                 `json:"remark" example:"主要DNS提供商" comment:"备注"`
	DomainCount int                    `json:"domain_count" example:"5" comment:"关联域名数量"`
	LastTestResult *TestConnectionResult `json:"last_test_result,omitempty" comment:"最后测试结果"`
	TenantID    uint                   `json:"tenant_id" example:"1" comment:"租户ID"`
	CreatedBy   uint                   `json:"created_by" example:"1" comment:"创建人ID"`
	UpdatedBy   uint                   `json:"updated_by" example:"1" comment:"更新人ID"`
	CreatedAt   time.Time              `json:"created_at" example:"2024-01-01T00:00:00Z" comment:"创建时间"`
	UpdatedAt   time.Time              `json:"updated_at" example:"2024-01-01T00:00:00Z" comment:"更新时间"`
}

// TestConnectionResult 测试连接结果
type TestConnectionResult struct {
	Success bool   `json:"success" comment:"是否成功"`
	Message string `json:"message" comment:"消息"`
	Latency string `json:"latency" comment:"延迟"`
}

// CertificateResponse 证书响应结构体
// @Description 证书信息响应
type CertificateResponse struct {
	ID              uint            `json:"id" example:"1" comment:"证书ID"`
	DomainID        uint            `json:"domain_id" example:"1" comment:"域名ID"`
	Domain          *DomainResponse `json:"domain,omitempty" comment:"域名信息"`
	CommonName      string          `json:"common_name" example:"example.com" comment:"主域名"`
	SubjectAltNames []string        `json:"subject_alt_names" example:"[\"www.example.com\",\"api.example.com\"]" comment:"SAN域名列表"`
	CAType          string          `json:"ca_type" example:"letsencrypt" comment:"CA类型"`
	Status          string          `json:"status" example:"issued" comment:"状态"`
	SerialNumber    string          `json:"serial_number" example:"03:A3:B2:C1:D4:E5:F6:07:08:09:0A:0B:0C:0D:0E:0F" comment:"证书序列号"`
	Fingerprint     string          `json:"fingerprint" example:"SHA256:1234567890ABCDEF..." comment:"证书指纹"`
	IssuedAt        *time.Time      `json:"issued_at" example:"2024-01-01T00:00:00Z" comment:"签发时间"`
	ExpiresAt       *time.Time      `json:"expires_at" example:"2024-04-01T00:00:00Z" comment:"过期时间"`
	AutoRenew       bool            `json:"auto_renew" example:"true" comment:"是否自动续期"`
	RenewDays       int             `json:"renew_days" example:"30" comment:"提前续期天数"`
	LastRenewAt     *time.Time      `json:"last_renew_at" example:"2024-01-01T00:00:00Z" comment:"最后续期时间"`
	Remark          string          `json:"remark" example:"主站证书" comment:"备注"`
	TenantID        uint            `json:"tenant_id" example:"1" comment:"租户ID"`
	CreatedBy       uint            `json:"created_by" example:"1" comment:"创建人ID"`
	UpdatedBy       uint            `json:"updated_by" example:"1" comment:"更新人ID"`
	CreatedAt       time.Time       `json:"created_at" example:"2024-01-01T00:00:00Z" comment:"创建时间"`
	UpdatedAt       time.Time       `json:"updated_at" example:"2024-01-01T00:00:00Z" comment:"更新时间"`
}

// StatisticsResponse 统计信息响应结构体
// @Description 统计信息响应
type StatisticsResponse struct {
	Total      int64                  `json:"total" example:"100" comment:"总数"`
	ByStatus   map[string]int64       `json:"by_status" example:"{\"active\":80,\"inactive\":20}" comment:"按状态统计"`
	ByType     map[string]int64       `json:"by_type,omitempty" example:"{\"A\":50,\"CNAME\":30,\"TXT\":20}" comment:"按类型统计"`
	Expiring   int                    `json:"expiring,omitempty" example:"5" comment:"即将过期数量"`
	Additional map[string]interface{} `json:"additional,omitempty" comment:"其他统计信息"`
}

// BatchOperationRequest 批量操作请求结构体
// @Description 批量操作请求参数
type BatchOperationRequest struct {
	IDs    []uint `json:"ids" binding:"required" example:"[1,2,3]" comment:"ID列表"`
	Action string `json:"action" binding:"required" example:"delete" comment:"操作类型"`
	Data   map[string]interface{} `json:"data,omitempty" comment:"操作数据"`
}

// BatchOperationResponse 批量操作响应结构体
// @Description 批量操作响应
type BatchOperationResponse struct {
	Success     int      `json:"success" example:"2" comment:"成功数量"`
	Failed      int      `json:"failed" example:"1" comment:"失败数量"`
	Total       int      `json:"total" example:"3" comment:"总数量"`
	FailedItems []string `json:"failed_items,omitempty" example:"[\"域名不存在\"]" comment:"失败项目"`
}

// ImportRequest 导入请求结构体
// @Description 导入请求参数
type ImportRequest struct {
	Data   []map[string]interface{} `json:"data" binding:"required" comment:"导入数据"`
	Format string                   `json:"format" example:"json" comment:"数据格式"`
	Options map[string]interface{}  `json:"options,omitempty" comment:"导入选项"`
}

// ImportResponse 导入响应结构体
// @Description 导入响应
type ImportResponse struct {
	Success     int      `json:"success" example:"8" comment:"成功数量"`
	Failed      int      `json:"failed" example:"2" comment:"失败数量"`
	Total       int      `json:"total" example:"10" comment:"总数量"`
	FailedItems []string `json:"failed_items,omitempty" example:"[\"域名格式错误\",\"域名已存在\"]" comment:"失败项目"`
}

// ExportRequest 导出请求结构体
// @Description 导出请求参数
type ExportRequest struct {
	Format  string                 `json:"format" example:"csv" comment:"导出格式(csv/excel/json)"`
	Filters map[string]interface{} `json:"filters,omitempty" comment:"筛选条件"`
	Fields  []string               `json:"fields,omitempty" example:"[\"name\",\"status\",\"expires_at\"]" comment:"导出字段"`
}

// SyncRequest 同步请求结构体
// @Description 同步请求参数
type SyncRequest struct {
	ProviderID *uint  `json:"provider_id" example:"1" comment:"提供商ID"`
	DomainID   *uint  `json:"domain_id" example:"1" comment:"域名ID"`
	SyncType   string `json:"sync_type" example:"full" comment:"同步类型(full/incremental)"`
	Force      bool   `json:"force" example:"false" comment:"是否强制同步"`
}

// SyncResponse 同步响应结构体
// @Description 同步响应
type SyncResponse struct {
	Success        bool     `json:"success" example:"true" comment:"是否成功"`
	TotalRecords   int      `json:"total_records" example:"100" comment:"总记录数"`
	AddedRecords   int      `json:"added_records" example:"5" comment:"新增记录数"`
	UpdatedRecords int      `json:"updated_records" example:"3" comment:"更新记录数"`
	DeletedRecords int      `json:"deleted_records" example:"2" comment:"删除记录数"`
	FailedRecords  int      `json:"failed_records" example:"1" comment:"失败记录数"`
	Duration       string   `json:"duration" example:"2.5s" comment:"耗时"`
	Errors         []string `json:"errors,omitempty" example:"[\"记录冲突\"]" comment:"错误信息"`
}

// TestConnectionRequest 测试连接请求结构体
// @Description 测试连接请求参数
type TestConnectionRequest struct {
	Type        string            `json:"type" binding:"required" example:"aliyun" comment:"提供商类型"`
	Credentials map[string]string `json:"credentials" binding:"required" comment:"凭证信息"`
	Endpoint    string            `json:"endpoint,omitempty" example:"https://alidns.aliyuncs.com" comment:"API端点"`
	Region      string            `json:"region,omitempty" example:"cn-hangzhou" comment:"区域"`
}

// TestConnectionResponse 测试连接响应结构体
// @Description 测试连接响应
type TestConnectionResponse struct {
	Success    bool              `json:"success" example:"true" comment:"是否成功"`
	Latency    string            `json:"latency" example:"150ms" comment:"延迟"`
	ErrorMsg   string            `json:"error_msg,omitempty" example:"" comment:"错误信息"`
	Details    map[string]string `json:"details,omitempty" comment:"详细信息"`
	TestedAt   time.Time         `json:"tested_at" example:"2024-01-01T00:00:00Z" comment:"测试时间"`
	TestType   string            `json:"test_type" example:"connection" comment:"测试类型"`
	Endpoint   string            `json:"endpoint" example:"https://alidns.aliyuncs.com" comment:"端点"`
	StatusCode int               `json:"status_code" example:"200" comment:"状态码"`
}
