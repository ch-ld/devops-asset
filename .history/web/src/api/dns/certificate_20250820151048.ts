import api from '../client'

interface Certificate {
  id: number
  domain_id: number
  common_name: string
  domain_name?: string // 兼容旧字段
  subject_alt_names?: string[]
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
  domains: string[]
  email: string
  provider_id: number
  key_type?: string
  valid_days?: number
  auto_renew?: boolean
  renew_days?: number
  challenge_type?: string
  schedule_time?: string
  notification_type?: string
  notification_email?: string
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

interface CertificateUploadRequest {
  cert_content: string
  key_content: string
  chain_content?: string
  auto_deploy?: boolean
  deploy_hosts?: number[]
  deploy_path?: string
  restart_command?: string
  email_notification?: boolean
  notification_email?: string
  remark?: string
}

interface ValidateCSRRequest {
  csrContent: string
}

interface ValidateCertificateRequest {
  certContent: string
  keyContent: string
  chainContent?: string
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
      data,
      timeout: 120000 // 2分钟超时，因为证书申请可能需要较长时间
    }),

  // 使用CSR申请证书
  createWithCSR: (data: CSRUploadRequest) =>
    api.post<{ code: number; data: Certificate; message: string }>({
      url: '/api/v1/dns/certificates/csr',
      data
    }),

  // 上传证书
  upload: (data: CertificateUploadRequest) =>
    api.post<{ code: number; data: Certificate; message: string }>({
      url: '/api/v1/dns/certificates/upload',
      data
    }),

  // 获取证书详情
  get: (id: number) =>
    api.get<{ code: number; data: Certificate; message: string }>({
      url: `/api/v1/dns/certificates/${id}`
    }),

  // 验证CSR
  validateCSR: (data: ValidateCSRRequest) =>
    api.post<{ code: number; data: any; message: string }>({
      url: '/api/v1/dns/certificates/validate-csr',
      data
    }),

  // 验证证书
  validateCertificate: (data: ValidateCertificateRequest) =>
    api.post<{ code: number; data: any; message: string }>({
      url: '/api/v1/dns/certificates/validate-certificate',
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
  download: (id: number, format?: string) =>
    api.get<Blob>({
      url: `/api/v1/dns/certificates/${id}/download`,
      params: { format: format || 'pem' },
      responseType: 'blob'
    }),

  // 部署证书
  deploy: (id: number, hostIds: number[]) =>
    api.post<{ code: number; message: string }>({
      url: `/api/v1/dns/certificates/${id}/deploy`,
      data: { host_ids: hostIds }
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
    }),

  // 批量续期
  batchRenew: (ids: number[]) =>
    api.post<{ code: number; message: string }>({
      url: '/api/v1/dns/certificates/batch-renew',
      data: { ids }
    }),

  // 批量删除
  batchDelete: (ids: number[]) =>
    api.post<{ code: number; message: string }>({
      url: '/api/v1/dns/certificates/batch-delete',
      data: { ids }
    }),

  // 批量导出
  batchExport: (ids: number[]) =>
    api.get<Blob>({
      url: '/api/v1/dns/certificates/batch-export',
      params: { ids: ids.join(',') },
      responseType: 'blob'
    }),

  // 批量部署
  batchDeploy: (ids: number[], hostIds: number[]) =>
    api.post<{ code: number; message: string }>({
      url: '/api/v1/dns/certificates/batch-deploy',
      data: { cert_ids: ids, host_ids: hostIds }
    })
}
