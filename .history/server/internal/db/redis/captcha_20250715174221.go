package redis

import (
	"fmt"
	"sync"
	"time"

	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

const (
	// DefaultRedisExpiration 验证码默认过期时间（5分钟）
	DefaultRedisExpiration = 5 * time.Minute
	// DefaultRedisPrefixKey Redis中验证码的键前缀
	DefaultRedisPrefixKey = "captcha"
)

// CaptchaStore 验证码存储结构
type CaptchaStore struct {
	expiration time.Duration
	prefixKey  string
	// 备选内存存储
	memoryStore *base64Captcha.MemoryStore
}

var store base64Captcha.Store
var once sync.Once

// GetCaptchaStore 获取验证码存储实例
func GetCaptchaStore() base64Captcha.Store {
	once.Do(func() {
		store = NewCaptchaStore(DefaultRedisExpiration, DefaultRedisPrefixKey)
	})
	return store
}

// NewCaptchaStore 创建新的验证码存储实例
func NewCaptchaStore(expiration time.Duration, prefixKey string) base64Captcha.Store {
	s := &CaptchaStore{
		expiration:  expiration,
		prefixKey:   prefixKey,
		memoryStore: base64Captcha.NewMemoryStore(20), // 默认存储20个验证码
	}

	if s.prefixKey == "" {
		s.prefixKey = DefaultRedisPrefixKey
	}
	if s.expiration == 0 {
		s.expiration = DefaultRedisExpiration
	}

	return s
}

// Set 存储验证码
func (s *CaptchaStore) Set(id string, value string) error {
	client := GetClient()
	if client == nil {
		zap.L().Warn("Redis不可用，使用内存存储验证码")
		s.memoryStore.Set(id, value)
		return nil
	}

	key := fmt.Sprintf("%s:%s", s.prefixKey, id)
	_, err := client.SetNX(key, value, s.expiration).Result()
	if err != nil {
		zap.L().Error("存储验证码失败", zap.Error(err), zap.String("key", key))
		// 降级到内存存储
		zap.L().Info("降级到内存存储验证码")
		s.memoryStore.Set(id, value)
		return nil
	}

	return err
}

// Get 获取验证码
func (s *CaptchaStore) Get(id string, clear bool) string {
	client := GetClient()
	if client == nil {
		zap.L().Warn("Redis不可用，从内存获取验证码")
		return s.memoryStore.Get(id, clear)
	}

	key := fmt.Sprintf("%s:%s", s.prefixKey, id)
	value, err := client.Get(key).Result()
	if err != nil {
		zap.L().Error("获取验证码失败", zap.Error(err), zap.String("key", key))
		// 降级到内存存储
		return s.memoryStore.Get(id, clear)
	}

	if clear {
		if _, delErr := client.Del(key).Result(); delErr != nil {
			zap.L().Error("删除验证码失败", zap.Error(delErr), zap.String("key", key))
		}
	}

	return value
}

// Verify 验证验证码
func (s *CaptchaStore) Verify(id, answer string, clear bool) bool {
	client := GetClient()
	if client == nil {
		zap.L().Warn("Redis不可用，从内存验证验证码")
		return s.memoryStore.Verify(id, answer, clear)
	}

	key := fmt.Sprintf("%s:%s", s.prefixKey, id)
	value, err := client.Get(key).Result()
	if err != nil {
		zap.L().Error("验证码验证失败", zap.Error(err), zap.String("key", key))
		// 降级到内存存储
		return s.memoryStore.Verify(id, answer, clear)
	}

	if value == answer {
		if clear {
			if _, delErr := client.Del(key).Result(); delErr != nil {
				zap.L().Error("验证后删除验证码失败", zap.Error(delErr), zap.String("key", key))
				return false
			}
		}
		return true
	}

	return false
}
