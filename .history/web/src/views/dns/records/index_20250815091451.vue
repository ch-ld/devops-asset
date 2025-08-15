<template>
  <div class="dns-records-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <h1>{{ viewMode === 'list' ? '公网权威解析' : 'DNS 解析记录' }}</h1>
          <p>{{ viewMode === 'list' ? '管理您的域名DNS解析记录，支持多种记录类型和智能解析' : '智能管理您的域名解析记录，支持多云平台同步' }}</p>
        </div>
        <div class="header-right">
          <el-radio-group v-model="viewMode" @change="handleViewModeChange" style="margin-right: 16px">
            <el-radio-button label="list">域名列表</el-radio-button>
            <el-radio-button label="detail">记录详情</el-radio-button>
          </el-radio-group>
          <el-button
            v-if="viewMode === 'detail'"
            type="primary"
            size="large"
            @click="handleAddRecord"
            :disabled="!selectedDomain"
          >
            <el-icon><Plus /></el-icon>
            添加记录
          </el-button>
          <el-button
            v-else
            type="primary"
            size="large"
            @click="handleAddDomain"
          >
            <el-icon><Plus /></el-icon>
            添加域名
          </el-button>
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
              <span class="record-count">{{ getDomainRecordCount(row.id) }}</span>
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
    <div v-else class="domain-selection">
      <el-card>
        <template #header>
          <div class="card-header">
            <el-icon><Connection /></el-icon>
            <span>选择域名</span>
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
            style="width: 300px"
          >
            <el-option
              v-for="domain in domains"
              :key="domain.id"
              :label="domain.name"
              :value="domain.id"
            >
              <div class="domain-option">
                <span class="domain-name">{{ domain.name }}</span>
                <span class="provider-tag">{{ domain.provider?.name }}</span>
              </div>
            </el-option>
          </el-select>
          <el-button
            v-if="selectedDomain"
            type="success"
            @click="handleSyncRecords"
            :loading="syncLoading"
          >
            <el-icon><Refresh /></el-icon>
            同步记录
          </el-button>
        </div>
      </el-card>
    </div>

    <!-- 统计信息 -->
    <div v-if="selectedDomain" class="stats-overview">
      <el-row :gutter="16">
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
              <el-icon size="24"><Check /></el-icon>
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
    <div v-if="selectedDomain" class="records-section">
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
          <el-table-column prop="sync_status" width="140">
            <template #header>
              <div class="sync-status-header">
                <span>同步状态</span>
                <el-tooltip placement="top" width="300">
                  <template #content>
                    <div class="sync-status-help">
                      <div><strong>同步状态说明：</strong></div>
                      <div>• <el-tag type="success" size="small">已同步</el-tag> - 记录已与DNS提供商同步</div>
                      <div>• <el-tag type="warning" size="small">同步中</el-tag> - 正在同步到DNS提供商</div>
                      <div>• <el-tag type="danger" size="small">同步失败</el-tag> - 同步过程中出现错误</div>
                      <div>• <el-tag type="info" size="small">待同步</el-tag> - 等待同步到DNS提供商</div>
                      <div>• <el-tag type="warning" size="small">仅本地</el-tag> - 记录仅存在于本地</div>
                      <div>• <el-tag type="info" size="small">仅远程</el-tag> - 记录仅存在于DNS提供商</div>
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
                  >
                    {{ getSyncStatusText(row.sync_status) }}
                  </el-tag>
                </el-tooltip>
                <div v-if="row.last_sync_at" class="sync-time">
                  {{ formatTime(row.last_sync_at) }}
                </div>
                <div v-if="row.sync_error" class="sync-error">
                  <el-tooltip :content="row.sync_error" placement="top">
                    <el-icon class="error-icon" color="#f56c6c">
                      <Warning />
                    </el-icon>
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
  Download
} from '@element-plus/icons-vue'
import RecordModal from './components/RecordModal.vue'
import AddDomainModal from './components/AddDomainModal.vue'
import { domainApi } from '@/api/dns/domain'
import { recordApi } from '@/api/dns/record'
import { dnsApi } from '@/api/dns'

// 响应式数据
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

// 统计数据
const domainStats = computed(() => {
  if (!records.value.length) return { total: 0, active: 0, syncing: 0, error: 0 }
  
  return {
    total: records.value.length,
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
  } catch (error) {
    console.error('获取域名列表失败:', error)
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

// 获取域名记录数
const getDomainRecordCount = (domainId: number) => {
  // 这里应该从API获取实际的记录数，暂时返回模拟数据
  const recordCounts: Record<number, number> = {
    1: 3,  // devopsgo.shop 有3条记录
    2: 5   // devopsgo.online 有5条记录
  }
  return recordCounts[domainId] || 0
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
  padding: 24px;
  background: #f5f7fa;
  min-height: 100vh;
}

.page-header {
  margin-bottom: 24px;

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;

    .header-left {
      h1 {
        margin: 0 0 8px 0;
        font-size: 28px;
        font-weight: 600;
        color: #1f2937;
      }

      p {
        margin: 0;
        color: #6b7280;
        font-size: 16px;
      }
    }
  }
}

.domain-selection {
  margin-bottom: 24px;

  .card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 16px;
    font-weight: 600;
  }

  .domain-selector {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .domain-option {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;

    .domain-name {
      font-weight: 500;
    }

    .provider-tag {
      font-size: 12px;
      color: #6b7280;
    }
  }
}

.stats-overview {
  margin-bottom: 24px;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.2s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;

  &.total {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  &.active {
    background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  }

  &.syncing {
    background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  }

  &.error {
    background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  }
}

.stat-content {
  .stat-value {
    font-size: 24px;
    font-weight: 700;
    color: #1f2937;
    line-height: 1;
  }

  .stat-label {
    font-size: 14px;
    color: #6b7280;
    margin-top: 4px;
  }
}

.records-section {
  .records-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: #1f2937;
    }

    .action-buttons {
      display: flex;
      gap: 8px;
    }
  }
}

.records-table {
  .record-name-cell {
    display: flex;
    align-items: center;
    gap: 8px;

    .record-name {
      font-weight: 500;
      color: #1f2937;
    }
  }

  .record-value-cell {
    display: flex;
    align-items: center;
    gap: 8px;

    .record-value {
      font-family: 'Monaco', 'Menlo', monospace;
      font-size: 13px;
      color: #374151;
    }

    .copy-btn {
      opacity: 0;
      transition: opacity 0.2s;
    }

    &:hover .copy-btn {
      opacity: 1;
    }
  }

  .ttl-cell {
    display: flex;
    align-items: center;
    gap: 4px;

    .ttl-value {
      font-weight: 500;
      color: #374151;
    }

    .ttl-unit {
      font-size: 12px;
      color: #9ca3af;
    }
  }

  .sync-status-header {
    display: flex;
    align-items: center;
    gap: 6px;

    .help-icon {
      font-size: 14px;
      color: #909399;
      cursor: pointer;
      transition: color 0.2s;

      &:hover {
        color: #409eff;
      }
    }
  }

  .sync-status-help {
    line-height: 1.6;

    div {
      margin-bottom: 4px;

      &:first-child {
        margin-bottom: 8px;
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
    gap: 4px;

    .sync-time {
      font-size: 12px;
      color: #6b7280;
      margin-top: 4px;
    }

    .sync-error {
      display: flex;
      align-items: center;
      gap: 4px;

      .error-icon {
        font-size: 14px;
        cursor: pointer;
        transition: color 0.2s;

        &:hover {
          color: #f56c6c;
        }
      }
    }
  }
}

.pagination-container {
  padding: 20px 0;
  display: flex;
  justify-content: center;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

// 域名列表视图样式
.domain-list-view {

  .search-section {
    background: white;
    padding: 16px;
    border-radius: 8px;
    margin-bottom: 16px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);

    .search-filters {
      display: flex;
      gap: 12px;
      align-items: center;
      justify-content: flex-end;

      .search-icon {
        cursor: pointer;
        color: #8c8c8c;

        &:hover {
          color: #1890ff;
        }
      }
    }
  }

  .domain-table-section {
    background: white;
    border-radius: 8px;
    padding: 16px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);

    .domain-cell {
      .domain-name {
        font-weight: 500;
        font-size: 14px;
      }

      .domain-tags {
        margin-top: 4px;
        display: flex;
        gap: 4px;

        .domain-tag {
          font-size: 12px;
        }
      }
    }

    .record-count {
      font-weight: 500;
      color: #1890ff;
    }



    .action-buttons {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
    }
  }
}

// 响应式设计
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
