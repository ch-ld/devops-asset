package redis

import (
	"fmt"
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
}

var store base64Captcha.Store

// GetCaptchaStore 获取验证码存储实例
func GetCaptchaStore() base64Captcha.Store {
	if store == nil {
		store = NewCaptchaStore(DefaultRedisExpiration, DefaultRedisPrefixKey)
	}
	return store
}

// NewCaptchaStore 创建新的验证码存储实例
func NewCaptchaStore(expiration time.Duration, prefixKey string) base64Captcha.Store {
	s := &CaptchaStore{
		expiration: expiration,
		prefixKey:  prefixKey,
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
		zap.L().Error("Redis客户端不可用")
		return fmt.Errorf("Redis客户端不可用")
	}

	key := fmt.Sprintf("%s:%s", s.prefixKey, id)
	_, err := client.SetNX(key, value, s.expiration).Result()
	if err != nil {
		zap.L().Error("存储验证码失败", zap.Error(err), zap.String("key", key))
	}

	return err
}

// Get 获取验证码
func (s *CaptchaStore) Get(id string, clear bool) string {
	client := GetClient()
	if client == nil {
		zap.L().Error("Redis客户端不可用")
		return ""
	}

	key := fmt.Sprintf("%s:%s", s.prefixKey, id)
	value, err := client.Get(key).Result()
	if err != nil {
		zap.L().Error("获取验证码失败", zap.Error(err), zap.String("key", key))
		return ""
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
		zap.L().Error("Redis客户端不可用")
		return false
	}

	key := fmt.Sprintf("%s:%s", s.prefixKey, id)
	value, err := client.Get(key).Result()
	if err != nil {
		zap.L().Error("验证码验证失败", zap.Error(err), zap.String("key", key))
		return false
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
