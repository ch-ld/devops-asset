import { api } from '../client'

interface Certificate {
  id: number
  domain_id: number
  common_name: string
  san_domains: string[]
  issuer: string
  cert_type: string
  status: string
  valid_from: string
  valid_to: string
  auto_renew: boolean
  fingerprint: string
}

interface CreateCertificateRequest {
  domain_id?: number
  cert_type: 'manual' | 'letsencrypt' | 'acme'
  certificate?: string
  private_key?: string
  certificate_chain?: string
  domains?: string
  challenge_type?: string
  acme_server?: string
  auto_renew?: boolean
  email?: string
}

interface UpdateCertificateRequest {
  auto_renew?: boolean
  [key: string]: any
}

interface ListCertificatesParams {
  domain_id?: number
  page?: number
  page_size?: number
  status?: string
}

interface ListCertificatesResponse {
  code: number
  data: {
    items: Certificate[]
    total: number
  }
  message: string
}

interface CheckCertificateResponse {
  code: number
  data: {
    status: string
    https_status: string
    chain_valid: boolean
    expires_at: string
    issuer: string
    issues: string[]
  }
  message: string
}

export const certificateApi = {
  // 获取证书列表
  list: (params?: ListCertificatesParams) =>
    api.get<ListCertificatesResponse>({
      url: '/api/dns/certificates',
      params
    }),

  // 创建证书
  create: (data: CreateCertificateRequest) =>
    api.post<{ code: number; data: Certificate; message: string }>({
      url: '/api/dns/certificates',
      data
    }),

  // 更新证书
  update: (id: number, data: UpdateCertificateRequest) =>
    api.put<{ code: number; data: Certificate; message: string }>({
      url: `/api/dns/certificates/${id}`,
      data
    }),

  // 删除证书
  delete: (id: number) =>
    api.delete<{ code: number; message: string }>({
      url: `/api/dns/certificates/${id}`
    }),

  // 获取证书详情
  get: (id: number) =>
    api.get<{ code: number; data: Certificate; message: string }>({
      url: `/api/dns/certificates/${id}`
    }),

  // 检查证书
  check: (domain: string) =>
    api.get<CheckCertificateResponse>({
      url: `/api/dns/certificates/check`,
      params: { domain }
    }),

  // 续费证书
  renew: (id: number) =>
    api.post<{ code: number; message: string }>({
      url: `/api/dns/certificates/${id}/renew`
    }),

  // 吊销证书
  revoke: (id: number) =>
    api.post<{ code: number; message: string }>({
      url: `/api/dns/certificates/${id}/revoke`
    }),

  // 下载证书
  download: (id: number) =>
    api.get<Blob>({
      url: `/api/dns/certificates/${id}/download`,
      responseType: 'blob'
    }),

  // 导出证书报告
  exportReport: (domainId: number) =>
    api.get<Blob>({
      url: `/api/dns/certificates/export-report`,
      params: { domain_id: domainId },
      responseType: 'blob'
    })
}