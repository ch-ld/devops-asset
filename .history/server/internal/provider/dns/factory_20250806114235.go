package dns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// DriverFactory DNS驱动工厂接口
type DriverFactory interface {
	// CreateDriver 创建驱动实例
	CreateDriver(config Config) (Driver, error)

	// GetSupportedTypes 获取支持的提供商类型
	GetSupportedTypes() []string

	// ValidateConfig 验证配置
	ValidateConfig(providerType string, config map[string]string) error

	// GetDefaultConfig 获取默认配置
	GetDefaultConfig(providerType string) map[string]interface{}
}

// DriverConstructor 驱动构造函数类型
type DriverConstructor func(config Config) (Driver, error)

// DefaultDriverFactory 默认驱动工厂实现
type DefaultDriverFactory struct {
	drivers map[string]DriverConstructor
	mutex   sync.RWMutex
}

// NewDriverFactory 创建驱动工厂
func NewDriverFactory() *DefaultDriverFactory {
	factory := &DefaultDriverFactory{
		drivers: make(map[string]DriverConstructor),
	}

	// 注册内置驱动
	factory.registerBuiltinDrivers()

	return factory
}

// registerBuiltinDrivers 注册内置驱动
func (f *DefaultDriverFactory) registerBuiltinDrivers() {
	// 注册Route53驱动
	f.RegisterDriver("route53", NewRoute53DriverFromConfig)

	// 注册阿里云DNS驱动
	f.RegisterDriver("aliyun", NewAliyunDriverFromConfig)

	// 注册腾讯云DNS驱动
	f.RegisterDriver("tencent", NewTencentDriverFromConfig)

	// 注册GoDaddy驱动
	f.RegisterDriver("godaddy", NewGoDaddyDriverFromConfig)

	// 注册Cloudflare驱动
	f.RegisterDriver("cloudflare", NewCloudflareDriver)

	// 注册DNSPod驱动
	f.RegisterDriver("dnspod", NewDNSPodDriver)
}

// RegisterDriver 注册驱动
func (f *DefaultDriverFactory) RegisterDriver(providerType string, constructor DriverConstructor) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.drivers[providerType] = constructor
}

// UnregisterDriver 注销驱动
func (f *DefaultDriverFactory) UnregisterDriver(providerType string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	delete(f.drivers, providerType)
}

// CreateDriver 创建驱动实例
func (f *DefaultDriverFactory) CreateDriver(config Config) (Driver, error) {
	f.mutex.RLock()
	constructor, exists := f.drivers[config.GetType()]
	f.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unsupported provider type: %s", config.GetType())
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// 创建驱动实例
	driver, err := constructor(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	return driver, nil
}

// GetSupportedTypes 获取支持的提供商类型
func (f *DefaultDriverFactory) GetSupportedTypes() []string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	types := make([]string, 0, len(f.drivers))
	for providerType := range f.drivers {
		types = append(types, providerType)
	}

	return types
}

// ValidateConfig 验证配置
func (f *DefaultDriverFactory) ValidateConfig(providerType string, config map[string]string) error {
	f.mutex.RLock()
	_, exists := f.drivers[providerType]
	f.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("unsupported provider type: %s", providerType)
	}

	// 根据提供商类型验证必需的配置项
	switch providerType {
	case "route53":
		return f.validateRoute53Config(config)
	case "aliyun":
		return f.validateAliyunConfig(config)
	case "godaddy":
		return f.validateGoDaddyConfig(config)
	case "cloudflare":
		return f.validateCloudflareConfig(config)
	case "dnspod":
		return f.validateDNSPodConfig(config)
	default:
		return fmt.Errorf("validation not implemented for provider type: %s", providerType)
	}
}

// validateRoute53Config 验证Route53配置
func (f *DefaultDriverFactory) validateRoute53Config(config map[string]string) error {
	required := []string{"access_key_id", "secret_access_key", "region"}
	for _, key := range required {
		if config[key] == "" {
			return fmt.Errorf("missing required config: %s", key)
		}
	}
	return nil
}

// validateAliyunConfig 验证阿里云配置
func (f *DefaultDriverFactory) validateAliyunConfig(config map[string]string) error {
	required := []string{"access_key_id", "access_key_secret", "region"}
	for _, key := range required {
		if config[key] == "" {
			return fmt.Errorf("missing required config: %s", key)
		}
	}
	return nil
}

// validateGoDaddyConfig 验证GoDaddy配置
func (f *DefaultDriverFactory) validateGoDaddyConfig(config map[string]string) error {
	required := []string{"api_key", "api_secret"}
	for _, key := range required {
		if config[key] == "" {
			return fmt.Errorf("missing required config: %s", key)
		}
	}
	return nil
}

// validateCloudflareConfig 验证Cloudflare配置
func (f *DefaultDriverFactory) validateCloudflareConfig(config map[string]string) error {
	// 支持两种认证方式：API Token 或 API Key + Email
	if config["api_token"] != "" {
		return nil
	}

	if config["api_key"] != "" && config["email"] != "" {
		return nil
	}

	return fmt.Errorf("missing required config: either 'api_token' or both 'api_key' and 'email'")
}

// validateDNSPodConfig 验证DNSPod配置
func (f *DefaultDriverFactory) validateDNSPodConfig(config map[string]string) error {
	required := []string{"login_token"}
	for _, key := range required {
		if config[key] == "" {
			return fmt.Errorf("missing required config: %s", key)
		}
	}
	return nil
}

// GetDefaultConfig 获取默认配置
func (f *DefaultDriverFactory) GetDefaultConfig(providerType string) map[string]interface{} {
	defaults := map[string]interface{}{
		"timeout":    30,
		"rate_limit": 10,
		"concurrent": 5,
	}

	switch providerType {
	case "route53":
		defaults["region"] = "us-east-1"
		defaults["endpoint"] = "https://route53.amazonaws.com"
	case "aliyun":
		defaults["region"] = "cn-hangzhou"
		defaults["endpoint"] = "https://alidns.aliyuncs.com"
	case "godaddy":
		defaults["endpoint"] = "https://api.godaddy.com"
	case "cloudflare":
		defaults["endpoint"] = "https://api.cloudflare.com/client/v4"
	case "dnspod":
		defaults["endpoint"] = "https://dnsapi.cn"
	}

	return defaults
}

// 全局工厂实例
var globalFactory *DefaultDriverFactory
var factoryOnce sync.Once

// GetFactory 获取全局工厂实例
func GetFactory() *DefaultDriverFactory {
	factoryOnce.Do(func() {
		globalFactory = NewDriverFactory()
	})
	return globalFactory
}

// CreateDriver 创建驱动实例（便捷函数）
func CreateDriver(providerType, name string, credentials map[string]string) (Driver, error) {
	config := &BaseConfig{
		Type:        providerType,
		Name:        name,
		Credentials: credentials,
		Timeout:     30 * time.Second,
		RateLimit:   10,
		Concurrent:  5,
	}

	// 设置默认配置
	defaults := GetFactory().GetDefaultConfig(providerType)
	if endpoint, ok := defaults["endpoint"].(string); ok {
		config.Endpoint = endpoint
	}
	if region, ok := defaults["region"].(string); ok {
		config.Region = region
	}

	return GetFactory().CreateDriver(config)
}

// TestDriver 测试驱动连接（便捷函数）
func TestDriver(providerType, name string, credentials map[string]string) *TestResult {
	driver, err := CreateDriver(providerType, name, credentials)
	if err != nil {
		return &TestResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("Failed to create driver: %v", err),
			TestedAt: time.Now(),
			TestType: "connection",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return driver.Test(ctx)
}

// NewRoute53DriverFromConfig 从通用配置创建Route53驱动
func NewRoute53DriverFromConfig(config Config) (Driver, error) {
	baseConfig, ok := config.(*BaseConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type for Route53 driver")
	}

	// 转换为Route53特定配置
	route53Config := &Route53Config{
		BaseConfig:      baseConfig,
		AccessKeyID:     baseConfig.Credentials["access_key_id"],
		SecretAccessKey: baseConfig.Credentials["secret_access_key"],
		SessionToken:    baseConfig.Credentials["session_token"],
	}

	return NewRoute53Driver(route53Config)
}

func NewAliyunDriverFromConfig(config Config) (Driver, error) {
	baseConfig, ok := config.(*BaseConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type for Aliyun driver")
	}

	// 转换为阿里云特定配置
	aliyunConfig := &AliyunConfig{
		BaseConfig:      baseConfig,
		AccessKeyID:     baseConfig.Credentials["access_key_id"],
		AccessKeySecret: baseConfig.Credentials["access_key_secret"],
	}

	return NewAliyunDriver(aliyunConfig)
}

func NewTencentDriverFromConfig(config Config) (Driver, error) {
	baseConfig, ok := config.(*BaseConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type for Tencent driver")
	}

	// 转换为腾讯云特定配置
	tencentConfig := &TencentConfig{
		BaseConfig: baseConfig,
		SecretID:   baseConfig.Credentials["secret_id"],
		SecretKey:  baseConfig.Credentials["secret_key"],
	}

	return NewTencentDriver(tencentConfig)
}

func NewGoDaddyDriverFromConfig(config Config) (Driver, error) {
	baseConfig, ok := config.(*BaseConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type for GoDaddy driver")
	}

	// 转换为GoDaddy特定配置
	goDaddyConfig := &GoDaddyConfig{
		BaseConfig: baseConfig,
		APIKey:     baseConfig.Credentials["api_key"],
		APISecret:  baseConfig.Credentials["api_secret"],
	}

	return NewGoDaddyDriver(goDaddyConfig)
}

func NewCloudflareDriver(config Config) (Driver, error) {
	return NewStubDriver("cloudflare", "Cloudflare"), nil
}

func NewDNSPodDriver(config Config) (Driver, error) {
	return NewStubDriver("dnspod", "DNSPod"), nil
}
