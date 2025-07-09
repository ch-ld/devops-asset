import client from '@/utils/http'
import type { Provider, Host } from '@/types/api/host'
import { useUserStore } from '@/store/modules/user'

const base = '/cmdb'

// ==============================
// Provider API
// ==============================

export function getProviderList() {
  return client.get<Provider[]>(`${base}/providers`)
}

export function getProvider(id: number) {
  return client.get<Provider>(`${base}/providers/${id}`)
}

export function createProvider(data: Partial<Provider>) {
  return client.post(`${base}/providers`, data)
}

export function updateProvider(id: number, data: Partial<Provider>) {
  return client.put(`${base}/providers/${id}`, data)
}

export function deleteProvider(id: number) {
  return client.delete(`${base}/providers/${id}`)
}

export function syncResources(providerId: number) {
  return client.post({
    url: `/cmdb/providers/${providerId}/sync`,
  });
}

// SFTP Operations
export function listFiles(hostId: number, path: string) {
  return client.get<any>({
    url: '/cmdb/sftp/list',
    params: { host_id: hostId, path },
  });
}

export function getSftpDownloadUrl(hostId: number, path:string): string {
  const token = useUserStore().token;
  // This constructs the URL for direct browser download.
  // The backend will handle the file streaming.
  // We need to pass the token for auth, assuming the download endpoint is protected.
  return `/api/v1/cmdb/sftp/download?host_id=${hostId}&path=${encodeURIComponent(path)}&token=${token}`;
}

export function deleteSftpFile(hostId: number, path: string) {
  return client.delete<any>({
    url: '/cmdb/sftp/delete',
    params: { host_id: hostId, path },
  });
}

// ==============================
// Host API
// ==============================

export function getHostList() {
  return client.get<Host[]>(`${base}/hosts`)
}

export function getHost(id: number) {
  return client.get<Host>(`${base}/hosts/${id}`)
}

export function createHost(data: Partial<Host>) {
  return client.post(`${base}/hosts`, data)
}

export function updateHost(id: number, data: Partial<Host>) {
  return client.put(`${base}/hosts/${id}`, data)
}

export function deleteHost(id: number) {
  return client.delete(`${base}/hosts/${id}`)
} 