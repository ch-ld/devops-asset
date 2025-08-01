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
          <a-button type="primary" @click="handleAdd">
            <template #icon>
              <PlusOutlined />
            </template>
            添加记录
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
              title="总记录数"
              :value="statistics.total"
              :value-style="{ color: '#1890ff' }"
            >
              <template #suffix>
                <DatabaseOutlined />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic
              title="正常记录"
              :value="statistics.active"
              :value-style="{ color: '#52c41a' }"
            >
              <template #suffix>
                <CheckCircleOutlined />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic
              title="异常记录"
              :value="statistics.error"
              :value-style="{ color: '#ff4d4f' }"
            >
              <template #suffix>
                <ExclamationCircleOutlined />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic
              title="今日更新"
              :value="statistics.today"
              :value-style="{ color: '#fa8c16' }"
            >
              <template #suffix>
                <CalendarOutlined />
              </template>
            </a-statistic>
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
          <a-form-item label="域名" name="domain_id">
            <a-select
              v-model:value="searchForm.domain_id"
              placeholder="请选择域名"
              allow-clear
              style="width: 200px"
              show-search
              :filter-option="filterOption"
            >
              <a-select-option
                v-for="domain in domainOptions"
                :key="domain.id"
                :value="domain.id"
              >
                {{ domain.name }}
              </a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="记录名" name="keyword">
            <a-input
              v-model:value="searchForm.keyword"
              placeholder="请输入记录名关键词"
              allow-clear
              style="width: 200px"
            />
          </a-form-item>
          <a-form-item label="记录类型" name="type">
            <a-select
              v-model:value="searchForm.type"
              placeholder="请选择记录类型"
              allow-clear
              style="width: 120px"
            >
              <a-select-option value="A">A</a-select-option>
              <a-select-option value="AAAA">AAAA</a-select-option>
              <a-select-option value="CNAME">CNAME</a-select-option>
              <a-select-option value="MX">MX</a-select-option>
              <a-select-option value="TXT">TXT</a-select-option>
              <a-select-option value="NS">NS</a-select-option>
              <a-select-option value="SRV">SRV</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="状态" name="status">
            <a-select
              v-model:value="searchForm.status"
              placeholder="请选择状态"
              allow-clear
              style="width: 120px"
            >
              <a-select-option value="active">正常</a-select-option>
              <a-select-option value="inactive">停用</a-select-option>
              <a-select-option value="error">异常</a-select-option>
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

    <!-- 记录列表 -->
    <div class="table-container">
      <a-card :bordered="false">
        <template #title>
          <div class="table-title">
            <span>解析记录列表</span>
            <div class="table-actions">
              <a-button @click="handleRefresh">
                <template #icon>
                  <ReloadOutlined />
                </template>
                刷新
              </a-button>
              <a-button @click="handleBatchSync" :disabled="!hasSelected">
                <template #icon>
                  <SyncOutlined />
                </template>
                批量同步
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
          <!-- 记录信息列 -->
          <template #record_info="{ record }">
            <div class="record-info">
              <div class="record-name">
                <span class="name">{{ record.name }}</span>
                <a-tag :color="getTypeColor(record.type)" size="small">
                  {{ record.type }}
                </a-tag>
              </div>
              <div class="record-domain">{{ record.domain?.name }}</div>
            </div>
          </template>

          <!-- 记录值列 -->
          <template #value="{ record }">
            <div class="record-value">
              <a-tooltip :title="record.value">
                <span class="value-text">{{ truncateText(record.value, 50) }}</span>
              </a-tooltip>
              <a-button
                type="link"
                size="small"
                @click="copyToClipboard(record.value)"
                style="padding: 0; margin-left: 4px"
              >
                <CopyOutlined />
              </a-button>
            </div>
          </template>

          <!-- 状态列 -->
          <template #status="{ record }">
            <a-badge
              :status="getStatusBadge(record.status)"
              :text="getStatusText(record.status)"
            />
          </template>

          <!-- TTL列 -->
          <template #ttl="{ record }">
            <span>{{ record.ttl }}s</span>
          </template>

          <!-- 优先级列 -->
          <template #priority="{ record }">
            <span v-if="record.priority !== undefined">{{ record.priority }}</span>
            <span v-else>-</span>
          </template>

          <!-- 最后同步列 -->
          <template #last_sync_at="{ record }">
            <div class="sync-info">
              <span>{{ formatDate(record.last_sync_at) }}</span>
              <div v-if="record.sync_status" class="sync-status">
                <a-tag
                  :color="record.sync_status === 'success' ? 'green' : 'red'"
                  size="small"
                >
                  {{ record.sync_status === 'success' ? '成功' : '失败' }}
                </a-tag>
              </div>
            </div>
          </template>

          <!-- 操作列 -->
          <template #action="{ record }">
            <a-space>
              <a @click="handleView(record)">查看</a>
              <a @click="handleEdit(record)">编辑</a>
              <a @click="handleSync(record)" :loading="syncingIds.includes(record.id)">
                同步
              </a>
              <a-dropdown>
                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="handleToggleStatus(record)">
                      {{ record.status === 'active' ? '停用' : '启用' }}
                    </a-menu-item>
                    <a-menu-item @click="handleClone(record)">
                      克隆记录
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

    <!-- 记录表单弹窗 -->
    <RecordModal
      v-model:visible="modalVisible"
      :mode="modalMode"
      :record="currentRecord"
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
  SyncOutlined,
  DownOutlined,
  DatabaseOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  CalendarOutlined,
  CopyOutlined
} from '@ant-design/icons-vue'
import type { TableColumnsType } from 'ant-design-vue'
import { usePagination } from '@/hooks/usePagination'
import { useSelection } from '@/hooks/useSelection'
import RecordModal from './components/RecordModal.vue'
import { dnsRecordApi } from '@/api/dns/record'
import { domainApi } from '@/api/dns/domain'
import type { DNSRecord, Domain } from '@/types/dns'

// 响应式数据
const loading = ref(false)
const modalVisible = ref(false)
const modalMode = ref<'add' | 'edit' | 'view'>('add')
const currentRecord = ref<DNSRecord | null>(null)
const tableData = ref<DNSRecord[]>([])
const domainOptions = ref<Domain[]>([])
const syncingIds = ref<number[]>([])
const statistics = ref({
  total: 0,
  active: 0,
  error: 0,
  today: 0
})

// 搜索表单
const searchFormRef = ref()
const searchForm = reactive({
  domain_id: undefined,
  keyword: '',
  type: undefined,
  status: undefined
})

// 分页和选择
const { pagination, handleTableChange } = usePagination()
const { selectedRowKeys, rowSelection, hasSelected } = useSelection()

// 表格列定义
const columns: TableColumnsType = [
  {
    title: '记录信息',
    dataIndex: 'record_info',
    key: 'record_info',
    slots: { customRender: 'record_info' },
    width: 250
  },
  {
    title: '记录值',
    dataIndex: 'value',
    key: 'value',
    slots: { customRender: 'value' }
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    slots: { customRender: 'status' },
    width: 100
  },
  {
    title: 'TTL',
    dataIndex: 'ttl',
    key: 'ttl',
    slots: { customRender: 'ttl' },
    width: 80
  },
  {
    title: '优先级',
    dataIndex: 'priority',
    key: 'priority',
    slots: { customRender: 'priority' },
    width: 80
  },
  {
    title: '最后同步',
    dataIndex: 'last_sync_at',
    key: 'last_sync_at',
    slots: { customRender: 'last_sync_at' },
    width: 150
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'action' },
    width: 180,
    fixed: 'right'
  }
]

// 工具方法
const getTypeColor = (type: string) => {
  const colorMap = {
    A: 'blue',
    AAAA: 'purple',
    CNAME: 'green',
    MX: 'orange',
    TXT: 'cyan',
    NS: 'red',
    SRV: 'magenta'
  }
  return colorMap[type] || 'default'
}

const getStatusBadge = (status: string) => {
  const statusMap = {
    active: 'success',
    inactive: 'default',
    error: 'error'
  }
  return statusMap[status] || 'default'
}

const getStatusText = (status: string) => {
  const statusMap = {
    active: '正常',
    inactive: '停用',
    error: '异常'
  }
  return statusMap[status] || status
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const truncateText = (text: string, maxLength: number) => {
  if (!text) return ''
  return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
}

const filterOption = (input: string, option: any) => {
  return option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
}

const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    message.success('已复制到剪贴板')
  } catch (error) {
    message.error('复制失败')
  }
}

// 事件处理
const handleAdd = () => {
  modalMode.value = 'add'
  currentRecord.value = null
  modalVisible.value = true
}

const handleEdit = (record: DNSRecord) => {
  modalMode.value = 'edit'
  currentRecord.value = record
  modalVisible.value = true
}

const handleView = (record: DNSRecord) => {
  modalMode.value = 'view'
  currentRecord.value = record
  modalVisible.value = true
}

const handleSync = async (record: DNSRecord) => {
  try {
    syncingIds.value.push(record.id)
    await dnsRecordApi.sync(record.id)
    message.success('记录同步成功')
    await fetchData()
  } catch (error) {
    message.error('记录同步失败')
  } finally {
    syncingIds.value = syncingIds.value.filter(id => id !== record.id)
  }
}

const handleBatchSync = async () => {
  if (!hasSelected.value) return
  
  try {
    await Promise.all(
      selectedRowKeys.value.map(id => dnsRecordApi.sync(id))
    )
    message.success('批量同步成功')
    selectedRowKeys.value = []
    await fetchData()
  } catch (error) {
    message.error('批量同步失败')
  }
}

const handleToggleStatus = async (record: DNSRecord) => {
  try {
    const newStatus = record.status === 'active' ? 'inactive' : 'active'
    await dnsRecordApi.update(record.id, { status: newStatus })
    message.success(`记录已${newStatus === 'active' ? '启用' : '停用'}`)
    await fetchData()
  } catch (error) {
    message.error('操作失败')
  }
}

const handleClone = (record: DNSRecord) => {
  modalMode.value = 'add'
  currentRecord.value = {
    ...record,
    id: 0,
    name: `${record.name}-copy`,
    created_at: '',
    updated_at: ''
  }
  modalVisible.value = true
}

const handleDelete = (record: DNSRecord) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除记录 "${record.name}" 吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await dnsRecordApi.delete(record.id)
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
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 个记录吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await dnsRecordApi.batchDelete(selectedRowKeys.value)
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
    const response = await dnsRecordApi.list(params)
    tableData.value = response.items || []
    pagination.total = response.total || 0
  } catch (error) {
    message.error('获取解析记录列表失败')
  } finally {
    loading.value = false
  }
}

const fetchDomainOptions = async () => {
  try {
    const response = await domainApi.list({ page: 1, size: 1000 })
    domainOptions.value = response.items || []
  } catch (error) {
    console.error('获取域名列表失败:', error)
  }
}

const fetchStatistics = async () => {
  try {
    // TODO: 实现统计数据获取
    // 暂时使用模拟数据
    const total = tableData.value.length
    const active = tableData.value.filter(item => item.status === 'active').length
    const error = tableData.value.filter(item => item.status === 'error').length
    const today = tableData.value.filter(item => {
      const today = new Date().toDateString()
      return new Date(item.updated_at).toDateString() === today
    }).length
    
    statistics.value = { total, active, error, today }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 生命周期
onMounted(() => {
  fetchData()
  fetchDomainOptions()
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

.record-info {
  .record-name {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 4px;
    
    .name {
      font-weight: 500;
      color: #262626;
    }
  }
  
  .record-domain {
    font-size: 12px;
    color: #8c8c8c;
  }
}

.record-value {
  display: flex;
  align-items: center;
  
  .value-text {
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 13px;
    color: #262626;
  }
}

.sync-info {
  .sync-status {
    margin-top: 4px;
  }
}

.danger {
  color: #ff4d4f !important;
}
</style>
