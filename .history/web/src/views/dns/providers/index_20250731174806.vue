<template>
  <div class="dns-provider-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title">
          <h1>DNS提供商</h1>
          <p>管理DNS服务提供商，包括阿里云、腾讯云、Cloudflare等</p>
        </div>
        <div class="page-actions">
          <a-button type="primary" @click="handleAdd">
            <template #icon>
              <PlusOutlined />
            </template>
            添加提供商
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
          <a-form-item label="名称" name="keyword">
            <a-input
              v-model:value="searchForm.keyword"
              placeholder="请输入提供商名称"
              allow-clear
              style="width: 200px"
            />
          </a-form-item>
          <a-form-item label="类型" name="type">
            <a-select
              v-model:value="searchForm.type"
              placeholder="请选择类型"
              allow-clear
              style="width: 150px"
            >
              <a-select-option value="aliyun">阿里云</a-select-option>
              <a-select-option value="tencent">腾讯云</a-select-option>
              <a-select-option value="cloudflare">Cloudflare</a-select-option>
              <a-select-option value="dnspod">DNSPod</a-select-option>
              <a-select-option value="godaddy">GoDaddy</a-select-option>
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

    <!-- 提供商列表 -->
    <div class="table-container">
      <a-card :bordered="false">
        <template #title>
          <div class="table-title">
            <span>提供商列表</span>
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
            <div class="provider-name">
              <a @click="handleView(record)" class="provider-link">
                {{ record.name }}
              </a>
              <a-tag v-if="record.is_default" color="blue" size="small">默认</a-tag>
            </div>
          </template>

          <!-- 类型列 -->
          <template #type="{ record }">
            <div class="provider-type">
              <a-avatar :size="24" :src="getProviderIcon(record.type)" />
              <span style="margin-left: 8px">{{ getProviderName(record.type) }}</span>
            </div>
          </template>

          <!-- 状态列 -->
          <template #status="{ record }">
            <a-badge
              :status="getStatusBadge(record.status)"
              :text="getStatusText(record.status)"
            />
          </template>

          <!-- 优先级列 -->
          <template #priority="{ record }">
            <a-tag :color="getPriorityColor(record.priority)">
              {{ record.priority }}
            </a-tag>
          </template>

          <!-- 最后测试列 -->
          <template #last_test_at="{ record }">
            <div class="test-result">
              <span>{{ formatDate(record.last_test_at) }}</span>
              <div v-if="record.test_result" class="test-status">
                <a-tag
                  :color="record.test_result === 'success' ? 'green' : 'red'"
                  size="small"
                >
                  {{ record.test_result === 'success' ? '成功' : '失败' }}
                </a-tag>
              </div>
            </div>
          </template>

          <!-- 操作列 -->
          <template #action="{ record }">
            <a-space>
              <a @click="handleView(record)">查看</a>
              <a @click="handleEdit(record)">编辑</a>
              <a @click="handleTest(record)">测试</a>
              <a-dropdown>
                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="handleSetDefault(record)" :disabled="record.is_default">
                      设为默认
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

    <!-- 提供商表单弹窗 -->
    <ProviderModal
      v-model:visible="modalVisible"
      :mode="modalMode"
      :provider="currentProvider"
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
import ProviderModal from './components/ProviderModal.vue'
import { dnsProviderApi } from '@/api/dns/provider'
import type { DNSProvider } from '@/types/dns'

// 响应式数据
const loading = ref(false)
const modalVisible = ref(false)
const modalMode = ref<'add' | 'edit' | 'view'>('add')
const currentProvider = ref<DNSProvider | null>(null)
const tableData = ref<DNSProvider[]>([])

// 搜索表单
const searchFormRef = ref()
const searchForm = reactive({
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
    title: '名称',
    dataIndex: 'name',
    key: 'name',
    slots: { customRender: 'name' }
  },
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type',
    slots: { customRender: 'type' },
    width: 150
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    slots: { customRender: 'status' },
    width: 100
  },
  {
    title: '优先级',
    dataIndex: 'priority',
    key: 'priority',
    slots: { customRender: 'priority' },
    width: 100
  },
  {
    title: '限流',
    dataIndex: 'rate_limit',
    key: 'rate_limit',
    render: (text: number) => `${text}/min`,
    width: 100
  },
  {
    title: '最后测试',
    dataIndex: 'last_test_at',
    key: 'last_test_at',
    slots: { customRender: 'last_test_at' },
    width: 150
  },
  {
    title: '备注',
    dataIndex: 'remark',
    key: 'remark',
    ellipsis: true
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
const getProviderIcon = (type: string) => {
  const iconMap = {
    aliyun: '/icons/aliyun.png',
    tencent: '/icons/tencent.png',
    cloudflare: '/icons/cloudflare.png',
    dnspod: '/icons/dnspod.png',
    godaddy: '/icons/godaddy.png'
  }
  return iconMap[type] || '/icons/default.png'
}

const getProviderName = (type: string) => {
  const nameMap = {
    aliyun: '阿里云',
    tencent: '腾讯云',
    cloudflare: 'Cloudflare',
    dnspod: 'DNSPod',
    godaddy: 'GoDaddy'
  }
  return nameMap[type] || type
}

const getStatusBadge = (status: string) => {
  const statusMap = {
    active: 'success',
    inactive: 'default'
  }
  return statusMap[status] || 'default'
}

const getStatusText = (status: string) => {
  const statusMap = {
    active: '正常',
    inactive: '停用'
  }
  return statusMap[status] || status
}

const getPriorityColor = (priority: number) => {
  if (priority <= 3) return 'red'
  if (priority <= 6) return 'orange'
  return 'green'
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 事件处理
const handleAdd = () => {
  modalMode.value = 'add'
  currentProvider.value = null
  modalVisible.value = true
}

const handleEdit = (record: DNSProvider) => {
  modalMode.value = 'edit'
  currentProvider.value = record
  modalVisible.value = true
}

const handleView = (record: DNSProvider) => {
  modalMode.value = 'view'
  currentProvider.value = record
  modalVisible.value = true
}

const handleTest = async (record: DNSProvider) => {
  try {
    loading.value = true
    const result = await dnsProviderApi.test(record.id)
    if (result.success) {
      message.success('连接测试成功')
    } else {
      message.error(`连接测试失败: ${result.message}`)
    }
    await fetchData()
  } catch (error) {
    message.error('连接测试失败')
  } finally {
    loading.value = false
  }
}

const handleSetDefault = async (record: DNSProvider) => {
  try {
    await dnsProviderApi.setDefault(record.id)
    message.success('已设置为默认提供商')
    await fetchData()
  } catch (error) {
    message.error('设置默认提供商失败')
  }
}

const handleDelete = (record: DNSProvider) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除提供商 "${record.name}" 吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await dnsProviderApi.delete(record.id)
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
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 个提供商吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await dnsProviderApi.batchDelete(selectedRowKeys.value)
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
    const response = await dnsProviderApi.list(params)
    tableData.value = response.items || []
    pagination.total = response.total || 0
  } catch (error) {
    message.error('获取提供商列表失败')
  } finally {
    loading.value = false
  }
}

// 生命周期
onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.dns-provider-container {
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

.provider-name {
  .provider-link {
    font-weight: 500;
    color: #1890ff;
    text-decoration: none;
    margin-right: 8px;
    
    &:hover {
      color: #40a9ff;
    }
  }
}

.provider-type {
  display: flex;
  align-items: center;
}

.test-result {
  .test-status {
    margin-top: 4px;
  }
}

.danger {
  color: #ff4d4f !important;
}
</style>
