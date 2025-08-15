/**
 * DNS模块统一组件导出
 * 提供DNS管理相关的复用组件
 */

// 基础组件
export { default as DNSStatusTag } from './base/DNSStatusTag.vue'
export { default as DNSTable } from './base/DNSTable.vue' 
export { default as DNSStatCard } from './base/DNSStatCard.vue'
export { default as DNSSearchForm } from './base/DNSSearchForm.vue'

// 表单组件
export { default as DomainModal } from './forms/DomainModal.vue'
export { default as RecordModal } from './forms/RecordModal.vue'
export { default as CertificateModal } from './forms/CertificateModal.vue'
export { default as ProviderModal } from './forms/ProviderModal.vue'
export { default as GroupModal } from './forms/GroupModal.vue'
export { default as TagModal } from './forms/TagModal.vue'

// 业务组件
export { default as DomainDetail } from './business/DomainDetail.vue'
export { default as CertificateDetail } from './business/CertificateDetail.vue'
export { default as RecordDiff } from './business/RecordDiff.vue'
export { default as MonitorChart } from './business/MonitorChart.vue'
