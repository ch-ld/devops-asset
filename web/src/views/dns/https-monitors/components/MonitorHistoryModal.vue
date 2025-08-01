<template>
  <a-modal
    :visible="visible"
    :title="`${monitor?.name} - 监控历史`"
    :width="1200"
    :footer="null"
    @cancel="handleCancel"
  >
    <div class="monitor-history-content" v-if="monitor">
      <!-- 统计概览 -->
      <div class="history-stats">
        <a-row :gutter="16">
          <a-col :span="6">
            <a-statistic
              title="可用率"
              :value="stats.uptime"
              suffix="%"
              :precision="2"
              :value-style="{ color: stats.uptime >= 99 ? '#52c41a' : stats.uptime >= 95 ? '#fa8c16' : '#ff4d4f' }"
            />
          </a-col>
          <a-col :span="6">
            <a-statistic
              title="平均响应时间"
              :value="stats.avgResponseTime"
              suffix="ms"
              :value-style="{ color: '#1890ff' }"
            />
          </a-col>
          <a-col :span="6">
            <a-statistic
              title="总检查次数"
              :value="stats.totalChecks"
              :value-style="{ color: '#722ed1' }"
            />
          </a-col>
          <a-col :span="6">
            <a-statistic
              title="失败次数"
              :value="stats.failedChecks"
              :value-style="{ color: '#ff4d4f' }"
            />
          </a-col>
        </a-row>
      </div>

      <!-- 时间范围选择 -->
      <div class="time-range-selector">
        <a-space>
          <span>时间范围：</span>
          <a-radio-group v-model:value="timeRange" @change="handleTimeRangeChange">
            <a-radio-button value="1h">最近1小时</a-radio-button>
            <a-radio-button value="24h">最近24小时</a-radio-button>
            <a-radio-button value="7d">最近7天</a-radio-button>
            <a-radio-button value="30d">最近30天</a-radio-button>
            <a-radio-button value="custom">自定义</a-radio-button>
          </a-radio-group>
          <a-range-picker
            v-if="timeRange === 'custom'"
            v-model:value="customTimeRange"
            show-time
            @change="handleCustomTimeChange"
          />
          <a-button @click="handleRefresh" :loading="loading">
            <ReloadOutlined />
            刷新
          </a-button>
        </a-space>
      </div>

      <!-- 响应时间图表 -->
      <div class="chart-container">
        <a-card title="响应时间趋势" size="small">
          <div ref="responseTimeChart" class="chart" style="height: 300px;"></div>
        </a-card>
      </div>

      <!-- 可用性图表 -->
      <div class="chart-container">
        <a-card title="可用性状态" size="small">
          <div ref="uptimeChart" class="chart" style="height: 200px;"></div>
        </a-card>
      </div>

      <!-- 历史记录表格 -->
      <div class="history-table">
        <a-card title="检查历史" size="small">
          <template #extra>
            <a-space>
              <a-select
                v-model:value="statusFilter"
                placeholder="状态筛选"
                allow-clear
                style="width: 120px"
                @change="handleStatusFilterChange"
              >
                <a-select-option value="online">在线</a-select-option>
                <a-select-option value="offline">离线</a-select-option>
                <a-select-option value="warning">警告</a-select-option>
              </a-select>
              <a-button @click="handleExportHistory">
                <ExportOutlined />
                导出
              </a-button>
            </a-space>
          </template>

          <a-table
            :columns="historyColumns"
            :data-source="historyData"
            :loading="loading"
            :pagination="historyPagination"
            row-key="id"
            size="small"
            @change="handleHistoryTableChange"
          >
            <!-- 检查时间列 -->
            <template #checked_at="{ record }">
              <span>{{ formatDateTime(record.checked_at) }}</span>
            </template>

            <!-- 状态列 -->
            <template #status="{ record }">
              <a-badge
                :status="getStatusBadge(record.status)"
                :text="getStatusText(record.status)"
              />
            </template>

            <!-- 响应时间列 -->
            <template #response_time="{ record }">
              <span v-if="record.response_time" :class="getResponseTimeClass(record.response_time)">
                {{ record.response_time }}ms
              </span>
              <span v-else class="no-data">-</span>
            </template>

            <!-- HTTP状态码列 -->
            <template #http_status="{ record }">
              <a-tag
                v-if="record.http_status"
                :color="getHttpStatusColor(record.http_status)"
              >
                {{ record.http_status }}
              </a-tag>
              <span v-else class="no-data">-</span>
            </template>

            <!-- 错误信息列 -->
            <template #error_message="{ record }">
              <div v-if="record.error_message" class="error-message">
                <a-tooltip :title="record.error_message">
                  <ExclamationCircleOutlined style="color: #ff4d4f; margin-right: 4px;" />
                  {{ truncateText(record.error_message, 30) }}
                </a-tooltip>
              </div>
              <span v-else class="no-data">-</span>
            </template>

            <!-- SSL信息列 -->
            <template #ssl_info="{ record }">
              <div v-if="record.ssl_cert_expires_at" class="ssl-info">
                <div>过期: {{ formatDate(record.ssl_cert_expires_at) }}</div>
                <div v-if="record.ssl_cert_issuer" class="ssl-issuer">
                  {{ truncateText(record.ssl_cert_issuer, 20) }}
                </div>
              </div>
              <span v-else class="no-data">-</span>
            </template>
          </a-table>
        </a-card>
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { message } from 'ant-design-vue'
import {
  ReloadOutlined,
  ExportOutlined,
  ExclamationCircleOutlined
} from '@ant-design/icons-vue'
import type { TableColumnsType } from 'ant-design-vue'
import type { Dayjs } from 'dayjs'
import * as echarts from 'echarts'
import type { HTTPSMonitor } from '@/types/dns'

interface Props {
  visible: boolean
  monitor?: HTTPSMonitor | null
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
}

interface HistoryRecord {
  id: number
  checked_at: string
  status: 'online' | 'offline' | 'warning'
  response_time?: number
  http_status?: number
  error_message?: string
  ssl_cert_expires_at?: string
  ssl_cert_issuer?: string
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const timeRange = ref('24h')
const customTimeRange = ref<[Dayjs, Dayjs] | null>(null)
const statusFilter = ref<string | undefined>(undefined)
const historyData = ref<HistoryRecord[]>([])
const responseTimeChart = ref<HTMLElement>()
const uptimeChart = ref<HTMLElement>()

// 图表实例
let responseTimeChartInstance: echarts.ECharts | null = null
let uptimeChartInstance: echarts.ECharts | null = null

// 统计数据
const stats = reactive({
  uptime: 0,
  avgResponseTime: 0,
  totalChecks: 0,
  failedChecks: 0
})

// 分页配置
const historyPagination = reactive({
  current: 1,
  pageSize: 50,
  total: 0,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number) => `共 ${total} 条记录`
})

// 表格列定义
const historyColumns: TableColumnsType = [
  {
    title: '检查时间',
    dataIndex: 'checked_at',
    key: 'checked_at',
    slots: { customRender: 'checked_at' },
    width: 150,
    sorter: true
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    slots: { customRender: 'status' },
    width: 80,
    filters: [
      { text: '在线', value: 'online' },
      { text: '离线', value: 'offline' },
      { text: '警告', value: 'warning' }
    ]
  },
  {
    title: '响应时间',
    dataIndex: 'response_time',
    key: 'response_time',
    slots: { customRender: 'response_time' },
    width: 100,
    sorter: true
  },
  {
    title: 'HTTP状态',
    dataIndex: 'http_status',
    key: 'http_status',
    slots: { customRender: 'http_status' },
    width: 100
  },
  {
    title: 'SSL信息',
    dataIndex: 'ssl_info',
    key: 'ssl_info',
    slots: { customRender: 'ssl_info' },
    width: 150
  },
  {
    title: '错误信息',
    dataIndex: 'error_message',
    key: 'error_message',
    slots: { customRender: 'error_message' },
    ellipsis: true
  }
]

// 方法
const getStatusBadge = (status: string) => {
  const statusMap = {
    online: 'success',
    offline: 'error',
    warning: 'warning'
  }
  return statusMap[status] || 'default'
}

const getStatusText = (status: string) => {
  const statusMap = {
    online: '在线',
    offline: '离线',
    warning: '警告'
  }
  return statusMap[status] || status
}

const getResponseTimeClass = (responseTime: number) => {
  if (responseTime < 500) return 'response-fast'
  if (responseTime < 2000) return 'response-normal'
  return 'response-slow'
}

const getHttpStatusColor = (status: number) => {
  if (status >= 200 && status < 300) return 'green'
  if (status >= 300 && status < 400) return 'blue'
  if (status >= 400 && status < 500) return 'orange'
  if (status >= 500) return 'red'
  return 'default'
}

const formatDateTime = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const truncateText = (text: string, maxLength: number) => {
  if (!text) return ''
  return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
}

// 初始化图表
const initCharts = async () => {
  await nextTick()
  
  if (responseTimeChart.value) {
    responseTimeChartInstance = echarts.init(responseTimeChart.value)
    updateResponseTimeChart()
  }
  
  if (uptimeChart.value) {
    uptimeChartInstance = echarts.init(uptimeChart.value)
    updateUptimeChart()
  }
}

// 更新响应时间图表
const updateResponseTimeChart = () => {
  if (!responseTimeChartInstance) return
  
  // 模拟数据，实际应该从API获取
  const data = historyData.value
    .filter(record => record.response_time)
    .map(record => ({
      time: record.checked_at,
      value: record.response_time
    }))
    .slice(-100) // 最近100条记录
  
  const option = {
    title: {
      text: '响应时间趋势',
      left: 'center',
      textStyle: { fontSize: 14 }
    },
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const point = params[0]
        return `${point.name}<br/>响应时间: ${point.value}ms`
      }
    },
    xAxis: {
      type: 'category',
      data: data.map(item => new Date(item.time).toLocaleTimeString('zh-CN'))
    },
    yAxis: {
      type: 'value',
      name: '响应时间 (ms)'
    },
    series: [{
      data: data.map(item => item.value),
      type: 'line',
      smooth: true,
      itemStyle: { color: '#1890ff' },
      areaStyle: { opacity: 0.3 }
    }]
  }
  
  responseTimeChartInstance.setOption(option)
}

// 更新可用性图表
const updateUptimeChart = () => {
  if (!uptimeChartInstance) return
  
  // 模拟数据，实际应该从API获取
  const data = historyData.value.slice(-50).map(record => ({
    time: record.checked_at,
    status: record.status
  }))
  
  const option = {
    title: {
      text: '可用性状态',
      left: 'center',
      textStyle: { fontSize: 14 }
    },
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const point = params[0]
        const statusText = getStatusText(point.value === 1 ? 'online' : point.value === 0.5 ? 'warning' : 'offline')
        return `${point.name}<br/>状态: ${statusText}`
      }
    },
    xAxis: {
      type: 'category',
      data: data.map(item => new Date(item.time).toLocaleTimeString('zh-CN'))
    },
    yAxis: {
      type: 'value',
      min: 0,
      max: 1,
      axisLabel: {
        formatter: (value: number) => {
          if (value === 1) return '在线'
          if (value === 0.5) return '警告'
          return '离线'
        }
      }
    },
    series: [{
      data: data.map(item => {
        if (item.status === 'online') return 1
        if (item.status === 'warning') return 0.5
        return 0
      }),
      type: 'line',
      step: 'end',
      itemStyle: { color: '#52c41a' },
      lineStyle: { width: 3 }
    }]
  }
}

// 事件处理
const handleTimeRangeChange = () => {
  fetchHistoryData()
}

const handleCustomTimeChange = () => {
  if (customTimeRange.value) {
    fetchHistoryData()
  }
}

const handleStatusFilterChange = () => {
  fetchHistoryData()
}

const handleRefresh = () => {
  fetchHistoryData()
  fetchStats()
}

const handleHistoryTableChange = (pagination: any) => {
  historyPagination.current = pagination.current
  historyPagination.pageSize = pagination.pageSize
  fetchHistoryData()
}

const handleExportHistory = () => {
  try {
    const csvContent = [
      ['检查时间', '状态', '响应时间(ms)', 'HTTP状态', 'SSL过期时间', '错误信息'].join(','),
      ...historyData.value.map(record => [
        formatDateTime(record.checked_at),
        getStatusText(record.status),
        record.response_time || '',
        record.http_status || '',
        formatDate(record.ssl_cert_expires_at || ''),
        (record.error_message || '').replace(/,/g, '；')
      ].join(','))
    ].join('\n')
    
    const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    link.setAttribute('href', url)
    link.setAttribute('download', `${props.monitor?.name}_监控历史_${new Date().toISOString().split('T')[0]}.csv`)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    
    message.success('监控历史已导出')
  } catch (error) {
    message.error('导出失败')
  }
}

// 数据获取
const fetchHistoryData = async () => {
  if (!props.monitor) return
  
  try {
    loading.value = true
    
    // 这里应该调用实际的API获取历史数据
    // 暂时使用模拟数据
    const mockData: HistoryRecord[] = []
    const now = new Date()
    
    for (let i = 0; i < 200; i++) {
      const time = new Date(now.getTime() - i * 5 * 60 * 1000) // 每5分钟一条记录
      const isOnline = Math.random() > 0.1 // 90%在线率
      const responseTime = isOnline ? Math.floor(Math.random() * 2000) + 100 : undefined
      
      mockData.push({
        id: i,
        checked_at: time.toISOString(),
        status: isOnline ? (responseTime && responseTime > 1500 ? 'warning' : 'online') : 'offline',
        response_time: responseTime,
        http_status: isOnline ? (Math.random() > 0.05 ? 200 : 500) : undefined,
        error_message: isOnline ? undefined : 'Connection timeout',
        ssl_cert_expires_at: props.monitor.ssl_cert_expires_at,
        ssl_cert_issuer: props.monitor.ssl_cert_issuer
      })
    }
    
    historyData.value = mockData
    historyPagination.total = mockData.length
    
    // 更新图表
    updateResponseTimeChart()
    updateUptimeChart()
  } catch (error) {
    message.error('获取历史数据失败')
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  if (!props.monitor) return
  
  try {
    // 这里应该调用实际的API获取统计数据
    // 暂时使用模拟数据
    const onlineCount = historyData.value.filter(record => record.status === 'online').length
    const totalCount = historyData.value.length
    const responseTimes = historyData.value
      .filter(record => record.response_time)
      .map(record => record.response_time!)
    
    stats.uptime = totalCount > 0 ? (onlineCount / totalCount) * 100 : 0
    stats.avgResponseTime = responseTimes.length > 0 
      ? Math.round(responseTimes.reduce((sum, time) => sum + time, 0) / responseTimes.length)
      : 0
    stats.totalChecks = totalCount
    stats.failedChecks = totalCount - onlineCount
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const handleCancel = () => {
  emit('update:visible', false)
}

// 生命周期
onMounted(() => {
  if (props.visible && props.monitor) {
    fetchHistoryData()
    fetchStats()
    initCharts()
  }
})

onUnmounted(() => {
  if (responseTimeChartInstance) {
    responseTimeChartInstance.dispose()
  }
  if (uptimeChartInstance) {
    uptimeChartInstance.dispose()
  }
})

// 监听visible变化
watch(() => props.visible, (visible) => {
  if (visible && props.monitor) {
    fetchHistoryData()
    fetchStats()
    nextTick(() => {
      initCharts()
    })
  }
})
</script>

<style scoped lang="scss">
.monitor-history-content {
  .history-stats {
    margin-bottom: 24px;
    padding: 16px;
    background: #fafafa;
    border-radius: 6px;
  }
  
  .time-range-selector {
    margin-bottom: 24px;
    padding: 16px;
    background: #fff;
    border: 1px solid #f0f0f0;
    border-radius: 6px;
  }
  
  .chart-container {
    margin-bottom: 24px;
    
    .chart {
      width: 100%;
    }
  }
  
  .history-table {
    .error-message {
      display: flex;
      align-items: center;
    }
    
    .ssl-info {
      font-size: 12px;
      
      .ssl-issuer {
        color: #8c8c8c;
        margin-top: 2px;
      }
    }
    
    .response-fast {
      color: #52c41a;
      font-weight: 500;
    }
    
    .response-normal {
      color: #1890ff;
    }
    
    .response-slow {
      color: #ff4d4f;
      font-weight: 500;
    }
    
    .no-data {
      color: #d9d9d9;
      font-style: italic;
    }
  }
}
</style>
