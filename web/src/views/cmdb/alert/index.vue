<template>
  <div class="alert-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title">
          <h1>告警管理</h1>
          <p class="page-description">查看和管理系统告警信息</p>
        </div>
        <div class="page-actions">
          <el-button @click="refreshData" :loading="loading">
            <el-icon><Refresh /></el-icon>
            刷新数据
          </el-button>
          <el-button type="primary" @click="triggerCheck">
            <el-icon><Search /></el-icon>
            手动检查
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
                <el-icon><Bell /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-number">{{ totalAlerts }}</div>
                <div class="stats-label">总告警数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :xs="12" :sm="6" :md="6" :lg="6">
          <el-card class="stats-card expired" shadow="hover">
            <div class="stats-content">
              <div class="stats-icon">
                <el-icon><CircleClose /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-number">{{ expiredAlerts }}</div>
                <div class="stats-label">已过期</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :xs="12" :sm="6" :md="6" :lg="6">
          <el-card class="stats-card expiring" shadow="hover">
            <div class="stats-content">
              <div class="stats-icon">
                <el-icon><Warning /></el-icon>
              </div>
              <div class="stats-info">
                <div class="stats-number">{{ expiringAlerts }}</div>
                <div class="stats-label">即将过期</div>
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
                <div class="stats-number">{{ errorAlerts }}</div>
                <div class="stats-label">状态异常</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 告警列表 -->
    <el-card class="alert-list-card" shadow="never">
      <template #header>
        <div class="card-header">
          <h3>告警列表</h3>
          <div class="filter-controls">
            <div class="custom-select-wrapper">
              <select
                v-model="filterLevel"
                class="custom-select"
                @change="handleFilterChange"
              >
                <option value="">全部</option>
                <option value="expired">已过期</option>
                <option value="expiring">即将过期</option>
                <option value="error">状态异常</option>
                <option value="abnormal">连接异常</option>
              </select>
              <div class="select-arrow">
                <el-icon><ArrowDown /></el-icon>
              </div>
            </div>
            <el-input
              v-model="searchKeyword"
              placeholder="搜索主机名称"
              clearable
              @input="handleSearch"
              style="width: 200px; margin-left: 12px;"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </div>
      </template>
      
      <div class="alert-list">
        <div v-for="alert in filteredAlerts" :key="alert.id" class="alert-item">
          <div class="alert-icon" :class="alert.alert_type">
            <el-icon><component :is="getAlertIcon(alert.alert_type)" /></el-icon>
          </div>
          <div class="alert-content">
            <div class="alert-header">
              <div class="alert-title">{{ alert.title }}</div>
              <div class="alert-time">{{ formatTime(alert.created_at) }}</div>
            </div>
            <div class="alert-message">{{ alert.message }}</div>
            <div class="alert-meta" v-if="alert.host">
              <el-tag size="small" type="info">{{ alert.host.name }}</el-tag>
              <el-tag size="small" type="warning" v-if="alert.host.provider_type">
                {{ getProviderDisplayName(alert.host.provider_type) }}
              </el-tag>
              <el-tag size="small" v-if="alert.host.region">{{ alert.host.region }}</el-tag>
            </div>
          </div>
          <div class="alert-actions">
            <el-button size="small" @click="viewHostDetail(alert.host)" v-if="alert.host">
              查看主机
            </el-button>
            <el-button size="small" type="primary" @click="markAsResolved(alert)" v-if="!alert.is_resolved">
              标记已解决
            </el-button>
          </div>
        </div>
        
        <div v-if="filteredAlerts.length === 0" class="no-alerts">
          <el-empty description="暂无告警信息">
            <el-button type="primary" @click="triggerCheck">手动检查告警</el-button>
          </el-empty>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
  import { ref, computed, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage, ElLoading } from 'element-plus'
  import {
    Refresh, Search, Bell, Warning, CircleClose, InfoFilled, ArrowDown
  } from '@element-plus/icons-vue'
  import dayjs from 'dayjs'
  
  import { alertApi } from '@/api/resource'

  const router = useRouter()
  
  // 响应式数据
  const loading = ref(false)
  const filterLevel = ref('')
  const searchKeyword = ref('')
  
  // 告警数据
  const alerts = ref([])
  const alertStats = ref({
    total: 0,
    expired: 0,
    expiring: 0,
    error: 0,
    abnormal: 0
  })

  // 计算属性
  const totalAlerts = computed(() => alertStats.value.total)
  const expiredAlerts = computed(() => alertStats.value.expired)
  const expiringAlerts = computed(() => alertStats.value.expiring)
  const errorAlerts = computed(() => alertStats.value.error + alertStats.value.abnormal)

  const filteredAlerts = computed(() => {
    let filtered = alerts.value

    // 按级别过滤
    if (filterLevel.value) {
      filtered = filtered.filter(alert => alert.alert_type === filterLevel.value)
    }

    // 按关键词搜索
    if (searchKeyword.value) {
      const keyword = searchKeyword.value.toLowerCase()
      filtered = filtered.filter(alert => 
        alert.title.toLowerCase().includes(keyword) ||
        alert.message.toLowerCase().includes(keyword) ||
        (alert.host && alert.host.name.toLowerCase().includes(keyword))
      )
    }

    return filtered
  })

  // 生命周期
  onMounted(async () => {
    await refreshData()
  })

  // 方法
  const refreshData = async () => {
    try {
      loading.value = true
      
      const [alertsRes, statsRes] = await Promise.all([
        alertApi.getAlerts(),
        alertApi.getAlertStatistics()
      ])
      
      alerts.value = alertsRes.data || []
      alertStats.value = statsRes.data || {}
      
    } catch (error) {
      ElMessage.error('获取告警数据失败: ' + error.message)
    } finally {
      loading.value = false
    }
  }

  const triggerCheck = async () => {
    try {
      const loading = ElLoading.service({
        lock: true,
        text: '正在检查告警...',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      
      await alertApi.triggerAlertCheck()
      
      ElMessage.success('告警检查完成')
      await refreshData()
      loading.close()
    } catch (error) {
      ElMessage.error('告警检查失败: ' + error.message)
    }
  }

  const handleFilterChange = () => {
    // 过滤变化时的处理逻辑
  }

  const handleSearch = () => {
    // 搜索时的处理逻辑
  }

  const viewHostDetail = (host) => {
    if (host && host.id) {
      router.push(`/cmdb/host/detail/${host.id}`)
    }
  }

  const markAsResolved = (alert) => {
    ElMessage.info('标记已解决功能开发中...')
  }

  // 工具函数
  const getAlertIcon = (type) => {
    const iconMap = {
      'expired': 'CircleClose',
      'expiring': 'Warning',
      'error': 'CircleClose',
      'abnormal': 'Warning'
    }
    return iconMap[type] || 'InfoFilled'
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

  const formatTime = (time) => {
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
  }
</script>

<style lang="scss" scoped>
.alert-management {
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
          }
        }
      }

      &.total .stats-icon {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      }

      &.expired .stats-icon {
        background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
      }

      &.expiring .stats-icon {
        background: linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%);
        color: #d97706;
      }

      &.error .stats-icon {
        background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%);
        color: #dc2626;
      }
    }
  }

  .alert-list-card {
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

      .filter-controls {
        display: flex;
        align-items: center;
      }
    }

    .alert-list {
      .alert-item {
        display: flex;
        align-items: flex-start;
        gap: 16px;
        padding: 16px 0;
        border-bottom: 1px solid #f3f4f6;

        &:last-child {
          border-bottom: none;
        }

        .alert-icon {
          width: 40px;
          height: 40px;
          border-radius: 8px;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 18px;
          color: white;
          flex-shrink: 0;

          &.expired {
            background: #ef4444;
          }

          &.expiring {
            background: #f59e0b;
          }

          &.error {
            background: #dc2626;
          }

          &.abnormal {
            background: #f97316;
          }
        }

        .alert-content {
          flex: 1;

          .alert-header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: 8px;

            .alert-title {
              font-size: 16px;
              font-weight: 500;
              color: #374151;
            }

            .alert-time {
              font-size: 12px;
              color: #9ca3af;
              white-space: nowrap;
            }
          }

          .alert-message {
            font-size: 14px;
            color: #6b7280;
            margin-bottom: 8px;
            line-height: 1.5;
          }

          .alert-meta {
            display: flex;
            gap: 8px;
            flex-wrap: wrap;
          }
        }

        .alert-actions {
          display: flex;
          flex-direction: column;
          gap: 8px;
          flex-shrink: 0;
        }
      }

      .no-alerts {
        padding: 40px 0;
      }
    }
  }
}

// 全局样式覆盖
:deep(.el-card) {
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
}

@media (max-width: 768px) {
  .alert-management {
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

    .alert-list-card {
      .card-header {
        flex-direction: column;
        align-items: flex-start;
        gap: 16px;

        .filter-controls {
          width: 100%;
          justify-content: space-between;
        }
      }

      .alert-list .alert-item {
        flex-direction: column;
        align-items: flex-start;

        .alert-actions {
          flex-direction: row;
          width: 100%;
          justify-content: flex-end;
        }
      }
    }
  }
}

/* 自定义Select样式 */
.custom-select-wrapper {
  position: relative;
  width: 100%;
}

.custom-select {
  width: 100%;
  height: 40px;
  padding: 8px 32px 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  background-color: #ffffff;
  font-size: 14px;
  color: #606266;
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  cursor: pointer;
  transition: all 0.3s ease;
  outline: none;
}

.custom-select:hover {
  border-color: #c0c4cc;
}

.custom-select:focus {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.select-arrow {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  color: #c0c4cc;
  transition: transform 0.3s ease;
}

.custom-select:focus + .select-arrow {
  color: #409eff;
  transform: translateY(-50%) rotate(180deg);
}
</style>
