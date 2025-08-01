import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'
import type {
  Provider,
  Host,
  HostGroup,
  HostAlert,
  HostMetrics,
  HostMetricsHistory
} from '@/types/api/host'
import * as hostApi from '@/api/system/host'

export const useHostStore = defineStore('hostStore', () => {
  // State
  const providerList = ref<Provider[]>([])
  const hostList = ref<Host[]>([])
  const hostGroupTree = ref<HostGroup[]>([])
  // 标记是否已加载过主机组数据，避免重复请求
  const hasLoadedGroups = ref(false)
  const hostAlerts = ref<HostAlert[]>([])
  const hostMetrics = ref<HostMetrics | null>(null)
  const hostMetricsHistory = ref<HostMetricsHistory | null>(null)
  const isLoading = ref(false)
  const isSubmitting = ref(false)
  const error = ref<Error | null>(null)
  const pagination = ref({
    page: 1,
    pageSize: 20,
    total: 0
  })

  // Actions for Providers
  const fetchProviders = async () => {
    isLoading.value = true
    try {
      const response = await hostApi.getProviderList()
      providerList.value = response as Provider[]
    } catch (error) {
      console.error('Failed to fetch providers:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const addProvider = async (provider: Partial<Provider>) => {
    isSubmitting.value = true
    try {
      await hostApi.createProvider(provider)
      await fetchProviders()
    } catch (error) {
      console.error('Failed to add provider:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  const updateProvider = async (id: number, provider: Partial<Provider>) => {
    isSubmitting.value = true
    try {
      await hostApi.updateProvider(id, provider)
      await fetchProviders()
    } catch (error) {
      console.error('Failed to update provider:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  const deleteProvider = async (id: number) => {
    isSubmitting.value = true
    try {
      await hostApi.deleteProvider(id)
      await fetchProviders()
    } catch (error) {
      console.error('Failed to delete provider:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  const syncProviderResources = async (id: number) => {
    isLoading.value = true
    try {
      await hostApi.syncResources(id)
      await fetchHosts()
    } catch (error) {
      console.error('Failed to sync resources:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // Actions for Host Groups
  const fetchHostGroupTree = async () => {
    isLoading.value = true
    try {
      const response = await hostApi.getHostGroupTree()
      // 处理后端返回的数据结构 {code, status, message, data, timestamp}
      hostGroupTree.value = (response as any)?.data || response as HostGroup[] || []
      hasLoadedGroups.value = true
    } catch (error) {
      console.error('Failed to fetch host group tree:', error)
      // 确保即使API调用失败，hostGroupTree也是一个空数组而不是undefined
      hostGroupTree.value = []
      // 不再向上抛出错误，避免阻断UI渲染
      // throw error
    } finally {
      isLoading.value = false
    }
  }

  // Set host group tree directly
  const setHostGroupTree = (groups: HostGroup[]) => {
    hostGroupTree.value = groups
  }

  const addHostGroup = async (group: Partial<HostGroup>) => {
    isSubmitting.value = true
    try {
      await hostApi.createHostGroup(group)
      await fetchHostGroupTree()
    } catch (error) {
      console.error('Failed to add host group:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  const updateHostGroup = async (id: number, group: Partial<HostGroup>) => {
    isSubmitting.value = true
    try {
      await hostApi.updateHostGroup(id, group)
      await fetchHostGroupTree()
    } catch (error) {
      console.error('Failed to update host group:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  const deleteHostGroup = async (id: number) => {
    isSubmitting.value = true
    try {
      await hostApi.deleteHostGroup(id)
      await fetchHostGroupTree()
    } catch (error) {
      console.error('Failed to delete host group:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  const moveHostGroup = async (id: number, parentId?: number) => {
    isSubmitting.value = true
    try {
      await hostApi.moveHostGroup(id, parentId)
      await fetchHostGroupTree()
    } catch (error) {
      console.error('Failed to move host group:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  /**
   * 供外部调用：确保 hostGroupTree 已加载完毕
   * 若已加载过则立即 resolve；否则等待 fetch 完成
   */
  const ensureHostGroupsReady = async () => {
    if (hasLoadedGroups.value && hostGroupTree.value.length) return
    await fetchHostGroupTree()
  }

  // Actions for Hosts
  const fetchHosts = async (params?: {
    page?: number
    page_size?: number
    keyword?: string
    status?: string
    group_id?: number
    region?: string
  }) => {
    isLoading.value = true
    error.value = null

    try {
      // 构建API参数，只传递有值的参数
      const apiParams: any = {
        page: 1,
        page_size: 100  // 后端限制最大值为100
      }

      // 只有当参数有值时才添加到请求中
      if (params?.keyword) apiParams.keyword = params.keyword
      if (params?.status) apiParams.status = params.status
      if (params?.region) apiParams.region = params.region
      if (params?.group_id) apiParams.group_id = params.group_id

      console.log('发送API请求参数:', apiParams)
      // 使用修复后的API函数
      const response = await hostApi.getHostList(apiParams)

      // 处理后端返回的数据结构
      console.log('API响应:', response)
      console.log('response.data:', response.data)
      console.log('response.data.data:', response.data?.data)

      if (response && response.data && Array.isArray(response.data)) {
        // API返回的数据结构：response.data 直接是主机数组
        console.log('解析的主机数据:', response.data)
        console.log('主机数据长度:', response.data.length)

        hostList.value = response.data
        pagination.value = {
          page: 1,
          pageSize: response.data.length,
          total: response.data.length
        }
        console.log('主机数据设置成功，数量:', response.data.length)
      } else {
        console.log('API响应格式错误或数据为空')
        console.log('response.data类型:', typeof response.data)
        console.log('response.data是否为数组:', Array.isArray(response.data))
        hostList.value = []
        pagination.value = {
          page: 1,
          pageSize: 0,
          total: 0
        }
      }
    } catch (err) {
      console.error('Failed to fetch hosts:', err)
      error.value = err as Error
      hostList.value = []
    } finally {
      isLoading.value = false
    }
  }

  const getHost = async (id: number) => {
    isLoading.value = true
    try {
      const response = await hostApi.getHost(id)
      return response
    } catch (error) {
      console.error('Failed to get host:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const addHost = async (host: Partial<Host>) => {
    isSubmitting.value = true
    try {
      await hostApi.createHost(host)
      await fetchHosts()
    } catch (error) {
      console.error('Failed to add host:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  const updateHost = async (id: number, host: Partial<Host>) => {
    isSubmitting.value = true
    try {
      await hostApi.updateHost(id, host)
      await fetchHosts()
    } catch (error) {
      console.error('Failed to update host:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  const deleteHost = async (id: number) => {
    isSubmitting.value = true
    try {
      await hostApi.deleteHost(id)
      await fetchHosts()
    } catch (error) {
      console.error('Failed to delete host:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  // 添加一些批量操作的方法
  const batchDeleteHosts = async (ids: number[]) => {
    isSubmitting.value = true
    try {
      await hostApi.batchDeleteHosts(ids)
      await fetchHosts()
    } catch (error) {
      console.error('Failed to batch delete hosts:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  const batchMoveHosts = async (data: { ids: number[]; group_id?: number }) => {
    isSubmitting.value = true
    try {
      await hostApi.batchMoveHosts(data)
      await fetchHosts()
    } catch (error) {
      console.error('Failed to batch move hosts:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  const syncHostStatus = async (hostId: number) => {
    isLoading.value = true
    try {
      await hostApi.syncHostStatus(hostId)
      await fetchHosts()
    } catch (error) {
      console.error('Failed to sync host status:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 同步全部主机（云资源）
  const syncHosts = async (providerId?: number) => {
    await hostApi.syncHosts(providerId)
    await fetchHosts()
  }

  return {
    // State
    providerList,
    hostList,
    hostGroupTree,
    hasLoadedGroups,
    hostAlerts,
    hostMetrics,
    hostMetricsHistory,
    isLoading,
    isSubmitting,
    error,
    pagination,

    // Actions for Providers
    fetchProviders,
    addProvider,
    updateProvider,
    deleteProvider,
    syncProviderResources,

    // Actions for Host Groups
    fetchHostGroupTree,
    setHostGroupTree,
    addHostGroup,
    updateHostGroup,
    deleteHostGroup,
    moveHostGroup,
    ensureHostGroupsReady,

    // Actions for Hosts
    fetchHosts,
    getHost,
    addHost,
    createHost: addHost, // 别名，兼容主机管理页面的调用
    updateHost,
    deleteHost,

    // Extra helpers for HostDetail compatibility
    // 兼容旧组件调用
    fetchHost: getHost,
    fetchHostMetrics: async (hostId: number) => {
      return await hostApi.getHostMetrics(hostId)
    },
    fetchHostMetricsHistory: async (
      hostId: number,
      options?: {
        period?: 'last_hour' | 'last_day' | 'last_week' | 'last_month'
        startTime?: string
        endTime?: string
        metricType?: 'cpu' | 'memory' | 'disk' | 'network' | 'all'
      }
    ) => {
      return await hostApi.getHostMetricsHistory(hostId, options)
    },
    batchUpdateTags: async (
      ids: number[],
      tags: string[],
      action: 'add' | 'remove' | 'replace' = 'replace'
    ) => {
      await hostApi.batchUpdateTags({ ids, tags, action })
      await fetchHosts()
    },
    moveHost: async (hostId: number, groupId?: number) => {
      await hostApi.moveHost(hostId, groupId)
      await fetchHosts()
    },

    // Batch operations
    batchDeleteHosts,
    batchMoveHosts,
    syncHostStatus,

    // 同步全部主机（云资源）
    syncHosts: async (providerId?: number) => {
      await hostApi.syncHosts(providerId)
      await fetchHosts()
    },

    // 兼容 simple-index.vue 旧调用
    getHostList: async (params?: {
      page?: number
      page_size?: number
      keyword?: string
      status?: string
      group_id?: number
      region?: string
    }) => {
      return await hostApi.getHostList(params)
    }
  }
})
