<template>
  <div class="dns-domain-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title">
          <h1>域名管理</h1>
          <p>管理您的域名资源，包括域名注册信息、分组管理和状态监控</p>
        </div>
        <div class="page-actions">
          <a-button type="primary" @click="handleAdd">
            <template #icon>
              <PlusOutlined />
            </template>
            添加域名
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
          <a-form-item label="域名" name="keyword">
            <a-input
              v-model:value="searchForm.keyword"
              placeholder="请输入域名关键词"
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
              <a-select-option value="active">正常</a-select-option>
              <a-select-option value="inactive">停用</a-select-option>
              <a-select-option value="expired">已过期</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="分组" name="group_id">
            <a-tree-select
              v-model:value="searchForm.group_id"
              :tree-data="groupTreeData"
              placeholder="请选择分组"
              allow-clear
              tree-default-expand-all
              style="width: 200px"
            />
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

    <!-- 域名列表 -->
    <div class="table-container">
      <a-card :bordered="false">
        <template #title>
          <div class="table-title">
            <span>域名列表</span>
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
          <!-- 域名名称列 -->
          <template #name="{ record }">
            <div class="domain-name">
              <a @click="handleView(record)" class="domain-link">
                {{ record.name }}
              </a>
              <div class="domain-tags" v-if="record.tags && record.tags.length">
                <a-tag
                  v-for="tag in record.tags"
                  :key="tag.id"
                  :color="tag.color"
                  size="small"
                >
                  {{ tag.name }}
                </a-tag>
              </div>
            </div>
          </template>

          <!-- 状态列 -->
          <template #status="{ record }">
            <a-badge
              :status="getStatusBadge(record.status)"
              :text="getStatusText(record.status)"
            />
          </template>

          <!-- 注册商列 -->
          <template #registrar="{ record }">
            <div class="registrar-info">
              <span>{{ getRegistrarName(record.registrar_type) }}</span>
            </div>
          </template>

          <!-- 过期时间列 -->
          <template #expires_at="{ record }">
            <div class="expire-time">
              <span :class="getExpireClass(record.expires_at)">
                {{ formatDate(record.expires_at) }}
              </span>
              <a-tag
                v-if="isExpiringSoon(record.expires_at)"
                color="orange"
                size="small"
              >
                即将过期
              </a-tag>
            </div>
          </template>

          <!-- 分组列 -->
          <template #group="{ record }">
            <span v-if="record.group">{{ record.group.name }}</span>
            <span v-else class="text-gray">未分组</span>
          </template>

          <!-- 操作列 -->
          <template #action="{ record }">
            <a-space>
              <a @click="handleView(record)">查看</a>
              <a @click="handleEdit(record)">编辑</a>
              <a-dropdown>
                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="handleManageRecords(record)">
                      管理解析
                    </a-menu-item>
                    <a-menu-item @click="handleManageCertificates(record)">
                      管理证书
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

    <!-- 域名表单弹窗 -->
    <DomainModal
      v-model:visible="modalVisible"
      :mode="modalMode"
      :domain="currentDomain"
      :group-options="groupOptions"
      @success="handleModalSuccess"
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
  DownOutlined
} from '@ant-design/icons-vue'
import type { TableColumnsType, TableProps } from 'ant-design-vue'
import { usePagination } from '@/hooks/usePagination'
import { useSelection } from '@/hooks/useSelection'
import DomainModal from './components/DomainModal.vue'
import { domainApi } from '@/api/dns/domain'
import { domainGroupApi } from '@/api/dns/domainGroup'
import type { Domain, DomainGroup } from '@/types/dns'

// 响应式数据
const loading = ref(false)
const modalVisible = ref(false)
const modalMode = ref<'add' | 'edit' | 'view'>('add')
const currentDomain = ref<Domain | null>(null)
const tableData = ref<Domain[]>([])
const groupOptions = ref<DomainGroup[]>([])
const groupTreeData = ref([])

// 搜索表单
const searchFormRef = ref()
const searchForm = reactive({
  keyword: '',
  status: undefined,
  group_id: undefined
})

// 分页和选择
const { pagination, handleTableChange } = usePagination()
const { selectedRowKeys, rowSelection, hasSelected } = useSelection()

// 表格列定义
const columns: TableColumnsType = [
  {
    title: '域名',
    dataIndex: 'name',
    key: 'name',
    slots: { customRender: 'name' },
    sorter: true
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    slots: { customRender: 'status' },
    width: 100
  },
  {
    title: '注册商',
    dataIndex: 'registrar_type',
    key: 'registrar_type',
    slots: { customRender: 'registrar' },
    width: 120
  },
  {
    title: '过期时间',
    dataIndex: 'expires_at',
    key: 'expires_at',
    slots: { customRender: 'expires_at' },
    sorter: true,
    width: 150
  },
  {
    title: '分组',
    dataIndex: 'group',
    key: 'group',
    slots: { customRender: 'group' },
    width: 120
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
    width: 150,
    fixed: 'right'
  }
]

// 计算属性和方法
const getStatusBadge = (status: string) => {
  const statusMap = {
    active: 'success',
    inactive: 'default',
    expired: 'error'
  }
  return statusMap[status] || 'default'
}

const getStatusText = (status: string) => {
  const statusMap = {
    active: '正常',
    inactive: '停用',
    expired: '已过期'
  }
  return statusMap[status] || status
}

const getRegistrarName = (type: string) => {
  const registrarMap = {
    godaddy: 'GoDaddy',
    aliyun: '阿里云',
    tencent: '腾讯云',
    dnspod: 'DNSPod',
    cloudflare: 'Cloudflare'
  }
  return registrarMap[type] || type
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const getExpireClass = (expiresAt: string) => {
  if (!expiresAt) return ''
  const expireDate = new Date(expiresAt)
  const now = new Date()
  const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  
  if (diffDays < 0) return 'text-red' // 已过期
  if (diffDays <= 30) return 'text-orange' // 30天内过期
  return ''
}

const isExpiringSoon = (expiresAt: string) => {
  if (!expiresAt) return false
  const expireDate = new Date(expiresAt)
  const now = new Date()
  const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  return diffDays > 0 && diffDays <= 30
}

// 事件处理方法
const handleAdd = () => {
  modalMode.value = 'add'
  currentDomain.value = null
  modalVisible.value = true
}

const handleEdit = (record: Domain) => {
  modalMode.value = 'edit'
  currentDomain.value = record
  modalVisible.value = true
}

const handleView = (record: Domain) => {
  modalMode.value = 'view'
  currentDomain.value = record
  modalVisible.value = true
}

const handleDelete = (record: Domain) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除域名 "${record.name}" 吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await domainApi.delete(record.id)
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
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 个域名吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await domainApi.batchDelete(selectedRowKeys.value)
        message.success('批量删除成功')
        selectedRowKeys.value = []
        await fetchData()
      } catch (error) {
        message.error('批量删除失败')
      }
    }
  })
}

const handleManageRecords = (record: Domain) => {
  // 跳转到DNS记录管理页面
  // router.push(`/dns/records?domain_id=${record.id}`)
}

const handleManageCertificates = (record: Domain) => {
  // 跳转到证书管理页面
  // router.push(`/dns/certificates?domain_id=${record.id}`)
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
    const response = await domainApi.list(params)
    tableData.value = response.items || []
    pagination.total = response.total || 0
  } catch (error) {
    message.error('获取域名列表失败')
  } finally {
    loading.value = false
  }
}

const fetchGroupOptions = async () => {
  try {
    const response = await domainGroupApi.tree()
    groupOptions.value = response || []
    groupTreeData.value = buildTreeData(response || [])
  } catch (error) {
    console.error('获取分组选项失败:', error)
  }
}

const buildTreeData = (groups: DomainGroup[]) => {
  return groups.map(group => ({
    title: group.name,
    value: group.id,
    key: group.id,
    children: group.children ? buildTreeData(group.children) : undefined
  }))
}

// 生命周期
onMounted(() => {
  fetchData()
  fetchGroupOptions()
})
</script>

<style scoped lang="scss">
.dns-domain-container {
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

.domain-name {
  .domain-link {
    font-weight: 500;
    color: #1890ff;
    text-decoration: none;
    
    &:hover {
      color: #40a9ff;
    }
  }
  
  .domain-tags {
    margin-top: 4px;
    
    .ant-tag {
      margin-right: 4px;
    }
  }
}

.registrar-info {
  display: flex;
  align-items: center;
}

.expire-time {
  .text-red {
    color: #ff4d4f;
  }
  
  .text-orange {
    color: #fa8c16;
  }
}

.text-gray {
  color: #8c8c8c;
}

.danger {
  color: #ff4d4f !important;
}
</style>
