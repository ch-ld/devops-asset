import api from '../client'

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
  domain_id: number
  domains: string[]
  email: string
  provider_id: number
  key_type?: string
  valid_days?: number
  auto_renew?: boolean
  deploy_hosts?: number[]
  remark?: string
}

interface CSRUploadRequest {
  domains: string[]
  email: string
  provider_id: number
  csr_content: string
  auto_renew?: boolean
  deploy_hosts?: number[]
  remark?: string
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
      url: '/api/v1/dns/certificates',
      params
    }),

  // 创建证书
  create: (data: CreateCertificateRequest) =>
    api.post<{ code: number; data: Certificate; message: string }>({
      url: '/api/v1/dns/certificates',
      data
    }),

  // 使用CSR申请证书
  createWithCSR: (data: CSRUploadRequest) =>
    api.post<{ code: number; data: Certificate; message: string }>({
      url: '/api/v1/dns/certificates/csr',
      data
    }),

  // 更新证书
  update: (id: number, data: UpdateCertificateRequest) =>
    api.put<{ code: number; data: Certificate; message: string }>({
      url: `/api/v1/dns/certificates/${id}`,
      data
    }),

  // 删除证书
  delete: (id: number) =>
    api.delete<{ code: number; message: string }>({
      url: `/api/v1/dns/certificates/${id}`
    }),

  // 获取证书详情
  get: (id: number) =>
    api.get<{ code: number; data: Certificate; message: string }>({
      url: `/api/v1/dns/certificates/${id}`
    }),

  // 检查证书
  check: (domain: string) =>
    api.get<CheckCertificateResponse>({
      url: `/api/v1/dns/certificates/check`,
      params: { domain }
    }),

  // 续费证书
  renew: (id: number) =>
    api.post<{ code: number; message: string }>({
      url: `/api/v1/dns/certificates/${id}/renew`
    }),

  // 吊销证书
  revoke: (id: number) =>
    api.post<{ code: number; message: string }>({
      url: `/api/v1/dns/certificates/${id}/revoke`
    }),

  // 下载证书
  download: (id: number) =>
    api.get<Blob>({
      url: `/api/v1/dns/certificates/${id}/download`,
      responseType: 'blob'
    }),

  // 导出证书报告
  exportReport: (domainId: number) =>
    api.get<Blob>({
      url: `/api/v1/dns/certificates/export-report`,
      params: { domain_id: domainId },
      responseType: 'blob'
    }),

  // 获取证书统计信息
  stats: () =>
    api.get<{ code: number; data: any; message: string }>({
      url: '/api/v1/dns/certificates/stats'
    })
}
