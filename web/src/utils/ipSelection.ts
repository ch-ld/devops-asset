/**
 * IP选择工具函数
 */

export interface Host {
  id: number
  name: string
  public_ip: string | string[]
  private_ip: string | string[]
}

export interface IPSelectionResult {
  needsSelection: boolean
  publicIPs: string[]
  privateIPs: string[]
  autoSelectedIP?: string
  autoSelectedType?: 'public' | 'private'
}

/**
 * 解析IP地址数组
 */
export function parseIPs(ipData: string | string[]): string[] {
  if (!ipData) return []
  
  try {
    if (Array.isArray(ipData)) {
      return ipData.filter(ip => ip && ip.trim())
    }
    
    if (typeof ipData === 'string') {
      // 尝试解析JSON数组
      const parsed = JSON.parse(ipData)
      if (Array.isArray(parsed)) {
        return parsed.filter(ip => ip && ip.trim())
      }
      // 如果不是JSON数组，当作单个IP处理
      return ipData.trim() ? [ipData.trim()] : []
    }
  } catch (error) {
    // JSON解析失败，当作单个IP处理
    if (typeof ipData === 'string' && ipData.trim()) {
      return [ipData.trim()]
    }
  }
  
  return []
}

/**
 * 分析主机IP配置，决定是否需要用户选择
 */
export function analyzeHostIPs(host: Host): IPSelectionResult {
  const publicIPs = parseIPs(host.public_ip)
  const privateIPs = parseIPs(host.private_ip)
  
  const result: IPSelectionResult = {
    needsSelection: false,
    publicIPs,
    privateIPs
  }
  
  // 如果没有任何IP，返回错误状态
  if (publicIPs.length === 0 && privateIPs.length === 0) {
    return result
  }
  
  // 如果只有一种类型的IP，自动选择
  if (publicIPs.length > 0 && privateIPs.length === 0) {
    // 只有公网IP
    result.autoSelectedIP = publicIPs[0]
    result.autoSelectedType = 'public'
    return result
  }
  
  if (privateIPs.length > 0 && publicIPs.length === 0) {
    // 只有私网IP
    result.autoSelectedIP = privateIPs[0]
    result.autoSelectedType = 'private'
    return result
  }
  
  // 如果同时有公网IP和私网IP，需要用户选择
  if (publicIPs.length > 0 && privateIPs.length > 0) {
    result.needsSelection = true
    return result
  }
  
  return result
}

/**
 * 获取主机的主要显示IP（用于列表显示）
 */
export function getPrimaryDisplayIP(host: Host): string {
  const analysis = analyzeHostIPs(host)
  
  if (analysis.autoSelectedIP) {
    return analysis.autoSelectedIP
  }
  
  // 如果需要选择，优先显示公网IP
  if (analysis.publicIPs.length > 0) {
    return analysis.publicIPs[0]
  }
  
  if (analysis.privateIPs.length > 0) {
    return analysis.privateIPs[0]
  }
  
  return '无IP'
}

/**
 * 格式化IP显示（带类型标识）
 */
export function formatIPWithType(host: Host): string {
  const analysis = analyzeHostIPs(host)
  
  if (analysis.autoSelectedIP) {
    const type = analysis.autoSelectedType === 'public' ? '公网' : '私网'
    return `${analysis.autoSelectedIP} (${type})`
  }
  
  // 如果需要选择，显示可用选项数量
  if (analysis.needsSelection) {
    const publicCount = analysis.publicIPs.length
    const privateCount = analysis.privateIPs.length
    return `${publicCount}个公网IP, ${privateCount}个私网IP`
  }
  
  return '无IP'
}

/**
 * 检查主机是否有可用的IP地址
 */
export function hasAvailableIPs(host: Host): boolean {
  const analysis = analyzeHostIPs(host)
  return analysis.publicIPs.length > 0 || analysis.privateIPs.length > 0
}

/**
 * 获取连接参数
 */
export interface ConnectionParams {
  ip: string
  ipType: 'public' | 'private'
  preferredIPType: string // 传给后端的参数
}

export function getConnectionParams(ipType: 'public' | 'private', ip: string): ConnectionParams {
  return {
    ip,
    ipType,
    preferredIPType: ipType
  }
}

/**
 * 自动选择连接IP（不需要用户交互的情况）
 */
export function autoSelectConnectionIP(host: Host): ConnectionParams | null {
  const analysis = analyzeHostIPs(host)
  
  if (analysis.autoSelectedIP && analysis.autoSelectedType) {
    return getConnectionParams(analysis.autoSelectedType, analysis.autoSelectedIP)
  }
  
  return null
}
