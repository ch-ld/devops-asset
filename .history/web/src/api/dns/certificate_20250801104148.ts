import { request } from '@/api/client'
import type { Certificate, CertificateCreateRequest, CertificateListParams } from '@/types/dns'

export const certificateApi = {
  // 获取证书列表
  list: (params: CertificateListParams) => {
    return request.get<{
      list: Certificate[]
      total: number
    }>({ url: '/api/v1/dns/certificates', params })
  },

  // 获取证书详情
  get: (id: number) => {
    return request.get<Certificate>({ url: `/api/v1/dns/certificates/${id}` })
  },

  // 申请证书
  issue: (data: CertificateCreateRequest) => {
    return request.post<Certificate>({ url: '/api/v1/dns/certificates', data })
  },

  // 续期证书
  renew: (id: number) => {
    return request.post<Certificate>({ url: `/api/v1/dns/certificates/${id}/renew` })
  },

  // 吊销证书
  revoke: (id: number, reason?: string) => {
    return request.post({ url: `/api/v1/dns/certificates/${id}/revoke`, data: { reason } })
  },

  // 部署证书
  deploy: (id: number, targets: string[]) => {
    return request.post(`/api/v1/dns/certificates/${id}/deploy`, { targets })
  },

  // 下载证书
  download: (id: number, format: 'pem' | 'pfx' | 'jks' = 'pem') => {
    return request.get(`/api/v1/dns/certificates/${id}/download`, {
      params: { format },
      responseType: 'blob'
    })
  },

  // 删除证书
  delete: (id: number) => {
    return request.delete(`/api/v1/dns/certificates/${id}`)
  },

  // 批量删除证书
  batchDelete: (ids: number[]) => {
    return request.delete('/api/v1/dns/certificates/batch', { data: { ids } })
  }
}
