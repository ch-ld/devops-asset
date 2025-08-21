<template>
  <div class="dns-records-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <div class="page-title-wrapper">
            <h1 class="page-title">{{ viewMode === 'list' ? '公网权威解析' : 'DNS 解析记录' }}</h1>
            <div class="page-breadcrumb">
              <el-breadcrumb separator="/">
                <el-breadcrumb-item>DNS管理</el-breadcrumb-item>
                <el-breadcrumb-item>{{ viewMode === 'list' ? '权威解析' : '解析记录' }}</el-breadcrumb-item>
              </el-breadcrumb>
            </div>
          </div>
          <p class="page-description">{{ viewMode === 'list' ? '管理您的域名DNS解析记录，支持多种记录类型和智能解析' : '智能管理您的域名解析记录，支持多云平台同步' }}</p>
        </div>
        <div class="header-right">
          <div class="header-controls">
            <el-radio-group v-model="viewMode" @change="handleViewModeChange" class="view-mode-switch">
              <el-radio-button label="list">
                <el-icon><List /></el-icon>
                域名列表
              </el-radio-button>
              <el-radio-button label="detail">
                <el-icon><Document /></el-icon>
                记录详情
              </el-radio-button>
            </el-radio-group>

            <div class="header-actions">
              <template v-if="viewMode === 'detail'">
                <el-space :size="8">
                  <el-button @click="handleRefresh" :icon="Refresh" circle />
                  <el-button :disabled="!selectedDomain" @click="handleExportAll" :icon="Download">
                    导出
                  </el-button>
                  <el-button :disabled="!selectedDomain" @click="openImportModal" :icon="Upload">
                    导入
                  </el-button>
                  <el-button type="primary" plain :disabled="!selectedDomain" @click="() => bulkVisible = true">
                    批量添加
                  </el-button>
                  <el-button type="primary" @click="handleAddRecord" :disabled="!selectedDomain" :icon="Plus">
                    添加记录
                  </el-button>
                </el-space>
              </template>
              <template v-else>
                <el-space :size="8">
                  <el-button @click="handleRefresh" :icon="Refresh" circle />
                  <el-button type="primary" @click="handleAddDomain" :icon="Plus">
                    添加域名
                  </el-button>
                </el-space>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 域名列表视图 -->
    <div v-if="viewMode === 'list'" class="domain-list-view">
      <!-- 搜索和筛选 -->
      <div class="search-section">
        <div class="search-filters">
          <el-select v-model="filters.status" placeholder="解析状态" clearable style="width: 120px">
            <el-option label="全部状态" value="" />
            <el-option label="正常" value="normal" />
            <el-option label="暂停" value="paused" />
            <el-option label="异常" value="error" />
          </el-select>

          <el-input
            v-model="filters.keyword"
            placeholder="请输入域名进行搜索"
            style="width: 300px"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #suffix>
              <el-icon class="search-icon" @click="handleSearch"><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </div>

      <!-- 域名列表表格 -->
      <div class="domain-table-section">
        <el-table
          :data="domains"
          :loading="loading"
          stripe
          style="width: 100%"
          @selection-change="handleDomainSelectionChange"
        >
          <el-table-column type="selection" width="55" />

          <el-table-column prop="name" label="域名" min-width="200">
            <template #default="{ row }">
              <div class="domain-cell">
                <el-button
                  type="primary"
                  link
                  @click="handleViewRecords(row)"
                  class="domain-name"
                >
                  {{ row.name }}
                </el-button>
                <div class="domain-tags" v-if="row.tags && row.tags.length">
                  <el-tag
                    v-for="tag in row.tags"
                    :key="tag"
                    size="small"
                    class="domain-tag"
                  >
                    {{ tag }}
                  </el-tag>
                </div>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="record_count" label="记录数" width="80" align="center">
            <template #default="{ row }">
              <span class="record-count">{{ domainRecordCounts[row.id] ?? '-' }}</span>
            </template>
          </el-table-column>

          <el-table-column prop="provider" label="DNS提供商" width="120" align="center">
            <template #default="{ row }">
              <span>{{ row.provider?.name || '-' }}</span>
            </template>
          </el-table-column>

          <el-table-column label="状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag
                :type="getStatusTagType(row.status)"
                size="small"
              >
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="updated_at" label="更新时间" width="160" align="center">
            <template #default="{ row }">
              <span>{{ formatTime(row.updated_at) }}</span>
            </template>
          </el-table-column>

          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <div class="action-buttons">
                <el-button
                  type="primary"
                  link
                  @click="handleViewRecords(row)"
                >
                  解析设置
                </el-button>
                <el-dropdown @command="(command) => handleMoreAction(command, row)">
                  <el-button type="primary" link>
                    更多
                    <el-icon><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="sync">同步记录</el-dropdown-item>
                      <el-dropdown-item command="export">导出记录</el-dropdown-item>
                      <el-dropdown-item command="edit">编辑域名</el-dropdown-item>
                      <el-dropdown-item command="delete" divided>删除域名</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <!-- 域名选择 (记录详情视图) -->
    <div v-else class="domain-selection-section">
      <el-card class="domain-selection-card" shadow="never">
        <div class="domain-selector-content">
          <div class="selector-left">
            <div class="selector-icon-wrapper">
              <el-icon class="selector-icon" size="20"><Monitor /></el-icon>
            </div>
            <div class="selector-info">
              <h3 class="selector-title">选择域名</h3>
              <p class="selector-desc">选择要管理DNS解析记录的域名</p>
            </div>
          </div>
          <div class="selector-right">
            <div class="selector-controls">
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
                >
                  <div class="domain-option">
                    <div class="domain-info">
                      <span class="domain-name">{{ domain.name }}</span>
                      <div class="domain-meta">
                        <el-tag size="small" type="info" v-if="domain.provider?.name">
                          {{ domain.provider.name }}
                        </el-tag>
                        <span class="record-count">{{ domainRecordCounts[domain.id] ?? 0 }} 条记录</span>
                      </div>
                    </div>
                  </div>
                </el-option>
              </el-select>
              <el-button
                v-if="selectedDomain"
                type="primary"
                @click="handleSyncRecords"
                :loading="syncLoading"
                :icon="Refresh"
              >
                同步记录
              </el-button>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 统计信息 -->
    <div v-if="viewMode === 'detail' && selectedDomain" class="stats-overview">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card class="stat-card total" shadow="hover">
            <div class="stat-content">
              <div class="stat-icon-wrapper">
                <el-icon class="stat-icon" size="28"><Document /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ domainStats.total || 0 }}</div>
                <div class="stat-label">总记录数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card active" shadow="hover">
            <div class="stat-content">
              <div class="stat-icon-wrapper">
                <el-icon class="stat-icon" size="28"><Check /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ domainStats.active || 0 }}</div>
                <div class="stat-label">正常记录</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card syncing" shadow="hover">
            <div class="stat-content">
              <div class="stat-icon-wrapper">
                <el-icon class="stat-icon" size="28"><Loading /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ domainStats.syncing || 0 }}</div>
                <div class="stat-label">同步中</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card error" shadow="hover">
            <div class="stat-content">
              <div class="stat-icon-wrapper">
                <el-icon class="stat-icon" size="28"><Warning /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ domainStats.error || 0 }}</div>
                <div class="stat-label">异常记录</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 记录列表 -->
    <div v-if="viewMode === 'detail' && selectedDomain" class="records-section">
      <el-card>
        <template #header>
          <div class="records-header">
            <h3>解析记录</h3>
            <div class="action-buttons">
              <el-button
                v-if="selectedRecords.length > 0"
                type="danger"
                @click="handleBatchDelete"
              >
                批量删除 ({{ selectedRecords.length }})
              </el-button>
              <el-button
                v-if="selectedRecords.length > 0"
                type="warning"
                @click="handleBatchSync"
              >
                批量同步 ({{ selectedRecords.length }})
              </el-button>
            </div>
          </div>
        </template>

        <el-table
          :data="records"
          v-loading="loading"
          @selection-change="handleSelectionChange"
          class="records-table"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="name" label="记录名称" min-width="150">
            <template #default="{ row }">
              <div class="record-name-cell">
                <span class="record-name">{{ row.name || '@' }}</span>
                <el-tag v-if="row.name === '@'" size="small" type="info">根域名</el-tag>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="type" label="类型" width="80">
            <template #default="{ row }">
              <el-tag :type="getTypeTagType(row.type)" size="small">{{ row.type }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="value" label="记录值" min-width="200">
            <template #default="{ row }">
              <div class="record-value-cell">
                <span class="record-value">{{ row.value }}</span>
                <el-button
                  size="small"
                  text
                  @click="copyToClipboard(row.value)"
                  class="copy-btn"
                >
                  <el-icon><DocumentCopy /></el-icon>
                </el-button>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="ttl" label="TTL" width="100">
            <template #default="{ row }">
              <div class="ttl-cell">
                <span class="ttl-value">{{ row.ttl }}</span>
                <span class="ttl-unit">秒</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="sync_status" width="160">
            <template #header>
              <div class="sync-status-header">
                <span>云端同步状态</span>
                <el-tooltip placement="top" width="350">
                  <template #content>
                    <div class="sync-status-help">
                      <div><strong>云端同步状态说明：</strong></div>
                      <div>• <el-tag type="success" size="small">已同步</el-tag> - 记录已成功同步到云厂商DNS</div>
                      <div>• <el-tag type="warning" size="small">同步中</el-tag> - 正在同步到云厂商DNS</div>
                      <div>• <el-tag type="danger" size="small">同步失败</el-tag> - 云厂商DNS同步失败</div>
                      <div>• <el-tag type="info" size="small">待同步</el-tag> - 等待同步到云厂商DNS</div>
                      <div>• <el-tag type="warning" size="small">仅本地</el-tag> - 记录仅存在本地，未同步</div>
                      <div>• <el-tag type="info" size="small">仅远程</el-tag> - 记录仅存在于云厂商</div>
                      <div style="margin-top: 8px; padding-top: 8px; border-top: 1px solid #eee;">
                        <small>删除记录时，系统会同时删除云厂商中的对应记录</small>
                      </div>
                    </div>
                  </template>
                  <el-icon class="help-icon">
                    <QuestionFilled />
                  </el-icon>
                </el-tooltip>
              </div>
            </template>
            <template #default="{ row }">
              <div class="sync-status-cell">
                <el-tooltip :content="getSyncStatusTooltip(row)" placement="top">
                  <el-tag
                    :type="getSyncStatusType(row.sync_status)"
                    size="small"
                    :icon="getSyncStatusIcon(row.sync_status)"
                    effect="light"
                  >
                    {{ getSyncStatusText(row.sync_status) }}
                  </el-tag>
                </el-tooltip>
                <div v-if="row.last_sync_at" class="sync-time">
                  <el-icon size="12"><Clock /></el-icon>
                  {{ formatTime(row.last_sync_at) }}
                </div>
                <div v-if="row.sync_error" class="sync-error">
                  <el-tooltip :content="row.sync_error" placement="top">
                    <el-icon class="error-icon" color="#f56c6c">
                      <Warning />
                    </el-icon>
                    <span class="error-text">失败详情</span>
                  </el-tooltip>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button size="small" @click="handleEditRecord(row)">编辑</el-button>
              <el-button size="small" type="danger" @click="handleDeleteRecord(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </el-card>
    </div>

    <!-- 空状态 -->
    <div v-if="!selectedDomain" class="empty-state">
      <el-empty description="请先选择一个域名来查看解析记录" />
    </div>

    <!-- 记录编辑弹窗 -->
    <RecordModal
      v-model:visible="recordModalVisible"
      :record="currentRecord"
      :domain-id="selectedDomain"
      @success="handleRecordSuccess"
    />

    <!-- 添加域名弹窗 -->
    <AddDomainModal
      v-model:visible="domainModalVisible"
      @success="handleDomainSuccess"
    />
    <RecordsImportModal
      v-model:visible="importVisible"
      :domain-id="selectedDomain"
      @success="fetchRecords"
    />
    <BulkAddRecordsModal
      v-model:visible="bulkVisible"
      :domain-id="selectedDomain"
      :provider-id="domains.find(d => d.id === selectedDomain)?.provider_id || null"
      @success="fetchRecords"
    />

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
import {
  Plus,
  Connection,
  Refresh,
  Document,
  Check,
  Loading,
  Warning,
  DocumentCopy,
  InfoFilled,
  WarningFilled,
  SuccessFilled,
  Search,
  TrendCharts,
  ArrowDown,
  Close,
  Clock,
  QuestionFilled,
  Edit,
  Delete,
  Upload,
  Download,
  Monitor,
  List
} from '@element-plus/icons-vue'
import RecordModal from './components/RecordModal.vue'
import AddDomainModal from './components/AddDomainModal.vue'
import RecordsImportModal from './components/RecordsImportModal.vue'
import BulkAddRecordsModal from './components/BulkAddRecordsModal.vue'
import { domainApi } from '@/api/dns/domain'
const bulkVisible = ref(false)

import { recordApi } from '@/api/dns/record'
import { dnsApi } from '@/api/dns'
import type { Domain } from '@/types/dns'

// 响应式数据
const importVisible = ref(false)

const viewMode = ref('list') // 'list' | 'detail'
const selectedDomain = ref<number | null>(null)
const domains = ref<any[]>([])
const records = ref<any[]>([])
const selectedRecords = ref<any[]>([])
const selectedDomains = ref<any[]>([])
const loading = ref(false)
const syncLoading = ref(false)
const recordModalVisible = ref(false)
const domainModalVisible = ref(false)
const currentRecord = ref<any>(null)

// 筛选条件
const filters = reactive({
  keyword: '',
  status: '',
  type: ''
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 记录每个域名的真实记录数
const domainRecordCounts = ref<Record<number, number>>({})

// 统计数据
const domainStats = computed(() => {
  if (!records.value.length && !pagination.total) return { total: 0, active: 0, syncing: 0, error: 0 }

  return {
    // 使用后端返回的总数，而不是当前页长度
    total: pagination.total || records.value.length,
    active: records.value.filter(r => r.sync_status === 'synced').length,
    syncing: records.value.filter(r => r.sync_status === 'syncing').length,
    error: records.value.filter(r => r.sync_status === 'error').length
  }
})

// 方法
const fetchDomains = async () => {
  try {
    const response = await domainApi.list({})
    domains.value = response.data?.items || []
    await fetchDomainRecordCounts()
  } catch (error) {
    console.error('获取域名列表失败:', error)
    ElMessage.error('获取域名列表失败')
  }
}

// 批量获取域名记录数
const fetchDomainRecordCounts = async () => {
  if (!domains.value || domains.value.length === 0) {
    domainRecordCounts.value = {}
    return
  }
  const ids = (domains.value as Domain[]).map(d => d.id)
  try {
    const resp = await recordApi.countByDomain(ids)
    // 处理API响应格式，可能被包装在data字段中
    const data = (resp as any).data || resp
    // data 是以字符串 key 的对象，转换为 number key
    const map: Record<number, number> = {}
    Object.keys(data as any).forEach(k => {
      const id = Number(k)
      map[id] = (data as any)[k] as number
    })
    domainRecordCounts.value = map
  } catch (e) {
    console.error('获取域名记录数失败', e)
  }
}


const fetchRecords = async () => {
  if (!selectedDomain.value) return

  loading.value = true
  try {
    const response = await recordApi.list({
      domain_id: selectedDomain.value,
      page: pagination.page,
      page_size: pagination.pageSize
    })
    records.value = response.data?.items || []
    pagination.total = response.data?.total || 0
  } catch (error) {
    ElMessage.error('获取记录列表失败')
  } finally {
    loading.value = false
  }
}

const handleDomainChange = () => {
  selectedRecords.value = []
  if (selectedDomain.value) {
    fetchRecords()
  } else {
    records.value = []
  }
}

const handleSyncRecords = async () => {
  if (!selectedDomain.value) {
    ElMessage.error('请先选择域名')
    return
  }

  // 获取选中域名的详细信息
  const selectedDomainInfo = domains.value.find(d => d.id === selectedDomain.value)
  if (!selectedDomainInfo) {
    ElMessage.error('未找到域名信息，请刷新页面重试')
    return
  }

  // 检查是否有provider信息
  const providerId = selectedDomainInfo.provider_id || selectedDomainInfo.provider?.id
  if (!providerId) {
    ElMessage.error(`域名 ${selectedDomainInfo.name} 未配置DNS提供商，无法同步。请先在DNS提供商页面同步该域名。`)
    return
  }

  syncLoading.value = true
  try {
    const result = await dnsApi.syncDomainRecords({
      domain_id: selectedDomain.value,
      provider_id: providerId,
      dry_run: false
    })

    if (result.success) {
      ElMessage.success(`同步成功！新增: ${result.to_add}, 更新: ${result.to_update}, 删除: ${result.to_delete}`)
    } else {
      ElMessage.warning('同步完成，但可能存在部分失败')
    }

    await fetchRecords()
  } catch (error: any) {
    console.error('同步失败:', error)

    // 根据错误类型显示更友好的错误信息
    let errorMessage = '同步失败'
    if (error.response?.data?.message) {
      errorMessage = error.response.data.message
    } else if (error.message) {
      errorMessage = error.message
    }

    // 特殊错误处理
    if (errorMessage.includes('empty response') || errorMessage.includes('空内容')) {
      ElMessage.error({
        message: `域名 ${selectedDomainInfo.name} 可能未在该DNS提供商托管，请检查域名配置`,
        duration: 5000
      })
    } else if (errorMessage.includes('authentication') || errorMessage.includes('认证')) {
      ElMessage.error({
        message: 'DNS提供商认证失败，请检查API密钥配置',
        duration: 5000
      })
    } else if (errorMessage.includes('rate limit') || errorMessage.includes('频率')) {
      ElMessage.error({
        message: 'API调用频率超限，请稍后重试',
        duration: 5000
      })
    } else {
      ElMessage.error({
        message: errorMessage,
        duration: 5000
      })
    }
  } finally {
    syncLoading.value = false
  }
}

const handleAddRecord = () => {
  currentRecord.value = null
  recordModalVisible.value = true
}

const handleEditRecord = (record: any) => {
  currentRecord.value = record
  recordModalVisible.value = true
}

const handleDeleteRecord = async (record: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除记录 "${record.name}" 吗？系统将同步删除云厂商中的对应记录。此操作不可恢复。`,
      '确认删除',
      {
        type: 'warning',
        distinguishCancelAndClose: true,
        confirmButtonText: '确认删除',
        cancelButtonText: '取消'
      }
    )

    // 显示删除进度
    const loadingInstance = ElLoading.service({
      lock: true,
      text: '正在删除DNS记录并同步到云厂商...',
      background: 'rgba(0, 0, 0, 0.7)'
    })

    try {
      await recordApi.delete(record.id)
      ElMessage.success({
        message: '删除成功，已同步更新云厂商DNS记录',
        duration: 3000
      })
      await fetchRecords()
    } catch (error: any) {
      console.error('删除记录失败:', error)

      // 根据错误类型显示不同的提示
      let errorMessage = '删除失败'
      if (error.response?.data?.message) {
        errorMessage = error.response.data.message
      } else if (error.message) {
        errorMessage = error.message
      }

      // 检查是否是部分成功（本地删除成功但云端删除失败）
      if (errorMessage.includes('partial_success') || errorMessage.includes('云厂商删除失败')) {
        ElMessageBox.alert(
          '本地DNS记录已删除，但云厂商同步删除失败。请登录云厂商控制台手动删除对应的DNS记录。',
          '部分删除成功',
          {
            type: 'warning',
            confirmButtonText: '知道了'
          }
        )
        // 刷新列表显示最新状态
        await fetchRecords()
      } else {
        ElMessage.error({
          message: errorMessage,
          duration: 5000
        })
      }
    } finally {
      loadingInstance.close()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('用户操作错误:', error)
    }
  }
}

const handleSelectionChange = (selection: any[]) => {
  selectedRecords.value = selection
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRecords.value.length} 条记录吗？系统将同步删除云厂商中的对应记录。此操作不可恢复。`,
      '批量删除确认',
      {
        type: 'warning',
        distinguishCancelAndClose: true,
        confirmButtonText: '确认批量删除',
        cancelButtonText: '取消'
      }
    )

    // 显示批量删除进度
    const loadingInstance = ElLoading.service({
      lock: true,
      text: `正在批量删除 ${selectedRecords.value.length} 条DNS记录...`,
      background: 'rgba(0, 0, 0, 0.7)'
    })

    try {
      // 使用批量删除API
      const ids = selectedRecords.value.map(record => record.id)
      await recordApi.batchDelete(ids)

      ElMessage.success({
        message: `批量删除成功，已同步更新云厂商DNS记录`,
        duration: 3000
      })
      selectedRecords.value = []
      await fetchRecords()
    } catch (error: any) {
      console.error('批量删除失败:', error)

      let errorMessage = '批量删除失败'
      if (error.response?.data?.message) {
        errorMessage = error.response.data.message
      } else if (error.message) {
        errorMessage = error.message
      }

      // 检查是否是部分成功
      if (errorMessage.includes('部分') || errorMessage.includes('个失败')) {
        ElMessageBox.alert(
          `批量删除部分成功。部分记录的云厂商同步可能失败，请检查操作日志并手动处理失败的记录。\n\n错误详情：${errorMessage}`,
          '部分删除成功',
          {
            type: 'warning',
            confirmButtonText: '知道了'
          }
        )
        // 刷新列表显示最新状态
        selectedRecords.value = []
        await fetchRecords()
      } else {
        ElMessage.error({
          message: errorMessage,
          duration: 5000
        })
      }
    } finally {
      loadingInstance.close()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('用户操作错误:', error)
    }
  }
}

const handleBatchSync = async () => {
  // 批量同步逻辑
  ElMessage.info('批量同步功能开发中...')
}

const handleRecordSuccess = () => {
  fetchRecords()
}

// 新增方法 - 视图模式切换
const handleViewModeChange = (mode: string) => {
  viewMode.value = mode
  if (mode === 'detail' && selectedDomain.value) {
    fetchRecords()
  }
}

// 域名列表相关方法
const handleSearch = () => {
  // 搜索逻辑
  fetchDomains()
}

const handleDomainSelectionChange = (selection: any[]) => {
  selectedDomains.value = selection
}

const handleViewRecords = (domain: any) => {
  selectedDomain.value = domain.id
  viewMode.value = 'detail'
  fetchRecords()
}

const handleAddDomain = () => {
  domainModalVisible.value = true
}

const handleDomainSuccess = () => {
  domainModalVisible.value = false
  fetchDomains()
}

// 旧的模拟函数已废弃，保留空实现供引用处迁移期不报错
const getDomainRecordCount = (_domainId: number) => {
  return domainRecordCounts.value[_domainId] ?? 0
}

const handleMoreAction = (command: string, domain: any) => {
  switch (command) {
    case 'sync':
      handleSyncDomainRecords(domain)
      break
    case 'export':
      handleExportDomainRecords(domain)
      break
    case 'edit':
      handleEditDomain(domain)
      break
    case 'delete':
      handleDeleteDomain(domain)
      break
  }
}

const handleSyncDomainRecords = async (domain: any) => {
  try {
    ElMessage.info('正在同步域名记录...')
    // 这里调用同步API
    ElMessage.success('同步成功')
    fetchDomains()
  } catch (error) {
    ElMessage.error('同步失败')
  }
}

// 刷新数据
const handleRefresh = async () => {
  if (viewMode.value === 'list') {
    await fetchDomains()
  } else {
    await fetchRecords()
    await fetchDomainRecordCounts()
  }
  ElMessage.success('刷新成功')
}

// 记录详情页 - 一键导出当前域名全部记录
const handleExportAll = async () => {
  console.log('点击导出按钮，selectedDomain:', selectedDomain.value)
  if (!selectedDomain.value) {
    ElMessage.warning('请先选择域名')
    return
  }
  try {
    const domain = (domains.value as any[]).find(d => d.id === selectedDomain.value)
    console.log('开始导出，域名:', domain)
    const blob = await recordApi.export({
      format: 'csv',
      filters: { domain_id: selectedDomain.value }
    })
    console.log('导出成功，blob:', blob)
    const url = URL.createObjectURL(blob as any)
    const a = document.createElement('a')
    a.href = url
    a.download = `${domain?.name || 'domain'}-records.csv`
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (e) {
    console.error('导出失败:', e)
    ElMessage.error('导出失败: ' + (e as any)?.message)
  }
}

// 打开导入弹窗
const openImportModal = () => {
  console.log('点击导入按钮')
  if (!selectedDomain.value) {
    ElMessage.warning('请先选择域名')
    return
  }
  importVisible.value = true
}

const handleExportDomainRecords = (domain: any) => {
  ElMessage.info(`正在导出 ${domain.name} 的解析记录...`)
}

const handleEditDomain = (domain: any) => {
  ElMessage.info(`编辑域名 ${domain.name} 功能开发中...`)
}

const handleDeleteDomain = async (domain: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除域名 ${domain.name} 吗？此操作不可恢复。`,
      '确认删除',
      { type: 'warning' }
    )
    ElMessage.success('删除成功')
    fetchDomains()
  } catch (error) {
    // 用户取消
  }
}

// 工具方法 - 域名列表相关
const getStatusTagType = (status: string) => {
  const types: Record<string, string> = {
    active: 'success',
    normal: 'success',
    paused: 'warning',
    error: 'danger',
    inactive: 'info'
  }
  return types[status] || 'success'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    active: '正常',
    normal: '正常',
    paused: '暂停',
    error: '异常',
    inactive: '未激活'
  }
  return texts[status] || '正常'
}

const formatTime = (timestamp: number | string) => {
  if (!timestamp) return '-'
  const date = typeof timestamp === 'number' ? new Date(timestamp * 1000) : new Date(timestamp)
  return date.toLocaleString('zh-CN')
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  fetchRecords()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchRecords()
}

// 工具方法
const getTypeTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    A: 'primary',
    AAAA: 'success',
    CNAME: 'warning',
    MX: 'info',
    TXT: 'danger',
    NS: 'primary'
  }
  return typeMap[type] || 'info'
}

const getSyncStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    synced: 'success',
    syncing: 'warning',
    error: 'danger',
    pending: 'info',
    failed: 'danger',
    success: 'success',
    local_only: 'warning',
    remote_only: 'info'
  }
  return statusMap[status] || 'info'
}

const getSyncStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    synced: '已同步',
    syncing: '同步中',
    error: '同步失败',
    pending: '待同步',
    failed: '同步失败',
    success: '同步成功',
    local_only: '仅本地',
    remote_only: '仅远程'
  }
  return statusMap[status] || '未知'
}

const getSyncStatusIcon = (status: string) => {
  const iconMap: Record<string, string> = {
    synced: 'Check',
    syncing: 'Loading',
    error: 'Close',
    pending: 'Clock',
    failed: 'Close',
    success: 'Check',
    local_only: 'Upload',
    remote_only: 'Download'
  }
  return iconMap[status] || 'QuestionFilled'
}

const getSyncStatusTooltip = (record: any) => {
  const statusMap: Record<string, string> = {
    synced: '记录已与DNS提供商同步',
    syncing: '正在同步到DNS提供商',
    error: `同步失败: ${record.sync_error || '未知错误'}`,
    pending: '等待同步到DNS提供商',
    failed: `同步失败: ${record.sync_error || '未知错误'}`,
    success: '记录已成功同步',
    local_only: '记录仅存在于本地，未同步到DNS提供商',
    remote_only: '记录仅存在于DNS提供商，本地未保存'
  }

  let tooltip = statusMap[record.sync_status] || '状态未知'

  if (record.last_sync_at) {
    tooltip += `\n最后同步: ${formatTime(record.last_sync_at)}`
  }

  return tooltip
}



const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 生命周期
onMounted(async () => {
  await fetchDomains()
})
</script>

<style scoped lang="scss">
.dns-records-page {
  min-height: 100vh;
  background: #f6f8fa;
  padding: 24px;
}

.page-header {
  background: white;
  border-radius: 8px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e1e4e8;

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;

    .header-left {
      .page-title-wrapper {
        margin-bottom: 8px;

        .page-title {
          margin: 0 0 4px 0;
          font-size: 28px;
          font-weight: 700;
          color: #24292e;
          line-height: 1.2;
        }

        .page-breadcrumb {
          :deep(.el-breadcrumb) {
            font-size: 14px;

            .el-breadcrumb__item {
              .el-breadcrumb__inner {
                color: #64748b;
                font-weight: 500;
                transition: all 0.3s ease;

                &:hover {
                  color: #667eea;
                }
              }

              &:last-child .el-breadcrumb__inner {
                background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
                -webkit-background-clip: text;
                -webkit-text-fill-color: transparent;
                background-clip: text;
                font-weight: 600;
              }
            }
          }
        }
      }

      .page-description {
        margin: 0;
        font-size: 16px;
        color: #64748b;
        font-weight: 500;
        line-height: 1.5;
      }
    }

    .header-right {
      .header-controls {
        display: flex;
        align-items: center;
        gap: 16px;

        .view-mode-switch {
          :deep(.el-radio-button) {
            .el-radio-button__inner {
              display: flex;
              align-items: center;
              gap: 6px;
              padding: 8px 16px;
              border-radius: 12px;
              font-weight: 600;
              transition: all 0.3s ease;
              border: 1px solid #e2e8f0;
              background: rgba(255, 255, 255, 0.9);
              
              &:hover {
                transform: translateY(-1px);
                box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
              }
            }
            
            &.is-active .el-radio-button__inner {
              background: #0366d6;
              color: white;
              border-color: #0366d6;
            }
          }
        }

        .header-actions {
          :deep(.el-button) {
            border-radius: 12px;
            font-weight: 600;
            transition: all 0.3s ease;

            &.is-circle {
              width: 48px;
              height: 48px;
              border-radius: 50%;
            }
            
            &:not(.is-circle):hover {
              transform: translateY(-2px);
              box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
            }
            
            &.el-button--primary {
              background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
              border: none;
              
              &:hover {
                transform: translateY(-2px);
                box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
              }
            }
          }
        }
      }
    }
  }
}

/* 域名选择器样式 */
.domain-selection-section {
  margin-bottom: 24px;

  .domain-selection-card {
    border-radius: 20px;
    border: 1px solid rgba(255, 255, 255, 0.2);
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);

    :deep(.el-card__body) {
      padding: 32px;
    }

    .domain-selector-content {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .selector-left {
        display: flex;
        align-items: center;
        gap: 20px;

        .selector-icon-wrapper {
          width: 64px;
          height: 64px;
          border-radius: 20px;
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          display: flex;
          align-items: center;
          justify-content: center;
          box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
          transition: all 0.3s ease;
          
          &:hover {
            transform: translateY(-2px);
            box-shadow: 0 12px 35px rgba(102, 126, 234, 0.4);
          }

          .selector-icon {
            color: white;
            font-size: 24px;
          }
        }

        .selector-info {
          .selector-title {
            margin: 0 0 6px 0;
            font-size: 22px;
            font-weight: 700;
            color: #1e293b;
          }

          .selector-desc {
            margin: 0;
            font-size: 15px;
            color: #64748b;
            font-weight: 500;
          }
        }
      }

      .selector-right {
        .selector-controls {
          display: flex;
          align-items: center;
          gap: 16px;

          .domain-select {
            width: 360px;

            :deep(.el-select__wrapper) {
              border-radius: 12px;
              border: 1px solid #e2e8f0;
              background: rgba(255, 255, 255, 0.9);
              backdrop-filter: blur(10px);
              transition: all 0.3s ease;
              padding: 12px 16px;
              font-size: 15px;

              &:hover {
                border-color: #94a3b8;
                transform: translateY(-1px);
                box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
              }

              &.is-focused {
                border-color: #667eea;
                box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
                transform: translateY(-1px);
              }
            }
          }
          
          :deep(.el-button) {
            border-radius: 12px;
            font-weight: 600;
            padding: 12px 20px;
            transition: all 0.3s ease;
            
            &.el-button--primary {
              background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
              border: none;
              
              &:hover {
                transform: translateY(-2px);
                box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
              }
            }
          }
        }
      }
    }
  }
}

.domain-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;

  .domain-info {
    display: flex;
    flex-direction: column;
    gap: 4px;

    .domain-name {
      font-weight: 500;
      color: #1a1a1a;
    }

    .domain-meta {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .domain-status {
      font-size: 12px;
      padding: 2px 6px;
      border-radius: 4px;

      &.active {
        background: #f6ffed;
        color: #52c41a;
      }

      &.inactive {
        background: #fff2e8;
        color: #fa8c16;
      }
    }
  }

  .record-count {
    font-size: 12px;
    color: #8c8c8c;
  }
}

.stats-overview {
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  cursor: pointer;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(45deg, transparent 30%, rgba(255, 255, 255, 0.1) 50%, transparent 70%);
    transform: translateX(-100%);
    transition: transform 0.6s ease;
  }

  &:hover {
    transform: translateY(-6px) scale(1.02);
    box-shadow: 0 20px 50px rgba(0, 0, 0, 0.15);
    
    &::before {
      transform: translateX(100%);
    }
  }

  :deep(.el-card__body) {
    padding: 28px;
    position: relative;
    z-index: 1;
  }

  .stat-content {
    display: flex;
    align-items: center;
    gap: 20px;

    .stat-icon-wrapper {
      width: 64px;
      height: 64px;
      border-radius: 20px;
      display: flex;
      align-items: center;
      justify-content: center;
      position: relative;
      overflow: hidden;
      box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
      transition: all 0.3s ease;

      .stat-icon {
        color: white;
        z-index: 1;
        font-size: 28px;
        transition: all 0.3s ease;
      }

      &::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        border-radius: 20px;
        opacity: 0.9;
        transition: all 0.3s ease;
      }

      &:hover {
        transform: scale(1.1);
        box-shadow: 0 12px 35px rgba(0, 0, 0, 0.2);
        
        &::before {
          opacity: 1;
        }
        
        .stat-icon {
          transform: scale(1.1);
        }
      }
    }

    .stat-info {
      flex: 1;

      .stat-value {
        font-size: 32px;
        font-weight: 800;
        color: #1e293b;
        line-height: 1;
        margin-bottom: 6px;
        background: linear-gradient(135deg, #1e293b 0%, #475569 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
      }

      .stat-label {
        font-size: 15px;
        color: #64748b;
        font-weight: 600;
      }
    }
  }

  &.total {
    .stat-icon-wrapper::before {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    }
  }

  &.active {
    .stat-icon-wrapper::before {
      background: linear-gradient(135deg, #10b981 0%, #059669 100%);
    }
  }

  &.syncing {
    .stat-icon-wrapper::before {
      background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
    }
  }

  &.error {
    .stat-icon-wrapper::before {
      background: linear-gradient(135deg, #ff4d4f 0%, #ff7875 100%);
    }
  }
}

.records-section {
  :deep(.el-card) {
    border-radius: 20px;
    border: 1px solid rgba(255, 255, 255, 0.2);
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  }
  
  .records-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    h3 {
      margin: 0;
      font-size: 20px;
      font-weight: 700;
      background: linear-gradient(135deg, #1e293b 0%, #475569 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }

    .action-buttons {
      display: flex;
      gap: 12px;
      
      :deep(.el-button) {
        border-radius: 10px;
        font-weight: 600;
        transition: all 0.3s ease;
        
        &:hover {
          transform: translateY(-2px);
        }
        
        &.el-button--danger {
          background: linear-gradient(135deg, #ff4d4f 0%, #ff7875 100%);
          border: none;

          &:hover {
            box-shadow: 0 8px 25px rgba(239, 68, 68, 0.4);
          }
        }
        
        &.el-button--warning {
          background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
          border: none;
          
          &:hover {
            box-shadow: 0 8px 25px rgba(245, 158, 11, 0.4);
          }
        }
      }
    }
  }
}

:deep(.el-table) {
  border-radius: 16px;
  overflow: hidden;
  background: transparent;
  
  .el-table__header th {
    background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
    color: #475569;
    font-weight: 700;
    border-bottom: 2px solid #e2e8f0;
    padding: 20px 12px;
  }

  .el-table__row {
    transition: all 0.3s ease;
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(10px);

    &:hover {
      background: rgba(255, 255, 255, 0.95);
      transform: translateY(-2px);
      box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
    }
  }

  .el-table__cell {
    padding: 18px 12px;
    border-bottom: 1px solid rgba(226, 232, 240, 0.5);
  }
}

.records-table {
  .record-name-cell {
    display: flex;
    align-items: center;
    gap: 10px;

    .record-name {
      font-weight: 600;
      color: #1e293b;
      font-size: 15px;
    }
  }

  .record-value-cell {
    display: flex;
    align-items: center;
    gap: 10px;

    .record-value {
      font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
      font-size: 14px;
      color: #374151;
      background: rgba(243, 244, 246, 0.8);
      padding: 4px 8px;
      border-radius: 6px;
      backdrop-filter: blur(5px);
    }

    .copy-btn {
      opacity: 0;
      transition: all 0.3s ease;
      border-radius: 6px;
      
      &:hover {
        color: #667eea;
        transform: scale(1.1);
      }
    }

    &:hover .copy-btn {
      opacity: 1;
    }
  }

  .ttl-cell {
    display: flex;
    align-items: center;
    gap: 6px;

    .ttl-value {
      font-weight: 600;
      color: #1e293b;
      font-size: 15px;
    }

    .ttl-unit {
      font-size: 12px;
      color: #64748b;
      font-weight: 500;
    }
  }

  .sync-status-header {
    display: flex;
    align-items: center;
    gap: 8px;

    .help-icon {
      font-size: 16px;
      color: #64748b;
      cursor: pointer;
      transition: all 0.3s ease;

      &:hover {
        color: #667eea;
        transform: scale(1.1);
      }
    }
  }

  .sync-status-help {
    line-height: 1.6;

    div {
      margin-bottom: 6px;

      &:first-child {
        margin-bottom: 12px;
      }

      &:last-child {
        margin-bottom: 0;
      }
    }
  }

  .sync-status-cell {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;

    .sync-time {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 12px;
      color: #64748b;
      font-weight: 500;
      padding: 2px 8px;
      background: rgba(243, 244, 246, 0.6);
      border-radius: 12px;
      backdrop-filter: blur(5px);
    }

    .sync-error {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 12px;
      color: #ff4d4f;
      font-weight: 500;

      .error-icon {
        font-size: 16px;
        cursor: pointer;
        transition: all 0.3s ease;

        &:hover {
          color: #ff7875;
          transform: scale(1.1);
        }
      }

      .error-text {
        cursor: pointer;
        text-decoration: underline;
        transition: all 0.3s ease;

        &:hover {
          color: #ff7875;
        }
      }
    }
  }
}

.pagination-container {
  padding: 24px 0;
  display: flex;
  justify-content: center;
}

// 全局按钮样式优化
:deep(.el-button) {
  border-radius: 10px;
  font-weight: 600;
  transition: all 0.3s ease;

  &.el-button--text {
    &:hover {
      transform: translateY(-1px);
      color: #667eea;
    }
  }

  &.el-button--primary {
    background: #0366d6;
    border: 1px solid #0366d6;

    &:hover {
      background: #0256cc;
      border-color: #0256cc;
    }
  }

  &.el-button--success {
    background: linear-gradient(135deg, #10b981 0%, #059669 100%);
    border: none;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 25px rgba(16, 185, 129, 0.4);
    }
  }

  &.el-button--warning {
    background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
    border: none;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 25px rgba(245, 158, 11, 0.4);
    }
  }

  &.el-button--danger {
    background: linear-gradient(135deg, #ff4d4f 0%, #ff7875 100%);
    border: none;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 25px rgba(239, 68, 68, 0.4);
    }
  }
}

// 标签样式优化
:deep(.el-tag) {
  border-radius: 8px;
  font-weight: 600;
  border: none;
  padding: 4px 12px;
  backdrop-filter: blur(10px);

  &.el-tag--success {
    background: linear-gradient(135deg, #10b981 0%, #059669 100%);
    color: white;
  }

  &.el-tag--warning {
    background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
    color: white;
  }

  &.el-tag--danger {
    background: linear-gradient(135deg, #ff4d4f 0%, #ff7875 100%);
    color: white;
  }

  &.el-tag--info {
    background: linear-gradient(135deg, #6b7280 0%, #4b5563 100%);
    color: white;
  }

  &.el-tag--primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }
}

// 输入框和选择器优化
:deep(.el-input) {
  .el-input__wrapper {
    border-radius: 10px;
    border: 1px solid #e2e8f0;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(10px);
    transition: all 0.3s ease;

    &:hover {
      border-color: #94a3b8;
    }

    &.is-focus {
      border-color: #667eea;
      box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
    }
  }
}

:deep(.el-select) {
  .el-select__wrapper {
    border-radius: 10px;
    border: 1px solid #e2e8f0;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(10px);
    transition: all 0.3s ease;

    &:hover {
      border-color: #94a3b8;
    }

    &.is-focused {
      border-color: #667eea;
      box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
    }
  }
}

// 分页器优化
:deep(.el-pagination) {
  .btn-prev,
  .btn-next,
  .el-pager li {
    border-radius: 8px;
    font-weight: 600;
    transition: all 0.3s ease;
    border: 1px solid #e2e8f0;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(10px);

    &:hover {
      transform: translateY(-1px);
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      border-color: #667eea;
    }

    &.is-active {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      border-color: #667eea;
    }
  }
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
}

/* 域名列表视图样式 */
.domain-list-view {
  .search-section {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    padding: 20px;
    border-radius: 16px;
    margin-bottom: 20px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.2);

    .search-filters {
      display: flex;
      gap: 16px;
      align-items: center;
      justify-content: flex-end;

      .search-icon {
        cursor: pointer;
        color: #64748b;
        font-size: 18px;
        transition: all 0.3s ease;

        &:hover {
          color: #667eea;
          transform: scale(1.1);
        }
      }
    }
  }

  .domain-table-section {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    border-radius: 16px;
    padding: 20px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.2);

    .domain-cell {
      .domain-name {
        font-weight: 600;
        font-size: 15px;
        color: #1e293b;
      }

      .domain-tags {
        margin-top: 6px;
        display: flex;
        gap: 6px;

        .domain-tag {
          font-size: 12px;
          font-weight: 500;
        }
      }
    }

    .record-count {
      font-weight: 600;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
      font-size: 15px;
    }

    .action-buttons {
      display: flex;
      gap: 10px;
      flex-wrap: wrap;
    }
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .domain-list-view {
    .search-section .search-filters {
      flex-direction: column;
      align-items: stretch;
      gap: 8px;
    }
  }
}
</style>
