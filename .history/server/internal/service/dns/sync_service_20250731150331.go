package dns

import (
	"context"
	"fmt"
	"time"

	"api-server/internal/model/dns"
	dnsprovider "api-server/internal/provider/dns"
	dnsrepo "api-server/internal/repository/dns"

	"go.uber.org/zap"
)

// SyncService DNS同步服务
type SyncService struct {
	domainRepo   *dnsrepo.DomainRepository
	recordRepo   *dnsrepo.RecordRepository
	providerRepo *dnsrepo.ProviderRepository
	metricsService *MetricsService
}

// NewSyncService 创建DNS同步服务
func NewSyncService(
	domainRepo *dnsrepo.DomainRepository,
	recordRepo *dnsrepo.RecordRepository,
	providerRepo *dnsrepo.ProviderRepository,
	metricsService *MetricsService,
) *SyncService {
	return &SyncService{
		domainRepo:     domainRepo,
		recordRepo:     recordRepo,
		providerRepo:   providerRepo,
		metricsService: metricsService,
	}
}

// RecordDiff DNS记录差异
type RecordDiff struct {
	ToAdd    []*dns.Record `json:"to_add"`    // 需要添加的记录
	ToUpdate []*dns.Record `json:"to_update"` // 需要更新的记录
	ToDelete []*dns.Record `json:"to_delete"` // 需要删除的记录
}

// SyncResult 同步结果
type SyncResult struct {
	DomainName     string        `json:"domain_name"`
	Provider       string        `json:"provider"`
	TotalLocal     int           `json:"total_local"`
	TotalRemote    int           `json:"total_remote"`
	AddedCount     int           `json:"added_count"`
	UpdatedCount   int           `json:"updated_count"`
	DeletedCount   int           `json:"deleted_count"`
	ErrorCount     int           `json:"error_count"`
	Duration       time.Duration `json:"duration"`
	Errors         []string      `json:"errors"`
	DryRun         bool          `json:"dry_run"`
}

// CompareRecords 比较本地和远程DNS记录
func (s *SyncService) CompareRecords(ctx context.Context, domainID uint, providerID uint) (*RecordDiff, error) {
	// 获取域名信息
	domain, err := s.domainRepo.FindByID(domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	// 获取DNS提供商信息
	provider, err := s.providerRepo.FindByID(providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	// 创建DNS驱动
	driver, err := s.createDriver(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	// 获取本地DNS记录
	localRecords, err := s.recordRepo.FindByDomainID(domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get local records: %w", err)
	}

	// 获取远程DNS记录
	remoteRecords, err := s.getRemoteRecords(ctx, driver, domain.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get remote records: %w", err)
	}

	// 比较记录差异
	diff := s.compareRecordLists(localRecords, remoteRecords)

	return diff, nil
}

// SyncDomainRecords 同步域名DNS记录
func (s *SyncService) SyncDomainRecords(ctx context.Context, domainID uint, providerID uint, dryRun bool) (*SyncResult, error) {
	start := time.Now()
	
	// 获取域名信息
	domain, err := s.domainRepo.FindByID(domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	// 获取DNS提供商信息
	provider, err := s.providerRepo.FindByID(providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	result := &SyncResult{
		DomainName: domain.Name,
		Provider:   provider.Name,
		DryRun:     dryRun,
		Errors:     []string{},
	}

	// 比较记录差异
	diff, err := s.CompareRecords(ctx, domainID, providerID)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to compare records: %v", err))
		result.Duration = time.Since(start)
		return result, err
	}

	result.TotalLocal = len(diff.ToAdd) + len(diff.ToUpdate) + len(diff.ToDelete)
	result.TotalRemote = result.TotalLocal // 简化计算

	if dryRun {
		// 试运行模式，只返回差异信息
		result.AddedCount = len(diff.ToAdd)
		result.UpdatedCount = len(diff.ToUpdate)
		result.DeletedCount = len(diff.ToDelete)
		result.Duration = time.Since(start)
		return result, nil
	}

	// 创建DNS驱动
	driver, err := s.createDriver(provider)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to create driver: %v", err))
		result.Duration = time.Since(start)
		return result, err
	}

	// 执行同步操作
	s.executeSyncOperations(ctx, driver, domain.Name, diff, result)

	result.Duration = time.Since(start)

	// 记录指标
	if s.metricsService != nil {
		status := "success"
		if result.ErrorCount > 0 {
			status = "partial_failure"
		}
		s.metricsService.RecordDNSSync(provider.Name, domain.Name, status, result.Duration)
	}

	return result, nil
}

// createDriver 创建DNS驱动
func (s *SyncService) createDriver(provider *dns.Provider) (SimpleDNSDriver, error) {
	// 解析凭证信息
	credentials := make(map[string]interface{})
	// TODO: 解密凭证信息
	// credentials = decryptCredentials(provider.CredentialsEnc)

	switch provider.Type {
	case "aliyun":
		return NewAliyunDriver(credentials)
	case "tencent":
		return NewTencentDriver(credentials)
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", provider.Type)
	}
}

// SimpleDNSDriver 简化的DNS驱动接口
type SimpleDNSDriver interface {
	GetName() string
	TestConnection(ctx context.Context) *TestResult
	ListRecords(ctx context.Context, domainName string) ([]*Record, error)
	CreateRecord(ctx context.Context, domainName string, record *Record) (*Record, error)
	UpdateRecord(ctx context.Context, domainName string, record *Record) (*Record, error)
	DeleteRecord(ctx context.Context, domainName string, recordID string) error
}

// Record 简化的DNS记录结构
type Record struct {
	ID       string
	Name     string
	Type     string
	Value    string
	TTL      int
	Priority int
	Status   string
}

// getRemoteRecords 获取远程DNS记录
func (s *SyncService) getRemoteRecords(ctx context.Context, driver SimpleDNSDriver, domainName string) ([]*Record, error) {
	records, err := driver.ListRecords(ctx, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to list remote records: %w", err)
	}

	return records, nil
}

// compareRecordLists 比较记录列表
func (s *SyncService) compareRecordLists(localRecords []*dns.Record, remoteRecords []*Record) *RecordDiff {
	diff := &RecordDiff{
		ToAdd:    []*dns.Record{},
		ToUpdate: []*dns.Record{},
		ToDelete: []*dns.Record{},
	}

	// 创建远程记录映射
	remoteMap := make(map[string]*Record)
	for _, record := range remoteRecords {
		key := fmt.Sprintf("%s:%s", record.Name, record.Type)
		remoteMap[key] = record
	}

	// 检查本地记录
	for _, localRecord := range localRecords {
		key := fmt.Sprintf("%s:%s", localRecord.Name, localRecord.Type)
		remoteRecord, exists := remoteMap[key]
		
		if !exists {
			// 远程不存在，需要添加
			diff.ToAdd = append(diff.ToAdd, localRecord)
		} else {
			// 检查是否需要更新
			if s.needsUpdate(localRecord, remoteRecord) {
				diff.ToUpdate = append(diff.ToUpdate, localRecord)
			}
			// 从映射中删除，剩下的就是需要删除的
			delete(remoteMap, key)
		}
	}

	// 剩余的远程记录需要删除
	for _, remoteRecord := range remoteMap {
		// 转换为本地记录格式
		localRecord := &dns.Record{
			Name:  remoteRecord.Name,
			Type:  remoteRecord.Type,
			Value: remoteRecord.Value,
		}
		diff.ToDelete = append(diff.ToDelete, localRecord)
	}

	return diff
}

// needsUpdate 检查记录是否需要更新
func (s *SyncService) needsUpdate(localRecord *dns.Record, remoteRecord *Record) bool {
	if localRecord.Value != remoteRecord.Value {
		return true
	}
	if localRecord.TTL != remoteRecord.TTL {
		return true
	}
	if localRecord.Priority != nil && *localRecord.Priority != remoteRecord.Priority {
		return true
	}
	return false
}

// executeSyncOperations 执行同步操作
func (s *SyncService) executeSyncOperations(ctx context.Context, driver SimpleDNSDriver, domainName string, diff *RecordDiff, result *SyncResult) {
	// 添加记录
	for _, record := range diff.ToAdd {
		remoteRecord := &Record{
			Name:     record.Name,
			Type:     record.Type,
			Value:    record.Value,
			TTL:      record.TTL,
			Priority: 0,
		}
		if record.Priority != nil {
			remoteRecord.Priority = *record.Priority
		}

		_, err := driver.CreateRecord(ctx, domainName, remoteRecord)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to create record %s: %v", record.Name, err))
			result.ErrorCount++
		} else {
			result.AddedCount++
		}
	}

	// 更新记录
	for _, record := range diff.ToUpdate {
		remoteRecord := &Record{
			Name:     record.Name,
			Type:     record.Type,
			Value:    record.Value,
			TTL:      record.TTL,
			Priority: 0,
		}
		if record.Priority != nil {
			remoteRecord.Priority = *record.Priority
		}

		_, err := driver.UpdateRecord(ctx, domainName, remoteRecord)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to update record %s: %v", record.Name, err))
			result.ErrorCount++
		} else {
			result.UpdatedCount++
		}
	}

	// 删除记录
	for _, record := range diff.ToDelete {
		err := driver.DeleteRecord(ctx, domainName, record.Name) // 简化，使用名称作为ID
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to delete record %s: %v", record.Name, err))
			result.ErrorCount++
		} else {
			result.DeletedCount++
		}
	}
}
