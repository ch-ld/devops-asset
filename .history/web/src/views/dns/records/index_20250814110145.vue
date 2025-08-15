<template>
  <div class="dns-records-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <h1>DNS 解析记录</h1>
          <p>智能管理您的域名解析记录，支持多云平台同步</p>
        </div>
        <div class="header-right">
          <el-button
            type="primary"
            size="large"
            @click="handleAddRecord"
            :disabled="!selectedDomain"
          >
            <el-icon><Plus /></el-icon>
            添加记录
          </el-button>
        </div>
      </div>
    </div>

    <!-- 域名选择 -->
    <div class="domain-selection">
      <el-card>
        <template #header>
          <div class="card-header">
            <el-icon><Globe /></el-icon>
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
                  <el-icon><CopyDocument /></el-icon>
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
          <el-table-column prop="sync_status" label="同步状态" width="120">
            <template #default="{ row }">
              <div class="sync-status-cell">
                <el-tag :type="getSyncStatusType(row.sync_status)" size="small">
                  {{ getSyncStatusText(row.sync_status) }}
                </el-tag>
                <div v-if="row.last_sync_at" class="sync-time">
                  {{ formatTime(row.last_sync_at) }}
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Connection,
  Refresh,
  Document,
  Check,
  Loading,
  Warning,
  CopyDocument
} from '@element-plus/icons-vue'
import RecordModal from './components/RecordModal.vue'
import { getDomains } from '@/api/dns'
import { getRecords, deleteRecord, batchSyncRecords } from '@/api/dns/record'

// 响应式数据
const selectedDomain = ref<number | null>(null)
const domains = ref<any[]>([])
const records = ref<any[]>([])
const selectedRecords = ref<any[]>([])
const loading = ref(false)
const syncLoading = ref(false)
const recordModalVisible = ref(false)
const currentRecord = ref<any>(null)

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
    const response = await getDomains()
    domains.value = response.data.items || []
  } catch (error) {
    ElMessage.error('获取域名列表失败')
  }
}

const fetchRecords = async () => {
  if (!selectedDomain.value) return
  
  loading.value = true
  try {
    const response = await getRecords({
      domain_id: selectedDomain.value,
      page: pagination.page,
      page_size: pagination.pageSize
    })
    records.value = response.data.items || []
    pagination.total = response.data.total || 0
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
  if (!selectedDomain.value) return
  
  syncLoading.value = true
  try {
    await batchSyncRecords({
      domain_id: selectedDomain.value,
      dry_run: false
    })
    ElMessage.success('同步成功')
    await fetchRecords()
  } catch (error) {
    ElMessage.error('同步失败')
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
      `确定要删除记录 "${record.name}" 吗？此操作不可恢复。`,
      '确认删除',
      { type: 'warning' }
    )
    
    await deleteRecord(record.id)
    ElMessage.success('删除成功')
    await fetchRecords()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleSelectionChange = (selection: any[]) => {
  selectedRecords.value = selection
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRecords.value.length} 条记录吗？此操作不可恢复。`,
      '批量删除确认',
      { type: 'warning' }
    )
    
    // 批量删除逻辑
    for (const record of selectedRecords.value) {
      await deleteRecord(record.id)
    }
    
    ElMessage.success('批量删除成功')
    selectedRecords.value = []
    await fetchRecords()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
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
    pending: 'info'
  }
  return statusMap[status] || 'info'
}

const getSyncStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    synced: '已同步',
    syncing: '同步中',
    error: '同步失败',
    pending: '待同步'
  }
  return statusMap[status] || '未知'
}

const formatTime = (time: string) => {
  if (!time) return ''
  return new Date(time).toLocaleString()
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

  .sync-status-cell {
    .sync-time {
      font-size: 12px;
      color: #6b7280;
      margin-top: 4px;
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
</style>
