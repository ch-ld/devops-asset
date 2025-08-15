import api from '@/api/client'
import type { DNSProvider, DNSProviderCreateRequest, DNSProviderListParams } from '@/types/dns'

export const dnsProviderApi = {
  // 获取DNS提供商列表
  list: (params: DNSProviderListParams) => {
    return api.get<{
      items: DNSProvider[]
      total: number
    }>({ url: '/api/v1/dns/providers', params })
  },

  // 获取DNS提供商详情
  get: (id: number) => {
    return api.get<DNSProvider>({ url: `/api/v1/dns/providers/${id}` })
  },

  // 创建DNS提供商
  create: (data: DNSProviderCreateRequest) => {
    return api.post<DNSProvider>({ url: '/api/v1/dns/providers', data })
  },

  // 更新DNS提供商
  update: (id: number, data: Partial<DNSProviderCreateRequest>) => {
    return api.put<DNSProvider>({ url: `/api/v1/dns/providers/${id}`, data })
  },

  // 删除DNS提供商
  delete: (id: number) => {
    return api.del({ url: `/api/v1/dns/providers/${id}` })
  },

  // 测试DNS提供商连接
  test: (id: number) => {
    return api.post<{ success: boolean; message: string }>({ url: `/api/v1/dns/providers/${id}/test` })
  },

  // 临时测试连接（不保存提供商）
  testConnection: (data: any) => {
    return api.post<{ success: boolean; message: string }>({ url: '/api/v1/dns/providers/test-connection', data })
  },

  // 设置默认DNS提供商
  setDefault: (id: number) => {
    return api.post({ url: `/api/v1/dns/providers/${id}/default` })
  },

  // 批量删除DNS提供商
  batchDelete: (ids: number[]) => {
    return api.del({ url: '/api/v1/dns/providers/batch', data: { ids } })
  },

  // 批量测试DNS提供商连接
  batchTest: (ids: number[]) => {
    return api.post({ url: '/api/v1/dns/providers/batch-test', data: { ids } })
  },

  // 同步单个提供商的域名
  syncDomains: (id: number) => {
    return api.post({ url: `/api/v1/dns/providers/${id}/sync-domains` })
  },

  // 同步所有提供商的域名
  syncAllDomains: () => {
    return api.post({ url: '/api/v1/dns/providers/sync-all-domains' })
  },

  // 获取DNS提供商统计信息
  statistics: () => {
    return api.get<{
      total: number
      by_status: Record<string, number>
      by_type: Record<string, number>
    }>({ url: '/api/v1/dns/stats/providers' })
  }
}

export default dnsProviderApi
