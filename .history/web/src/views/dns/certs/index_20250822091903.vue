<template>
  <div class="modern-page-container">
    <!-- 现代化页面头部 -->
    <div class="modern-page-header">
      <div class="header-container">
        <div class="header-main">
          <div class="title-section">
            <div class="title-group">
              <div class="title-icon-wrapper">
                <el-icon class="title-icon"><Lock /></el-icon>
              </div>
              <div class="title-content">
                <h1 class="page-title">证书管理</h1>
                <div class="title-badge">
                  <el-tag type="info" size="small">SSL/TLS</el-tag>
                </div>
              </div>
            </div>
            <p class="page-subtitle">管理SSL/TLS证书，包括申请、续期、部署和监控</p>
          </div>

          <div class="header-actions">
            <el-button @click="handleRefresh" :icon="Refresh" circle class="refresh-btn" />
            <el-divider direction="vertical" />
            <el-dropdown @command="handleCreateCommand" class="create-dropdown">
              <el-button type="primary" class="primary-btn">
                <el-icon><Plus /></el-icon>
                申请证书
                <el-icon><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="simple">
                    <el-icon><Lightning /></el-icon>
                    快速申请
                  </el-dropdown-item>
                  <el-dropdown-item command="advanced">
                    <el-icon><Setting /></el-icon>
                    高级配置
                  </el-dropdown-item>
                  <el-dropdown-item command="upload" divided>
                    <el-icon><Upload /></el-icon>
                    导入证书
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </div>
    </div>

    <!-- 现代化统计卡片 -->
    <div class="modern-stats-section">
      <div class="stats-grid">
        <div class="stat-card total-card">
          <div class="card-content">
            <div class="stat-header">
              <div class="stat-icon-container total">
                <el-icon class="stat-icon"><Lock /></el-icon>
              </div>
              <div class="stat-trend">
                <el-icon class="trend-icon"><TrendCharts /></el-icon>
              </div>
            </div>
            <div class="stat-body">
              <div class="stat-number">{{ statistics.total || 0 }}</div>
              <div class="stat-label">总证书数</div>
              <div class="stat-description">所有SSL/TLS证书</div>
            </div>
          </div>
        </div>

        <div class="stat-card valid-card">
          <div class="card-content">
            <div class="stat-header">
              <div class="stat-icon-container valid">
                <el-icon class="stat-icon"><Check /></el-icon>
              </div>
              <div class="stat-badge success">
                <span>{{ Math.round((statistics.valid / statistics.total) * 100) || 0 }}%</span>
              </div>
            </div>
            <div class="stat-body">
              <div class="stat-number">{{ statistics.valid || 0 }}</div>
              <div class="stat-label">有效证书</div>
              <div class="stat-description">正常使用中</div>
            </div>
          </div>
        </div>

        <div class="stat-card expiring-card">
          <div class="card-content">
            <div class="stat-header">
              <div class="stat-icon-container expiring">
                <el-icon class="stat-icon"><Warning /></el-icon>
              </div>
              <div class="stat-alert" v-if="statistics.expiring > 0">
                <el-icon><Bell /></el-icon>
              </div>
            </div>
            <div class="stat-body">
              <div class="stat-number">{{ statistics.expiring || 0 }}</div>
              <div class="stat-label">即将过期</div>
              <div class="stat-description">30天内到期</div>
            </div>
          </div>
        </div>

        <div class="stat-card expired-card">
          <div class="card-content">
            <div class="stat-header">
              <div class="stat-icon-container expired">
                <el-icon class="stat-icon"><CircleClose /></el-icon>
              </div>
              <div class="stat-urgent" v-if="statistics.expired > 0">
                <el-icon><Warning /></el-icon>
              </div>
            </div>
            <div class="stat-body">
              <div class="stat-number">{{ statistics.expired || 0 }}</div>
              <div class="stat-label">已过期</div>
              <div class="stat-description">需要立即处理</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 现代化搜索和筛选 -->
    <div class="modern-search-section">
      <el-card class="search-card" shadow="never">
        <div class="search-content">
          <div class="search-left">
            <el-input
              v-model="searchForm.keyword"
              placeholder="搜索证书域名..."
              size="large"
              clearable
              class="search-input"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
          <div class="search-right">
            <el-select
              v-model="searchForm.status"
              placeholder="证书状态"
              size="large"
              clearable
              class="status-filter"
              style="width: 150px"
            >
              <el-option label="全部" value="" />
              <el-option label="申请中" value="pending" />
              <el-option label="处理中" value="processing" />
              <el-option label="验证中" value="validating" />
              <el-option label="已签发" value="issued" />
              <el-option label="已过期" value="expired" />
              <el-option label="申请失败" value="failed" />
              <el-option label="已吊销" value="revoked" />
            </el-select>
            <el-select
              v-model="searchForm.ca_type"
              placeholder="CA类型"
              size="large"
              clearable
              class="status-filter"
              style="width: 150px"
            >
              <el-option label="全部" value="" />
              <el-option label="Let's Encrypt" value="letsencrypt" />
              <el-option label="ZeroSSL" value="zerossl" />
              <el-option label="自定义" value="custom" />
            </el-select>
            <el-button type="primary" @click="handleSearch" size="large" :icon="Search">
              搜索
            </el-button>
            <el-button @click="handleReset" size="large" :icon="Refresh">
              重置
            </el-button>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 证书列表 -->
    <div class="table-container">
      <el-card shadow="never">
        <template #header>
          <div class="card-header">
            <span>证书列表</span>
            <div class="table-actions">
              <el-space>
                <el-button @click="handleRefresh">
                  <el-icon><Refresh /></el-icon>
                  刷新
                </el-button>
                <el-button type="primary" @click="handleImportCert">
                  <el-icon><Upload /></el-icon>
                  导入证书
                </el-button>
                <el-dropdown @command="handleBatchCommand" :disabled="!hasSelected">
                  <el-button type="warning" :disabled="!hasSelected">
                    批量操作 ({{ selectedRows.length }})
                    <el-icon><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="renew">
                        <el-icon><Refresh /></el-icon>
                        批量续期
                      </el-dropdown-item>
                      <el-dropdown-item command="download">
                        <el-icon><Download /></el-icon>
                        批量下载
                      </el-dropdown-item>
                      <el-dropdown-item command="deploy">
                        <el-icon><Position /></el-icon>
                        批量部署
                      </el-dropdown-item>
                      <el-dropdown-item command="export">
                        <el-icon><DocumentCopy /></el-icon>
                        导出报告
                      </el-dropdown-item>
                      <el-dropdown-item command="delete" divided>
                        <el-icon><Delete /></el-icon>
                        批量删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </el-space>
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
          <el-table-column label="域名" min-width="200">
            <template #default="{ row }">
              <div class="domain-cell">
                <div class="domain-name">{{ row.common_name || row.domain_name || '-' }}</div>
                <div class="cert-type" v-if="row.subject_alt_names && row.subject_alt_names.length > 0">
                  <el-tag
                    v-for="san in row.subject_alt_names.slice(0, 2)"
                    :key="san"
                    type="info"
                    size="small"
                    style="margin-right: 4px;"
                  >
                    {{ san }}
                  </el-tag>
                  <el-tag
                    v-if="row.subject_alt_names.length > 2"
                    type="info"
                    size="small"
                  >
                    +{{ row.subject_alt_names.length - 2 }}
                  </el-tag>
                </div>
                <div class="cert-type" v-if="row.common_name && row.common_name.startsWith('*.')">
                  <el-tag type="success" size="small">泛域名</el-tag>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="120">
            <template #default="{ row }">
              <div class="status-cell">
                <el-tag :type="getStatusType(row.status)" class="status-tag">
                  <el-icon class="status-icon">
                    <component :is="getStatusIcon(row.status)" />
                  </el-icon>
                  {{ getStatusText(row.status) }}
                </el-tag>
                <!-- 申请中状态显示进度 -->
                <div v-if="row.status === 'pending'" class="progress-info">
                  <el-tooltip :content="getProgressText(row.progress || 0)" placement="top">
                    <el-progress 
                      :percentage="row.progress || 0" 
                      :stroke-width="3" 
                      :show-text="false"
                      status="success"
                    />
                  </el-tooltip>
                </div>
                <!-- 错误状态显示错误信息 -->
                <div v-if="row.status === 'failed' && row.error_message" class="error-info">
                  <el-tooltip :content="row.error_message" placement="top">
                    <el-text type="danger" size="small" truncated>
                      {{ row.error_message }}
                    </el-text>
                  </el-tooltip>
                </div>
              </div>
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
                :model-value="row.auto_renew"
                @change="val => handleAutoRenewChange(row, Boolean(val))"
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
                :disabled="row.status !== 'issued' || renewingCertIds.has(row.id)"
                :loading="renewingCertIds.has(row.id)"
              >
                {{ renewingCertIds.has(row.id) ? '续期中...' : '续期' }}
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
    <CertificateDetail
      v-if="currentCertificate"
      v-model:visible="drawerVisible"
      :certificate="currentCertificate"
      @refresh="fetchData"
    />

    <!-- 证书下载对话框 -->
    <CertificateDownloadDialog
      v-model:visible="showDownloadDialog"
      :certificate="selectedCertificate"
    />

    <!-- 证书部署对话框 -->
    <CertificateDeployDialog
      v-model:visible="showDeployDialog"
      :certificate="selectedCertificate"
      @success="fetchData"
    />

    <!-- 证书导出配置对话框 -->
    <CertificateExportDialog
      v-model:visible="showExportDialog"
      :certificate="selectedCertificate"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  ArrowDown,
  Refresh,
  Upload,
  Download,
  Position,
  DocumentCopy,
  Delete,
  Lightning,
  Setting,
  Lock,
  Search,
  Loading,
  CircleCheck,
  CircleClose,
  Warning,
  Clock
} from '@element-plus/icons-vue'
import { DNSStatCard } from '@/components/dns'
import CertificateModal from './components/CertificateModal.vue'
import CertificateDetail from './components/CertificateDetail.vue'
import CertificateDownloadDialog from '@/components/dns/CertificateDownloadDialog.vue'
import CertificateDeployDialog from '@/components/dns/CertificateDeployDialog.vue'
import CertificateExportDialog from '@/components/dns/CertificateExportDialog.vue'
import { certificateApi } from '@/api/dns/certificate'
import type { Certificate } from '@/types/dns'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const statisticsLoading = ref(false)
const modalVisible = ref(false)
const drawerVisible = ref(false)
const showDownloadDialog = ref(false)
const showDeployDialog = ref(false)
const showExportDialog = ref(false)
const currentCertificate = ref<Certificate | null>(null)
const selectedCertificate = ref<Certificate | null>(null)
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
    processing: 'warning',
    validating: 'warning',
    issued: 'success',
    active: 'success',
    expired: 'danger',
    failed: 'danger',
    revoked: 'info'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '申请中',
    processing: '处理中',
    validating: '验证中',
    issued: '已签发',
    active: '已签发',
    expired: '已过期',
    failed: '申请失败',
    revoked: '已吊销'
  }
  return statusMap[status] || status
}

const getStatusIcon = (status: string) => {
  const iconMap: Record<string, any> = {
    pending: Clock,
    processing: Loading,
    validating: Loading,
    issued: CircleCheck,
    active: CircleCheck,
    expired: Warning,
    failed: CircleClose,
    revoked: CircleClose
  }
  return iconMap[status] || CircleCheck
}

const getProgressText = (progress: number) => {
  if (progress < 20) return '初始化申请...'
  if (progress < 40) return '验证域名所有权...'
  if (progress < 60) return '生成密钥对...'
  if (progress < 80) return '申请证书...'
  if (progress < 100) return '处理证书...'
  return '申请完成'
}

const getCATypeName = (caType: string) => {
  const typeMap: Record<string, string> = {
    letsencrypt: "Let's Encrypt",
    zerossl: 'ZeroSSL',
    custom: '自定义'
  }
  return typeMap[caType] || caType
}

const formatDate = (dateValue: string | number) => {
  if (!dateValue) return '-'

  try {
    let date: Date

    if (typeof dateValue === 'number') {
      // 处理时间戳
      if (dateValue < 10000000000) {
        // 秒级时间戳，转换为毫秒
        date = new Date(dateValue * 1000)
      } else {
        // 毫秒级时间戳
        date = new Date(dateValue)
      }
    } else {
      // 字符串格式
      date = new Date(dateValue)
    }

    // 检查日期是否有效
    if (isNaN(date.getTime())) {
      return '-'
    }

    // 检查是否是1970年（通常表示无效时间戳）
    if (date.getFullYear() === 1970) {
      return '-'
    }

    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    })
  } catch (error) {
    console.warn('时间格式化错误:', error, dateValue)
    return '-'
  }
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

  // 尝试多种时间格式解析
  let expireDate: Date

  // 如果是时间戳（数字）
  if (/^\d+$/.test(expiresAt)) {
    const timestamp = parseInt(expiresAt)
    // 判断是秒级还是毫秒级时间戳
    expireDate = new Date(timestamp < 10000000000 ? timestamp * 1000 : timestamp)
  } else {
    // 字符串格式
    expireDate = new Date(expiresAt)
  }

  // 检查日期是否有效
  if (isNaN(expireDate.getTime())) {
    console.warn('Invalid date:', expiresAt)
    return '日期格式错误'
  }

  const now = new Date()
  const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))

  return diffDays > 0 ? `${diffDays}天后过期` : `已过期${Math.abs(diffDays)}天`
}

// 事件处理
const handleAdd = () => {
  // 跳转到新的证书申请页面（快速申请）
  router.push('/dns/certs/create')
}

const handleAdvancedCreate = () => {
  // 主按钮默认跳转到高级配置页面
  router.push('/dns/certs/create-advanced')
}

const handleCreateCommand = (command: string) => {
  switch (command) {
    case 'simple':
      router.push('/dns/certs/create')
      break
    case 'advanced':
      router.push('/dns/certs/create-advanced')
      break
    case 'upload':
      router.push('/dns/certs/import')
      break
  }
}

const handleView = (row: Certificate) => {
  currentCertificate.value = row
  drawerVisible.value = true
}

// 续期状态管理
const renewingCertIds = ref<Set<number>>(new Set())

const handleRenew = async (row: Certificate) => {
  try {
    const domainName = row.common_name || row.domain_name || `证书ID: ${row.id}`
    await ElMessageBox.confirm(
      `确定要续期证书 "${domainName}" 吗？`,
      '确认续期',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    // 设置续期状态
    renewingCertIds.value.add(row.id)

    // 创建超时Promise
    const timeoutPromise = new Promise((_, reject) => {
      setTimeout(() => reject(new Error('续期操作超时')), 30000) // 30秒超时
    })

    // 执行续期操作
    const renewPromise = certificateApi.renew(row.id)

    try {
      await Promise.race([renewPromise, timeoutPromise])
      ElMessage.success('证书续期成功')
      await fetchData()
    } catch (error: any) {
      if (error.message === '续期操作超时') {
        ElMessage.warning('续期操作超时，请稍后查看证书状态')
      } else {
        throw error
      }
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('续期失败:', error)
      ElMessage.error('续期失败，请稍后重试')
    }
  } finally {
    // 清除续期状态
    renewingCertIds.value.delete(row.id)
  }
}

const handleDownload = async (row: Certificate) => {
  try {
    const domainName = row.common_name || row.domain_name || `cert_${row.id}`

    // 显示下载格式选择对话框
    showDownloadDialog.value = true
    selectedCertificate.value = row
  } catch (error: any) {
    console.error('打开下载对话框失败:', error)
    ElMessage.error('打开下载对话框失败')
  }
}

const handleCommand = (command: string, row: Certificate) => {
  const domainName = row.common_name || row.domain_name || `证书ID: ${row.id}`

  switch (command) {
    case 'deploy':
      handleDeploy(row)
      break
    case 'revoke':
      handleRevoke(row)
      break
    case 'export':
      handleExportConfig(row)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const handleDeploy = (row: Certificate) => {
  selectedCertificate.value = row
  showDeployDialog.value = true
}

const handleExportConfig = (row: Certificate) => {
  selectedCertificate.value = row
  showExportDialog.value = true
}

const handleRevoke = async (row: Certificate) => {
  try {
    const domainName = row.common_name || row.domain_name || `证书ID: ${row.id}`
    await ElMessageBox.confirm(
      `确定要吊销证书 "${domainName}" 吗？此操作不可恢复。`,
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
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('吊销证书失败:', error)
      const errorMessage = error?.response?.data?.message || error?.message || '吊销失败，请重试'
      ElMessage.error(errorMessage)
    }
  }
}

const handleDelete = async (row: Certificate) => {
  try {
    const domainName = row.common_name || row.domain_name || `证书ID: ${row.id}`
    await ElMessageBox.confirm(
      `确定要删除证书 "${domainName}" 吗？此操作不可恢复。`,
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
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除证书失败:', error)
      const errorMessage = error?.response?.data?.message || error?.message || '删除失败，请重试'
      ElMessage.error(errorMessage)
    }
  }
}

// 导入证书
const handleImportCert = () => {
  // 跳转到上传证书页面
  router.push('/dns/certs/import')
}

// 批量操作处理
const handleBatchCommand = (command: string) => {
  switch (command) {
    case 'renew':
      handleBatchRenew()
      break
    case 'download':
      handleBatchDownload()
      break
    case 'deploy':
      handleBatchDeploy()
      break
    case 'export':
      handleBatchExport()
      break
    case 'delete':
      handleBatchDelete()
      break
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

// 批量下载
const handleBatchDownload = async () => {
  if (!hasSelected.value) return

  try {
    const validCerts = selectedRows.value.filter(row => row.status === 'issued')
    if (validCerts.length === 0) {
      ElMessage.warning('请选择已签发的证书')
      return
    }

    ElMessage.info('正在准备下载，请稍候...')

    for (const cert of validCerts) {
      try {
        const response = await certificateApi.download(cert.id, 'pem')
        const url = window.URL.createObjectURL(response as Blob)
        const link = document.createElement('a')
        link.href = url
        link.download = `${cert.common_name}_${cert.id}.pem`
        link.click()
        window.URL.revokeObjectURL(url)

        // 添加延迟避免浏览器阻止多个下载
        await new Promise(resolve => setTimeout(resolve, 500))
      } catch (error) {
        console.error(`下载证书 ${cert.common_name} 失败:`, error)
      }
    }

    ElMessage.success(`已下载 ${validCerts.length} 个证书`)
    selectedRows.value = []
  } catch (error) {
    ElMessage.error('批量下载失败')
  }
}

// 批量部署
const handleBatchDeploy = async () => {
  if (!hasSelected.value) return

  const validCerts = selectedRows.value.filter(row => row.status === 'issued')
  if (validCerts.length === 0) {
    ElMessage.warning('请选择已签发的证书')
    return
  }

  try {
    // 简化的批量部署，可以后续改为打开批量部署对话框
    const hostIds = [1] // 默认主机ID，实际应用中应该让用户选择
    
    for (const cert of validCerts) {
      await certificateApi.deploy(cert.id, hostIds)
    }
    
    ElMessage.success(`已提交 ${validCerts.length} 个证书的部署任务`)
    selectedRows.value = []
    await fetchData()
  } catch (error) {
    ElMessage.error('批量部署失败')
  }
}

// 批量导出
const handleBatchExport = async () => {
  if (!hasSelected.value) return

  try {
    const ids = selectedRows.value.map(row => row.id)
    const response = await certificateApi.batchExport(ids)

    const blobPart = (response as any).content || response
    const blob = new Blob([blobPart], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `certificates_report_${new Date().toISOString().split('T')[0]}.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)

    ElMessage.success('证书报告导出成功')
    selectedRows.value = []
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

// 批量删除
const handleBatchDelete = async () => {
  if (!hasSelected.value) return

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 个证书吗？此操作不可恢复。`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const ids = selectedRows.value.map(row => row.id)
    await certificateApi.batchDelete(ids)
    ElMessage.success('批量删除成功')
    selectedRows.value = []
    await fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

const handleSelectionChange = (rows: Certificate[]) => {
  selectedRows.value = rows
}

const handleAutoRenewChange = async (row: Certificate, value: boolean) => {
  try {
    await certificateApi.update(row.id, { auto_renew: value })
    row.auto_renew = value // 更新本地状态
    ElMessage.success(`${value ? '启用' : '禁用'}自动续期成功`)
  } catch (error) {
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
    const resp: any = await certificateApi.list(params)

    console.log('证书列表API响应:', resp)

    let certificateList: any[] = []
    let total = 0

    if (resp.data?.items) {
      certificateList = resp.data.items
      total = resp.data.total || 0
    } else if (resp.data?.list) {
      certificateList = resp.data.list
      total = resp.data.total || 0
    } else if (resp.list) {
      certificateList = resp.list
      total = resp.total || 0
    } else if (Array.isArray(resp.data)) {
      certificateList = resp.data
      total = resp.data.length
    } else if (Array.isArray(resp)) {
      certificateList = resp
      total = resp.length
    }

    certificates.value = certificateList
    pagination.total = total

    // 数据获取完成后计算统计数据
    calculateStatistics()

    console.log('解析后的证书列表:', certificateList)
    console.log('总数:', total)
  } catch (error) {
    console.error('获取证书列表失败:', error)
    ElMessage.error('获取证书列表失败')
  } finally {
    loading.value = false
  }
}

const fetchStatistics = async () => {
  try {
    // 尝试从API获取统计数据
    const stats: any = await certificateApi.stats()
    if (stats && typeof stats.total === 'number') {
      statistics.value = stats as any
      return
    }
  } catch (error) {
    console.warn('API统计数据获取失败，使用本地计算:', error)
  }

  // 如果API失败，基于当前证书列表计算统计数据
  calculateStatistics()
}

// 基于证书列表计算统计数据
const calculateStatistics = () => {
  const now = new Date()
  const stats = {
    total: certificates.value.length,
    valid: 0,
    expiring: 0,
    expired: 0
  }

  certificates.value.forEach(cert => {
    if (!cert.expires_at) return

    // 解析过期时间
    let expireDate: Date
    if (typeof cert.expires_at === 'number' || /^\d+$/.test(cert.expires_at.toString())) {
      const timestamp = typeof cert.expires_at === 'number' ? cert.expires_at : parseInt(cert.expires_at.toString())
      expireDate = new Date(timestamp < 10000000000 ? timestamp * 1000 : timestamp)
    } else {
      expireDate = new Date(cert.expires_at)
    }

    if (isNaN(expireDate.getTime())) return

    const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))

    if (diffDays < 0) {
      stats.expired++
    } else if (diffDays <= 30) {
      stats.expiring++
    } else {
      stats.valid++
    }
  })

  statistics.value = stats
}

// 生命周期
onMounted(() => {
  fetchData()
  fetchStatistics()
})
</script>

<style scoped lang="scss">
/* 现代化证书管理页面样式 */
.modern-cert-container {
  min-height: 100vh;
  background: #f8fafc;
  padding: 0;
  margin: 0;

  /* 现代化页面头部 */
  .modern-page-header {
    background: #ffffff;
    border-bottom: 1px solid #e2e8f0;
    padding: 24px 0;
    margin-bottom: 24px;

    .header-container {
      max-width: 1400px;
      margin: 0 auto;
      padding: 0 24px;

      .header-main {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        gap: 32px;

        .title-section {
          flex: 1;

          .title-group {
            display: flex;
            align-items: center;
            gap: 16px;
            margin-bottom: 8px;

            .title-icon-wrapper {
              width: 48px;
              height: 48px;
              border-radius: 12px;
              background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
              display: flex;
              align-items: center;
              justify-content: center;

              .title-icon {
                font-size: 24px;
                color: #ffffff;
              }
            }

            .title-content {
              display: flex;
              align-items: center;
              gap: 12px;

              .page-title {
                font-size: 28px;
                font-weight: 700;
                color: #1e293b;
                margin: 0;
                line-height: 1.2;
              }

              .title-badge {
                :deep(.el-tag) {
                  background: var(--chip-bg-primary);
                  color: var(--chip-text-primary);
                  border: 1px solid var(--chip-border-primary);
                  font-weight: 500;
                }
              }
            }
          }

          .page-subtitle {
            font-size: 14px;
            color: #64748b;
            line-height: 1.5;
            margin: 0;
            max-width: 500px;
          }
        }

        .header-actions {
          display: flex;
          align-items: center;
          gap: 8px;

          .refresh-btn {
            width: 36px;
            height: 36px;
            border: 1px solid #e2e8f0;
            background: #ffffff;
            color: #64748b;

            &:hover {
              border-color: var(--brand-primary);
              color: var(--brand-primary);
            }
          }

          .primary-btn {
            background: var(--brand-primary);
            border: 1px solid var(--brand-primary);
            color: #ffffff;
            font-weight: 600;

            &:hover {
              background: var(--brand-primary-hover);
              border-color: var(--brand-primary-hover);
            }
          }

          :deep(.el-divider--vertical) {
            height: 20px;
            border-color: #e2e8f0;
          }
        }
      }
    }
  }

  /* 现代化统计卡片 */
  .modern-stats-section {
    max-width: 1400px;
    margin: 0 auto 24px;
    padding: 0 24px;

    .stats-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: 20px;

      .stat-card {
        background: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 16px;
        padding: 24px;
        transition: all 0.3s ease;
        position: relative;
        overflow: hidden;

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 8px 25px rgba(0, 0, 0, 0.08);
          border-color: #cbd5e1;
        }

        .card-content {
          .stat-header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: 16px;

            .stat-icon-container {
              width: 48px;
              height: 48px;
              border-radius: 12px;
              display: flex;
              align-items: center;
              justify-content: center;
              position: relative;

              .stat-icon {
                font-size: 24px;
              }

              &.total {
                background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
                color: #ffffff;
              }

              &.valid {
                background: var(--chip-bg-success);
                color: var(--chip-text-success);
              }

              &.expiring {
                background: var(--chip-bg-warning);
                color: var(--chip-text-warning);
              }

              &.expired {
                background: var(--chip-bg-danger);
                color: var(--chip-text-danger);
              }
            }

            .stat-trend {
              color: #10b981;
              font-size: 16px;
            }

            .stat-badge {
              padding: 4px 8px;
              border-radius: 6px;
              font-size: 12px;
              font-weight: 600;

              &.success {
                background: var(--chip-bg-success);
                color: var(--chip-text-success);
              }
            }

            .stat-alert {
              color: var(--chip-text-warning);
              font-size: 16px;
              animation: bounce 2s infinite;
            }

            .stat-urgent {
              color: var(--chip-text-danger);
              font-size: 16px;
              animation: bounce 2s infinite;
            }
          }

          .stat-body {
            .stat-number {
              font-size: 32px;
              font-weight: 700;
              color: #1e293b;
              line-height: 1;
              margin-bottom: 4px;
            }

            .stat-label {
              font-size: 14px;
              font-weight: 600;
              color: #475569;
              margin-bottom: 4px;
            }

            .stat-description {
              font-size: 12px;
              color: #64748b;
              line-height: 1.4;
            }
          }
        }

        &.total-card {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          color: #ffffff;
          border: none;

          .stat-number,
          .stat-label,
          .stat-description {
            color: #ffffff;
          }

          .stat-description {
            opacity: 0.8;
          }
        }
      }
    }
  }

  /* 现代化搜索区域 */
  .modern-search-section {
    max-width: 1400px;
    margin: 0 auto 24px;
    padding: 0 24px;

    .search-card {
      border: 1px solid #e2e8f0;
      border-radius: 12px;
      background: #ffffff;

      :deep(.el-card__body) {
        padding: 20px;
      }

      .search-content {
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 16px;

        .search-left {
          flex: 1;
          max-width: 400px;

          .search-input {
            :deep(.el-input__wrapper) {
              border-radius: 8px;
              border: 1px solid #e2e8f0;
              box-shadow: none;

              &:hover {
                border-color: var(--brand-primary);
              }

              &.is-focus {
                border-color: var(--brand-primary);
                box-shadow: 0 0 0 2px rgba(3, 102, 214, 0.1);
              }
            }
          }
        }

        .search-right {
          display: flex;
          align-items: center;
          gap: 12px;

          .status-filter {
            width: 140px;

            :deep(.el-select__wrapper) {
              border-radius: 8px;
              border: 1px solid #e2e8f0;

              &:hover {
                border-color: var(--brand-primary);
              }
            }
          }
        }
      }
    }
  }

  /* 动画效果 */
  @keyframes bounce {
    0%, 20%, 50%, 80%, 100% { transform: translateY(0); }
    40% { transform: translateY(-4px); }
    60% { transform: translateY(-2px); }
  }
}

.page-header {
  background: white;
  border-radius: 8px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e1e4e8;
  
  .page-header-content {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    
    .page-title {
      h1 {
        margin: 0 0 8px 0;
        font-size: 28px;
        font-weight: 700;
        color: #24292e;
        display: flex;
        align-items: center;
        gap: 12px;

        .title-icon {
          font-size: 28px;
          color: #0366d6;
        }
      }
      
      p {
        margin: 0;
        color: #64748b;
        font-size: 16px;
        font-weight: 500;
      }
    }

    .page-actions {
      .create-dropdown {
        background: #0366d6;
        border: 1px solid #0366d6;
        color: white;
        font-weight: 600;
        border-radius: 6px;
        transition: all 0.2s ease;

        &:hover {
          background: #0256cc;
          border-color: #0256cc;
        }
      }
    }
  }
}

.stats-container {
  margin-bottom: 24px;
}

.search-container {
  margin-bottom: 24px;

  .el-card {
    border-radius: 8px;
    border: 1px solid #e1e4e8;
    background: white;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }
}

.table-container {
  .el-card {
    border-radius: 8px;
    border: 1px solid #e1e4e8;
    background: white;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 700;
    font-size: 18px;
    color: #1e293b;
  }
  
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}

.domain-cell {
  .domain-name {
    font-weight: 600;
    margin-bottom: 4px;
    color: #1e293b;
    font-size: 15px;
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
    font-weight: 600;
    color: #1e293b;
  }
  
  .expires-days {
    font-size: 12px;
    font-weight: 600;
    padding: 2px 8px;
    border-radius: 12px;
    
    &.text-danger {
      background: rgba(255, 77, 79, 0.12);
      color: #ff4d4f;
      border: 1px solid rgba(255, 77, 79, 0.2);
    }

    &.text-warning {
      background: rgba(250, 140, 22, 0.12);
      color: #fa8c16;
      border: 1px solid rgba(250, 140, 22, 0.2);
    }

    &.text-success {
      background: rgba(16, 185, 129, 0.12);
      color: #059669;
      border: 1px solid rgba(16, 185, 129, 0.2);
    }
  }
}

// 表格样式优化
:deep(.el-table) {
  border-radius: 16px;
  overflow: hidden;
  background: transparent;
  
  .el-table__header th {
    background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
    color: #475569;
    font-weight: 700;
    border-bottom: 2px solid #e2e8f0;
    padding: 20px 12px;
  }

  .el-table__row {
    transition: all 0.3s ease;
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(10px);

    &:hover {
      background: rgba(255, 255, 255, 0.95);
      transform: translateY(-2px);
      box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
    }
  }

  .el-table__cell {
    padding: 18px 12px;
    border-bottom: 1px solid rgba(226, 232, 240, 0.5);
  }
}

// 按钮样式优化
:deep(.el-button) {
  border-radius: 10px;
  font-weight: 600;
  transition: all 0.3s ease;

  &.el-button--text {
    &:hover {
      transform: translateY(-2px);
      color: #667eea;
    }
  }

  &.el-button--primary:not(.is-link) {
    background: #0366d6;
    border: 1px solid #0366d6;
    color: #fff;

    &:hover {
      background: #0256cc;
      border-color: #0256cc;
      transform: translateY(-2px);
      box-shadow: none;
    }
  }

  &.is-link {
    background: transparent;
    border: none;
    box-shadow: none;
    color: #0366d6;
    padding: 0 4px;
  }

  &.el-button--success {
    background: rgba(16, 185, 129, 0.12);
    color: #059669;
    border: 1px solid rgba(16, 185, 129, 0.3);

    &:hover {
      background: rgba(16, 185, 129, 0.18);
      box-shadow: none;
    }
  }

  &.el-button--warning {
    background: rgba(250, 140, 22, 0.12);
    color: #fa8c16;
    border: 1px solid rgba(250, 140, 22, 0.3);

    &:hover {
      background: rgba(250, 140, 22, 0.18);
      box-shadow: none;
    }
  }
}

// 标签样式优化
:deep(.el-tag) {
  border-radius: 8px;
  font-weight: 600;
  border: 1px solid rgba(0,0,0,0.06);
  padding: 4px 10px;

  &.el-tag--success {
    background: rgba(16, 185, 129, 0.1);
    color: #059669;
    border-color: rgba(16, 185, 129, 0.2);
  }

  &.el-tag--warning {
    background: rgba(250, 140, 22, 0.1);
    color: #fa8c16;
    border-color: rgba(250, 140, 22, 0.2);
  }

  &.el-tag--danger {
    background: rgba(255, 77, 79, 0.1);
    color: #ff4d4f;
    border-color: rgba(255,77,79,0.2);
  }

  &.el-tag--info {
    background: rgba(107, 114, 128, 0.08);
    color: #4b5563;
    border-color: rgba(107,114,128,0.2);
  }
}

// 状态列样式
.status-cell {
  .status-tag {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-bottom: 4px;
    padding: 6px 12px;
    border-radius: 12px;
    font-weight: 600;
    
    .status-icon {
      font-size: 14px;
    }
  }
  
  .progress-info {
    margin-top: 8px;
    
    :deep(.el-progress) {
      .el-progress-bar {
        .el-progress-bar__outer {
          background: rgba(226, 232, 240, 0.5);
          border-radius: 8px;
        }
        
        .el-progress-bar__inner {
          background: linear-gradient(135deg, #10b981 0%, #059669 100%);
          border-radius: 8px;
        }
      }
    }
  }
  
  .error-info {
    margin-top: 4px;
  }
}

// DNS统计卡片组件优化
:deep(.dns-stat-card) {
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 16px 48px rgba(0, 0, 0, 0.12);
  }
}

// 输入框和选择器优化
:deep(.el-input) {
  .el-input__wrapper {
    border-radius: 10px;
    border: 1px solid #e2e8f0;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(10px);
    transition: all 0.3s ease;

    &:hover {
      border-color: #94a3b8;
    }

    &.is-focus {
      border-color: #667eea;
      box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
    }
  }
}

:deep(.el-select) {
  .el-select__wrapper {
    border-radius: 10px;
    border: 1px solid #e2e8f0;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(10px);
    transition: all 0.3s ease;

    &:hover {
      border-color: #94a3b8;
    }

    &.is-focused {
      border-color: #667eea;
      box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
    }
  }
}

// 分页器优化
:deep(.el-pagination) {
  .btn-prev,
  .btn-next,
  .el-pager li {
    border-radius: 8px;
    font-weight: 600;
    transition: all 0.3s ease;
    border: 1px solid #e2e8f0;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(10px);

    &:hover {
      transform: translateY(-1px);
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      border-color: #667eea;
    }

    &.is-active {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      border-color: #667eea;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .dns-cert-container {
    padding: 16px;
  }

  .page-header {
    padding: 24px;

    .page-header-content {
      flex-direction: column;
      gap: 20px;
      align-items: flex-start;

      .page-title h1 {
        font-size: 28px;
      }

      .page-actions {
        width: 100%;
        display: flex;
        justify-content: flex-end;
      }
    }
  }

  .table-container {
    overflow-x: auto;
  }
}

// 加载动画
@keyframes shimmer {
  0% {
    background-position: -200px 0;
  }
  100% {
    background-position: calc(200px + 100%) 0;
  }
}

.loading-skeleton {
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200px 100%;
  animation: shimmer 1.5s infinite;
}
</style>
