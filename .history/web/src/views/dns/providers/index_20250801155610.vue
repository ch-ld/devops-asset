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
      <el-row :gutter="16">
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="提供商总数"
              :value="statistics.total"
              :value-style="{ color: '#409eff' }"
            >
              <template #suffix>
                <el-icon><Setting /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="已启用"
              :value="statistics.enabled"
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
              title="已禁用"
              :value="statistics.disabled"
              :value-style="{ color: '#f56c6c' }"
            >
              <template #suffix>
                <el-icon><CircleClose /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="连接异常"
              :value="statistics.error"
              :value-style="{ color: '#e6a23c' }"
            >
              <template #suffix>
                <el-icon><Warning /></el-icon>
              </template>
            </el-statistic>
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
          <span>提供商列表</span>
          <div class="table-actions">
            <el-button @click="handleTestAll">
              <el-icon><Connection /></el-icon>
              测试连接
            </el-button>
            <el-button 
              type="danger" 
              @click="handleBatchDelete" 
              :disabled="selectedRowKeys.length === 0"
            >
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="providerList"
        @selection-change="handleSelectionChange"
        stripe
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="name" label="提供商名称" width="200">
          <template #default="{ row }">
            <div class="provider-info">
              <el-avatar :src="getProviderLogo(row.type)" :size="24" class="provider-logo" />
              <div class="provider-details">
                <div class="provider-name">{{ row.name }}</div>
                <div class="provider-type">{{ getProviderTypeName(row.type) }}</div>
              </div>
            </div>
          </template>
        </el-table-column>

                 <el-table-column label="状态" width="100">
           <template #default="{ row }">
             <el-tag :type="getStatusType(row.status) as any">
               {{ getStatusText(row.status) }}
             </el-tag>
           </template>
         </el-table-column>

        <el-table-column label="连接状态" width="120">
          <template #default="{ row }">
            <div class="connection-status">
              <el-tag 
                v-if="row.lastTestResult" 
                :type="row.lastTestResult.success ? 'success' : 'danger'"
                size="small"
              >
                {{ row.lastTestResult.success ? '连接正常' : '连接失败' }}
              </el-tag>
              <span v-else class="no-test">未测试</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="description" label="描述" show-overflow-tooltip />

        <el-table-column label="域名数量" width="100">
          <template #default="{ row }">
            <el-link @click="viewDomains(row)" type="primary">{{ row.domainCount || 0 }}</el-link>
          </template>
        </el-table-column>

        <el-table-column prop="createdAt" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>

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
  Warning
} from '@element-plus/icons-vue'
import ProviderModal from './components/ProviderModal.vue'
import { dnsProviderApi } from '@/api/dns/provider'

// 响应式数据
const loading = ref(false)
const syncLoading = ref(false)
const providerList = ref<any[]>([])
const selectedRowKeys = ref<number[]>([])
const providerModalVisible = ref(false)
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
    providerList.value = response.items || []
    pagination.total = response.total || 0
    updateStatistics()
  } catch (error) {
    console.error('获取提供商列表失败:', error)
    ElMessage.error('获取提供商列表失败')
  } finally {
    loading.value = false
  }
}

// 更新统计数据
const updateStatistics = () => {
  statistics.total = providerList.value.length
  statistics.enabled = providerList.value.filter(item => item.status === 'enabled').length
  statistics.disabled = providerList.value.filter(item => item.status === 'disabled').length
  statistics.error = providerList.value.filter(item => 
    item.lastTestResult && !item.lastTestResult.success
  ).length
}

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
    // TODO: 调用API批量删除
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
    // TODO: 调用API测试所有连接
    await new Promise(resolve => setTimeout(resolve, 3000))
    ElMessage.success('连接测试完成')
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
    await ElMessageBox.confirm(
      '确定要同步所有DNS提供商的域名吗？这可能需要较长时间。',
      '确认同步',
      {
        type: 'warning',
        confirmButtonText: '开始同步',
        cancelButtonText: '取消'
      }
    )

    syncLoading.value = true
    ElMessage.info('正在同步所有域名，请稍候...')
    await dnsProviderApi.syncAllDomains()
    ElMessage.success('所有域名同步成功')
    fetchProviderList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('同步所有域名失败:', error)
      ElMessage.error('同步所有域名失败')
    }
  } finally {
    syncLoading.value = false
  }
}

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

const getProviderTypeName = (type: string) => {
  const nameMap: Record<string, string> = {
    'aliyun': '阿里云',
    'cloudflare': 'Cloudflare',
    'tencent': '腾讯云',
    'dnspod': 'DNSPod',
    'godaddy': 'GoDaddy',
    'route53': 'AWS Route53'
  }
  return nameMap[type] || type
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    'enabled': 'success',
    'disabled': 'info'
  }
  return statusMap[status] || 'info'
}

const getStatusIcon = (status: string) => {
  // Element Plus中的图标组件需要单独处理
  return undefined
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    'enabled': '已启用',
    'disabled': '已禁用'
  }
  return statusMap[status] || '未知'
}

const formatDateTime = (dateTime: string) => {
  return new Date(dateTime).toLocaleString()
}

onMounted(() => {
  fetchProviderList()
})
</script>

<style scoped>
.dns-provider-container {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.page-title h1 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
}

.page-title p {
  margin: 0;
  color: var(--el-text-color-regular);
}

.stats-container {
  margin-bottom: 20px;
}

.search-card {
  margin-bottom: 20px;
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

.provider-logo {
  flex-shrink: 0;
}

.provider-details {
  flex: 1;
}

.provider-name {
  font-weight: 500;
  margin-bottom: 2px;
}

.provider-type {
  font-size: 12px;
  color: var(--el-text-color-regular);
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
