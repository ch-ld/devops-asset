<template>
  <div class="certificate-detail">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>证书详情</span>
          <div>
            <el-button 
              type="success" 
              size="small"
              @click="$emit('download')"
              v-if="certificate.status === 'issued'"
            >
              下载证书
            </el-button>
            <el-button 
              type="warning" 
              size="small"
              @click="$emit('renew')"
              v-if="certificate.status === 'issued' && needsRenewal"
            >
              续期证书
            </el-button>
            <el-button 
              type="danger" 
              size="small"
              @click="$emit('revoke')"
              v-if="certificate.status === 'issued'"
            >
              吊销证书
            </el-button>
          </div>
        </div>
      </template>
      
      <el-descriptions :column="2" border>
        <el-descriptions-item label="证书状态">
          <el-tag :type="getStatusType(certificate.status)">
            {{ getStatusText(certificate.status) }}
          </el-tag>
        </el-descriptions-item>
        
        <el-descriptions-item label="证书类型">
          {{ certificate.key_type }}
        </el-descriptions-item>
        
        <el-descriptions-item label="主域名">
          {{ certificate.common_name }}
        </el-descriptions-item>
        
        <el-descriptions-item label="备用域名">
          <el-tag
            v-for="domain in certificate.subject_alt_names"
            :key="domain"
            size="small"
            style="margin-right: 4px;"
          >
            {{ domain }}
          </el-tag>
          <span v-if="!certificate.subject_alt_names?.length">-</span>
        </el-descriptions-item>
        
        <el-descriptions-item label="颁发机构">
          {{ certificate.issuer || 'Let\'s Encrypt' }}
        </el-descriptions-item>
        
        <el-descriptions-item label="序列号">
          {{ certificate.serial_number || '-' }}
        </el-descriptions-item>
        
        <el-descriptions-item label="有效期开始">
          {{ formatDate(certificate.not_before) }}
        </el-descriptions-item>
        
        <el-descriptions-item label="有效期结束">
          <span :class="getExpiryClass(certificate.not_after)">
            {{ formatDate(certificate.not_after) }}
          </span>
        </el-descriptions-item>
        
        <el-descriptions-item label="剩余天数">
          <span :class="getDaysLeftClass(certificate.not_after)">
            {{ getDaysLeft(certificate.not_after) }} 天
          </span>
        </el-descriptions-item>
        
        <el-descriptions-item label="自动续期">
          <el-tag :type="certificate.auto_renew ? 'success' : 'info'" size="small">
            {{ certificate.auto_renew ? '已启用' : '未启用' }}
          </el-tag>
        </el-descriptions-item>
        
        <el-descriptions-item label="申请邮箱">
          {{ certificate.email }}
        </el-descriptions-item>
        
        <el-descriptions-item label="DNS提供商">
          {{ certificate.provider?.name || '-' }}
        </el-descriptions-item>
        
        <el-descriptions-item label="申请时间">
          {{ formatDate(certificate.created_at) }}
        </el-descriptions-item>
        
        <el-descriptions-item label="更新时间">
          {{ formatDate(certificate.updated_at) }}
        </el-descriptions-item>
        
        <el-descriptions-item label="备注" :span="2">
          {{ certificate.remark || '-' }}
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 部署信息 -->
    <el-card style="margin-top: 16px;" v-if="deployments?.length">
      <template #header>
        <div class="card-header">
          <span>部署信息</span>
          <el-button 
            type="primary" 
            size="small"
            @click="$emit('deploy')"
            v-if="certificate.status === 'issued'"
          >
            重新部署
          </el-button>
        </div>
      </template>
      
      <el-table :data="deployments" empty-text="暂无部署记录">
        <el-table-column prop="host_name" label="主机名称" />
        <el-table-column prop="host_ip" label="主机IP" />
        <el-table-column prop="deploy_path" label="部署路径" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="getDeployStatusType(row.status)" size="small">
              {{ getDeployStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="deployed_at" label="部署时间">
          <template #default="{ row }">
            {{ formatDate(row.deployed_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button 
              type="primary" 
              size="small" 
              text
              @click="$emit('redeploy', row)"
            >
              重新部署
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 操作日志 -->
    <el-card style="margin-top: 16px;" v-if="logs?.length">
      <template #header>
        <span>操作日志</span>
      </template>
      
      <el-timeline>
        <el-timeline-item
          v-for="log in logs"
          :key="log.id"
          :timestamp="formatDateTime(log.created_at)"
          :type="getLogType(log.action)"
        >
          <div>
            <strong>{{ getLogActionText(log.action) }}</strong>
            <div v-if="log.message" style="color: #666; font-size: 12px; margin-top: 4px;">
              {{ log.message }}
            </div>
            <div v-if="log.error" style="color: #f56c6c; font-size: 12px; margin-top: 4px;">
              错误：{{ log.error }}
            </div>
          </div>
        </el-timeline-item>
      </el-timeline>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Certificate, CertificateDeployment, CertificateLog } from '@/types/dns'

interface Props {
  certificate: Certificate
  deployments?: CertificateDeployment[]
  logs?: CertificateLog[]
}

const props = withDefaults(defineProps<Props>(), {
  deployments: () => [],
  logs: () => []
})

defineEmits<{
  download: []
  renew: []
  revoke: []
  deploy: []
  redeploy: [deployment: CertificateDeployment]
}>()

// 计算属性
const needsRenewal = computed(() => {
  return getDaysLeft(props.certificate.not_after) <= 30
})

// 工具方法
const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: 'warning',
    processing: 'info',
    issued: 'success',
    expired: 'danger',
    failed: 'danger',
    revoked: 'info'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待申请',
    processing: '申请中',
    issued: '已签发',
    expired: '已过期',
    failed: '申请失败',
    revoked: '已吊销'
  }
  return statusMap[status] || status
}

const getDeployStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    success: 'success',
    failed: 'danger',
    pending: 'warning'
  }
  return statusMap[status] || 'info'
}

const getDeployStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    success: '部署成功',
    failed: '部署失败',
    pending: '部署中'
  }
  return statusMap[status] || status
}

const getLogType = (action: string) => {
  const typeMap: Record<string, string> = {
    create: 'primary',
    renew: 'success',
    deploy: 'warning',
    revoke: 'danger'
  }
  return typeMap[action] || 'info'
}

const getLogActionText = (action: string) => {
  const actionMap: Record<string, string> = {
    create: '证书申请',
    renew: '证书续期',
    deploy: '证书部署',
    revoke: '证书吊销',
    download: '证书下载'
  }
  return actionMap[action] || action
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const formatDateTime = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const getDaysLeft = (expiresAt: string) => {
  if (!expiresAt) return 0
  const expires = new Date(expiresAt)
  const now = new Date()
  const diffTime = expires.getTime() - now.getTime()
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
}

const getDaysLeftClass = (expiresAt: string) => {
  const days = getDaysLeft(expiresAt)
  if (days <= 7) return 'text-red-500'
  if (days <= 30) return 'text-orange-500'
  return 'text-green-500'
}

const getExpiryClass = (expiresAt: string) => {
  const days = getDaysLeft(expiresAt)
  if (days <= 7) return 'text-red-500'
  if (days <= 30) return 'text-orange-500'
  return ''
}
</script>

<style scoped lang="scss">
.certificate-detail {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .text-red-500 {
    color: #f56565;
  }
  
  .text-orange-500 {
    color: #ed8936;
  }
  
  .text-green-500 {
    color: #48bb78;
  }
}
</style>
