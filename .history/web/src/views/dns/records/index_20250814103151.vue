<template>
  <div class="dns-records-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <div class="page-title">
            <h1>DNS 解析记录</h1>
            <p class="page-subtitle">智能管理您的域名解析记录，支持多云平台同步</p>
          </div>
        </div>
        <div class="header-right">
          <el-button
            type="primary"
            size="large"
            @click="handleAddRecord"
            :disabled="!selectedDomain"
            class="add-record-btn"
          >
            <el-icon><Plus /></el-icon>
            添加记录
          </el-button>
        </div>
      </div>
    </div>

    <!-- 域名选择区域 -->
    <div class="domain-selection-section">
      <el-card class="domain-selector-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <div class="header-icon">
              <el-icon size="20"><Globe /></el-icon>
            </div>
            <div class="header-text">
              <h3>选择域名</h3>
              <p>请先选择要管理的域名，系统将自动同步该域名下的所有解析记录</p>
            </div>
          </div>
        </template>

        <div class="domain-selector">
          <el-select
            v-model="selectedDomain"
            placeholder="请选择要管理的域名"
            size="large"
            filterable
            clearable
            @change="handleDomainChange"
            class="domain-select"
          >
            <el-option
              v-for="domain in domains"
              :key="domain.id"
              :label="domain.name"
              :value="domain.id"
              class="domain-option"
            >
              <div class="domain-option-content">
                <div class="domain-name">{{ domain.name }}</div>
                <div class="domain-info">
                  <el-tag
                    :type="getDomainStatusType(domain.status)"
                    size="small"
                  >
                    {{ getDomainStatusText(domain.status) }}
                  </el-tag>
                  <span class="provider-info">
                    {{ domain.provider?.name || '未配置' }}
                  </span>
                </div>
              </div>
            </el-option>
          </el-select>

          <div class="domain-actions" v-if="selectedDomain">
            <el-button
              type="primary"
              :loading="syncLoading"
              @click="handleSyncDomain"
              class="sync-btn"
            >
              <el-icon><Refresh /></el-icon>
              {{ syncLoading ? '同步中...' : '立即同步' }}
            </el-button>
            <el-button @click="handleDomainSettings">
              <el-icon><Setting /></el-icon>
              域名设置
            </el-button>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 同步状态提示 -->
    <div class="sync-status-section" v-if="selectedDomain && syncStatus">
      <el-alert
        :title="syncStatus.title"
        :description="syncStatus.description"
        :type="syncStatus.type"
        :closable="false"
        show-icon
        class="sync-alert"
      >
        <template #default>
          <div class="sync-details">
            <div class="sync-info">
              <span>{{ syncStatus.title }}</span>
              <span class="sync-time" v-if="syncStatus.lastSyncTime">
                最后同步：{{ formatTime(syncStatus.lastSyncTime) }}
              </span>
            </div>
            <div class="sync-progress" v-if="syncStatus.progress !== undefined">
              <el-progress
                :percentage="syncStatus.progress"
                :status="syncStatus.type === 'error' ? 'exception' : undefined"
                :stroke-width="6"
              />
            </div>
          </div>
        </template>
      </el-alert>
    </div>

    <!-- 统计概览 -->
    <div class="stats-overview" v-if="selectedDomain && domainStats">
      <el-row :gutter="20">
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon total">
              <el-icon size="24"><Document /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ domainStats.total || 0 }}</div>
              <div class="stat-label">总记录数</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon active">
              <el-icon size="24"><CircleCheck /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ domainStats.active || 0 }}</div>
              <div class="stat-label">正常记录</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon syncing">
              <el-icon size="24"><Loading /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ domainStats.syncing || 0 }}</div>
              <div class="stat-label">同步中</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon error">
              <el-icon size="24"><Warning /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ domainStats.error || 0 }}</div>
              <div class="stat-label">异常记录</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 记录列表 -->
    <div class="records-section" v-if="selectedDomain">
      <el-card class="records-card" shadow="never">
        <template #header>
          <div class="records-header">
            <div class="header-left">
              <h3>解析记录</h3>
              <el-tag v-if="selectedDomainInfo" type="info" size="large">
                {{ selectedDomainInfo.name }}
              </el-tag>
            </div>
            <div class="header-right">
              <div class="record-filters">
                <el-select
                  v-model="recordFilter.type"
                  placeholder="记录类型"
                  clearable
                  size="default"
                  style="width: 120px; margin-right: 12px;"
                >
                  <el-option label="全部" value="" />
                  <el-option label="A" value="A" />
                  <el-option label="AAAA" value="AAAA" />
                  <el-option label="CNAME" value="CNAME" />
                  <el-option label="MX" value="MX" />
                  <el-option label="TXT" value="TXT" />
                  <el-option label="NS" value="NS" />
                </el-select>

                <el-select
                  v-model="recordFilter.status"
                  placeholder="同步状态"
                  clearable
                  size="default"
                  style="width: 120px; margin-right: 12px;"
                >
                  <el-option label="全部" value="" />
                  <el-option label="已同步" value="synced" />
                  <el-option label="待同步" value="pending" />
                  <el-option label="同步失败" value="failed" />
                </el-select>

                <el-input
                  v-model="recordFilter.search"
                  placeholder="搜索记录名或值"
                  clearable
                  size="default"
                  style="width: 200px; margin-right: 12px;"
                >
                  <template #prefix>
                    <el-icon><Search /></el-icon>
                  </template>
                </el-input>
              </div>

              <div class="action-buttons">
                <el-button @click="handleRefreshRecords" :loading="loading">
                  <el-icon><Refresh /></el-icon>
                  刷新
                </el-button>
                <el-button @click="handleExportRecords">
                  <el-icon><Download /></el-icon>
                  导出
                </el-button>
              </div>
            </div>
          </div>
        </template>

        <!-- 批量操作栏 -->
        <div class="batch-actions-bar" v-if="selectedRecords.length > 0">
          <div class="batch-info">
            <el-icon><InfoFilled /></el-icon>
            已选择 {{ selectedRecords.length }} 条记录
          </div>
          <div class="batch-buttons">
            <el-button size="small" @click="handleBatchSync">
              <el-icon><Refresh /></el-icon>
              批量同步
            </el-button>
            <el-button size="small" type="danger" @click="handleBatchDelete">
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
          </div>
        </div>

        <el-table
          ref="tableRef"
          :data="records"
          :loading="loading"
          @selection-change="handleSelectionChange"
          row-key="id"
          empty-text="暂无解析记录数据"
          class="records-table"
        >
          <el-table-column type="selection" width="55" />

          <el-table-column prop="name" label="记录名" min-width="150">
            <template #default="{ row }">
              <div class="record-name-cell">
                <div class="record-type-icon">
                  <el-icon :color="getRecordTypeColor(row.type)" size="16">
                    <component :is="getRecordTypeIcon(row.type)" />
                  </el-icon>
                </div>
                <div class="record-name-content">
                  <span class="name">{{ row.name || '@' }}</span>
                  <el-tag v-if="row.name === '@'" type="info" size="small" class="root-tag">根域名</el-tag>
                </div>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="type" label="类型" width="80">
            <template #default="{ row }">
              <el-tag :type="getRecordTypeTagType(row.type)" size="small" class="type-tag">
                {{ row.type }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="value" label="记录值" min-width="250">
            <template #default="{ row }">
              <div class="record-value-cell">
                <div class="value-content">
                  <span :title="row.value" class="value-text">{{ truncateValue(row.value) }}</span>
                  <el-button
                    v-if="row.value && row.value.length > 30"
                    type="primary"
                    text
                    size="small"
                    @click="copyToClipboard(row.value)"
                    class="copy-btn"
                  >
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </div>
                <div class="value-meta" v-if="row.priority || row.weight">
                  <span v-if="row.priority" class="priority">优先级: {{ row.priority }}</span>
                  <span v-if="row.weight" class="weight">权重: {{ row.weight }}</span>
                </div>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="ttl" label="TTL" width="100">
            <template #default="{ row }">
              <div class="ttl-cell">
                <span class="ttl-value">{{ formatTTL(row.ttl) }}</span>
                <span class="ttl-unit">秒</span>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="sync_status" label="同步状态" width="120">
            <template #default="{ row }">
              <div class="sync-status-cell">
                <el-tag
                  :type="getSyncStatusTagType(row.sync_status)"
                  size="small"
                  class="status-tag"
                >
                  <el-icon class="status-icon">
                    <component :is="getSyncStatusIcon(row.sync_status)" />
                  </el-icon>
                  {{ getSyncStatusText(row.sync_status) }}
                </el-tag>
                <div class="sync-time" v-if="row.last_sync_at">
                  {{ formatSyncTime(row.last_sync_at) }}
                </div>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="status" label="状态" width="80">
            <template #default="{ row }">
              <el-switch
                v-model="row.status"
                active-value="active"
                inactive-value="inactive"
                @change="handleStatusChange(row)"
                :loading="row.statusLoading"
              />
            </template>
          </el-table-column>

          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <div class="action-buttons">
                <el-button
                  type="primary"
                  size="small"
                  text
                  @click="handleEditRecord(row)"
                  class="action-btn"
                >
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-button>
                <el-button
                  type="success"
                  size="small"
                  text
                  @click="handleSyncRecord(row)"
                  :loading="row.syncLoading"
                  class="action-btn"
                >
                  <el-icon><Refresh /></el-icon>
                  同步
                </el-button>
                <el-dropdown @command="(command) => handleRecordCommand(command, row)" trigger="click">
                  <el-button type="info" size="small" text class="action-btn">
                    <el-icon><MoreFilled /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="test">
                        <el-icon><Connection /></el-icon>
                        测试解析
                      </el-dropdown-item>
                      <el-dropdown-item command="copy">
                        <el-icon><CopyDocument /></el-icon>
                        复制记录
                      </el-dropdown-item>
                      <el-dropdown-item command="history">
                        <el-icon><Clock /></el-icon>
                        同步历史
                      </el-dropdown-item>
                      <el-dropdown-item divided command="delete" class="danger-item">
                        <el-icon><Delete /></el-icon>
                        删除记录
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-container" v-if="pagination.total > 0">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            class="custom-pagination"
          />
        </div>
      </el-card>
    </div>

    <!-- 空状态 -->
    <div class="empty-state" v-if="!selectedDomain">
      <el-empty
        description="请先选择要管理的域名"
        :image-size="120"
      >
        <template #image>
          <el-icon size="120" color="#dcdfe6"><Globe /></el-icon>
        </template>
        <el-button type="primary" @click="scrollToDomainSelector">
          选择域名
        </el-button>
      </el-empty>
    </div>

    <!-- 记录表单弹窗 -->
    <RecordModal
      v-model:visible="modalVisible"
      :record="currentRecord"
      :domain-id="selectedDomain"
      @success="handleModalSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch, nextTick } from 'vue'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'
import {
  Plus,
  Search,
  Refresh,
  Delete,
  Edit,
  Globe,
  Setting,
  Document,
  CircleCheck,
  Loading,
  Warning,
  CopyDocument,
  MoreFilled,
  Connection,
  Clock,
  InfoFilled
} from '@element-plus/icons-vue'
import RecordModal from './components/RecordModal.vue'
import { recordApi } from '@/api/dns/record'
import { domainApi } from '@/api/dns/domain'
import { dnsApi } from '@/api/dns'
import type { DnsRecord, Domain } from '@/types/dns'
import { formatDistanceToNow } from 'date-fns'
import { zhCN } from 'date-fns/locale'

// 响应式数据
const loading = ref(false)
const syncLoading = ref(false)
const modalVisible = ref(false)
const currentRecord = ref<DnsRecord | null>(null)
const selectedRecords = ref<DnsRecord[]>([])
const records = ref<DnsRecord[]>([])
const domains = ref<Domain[]>([])
const selectedDomain = ref<number | null>(null)
const selectedDomainInfo = ref<Domain | null>(null)

// 同步状态
const syncStatus = ref<{
  title: string
  description: string
  type: 'success' | 'info' | 'warning' | 'error'
  progress?: number
  lastSyncTime?: string
} | null>(null)

// 记录筛选
const recordFilter = reactive({
  type: '',
  status: '',
  search: ''
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 域名统计数据
const domainStats = ref<{
  total: number
  active: number
  syncing: number
  error: number
} | null>(null)

// 计算属性
const filteredRecords = computed(() => {
  let filtered = records.value

  if (recordFilter.type) {
    filtered = filtered.filter(record => record.type === recordFilter.type)
  }

  if (recordFilter.status) {
    filtered = filtered.filter(record => record.sync_status === recordFilter.status)
  }

  if (recordFilter.search) {
    const search = recordFilter.search.toLowerCase()
    filtered = filtered.filter(record =>
      record.name?.toLowerCase().includes(search) ||
      record.value?.toLowerCase().includes(search)
    )
  }

  return filtered
})

// 工具方法
const getRecordTypeIcon = (type: string) => {
  const iconMap: Record<string, any> = {
    A: 'Connection',
    AAAA: 'Connection',
    CNAME: 'Link',
    MX: 'Message',
    TXT: 'Document',
    NS: 'Globe',
    SRV: 'Service'
  }
  return iconMap[type] || 'Document'
}

const getRecordTypeColor = (type: string) => {
  const colorMap: Record<string, string> = {
    A: '#67c23a',
    AAAA: '#67c23a',
    CNAME: '#e6a23c',
    MX: '#f56c6c',
    TXT: '#409eff',
    NS: '#909399',
    SRV: '#909399'
  }
  return colorMap[type] || '#909399'
}

const getRecordTypeTagType = (type: string) => {
  const typeMap: Record<string, any> = {
    A: 'success',
    AAAA: 'success',
    CNAME: 'warning',
    MX: 'danger',
    TXT: 'primary',
    NS: 'info',
    SRV: 'info'
  }
  return typeMap[type] || 'info'
}

const getSyncStatusIcon = (status: string) => {
  const iconMap: Record<string, any> = {
    synced: 'CircleCheck',
    pending: 'Clock',
    syncing: 'Loading',
    failed: 'Warning'
  }
  return iconMap[status] || 'Clock'
}

const getSyncStatusTagType = (status: string) => {
  const statusMap: Record<string, any> = {
    synced: 'success',
    pending: 'warning',
    syncing: 'primary',
    failed: 'danger'
  }
  return statusMap[status] || 'info'
}

const getSyncStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    synced: '已同步',
    pending: '待同步',
    syncing: '同步中',
    failed: '同步失败'
  }
  return statusMap[status] || status
}

const getDomainStatusType = (status: string) => {
  const statusMap: Record<string, any> = {
    active: 'success',
    inactive: 'info',
    expired: 'danger'
  }
  return statusMap[status] || 'info'
}

const getDomainStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '正常',
    inactive: '未激活',
    expired: '已过期'
  }
  return statusMap[status] || status
}

const truncateValue = (value: string) => {
  if (!value) return '-'
  return value.length > 30 ? value.substring(0, 30) + '...' : value
}

const formatTTL = (ttl: number) => {
  if (!ttl) return '-'
  if (ttl >= 3600) {
    return `${Math.floor(ttl / 3600)}h`
  } else if (ttl >= 60) {
    return `${Math.floor(ttl / 60)}m`
  }
  return `${ttl}s`
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return formatDistanceToNow(new Date(time), {
    addSuffix: true,
    locale: zhCN
  })
}

const formatSyncTime = (time: string) => {
  if (!time) return '未同步'
  return formatDistanceToNow(new Date(time), {
    addSuffix: true,
    locale: zhCN
  })
}

const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 事件处理方法
const handleDomainChange = async (domainId: number | null) => {
  if (!domainId) {
    selectedDomainInfo.value = null
    records.value = []
    domainStats.value = null
    syncStatus.value = null
    return
  }

  // 获取域名信息
  const domain = domains.value.find(d => d.id === domainId)
  selectedDomainInfo.value = domain || null

  // 自动同步该域名的记录
  await handleSyncDomain()
}

const handleSyncDomain = async () => {
  if (!selectedDomain.value) return

  syncLoading.value = true
  syncStatus.value = {
    title: '正在同步域名记录...',
    description: '正在从DNS提供商获取最新记录',
    type: 'info',
    progress: 0
  }

  try {
    // 模拟同步进度
    const progressInterval = setInterval(() => {
      if (syncStatus.value && syncStatus.value.progress !== undefined) {
        syncStatus.value.progress = Math.min(syncStatus.value.progress + 10, 90)
      }
    }, 200)

    // 调用同步API
    const response = await dnsApi.syncDomainRecords({
      domain_id: selectedDomain.value,
      provider_id: selectedDomainInfo.value?.provider_id,
      dry_run: false
    })

    clearInterval(progressInterval)

    syncStatus.value = {
      title: '同步完成',
      description: `成功同步 ${response.data.to_add + response.data.to_update} 条记录`,
      type: 'success',
      progress: 100,
      lastSyncTime: new Date().toISOString()
    }

    // 刷新记录列表
    await fetchRecords()

    ElNotification({
      title: '同步成功',
      message: `域名 ${selectedDomainInfo.value?.name} 的记录已同步完成`,
      type: 'success'
    })

  } catch (error) {
    syncStatus.value = {
      title: '同步失败',
      description: '同步过程中发生错误，请稍后重试',
      type: 'error'
    }
    ElMessage.error('同步失败')
  } finally {
    syncLoading.value = false
  }
}

const handleAddRecord = () => {
  if (!selectedDomain.value) {
    ElMessage.warning('请先选择域名')
    return
  }
  currentRecord.value = null
  modalVisible.value = true
}

const handleEditRecord = (record: DnsRecord) => {
  currentRecord.value = record
  modalVisible.value = true
}

const handleSyncRecord = async (record: DnsRecord) => {
  record.syncLoading = true
  try {
    await recordApi.sync(record.id)
    ElMessage.success('记录同步成功')
    await fetchRecords()
  } catch (error) {
    ElMessage.error('记录同步失败')
  } finally {
    record.syncLoading = false
  }
}

const handleStatusChange = async (record: DnsRecord) => {
  record.statusLoading = true
  try {
    await recordApi.updateStatus(record.id, record.status)
    ElMessage.success('状态更新成功')
  } catch (error) {
    ElMessage.error('状态更新失败')
    // 恢复原状态
    record.status = record.status === 'active' ? 'inactive' : 'active'
  } finally {
    record.statusLoading = false
  }
}

const handleRecordCommand = async (command: string, record: DnsRecord) => {
  switch (command) {
    case 'test':
      await testRecord(record)
      break
    case 'copy':
      await copyRecord(record)
      break
    case 'history':
      viewSyncHistory(record)
      break
    case 'delete':
      await deleteRecord(record)
      break
  }
}

const testRecord = async (record: DnsRecord) => {
  try {
    const result = await recordApi.test(record.id)
    ElNotification({
      title: '解析测试结果',
      message: result.success ? '解析正常' : '解析异常',
      type: result.success ? 'success' : 'error'
    })
  } catch (error) {
    ElMessage.error('测试失败')
  }
}

const copyRecord = async (record: DnsRecord) => {
  const recordText = `${record.name} ${record.type} ${record.value} ${record.ttl}`
  await copyToClipboard(recordText)
}

const viewSyncHistory = (record: DnsRecord) => {
  // TODO: 打开同步历史弹窗
  ElMessage.info('同步历史功能开发中')
}

const deleteRecord = async (record: DnsRecord) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除记录 "${record.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        type: 'warning',
        confirmButtonText: '确定删除',
        cancelButtonText: '取消'
      }
    )

    await recordApi.delete(record.id)
    ElMessage.success('记录删除成功')
    await fetchRecords()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('记录删除失败')
    }
  }
}

const handleSelectionChange = (selection: DnsRecord[]) => {
  selectedRecords.value = selection
}

const handleBatchSync = async () => {
  if (selectedRecords.value.length === 0) {
    ElMessage.warning('请先选择要同步的记录')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要同步选中的 ${selectedRecords.value.length} 条记录吗？`,
      '批量同步确认',
      { type: 'info' }
    )

    const promises = selectedRecords.value.map(record => recordApi.sync(record.id))
    await Promise.all(promises)

    ElMessage.success('批量同步请求已提交')
    await fetchRecords()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量同步失败')
    }
  }
}

const handleBatchDelete = async () => {
  if (selectedRecords.value.length === 0) {
    ElMessage.warning('请先选择要删除的记录')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRecords.value.length} 条记录吗？此操作不可恢复。`,
      '批量删除确认',
      {
        type: 'warning',
        confirmButtonText: '确定删除',
        cancelButtonText: '取消'
      }
    )

    const promises = selectedRecords.value.map(record => recordApi.delete(record.id))
    await Promise.all(promises)

    ElMessage.success('批量删除成功')
    await fetchRecords()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

const handleRefreshRecords = async () => {
  await fetchRecords()
}

const handleExportRecords = () => {
  // TODO: 实现记录导出功能
  ElMessage.info('导出功能开发中')
}

const handleDomainSettings = () => {
  if (selectedDomainInfo.value) {
    // TODO: 跳转到域名设置页面
    ElMessage.info('域名设置功能开发中')
  }
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchRecords()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchRecords()
}

const handleModalSuccess = () => {
  modalVisible.value = false
  fetchRecords()
}

const scrollToDomainSelector = () => {
  const element = document.querySelector('.domain-selector-card')
  element?.scrollIntoView({ behavior: 'smooth' })
}

// 数据获取方法
const fetchDomains = async () => {
  try {
    const response = await domainApi.list({ page: 1, pageSize: 100 })
    domains.value = response.data.items || []
  } catch (error) {
    ElMessage.error('获取域名列表失败')
  }
}

const fetchRecords = async () => {
  if (!selectedDomain.value) return

  loading.value = true
  try {
    const response = await recordApi.list({
      domain_id: selectedDomain.value,
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...recordFilter
    })

    records.value = response.data.items || []
    pagination.total = response.data.total || 0

    // 更新统计数据
    updateDomainStats()
  } catch (error) {
    ElMessage.error('获取记录列表失败')
  } finally {
    loading.value = false
  }
}

const updateDomainStats = () => {
  if (!records.value.length) {
    domainStats.value = { total: 0, active: 0, syncing: 0, error: 0 }
    return
  }

  domainStats.value = {
    total: records.value.length,
    active: records.value.filter(r => r.sync_status === 'synced').length,
    syncing: records.value.filter(r => r.sync_status === 'syncing').length,
    error: records.value.filter(r => r.sync_status === 'failed').length
  }
}

// 监听器
watch(selectedDomain, (newDomainId) => {
  if (newDomainId) {
    fetchRecords()
  }
})

watch(recordFilter, () => {
  pagination.page = 1
  fetchRecords()
}, { deep: true })

// 生命周期钩子
onMounted(async () => {
  await fetchDomains()
})
</script>

<style scoped lang="scss">
.dns-records-page {
  padding: 24px;
  background: #f5f7fa;
  min-height: 100vh;

  .page-header {
    margin-bottom: 24px;

    .header-content {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;

      .header-left {
        .page-title {
          h1 {
            font-size: 28px;
            font-weight: 600;
            color: #1f2937;
            margin: 0 0 8px 0;
          }

          .page-subtitle {
            font-size: 16px;
            color: #6b7280;
            margin: 0;
          }
        }
      }

      .header-right {
        .add-record-btn {
          height: 44px;
          padding: 0 24px;
          font-size: 16px;
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(59, 130, 246, 0.15);

          &:hover {
            transform: translateY(-1px);
            box-shadow: 0 4px 8px rgba(59, 130, 246, 0.25);
          }
        }
      }
    }
  }

  .domain-selection-section {
    margin-bottom: 24px;

    .domain-selector-card {
      border-radius: 12px;
      border: 1px solid #e5e7eb;

      :deep(.el-card__header) {
        padding: 20px 24px;
        border-bottom: 1px solid #f3f4f6;

        .card-header {
          display: flex;
          align-items: center;
          gap: 16px;

          .header-icon {
            width: 48px;
            height: 48px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 12px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
          }

          .header-text {
            h3 {
              font-size: 20px;
              font-weight: 600;
              color: #1f2937;
              margin: 0 0 4px 0;
            }

            p {
              font-size: 14px;
              color: #6b7280;
              margin: 0;
            }
          }
        }
      }

      :deep(.el-card__body) {
        padding: 24px;
      }
    }

    .domain-selector {
      display: flex;
      align-items: center;
      gap: 16px;

      .domain-select {
        flex: 1;
        max-width: 400px;

        :deep(.el-select__wrapper) {
          border-radius: 8px;
          border: 2px solid #e5e7eb;
          transition: all 0.3s ease;

          &:hover {
            border-color: #d1d5db;
          }

          &.is-focused {
            border-color: #3b82f6;
            box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
          }
        }
      }

      .domain-actions {
        display: flex;
        gap: 12px;

        .sync-btn {
          background: linear-gradient(135deg, #10b981 0%, #059669 100%);
          border: none;
          color: white;
          padding: 0 20px;
          height: 40px;
          border-radius: 8px;
          font-weight: 500;

          &:hover {
            transform: translateY(-1px);
            box-shadow: 0 4px 8px rgba(16, 185, 129, 0.25);
          }
        }
      }
    }

    .domain-option-content {
      display: flex;
      justify-content: space-between;
      align-items: center;
      width: 100%;

      .domain-name {
        font-weight: 500;
        color: #1f2937;
      }

      .domain-info {
        display: flex;
        align-items: center;
        gap: 8px;

        .provider-info {
          font-size: 12px;
          color: #6b7280;
        }
      }
    }
  }
      page: pagination.page,
      page_size: pagination.pageSize,
      ...searchForm
    }
    const response = await recordApi.list(params)

    // 处理后端返回的包装响应结构
    const data = (response as any).data || response
    records.value = data.items || data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取记录列表失败:', error)
    ElMessage.error('获取记录列表失败')
  } finally {
    loading.value = false
  }
}

const fetchDomains = async () => {
  try {
    const response = await domainApi.list({ page: 1, page_size: 100 })
    // 处理后端返回的包装响应结构
    const data = (response as any).data || response
    domains.value = data.items || data.list || []
  } catch (error) {
    console.error('获取域名列表失败:', error)
  }
}

const fetchStatistics = async () => {
  try {
    // 简化统计实现
    const total = records.value.length
    const active = records.value.filter(item => item.status === 'active').length
    const syncing = records.value.filter(item => item.status === 'syncing').length
    const error = records.value.filter(item => item.status === 'error').length
    
    statistics.value = { total, active, syncing, error }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 生命周期
onMounted(() => {
  fetchData()
  fetchDomains()
  fetchStatistics()
})
</script>

<style scoped lang="scss">
.dns-record-container {
  padding: 24px;
  background: #f5f5f5;
  min-height: 100vh;
}

.page-header {
  margin-bottom: 24px;
  
  .page-header-content {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    
    .page-title {
      h1 {
        margin: 0 0 8px 0;
        font-size: 24px;
        font-weight: 600;
        color: #262626;
      }
      
      p {
        margin: 0;
        color: #8c8c8c;
        font-size: 14px;
      }
    }
  }
}

.stats-container {
  margin-bottom: 24px;
}

.search-container {
  margin-bottom: 24px;
}

.table-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}

.records-table {
  .domain-cell {
    .domain-name {
      font-weight: 600;
      color: #1890ff;
      margin-bottom: 2px;
    }
    
    .domain-id {
      font-size: 11px;
      color: #8c8c8c;
      font-family: monospace;
    }
  }

  .provider-cell {
    .el-tag {
      font-size: 11px;
    }
  }

  .record-name {
    display: flex;
    align-items: center;
    gap: 8px;

    .name {
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
      font-size: 13px;
      color: #262626;
    }

    .root-tag {
      margin-left: 4px;
    }
  }

  .record-value {
    display: flex;
    align-items: center;
    gap: 8px;

    .value-text {
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
      font-size: 12px;
      color: #595959;
      max-width: 180px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .copy-btn {
      opacity: 0;
      transition: opacity 0.2s;
      padding: 2px 4px;
      min-height: auto;
    }

    &:hover .copy-btn {
      opacity: 1;
    }
  }

  .ttl-value, .priority-value {
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 12px;
    color: #8c8c8c;
  }
}

// 空状态样式
.el-table__empty-block {
  padding: 60px 0;

  .el-table__empty-text {
    color: #8c8c8c;
    font-size: 14px;
  }
}
</style>
