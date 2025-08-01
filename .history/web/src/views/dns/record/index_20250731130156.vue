<template>
  <div class="dns-record-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title">
          <h1>DNS记录管理</h1>
          <p>管理域名的DNS解析记录，支持A、AAAA、CNAME、TXT、MX、SRV等记录类型</p>
        </div>
        <div class="page-actions">
          <a-button @click="handleSync" :loading="syncLoading">
            <template #icon>
              <SyncOutlined />
            </template>
            同步记录
          </a-button>
          <a-button type="primary" @click="handleAdd">
            <template #icon>
              <PlusOutlined />
            </template>
            添加记录
          </a-button>
        </div>
      </div>
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
              show-search
              :filter-option="filterDomainOption"
              style="width: 200px"
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
          <a-form-item label="记录名称" name="keyword">
            <a-input
              v-model:value="searchForm.keyword"
              placeholder="请输入记录名称"
              allow-clear
              style="width: 150px"
            />
          </a-form-item>
          <a-form-item label="记录类型" name="type">
            <a-select
              v-model:value="searchForm.type"
              placeholder="请选择类型"
              allow-clear
              style="width: 120px"
            >
              <a-select-option value="A">A</a-select-option>
              <a-select-option value="AAAA">AAAA</a-select-option>
              <a-select-option value="CNAME">CNAME</a-select-option>
              <a-select-option value="TXT">TXT</a-select-option>
              <a-select-option value="MX">MX</a-select-option>
              <a-select-option value="SRV">SRV</a-select-option>
              <a-select-option value="NS">NS</a-select-option>
              <a-select-option value="PTR">PTR</a-select-option>
            </a-select>
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

    <!-- DNS记录列表 -->
    <div class="table-container">
      <a-card :bordered="false">
        <template #title>
          <div class="table-title">
            <span>DNS记录列表</span>
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
          <!-- 记录名称列 -->
          <template #name="{ record }">
            <div class="record-name">
              <span class="name-text">{{ record.name }}</span>
              <a-tag
                v-if="record.name === '@'"
                color="blue"
                size="small"
              >
                根域名
              </a-tag>
            </div>
          </template>

          <!-- 记录类型列 -->
          <template #type="{ record }">
            <a-tag :color="getTypeColor(record.type)">
              {{ record.type }}
            </a-tag>
          </template>

          <!-- 记录值列 -->
          <template #value="{ record }">
            <div class="record-value">
              <a-tooltip :title="record.value">
                <span class="value-text">{{ truncateValue(record.value) }}</span>
              </a-tooltip>
              <a-button
                type="text"
                size="small"
                @click="copyToClipboard(record.value)"
              >
                <template #icon>
                  <CopyOutlined />
                </template>
              </a-button>
            </div>
          </template>

          <!-- TTL列 -->
          <template #ttl="{ record }">
            <span>{{ formatTTL(record.ttl) }}</span>
          </template>

          <!-- 优先级列 -->
          <template #priority="{ record }">
            <span v-if="record.priority !== null && record.priority !== undefined">
              {{ record.priority }}
            </span>
            <span v-else class="text-gray">-</span>
          </template>

          <!-- 状态列 -->
          <template #status="{ record }">
            <a-badge
              :status="record.status === 'active' ? 'success' : 'default'"
              :text="record.status === 'active' ? '启用' : '停用'"
            />
          </template>

          <!-- 同步状态列 -->
          <template #sync_status="{ record }">
            <a-tag
              :color="getSyncStatusColor(record.sync_status)"
              size="small"
            >
              {{ getSyncStatusText(record.sync_status) }}
            </a-tag>
          </template>

          <!-- 域名列 -->
          <template #domain="{ record }">
            <span v-if="record.domain">{{ record.domain.name }}</span>
            <span v-else class="text-gray">-</span>
          </template>

          <!-- 操作列 -->
          <template #action="{ record }">
            <a-space>
              <a @click="handleView(record)">查看</a>
              <a @click="handleEdit(record)">编辑</a>
              <a-dropdown>
                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="handleTest(record)">
                      测试解析
                    </a-menu-item>
                    <a-menu-item @click="handleSyncSingle(record)">
                      同步记录
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

    <!-- DNS记录表单弹窗 -->
    <RecordModal
      v-model:visible="modalVisible"
      :mode="modalMode"
      :record="currentRecord"
      :domain-options="domainOptions"
      :provider-options="providerOptions"
      @success="handleModalSuccess"
    />

    <!-- 同步结果弹窗 -->
    <SyncResultModal
      v-model:visible="syncResultVisible"
      :result="syncResult"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  ReloadOutlined,
  DeleteOutlined,
  DownOutlined,
  SyncOutlined,
  CopyOutlined
} from '@ant-design/icons-vue'
import type { TableColumnsType } from 'ant-design-vue'
import { usePagination } from '@/hooks/usePagination'
import { useSelection } from '@/hooks/useSelection'
import RecordModal from './components/RecordModal.vue'
import SyncResultModal from './components/SyncResultModal.vue'
import { recordApi } from '@/api/dns/record'
import { domainApi } from '@/api/dns/domain'
import { providerApi } from '@/api/dns/provider'
import type { DNSRecord, Domain, DNSProvider, SyncResponse } from '@/types/dns'

// 响应式数据
const loading = ref(false)
const syncLoading = ref(false)
const modalVisible = ref(false)
const syncResultVisible = ref(false)
const modalMode = ref<'add' | 'edit' | 'view'>('add')
const currentRecord = ref<DNSRecord | null>(null)
const tableData = ref<DNSRecord[]>([])
const domainOptions = ref<Domain[]>([])
const providerOptions = ref<DNSProvider[]>([])
const syncResult = ref<SyncResponse | null>(null)

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
    title: '记录名称',
    dataIndex: 'name',
    key: 'name',
    slots: { customRender: 'name' },
    width: 150
  },
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type',
    slots: { customRender: 'type' },
    width: 80
  },
  {
    title: '记录值',
    dataIndex: 'value',
    key: 'value',
    slots: { customRender: 'value' },
    ellipsis: true
  },
  {
    title: 'TTL',
    dataIndex: 'ttl',
    key: 'ttl',
    slots: { customRender: 'ttl' },
    width: 100
  },
  {
    title: '优先级',
    dataIndex: 'priority',
    key: 'priority',
    slots: { customRender: 'priority' },
    width: 80
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    slots: { customRender: 'status' },
    width: 80
  },
  {
    title: '同步状态',
    dataIndex: 'sync_status',
    key: 'sync_status',
    slots: { customRender: 'sync_status' },
    width: 100
  },
  {
    title: '域名',
    dataIndex: 'domain',
    key: 'domain',
    slots: { customRender: 'domain' },
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
const getTypeColor = (type: string) => {
  const colorMap = {
    A: 'blue',
    AAAA: 'cyan',
    CNAME: 'green',
    TXT: 'orange',
    MX: 'purple',
    SRV: 'magenta',
    NS: 'red',
    PTR: 'lime'
  }
  return colorMap[type] || 'default'
}

const getSyncStatusColor = (status: string) => {
  const colorMap = {
    synced: 'success',
    pending: 'processing',
    failed: 'error'
  }
  return colorMap[status] || 'default'
}

const getSyncStatusText = (status: string) => {
  const textMap = {
    synced: '已同步',
    pending: '同步中',
    failed: '同步失败'
  }
  return textMap[status] || status
}

const formatTTL = (ttl: number) => {
  if (ttl >= 3600) {
    return `${Math.floor(ttl / 3600)}h`
  } else if (ttl >= 60) {
    return `${Math.floor(ttl / 60)}m`
  }
  return `${ttl}s`
}

const truncateValue = (value: string, maxLength = 50) => {
  if (value.length <= maxLength) return value
  return value.substring(0, maxLength) + '...'
}

const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    message.success('已复制到剪贴板')
  } catch (error) {
    message.error('复制失败')
  }
}

const filterDomainOption = (input: string, option: any) => {
  return option.children.toLowerCase().includes(input.toLowerCase())
}

// 事件处理方法
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

const handleDelete = (record: DNSRecord) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除DNS记录 "${record.name}" 吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await recordApi.delete(record.id)
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
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 条DNS记录吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await recordApi.batchDelete(selectedRowKeys.value)
        message.success('批量删除成功')
        selectedRowKeys.value = []
        await fetchData()
      } catch (error) {
        message.error('批量删除失败')
      }
    }
  })
}

const handleTest = (record: DNSRecord) => {
  // 测试DNS解析
  message.info('DNS解析测试功能开发中...')
}

const handleSync = async () => {
  try {
    syncLoading.value = true
    const result = await recordApi.sync({
      sync_type: 'full'
    })
    syncResult.value = result
    syncResultVisible.value = true
    await fetchData()
  } catch (error) {
    message.error('同步失败')
  } finally {
    syncLoading.value = false
  }
}

const handleSyncSingle = async (record: DNSRecord) => {
  try {
    await recordApi.sync({
      domain_id: record.domain_id,
      sync_type: 'incremental'
    })
    message.success('记录同步成功')
    await fetchData()
  } catch (error) {
    message.error('记录同步失败')
  }
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
}

const handleModalSuccess = () => {
  modalVisible.value = false
  fetchData()
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
    const response = await recordApi.list(params)
    tableData.value = response.items || []
    pagination.total = response.total || 0
  } catch (error) {
    message.error('获取DNS记录列表失败')
  } finally {
    loading.value = false
  }
}

const fetchDomainOptions = async () => {
  try {
    const response = await domainApi.list({ size: 1000 })
    domainOptions.value = response.items || []
  } catch (error) {
    console.error('获取域名选项失败:', error)
  }
}

const fetchProviderOptions = async () => {
  try {
    const response = await providerApi.list({ size: 1000 })
    providerOptions.value = response.items || []
  } catch (error) {
    console.error('获取提供商选项失败:', error)
  }
}

// 生命周期
onMounted(() => {
  fetchData()
  fetchDomainOptions()
  fetchProviderOptions()
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
    
    .page-actions {
      .ant-btn {
        margin-left: 8px;
      }
    }
  }
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

.record-name {
  display: flex;
  align-items: center;
  gap: 8px;
  
  .name-text {
    font-weight: 500;
  }
}

.record-value {
  display: flex;
  align-items: center;
  gap: 4px;
  
  .value-text {
    flex: 1;
    min-width: 0;
  }
}

.text-gray {
  color: #8c8c8c;
}

.danger {
  color: #ff4d4f !important;
}
</style>
