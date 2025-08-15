<template>
  <el-dialog
    :model-value="visible"
    title="SSL证书管理"
    width="900px"
    :before-close="handleClose"
    class="certificate-modal"
  >
    <div class="certificate-container">
      <!-- 证书头部 -->
      <div class="certificate-header">
        <div class="domain-info">
          <div class="domain-icon">
            <el-icon size="24"><Lock /></el-icon>
          </div>
          <div class="domain-content">
            <h3 class="domain-name">{{ domain?.name }}</h3>
            <p class="domain-desc">SSL证书信息管理</p>
          </div>
        </div>
        <div class="header-actions">
          <el-button 
            type="success" 
            @click="handleCheck"
            :loading="checkLoading"
            class="check-btn"
          >
            <el-icon><View /></el-icon>
            检查证书
          </el-button>
          <el-button 
            type="primary" 
            @click="handleAdd"
            class="add-btn"
          >
            <el-icon><Plus /></el-icon>
            添加证书
          </el-button>
        </div>
      </div>

      <!-- 证书列表 -->
      <div class="certificate-content">
        <el-card class="certificate-list-card">
          <template #header>
            <div class="card-header">
              <span class="header-title">证书列表</span>
              <div class="header-actions">
                <el-button 
                  size="small" 
                  @click="refreshCertificates"
                  :loading="loading"
                  class="refresh-btn"
                >
                  <el-icon><Refresh /></el-icon>
                  刷新
                </el-button>
              </div>
            </div>
          </template>

          <el-table 
            :data="certificates" 
            :loading="loading"
            class="certificate-table"
            empty-text="暂无证书记录"
          >
            <el-table-column prop="common_name" label="域名" min-width="150">
              <template #default="{ row }">
                <div class="domain-cell">
                  <span class="domain-name">{{ row.common_name }}</span>
                  <div v-if="row.san_domains && row.san_domains.length > 0" class="san-domains">
                    <el-tag 
                      v-for="san in row.san_domains.slice(0, 3)" 
                      :key="san" 
                      size="small" 
                      type="info"
                      class="san-tag"
                    >
                      {{ san }}
                    </el-tag>
                    <el-tag 
                      v-if="row.san_domains.length > 3" 
                      size="small" 
                      type="info"
                    >
                      +{{ row.san_domains.length - 3 }}
                    </el-tag>
                  </div>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="issuer" label="颁发机构" width="150">
              <template #default="{ row }">
                <div class="issuer-cell">
                  <span class="issuer-name">{{ getIssuerName(row.issuer) }}</span>
                  <span class="issuer-type">{{ getCertType(row.cert_type) }}</span>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag 
                  :type="getCertStatusType(row.status)" 
                  effect="light"
                  size="small"
                >
                  {{ getCertStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            
            <el-table-column prop="valid_from" label="有效期" width="200">
              <template #default="{ row }">
                <div class="validity-cell">
                  <div class="validity-dates">
                    <span class="from-date">{{ formatDate(row.valid_from) }}</span>
                    <span class="separator">至</span>
                    <span class="to-date">{{ formatDate(row.valid_to) }}</span>
                  </div>
                  <div class="validity-status" :class="getValidityClass(row.valid_to)">
                    {{ getValidityText(row.valid_to) }}
                  </div>
                </div>
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
            
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <div class="action-buttons">
                  <el-button type="primary" size="small" text @click="handleView(row)">
                    查看
                  </el-button>
                  <el-button type="success" size="small" text @click="handleDownload(row)">
                    下载
                  </el-button>
                  <el-dropdown @command="(command) => handleCommand(command, row)">
                    <el-button type="primary" size="small" text>
                      更多
                      <el-icon><ArrowDown /></el-icon>
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item command="renew">续费</el-dropdown-item>
                        <el-dropdown-item command="revoke">吊销</el-dropdown-item>
                        <el-dropdown-item command="delete" divided>删除</el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <!-- 证书检查结果 -->
        <el-card v-if="checkResult" class="check-result-card">
          <template #header>
            <div class="card-header">
              <span class="header-title">证书检查结果</span>
              <div class="header-actions">
                <el-tag :type="getCheckResultType(checkResult.status)" effect="light">
                  {{ checkResult.status }}
                </el-tag>
              </div>
            </div>
          </template>

          <div class="check-result-content">
            <div class="result-grid">
              <div class="result-item">
                <label>HTTPS状态</label>
                <span :class="getHttpsStatusClass(checkResult.https_status)">
                  {{ checkResult.https_status }}
                </span>
              </div>
              <div class="result-item">
                <label>证书链完整性</label>
                <span :class="getChainStatusClass(checkResult.chain_valid)">
                  {{ checkResult.chain_valid ? '完整' : '不完整' }}
                </span>
              </div>
              <div class="result-item">
                <label>证书有效期</label>
                <span :class="getValidityClass(checkResult.expires_at)">
                  {{ formatDate(checkResult.expires_at) }}
                </span>
              </div>
              <div class="result-item">
                <label>证书颁发机构</label>
                <span>{{ checkResult.issuer }}</span>
              </div>
            </div>
            
            <div v-if="checkResult.issues && checkResult.issues.length > 0" class="issues-section">
              <h4 class="issues-title">发现的问题</h4>
              <ul class="issues-list">
                <li v-for="issue in checkResult.issues" :key="issue" class="issue-item">
                  <el-icon class="issue-icon"><Warning /></el-icon>
                  {{ issue }}
                </li>
              </ul>
            </div>
          </div>
        </el-card>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose" class="close-btn">
          关闭
        </el-button>
        <el-button 
          type="primary" 
          @click="exportReport"
          :disabled="!certificates.length"
          class="export-btn"
        >
          <el-icon><Download /></el-icon>
          导出报告
        </el-button>
      </div>
    </template>

    <!-- 证书详情弹窗 -->
    <CertificateDetailModal
      v-model:visible="detailModalVisible"
      :certificate="currentCertificate"
    />

    <!-- 添加证书弹窗 -->
    <AddCertificateModal
      v-model:visible="addModalVisible"
      :domain="domain"
      @success="handleAddSuccess"
    />
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Lock, 
  View, 
  Plus, 
  Refresh, 
  ArrowDown, 
  Warning, 
  Download 
} from '@element-plus/icons-vue'
import CertificateDetailModal from './CertificateDetailModal.vue'
import AddCertificateModal from './AddCertificateModal.vue'
import { certificateApi } from '@/api/dns/certificate'
import type { Domain } from '@/types/dns'

interface Props {
  visible: boolean
  domain?: Domain | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
}

interface Certificate {
  id: number
  common_name: string
  san_domains: string[]
  issuer: string
  cert_type: string
  status: string
  valid_from: string
  valid_to: string
  auto_renew: boolean
  fingerprint: string
}

interface CheckResult {
  status: string
  https_status: string
  chain_valid: boolean
  expires_at: string
  issuer: string
  issues: string[]
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const checkLoading = ref(false)
const certificates = ref<Certificate[]>([])
const checkResult = ref<CheckResult | null>(null)
const detailModalVisible = ref(false)
const addModalVisible = ref(false)
const currentCertificate = ref<Certificate | null>(null)

// 方法
const handleClose = () => {
  emit('update:visible', false)
}

const refreshCertificates = async () => {
  if (!props.domain) return
  
  try {
    loading.value = true
    const response = await certificateApi.list({ domain_id: props.domain.id })
    certificates.value = response.data?.items || []
  } catch (error) {
    ElMessage.error('获取证书列表失败')
  } finally {
    loading.value = false
  }
}

const handleCheck = async () => {
  if (!props.domain) return
  
  try {
    checkLoading.value = true
    const response = await certificateApi.check(props.domain.name)
    checkResult.value = response.data
    ElMessage.success('证书检查完成')
  } catch (error) {
    ElMessage.error('证书检查失败')
    checkResult.value = null
  } finally {
    checkLoading.value = false
  }
}

const handleAdd = () => {
  addModalVisible.value = true
}

const handleView = (certificate: Certificate) => {
  currentCertificate.value = certificate
  detailModalVisible.value = true
}

const handleDownload = async (certificate: Certificate) => {
  try {
    const blob = await certificateApi.download(certificate.id)
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${certificate.common_name}-certificate.zip`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    ElMessage.success('下载成功')
  } catch (error) {
    ElMessage.error('下载失败')
  }
}

const handleCommand = async (command: string, certificate: Certificate) => {
  switch (command) {
    case 'renew':
      handleRenew(certificate)
      break
    case 'revoke':
      handleRevoke(certificate)
      break
    case 'delete':
      handleDelete(certificate)
      break
  }
}

const handleRenew = async (certificate: Certificate) => {
  try {
    await ElMessageBox.confirm(
      `确定要续费证书 "${certificate.common_name}" 吗？`,
      '确认续费',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info',
      }
    )
    
    await certificateApi.renew(certificate.id)
    ElMessage.success('续费成功')
    await refreshCertificates()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('续费失败')
    }
  }
}

const handleRevoke = async (certificate: Certificate) => {
  try {
    await ElMessageBox.confirm(
      `确定要吊销证书 "${certificate.common_name}" 吗？此操作不可恢复。`,
      '确认吊销',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await certificateApi.revoke(certificate.id)
    ElMessage.success('吊销成功')
    await refreshCertificates()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('吊销失败')
    }
  }
}

const handleDelete = async (certificate: Certificate) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除证书 "${certificate.common_name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await certificateApi.delete(certificate.id)
    ElMessage.success('删除成功')
    await refreshCertificates()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleAutoRenewChange = async (certificate: Certificate, value: boolean) => {
  try {
    await certificateApi.update(certificate.id, { auto_renew: value })
    certificate.auto_renew = value
    ElMessage.success(`${value ? '启用' : '禁用'}自动续费成功`)
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleAddSuccess = () => {
  addModalVisible.value = false
  refreshCertificates()
}

const exportReport = async () => {
  if (!props.domain) return
  
  try {
    const blob = await certificateApi.exportReport(props.domain.id)
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${props.domain.name}-certificates-${new Date().toISOString().split('T')[0]}.pdf`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

// 工具方法
const formatDate = (date: string) => {
  if (!date) return '-'
  try {
    return new Date(date).toLocaleDateString('zh-CN')
  } catch {
    return date
  }
}

const getIssuerName = (issuer: string) => {
  const issuerMap: Record<string, string> = {
    'Let\'s Encrypt': 'Let\'s Encrypt',
    'DigiCert': 'DigiCert',
    'Comodo': 'Comodo',
    'GeoTrust': 'GeoTrust',
    'Symantec': 'Symantec'
  }
  
  for (const [key, value] of Object.entries(issuerMap)) {
    if (issuer.includes(key)) return value
  }
  
  return issuer.split(',')[0] || issuer
}

const getCertType = (type: string) => {
  const typeMap: Record<string, string> = {
    'DV': 'DV',
    'OV': 'OV', 
    'EV': 'EV',
    'Wildcard': '通配符'
  }
  return typeMap[type] || type
}

const getCertStatusType = (status: string) => {
  const typeMap: Record<string, any> = {
    'valid': 'success',
    'expired': 'danger',
    'expiring': 'warning',
    'revoked': 'danger',
    'pending': 'info'
  }
  return typeMap[status] || 'info'
}

const getCertStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    'valid': '有效',
    'expired': '已过期',
    'expiring': '即将过期',
    'revoked': '已吊销',
    'pending': '待处理'
  }
  return textMap[status] || status
}

const getValidityClass = (date: string) => {
  if (!date) return ''
  try {
    const expireDate = new Date(date)
    const now = new Date()
    const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
    
    if (diffDays < 0) return 'text-danger'
    if (diffDays <= 30) return 'text-warning'
    return 'text-success'
  } catch {
    return ''
  }
}

const getValidityText = (date: string) => {
  if (!date) return ''
  try {
    const expireDate = new Date(date)
    const now = new Date()
    const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
    
    if (diffDays < 0) return `已过期 ${Math.abs(diffDays)} 天`
    if (diffDays <= 30) return `${diffDays} 天后过期`
    return `还有 ${diffDays} 天`
  } catch {
    return ''
  }
}

const getCheckResultType = (status: string) => {
  const typeMap: Record<string, any> = {
    'valid': 'success',
    'warning': 'warning',
    'error': 'danger'
  }
  return typeMap[status] || 'info'
}

const getHttpsStatusClass = (status: string) => {
  return status === 'valid' ? 'text-success' : 'text-danger'
}

const getChainStatusClass = (valid: boolean) => {
  return valid ? 'text-success' : 'text-danger'
}

// 监听弹窗显示状态
watch(() => props.visible, (visible) => {
  if (visible && props.domain) {
    refreshCertificates()
  } else if (!visible) {
    checkResult.value = null
    certificates.value = []
  }
})
</script>

<style scoped lang="scss">
.certificate-modal {
  :deep(.el-dialog) {
    border-radius: 16px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  }
  
  :deep(.el-dialog__header) {
    padding: 24px 24px 0;
    border-bottom: 1px solid #f0f2f5;
    
    .el-dialog__title {
      font-size: 18px;
      font-weight: 600;
      color: #1f2937;
    }
  }
  
  :deep(.el-dialog__body) {
    padding: 24px;
    max-height: 70vh;
    overflow-y: auto;
  }
  
  :deep(.el-dialog__footer) {
    padding: 0 24px 24px;
    border-top: 1px solid #f0f2f5;
  }
}

.certificate-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.certificate-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

.domain-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.domain-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.domain-content {
  .domain-name {
    margin: 0 0 4px 0;
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
  }
  
  .domain-desc {
    margin: 0;
    color: #6b7280;
    font-size: 14px;
  }
}

.header-actions {
  display: flex;
  gap: 12px;
  
  .check-btn {
    background: linear-gradient(135deg, #10b981 0%, #059669 100%);
    border: none;
    border-radius: 8px;
    padding: 12px 24px;
    font-weight: 600;
  }
  
  .add-btn {
    background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
    border: none;
    border-radius: 8px;
    padding: 12px 24px;
    font-weight: 600;
  }
}

.certificate-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.certificate-list-card,
.check-result-card {
  border-radius: 12px;
  border: none;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  
  :deep(.el-card__header) {
    padding: 20px 24px;
    border-bottom: 1px solid #f0f2f5;
  }
  
  :deep(.el-card__body) {
    padding: 24px;
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  .header-title {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
  }
  
  .refresh-btn {
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    color: #475569;
  }
}

.certificate-table {
  .domain-cell {
    .domain-name {
      font-weight: 600;
      margin-bottom: 4px;
    }
    
    .san-domains {
      display: flex;
      gap: 4px;
      flex-wrap: wrap;
      
      .san-tag {
        font-size: 10px;
        padding: 1px 6px;
      }
    }
  }
  
  .issuer-cell {
    display: flex;
    flex-direction: column;
    gap: 2px;
    
    .issuer-name {
      font-weight: 500;
    }
    
    .issuer-type {
      font-size: 12px;
      color: #6b7280;
    }
  }
  
  .validity-cell {
    .validity-dates {
      display: flex;
      flex-direction: column;
      gap: 2px;
      font-size: 12px;
      margin-bottom: 4px;
      
      .separator {
        color: #6b7280;
        font-size: 10px;
      }
    }
    
    .validity-status {
      font-size: 11px;
      font-weight: 500;
      
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
}

.check-result-content {
  .result-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 16px;
    margin-bottom: 24px;
  }
  
  .result-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
    
    label {
      font-size: 12px;
      font-weight: 500;
      color: #6b7280;
    }
    
    span {
      font-size: 14px;
      font-weight: 500;
      
      &.text-success {
        color: #10b981;
      }
      
      &.text-warning {
        color: #f59e0b;
      }
      
      &.text-danger {
        color: #ef4444;
      }
    }
  }
  
  .issues-section {
    padding: 16px;
    background: #fef3c7;
    border: 1px solid #fbbf24;
    border-radius: 8px;
    
    .issues-title {
      margin: 0 0 12px 0;
      font-size: 14px;
      font-weight: 600;
      color: #92400e;
    }
    
    .issues-list {
      margin: 0;
      padding: 0;
      list-style: none;
      
      .issue-item {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 8px;
        font-size: 13px;
        color: #92400e;
        
        &:last-child {
          margin-bottom: 0;
        }
        
        .issue-icon {
          color: #f59e0b;
          font-size: 14px;
        }
      }
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 20px;
  
  .close-btn {
    padding: 10px 20px;
    border-radius: 8px;
    border: 1px solid #d1d5db;
    background: #f9fafb;
    color: #374151;
    font-weight: 500;
  }
  
  .export-btn {
    padding: 10px 20px;
    border-radius: 8px;
    background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
    border: none;
    font-weight: 600;
  }
}

/* 状态颜色 */
.text-success {
  color: #10b981 !important;
}

.text-warning {
  color: #f59e0b !important;
}

.text-danger {
  color: #ef4444 !important;
}
</style>
