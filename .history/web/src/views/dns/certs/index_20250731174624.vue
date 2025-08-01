<template>
  <div class="dns-cert-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title">
          <h1>证书管理</h1>
          <p>管理SSL/TLS证书，包括申请、续期、部署和监控</p>
        </div>
        <div class="page-actions">
          <a-button type="primary" @click="handleAdd">
            <template #icon>
              <PlusOutlined />
            </template>
            申请证书
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
              <a-select-option value="pending">申请中</a-select-option>
              <a-select-option value="issued">已签发</a-select-option>
              <a-select-option value="expired">已过期</a-select-option>
              <a-select-option value="revoked">已吊销</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="CA类型" name="ca_type">
            <a-select
              v-model:value="searchForm.ca_type"
              placeholder="请选择CA类型"
              allow-clear
              style="width: 150px"
            >
              <a-select-option value="letsencrypt">Let's Encrypt</a-select-option>
              <a-select-option value="zerossl">ZeroSSL</a-select-option>
              <a-select-option value="buypass">Buypass</a-select-option>
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

    <!-- 证书列表 -->
    <div class="table-container">
      <a-card :bordered="false">
        <template #title>
          <div class="table-title">
            <span>证书列表</span>
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
          <!-- 域名列 -->
          <template #common_name="{ record }">
            <div class="cert-domain">
              <div class="primary-domain">{{ record.common_name }}</div>
              <div v-if="record.subject_alt_names && record.subject_alt_names.length" class="alt-domains">
                <a-tag
                  v-for="domain in record.subject_alt_names.slice(0, 3)"
                  :key="domain"
                  size="small"
                >
                  {{ domain }}
                </a-tag>
                <a-tag v-if="record.subject_alt_names.length > 3" size="small">
                  +{{ record.subject_alt_names.length - 3 }}
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

          <!-- CA类型列 -->
          <template #ca_type="{ record }">
            <span>{{ getCATypeName(record.ca_type) }}</span>
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

          <!-- 自动续期列 -->
          <template #auto_renew="{ record }">
            <a-switch
              :checked="record.auto_renew"
              @change="(checked) => handleToggleAutoRenew(record, checked)"
              size="small"
            />
          </template>

          <!-- 操作列 -->
          <template #action="{ record }">
            <a-space>
              <a @click="handleView(record)">查看</a>
              <a @click="handleDownload(record)">下载</a>
              <a-dropdown>
                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="handleRenew(record)" :disabled="record.status !== 'issued'">
                      续期
                    </a-menu-item>
                    <a-menu-item @click="handleDeploy(record)" :disabled="record.status !== 'issued'">
                      部署
                    </a-menu-item>
                    <a-menu-divider />
                    <a-menu-item @click="handleRevoke(record)" :disabled="record.status !== 'issued'" class="danger">
                      吊销
                    </a-menu-item>
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

    <!-- 证书表单弹窗 -->
    <CertificateModal
      v-model:visible="modalVisible"
      :mode="modalMode"
      :certificate="currentCertificate"
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
import CertificateModal from './components/CertificateModal.vue'
import { certificateApi } from '@/api/dns/certificate'
import type { Certificate } from '@/types/dns'

// 响应式数据
const loading = ref(false)
const modalVisible = ref(false)
const modalMode = ref<'add' | 'edit' | 'view'>('add')
const currentCertificate = ref<Certificate | null>(null)
const tableData = ref<Certificate[]>([])

// 搜索表单
const searchFormRef = ref()
const searchForm = reactive({
  keyword: '',
  status: undefined,
  ca_type: undefined
})

// 分页和选择
const { pagination, handleTableChange } = usePagination()
const { selectedRowKeys, rowSelection, hasSelected } = useSelection()

// 表格列定义
const columns: TableColumnsType = [
  {
    title: '域名',
    dataIndex: 'common_name',
    key: 'common_name',
    slots: { customRender: 'common_name' }
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    slots: { customRender: 'status' },
    width: 100
  },
  {
    title: 'CA类型',
    dataIndex: 'ca_type',
    key: 'ca_type',
    slots: { customRender: 'ca_type' },
    width: 120
  },
  {
    title: '签发时间',
    dataIndex: 'issued_at',
    key: 'issued_at',
    render: (text: string) => formatDate(text),
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
    title: '自动续期',
    dataIndex: 'auto_renew',
    key: 'auto_renew',
    slots: { customRender: 'auto_renew' },
    width: 100
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

// 工具方法
const getStatusBadge = (status: string) => {
  const statusMap = {
    pending: 'processing',
    issued: 'success',
    expired: 'error',
    revoked: 'default'
  }
  return statusMap[status] || 'default'
}

const getStatusText = (status: string) => {
  const statusMap = {
    pending: '申请中',
    issued: '已签发',
    expired: '已过期',
    revoked: '已吊销'
  }
  return statusMap[status] || status
}

const getCATypeName = (type: string) => {
  const typeMap = {
    letsencrypt: "Let's Encrypt",
    zerossl: 'ZeroSSL',
    buypass: 'Buypass'
  }
  return typeMap[type] || type
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
  
  if (diffDays < 0) return 'text-red'
  if (diffDays <= 30) return 'text-orange'
  return ''
}

const isExpiringSoon = (expiresAt: string) => {
  if (!expiresAt) return false
  const expireDate = new Date(expiresAt)
  const now = new Date()
  const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  return diffDays > 0 && diffDays <= 30
}

// 事件处理
const handleAdd = () => {
  modalMode.value = 'add'
  currentCertificate.value = null
  modalVisible.value = true
}

const handleView = (record: Certificate) => {
  modalMode.value = 'view'
  currentCertificate.value = record
  modalVisible.value = true
}

const handleRenew = async (record: Certificate) => {
  try {
    await certificateApi.renew(record.id)
    message.success('证书续期请求已提交')
    await fetchData()
  } catch (error) {
    message.error('证书续期失败')
  }
}

const handleRevoke = (record: Certificate) => {
  Modal.confirm({
    title: '确认吊销',
    content: `确定要吊销证书 "${record.common_name}" 吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await certificateApi.revoke(record.id)
        message.success('证书吊销成功')
        await fetchData()
      } catch (error) {
        message.error('证书吊销失败')
      }
    }
  })
}

const handleDownload = async (record: Certificate) => {
  try {
    const blob = await certificateApi.download(record.id)
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${record.common_name}.pem`
    a.click()
    window.URL.revokeObjectURL(url)
  } catch (error) {
    message.error('证书下载失败')
  }
}

const handleDeploy = (record: Certificate) => {
  // TODO: 实现证书部署功能
  message.info('证书部署功能开发中')
}

const handleToggleAutoRenew = async (record: Certificate, checked: boolean) => {
  try {
    // TODO: 实现自动续期切换
    message.success(`已${checked ? '开启' : '关闭'}自动续期`)
  } catch (error) {
    message.error('操作失败')
  }
}

const handleDelete = (record: Certificate) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除证书 "${record.common_name}" 吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await certificateApi.delete(record.id)
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
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 个证书吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await certificateApi.batchDelete(selectedRowKeys.value)
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
    const response = await certificateApi.list(params)
    tableData.value = response.items || []
    pagination.total = response.total || 0
  } catch (error) {
    message.error('获取证书列表失败')
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
.dns-cert-container {
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

.cert-domain {
  .primary-domain {
    font-weight: 500;
    color: #262626;
  }
  
  .alt-domains {
    margin-top: 4px;
    
    .ant-tag {
      margin-right: 4px;
    }
  }
}

.expire-time {
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
