import { defineStore } from 'pinia'
import { ref } from 'vue'
import type {
  Provider,
  Host,
  HostGroup,
  HostAlert,
  HostMetrics,
  HostMetricsHistory
} from '@/types/api/host'
import * as hostApi from '@/api/system/host'
import { ApiResponse } from '@/api/client'

interface HostListResponse {
  data: Host[]
  total: number
}

export const useHostStore = defineStore('hostStore', () => {
  // State
  const providerList = ref<Provider[]>([])
  const hostList = ref<Host[]>([])
  const hostGroupTree = ref<HostGroup[]>([])
  const hostAlerts = ref<HostAlert[]>([])
  const hostMetrics = ref<HostMetrics | null>(null)
  const hostMetricsHistory = ref<HostMetricsHistory | null>(null)
  const isLoading = ref(false)
  const isSubmitting = ref(false)
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
      hostGroupTree.value = response as HostGroup[]
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
    try {
      const response = (await hostApi.getHostList(params)) as ApiResponse<HostListResponse>
      hostList.value = response.data.data
      pagination.value.total = response.data.total
      if (params?.page) {
        pagination.value.page = params.page
      }
      if (params?.page_size) {
        pagination.value.pageSize = params.page_size
      }
    } catch (error) {
      console.error('Failed to fetch hosts:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const addHost = async (hostData: Partial<Host>) => {
    isSubmitting.value = true
    try {
      const response = await hostApi.createHost(hostData)
      await fetchHosts()
      return response
    } catch (error) {
      console.error('Failed to add host:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  const addManualHost = async (hostData: Partial<Host>) => {
    isSubmitting.value = true
    try {
      const response = await hostApi.createManualHost(hostData)
      await fetchHosts()
      return response
    } catch (error) {
      console.error('Failed to add manual host:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  const updateHost = async (id: number, hostData: Partial<Host>) => {
    isSubmitting.value = true
    try {
      const response = await hostApi.updateHost(id, hostData)
      await fetchHosts()
      return response
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

  const moveHost = async (id: number, groupId?: number) => {
    isSubmitting.value = true
    try {
      await hostApi.moveHost(id, groupId)
      await fetchHosts()
    } catch (error) {
      console.error('Failed to move host:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  // Sync operations
  const syncHosts = async (providerId?: number) => {
    isLoading.value = true
    try {
      await hostApi.syncHosts(providerId)
      await fetchHosts()
    } catch (error) {
      console.error('Failed to sync hosts:', error)
      throw error
    } finally {
      isLoading.value = false
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

  // Batch operations
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

  const batchChangeStatus = async (ids: number[], status: string) => {
    isSubmitting.value = true
    try {
      await hostApi.batchChangeStatus({ ids, status })
      await fetchHosts()
    } catch (error) {
      console.error('Failed to batch change host status:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  const batchImportHosts = async (file: File) => {
    isSubmitting.value = true
    try {
      const response = await hostApi.batchImportHosts(file)
      await fetchHosts()
      return response
    } catch (error) {
      console.error('Failed to batch import hosts:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  // 新增: 批量管理标签
  const batchUpdateTags = async (
    ids: number[],
    tags: string[],
    action: 'add' | 'remove' | 'replace'
  ) => {
    isSubmitting.value = true
    try {
      await hostApi.batchUpdateTags({ ids, tags, action })
      await fetchHosts()
    } catch (error) {
      console.error('Failed to update tags:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  // 新增: 批量设置过期时间
  const batchSetExpired = async (ids: number[], expiredAt: string) => {
    isSubmitting.value = true
    try {
      await hostApi.batchLifecycleHosts({ ids, expired_at: expiredAt })
      await fetchHosts()
    } catch (error) {
      console.error('Failed to set expire date:', error)
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  // 新增: 获取主机监控指标
  const fetchHostMetrics = async (hostId: number) => {
    isLoading.value = true
    try {
      const response = (await hostApi.getHostMetrics(hostId)) as ApiResponse<HostMetrics>
      hostMetrics.value = response.data
      return response.data
    } catch (error) {
      console.error('Failed to fetch host metrics:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 新增: 获取主机历史监控数据
  const fetchHostMetricsHistory = async (
    hostId: number,
    options?: {
      period?: 'last_hour' | 'last_day' | 'last_week' | 'last_month'
      startTime?: string
      endTime?: string
      metricType?: 'cpu' | 'memory' | 'disk' | 'network' | 'all'
    }
  ) => {
    isLoading.value = true
    try {
      const response = (await hostApi.getHostMetricsHistory(
        hostId,
        options
      )) as ApiResponse<HostMetricsHistory>
      hostMetricsHistory.value = response.data
      return response.data
    } catch (error) {
      console.error('Failed to fetch host metrics history:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 新增: 获取主机告警
  const loadHostAlerts = async (days = 7) => {
    isLoading.value = true
    try {
      const response = (await hostApi.getHostAlerts(days)) as ApiResponse<HostAlert[]>
      hostAlerts.value = response.data
      return response.data
    } catch (error) {
      console.error('Failed to load host alerts:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  return {
    // State
    providerList,
    hostList,
    hostGroupTree,
    hostAlerts,
    hostMetrics,
    hostMetricsHistory,
    isLoading,
    isSubmitting,
    pagination,

    // Provider actions
    fetchProviders,
    addProvider,
    updateProvider,
    deleteProvider,
    syncProviderResources,

    // Host Group actions
    fetchHostGroupTree,
    setHostGroupTree,
    addHostGroup,
    updateHostGroup,
    deleteHostGroup,
    moveHostGroup,

    // Host actions
    fetchHosts,
    addHost,
    addManualHost,
    updateHost,
    deleteHost,
    moveHost,

    // Sync actions
    syncHosts,
    syncHostStatus,

    // Batch operations
    batchDeleteHosts,
    batchChangeStatus,
    batchImportHosts,
    batchUpdateTags,
    batchSetExpired,

    // Metrics and alerts
    fetchHostMetrics,
    fetchHostMetricsHistory,
    loadHostAlerts
  }
})
