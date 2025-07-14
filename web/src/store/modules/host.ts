import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Provider, Host, HostGroup } from '@/types/api/host'
import * as hostApi from '@/api/system/host'

export const useHostStore = defineStore('hostStore', () => {
  // State
  const providerList = ref<Provider[]>([])
  const hostList = ref<Host[]>([])
  const hostGroupTree = ref<HostGroup[]>([])
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
      const { data } = await hostApi.getProviderList()
      providerList.value = data
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
      const { data } = await hostApi.getHostGroupTree()
      hostGroupTree.value = data
    } catch (error) {
      console.error('Failed to fetch host group tree:', error)
      throw error
    } finally {
      isLoading.value = false
    }
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
      const { data } = await hostApi.getHostList(params)
      hostList.value = data.data
      pagination.value.total = data.total
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
      console.error('Failed to batch change status:', error)
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

  return {
    // State
    providerList,
    hostList,
    hostGroupTree,
    isLoading,
    isSubmitting,
    pagination,
    
    // Provider actions
    fetchProviders,
    addProvider,
    updateProvider,
    deleteProvider,
    syncProviderResources,
    
    // Host group actions
    fetchHostGroupTree,
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
    
    // Batch actions
    batchDeleteHosts,
    batchChangeStatus,
    batchImportHosts,
  }
}) 