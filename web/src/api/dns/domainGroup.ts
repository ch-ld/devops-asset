import { request } from '@/utils/http'
import type { DomainGroup, DomainGroupListParams, DomainGroupListResponse } from '@/types/dns'

/**
 * 域名分组管理API
 */
export const domainGroupApi = {
  /**
   * 获取域名分组列表
   */
  list: (params: DomainGroupListParams): Promise<DomainGroupListResponse> => {
    return request.get('/dns/groups', { params })
  },

  /**
   * 获取域名分组树形结构
   */
  tree: (): Promise<DomainGroup[]> => {
    return request.get('/dns/groups/tree')
  },

  /**
   * 获取域名分组详情
   */
  get: (id: number): Promise<DomainGroup> => {
    return request.get(`/dns/groups/${id}`)
  },

  /**
   * 创建域名分组
   */
  create: (data: Partial<DomainGroup>): Promise<DomainGroup> => {
    return request.post('/dns/groups', data)
  },

  /**
   * 更新域名分组
   */
  update: (id: number, data: Partial<DomainGroup>): Promise<DomainGroup> => {
    return request.put(`/dns/groups/${id}`, data)
  },

  /**
   * 删除域名分组
   */
  delete: (id: number): Promise<void> => {
    return request.delete(`/dns/groups/${id}`)
  },

  /**
   * 更新域名分组排序
   */
  updateSort: (id: number, sort: number): Promise<void> => {
    return request.put(`/dns/groups/${id}/sort`, { sort })
  }
}
