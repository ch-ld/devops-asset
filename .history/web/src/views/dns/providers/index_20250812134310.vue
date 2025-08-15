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
          <el-button type="success" @click="handleSyncAllDomains" :loading="syncLoading">
            <el-icon><Refresh /></el-icon>
            一键同步域名
          </el-button>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加提供商
          </el-button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-container">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card class="stats-card stats-card-primary">
            <div class="stats-content">
              <div class="stats-icon">
                <el-icon size="24"><Setting /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-value">{{ statistics.total }}</div>
                <div class="stats-label">提供商总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stats-card stats-card-success">
            <div class="stats-content">
              <div class="stats-icon">
                <el-icon size="24"><CircleCheck /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-value">{{ statistics.enabled }}</div>
                <div class="stats-label">已启用</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stats-card stats-card-danger">
            <div class="stats-content">
              <div class="stats-icon">
                <el-icon size="24"><CircleClose /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-value">{{ statistics.disabled }}</div>
                <div class="stats-label">已禁用</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stats-card stats-card-warning">
            <div class="stats-content">
              <div class="stats-icon">
                <el-icon size="24"><Warning /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-value">{{ statistics.error }}</div>
                <div class="stats-label">连接异常</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 搜索和筛选 -->
    <el-card class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="提供商名称">
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入提供商名称"
            clearable
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="提供商类型">
          <el-select 
            v-model="searchForm.type" 
            placeholder="请选择类型"
            clearable
            style="width: 150px"
          >
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="Cloudflare" value="cloudflare" />
            <el-option label="DNSPod" value="dnspod" />
            <el-option label="GoDaddy" value="godaddy" />
            <el-option label="AWS Route53" value="route53" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select 
            v-model="searchForm.status" 
            placeholder="请选择状态"
            clearable
            style="width: 120px"
          >
            <el-option label="全部" value="" />
            <el-option label="启用" value="enabled" />
            <el-option label="禁用" value="disabled" />
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

    <!-- 提供商列表 -->
    <el-card class="table-card">
      <template #header>
        <div class="table-header">
          <div class="header-left">
            <h3 class="table-title">提供商列表</h3>
            <span class="table-subtitle">管理您的DNS提供商配置</span>
          </div>
          <div class="table-actions">
            <el-button type="primary" @click="handleAdd" class="action-btn primary-btn">
              <el-icon><Plus /></el-icon>
              新增提供商
            </el-button>
            <el-button @click="handleSyncAllDomains" class="action-btn sync-btn">
              <el-icon><Refresh /></el-icon>
              一键同步域名
            </el-button>
            <el-button @click="handleTestAll" class="action-btn">
              <el-icon><Connection /></el-icon>
              批量测试
            </el-button>
            <el-button 
              type="danger" 
              @click="handleBatchDelete" 
              :disabled="selectedRowKeys.length === 0"
              class="action-btn danger-btn"
            >
              <el-icon><Delete /></el-icon>
              批量删除 <span v-if="selectedRowKeys.length > 0">({{ selectedRowKeys.length }})</span>
            </el-button>
          </div>
        </div>
      </template>

      <template v-if="loading">
        <el-skeleton :rows="5" animated style="padding: 32px" />
      </template>

      <el-table
        v-else
        v-loading="loading"
        :data="providerList"
        @selection-change="handleSelectionChange"
        stripe
        style="width: 100%"
        class="provider-table"
        :header-cell-style="{ background: '#f8f9fa', color: '#606266', fontWeight: '600' }"
        :row-style="{ height: '60px' }"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="name" label="提供商信息" width="250">
          <template #default="{ row }">
            <div class="provider-info">
              <div class="provider-avatar">
                <el-avatar :size="40" class="provider-logo">
                  <el-icon size="20">
                    <component :is="getProviderIcon(row.type)" />
                  </el-icon>
                </el-avatar>
              </div>
              <div class="provider-details">
                <div class="provider-name">{{ row.name }}</div>
                                 <div class="provider-type">
                   <el-tag size="small" :type="getProviderTagType(row.type) as any">
                     {{ getProviderTypeName(row.type) }}
                   </el-tag>
                 </div>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <div class="status-wrapper">
              <el-tag 
                :type="getStatusType(row.status) as any" 
                size="default"
                effect="light"
                class="status-tag"
              >
                <el-icon class="status-icon">
                  <component :is="getStatusIcon(row.status)" />
                </el-icon>
                {{ getStatusText(row.status) }}
              </el-tag>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="连接状态" width="140">
          <template #default="{ row }">
            <div class="connection-status">
              <el-tag 
                v-if="row.last_test_result" 
                :type="row.last_test_result.success ? 'success' : 'danger'"
                size="default"
                effect="light"
                class="connection-tag"
              >
                <el-icon class="connection-icon">
                  <component :is="row.last_test_result.success ? 'CircleCheck' : 'CircleClose'" />
                </el-icon>
                {{ row.last_test_result.success ? '连接正常' : '连接失败' }}
              </el-tag>
              <el-tag 
                v-else-if="row.test_result === 'success'" 
                type="success"
                size="default"
                effect="light"
                class="connection-tag"
              >
                <el-icon class="connection-icon">
                  <CircleCheck />
                </el-icon>
                连接正常
              </el-tag>
              <el-tag 
                v-else-if="row.test_result === 'failed'" 
                type="danger"
                size="default"
                effect="light"
                class="connection-tag"
              >
                <el-icon class="connection-icon">
                  <CircleClose />
                </el-icon>
                连接失败
              </el-tag>
              <el-tag v-else type="info" effect="light" class="connection-tag">
                <el-icon class="connection-icon">
                  <Warning />
                </el-icon>
                未测试
              </el-tag>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="域名数量" width="100" align="center">
          <template #default="{ row }">
            <div class="domain-count">
              <el-badge :value="row.domain_count || 0" :max="999" class="domain-badge">
                <el-icon size="16" color="#409eff">
                  <Setting />
                </el-icon>
              </el-badge>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            <div class="time-info">
              <div class="created-time">
                <el-icon size="14" class="time-icon">
                  <Clock />
                </el-icon>
                <span class="time-text">{{ formatTime(row.created_at) }}</span>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="remark" label="备注" show-overflow-tooltip />

        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button text @click="handleTest(row)">
              <el-icon><Connection /></el-icon>
              测试
            </el-button>
            <el-button text @click="handleEdit(row)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-dropdown>
              <el-button text>
                更多<el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleSyncDomains(row)">
                    <el-icon><Refresh /></el-icon>
                    同步域名
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleToggleStatus(row)">
                    {{ row.status === 'active' ? '禁用' : '启用' }}
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleClone(row)">克隆配置</el-dropdown-item>
                  <el-dropdown-item divided @click="handleDelete(row)">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handlePageSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 提供商弹窗 -->
    <ProviderModal
      v-model:visible="providerModalVisible"
      :mode="currentProvider ? 'edit' : 'add'"
      :provider="currentProvider"
      @success="handleProviderSuccess"
    />

    <!-- 同步进度对话框 -->
    <!-- 
    <SyncProgressDialog
      v-model:visible="syncProgressVisible"
      :providers="enabledProviders"
      @cancel="handleSyncCancel"
    />
    -->
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
  Edit,
  Connection,
  Setting,
  CircleCheck,
  CircleClose,
  Warning,
  Clock,
  Cloudy
} from '@element-plus/icons-vue'
import ProviderModal from './components/ProviderModal.vue'
// import SyncProgressDialog from './components/SyncProgressDialog.vue'
import { dnsProviderApi } from '@/api/dns/provider'

// 响应式数据
const loading = ref(false)
const syncLoading = ref(false)
const providerList = ref<any[]>([])
const selectedRowKeys = ref<number[]>([])
const providerModalVisible = ref(false)
// const syncProgressVisible = ref(false)
const currentProvider = ref<any>(null)

// 搜索表单
const searchForm = reactive({
  keyword: '',
  type: '',
  status: ''
})

// 分页
const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

// 统计数据
const statistics = reactive({
  total: 0,
  enabled: 0,
  disabled: 0,
  error: 0
})

// 获取提供商列表
const fetchProviderList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: searchForm.keyword,
      type: searchForm.type,
      status: searchForm.status
    }

    const response = await dnsProviderApi.list(params)
    // 后端返回的数据结构: { code, status, message, data: { items, total }, timestamp }
    const data = (response as any).data || response
    providerList.value = data.items || []
    pagination.total = data.total || 0
    fetchStatistics()
  } catch (error) {
    console.error('获取提供商列表失败:', error)
    ElMessage.error('获取提供商列表失败')
  } finally {
    loading.value = false
  }
}

// 获取统计数据
const fetchStatistics = async () => {
  try {
    const response = await dnsProviderApi.statistics()
    statistics.total = response.total || 0
    statistics.enabled = response.by_status?.active || 0
    statistics.disabled = response.by_status?.inactive || 0
    statistics.error = providerList.value.filter(item => 
      item.lastTestResult && !item.lastTestResult.success
    ).length
  } catch (err) {
    console.error('获取提供商统计数据失败:', err)
    // 降级处理：使用当前页面数据计算统计
    statistics.total = providerList.value.length
    statistics.enabled = providerList.value.filter(item => item.status === 'active').length
    statistics.disabled = providerList.value.filter(item => item.status === 'inactive').length
    statistics.error = providerList.value.filter(item => 
      item.lastTestResult && !item.lastTestResult.success
    ).length
  }
}

// 计算启用的提供商列表
const enabledProviders = computed(() => {
  return providerList.value.filter(provider => provider.status === 'active')
})

// 处理选择变化
const handleSelectionChange = (selection: any[]) => {
  selectedRowKeys.value = selection.map(item => item.id)
}

// 处理搜索
const handleSearch = () => {
  pagination.current = 1
  fetchProviderList()
}

// 处理重置
const handleReset = () => {
  Object.assign(searchForm, {
    keyword: '',
    type: '',
    status: ''
  })
  handleSearch()
}

// 处理添加
const handleAdd = () => {
  currentProvider.value = null
  providerModalVisible.value = true
}

// 处理编辑
const handleEdit = (record: any) => {
  currentProvider.value = { ...record }
  providerModalVisible.value = true
}

// 处理测试连接
const handleTest = async (record: any) => {
  try {
    ElMessage.info('正在测试连接...')
    await dnsProviderApi.test(record.id)
    ElMessage.success('连接测试成功')
    fetchProviderList()
  } catch (error) {
    console.error('连接测试失败:', error)
    ElMessage.error('连接测试失败')
  }
}

// 处理删除
const handleDelete = async (record: any) => {
  try {
    await ElMessageBox.confirm('确定要删除此提供商吗？', '确认删除', {
      type: 'warning'
    })
    await dnsProviderApi.delete(record.id)
    ElMessage.success('删除成功')
    fetchProviderList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除提供商失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 处理批量删除
const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(`确定要删除选中的${selectedRowKeys.value.length}个提供商吗？`, '确认删除', {
      type: 'warning'
    })
    await dnsProviderApi.batchDelete(selectedRowKeys.value)
    ElMessage.success('批量删除成功')
    selectedRowKeys.value = []
    fetchProviderList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 处理测试所有连接
const handleTestAll = async () => {
  try {
    ElMessage.info('正在测试所有提供商连接...')
    const ids = providerList.value.map(item => item.id)
    const response = await dnsProviderApi.batchTest(ids)
    const result = (response as any)?.data || response
    ElMessage.success(`测试完成，成功 ${result.success}，失败 ${result.failed}`)
    fetchProviderList()
  } catch (error) {
    ElMessage.error('连接测试失败')
  }
}

// 处理启用/禁用
const handleToggleStatus = async (record: any) => {
  try {
    const newStatus = record.status === 'active' ? 'inactive' : 'active'
    await dnsProviderApi.update(record.id, { status: newStatus } as any)
    ElMessage.success(`${record.status === 'active' ? '禁用' : '启用'}成功`)
    fetchProviderList()
  } catch (error) {
    console.error('切换状态失败:', error)
    ElMessage.error('操作失败')
  }
}

// 处理克隆
const handleClone = (record: any) => {
  currentProvider.value = {
    ...record,
    id: null,
    name: `${record.name} - 副本`
  }
  providerModalVisible.value = true
}

// 查看域名
const viewDomains = (record: any) => {
  // TODO: 跳转到域名列表页面，筛选此提供商的域名
  ElMessage.info('跳转到域名列表页面')
}

// 同步单个提供商的域名
const handleSyncDomains = async (record: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要同步提供商 "${record.name}" 的所有域名吗？这可能需要一些时间。`,
      '确认同步',
      {
        type: 'info',
        confirmButtonText: '开始同步',
        cancelButtonText: '取消'
      }
    )

    ElMessage.info('正在同步域名，请稍候...')
    const response = await dnsProviderApi.syncDomains(record.id)
    
    // 显示同步结果
    const result = (response as any)?.data || response
    const syncedCount = (result as any)?.synced_count || 0
    const errorCount = (result as any)?.error_count || 0
    
    if (errorCount > 0) {
      ElMessage.warning(`域名同步完成，成功: ${syncedCount}，失败: ${errorCount}`)
    } else {
      ElMessage.success(`域名同步成功，共同步 ${syncedCount} 个域名`)
    }
    
    fetchProviderList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('同步域名失败:', error)
      ElMessage.error('同步域名失败')
    }
  }
}

// 一键同步所有提供商的域名
const handleSyncAllDomains = async () => {
  try {
    if (enabledProviders.value.length === 0) {
      ElMessage.warning('没有已启用的DNS提供商，请先添加并启用提供商')
      return
    }

    await ElMessageBox.confirm(
      `确定要同步所有 ${enabledProviders.value.length} 个已启用的DNS提供商域名吗？这可能需要较长时间。`,
      '确认同步',
      {
        type: 'warning',
        confirmButtonText: '开始同步',
        cancelButtonText: '取消'
      }
    )

    // 显示同步进度对话框（暂时注释）
    // syncProgressVisible.value = true
    
    // 临时使用简单的同步方式
    syncLoading.value = true
    ElMessage.info('正在同步所有域名，请稍候...')
    
    try {
      const response = await dnsProviderApi.syncAllDomains()
      
      // 显示同步结果
      const result = (response as any)?.data || response
      const totalProviders = (result as any)?.total_providers || 0
      const totalSynced = (result as any)?.total_synced || 0
      const totalErrors = (result as any)?.total_errors || 0
      
      if (totalErrors > 0) {
        ElMessage.warning(`同步完成，共处理 ${totalProviders} 个提供商，成功同步 ${totalSynced} 个域名，失败 ${totalErrors} 个`)
      } else {
        ElMessage.success(`同步成功，共从 ${totalProviders} 个提供商同步了 ${totalSynced} 个域名`)
      }
      
      fetchProviderList()
    } catch (error) {
      console.error('同步所有域名失败:', error)
      ElMessage.error('同步所有域名失败')
    } finally {
      syncLoading.value = false
    }
    
  } catch (error) {
    if (error !== 'cancel') {
      console.error('同步所有域名失败:', error)
      ElMessage.error('同步所有域名失败')
    }
  }
}

// 处理同步取消（暂时注释）
// const handleSyncCancel = () => {
//   syncProgressVisible.value = false
//   ElMessage.info('已取消域名同步')
// }

// 处理分页
const handlePageChange = (page: number) => {
  pagination.current = page
  fetchProviderList()
}

const handlePageSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.current = 1
  fetchProviderList()
}

// 处理提供商成功
const handleProviderSuccess = () => {
  fetchProviderList()
}

// 时间格式化函数 - 优化显示完整时间
const formatTime = (time: string) => {
  if (!time) return '-'
  try {
    const date = new Date(time)
    if (isNaN(date.getTime())) return '-'
    
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    })
  } catch (error) {
    console.error('时间格式化错误:', error)
    return '-'
  }
}

const formatFullTime = (time: string) => {
  if (!time) return '未知时间'
  try {
    const date = new Date(time)
    if (isNaN(date.getTime())) return '未知时间'
    
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    })
  } catch (error) {
    console.error('完整时间格式化错误:', error)
    return '未知时间'
  }
}

// 获取提供商图标
const getProviderIcon = (type: string) => {
  const iconMap: Record<string, string> = {
    aliyun: 'Cloudy',
    tencent: 'Cloudy',
    route53: 'Cloudy',
    cloudflare: 'Cloudy',
    dnspod: 'Cloudy',
    godaddy: 'Cloudy'
  }
  return iconMap[type] || 'Setting'
}

// 获取提供商类型名称
const getProviderTypeName = (type: string) => {
  const nameMap: Record<string, string> = {
    aliyun: '阿里云',
    tencent: '腾讯云',
    route53: 'AWS Route53',
    cloudflare: 'Cloudflare',
    dnspod: 'DNSPod',
    godaddy: 'GoDaddy'
  }
  return nameMap[type] || type
}

// 获取提供商标签类型
const getProviderTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    aliyun: 'primary',
    tencent: 'success',
    route53: 'warning',
    cloudflare: 'info',
    dnspod: 'danger',
    godaddy: ''
  }
  return typeMap[type] || ''
}

// 获取状态类型
const getStatusType = (status: string) => {
  return status === 'active' ? 'success' : 'danger'
}

// 获取状态图标
const getStatusIcon = (status: string) => {
  return status === 'active' ? 'CircleCheck' : 'CircleClose'
}

// 获取状态文本
const getStatusText = (status: string) => {
  return status === 'active' ? '启用' : '禁用'
}

// 工具函数
const getProviderLogo = (type: string) => {
  const logoMap: Record<string, string> = {
    'aliyun': '/src/assets/img/providers/aliyun.png',
    'cloudflare': '/src/assets/img/providers/cloudflare.png',
    'tencent': '/src/assets/img/providers/tencent.png',
    'dnspod': '/src/assets/img/providers/dnspod.png',
    'godaddy': '/src/assets/img/providers/godaddy.png',
    'route53': '/src/assets/img/providers/aws.png'
  }
  return logoMap[type] || '/src/assets/img/providers/default.png'
}

const formatDateTime = (dateTime: string) => {
  if (!dateTime) return '-'
  const date = new Date(dateTime)
  if (isNaN(date.getTime())) return '-'
  return date.toLocaleString()
}

onMounted(() => {
  fetchProviderList()
})
</script>

<style scoped>
.dns-provider-container {
  padding: 24px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  min-height: calc(100vh - 64px);
}

.page-header {
  margin-bottom: 24px;
  background: white;
  padding: 24px;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
}

.page-header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.page-title h1 {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-title p {
  margin: 0;
  color: #64748b;
  font-size: 15px;
}

.stats-container {
  margin-bottom: 24px;
}

.stats-card {
  border-radius: 12px;
  border: 1px solid var(--el-border-color-lighter);
  transition: all 0.3s ease;
  overflow: hidden;
}

.stats-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.08);
}

.stats-content {
  display: flex;
  align-items: center;
  padding: 8px 0;
}

.stats-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  color: white;
}

.stats-card-primary .stats-icon {
  background: linear-gradient(135deg, #409eff, #66b1ff);
}

.stats-card-success .stats-icon {
  background: linear-gradient(135deg, #67c23a, #85ce61);
}

.stats-card-danger .stats-icon {
  background: linear-gradient(135deg, #f56c6c, #f78989);
}

.stats-card-warning .stats-icon {
  background: linear-gradient(135deg, #e6a23c, #ebb563);
}

.stats-info {
  flex: 1;
}

.stats-value {
  font-size: 28px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  line-height: 1;
  margin-bottom: 4px;
}

.stats-label {
  font-size: 14px;
  color: var(--el-text-color-regular);
  font-weight: 500;
}

.search-card {
  margin-bottom: 24px;
  border-radius: 12px;
  border: 1px solid var(--el-border-color-lighter);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.search-form {
  margin: 0;
}

.table-card {
  margin-bottom: 20px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.table-actions {
  display: flex;
  gap: 8px;
}

.provider-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.provider-avatar {
  flex-shrink: 0;
}

.provider-logo {
  background: linear-gradient(135deg, var(--el-color-primary), var(--el-color-primary-light-3));
  color: white;
  border: none;
}

.provider-details {
  flex: 1;
}

.provider-name {
  font-weight: 600;
  font-size: 15px;
  color: var(--el-text-color-primary);
  margin-bottom: 6px;
  line-height: 1.2;
}

.provider-type {
  margin-top: 4px;
}

.status-wrapper {
  display: flex;
  align-items: center;
}

.status-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  border-radius: 16px;
  padding: 4px 12px;
}

.status-icon {
  font-size: 12px;
}

.connection-status {
  display: flex;
  align-items: center;
}

.connection-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  border-radius: 16px;
  padding: 4px 12px;
}

.connection-icon {
  font-size: 12px;
}

.action-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.action-buttons .el-button {
  border-radius: 6px;
  font-size: 12px;
  padding: 6px 12px;
  transition: all 0.3s ease;
}

.action-buttons .el-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.provider-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.provider-table .el-table__header {
  border-radius: 12px 12px 0 0;
}

.provider-table .el-table__body tr:hover {
  background-color: #f8f9ff;
}

.provider-table .el-table__row {
  transition: background-color 0.3s ease;
}

.provider-table .el-table__cell {
  border-bottom: 1px solid #f0f0f0;
  padding: 16px 12px;
}

.table-card {
  border-radius: 12px;
  border: 1px solid var(--el-border-color-lighter);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 0;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.table-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.table-subtitle {
  font-size: 14px;
  color: var(--el-text-color-regular);
  margin: 0;
}

.table-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.action-btn {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.primary-btn {
  background: linear-gradient(135deg, #409eff, #66b1ff);
  border: none;
}

.sync-btn {
  background: linear-gradient(135deg, #67c23a, #85ce61);
  color: white;
  border: none;
}

.danger-btn:not(:disabled) {
  background: linear-gradient(135deg, #f56c6c, #f78989);
  border: none;
}

.connection-status .no-test {
  color: var(--el-text-color-placeholder);
  font-size: 12px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
