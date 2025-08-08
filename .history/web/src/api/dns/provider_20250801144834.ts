import { request } from '@/api/client'
import type { DNSProvider, DNSProviderCreateRequest, DNSProviderListParams } from '@/types/dns'

export const dnsProviderApi = {
  // 获取DNS提供商列表
  list: (params: DNSProviderListParams) => {
    return request.get<{
      items: DNSProvider[]
      total: number
    }>({ url: '/api/v1/dns/providers', params })
  },

  // 获取DNS提供商详情
  get: (id: number) => {
    return request.get<DNSProvider>({ url: `/api/v1/dns/providers/${id}` })
  },

  // 创建DNS提供商
  create: (data: DNSProviderCreateRequest) => {
    return request.post<DNSProvider>({ url: '/api/v1/dns/providers', data })
  },

  // 更新DNS提供商
  update: (id: number, data: Partial<DNSProviderCreateRequest>) => {
    return request.put<DNSProvider>({ url: `/api/v1/dns/providers/${id}`, data })
  },

  // 删除DNS提供商
  delete: (id: number) => {
    return request.del({ url: `/api/v1/dns/providers/${id}` })
  },

  // 测试DNS提供商连接
  test: (id: number) => {
    return request.post<{ success: boolean; message: string }>({ url: `/api/v1/dns/providers/${id}/test` })
  },

  // 设置默认DNS提供商
  setDefault: (id: number) => {
    return request.post({ url: `/api/v1/dns/providers/${id}/default` })
  },

  // 批量删除DNS提供商
  batchDelete: (ids: number[]) => {
    return request.del({ url: '/api/v1/dns/providers/batch', data: { ids } })
  },

  // 同步单个提供商的域名
  syncDomains: (id: number) => {
    return request.post({ url: `/api/v1/dns/providers/${id}/sync-domains` })
  },

  // 同步所有提供商的域名
  syncAllDomains: () => {
    return request.post({ url: '/api/v1/dns/providers/sync-all-domains' })
  }
}
