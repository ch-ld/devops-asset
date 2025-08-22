// 测试数据 - 用于开发和测试证书管理页面
export const mockCertificates = [
  {
    id: 1,
    common_name: 'example.com',
    subject_alt_names: ['www.example.com', 'api.example.com'],
    san_domains: ['www.example.com', 'api.example.com'],
    issuer: "Let's Encrypt Authority X3",
    cert_type: 'letsencrypt',
    status: 'issued',
    valid_from: '2024-01-01T00:00:00Z',
    valid_to: '2024-04-01T00:00:00Z',
    auto_renew: true,
    fingerprint: 'SHA256:1234567890abcdef...'
  },
  {
    id: 2,
    common_name: 'test.com',
    subject_alt_names: ['*.test.com'],
    san_domains: ['*.test.com'],
    issuer: 'ZeroSSL RSA Domain Secure Site CA',
    cert_type: 'zerossl',
    status: 'issued',
    valid_from: '2024-02-01T00:00:00Z',
    valid_to: '2024-05-01T00:00:00Z',
    auto_renew: false,
    fingerprint: 'SHA256:abcdef1234567890...'
  },
  {
    id: 3,
    common_name: 'expired.com',
    subject_alt_names: ['www.expired.com'],
    san_domains: ['www.expired.com'],
    issuer: "Let's Encrypt Authority X3",
    cert_type: 'letsencrypt',
    status: 'expired',
    valid_from: '2023-01-01T00:00:00Z',
    valid_to: '2023-04-01T00:00:00Z',
    auto_renew: true,
    fingerprint: 'SHA256:expired123456...'
  },
  {
    id: 4,
    common_name: 'manual.com',
    subject_alt_names: ['www.manual.com'],
    san_domains: ['www.manual.com'],
    issuer: 'Custom CA',
    cert_type: 'manual',
    status: 'issued',
    valid_from: '2024-01-15T00:00:00Z',
    valid_to: '2025-01-15T00:00:00Z',
    auto_renew: false,
    fingerprint: 'SHA256:manual123456...'
  },
  {
    id: 5,
    common_name: 'pending.com',
    subject_alt_names: ['www.pending.com'],
    san_domains: ['www.pending.com'],
    issuer: '',
    cert_type: 'letsencrypt',
    status: 'pending',
    valid_from: null,
    valid_to: null,
    auto_renew: true,
    fingerprint: ''
  }
]

export const mockStats = {
  total: 5,
  valid: 3,
  expiring: 1,
  expired: 1
}

// 模拟API响应
export const mockApiResponse = {
  code: 200,
  data: {
    list: mockCertificates,
    total: mockCertificates.length,
    page: 1,
    page_size: 20,
    total_page: 1
  },
  message: 'success'
}
