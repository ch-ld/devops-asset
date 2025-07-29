<template>
  <div class="host-dashboard">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title">
          <h1>主机概览</h1>
          <p class="page-description">查看主机资源的整体状况和统计信息</p>
        </div>
        <div class="page-actions">
          <el-button @click="refreshData" :loading="loading">
            <el-icon><Refresh /></el-icon>
            刷新数据
          </el-button>
          <el-button type="primary" @click="goToHostList">
            <el-icon><List /></el-icon>
            主机列表
          </el-button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-overview">
      <el-row :gutter="24">
        <el-col :xs="12" :sm="6" :md="6" :lg="6">
          <el-card class="stats-card total" shadow="hover">
            <div class="stats-content">
              <div class="stats-icon">
                <el-icon><Monitor /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-number">{{ hostStats.total }}</div>
                <div class="stats-label">总主机数</div>
                <div class="stats-trend">
                  <span class="trend-text">总计</span>
                  <span class="trend-value">{{ hostStats.total }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :xs="12" :sm="6" :md="6" :lg="6">
          <el-card class="stats-card running" shadow="hover">
            <div class="stats-content">
              <div class="stats-icon">
                <el-icon><CircleCheck /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-number">{{ hostStats.running }}</div>
                <div class="stats-label">运行中</div>
                <div class="stats-trend">
                  <span class="trend-text">运行率</span>
                  <span class="trend-value">{{ runningRate }}%</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :xs="12" :sm="6" :md="6" :lg="6">
          <el-card class="stats-card stopped" shadow="hover">
            <div class="stats-content">
              <div class="stats-icon">
                <el-icon><VideoPause /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-number">{{ hostStats.stopped }}</div>
                <div class="stats-label">已停止</div>
                <div class="stats-trend">
                  <span class="trend-text">停止率</span>
                  <span class="trend-value">{{ stoppedRate }}%</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :xs="12" :sm="6" :md="6" :lg="6">
          <el-card class="stats-card error" shadow="hover">
            <div class="stats-content">
              <div class="stats-icon">
                <el-icon><CircleClose /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-number">{{ hostStats.error }}</div>
                <div class="stats-label">异常</div>
                <div class="stats-trend">
                  <span class="trend-text">异常率</span>
                  <span class="trend-value">{{ errorRate }}%</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 主要内容区域 -->
    <el-row :gutter="24" class="main-content">
      <!-- 左侧：主机状态分布 -->
      <el-col :xs="24" :sm="12" :md="12" :lg="12">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-header">
              <h3>主机状态分布</h3>
            </div>
          </template>
          <div class="chart-container">
            <div ref="statusChartRef" class="chart"></div>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧：主机组统计 -->
      <el-col :xs="24" :sm="12" :md="12" :lg="12">
        <el-card class="group-card" shadow="never">
          <template #header>
            <div class="card-header">
              <h3>主机组统计</h3>
            </div>
          </template>
          <div class="group-list">
            <div v-for="group in hostGroups" :key="group.id" class="group-item">
              <div class="group-info">
                <div class="group-name">{{ group.name }}</div>
                <div class="group-count">{{ group.hostCount }} 台主机</div>
              </div>
              <div class="group-progress">
                <el-progress
                  :percentage="Math.round(group.runningRate)"
                  :color="getProgressColor(group.runningRate)"
                  :show-text="false"
                  :stroke-width="6"
                />
                <span class="progress-text">{{ Math.round(group.runningRate) }}%</span>
              </div>
            </div>
            <div v-if="hostGroups.length === 0" class="no-groups">
              <el-icon><FolderOpened /></el-icon>
              <span>暂无主机组</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近添加的主机 -->
    <el-row :gutter="24" class="recent-hosts-section">
      <el-col :span="24">
        <el-card class="recent-hosts-card" shadow="never">
          <template #header>
            <div class="card-header">
              <h3>最近添加的主机</h3>
              <el-button text type="primary" @click="goToHostList">
                查看全部
                <el-icon><ArrowRight /></el-icon>
              </el-button>
            </div>
          </template>
          <div class="recent-hosts-list">
            <div v-for="host in recentHosts" :key="host.id" class="host-item">
              <div class="host-avatar">
                <el-icon><Monitor /></el-icon>
              </div>
              <div class="host-info">
                <div class="host-name">{{ host.name }}</div>
                <div class="host-details">
                  <span class="host-ip">{{ getHostIP(host) }}</span>
                  <span class="host-region">{{ host.region }}</span>
                  <span class="host-os">{{ host.os }}</span>
                </div>
              </div>
              <div class="host-status">
                <el-tag 
                  :type="getStatusType(host.status)"
                  :effect="host.status === 'running' ? 'dark' : 'plain'"
                  size="small"
                >
                  {{ getStatusText(host.status) }}
                </el-tag>
              </div>
              <div class="host-time">
                <span class="time-text">{{ formatTime(host.created_at) }}</span>
              </div>
            </div>
            <div v-if="recentHosts.length === 0" class="no-hosts">
              <el-icon><Monitor /></el-icon>
              <span>暂无主机数据</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
  import { ref, computed, onMounted, nextTick } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { 
    Monitor, Refresh, List, CircleCheck, VideoPause, CircleClose, 
    FolderOpened, ArrowRight
  } from '@element-plus/icons-vue'
  import * as echarts from 'echarts'
  import dayjs from 'dayjs'
  import { getHostStats, getHostGroupStats, getRecentHosts } from '@/api/system/host'

  const router = useRouter()

  // 响应式数据
  const loading = ref(false)
  const statusChartRef = ref(null)
  let statusChart = null

  // 统计数据
  const hostStats = ref({
    total: 0,
    running: 0,
    stopped: 0,
    error: 0
  })

  // 主机组数据
  const hostGroups = ref([])

  // 最近主机数据
  const recentHosts = ref([])

  // 计算属性
  const runningRate = computed(() => {
    if (hostStats.value.total === 0) return 0
    return Math.round((hostStats.value.running / hostStats.value.total) * 100)
  })

  const stoppedRate = computed(() => {
    if (hostStats.value.total === 0) return 0
    return Math.round((hostStats.value.stopped / hostStats.value.total) * 100)
  })

  const errorRate = computed(() => {
    if (hostStats.value.total === 0) return 0
    return Math.round((hostStats.value.error / hostStats.value.total) * 100)
  })

  // 生命周期
  onMounted(async () => {
    await refreshData()
    await nextTick()
    initCharts()
  })

  // 方法
  const refreshData = async () => {
    try {
      loading.value = true
      await Promise.all([
        fetchHostStats(),
        fetchHostGroups(),
        fetchRecentHosts()
      ])
    } catch (error) {
      console.error('刷新数据失败:', error)
      ElMessage.error('刷新数据失败')
    } finally {
      loading.value = false
    }
  }

  const fetchHostStats = async () => {
    try {
      const response = await getHostStats()
      hostStats.value = response.data || hostStats.value
      updateStatusChart()
    } catch (error) {
      console.error('获取统计数据失败:', error)
    }
  }

  const fetchHostGroups = async () => {
    try {
      const response = await getHostGroupStats()
      hostGroups.value = response.data || []
    } catch (error) {
      console.error('获取主机组统计失败:', error)
    }
  }

  const fetchRecentHosts = async () => {
    try {
      const response = await getRecentHosts(8)
      recentHosts.value = response.data || []
    } catch (error) {
      console.error('获取最近主机列表失败:', error)
    }
  }

  const initCharts = () => {
    if (statusChartRef.value) {
      statusChart = echarts.init(statusChartRef.value)
      updateStatusChart()
    }
  }

  const updateStatusChart = () => {
    if (!statusChart) return

    const option = {
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        left: 10,
        data: ['运行中', '已停止', '异常']
      },
      series: [
        {
          name: '主机状态',
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['60%', '50%'],
          avoidLabelOverlap: false,
          label: {
            show: false,
            position: 'center'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: '18',
              fontWeight: 'bold'
            }
          },
          labelLine: {
            show: false
          },
          data: [
            { 
              value: hostStats.value.running, 
              name: '运行中',
              itemStyle: { color: '#67C23A' }
            },
            { 
              value: hostStats.value.stopped, 
              name: '已停止',
              itemStyle: { color: '#E6A23C' }
            },
            { 
              value: hostStats.value.error, 
              name: '异常',
              itemStyle: { color: '#F56C6C' }
            }
          ]
        }
      ]
    }

    statusChart.setOption(option)
  }

  const goToHostList = () => {
    router.push('/cmdb/host')
  }

  const getProgressColor = (rate) => {
    if (rate >= 80) return '#67C23A'
    if (rate >= 60) return '#E6A23C'
    return '#F56C6C'
  }

  const getHostIP = (host) => {
    if (host.public_ip && host.public_ip.length > 0) {
      return Array.isArray(host.public_ip) ? host.public_ip[0] : host.public_ip
    }
    if (host.private_ip && host.private_ip.length > 0) {
      return Array.isArray(host.private_ip) ? host.private_ip[0] : host.private_ip
    }
    return '未知'
  }

  const getStatusType = (status) => {
    const statusMap = {
      running: 'success',
      stopped: 'warning',
      error: 'danger',
      unknown: 'info'
    }
    return statusMap[status] || 'info'
  }

  const getStatusText = (status) => {
    const statusMap = {
      running: '运行中',
      stopped: '已停止',
      error: '异常',
      unknown: '未知'
    }
    return statusMap[status] || '未知'
  }

  const formatTime = (time) => {
    return dayjs(time).format('MM-DD HH:mm')
  }

  const goToHostDetail = (hostId) => {
    router.push(`/cmdb/host/detail/${hostId}`)
  }
</script>

<style lang="scss" scoped>
.host-dashboard {
  padding: 16px;
  background-color: #f5f5f5;
  min-height: 100vh;

  .page-header {
    margin-bottom: 24px;

    .page-header-content {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 20px 24px;
      background: white;
      border-radius: 8px;
      box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);

      .page-title {
        h1 {
          margin: 0;
          font-size: 24px;
          font-weight: 600;
          color: #1f2937;
        }

        .page-description {
          margin: 4px 0 0 0;
          color: #6b7280;
          font-size: 14px;
        }
      }

      .page-actions {
        display: flex;
        gap: 12px;
      }
    }
  }

  .stats-overview {
    margin-bottom: 24px;

    .stats-card {
      border: none;
      border-radius: 12px;
      overflow: hidden;
      transition: all 0.3s ease;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 8px 25px 0 rgba(0, 0, 0, 0.1);
      }

      &.total {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
      }

      &.running {
        background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
        color: white;
      }

      &.stopped {
        background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
        color: white;
      }

      &.error {
        background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
        color: white;
      }

      .stats-content {
        display: flex;
        align-items: center;
        padding: 20px;

        .stats-icon {
          font-size: 48px;
          margin-right: 16px;
          opacity: 0.8;
        }

        .stats-info {
          flex: 1;

          .stats-number {
            font-size: 32px;
            font-weight: 700;
            line-height: 1;
            margin-bottom: 4px;
          }

          .stats-label {
            font-size: 14px;
            opacity: 0.9;
            margin-bottom: 8px;
          }

          .stats-trend {
            display: flex;
            justify-content: space-between;
            font-size: 12px;
            opacity: 0.8;

            .trend-value {
              font-weight: 600;
            }
          }
        }
      }
    }
  }

  .main-content {
    margin-bottom: 24px;

    .chart-card, .group-card {
      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;

        h3 {
          margin: 0;
          font-size: 16px;
          font-weight: 600;
          color: #1f2937;
        }
      }

      .chart-container {
        height: 300px;

        .chart {
          width: 100%;
          height: 100%;
        }
      }

      .group-list {
        max-height: 300px;
        overflow-y: auto;

        .group-item {
          display: flex;
          align-items: center;
          padding: 12px 0;
          border-bottom: 1px solid #f3f4f6;

          &:last-child {
            border-bottom: none;
          }

          .group-info {
            flex: 1;
            margin-right: 16px;

            .group-name {
              font-size: 14px;
              font-weight: 500;
              color: #1f2937;
              margin-bottom: 4px;
            }

            .group-count {
              font-size: 12px;
              color: #6b7280;
            }
          }

          .group-progress {
            display: flex;
            align-items: center;
            min-width: 120px;

            .el-progress {
              flex: 1;
              margin-right: 8px;
            }

            .progress-text {
              font-size: 12px;
              color: #6b7280;
              min-width: 30px;
              text-align: right;
            }
          }
        }

        .no-groups {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          padding: 40px 20px;
          color: #9ca3af;

          .el-icon {
            font-size: 48px;
            margin-bottom: 12px;
          }
        }
      }
    }
  }

  .recent-hosts-section {
    .recent-hosts-card {
      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;

        h3 {
          margin: 0;
          font-size: 16px;
          font-weight: 600;
          color: #1f2937;
        }
      }

      .recent-hosts-list {
        .host-item {
          display: flex;
          align-items: center;
          padding: 16px 0;
          border-bottom: 1px solid #f3f4f6;

          &:last-child {
            border-bottom: none;
          }

          .host-avatar {
            width: 40px;
            height: 40px;
            border-radius: 8px;
            background: #f3f4f6;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-right: 12px;

            .el-icon {
              font-size: 20px;
              color: #6b7280;
            }
          }

          .host-info {
            flex: 1;
            margin-right: 16px;

            .host-name {
              font-size: 14px;
              font-weight: 500;
              color: #1f2937;
              margin-bottom: 4px;
            }

            .host-details {
              display: flex;
              gap: 12px;
              font-size: 12px;
              color: #6b7280;

              .host-ip, .host-region, .host-os {
                padding: 2px 6px;
                background: #f9fafb;
                border-radius: 4px;
              }
            }
          }

          .host-status {
            margin-right: 16px;
          }

          .host-time {
            min-width: 80px;
            text-align: right;

            .time-text {
              font-size: 12px;
              color: #9ca3af;
            }
          }
        }

        .no-hosts {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          padding: 40px 20px;
          color: #9ca3af;

          .el-icon {
            font-size: 48px;
            margin-bottom: 12px;
          }
        }
      }
    }
  }
}
</style>
