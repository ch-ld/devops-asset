<template>
  <div class="dns-domain-management">
    <div class="page-container">
      <!-- 侧边栏折叠按钮 -->
      <div class="sidebar-toggle" v-show="!sidebarCollapsed">
        <el-button
          class="toggle-btn"
          circle
          size="small"
          @click="toggleSidebar"
          :icon="Fold"
          title="隐藏分组面板"
        />
      </div>

      <!-- 侧边栏展开按钮 -->
      <div class="sidebar-expand" v-show="sidebarCollapsed">
        <el-button
          class="expand-btn"
          circle
          size="small"
          @click="toggleSidebar"
          :icon="Expand"
          title="显示分组面板"
        />
      </div>

      <!-- 左侧域名分组边栏 -->
      <transition name="sidebar-slide">
        <DomainGroupSidebar
          v-show="!sidebarCollapsed"
          :groups="groups"
          :selected-group-id="selectedGroupId"
          :loading="groupLoading"
          :total-domain-count="allDomains.length"
          @group-select="handleGroupSelect"
          @add-group="handleAddGroup"
          @edit-group="handleEditGroup"
          @delete-group="handleDeleteGroup"
          @refresh="refreshGroups"
        />
      </transition>

      <!-- 主内容区 -->
      <div class="main-content" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
        <!-- 页面头部 -->
        <div class="page-header">
          <div class="header-top">
            <div class="title-section">
              <div class="title-wrapper">
                <div class="title-icon">
                  <el-icon size="24"><Monitor /></el-icon>
                </div>
                <div class="title-content">
                  <h1 class="page-title">域名管理</h1>
                  <div class="breadcrumb">
                    <span class="breadcrumb-item">DNS管理</span>
                    <el-icon class="breadcrumb-separator"><ArrowRight /></el-icon>
                    <span class="breadcrumb-current">{{ currentGroupName }}</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="stats-section">
              <div class="stat-card primary">
                <div class="stat-icon">
                  <el-icon><Operation /></el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-number">{{ statistics.total }}</div>
                  <div class="stat-label">总域名数</div>
                </div>
              </div>
              <div class="stat-card success">
                <div class="stat-icon">
                  <el-icon><CircleCheck /></el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-number">{{ statistics.active }}</div>
                  <div class="stat-label">正常解析</div>
                </div>
              </div>
              <div class="stat-card danger">
                <div class="stat-icon">
                  <el-icon><Warning /></el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-number">{{ statistics.error }}</div>
                  <div class="stat-label">解析异常</div>
                </div>
              </div>
              <div class="stat-card warning">
                <div class="stat-icon">
                  <el-icon><Clock /></el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-number">{{ statistics.expiring }}</div>
                  <div class="stat-label">即将过期</div>
                </div>
              </div>
            </div>
          </div>

          <!-- 操作按钮区域 -->
          <div class="header-actions">
            <div class="action-group">
              <el-button @click="handleImport" class="action-btn import-btn">
                <el-icon><Upload /></el-icon>
                <span>批量导入</span>
              </el-button>
              <el-button @click="handleExport" :loading="exportLoading" class="action-btn export-btn">
                <el-icon><Download /></el-icon>
                <span>导出</span>
              </el-button>
              <el-button type="primary" @click="handleAdd" class="action-btn primary-btn">
                <el-icon><Plus /></el-icon>
                <span>添加域名</span>
              </el-button>
            </div>
          </div>
        </div>

        <!-- 搜索和筛选区域 -->
        <div class="search-section">
          <el-card class="search-card">
            <el-form :inline="true" :model="searchForm" class="search-form">
              <el-form-item label="域名关键词">
                <el-input
                  v-model="searchForm.keyword"
                  placeholder="请输入域名关键词"
                  clearable
                  style="width: 200px"
                >
                  <template #prefix>
                    <el-icon><Search /></el-icon>
                  </template>
                </el-input>
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
              <el-form-item label="注册商">
                <el-select 
                  v-model="searchForm.registrar_type" 
                  placeholder="请选择注册商"
                  clearable
                  style="width: 150px"
                >
                  <el-option label="全部" value="" />
                  <el-option label="阿里云" value="aliyun" />
                  <el-option label="腾讯云" value="tencent" />
                  <el-option label="AWS Route53" value="route53" />
                  <el-option label="Cloudflare" value="cloudflare" />
                  <el-option label="DNSPod" value="dnspod" />
                  <el-option label="GoDaddy" value="godaddy" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleSearch" class="search-btn">
                  <el-icon><Search /></el-icon>
                  搜索
                </el-button>
                <el-button @click="handleReset" class="reset-btn">
                  <el-icon><Refresh /></el-icon>
                  重置
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </div>

        <!-- 域名列表区域 -->
        <div class="table-section">
          <el-card class="table-card">
            <template #header>
              <div class="table-header">
                <div class="header-left">
                  <h3 class="table-title">域名列表</h3>
                  <span class="table-subtitle">共 {{ pagination.total }} 个域名，当前显示第 {{ pagination.page }} 页</span>
                </div>
                <div class="table-actions">
                  <el-button @click="handleRefresh" class="refresh-btn">
                    <el-icon><Refresh /></el-icon>
                    刷新
                  </el-button>
                  <div v-if="hasSelected" class="selected-actions">
                    <span class="selected-info">已选择 {{ selectedRows.length }} 项</span>
                    <el-button 
                      type="danger" 
                      @click="handleBatchDelete" 
                      size="small"
                      class="danger-btn"
                    >
                      <el-icon><Delete /></el-icon>
                      批量删除
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
              class="domain-table"
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
                    :type="getStatusType(row.status)" 
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
                          <el-dropdown-item command="whois">WHOIS查询</el-dropdown-item>
                          <el-dropdown-item command="certificates">证书管理</el-dropdown-item>
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
        </div>
      </div>
    </div>

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

    <!-- 分组管理弹窗 -->
    <DomainGroupModal
      v-model:visible="groupModalVisible"
      :group="currentGroup"
      :parent-groups="groups"
      @success="handleGroupModalSuccess"
    />

    <!-- WHOIS查询弹窗 -->
    <WhoisModal
      v-model:visible="whoisModalVisible"
      :domain="currentDomain"
    />

    <!-- 证书管理弹窗 -->
    <CertificateModal
      v-model:visible="certificateModalVisible"
      :domain="currentDomain"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Upload,
  Download,
  ArrowDown,
  ArrowRight,
  Search,
  Refresh,
  Delete,
  Operation,
  CircleCheck,
  Warning,
  Clock,
  Monitor,
  Expand,
  Fold
} from '@element-plus/icons-vue'
import DomainGroupSidebar from './components/DomainGroupSidebar.vue'
import DomainModal from './components/DomainModal.vue'
import DomainImportModal from './components/DomainImportModal.vue'
import DomainGroupModal from './components/DomainGroupModal.vue'
import WhoisModal from './components/WhoisModal.vue'
import CertificateModal from './components/CertificateModal.vue'
import { domainApi } from '@/api/dns/domain'
import { domainGroupApi } from '@/api/dns/domainGroup'
import type { Domain, DomainGroup } from '@/types/dns'

// 响应式数据
const loading = ref(false)
const groupLoading = ref(false)
const exportLoading = ref(false)
const modalVisible = ref(false)
const importModalVisible = ref(false)
const groupModalVisible = ref(false)
const whoisModalVisible = ref(false)
const certificateModalVisible = ref(false)
const sidebarCollapsed = ref(false)

const currentDomain = ref<Domain | null>(null)
const currentGroup = ref<DomainGroup | null>(null)
const domains = ref<Domain[]>([])
const allDomains = ref<Domain[]>([])
const groups = ref<DomainGroup[]>([])
const selectedGroupId = ref<number | null>(null)

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
  registrar_type: '',
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

// 计算属性
const currentGroupName = computed(() => {
  if (selectedGroupId.value === null) return '全部域名'
  const group = groups.value.find(g => g.id === selectedGroupId.value)
  return group?.name || '未知分组'
})

// 工具方法
const formatDate = (date: string | number) => {
  if (!date) return '-'
  try {
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
  const typeMap: Record<string, any> = {
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
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

const handleAdd = () => {
  currentDomain.value = null
  modalVisible.value = true
}

const handleEdit = (row: Domain) => {
  currentDomain.value = row
  modalVisible.value = true
}

const handleRecords = (row: Domain) => {
  // 跳转到解析记录页面
  window.open(`/dns/records?domain=${row.name}`, '_blank')
}

const handleCommand = (command: string, row: Domain) => {
  currentDomain.value = row
  switch (command) {
    case 'whois':
      whoisModalVisible.value = true
      break
    case 'certificates':
      certificateModalVisible.value = true
      break
    case 'export':
      handleExportSingle(row)
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
    row.auto_renew = value
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
    registrar_type: '',
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

const handleExportSingle = async (domain: Domain) => {
  try {
    const blob = await domainApi.exportSingle(domain.id)
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${domain.name}-config.json`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
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
const handleGroupSelect = (group: DomainGroup) => {
  selectedGroupId.value = group.id
  searchForm.group_id = group.id || ''
  pagination.page = 1
  fetchData()
}

const handleAddGroup = () => {
  currentGroup.value = null
  groupModalVisible.value = true
}

const handleEditGroup = (group: DomainGroup) => {
  currentGroup.value = group
  groupModalVisible.value = true
}

const handleDeleteGroup = async (group: DomainGroup) => {
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
    
    await domainGroupApi.delete(group.id!)
    ElMessage.success('删除成功')
    await refreshGroups()
    await fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const refreshGroups = async () => {
  try {
    groupLoading.value = true
    const response = await domainGroupApi.list()
    groups.value = response.data?.items || []
  } catch (error) {
    ElMessage.error('获取分组列表失败')
  } finally {
    groupLoading.value = false
  }
}

const handleGroupModalSuccess = () => {
  groupModalVisible.value = false
  refreshGroups()
  fetchData()
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
      registrar_type: searchForm.registrar_type,
      group_id: typeof searchForm.group_id === 'string' && searchForm.group_id === '' ? undefined : Number(searchForm.group_id)
    }
    
    const response = await domainApi.list(params)
    
    if (response && (response as any).data && (response as any).data.items) {
      domains.value = (response as any).data.items || []
      pagination.total = (response as any).data.total || 0
    } else if (response && (response as any).items) {
      domains.value = (response as any).items || []
      pagination.total = (response as any).total || 0
    } else {
      domains.value = []
      pagination.total = 0
    }
  } catch (err) {
    console.error('获取域名列表失败:', err)
    ElMessage.error('获取域名列表失败')
  } finally {
    loading.value = false
  }
}

const fetchStatistics = async () => {
  try {
    const response = await domainApi.statistics()
    statistics.value = {
      total: response.total || 0,
      active: response.by_status?.active || 0,
      error: response.by_status?.error || 0,
      expiring: response.expiring || 0
    }
  } catch (err) {
    console.error('获取统计数据失败:', err)
    const total = pagination.total || domains.value.length
    const active = domains.value.filter(item => item.status === 'active').length
    const errorCount = domains.value.filter(item => item.status === 'inactive').length
    const expiring = domains.value.filter(item => 
      item.expires_at && isExpiringSoon(item.expires_at)
    ).length
    statistics.value = { total, active, error: errorCount, expiring }
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
onMounted(async () => {
  await Promise.all([
    refreshGroups(),
    fetchData(),
    fetchStatistics()
  ])
})
</script>

<style scoped lang="scss">
.dns-domain-management {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  overflow: hidden;
}

.page-container {
  display: flex;
  height: 100%;
  position: relative;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 24px;
  transition: margin-left 0.3s ease;

  &.sidebar-collapsed {
    margin-left: -320px; /* 侧边栏宽度 */
  }
  padding-left: 0;
}

.page-header {
  background: white;
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.title-section {
  flex: 1;
}

.title-wrapper {
  display: flex;
  align-items: center;
  gap: 16px;
}

.title-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.title-content {
  flex: 1;
}

.page-title {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6b7280;
  font-size: 14px;
}

.breadcrumb-separator {
  color: #d1d5db;
}

.breadcrumb-current {
  color: #3b82f6;
  font-weight: 500;
}

.stats-section {
  display: flex;
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-radius: 12px;
  min-width: 120px;
  
  &.primary {
    background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
    color: white;
  }
  
  &.success {
    background: linear-gradient(135deg, #10b981 0%, #059669 100%);
    color: white;
  }
  
  &.danger {
    background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
    color: white;
  }
  
  &.warning {
    background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
    color: white;
  }
}

.stat-icon {
  font-size: 24px;
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  opacity: 0.9;
}

.header-actions {
  margin-top: 16px;
}

.action-group {
  display: flex;
  gap: 12px;
}

.action-btn {
  padding: 12px 24px;
  border-radius: 12px;
  font-weight: 600;
  transition: all 0.3s ease;
  
  &.import-btn {
    background: linear-gradient(135deg, #10b981 0%, #059669 100%);
    color: white;
    border: none;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 20px rgba(16, 185, 129, 0.4);
    }
  }
  
  &.export-btn {
    background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
    color: white;
    border: none;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 20px rgba(245, 158, 11, 0.4);
    }
  }
  
  &.primary-btn {
    background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 20px rgba(59, 130, 246, 0.4);
    }
  }
}

.search-section {
  margin-bottom: 20px;
}

.search-card {
  border-radius: 16px;
  border: none;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.search-form {
  padding: 8px 0;
}

.search-btn {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  border: none;
}

.reset-btn {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #475569;
}

.table-section {
  flex: 1;
  overflow: hidden;
}

.table-card {
  height: 100%;
  border-radius: 16px;
  border: none;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
}

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
}

.refresh-btn {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #475569;
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

.danger-btn {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: white;
  border: none;
}

.domain-table {
  flex: 1;
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

/* 侧边栏折叠相关样式 */
.sidebar-toggle {
  position: fixed;
  top: 50%;
  left: 240px; /* 展开状态的侧边栏宽度 */
  z-index: 1000;
  transform: translate(-50%, -50%);
  transition: left 0.3s ease;

  .toggle-btn {
    background: #ffffff;
    border: 1px solid #e2e8f0;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    transition: all 0.3s ease;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;

    &:hover {
      background: #f8fafc;
      border-color: #3b82f6;
      color: #3b82f6;
      transform: scale(1.1);
    }
  }
}

.sidebar-expand {
  position: fixed;
  top: 50%;
  left: 70px; /* 折叠状态的侧边栏宽度 */
  z-index: 1000;
  transform: translate(-50%, -50%);

  .expand-btn {
    background: #3b82f6;
    border: 1px solid #3b82f6;
    color: #ffffff;
    box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
    transition: all 0.3s ease;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;

    &:hover {
      background: #2563eb;
      border-color: #2563eb;
      transform: scale(1.1);
    }
  }
}



/* 侧边栏滑动动画 */
.sidebar-slide-enter-active,
.sidebar-slide-leave-active {
  transition: all 0.3s ease;
}

.sidebar-slide-enter-from,
.sidebar-slide-leave-to {
  transform: translateX(-100%);
  opacity: 0;
}

.sidebar-slide-enter-to,
.sidebar-slide-leave-from {
  transform: translateX(0);
  opacity: 1;
}
</style>
