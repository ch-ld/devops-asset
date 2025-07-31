package adapter

import (
	"api-server/internal/model/cmdb"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"gorm.io/datatypes"
)

// AliyunAdapter 阿里云适配器
type AliyunAdapter struct {
	*BaseAdapter
	client *ecs.Client
}

// NewAliyunAdapter 创建阿里云适配器实例
func NewAliyunAdapter(accessKey, secretKey string) (*AliyunAdapter, error) {
	// 使用默认区域创建客户端，用于获取区域列表等全局操作
	client, err := ecs.NewClientWithAccessKey("cn-hangzhou", accessKey, secretKey)
	if err != nil {
		return nil, err
	}

	return &AliyunAdapter{
		BaseAdapter: NewBaseAdapter(accessKey, secretKey),
		client:      client,
	}, nil
}

// NewAliyunAdapterWithRegion 创建指定区域的阿里云适配器实例
func NewAliyunAdapterWithRegion(accessKey, secretKey, region string) (*AliyunAdapter, error) {
	if region == "" {
		region = "cn-hangzhou" // 默认区域
	}

	client, err := ecs.NewClientWithAccessKey(region, accessKey, secretKey)
	if err != nil {
		return nil, err
	}

	return &AliyunAdapter{
		BaseAdapter: NewBaseAdapter(accessKey, secretKey),
		client:      client,
	}, nil
}

// 阿里云主机状态映射
var aliyunStatusMap = map[string]string{
	"Running":         "running",
	"Stopped":         "stopped",
	"Starting":        "starting",
	"Stopping":        "stopping",
	"Pending":         "pending",
	"Creating":        "creating",
	"Rebooting":       "rebooting",
	"Deleted":         "deleted",
	"Recycling":       "recycling",
	"Expired":         "expired",
	"Error":           "error",
	"Upgrading":       "upgrading",
	"Maintaining":     "maintaining",
	"ImageProcessing": "processing",
}

// convertStatus 转换阿里云主机状态为标准状态
func (a *AliyunAdapter) convertStatus(status string) string {
	if standardStatus, ok := aliyunStatusMap[status]; ok {
		return standardStatus
	}
	return "unknown"
}

// ListInstances 获取阿里云主机列表
func (a *AliyunAdapter) ListInstances() ([]cmdb.Host, error) {
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	request.PageSize = requests.NewInteger(100)

	var hosts []cmdb.Host
	pageNum := 1

	for {
		request.PageNumber = requests.NewInteger(pageNum)
		response, err := a.client.DescribeInstances(request)
		if err != nil {
			return nil, fmt.Errorf("获取阿里云主机列表失败: %w", err)
		}

		for _, instance := range response.Instances.Instance {
			// 转换公网IP和私网IP为JSON
			publicIPs, _ := json.Marshal(instance.PublicIpAddress.IpAddress)
			privateIPs, _ := json.Marshal(instance.VpcAttributes.PrivateIpAddress.IpAddress)

			// 构建配置信息
			config := map[string]interface{}{
				"cpu_cores":     instance.Cpu,
				"memory_size":   instance.Memory,
				"zone_id":       instance.ZoneId,
				"vpc_id":        instance.VpcAttributes.VpcId,
				"instance_type": instance.InstanceType,
				"charge_type":   instance.InstanceChargeType,
			}
			configJSON, _ := json.Marshal(config)

			// 解析过期时间
			var expiredAt *time.Time
			if instance.ExpiredTime != "" {
				t, err := time.Parse("2006-01-02T15:04Z", instance.ExpiredTime)
				if err == nil {
					expiredAt = &t
				}
			}

			host := cmdb.Host{
				InstanceID:    instance.InstanceId,
				Name:          instance.InstanceName,
				Region:        instance.RegionId,
				PublicIP:      datatypes.JSON(publicIPs),
				PrivateIP:     datatypes.JSON(privateIPs),
				Configuration: datatypes.JSON(configJSON),
				OS:            instance.OSName, // 使用OSName字段
				Status:        a.convertStatus(instance.Status),
				ProviderType:  "aliyun",
				ResourceType:  cmdb.ResourceTypeECS,
				ExpiredAt:     expiredAt,
			}
			hosts = append(hosts, host)
		}

		if pageNum*100 >= response.TotalCount {
			break
		}
		pageNum++
	}

	return hosts, nil
}

// GetInstanceStatus 获取阿里云主机状态
func (a *AliyunAdapter) GetInstanceStatus(instanceID string) (string, error) {
	request := ecs.CreateDescribeInstanceAttributeRequest()
	request.Scheme = "https"
	request.InstanceId = instanceID

	response, err := a.client.DescribeInstanceAttribute(request)
	if err != nil {
		return "", fmt.Errorf("获取阿里云主机状态失败: %w", err)
	}

	return a.convertStatus(response.Status), nil
}

// GetInstanceInfo 获取阿里云主机详细信息
func (a *AliyunAdapter) GetInstanceInfo(instanceID string) (*cmdb.Host, error) {
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	request.InstanceIds = fmt.Sprintf("[\"%s\"]", instanceID)

	response, err := a.client.DescribeInstances(request)
	if err != nil {
		return nil, fmt.Errorf("获取阿里云主机信息失败: %w", err)
	}

	if len(response.Instances.Instance) == 0 {
		return nil, fmt.Errorf("阿里云实例不存在: %s", instanceID)
	}

	instance := response.Instances.Instance[0]

	// 转换公网IP和私网IP为JSON
	publicIPs, _ := json.Marshal(instance.PublicIpAddress.IpAddress)
	privateIPs, _ := json.Marshal(instance.VpcAttributes.PrivateIpAddress.IpAddress)

	// 构建配置信息
	config := map[string]interface{}{
		"cpu_cores":     instance.Cpu,
		"memory_size":   instance.Memory,
		"zone_id":       instance.ZoneId,
		"vpc_id":        instance.VpcAttributes.VpcId,
		"instance_type": instance.InstanceType,
		"charge_type":   instance.InstanceChargeType,
	}
	configJSON, _ := json.Marshal(config)

	// 解析过期时间
	var expiredAt *time.Time
	if instance.ExpiredTime != "" {
		t, err := time.Parse("2006-01-02T15:04Z", instance.ExpiredTime)
		if err == nil {
			expiredAt = &t
		}
	}

	return &cmdb.Host{
		InstanceID:    instance.InstanceId,
		Name:          instance.InstanceName,
		Status:        a.convertStatus(instance.Status),
		PublicIP:      datatypes.JSON(publicIPs),
		PrivateIP:     datatypes.JSON(privateIPs),
		Configuration: datatypes.JSON(configJSON),
		OS:            instance.OSName, // 使用OSName字段
		ProviderType:  "aliyun",
		ResourceType:  cmdb.ResourceTypeECS,
		ExpiredAt:     expiredAt,
	}, nil
}

// GetRegions 获取阿里云支持的区域列表
func (a *AliyunAdapter) GetRegions() ([]Region, error) {
	request := ecs.CreateDescribeRegionsRequest()
	request.Scheme = "https"

	response, err := a.client.DescribeRegions(request)
	if err != nil {
		return nil, fmt.Errorf("获取阿里云区域列表失败: %w", err)
	}

	var regions []Region
	for _, region := range response.Regions.Region {
		regions = append(regions, Region{
			ID:   region.RegionId,
			Name: region.LocalName,
		})
	}

	return regions, nil
}
