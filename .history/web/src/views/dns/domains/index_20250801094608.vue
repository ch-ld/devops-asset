<template>
  <div class="dns-domain-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title">
          <h1>域名管理</h1>
          <p>管理DNS域名，包括域名添加、解析配置、状态监控等</p>
        </div>
        <div class="page-actions">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加域名
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
              title="总域名数"
              :value="statistics.total"
              :value-style="{ color: '#409eff' }"
            >
              <template #suffix>
                <el-icon><Globe /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="正常解析"
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
              title="解析异常"
              :value="statistics.error"
              :value-style="{ color: '#f56c6c' }"
            >
              <template #suffix>
                <el-icon><Warning /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="即将过期"
              :value="statistics.expiring"
              :value-style="{ color: '#e6a23c' }"
            >
              <template #suffix>
                <el-icon><Clock /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-container">
      <el-card :shadow="'never'">
        <el-form
          ref="searchFormRef"
          :model="searchForm"
          inline
          class="search-form"
        >
          <el-form-item label="域名" prop="keyword">
            <el-input
              v-model="searchForm.keyword"
              placeholder="请输入域名关键词"
              clearable
              style="width: 200px"
            />
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-select
              v-model="searchForm.status"
              placeholder="请选择状态"
              clearable
              style="width: 150px"
            >
              <el-option label="全部" value="" />
              <el-option label="正常" value="active" />
              <el-option label="异常" value="error" />
              <el-option label="过期" value="expired" />
            </el-select>
          </el-form-item>
          <el-form-item label="分组" prop="group_id">
            <el-select
              v-model="searchForm.group_id"
              placeholder="请选择分组"
              clearable
              style="width: 150px"
            >
              <el-option label="全部" value="" />
              <el-option 
                v-for="group in groups" 
                :key="group.id" 
                :label="group.name" 
                :value="group.id" 
              />
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

    <!-- 操作工具栏 -->
    <div class="toolbar-container">
      <el-card :shadow="'never'">
        <div class="toolbar-content">
          <div class="toolbar-left">
            <el-button @click="handleRefresh">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button @click="handleBatchDelete" :disabled="!hasSelected">
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
            <el-button @click="handleExport">
              <el-icon><Download /></el-icon>
              导出
            </el-button>
            <el-button @click="handleImport">
              <el-icon><Upload /></el-icon>
              导入
            </el-button>
          </div>
          <div class="toolbar-right">
            <el-button-group>
              <el-button :type="viewMode === 'table' ? 'primary' : ''" @click="viewMode = 'table'">
                <el-icon><List /></el-icon>
              </el-button>
              <el-button :type="viewMode === 'card' ? 'primary' : ''" @click="viewMode = 'card'">
                <el-icon><Grid /></el-icon>
              </el-button>
            </el-button-group>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 域名列表 -->
    <div class="table-container">
      <el-card :shadow="'never'">
        <el-table
          ref="tableRef"
          :data="domains"
          :loading="loading"
          @selection-change="handleSelectionChange"
          @sort-change="handleSortChange"
          row-key="id"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="name" label="域名" sortable min-width="200">
            <template #default="{ row }">
              <div class="domain-cell">
                <div class="domain-name">{{ row.name }}</div>
                <div class="domain-tags" v-if="row.tags && row.tags.length">
                  <el-tag
                    v-for="tag in row.tags"
                    :key="tag.id"
                    :color="tag.color"
                    size="small"
                  >
                    {{ tag.name }}
                  </el-tag>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="group" label="分组" width="120">
            <template #default="{ row }">
              {{ row.group?.name || '未分组' }}
            </template>
          </el-table-column>
          <el-table-column prop="registrar_type" label="注册商" width="120">
            <template #default="{ row }">
              <el-tag v-if="row.registrar_type" type="info" size="small">
                {{ row.registrar_type }}
              </el-tag>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column prop="expires_at" label="到期时间" width="140" sortable>
            <template #default="{ row }">
              <div v-if="row.expires_at" class="expires-cell">
                <div class="expires-date">{{ formatDate(row.expires_at) }}</div>
                <div class="expires-days" :class="getExpiresClass(row.expires_at)">
                  {{ getExpiringDays(row.expires_at) }}
                </div>
              </div>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column prop="auto_renew" label="自动续费" width="100">
            <template #default="{ row }">
              <el-switch
                v-model="row.auto_renew"
                @change="handleAutoRenewChange(row)"
              />
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="140" sortable>
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
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
                @click="handleRecords(row)"
              >
                解析记录
              </el-button>
              <el-button
                type="warning"
                size="small"
                text
                @click="handleCertificates(row)"
              >
                证书
              </el-button>
              <el-dropdown @command="(command) => handleCommand(command, row)">
                <el-button type="primary" size="small" text>
                  更多
                  <el-icon><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="whois">WHOIS查询</el-dropdown-item>
                    <el-dropdown-item command="ping">Ping检测</el-dropdown-item>
                    <el-dropdown-item command="export">导出配置</el-dropdown-item>
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

    <!-- 域名表单弹窗 -->
    <DomainModal
      v-model:visible="modalVisible"
      :domain="currentDomain"
      :groups="groups"
      @success="handleModalSuccess"
    />

    <!-- 批量导入弹窗 -->
    <ImportModal
      v-model:visible="importVisible"
      @success="handleImportSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, 
  Search, 
  Refresh, 
  Delete, 
  ArrowDown, 
  Operation, 
  CircleCheck, 
  Warning, 
  Clock,
  Download,
  Upload,
  List,
  Grid
} from '@element-plus/icons-vue'
import { usePagination } from '@/hooks/usePagination'
import { useSelection } from '@/hooks/useSelection'
import DomainModal from './components/DomainModal.vue'
import ImportModal from './components/ImportModal.vue'
import { domainApi } from '@/api/dns/domain'
import { dnsProviderApi } from '@/api/dns/provider'
import type { Domain, DNSProvider } from '@/types/dns'

// 响应式数据
const loading = ref(false)
const modalVisible = ref(false)
const modalMode = ref<'add' | 'edit' | 'view'>('add')
const currentDomain = ref<Domain | null>(null)
const tableData = ref<Domain[]>([])
const providerOptions = ref<DNSProvider[]>([])
const statistics = ref({
  total: 0,
  active: 0,
  error: 0,
  expiring: 0
})

// 搜索表单
const searchFormRef = ref()
const searchForm = reactive({
  keyword: '',
  status: undefined,
  provider_id: undefined
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
    title: 'DNS提供商',
    dataIndex: 'provider',
    key: 'provider',
    slots: { customRender: 'provider' },
    width: 150
  },
  {
    title: '解析记录',
    dataIndex: 'record_count',
    key: 'record_count',
    slots: { customRender: 'record_count' },
    width: 100
  },
  {
    title: '过期时间',
    dataIndex: 'expires_at',
    key: 'expires_at',
    slots: { customRender: 'expires_at' },
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
    width: 200,
    fixed: 'right'
  }
]

// 工具方法
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

const getProviderTypeName = (type: string) => {
  const typeMap = {
    aliyun: '阿里云',
    tencent: '腾讯云',
    cloudflare: 'Cloudflare',
    dnspod: 'DNSPod',
    godaddy: 'GoDaddy'
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

const handleViewRecords = (record: Domain) => {
  // TODO: 跳转到解析记录页面
  message.info(`查看 ${record.name} 的解析记录`)
}

const handleSync = async (record: Domain) => {
  try {
    await domainApi.sync(record.id)
    message.success('域名记录同步成功')
    await fetchData()
  } catch (error) {
    message.error('域名记录同步失败')
  }
}

const handleToggleStatus = async (record: Domain) => {
  try {
    const newStatus = record.status === 'active' ? 'inactive' : 'active'
    await domainApi.update(record.id, { status: newStatus })
    message.success(`域名已${newStatus === 'active' ? '启用' : '停用'}`)
    await fetchData()
  } catch (error) {
    message.error('操作失败')
  }
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
    const response = await domainApi.list(params)
    tableData.value = response.items || []
    pagination.total = response.total || 0
  } catch (error) {
    message.error('获取域名列表失败')
  } finally {
    loading.value = false
  }
}

const fetchProviderOptions = async () => {
  try {
    const response = await dnsProviderApi.list({ page: 1, size: 100 })
    providerOptions.value = response.items || []
  } catch (error) {
    console.error('获取DNS提供商列表失败:', error)
  }
}

const fetchStatistics = async () => {
  try {
    // TODO: 实现统计数据获取
    // 暂时使用模拟数据
    const total = tableData.value.length
    const active = tableData.value.filter(item => item.status === 'active').length
    const error = tableData.value.filter(item => item.status === 'error').length
    const expiring = tableData.value.filter(item => 
      item.expires_at && isExpiringSoon(item.expires_at)
    ).length
    
    statistics.value = { total, active, error, expiring }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 生命周期
onMounted(() => {
  fetchData()
  fetchProviderOptions()
  fetchStatistics()
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

.domain-name {
  .domain-link {
    font-weight: 500;
    color: #1890ff;
    text-decoration: none;
    
    &:hover {
      color: #40a9ff;
    }
  }
  
  .domain-info {
    margin-top: 4px;
    
    .ant-tag {
      margin-right: 4px;
    }
  }
}

.provider-info {
  .provider-type {
    font-size: 12px;
    color: #8c8c8c;
    margin-top: 2px;
  }
}

.record-count {
  color: #1890ff;
  text-decoration: none;
  
  &:hover {
    color: #40a9ff;
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
