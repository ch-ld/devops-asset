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
          <el-dropdown @command="handleCreateCommand" split-button type="primary" @click="handleAdvancedCreate">
            <el-icon><Plus /></el-icon>
            申请证书
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
                :model-value="row.auto_renew"
                @change="(value: boolean) => handleAutoRenewChange(row, value)"
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
    <CertificateDetail
      v-if="currentCertificate"
      v-model:visible="drawerVisible"
      :certificate="currentCertificate"
      @refresh="fetchData"
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
  Setting
} from '@element-plus/icons-vue'
import { DNSStatCard } from '@/components/dns'
import CertificateModal from './components/CertificateModal.vue'
import CertificateDetail from './components/CertificateDetail.vue'
import { certificateApi } from '@/api/dns/certificate'
import type { Certificate } from '@/types/dns'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const statisticsLoading = ref(false)
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
      router.push('/dns/certs/create?type=upload')
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

    // 格式选择选项
    const formatOptions = [
      { label: 'PEM格式 (通用)', value: 'pem', icon: '📄', desc: '标准PEM格式，适用于大多数服务器' },
      { label: 'Nginx配置', value: 'nginx', icon: '🌐', desc: 'Nginx服务器配置文件' },
      { label: 'Apache配置', value: 'apache', icon: '🔧', desc: 'Apache服务器配置文件' },
      { label: 'IIS配置', value: 'iis', icon: '🖥️', desc: 'Windows IIS服务器配置' },
      { label: 'Tomcat配置', value: 'tomcat', icon: '☕', desc: 'Tomcat服务器配置文件' },
      { label: '私钥文件', value: 'key', icon: '🔑', desc: '证书私钥文件' },
      { label: '证书文件', value: 'crt', icon: '📜', desc: '证书公钥文件' },
      { label: '证书链文件', value: 'chain', icon: '🔗', desc: '完整证书链文件' }
    ]

    // 使用更美观的选择对话框
    const { value: format } = await ElMessageBox({
      title: '选择下载格式',
      message: h('div', { style: 'max-height: 400px; overflow-y: auto;' }, [
        h('div', { style: 'margin-bottom: 16px; color: #606266; font-size: 14px;' },
          `为证书 "${domainName}" 选择下载格式：`
        ),
        ...formatOptions.map(option =>
          h('div', {
            key: option.value,
            style: {
              padding: '12px 16px',
              margin: '8px 0',
              border: '1px solid #e4e7ed',
              borderRadius: '6px',
              cursor: 'pointer',
              transition: 'all 0.3s',
              display: 'flex',
              alignItems: 'center',
              gap: '12px'
            },
            onClick: () => {
              (window as any).selectedFormat = option.value
              // 触发确认按钮点击
              setTimeout(() => {
                const confirmBtn = document.querySelector('.el-message-box__btns .el-button--primary') as HTMLElement
                if (confirmBtn) confirmBtn.click()
              }, 100)
            },
            onMouseenter: (e: Event) => {
              const target = e.target as HTMLElement
              target.style.backgroundColor = '#f5f7fa'
              target.style.borderColor = '#409eff'
            },
            onMouseleave: (e: Event) => {
              const target = e.target as HTMLElement
              target.style.backgroundColor = ''
              target.style.borderColor = '#e4e7ed'
            }
          }, [
            h('span', { style: 'fontSize: 20px;' }, option.icon),
            h('div', { style: 'flex: 1;' }, [
              h('div', { style: 'fontWeight: 500; marginBottom: 4px;' }, option.label),
              h('div', { style: 'fontSize: 12px; color: #909399;' }, option.desc)
            ])
          ])
        )
      ]),
      confirmButtonText: '确认下载',
      cancelButtonText: '取消',
      showCancelButton: true,
      beforeClose: (action, instance, done) => {
        if (action === 'confirm' && !(window as any).selectedFormat) {
          ElMessage.warning('请选择下载格式')
          return false
        }
        done()
      }
    })

    const selectedFormat = (window as any).selectedFormat || 'pem'
    delete (window as any).selectedFormat

    if (!selectedFormat) {
      return
    }

    // 下载证书
    const response = await certificateApi.download(row.id, selectedFormat)

    // 获取文件扩展名
    const getFileExtension = (format: string) => {
      const extensionMap: Record<string, string> = {
        'pem': 'pem',
        'nginx': 'conf',
        'apache': 'conf',
        'iis': 'txt',
        'tomcat': 'jks',
        'key': 'key',
        'crt': 'crt',
        'chain': 'pem'
      }
      return extensionMap[format] || 'txt'
    }

    // 创建下载链接
    const url = window.URL.createObjectURL(response as Blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${domainName}_${selectedFormat}.${getFileExtension(selectedFormat)}`
    link.click()
    window.URL.revokeObjectURL(url)

    const formatLabel = formatOptions.find(opt => opt.value === selectedFormat)?.label || selectedFormat
    ElMessage.success(`${formatLabel}下载成功`)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('证书下载失败:', error)
      ElMessage.error('证书下载失败，请稍后重试')
    }
  }
}

const handleCommand = (command: string, row: Certificate) => {
  const domainName = row.common_name || row.domain_name || `证书ID: ${row.id}`

  switch (command) {
    case 'deploy':
      ElMessage.info(`部署证书 ${domainName}`)
      break
    case 'revoke':
      handleRevoke(row)
      break
    case 'export':
      ElMessage.info(`导出证书 ${domainName} 的配置`)
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

// 导入证书
const handleImportCert = () => {
  // 跳转到申请证书页面，默认选择上传模式
  router.push('/dns/certs/create?type=upload')
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

  // 这里应该打开部署配置弹窗
  ElMessage.info('批量部署功能开发中...')
}

// 批量导出
const handleBatchExport = async () => {
  if (!hasSelected.value) return

  try {
    const ids = selectedRows.value.map(row => row.id)
    const response = await certificateApi.batchExport(ids)

    const blob = new Blob([response.content], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
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
    const response = await certificateApi.list(params)

    // 处理不同的响应格式
    console.log('证书列表API响应:', response)

    // 尝试多种可能的数据结构
    let certificateList = []
    let total = 0

    if (response.data?.items) {
      // 标准格式: { data: { items: [], total: number } }
      certificateList = response.data.items
      total = response.data.total || 0
    } else if (response.data?.list) {
      // 格式: { data: { list: [], total: number } }
      certificateList = response.data.list
      total = response.data.total || 0
    } else if (response.list) {
      // 格式: { list: [], total: number }
      certificateList = response.list
      total = response.total || 0
    } else if (Array.isArray(response.data)) {
      // 格式: { data: [] }
      certificateList = response.data
      total = response.data.length
    } else if (Array.isArray(response)) {
      // 格式: []
      certificateList = response
      total = response.length
    }

    certificates.value = certificateList
    pagination.total = total

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
