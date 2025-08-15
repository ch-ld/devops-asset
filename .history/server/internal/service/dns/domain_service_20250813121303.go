package dns

import (
	"api-server/internal/model/dns"
	repo "api-server/internal/repository/dns"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DomainService 域名业务逻辑服务
type DomainService struct {
	domainRepo      *repo.DomainRepository
	domainGroupRepo *repo.DomainGroupRepository
	changeLogRepo   *repo.ChangeLogRepository
	db              *gorm.DB
}

// NewDomainService 创建域名业务服务实例
func NewDomainService(
	domainRepo *repo.DomainRepository,
	domainGroupRepo *repo.DomainGroupRepository,
	changeLogRepo *repo.ChangeLogRepository,
	db *gorm.DB,
) *DomainService {
	return &DomainService{
		domainRepo:      domainRepo,
		domainGroupRepo: domainGroupRepo,
		changeLogRepo:   changeLogRepo,
		db:              db,
	}
}

// CreateDomain 创建域名
func (s *DomainService) CreateDomain(domain *dns.Domain, actorID uint, clientIP string) error {
	// 验证域名格式
	if err := s.validateDomainName(domain.Name); err != nil {
		return fmt.Errorf("域名格式验证失败: %w", err)
	}

	// 检查域名是否已存在
	exists, err := s.domainRepo.ExistsByName(domain.Name)
	if err != nil {
		return fmt.Errorf("检查域名是否存在失败: %w", err)
	}
	if exists {
		return errors.New("域名已存在")
	}

	// 验证分组是否存在
	if domain.GroupID != nil {
		_, err := s.domainGroupRepo.FindByID(*domain.GroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("指定的域名分组不存在")
			}
			return fmt.Errorf("验证域名分组失败: %w", err)
		}
	}

	// 设置创建信息
	domain.CreatedBy = actorID
	domain.UpdatedBy = actorID

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 创建域名
		if err := s.domainRepo.Create(domain); err != nil {
			return fmt.Errorf("创建域名失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "domain",
			ResourceID:   domain.ID,
			Action:       "create",
			Description:  fmt.Sprintf("创建域名: %s", domain.Name),
			NewData:      s.domainToJSON(domain),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     domain.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录域名创建日志失败", zap.Error(err))
		}

		return nil
	})
}

// UpdateDomain 更新域名
func (s *DomainService) UpdateDomain(domain *dns.Domain, actorID uint, clientIP string) error {
	// 获取原始数据
	oldDomain, err := s.domainRepo.FindByID(domain.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("域名不存在")
		}
		return fmt.Errorf("获取域名信息失败: %w", err)
	}

	// 验证域名格式
	if err := s.validateDomainName(domain.Name); err != nil {
		return fmt.Errorf("域名格式验证失败: %w", err)
	}

	// 如果域名名称发生变化，检查新名称是否已存在
	if oldDomain.Name != domain.Name {
		exists, err := s.domainRepo.ExistsByName(domain.Name)
		if err != nil {
			return fmt.Errorf("检查域名是否存在失败: %w", err)
		}
		if exists {
			return errors.New("域名已存在")
		}
	}

	// 验证分组是否存在
	if domain.GroupID != nil {
		_, err := s.domainGroupRepo.FindByID(*domain.GroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("指定的域名分组不存在")
			}
			return fmt.Errorf("验证域名分组失败: %w", err)
		}
	}

	// 设置更新信息
	domain.UpdatedBy = actorID

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 更新域名
		if err := s.domainRepo.Update(domain); err != nil {
			return fmt.Errorf("更新域名失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "domain",
			ResourceID:   domain.ID,
			Action:       "update",
			Description:  fmt.Sprintf("更新域名: %s", domain.Name),
			OldData:      s.domainToJSON(oldDomain),
			NewData:      s.domainToJSON(domain),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     domain.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录域名更新日志失败", zap.Error(err))
		}

		return nil
	})
}

// DeleteDomain 删除域名
func (s *DomainService) DeleteDomain(id uint, actorID uint, clientIP string) error {
	// 获取域名信息
	domain, err := s.domainRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("域名不存在")
		}
		return fmt.Errorf("获取域名信息失败: %w", err)
	}

	// TODO: 检查是否有关联的DNS记录或证书
	// 这里可以添加业务规则，比如有关联记录时不允许删除

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除域名
		if err := s.domainRepo.Delete(id); err != nil {
			return fmt.Errorf("删除域名失败: %w", err)
		}

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "domain",
			ResourceID:   id,
			Action:       "delete",
			Description:  fmt.Sprintf("删除域名: %s", domain.Name),
			OldData:      s.domainToJSON(domain),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     domain.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录域名删除日志失败", zap.Error(err))
		}

		return nil
	})
}

// GetDomain 获取域名详情
func (s *DomainService) GetDomain(id uint) (*dns.Domain, error) {
	domain, err := s.domainRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("域名不存在")
		}
		return nil, fmt.Errorf("获取域名信息失败: %w", err)
	}
	return domain, nil
}

// GetDomainByName 根据名称获取域名
func (s *DomainService) GetDomainByName(name string) (*dns.Domain, error) {
	domain, err := s.domainRepo.FindByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("域名不存在")
		}
		return nil, fmt.Errorf("获取域名信息失败: %w", err)
	}
	return domain, nil
}

// ListDomains 获取域名列表
func (s *DomainService) ListDomains(tenantID uint, filters map[string]interface{}, limit, offset int) ([]*dns.Domain, int64, error) {
	// 添加租户过滤
	if filters == nil {
		filters = make(map[string]interface{})
	}
	filters["tenant_id"] = tenantID

	domains, total, err := s.domainRepo.SearchWithFilters(filters, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("获取域名列表失败: %w", err)
	}

	return domains, total, nil
}

// GetExpiringDomains 获取即将过期的域名
func (s *DomainService) GetExpiringDomains(days int) ([]*dns.Domain, error) {
	domains, err := s.domainRepo.FindExpiring(days)
	if err != nil {
		return nil, fmt.Errorf("获取即将过期域名失败: %w", err)
	}
	return domains, nil
}

// GetDomainStatistics 获取域名统计信息
func (s *DomainService) GetDomainStatistics(tenantID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总数统计
	total, err := s.domainRepo.CountByTenantID(tenantID)
	if err != nil {
		return nil, fmt.Errorf("统计域名总数失败: %w", err)
	}
	stats["total"] = total

	// 状态统计（按租户过滤）
	statusCounts, err := s.domainRepo.CountByStatusAndTenant(tenantID)
	if err != nil {
		return nil, fmt.Errorf("统计域名状态失败: %w", err)
	}
	stats["by_status"] = statusCounts

	// 即将过期统计（30天内）
	expiringDomains, err := s.GetExpiringDomains(30)
	if err != nil {
		return nil, fmt.Errorf("统计即将过期域名失败: %w", err)
	}
	stats["expiring_count"] = len(expiringDomains)

	return stats, nil
}

// CountByProviderID 根据提供商ID统计域名数量
func (s *DomainService) CountByProviderID(providerID uint) (int64, error) {
	count, err := s.domainRepo.CountByProviderID(providerID)
	if err != nil {
		return 0, fmt.Errorf("统计提供商域名数量失败: %w", err)
	}
	return count, nil
}

// validateDomainName 验证域名格式
func (s *DomainService) validateDomainName(name string) error {
	if name == "" {
		return errors.New("域名不能为空")
	}

	name = strings.TrimSpace(strings.ToLower(name))
	if len(name) > 253 {
		return errors.New("域名长度不能超过253个字符")
	}

	// 简单的域名格式验证
	if !strings.Contains(name, ".") {
		return errors.New("域名格式不正确")
	}

	return nil
}

// domainToJSON 将域名对象转换为JSON
func (s *DomainService) domainToJSON(domain *dns.Domain) []byte {
	data, err := json.Marshal(domain)
	if err != nil {
		zap.L().Error("域名对象转JSON失败", zap.Error(err))
		return []byte("{}")
	}
	return data
}

// 在文件末尾添加导入导出功能

// DomainImportResult 域名导入结果
type DomainImportResult struct {
	Total    int                  `json:"total"`
	Success  int                  `json:"success"`
	Failed   int                  `json:"failed"`
	Skipped  int                  `json:"skipped"`
	Errors   []DomainImportError  `json:"errors"`
	Details  []DomainImportDetail `json:"details"`
	Duration time.Duration        `json:"duration"`
}

// DomainImportError 导入错误信息
type DomainImportError struct {
	Row     int    `json:"row"`
	Domain  string `json:"domain"`
	Message string `json:"message"`
}

// DomainImportDetail 导入详情
type DomainImportDetail struct {
	Row      int    `json:"row"`
	Domain   string `json:"domain"`
	Status   string `json:"status"` // success, failed, skipped
	Message  string `json:"message"`
	DomainID uint   `json:"domain_id,omitempty"`
}

// DomainExportData 域名导出数据
type DomainExportData struct {
	Name          string `json:"name" csv:"域名"`
	Status        string `json:"status" csv:"状态"`
	RegistrarType string `json:"registrar_type" csv:"注册商类型"`
	ExpiresAt     string `json:"expires_at" csv:"过期时间"`
	AutoRenew     string `json:"auto_renew" csv:"自动续费"`
	GroupName     string `json:"group_name" csv:"分组"`
	Remark        string `json:"remark" csv:"备注"`
	CreatedAt     string `json:"created_at" csv:"创建时间"`
}

// ImportDomainsFromCSV 从CSV文件导入域名
func (s *DomainService) ImportDomainsFromCSV(reader io.Reader, tenantID, actorID uint, clientIP string) (*DomainImportResult, error) {
	startTime := time.Now()
	result := &DomainImportResult{
		Errors:  make([]DomainImportError, 0),
		Details: make([]DomainImportDetail, 0),
	}

	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1 // 允许不同记录有不同字段数

	// 读取表头
	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("读取CSV表头失败: %w", err)
	}

	// 验证表头
	requiredHeaders := []string{"域名", "状态", "注册商类型", "过期时间", "自动续费", "分组", "备注"}
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.TrimSpace(header)] = i
	}

	for _, required := range requiredHeaders {
		if _, exists := headerMap[required]; !exists {
			return nil, fmt.Errorf("缺少必需的列: %s", required)
		}
	}

	// 读取数据行
	rowIndex := 1 // 从1开始，因为第0行是表头
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, DomainImportError{
				Row:     rowIndex,
				Message: fmt.Sprintf("读取CSV行失败: %v", err),
			})
			rowIndex++
			continue
		}

		result.Total++
		if err := s.importDomainRecord(record, headerMap, tenantID, actorID, clientIP, rowIndex, result); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, DomainImportError{
				Row:     rowIndex,
				Domain:  s.getFieldValue(record, headerMap, "域名"),
				Message: err.Error(),
			})
		} else {
			result.Success++
		}
		rowIndex++
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

// ImportDomainsFromExcel 从Excel文件导入域名
func (s *DomainService) ImportDomainsFromExcel(reader io.Reader, tenantID, actorID uint, clientIP string) (*DomainImportResult, error) {
	startTime := time.Now()
	result := &DomainImportResult{
		Errors:  make([]DomainImportError, 0),
		Details: make([]DomainImportDetail, 0),
	}

	// 读取Excel文件
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("读取Excel文件失败: %w", err)
	}
	defer f.Close()

	// 获取第一个工作表
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取Excel工作表失败: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel文件至少需要包含表头和一行数据")
	}

	// 处理表头
	headers := rows[0]
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.TrimSpace(header)] = i
	}

	// 验证必需的列
	requiredHeaders := []string{"域名", "状态", "注册商类型", "过期时间", "自动续费", "分组", "备注"}
	for _, required := range requiredHeaders {
		if _, exists := headerMap[required]; !exists {
			return nil, fmt.Errorf("缺少必需的列: %s", required)
		}
	}

	// 处理数据行
	for rowIndex, record := range rows[1:] {
		actualRowIndex := rowIndex + 2 // Excel行号从1开始，加上表头行
		result.Total++

		if err := s.importDomainRecord(record, headerMap, tenantID, actorID, clientIP, actualRowIndex, result); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, DomainImportError{
				Row:     actualRowIndex,
				Domain:  s.getFieldValue(record, headerMap, "域名"),
				Message: err.Error(),
			})
		} else {
			result.Success++
		}
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

// ExportDomainsToCSV 导出域名到CSV格式
func (s *DomainService) ExportDomainsToCSV(writer io.Writer, tenantID uint, filters map[string]interface{}) error {
	// 获取域名列表
	domains, err := s.domainRepo.FindAll()
	if err != nil {
		return fmt.Errorf("获取域名列表失败: %w", err)
	}

	// TODO: 根据filters参数过滤结果

	// 准备导出数据
	exportData := make([]DomainExportData, len(domains))
	for i, domain := range domains {
		exportData[i] = s.prepareDomainExportData(domain)
	}

	// 写入CSV
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// 写入表头
	headers := []string{"域名", "状态", "注册商类型", "过期时间", "自动续费", "分组", "备注", "创建时间"}
	if err := csvWriter.Write(headers); err != nil {
		return fmt.Errorf("写入CSV表头失败: %w", err)
	}

	// 写入数据
	for _, data := range exportData {
		record := []string{
			data.Name,
			data.Status,
			data.RegistrarType,
			data.ExpiresAt,
			data.AutoRenew,
			data.GroupName,
			data.Remark,
			data.CreatedAt,
		}
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("写入CSV数据失败: %w", err)
		}
	}

	return nil
}

// ExportDomainsToExcel 导出域名到Excel格式
func (s *DomainService) ExportDomainsToExcel(writer io.Writer, tenantID uint, filters map[string]interface{}) error {
	// 获取域名列表
	domains, err := s.domainRepo.FindAll()
	if err != nil {
		return fmt.Errorf("获取域名列表失败: %w", err)
	}

	// TODO: 根据filters参数过滤结果

	// 创建Excel文件
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "域名列表"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("创建Excel工作表失败: %w", err)
	}

	// 设置表头
	headers := []string{"域名", "状态", "注册商类型", "过期时间", "自动续费", "分组", "备注", "创建时间"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
	}

	// 设置数据
	for i, domain := range domains {
		row := i + 2
		data := s.prepareDomainExportData(domain)

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), data.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), data.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), data.RegistrarType)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), data.ExpiresAt)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), data.AutoRenew)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), data.GroupName)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), data.Remark)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), data.CreatedAt)
	}

	// 设置活动工作表
	f.SetActiveSheet(index)

	// 写入文件
	if err := f.Write(writer); err != nil {
		return fmt.Errorf("写入Excel文件失败: %w", err)
	}

	return nil
}

// importDomainRecord 导入单个域名记录
func (s *DomainService) importDomainRecord(record []string, headerMap map[string]int, tenantID, actorID uint, clientIP string, rowIndex int, result *DomainImportResult) error {
	domainName := s.getFieldValue(record, headerMap, "域名")
	if domainName == "" {
		return fmt.Errorf("域名不能为空")
	}

	// 检查域名是否已存在
	exists, err := s.domainRepo.ExistsByName(domainName)
	if err != nil {
		return fmt.Errorf("检查域名是否存在失败: %w", err)
	}
	if exists {
		result.Skipped++
		result.Details = append(result.Details, DomainImportDetail{
			Row:     rowIndex,
			Domain:  domainName,
			Status:  "skipped",
			Message: "域名已存在",
		})
		return nil
	}

	// 解析过期时间
	var expiresAt *time.Time
	if expiresStr := s.getFieldValue(record, headerMap, "过期时间"); expiresStr != "" {
		if parsed, err := time.Parse("2006-01-02", expiresStr); err == nil {
			expiresAt = &parsed
		} else if parsed, err := time.Parse("2006-01-02 15:04:05", expiresStr); err == nil {
			expiresAt = &parsed
		}
	}

	// 解析自动续费
	autoRenew := false
	if autoRenewStr := s.getFieldValue(record, headerMap, "自动续费"); autoRenewStr != "" {
		autoRenew = strings.ToLower(autoRenewStr) == "true" || autoRenewStr == "1" || strings.ToLower(autoRenewStr) == "是"
	}

	// 查找分组ID
	var groupID *uint
	if groupName := s.getFieldValue(record, headerMap, "分组"); groupName != "" {
		if group, err := s.domainGroupRepo.FindByName(groupName); err == nil {
			groupID = &group.ID
		}
	}

	// 创建域名
	domain := &dns.Domain{
		Name:          domainName,
		Status:        s.getFieldValue(record, headerMap, "状态"),
		RegistrarType: s.getFieldValue(record, headerMap, "注册商类型"),
		ExpiresAt:     expiresAt,
		AutoRenew:     autoRenew,
		GroupID:       groupID,
		Remark:        s.getFieldValue(record, headerMap, "备注"),
		TenantID:      tenantID,
		CreatedBy:     actorID,
		UpdatedBy:     actorID,
	}

	if err := s.domainRepo.Create(domain); err != nil {
		return fmt.Errorf("创建域名失败: %w", err)
	}

	result.Details = append(result.Details, DomainImportDetail{
		Row:      rowIndex,
		Domain:   domainName,
		Status:   "success",
		Message:  "导入成功",
		DomainID: domain.ID,
	})

	return nil
}

// getFieldValue 从记录中获取指定字段的值
func (s *DomainService) getFieldValue(record []string, headerMap map[string]int, fieldName string) string {
	if index, exists := headerMap[fieldName]; exists && index < len(record) {
		return strings.TrimSpace(record[index])
	}
	return ""
}

// prepareDomainExportData 准备域名导出数据
func (s *DomainService) prepareDomainExportData(domain *dns.Domain) DomainExportData {
	data := DomainExportData{
		Name:          domain.Name,
		Status:        domain.Status,
		RegistrarType: domain.RegistrarType,
		Remark:        domain.Remark,
		CreatedAt:     domain.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if domain.ExpiresAt != nil {
		data.ExpiresAt = domain.ExpiresAt.Format("2006-01-02")
	}

	if domain.AutoRenew {
		data.AutoRenew = "是"
	} else {
		data.AutoRenew = "否"
	}

	// 获取分组名称
	if domain.GroupID != nil && domain.Group != nil {
		data.GroupName = domain.Group.Name
	}

	return data
}

// GetDomainWhois 获取域名WHOIS信息
func (s *DomainService) GetDomainWhois(domainName string) (map[string]interface{}, error) {
	// 实际的WHOIS查询实现
	whoisInfo, err := s.queryWhoisInfo(domainName)
	if err != nil {
		return nil, fmt.Errorf("WHOIS查询失败: %w", err)
	}

	return whoisInfo, nil
}

// queryWhoisInfo 执行实际的WHOIS查询
func (s *DomainService) queryWhoisInfo(domainName string) (map[string]interface{}, error) {
	// 使用系统whois命令查询
	cmd := exec.Command("whois", domainName)
	output, err := cmd.Output()
	if err != nil {
		// 如果系统没有whois命令，使用网络查询
		return s.queryWhoisViaNetwork(domainName)
	}

	// 解析whois输出
	return s.parseWhoisOutput(domainName, string(output))
}

// queryWhoisViaNetwork 通过网络查询WHOIS信息
func (s *DomainService) queryWhoisViaNetwork(domainName string) (map[string]interface{}, error) {
	// 获取域名的TLD
	parts := strings.Split(domainName, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("无效的域名格式")
	}
	tld := parts[len(parts)-1]

	// 根据TLD选择WHOIS服务器
	whoisServer := s.getWhoisServer(tld)
	if whoisServer == "" {
		return nil, fmt.Errorf("不支持的TLD: %s", tld)
	}

	// 连接WHOIS服务器
	conn, err := net.DialTimeout("tcp", whoisServer+":43", 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("连接WHOIS服务器失败: %w", err)
	}
	defer conn.Close()

	// 发送查询
	_, err = conn.Write([]byte(domainName + "\r\n"))
	if err != nil {
		return nil, fmt.Errorf("发送查询失败: %w", err)
	}

	// 读取响应
	response, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	return s.parseWhoisOutput(domainName, string(response))
}

// getWhoisServer 根据TLD获取WHOIS服务器
func (s *DomainService) getWhoisServer(tld string) string {
	whoisServers := map[string]string{
		"com":    "whois.verisign-grs.com",
		"net":    "whois.verisign-grs.com",
		"org":    "whois.pir.org",
		"info":   "whois.afilias.net",
		"biz":    "whois.neulevel.biz",
		"name":   "whois.nic.name",
		"cn":     "whois.cnnic.cn",
		"tw":     "whois.twnic.net.tw",
		"hk":     "whois.hkirc.hk",
		"jp":     "whois.jprs.jp",
		"kr":     "whois.nida.or.kr",
		"uk":     "whois.nominet.uk",
		"de":     "whois.denic.de",
		"fr":     "whois.afnic.fr",
		"it":     "whois.nic.it",
		"nl":     "whois.domain-registry.nl",
		"be":     "whois.dns.be",
		"ch":     "whois.nic.ch",
		"at":     "whois.nic.at",
		"se":     "whois.iis.se",
		"no":     "whois.norid.no",
		"dk":     "whois.dk-hostmaster.dk",
		"fi":     "whois.fi",
		"pl":     "whois.dns.pl",
		"cz":     "whois.nic.cz",
		"sk":     "whois.sk-nic.sk",
		"hu":     "whois.nic.hu",
		"ro":     "whois.rotld.ro",
		"bg":     "whois.register.bg",
		"hr":     "whois.dns.hr",
		"si":     "whois.arnes.si",
		"lt":     "whois.domreg.lt",
		"lv":     "whois.nic.lv",
		"ee":     "whois.eenet.ee",
		"ru":     "whois.tcinet.ru",
		"ua":     "whois.ua",
		"by":     "whois.cctld.by",
		"kz":     "whois.nic.kz",
		"kg":     "whois.domain.kg",
		"uz":     "whois.cctld.uz",
		"tj":     "whois.nic.tj",
		"tm":     "whois.nic.tm",
		"am":     "whois.amnic.net",
		"az":     "whois.az",
		"ge":     "whois.nic.ge",
		"md":     "whois.nic.md",
		"ca":     "whois.cira.ca",
		"us":     "whois.nic.us",
		"mx":     "whois.mx",
		"br":     "whois.registro.br",
		"ar":     "whois.nic.ar",
		"cl":     "whois.nic.cl",
		"co":     "whois.nic.co",
		"ve":     "whois.nic.ve",
		"pe":     "kero.yachay.pe",
		"ec":     "whois.nic.ec",
		"uy":     "whois.nic.org.uy",
		"py":     "whois.nic.py",
		"bo":     "whois.nic.bo",
		"au":     "whois.auda.org.au",
		"nz":     "whois.srs.net.nz",
		"in":     "whois.registry.in",
		"sg":     "whois.sgnic.sg",
		"my":     "whois.mynic.my",
		"th":     "whois.thnic.co.th",
		"id":     "whois.id",
		"ph":     "whois.dot.ph",
		"vn":     "whois.nic.vn",
		"za":     "whois.registry.net.za",
		"eg":     "whois.ripe.net",
		"ma":     "whois.registre.ma",
		"tn":     "whois.ati.tn",
		"ke":     "whois.kenic.or.ke",
		"ng":     "whois.nic.net.ng",
		"gh":     "whois.nic.gh",
		"mu":     "whois.nic.mu",
		"shop":   "whois.nic.shop",
		"online": "whois.nic.online",
		"site":   "whois.nic.site",
		"store":  "whois.nic.store",
		"tech":   "whois.nic.tech",
		"app":    "whois.nic.google",
		"dev":    "whois.nic.google",
		"page":   "whois.nic.google",
		"cloud":  "whois.nic.cloud",
		"ai":     "whois.nic.ai",
		"io":     "whois.nic.io",
		"me":     "whois.nic.me",
		"tv":     "whois.nic.tv",
		"cc":     "whois.nic.cc",
		"ws":     "whois.website.ws",
		"mobi":   "whois.dotmobiregistry.net",
		"tel":    "whois.nic.tel",
		"travel": "whois.nic.travel",
		"museum": "whois.museum",
		"aero":   "whois.information.aero",
		"coop":   "whois.nic.coop",
		"jobs":   "whois.nic.jobs",
		"cat":    "whois.nic.cat",
		"pro":    "whois.registrypro.pro",
		"asia":   "whois.nic.asia",
		"xxx":    "whois.nic.xxx",
		"post":   "whois.dotpostregistry.net",
	}

	return whoisServers[strings.ToLower(tld)]
}

// parseWhoisOutput 解析WHOIS输出
func (s *DomainService) parseWhoisOutput(domainName, output string) (map[string]interface{}, error) {
	lines := strings.Split(output, "\n")
	result := map[string]interface{}{
		"domain_name":     domainName,
		"raw_data":        output,
		"query_time":      time.Now().Format("2006-01-02 15:04:05"),
		"registrar":       "",
		"creation_date":   "",
		"expiration_date": "",
		"updated_date":    "",
		"status":          []string{},
		"name_servers":    []string{},
		"registrant":      map[string]interface{}{},
		"admin_contact":   map[string]interface{}{},
		"tech_contact":    map[string]interface{}{},
	}

	var nameServers []string
	var statusList []string
	registrant := make(map[string]interface{})
	adminContact := make(map[string]interface{})
	techContact := make(map[string]interface{})

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "%") || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析各种字段
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(strings.ToLower(parts[0]))
			value := strings.TrimSpace(parts[1])

			switch {
			// 注册商信息
			case strings.Contains(key, "registrar") && !strings.Contains(key, "url") && !strings.Contains(key, "phone"):
				if value != "" {
					result["registrar"] = value
				}

			// 日期信息
			case strings.Contains(key, "creation") || strings.Contains(key, "created") || strings.Contains(key, "registered"):
				if value != "" {
					result["creation_date"] = s.parseDate(value)
				}
			case strings.Contains(key, "expir") || strings.Contains(key, "expire"):
				if value != "" {
					result["expiration_date"] = s.parseDate(value)
				}
			case strings.Contains(key, "updated") || strings.Contains(key, "modified") || strings.Contains(key, "last"):
				if value != "" {
					result["updated_date"] = s.parseDate(value)
				}

			// 状态信息
			case strings.Contains(key, "status") || strings.Contains(key, "state"):
				if value != "" {
					statusList = append(statusList, value)
				}

			// 域名服务器
			case strings.Contains(key, "name server") || strings.Contains(key, "nserver") || key == "ns":
				if value != "" {
					nameServers = append(nameServers, value)
				}

			// 注册人信息
			case strings.Contains(key, "registrant") && strings.Contains(key, "name"):
				registrant["name"] = value
			case strings.Contains(key, "registrant") && strings.Contains(key, "org"):
				registrant["organization"] = value
			case strings.Contains(key, "registrant") && strings.Contains(key, "email"):
				registrant["email"] = value
			case strings.Contains(key, "registrant") && strings.Contains(key, "country"):
				registrant["country"] = value

			// 管理联系人信息
			case strings.Contains(key, "admin") && strings.Contains(key, "name"):
				adminContact["name"] = value
			case strings.Contains(key, "admin") && strings.Contains(key, "org"):
				adminContact["organization"] = value
			case strings.Contains(key, "admin") && strings.Contains(key, "email"):
				adminContact["email"] = value
			case strings.Contains(key, "admin") && strings.Contains(key, "country"):
				adminContact["country"] = value

			// 技术联系人信息
			case strings.Contains(key, "tech") && strings.Contains(key, "name"):
				techContact["name"] = value
			case strings.Contains(key, "tech") && strings.Contains(key, "org"):
				techContact["organization"] = value
			case strings.Contains(key, "tech") && strings.Contains(key, "email"):
				techContact["email"] = value
			case strings.Contains(key, "tech") && strings.Contains(key, "country"):
				techContact["country"] = value
			}
		}
	}

	// 设置解析结果
	if len(statusList) > 0 {
		result["status"] = statusList
	}
	if len(nameServers) > 0 {
		result["name_servers"] = nameServers
	}
	if len(registrant) > 0 {
		result["registrant"] = registrant
	}
	if len(adminContact) > 0 {
		result["admin_contact"] = adminContact
	}
	if len(techContact) > 0 {
		result["tech_contact"] = techContact
	}

	return result, nil
}

// parseDate 解析日期字符串
func (s *DomainService) parseDate(dateStr string) string {
	// 常见的日期格式
	formats := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02-Jan-2006",
		"2006.01.02",
		"01/02/2006",
		"02/01/2006",
		"2006/01/02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("2006-01-02T15:04:05Z")
		}
	}

	// 如果无法解析，返回原始字符串
	return dateStr
}

// GenerateWhoisReport 生成WHOIS报告
func (s *DomainService) GenerateWhoisReport(domainName string) ([]byte, error) {
	// 获取WHOIS信息
	whoisInfo, err := s.GetDomainWhois(domainName)
	if err != nil {
		return nil, fmt.Errorf("获取WHOIS信息失败: %w", err)
	}

	// TODO: 实现PDF报告生成
	// 这里可以使用PDF生成库如gofpdf或wkhtmltopdf

	// 模拟PDF内容
	reportContent := fmt.Sprintf(`WHOIS Report for %s

Domain Name: %s
Registrar: %v
Creation Date: %v
Expiration Date: %v
Updated Date: %v

Status: %v
Name Servers: %v

Registrant Information:
%v

Admin Contact:
%v

Tech Contact:
%v

Generated at: %s
`,
		domainName,
		whoisInfo["domain_name"],
		whoisInfo["registrar"],
		whoisInfo["creation_date"],
		whoisInfo["expiration_date"],
		whoisInfo["updated_date"],
		whoisInfo["status"],
		whoisInfo["name_servers"],
		whoisInfo["registrant"],
		whoisInfo["admin_contact"],
		whoisInfo["tech_contact"],
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return []byte(reportContent), nil
}

// BatchDeleteDomains 批量删除域名
func (s *DomainService) BatchDeleteDomains(domainIDs []uint, actorID uint, clientIP string) error {
	if len(domainIDs) == 0 {
		return errors.New("域名ID列表不能为空")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var deletedDomains []string

	for _, domainID := range domainIDs {
		// 获取域名信息
		domain, err := s.domainRepo.FindByID(domainID)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("域名ID %d 不存在", domainID)
			}
			return fmt.Errorf("获取域名信息失败: %w", err)
		}

		// 删除域名
		if err := s.domainRepo.Delete(domainID); err != nil {
			tx.Rollback()
			return fmt.Errorf("删除域名 %s 失败: %w", domain.Name, err)
		}

		deletedDomains = append(deletedDomains, domain.Name)

		// 记录变更日志
		changeLog := &dns.ChangeLog{
			ResourceType: "domain",
			ResourceID:   domainID,
			Action:       "delete",
			Description:  fmt.Sprintf("批量删除域名: %s", domain.Name),
			Status:       "success",
			ClientIP:     clientIP,
			TenantID:     domain.TenantID,
			ActorID:      actorID,
		}
		if err := s.changeLogRepo.Create(changeLog); err != nil {
			zap.L().Error("记录域名删除日志失败", zap.Error(err))
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	zap.L().Info("批量删除域名成功",
		zap.Strings("domains", deletedDomains),
		zap.Uint("actor_id", actorID),
		zap.String("client_ip", clientIP))

	return nil
}

// GenerateExcelTemplate 生成Excel导入模板
func (s *DomainService) GenerateExcelTemplate() ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			zap.L().Error("关闭Excel文件失败", zap.Error(err))
		}
	}()

	// 设置工作表名称
	sheetName := "域名导入模板"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("创建工作表失败: %w", err)
	}

	// 设置表头
	headers := []string{"域名", "状态", "注册商类型", "过期时间", "自动续费", "分组ID", "备注"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// 添加示例数据
	examples := [][]interface{}{
		{"example.com", "active", "godaddy", "2025-12-31", "true", "1", "示例域名"},
		{"test.com", "active", "dnspod", "2025-06-30", "false", "2", "测试域名"},
	}

	for rowIndex, example := range examples {
		for colIndex, value := range example {
			cell := fmt.Sprintf("%c%d", 'A'+colIndex, rowIndex+2)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// 设置列宽
	f.SetColWidth(sheetName, "A", "A", 20) // 域名
	f.SetColWidth(sheetName, "B", "B", 10) // 状态
	f.SetColWidth(sheetName, "C", "C", 15) // 注册商类型
	f.SetColWidth(sheetName, "D", "D", 15) // 过期时间
	f.SetColWidth(sheetName, "E", "E", 10) // 自动续费
	f.SetColWidth(sheetName, "F", "F", 10) // 分组ID
	f.SetColWidth(sheetName, "G", "G", 20) // 备注

	// 设置活动工作表
	f.SetActiveSheet(index)

	// 删除默认的Sheet1
	f.DeleteSheet("Sheet1")

	// 保存到字节数组
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("生成Excel文件失败: %w", err)
	}

	return buffer.Bytes(), nil
}
