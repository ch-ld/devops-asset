import { request } from '@/api/client'
import type { Domain, DomainListParams, DomainListResponse } from '@/types/dns'

/**
 * 域名管理API
 */
export const domainApi = {
  /**
   * 获取域名列表
   */
  list: (params: DomainListParams): Promise<DomainListResponse> => {
    return request.get({ url: '/api/v1/dns/domains', params })
  },

  /**
   * 获取域名详情
   */
  get: (id: number): Promise<Domain> => {
    return request.get({ url: `/api/v1/dns/domains/${id}` })
  },

  /**
   * 创建域名
   */
  create: (data: Partial<Domain>): Promise<Domain> => {
    return request.post({ url: '/api/v1/dns/domains', data })
  },

  /**
   * 更新域名
   */
  update: (id: number, data: Partial<Domain>): Promise<Domain> => {
    return request.put({ url: `/api/v1/dns/domains/${id}`, data })
  },

  /**
   * 删除域名
   */
  delete: (id: number): Promise<void> => {
    return request.delete(`/api/v1/dns/domains/${id}`)
  },

  /**
   * 批量删除域名
   */
  batchDelete: (ids: number[]): Promise<void> => {
    return request.post('/api/v1/dns/domains/batch', {
      action: 'delete',
      ids
    })
  },

  /**
   * 获取域名统计信息
   */
  statistics: (): Promise<{
    total: number
    by_status: Record<string, number>
    expiring: number
  }> => {
    return request.get('/api/v1/dns/stats/domains')
  },

  /**
   * 导入域名
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
    return request.post('/api/v1/dns/import-export/domains/import', data)
  },

  /**
   * 导出域名
   */
  export: (params: {
    format: string
    filters?: Record<string, any>
    fields?: string[]
  }): Promise<Blob> => {
    return request.post('/api/v1/dns/import-export/domains/export', params, {
      responseType: 'blob'
    })
  }
}
