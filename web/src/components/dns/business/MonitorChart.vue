<template>
  <div class="monitor-chart">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ title }}</span>
          <el-select 
            v-model="timeRange" 
            size="small" 
            style="width: 120px;"
            @change="handleTimeRangeChange"
          >
            <el-option label="1小时" value="1h" />
            <el-option label="6小时" value="6h" />
            <el-option label="24小时" value="24h" />
            <el-option label="7天" value="7d" />
            <el-option label="30天" value="30d" />
          </el-select>
        </div>
      </template>

      <div class="chart-container" v-loading="loading">
        <div v-if="!chartData.length" class="empty-chart">
          <el-empty description="暂无监控数据" />
        </div>
        <div v-else>
          <!-- 这里应该集成实际的图表库，如 ECharts 或 Chart.js -->
          <!-- 为了演示，先用简单的表格展示数据 -->
          <div class="metric-summary">
            <el-row :gutter="16" style="margin-bottom: 16px;">
              <el-col :span="6">
                <el-statistic
                  title="平均响应时间"
                  :value="avgResponseTime"
                  suffix="ms"
                  :precision="2"
                />
              </el-col>
              <el-col :span="6">
                <el-statistic
                  title="成功率"
                  :value="successRate"
                  suffix="%"
                  :precision="2"
                />
              </el-col>
              <el-col :span="6">
                <el-statistic
                  title="总请求数"
                  :value="totalRequests"
                />
              </el-col>
              <el-col :span="6">
                <el-statistic
                  title="失败次数"
                  :value="failureCount"
                />
              </el-col>
            </el-row>
          </div>

          <!-- 简化的数据表格 -->
          <el-table 
            :data="recentData" 
            size="small" 
            max-height="300"
            style="margin-top: 16px;"
          >
            <el-table-column prop="timestamp" label="时间" width="180">
              <template #default="{ row }">
                {{ formatTime(row.timestamp) }}
              </template>
            </el-table-column>
            <el-table-column prop="response_time" label="响应时间(ms)" width="120">
              <template #default="{ row }">
                <span :class="getResponseTimeClass(row.response_time)">
                  {{ row.response_time }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="status_code" label="状态码" width="100">
              <template #default="{ row }">
                <el-tag 
                  :type="getStatusCodeType(row.status_code)" 
                  size="small"
                >
                  {{ row.status_code }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="ssl_days_left" label="SSL剩余天数" width="120" v-if="type === 'https'">
              <template #default="{ row }">
                <span :class="getSslDaysClass(row.ssl_days_left)">
                  {{ row.ssl_days_left }} 天
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="error_message" label="错误信息" show-overflow-tooltip>
              <template #default="{ row }">
                <span v-if="row.error_message" style="color: #f56c6c;">
                  {{ row.error_message }}
                </span>
                <span v-else style="color: #67c23a;">正常</span>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

interface MonitorData {
  timestamp: string
  response_time: number
  status_code: number
  ssl_days_left?: number
  error_message?: string
}

interface Props {
  title?: string
  type?: 'dns' | 'https' | 'ping'
  chartData: MonitorData[]
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: '监控图表',
  type: 'https',
  chartData: () => [],
  loading: false
})

const emit = defineEmits<{
  'time-range-change': [range: string]
}>()

// 响应式数据
const timeRange = ref('24h')

// 计算属性
const recentData = computed(() => {
  // 根据时间范围返回最近的数据，最多显示50条
  return props.chartData
    .slice(-50)
    .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
})

const avgResponseTime = computed(() => {
  if (!props.chartData.length) return 0
  const total = props.chartData.reduce((sum, item) => sum + item.response_time, 0)
  return total / props.chartData.length
})

const successRate = computed(() => {
  if (!props.chartData.length) return 0
  const successCount = props.chartData.filter(item => item.status_code >= 200 && item.status_code < 400).length
  return (successCount / props.chartData.length) * 100
})

const totalRequests = computed(() => props.chartData.length)

const failureCount = computed(() => {
  return props.chartData.filter(item => item.status_code >= 400 || item.error_message).length
})

// 事件处理
const handleTimeRangeChange = (range: string) => {
  emit('time-range-change', range)
}

// 工具方法
const formatTime = (timestamp: string) => {
  return new Date(timestamp).toLocaleString('zh-CN')
}

const getResponseTimeClass = (responseTime: number) => {
  if (responseTime > 3000) return 'text-red-500'
  if (responseTime > 1000) return 'text-orange-500'
  return 'text-green-500'
}

const getStatusCodeType = (statusCode: number) => {
  if (statusCode >= 200 && statusCode < 300) return 'success'
  if (statusCode >= 300 && statusCode < 400) return 'warning'
  if (statusCode >= 400 && statusCode < 500) return 'danger'
  if (statusCode >= 500) return 'danger'
  return 'info'
}

const getSslDaysClass = (days: number) => {
  if (days <= 7) return 'text-red-500'
  if (days <= 30) return 'text-orange-500'
  return 'text-green-500'
}
</script>

<style scoped lang="scss">
.monitor-chart {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .chart-container {
    min-height: 200px;
  }

  .empty-chart {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 200px;
  }

  .metric-summary {
    padding: 16px;
    background-color: #f8f9fa;
    border-radius: 8px;
  }

  .text-red-500 {
    color: #f56565;
    font-weight: bold;
  }

  .text-orange-500 {
    color: #ed8936;
    font-weight: bold;
  }

  .text-green-500 {
    color: #48bb78;
    font-weight: bold;
  }
}
</style>
