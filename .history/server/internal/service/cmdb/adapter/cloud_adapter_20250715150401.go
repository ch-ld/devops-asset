package adapter

import (
	"api-server/internal/model/cmdb"
	"errors"
	"fmt"
)

// CloudAdapter 云服务提供商适配器接口
// 提供统一的云资源操作方法
type CloudAdapter interface {
	// ListInstances 获取云主机列表
	ListInstances() ([]cmdb.Host, error)
	// GetInstanceStatus 获取云主机状态
	GetInstanceStatus(instanceID string) (string, error)
	// GetInstanceInfo 获取云主机详细信息
	GetInstanceInfo(instanceID string) (*cmdb.Host, error)
}

// BaseAdapter 基础适配器实现
type BaseAdapter struct {
	accessKey string
	secretKey string
}

// NewBaseAdapter 创建基础适配器
func NewBaseAdapter(accessKey, secretKey string) *BaseAdapter {
	return &BaseAdapter{
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

// GetCloudAdapter 根据提供商类型获取对应适配器
func GetCloudAdapter(provider *cmdb.Provider) (CloudAdapter, error) {
	if provider == nil {
		return nil, errors.New("provider 不能为空")
	}

	switch provider.Type {
	case "aliyun":
		return NewAliyunAdapter(provider.AccessKey, provider.SecretKey)
	case "aws":
		return NewAWSAdapter(provider.AccessKey, provider.SecretKey, provider.Region)
	case "tencent":
		return NewTencentAdapter(provider.AccessKey, provider.SecretKey, provider.Region)
	default:
		return nil, fmt.Errorf("不支持的云服务提供商类型: %s", provider.Type)
	}
}
