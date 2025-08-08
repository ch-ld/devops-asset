import { request } from '@/utils/http'
import type { DomainGroup, DomainGroupListParams, DomainGroupListResponse } from '@/types/dns'

/**
 * 域名分组管理API
 */
export const domainGroupApi = {
  /**
   * 获取分组列表
   */
  list: (params: DomainGroupListParams): Promise<DomainGroupListResponse> => {
    return request.get({ url: '/dns/domain-groups', params })
  },

  /**
   * 获取分组详情
   */
  get: (id: number): Promise<DomainGroup> => {
    return request.get({ url: `/dns/domain-groups/${id}` })
  },

  /**
   * 创建分组
   */
  create: (data: Partial<DomainGroup>): Promise<DomainGroup> => {
    return request.post({ url: '/dns/domain-groups', data })
  },

  /**
   * 更新分组
   */
  update: (id: number, data: Partial<DomainGroup>): Promise<DomainGroup> => {
    return request.put({ url: `/dns/domain-groups/${id}`, data })
  },

  /**
   * 删除分组
   */
  delete: (id: number): Promise<void> => {
    return request.del({ url: `/dns/domain-groups/${id}` })
  },

  /**
   * 批量删除分组
   */
  batchDelete: (ids: number[]): Promise<void> => {
    return request.post({ url: '/dns/domain-groups/batch-delete', data: { ids } })
  }
}
