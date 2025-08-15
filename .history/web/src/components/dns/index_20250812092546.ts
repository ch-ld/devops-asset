/**
 * DNS模块统一组件导出
 * 提供DNS管理相关的复用组件
 */

// 布局组件
export { default as DNSPageWrapper } from './layout/DNSPageWrapper.vue'

// 基础组件
export { default as DNSStatusTag } from './base/DNSStatusTag.vue'
export { default as DNSTable } from './base/DNSTable.vue' 
export { default as DNSStatCard } from './base/DNSStatCard.vue'
export { default as DNSSearchForm } from './base/DNSSearchForm.vue'

// 表单组件 - 引用实际存在的文件
export { default as DomainModal } from '../../views/dns/domains/components/DomainModal.vue'
export { default as RecordModal } from '../../views/dns/records/components/RecordModal.vue'
export { default as CertificateModal } from '../../views/dns/certs/components/CertificateModal.vue'
export { default as ProviderModal } from '../../views/dns/providers/components/ProviderModal.vue'
export { default as GroupModal } from '../../views/dns/domains/components/GroupModal.vue'
export { default as TagModal } from '../../views/dns/domains/components/TagModal.vue'
