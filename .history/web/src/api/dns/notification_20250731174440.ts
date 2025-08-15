import { request } from '@/api/client'
import type { Notification, NotificationCreateRequest, NotificationListParams } from '@/types/dns'

export const notificationApi = {
  // 获取通知列表
  list: (params: NotificationListParams) => {
    return request.get<{
      items: Notification[]
      total: number
    }>('/api/v1/dns/notifications', { params })
  },

  // 获取通知详情
  get: (id: number) => {
    return request.get<Notification>(`/api/v1/dns/notifications/${id}`)
  },

  // 创建通知
  create: (data: NotificationCreateRequest) => {
    return request.post<Notification>('/api/v1/dns/notifications', data)
  },

  // 更新通知
  update: (id: number, data: Partial<NotificationCreateRequest>) => {
    return request.put<Notification>(`/api/v1/dns/notifications/${id}`, data)
  },

  // 删除通知
  delete: (id: number) => {
    return request.delete(`/api/v1/dns/notifications/${id}`)
  },

  // 标记为已读
  markAsRead: (id: number) => {
    return request.post(`/api/v1/dns/notifications/${id}/read`)
  },

  // 批量标记为已读
  batchMarkAsRead: (ids: number[]) => {
    return request.post('/api/v1/dns/notifications/batch-read', { ids })
  },

  // 批量删除通知
  batchDelete: (ids: number[]) => {
    return request.delete('/api/v1/dns/notifications/batch', { data: { ids } })
  },

  // 获取未读通知数量
  getUnreadCount: () => {
    return request.get<{ count: number }>('/api/v1/dns/notifications/unread-count')
  }
}
