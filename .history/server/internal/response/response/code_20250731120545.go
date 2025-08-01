package response

// https://google-cloud.gitbook.io/api-design-guide/errors

type responseData struct {
	Code      int         `json:"code"`
	Status    string      `json:"status"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
	Count     *int        `json:"count,omitempty"`
}

// 通用的错误信息

// 没有错误
var Success = responseData{
	Code:    200,
	Status:  "OK",
	Message: "请求成功",
}

// 客户端发送的数据包含非法参数。查看错误消息和错误详情来获取更多的信息
var INVALID_ARGUMENT = responseData{
	Code:    400,
	Status:  "INVALID_ARGUMENT",
	Message: "请求参数错误",
}

// 现在的系统状态不可以执行当前的请求，例如删除一个非空的目录
var FAILED_PRECONDITION = responseData{
	Code:    400,
	Status:  "FAILED_PRECONDITION",
	Message: "无法执行客户端请求",
}

// 客户端指定了一个非法的范围
var OUT_OF_RANGE = responseData{
	Code:    400,
	Status:  "OUT_OF_RANGE",
	Message: "客户端越限访问",
}

// 因为缺失的，失效的或者过期的OAuth令牌，请求未能通过身份认证
var UNAUTHENTICATED = responseData{
	Code:    401,
	Status:  "UNAUTHENTICATED",
	Message: "身份验证失败",
}

// 客户端没有足够的权限。这可能是因为OAuth令牌没有正确的作用域，或者客户端没有权限，或者是API对客户端代码禁用了
var PERMISSION_DENIED = responseData{
	Code:    403,
	Status:  "PERMISSION_DENIED",
	Message: "客户端权限不足",
}

// 特定的资源没有被找到或者请求因为某些未被公开的原因拒绝（例如白名单）
var NOT_FOUND = responseData{
	Code:    404,
	Status:  "NOT_FOUND",
	Message: "资源不存在",
}

// 并发冲突，如读-修改-写冲突
var ABORTED = responseData{
	Code:    409,
	Status:  "ABORTED",
	Message: "数据处理冲突",
}

// 客户端尝试新建的资源已经存在了
var ALREADY_EXISTS = responseData{
	Code:    409,
	Status:  "ALREADY_EXISTS",
	Message: "资源已存在",
}

// 资源配额不足或达不到速率限制
var RESOURCE_EXHAUSTED = responseData{
	Code:    429,
	Status:  "RESOURCE_EXHAUSTED",
	Message: "资源配额不足或达不到速率限制",
}

// 请求被客户端取消了
var CANCELLED = responseData{
	Code:    499,
	Status:  "CANCELLED",
	Message: "请求被客户端取消",
}

// 不可恢复的数据丢失或数据损坏。客户端应该向用户报告错误
var DATA_LOSS = responseData{
	Code:    500,
	Status:  "DATA_LOSS",
	Message: "处理数据发生错误",
}

// 未知的服务端出错，通常是由于服务器出现bug了
var UNKNOWN = responseData{
	Code:    500,
	Status:  "UNKNOWN",
	Message: "服务器未知错误",
}

// 服务器内部错误。通常是由于服务器出现bug了
var INTERNAL = responseData{
	Code:    500,
	Status:  "INTERNAL",
	Message: "服务器内部错误",
}

// API方法没有被服务器实现
var NOT_IMPLEMENTED = responseData{
	Code:    501,
	Status:  "NOT_IMPLEMENTED",
	Message: "API不存在",
}

// 服务不可用。通常是由于服务器宕机了
var UNAVAILABLE = responseData{
	Code:    503,
	Status:  "UNAVAILABLE",
	Message: "服务不可用",
}

// 请求超过了截止日期。只有当调用者设置的截止日期比方法的默认截止日期更短（服务器没能够在截止日期之前处理完请求）并且请求没有在截止日期内完成时，才会发生这种情况
var DEALINE_EXCEED = responseData{
	Code:    504,
	Status:  "DEALINE_EXCEED",
	Message: "请求超时",
}

// 根据业务自定义的常用状态码

// 主机管理相关响应码定义
// 包含主机、云账号、同步、告警等模块的响应码

// DNS管理相关响应码定义 (67001-67201)
// 域名管理错误码 (67001-67050)

// 域名不存在
var DNS_DOMAIN_NOT_FOUND = responseData{
	Code:    67001,
	Status:  "DNS_DOMAIN_NOT_FOUND",
	Message: "域名不存在",
}

// 域名已存在
var DNS_DOMAIN_ALREADY_EXISTS = responseData{
	Code:    67002,
	Status:  "DNS_DOMAIN_ALREADY_EXISTS",
	Message: "域名已存在",
}

// 域名格式错误
var DNS_DOMAIN_INVALID_FORMAT = responseData{
	Code:    67003,
	Status:  "DNS_DOMAIN_INVALID_FORMAT",
	Message: "域名格式错误",
}

// 域名已过期
var DNS_DOMAIN_EXPIRED = responseData{
	Code:    67004,
	Status:  "DNS_DOMAIN_EXPIRED",
	Message: "域名已过期",
}

// 域名分组不存在
var DNS_DOMAIN_GROUP_NOT_FOUND = responseData{
	Code:    67005,
	Status:  "DNS_DOMAIN_GROUP_NOT_FOUND",
	Message: "域名分组不存在",
}

// 域名分组已存在
var DNS_DOMAIN_GROUP_ALREADY_EXISTS = responseData{
	Code:    67006,
	Status:  "DNS_DOMAIN_GROUP_ALREADY_EXISTS",
	Message: "域名分组已存在",
}

// 域名分组有子分组，无法删除
var DNS_DOMAIN_GROUP_HAS_CHILDREN = responseData{
	Code:    67007,
	Status:  "DNS_DOMAIN_GROUP_HAS_CHILDREN",
	Message: "域名分组存在子分组，无法删除",
}

// 域名分组有关联域名，无法删除
var DNS_DOMAIN_GROUP_HAS_DOMAINS = responseData{
	Code:    67008,
	Status:  "DNS_DOMAIN_GROUP_HAS_DOMAINS",
	Message: "域名分组存在关联域名，无法删除",
}

// DNS记录管理错误码 (67051-67100)

// DNS记录不存在
var DNS_RECORD_NOT_FOUND = responseData{
	Code:    67051,
	Status:  "DNS_RECORD_NOT_FOUND",
	Message: "DNS记录不存在",
}

// DNS记录已存在
var DNS_RECORD_ALREADY_EXISTS = responseData{
	Code:    67052,
	Status:  "DNS_RECORD_ALREADY_EXISTS",
	Message: "DNS记录已存在",
}

// DNS记录类型不支持
var DNS_RECORD_TYPE_UNSUPPORTED = responseData{
	Code:    67053,
	Status:  "DNS_RECORD_TYPE_UNSUPPORTED",
	Message: "不支持的DNS记录类型",
}

// DNS记录值格式错误
var DNS_RECORD_VALUE_INVALID = responseData{
	Code:    67054,
	Status:  "DNS_RECORD_VALUE_INVALID",
	Message: "DNS记录值格式错误",
}

// DNS记录TTL值无效
var DNS_RECORD_TTL_INVALID = responseData{
	Code:    67055,
	Status:  "DNS_RECORD_TTL_INVALID",
	Message: "DNS记录TTL值无效",
}

// DNS记录冲突
var DNS_RECORD_CONFLICT = responseData{
	Code:    67056,
	Status:  "DNS_RECORD_CONFLICT",
	Message: "DNS记录冲突",
}

// DNS记录同步失败
var DNS_RECORD_SYNC_FAILED = responseData{
	Code:    67057,
	Status:  "DNS_RECORD_SYNC_FAILED",
	Message: "DNS记录同步失败",
}

// DNS提供商管理错误码 (67101-67150)

// DNS提供商不存在
var DNS_PROVIDER_NOT_FOUND = responseData{
	Code:    67101,
	Status:  "DNS_PROVIDER_NOT_FOUND",
	Message: "DNS提供商不存在",
}

// DNS提供商已存在
var DNS_PROVIDER_ALREADY_EXISTS = responseData{
	Code:    67102,
	Status:  "DNS_PROVIDER_ALREADY_EXISTS",
	Message: "DNS提供商已存在",
}

// DNS提供商类型不支持
var DNS_PROVIDER_TYPE_UNSUPPORTED = responseData{
	Code:    67103,
	Status:  "DNS_PROVIDER_TYPE_UNSUPPORTED",
	Message: "不支持的DNS提供商类型",
}

// DNS提供商凭证无效
var DNS_PROVIDER_CREDENTIALS_INVALID = responseData{
	Code:    67104,
	Status:  "DNS_PROVIDER_CREDENTIALS_INVALID",
	Message: "DNS提供商凭证无效",
}

// DNS提供商连接失败
var DNS_PROVIDER_CONNECTION_FAILED = responseData{
	Code:    67105,
	Status:  "DNS_PROVIDER_CONNECTION_FAILED",
	Message: "DNS提供商连接失败",
}

// DNS提供商认证失败
var DNS_PROVIDER_AUTH_FAILED = responseData{
	Code:    67106,
	Status:  "DNS_PROVIDER_AUTH_FAILED",
	Message: "DNS提供商认证失败",
}

// DNS提供商速率限制
var DNS_PROVIDER_RATE_LIMITED = responseData{
	Code:    67107,
	Status:  "DNS_PROVIDER_RATE_LIMITED",
	Message: "DNS提供商请求频率过高",
}

// DNS提供商配额不足
var DNS_PROVIDER_QUOTA_EXCEEDED = responseData{
	Code:    67108,
	Status:  "DNS_PROVIDER_QUOTA_EXCEEDED",
	Message: "DNS提供商配额不足",
}

// 证书管理错误码 (67151-67200)

// 证书不存在
var DNS_CERTIFICATE_NOT_FOUND = responseData{
	Code:    67151,
	Status:  "DNS_CERTIFICATE_NOT_FOUND",
	Message: "证书不存在",
}

// 证书已存在
var DNS_CERTIFICATE_ALREADY_EXISTS = responseData{
	Code:    67152,
	Status:  "DNS_CERTIFICATE_ALREADY_EXISTS",
	Message: "证书已存在",
}

// 证书已过期
var DNS_CERTIFICATE_EXPIRED = responseData{
	Code:    67153,
	Status:  "DNS_CERTIFICATE_EXPIRED",
	Message: "证书已过期",
}

// 证书签发失败
var DNS_CERTIFICATE_ISSUE_FAILED = responseData{
	Code:    67154,
	Status:  "DNS_CERTIFICATE_ISSUE_FAILED",
	Message: "证书签发失败",
}

// 证书续期失败
var DNS_CERTIFICATE_RENEW_FAILED = responseData{
	Code:    67155,
	Status:  "DNS_CERTIFICATE_RENEW_FAILED",
	Message: "证书续期失败",
}

// 证书部署失败
var DNS_CERTIFICATE_DEPLOY_FAILED = responseData{
	Code:    67156,
	Status:  "DNS_CERTIFICATE_DEPLOY_FAILED",
	Message: "证书部署失败",
}

// 证书验证失败
var DNS_CERTIFICATE_VALIDATION_FAILED = responseData{
	Code:    67157,
	Status:  "DNS_CERTIFICATE_VALIDATION_FAILED",
	Message: "证书验证失败",
}

// DNS验证失败
var DNS_CHALLENGE_FAILED = responseData{
	Code:    67158,
	Status:  "DNS_CHALLENGE_FAILED",
	Message: "DNS验证失败",
}
