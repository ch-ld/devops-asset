package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"api-server/internal/model/dns"
	dnsrepo "api-server/internal/repository/dns"

	"go.uber.org/zap"
)

// NotificationService 通知服务
type NotificationService struct {
	notificationRepo *dnsrepo.NotificationRepository
	logger           *zap.Logger
	emailConfig      *EmailConfig
	webhookConfig    *WebhookConfig
	slackConfig      *SlackConfig
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPHost    string `json:"smtp_host"`
	SMTPPort    int    `json:"smtp_port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	FromAddress string `json:"from_address"`
	FromName    string `json:"from_name"`
	UseTLS      bool   `json:"use_tls"`
	UseStartTLS bool   `json:"use_starttls"`
}

// WebhookConfig Webhook配置
type WebhookConfig struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Timeout int               `json:"timeout"`
}

// SlackConfig Slack配置
type SlackConfig struct {
	WebhookURL string `json:"webhook_url"`
	Channel    string `json:"channel"`
	Username   string `json:"username"`
	IconEmoji  string `json:"icon_emoji"`
}

// NotificationTemplate 通知模板
type NotificationTemplate struct {
	Type        string `json:"type"`        // email, webhook, slack
	Subject     string `json:"subject"`     // 邮件主题或通知标题
	Body        string `json:"body"`        // 邮件正文或通知内容
	HTMLBody    string `json:"html_body"`   // HTML格式邮件正文
	Attachments []byte `json:"attachments"` // 附件
}

// NotificationData 通知数据
type NotificationData struct {
	Type         string                 `json:"type"` // domain_expiry, cert_expiry, https_down, dns_sync_failed
	Title        string                 `json:"title"`
	Message      string                 `json:"message"`
	Severity     string                 `json:"severity"` // info, warning, error, critical
	Domain       string                 `json:"domain"`
	URL          string                 `json:"url"`
	ExpiryDate   *time.Time             `json:"expiry_date"`
	DaysLeft     *int                   `json:"days_left"`
	ErrorDetails string                 `json:"error_details"`
	Metadata     map[string]interface{} `json:"metadata"`
	Timestamp    time.Time              `json:"timestamp"`
}

// NewNotificationService 创建通知服务
func NewNotificationService(
	notificationRepo *dnsrepo.NotificationRepository,
	emailConfig *EmailConfig,
	webhookConfig *WebhookConfig,
	slackConfig *SlackConfig,
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		logger:           zap.L().Named("notification-service"),
		emailConfig:      emailConfig,
		webhookConfig:    webhookConfig,
		slackConfig:      slackConfig,
	}
}

// SendNotification 发送通知
func (s *NotificationService) SendNotification(ctx context.Context, data *NotificationData, channels []string) error {
	s.logger.Info("Sending notification",
		zap.String("type", data.Type),
		zap.String("title", data.Title),
		zap.Strings("channels", channels))

	var errors []string

	for _, channel := range channels {
		var err error

		switch channel {
		case "email":
			err = s.sendEmailNotification(ctx, data)
		case "webhook":
			err = s.sendWebhookNotification(ctx, data)
		case "slack":
			err = s.sendSlackNotification(ctx, data)
		default:
			err = fmt.Errorf("unsupported notification channel: %s", channel)
		}

		if err != nil {
			s.logger.Error("Failed to send notification",
				zap.String("channel", channel),
				zap.Error(err))
			errors = append(errors, fmt.Sprintf("%s: %v", channel, err))
		}
	}

	// 记录通知历史
	now := time.Now()
	notification := &dns.Notification{
		Type:     data.Type,
		Title:    data.Title,
		Message:  data.Message,
		Severity: data.Severity,
		Channels: strings.Join(channels, ","),
		Status:   "sent",
		SentAt:   &now,
		Metadata: s.serializeMetadata(data.Metadata),
	}

	if len(errors) > 0 {
		notification.Status = "failed"
		notification.ErrorMsg = strings.Join(errors, "; ")
	}

	err := s.notificationRepo.Create(notification)
	if err != nil {
		s.logger.Error("Failed to save notification history", zap.Error(err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification failed on some channels: %s", strings.Join(errors, "; "))
	}

	return nil
}

// sendEmailNotification 发送邮件通知
func (s *NotificationService) sendEmailNotification(ctx context.Context, data *NotificationData) error {
	if s.emailConfig == nil {
		return fmt.Errorf("email configuration not provided")
	}

	template := s.generateEmailTemplate(data)

	// 构建邮件
	to := []string{} // TODO: 从配置或用户设置中获取收件人

	msg := s.buildEmailMessage(s.emailConfig.FromAddress, to, template.Subject, template.Body, template.HTMLBody)

	// 发送邮件
	auth := smtp.PlainAuth("", s.emailConfig.Username, s.emailConfig.Password, s.emailConfig.SMTPHost)
	addr := fmt.Sprintf("%s:%d", s.emailConfig.SMTPHost, s.emailConfig.SMTPPort)

	err := smtp.SendMail(addr, auth, s.emailConfig.FromAddress, to, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info("Email notification sent successfully",
		zap.String("type", data.Type),
		zap.Strings("to", to))

	return nil
}

// sendWebhookNotification 发送Webhook通知
func (s *NotificationService) sendWebhookNotification(ctx context.Context, data *NotificationData) error {
	if s.webhookConfig == nil {
		return fmt.Errorf("webhook configuration not provided")
	}

	// 序列化数据
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, s.webhookConfig.Method, s.webhookConfig.URL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	for key, value := range s.webhookConfig.Headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	client := &http.Client{
		Timeout: time.Duration(s.webhookConfig.Timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned error status: %d", resp.StatusCode)
	}

	s.logger.Info("Webhook notification sent successfully",
		zap.String("type", data.Type),
		zap.String("url", s.webhookConfig.URL))

	return nil
}

// sendSlackNotification 发送Slack通知
func (s *NotificationService) sendSlackNotification(ctx context.Context, data *NotificationData) error {
	if s.slackConfig == nil {
		return fmt.Errorf("slack configuration not provided")
	}

	// 构建Slack消息
	message := s.buildSlackMessage(data)

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal slack payload: %w", err)
	}

	// 发送到Slack
	req, err := http.NewRequestWithContext(ctx, "POST", s.slackConfig.WebhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create slack request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("slack returned error status: %d", resp.StatusCode)
	}

	s.logger.Info("Slack notification sent successfully",
		zap.String("type", data.Type),
		zap.String("channel", s.slackConfig.Channel))

	return nil
}

// generateEmailTemplate 生成邮件模板
func (s *NotificationService) generateEmailTemplate(data *NotificationData) *NotificationTemplate {
	var subject, body, htmlBody string

	switch data.Type {
	case "domain_expiry":
		subject = fmt.Sprintf("域名即将过期提醒 - %s", data.Domain)
		body = fmt.Sprintf("您的域名 %s 将在 %d 天后过期，请及时续费。", data.Domain, *data.DaysLeft)
		htmlBody = fmt.Sprintf(`
			<h2>域名过期提醒</h2>
			<p>您的域名 <strong>%s</strong> 将在 <strong>%d</strong> 天后过期。</p>
			<p>过期时间：%s</p>
			<p>请及时续费以避免服务中断。</p>
		`, data.Domain, *data.DaysLeft, data.ExpiryDate.Format("2006-01-02 15:04:05"))

	case "cert_expiry":
		subject = fmt.Sprintf("SSL证书即将过期提醒 - %s", data.Domain)
		body = fmt.Sprintf("您的SSL证书 %s 将在 %d 天后过期，请及时更新。", data.Domain, *data.DaysLeft)
		htmlBody = fmt.Sprintf(`
			<h2>SSL证书过期提醒</h2>
			<p>您的SSL证书 <strong>%s</strong> 将在 <strong>%d</strong> 天后过期。</p>
			<p>过期时间：%s</p>
			<p>请及时更新证书以确保网站安全。</p>
		`, data.Domain, *data.DaysLeft, data.ExpiryDate.Format("2006-01-02 15:04:05"))

	case "https_down":
		subject = fmt.Sprintf("HTTPS服务异常 - %s", data.Domain)
		body = fmt.Sprintf("您的HTTPS服务 %s 出现异常：%s", data.URL, data.ErrorDetails)
		htmlBody = fmt.Sprintf(`
			<h2>HTTPS服务异常</h2>
			<p>您的HTTPS服务出现异常：</p>
			<p><strong>URL:</strong> %s</p>
			<p><strong>错误信息:</strong> %s</p>
			<p><strong>检测时间:</strong> %s</p>
		`, data.URL, data.ErrorDetails, data.Timestamp.Format("2006-01-02 15:04:05"))

	default:
		subject = data.Title
		body = data.Message
		htmlBody = fmt.Sprintf("<p>%s</p>", data.Message)
	}

	return &NotificationTemplate{
		Type:     "email",
		Subject:  subject,
		Body:     body,
		HTMLBody: htmlBody,
	}
}

// buildEmailMessage 构建邮件消息
func (s *NotificationService) buildEmailMessage(from string, to []string, subject, body, htmlBody string) string {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ",")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")

	if htmlBody != "" {
		msg.WriteString("Content-Type: multipart/alternative; boundary=\"boundary123\"\r\n\r\n")
		msg.WriteString("--boundary123\r\n")
		msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n")
		msg.WriteString(body)
		msg.WriteString("\r\n--boundary123\r\n")
		msg.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n\r\n")
		msg.WriteString(htmlBody)
		msg.WriteString("\r\n--boundary123--\r\n")
	} else {
		msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n")
		msg.WriteString(body)
	}

	return msg.String()
}

// buildSlackMessage 构建Slack消息
func (s *NotificationService) buildSlackMessage(data *NotificationData) map[string]interface{} {
	color := "good"
	switch data.Severity {
	case "warning":
		color = "warning"
	case "error", "critical":
		color = "danger"
	}

	message := map[string]interface{}{
		"channel":    s.slackConfig.Channel,
		"username":   s.slackConfig.Username,
		"icon_emoji": s.slackConfig.IconEmoji,
		"attachments": []map[string]interface{}{
			{
				"color":     color,
				"title":     data.Title,
				"text":      data.Message,
				"timestamp": data.Timestamp.Unix(),
				"fields": []map[string]interface{}{
					{
						"title": "域名",
						"value": data.Domain,
						"short": true,
					},
					{
						"title": "严重程度",
						"value": data.Severity,
						"short": true,
					},
				},
			},
		},
	}

	if data.URL != "" {
		attachment := message["attachments"].([]map[string]interface{})[0]
		fields := attachment["fields"].([]map[string]interface{})
		fields = append(fields, map[string]interface{}{
			"title": "URL",
			"value": data.URL,
			"short": false,
		})
		attachment["fields"] = fields
	}

	return message
}

// serializeMetadata 序列化元数据
func (s *NotificationService) serializeMetadata(metadata map[string]interface{}) string {
	if metadata == nil {
		return ""
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		s.logger.Error("Failed to serialize metadata", zap.Error(err))
		return ""
	}

	return string(data)
}
