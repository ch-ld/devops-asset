import api from '@/api/client'
import type { HTTPSMonitor, HTTPSMonitorCreateRequest, HTTPSMonitorListParams } from '@/types/dns'

export const httpsMonitorApi = {
  // 获取HTTPS监控列表
  list: (params: HTTPSMonitorListParams) => {
    return api.get<{
      items: HTTPSMonitor[]
      total: number
    }>({ url: '/api/v1/dns/monitors', params })
  },

  // 获取HTTPS监控详情
  get: (id: number) => {
    return api.get<HTTPSMonitor>({ url: `/api/v1/dns/monitors/${id}` })
  },

  // 创建HTTPS监控
  create: (data: HTTPSMonitorCreateRequest) => {
    return api.post<HTTPSMonitor>({ url: '/api/v1/dns/monitors', data })
  },

  // 更新HTTPS监控
  update: (id: number, data: Partial<HTTPSMonitorCreateRequest>) => {
    return api.put<HTTPSMonitor>({ url: `/api/v1/dns/monitors/${id}`, data })
  },

  // 删除HTTPS监控
  delete: (id: number) => {
    return api.del({ url: `/api/v1/dns/monitors/${id}` })
  },

  // 手动检查监控
  check: (id: number) => {
    return api.post<{
      status: string
      response_time: number
      ssl_info: any
      error?: string
    }>({ url: `/api/v1/dns/monitors/${id}/check` })
  },

  // 获取监控统计
  getStatistics: () => {
    return api.get<{
      total: number
      online: number
      offline: number
      warning: number
    }>({ url: '/api/v1/dns/monitors/statistics' })
  },

  // 获取即将过期的证书
  getExpiringCertificates: (days: number = 30) => {
    return api.get<HTTPSMonitor[]>({ url: '/api/v1/dns/monitors/expiring', params: { days } })
  },

  // 批量删除监控
  batchDelete: (ids: number[]) => {
    return api.del({ url: '/api/v1/dns/monitors/batch', data: { ids } })
  }
}
