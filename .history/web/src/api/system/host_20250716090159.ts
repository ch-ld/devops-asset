import client from '@/utils/http'
import type {
  Provider,
  Host,
  HostGroup,
  BatchOperationResult,
  SSHResult,
  SFTPResult,
  HostAlert
} from '@/types/api/host'
import { useUserStore } from '@/store/modules/user'

const base = '/api/v1/cmdb'

// ==============================
// Provider API
// ==============================

export function getProviderList() {
  return client.get<Provider[]>({ url: `${base}/providers` })
}

export function getProvider(id: number) {
  return client.get<Provider>({ url: `${base}/providers/${id}` })
}

export function createProvider(data: Partial<Provider>) {
  return client.post({ url: `${base}/providers`, data })
}

export function updateProvider(id: number, data: Partial<Provider>) {
  return client.put({ url: `${base}/providers/${id}`, data })
}

export function deleteProvider(id: number) {
  return client.del({ url: `${base}/providers/${id}` })
}

export function syncResources(providerId: number) {
  return client.post({ url: `${base}/providers/${providerId}/sync` })
}

// ==============================
// Host Group API
// ==============================

export function getHostGroupTree() {
  return client.get<HostGroup[]>({ url: `${base}/host_groups/tree` })
}

export function getHostGroup(id: number) {
  return client.get<HostGroup>({ url: `${base}/host_groups/${id}` })
}

export function createHostGroup(data: Partial<HostGroup>) {
  return client.post({ url: `${base}/host_groups`, data })
}

export function updateHostGroup(id: number, data: Partial<HostGroup>) {
  return client.put({ url: `${base}/host_groups/${id}`, data })
}

export function deleteHostGroup(id: number) {
  return client.del({ url: `${base}/host_groups/${id}` })
}

export function moveHostGroup(id: number, parentId?: number) {
  return client.put({
    url: `${base}/host_groups/${id}/move`,
    params: { parent_id: parentId }
  })
}

export function getGroupHosts(groupId: number, page = 1, pageSize = 20, keyword = '') {
  return client.get({
    url: `${base}/groups/${groupId}/hosts`,
    params: { page, page_size: pageSize, keyword }
  })
}

// ==============================
// Host API
// ==============================

export function getHostList(params?: {
  page?: number
  page_size?: number
  keyword?: string
  status?: string
  group_id?: number
  region?: string
}) {
  return client.get<{ total: number; data: Host[] }>({
    url: `${base}/hosts`,
    params
  })
}

export function getHost(id: number) {
  return client.get<Host>({ url: `${base}/hosts/${id}` })
}

export function createHost(data: Partial<Host>) {
  return client.post({ url: `${base}/hosts`, data })
}

export function updateHost(id: number, data: Partial<Host>) {
  return client.put({ url: `${base}/hosts/${id}`, data })
}

export function deleteHost(id: number) {
  return client.del({ url: `${base}/hosts/${id}` })
}

export function createManualHost(data: Partial<Host>) {
  return client.post({ url: `${base}/hosts/manual`, data })
}

export function moveHost(id: number, groupId?: number) {
  return client.post({
    url: `${base}/hosts/${id}/move`,
    params: { group_id: groupId }
  })
}

// ==============================
// Host Sync API
// ==============================

export function syncHosts(providerId?: number) {
  return client.post(`${base}/hosts/sync`, null, {
    params: { provider_id: providerId }
  })
}

export function syncHostStatus(hostId: number) {
  return client.post(`${base}/hosts/${hostId}/sync_status`)
}

// ==============================
// Host Batch Operations API
// ==============================

export function batchImportHosts(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return client.post<BatchOperationResult>(`${base}/hosts/batch_import`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export function batchExportHosts(params?: {
  format?: 'excel' | 'csv'
  name?: string
  group?: string
  tag?: string
}) {
  return client.get(`${base}/hosts/batch_export`, { params, responseType: 'blob' })
}

export function batchDeleteHosts(ids: number[]) {
  return client.delete(`${base}/hosts/batch_delete`, { data: ids })
}

export function batchUpdateHosts(hosts: Host[]) {
  return client.put(`${base}/hosts/batch_update`, hosts)
}

export function batchAssignHosts(data: {
  ids: number[]
  group?: string
  tags?: string[]
  owner?: string
}) {
  return client.put(`${base}/hosts/batch_assign`, data)
}

export function batchLifecycleHosts(data: {
  ids: number[]
  expired_at?: string
  status?: string
  recycle?: boolean
}) {
  return client.put(`${base}/hosts/lifecycle`, data)
}

export function batchSetCustomFields(data: {
  ids: number[]
  extra_fields: { [key: string]: any }
}) {
  return client.put(`${base}/hosts/custom_fields`, data)
}

export function batchChangeStatus(data: { ids: number[]; status: string }) {
  return client.put(`${base}/hosts/batch_status`, data)
}

// 新增批量标签操作API
export function batchUpdateTags(data: {
  ids: number[]
  tags: string[]
  action: 'add' | 'remove' | 'replace'
}) {
  return client.put(`${base}/hosts/batch_tags`, data)
}

export function batchSSH(data: { ids: number[]; cmd: string; timeout?: number }) {
  return client.post<SSHResult[]>(`${base}/hosts/batch_ssh`, data)
}

export function batchSFTP(data: { ids: number[]; remote_path: string; file: File }) {
  const formData = new FormData()
  formData.append('file', data.file)
  data.ids.forEach((id) => formData.append('ids', id.toString()))
  formData.append('remote_path', data.remote_path)
  return client.post<SFTPResult[]>(`${base}/hosts/batch_sftp`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// ==============================
// Host Monitoring API
// ==============================

export function getHostAlerts(days = 7) {
  return client.get<HostAlert[]>(`${base}/hosts/alert`, { params: { days } })
}

export function getHostHistory(id: string) {
  return client.get(`${base}/hosts/history`, { params: { id } })
}

// 新增主机监控API
export function getHostMetrics(hostId: number) {
  return client.get(`${base}/hosts/${hostId}/metrics`)
}

export function getHostMetricsHistory(
  hostId: number,
  options?: {
    period?: 'last_hour' | 'last_day' | 'last_week' | 'last_month'
    startTime?: string
    endTime?: string
    metricType?: 'cpu' | 'memory' | 'disk' | 'network' | 'all'
  }
) {
  return client.get(`${base}/hosts/${hostId}/metrics/history`, { params: options })
}

export function getHostsOverallMetrics() {
  return client.get(`${base}/hosts/metrics/overall`)
}

// ==============================
// SFTP Operations API
// ==============================

export function listFiles(hostId: number, path: string) {
  return client.get(`${base}/sftp/list`, {
    params: { host_id: hostId, path }
  })
}

export function getSftpDownloadUrl(hostId: number, path: string): string {
  const token = useUserStore().accessToken
  return `/api/v1${base}/sftp/download?host_id=${hostId}&path=${encodeURIComponent(path)}&token=${token}`
}

export function deleteSftpFile(hostId: number, path: string) {
  return client.delete(`${base}/sftp/delete`, {
    params: { host_id: hostId, path }
  })
}

export function uploadSftpFile(hostId: number, remotePath: string, file: File) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('path', remotePath)
  return client.post(`${base}/sftp/upload`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    params: { host_id: hostId }
  })
}
