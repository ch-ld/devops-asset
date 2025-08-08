/**
 * DNS模块相关类型定义
 */

// 基础分页参数
export interface BaseListParams {
  page?: number
  size?: number
  keyword?: string
}

// 基础列表响应
export interface BaseListResponse<T> {
  total: number
  list: T[]
  page?: number
  page_size?: number
  total_page?: number
}

// 域名标签
export interface DomainTag {
  id: number
  name: string
  color: string
  description?: string
  tenant_id: number
  created_by: number
  updated_by: number
  created_at: string
  updated_at: string
}

// 域名分组
export interface DomainGroup {
  id: number
  name: string
  parent_id?: number
  parent?: DomainGroup
  children?: DomainGroup[]
  description?: string
  sort: number
  tenant_id: number
  created_by: number
  updated_by: number
  created_at: string
  updated_at: string
}

// 域名
export interface Domain {
  id: number
  name: string
  status: 'active' | 'inactive' | 'expired'
  registrar_type: string
  registrar_id?: number
  expires_at?: string
  auto_renew: boolean
  group_id?: number
  group?: DomainGroup
  tags?: DomainTag[]
  configuration: Record<string, any>
  remark?: string
  tenant_id: number
  created_by: number
  updated_by: number
  created_at: string
  updated_at: string
}

// DNS记录
export interface DNSRecord {
  id: number
  domain_id: number
  domain?: Domain
  provider_id: number
  provider?: DNSProvider
  name: string
  type: string
  value: string
  ttl: number
  priority?: number
  weight?: number
  port?: number
  status: 'active' | 'inactive' | 'syncing' | 'error'
  cloud_record_id?: string
  sync_status: 'synced' | 'pending' | 'failed'
  last_sync_at?: string
  configuration: Record<string, any>
  remark?: string
  tenant_id: number
  created_by: number
  updated_by: number
  created_at: string
  updated_at: string
}

// DNS记录类型别名（用于简化引用）
export type DnsRecord = DNSRecord

// DNS提供商
export interface DNSProvider {
  id: number
  name: string
  type: string
  status: 'active' | 'inactive'
  is_default: boolean
  priority: number
  rate_limit: number
  concurrent: number
  timeout: number
  last_test_at?: string
  test_result?: string
  error_message?: string
  remark?: string
  tenant_id: number
  created_by: number
  updated_by: number
  created_at: string
  updated_at: string
}

// 证书
export interface Certificate {
  id: number
  domain_id: number
  domain?: Domain
  common_name: string
  subject_alt_names: string[]
  ca_type: string
  status: 'pending' | 'issued' | 'expired' | 'revoked'
  serial_number?: string
  fingerprint?: string
  issued_at?: string
  expires_at?: string
  auto_renew: boolean
  renew_days: number
  last_renew_at?: string
  remark?: string
  tenant_id: number
  created_by: number
  updated_by: number
  created_at: string
  updated_at: string
}

// 域名列表查询参数
export interface DomainListParams extends BaseListParams {
  status?: string
  group_id?: number
  registrar_type?: string
}

// 域名列表响应
export interface DomainListResponse extends BaseListResponse<Domain> {}

// 域名分组列表查询参数
export interface DomainGroupListParams extends BaseListParams {
  parent_id?: number
}

// 域名分组列表响应
export interface DomainGroupListResponse extends BaseListResponse<DomainGroup> {}

// DNS记录列表查询参数
export interface DNSRecordListParams extends BaseListParams {
  domain_id?: number
  type?: string
  status?: string
  provider_id?: number
}

// DNS记录列表响应
export interface DNSRecordListResponse extends BaseListResponse<DNSRecord> {}

// DNS提供商列表查询参数
export interface DNSProviderListParams extends BaseListParams {
  type?: string
  status?: string
}

// DNS提供商列表响应
export interface DNSProviderListResponse extends BaseListResponse<DNSProvider> {}

// 证书列表查询参数
export interface CertificateListParams extends BaseListParams {
  domain_id?: number
  status?: string
  ca_type?: string
}

// 证书列表响应
export interface CertificateListResponse extends BaseListResponse<Certificate> {}

// 统计信息
export interface StatisticsResponse {
  total: number
  by_status: Record<string, number>
  by_type?: Record<string, number>
  expiring?: number
  additional?: Record<string, any>
}

// 批量操作请求
export interface BatchOperationRequest {
  ids: number[]
  action: string
  data?: Record<string, any>
}

// 批量操作响应
export interface BatchOperationResponse {
  success: number
  failed: number
  total: number
  failed_items?: string[]
}

// 导入请求
export interface ImportRequest {
  data: any[]
  format: string
  options?: Record<string, any>
}

// 导入响应
export interface ImportResponse {
  success: number
  failed: number
  total: number
  failed_items?: string[]
}

// 导出请求
export interface ExportRequest {
  format: string
  filters?: Record<string, any>
  fields?: string[]
}

// 同步请求
export interface SyncRequest {
  provider_id?: number
  domain_id?: number
  sync_type: 'full' | 'incremental'
  force?: boolean
}

// 同步响应
export interface SyncResponse {
  success: boolean
  total_records: number
  added_records: number
  updated_records: number
  deleted_records: number
  failed_records: number
  duration: string
  errors?: string[]
}

// 测试连接请求
export interface TestConnectionRequest {
  type: string
  credentials: Record<string, string>
  endpoint?: string
  region?: string
}

// 测试连接响应
export interface TestConnectionResponse {
  success: boolean
  latency: string
  error_msg?: string
  details?: Record<string, string>
  tested_at: string
  test_type: string
  endpoint: string
  status_code: number
}

// HTTPS监控
export interface HTTPSMonitor {
  id: number
  name: string
  url: string
  check_interval: number
  timeout: number
  status: 'active' | 'inactive'
  last_status: 'online' | 'offline' | 'warning'
  last_checked?: string
  last_response_time?: number
  ssl_cert_expires_at?: string
  ssl_cert_issuer?: string
  ssl_cert_subject?: string
  notification_enabled: boolean
  notification_channels: string[]
  alert_threshold: number
  consecutive_failures: number
  max_failures: number
  configuration: Record<string, any>
  remark?: string
  tenant_id: number
  created_by: number
  updated_by: number
  created_at: string
  updated_at: string
}

// 通知
export interface Notification {
  id: number
  type: 'certificate_expiring' | 'certificate_expired' | 'monitor_down' | 'monitor_up' | 'system'
  title: string
  message: string
  severity: 'info' | 'warning' | 'error' | 'critical'
  status: 'unread' | 'read'
  resource_type?: string
  resource_id?: number
  channels: string[]
  sent_at?: string
  read_at?: string
  metadata: Record<string, any>
  tenant_id: number
  created_by: number
  created_at: string
  updated_at: string
}

// HTTPS监控列表查询参数
export interface HTTPSMonitorListParams extends BaseListParams {
  status?: string
  last_status?: string
  url?: string
}

// 通知列表查询参数
export interface NotificationListParams extends BaseListParams {
  type?: string
  severity?: string
  status?: string
  resource_type?: string
}

// 创建请求类型
export interface DomainCreateRequest {
  name: string
  registrar_type: string
  registrar_id?: number
  expires_at?: string
  auto_renew?: boolean
  group_id?: number
  tag_ids?: number[]
  configuration?: Record<string, any>
  remark?: string
}

export interface DNSProviderCreateRequest {
  name: string
  type: string
  priority?: number
  rate_limit?: number
  concurrent?: number
  timeout?: number
  credentials: Record<string, string>
  remark?: string
}

export interface CertificateCreateRequest {
  domain_id: number
  common_name: string
  subject_alt_names?: string[]
  ca_type: string
  auto_renew?: boolean
  renew_days?: number
  remark?: string
}

export interface HTTPSMonitorCreateRequest {
  name: string
  url: string
  check_interval?: number
  timeout?: number
  notification_enabled?: boolean
  notification_channels?: string[]
  alert_threshold?: number
  max_failures?: number
  configuration?: Record<string, any>
  remark?: string
}

export interface NotificationCreateRequest {
  type: string
  title: string
  message: string
  severity: string
  resource_type?: string
  resource_id?: number
  channels: string[]
  metadata?: Record<string, any>
}
