import client from '@/api/client'
import axios from 'axios'
import type {
  Provider,
  Host,
  HostGroup,
  BatchOperationResult,
  SSHResult,
  SFTPResult,
  HostAlert
} from '@/types/api/host'

const base = '/api/v1/cmdb'

// ==============================
// Provider API
// ==============================

export function getProviderList(params?: any) {
  return client.get<any>({ url: `${base}/providers`, params })
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

export function validateProviderCredentials(data: {
  type: string
  access_key: string
  secret_key: string
  region?: string
}) {
  return client.post({ url: `${base}/providers/validate`, data })
}

export function getProviderTypes() {
  return client.get({ url: `${base}/providers/types` })
}

export function getProviderRegions(type: string) {
  return client.get({ url: `${base}/providers/types/${type}/regions` })
}

export function getProviderRegionsWithCredentials(data: {
  type: string
  access_key: string
  secret_key: string
}) {
  return client.post({ url: `${base}/providers/regions-with-credentials`, data })
}

// 获取指定云账号的可用区域
export function getProviderRegionsById(providerId: number) {
  return client.get({ url: `${base}/providers/${providerId}/regions` })
}

export function syncResources(providerId: number, groupId?: number) {
  const params = groupId ? { group_id: groupId } : {}
  return client.post({
    url: `${base}/providers/${providerId}/sync`,
    params
  })
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

export async function getHostList(params?: {
  page?: number
  page_size?: number
  keyword?: string
  status?: string
  group_id?: number
  region?: string
}) {
  // 直接使用axios调用，避免client的问题
  const response = await axios.get(`${base}/hosts`, { params })
  return response.data
}

// 获取主机列表（包含完整响应信息）
export async function getHostListWithCount(params?: {
  page?: number
  page_size?: number
  keyword?: string
  status?: string
  group_id?: number
  region?: string
}) {
  const response = await client.get<any>({
    url: `${base}/hosts`,
    params
  })
  return response  // 返回完整的API响应 {code, data, count}
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

export function getHostFilterOptions() {
  return client.get<{
    statuses: string[]
    regions: string[]
    providers: string[]
  }>({ url: `${base}/hosts/filter_options` })
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
  return client.post({
    url: `${base}/hosts/sync`,
    params: { provider_id: providerId }
  })
}

export function syncHostStatus(hostId: number) {
  return client.post({ url: `${base}/hosts/${hostId}/sync_status` })
}

// ==============================
// Host Batch Operations API
// ==============================

export function batchImportHosts(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return client.post<BatchOperationResult>({
    url: `${base}/hosts/batch_import`,
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' },
    transformRequest: [(data) => data] // 不转换FormData
  })
}

export function batchExportHosts(params?: {
  format?: 'excel' | 'csv'
  scope?: 'all' | 'group' | 'current'
  group_id?: number
  fields?: string[]
  name?: string
  status?: string
  region?: string
  provider?: string
}) {
  // 将fields数组转换为逗号分隔的字符串
  const queryParams = { ...params }
  if (queryParams.fields && Array.isArray(queryParams.fields)) {
    queryParams.fields = queryParams.fields.join(',') as any
  }

  return client.get({
    url: `${base}/hosts/batch_export`,
    params: queryParams,
    responseType: 'blob'
  })
}

export function batchDeleteHosts(ids: number[]) {
  return client.del({ url: `${base}/hosts/batch_delete`, data: { ids } })
}

export function batchUpdateHosts(hosts: Host[]) {
  return client.put({ url: `${base}/hosts/batch_update`, data: hosts })
}

export function batchAssignHosts(data: {
  ids: number[]
  group?: string
  tags?: string[]
  owner?: string
}) {
  return client.put({ url: `${base}/hosts/batch_assign`, data })
}

export function batchMoveHosts(data: { ids: number[]; group_id?: number }) {
  return client.put({ url: `${base}/hosts/batch_move`, data })
}

export function batchLifecycleHosts(data: {
  ids: number[]
  expired_at?: string
  status?: string
  recycle?: boolean
}) {
  return client.put({ url: `${base}/hosts/lifecycle`, data })
}

export function batchSetCustomFields(data: {
  ids: number[]
  extra_fields: { [key: string]: any }
}) {
  return client.put({ url: `${base}/hosts/custom_fields`, data })
}

export function batchChangeStatus(data: { ids: number[]; status: string }) {
  return client.put({ url: `${base}/hosts/batch_status`, data })
}

// 新增批量标签操作API
export function batchUpdateTags(data: {
  ids: number[]
  tags: string[]
  action: 'add' | 'remove' | 'replace'
}) {
  return client.put({ url: `${base}/hosts/batch_tags`, data })
}

export function batchSSH(data: { ids: number[]; cmd: string; timeout?: number }) {
  return client.post<SSHResult[]>({ url: `${base}/hosts/batch_ssh`, data })
}

export function batchStart(data: { ids: number[] }) {
  return client.post({ url: `${base}/hosts/batch_start`, data })
}

export function batchStop(data: { ids: number[] }) {
  return client.post({ url: `${base}/hosts/batch_stop`, data })
}

export function batchReboot(data: { ids: number[] }) {
  return client.post({ url: `${base}/hosts/batch_reboot`, data })
}

export function batchSFTP(data: { ids: number[]; remote_path: string; file: File }) {
  const formData = new FormData()
  formData.append('file', data.file)
  data.ids.forEach((id) => formData.append('ids', id.toString()))
  formData.append('remote_path', data.remote_path)
  return client.post<SFTPResult[]>({
    url: `${base}/hosts/batch_sftp`,
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// ==============================
// Host Monitoring API
// ==============================

export function getHostAlerts(days = 7) {
  return client.get<HostAlert[]>({ url: `${base}/hosts/alert`, params: { days } })
}

export function getHostHistory(id: string) {
  return client.get({ url: `${base}/hosts/history`, params: { id } })
}

// 新增主机监控API
export function getHostMetrics(hostId: number) {
  return client.get({ url: `${base}/hosts/${hostId}/metrics` })
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
  return client.get({ url: `${base}/hosts/${hostId}/metrics/history`, params: options })
}

export function getHostsOverallMetrics() {
  return client.get({ url: `${base}/hosts/metrics/overall` })
}

// ==============================
// Host Statistics API
// ==============================

export function getHostStats(params?: any) {
  return client.get({ url: `${base}/dashboard/summary`, params })
}

// 获取主机统计数据（包含完整响应信息）
export async function getHostStatsWithResponse(params?: any) {
  const response = await axios.get(`${base}/dashboard/summary`, { params })
  return response.data  // 返回完整的API响应 {code, data}
}

// 获取主机分布统计
export function getHostDistribution() {
  return client.get({ url: `${base}/dashboard/distribution` })
}

// 获取主机组统计
export function getHostGroupStats() {
  return client.get({ url: `${base}/dashboard/host_groups` })
}

// 获取最近添加的主机列表
export function getRecentHosts(limit?: number) {
  return client.get({ url: `${base}/dashboard/recent_hosts`, params: { limit } })
}

// ==============================
// Host Groups API (alias for compatibility)
// ==============================

export function getHostGroups() {
  return getHostGroupTree()
}

// ==============================
// Export API (alias for compatibility)
// ==============================

export function exportHosts(params?: any) {
  return batchExportHosts(params)
}

// ==============================
// SFTP Operations API
// ==============================

export function listFiles(hostId: number, path: string) {
  return client.get({ url: `${base}/sftp/${hostId}/list`, params: { path } })
}

export function getSftpDownloadUrl(hostId: number, path: string): string {
  return `${import.meta.env.VITE_API_URL}${base}/sftp/${hostId}/download?path=${encodeURIComponent(path)}`
}

export function deleteSftpFile(hostId: number, path: string) {
  return client.del({ url: `${base}/sftp/${hostId}/delete`, params: { path } })
}

export function uploadSftpFile(hostId: number, remotePath: string, file: File) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('path', remotePath)
  return client.post({
    url: `${base}/sftp/${hostId}/upload`,
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}
