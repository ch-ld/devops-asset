import api from '@/api/client'
import { ApiResponse } from '@/api/client'

// 获取公告列表
export const getAnnouncementList = (params?: {
  title?: string
  type?: number
  status?: number
  page?: number
  page_size?: number
}): Promise<ApiResponse<any[]>> => {
  return api.get({ url: '/api/v1/announcements', params })
}

// 获取公告详情
export const getAnnouncementDetail = (id: number): Promise<ApiResponse<any>> => {
  return api.get({ url: `/api/v1/announcements/${id}` })
}

// 创建公告
export const createAnnouncement = (data: {
  title: string
  content: string
  type: number
  status: number
  is_pinned: number
  priority: number
  publish_time?: string
  expire_time?: string
}): Promise<ApiResponse<any>> => {
  return api.post({ url: '/api/v1/announcements', data })
}

// 更新公告
export const updateAnnouncement = (id: number, data: {
  title: string
  content: string
  type: number
  status: number
  is_pinned: number
  priority: number
  publish_time?: string
  expire_time?: string
}): Promise<ApiResponse<any>> => {
  return api.put({ url: `/api/v1/announcements/${id}`, data })
}

// 删除公告
export const deleteAnnouncement = (id: number): Promise<ApiResponse<any>> => {
  return api.del({ url: `/api/v1/announcements/${id}` })
}

// 获取活跃公告(用于首页显示)
export const getActiveAnnouncements = (limit: number = 5): Promise<ApiResponse<any[]>> => {
  return api.get({ url: '/api/v1/announcements/active', params: { limit } })
} 