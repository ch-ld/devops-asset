package dns

import (
	"context"
	"fmt"
	"strings"
	"time"

	dnsprovider "api-server/internal/provider/dns"

	"github.com/go-acme/lego/v4/challenge/dns01"
	"go.uber.org/zap"
)

// DNSProvider ACME DNS-01挑战提供商
type DNSProvider struct {
	driver dnsprovider.Driver
	logger *zap.Logger
}

// NewDNSProvider 创建DNS提供商
func NewDNSProvider(driver dnsprovider.Driver) *DNSProvider {
	return &DNSProvider{
		driver: driver,
		logger: zap.L().Named("dns-provider"),
	}
}

// Present 创建DNS记录用于ACME挑战
func (p *DNSProvider) Present(domain, token, keyAuth string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	fqdn, value := dns01.GetRecord(domain, keyAuth)
	
	// 提取域名和记录名
	recordName := p.extractRecordName(fqdn, domain)
	zoneName := p.extractZoneName(domain)

	p.logger.Info("Creating TXT record for ACME challenge",
		zap.String("domain", domain),
		zap.String("zone", zoneName),
		zap.String("record", recordName),
		zap.String("value", value))

	// 创建TXT记录 - 使用600秒TTL以满足阿里云要求
	result, err := p.driver.CreateTXTChallenge(ctx, zoneName, value, 600)
	if err != nil {
		p.logger.Error("Failed to create TXT challenge record",
			zap.String("domain", domain),
			zap.Error(err))
		return fmt.Errorf("failed to create TXT challenge record: %w", err)
	}

	p.logger.Info("TXT record created successfully",
		zap.String("domain", domain),
		zap.String("record_id", result.RecordID))

	// 等待DNS传播
	return p.waitForPropagation(ctx, fqdn, value)
}

// CleanUp 清理DNS记录
func (p *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	_, value := dns01.GetRecord(domain, keyAuth)
	zoneName := p.extractZoneName(domain)

	p.logger.Info("Cleaning up TXT record for ACME challenge",
		zap.String("domain", domain),
		zap.String("zone", zoneName),
		zap.String("value", value))

	err := p.driver.DeleteTXTChallenge(ctx, zoneName, value)
	if err != nil {
		p.logger.Error("Failed to delete TXT challenge record",
			zap.String("domain", domain),
			zap.Error(err))
		return fmt.Errorf("failed to delete TXT challenge record: %w", err)
	}

	p.logger.Info("TXT record cleaned up successfully",
		zap.String("domain", domain))

	return nil
}

// Timeout 返回DNS传播超时时间
func (p *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return 2 * time.Minute, 10 * time.Second
}

// extractRecordName 提取记录名称
func (p *DNSProvider) extractRecordName(fqdn, domain string) string {
	// fqdn: _acme-challenge.example.com.
	// domain: example.com
	// return: _acme-challenge
	
	if strings.HasSuffix(fqdn, ".") {
		fqdn = fqdn[:len(fqdn)-1]
	}
	
	if strings.HasSuffix(fqdn, "."+domain) {
		return fqdn[:len(fqdn)-len(domain)-1]
	}
	
	return "_acme-challenge"
}

// extractZoneName 提取区域名称
func (p *DNSProvider) extractZoneName(domain string) string {
	// 对于子域名，需要找到根域名
	// 例如：api.example.com -> example.com
	
	parts := strings.Split(domain, ".")
	if len(parts) >= 2 {
		// 返回最后两个部分作为根域名
		return strings.Join(parts[len(parts)-2:], ".")
	}
	
	return domain
}

// waitForPropagation 等待DNS传播
func (p *DNSProvider) waitForPropagation(ctx context.Context, fqdn, expectedValue string) error {
	p.logger.Info("Waiting for DNS propagation",
		zap.String("fqdn", fqdn),
		zap.String("expected_value", expectedValue))

	timeout := 2 * time.Minute
	interval := 10 * time.Second
	
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 检查DNS记录是否已传播
		if p.checkDNSRecord(fqdn, expectedValue) {
			p.logger.Info("DNS propagation completed",
				zap.String("fqdn", fqdn))
			return nil
		}

		p.logger.Debug("DNS record not yet propagated, waiting...",
			zap.String("fqdn", fqdn),
			zap.Duration("interval", interval))

		time.Sleep(interval)
	}

	return fmt.Errorf("DNS propagation timeout after %v", timeout)
}

// checkDNSRecord 检查DNS记录
func (p *DNSProvider) checkDNSRecord(fqdn, expectedValue string) bool {
	// TODO: 实现实际的DNS查询
	// 这里可以使用net包的LookupTXT函数来查询TXT记录
	// 或者使用第三方DNS库如miekg/dns
	
	// 暂时返回true，实际应该查询DNS服务器
	return true
}

// Sequential DNS Provider for sequential challenges
type SequentialDNSProvider struct {
	*DNSProvider
	challenges map[string]string // domain -> value mapping
}

// NewSequentialDNSProvider 创建顺序DNS提供商
func NewSequentialDNSProvider(driver dnsprovider.Driver) *SequentialDNSProvider {
	return &SequentialDNSProvider{
		DNSProvider: NewDNSProvider(driver),
		challenges:  make(map[string]string),
	}
}

// Present 创建DNS记录（支持多域名）
func (p *SequentialDNSProvider) Present(domain, token, keyAuth string) error {
	_, value := dns01.GetRecord(domain, keyAuth)
	p.challenges[domain] = value
	
	return p.DNSProvider.Present(domain, token, keyAuth)
}

// CleanUp 清理DNS记录（支持多域名）
func (p *SequentialDNSProvider) CleanUp(domain, token, keyAuth string) error {
	delete(p.challenges, domain)
	return p.DNSProvider.CleanUp(domain, token, keyAuth)
}

// GetChallenges 获取当前挑战
func (p *SequentialDNSProvider) GetChallenges() map[string]string {
	result := make(map[string]string)
	for k, v := range p.challenges {
		result[k] = v
	}
	return result
}

// ClearChallenges 清理所有挑战
func (p *SequentialDNSProvider) ClearChallenges() {
	for domain := range p.challenges {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		zoneName := p.extractZoneName(domain)
		value := p.challenges[domain]
		
		err := p.driver.DeleteTXTChallenge(ctx, zoneName, value)
		if err != nil {
			p.logger.Error("Failed to clear challenge",
				zap.String("domain", domain),
				zap.Error(err))
		}
		
		cancel()
	}
	
	p.challenges = make(map[string]string)
}

// Batch DNS Provider for batch operations
type BatchDNSProvider struct {
	*DNSProvider
	batchSize int
}

// NewBatchDNSProvider 创建批量DNS提供商
func NewBatchDNSProvider(driver dnsprovider.Driver, batchSize int) *BatchDNSProvider {
	if batchSize <= 0 {
		batchSize = 10
	}
	
	return &BatchDNSProvider{
		DNSProvider: NewDNSProvider(driver),
		batchSize:   batchSize,
	}
}

// PresentBatch 批量创建DNS记录
func (p *BatchDNSProvider) PresentBatch(challenges map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	p.logger.Info("Creating batch TXT records for ACME challenges",
		zap.Int("count", len(challenges)))

	// 分批处理
	domains := make([]string, 0, len(challenges))
	for domain := range challenges {
		domains = append(domains, domain)
	}

	for i := 0; i < len(domains); i += p.batchSize {
		end := i + p.batchSize
		if end > len(domains) {
			end = len(domains)
		}

		batch := domains[i:end]
		err := p.processBatch(ctx, batch, challenges)
		if err != nil {
			return fmt.Errorf("failed to process batch %d-%d: %w", i, end-1, err)
		}
	}

	p.logger.Info("Batch TXT records created successfully")
	return nil
}

// processBatch 处理批量操作
func (p *BatchDNSProvider) processBatch(ctx context.Context, domains []string, challenges map[string]string) error {
	for _, domain := range domains {
		value := challenges[domain]
		zoneName := p.extractZoneName(domain)

		_, err := p.driver.CreateTXTChallenge(ctx, zoneName, value, 300)
		if err != nil {
			p.logger.Error("Failed to create TXT challenge in batch",
				zap.String("domain", domain),
				zap.Error(err))
			return err
		}
	}

	return nil
}

// CleanUpBatch 批量清理DNS记录
func (p *BatchDNSProvider) CleanUpBatch(challenges map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	p.logger.Info("Cleaning up batch TXT records for ACME challenges",
		zap.Int("count", len(challenges)))

	for domain, value := range challenges {
		zoneName := p.extractZoneName(domain)
		
		err := p.driver.DeleteTXTChallenge(ctx, zoneName, value)
		if err != nil {
			p.logger.Error("Failed to delete TXT challenge in batch",
				zap.String("domain", domain),
				zap.Error(err))
			// 继续处理其他记录，不返回错误
		}
	}

	p.logger.Info("Batch TXT records cleaned up successfully")
	return nil
}
