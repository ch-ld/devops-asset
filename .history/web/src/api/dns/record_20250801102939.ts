import { request } from '@/utils/http'
import type { 
  DNSRecord, 
  DNSRecordListParams, 
  DNSRecordListResponse,
  SyncRequest,
  SyncResponse,
  BatchOperationRequest,
  BatchOperationResponse
} from '@/types/dns'

/**
 * DNS记录管理API
 */
export const recordApi = {
  /**
   * 获取DNS记录列表
   */
  list: (params: DNSRecordListParams): Promise<DNSRecordListResponse> => {
    return request.get({ url: '/dns/records', params })
  },

  /**
   * 获取DNS记录详情
   */
  get: (id: number): Promise<DNSRecord> => {
    return request.get({ url: `/dns/records/${id}` })
  },

  /**
   * 创建DNS记录
   */
  create: (data: Partial<DNSRecord>): Promise<DNSRecord> => {
    return request.post({ url: '/dns/records', data })
  },

  /**
   * 更新DNS记录
   */
  update: (id: number, data: Partial<DNSRecord>): Promise<DNSRecord> => {
    return request.put({ url: `/dns/records/${id}`, data })
  },

  /**
   * 删除DNS记录
   */
  delete: (id: number): Promise<void> => {
    return request.del({ url: `/dns/records/${id}` })
  },

  /**
   * 批量删除DNS记录
   */
  batchDelete: (ids: (string | number)[]): Promise<BatchOperationResponse> => {
    return request.post({ 
      url: '/dns/batch/records', 
      data: {
        action: 'delete',
        ids
      }
    })
  },

  /**
   * 批量操作DNS记录
   */
  batchOperate: (data: BatchOperationRequest): Promise<BatchOperationResponse> => {
    return request.post({ url: '/dns/batch/records', data })
  },

  /**
   * 批量同步DNS记录
   */
  batchSync: (data: { domain_id: number; dry_run?: boolean }): Promise<{
    success: boolean
    message: string
    data?: {
      added: number
      updated: number
      deleted: number
      failed: number
    }
  }> => {
    return request.post({ url: '/dns/domains/sync', data })
  },

  /**
   * 同步DNS记录
   */
  sync: (data: SyncRequest): Promise<SyncResponse> => {
    return request.post({ url: '/dns/sync/records', data })
  },

  /**
   * 获取同步状态
   */
  getSyncStatus: (): Promise<{
    is_syncing: boolean
    last_sync_at?: string
    sync_progress?: number
  }> => {
    return request.get({ url: '/dns/sync/status' })
  },

  /**
   * 测试DNS解析
   */
  testResolution: (data: {
    name: string
    type: string
    server?: string
  }): Promise<{
    success: boolean
    result?: any[]
    error_msg?: string
    query_time: number
  }> => {
    return request.post({ url: '/dns/test/dns-resolution', data })
  },

  /**
   * 获取DNS记录统计信息
   */
  statistics: (): Promise<{
    total: number
    by_status: Record<string, number>
    by_type: Record<string, number>
  }> => {
    return request.get({ url: '/dns/stats/records' })
  },

  /**
   * 导入DNS记录
   */
  import: (data: {
    data: any[]
    format: string
    options?: Record<string, any>
  }): Promise<{
    success: number
    failed: number
    total: number
    failed_items: string[]
  }> => {
    return request.post({ url: '/dns/import-export/records/import', data })
  },

  /**
   * 导出DNS记录
   */
  export: (params: {
    format: string
    filters?: Record<string, any>
    fields?: string[]
  }): Promise<Blob> => {
    return request.post({
      url: '/dns/import-export/records/export',
      data: params,
      responseType: 'blob'
    })
  }
}
