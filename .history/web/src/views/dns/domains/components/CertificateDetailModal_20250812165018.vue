<template>
  <el-dialog
    :model-value="visible"
    title="证书详情"
    width="700px"
    :before-close="handleClose"
    class="certificate-detail-modal"
  >
    <div v-if="certificate" class="certificate-detail">
      <!-- 证书基本信息 -->
      <el-card class="detail-card">
        <template #header>
          <span class="card-title">证书信息</span>
        </template>
        
        <div class="detail-grid">
          <div class="detail-item">
            <label>通用名称</label>
            <span>{{ certificate.common_name }}</span>
          </div>
          <div class="detail-item">
            <label>SAN域名</label>
            <div v-if="certificate.san_domains && certificate.san_domains.length > 0" class="san-list">
              <el-tag 
                v-for="san in certificate.san_domains" 
                :key="san" 
                size="small" 
                type="info"
                class="san-tag"
              >
                {{ san }}
              </el-tag>
            </div>
            <span v-else>-</span>
          </div>
          <div class="detail-item">
            <label>颁发机构</label>
            <span>{{ certificate.issuer }}</span>
          </div>
          <div class="detail-item">
            <label>证书类型</label>
            <span>{{ certificate.cert_type }}</span>
          </div>
          <div class="detail-item">
            <label>证书状态</label>
            <el-tag :type="getCertStatusType(certificate.status)" effect="light">
              {{ getCertStatusText(certificate.status) }}
            </el-tag>
          </div>
          <div class="detail-item">
            <label>证书指纹</label>
            <span class="fingerprint">{{ certificate.fingerprint }}</span>
          </div>
        </div>
      </el-card>

      <!-- 有效期信息 -->
      <el-card class="detail-card">
        <template #header>
          <span class="card-title">有效期</span>
        </template>
        
        <div class="validity-info">
          <div class="validity-item">
            <label>生效时间</label>
            <span>{{ formatDateTime(certificate.valid_from) }}</span>
          </div>
          <div class="validity-item">
            <label>过期时间</label>
            <span :class="getValidityClass(certificate.valid_to)">
              {{ formatDateTime(certificate.valid_to) }}
            </span>
          </div>
          <div class="validity-item">
            <label>剩余天数</label>
            <span :class="getValidityClass(certificate.valid_to)">
              {{ getRemainingDays(certificate.valid_to) }}
            </span>
          </div>
        </div>
      </el-card>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
interface Props {
  visible: boolean
  certificate?: any | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const handleClose = () => {
  emit('update:visible', false)
}

const getCertStatusType = (status: string) => {
  const typeMap: Record<string, any> = {
    'valid': 'success',
    'expired': 'danger',
    'expiring': 'warning',
    'revoked': 'danger'
  }
  return typeMap[status] || 'info'
}

const getCertStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    'valid': '有效',
    'expired': '已过期',
    'expiring': '即将过期',
    'revoked': '已吊销'
  }
  return textMap[status] || status
}

const formatDateTime = (date: string) => {
  if (!date) return '-'
  try {
    return new Date(date).toLocaleString('zh-CN')
  } catch {
    return date
  }
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

const getRemainingDays = (date: string) => {
  if (!date) return '-'
  try {
    const expireDate = new Date(date)
    const now = new Date()
    const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
    
    if (diffDays < 0) return `已过期 ${Math.abs(diffDays)} 天`
    return `${diffDays} 天`
  } catch {
    return '-'
  }
}
</script>

<style scoped lang="scss">
.certificate-detail-modal {
  :deep(.el-dialog) {
    border-radius: 16px;
  }
}

.certificate-detail {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.detail-card {
  border-radius: 12px;
  border: none;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  
  .card-title {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
  }
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
}

.detail-item {
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
    color: #1f2937;
    
    &.fingerprint {
      font-family: monospace;
      word-break: break-all;
    }
  }
}

.san-list {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  
  .san-tag {
    font-size: 11px;
  }
}

.validity-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.validity-item {
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
  }
}

.text-success {
  color: #10b981 !important;
}

.text-warning {
  color: #f59e0b !important;
}

.text-danger {
  color: #ef4444 !important;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
}
</style>
