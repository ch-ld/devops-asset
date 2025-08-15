<template>
  <div class="domain-detail">
    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>域名信息</span>
          <el-button 
            type="primary" 
            size="small"
            @click="$emit('edit')"
            v-if="!readonly"
          >
            编辑
          </el-button>
        </div>
      </template>
      
      <el-descriptions :column="2" border>
        <el-descriptions-item label="域名">
          {{ domain.name }}
        </el-descriptions-item>
        
        <el-descriptions-item label="注册商">
          {{ getRegistrarName(domain.registrar_type) }}
        </el-descriptions-item>
        
        <el-descriptions-item label="分组">
          <el-tag 
            v-if="domain.group"
            :color="domain.group.color"
            size="small"
          >
            {{ domain.group.name }}
          </el-tag>
          <span v-else>-</span>
        </el-descriptions-item>
        
        <el-descriptions-item label="标签">
          <el-tag
            v-for="tag in domain.tags"
            :key="tag.id"
            :color="tag.background_color"
            :style="{ color: tag.color }"
            size="small"
            style="margin-right: 4px;"
          >
            {{ tag.name }}
          </el-tag>
          <span v-if="!domain.tags?.length">-</span>
        </el-descriptions-item>
        
        <el-descriptions-item label="自动续期">
          <el-tag :type="domain.auto_renew ? 'success' : 'info'" size="small">
            {{ domain.auto_renew ? '已启用' : '未启用' }}
          </el-tag>
        </el-descriptions-item>
        
        <el-descriptions-item label="过期时间">
          <span :class="getExpiryClass(domain.expires_at)">
            {{ formatDate(domain.expires_at) }}
          </span>
          <el-tag 
            v-if="getDaysUntilExpiry(domain.expires_at) <= 30"
            :type="getDaysUntilExpiry(domain.expires_at) <= 7 ? 'danger' : 'warning'"
            size="small"
            style="margin-left: 8px;"
          >
            {{ getDaysUntilExpiry(domain.expires_at) }}天后过期
          </el-tag>
        </el-descriptions-item>
        
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(domain.status)" size="small">
            {{ getStatusText(domain.status) }}
          </el-tag>
        </el-descriptions-item>
        
        <el-descriptions-item label="创建时间">
          {{ formatDate(domain.created_at) }}
        </el-descriptions-item>
        
        <el-descriptions-item label="备注" span="2">
          {{ domain.remark || '-' }}
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 证书列表 -->
    <el-card class="certificates-card" style="margin-top: 16px;">
      <template #header>
        <div class="card-header">
          <span>关联证书</span>
          <el-button 
            type="primary" 
            size="small"
            @click="$emit('create-certificate')"
            v-if="!readonly"
          >
            申请证书
          </el-button>
        </div>
      </template>
      
      <el-table 
        :data="certificates" 
        v-loading="certificatesLoading"
        empty-text="暂无证书"
      >
        <el-table-column prop="common_name" label="主域名" />
        <el-table-column prop="subject_alt_names" label="备用域名">
          <template #default="{ row }">
            <el-tag
              v-for="domain in row.subject_alt_names"
              :key="domain"
              size="small"
              style="margin-right: 4px;"
            >
              {{ domain }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="getCertStatusType(row.status)" size="small">
              {{ getCertStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="not_after" label="过期时间">
          <template #default="{ row }">
            <span :class="getExpiryClass(row.not_after)">
              {{ formatDate(row.not_after) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button 
              type="primary" 
              size="small" 
              text
              @click="$emit('view-certificate', row)"
            >
              查看
            </el-button>
            <el-button 
              type="success" 
              size="small" 
              text
              @click="$emit('download-certificate', row)"
              v-if="row.status === 'issued'"
            >
              下载
            </el-button>
            <el-button 
              type="warning" 
              size="small" 
              text
              @click="$emit('renew-certificate', row)"
              v-if="row.status === 'issued' && getDaysUntilExpiry(row.not_after) <= 30"
            >
              续期
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- DNS记录 -->
    <el-card class="records-card" style="margin-top: 16px;">
      <template #header>
        <div class="card-header">
          <span>DNS记录</span>
          <div>
            <el-button 
              type="default" 
              size="small"
              @click="$emit('sync-records')"
              :loading="syncLoading"
            >
              同步记录
            </el-button>
            <el-button 
              type="primary" 
              size="small"
              @click="$emit('create-record')"
              v-if="!readonly"
            >
              添加记录
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table 
        :data="records" 
        v-loading="recordsLoading"
        empty-text="暂无记录"
      >
        <el-table-column prop="name" label="记录名" />
        <el-table-column prop="type" label="类型" width="80" />
        <el-table-column prop="value" label="记录值" show-overflow-tooltip />
        <el-table-column prop="ttl" label="TTL" width="100" />
        <el-table-column prop="sync_status" label="同步状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getSyncStatusType(row.sync_status)" size="small">
              {{ getSyncStatusText(row.sync_status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button 
              type="primary" 
              size="small" 
              text
              @click="$emit('edit-record', row)"
              v-if="!readonly"
            >
              编辑
            </el-button>
            <el-button 
              type="danger" 
              size="small" 
              text
              @click="$emit('delete-record', row)"
              v-if="!readonly"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Domain, Certificate, DnsRecord } from '@/types/dns'

interface Props {
  domain: Domain
  certificates?: Certificate[]
  records?: DnsRecord[]
  certificatesLoading?: boolean
  recordsLoading?: boolean
  syncLoading?: boolean
  readonly?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  certificates: () => [],
  records: () => [],
  certificatesLoading: false,
  recordsLoading: false,
  syncLoading: false,
  readonly: false
})

defineEmits<{
  edit: []
  'create-certificate': []
  'view-certificate': [certificate: Certificate]
  'download-certificate': [certificate: Certificate]
  'renew-certificate': [certificate: Certificate]
  'sync-records': []
  'create-record': []
  'edit-record': [record: DnsRecord]
  'delete-record': [record: DnsRecord]
}>()

// 工具方法
const getRegistrarName = (type: string) => {
  const names: Record<string, string> = {
    aliyun: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS',
    godaddy: 'GoDaddy',
    other: '其他'
  }
  return names[type] || type
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    active: 'success',
    inactive: 'info',
    expired: 'danger',
    expiring: 'warning'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '正常',
    inactive: '停用',
    expired: '已过期',
    expiring: '即将过期'
  }
  return statusMap[status] || status
}

const getCertStatusType = (status: string) => {
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

const getCertStatusText = (status: string) => {
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

const getSyncStatusType = (syncStatus: string) => {
  const statusMap: Record<string, string> = {
    synced: 'success',
    pending: 'warning',
    failed: 'danger'
  }
  return statusMap[syncStatus] || 'info'
}

const getSyncStatusText = (syncStatus: string) => {
  const statusMap: Record<string, string> = {
    synced: '已同步',
    pending: '待同步',
    failed: '同步失败'
  }
  return statusMap[syncStatus] || syncStatus
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const getDaysUntilExpiry = (expiryDate: string) => {
  if (!expiryDate) return 0
  const expiry = new Date(expiryDate)
  const now = new Date()
  const diffTime = expiry.getTime() - now.getTime()
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
}

const getExpiryClass = (expiryDate: string) => {
  const days = getDaysUntilExpiry(expiryDate)
  if (days <= 7) return 'text-red-500'
  if (days <= 30) return 'text-orange-500'
  return ''
}
</script>

<style scoped lang="scss">
.domain-detail {
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
}
</style>
