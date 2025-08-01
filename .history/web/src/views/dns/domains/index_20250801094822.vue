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
                <el-icon><Operation /></el-icon>
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
      <el-card shadow="never">
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

    <!-- 域名列表 -->
    <div class="table-container">
      <el-card shadow="never">
        <template #header>
          <div class="card-header">
            <span>域名列表</span>
            <div>
              <el-button @click="handleRefresh">
                <el-icon><Refresh /></el-icon>
                刷新
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
          :data="domains"
          :loading="loading"
          @selection-change="handleSelectionChange"
          row-key="id"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="name" label="域名" min-width="200">
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
          <el-table-column prop="expires_at" label="到期时间" width="180">
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
          <el-table-column prop="created_at" label="创建时间" width="140">
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
              <el-dropdown @command="(command) => handleCommand(command, row)">
                <el-button type="primary" size="small" text>
                  更多
                  <el-icon><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="certificates">证书管理</el-dropdown-item>
                    <el-dropdown-item command="whois">WHOIS查询</el-dropdown-item>
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
  Operation, 
  CircleCheck, 
  Warning, 
  Clock
} from '@element-plus/icons-vue'
import DomainModal from './components/DomainModal.vue'
import { domainApi } from '@/api/dns/domain'
import type { Domain } from '@/types/dns'

// 响应式数据
const loading = ref(false)
const modalVisible = ref(false)
const currentDomain = ref<Domain | null>(null)
const domains = ref<Domain[]>([])
const groups = ref<any[]>([])
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
  status: '',
  group_id: ''
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 选择
const selectedRows = ref<Domain[]>([])
const hasSelected = computed(() => selectedRows.value.length > 0)

// 工具方法
const getStatusType = (status: string) => {
  const statusMap: Record<string, any> = {
    active: 'success',
    inactive: 'info',
    error: 'danger',
    expired: 'warning'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '正常',
    inactive: '停用',
    error: '异常',
    expired: '过期'
  }
  return statusMap[status] || status
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const getExpiresClass = (expiresAt: string) => {
  if (!expiresAt) return ''
  const expireDate = new Date(expiresAt)
  const now = new Date()
  const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  
  if (diffDays < 0) return 'text-danger'
  if (diffDays <= 7) return 'text-danger'
  if (diffDays <= 30) return 'text-warning'
  return 'text-success'
}

const getExpiringDays = (expiresAt: string) => {
  if (!expiresAt) return ''
  const expireDate = new Date(expiresAt)
  const now = new Date()
  const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  return diffDays > 0 ? `${diffDays}天后过期` : `已过期${Math.abs(diffDays)}天`
}

// 事件处理
const handleAdd = () => {
  currentDomain.value = null
  modalVisible.value = true
}

const handleEdit = (row: Domain) => {
  currentDomain.value = row
  modalVisible.value = true
}

const handleRecords = (row: Domain) => {
  // TODO: 跳转到解析记录页面
  ElMessage.info(`查看 ${row.name} 的解析记录`)
}

const handleCommand = (command: string, row: Domain) => {
  switch (command) {
    case 'certificates':
      ElMessage.info(`查看 ${row.name} 的证书`)
      break
    case 'whois':
      ElMessage.info(`查询 ${row.name} 的WHOIS信息`)
      break
    case 'export':
      ElMessage.info(`导出 ${row.name} 的配置`)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const handleDelete = async (row: Domain) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除域名 "${row.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await domainApi.delete(row.id)
    ElMessage.success('删除成功')
    await fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleBatchDelete = async () => {
  if (!hasSelected.value) return
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 个域名吗？此操作不可恢复。`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const ids = selectedRows.value.map(row => row.id)
    await domainApi.batchDelete(ids)
    ElMessage.success('批量删除成功')
    selectedRows.value = []
    await fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

const handleSelectionChange = (rows: Domain[]) => {
  selectedRows.value = rows
}

const handleAutoRenewChange = async (row: Domain) => {
  try {
    await domainApi.update(row.id, { auto_renew: row.auto_renew })
    ElMessage.success(`${row.auto_renew ? '启用' : '禁用'}自动续费成功`)
  } catch (error) {
    row.auto_renew = !row.auto_renew // 回滚状态
    ElMessage.error('操作失败')
  }
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
    const response = await domainApi.list(params)
    domains.value = response.list || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('获取域名列表失败')
  } finally {
    loading.value = false
  }
}

const fetchStatistics = async () => {
  try {
    // 简化统计实现
    const total = domains.value.length
    const active = domains.value.filter(item => item.status === 'active').length
    const error = domains.value.filter(item => item.status === 'error').length
    const expiring = domains.value.filter(item => 
      item.expires_at && isExpiringSoon(item.expires_at)
    ).length
    
    statistics.value = { total, active, error, expiring }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const isExpiringSoon = (expiresAt: string) => {
  if (!expiresAt) return false
  const expireDate = new Date(expiresAt)
  const now = new Date()
  const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  return diffDays > 0 && diffDays <= 30
}

// 生命周期
onMounted(() => {
  fetchData()
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

.domain-cell {
  .domain-name {
    font-weight: 500;
    margin-bottom: 4px;
  }
  
  .domain-tags {
    display: flex;
    gap: 4px;
    flex-wrap: wrap;
  }
}

.expires-cell {
  .expires-date {
    margin-bottom: 2px;
  }
  
  .expires-days {
    font-size: 12px;
    
    &.text-danger {
      color: #f56c6c;
    }
    
    &.text-warning {
      color: #e6a23c;
    }
    
    &.text-success {
      color: #67c23a;
    }
  }
}
</style>
