package cmdb

import (
	"api-server/internal/config"
	model "api-server/internal/model/cmdb"
	repo "api-server/internal/repository/cmdb"
	"api-server/internal/service/cmdb/adapter"
	"api-server/pkg/crypto/encryption"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"encoding/json"

	gosshtool "golang.org/x/crypto/ssh"
	"gorm.io/datatypes"
)

// HostService 主机业务逻辑服务
// 提供主机的增删改查、加密解密等核心业务逻辑
type HostService struct {
	hostRepo     *repo.HostRepository
	providerRepo *repo.ProviderRepository
	groupRepo    *repo.HostGroupRepository
}

// NewHostService 创建主机业务服务实例
func NewHostService(hostRepo *repo.HostRepository, providerRepo *repo.ProviderRepository, groupRepo *repo.HostGroupRepository) *HostService {
	return &HostService{
		hostRepo:     hostRepo,
		providerRepo: providerRepo,
		groupRepo:    groupRepo,
	}
}

// CreateHost 创建主机，自动加密敏感字段
func (s *HostService) CreateHost(host *model.Host) error {
	// 多密钥加密主机密码
	if host.Password != "" {
		keys := getAllAESKeys()
		enc, err := encryption.EncryptAESWithKeys([]byte(host.Password), keys)
		if err != nil {
			return fmt.Errorf("主机密码加密失败: %w", err)
		}
		host.Password = enc
	}

	// 加密私钥
	if host.PrivateKey != "" {
		keys := getAllAESKeys()
		enc, err := encryption.EncryptAESWithKeys([]byte(host.PrivateKey), keys)
		if err != nil {
			return fmt.Errorf("主机私钥加密失败: %w", err)
		}
		host.PrivateKey = enc
	}

	// 设置默认端口
	if host.Port == 0 {
		host.Port = 22
	}

	return s.hostRepo.Create(host)
}

// UpdateHost 更新主机，自动加密敏感字段，未传递密码时保留原密码
func (s *HostService) UpdateHost(host *model.Host) error {
	// 获取现有主机信息，用于保留未更新的敏感字段
	existingHost, err := s.hostRepo.FindByID(host.ID)
	if err != nil {
		return err
	}

	// 处理密码字段
	if host.Password != "" {
		key := []byte(config.GlobalConfig.App.AesKey)
		encryptedPassword, err := encryption.EncryptAES([]byte(host.Password), key)
		if err != nil {
			return errors.New("failed to encrypt password")
		}
		host.Password = encryptedPassword
	} else {
		// 保留原密码
		host.Password = existingHost.Password
	}

	// 处理私钥字段
	if host.PrivateKey != "" {
		key := []byte(config.GlobalConfig.App.AesKey)
		encryptedPrivateKey, err := encryption.EncryptAES([]byte(host.PrivateKey), key)
		if err != nil {
			return errors.New("failed to encrypt private key")
		}
		host.PrivateKey = encryptedPrivateKey
	} else {
		// 保留原私钥
		host.PrivateKey = existingHost.PrivateKey
	}

	// 处理端口字段
	if host.Port == 0 {
		if existingHost.Port != 0 {
			host.Port = existingHost.Port
		} else {
			host.Port = 22 // 默认端口
		}
	}

	return s.hostRepo.Update(host)
}

// ListHosts 查询主机列表，支持分页和过滤
func (s *HostService) ListHosts(params map[string]interface{}) ([]model.Host, error) {
	page := 1
	pageSize := 10
	if p, ok := params["page"].(int); ok {
		page = p
	}
	if ps, ok := params["page_size"].(int); ok {
		pageSize = ps
	}

	// 构建查询条件
	conditions := make(map[string]interface{})
	if keyword, ok := params["keyword"].(string); ok && keyword != "" {
		conditions["keyword"] = keyword
	}
	if name, ok := params["name"].(string); ok && name != "" {
		conditions["name"] = name
	}
	if ip, ok := params["ip"].(string); ok && ip != "" {
		conditions["ip"] = ip
	}
	if status, ok := params["status"].(string); ok && status != "" {
		conditions["status"] = status
	}
	if groupID, ok := params["group_id"].(uint); ok {
		conditions["group_id"] = groupID
	}
	if region, ok := params["region"].(string); ok && region != "" {
		conditions["region"] = region
	}

	hosts, err := s.hostRepo.FindByConditions(conditions, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 解密敏感信息
	for i := range hosts {
		if hosts[i].Password != "" {
			keys := getAllAESKeys()
			plain, err := encryption.DecryptAESWithKeys(hosts[i].Password, keys)
			if err == nil {
				hosts[i].Password = plain
			}
		}
	}

	return hosts, nil
}

// ListAllHosts 查询所有主机列表，不分页，用于导出功能
func (s *HostService) ListAllHosts(params map[string]interface{}) ([]model.Host, error) {
	// 构建查询条件
	conditions := make(map[string]interface{})
	if keyword, ok := params["keyword"].(string); ok && keyword != "" {
		conditions["keyword"] = keyword
	}
	if name, ok := params["name"].(string); ok && name != "" {
		conditions["name"] = name
	}
	if ip, ok := params["ip"].(string); ok && ip != "" {
		conditions["ip"] = ip
	}
	if status, ok := params["status"].(string); ok && status != "" {
		conditions["status"] = status
	}
	if groupID, ok := params["group_id"].(uint); ok {
		conditions["group_id"] = groupID
	}
	if region, ok := params["region"].(string); ok && region != "" {
		conditions["region"] = region
	}
	if providerType, ok := params["provider_type"].(string); ok && providerType != "" {
		conditions["provider_type"] = providerType
	}

	// 使用不分页的查询方法
	hosts, err := s.hostRepo.FindAllByConditions(conditions)
	if err != nil {
		return nil, err
	}

	// 解密敏感信息
	for i := range hosts {
		if hosts[i].Password != "" {
			keys := getAllAESKeys()
			plain, err := encryption.DecryptAESWithKeys(hosts[i].Password, keys)
			if err == nil {
				hosts[i].Password = plain
			}
		}
	}

	return hosts, nil
}

// CountHosts 统计主机总数，支持过滤条件
func (s *HostService) CountHosts(params map[string]interface{}) (int64, error) {
	// 构建查询条件
	conditions := make(map[string]interface{})
	if keyword, ok := params["keyword"].(string); ok && keyword != "" {
		conditions["keyword"] = keyword
	}
	if name, ok := params["name"].(string); ok && name != "" {
		conditions["name"] = name
	}
	if ip, ok := params["ip"].(string); ok && ip != "" {
		conditions["ip"] = ip
	}
	if status, ok := params["status"].(string); ok && status != "" {
		conditions["status"] = status
	}
	if groupID, ok := params["group_id"].(uint); ok {
		conditions["group_id"] = groupID
	}
	if region, ok := params["region"].(string); ok && region != "" {
		conditions["region"] = region
	}

	return s.hostRepo.CountByConditions(conditions)
}

// GetHost 查询单个主机，自动解密敏感字段
func (s *HostService) GetHost(id uint) (*model.Host, error) {
	host, err := s.hostRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	// 多密钥解密主机密码
	if host.Password != "" {
		keys := getAllAESKeys()
		plain, err := encryption.DecryptAESWithKeys(host.Password, keys)
		if err == nil {
			host.Password = plain
		}
	}
	return host, nil
}

// DeleteHost 删除主机
func (s *HostService) DeleteHost(id uint) error {
	return s.hostRepo.DeleteByID(id)
}

// SyncHostsFromCloud 从云端同步主机资源到本地数据库
func (s *HostService) SyncHostsFromCloud(providerID uint) error {
	// 获取云账号信息
	provider, err := s.providerRepo.FindByID(providerID)
	if err != nil {
		return fmt.Errorf("获取云账号信息失败: %w", err)
	}

	// 获取适配器
	adapter, err := s.getCloudAdapter(provider)
	if err != nil {
		return fmt.Errorf("初始化云适配器失败: %w", err)
	}

	// 从云端获取主机列表
	cloudHosts, err := adapter.ListInstances()
	if err != nil {
		return fmt.Errorf("获取云主机列表失败: %w", err)
	}

	// 获取本地已存在的主机
	localHosts, err := s.hostRepo.FindByProviderID(providerID)
	if err != nil {
		return fmt.Errorf("获取本地主机列表失败: %w", err)
	}

	// 构建本地主机映射，用于快速查找
	localHostMap := make(map[string]*model.Host)
	for i, host := range localHosts {
		localHostMap[host.InstanceID] = &localHosts[i]
	}

	// 处理需要新增和更新的主机
	var (
		hostsToCreate []model.Host
		hostsToUpdate []model.Host
	)

	for _, cloudHost := range cloudHosts {
		if localHost, exists := localHostMap[cloudHost.InstanceID]; exists {
			// 更新已存在的主机
			localHost.Name = cloudHost.Name
			localHost.Status = cloudHost.Status
			localHost.PublicIP = cloudHost.PublicIP
			localHost.PrivateIP = cloudHost.PrivateIP
			localHost.Configuration = cloudHost.Configuration
			localHost.OS = cloudHost.OS
			localHost.Region = cloudHost.Region
			localHost.ExpiredAt = cloudHost.ExpiredAt
			hostsToUpdate = append(hostsToUpdate, *localHost)
			delete(localHostMap, cloudHost.InstanceID) // 从映射中删除已处理的主机
		} else {
			// 添加新主机
			cloudHost.ProviderID = &providerID
			hostsToCreate = append(hostsToCreate, cloudHost)
		}
	}

	// 批量创建新主机
	if len(hostsToCreate) > 0 {
		if err := s.hostRepo.BatchCreate(hostsToCreate); err != nil {
			return fmt.Errorf("批量创建主机失败: %w", err)
		}
	}

	// 批量更新主机
	if len(hostsToUpdate) > 0 {
		if err := s.BatchUpdateHosts(hostsToUpdate); err != nil {
			return fmt.Errorf("批量更新主机失败: %w", err)
		}
	}

	// 处理已删除的主机（在本地存在但云端不存在的主机）
	var deletedHostIDs []uint
	for _, host := range localHostMap {
		deletedHostIDs = append(deletedHostIDs, host.ID)
	}
	if len(deletedHostIDs) > 0 {
		if err := s.hostRepo.BatchDelete(deletedHostIDs); err != nil {
			return fmt.Errorf("删除已不存在的主机失败: %w", err)
		}
	}

	return nil
}

// SyncAllProviderHosts 同步所有云账号的主机
func (s *HostService) SyncAllProviderHosts() error {
	// 获取所有启用的云账号
	providers, err := s.providerRepo.FindEnabled()
	if err != nil {
		return fmt.Errorf("获取云账号列表失败: %w", err)
	}

	// 遍历同步每个云账号的主机
	for _, provider := range providers {
		if err := s.SyncHostsFromCloud(provider.ID); err != nil {
			// 记录错误但继续执行其他云账号的同步
			log.Printf("同步云账号[%s]的主机失败: %s", provider.Name, err.Error())
			continue
		}
	}

	return nil
}

// SyncHostStatus 同步单个主机状态
func (s *HostService) SyncHostStatus(hostID uint) error {
	// 获取主机信息
	host, err := s.hostRepo.FindByID(hostID)
	if err != nil {
		return fmt.Errorf("获取主机信息失败: %w", err)
	}

	// 获取云账号信息
	if host.ProviderID == nil {
		return fmt.Errorf("主机未关联云账号")
	}
	provider, err := s.providerRepo.FindByID(*host.ProviderID)
	if err != nil {
		return fmt.Errorf("获取云账号信息失败: %w", err)
	}

	// 获取适配器
	adapter, err := s.getCloudAdapter(provider)
	if err != nil {
		return fmt.Errorf("初始化云适配器失败: %w", err)
	}

	// 获取最新状态
	status, err := adapter.GetInstanceStatus(host.InstanceID)
	if err != nil {
		return fmt.Errorf("获取主机状态失败: %w", err)
	}

	// 更新状态
	if host.Status != status {
		host.Status = status
		if err := s.hostRepo.Update(host); err != nil {
			return fmt.Errorf("更新主机状态失败: %w", err)
		}
	}

	return nil
}

// AlertHostStatus 检查主机状态并触发异常/到期告警
func (s *HostService) AlertHostStatus() error {
	// TODO: Implement logic to check host status and trigger alerts.
	// This might involve querying hosts from the repository,
	// checking their expiration dates, and sending alerts if necessary.
	return nil
}

// AlertHosts 查询异常、到期、即将到期主机
func (s *HostService) AlertHosts(days int) ([]model.Host, error) {
	var alerts []model.Host
	hosts, err := s.hostRepo.FindAll()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	soon := now.Add(time.Duration(days) * 24 * time.Hour)
	for _, h := range hosts {
		if h.Status == "error" || h.Status == "abnormal" {
			alerts = append(alerts, *h)
			continue
		}
		if h.ExpiredAt != nil && (h.ExpiredAt.Before(now) || h.ExpiredAt.Before(soon)) {
			alerts = append(alerts, *h)
		}
	}
	return alerts, nil
}

// BatchCreateHosts 批量创建主机，自动加密敏感字段
func (s *HostService) BatchCreateHosts(hosts []model.Host) error {
	keys := getAllAESKeys()
	for i := range hosts {
		if hosts[i].Password != "" {
			enc, err := encryption.EncryptAESWithKeys([]byte(hosts[i].Password), keys)
			if err != nil {
				return fmt.Errorf("主机密码加密失败: %w", err)
			}
			hosts[i].Password = enc
		}
	}
	return s.hostRepo.BatchCreate(hosts)
}

// BatchDeleteHosts 批量删除主机
func (s *HostService) BatchDeleteHosts(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.hostRepo.BatchDelete(ids)
}

// BatchUpdateHosts 批量更新主机属性
func (s *HostService) BatchUpdateHosts(hosts []model.Host) error {
	for i := range hosts {
		if err := s.UpdateHost(&hosts[i]); err != nil {
			return err
		}
	}
	return nil
}

// BatchAssignRequest 批量分配请求
type BatchAssignRequest struct {
	IDs   []uint
	Group string
	Tags  []string
	Owner string
}

// BatchAssignHosts 批量分配主机
func (s *HostService) BatchAssignHosts(req BatchAssignRequest) error {
	if len(req.IDs) == 0 {
		return nil
	}
	for _, id := range req.IDs {
		host, err := s.hostRepo.FindByID(id)
		if err != nil {
			return err
		}
		if req.Group != "" {
			// 根据组名查找组ID，这里简化处理，实际应该查找组ID
			// TODO: 实现根据组名查找组ID的逻辑
			// 暂时跳过组分配逻辑
		}
		if len(req.Tags) > 0 {
			b, _ := json.Marshal(req.Tags)
			host.Tags = datatypes.JSON(b)
		}
		// 移除host.Owner相关逻辑（如需请补充模型）
		if err := s.hostRepo.Update(host); err != nil {
			return err
		}
	}
	return nil
}

// BatchLifecycleHosts 批量设置主机生命周期
func (s *HostService) BatchLifecycleHosts(req model.BatchLifecycleRequest) error {
	if len(req.IDs) == 0 {
		return nil
	}
	for _, id := range req.IDs {
		host, err := s.hostRepo.FindByID(id)
		if err != nil {
			return err
		}
		if req.ExpiredAt != nil {
			host.ExpiredAt = req.ExpiredAt
		}
		if req.Status != "" {
			host.Status = req.Status
		}
		if req.Recycle {
			host.Status = "recycled"
		}
		if err := s.hostRepo.Update(host); err != nil {
			return err
		}
	}
	return nil
}

// BatchSetCustomFields 批量设置自定义字段
func (s *HostService) BatchSetCustomFields(req model.BatchSetCustomFieldsRequest) error {
	if len(req.IDs) == 0 || len(req.ExtraFields) == 0 {
		return nil
	}
	for _, id := range req.IDs {
		host, err := s.hostRepo.FindByID(id)
		if err != nil {
			return err
		}
		// 合并原有ExtraFields与新字段
		var extra map[string]interface{}
		if host.ExtraFields != nil && len(host.ExtraFields) > 0 {
			_ = json.Unmarshal(host.ExtraFields, &extra)
		}
		if extra == nil {
			extra = map[string]interface{}{}
		}
		for k, v := range req.ExtraFields {
			extra[k] = v
		}
		b, _ := json.Marshal(extra)
		host.ExtraFields = datatypes.JSON(b)
		if err := s.hostRepo.Update(host); err != nil {
			return err
		}
	}
	return nil
}

// BatchChangeStatus 批量变更主机状态
func (s *HostService) BatchChangeStatus(req model.BatchChangeStatusRequest) error {
	if len(req.IDs) == 0 || req.Status == "" {
		return nil
	}
	for _, id := range req.IDs {
		host, err := s.hostRepo.FindByID(id)
		if err != nil {
			return err
		}
		host.Status = req.Status
		if err := s.hostRepo.Update(host); err != nil {
			return err
		}
	}
	return nil
}

// SSHResult SSH命令执行结果
type SSHResult struct {
	HostID  uint
	Success bool
	Output  string
	Error   string
}

// BatchSSH 批量执行SSH命令
func (s *HostService) BatchSSH(req model.BatchSSHRequest) []SSHResult {
	var results []SSHResult
	for _, id := range req.IDs {
		host, err := s.hostRepo.FindByID(id)
		if err != nil {
			results = append(results, SSHResult{HostID: id, Success: false, Error: err.Error()})
			continue
		}
		output, err := execSSHCommand(host, req.Cmd, req.Timeout)
		if err != nil {
			results = append(results, SSHResult{HostID: id, Success: false, Error: err.Error(), Output: output})
		} else {
			results = append(results, SSHResult{HostID: id, Success: true, Output: output})
		}
	}
	return results
}

// execSSHCommand 执行单台主机SSH命令
func execSSHCommand(host *model.Host, cmd string, timeout int) (string, error) {
	if len(host.PublicIP) < 3 {
		return "", fmt.Errorf("主机公网IP未配置")
	}

	// 解析IP地址
	targetIP := string(host.PublicIP[1 : len(host.PublicIP)-1])

	if host.Username == "" || host.Password == "" {
		return "", fmt.Errorf("主机SSH凭证未配置")
	}

	// 解密密码
	key := []byte(config.GlobalConfig.App.AesKey)
	decryptedPassword, err := encryption.DecryptAES(host.Password, key)
	if err != nil {
		return "", fmt.Errorf("解密密码失败: %w", err)
	}

	// 创建SSH连接
	sshClient, err := createSSHConnection(host.Username, decryptedPassword, targetIP, 22)
	if err != nil {
		return "", fmt.Errorf("SSH连接失败: %w", err)
	}
	defer sshClient.Close()

	// 创建会话
	session, err := sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建SSH会话失败: %w", err)
	}
	defer session.Close()

	// 执行命令
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(output), fmt.Errorf("命令执行失败: %w", err)
	}

	return string(output), nil
}

// uploadFileToHost 上传文件到单台主机
func uploadFileToHost(host *model.Host, remotePath string, file io.Reader) error {
	if len(host.PublicIP) < 3 {
		return fmt.Errorf("主机公网IP未配置")
	}

	// 解析IP地址
	targetIP := string(host.PublicIP[1 : len(host.PublicIP)-1])

	if host.Username == "" || host.Password == "" {
		return fmt.Errorf("主机SSH凭证未配置")
	}

	// 解密密码
	key := []byte(config.GlobalConfig.App.AesKey)
	decryptedPassword, err := encryption.DecryptAES(host.Password, key)
	if err != nil {
		return fmt.Errorf("解密密码失败: %w", err)
	}

	// 创建SSH连接
	sshClient, err := createSSHConnection(host.Username, decryptedPassword, targetIP, 22)
	if err != nil {
		return fmt.Errorf("SSH连接失败: %w", err)
	}
	defer sshClient.Close()

	// TODO: 使用新的SFTP客户端重构此功能
	// 暂时返回错误，避免编译问题
	return fmt.Errorf("文件上传功能正在重构中，请使用SFTP管理界面")
}

// createSSHConnection 创建SSH连接的通用方法
func createSSHConnection(username, password, host string, port int) (*gosshtool.Client, error) {
	config := &gosshtool.ClientConfig{
		User: username,
		Auth: []gosshtool.AuthMethod{
			gosshtool.Password(password),
		},
		HostKeyCallback: gosshtool.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := gosshtool.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("SSH连接失败: %w", err)
	}

	return client, nil
}

// BatchSFTPRequest 批量SFTP文件上传请求
type BatchSFTPRequest struct {
	IDs        []uint
	RemotePath string
	File       io.Reader
}
type SFTPResult struct {
	HostID  uint
	Success bool
	Error   string
}

// BatchSFTP 对多台主机批量上传文件
func (s *HostService) BatchSFTP(ids []uint, remotePath string, file io.Reader) []SFTPResult {
	var results []SFTPResult
	for _, id := range ids {
		host, err := s.hostRepo.FindByID(id)
		if err != nil {
			results = append(results, SFTPResult{HostID: id, Success: false, Error: err.Error()})
			continue
		}
		err = uploadFileToHost(host, remotePath, file)
		if err != nil {
			results = append(results, SFTPResult{HostID: id, Success: false, Error: err.Error()})
		} else {
			results = append(results, SFTPResult{HostID: id, Success: true})
		}
	}
	return results
}

// CreateManualHost 创建自建主机
func (s *HostService) CreateManualHost(host *model.Host) error {
	// 设置提供商类型为manual
	host.ProviderType = "manual"
	host.ProviderID = nil

	// 如果指定了组，检查组是否存在
	if host.GroupID != nil {
		if _, err := s.groupRepo.FindByID(*host.GroupID); err != nil {
			return fmt.Errorf("主机组不存在: %w", err)
		}
	}

	return s.hostRepo.CreateManualHost(host)
}

// GetGroupHosts 获取主机组下的主机列表
func (s *HostService) GetGroupHosts(groupID uint, page, pageSize int, keyword string) ([]model.Host, int64, error) {
	// 检查组是否存在
	if _, err := s.groupRepo.FindByID(groupID); err != nil {
		return nil, 0, fmt.Errorf("主机组不存在: %w", err)
	}

	return s.hostRepo.FindByGroupID(groupID, page, pageSize, keyword)
}

// MoveHost 移动主机到指定组
func (s *HostService) MoveHost(hostID uint, groupID *uint) error {
	// 检查主机是否存在
	host, err := s.hostRepo.FindByID(hostID)
	if err != nil {
		return fmt.Errorf("主机不存在: %w", err)
	}

	// 如果指定了组，检查组是否存在
	if groupID != nil {
		if _, err := s.groupRepo.FindByID(*groupID); err != nil {
			return fmt.Errorf("主机组不存在: %w", err)
		}
	}

	// 更新主机所属组
	host.GroupID = groupID
	return s.hostRepo.Update(host)
}

// getAllAESKeys 获取所有配置的AES密钥（主密钥+备份密钥）
func getAllAESKeys() [][]byte {
	var keys [][]byte
	if config.GlobalConfig.App.AesKey != "" {
		keys = append(keys, []byte(config.GlobalConfig.App.AesKey))
	}
	for _, k := range config.GlobalConfig.App.AesKeys {
		if k != "" {
			keys = append(keys, []byte(k))
		}
	}
	return keys
}

// FilterOptions 筛选选项结构
type FilterOptions struct {
	Statuses  []string `json:"statuses"`
	Regions   []string `json:"regions"`
	Providers []string `json:"providers"`
}

// GetFilterOptions 获取主机筛选选项
func (s *HostService) GetFilterOptions() (*FilterOptions, error) {
	// 获取所有主机的状态、地域、云厂商信息
	hosts, err := s.hostRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("获取主机列表失败: %w", err)
	}

	// 使用map去重
	statusMap := make(map[string]bool)
	regionMap := make(map[string]bool)
	providerMap := make(map[string]bool)

	for _, host := range hosts {
		if host.Status != "" {
			statusMap[host.Status] = true
		}
		if host.Region != "" {
			regionMap[host.Region] = true
		}
		if host.ProviderType != "" {
			providerMap[host.ProviderType] = true
		}
	}

	// 转换为切片
	var statuses, regions, providers []string
	for status := range statusMap {
		statuses = append(statuses, status)
	}
	for region := range regionMap {
		regions = append(regions, region)
	}
	for provider := range providerMap {
		providers = append(providers, provider)
	}

	return &FilterOptions{
		Statuses:  statuses,
		Regions:   regions,
		Providers: providers,
	}, nil
}

// getCloudAdapter 根据云账号类型获取对应的适配器
func (s *HostService) getCloudAdapter(provider *model.Provider) (adapter.CloudAdapter, error) {
	switch provider.Type {
	case model.ProviderTypeAliyun:
		return adapter.NewAliyunAdapter(provider.AccessKey, provider.SecretKey)
	case model.ProviderTypeAWS:
		return adapter.NewAWSAdapter(provider.AccessKey, provider.SecretKey, provider.Region)
	case model.ProviderTypeTencent:
		return adapter.NewTencentAdapter(provider.AccessKey, provider.SecretKey, provider.Region)
	default:
		return nil, fmt.Errorf("不支持的云服务商类型: %s", provider.Type)
	}
}
