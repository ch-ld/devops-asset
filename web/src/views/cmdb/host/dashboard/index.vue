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
                <el-icon><CircleClose /></el-icon>
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
          <el-card class="stats-card warning" shadow="hover">
            <div class="stats-content">
              <div class="stats-icon">
                <el-icon><Warning /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-number">{{ hostStats.error }}</div>
                <div class="stats-label">异常</div>
                <div class="stats-trend">
                  <span class="trend-text">错误</span>
                  <span class="trend-value warning">{{ hostStats.error }}台</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 图表和详细信息 -->
    <el-row :gutter="24">
      <!-- 左侧图表区域 -->
      <el-col :lg="16" :md="24">
        <!-- 主机状态分布图 -->
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-header">
              <h3>主机状态分布</h3>
              <el-radio-group v-model="chartTimeRange" size="small">
                <el-radio-button label="7d">7天</el-radio-button>
                <el-radio-button label="30d">30天</el-radio-button>
                <el-radio-button label="90d">90天</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container">
            <div ref="statusChartRef" class="chart" style="height: 300px;"></div>
          </div>
        </el-card>

        <!-- 云厂商分布图 -->
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-header">
              <h3>云厂商分布</h3>
            </div>
          </template>
          <div class="chart-container">
            <div ref="providerChartRef" class="chart" style="height: 300px;"></div>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧信息面板 -->
      <el-col :lg="8" :md="24">
        <!-- 告警信息 -->
        <el-card class="alert-card" shadow="never">
          <template #header>
            <div class="card-header">
              <h3>告警信息</h3>
              <el-badge :value="alertCount" class="alert-badge">
                <el-icon><Bell /></el-icon>
              </el-badge>
            </div>
          </template>
          <div class="alert-list">
            <div v-for="alert in alerts" :key="alert.id" class="alert-item">
              <div class="alert-icon" :class="alert.level">
                <el-icon><component :is="getAlertIcon(alert.level)" /></el-icon>
              </div>
              <div class="alert-content">
                <div class="alert-title">{{ alert.title }}</div>
                <div class="alert-time">{{ formatTime(alert.time) }}</div>
              </div>
            </div>
            <div v-if="alerts.length === 0" class="no-alerts">
              <el-icon><CircleCheck /></el-icon>
              <span>暂无告警信息</span>
            </div>
          </div>
        </el-card>

        <!-- 主机组统计 -->
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
                  :percentage="group.runningRate"
                  :color="getProgressColor(group.runningRate)"
                  :show-text="false"
                  :stroke-width="6"
                />
                <span class="progress-text">{{ group.runningRate }}%</span>
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
  </div>
</template>

<script setup>
  import { ref, computed, onMounted, nextTick } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { getHostStats, getHostList } from '@/api/system/host'

  const router = useRouter()

  // 响应式数据
  const loading = ref(false)

  // 统计数据
  const hostStats = ref({
    total: 0,
    running: 0,
    stopped: 0,
    error: 0
  })

  const hostList = ref([])

  // 告警数据
  const alerts = ref([])

  // 主机组数据
  const hostGroups = ref([])

  // 计算属性
  const runningRate = computed(() => {
    if (hostStats.value.total === 0) return 0
    return Math.round((hostStats.value.running / hostStats.value.total) * 100)
  })

  const stoppedRate = computed(() => {
    if (hostStats.value.total === 0) return 0
    return Math.round((hostStats.value.stopped / hostStats.value.total) * 100)
  })

  const alertCount = computed(() => alerts.value.length)

  // 生命周期
  onMounted(async () => {
    await refreshData()
  })

  // 方法
  const refreshData = async () => {
    try {
      loading.value = true
      await Promise.all([
        fetchHostStats(),
        fetchHostList()
      ])
    } catch (error) {
      console.error('刷新数据失败:', error)
    } finally {
      loading.value = false
    }
  }

  const fetchHostStats = async () => {
    try {
      const response = await getHostStats()
      hostStats.value = response.data || hostStats.value
    } catch (error) {
      console.error('获取统计数据失败:', error)
    }
  }

  const fetchHostList = async () => {
    try {
      const response = await getHostList({ page: 1, page_size: 100 })
      hostList.value = response.data?.list || []
    } catch (error) {
      console.error('获取主机列表失败:', error)
    }
  }

  const goToHostList = () => {
    router.push('/cmdb/host')
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
      border-radius: 8px;
      border: none;
      transition: all 0.3s ease;

      &:hover {
        transform: translateY(-2px);
      }

      .stats-content {
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 20px;

        .stats-icon {
          width: 60px;
          height: 60px;
          border-radius: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 28px;
          color: white;
        }

        .stats-info {
          flex: 1;

          .stats-number {
            font-size: 32px;
            font-weight: 700;
            color: #1f2937;
            line-height: 1;
            margin-bottom: 4px;
          }

          .stats-label {
            font-size: 14px;
            color: #6b7280;
            margin-bottom: 8px;
          }

          .stats-trend {
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 12px;

            .trend-text {
              color: #9ca3af;
            }

            .trend-value {
              font-weight: 600;

              &.positive {
                color: #10b981;
              }

              &.warning {
                color: #f59e0b;
              }
            }
          }
        }
      }

      &.total .stats-icon {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      }

      &.running .stats-icon {
        background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
      }

      &.stopped .stats-icon {
        background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
      }

      &.warning .stats-icon {
        background: linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%);
        color: #d97706;
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

  .alert-card, .group-card {
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

      .alert-badge {
        .el-icon {
          font-size: 18px;
          color: #6b7280;
        }
      }
    }

    .alert-list {
      .alert-item {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 12px 0;
        border-bottom: 1px solid #f3f4f6;

        &:last-child {
          border-bottom: none;
        }

        .alert-icon {
          width: 32px;
          height: 32px;
          border-radius: 6px;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 16px;

          &.error {
            background: #fef2f2;
            color: #dc2626;
          }

          &.warning {
            background: #fffbeb;
            color: #d97706;
          }

          &.info {
            background: #eff6ff;
            color: #2563eb;
          }
        }

        .alert-content {
          flex: 1;

          .alert-title {
            font-size: 14px;
            color: #374151;
            margin-bottom: 4px;
          }

          .alert-time {
            font-size: 12px;
            color: #9ca3af;
          }
        }
      }

      .no-alerts {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        padding: 40px 0;
        color: #9ca3af;
        font-size: 14px;
      }
    }

    .group-list {
      .group-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 16px 0;
        border-bottom: 1px solid #f3f4f6;

        &:last-child {
          border-bottom: none;
        }

        .group-info {
          .group-name {
            font-size: 14px;
            font-weight: 500;
            color: #374151;
            margin-bottom: 4px;
          }

          .group-count {
            font-size: 12px;
            color: #9ca3af;
          }
        }

        .group-progress {
          display: flex;
          align-items: center;
          gap: 8px;
          min-width: 100px;

          .el-progress {
            flex: 1;
          }

          .progress-text {
            font-size: 12px;
            color: #6b7280;
            min-width: 32px;
            text-align: right;
          }
        }
      }

      .no-groups {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        padding: 40px 0;
        color: #9ca3af;
        font-size: 14px;
      }
    }
  }
}

// 全局样式覆盖
:deep(.el-card) {
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
}

:deep(.el-progress-bar__outer) {
  border-radius: 3px;
}

:deep(.el-progress-bar__inner) {
  border-radius: 3px;
}

@media (max-width: 768px) {
  .host-dashboard {
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

    .stats-overview {
      .stats-card .stats-content {
        padding: 16px;

        .stats-icon {
          width: 48px;
          height: 48px;
          font-size: 24px;
        }

        .stats-info .stats-number {
          font-size: 24px;
        }
      }
    }
  }
}
</style>
