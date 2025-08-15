import { api } from '@/utils/http'

/**
 * DNS批量操作和同步API
 */
export const dnsApi = {
  /**
   * 同步域名记录
   */
  syncDomainRecords: (data: {
    domain_id: number
    provider_id?: number
    dry_run?: boolean
  }): Promise<{
    success: boolean
    to_add: number
    to_update: number
    to_delete: number
    errors: string[]
  }> => {
    return api.post({ url: '/api/v1/dns/batch/sync', data })
  },

  /**
   * 批量同步多个域名
   */
  batchSyncDomains: (data: {
    domain_ids: number[]
    provider_id?: number
    dry_run?: boolean
  }): Promise<{
    success: boolean
    results: Array<{
      domain_id: number
      success: boolean
      to_add: number
      to_update: number
      to_delete: number
      errors: string[]
    }>
  }> => {
    return api.post({ url: '/api/v1/dns/batch/sync-domains', data })
  },

  /**
   * 获取同步状态
   */
  getSyncStatus: (domain_id?: number): Promise<{
    is_syncing: boolean
    last_sync_at?: string
    sync_progress?: number
    current_domain?: string
  }> => {
    const params = domain_id ? { domain_id } : {}
    return api.get({ url: '/api/v1/dns/sync/status', params })
  },

  /**
   * 取消同步
   */
  cancelSync: (domain_id?: number): Promise<{
    success: boolean
    message: string
  }> => {
    const data = domain_id ? { domain_id } : {}
    return api.post({ url: '/api/v1/dns/sync/cancel', data })
  },

  /**
   * 获取同步历史
   */
  getSyncHistory: (params: {
    domain_id?: number
    page?: number
    pageSize?: number
  }): Promise<{
    items: Array<{
      id: number
      domain_id: number
      domain_name: string
      provider_name: string
      sync_type: 'manual' | 'auto'
      status: 'success' | 'failed' | 'partial'
      records_added: number
      records_updated: number
      records_deleted: number
      error_message?: string
      started_at: string
      completed_at?: string
      duration?: number
    }>
    total: number
    page: number
    pageSize: number
  }> => {
    return api.get({ url: '/api/v1/dns/sync/history', params })
  },

  /**
   * 测试DNS提供商连接
   */
  testProviderConnection: (provider_id: number): Promise<{
    success: boolean
    message: string
    latency?: number
    zones_count?: number
  }> => {
    return api.post({ url: `/api/v1/dns/providers/${provider_id}/test` })
  },

  /**
   * 获取DNS统计信息
   */
  getStatistics: (domain_id?: number): Promise<{
    total_records: number
    records_by_type: Record<string, number>
    records_by_status: Record<string, number>
    sync_status: {
      synced: number
      pending: number
      failed: number
    }
    last_sync_time?: string
  }> => {
    const params = domain_id ? { domain_id } : {}
    return api.get({ url: '/api/v1/dns/statistics', params })
  }
}

// 重新导出其他API
export { recordApi } from './record'
export { domainApi } from './domain'
export { providerApi } from './provider'
