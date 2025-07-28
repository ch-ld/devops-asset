import api from '@/api/client'
import { ApiResponse } from '@/api/client'

// 主机组接口类型定义
export interface HostGroup {
  id?: number
  name: string
  description?: string
  parent_id?: number | null
  path?: string
  sort?: number
  host_count?: number
  children?: HostGroup[]
  created_at?: string
  updated_at?: string
}

// 创建主机组请求参数
export interface CreateHostGroupRequest {
  name: string
  description?: string
  parent_id?: number | null
  sort?: number
}

// 更新主机组请求参数
export interface UpdateHostGroupRequest {
  name: string
  description?: string
  parent_id?: number | null
  sort?: number
}

// 移动主机组请求参数
export interface MoveHostGroupRequest {
  parent_id?: number | null
}

const base = '/api/v1/cmdb/host_groups'

// ==============================
// 主机组 API
// ==============================

/**
 * 获取主机组树形结构
 */
export const getHostGroupTree = (): Promise<ApiResponse<HostGroup[]>> => {
  return api.get({ url: `${base}/tree` })
}

/**
 * 获取单个主机组信息
 */
export const getHostGroup = (id: number): Promise<ApiResponse<HostGroup>> => {
  return api.get({ url: `${base}/${id}` })
}

/**
 * 创建主机组
 */
export const createHostGroup = (data: CreateHostGroupRequest): Promise<ApiResponse<HostGroup>> => {
  return api.post({ url: base, data })
}

/**
 * 更新主机组
 */
export const updateHostGroup = (id: number, data: UpdateHostGroupRequest): Promise<ApiResponse<HostGroup>> => {
  return api.put({ url: `${base}/${id}`, data })
}

/**
 * 删除主机组
 */
export const deleteHostGroup = (id: number): Promise<ApiResponse<void>> => {
  return api.del({ url: `${base}/${id}` })
}

/**
 * 移动主机组
 */
export const moveHostGroup = (id: number, parentId?: number | null): Promise<ApiResponse<void>> => {
  return api.put({ 
    url: `${base}/${id}/move`,
    params: { parent_id: parentId }
  })
}

/**
 * 获取主机组下的主机列表
 */
export const getGroupHosts = (groupId: number, params?: {
  page?: number
  page_size?: number
  keyword?: string
}): Promise<ApiResponse<any[]>> => {
  return api.get({
    url: `/api/v1/cmdb/groups/${groupId}/hosts`,
    params
  })
}
