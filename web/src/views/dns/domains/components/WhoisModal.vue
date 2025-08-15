<template>
  <el-dialog
    :model-value="visible"
    title="WHOIS查询"
    width="800px"
    :before-close="handleClose"
    class="whois-modal"
  >
    <div class="whois-container">
      <!-- 查询头部 -->
      <div class="query-header">
        <div class="domain-info">
          <div class="domain-icon">
            <el-icon size="24"><Monitor /></el-icon>
          </div>
          <div class="domain-content">
            <h3 class="domain-name">{{ domain?.name }}</h3>
            <p class="domain-desc">域名WHOIS信息查询</p>
          </div>
        </div>
        <div class="query-actions">
          <el-button 
            type="primary" 
            @click="handleQuery"
            :loading="queryLoading"
            class="query-btn"
          >
            <el-icon><Search /></el-icon>
            {{ queryLoading ? '查询中...' : '重新查询' }}
          </el-button>
        </div>
      </div>

      <!-- 查询结果 -->
      <div class="whois-content">
        <el-card v-if="whoisData" class="result-card">
          <template #header>
            <div class="card-header">
              <span class="header-title">查询结果</span>
              <div class="header-actions">
                <el-tag :type="getStatusType(whoisData.status)" effect="light">
                  {{ getStatusText(whoisData.status) }}
                </el-tag>
                <span class="query-time">{{ formatTime(whoisData.query_time) }}</span>
              </div>
            </div>
          </template>

          <!-- 基本信息 -->
          <div class="info-section">
            <h4 class="section-title">基本信息</h4>
            <div class="info-grid">
              <div class="info-item">
                <label>域名状态</label>
                <span>{{ whoisData.domain_status || '-' }}</span>
              </div>
              <div class="info-item">
                <label>注册商</label>
                <span>{{ whoisData.registrar || '-' }}</span>
              </div>
              <div class="info-item">
                <label>注册时间</label>
                <span>{{ formatDate(whoisData.creation_date) }}</span>
              </div>
              <div class="info-item">
                <label>更新时间</label>
                <span>{{ formatDate(whoisData.updated_date) }}</span>
              </div>
              <div class="info-item">
                <label>过期时间</label>
                <span :class="getExpirationClass(whoisData.expiration_date)">
                  {{ formatDate(whoisData.expiration_date) }}
                </span>
              </div>
              <div class="info-item">
                <label>DNS服务器</label>
                <div v-if="whoisData.name_servers && whoisData.name_servers.length > 0" class="dns-servers">
                  <el-tag 
                    v-for="server in whoisData.name_servers" 
                    :key="server" 
                    size="small"
                    type="info"
                    class="dns-tag"
                  >
                    {{ server }}
                  </el-tag>
                </div>
                <span v-else>-</span>
              </div>
            </div>
          </div>

          <!-- 联系人信息 -->
          <div class="info-section" v-if="hasContactInfo">
            <h4 class="section-title">联系人信息</h4>
            <el-tabs type="border-card" class="contact-tabs">
              <el-tab-pane label="注册人" name="registrant" v-if="whoisData.registrant">
                <div class="contact-info">
                  <div class="contact-item" v-if="whoisData.registrant.name">
                    <label>姓名</label>
                    <span>{{ whoisData.registrant.name }}</span>
                  </div>
                  <div class="contact-item" v-if="whoisData.registrant.organization">
                    <label>组织</label>
                    <span>{{ whoisData.registrant.organization }}</span>
                  </div>
                  <div class="contact-item" v-if="whoisData.registrant.email">
                    <label>邮箱</label>
                    <span>{{ whoisData.registrant.email }}</span>
                  </div>
                  <div class="contact-item" v-if="whoisData.registrant.country">
                    <label>国家</label>
                    <span>{{ whoisData.registrant.country }}</span>
                  </div>
                </div>
              </el-tab-pane>
              <el-tab-pane label="管理员" name="admin" v-if="whoisData.admin">
                <div class="contact-info">
                  <div class="contact-item" v-if="whoisData.admin.name">
                    <label>姓名</label>
                    <span>{{ whoisData.admin.name }}</span>
                  </div>
                  <div class="contact-item" v-if="whoisData.admin.email">
                    <label>邮箱</label>
                    <span>{{ whoisData.admin.email }}</span>
                  </div>
                </div>
              </el-tab-pane>
              <el-tab-pane label="技术联系人" name="tech" v-if="whoisData.tech">
                <div class="contact-info">
                  <div class="contact-item" v-if="whoisData.tech.name">
                    <label>姓名</label>
                    <span>{{ whoisData.tech.name }}</span>
                  </div>
                  <div class="contact-item" v-if="whoisData.tech.email">
                    <label>邮箱</label>
                    <span>{{ whoisData.tech.email }}</span>
                  </div>
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>

          <!-- 原始数据 -->
          <div class="info-section">
            <h4 class="section-title">
              原始WHOIS数据
              <el-button 
                type="primary" 
                size="small" 
                text 
                @click="copyRawData"
                class="copy-btn"
              >
                <el-icon><CopyDocument /></el-icon>
                复制
              </el-button>
            </h4>
            <div class="raw-data">
              <pre>{{ whoisData.raw_data || '暂无原始数据' }}</pre>
            </div>
          </div>
        </el-card>

        <!-- 加载状态 -->
        <div v-else-if="queryLoading" class="loading-container">
          <el-icon class="loading-icon"><Loading /></el-icon>
          <p>正在查询WHOIS信息...</p>
        </div>

        <!-- 空状态 -->
        <div v-else class="empty-container">
          <el-icon class="empty-icon"><Search /></el-icon>
          <p>点击"重新查询"按钮获取WHOIS信息</p>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose" class="close-btn">
          关闭
        </el-button>
        <el-button 
          type="primary" 
          @click="exportWhois"
          :disabled="!whoisData"
          class="export-btn"
        >
          <el-icon><Download /></el-icon>
          导出报告
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Monitor, 
  Search, 
  Loading, 
  CopyDocument, 
  Download 
} from '@element-plus/icons-vue'
import { domainApi } from '@/api/dns/domain'
import type { Domain } from '@/types/dns'

interface Props {
  visible: boolean
  domain?: Domain | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
}

interface WhoisData {
  domain: string
  status: string
  registrar: string
  creation_date: string
  updated_date: string
  expiration_date: string
  name_servers: string[]
  domain_status: string
  registrant?: {
    name?: string
    organization?: string
    email?: string
    country?: string
  }
  admin?: {
    name?: string
    email?: string
  }
  tech?: {
    name?: string
    email?: string
  }
  raw_data: string
  query_time: string
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 响应式数据
const queryLoading = ref(false)
const whoisData = ref<WhoisData | null>(null)

// 计算属性
const hasContactInfo = computed(() => {
  return whoisData.value && (
    whoisData.value.registrant || 
    whoisData.value.admin || 
    whoisData.value.tech
  )
})

// 方法
const handleClose = () => {
  emit('update:visible', false)
}

const handleQuery = async () => {
  if (!props.domain) return
  
  try {
    queryLoading.value = true
    const response = await domainApi.whois(props.domain.id)
    whoisData.value = response.data
    ElMessage.success('查询成功')
  } catch (error) {
    ElMessage.error('查询失败')
    whoisData.value = null
  } finally {
    queryLoading.value = false
  }
}

const getStatusType = (status: string) => {
  const typeMap: Record<string, any> = {
    'active': 'success',
    'inactive': 'danger',
    'pending': 'warning',
    'expired': 'danger'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    'active': '正常',
    'inactive': '非活跃',
    'pending': '待处理',
    'expired': '已过期'
  }
  return textMap[status] || status
}

const formatDate = (date: string) => {
  if (!date) return '-'
  try {
    return new Date(date).toLocaleString('zh-CN')
  } catch {
    return date
  }
}

const formatTime = (time: string) => {
  if (!time) return ''
  try {
    return new Date(time).toLocaleString('zh-CN')
  } catch {
    return time
  }
}

const getExpirationClass = (date: string) => {
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

const copyRawData = async () => {
  if (!whoisData.value?.raw_data) return
  
  try {
    await navigator.clipboard.writeText(whoisData.value.raw_data)
    ElMessage.success('复制成功')
  } catch {
    ElMessage.error('复制失败')
  }
}

const exportWhois = async () => {
  if (!whoisData.value || !props.domain) return
  
  try {
    const blob = await domainApi.exportWhois(props.domain.id)
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${props.domain.name}-whois-${new Date().toISOString().split('T')[0]}.txt`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

// 监听弹窗显示状态
watch(() => props.visible, (visible) => {
  if (visible && props.domain) {
    // 自动查询
    handleQuery()
  } else if (!visible) {
    // 清空数据
    whoisData.value = null
  }
})
</script>

<style scoped lang="scss">
.whois-modal {
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

.whois-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.query-header {
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
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
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

.query-actions {
  .query-btn {
    background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
    border: none;
    border-radius: 8px;
    padding: 12px 24px;
    font-weight: 600;
  }
}

.whois-content {
  min-height: 400px;
}

.result-card {
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
  
  .header-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    
    .query-time {
      font-size: 12px;
      color: #6b7280;
    }
  }
}

.info-section {
  margin-bottom: 32px;
  
  &:last-child {
    margin-bottom: 0;
  }
  
  .section-title {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin: 0 0 16px 0;
    font-size: 14px;
    font-weight: 600;
    color: #374151;
    padding-bottom: 8px;
    border-bottom: 1px solid #e5e7eb;
    
    .copy-btn {
      font-size: 12px;
      padding: 4px 8px;
    }
  }
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  
  label {
    font-size: 12px;
    font-weight: 500;
    color: #6b7280;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  
  span {
    font-size: 14px;
    color: #1f2937;
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

.dns-servers {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 4px;
  
  .dns-tag {
    font-size: 11px;
    padding: 2px 8px;
  }
}

.contact-tabs {
  :deep(.el-tabs__header) {
    margin-bottom: 16px;
  }
  
  :deep(.el-tabs__content) {
    padding: 16px;
  }
}

.contact-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.contact-item {
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
  }
}

.raw-data {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 16px;
  max-height: 300px;
  overflow-y: auto;
  
  pre {
    margin: 0;
    font-family: 'Courier New', monospace;
    font-size: 12px;
    line-height: 1.5;
    color: #374151;
    white-space: pre-wrap;
    word-break: break-all;
  }
}

.loading-container,
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #6b7280;
  
  .loading-icon,
  .empty-icon {
    font-size: 48px;
    margin-bottom: 16px;
    color: #9ca3af;
  }
  
  .loading-icon {
    animation: spin 1s linear infinite;
  }
  
  p {
    margin: 0;
    font-size: 14px;
  }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
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
    background: linear-gradient(135deg, #10b981 0%, #059669 100%);
    border: none;
    font-weight: 600;
  }
}
</style>
