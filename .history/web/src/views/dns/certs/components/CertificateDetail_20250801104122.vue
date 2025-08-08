<template>
  <el-drawer
    v-model="drawerVisible"
    title="证书详情"
    direction="rtl"
    size="50%"
  >
    <div v-if="certificate" class="cert-detail">
      <el-descriptions title="基本信息" :column="2" border>
        <el-descriptions-item label="证书名称">
          {{ certificate.common_name }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(certificate.status)">
            {{ getStatusText(certificate.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="CA类型">
          {{ certificate.ca_type }}
        </el-descriptions-item>
        <el-descriptions-item label="序列号">
          {{ certificate.serial_number || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="指纹">
          {{ certificate.fingerprint || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="自动续期">
          <el-tag :type="certificate.auto_renew ? 'success' : 'info'">
            {{ certificate.auto_renew ? '是' : '否' }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>

      <el-descriptions title="时间信息" :column="2" border style="margin-top: 20px;">
        <el-descriptions-item label="签发时间">
          {{ formatDate(certificate.issued_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="过期时间">
          {{ formatDate(certificate.expires_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="最后续期时间">
          {{ formatDate(certificate.last_renew_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="续期提前天数">
          {{ certificate.renew_days }} 天
        </el-descriptions-item>
      </el-descriptions>

      <el-descriptions title="域名信息" :column="1" border style="margin-top: 20px;">
        <el-descriptions-item label="主域名">
          {{ certificate.common_name }}
        </el-descriptions-item>
        <el-descriptions-item label="备用域名" v-if="certificate.subject_alt_names && certificate.subject_alt_names.length">
          <el-tag
            v-for="domain in certificate.subject_alt_names"
            :key="domain"
            style="margin: 2px;"
          >
            {{ domain }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>

      <div v-if="certificate.remark" style="margin-top: 20px;">
        <h4>备注信息</h4>
        <el-card shadow="never">
          {{ certificate.remark }}
        </el-card>
      </div>
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Certificate } from '@/types/dns'

interface Props {
  visible: boolean
  certificate?: Certificate | null
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
}

const props = withDefaults(defineProps<Props>(), {
  certificate: null
})

const emit = defineEmits<Emits>()

// 对话框可见性
const drawerVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

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

const formatDate = (date?: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}
</script>

<style scoped lang="scss">
.cert-detail {
  padding: 16px;
}
</style> 
