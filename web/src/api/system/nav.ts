import api from '@/api/client'
import { ApiResponse } from '@/api/client'

// 获取导航列表
export const getNavList = (params?: any): Promise<ApiResponse<ApiResponse>> => {
  return api.get({ url: '/api/v1/navs', params })
}

// 新增导航
export const addNav = (data: any): Promise<ApiResponse<ApiResponse>> => {
  console.log('API调用 addNav, 发送数据:', data)
  return api.post({ url: '/api/v1/navs', data })
}

// 更新导航
export const updateNav = (data: any): Promise<ApiResponse<ApiResponse>> => {
  console.log('API调用 updateNav, 发送数据:', data)
  const { id, ...requestData } = data  // 提取ID，其余数据作为请求体
  console.log('更新导航 ID:', id, '请求体数据:', requestData)
  return api.put({ url: `/api/v1/navs/${id}`, data: requestData })
}

// 删除导航
export const deleteNav = (id: number): Promise<ApiResponse<ApiResponse>> => {
  return api.del({ url: `/api/v1/navs/${id}` })
} 