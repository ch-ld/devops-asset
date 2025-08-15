<template>
  <div class="dns-domain-container">
    <!-- 侧边分组管理 -->
    <div class="sidebar" v-if="showSidebar">
      <div class="sidebar-header">
        <h3>
          <el-icon><FolderOpened /></el-icon>
          域名分组
        </h3>
        <div class="sidebar-actions">
          <el-button @click="handleAddGroup" type="primary" size="small" text>
            <el-icon><Plus /></el-icon>
          </el-button>
          <el-button @click="showSidebar = false" size="small" text>
            <el-icon><Close /></el-icon>
          </el-button>
        </div>
      </div>
      
      <div class="sidebar-content">
        <div class="group-list">
          <div 
            class="group-item" 
            :class="{ active: !searchForm.group_id }"
            @click="selectGroup(null)"
          >
            <el-icon><Folder /></el-icon>
            <span>全部域名</span>
            <span class="group-count">{{ statistics.total }}</span>
          </div>
          
          <div 
            v-for="group in groups" 
            :key="group.id"
            class="group-item"
            :class="{ active: searchForm.group_id === group.id }"
            @click="selectGroup(group.id)"
          >
            <el-icon><FolderOpened /></el-icon>
            <span>{{ group.name }}</span>
            <span class="group-count">{{ group.domain_count || 0 }}</span>
            <div class="group-actions">
              <el-button @click.stop="handleEditGroup(group)" size="small" text>
                <el-icon><Edit /></el-icon>
              </el-button>
              <el-button @click.stop="handleDeleteGroup(group)" size="small" text type="danger">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 主内容区域 -->
    <div class="main-content" :class="{ 'sidebar-open': showSidebar }">
      <!-- 页面头部 -->
      <div class="page-header">
        <div class="page-header-content">
          <div class="page-title">
            <h1>域名管理</h1>
            <p>管理DNS域名，包括域名添加、解析配置、状态监控等</p>
          </div>
          <div class="page-actions">
            <el-button @click="showSidebar = !showSidebar" class="page-action-btn sidebar-btn">
              <el-icon><Menu /></el-icon>
              <span>分组管理</span>
            </el-button>
            <el-button @click="handleImport" class="page-action-btn import-btn">
              <el-icon><Upload /></el-icon>
              <span>批量导入</span>
            </el-button>
            <el-button type="primary" @click="handleAdd" class="page-action-btn primary-btn">
              <el-icon><Plus /></el-icon>
              <span>添加域名</span>
            </el-button>
          </div>
        </div>
      </div>

      <!-- 统计卡片 -->
      <div class="stats-container">
        <el-row :gutter="24">
          <el-col :span="6">
            <el-card class="stats-card stats-card-primary">
              <div class="stats-content">
                <div class="stats-icon">
                  <el-icon size="32"><Operation /></el-icon>
                </div>
                <div class="stats-info">
                  <div class="stats-value">{{ statistics.total }}</div>
                  <div class="stats-label">总域名数</div>
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card stats-card-success">
              <div class="stats-content">
                <div class="stats-icon">
                  <el-icon size="32"><CircleCheck /></el-icon>
                </div>
                <div class="stats-info">
                  <div class="stats-value">{{ statistics.active }}</div>
                  <div class="stats-label">正常解析</div>
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card stats-card-danger">
              <div class="stats-content">
                <div class="stats-icon">
                  <el-icon size="32"><Warning /></el-icon>
                </div>
                <div class="stats-info">
                  <div class="stats-value">{{ statistics.error }}</div>
                  <div class="stats-label">解析异常</div>
                </div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card stats-card-warning">
              <div class="stats-content">
                <div class="stats-icon">
                  <el-icon size="32"><Clock /></el-icon>
                </div>
                <div class="stats-info">
                  <div class="stats-value">{{ statistics.expiring }}</div>
                  <div class="stats-label">即将过期</div>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 搜索和筛选 -->
      <el-card class="search-card">
        <el-form :inline="true" :model="searchForm" class="search-form">
          <el-form-item label="域名关键词">
            <el-input
              v-model="searchForm.keyword"
              placeholder="请输入域名关键词"
              clearable
              style="width: 200px"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select 
              v-model="searchForm.status" 
              placeholder="请选择状态"
              clearable
              style="width: 150px"
            >
              <el-option label="全部" value="" />
              <el-option label="正常" value="active" />
              <el-option label="禁用" value="inactive" />
              <el-option label="过期" value="expired" />
            </el-select>
          </el-form-item>
          <el-form-item label="分组">
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
            <el-button @click="handleExport" :loading="exportLoading">
              <el-icon><Download /></el-icon>
              导出
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 域名列表 -->
      <el-card class="table-card">
        <template #header>
          <div class="table-header">
            <div class="header-left">
              <h3 class="table-title">域名列表</h3>
              <span class="table-subtitle">共 {{ pagination.total }} 个域名，当前显示第 {{ pagination.page }} 页</span>
            </div>
            <div class="table-actions">
              <el-button type="primary" @click="handleAdd" class="modern-btn primary-btn">
                <el-icon><Plus /></el-icon>
                <span>新增域名</span>
              </el-button>
              <el-button @click="handleRefresh" class="modern-btn refresh-btn">
                <el-icon><Refresh /></el-icon>
                <span>刷新</span>
              </el-button>
              <div v-if="hasSelected" class="selected-actions">
                <span class="selected-info">已选择 {{ selectedRows.length }} 项</span>
                <el-button 
                  type="danger" 
                  @click="handleBatchDelete" 
                  size="small"
                  class="modern-btn danger-btn"
                >
                  <el-icon><Delete /></el-icon>
                  <span>批量删除</span>
                </el-button>
              </div>
            </div>
          </div>
        </template>

        <el-table 
          :data="domains" 
          :loading="loading"
          @selection-change="handleSelectionChange"
          row-key="id"
          style="width: 100%"
        >
          <el-table-column type="selection" width="55" />
          
          <el-table-column prop="name" label="域名" min-width="200">
            <template #default="{ row }">
              <div class="domain-cell">
                <div class="domain-name">{{ row.name }}</div>
                <div v-if="row.tags && row.tags.length > 0" class="domain-tags">
                  <el-tag 
                    v-for="tag in row.tags.slice(0, 2)" 
                    :key="tag" 
                    size="small" 
                    type="info"
                  >
                    {{ tag }}
                  </el-tag>
                  <el-tag v-if="row.tags.length > 2" size="small" type="info">
                    +{{ row.tags.length - 2 }}
                  </el-tag>
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag 
                :type="getStatusType(row.status) as any" 
                size="default"
                effect="light"
              >
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
                {{ getRegistrarName(row.registrar_type) }}
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
                :model-value="row.auto_renew"
                @change="(value) => handleAutoRenewChange(row, value as boolean)"
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
              <div class="action-buttons">
                <el-button type="primary" size="small" text @click="handleEdit(row)">
                  编辑
                </el-button>
                <el-button type="success" size="small" text @click="handleRecords(row)">
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
              </div>
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

      <!-- 域名表单弹窗 -->
      <DomainModal
        v-model:visible="modalVisible"
        :domain="currentDomain"
        :groups="groups"
        @success="handleModalSuccess"
      />
      
      <!-- 批量导入弹窗 -->
      <DomainImportModal
        v-model:visible="importModalVisible"
        @success="handleImportSuccess"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, 
  Upload,
  ArrowDown,
  Search,
  Refresh,
  Download,
  Delete,
  Operation,
  CircleCheck,
  Warning,
  Clock,
  Menu,
  FolderOpened,
  Folder,
  Edit,
  Close
} from '@element-plus/icons-vue'
import DomainModal from './components/DomainModal.vue'
import DomainImportModal from './components/DomainImportModal.vue'
import { domainApi } from '@/api/dns/domain'
import type { Domain } from '@/types/dns'

// 响应式数据
const loading = ref(false)
const statisticsLoading = ref(false)
const exportLoading = ref(false)
const modalVisible = ref(false)
const importModalVisible = ref(false)
const showSidebar = ref(false)
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
const searchForm = reactive({
  keyword: '',
  status: '',
  group_id: '' as string | number
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
const formatDate = (date: string | number) => {
  if (!date) return '-'
  try {
    // 如果是Unix时间戳（数字），先转换为毫秒
    const dateObj = typeof date === 'number' ? new Date(date * 1000) : new Date(date)
    return dateObj.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    })
  } catch (error) {
    console.error('时间格式化错误:', error, date)
    return '-'
  }
}

const getStatusType = (status: string) => {
  const typeMap: Record<string, string> = {
    'active': 'success',
    'inactive': 'danger',
    'expired': 'danger',
    'error': 'danger'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    'active': '正常',
    'inactive': '禁用',
    'expired': '过期',
    'error': '异常'
  }
  return textMap[status] || status
}

const getRegistrarName = (type: string) => {
  const nameMap: Record<string, string> = {
    'aliyun': '阿里云',
    'tencent': '腾讯云',
    'route53': 'AWS Route53',
    'cloudflare': 'Cloudflare',
    'dnspod': 'DNSPod',
    'godaddy': 'GoDaddy'
  }
  return nameMap[type] || type
}

const getExpiresClass = (expiresAt: string | number) => {
  if (!expiresAt) return ''
  const expireDate = typeof expiresAt === 'number' ? new Date(expiresAt * 1000) : new Date(expiresAt)
  const now = new Date()
  const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  
  if (diffDays < 0) return 'text-danger'
  if (diffDays <= 7) return 'text-danger'
  if (diffDays <= 30) return 'text-warning'
  return 'text-success'
}

const getExpiringDays = (expiresAt: string | number) => {
  if (!expiresAt) return ''
  const expireDate = typeof expiresAt === 'number' ? new Date(expiresAt * 1000) : new Date(expiresAt)
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

const handleAutoRenewChange = async (row: Domain, value: boolean) => {
  try {
    await domainApi.update(row.id, { auto_renew: value })
    row.auto_renew = value // 更新本地状态
    ElMessage.success(`${value ? '启用' : '禁用'}自动续费成功`)
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  Object.assign(searchForm, {
    keyword: '',
    status: '',
    group_id: ''
  })
  pagination.page = 1
  fetchData()
}

const handleImport = () => {
  importModalVisible.value = true
}

const handleExport = async () => {
  try {
    exportLoading.value = true
    const blob = await domainApi.export({
      format: 'xlsx',
      filters: { ...searchForm }
    })
    
    // 创建下载链接
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `domains-${new Date().toISOString().split('T')[0]}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  } finally {
    exportLoading.value = false
  }
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

const handleImportSuccess = () => {
  importModalVisible.value = false
  fetchData()
  fetchStatistics()
}

// 分组管理
const selectGroup = (groupId: number | null) => {
  searchForm.group_id = groupId || ''
  pagination.page = 1
  fetchData()
}

const handleAddGroup = () => {
  ElMessage.info('分组创建功能开发中')
}

const handleEditGroup = (group: any) => {
  ElMessage.info(`编辑分组: ${group.name}`)
}

const handleDeleteGroup = async (group: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除分组 "${group.name}" 吗？删除后该分组下的域名将变为未分组状态。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 数据获取
const fetchData = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchForm.keyword,
      status: searchForm.status,
      group_id: typeof searchForm.group_id === 'string' && searchForm.group_id === '' ? undefined : Number(searchForm.group_id)
    }
    console.log('发送请求参数:', params)
    const response = await domainApi.list(params)
    console.log('获取到的响应:', response)
    
    // 处理响应数据
    if (response && (response as any).data) {
      // 如果响应有data字段，取data.items
      domains.value = (response as any).data.items || []
      pagination.total = (response as any).data.total || 0
    } else {
      // 如果直接返回items，取items
      domains.value = (response as any).items || []
      pagination.total = (response as any).total || 0
    }
    
    console.log('设置的域名列表:', domains.value)
    console.log('设置的总数:', pagination.total)
  } catch (err) {
    console.error('获取域名列表失败:', err)
    ElMessage.error('获取域名列表失败')
  } finally {
    loading.value = false
  }
}

const fetchStatistics = async () => {
  try {
    statisticsLoading.value = true
    const response = await domainApi.statistics()
    statistics.value = {
      total: response.total || 0,
      active: response.by_status?.active || 0,
      error: response.by_status?.error || 0,
      expiring: response.expiring || 0
    }
  } catch (err) {
    console.error('获取统计数据失败:', err)
    // 降级处理：使用当前页面数据计算统计
    const total = pagination.total || domains.value.length
    const active = domains.value.filter(item => item.status === 'active').length
    const errorCount = domains.value.filter(item => item.status === 'inactive').length
    const expiring = domains.value.filter(item => 
      item.expires_at && isExpiringSoon(item.expires_at)
    ).length
    statistics.value = { total, active, error: errorCount, expiring }
  } finally {
    statisticsLoading.value = false
  }
}

const isExpiringSoon = (expiresAt: string | number) => {
  if (!expiresAt) return false
  const expireDate = typeof expiresAt === 'number' ? new Date(expiresAt * 1000) : new Date(expiresAt)
  const now = new Date()
  const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  return diffDays > 0 && diffDays <= 30
}

// 生命周期
onMounted(() => {
  // 初始化测试分组数据
  groups.value = [
    { id: 1, name: '生产环境', domain_count: 1 },
    { id: 2, name: '测试环境', domain_count: 1 },
    { id: 3, name: '开发环境', domain_count: 0 }
  ]
  
  fetchData()
  fetchStatistics()
})
</script>

<style scoped lang="scss">
.dns-domain-container {
  position: relative;
  display: flex;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  
  .sidebar {
    position: fixed;
    top: 0;
    left: 0;
    width: 280px;
    height: 100vh;
    background: white;
    border-right: 1px solid #e5e7eb;
    box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
    z-index: 2000;
    overflow-y: auto;
    
    .sidebar-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 20px 24px;
      border-bottom: 1px solid #e5e7eb;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      
      h3 {
        display: flex;
        align-items: center;
        gap: 8px;
        margin: 0;
        font-size: 16px;
        font-weight: 600;
      }
    }
    
    .sidebar-content {
      padding: 16px 0;
      
      .group-list {
        .group-item {
          display: flex;
          align-items: center;
          gap: 12px;
          padding: 12px 24px;
          cursor: pointer;
          transition: all 0.3s ease;
          border-left: 3px solid transparent;
          color: #374151;
          
          span {
            color: #374151;
            font-weight: 500;
          }
          
          .el-icon {
            color: #6b7280;
          }
          
          &:hover {
            background: #f3f4f6;
            
            .group-actions {
              display: flex;
            }
          }
          
          &.active {
            background: #eff6ff;
            border-left-color: #3b82f6;
            color: #3b82f6;
            font-weight: 600;
            
            .group-count {
              background: #dbeafe;
              color: #3b82f6;
            }
          }
          
          .group-count {
            font-size: 12px;
            color: #9ca3af;
            background: #f3f4f6;
            padding: 2px 8px;
            border-radius: 12px;
            min-width: 20px;
            text-align: center;
          }
          
          .group-actions {
            display: none;
            gap: 4px;
          }
        }
      }
    }
  }
  
  .main-content {
    flex: 1;
    padding: 24px;
    transition: margin-left 0.3s ease;
    
    &.sidebar-open {
      margin-left: 280px;
    }
  }
}

.page-header {
  background: white;
  border-radius: 16px;
  padding: 32px;
  margin-bottom: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  backdrop-filter: blur(10px);
  
  .page-header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    .page-title {
      h1 {
        font-size: 28px;
        font-weight: 700;
        color: #1f2937;
        margin: 0 0 8px 0;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
      }
      
      p {
        color: #6b7280;
        margin: 0;
        font-size: 16px;
      }
    }
    
    .page-actions {
      display: flex;
      gap: 12px;
      
      .page-action-btn {
        padding: 12px 24px;
        border-radius: 12px;
        font-weight: 600;
        transition: all 0.3s ease;
        border: 2px solid transparent;
        
        &.sidebar-btn {
          background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
          color: white;
          
          &:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 20px rgba(99, 102, 241, 0.4);
          }
        }
        
        &.import-btn {
          background: linear-gradient(135deg, #10b981 0%, #059669 100%);
          color: white;
          
          &:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 20px rgba(16, 185, 129, 0.4);
          }
        }
        
        &.primary-btn {
          background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
          color: white;
          
          &:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 20px rgba(59, 130, 246, 0.4);
          }
        }
        
        .el-icon {
          margin-right: 8px;
        }
      }
    }
  }
}

.stats-container {
  margin-bottom: 24px;
  
  .stats-card {
    border-radius: 16px;
    border: none;
    overflow: hidden;
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 12px 28px rgba(0, 0, 0, 0.12);
    }
    
    .stats-content {
      display: flex;
      align-items: center;
      gap: 20px;
      padding: 24px;
      
      .stats-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 64px;
        height: 64px;
        border-radius: 16px;
        color: white;
      }
      
      .stats-info {
        flex: 1;
        
        .stats-value {
          font-size: 32px;
          font-weight: 700;
          margin-bottom: 4px;
          background: linear-gradient(135deg, #1f2937 0%, #374151 100%);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
        }
        
        .stats-label {
          color: #6b7280;
          font-size: 14px;
          font-weight: 500;
        }
      }
    }
    
    &.stats-card-primary .stats-icon {
      background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
    }
    
    &.stats-card-success .stats-icon {
      background: linear-gradient(135deg, #10b981 0%, #059669 100%);
    }
    
    &.stats-card-danger .stats-icon {
      background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
    }
    
    &.stats-card-warning .stats-icon {
      background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
    }
  }
}

.search-card, .table-card {
  border-radius: 16px;
  border: none;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  margin-bottom: 24px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  .header-left {
    .table-title {
      font-size: 18px;
      font-weight: 600;
      color: #1e293b;
      margin: 0 0 4px 0;
    }
    
    .table-subtitle {
      color: #64748b;
      font-size: 14px;
    }
  }
  
  .table-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    
    .modern-btn {
      padding: 8px 16px;
      border-radius: 8px;
      font-weight: 500;
      transition: all 0.3s ease;
      
      &.primary-btn {
        background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
        border: none;
        color: white;
        
        &:hover {
          transform: translateY(-1px);
          box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
        }
      }
      
      &.refresh-btn {
        background: #f8fafc;
        border: 1px solid #e2e8f0;
        color: #475569;
        
        &:hover {
          background: #f1f5f9;
          border-color: #cbd5e1;
        }
      }
      
      &.danger-btn {
        background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
        border: none;
        color: white;
        
        &:hover {
          transform: translateY(-1px);
          box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
        }
      }
    }
    
    .selected-actions {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 8px 16px;
      background: #fef3c7;
      border: 1px solid #fbbf24;
      border-radius: 8px;
      
      .selected-info {
        color: #92400e;
        font-size: 14px;
        font-weight: 500;
      }
    }
  }
}

.domain-cell {
  .domain-name {
    font-weight: 600;
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
    font-weight: 500;
    margin-bottom: 2px;
  }
  
  .expires-days {
    font-size: 12px;
    
    &.text-danger {
      color: #ef4444;
    }
    
    &.text-warning {
      color: #f59e0b;
    }
    
    &.text-success {
      color: #10b981;
    }
  }
}

.action-buttons {
  display: flex;
  gap: 8px;
  align-items: center;
}

.pagination-container {
  padding: 20px 24px;
  background: #f8fafc;
  border-top: 1px solid #e2e8f0;
  
  :deep(.el-pagination) {
    justify-content: center;
  }
}
</style>
