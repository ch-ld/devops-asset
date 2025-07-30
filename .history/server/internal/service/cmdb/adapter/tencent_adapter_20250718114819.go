package adapter

import (
	"api-server/internal/model/cmdb"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"gorm.io/datatypes"
)

// TencentAdapter 腾讯云适配器
type TencentAdapter struct {
	*BaseAdapter
	client *cvm.Client
	region string
}

// NewTencentAdapter 创建腾讯云适配器实例
func NewTencentAdapter(secretID, secretKey, region string) (*TencentAdapter, error) {
	if region == "" {
		region = regions.Guangzhou
	}

	credential := common.NewCredential(secretID, secretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"

	client, err := cvm.NewClient(credential, region, cpf)
	if err != nil {
		return nil, err
	}

	return &TencentAdapter{
		BaseAdapter: NewBaseAdapter(secretID, secretKey),
		client:      client,
		region:      region,
	}, nil
}

// 腾讯云主机状态映射
var tencentStatusMap = map[string]string{
	"PENDING":       "pending",
	"LAUNCH_FAILED": "error",
	"RUNNING":       "running",
	"STOPPED":       "stopped",
	"STARTING":      "starting",
	"STOPPING":      "stopping",
	"REBOOTING":     "rebooting",
	"SHUTDOWN":      "stopped",
	"TERMINATING":   "stopping",
}

// convertStatus 转换腾讯云主机状态为标准状态
func (a *TencentAdapter) convertStatus(status string) string {
	if standardStatus, ok := tencentStatusMap[status]; ok {
		return standardStatus
	}
	return "unknown"
}

// ListInstances 获取腾讯云主机列表
func (a *TencentAdapter) ListInstances() ([]cmdb.Host, error) {
	var hosts []cmdb.Host
	var offset int64 = 0
	var limit int64 = 100
	var total int64 = 0

	request := cvm.NewDescribeInstancesRequest()
	request.Limit = common.Int64Ptr(limit)

	for {
		request.Offset = common.Int64Ptr(offset)
		response, err := a.client.DescribeInstances(request)
		if err != nil {
			return nil, fmt.Errorf("获取腾讯云主机列表失败: %w", err)
		}

		for _, instance := range response.Response.InstanceSet {
			// 构建公网IP和私网IP
			var publicIPs, privateIPs []string

			for _, ip := range instance.PublicIpAddresses {
				publicIPs = append(publicIPs, *ip)
			}

			for _, ip := range instance.PrivateIpAddresses {
				privateIPs = append(privateIPs, *ip)
			}

			publicIPsJSON, _ := json.Marshal(publicIPs)
			privateIPsJSON, _ := json.Marshal(privateIPs)

			// 构建配置信息
			config := map[string]interface{}{
				"cpu_cores":     *instance.CPU,
				"memory_size":   *instance.Memory,
				"instance_type": *instance.InstanceType,
				"zone_id":       *instance.Placement.Zone,
				"vpc_id":        *instance.VirtualPrivateCloud.VpcId,
				"subnet_id":     *instance.VirtualPrivateCloud.SubnetId,
				"charge_type":   *instance.InstanceChargeType,
			}
			configJSON, _ := json.Marshal(config)

			// 解析过期时间
			var expiredAt *time.Time
			if instance.ExpiredTime != nil {
				t, err := time.Parse("2006-01-02T15:04:05Z", *instance.ExpiredTime)
				if err == nil {
					expiredAt = &t
				}
			}

			host := cmdb.Host{
				InstanceID:    *instance.InstanceId,
				Name:          *instance.InstanceName,
				Region:        a.region,
				PublicIP:      datatypes.JSON(publicIPsJSON),
				PrivateIP:     datatypes.JSON(privateIPsJSON),
				Configuration: datatypes.JSON(configJSON),
				OS:            *instance.OsName,
				Status:        a.convertStatus(*instance.InstanceState),
				ProviderType:  "tencent",
				ResourceType:  cmdb.ResourceTypeCVM,
				ExpiredAt:     expiredAt,
			}
			hosts = append(hosts, host)
		}

		total = *response.Response.TotalCount
		offset += limit
		if offset >= total {
			break
		}
	}

	return hosts, nil
}

// GetInstanceStatus 获取腾讯云主机状态
func (a *TencentAdapter) GetInstanceStatus(instanceID string) (string, error) {
	request := cvm.NewDescribeInstancesRequest()
	request.InstanceIds = []*string{common.StringPtr(instanceID)}

	response, err := a.client.DescribeInstances(request)
	if err != nil {
		return "", fmt.Errorf("获取腾讯云主机状态失败: %w", err)
	}

	if len(response.Response.InstanceSet) == 0 {
		return "", fmt.Errorf("腾讯云实例不存在: %s", instanceID)
	}

	instance := response.Response.InstanceSet[0]
	return a.convertStatus(*instance.InstanceState), nil
}

// GetInstanceInfo 获取腾讯云主机详细信息
func (a *TencentAdapter) GetInstanceInfo(instanceID string) (*cmdb.Host, error) {
	request := cvm.NewDescribeInstancesRequest()
	request.InstanceIds = []*string{common.StringPtr(instanceID)}

	response, err := a.client.DescribeInstances(request)
	if err != nil {
		return nil, fmt.Errorf("获取腾讯云主机信息失败: %w", err)
	}

	if len(response.Response.InstanceSet) == 0 {
		return nil, fmt.Errorf("腾讯云实例不存在: %s", instanceID)
	}

	instance := response.Response.InstanceSet[0]

	// 构建公网IP和私网IP
	var publicIPs, privateIPs []string

	for _, ip := range instance.PublicIpAddresses {
		publicIPs = append(publicIPs, *ip)
	}

	for _, ip := range instance.PrivateIpAddresses {
		privateIPs = append(privateIPs, *ip)
	}

	publicIPsJSON, _ := json.Marshal(publicIPs)
	privateIPsJSON, _ := json.Marshal(privateIPs)

	// 构建配置信息
	config := map[string]interface{}{
		"cpu_cores":     *instance.CPU,
		"memory_size":   *instance.Memory,
		"instance_type": *instance.InstanceType,
		"zone_id":       *instance.Placement.Zone,
		"vpc_id":        *instance.VirtualPrivateCloud.VpcId,
		"subnet_id":     *instance.VirtualPrivateCloud.SubnetId,
		"charge_type":   *instance.InstanceChargeType,
	}
	configJSON, _ := json.Marshal(config)

	// 解析过期时间
	var expiredAt *time.Time
	if instance.ExpiredTime != nil {
		t, err := time.Parse("2006-01-02T15:04:05Z", *instance.ExpiredTime)
		if err == nil {
			expiredAt = &t
		}
	}

	return &cmdb.Host{
		InstanceID:    *instance.InstanceId,
		Name:          *instance.InstanceName,
		Region:        a.region,
		PublicIP:      datatypes.JSON(publicIPsJSON),
		PrivateIP:     datatypes.JSON(privateIPsJSON),
		Configuration: datatypes.JSON(configJSON),
		OS:            *instance.OsName,
		Status:        a.convertStatus(*instance.InstanceState),
		ProviderType:  "tencent",
		ResourceType:  cmdb.ResourceTypeCVM,
		ExpiredAt:     expiredAt,
	}, nil
}
