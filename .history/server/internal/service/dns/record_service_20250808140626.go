package dns

import (
	"api-server/internal/model/dns"
	repo "api-server/internal/repository/dns"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	dnsprovider "api-server/internal/provider/dns"

	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RecordService DNS记录业务逻辑服务
type RecordService struct {
	recordRepo    *repo.RecordRepository
	domainRepo    *repo.DomainRepository
	providerRepo  *repo.ProviderRepository
	changeLogRepo *repo.ChangeLogRepository
	db            *gorm.DB
}

// NewRecordService 创建DNS记录业务服务实例
func NewRecordService(
	recordRepo *repo.RecordRepository,
	domainRepo *repo.DomainRepository,
	providerRepo *repo.ProviderRepository,
	changeLogRepo *repo.ChangeLogRepository,
	db *gorm.DB,
) *RecordService {
	return &RecordService{
		recordRepo:    recordRepo,
		domainRepo:    domainRepo,
		providerRepo:  providerRepo,
		changeLogRepo: changeLogRepo,
		db:            db,
	}
}

// CreateRecord 创建DNS记录
func (s *RecordService) CreateRecord(record *dns.Record, actorID uint, clientIP string) error {
	// 验证域名是否存在
	domain, err := s.domainRepo.FindByID(record.DomainID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("指定的域名不存在")
		}
		return fmt.Errorf("验证域名失败: %w", err)
	}

	// 验证DNS提供商是否存在
	provider, err := s.providerRepo.FindByID(record.ProviderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("指定的DNS提供商不存在")
		}
		return fmt.Errorf("验证DNS提供商失败: %w", err)
	}

	// 验证记录格式
	if err := s.validateRecord(record); err != nil {
		return fmt.Errorf("DNS记录验证失败: %w", err)
	}

	// 检查记录是否已存在
	exists, err := s.recordRepo.ExistsByDomainAndNameAndType(record.DomainID, record.Name, record.Type)
	if err != nil {
		return fmt.Errorf("检查DNS记录是否存在失败: %w", err)
	}
	if exists {
		return errors.New("相同名称和类型的DNS记录已存在")
	}

	// 设置创建信息
	record.TenantID = domain.TenantID
	record.CreatedBy = actorID
	record.UpdatedBy = actorID
	record.SyncStatus = "pending"

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 创建记录
		if err := s.recordRepo.Create(record); err != nil {
			return fmt.Errorf("创建DNS记录失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "record",
			ResourceID:   record.ID,
			Action:       "create",
			Description:  fmt.Sprintf("创建DNS记录: %s.%s (%s)", record.Name, domain.Name, record.Type),
			NewData:      s.recordToJSON(record),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     record.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录DNS记录创建日志失败", zap.Error(err))
		}

		// TODO: 异步同步到DNS提供商
		go s.syncRecordToProvider(record, provider, "create")

		return nil
	})
}

// UpdateRecord 更新DNS记录
func (s *RecordService) UpdateRecord(record *dns.Record, actorID uint, clientIP string) error {
	// 获取原始数据
	oldRecord, err := s.recordRepo.FindByID(record.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("DNS记录不存在")
		}
		return fmt.Errorf("获取DNS记录信息失败: %w", err)
	}

	// 验证域名是否存在
	domain, err := s.domainRepo.FindByID(record.DomainID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("指定的域名不存在")
		}
		return fmt.Errorf("验证域名失败: %w", err)
	}

	// 验证DNS提供商是否存在
	provider, err := s.providerRepo.FindByID(record.ProviderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("指定的DNS提供商不存在")
		}
		return fmt.Errorf("验证DNS提供商失败: %w", err)
	}

	// 验证记录格式
	if err := s.validateRecord(record); err != nil {
		return fmt.Errorf("DNS记录验证失败: %w", err)
	}

	// 如果名称或类型发生变化，检查是否冲突
	if oldRecord.Name != record.Name || oldRecord.Type != record.Type {
		exists, err := s.recordRepo.ExistsByDomainAndNameAndType(record.DomainID, record.Name, record.Type)
		if err != nil {
			return fmt.Errorf("检查DNS记录是否存在失败: %w", err)
		}
		if exists {
			return errors.New("相同名称和类型的DNS记录已存在")
		}
	}

	// 设置更新信息
	record.UpdatedBy = actorID
	record.SyncStatus = "pending"

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 更新记录
		if err := s.recordRepo.Update(record); err != nil {
			return fmt.Errorf("更新DNS记录失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "record",
			ResourceID:   record.ID,
			Action:       "update",
			Description:  fmt.Sprintf("更新DNS记录: %s.%s (%s)", record.Name, domain.Name, record.Type),
			OldData:      s.recordToJSON(oldRecord),
			NewData:      s.recordToJSON(record),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     record.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录DNS记录更新日志失败", zap.Error(err))
		}

		// TODO: 异步同步到DNS提供商
		go s.syncRecordToProvider(record, provider, "update")

		return nil
	})
}

// DeleteRecord 删除DNS记录
func (s *RecordService) DeleteRecord(id uint, actorID uint, clientIP string) error {
	// 获取记录信息
	record, err := s.recordRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("DNS记录不存在")
		}
		return fmt.Errorf("获取DNS记录信息失败: %w", err)
	}

	// 获取域名信息
	domain, err := s.domainRepo.FindByID(record.DomainID)
	if err != nil {
		return fmt.Errorf("获取域名信息失败: %w", err)
	}

	// 获取提供商信息
	provider, err := s.providerRepo.FindByID(record.ProviderID)
	if err != nil {
		return fmt.Errorf("获取DNS提供商信息失败: %w", err)
	}

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除记录
		if err := s.recordRepo.Delete(id); err != nil {
			return fmt.Errorf("删除DNS记录失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "record",
			ResourceID:   id,
			Action:       "delete",
			Description:  fmt.Sprintf("删除DNS记录: %s.%s (%s)", record.Name, domain.Name, record.Type),
			OldData:      s.recordToJSON(record),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     record.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录DNS记录删除日志失败", zap.Error(err))
		}

		// TODO: 异步从DNS提供商删除
		go s.syncRecordToProvider(record, provider, "delete")

		return nil
	})
}

// GetRecord 获取DNS记录详情
func (s *RecordService) GetRecord(id uint) (*dns.Record, error) {
	record, err := s.recordRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("DNS记录不存在")
		}
		return nil, fmt.Errorf("获取DNS记录信息失败: %w", err)
	}
	return record, nil
}

// ListRecords 获取DNS记录列表
func (s *RecordService) ListRecords(tenantID uint, filters map[string]interface{}, limit, offset int) ([]*dns.Record, int64, error) {
	// 添加租户过滤
	if filters == nil {
		filters = make(map[string]interface{})
	}
	filters["tenant_id"] = tenantID

	records, total, err := s.recordRepo.SearchWithFilters(filters, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("获取DNS记录列表失败: %w", err)
	}

	return records, total, nil
}

// GetRecordsByDomain 获取指定域名的DNS记录
func (s *RecordService) GetRecordsByDomain(domainID uint) ([]*dns.Record, error) {
	records, err := s.recordRepo.FindByDomainID(domainID)
	if err != nil {
		return nil, fmt.Errorf("获取域名DNS记录失败: %w", err)
	}
	return records, nil
}

// GetRecordStatistics 获取DNS记录统计信息
func (s *RecordService) GetRecordStatistics(tenantID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总数统计
	total, err := s.recordRepo.CountByTenantID(tenantID)
	if err != nil {
		return nil, fmt.Errorf("统计DNS记录总数失败: %w", err)
	}
	stats["total"] = total

	// 类型统计
	typeCounts, err := s.recordRepo.CountByType()
	if err != nil {
		return nil, fmt.Errorf("统计DNS记录类型失败: %w", err)
	}
	stats["by_type"] = typeCounts

	// 状态统计
	statusCounts, err := s.recordRepo.CountByStatus()
	if err != nil {
		return nil, fmt.Errorf("统计DNS记录状态失败: %w", err)
	}
	stats["by_status"] = statusCounts

	return stats, nil
}

// validateRecord 验证DNS记录
func (s *RecordService) validateRecord(record *dns.Record) error {
	if record.Name == "" {
		return errors.New("记录名称不能为空")
	}

	if record.Type == "" {
		return errors.New("记录类型不能为空")
	}

	if record.Value == "" {
		return errors.New("记录值不能为空")
	}

	// 验证记录类型
	validTypes := []string{"A", "AAAA", "CNAME", "MX", "TXT", "SRV", "NS", "PTR", "CAA"}
	isValidType := false
	for _, validType := range validTypes {
		if record.Type == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("不支持的记录类型: %s", record.Type)
	}

	// 验证TTL
	if record.TTL < 60 || record.TTL > 86400 {
		return errors.New("TTL值必须在60-86400秒之间")
	}

	// TODO: 根据记录类型验证记录值格式

	return nil
}

// syncRecordToProvider 同步记录到DNS提供商
func (s *RecordService) syncRecordToProvider(record *dns.Record, provider *dns.Provider, action string) {
	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 解密凭证并创建驱动
	creds, err := s.getProviderCredentials(provider)
	if err != nil {
		zap.L().Error("获取DNS提供商凭证失败", zap.Error(err))
		_ = s.recordRepo.UpdateSyncStatus(record.ID, "error")
		return
	}
	driver, err := dnsprovider.CreateDriver(provider.Type, provider.Name, creds)
	if err != nil {
		zap.L().Error("创建DNS驱动失败", zap.Error(err))
		_ = s.recordRepo.UpdateSyncStatus(record.ID, "error")
		return
	}

	zoneName, err := s.getZoneName(record.DomainID)
	if err != nil {
		zap.L().Error("获取域名失败", zap.Error(err))
		_ = s.recordRepo.UpdateSyncStatus(record.ID, "error")
		return
	}

	switch action {
	case "create":
		prio := optionalInt(record.Priority)
		created, err2 := driver.CreateRecord(ctx, zoneName, &dnsprovider.Record{
			Name:     record.Name,
			Type:     record.Type,
			Value:    record.Value,
			TTL:      record.TTL,
			Priority: prio,
		})
		if err2 == nil && created != nil && created.ID != "" {
			_ = s.recordRepo.UpdateCloudRecordID(record.ID, created.ID)
		}
		err = err2
	case "update":
		prio := optionalInt(record.Priority)
		updated, err2 := driver.UpdateRecord(ctx, zoneName, &dnsprovider.Record{
			ID:       record.CloudRecordID,
			Name:     record.Name,
			Type:     record.Type,
			Value:    record.Value,
			TTL:      record.TTL,
			Priority: prio,
		})
		if err2 == nil && updated != nil && updated.ID != "" && updated.ID != record.CloudRecordID {
			_ = s.recordRepo.UpdateCloudRecordID(record.ID, updated.ID)
		}
		err = err2
	case "delete":
		targetID := record.CloudRecordID
		if targetID == "" {
			// 兼容：若无云端ID，回退到 name:type 组合
			targetID = fmt.Sprintf("%s:%s", record.Name, record.Type)
		}
		err = driver.DeleteRecord(ctx, zoneName, targetID)
	}

	if err != nil {
		zap.L().Error("同步DNS记录到提供商失败", zap.String("action", action), zap.Error(err))
		_ = s.recordRepo.UpdateSyncStatus(record.ID, "error")
		return
	}

	_ = s.recordRepo.UpdateSyncStatus(record.ID, "synced")
	zap.L().Info("DNS记录同步到提供商成功",
		zap.Uint("record_id", record.ID),
		zap.String("provider", provider.Name),
		zap.String("action", action))
}

func (s *RecordService) getZoneName(domainID uint) (string, error) {
	domain, err := s.domainRepo.FindByID(domainID)
	if err != nil {
		return "", err
	}
	return domain.Name, nil
}

func (s *RecordService) getProviderCredentials(provider *dns.Provider) (map[string]string, error) {
	// 这里复用 ProviderService 中的逻辑简化处理：凭证字段已加密为 JSON
	// 简化：允许明文 JSON 兼容
	var creds map[string]string
	if err := json.Unmarshal([]byte(provider.CredentialsEnc), &creds); err == nil && len(creds) > 0 {
		return creds, nil
	}
	// 若为加密，留给驱动内部失败；此处返回空避免崩溃
	return map[string]string{}, nil
}

func optionalInt(p *int) *int {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

// recordToJSON 将DNS记录对象转换为JSON
func (s *RecordService) recordToJSON(record *dns.Record) []byte {
	data, err := json.Marshal(record)
	if err != nil {
		zap.L().Error("DNS记录对象转JSON失败", zap.Error(err))
		return []byte("{}")
	}
	return data
}

// SyncDomainRecordsParams DNS记录同步参数
type SyncDomainRecordsParams struct {
	DomainID   uint   `json:"domain_id"`
	ProviderID uint   `json:"provider_id"`
	DryRun     bool   `json:"dry_run"`
	ActorID    uint   `json:"actor_id"`
	ClientIP   string `json:"client_ip"`
}

// RecordSyncResult DNS记录同步结果
type RecordSyncResult struct {
	DomainName     string        `json:"domain_name"`
	Provider       string        `json:"provider"`
	TotalLocal     int           `json:"total_local"`
	TotalRemote    int           `json:"total_remote"`
	ToAdd          int           `json:"to_add"`
	ToUpdate       int           `json:"to_update"`
	ToDelete       int           `json:"to_delete"`
	AddedRecords   []*dns.Record `json:"added_records"`
	UpdatedRecords []*dns.Record `json:"updated_records"`
	DeletedRecords []string      `json:"deleted_records"`
	Errors         []string      `json:"errors"`
	DryRun         bool          `json:"dry_run"`
}

// RemoteRecord 远程DNS记录结构
type RemoteRecord struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	TTL      int    `json:"ttl"`
	Priority int    `json:"priority"`
	Weight   int    `json:"weight"`
	Port     int    `json:"port"`
	RecordID string `json:"record_id"` // 远程记录ID
}

// SyncDomainRecords 同步域名DNS记录
func (s *RecordService) SyncDomainRecords(params *SyncDomainRecordsParams) (*RecordSyncResult, error) {
	zap.L().Info("开始同步DNS记录",
		zap.Uint("domain_id", params.DomainID),
		zap.Uint("provider_id", params.ProviderID),
		zap.Bool("dry_run", params.DryRun))

	// 获取域名信息
	domain, err := s.domainRepo.FindByID(params.DomainID)
	if err != nil {
		return nil, fmt.Errorf("获取域名信息失败: %w", err)
	}

	// 获取DNS提供商信息
	provider, err := s.providerRepo.FindByID(params.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("获取DNS提供商信息失败: %w", err)
	}

	// 初始化同步结果
	result := &RecordSyncResult{
		DomainName:     domain.Name,
		Provider:       provider.Name,
		AddedRecords:   []*dns.Record{},
		UpdatedRecords: []*dns.Record{},
		DeletedRecords: []string{},
		Errors:         []string{},
		DryRun:         params.DryRun,
	}

	// 获取本地DNS记录
	localRecords, err := s.recordRepo.FindByDomainID(params.DomainID)
	if err != nil {
		return nil, fmt.Errorf("获取本地DNS记录失败: %w", err)
	}
	result.TotalLocal = len(localRecords)

	// 获取远程DNS记录
	remoteRecords, err := s.getRemoteRecords(domain, provider)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("获取远程DNS记录失败: %v", err))
		return result, nil
	}
	result.TotalRemote = len(remoteRecords)

	// 比较记录差异
	toAdd, toUpdate, toDelete := s.compareRecords(localRecords, remoteRecords)
	result.ToAdd = len(toAdd)
	result.ToUpdate = len(toUpdate)
	result.ToDelete = len(toDelete)

	// 如果是试运行，直接返回差异信息
	if params.DryRun {
		zap.L().Info("DNS记录同步试运行结果",
			zap.Int("to_add", result.ToAdd),
			zap.Int("to_update", result.ToUpdate),
			zap.Int("to_delete", result.ToDelete))
		return result, nil
	}

	// 执行实际同步操作
	if err := s.executeSyncOperations(domain, provider, toAdd, toUpdate, toDelete, result, params); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("执行同步操作失败: %v", err))
	}

	// 记录同步操作日志
	s.logSyncOperation(domain, provider, result, params.ActorID, params.ClientIP)

	return result, nil
}

// getRemoteRecords 获取远程DNS记录
func (s *RecordService) getRemoteRecords(domain *dns.Domain, provider *dns.Provider) ([]*RemoteRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	creds, err := s.getProviderCredentials(provider)
	if err != nil {
		return nil, fmt.Errorf("获取凭证失败: %w", err)
	}
	driver, err := dnsprovider.CreateDriver(provider.Type, provider.Name, creds)
	if err != nil {
		return nil, fmt.Errorf("创建驱动失败: %w", err)
	}

	recs, err := driver.ListRecords(ctx, domain.Name, nil)
	if err != nil {
		return nil, fmt.Errorf("查询远程记录失败: %w", err)
	}

	out := make([]*RemoteRecord, 0, len(recs))
	for i := range recs {
		r := recs[i]
		out = append(out, &RemoteRecord{
			Name:     r.Name,
			Type:     r.Type,
			Value:    r.Value,
			TTL:      r.TTL,
			Priority: derefInt(r.Priority),
			Weight:   derefInt(r.Weight),
			Port:     derefInt(r.Port),
			RecordID: r.ID,
		})
	}
	return out, nil
}

func derefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// compareRecords 比较本地和远程记录的差异
func (s *RecordService) compareRecords(localRecords []*dns.Record, remoteRecords []*RemoteRecord) (
	toAdd []*RemoteRecord, toUpdate []*dns.Record, toDelete []*dns.Record) {

	// 创建本地记录映射，便于查找
	localMap := make(map[string]*dns.Record)
	for _, record := range localRecords {
		key := fmt.Sprintf("%s:%s", record.Name, record.Type)
		localMap[key] = record
	}

	// 创建远程记录映射
	remoteMap := make(map[string]*RemoteRecord)
	for _, record := range remoteRecords {
		key := fmt.Sprintf("%s:%s", record.Name, record.Type)
		remoteMap[key] = record
	}

	// 找出需要添加的记录（远程有，本地没有）
	for key, remoteRecord := range remoteMap {
		if _, exists := localMap[key]; !exists {
			toAdd = append(toAdd, remoteRecord)
		}
	}

	// 找出需要更新的记录（本地和远程都有，但内容不同）
	for key, localRecord := range localMap {
		if remoteRecord, exists := remoteMap[key]; exists {
			if s.recordNeedsUpdate(localRecord, remoteRecord) {
				toUpdate = append(toUpdate, localRecord)
			}
		}
	}

	// 找出需要删除的记录（本地有，远程没有）
	for key, localRecord := range localMap {
		if _, exists := remoteMap[key]; !exists {
			toDelete = append(toDelete, localRecord)
		}
	}

	return toAdd, toUpdate, toDelete
}

// recordNeedsUpdate 判断记录是否需要更新
func (s *RecordService) recordNeedsUpdate(local *dns.Record, remote *RemoteRecord) bool {
	if local.Value != remote.Value || local.TTL != remote.TTL {
		return true
	}
	if local.Priority != nil && *local.Priority != remote.Priority {
		return true
	}
	if local.Weight != nil && *local.Weight != remote.Weight {
		return true
	}
	if local.Port != nil && *local.Port != remote.Port {
		return true
	}
	return false
}

// executeSyncOperations 执行同步操作
func (s *RecordService) executeSyncOperations(
	domain *dns.Domain,
	provider *dns.Provider,
	toAdd []*RemoteRecord,
	toUpdate []*dns.Record,
	toDelete []*dns.Record,
	result *RecordSyncResult,
	params *SyncDomainRecordsParams,
) error {

	// 添加新记录
	for _, remoteRecord := range toAdd {
		newRecord := &dns.Record{
			DomainID:   domain.ID,
			Name:       remoteRecord.Name,
			Type:       remoteRecord.Type,
			Value:      remoteRecord.Value,
			TTL:        remoteRecord.TTL,
			Priority:   &remoteRecord.Priority,
			Weight:     &remoteRecord.Weight,
			Port:       &remoteRecord.Port,
			Status:     "active",
			SyncStatus: "synced",
			ProviderID: provider.ID,
			TenantID:   domain.TenantID,
			CreatedBy:  params.ActorID,
			UpdatedBy:  params.ActorID,
			Remark:     fmt.Sprintf("从%s同步", provider.Name),
		}

		if err := s.recordRepo.Create(newRecord); err != nil {
			result.Errors = append(result.Errors,
				fmt.Sprintf("创建记录失败 %s:%s - %v", remoteRecord.Name, remoteRecord.Type, err))
			continue
		}

		result.AddedRecords = append(result.AddedRecords, newRecord)
		zap.L().Info("添加DNS记录成功",
			zap.String("name", newRecord.Name),
			zap.String("type", newRecord.Type))
	}

	// 更新现有记录
	for _, localRecord := range toUpdate {
		// 这里需要重新获取远程记录来更新，简化处理
		localRecord.SyncStatus = "synced"
		localRecord.UpdatedBy = params.ActorID

		if err := s.recordRepo.Update(localRecord); err != nil {
			result.Errors = append(result.Errors,
				fmt.Sprintf("更新记录失败 %s:%s - %v", localRecord.Name, localRecord.Type, err))
			continue
		}

		result.UpdatedRecords = append(result.UpdatedRecords, localRecord)
		zap.L().Info("更新DNS记录成功",
			zap.String("name", localRecord.Name),
			zap.String("type", localRecord.Type))
	}

	// 删除多余的记录
	for _, localRecord := range toDelete {
		if err := s.recordRepo.Delete(localRecord.ID); err != nil {
			result.Errors = append(result.Errors,
				fmt.Sprintf("删除记录失败 %s:%s - %v", localRecord.Name, localRecord.Type, err))
			continue
		}

		result.DeletedRecords = append(result.DeletedRecords, fmt.Sprintf("%d", localRecord.ID))
		zap.L().Info("删除DNS记录成功",
			zap.String("name", localRecord.Name),
			zap.String("type", localRecord.Type))
	}

	return nil
}

// logSyncOperation 记录同步操作日志
func (s *RecordService) logSyncOperation(
	domain *dns.Domain,
	provider *dns.Provider,
	result *RecordSyncResult,
	actorID uint,
	clientIP string,
) {
	description := fmt.Sprintf("同步域名%s的DNS记录，提供商：%s，添加：%d，更新：%d，删除：%d",
		domain.Name, provider.Name, result.ToAdd, result.ToUpdate, result.ToDelete)

	changeLog := &dns.ChangeLog{
		ResourceType: "domain",
		ResourceID:   domain.ID,
		Action:       "sync",
		Description:  description,
		Status:       "success",
		ClientIP:     clientIP,
		TenantID:     domain.TenantID,
		ActorID:      actorID,
	}

	if len(result.Errors) > 0 {
		changeLog.Status = "partial_success"
		changeLog.Description += fmt.Sprintf("，错误：%d", len(result.Errors))
	}

	if err := s.changeLogRepo.Create(changeLog); err != nil {
		zap.L().Error("记录DNS同步操作日志失败", zap.Error(err))
	}
}

// GetSyncStatus 获取域名同步状态
func (s *RecordService) GetSyncStatus(domainID uint) (map[string]interface{}, error) {
	// 获取域名下记录的同步状态统计
	records, err := s.recordRepo.FindByDomainID(domainID)
	if err != nil {
		return nil, fmt.Errorf("获取域名记录失败: %w", err)
	}

	syncStatusCounts := make(map[string]int)
	for _, record := range records {
		syncStatusCounts[record.SyncStatus]++
	}

	// 获取最近的同步记录
	// TODO: 这里可以从change_log表中获取最近的同步操作信息

	status := map[string]interface{}{
		"total_records":     len(records),
		"sync_status_count": syncStatusCounts,
		"last_sync_at":      nil,   // TODO: 实现
		"is_syncing":        false, // TODO: 实现
	}

	return status, nil
}
