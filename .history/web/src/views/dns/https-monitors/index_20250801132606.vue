<template>
  <div class="https-monitor-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title">
          <h1>HTTPS监控</h1>
          <p>监控网站HTTPS服务状态和SSL证书有效性</p>
        </div>
        <div class="page-actions">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加监控
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
              title="总监控数"
              :value="statistics.total"
              :value-style="{ color: '#409eff' }"
            />
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="在线"
              :value="statistics.online"
              :value-style="{ color: '#67c23a' }"
            />
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="离线"
              :value="statistics.offline"
              :value-style="{ color: '#f56c6c' }"
            />
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="证书过期"
              :value="statistics.expiring"
              :value-style="{ color: '#e6a23c' }"
            />
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 搜索和筛选 -->
    <el-card class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="监控名称">
          <el-input
            v-model="searchForm.name"
            placeholder="请输入监控名称"
            clearable
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="监控状态">
          <el-select 
            v-model="searchForm.status" 
            placeholder="请选择状态"
            clearable
            style="width: 120px"
          >
            <el-option label="全部" value="" />
            <el-option label="在线" value="up" />
            <el-option label="离线" value="down" />
            <el-option label="未知" value="unknown" />
          </el-select>
        </el-form-item>
        <el-form-item label="证书状态">
          <el-select 
            v-model="searchForm.certStatus" 
            placeholder="请选择证书状态"
            clearable
            style="width: 140px"
          >
            <el-option label="全部" value="" />
            <el-option label="正常" value="valid" />
            <el-option label="即将过期" value="expiring" />
            <el-option label="已过期" value="expired" />
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

    <!-- 监控列表 -->
    <el-card class="table-card">
      <template #header>
        <div class="table-header">
          <span>监控列表</span>
          <div class="table-actions">
            <el-button @click="handleRefreshAll">
              <el-icon><Refresh /></el-icon>
              刷新全部
            </el-button>
            <el-button @click="handleBatchCheck" :disabled="selectedRowKeys.length === 0">
              <el-icon><View /></el-icon>
              批量检查
            </el-button>
            <el-button 
              type="danger" 
              @click="handleBatchDelete" 
              :disabled="selectedRowKeys.length === 0"
            >
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="monitorList"
        @selection-change="handleSelectionChange"
        stripe
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="name" label="监控名称" width="200">
          <template #default="{ row }">
            <div class="monitor-name">
              <el-link @click="handleView(row)" type="primary">{{ row.name }}</el-link>
              <div class="monitor-url">{{ row.url }}</div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" :icon="getStatusIcon(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="响应时间" width="100">
          <template #default="{ row }">
            <span v-if="row.responseTime">{{ row.responseTime }}ms</span>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column label="HTTP状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.httpStatus" :type="getHttpStatusType(row.httpStatus)">
              {{ row.httpStatus }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column label="证书状态" width="120">
          <template #default="{ row }">
            <div v-if="row.sslInfo">
              <el-tag :type="getCertStatusType(row.sslInfo.status)">
                {{ getCertStatusText(row.sslInfo.status) }}
              </el-tag>
              <div class="cert-expiry" v-if="row.sslInfo.expiryDays !== null">
                剩余{{ row.sslInfo.expiryDays }}天
              </div>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column label="检查频率" width="100">
          <template #default="{ row }">
            {{ getIntervalText(row.checkInterval) }}
          </template>
        </el-table-column>

        <el-table-column label="最后检查" width="160">
          <template #default="{ row }">
            <div v-if="row.lastCheckedAt">
              {{ formatDateTime(row.lastCheckedAt) }}
            </div>
            <span v-else>从未检查</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button text @click="handleCheck(row)">
              <el-icon><View /></el-icon>
              检查
            </el-button>
            <el-button text @click="handleView(row)">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button text @click="handleEdit(row)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-dropdown>
              <el-button text>
                更多<el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleViewHistory(row)">查看历史</el-dropdown-item>
                  <el-dropdown-item @click="handleToggleStatus(row)">
                    {{ row.enabled ? '禁用' : '启用' }}
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="handleDelete(row)">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handlePageSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 监控弹窗 -->
    <MonitorModal
      v-model:visible="monitorModalVisible"
      :edit-data="currentMonitor"
      @success="handleMonitorSuccess"
    />

    <!-- 历史记录弹窗 -->
    <MonitorHistoryModal
      v-model:visible="historyModalVisible"
      :monitor-id="currentMonitorId"
    />

    <!-- 即将过期证书弹窗 -->
    <ExpiringCertsModal
      v-model:visible="expiringCertsModalVisible"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, 
  Search, 
  Refresh, 
  View, 
  Edit, 
  Delete, 
  ArrowDown 
} from '@element-plus/icons-vue'
import MonitorModal from './components/MonitorModal.vue'
import MonitorHistoryModal from './components/MonitorHistoryModal.vue'
import ExpiringCertsModal from './components/ExpiringCertsModal.vue'

// 响应式数据
const loading = ref(false)
const monitorList = ref<any[]>([])
const selectedRowKeys = ref<number[]>([])
const monitorModalVisible = ref(false)
const historyModalVisible = ref(false)
const expiringCertsModalVisible = ref(false)
const currentMonitor = ref<any>(null)
const currentMonitorId = ref<number | null>(null)

// 搜索表单
const searchForm = reactive({
  name: '',
  status: '',
  certStatus: ''
})

// 分页
const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

// 统计数据
const statistics = reactive({
  total: 0,
  online: 0,
  offline: 0,
  expiring: 0
})

// 获取监控列表
const fetchMonitorList = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取监控列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 模拟数据
    monitorList.value = [
      {
        id: 1,
        name: '主站监控',
        url: 'https://example.com',
        status: 'up',
        responseTime: 245,
        httpStatus: 200,
        checkInterval: 300,
        enabled: true,
        lastCheckedAt: new Date().toISOString(),
        sslInfo: {
          status: 'valid',
          expiryDays: 45,
          issuer: "Let's Encrypt"
        }
      },
      {
        id: 2,
        name: 'API监控',
        url: 'https://api.example.com',
        status: 'down',
        responseTime: null,
        httpStatus: null,
        checkInterval: 60,
        enabled: true,
        lastCheckedAt: new Date(Date.now() - 3600000).toISOString(),
        sslInfo: {
          status: 'expiring',
          expiryDays: 7,
          issuer: "Let's Encrypt"
        }
      }
    ]
    
    pagination.total = monitorList.value.length
    updateStatistics()
  } catch (error) {
    ElMessage.error('获取监控列表失败')
  } finally {
    loading.value = false
  }
}

// 更新统计数据
const updateStatistics = () => {
  statistics.total = monitorList.value.length
  statistics.online = monitorList.value.filter(item => item.status === 'up').length
  statistics.offline = monitorList.value.filter(item => item.status === 'down').length
  statistics.expiring = monitorList.value.filter(item => 
    item.sslInfo && item.sslInfo.status === 'expiring'
  ).length
}

// 处理选择变化
const handleSelectionChange = (selection: any[]) => {
  selectedRowKeys.value = selection.map(item => item.id)
}

// 处理搜索
const handleSearch = () => {
  pagination.current = 1
  fetchMonitorList()
}

// 处理重置
const handleReset = () => {
  Object.assign(searchForm, {
    name: '',
    status: '',
    certStatus: ''
  })
  handleSearch()
}

// 处理添加
const handleAdd = () => {
  currentMonitor.value = null
  monitorModalVisible.value = true
}

// 处理编辑
const handleEdit = (record: any) => {
  currentMonitor.value = { ...record }
  monitorModalVisible.value = true
}

// 处理查看
const handleView = (record: any) => {
  currentMonitorId.value = record.id
  historyModalVisible.value = true
}

// 处理检查
const handleCheck = async (record: any) => {
  try {
    ElMessage.info('正在检查监控状态...')
    // TODO: 调用API检查单个监控
    await new Promise(resolve => setTimeout(resolve, 2000))
    ElMessage.success('检查完成')
    fetchMonitorList()
  } catch (error) {
    ElMessage.error('检查失败')
  }
}

// 处理删除
const handleDelete = async (record: any) => {
  try {
    await ElMessageBox.confirm('确定要删除此监控吗？', '确认删除', {
      type: 'warning'
    })
    // TODO: 调用API删除监控
    ElMessage.success('删除成功')
    fetchMonitorList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 处理批量操作
const handleRefreshAll = async () => {
  try {
    ElMessage.info('正在刷新所有监控...')
    // TODO: 调用API刷新所有监控
    await new Promise(resolve => setTimeout(resolve, 3000))
    ElMessage.success('刷新完成')
    fetchMonitorList()
  } catch (error) {
    ElMessage.error('刷新失败')
  }
}

const handleBatchCheck = async () => {
  try {
    ElMessage.info(`正在检查${selectedRowKeys.value.length}个监控...`)
    // TODO: 调用API批量检查
    await new Promise(resolve => setTimeout(resolve, 2000))
    ElMessage.success('批量检查完成')
    fetchMonitorList()
  } catch (error) {
    ElMessage.error('批量检查失败')
  }
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(`确定要删除选中的${selectedRowKeys.value.length}个监控吗？`, '确认删除', {
      type: 'warning'
    })
    // TODO: 调用API批量删除
    ElMessage.success('批量删除成功')
    selectedRowKeys.value = []
    fetchMonitorList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 处理查看历史
const handleViewHistory = (record: any) => {
  currentMonitorId.value = record.id
  historyModalVisible.value = true
}

// 处理启用/禁用
const handleToggleStatus = async (record: any) => {
  try {
    // TODO: 调用API切换状态
    ElMessage.success(`${record.enabled ? '禁用' : '启用'}成功`)
    fetchMonitorList()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 处理分页
const handlePageChange = (page: number) => {
  pagination.current = page
  fetchMonitorList()
}

const handlePageSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.current = 1
  fetchMonitorList()
}

// 处理监控成功
const handleMonitorSuccess = () => {
  fetchMonitorList()
}

// 工具函数
const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    'up': 'success',
    'down': 'danger',
    'unknown': 'info'
  }
  return statusMap[status] || 'info'
}

const getStatusIcon = (status: string) => {
  // Element Plus中的图标组件需要单独处理
  return undefined
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    'up': '在线',
    'down': '离线',
    'unknown': '未知'
  }
  return statusMap[status] || '未知'
}

const getHttpStatusType = (status: number) => {
  if (status >= 200 && status < 300) return 'success'
  if (status >= 300 && status < 400) return 'warning'
  if (status >= 400) return 'danger'
  return 'info'
}

const getCertStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    'valid': 'success',
    'expiring': 'warning',
    'expired': 'danger'
  }
  return statusMap[status] || 'info'
}

const getCertStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    'valid': '正常',
    'expiring': '即将过期',
    'expired': '已过期'
  }
  return statusMap[status] || '未知'
}

const getIntervalText = (interval: number) => {
  if (interval < 60) return `${interval}秒`
  if (interval < 3600) return `${Math.floor(interval / 60)}分钟`
  return `${Math.floor(interval / 3600)}小时`
}

const formatDateTime = (dateTime: string) => {
  return new Date(dateTime).toLocaleString()
}

onMounted(() => {
  fetchMonitorList()
})
</script>

<style scoped>
.https-monitor-container {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.page-title h1 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
}

.page-title p {
  margin: 0;
  color: var(--el-text-color-regular);
}

.stats-container {
  margin-bottom: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.search-form {
  margin: 0;
}

.table-card {
  margin-bottom: 20px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.table-actions {
  display: flex;
  gap: 8px;
}

.monitor-name {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.monitor-url {
  font-size: 12px;
  color: var(--el-text-color-regular);
}

.cert-expiry {
  font-size: 12px;
  color: var(--el-text-color-regular);
  margin-top: 2px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
