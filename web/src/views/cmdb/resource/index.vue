<template>
  <div class="resource-overview">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title">
          <h1>资源概览</h1>
          <p class="page-description">查看所有云资源的整体状况和统计信息</p>
        </div>
        <div class="page-actions">
          <el-button @click="refreshData" :loading="loading">
            <el-icon><Refresh /></el-icon>
            刷新数据
          </el-button>
          <el-dropdown @command="handleExportCommand" trigger="click">
            <el-button>
              <el-icon><Download /></el-icon>
              导出报告
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="excel">导出Excel</el-dropdown-item>
                <el-dropdown-item command="pdf">导出PDF</el-dropdown-item>
                <el-dropdown-item command="csv">导出CSV</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </div>

    <!-- 健康状况卡片 -->
    <div class="health-card">
      <el-card shadow="never" class="health-status-card">
        <div class="health-content">
          <div class="health-score">
            <div class="score-circle" :class="healthLevel">
              <div class="score-number">{{ Math.round(healthScore) }}</div>
              <div class="score-label">健康分数</div>
            </div>
          </div>
          <div class="health-details">
            <div class="health-item">
              <div class="health-label">总资源数</div>
              <div class="health-value">{{ totalResources }}</div>
            </div>
            <div class="health-item">
              <div class="health-label">运行中</div>
              <div class="health-value success">{{ runningResources }}</div>
            </div>
            <div class="health-item">
              <div class="health-label">异常</div>
              <div class="health-value error">{{ errorResources }}</div>
            </div>
            <div class="health-item">
              <div class="health-label">即将过期</div>
              <div class="health-value warning">{{ expiringResources }}</div>
            </div>
          </div>
          <div class="health-trends">
            <div class="trend-item">
              <span class="trend-label">运行率</span>
              <el-progress 
                :percentage="runningRate" 
                :color="getProgressColor(runningRate)"
                :show-text="false"
                :stroke-width="8"
              />
              <span class="trend-value">{{ runningRate.toFixed(1) }}%</span>
            </div>
            <div class="trend-item">
              <span class="trend-label">异常率</span>
              <el-progress 
                :percentage="errorRate" 
                color="#f56c6c"
                :show-text="false"
                :stroke-width="8"
              />
              <span class="trend-value">{{ errorRate.toFixed(1) }}%</span>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 资源统计卡片 -->
    <div class="resource-stats">
      <el-row :gutter="24">
        <el-col :xs="12" :sm="6" :md="6" :lg="6" v-for="resource in resourceSummary" :key="resource.type">
          <el-card class="resource-card" shadow="hover" @click="navigateToResource(resource.type)">
            <div class="resource-content">
              <div class="resource-icon" :class="getResourceIconClass(resource.type)">
                <el-icon><component :is="getResourceIcon(resource.type)" /></el-icon>
              </div>
              <div class="resource-info">
                <div class="resource-name">{{ resource.name }}</div>
                <div class="resource-total">{{ resource.total }}</div>
                <div class="resource-status">
                  <span class="status-item running">运行: {{ resource.running }}</span>
                  <span class="status-item stopped">停止: {{ resource.stopped }}</span>
                  <span class="status-item error" v-if="resource.error > 0">异常: {{ resource.error }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 图表区域 -->
    <el-row :gutter="24">
      <!-- 左侧图表 -->
      <el-col :lg="16" :md="24">
        <!-- 云厂商分布 -->
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-header">
              <h3>云厂商分布</h3>
              <el-radio-group v-model="chartType" size="small">
                <el-radio-button label="pie">饼图</el-radio-button>
                <el-radio-button label="bar">柱状图</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container">
            <div ref="providerChartRef" class="chart" style="height: 300px;"></div>
          </div>
        </el-card>

        <!-- 地区分布 -->
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-header">
              <h3>地区分布</h3>
            </div>
          </template>
          <div class="chart-container">
            <div ref="regionChartRef" class="chart" style="height: 300px;"></div>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧信息面板 -->
      <el-col :lg="8" :md="24">
        <!-- 即将过期资源 -->
        <el-card class="expiring-card" shadow="never">
          <template #header>
            <div class="card-header">
              <h3>即将过期资源</h3>
              <el-badge :value="expiringResourcesList.length" class="expiring-badge">
                <el-icon><Warning /></el-icon>
              </el-badge>
            </div>
          </template>
          <div class="expiring-list">
            <div v-for="resource in expiringResourcesList" :key="`${resource.type}-${resource.id}`" class="expiring-item">
              <div class="expiring-icon" :class="getResourceIconClass(resource.type)">
                <el-icon><component :is="getResourceIcon(resource.type)" /></el-icon>
              </div>
              <div class="expiring-content">
                <div class="expiring-name">{{ resource.name }}</div>
                <div class="expiring-meta">
                  <span class="expiring-type">{{ getResourceTypeName(resource.type) }}</span>
                  <span class="expiring-provider">{{ getProviderDisplayName(resource.provider_type) }}</span>
                </div>
                <div class="expiring-time">
                  <el-tag :type="getDaysLeftTagType(resource.days_left)" size="small">
                    {{ resource.days_left }}天后过期
                  </el-tag>
                </div>
              </div>
            </div>
            <div v-if="expiringResourcesList.length === 0" class="no-expiring">
              <el-icon><CircleCheck /></el-icon>
              <span>暂无即将过期的资源</span>
            </div>
          </div>
        </el-card>

        <!-- 快速操作 -->
        <el-card class="quick-actions-card" shadow="never">
          <template #header>
            <h3>快速操作</h3>
          </template>
          <div class="quick-actions">
            <el-button type="primary" size="large" class="action-btn" @click="navigateToResource('host')">
              <el-icon><Monitor /></el-icon>
              管理主机
            </el-button>
            <el-button type="success" size="large" class="action-btn" @click="syncAllResources">
              <el-icon><Refresh /></el-icon>
              同步资源
            </el-button>
            <el-button size="large" class="action-btn" @click="showAlerts">
              <el-icon><Bell /></el-icon>
              查看告警
            </el-button>
            <el-button size="large" class="action-btn" @click="exportReport">
              <el-icon><Download /></el-icon>
              导出报告
            </el-button>
          </div>
        </el-card>

        <!-- 最近更新 -->
        <el-card class="recent-updates-card" shadow="never">
          <template #header>
            <h3>最近更新</h3>
          </template>
          <div class="recent-updates">
            <div class="update-item">
              <div class="update-time">{{ formatTime(updateTime) }}</div>
              <div class="update-content">数据最后更新时间</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
  import { ref, computed, onMounted, nextTick, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage, ElLoading } from 'element-plus'
  import {
    Refresh, Download, ArrowDown, Warning, CircleCheck, Monitor,
    Bell, Database, Coin, Share
  } from '@element-plus/icons-vue'
  import dayjs from 'dayjs'
  import * as echarts from 'echarts'

  import { resourceApi } from '@/api/resource'

  const router = useRouter()

  // 响应式数据
  const loading = ref(false)
  const chartType = ref('pie')

  // 图表引用
  const providerChartRef = ref(null)
  const regionChartRef = ref(null)

  // 统计数据
  const totalResources = ref(0)
  const runningResources = ref(0)
  const errorResources = ref(0)
  const expiringResources = ref(0)
  const healthScore = ref(100)
  const healthLevel = ref('excellent')
  const runningRate = ref(100)
  const errorRate = ref(0)
  const updateTime = ref(new Date())

  // 资源摘要
  const resourceSummary = ref([])
  const providerSummary = ref({})
  const regionSummary = ref({})
  const expiringResourcesList = ref([])

  // 生命周期
  onMounted(async () => {
    await refreshData()
    await nextTick()
    initCharts()
  })

  // 监听图表类型变化
  watch(chartType, () => {
    initProviderChart()
  })

  // 方法
  const refreshData = async () => {
    try {
      loading.value = true

      // 获取整体统计信息
      const [statisticsRes, healthRes] = await Promise.all([
        resourceApi.getOverallStatistics(),
        resourceApi.getResourceHealth()
      ])

      const statistics = statisticsRes.data
      const health = healthRes.data

      // 更新统计数据
      totalResources.value = statistics.total_resources
      resourceSummary.value = statistics.resource_summary
      providerSummary.value = statistics.provider_summary
      regionSummary.value = statistics.region_summary
      expiringResourcesList.value = statistics.expiring_resources

      // 更新健康数据
      healthScore.value = health.health_score
      healthLevel.value = health.health_level
      runningResources.value = health.running_resources
      errorResources.value = health.error_resources
      expiringResources.value = health.expiring_resources
      runningRate.value = health.running_rate || 0
      errorRate.value = health.error_rate || 0
      updateTime.value = new Date()

    } catch (error) {
      ElMessage.error('获取数据失败: ' + error.message)
    } finally {
      loading.value = false
    }
  }

  const initCharts = () => {
    initProviderChart()
    initRegionChart()
  }

  const initProviderChart = () => {
    if (!providerChartRef.value) return

    const chart = echarts.init(providerChartRef.value)

    const providerNames = {
      'aliyun': '阿里云',
      'tencent': '腾讯云',
      'aws': 'AWS',
      'manual': '自建'
    }

    const data = Object.entries(providerSummary.value).map(([key, value]) => ({
      name: providerNames[key] || key,
      value
    })).filter(item => item.value > 0)

    let option

    if (chartType.value === 'pie') {
      option = {
        tooltip: {
          trigger: 'item',
          formatter: '{a} <br/>{b}: {c} ({d}%)'
        },
        legend: {
          orient: 'vertical',
          left: 'left'
        },
        series: [
          {
            name: '云厂商',
            type: 'pie',
            radius: ['40%', '70%'],
            avoidLabelOverlap: false,
            data,
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
              }
            }
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
        xAxis: {
          type: 'category',
          data: data.map(item => item.name)
        },
        yAxis: {
          type: 'value'
        },
        series: [
          {
            name: '资源数量',
            type: 'bar',
            data: data.map(item => item.value),
            itemStyle: {
              color: '#409eff'
            }
          }
        ]
      }
    }

    chart.setOption(option)
  }

  const initRegionChart = () => {
    if (!regionChartRef.value) return

    const chart = echarts.init(regionChartRef.value)

    const data = Object.entries(regionSummary.value).map(([key, value]) => ({
      name: key,
      value
    })).filter(item => item.value > 0)

    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      xAxis: {
        type: 'category',
        data: data.map(item => item.name),
        axisLabel: {
          rotate: 45
        }
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '资源数量',
          type: 'bar',
          data: data.map(item => item.value),
          itemStyle: {
            color: '#67c23a'
          }
        }
      ]
    }

    chart.setOption(option)
  }

  // 工具函数
  const getResourceIcon = (type) => {
    const iconMap = {
      'host': 'Monitor',
      'load_balancer': 'Share',
      'database': 'Database',
      'redis': 'Coin'
    }
    return iconMap[type] || 'Monitor'
  }

  const getResourceIconClass = (type) => {
    return `resource-icon-${type}`
  }

  const getResourceTypeName = (type) => {
    const nameMap = {
      'host': '主机',
      'load_balancer': '负载均衡器',
      'database': '数据库',
      'redis': 'Redis'
    }
    return nameMap[type] || type
  }

  const getProviderDisplayName = (providerType) => {
    const nameMap = {
      'aliyun': '阿里云',
      'tencent': '腾讯云',
      'aws': 'AWS',
      'manual': '自建'
    }
    return nameMap[providerType] || providerType
  }

  const getDaysLeftTagType = (daysLeft) => {
    if (daysLeft <= 7) return 'danger'
    if (daysLeft <= 15) return 'warning'
    return 'info'
  }

  const getProgressColor = (percentage) => {
    if (percentage >= 90) return '#67c23a'
    if (percentage >= 70) return '#e6a23c'
    return '#f56c6c'
  }

  const formatTime = (time) => {
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
  }

  // 事件处理
  const navigateToResource = (type) => {
    const routeMap = {
      'host': '/cmdb/host',
      'load_balancer': '/cmdb/load-balancer',
      'database': '/cmdb/database',
      'redis': '/cmdb/redis'
    }
    const route = routeMap[type]
    if (route) {
      router.push(route)
    }
  }

  const handleExportCommand = (command) => {
    ElMessage.info(`导出${command.toUpperCase()}功能开发中...`)
  }

  const syncAllResources = async () => {
    try {
      const loading = ElLoading.service({
        lock: true,
        text: '正在同步所有资源...',
        background: 'rgba(0, 0, 0, 0.7)'
      })

      // 这里应该调用同步API
      await new Promise(resolve => setTimeout(resolve, 2000))

      ElMessage.success('资源同步成功')
      await refreshData()
      loading.close()
    } catch (error) {
      ElMessage.error('资源同步失败: ' + error.message)
    }
  }

  const showAlerts = () => {
    router.push('/cmdb/alerts')
  }

  const exportReport = () => {
    ElMessage.info('导出报告功能开发中...')
  }
</script>

<style lang="scss" scoped>
.resource-overview {
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

  .health-card {
    margin-bottom: 24px;

    .health-status-card {
      border-radius: 8px;
      border: none;

      .health-content {
        display: flex;
        align-items: center;
        gap: 32px;
        padding: 20px;

        .health-score {
          .score-circle {
            width: 120px;
            height: 120px;
            border-radius: 50%;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            position: relative;

            &.excellent {
              background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
            }

            &.good {
              background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
            }

            &.fair {
              background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
            }

            &.poor {
              background: linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%);
            }

            &.critical {
              background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%);
            }

            .score-number {
              font-size: 32px;
              font-weight: 700;
              color: white;
              line-height: 1;
            }

            .score-label {
              font-size: 12px;
              color: white;
              margin-top: 4px;
            }
          }
        }

        .health-details {
          display: flex;
          gap: 32px;

          .health-item {
            text-align: center;

            .health-label {
              font-size: 14px;
              color: #6b7280;
              margin-bottom: 8px;
            }

            .health-value {
              font-size: 24px;
              font-weight: 600;
              color: #1f2937;

              &.success {
                color: #10b981;
              }

              &.error {
                color: #ef4444;
              }

              &.warning {
                color: #f59e0b;
              }
            }
          }
        }

        .health-trends {
          flex: 1;

          .trend-item {
            display: flex;
            align-items: center;
            gap: 12px;
            margin-bottom: 16px;

            .trend-label {
              min-width: 60px;
              font-size: 14px;
              color: #6b7280;
            }

            .el-progress {
              flex: 1;
            }

            .trend-value {
              min-width: 50px;
              text-align: right;
              font-size: 14px;
              font-weight: 500;
              color: #374151;
            }
          }
        }
      }
    }
  }

  .resource-stats {
    margin-bottom: 24px;

    .resource-card {
      border-radius: 8px;
      border: none;
      cursor: pointer;
      transition: all 0.3s ease;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px 0 rgba(0, 0, 0, 0.15);
      }

      .resource-content {
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 20px;

        .resource-icon {
          width: 60px;
          height: 60px;
          border-radius: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 28px;
          color: white;

          &.resource-icon-host {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          }

          &.resource-icon-load_balancer {
            background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
          }

          &.resource-icon-database {
            background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
          }

          &.resource-icon-redis {
            background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
          }
        }

        .resource-info {
          flex: 1;

          .resource-name {
            font-size: 16px;
            font-weight: 600;
            color: #1f2937;
            margin-bottom: 4px;
          }

          .resource-total {
            font-size: 28px;
            font-weight: 700;
            color: #1f2937;
            line-height: 1;
            margin-bottom: 8px;
          }

          .resource-status {
            display: flex;
            gap: 12px;
            font-size: 12px;

            .status-item {
              &.running {
                color: #10b981;
              }

              &.stopped {
                color: #6b7280;
              }

              &.error {
                color: #ef4444;
              }
            }
          }
        }
      }
    }
  }

  .chart-card {
    margin-bottom: 24px;
    border-radius: 8px;
    border: none;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      h3 {
        margin: 0;
        font-size: 18px;
        font-weight: 600;
        color: #1f2937;
      }
    }

    .chart-container {
      .chart {
        width: 100%;
      }
    }
  }

  .expiring-card, .quick-actions-card, .recent-updates-card {
    margin-bottom: 24px;
    border-radius: 8px;
    border: none;

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

      .expiring-badge {
        .el-icon {
          font-size: 18px;
          color: #f59e0b;
        }
      }
    }

    .expiring-list {
      .expiring-item {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 12px 0;
        border-bottom: 1px solid #f3f4f6;

        &:last-child {
          border-bottom: none;
        }

        .expiring-icon {
          width: 40px;
          height: 40px;
          border-radius: 8px;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 18px;
          color: white;

          &.resource-icon-host {
            background: #667eea;
          }

          &.resource-icon-load_balancer {
            background: #f093fb;
          }

          &.resource-icon-database {
            background: #4facfe;
          }

          &.resource-icon-redis {
            background: #43e97b;
          }
        }

        .expiring-content {
          flex: 1;

          .expiring-name {
            font-size: 14px;
            font-weight: 500;
            color: #374151;
            margin-bottom: 4px;
          }

          .expiring-meta {
            display: flex;
            gap: 8px;
            font-size: 12px;
            color: #9ca3af;
            margin-bottom: 4px;

            .expiring-type, .expiring-provider {
              &::after {
                content: '•';
                margin-left: 8px;
              }

              &:last-child::after {
                display: none;
              }
            }
          }

          .expiring-time {
            .el-tag {
              font-size: 11px;
            }
          }
        }
      }

      .no-expiring {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        padding: 40px 0;
        color: #9ca3af;
        font-size: 14px;
      }
    }

    .quick-actions {
      display: flex;
      flex-direction: column;
      gap: 12px;

      .action-btn {
        width: 100%;
        height: 48px;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        font-size: 14px;
        font-weight: 500;
      }
    }

    .recent-updates {
      .update-item {
        padding: 12px 0;

        .update-time {
          font-size: 12px;
          color: #9ca3af;
          margin-bottom: 4px;
        }

        .update-content {
          font-size: 14px;
          color: #374151;
        }
      }
    }
  }
}

// 全局样式覆盖
:deep(.el-card) {
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
}

:deep(.el-progress-bar__outer) {
  border-radius: 4px;
}

:deep(.el-progress-bar__inner) {
  border-radius: 4px;
}

@media (max-width: 768px) {
  .resource-overview {
    padding: 12px;

    .page-header .page-header-content {
      flex-direction: column;
      align-items: flex-start;
      gap: 16px;

      .page-actions {
        width: 100%;
        justify-content: flex-end;
      }
    }

    .health-card .health-content {
      flex-direction: column;
      gap: 20px;

      .health-details {
        gap: 20px;
      }
    }

    .resource-stats .resource-card .resource-content {
      padding: 16px;

      .resource-icon {
        width: 48px;
        height: 48px;
        font-size: 24px;
      }

      .resource-info .resource-total {
        font-size: 24px;
      }
    }
  }
}
</style>
