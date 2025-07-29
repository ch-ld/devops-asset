/**
 * WebSocket 工具函数
 * 用于动态获取WebSocket连接URL
 */

/**
 * 获取WebSocket基础URL
 * 根据环境变量和当前运行环境动态确定WebSocket连接地址
 */
export function getWebSocketBaseUrl(): string {
  const { VITE_API_URL } = import.meta.env
  
  // 如果是开发环境且有配置API_URL，使用配置的URL
  if (VITE_API_URL && import.meta.env.DEV) {
    // 将HTTP/HTTPS协议转换为WS/WSS协议
    const apiUrl = new URL(VITE_API_URL)
    const protocol = apiUrl.protocol === 'https:' ? 'wss:' : 'ws:'
    return `${protocol}//${apiUrl.host}`
  }
  
  // 生产环境或没有配置API_URL时，使用当前页面的host
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${protocol}//${window.location.host}`
}

/**
 * 构建SSH WebSocket连接URL
 * @param hostId 主机ID
 * @param token 认证token
 * @param connectionParams 连接参数（如IP地址、端口等）
 */
export function buildSSHWebSocketUrl(
  hostId: string | number, 
  token: string,
  connectionParams?: {
    ip?: string
    port?: number
    username?: string
  }
): string {
  const baseUrl = getWebSocketBaseUrl()
  const params = new URLSearchParams({
    host_id: String(hostId),
    token: token
  })
  
  // 添加连接参数
  if (connectionParams?.ip) {
    params.append('ip', connectionParams.ip)
  }
  if (connectionParams?.port) {
    params.append('port', String(connectionParams.port))
  }
  if (connectionParams?.username) {
    params.append('username', connectionParams.username)
  }
  
  return `${baseUrl}/api/v1/cmdb/ws/ssh?${params.toString()}`
}

/**
 * 获取API基础URL（用于HTTP请求）
 */
export function getApiBaseUrl(): string {
  const { VITE_API_URL } = import.meta.env
  
  // 开发环境使用配置的API_URL
  if (VITE_API_URL && import.meta.env.DEV) {
    return VITE_API_URL
  }
  
  // 生产环境使用当前域名
  return `${window.location.protocol}//${window.location.host}`
}
