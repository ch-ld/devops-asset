package dns

import (
	"api-server/internal/model/dns"
	repo "api-server/internal/repository/dns"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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

	// 状态统计
	statusCounts, err := s.domainRepo.CountByStatus()
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
