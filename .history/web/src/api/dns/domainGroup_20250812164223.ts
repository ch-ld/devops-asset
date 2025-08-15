import { api } from '../client'
import type { DomainGroup } from '@/types/dns'

interface CreateDomainGroupRequest {
  name: string
  description?: string
  parent_id?: number | null
  sort?: number
  color?: string
  status?: 'active' | 'inactive'
}

interface UpdateDomainGroupRequest extends Partial<CreateDomainGroupRequest> {}

interface ListDomainGroupsParams {
  page?: number
  page_size?: number
  keyword?: string
  parent_id?: number | null
  status?: string
}

interface ListDomainGroupsResponse {
  code: number
  data: {
    items: DomainGroup[]
    total: number
  }
  message: string
}

export const domainGroupApi = {
  // 获取分组列表
  list: (params?: ListDomainGroupsParams) => 
    api.get<ListDomainGroupsResponse>({
      url: '/api/dns/domain-groups',
      params
    }),

  // 创建分组
  create: (data: CreateDomainGroupRequest) =>
    api.post<{ code: number; data: DomainGroup; message: string }>({
      url: '/api/dns/domain-groups',
      data
    }),

  // 更新分组
  update: (id: number, data: UpdateDomainGroupRequest) =>
    api.put<{ code: number; data: DomainGroup; message: string }>({
      url: `/api/dns/domain-groups/${id}`,
      data
    }),

  // 删除分组
  delete: (id: number) =>
    api.delete<{ code: number; message: string }>({
      url: `/api/dns/domain-groups/${id}`
    }),

  // 获取分组详情
  get: (id: number) =>
    api.get<{ code: number; data: DomainGroup; message: string }>({
      url: `/api/dns/domain-groups/${id}`
    }),

  // 获取分组树
  tree: () =>
    api.get<{ code: number; data: DomainGroup[]; message: string }>({
      url: '/api/dns/domain-groups/tree'
    }),

  // 移动分组
  move: (id: number, data: { parent_id?: number | null; sort?: number }) =>
    api.put<{ code: number; message: string }>({
      url: `/api/dns/domain-groups/${id}/move`,
      data
    })
}