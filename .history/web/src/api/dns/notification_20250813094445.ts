import api from '@/api/client'
import type { Notification, NotificationCreateRequest, NotificationListParams } from '@/types/dns'

export const notificationApi = {
  // 获取通知列表
  list: (params: NotificationListParams) => {
    return api.get<{
      items: Notification[]
      total: number
    }>({ url: '/api/v1/dns/notifications', params })
  },

  // 获取通知详情
  get: (id: number) => {
    return api.get<Notification>({ url: `/api/v1/dns/notifications/${id}` })
  },

  // 创建通知
  create: (data: NotificationCreateRequest) => {
    return api.post<Notification>({ url: '/api/v1/dns/notifications', data })
  },

  // 更新通知
  update: (id: number, data: Partial<NotificationCreateRequest>) => {
    return api.put<Notification>({ url: `/api/v1/dns/notifications/${id}`, data })
  },

  // 删除通知
  delete: (id: number) => {
    return api.del({ url: `/api/v1/dns/notifications/${id}` })
  },

  // 标记为已读
  markAsRead: (id: number) => {
    return api.post({ url: `/api/v1/dns/notifications/${id}/read` })
  },

  // 批量标记为已读
  batchMarkAsRead: (ids: number[]) => {
    return api.post({ url: '/api/v1/dns/notifications/batch-read', data: { ids } })
  },

  // 批量删除通知
  batchDelete: (ids: number[]) => {
    return api.del({ url: '/api/v1/dns/notifications/batch', data: { ids } })
  },

  // 获取未读通知数量
  getUnreadCount: () => {
    return api.get<{ count: number }>({ url: '/api/v1/dns/notifications/unread-count' })
  }
}
