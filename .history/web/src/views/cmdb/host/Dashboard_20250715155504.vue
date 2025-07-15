<template>
  <div class="host-dashboard-container">
    <a-row :gutter="[16, 16]">
      <!-- 统计卡片 -->
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card" :bordered="false">
          <div class="stat-content">
            <div class="stat-icon" style="background-color: #1890ff;">
              <desktop-outlined />
            </div>
            <div class="stat-info">
              <div class="stat-title">主机总数</div>
              <div class="stat-value">{{ dashboardData.totalHosts || 0 }}</div>
            </div>
          </div>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card" :bordered="false">
          <div class="stat-content">
            <div class="stat-icon" style="background-color: #52c41a;">
              <check-circle-outlined />
            </div>
            <div class="stat-info">
              <div class="stat-title">正常运行</div>
              <div class="stat-value">{{ dashboardData.runningHosts || 0 }}</div>
            </div>
          </div>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card" :bordered="false">
          <div class="stat-content">
            <div class="stat-icon" style="background-color: #f5222d;">
              <warning-outlined />
            </div>
            <div class="stat-info">
              <div class="stat-title">异常主机</div>
              <div class="stat-value">{{ dashboardData.errorHosts || 0 }}</div>
            </div>
          </div>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card" :bordered="false">
          <div class="stat-content">
            <div class="stat-icon" style="background-color: #faad14;">
              <clock-circle-outlined />
            </div>
            <div class="stat-info">
              <div class="stat-title">即将到期</div>
              <div class="stat-value">{{ dashboardData.expiringHosts || 0 }}</div>
            </div>
          </div>
        </a-card>
      </a-col>

      <!-- 状态分布 -->
      <a-col :xs="24" :sm="24" :md="12">
        <a-card title="主机状态分布" :bordered="false">
          <template #extra>
            <a-button type="link" size="small" @click="refreshData">
              <reload-outlined />
              刷新
            </a-button>
          </template>
          <div id="statusChart" ref="statusChartRef" style="height: 300px"></div>
        </a-card>
      </a-col>

      <!-- 地区分布 -->
      <a-col :xs="24" :sm="24" :md="12">
        <a-card title="主机地区分布" :bordered="false">
          <template #extra>
            <a-select v-model:value="regionDisplayMode" style="width: 120px" size="small">
              <a-select-option value="pie">饼图</a-select-option>
              <a-select-option value="bar">柱状图</a-select-option>
            </a-select>
          </template>
          <div id="regionChart" ref="regionChartRef" style="height: 300px"></div>
        </a-card>
      </a-col>

      <!-- 近期告警 -->
      <a-col :span="24">
        <a-card title="近期告警" :bordered="false">
          <template #extra>
            <a-space>
              <a-radio-group v-model:value="alertDays" size="small" @change="loadAlerts">
                <a-radio-button :value="7">7天</a-radio-button>
                <a-radio-button :value="15">15天</a-radio-button>
                <a-radio-button :value="30">30天</a-radio-button>
              </a-radio-group>
              <a-button type="primary" size="small" @click="handleViewAllAlerts">
                查看全部
              </a-button>
            </a-space>
          </template>

          <a-table
            :columns="alertColumns"
            :data-source="alertList"
            :pagination="{ pageSize: 5 }"
            :loading="alertsLoading"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <!-- 告警类型 -->
              <template v-if="column.key === 'alert_type'">
                <a-tag :color="getAlertTypeColor(record.alert_type)">
                  {{ getAlertTypeText(record.alert_type) }}
                </a-tag>
              </template>
              
              <!-- 告警级别 -->
              <template v-if="column.key === 'level'">
                <a-tag :color="getAlertLevelColor(record.level)">
                  {{ getAlertLevelText(record.level) }}
                </a-tag>
              </template>
              
              <!-- 主机名 -->
              <template v-if="column.key === 'host_name'">
                <a @click="viewHostDetail(record.host?.id)">{{ record.host?.name }}</a>
              </template>
              
              <!-- 处理状态 -->
              <template v-if="column.key === 'is_resolved'">
                <a-badge :status="record.is_resolved ? 'success' : 'processing'" :text="record.is_resolved ? '已处理' : '未处理'" />
              </template>
              
              <!-- 操作 -->
              <template v-if="column.key === 'action'">
                <a-space>
                  <a @click="viewAlertDetail(record)">详情</a>
                  <a v-if="!record.is_resolved" @click="resolveAlert(record)">标记已处理</a>
                </a-space>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>
    
    <!-- 告警处理弹窗 -->
    <a-modal
      v-model:visible="resolveModalVisible"
      title="处理告警"
      :confirm-loading="resolveLoading"
      @ok="confirmResolveAlert"
    >
      <a-form :label-col="{ span: 4 }" :wrapper-col="{ span: 20 }">
        <a-form-item label="处理备注">
          <a-textarea
            v-model:value="resolveComment"
            :rows="4"
            placeholder="请输入处理过程和结果..."
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useHostStore } from '@/store/modules/host'
import { useECharts } from '@/utils/echarts/useECharts'
import { message } from 'ant-design-vue'
import { 
  DesktopOutlined, 
  CheckCircleOutlined, 
  WarningOutlined, 
  ClockCircleOutlined, 
  ReloadOutlined 
} from '@ant-design/icons-vue'
import type { HostAlert } from '@/types/api/host'

const router = useRouter()
const hostStore = useHostStore()

// 状态变量
const loading = ref(false)
const alertsLoading = ref(false)
const regionDisplayMode = ref('pie')
const alertDays = ref(7)
const statusChartRef = ref()
const regionChartRef = ref()
const statusChart = ref()
const regionChart = ref()
const alertList = ref<HostAlert[]>([])

// 告警处理
const resolveModalVisible = ref(false)
const resolveLoading = ref(false)
const resolveComment = ref('')
const currentAlert = ref<HostAlert | null>(null)

// 仪表盘数据
const dashboardData = reactive({
  totalHosts: 0,
  runningHosts: 0,
  errorHosts: 0,
  expiringHosts: 0,
  statusDistribution: [] as { name: string; value: number }[],
  regionDistribution: [] as { name: string; value: number }[]
})

// 告警表格列定义
const alertColumns = [
  {
    title: '告警级别',
    dataIndex: 'level',
    key: 'level',
    width: 100,
  },
  {
    title: '告警类型',
    dataIndex: 'alert_type',
    key: 'alert_type',
    width: 120,
  },
  {
    title: '主机名',
    dataIndex: ['host', 'name'],
    key: 'host_name',
    width: 200,
  },
  {
    title: '告警内容',
    dataIndex: 'message',
    ellipsis: true,
  },
  {
    title: '告警时间',
    dataIndex: 'time',
    width: 180,
  },
  {
    title: '状态',
    dataIndex: 'is_resolved',
    key: 'is_resolved',
    width: 100,
  },
  {
    title: '操作',
    key: 'action',
    width: 120,
  },
]

// 获取告警类型文字和颜色
const getAlertTypeText = (type: string) => {
  const map: Record<string, string> = {
    expired: '已过期',
    expiring: '即将过期',
    error: '错误',
    abnormal: '异常',
    unreachable: '无法连接'
  }
  return map[type] || type
}

const getAlertTypeColor = (type: string) => {
  const map: Record<string, string> = {
    expired: 'red',
    expiring: 'orange',
    error: 'red',
    abnormal: 'volcano',
    unreachable: 'purple'
  }
  return map[type] || 'blue'
}

// 获取告警级别文字和颜色
const getAlertLevelText = (level: string) => {
  const map: Record<string, string> = {
    info: '提示',
    warning: '警告',
    error: '错误',
    critical: '严重'
  }
  return map[level] || level
}

const getAlertLevelColor = (level: string) => {
  const map: Record<string, string> = {
    info: 'blue',
    warning: 'orange',
    error: 'red',
    critical: 'purple'
  }
  return map[level] || 'default'
}

// 加载仪表盘数据
const loadDashboardData = async () => {
  try {
    loading.value = true
    const res = await hostStore.fetchHostDashboard()
    if (res) {
      dashboardData.totalHosts = res.total
      dashboardData.runningHosts = res.running
      dashboardData.errorHosts = res.error
      dashboardData.expiringHosts = res.expiring
      dashboardData.statusDistribution = res.status_distribution
      dashboardData.regionDistribution = res.region_distribution
      
      // 更新图表
      renderStatusChart()
      renderRegionChart()
    }
  } catch (error) {
    console.error('加载仪表盘数据失败:', error)
    message.error('加载仪表盘数据失败')
  } finally {
    loading.value = false
  }
}

// 加载告警数据
const loadAlerts = async () => {
  try {
    alertsLoading.value = true
    const alerts = await hostStore.loadHostAlerts(alertDays.value)
    alertList.value = alerts
  } catch (error) {
    console.error('加载告警数据失败:', error)
    message.error('加载告警数据失败')
  } finally {
    alertsLoading.value = false
  }
}

// 渲染状态分布图
const renderStatusChart = () => {
  if (!statusChartRef.value) return
  
  const { initChart } = useECharts()
  if (!statusChart.value) {
    statusChart.value = initChart(statusChartRef.value)
  }
  
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center',
      data: dashboardData.statusDistribution.map(item => item.name)
    },
    series: [
      {
        name: '主机状态',
        type: 'pie',
        radius: ['50%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: dashboardData.statusDistribution.map(item => {
          let itemColor = ''
          switch (item.name) {
            case '运行中':
              itemColor = '#52c41a'
              break
            case '已停止':
              itemColor = '#faad14'
              break
            case '错误':
              itemColor = '#f5222d'
              break
            case '已过期':
              itemColor = '#d9d9d9'
              break
            default:
              itemColor = '#1890ff'
          }
          
          return {
            ...item,
            itemStyle: {
              color: itemColor
            }
          }
        })
      }
    ]
  }
  
  statusChart.value.setOption(option)
}

// 渲染地区分布图
const renderRegionChart = () => {
  if (!regionChartRef.value) return
  
  const { initChart } = useECharts()
  if (!regionChart.value) {
    regionChart.value = initChart(regionChartRef.value)
  }
  
  let option
  
  if (regionDisplayMode.value === 'pie') {
    option = {
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        right: 10,
        top: 'center',
        data: dashboardData.regionDistribution.map(item => item.name)
      },
      series: [
        {
          name: '地区分布',
          type: 'pie',
          radius: ['50%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: false,
            position: 'center'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 16,
              fontWeight: 'bold'
            }
          },
          labelLine: {
            show: false
          },
          data: dashboardData.regionDistribution
        }
      ]
    }
  } else {
    option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: dashboardData.regionDistribution.map(item => item.name),
        axisLabel: {
          interval: 0,
          rotate: 30
        }
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '主机数量',
          type: 'bar',
          data: dashboardData.regionDistribution.map(item => item.value)
        }
      ]
    }
  }
  
  regionChart.value.setOption(option)
}

// 刷新数据
const refreshData = async () => {
  await loadDashboardData()
  await loadAlerts()
}

// 查看主机详情
const viewHostDetail = (id: number) => {
  router.push(`/cmdb/host/${id}`)
}

// 查看告警详情
const viewAlertDetail = (alert: HostAlert) => {
  // 实际项目中可能需要展示告警详情的弹窗或页面
  message.info('查看告警详情: ' + alert.message)
}

// 处理告警
const resolveAlert = (alert: HostAlert) => {
  currentAlert.value = alert
  resolveModalVisible.value = true
  resolveComment.value = ''
}

// 确认处理告警
const confirmResolveAlert = async () => {
  if (!currentAlert.value) return
  
  try {
    resolveLoading.value = true
    await hostStore.resolveAlert(currentAlert.value.host.id, resolveComment.value)
    message.success('告警已处理')
    resolveModalVisible.value = false
    await loadAlerts()
  } catch (error) {
    console.error('处理告警失败:', error)
    message.error('处理告警失败')
  } finally {
    resolveLoading.value = false
  }
}

// 查看所有告警
const handleViewAllAlerts = () => {
  router.push('/cmdb/host-alerts')
}

// 监听地区显示模式变化
watch(regionDisplayMode, () => {
  renderRegionChart()
})

// 初始化
onMounted(() => {
  loadDashboardData()
  loadAlerts()
  
  // 窗口大小变化时重绘图表
  window.addEventListener('resize', handleResize)
})

// 组件销毁前
onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  
  // 销毁图表实例
  if (statusChart.value) {
    statusChart.value.dispose()
  }
  if (regionChart.value) {
    regionChart.value.dispose()
  }
})

// 窗口大小变化处理
const handleResize = () => {
  statusChart.value?.resize()
  regionChart.value?.resize()
}

// 添加getHostDetail方法到hostStore
if (!hostStore.getHostDetail) {
  hostStore.getHostDetail = async (id: number) => {
    try {
      const response = await hostStore.fetchHost(id)
      return response
    } catch (error) {
      throw error
    }
  }
}

// 添加fetchHostDashboard方法到hostStore
if (!hostStore.fetchHostDashboard) {
  hostStore.fetchHostDashboard = async () => {
    // 模拟仪表盘数据，实际项目中应从API获取
    return {
      total: 128,
      running: 108,
      error: 6,
      expiring: 14,
      status_distribution: [
        { name: '运行中', value: 108 },
        { name: '已停止', value: 8 },
        { name: '错误', value: 6 },
        { name: '已过期', value: 6 }
      ],
      region_distribution: [
        { name: '华北', value: 42 },
        { name: '华东', value: 36 },
        { name: '华南', value: 24 },
        { name: '西南', value: 12 },
        { name: '海外', value: 14 }
      ]
    }
  }
}

// 添加resolveAlert方法到hostStore
if (!hostStore.resolveAlert) {
  hostStore.resolveAlert = async (hostId: number, comment: string) => {
    // 实际项目中应调用API处理告警
    console.log('处理主机ID为', hostId, '的告警，备注:', comment)
    return Promise.resolve()
  }
}
</script>

<style lang="scss" scoped>
.host-dashboard-container {
  padding: 0;
  
  .stat-card {
    height: 100%;
    
    .stat-content {
      display: flex;
      align-items: center;
      
      .stat-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 64px;
        height: 64px;
        border-radius: 8px;
        font-size: 32px;
        color: #fff;
      }
      
      .stat-info {
        margin-left: 16px;
        
        .stat-title {
          color: rgba(0, 0, 0, 0.45);
          font-size: 14px;
        }
        
        .stat-value {
          font-size: 24px;
          font-weight: 600;
          margin-top: 4px;
        }
      }
    }
  }
}
</style> 
