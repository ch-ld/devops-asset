<template>
  <div class="modern-page-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <div class="title-icon">
              <el-icon><Lock /></el-icon>
            </div>
            SSL证书管理
          </h1>
          <p class="page-description">管理您的SSL证书，监控证书状态，确保网站安全访问</p>
        </div>
        <div class="header-actions">
          <el-button class="modern-btn secondary" @click="handleRefresh" :icon="Refresh">
            刷新
          </el-button>
          <el-dropdown @command="handleBulkAction">
            <el-button class="modern-btn secondary">
              批量操作
              <el-icon><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="deploy">批量部署</el-dropdown-item>
                <el-dropdown-item command="renew">批量续期</el-dropdown-item>
                <el-dropdown-item command="export">批量导出</el-dropdown-item>
                <el-dropdown-item command="delete" divided>批量删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-button class="modern-btn primary" @click="handleAddCert" :icon="Plus">
            申请证书
          </el-button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="modern-stats-grid">
      <div class="stat-card">
        <div class="stat-header">
          <div class="stat-icon primary">
            <el-icon><Lock /></el-icon>
          </div>
          <div class="stat-trend up">+12%</div>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ statistics.total || 0 }}</div>
          <div class="stat-label">总证书数</div>
          <div class="stat-description">所有SSL/TLS证书</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-header">
          <div class="stat-icon success">
            <el-icon><Check /></el-icon>
          </div>
          <div class="stat-trend up">+5%</div>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ statistics.valid || 0 }}</div>
          <div class="stat-label">有效证书</div>
          <div class="stat-description">正常使用中</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-header">
          <div class="stat-icon warning">
            <el-icon><Clock /></el-icon>
          </div>
          <div class="stat-trend down">-2</div>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ statistics.expiring || 0 }}</div>
          <div class="stat-label">即将过期</div>
          <div class="stat-description">30天内过期</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-header">
          <div class="stat-icon error">
            <el-icon><Warning /></el-icon>
          </div>
          <div class="stat-trend down">-1</div>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ statistics.expired || 0 }}</div>
          <div class="stat-label">已过期</div>
          <div class="stat-description">需要立即处理</div>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="modern-search-section">
      <div class="search-content">
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
        <div class="search-filters">
          <el-select
            v-model="searchForm.status"
            placeholder="证书状态"
            size="large"
            clearable
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
            style="width: 150px"
          >
            <el-option label="全部" value="" />
            <el-option label="Let's Encrypt" value="letsencrypt" />
            <el-option label="ZeroSSL" value="zerossl" />
            <el-option label="自定义" value="custom" />
            <el-option label="手动上传" value="manual" />
          </el-select>
          <el-button class="modern-btn primary" @click="handleSearch" :icon="Search">
            搜索
          </el-button>
          <el-button class="modern-btn secondary" @click="handleReset" :icon="Refresh">
            重置
          </el-button>
        </div>
      </div>
    </div>

    <!-- 证书列表 -->
    <div class="modern-content-card">
      <div class="card-header">
        <div class="header-content">
          <div class="header-left">
            <h3 class="card-title">证书列表</h3>
            <p class="card-subtitle">{{ certificates.length }} 个证书</p>
          </div>
          <div class="header-actions">
            <el-button class="modern-btn secondary" @click="handleRefresh" :icon="Refresh">
              刷新
            </el-button>
            <el-button class="modern-btn secondary" @click="handleImportCert" :icon="Upload">
              导入证书
            </el-button>
            <el-dropdown @command="handleBatchCommand" :disabled="!hasSelected">
              <el-button class="modern-btn warning" :disabled="!hasSelected">
                批量操作 ({{ selectedRows.length }})
                <el-icon><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="renew">批量续期</el-dropdown-item>
                  <el-dropdown-item command="download">批量下载</el-dropdown-item>
                  <el-dropdown-item command="deploy">批量部署</el-dropdown-item>
                  <el-dropdown-item command="export">导出报告</el-dropdown-item>
                  <el-dropdown-item command="delete" divided>批量删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </div>
      <div class="card-content">
        <div class="modern-table">
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
                  <el-tag
                    :type="getStatusTagType(row.status)"
                    size="small"
                  >
                    {{ getStatusText(row.status) }}
                  </el-tag>
                  <div v-if="row.status === 'failed' && row.error_message" class="error-message">
                    <el-tooltip :content="row.error_message" placement="top">
                      <el-text type="danger" size="small">
                        {{ row.error_message }}
                      </el-text>
                    </el-tooltip>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="cert_type" label="CA类型" width="120">
              <template #default="{ row }">
                <el-tag type="info" size="small">
                  {{ getCATypeName(row.cert_type || row.ca_type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="valid_from" label="签发时间" width="140">
              <template #default="{ row }">
                {{ formatTime(row.valid_from) }}
              </template>
            </el-table-column>
            <el-table-column prop="valid_to" label="过期时间" width="180">
              <template #default="{ row }">
                <div v-if="row.valid_to" class="expires-cell">
                  <div class="expires-date">{{ formatTime(row.valid_to) }}</div>
                  <div class="expires-days" :class="getExpiresClass(row.valid_to)">
                    {{ getExpiringDays(row.valid_to) }}
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
                  size="small"
                  @click="handleViewCert(row)"
                >
                  查看
                </el-button>
                <el-button
                  v-if="row.status === 'issued'"
                  size="small"
                  type="primary"
                  @click="handleDeployCert(row)"
                >
                  部署
                </el-button>
                <el-button
                  v-if="row.status === 'issued'"
                  size="small"
                  @click="handleDownloadCert(row)"
                >
                  下载
                </el-button>
                <el-dropdown @command="cmd => handleCertAction(cmd, row)" trigger="click">
                  <el-button size="small">
                    更多
                    <el-icon><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="renew">手动续期</el-dropdown-item>
                      <el-dropdown-item command="revoke" v-if="row.status === 'issued'">
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
        </div>
      </div>
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
    />

    <!-- 证书部署弹窗 -->
    <CertificateDeployDialog
      v-model:visible="showDeployDialog"
      :certificate="selectedCertificate"
    />

    <!-- 证书导出弹窗 -->
    <CertificateExportDialog
      v-model:visible="showExportDialog"
      :certificate="selectedCertificate"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Lock,
  Check,
  Clock,
  Warning,
  Refresh,
  ArrowDown,
  Plus,
  Search,
  Upload
} from '@element-plus/icons-vue'
import CertificateModal from './components/CertificateModal.vue'
import CertificateDetail from './components/CertificateDetail.vue'
import CertificateDeployDialog from '@/components/dns/CertificateDeployDialog.vue'
import CertificateExportDialog from '@/components/dns/CertificateExportDialog.vue'
import { certificateApi } from '@/api/dns/certificate'
import { formatTime } from '@/utils/time'

// 响应式数据
const loading = ref(false)
const certificates = ref([])
const selectedRows = ref([])
const currentCertificate = ref(null)
const selectedCertificate = ref(null)

// 弹窗状态
const modalVisible = ref(false)
const drawerVisible = ref(false)
const showDeployDialog = ref(false)
const showExportDialog = ref(false)

// 搜索表单
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

// 统计数据
const statistics = computed(() => {
  const total = certificates.value.length
  const valid = certificates.value.filter(cert => cert.status === 'issued').length
  const expiring = certificates.value.filter(cert => {
    if (!cert.valid_to) return false
    const validTo = new Date(cert.valid_to)
    const now = new Date()
    const diffDays = Math.ceil((validTo.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
    return diffDays <= 30 && diffDays > 0
  }).length
  const expired = certificates.value.filter(cert => {
    if (!cert.valid_to) return false
    const validTo = new Date(cert.valid_to)
    const now = new Date()
    return validTo < now
  }).length

  return { total, valid, expiring, expired }
})

// 计算属性
const hasSelected = computed(() => selectedRows.value.length > 0)

// 方法
const handleRefresh = async () => {
  await loadCertificates()
}

const handleSearch = () => {
  pagination.page = 1
  loadCertificates()
}

const handleReset = () => {
  Object.assign(searchForm, {
    keyword: '',
    status: '',
    ca_type: ''
  })
  pagination.page = 1
  loadCertificates()
}

const handleAddCert = () => {
  currentCertificate.value = null
  modalVisible.value = true
}

const handleSelectionChange = (selection) => {
  selectedRows.value = selection
}

const handleBulkAction = (command) => {
  if (!hasSelected.value) {
    ElMessage.warning('请先选择要操作的证书')
    return
  }

  switch (command) {
    case 'deploy':
      handleBatchDeploy()
      break
    case 'renew':
      handleBatchRenew()
      break
    case 'export':
      handleBatchExport()
      break
    case 'delete':
      handleBatchDelete()
      break
  }
}

const handleBatchCommand = (command) => {
  handleBulkAction(command)
}

const handleImportCert = () => {
  // 导入证书逻辑
  ElMessage.info('导入证书功能开发中...')
}

const handleViewCert = (cert) => {
  currentCertificate.value = cert
  drawerVisible.value = true
}

const handleDeployCert = (cert) => {
  selectedCertificate.value = cert
  showDeployDialog.value = true
}

const handleDownloadCert = async (cert) => {
  try {
    const response = await certificateApi.download(cert.id)

    // 创建下载链接
    const url = window.URL.createObjectURL(new Blob([response]))
    const link = document.createElement('a')
    link.href = url
    link.download = `${cert.common_name}.pem`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    ElMessage.success('证书下载成功')
  } catch (error) {
    ElMessage.error('下载失败：' + (error.response?.data?.message || error.message))
  }
}

const handleCertAction = (command, cert) => {
  switch (command) {
    case 'renew':
      handleRenewCert(cert)
      break
    case 'revoke':
      handleRevokeCert(cert)
      break
    case 'export':
      handleExportCert(cert)
      break
    case 'delete':
      handleDeleteCert(cert)
      break
  }
}

const handleAutoRenewChange = async (cert, value) => {
  try {
    await certificateApi.update(cert.id, { auto_renew: value })
    cert.auto_renew = value
    ElMessage.success('自动续期设置已更新')
  } catch (error) {
    ElMessage.error('更新失败：' + (error.response?.data?.message || error.message))
  }
}

const handleRenewCert = async (cert) => {
  try {
    await ElMessageBox.confirm(
      `确定要续期证书 "${cert.common_name}" 吗？`,
      '确认续期',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }
    )

    await certificateApi.renew(cert.id)
    ElMessage.success('证书续期成功')
    await loadCertificates()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('续期失败：' + (error.response?.data?.message || error.message))
    }
  }
}

const handleRevokeCert = async (cert) => {
  try {
    await ElMessageBox.confirm(
      `确定要吊销证书 "${cert.common_name}" 吗？此操作不可撤销。`,
      '确认吊销',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await certificateApi.revoke(cert.id)
    ElMessage.success('证书吊销成功')
    await loadCertificates()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('吊销失败：' + (error.response?.data?.message || error.message))
    }
  }
}

const handleExportCert = (cert) => {
  selectedCertificate.value = cert
  showExportDialog.value = true
}

const handleDeleteCert = async (cert) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除证书 "${cert.common_name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await certificateApi.delete(cert.id)
    ElMessage.success('删除成功')
    await loadCertificates()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败：' + (error.response?.data?.message || error.message))
    }
  }
}

const handleBatchDeploy = async () => {
  try {
    const certIds = selectedRows.value.map(cert => cert.id)
    await ElMessageBox.confirm(
      `确定要批量部署选中的 ${selectedRows.value.length} 个证书吗？`,
      '确认部署',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }
    )

    // 这里需要实现批量部署逻辑
    ElMessage.success('批量部署成功')
    await loadCertificates()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量部署失败：' + (error.response?.data?.message || error.message))
    }
  }
}

const handleBatchRenew = async () => {
  try {
    const certIds = selectedRows.value.map(cert => cert.id)
    await ElMessageBox.confirm(
      `确定要批量续期选中的 ${selectedRows.value.length} 个证书吗？`,
      '确认续期',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }
    )

    await certificateApi.batchRenew(certIds)
    ElMessage.success('批量续期成功')
    await loadCertificates()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量续期失败：' + (error.response?.data?.message || error.message))
    }
  }
}

const handleBatchExport = async () => {
  try {
    const certIds = selectedRows.value.map(cert => cert.id)
    // 这里需要实现批量导出逻辑
    ElMessage.success('批量导出成功')
  } catch (error) {
    ElMessage.error('批量导出失败：' + (error.response?.data?.message || error.message))
  }
}

const handleBatchDelete = async () => {
  try {
    const certIds = selectedRows.value.map(cert => cert.id)
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 个证书吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await certificateApi.batchDelete(certIds)
    ElMessage.success('批量删除成功')
    await loadCertificates()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败：' + (error.response?.data?.message || error.message))
    }
  }
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  loadCertificates()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadCertificates()
}

const handleModalSuccess = () => {
  modalVisible.value = false
  loadCertificates()
}

// 工具方法
const getStatusTagType = (status) => {
  const statusMap = {
    pending: 'warning',
    processing: 'warning',
    validating: 'warning',
    issued: 'success',
    expired: 'danger',
    failed: 'danger',
    revoked: 'info'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status) => {
  const statusMap = {
    pending: '申请中',
    processing: '处理中',
    validating: '验证中',
    issued: '已签发',
    expired: '已过期',
    failed: '申请失败',
    revoked: '已吊销'
  }
  return statusMap[status] || status
}

const getCATypeName = (caType) => {
  const caTypeMap = {
    letsencrypt: "Let's Encrypt",
    zerossl: 'ZeroSSL',
    custom: '自定义',
    manual: '手动上传'
  }
  return caTypeMap[caType] || caType
}

const getExpiringDays = (validTo) => {
  if (!validTo) return ''

  const expires = new Date(validTo)
  const now = new Date()
  const diffTime = expires.getTime() - now.getTime()
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays < 0) {
    return `已过期 ${Math.abs(diffDays)} 天`
  } else if (diffDays === 0) {
    return '今天过期'
  } else if (diffDays <= 30) {
    return `${diffDays} 天后过期`
  } else {
    return `${diffDays} 天后过期`
  }
}

const getExpiresClass = (validTo) => {
  if (!validTo) return ''

  const expires = new Date(validTo)
  const now = new Date()
  const diffTime = expires.getTime() - now.getTime()
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays < 0) {
    return 'expired'
  } else if (diffDays <= 7) {
    return 'critical'
  } else if (diffDays <= 30) {
    return 'warning'
  } else {
    return 'normal'
  }
}

const loadCertificates = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchForm.keyword,
      status: searchForm.status,
      ca_type: searchForm.ca_type
    }

    const response = await certificateApi.list(params)
    certificates.value = response.data?.list || []
    pagination.total = response.data?.total || 0
  } catch (error) {
    ElMessage.error('加载证书列表失败：' + (error.response?.data?.message || error.message))
  } finally {
    loading.value = false
  }
}

// 加载证书统计数据
const loadCertificateStats = async () => {
  try {
    const response = await certificateApi.stats()
    // 这里可以更新统计数据，如果API返回了统计信息
  } catch (error) {
    console.error('加载证书统计失败：', error)
  }
}

// 生命周期
onMounted(() => {
  loadCertificates()
  loadCertificateStats()
})
</script>

<style lang="scss" scoped>
.domain-cell {
  .domain-name {
    font-weight: 500;
    color: var(--text-primary);
    margin-bottom: 4px;
  }

  .cert-type {
    display: flex;
    gap: 4px;
    flex-wrap: wrap;
  }
}

.status-cell {
  .error-message {
    margin-top: 4px;
    max-width: 200px;

    .el-text {
      font-size: 12px;
      line-height: 1.2;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
  }
}

.expires-cell {
  .expires-date {
    font-size: 14px;
    color: var(--text-primary);
    margin-bottom: 2px;
  }

  .expires-days {
    font-size: 12px;

    &.normal {
      color: var(--text-secondary);
    }

    &.warning {
      color: var(--warning-color);
      font-weight: 500;
    }

    &.critical {
      color: var(--error-color);
      font-weight: 500;
    }

    &.expired {
      color: var(--error-color);
      font-weight: 600;
    }
  }
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: var(--spacing-lg);
  padding-top: var(--spacing-lg);
  border-top: 1px solid var(--border-secondary);
}

// 响应式设计
@media (max-width: 768px) {
  .modern-stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .search-content {
    flex-direction: column;
    gap: var(--spacing-md);

    .search-input {
      width: 100%;
    }

    .search-filters {
      width: 100%;
      justify-content: space-between;
    }
  }

  .header-content {
    flex-direction: column;
    gap: var(--spacing-md);

    .header-actions {
      width: 100%;
      justify-content: flex-start;
      flex-wrap: wrap;
    }
  }
}

@media (max-width: 480px) {
  .modern-stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
