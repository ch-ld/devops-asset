<template>
  <div class="host-detail">
    <el-descriptions :column="2" border>
      <el-descriptions-item label="主机名称">
        {{ host?.name || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="实例ID">
        {{ host?.instance_id || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="公网IP">
        <span v-if="Array.isArray(host?.public_ip)">
          <el-tag v-for="ip in host.public_ip" :key="ip" size="small" class="ip-tag">{{ ip }}</el-tag>
        </span>
        <span v-else>{{ host?.public_ip || '-' }}</span>
      </el-descriptions-item>
      <el-descriptions-item label="私网IP">
        <span v-if="Array.isArray(host?.private_ip)">
          <el-tag v-for="ip in host.private_ip" :key="ip" size="small" class="ip-tag">{{ ip }}</el-tag>
        </span>
        <span v-else>{{ host?.private_ip || '-' }}</span>
      </el-descriptions-item>
      <el-descriptions-item label="状态">
        <el-tag :type="getStatusType(host?.status)" size="small">
          {{ getStatusText(host?.status) }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="操作系统">
        {{ host?.os || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="云厂商">
        {{ host?.provider?.name || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="地域">
        {{ host?.region || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="CPU核数">
        {{ host?.configuration?.cpu_cores || '-' }}核
      </el-descriptions-item>
      <el-descriptions-item label="内存大小">
        {{ host?.configuration?.memory_size || '-' }}GB
      </el-descriptions-item>
      <el-descriptions-item label="实例类型">
        {{ host?.configuration?.instance_type || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="可用区">
        {{ host?.configuration?.zone_id || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="VPC ID">
        {{ host?.configuration?.vpc_id || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="用户名">
        {{ host?.username || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="过期时间">
        {{ host?.expired_at ? new Date(host.expired_at).toLocaleString() : '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="创建时间">
        {{ host?.created_at ? new Date(host.created_at).toLocaleString() : '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="更新时间">
        {{ host?.updated_at ? new Date(host.updated_at).toLocaleString() : '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="备注" :span="2">
        {{ host?.remark || '-' }}
      </el-descriptions-item>
    </el-descriptions>

    <!-- 标签信息 -->
    <div v-if="host?.tags && host.tags.length > 0" class="tags-section">
      <h4>标签信息</h4>
      <el-tag v-for="tag in host.tags" :key="tag" size="small" class="tag-item">
        {{ tag }}
      </el-tag>
    </div>

    <!-- 自定义字段 -->
    <div v-if="host?.extra_fields && Object.keys(host.extra_fields).length > 0" class="extra-fields-section">
      <h4>自定义字段</h4>
      <el-descriptions :column="2" border>
        <el-descriptions-item 
          v-for="(value, key) in host.extra_fields" 
          :key="key" 
          :label="key"
        >
          {{ value }}
        </el-descriptions-item>
      </el-descriptions>
    </div>

    <!-- 操作按钮 -->
    <div class="action-buttons">
      <el-button type="primary" @click="handleSSH" :disabled="host?.status !== 'running'">
        <i class="el-icon-monitor"></i>
        SSH连接
      </el-button>
      <el-button @click="handleSync" :loading="syncing">
        <i class="el-icon-refresh"></i>
        同步状态
      </el-button>
      <el-button @click="handleEdit">
        <i class="el-icon-edit"></i>
        编辑
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { syncHostStatus } from '@/api/system/host'
import type { Host } from '@/types/api/host'

interface Props {
  host: Host | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  edit: [host: Host]
  refresh: []
}>()

const syncing = ref(false)

// 状态相关方法
const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    running: 'success',
    stopped: 'warning',
    error: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    running: '运行中',
    stopped: '已停止',
    error: '异常'
  }
  return statusMap[status] || status
}

// SSH连接
const handleSSH = () => {
  if (!props.host) return
  window.open(`/terminal?host=${props.host.id}`, '_blank')
}

// 同步状态
const handleSync = async () => {
  if (!props.host) return
  
  try {
    syncing.value = true
    await syncHostStatus(props.host.id)
    ElMessage.success('同步成功')
    emit('refresh')
  } catch (error) {
    ElMessage.error('同步失败')
  } finally {
    syncing.value = false
  }
}

// 编辑
const handleEdit = () => {
  if (!props.host) return
  emit('edit', props.host)
}
</script>

<style scoped>
.host-detail {
  padding: 20px 0;
}

.ip-tag {
  margin-right: 8px;
  margin-bottom: 4px;
}

.tags-section,
.extra-fields-section {
  margin-top: 24px;
}

.tags-section h4,
.extra-fields-section h4 {
  margin: 0 0 12px 0;
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.tag-item {
  margin-right: 8px;
  margin-bottom: 4px;
}

.action-buttons {
  margin-top: 24px;
  text-align: center;
}

.action-buttons .el-button {
  margin: 0 8px;
}
</style>
