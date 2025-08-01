package dns

import (
	"api-server/internal/model/dns"
	"fmt"

	"gorm.io/gorm"
)

// CertificateRepository 证书数据访问层
type CertificateRepository struct {
	DB *gorm.DB
}

// NewCertificateRepository 创建证书仓库实例
func NewCertificateRepository(db *gorm.DB) *CertificateRepository {
	return &CertificateRepository{DB: db}
}

// Create 创建证书
func (r *CertificateRepository) Create(cert *dns.Certificate) error {
	return r.DB.Create(cert).Error
}

// Update 更新证书
func (r *CertificateRepository) Update(cert *dns.Certificate) error {
	return r.DB.Model(cert).Updates(cert).Error
}

// Delete 删除证书
func (r *CertificateRepository) Delete(id uint) error {
	return r.DB.Delete(&dns.Certificate{}, id).Error
}

// FindByID 根据ID查找证书
func (r *CertificateRepository) FindByID(id uint) (*dns.Certificate, error) {
	var cert dns.Certificate
	err := r.DB.Preload("Domain").First(&cert, id).Error
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// FindByDomainID 根据域名ID查找证书
func (r *CertificateRepository) FindByDomainID(domainID uint) ([]*dns.Certificate, error) {
	var certs []*dns.Certificate
	err := r.DB.Preload("Domain").Where("domain_id = ?", domainID).Find(&certs).Error
	return certs, err
}

// FindByCommonName 根据主域名查找证书
func (r *CertificateRepository) FindByCommonName(commonName string) ([]*dns.Certificate, error) {
	var certs []*dns.Certificate
	err := r.DB.Preload("Domain").Where("common_name = ?", commonName).Find(&certs).Error
	return certs, err
}

// FindByStatus 根据状态查找证书
func (r *CertificateRepository) FindByStatus(status string) ([]*dns.Certificate, error) {
	var certs []*dns.Certificate
	err := r.DB.Preload("Domain").Where("status = ?", status).Find(&certs).Error
	return certs, err
}

// FindByCAType 根据CA类型查找证书
func (r *CertificateRepository) FindByCAType(caType string) ([]*dns.Certificate, error) {
	var certs []*dns.Certificate
	err := r.DB.Preload("Domain").Where("ca_type = ?", caType).Find(&certs).Error
	return certs, err
}

// FindExpiring 查找即将过期的证书
func (r *CertificateRepository) FindExpiring(days int) ([]*dns.Certificate, error) {
	var certs []*dns.Certificate
	err := r.DB.Preload("Domain").
		Where("expires_at IS NOT NULL AND expires_at <= DATE_ADD(NOW(), INTERVAL ? DAY)", days).
		Find(&certs).Error
	return certs, err
}

// FindAutoRenew 查找需要自动续期的证书
func (r *CertificateRepository) FindAutoRenew() ([]*dns.Certificate, error) {
	var certs []*dns.Certificate
	err := r.DB.Preload("Domain").
		Where("auto_renew = ? AND status = ? AND expires_at <= DATE_ADD(NOW(), INTERVAL renew_days DAY)",
			true, "issued").
		Find(&certs).Error
	return certs, err
}

// FindByTenantID 根据租户ID查找证书
func (r *CertificateRepository) FindByTenantID(tenantID uint) ([]*dns.Certificate, error) {
	var certs []*dns.Certificate
	err := r.DB.Preload("Domain").Where("tenant_id = ?", tenantID).Find(&certs).Error
	return certs, err
}

// Search 搜索证书
func (r *CertificateRepository) Search(keyword string, limit, offset int) ([]*dns.Certificate, int64, error) {
	var certs []*dns.Certificate
	var total int64

	query := r.DB.Preload("Domain").Model(&dns.Certificate{})

	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("common_name LIKE ? OR serial_number LIKE ? OR remark LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&certs).Error
	return certs, total, err
}

// SearchWithFilters 带过滤条件的搜索
func (r *CertificateRepository) SearchWithFilters(filters map[string]interface{}, limit, offset int) ([]*dns.Certificate, int64, error) {
	var certs []*dns.Certificate
	var total int64

	query := r.DB.Preload("Domain").Model(&dns.Certificate{})

	// 应用过滤条件
	for key, value := range filters {
		if value != nil && value != "" {
			switch key {
			case "keyword":
				searchPattern := "%" + fmt.Sprintf("%v", value) + "%"
				query = query.Where("common_name LIKE ? OR serial_number LIKE ? OR remark LIKE ?",
					searchPattern, searchPattern, searchPattern)
			case "status":
				query = query.Where("status = ?", value)
			case "ca_type":
				query = query.Where("ca_type = ?", value)
			case "domain_id":
				query = query.Where("domain_id = ?", value)
			case "auto_renew":
				query = query.Where("auto_renew = ?", value)
			case "tenant_id":
				query = query.Where("tenant_id = ?", value)
			}
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&certs).Error
	return certs, total, err
}

// CountByStatus 统计各状态证书数量
func (r *CertificateRepository) CountByStatus() (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	err := r.DB.Model(&dns.Certificate{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.Status] = result.Count
	}

	return counts, nil
}

// CountByCAType 统计各CA类型证书数量
func (r *CertificateRepository) CountByCAType() (map[string]int64, error) {
	var results []struct {
		CAType string `gorm:"column:ca_type"`
		Count  int64
	}

	err := r.DB.Model(&dns.Certificate{}).
		Select("ca_type, COUNT(*) as count").
		Group("ca_type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.CAType] = result.Count
	}

	return counts, nil
}

// BatchCreate 批量创建证书
func (r *CertificateRepository) BatchCreate(certs []*dns.Certificate) error {
	return r.DB.Create(&certs).Error
}

// BatchUpdate 批量更新证书
func (r *CertificateRepository) BatchUpdate(certs []*dns.Certificate) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		for _, cert := range certs {
			if err := tx.Model(cert).Updates(cert).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchDelete 批量删除证书
func (r *CertificateRepository) BatchDelete(ids []uint) error {
	return r.DB.Delete(&dns.Certificate{}, ids).Error
}

// ExistsByCommonName 检查证书是否存在
func (r *CertificateRepository) ExistsByCommonName(commonName string) (bool, error) {
	var count int64
	err := r.DB.Model(&dns.Certificate{}).Where("common_name = ?", commonName).Count(&count).Error
	return count > 0, err
}

// CountByTenantID 统计租户证书数量
func (r *CertificateRepository) CountByTenantID(tenantID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&dns.Certificate{}).Where("tenant_id = ?", tenantID).Count(&count).Error
	return count, err
}

// UpdateStatus 更新证书状�?func (r *CertificateRepository) UpdateStatus(id uint, status string) error {
	return r.DB.Model(&dns.Certificate{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateCertificateData 更新证书数据（加密存储）
func (r *CertificateRepository) UpdateCertificateData(id uint, certEnc, keyEnc, chainEnc string) error {
	return r.DB.Model(&dns.Certificate{}).Where("id = ?", id).Updates(map[string]interface{}{
		"certificate_enc": certEnc,
		"private_key_enc": keyEnc,
		"chain_enc":       chainEnc,
	}).Error
}
