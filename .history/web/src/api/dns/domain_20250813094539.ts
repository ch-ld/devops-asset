import api from '@/api/client'
import type { Domain, DomainListParams, DomainListResponse } from '@/types/dns'

/**
 * 域名管理API
 */
export const domainApi = {
  /**
   * 获取域名列表
   */
  list: (params: DomainListParams): Promise<DomainListResponse> => {
    return api.get({ url: '/api/v1/dns/domains', params })
  },

  /**
   * 获取域名详情
   */
  get: (id: number): Promise<Domain> => {
    return api.get({ url: `/api/v1/dns/domains/${id}` })
  },

  /**
   * 创建域名
   */
  create: (data: Partial<Domain>): Promise<Domain> => {
    return api.post({ url: '/api/v1/dns/domains', data })
  },

  /**
   * 更新域名
   */
  update: (id: number, data: Partial<Domain>): Promise<Domain> => {
    return api.put({ url: `/api/v1/dns/domains/${id}`, data })
  },

  /**
   * 删除域名
   */
  delete: (id: number): Promise<void> => {
    return api.del({ url: `/api/v1/dns/domains/${id}` })
  },

  /**
   * 批量删除域名
   */
  batchDelete: (ids: number[]): Promise<void> => {
    return api.post({ url: '/api/v1/dns/domains/batch', data: {
      action: 'delete',
      ids
    }})
  },

  /**
   * 获取域名统计信息
   */
  statistics: (): Promise<{
    total: number
    by_status: Record<string, number>
    expiring: number
  }> => {
    return api.get({ url: '/api/v1/dns/stats/domains' })
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
  },

  /**
   * 导出单个域名配置
   */
  exportSingle: (id: number): Promise<Blob> => {
    return request.get({ 
      url: `/api/v1/dns/domains/${id}/export`,
      responseType: 'blob'
    })
  },

  /**
   * WHOIS查询
   */
  whois: (id: number): Promise<{
    data: {
      domain: string
      status: string
      registrar: string
      creation_date: string
      updated_date: string
      expiration_date: string
      name_servers: string[]
      domain_status: string
      registrant?: {
        name?: string
        organization?: string
        email?: string
        country?: string
      }
      admin?: {
        name?: string
        email?: string
      }
      tech?: {
        name?: string
        email?: string
      }
      raw_data: string
      query_time: string
    }
  }> => {
    return request.get({ url: `/api/v1/dns/domains/${id}/whois` })
  },

  /**
   * 导出WHOIS信息
   */
  exportWhois: (id: number): Promise<Blob> => {
    return request.get({ 
      url: `/api/v1/dns/domains/${id}/whois/export`,
      responseType: 'blob'
    })
  }
}
