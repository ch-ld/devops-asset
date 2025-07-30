package adapter

import (
	"api-server/internal/model/cmdb"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gorm.io/datatypes"
)

// AWSAdapter AWS适配器
type AWSAdapter struct {
	*BaseAdapter
	client *ec2.EC2
}

// NewAWSAdapter 创建AWS适配器实例
func NewAWSAdapter(accessKey, secretKey, region string) (*AWSAdapter, error) {
	// 如果未指定区域，使用默认区域
	if region == "" {
		region = "us-east-1"
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, err
	}

	client := ec2.New(sess)
	return &AWSAdapter{
		BaseAdapter: NewBaseAdapter(accessKey, secretKey),
		client:      client,
	}, nil
}

// AWS主机状态映射
var awsStatusMap = map[string]string{
	"pending":       "pending",
	"running":       "running",
	"shutting-down": "stopping",
	"terminated":    "deleted",
	"stopping":      "stopping",
	"stopped":       "stopped",
}

// AWS实例类型配置映射
type InstanceTypeConfig struct {
	CPUCores   int     `json:"cpu_cores"`
	MemorySize float64 `json:"memory_size"`
	DiskSize   int     `json:"disk_size"`
	Network    string  `json:"network"`
}

var awsInstanceTypeMap = map[string]InstanceTypeConfig{
	// t2 系列 - 通用型
	"t2.nano":    {CPUCores: 1, MemorySize: 0.5, DiskSize: 8, Network: "低到中等"},
	"t2.micro":   {CPUCores: 1, MemorySize: 1, DiskSize: 8, Network: "低到中等"},
	"t2.small":   {CPUCores: 1, MemorySize: 2, DiskSize: 20, Network: "低到中等"},
	"t2.medium":  {CPUCores: 2, MemorySize: 4, DiskSize: 30, Network: "低到中等"},
	"t2.large":   {CPUCores: 2, MemorySize: 8, DiskSize: 30, Network: "低到中等"},
	"t2.xlarge":  {CPUCores: 4, MemorySize: 16, DiskSize: 40, Network: "中等"},
	"t2.2xlarge": {CPUCores: 8, MemorySize: 32, DiskSize: 40, Network: "中等"},

	// t3 系列 - 通用型
	"t3.nano":    {CPUCores: 2, MemorySize: 0.5, DiskSize: 8, Network: "最高5Gbps"},
	"t3.micro":   {CPUCores: 2, MemorySize: 1, DiskSize: 8, Network: "最高5Gbps"},
	"t3.small":   {CPUCores: 2, MemorySize: 2, DiskSize: 20, Network: "最高5Gbps"},
	"t3.medium":  {CPUCores: 2, MemorySize: 4, DiskSize: 30, Network: "最高5Gbps"},
	"t3.large":   {CPUCores: 2, MemorySize: 8, DiskSize: 30, Network: "最高5Gbps"},
	"t3.xlarge":  {CPUCores: 4, MemorySize: 16, DiskSize: 40, Network: "最高5Gbps"},
	"t3.2xlarge": {CPUCores: 8, MemorySize: 32, DiskSize: 40, Network: "最高5Gbps"},

	// m5 系列 - 通用型
	"m5.large":    {CPUCores: 2, MemorySize: 8, DiskSize: 50, Network: "最高10Gbps"},
	"m5.xlarge":   {CPUCores: 4, MemorySize: 16, DiskSize: 50, Network: "最高10Gbps"},
	"m5.2xlarge":  {CPUCores: 8, MemorySize: 32, DiskSize: 100, Network: "最高10Gbps"},
	"m5.4xlarge":  {CPUCores: 16, MemorySize: 64, DiskSize: 200, Network: "最高10Gbps"},
	"m5.8xlarge":  {CPUCores: 32, MemorySize: 128, DiskSize: 300, Network: "10Gbps"},
	"m5.12xlarge": {CPUCores: 48, MemorySize: 192, DiskSize: 450, Network: "12Gbps"},
	"m5.16xlarge": {CPUCores: 64, MemorySize: 256, DiskSize: 600, Network: "20Gbps"},
	"m5.24xlarge": {CPUCores: 96, MemorySize: 384, DiskSize: 900, Network: "25Gbps"},

	// c5 系列 - 计算优化型
	"c5.large":    {CPUCores: 2, MemorySize: 4, DiskSize: 50, Network: "最高10Gbps"},
	"c5.xlarge":   {CPUCores: 4, MemorySize: 8, DiskSize: 50, Network: "最高10Gbps"},
	"c5.2xlarge":  {CPUCores: 8, MemorySize: 16, DiskSize: 100, Network: "最高10Gbps"},
	"c5.4xlarge":  {CPUCores: 16, MemorySize: 32, DiskSize: 200, Network: "最高10Gbps"},
	"c5.9xlarge":  {CPUCores: 36, MemorySize: 72, DiskSize: 450, Network: "10Gbps"},
	"c5.12xlarge": {CPUCores: 48, MemorySize: 96, DiskSize: 600, Network: "12Gbps"},
	"c5.18xlarge": {CPUCores: 72, MemorySize: 144, DiskSize: 900, Network: "25Gbps"},
	"c5.24xlarge": {CPUCores: 96, MemorySize: 192, DiskSize: 1200, Network: "25Gbps"},

	// r5 系列 - 内存优化型
	"r5.large":    {CPUCores: 2, MemorySize: 16, DiskSize: 50, Network: "最高10Gbps"},
	"r5.xlarge":   {CPUCores: 4, MemorySize: 32, DiskSize: 50, Network: "最高10Gbps"},
	"r5.2xlarge":  {CPUCores: 8, MemorySize: 64, DiskSize: 100, Network: "最高10Gbps"},
	"r5.4xlarge":  {CPUCores: 16, MemorySize: 128, DiskSize: 200, Network: "最高10Gbps"},
	"r5.8xlarge":  {CPUCores: 32, MemorySize: 256, DiskSize: 300, Network: "10Gbps"},
	"r5.12xlarge": {CPUCores: 48, MemorySize: 384, DiskSize: 450, Network: "12Gbps"},
	"r5.16xlarge": {CPUCores: 64, MemorySize: 512, DiskSize: 600, Network: "20Gbps"},
	"r5.24xlarge": {CPUCores: 96, MemorySize: 768, DiskSize: 900, Network: "25Gbps"},
}

// convertStatus 转换AWS主机状态为标准状态
func (a *AWSAdapter) convertStatus(status string) string {
	if standardStatus, ok := awsStatusMap[status]; ok {
		return standardStatus
	}
	return "unknown"
}

// ListInstances 获取AWS主机列表
func (a *AWSAdapter) ListInstances() ([]cmdb.Host, error) {
	var hosts []cmdb.Host
	var nextToken *string

	for {
		input := &ec2.DescribeInstancesInput{
			NextToken: nextToken,
		}

		resp, err := a.client.DescribeInstances(input)
		if err != nil {
			return nil, fmt.Errorf("获取AWS主机列表失败: %w", err)
		}

		for _, reservation := range resp.Reservations {
			for _, instance := range reservation.Instances {
				// 跳过已终止的实例
				if *instance.State.Name == "terminated" {
					continue
				}

				// 构建公网IP和私网IP
				var publicIPs, privateIPs []string
				if instance.PublicIpAddress != nil {
					publicIPs = append(publicIPs, *instance.PublicIpAddress)
				}
				if instance.PrivateIpAddress != nil {
					privateIPs = append(privateIPs, *instance.PrivateIpAddress)
				}

				publicIPsJSON, _ := json.Marshal(publicIPs)
				privateIPsJSON, _ := json.Marshal(privateIPs)

				// 构建配置信息
				config := map[string]interface{}{
					"instance_type":     *instance.InstanceType,
					"availability_zone": *instance.Placement.AvailabilityZone,
				}

				// 根据实例类型添加CPU、内存、磁盘信息
				if typeConfig, exists := awsInstanceTypeMap[*instance.InstanceType]; exists {
					config["cpu_cores"] = typeConfig.CPUCores
					config["memory_size"] = typeConfig.MemorySize
					config["disk_size"] = typeConfig.DiskSize
					config["network_performance"] = typeConfig.Network
				} else {
					// 如果实例类型不在映射表中，尝试从API获取CPU信息
					if instance.CpuOptions != nil {
						config["cpu_cores"] = *instance.CpuOptions.CoreCount * *instance.CpuOptions.ThreadsPerCore
					}
				}

				// 添加VPC信息
				if instance.VpcId != nil {
					config["vpc_id"] = *instance.VpcId
				}

				// 添加子网信息
				if instance.SubnetId != nil {
					config["subnet_id"] = *instance.SubnetId
				}

				// 添加安全组信息
				if len(instance.SecurityGroups) > 0 {
					var securityGroups []string
					for _, sg := range instance.SecurityGroups {
						if sg.GroupId != nil {
							securityGroups = append(securityGroups, *sg.GroupId)
						}
					}
					config["security_groups"] = securityGroups
				}

				configJSON, _ := json.Marshal(config)

				// 解析实例名称和标签信息
				name := fmt.Sprintf("aws-instance-%s", *instance.InstanceId)
				var tags []map[string]string
				osInfo := "Amazon Linux"

				for _, tag := range instance.Tags {
					if tag.Key != nil && tag.Value != nil {
						tagMap := map[string]string{
							"key":   *tag.Key,
							"value": *tag.Value,
						}
						tags = append(tags, tagMap)

						// 特殊处理Name和OS标签
						switch *tag.Key {
						case "Name":
							name = *tag.Value
						case "OS", "OperatingSystem":
							osInfo = *tag.Value
						}
					}
				}

				// 序列化标签信息
				tagsJSON, _ := json.Marshal(tags)

				// 创建主机信息
				host := cmdb.Host{
					InstanceID:    *instance.InstanceId,
					Name:          name,
					Region:        *instance.Placement.AvailabilityZone,
					PublicIP:      datatypes.JSON(publicIPsJSON),
					PrivateIP:     datatypes.JSON(privateIPsJSON),
					Configuration: datatypes.JSON(configJSON),
					Tags:          datatypes.JSON(tagsJSON),
					OS:            osInfo,
					Status:        a.convertStatus(*instance.State.Name),
					ProviderType:  "aws",
					ResourceType:  cmdb.ResourceTypeEC2,
				}

				// 检查是否有过期时间标签
				for _, tag := range instance.Tags {
					if *tag.Key == "ExpirationDate" {
						if t, err := time.Parse("2006-01-02", *tag.Value); err == nil {
							host.ExpiredAt = &t
						}
					}
				}

				hosts = append(hosts, host)
			}
		}

		// 检查是否有下一页
		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	return hosts, nil
}

// GetInstanceStatus 获取AWS主机状态
func (a *AWSAdapter) GetInstanceStatus(instanceID string) (string, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}

	resp, err := a.client.DescribeInstances(input)
	if err != nil {
		return "", fmt.Errorf("获取AWS主机状态失败: %w", err)
	}

	if len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
		return "", fmt.Errorf("AWS实例不存在: %s", instanceID)
	}

	instance := resp.Reservations[0].Instances[0]
	return a.convertStatus(*instance.State.Name), nil
}

// GetInstanceInfo 获取AWS主机详细信息
func (a *AWSAdapter) GetInstanceInfo(instanceID string) (*cmdb.Host, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}

	resp, err := a.client.DescribeInstances(input)
	if err != nil {
		return nil, fmt.Errorf("获取AWS主机信息失败: %w", err)
	}

	if len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("AWS实例不存在: %s", instanceID)
	}

	instance := resp.Reservations[0].Instances[0]

	// 构建公网IP和私网IP
	var publicIPs, privateIPs []string
	if instance.PublicIpAddress != nil {
		publicIPs = append(publicIPs, *instance.PublicIpAddress)
	}
	if instance.PrivateIpAddress != nil {
		privateIPs = append(privateIPs, *instance.PrivateIpAddress)
	}

	publicIPsJSON, _ := json.Marshal(publicIPs)
	privateIPsJSON, _ := json.Marshal(privateIPs)

	// 构建配置信息
	config := map[string]interface{}{
		"instance_type":     *instance.InstanceType,
		"availability_zone": *instance.Placement.AvailabilityZone,
	}

	// 添加vCPU和内存信息
	if instance.CpuOptions != nil {
		config["cpu_cores"] = *instance.CpuOptions.CoreCount * *instance.CpuOptions.ThreadsPerCore
	}

	// 添加VPC信息
	if instance.VpcId != nil {
		config["vpc_id"] = *instance.VpcId
	}

	configJSON, _ := json.Marshal(config)

	// 解析实例名称
	name := fmt.Sprintf("aws-instance-%s", *instance.InstanceId)
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			name = *tag.Value
			break
		}
	}

	// 解析操作系统信息
	osInfo := "Amazon Linux"
	for _, tag := range instance.Tags {
		if *tag.Key == "OS" {
			osInfo = *tag.Value
			break
		}
	}

	// 创建主机信息
	host := &cmdb.Host{
		InstanceID:    *instance.InstanceId,
		Name:          name,
		Region:        *instance.Placement.AvailabilityZone,
		PublicIP:      datatypes.JSON(publicIPsJSON),
		PrivateIP:     datatypes.JSON(privateIPsJSON),
		Configuration: datatypes.JSON(configJSON),
		OS:            osInfo,
		Status:        a.convertStatus(*instance.State.Name),
		ProviderType:  "aws",
		ResourceType:  cmdb.ResourceTypeEC2,
	}

	// 检查是否有过期时间标签
	for _, tag := range instance.Tags {
		if *tag.Key == "ExpirationDate" {
			if t, err := time.Parse("2006-01-02", *tag.Value); err == nil {
				host.ExpiredAt = &t
			}
		}
	}

	return host, nil
}

// GetRegions 获取AWS支持的区域列表
func (a *AWSAdapter) GetRegions() ([]Region, error) {
	input := &ec2.DescribeRegionsInput{}

	response, err := a.client.DescribeRegions(input)
	if err != nil {
		return nil, fmt.Errorf("获取AWS区域列表失败: %w", err)
	}

	var regions []Region
	for _, region := range response.Regions {
		regions = append(regions, Region{
			ID:   *region.RegionName,
			Name: *region.RegionName, // AWS区域名称通常就是ID
		})
	}

	return regions, nil
}
