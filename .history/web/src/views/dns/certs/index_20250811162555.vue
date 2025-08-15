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
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            申请证书
          </el-button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-container">
      <el-row :gutter="16">
        <el-col :span="6">
          <DNSStatCard
            title="总证书数"
            :value="statistics.total"
            icon="lock"
            variant="primary"
            :loading="statisticsLoading"
          />
        </el-col>
        <el-col :span="6">
          <DNSStatCard
            title="有效证书"
            :value="statistics.valid"
            icon="check"
            variant="success"
            :loading="statisticsLoading"
          />
        </el-col>
        <el-col :span="6">
          <DNSStatCard
            title="即将过期"
            :value="statistics.expiring"
            icon="warning"
            variant="warning"
            :loading="statisticsLoading"
          />
        </el-col>
        <el-col :span="6">
          <DNSStatCard
            title="已过期"
            :value="statistics.expired"
            icon="warning"
            variant="danger"
            :loading="statisticsLoading"
          />
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
              <el-option label="申请中" value="pending" />
              <el-option label="已签发" value="issued" />
              <el-option label="已过期" value="expired" />
              <el-option label="已吊销" value="revoked" />
            </el-select>
          </el-form-item>
          <el-form-item label="CA类型" prop="ca_type">
            <el-select
              v-model="searchForm.ca_type"
              placeholder="请选择CA类型"
              clearable
              style="width: 150px"
            >
              <el-option label="全部" value="" />
              <el-option label="Let's Encrypt" value="letsencrypt" />
              <el-option label="ZeroSSL" value="zerossl" />
              <el-option label="自定义" value="custom" />
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

    <!-- 证书列表 -->
    <div class="table-container">
      <el-card shadow="never">
        <template #header>
          <div class="card-header">
            <span>证书列表</span>
            <div>
              <el-button @click="handleRefresh">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
              <el-button @click="handleBatchRenew" :disabled="!hasSelected">
                <el-icon><Refresh /></el-icon>
                批量续期
              </el-button>
            </div>
          </div>
        </template>

        <el-table
          ref="tableRef"
          :data="certificates"
          :loading="loading"
          @selection-change="handleSelectionChange"
          row-key="id"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="domain_name" label="域名" min-width="200">
            <template #default="{ row }">
              <div class="domain-cell">
                <div class="domain-name">{{ row.domain_name }}</div>
                <div class="cert-type" v-if="row.is_wildcard">
                  <el-tag type="success" size="small">泛域名</el-tag>
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
          <el-table-column prop="ca_type" label="CA类型" width="120">
            <template #default="{ row }">
              <el-tag type="info" size="small">
                {{ getCATypeName(row.ca_type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="issued_at" label="签发时间" width="140">
            <template #default="{ row }">
              {{ formatDate(row.issued_at) }}
            </template>
          </el-table-column>
          <el-table-column prop="expires_at" label="过期时间" width="180">
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
          <el-table-column prop="auto_renew" label="自动续期" width="100">
            <template #default="{ row }">
              <el-switch
                v-model="row.auto_renew"
                @change="handleAutoRenewChange(row)"
              />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="250" fixed="right">
            <template #default="{ row }">
              <el-button
                type="primary"
                size="small"
                text
                @click="handleView(row)"
              >
                查看
              </el-button>
              <el-button
                type="success"
                size="small"
                text
                @click="handleRenew(row)"
                :disabled="row.status !== 'issued'"
              >
                续期
              </el-button>
              <el-button
                type="warning"
                size="small"
                text
                @click="handleDownload(row)"
                :disabled="row.status !== 'issued'"
              >
                下载
              </el-button>
              <el-dropdown @command="(command) => handleCommand(command, row)">
                <el-button type="primary" size="small" text>
                  更多
                  <el-icon><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="deploy" :disabled="row.status !== 'issued'">
                      部署
                    </el-dropdown-item>
                    <el-dropdown-item command="revoke" :disabled="row.status !== 'issued'">
                      吊销
                    </el-dropdown-item>
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

    <!-- 证书申请弹窗 -->
    <CertificateModal
      v-model:visible="modalVisible"
      :certificate="currentCertificate"
      @success="handleModalSuccess"
    />

    <!-- 证书详情抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      :title="`证书详情 - ${currentCertificate?.domain_name}`"
      size="50%"
    >
      <CertificateDetail
        v-if="currentCertificate"
        :certificate="currentCertificate"
        @refresh="fetchData"
      />
    </el-drawer>
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
  Lock,
  CircleCheck, 
  Warning, 
  CircleClose
} from '@element-plus/icons-vue'
import CertificateModal from './components/CertificateModal.vue'
import CertificateDetail from './components/CertificateDetail.vue'
import { certificateApi } from '@/api/dns/certificate'
import type { Certificate } from '@/types/dns'

// 响应式数据
const loading = ref(false)
const modalVisible = ref(false)
const drawerVisible = ref(false)
const currentCertificate = ref<Certificate | null>(null)
const certificates = ref<Certificate[]>([])
const statistics = ref({
  total: 0,
  valid: 0,
  expiring: 0,
  expired: 0
})

// 搜索表单
const searchFormRef = ref()
const searchForm = reactive({
  keyword: '',
  status: '',
  ca_type: ''
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 选择
const selectedRows = ref<Certificate[]>([])
const hasSelected = computed(() => selectedRows.value.length > 0)

// 工具方法
const getStatusType = (status: string) => {
  const statusMap: Record<string, any> = {
    pending: 'warning',
    issued: 'success',
    expired: 'danger',
    revoked: 'info'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '申请中',
    issued: '已签发',
    expired: '已过期',
    revoked: '已吊销'
  }
  return statusMap[status] || status
}

const getCATypeName = (caType: string) => {
  const typeMap: Record<string, string> = {
    letsencrypt: "Let's Encrypt",
    zerossl: 'ZeroSSL',
    custom: '自定义'
  }
  return typeMap[caType] || caType
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
  currentCertificate.value = null
  modalVisible.value = true
}

const handleView = (row: Certificate) => {
  currentCertificate.value = row
  drawerVisible.value = true
}

const handleRenew = async (row: Certificate) => {
  try {
    await ElMessageBox.confirm(
      `确定要续期证书 "${row.domain_name}" 吗？`,
      '确认续期',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await certificateApi.renew(row.id)
    ElMessage.success('证书续期请求已提交')
    await fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('续期失败')
    }
  }
}

const handleDownload = async (row: Certificate) => {
  try {
    const response = await certificateApi.download(row.id, 'pem')
    // 创建下载链接
    const blob = new Blob([response.content], { type: response.mime_type })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = response.filename
    link.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('证书下载成功')
  } catch (error) {
    ElMessage.error('证书下载失败')
  }
}

const handleCommand = (command: string, row: Certificate) => {
  switch (command) {
    case 'deploy':
      ElMessage.info(`部署证书 ${row.domain_name}`)
      break
    case 'revoke':
      handleRevoke(row)
      break
    case 'export':
      ElMessage.info(`导出证书 ${row.domain_name} 的配置`)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const handleRevoke = async (row: Certificate) => {
  try {
    await ElMessageBox.confirm(
      `确定要吊销证书 "${row.domain_name}" 吗？此操作不可恢复。`,
      '确认吊销',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await certificateApi.revoke(row.id)
    ElMessage.success('证书吊销成功')
    await fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('吊销失败')
    }
  }
}

const handleDelete = async (row: Certificate) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除证书 "${row.domain_name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await certificateApi.delete(row.id)
    ElMessage.success('删除成功')
    await fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleBatchRenew = async () => {
  if (!hasSelected.value) return
  
  try {
    await ElMessageBox.confirm(
      `确定要续期选中的 ${selectedRows.value.length} 个证书吗？`,
      '确认批量续期',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const ids = selectedRows.value.map(row => row.id)
    await certificateApi.batchRenew(ids)
    ElMessage.success('批量续期请求已提交')
    selectedRows.value = []
    await fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量续期失败')
    }
  }
}

const handleSelectionChange = (rows: Certificate[]) => {
  selectedRows.value = rows
}

const handleAutoRenewChange = async (row: Certificate) => {
  try {
    await certificateApi.update(row.id, { auto_renew: row.auto_renew })
    ElMessage.success(`${row.auto_renew ? '启用' : '禁用'}自动续期成功`)
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
    const response = await certificateApi.list(params)
    certificates.value = response.list || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('获取证书列表失败')
  } finally {
    loading.value = false
  }
}

const fetchStatistics = async () => {
  try {
    const stats = await certificateApi.stats()
    statistics.value = stats || {
      total: 0,
      valid: 0,
      expiring: 0,
      expired: 0
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 生命周期
onMounted(() => {
  fetchData()
  fetchStatistics()
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
  
  .cert-type {
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
