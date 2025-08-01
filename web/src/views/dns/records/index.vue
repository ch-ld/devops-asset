<template>
  <div class="dns-record-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title">
          <h1>解析记录</h1>
          <p>管理DNS解析记录，包括A、AAAA、CNAME、MX、TXT等记录类型</p>
        </div>
        <div class="page-actions">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加记录
          </el-button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-container">
      <el-row :gutter="16">
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="总记录数"
              :value="statistics.total"
              :value-style="{ color: '#409eff' }"
            >
              <template #suffix>
                <el-icon><Document /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="正常记录"
              :value="statistics.active"
              :value-style="{ color: '#67c23a' }"
            >
              <template #suffix>
                <el-icon><CircleCheck /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="同步中"
              :value="statistics.syncing"
              :value-style="{ color: '#e6a23c' }"
            >
              <template #suffix>
                <el-icon><Loading /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="异常记录"
              :value="statistics.error"
              :value-style="{ color: '#f56c6c' }"
            >
              <template #suffix>
                <el-icon><Warning /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-container">
      <el-card shadow="never">
        <el-form
          ref="searchFormRef"
          :model="searchForm"
          inline
          class="search-form"
        >
          <el-form-item label="域名" prop="domain_id">
            <el-select
              v-model="searchForm.domain_id"
              placeholder="请选择域名"
              clearable
              style="width: 200px"
            >
              <el-option label="全部" value="" />
              <el-option 
                v-for="domain in domains" 
                :key="domain.id" 
                :label="domain.name" 
                :value="domain.id" 
              />
            </el-select>
          </el-form-item>
          <el-form-item label="记录名" prop="keyword">
            <el-input
              v-model="searchForm.keyword"
              placeholder="请输入记录名关键词"
              clearable
              style="width: 200px"
            />
          </el-form-item>
          <el-form-item label="记录类型" prop="type">
            <el-select
              v-model="searchForm.type"
              placeholder="请选择记录类型"
              clearable
              style="width: 120px"
            >
              <el-option label="全部" value="" />
              <el-option label="A" value="A" />
              <el-option label="AAAA" value="AAAA" />
              <el-option label="CNAME" value="CNAME" />
              <el-option label="MX" value="MX" />
              <el-option label="TXT" value="TXT" />
              <el-option label="NS" value="NS" />
              <el-option label="SRV" value="SRV" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-select
              v-model="searchForm.status"
              placeholder="请选择状态"
              clearable
              style="width: 120px"
            >
              <el-option label="全部" value="" />
              <el-option label="正常" value="active" />
              <el-option label="异常" value="error" />
              <el-option label="同步中" value="syncing" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button @click="handleReset">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>

    <!-- 记录列表 -->
    <div class="table-container">
      <el-card shadow="never">
        <template #header>
          <div class="card-header">
            <span>解析记录列表</span>
            <div>
              <el-button @click="handleRefresh">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
              <el-button @click="handleBatchSync" :disabled="!hasSelected">
                <el-icon><Refresh /></el-icon>
                批量同步
              </el-button>
              <el-button @click="handleBatchDelete" :disabled="!hasSelected">
                <el-icon><Delete /></el-icon>
                批量删除
              </el-button>
            </div>
          </div>
        </template>

        <el-table
          ref="tableRef"
          :data="records"
          :loading="loading"
          @selection-change="handleSelectionChange"
          row-key="id"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="name" label="记录名" min-width="150">
            <template #default="{ row }">
              <div class="record-name">
                {{ row.name || '@' }}
                <el-tag v-if="row.name === '@'" type="info" size="small">根域名</el-tag>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="type" label="类型" width="80">
            <template #default="{ row }">
              <el-tag :type="getRecordTypeColor(row.type)" size="small">
                {{ row.type }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="value" label="记录值" min-width="200">
            <template #default="{ row }">
              <div class="record-value">
                <span :title="row.value">{{ row.value }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="ttl" label="TTL" width="80" />
          <el-table-column prop="priority" label="优先级" width="80">
            <template #default="{ row }">
              {{ row.priority || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="sync_status" label="同步状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getSyncStatusType(row.sync_status)" size="small">
                {{ getSyncStatusText(row.sync_status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="last_sync_at" label="最后同步" width="140">
            <template #default="{ row }">
              {{ formatDate(row.last_sync_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button
                type="primary"
                size="small"
                text
                @click="handleEdit(row)"
              >
                编辑
              </el-button>
              <el-button
                type="success"
                size="small"
                text
                @click="handleSync(row)"
                :disabled="row.sync_status === 'syncing'"
              >
                同步
              </el-button>
              <el-dropdown @command="(command) => handleCommand(command, row)">
                <el-button type="primary" size="small" text>
                  更多
                  <el-icon><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="validate">验证记录</el-dropdown-item>
                    <el-dropdown-item command="history">变更历史</el-dropdown-item>
                    <el-dropdown-item command="copy">复制记录</el-dropdown-item>
                    <el-dropdown-item command="delete" divided>删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
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

    <!-- 记录表单弹窗 -->
    <RecordModal
      v-model:visible="modalVisible"
      :record="currentRecord"
      :domains="domains"
      @success="handleModalSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, 
  Search, 
  Refresh, 
  Delete, 
  ArrowDown,
  Document,
  CircleCheck, 
  Loading,
  Warning
} from '@element-plus/icons-vue'
import RecordModal from './components/RecordModal.vue'
import { recordApi } from '@/api/dns/record'
import { domainApi } from '@/api/dns/domain'
import type { Record, Domain } from '@/types/dns'

// 响应式数据
const loading = ref(false)
const modalVisible = ref(false)
const currentRecord = ref<Record | null>(null)
const records = ref<Record[]>([])
const domains = ref<Domain[]>([])
const statistics = ref({
  total: 0,
  active: 0,
  syncing: 0,
  error: 0
})

// 搜索表单
const searchFormRef = ref()
const searchForm = reactive({
  domain_id: '',
  keyword: '',
  type: '',
  status: ''
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 选择
const selectedRows = ref<Record[]>([])
const hasSelected = computed(() => selectedRows.value.length > 0)

// 工具方法
const getRecordTypeColor = (type: string) => {
  const colorMap: Record<string, string> = {
    A: 'success',
    AAAA: 'success',
    CNAME: 'warning',
    MX: 'danger',
    TXT: 'info',
    NS: 'primary',
    SRV: 'info'
  }
  return colorMap[type] || 'info'
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, any> = {
    active: 'success',
    error: 'danger',
    syncing: 'warning'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '正常',
    error: '异常',
    syncing: '同步中'
  }
  return statusMap[status] || status
}

const getSyncStatusType = (syncStatus: string) => {
  const statusMap: Record<string, any> = {
    synced: 'success',
    pending: 'warning',
    failed: 'danger'
  }
  return statusMap[syncStatus] || 'info'
}

const getSyncStatusText = (syncStatus: string) => {
  const statusMap: Record<string, string> = {
    synced: '已同步',
    pending: '待同步',
    failed: '同步失败'
  }
  return statusMap[syncStatus] || syncStatus
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

// 事件处理
const handleAdd = () => {
  currentRecord.value = null
  modalVisible.value = true
}

const handleEdit = (row: Record) => {
  currentRecord.value = row
  modalVisible.value = true
}

const handleSync = async (row: Record) => {
  try {
    await recordApi.sync(row.id)
    ElMessage.success(`记录 ${row.name} 同步请求已提交`)
    await fetchData()
  } catch (error) {
    ElMessage.error('记录同步失败')
  }
}

const handleCommand = (command: string, row: Record) => {
  switch (command) {
    case 'validate':
      ElMessage.info(`验证记录 ${row.name}`)
      break
    case 'history':
      ElMessage.info(`查看记录 ${row.name} 的变更历史`)
      break
    case 'copy':
      ElMessage.info(`复制记录 ${row.name}`)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const handleDelete = async (row: Record) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除记录 "${row.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await recordApi.delete(row.id)
    ElMessage.success('删除成功')
    await fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleBatchSync = async () => {
  if (!hasSelected.value) return
  
  try {
    await ElMessageBox.confirm(
      `确定要同步选中的 ${selectedRows.value.length} 条记录吗？`,
      '确认批量同步',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const ids = selectedRows.value.map(row => row.id)
    await recordApi.batchSync(ids)
    ElMessage.success('批量同步请求已提交')
    selectedRows.value = []
    await fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量同步失败')
    }
  }
}

const handleBatchDelete = async () => {
  if (!hasSelected.value) return
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 条记录吗？此操作不可恢复。`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const ids = selectedRows.value.map(row => row.id)
    await recordApi.batchDelete(ids)
    ElMessage.success('批量删除成功')
    selectedRows.value = []
    await fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

const handleSelectionChange = (rows: Record[]) => {
  selectedRows.value = rows
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchFormRef.value?.resetFields()
  pagination.page = 1
  fetchData()
}

const handleRefresh = () => {
  fetchData()
  fetchStatistics()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchData()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchData()
}

const handleModalSuccess = () => {
  modalVisible.value = false
  fetchData()
  fetchStatistics()
}

// 数据获取
const fetchData = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...searchForm
    }
    const response = await recordApi.list(params)
    records.value = response.list || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('获取记录列表失败')
  } finally {
    loading.value = false
  }
}

const fetchDomains = async () => {
  try {
    const response = await domainApi.list({ page: 1, page_size: 100 })
    domains.value = response.list || []
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

.record-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.record-value {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
