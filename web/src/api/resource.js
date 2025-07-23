import request from '@/utils/request'

// 资源管理API
export const resourceApi = {
  // 获取整体统计信息
  getOverallStatistics() {
    return request({
      url: '/api/v1/cmdb/resources/statistics',
      method: 'get'
    })
  },

  // 获取资源摘要信息
  getResourceSummary() {
    return request({
      url: '/api/v1/cmdb/resources/summary',
      method: 'get'
    })
  },

  // 获取云厂商统计信息
  getProviderSummary() {
    return request({
      url: '/api/v1/cmdb/resources/providers',
      method: 'get'
    })
  },

  // 获取地区统计信息
  getRegionSummary() {
    return request({
      url: '/api/v1/cmdb/resources/regions',
      method: 'get'
    })
  },

  // 获取状态统计信息
  getStatusSummary() {
    return request({
      url: '/api/v1/cmdb/resources/status',
      method: 'get'
    })
  },

  // 获取即将过期的资源
  getExpiringResources(days = 30) {
    return request({
      url: '/api/v1/cmdb/resources/expiring',
      method: 'get',
      params: { days }
    })
  },

  // 获取仪表盘数据
  getDashboardData() {
    return request({
      url: '/api/v1/cmdb/resources/dashboard',
      method: 'get'
    })
  },

  // 获取资源健康状况
  getResourceHealth() {
    return request({
      url: '/api/v1/cmdb/resources/health',
      method: 'get'
    })
  }
}

// 告警管理API
export const alertApi = {
  // 获取告警列表
  getAlerts(params = {}) {
    return request({
      url: '/api/v1/cmdb/alerts',
      method: 'get',
      params
    })
  },

  // 获取告警统计信息
  getAlertStatistics() {
    return request({
      url: '/api/v1/cmdb/alerts/statistics',
      method: 'get'
    })
  },

  // 获取告警摘要信息
  getAlertSummary() {
    return request({
      url: '/api/v1/cmdb/alerts/summary',
      method: 'get'
    })
  },

  // 检查即将过期的主机
  checkExpiringHosts(days = 30) {
    return request({
      url: '/api/v1/cmdb/alerts/expiring',
      method: 'get',
      params: { days }
    })
  },

  // 检查主机状态异常
  checkHostStatus() {
    return request({
      url: '/api/v1/cmdb/alerts/status',
      method: 'get'
    })
  },

  // 手动触发告警检查
  triggerAlertCheck() {
    return request({
      url: '/api/v1/cmdb/alerts/check',
      method: 'post'
    })
  },

  // 获取指定主机的告警信息
  getHostAlerts(hostId) {
    return request({
      url: `/api/v1/cmdb/hosts/${hostId}/alerts`,
      method: 'get'
    })
  }
}

export default {
  resourceApi,
  alertApi
}
