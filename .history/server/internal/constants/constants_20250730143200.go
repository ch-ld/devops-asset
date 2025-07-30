package constants

import "time"

// 应用相关常量
const (
	// AppName 应用名称
	AppName = "server"

	// 运行模式
	RunModeProduction  = "prod"
	RunModeDevelopment = "dev"
	RunModeRelease     = "release"
)

// HTTP相关常量
const (
	// 默认端口
	DefaultPort = 8080

	// API版本
	APIVersionV1 = "/api/v1"

	// 静态文件路径
	StaticPath = "/static"
	StaticDir  = "./static"
)

// 数据库相关常量
const (
	// 数据库类型
	DatabaseTypeMySQL = "mysql"

	// 默认数据库配置
	DefaultDatabaseHost = "localhost"
	DefaultRedisHost    = "127.0.0.1"
	DefaultRedisPort    = 6379
	DefaultRedisDB      = 0
)

// JWT相关常量
const (
	// 默认JWT过期时间
	DefaultJWTExpiration = 12 * time.Hour

	// JWT Header
	JWTHeaderKey = "Authorization"
	JWTPrefix    = "Bearer "
)

// 分页相关常量
const (
	// 默认分页参数
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100

	// 取消分页的标识
	CancelPage     = -1
	CancelPageSize = -1
)

// 日志相关常量
const (
	// 日志级别
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"

	// 默认日志配置
	DefaultLogMaxSize    = 50 // MB
	DefaultLogMaxBackups = 3
	DefaultLogMaxAge     = 30 // days
)

// 缓存相关常量
const (
	// 缓存键前缀
	CacheKeyPrefix = "devops_asset:"

	// 用户缓存相关
	UserCacheKeyPrefix = CacheKeyPrefix + "user:"
	UserListCacheKey   = CacheKeyPrefix + "user:list"

	// 缓存过期时间
	DefaultCacheExpiration = 12 * time.Hour
)

// 文件路径相关常量
const (
	// 配置文件路径
	ConfigDirPath         = "configs"
	DefaultConfigFile     = "config.yaml"
	DevelopmentConfigFile = "config.dev.yaml"

	// 日志目录
	LogDirPath = "log"

	// 静态文件目录
	StaticDirPath = "static"
)

// 系统相关常量
const (
	// 系统管理员ID
	SuperAdminUserID = 1

	// 默认状态
	StatusEnabled  = 1
	StatusDisabled = 2

	// 性别
	GenderMale   = 1
	GenderFemale = 2
)

// 错误码相关常量
const (
	// 成功
	CodeSuccess = 0

	// 通用错误码
	CodeInvalidParam    = 400
	CodeUnauthorized    = 401
	CodeForbidden       = 403
	CodeNotFound        = 404
	CodeInternalError   = 500
	CodeDataLoss        = 1001
	CodeInvalidArgument = 1002
)

// 消息相关常量
const (
	// 成功消息
	MsgSuccess = "操作成功"

	// 错误消息
	MsgInvalidParam  = "参数错误"
	MsgUnauthorized  = "未授权"
	MsgForbidden     = "禁止访问"
	MsgNotFound      = "资源不存在"
	MsgInternalError = "内部服务器错误"
)

// 主机管理相关常量定义
// 包含主机状态、云账号类型、告警类型等常量
