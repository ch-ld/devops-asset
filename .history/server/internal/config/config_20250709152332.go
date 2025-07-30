package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 全局配置结构
type Config struct {
	App        AppConfig        `yaml:"app"`
	Database   DatabaseConfig   `yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	JWT        JWTConfig        `yaml:"jwt"`
	Log        LogConfig        `yaml:"log"`
	Admin      AdminConfig      `yaml:"admin"`
	Pagination PaginationConfig `yaml:"pagination"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string   `yaml:"name"`
	Version string   `yaml:"version"`
	Mode    string   `yaml:"mode"`
	Port    int      `yaml:"port"`
	AesKey  string   `yaml:"aes_key"`  // 兼容单密钥
	AesKeys []string `yaml:"aes_keys"` // 支持多密钥轮换
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type  string      `yaml:"type"`
	MySQL MySQLConfig `yaml:"mysql"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Database  string `yaml:"database"`
	Charset   string `yaml:"charset"`
	ParseTime bool   `yaml:"parseTime"`
	Loc       string `yaml:"loc"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	Expiration string `yaml:"expiration"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	MaxAge     int    `yaml:"maxAge"`
	Compress   bool   `yaml:"compress"`
}

// AdminConfig 管理员配置
type AdminConfig struct {
	Password string `yaml:"password"`
	Salt     string `yaml:"salt"`
}

// PaginationConfig 分页配置
type PaginationConfig struct {
	DefaultPage     int `yaml:"defaultPage"`
	DefaultPageSize int `yaml:"defaultPageSize"`
	MaxPageSize     int `yaml:"maxPageSize"`
}

// 全局配置实例
var GlobalConfig *Config

// 兼容性变量（保持向后兼容）
var (
	// App相关
	ListenPort       int
	RunModel         string
	RunModelDevValue = "dev"
	RunModelRelease  = "release"

	// JWT相关
	JWTKey        string
	JWTExpiration time.Duration

	// Redis相关
	RedisHost     string
	RedisPassword string

	// 数据库相关
	DatabaseType   string
	MySQLConfigVar MySQLConfig

	// 管理员相关
	AdminPassword string
	PWDSalt       string

	// 分页相关
	DefaultPageSize = 20
	DefaultPage     = 1
	CancelPageSize  = -1
	CancelPage      = -1

	// 路径相关
	SelfName string
	AbsPath  string
	LogDir   string
	LogPath  string
)

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	// 如果没有指定配置文件路径，则根据运行模式自动选择
	if configPath == "" {
		mode := os.Getenv("APP_MODE")
		if mode == "" {
			mode = "prod"
		}

		if mode == "dev" {
			configPath = "configs/config.dev.yaml"
		} else {
			configPath = "configs/config.yaml"
		}
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置文件
	GlobalConfig = &Config{}
	if err := yaml.Unmarshal(data, GlobalConfig); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// 环境变量覆盖
	overrideWithEnv()

	// 设置兼容性变量
	setCompatibilityVars()

	return nil
}

// overrideWithEnv 使用环境变量覆盖配置
func overrideWithEnv() {
	if mode := os.Getenv("APP_MODE"); mode != "" {
		GlobalConfig.App.Mode = mode
	}
	if port := os.Getenv("APP_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &GlobalConfig.App.Port)
	}
	if dbType := os.Getenv("DATABASE_TYPE"); dbType != "" {
		GlobalConfig.Database.Type = dbType
	}
	if host := os.Getenv("MYSQL_HOST"); host != "" {
		GlobalConfig.Database.MySQL.Host = host
	}
	if port := os.Getenv("MYSQL_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &GlobalConfig.Database.MySQL.Port)
	}
	if user := os.Getenv("MYSQL_USER"); user != "" {
		GlobalConfig.Database.MySQL.Username = user
	}
	if password := os.Getenv("MYSQL_PASSWORD"); password != "" {
		GlobalConfig.Database.MySQL.Password = password
	}
	if database := os.Getenv("MYSQL_DB"); database != "" {
		GlobalConfig.Database.MySQL.Database = database
	}
	if host := os.Getenv("REDIS_HOST"); host != "" {
		GlobalConfig.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &GlobalConfig.Redis.Port)
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		GlobalConfig.Redis.Password = password
	}
}

// setCompatibilityVars 设置兼容性变量
func setCompatibilityVars() {
	// App相关
	ListenPort = GlobalConfig.App.Port
	RunModel = GlobalConfig.App.Mode

	// JWT相关
	JWTKey = GlobalConfig.JWT.Secret
	if duration, err := time.ParseDuration(GlobalConfig.JWT.Expiration); err == nil {
		JWTExpiration = duration
	} else {
		JWTExpiration = 12 * time.Hour
	}

	// Redis相关
	RedisHost = fmt.Sprintf("%s:%d", GlobalConfig.Redis.Host, GlobalConfig.Redis.Port)
	RedisPassword = GlobalConfig.Redis.Password

	// 数据库相关
	DatabaseType = GlobalConfig.Database.Type
	MySQLConfigVar = GlobalConfig.Database.MySQL

	// 管理员相关
	AdminPassword = GlobalConfig.Admin.Password
	PWDSalt = GlobalConfig.Admin.Salt

	// 分页相关
	DefaultPageSize = GlobalConfig.Pagination.DefaultPageSize
	DefaultPage = GlobalConfig.Pagination.DefaultPage

	// 路径相关
	SelfName = filepath.Base(os.Args[0])
	if wd, err := os.Getwd(); err == nil {
		AbsPath = wd
	}
	LogDir = filepath.Join(AbsPath, "log")
	LogPath = filepath.Join(LogDir, fmt.Sprintf("%s.log", SelfName))

	// 创建日志目录
	os.MkdirAll(LogDir, 0755)
}

// IsMySQL 判断是否使用MySQL
func IsMySQL() bool {
	return GlobalConfig.Database.Type == "mysql"
}

// IsDev 判断是否为开发模式
func IsDev() bool {
	return GlobalConfig.App.Mode == "dev"
}

// GetMySQLDSN 获取MySQL连接字符串
func GetMySQLDSN() string {
	cfg := GlobalConfig.Database.MySQL
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset, cfg.ParseTime, cfg.Loc)
}

// GetRedisAddr 获取Redis地址
func GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", GlobalConfig.Redis.Host, GlobalConfig.Redis.Port)
}
