import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Provider, Host } from '@/types/api/host'
import * as hostApi from '@/api/system/host'

export const useHostStore = defineStore('hostStore', () => {
  // State
  const providerList = ref<Provider[]>([])
  const hostList = ref<Host[]>([])
  const isLoading = ref(false)
  const isSubmitting = ref(false) // Added for new actions

  // Actions
  const fetchProviders = async () => {
    isLoading.value = true
    try {
      const { data } = await hostApi.getProviderList()
      providerList.value = data
    } catch (error) {
      console.error('Failed to fetch providers:', error)
      // TODO: Add user-friendly error notification
    } finally {
      isLoading.value = false
    }
  }

  const fetchHosts = async () => {
    isLoading.value = true
    try {
      const { data } = await hostApi.getHostList()
      hostList.value = data
    } catch (error) {
      console.error('Failed to fetch hosts:', error)
      // TODO: Add user-friendly error notification
    } finally {
      isLoading.value = false
    }
  }

  const addProvider = async (provider: Partial<Provider>) => {
    try {
      await hostApi.createProvider(provider)
      await fetchProviders() // Refresh list after adding
    } catch (error) {
      console.error('Failed to add provider:', error)
      // TODO: Add user-friendly error notification
    }
  }

  const updateProvider = async (id: number, provider: Partial<Provider>) => {
    try {
      await hostApi.updateProvider(id, provider)
      await fetchProviders() // Refresh list after updating
    } catch (error) {
      console.error('Failed to update provider:', error)
      // TODO: Add user-friendly error notification
    }
  }

  const deleteProvider = async (id: number) => {
    try {
      await hostApi.deleteProvider(id)
      await fetchProviders() // Refresh list after deleting
    } catch (error) {
      console.error('Failed to delete provider:', error)
      // TODO: Add user-friendly error notification
    }
  }

  const syncProviderResources = async (id: number) => {
    isLoading.value = true; // Use main loading state for simplicity
    try {
      await hostApi.syncResources(id);
      await fetchHosts(); // Refresh host list after sync
      // Optionally show a success message
    } catch (error) {
      console.error('Failed to sync resources:', error);
      // TODO: Add user-friendly error notification
    } finally {
      isLoading.value = false;
    }
  }

  const addHost = async (hostData: Partial<Host>) => {
    try {
      isSubmitting.value = true;
      const response = await hostApi.createHost(hostData);
      hostList.value.push(response.data);
    } catch (error) {
      console.error("Failed to add host:", error);
      throw error;
    } finally {
      isSubmitting.value = false;
    }
  }

  const updateHost = async (id: number, hostData: Partial<Host>) => {
    try {
      isSubmitting.value = true;
      const response = await hostApi.updateHost(id, hostData);
      const index = hostList.value.findIndex(h => h.id === id);
      if (index !== -1) {
        hostList.value[index] = response.data;
      }
    } catch (error) {
      console.error("Failed to update host:", error);
      throw error;
    } finally {
      isSubmitting.value = false;
    }
  }

  const deleteHost = async (id: number) => {
    try {
      isLoading.value = true;
      await hostApi.deleteHost(id);
      hostList.value = hostList.value.filter(h => h.id !== id);
    } catch (error) {
      console.error("Failed to delete host:", error);
      throw error;
    } finally {
      isLoading.value = false;
    }
  }

  // ... (actions for hosts can be added similarly)

  return {
    // State
    providerList,
    hostList,
    isLoading,
    isSubmitting, // Added to state
    // Actions
    fetchProviders,
    fetchHosts,
    addProvider,
    updateProvider,
    deleteProvider,
    syncProviderResources,
    addHost, // Added to return
    updateHost, // Added to return
    deleteHost, // Added to return
  }
}) 