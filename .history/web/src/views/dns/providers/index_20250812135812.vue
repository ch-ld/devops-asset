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

        <el-table-column label="域名数量" width="120" align="center">
          <template #default="{ row }">
            <div class="domain-count-chart">
              <div class="count-display">
                <div class="count-number">{{ row.domain_count || 0 }}</div>
                <div class="count-label">个域名</div>
              </div>
              <div class="mini-chart">
                <div 
                  class="chart-bar" 
                  :style="{ 
                    width: Math.min(100, (row.domain_count || 0) * 10) + '%',
                    backgroundColor: getChartColor(row.domain_count || 0)
                  }"
                ></div>
              </div>
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

        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons-wrapper">
              <!-- 主要操作按钮 -->
              <div class="primary-actions">
                <el-tooltip content="测试连接" placement="top">
                  <el-button 
                    :type="row.status === 'active' ? 'primary' : 'info'"
                    size="small" 
                    @click="handleTest(row)"
                    class="action-btn test-connection-btn"
                    :disabled="row.status !== 'active'"
                    circle
                  >
                    <el-icon><Connection /></el-icon>
                  </el-button>
                </el-tooltip>
                
                <el-tooltip content="编辑配置" placement="top">
                  <el-button 
                    type="success"
                    size="small" 
                    @click="handleEdit(row)"
                    class="action-btn edit-config-btn"
                    circle
                  >
                    <el-icon><Edit /></el-icon>
                  </el-button>
                </el-tooltip>
                
                <el-tooltip content="同步域名" placement="top">
                  <el-button 
                    type="warning"
                    size="small" 
                    @click="handleSyncDomains(row)"
                    class="action-btn sync-domains-btn"
                    :disabled="row.status !== 'active'"
                    circle
                  >
                    <el-icon><Refresh /></el-icon>
                  </el-button>
                </el-tooltip>
              </div>
              
              <!-- 更多操作下拉菜单 -->
              <el-dropdown 
                trigger="click" 
                @command="(command: string) => handleCommand(command, row)"
                class="more-actions-dropdown"
              >
                <el-button size="small" class="more-actions-btn" circle>
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu class="modern-dropdown">
                    <el-dropdown-item command="toggle" class="dropdown-item">
                      <div class="dropdown-item-content">
                        <el-icon class="dropdown-item-icon" :color="row.status === 'active' ? '#ff4d4f' : '#52c41a'">
                          <component :is="row.status === 'active' ? 'VideoPause' : 'VideoPlay'" />
                        </el-icon>
                        <span class="dropdown-item-text">{{ row.status === 'active' ? '禁用提供商' : '启用提供商' }}</span>
                      </div>
                    </el-dropdown-item>
                    <el-dropdown-item command="clone" class="dropdown-item">
                      <div class="dropdown-item-content">
                        <el-icon class="dropdown-item-icon" color="#1890ff">
                          <CopyDocument />
                        </el-icon>
                        <span class="dropdown-item-text">克隆配置</span>
                      </div>
                    </el-dropdown-item>
                    <el-dropdown-item command="delete" divided class="dropdown-item danger-item">
                      <div class="dropdown-item-content">
                        <el-icon class="dropdown-item-icon" color="#ff4d4f">
                          <Delete />
                        </el-icon>
                        <span class="dropdown-item-text">删除提供商</span>
                      </div>
                    </el-dropdown-item>
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
  Cloudy,
  CopyDocument,
  MoreFilled,
  VideoPause,
  VideoPlay
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

// 处理下拉菜单命令
const handleCommand = (command: string, record: any) => {
  switch (command) {
    case 'sync':
      handleSyncDomains(record)
      break
    case 'toggle':
      handleToggleStatus(record)
      break
    case 'clone':
      handleClone(record)
      break
    case 'delete':
      handleDelete(record)
      break
    default:
      console.warn('未知的命令:', command)
  }
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
    let date: Date
    
    // 处理不同的时间格式
    if (typeof time === 'string') {
      // 如果是 ISO 字符串或标准时间字符串
      if (time.includes('T') || time.includes('-')) {
        date = new Date(time)
      } else {
        // 如果是时间戳（毫秒或秒）
        const timestamp = parseInt(time)
        if (timestamp < 9999999999) {
          // 秒级时间戳，转换为毫秒
          date = new Date(timestamp * 1000)
        } else {
          // 毫秒级时间戳
          date = new Date(timestamp)
        }
      }
    } else {
      date = new Date(time)
    }
    
    // 检查时间是否有效且不是1970年
    if (isNaN(date.getTime()) || date.getFullYear() < 2020) {
      console.warn('无效的时间数据:', time)
      return '时间数据异常'
    }
    
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
    console.error('时间格式化错误:', error, '原始数据:', time)
    return '时间格式错误'
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

// 获取图表颜色
const getChartColor = (count: number) => {
  if (count === 0) return '#d9d9d9'
  if (count <= 2) return '#52c41a'
  if (count <= 5) return '#1890ff'
  if (count <= 10) return '#fa8c16'
  return '#f5222d'
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
  background: linear-gradient(135deg, #f0f8ff 0%, #e6f3ff 50%, #f8f9fa 100%);
  min-height: calc(100vh - 64px);
  position: relative;
}

.dns-provider-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 200px;
  background: linear-gradient(135deg, rgba(24, 144, 255, 0.05) 0%, rgba(64, 169, 255, 0.03) 100%);
  pointer-events: none;
  z-index: 0;
}

.dns-provider-container > * {
  position: relative;
  z-index: 1;
}

.page-header {
  margin-bottom: 24px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(255, 255, 255, 0.9) 100%);
  padding: 32px;
  border-radius: 20px;
  box-shadow: 
    0 8px 32px rgba(24, 144, 255, 0.12),
    0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(24, 144, 255, 0.1);
  backdrop-filter: blur(20px);
  position: relative;
  overflow: hidden;
}

.page-header::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(24, 144, 255, 0.1) 0%, transparent 70%);
  pointer-events: none;
}

.page-header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.page-title h1 {
  margin: 0 0 12px 0;
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #1890ff 0%, #722ed1 50%, #eb2f96 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  position: relative;
  z-index: 1;
  letter-spacing: -0.5px;
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
  border-radius: 16px;
  border: 1px solid rgba(24, 144, 255, 0.1);
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.8) 0%, rgba(255, 255, 255, 0.95) 100%);
  backdrop-filter: blur(10px);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  position: relative;
}

.stats-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 100%;
  background: linear-gradient(135deg, transparent 0%, rgba(24, 144, 255, 0.03) 100%);
  pointer-events: none;
}

.stats-card:hover {
  transform: translateY(-4px) scale(1.02);
  box-shadow: 
    0 12px 40px rgba(24, 144, 255, 0.15),
    0 4px 16px rgba(0, 0, 0, 0.08);
  border-color: rgba(24, 144, 255, 0.2);
}

.stats-content {
  display: flex;
  align-items: center;
  padding: 8px 0;
}

.stats-icon {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
  color: white;
  position: relative;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.stats-card-primary .stats-icon {
  background: linear-gradient(135deg, #1890ff, #40a9ff, #69c0ff);
  box-shadow: 0 4px 20px rgba(24, 144, 255, 0.3);
}

.stats-card-success .stats-icon {
  background: linear-gradient(135deg, #52c41a, #73d13d, #95de64);
  box-shadow: 0 4px 20px rgba(82, 196, 26, 0.3);
}

.stats-card-danger .stats-icon {
  background: linear-gradient(135deg, #ff4d4f, #ff7875, #ffa39e);
  box-shadow: 0 4px 20px rgba(255, 77, 79, 0.3);
}

.stats-card-warning .stats-icon {
  background: linear-gradient(135deg, #fa8c16, #ffa940, #ffc069);
  box-shadow: 0 4px 20px rgba(250, 140, 22, 0.3);
}

.stats-info {
  flex: 1;
}

.stats-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  line-height: 1;
  margin-bottom: 6px;
  background: linear-gradient(135deg, #1890ff, #722ed1);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
}

.stats-label {
  font-size: 15px;
  color: var(--el-text-color-regular);
  font-weight: 600;
  letter-spacing: 0.5px;
}

.search-card {
  margin-bottom: 24px;
  border-radius: 16px;
  border: 1px solid rgba(24, 144, 255, 0.1);
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.8) 0%, rgba(255, 255, 255, 0.95) 100%);
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 20px rgba(24, 144, 255, 0.08);
}

.search-form {
  margin: 0;
}

.table-card {
  margin-bottom: 24px;
  border-radius: 16px;
  border: 1px solid rgba(24, 144, 255, 0.1);
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.9) 0%, rgba(255, 255, 255, 0.95) 100%);
  backdrop-filter: blur(10px);
  box-shadow: 0 6px 30px rgba(24, 144, 255, 0.12);
  overflow: hidden;
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

/* 新的操作栏样式 */
.action-buttons-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: center;
}

.primary-actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.action-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.action-btn:hover:not(:disabled) {
  transform: translateY(-2px) scale(1.05);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
}

.action-btn:active {
  transform: translateY(0) scale(0.95);
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none !important;
  box-shadow: none !important;
}

.test-connection-btn {
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  border: none;
  color: white;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.3);
}

.test-connection-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #096dd9, #1890ff);
  box-shadow: 0 6px 20px rgba(24, 144, 255, 0.4);
}

.edit-config-btn {
  background: linear-gradient(135deg, #52c41a, #73d13d);
  border: none;
  color: white;
  box-shadow: 0 2px 8px rgba(82, 196, 26, 0.3);
}

.edit-config-btn:hover {
  background: linear-gradient(135deg, #389e0d, #52c41a);
  box-shadow: 0 6px 20px rgba(82, 196, 26, 0.4);
}

.sync-domains-btn {
  background: linear-gradient(135deg, #fa8c16, #ffa940);
  border: none;
  color: white;
  box-shadow: 0 2px 8px rgba(250, 140, 22, 0.3);
}

.sync-domains-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #d46b08, #fa8c16);
  box-shadow: 0 6px 20px rgba(250, 140, 22, 0.4);
}

.more-actions-dropdown {
  margin-left: 4px;
}

.more-actions-btn {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #8c8c8c, #bfbfbf);
  border: none;
  color: white;
  border-radius: 50%;
  box-shadow: 0 2px 8px rgba(140, 140, 140, 0.3);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.more-actions-btn:hover {
  background: linear-gradient(135deg, #595959, #8c8c8c);
  transform: translateY(-2px) scale(1.05);
  box-shadow: 0 6px 20px rgba(140, 140, 140, 0.4);
}

.provider-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: none;
  background: transparent;
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

/* 现代化下拉菜单样式 */
.modern-dropdown {
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  border: 1px solid rgba(0, 0, 0, 0.06);
  padding: 8px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
}

.dropdown-item {
  border-radius: 8px;
  padding: 0 !important;
  margin-bottom: 4px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.dropdown-item:last-child {
  margin-bottom: 0;
}

.dropdown-item-content {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  transition: all 0.3s ease;
}

.dropdown-item:hover .dropdown-item-content {
  background: linear-gradient(135deg, #f0f9ff, #e6f7ff);
  transform: translateX(4px);
}

.dropdown-item.danger-item:hover .dropdown-item-content {
  background: linear-gradient(135deg, #fff2f0, #ffece8);
}

.dropdown-item-icon {
  font-size: 16px;
  flex-shrink: 0;
  transition: all 0.3s ease;
}

.dropdown-item:hover .dropdown-item-icon {
  transform: scale(1.1);
}

.dropdown-item-text {
  font-size: 14px;
  font-weight: 500;
  color: var(--el-text-color-primary);
  transition: all 0.3s ease;
}

.dropdown-item:hover .dropdown-item-text {
  color: var(--el-color-primary);
}

.dropdown-item.danger-item:hover .dropdown-item-text {
  color: var(--el-color-danger);
}

/* 时间信息样式优化 */
.time-info {
  display: flex;
  align-items: center;
}

.created-time {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--el-text-color-regular);
  font-size: 13px;
}

.time-icon {
  color: var(--el-text-color-placeholder);
  flex-shrink: 0;
}

.time-text {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
  letter-spacing: 0.5px;
}

/* 域名计数图表样式 */
.domain-count-chart {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 8px;
}

.count-display {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.count-number {
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(135deg, #1890ff, #722ed1);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1;
}

.count-label {
  font-size: 11px;
  color: var(--el-text-color-placeholder);
  font-weight: 500;
}

.mini-chart {
  width: 60px;
  height: 4px;
  background: #f0f0f0;
  border-radius: 2px;
  overflow: hidden;
  position: relative;
}

.chart-bar {
  height: 100%;
  border-radius: 2px;
  transition: all 0.3s ease;
  min-width: 2px;
  position: relative;
}

.chart-bar::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, transparent 0%, rgba(255,255,255,0.3) 50%, transparent 100%);
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
