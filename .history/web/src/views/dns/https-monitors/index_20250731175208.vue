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
          <a-button type="primary" @click="handleAdd">
            <template #icon>
              <PlusOutlined />
            </template>
            添加监控
          </a-button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-container">
      <a-row :gutter="16">
        <a-col :span="6">
          <a-card>
            <a-statistic
              title="总监控数"
              :value="statistics.total"
              :value-style="{ color: '#1890ff' }"
            />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic
              title="在线"
              :value="statistics.online"
              :value-style="{ color: '#52c41a' }"
            />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic
              title="离线"
              :value="statistics.offline"
              :value-style="{ color: '#ff4d4f' }"
            />
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic
              title="警告"
              :value="statistics.warning"
              :value-style="{ color: '#fa8c16' }"
            />
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-container">
      <a-card :bordered="false">
        <a-form
          ref="searchFormRef"
          :model="searchForm"
          layout="inline"
          class="search-form"
        >
          <a-form-item label="名称" name="keyword">
            <a-input
              v-model:value="searchForm.keyword"
              placeholder="请输入监控名称或URL"
              allow-clear
              style="width: 200px"
            />
          </a-form-item>
          <a-form-item label="状态" name="status">
            <a-select
              v-model:value="searchForm.status"
              placeholder="请选择状态"
              allow-clear
              style="width: 120px"
            >
              <a-select-option value="active">启用</a-select-option>
              <a-select-option value="inactive">停用</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="运行状态" name="last_status">
            <a-select
              v-model:value="searchForm.last_status"
              placeholder="请选择运行状态"
              allow-clear
              style="width: 120px"
            >
              <a-select-option value="online">在线</a-select-option>
              <a-select-option value="offline">离线</a-select-option>
              <a-select-option value="warning">警告</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="handleSearch">
              <template #icon>
                <SearchOutlined />
              </template>
              搜索
            </a-button>
            <a-button @click="handleReset" style="margin-left: 8px">
              重置
            </a-button>
          </a-form-item>
        </a-form>
      </a-card>
    </div>

    <!-- 监控列表 -->
    <div class="table-container">
      <a-card :bordered="false">
        <template #title>
          <div class="table-title">
            <span>监控列表</span>
            <div class="table-actions">
              <a-button @click="handleRefresh">
                <template #icon>
                  <ReloadOutlined />
                </template>
                刷新
              </a-button>
              <a-button @click="handleBatchDelete" :disabled="!hasSelected">
                <template #icon>
                  <DeleteOutlined />
                </template>
                批量删除
              </a-button>
            </div>
          </div>
        </template>

        <a-table
          :columns="columns"
          :data-source="tableData"
          :loading="loading"
          :pagination="pagination"
          :row-selection="rowSelection"
          row-key="id"
          @change="handleTableChange"
        >
          <!-- 名称列 -->
          <template #name="{ record }">
            <div class="monitor-name">
              <a @click="handleView(record)" class="monitor-link">
                {{ record.name }}
              </a>
              <div class="monitor-url">{{ record.url }}</div>
            </div>
          </template>

          <!-- 状态列 -->
          <template #status="{ record }">
            <a-badge
              :status="getStatusBadge(record.status)"
              :text="getStatusText(record.status)"
            />
          </template>

          <!-- 运行状态列 -->
          <template #last_status="{ record }">
            <a-badge
              :status="getLastStatusBadge(record.last_status)"
              :text="getLastStatusText(record.last_status)"
            />
          </template>

          <!-- 响应时间列 -->
          <template #last_response_time="{ record }">
            <span v-if="record.last_response_time">
              {{ record.last_response_time }}ms
            </span>
            <span v-else>-</span>
          </template>

          <!-- 证书过期时间列 -->
          <template #ssl_cert_expires_at="{ record }">
            <div class="cert-expire" v-if="record.ssl_cert_expires_at">
              <span :class="getCertExpireClass(record.ssl_cert_expires_at)">
                {{ formatDate(record.ssl_cert_expires_at) }}
              </span>
              <a-tag
                v-if="isCertExpiringSoon(record.ssl_cert_expires_at)"
                color="orange"
                size="small"
              >
                即将过期
              </a-tag>
            </div>
            <span v-else>-</span>
          </template>

          <!-- 最后检查列 -->
          <template #last_checked="{ record }">
            <span>{{ formatDate(record.last_checked) }}</span>
          </template>

          <!-- 操作列 -->
          <template #action="{ record }">
            <a-space>
              <a @click="handleView(record)">查看</a>
              <a @click="handleEdit(record)">编辑</a>
              <a @click="handleCheck(record)">检查</a>
              <a-dropdown>
                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="handleToggleStatus(record)">
                      {{ record.status === 'active' ? '停用' : '启用' }}
                    </a-menu-item>
                    <a-menu-divider />
                    <a-menu-item @click="handleDelete(record)" class="danger">
                      删除
                    </a-menu-item>
                  </a-menu>
                </template>
                <a>
                  更多
                  <DownOutlined />
                </a>
              </a-dropdown>
            </a-space>
          </template>
        </a-table>
      </a-card>
    </div>

    <!-- 监控表单弹窗 -->
    <MonitorModal
      v-model:visible="modalVisible"
      :mode="modalMode"
      :monitor="currentMonitor"
      @success="handleModalSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  ReloadOutlined,
  DeleteOutlined,
  DownOutlined
} from '@ant-design/icons-vue'
import type { TableColumnsType } from 'ant-design-vue'
import { usePagination } from '@/hooks/usePagination'
import { useSelection } from '@/hooks/useSelection'
import MonitorModal from './components/MonitorModal.vue'
import { httpsMonitorApi } from '@/api/dns/httpsMonitor'
import type { HTTPSMonitor } from '@/types/dns'

// 响应式数据
const loading = ref(false)
const modalVisible = ref(false)
const modalMode = ref<'add' | 'edit' | 'view'>('add')
const currentMonitor = ref<HTTPSMonitor | null>(null)
const tableData = ref<HTTPSMonitor[]>([])
const statistics = ref({
  total: 0,
  online: 0,
  offline: 0,
  warning: 0
})

// 搜索表单
const searchFormRef = ref()
const searchForm = reactive({
  keyword: '',
  status: undefined,
  last_status: undefined
})

// 分页和选择
const { pagination, handleTableChange } = usePagination()
const { selectedRowKeys, rowSelection, hasSelected } = useSelection()

// 表格列定义
const columns: TableColumnsType = [
  {
    title: '监控名称',
    dataIndex: 'name',
    key: 'name',
    slots: { customRender: 'name' }
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    slots: { customRender: 'status' },
    width: 100
  },
  {
    title: '运行状态',
    dataIndex: 'last_status',
    key: 'last_status',
    slots: { customRender: 'last_status' },
    width: 100
  },
  {
    title: '响应时间',
    dataIndex: 'last_response_time',
    key: 'last_response_time',
    slots: { customRender: 'last_response_time' },
    width: 100
  },
  {
    title: '检查间隔',
    dataIndex: 'check_interval',
    key: 'check_interval',
    render: (text: number) => `${text}分钟`,
    width: 100
  },
  {
    title: '证书过期',
    dataIndex: 'ssl_cert_expires_at',
    key: 'ssl_cert_expires_at',
    slots: { customRender: 'ssl_cert_expires_at' },
    width: 150
  },
  {
    title: '最后检查',
    dataIndex: 'last_checked',
    key: 'last_checked',
    slots: { customRender: 'last_checked' },
    width: 150
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'action' },
    width: 150,
    fixed: 'right'
  }
]

// 工具方法
const getStatusBadge = (status: string) => {
  const statusMap = {
    active: 'success',
    inactive: 'default'
  }
  return statusMap[status] || 'default'
}

const getStatusText = (status: string) => {
  const statusMap = {
    active: '启用',
    inactive: '停用'
  }
  return statusMap[status] || status
}

const getLastStatusBadge = (status: string) => {
  const statusMap = {
    online: 'success',
    offline: 'error',
    warning: 'warning'
  }
  return statusMap[status] || 'default'
}

const getLastStatusText = (status: string) => {
  const statusMap = {
    online: '在线',
    offline: '离线',
    warning: '警告'
  }
  return statusMap[status] || status
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const getCertExpireClass = (expiresAt: string) => {
  if (!expiresAt) return ''
  const expireDate = new Date(expiresAt)
  const now = new Date()
  const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  
  if (diffDays < 0) return 'text-red'
  if (diffDays <= 30) return 'text-orange'
  return ''
}

const isCertExpiringSoon = (expiresAt: string) => {
  if (!expiresAt) return false
  const expireDate = new Date(expiresAt)
  const now = new Date()
  const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  return diffDays > 0 && diffDays <= 30
}

// 事件处理
const handleAdd = () => {
  modalMode.value = 'add'
  currentMonitor.value = null
  modalVisible.value = true
}

const handleEdit = (record: HTTPSMonitor) => {
  modalMode.value = 'edit'
  currentMonitor.value = record
  modalVisible.value = true
}

const handleView = (record: HTTPSMonitor) => {
  modalMode.value = 'view'
  currentMonitor.value = record
  modalVisible.value = true
}

const handleCheck = async (record: HTTPSMonitor) => {
  try {
    loading.value = true
    const result = await httpsMonitorApi.check(record.id)
    message.success(`检查完成，状态：${result.status}，响应时间：${result.response_time}ms`)
    await fetchData()
  } catch (error) {
    message.error('检查失败')
  } finally {
    loading.value = false
  }
}

const handleToggleStatus = async (record: HTTPSMonitor) => {
  try {
    const newStatus = record.status === 'active' ? 'inactive' : 'active'
    await httpsMonitorApi.update(record.id, { status: newStatus })
    message.success(`已${newStatus === 'active' ? '启用' : '停用'}监控`)
    await fetchData()
  } catch (error) {
    message.error('操作失败')
  }
}

const handleDelete = (record: HTTPSMonitor) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除监控 "${record.name}" 吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await httpsMonitorApi.delete(record.id)
        message.success('删除成功')
        await fetchData()
      } catch (error) {
        message.error('删除失败')
      }
    }
  })
}

const handleBatchDelete = () => {
  if (!hasSelected.value) return
  
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 个监控吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await httpsMonitorApi.batchDelete(selectedRowKeys.value)
        message.success('批量删除成功')
        selectedRowKeys.value = []
        await fetchData()
      } catch (error) {
        message.error('批量删除失败')
      }
    }
  })
}

const handleSearch = () => {
  pagination.current = 1
  fetchData()
}

const handleReset = () => {
  searchFormRef.value?.resetFields()
  pagination.current = 1
  fetchData()
}

const handleRefresh = () => {
  fetchData()
  fetchStatistics()
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
      page: pagination.current,
      size: pagination.pageSize,
      ...searchForm
    }
    const response = await httpsMonitorApi.list(params)
    tableData.value = response.items || []
    pagination.total = response.total || 0
  } catch (error) {
    message.error('获取监控列表失败')
  } finally {
    loading.value = false
  }
}

const fetchStatistics = async () => {
  try {
    const response = await httpsMonitorApi.getStatistics()
    statistics.value = response
  } catch (error) {
    console.error('获取统计信息失败:', error)
  }
}

// 表格列定义（完整版）
const columns: TableColumnsType = [
  {
    title: '监控名称',
    dataIndex: 'name',
    key: 'name',
    slots: { customRender: 'name' },
    width: 250,
    fixed: 'left'
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    slots: { customRender: 'status' },
    width: 80,
    filters: [
      { text: '启用', value: 'active' },
      { text: '停用', value: 'inactive' }
    ]
  },
  {
    title: '运行状态',
    dataIndex: 'last_status',
    key: 'last_status',
    slots: { customRender: 'last_status' },
    width: 120,
    filters: [
      { text: '在线', value: 'online' },
      { text: '离线', value: 'offline' },
      { text: '警告', value: 'warning' }
    ]
  },
  {
    title: '响应时间',
    dataIndex: 'last_response_time',
    key: 'last_response_time',
    slots: { customRender: 'last_response_time' },
    width: 100,
    sorter: true
  },
  {
    title: '检查间隔',
    dataIndex: 'check_interval',
    key: 'check_interval',
    slots: { customRender: 'check_interval' },
    width: 100,
    sorter: true
  },
  {
    title: 'SSL证书',
    dataIndex: 'ssl_info',
    key: 'ssl_info',
    slots: { customRender: 'ssl_info' },
    width: 180
  },
  {
    title: '最后检查',
    dataIndex: 'last_checked',
    key: 'last_checked',
    slots: { customRender: 'last_checked' },
    width: 150,
    sorter: true
  },
  {
    title: '通知设置',
    dataIndex: 'notification',
    key: 'notification',
    slots: { customRender: 'notification' },
    width: 120
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'action' },
    width: 200,
    fixed: 'right'
  }
]

// 扩展的工具方法
const getResponseTimeClass = (responseTime: number) => {
  if (responseTime < 500) return 'response-fast'
  if (responseTime < 2000) return 'response-normal'
  return 'response-slow'
}

const getTimeDiffColor = (lastChecked: string, interval: number) => {
  if (!lastChecked) return 'red'
  const diff = Date.now() - new Date(lastChecked).getTime()
  const expectedInterval = interval * 60 * 1000 // 转换为毫秒

  if (diff > expectedInterval * 2) return 'red'
  if (diff > expectedInterval * 1.5) return 'orange'
  return 'green'
}

const getTimeDiffText = (lastChecked: string, interval: number) => {
  if (!lastChecked) return '未检查'
  const diff = Date.now() - new Date(lastChecked).getTime()
  const minutes = Math.floor(diff / (1000 * 60))

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  return `${days}天前`
}

const getChannelColor = (channel: string) => {
  const colorMap = {
    email: 'blue',
    sms: 'green',
    webhook: 'purple',
    dingtalk: 'orange',
    wechat: 'cyan'
  }
  return colorMap[channel] || 'default'
}

const getChannelName = (channel: string) => {
  const nameMap = {
    email: '邮件',
    sms: '短信',
    webhook: 'Webhook',
    dingtalk: '钉钉',
    wechat: '微信'
  }
  return nameMap[channel] || channel
}

const truncateText = (text: string, maxLength: number) => {
  if (!text) return ''
  return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
}

const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    message.success('已复制到剪贴板')
  } catch (error) {
    message.error('复制失败')
  }
}

// 扩展的事件处理方法
const handleBatchCheck = async () => {
  if (!hasSelected.value) return

  try {
    checkingIds.value = [...selectedRowKeys.value]
    message.info(`开始批量检查 ${selectedRowKeys.value.length} 个监控...`)

    // 并发检查，但限制并发数
    const batchSize = 5
    for (let i = 0; i < selectedRowKeys.value.length; i += batchSize) {
      const batch = selectedRowKeys.value.slice(i, i + batchSize)
      await Promise.all(batch.map(id => httpsMonitorApi.check(id)))
    }

    message.success('批量检查完成')
    await fetchData()
  } catch (error) {
    message.error('批量检查失败')
  } finally {
    checkingIds.value = []
  }
}

const handleBatchEnable = async () => {
  if (!hasSelected.value) return

  try {
    await Promise.all(
      selectedRowKeys.value.map(id =>
        httpsMonitorApi.update(id, { status: 'active' })
      )
    )
    message.success('批量启用成功')
    selectedRowKeys.value = []
    await fetchData()
  } catch (error) {
    message.error('批量启用失败')
  }
}

const handleBatchDisable = async () => {
  if (!hasSelected.value) return

  try {
    await Promise.all(
      selectedRowKeys.value.map(id =>
        httpsMonitorApi.update(id, { status: 'inactive' })
      )
    )
    message.success('批量停用成功')
    selectedRowKeys.value = []
    await fetchData()
  } catch (error) {
    message.error('批量停用失败')
  }
}

const handleExport = () => {
  // TODO: 实现导出功能
  message.info('导出功能开发中')
}

const handleBatchExport = () => {
  if (!hasSelected.value) return
  // TODO: 实现批量导出功能
  message.info('批量导出功能开发中')
}

const handleToggleNotification = async (record: HTTPSMonitor, enabled: boolean) => {
  try {
    await httpsMonitorApi.update(record.id, { notification_enabled: enabled })
    message.success(`已${enabled ? '开启' : '关闭'}通知`)
    await fetchData()
  } catch (error) {
    message.error('操作失败')
  }
}

const handleViewHistory = (record: HTTPSMonitor) => {
  currentMonitor.value = record
  historyModalVisible.value = true
}

const handleClone = (record: HTTPSMonitor) => {
  modalMode.value = 'add'
  currentMonitor.value = {
    ...record,
    id: 0,
    name: `${record.name} - 副本`,
    created_at: '',
    updated_at: ''
  }
  modalVisible.value = true
}

const showExpiringCerts = () => {
  expiringModalVisible.value = true
}

// 获取即将过期的证书
const fetchExpiringCerts = async () => {
  try {
    const response = await httpsMonitorApi.getExpiringCertificates(30)
    expiringCerts.value = response || []
  } catch (error) {
    console.error('获取即将过期证书失败:', error)
  }
}

// 启动自动刷新
const startAutoRefresh = () => {
  refreshTimer = setInterval(() => {
    fetchData()
    fetchStatistics()
  }, 60000) // 每分钟刷新一次
}

// 停止自动刷新
const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

// 生命周期
onMounted(() => {
  fetchData()
  fetchStatistics()
  fetchExpiringCerts()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped lang="scss">
.https-monitor-container {
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
  
  .search-form {
    .ant-form-item {
      margin-bottom: 16px;
    }
  }
}

.table-container {
  .table-title {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    .table-actions {
      .ant-btn {
        margin-left: 8px;
      }
    }
  }
}

.monitor-name {
  .monitor-link {
    font-weight: 500;
    color: #1890ff;
    text-decoration: none;
    
    &:hover {
      color: #40a9ff;
    }
  }
  
  .monitor-url {
    font-size: 12px;
    color: #8c8c8c;
    margin-top: 2px;
  }
}

.cert-expire {
  .text-red {
    color: #ff4d4f;
  }
  
  .text-orange {
    color: #fa8c16;
  }
}

.danger {
  color: #ff4d4f !important;
}
</style>
