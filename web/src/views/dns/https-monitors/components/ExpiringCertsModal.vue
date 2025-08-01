<template>
  <a-modal
    :visible="visible"
    title="即将过期的SSL证书"
    :width="1000"
    :footer="null"
    @cancel="handleCancel"
  >
    <div class="expiring-certs-content">
      <a-alert
        :message="`发现 ${certificates.length} 个监控的SSL证书即将在30天内过期`"
        type="warning"
        show-icon
        style="margin-bottom: 16px"
      >
        <template #description>
          请及时关注证书状态，避免因证书过期导致服务中断。建议提前续期或更新证书。
        </template>
      </a-alert>

      <a-table
        :columns="columns"
        :data-source="certificates"
        :pagination="false"
        row-key="id"
        size="small"
      >
        <!-- 监控名称列 -->
        <template #name="{ record }">
          <div class="monitor-info">
            <div class="monitor-name">{{ record.name }}</div>
            <div class="monitor-url">{{ record.url }}</div>
          </div>
        </template>

        <!-- 证书信息列 -->
        <template #cert_info="{ record }">
          <div class="cert-info">
            <div class="cert-expire">
              <strong>过期时间：</strong>
              <span :class="getCertExpireClass(record.ssl_cert_expires_at)">
                {{ formatDate(record.ssl_cert_expires_at) }}
              </span>
            </div>
            <div class="cert-days" v-if="record.ssl_cert_expires_at">
              <strong>剩余天数：</strong>
              <a-tag :color="getDaysColor(getDaysRemaining(record.ssl_cert_expires_at))">
                {{ getDaysRemaining(record.ssl_cert_expires_at) }} 天
              </a-tag>
            </div>
            <div class="cert-issuer" v-if="record.ssl_cert_issuer">
              <strong>颁发机构：</strong>{{ record.ssl_cert_issuer }}
            </div>
          </div>
        </template>

        <!-- 状态列 -->
        <template #status="{ record }">
          <div class="status-info">
            <a-badge
              :status="getLastStatusBadge(record.last_status)"
              :text="getLastStatusText(record.last_status)"
            />
            <div class="last-check" v-if="record.last_checked">
              <small>{{ formatRelativeTime(record.last_checked) }}</small>
            </div>
          </div>
        </template>

        <!-- 操作列 -->
        <template #action="{ record }">
          <a-space>
            <a @click="handleViewMonitor(record)">
              <EyeOutlined />
              查看
            </a>
            <a @click="handleCheckNow(record)" :loading="checkingIds.includes(record.id)">
              <SyncOutlined :spin="checkingIds.includes(record.id)" />
              检查
            </a>
            <a :href="record.url" target="_blank">
              <LinkOutlined />
              访问
            </a>
          </a-space>
        </template>
      </a-table>

      <div class="modal-footer">
        <a-space>
          <a-button @click="handleRefreshAll" :loading="refreshing">
            <ReloadOutlined />
            刷新所有
          </a-button>
          <a-button @click="handleExportList">
            <ExportOutlined />
            导出列表
          </a-button>
          <a-button type="primary" @click="handleCancel">
            关闭
          </a-button>
        </a-space>
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { message } from 'ant-design-vue'
import {
  EyeOutlined,
  SyncOutlined,
  LinkOutlined,
  ReloadOutlined,
  ExportOutlined
} from '@ant-design/icons-vue'
import type { TableColumnsType } from 'ant-design-vue'
import { httpsMonitorApi } from '@/api/dns/httpsMonitor'
import type { HTTPSMonitor } from '@/types/dns'

interface Props {
  visible: boolean
  certificates: HTTPSMonitor[]
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'view-monitor', monitor: HTTPSMonitor): void
  (e: 'refresh'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 响应式数据
const checkingIds = ref<number[]>([])
const refreshing = ref(false)

// 表格列定义
const columns: TableColumnsType = [
  {
    title: '监控信息',
    dataIndex: 'name',
    key: 'name',
    slots: { customRender: 'name' },
    width: 250
  },
  {
    title: '证书信息',
    dataIndex: 'cert_info',
    key: 'cert_info',
    slots: { customRender: 'cert_info' },
    width: 300
  },
  {
    title: '监控状态',
    dataIndex: 'status',
    key: 'status',
    slots: { customRender: 'status' },
    width: 120
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'action' },
    width: 150
  }
]

// 计算属性
const sortedCertificates = computed(() => {
  return [...props.certificates].sort((a, b) => {
    const daysA = getDaysRemaining(a.ssl_cert_expires_at)
    const daysB = getDaysRemaining(b.ssl_cert_expires_at)
    return daysA - daysB
  })
})

// 方法
const getDaysRemaining = (expiresAt: string) => {
  if (!expiresAt) return 0
  const expireDate = new Date(expiresAt)
  const now = new Date()
  return Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
}

const getDaysColor = (days: number) => {
  if (days < 0) return 'red'
  if (days <= 7) return 'red'
  if (days <= 15) return 'orange'
  if (days <= 30) return 'gold'
  return 'green'
}

const getCertExpireClass = (expiresAt: string) => {
  const days = getDaysRemaining(expiresAt)
  if (days < 0) return 'text-red'
  if (days <= 7) return 'text-red'
  if (days <= 15) return 'text-orange'
  return 'text-warning'
}

const getLastStatusBadge = (status: string) => {
  const statusMap = {
    online: 'success',
    offline: 'error',
    warning: 'warning'
  }
  return statusMap[status] || 'default'
}

const getLastStatusText = (status: string) => {
  const statusMap = {
    online: '在线',
    offline: '离线',
    warning: '警告'
  }
  return statusMap[status] || status
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const formatRelativeTime = (date: string) => {
  if (!date) return ''
  const diff = Date.now() - new Date(date).getTime()
  const minutes = Math.floor(diff / (1000 * 60))
  
  if (minutes < 1) return '刚刚检查'
  if (minutes < 60) return `${minutes}分钟前检查`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前检查`
  const days = Math.floor(hours / 24)
  return `${days}天前检查`
}

// 事件处理
const handleViewMonitor = (monitor: HTTPSMonitor) => {
  emit('view-monitor', monitor)
  handleCancel()
}

const handleCheckNow = async (monitor: HTTPSMonitor) => {
  try {
    checkingIds.value.push(monitor.id)
    await httpsMonitorApi.check(monitor.id)
    message.success(`${monitor.name} 检查完成`)
    emit('refresh')
  } catch (error) {
    message.error(`${monitor.name} 检查失败`)
  } finally {
    checkingIds.value = checkingIds.value.filter(id => id !== monitor.id)
  }
}

const handleRefreshAll = async () => {
  try {
    refreshing.value = true
    
    // 批量检查所有即将过期的证书
    const batchSize = 3
    for (let i = 0; i < props.certificates.length; i += batchSize) {
      const batch = props.certificates.slice(i, i + batchSize)
      await Promise.all(batch.map(cert => httpsMonitorApi.check(cert.id)))
    }
    
    message.success('所有证书状态已刷新')
    emit('refresh')
  } catch (error) {
    message.error('刷新失败')
  } finally {
    refreshing.value = false
  }
}

const handleExportList = () => {
  try {
    const csvContent = [
      ['监控名称', 'URL', '证书过期时间', '剩余天数', '颁发机构', '监控状态'].join(','),
      ...sortedCertificates.value.map(cert => [
        cert.name,
        cert.url,
        formatDate(cert.ssl_cert_expires_at),
        getDaysRemaining(cert.ssl_cert_expires_at),
        cert.ssl_cert_issuer || '',
        getLastStatusText(cert.last_status)
      ].join(','))
    ].join('\n')
    
    const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    link.setAttribute('href', url)
    link.setAttribute('download', `即将过期证书列表_${new Date().toISOString().split('T')[0]}.csv`)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    
    message.success('证书列表已导出')
  } catch (error) {
    message.error('导出失败')
  }
}

const handleCancel = () => {
  emit('update:visible', false)
}
</script>

<style scoped lang="scss">
.expiring-certs-content {
  .monitor-info {
    .monitor-name {
      font-weight: 500;
      color: #262626;
      margin-bottom: 4px;
    }
    
    .monitor-url {
      font-size: 12px;
      color: #8c8c8c;
    }
  }
  
  .cert-info {
    div {
      margin-bottom: 4px;
      font-size: 12px;
      
      &:last-child {
        margin-bottom: 0;
      }
    }
    
    .cert-expire {
      .text-red {
        color: #ff4d4f;
        font-weight: 500;
      }
      
      .text-orange {
        color: #fa8c16;
        font-weight: 500;
      }
      
      .text-warning {
        color: #faad14;
        font-weight: 500;
      }
    }
  }
  
  .status-info {
    .last-check {
      margin-top: 4px;
      color: #8c8c8c;
    }
  }
  
  .modal-footer {
    margin-top: 16px;
    text-align: right;
    border-top: 1px solid #f0f0f0;
    padding-top: 16px;
  }
}
</style>
