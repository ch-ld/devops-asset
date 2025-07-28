<template>
  <div class="system-status-widget">
    <div class="widget-header">
      <h3 class="widget-title">
        <i class="iconfont-sys">&#xe642;</i>
        系统状态
      </h3>
      <div class="status-indicator" :class="overallStatus.toLowerCase()">
        <span class="status-dot"></span>
        {{ overallStatus }}
      </div>
    </div>
    
    <div class="status-metrics">
      <div class="metric-item">
        <div class="metric-header">
          <span class="metric-label">CPU使用率</span>
          <span class="metric-value">{{ systemMetrics.cpu }}%</span>
        </div>
        <el-progress 
          :percentage="systemMetrics.cpu" 
          :stroke-width="6"
          :color="getProgressColor(systemMetrics.cpu)"
          :show-text="false"
        />
      </div>
      
      <div class="metric-item">
        <div class="metric-header">
          <span class="metric-label">内存使用率</span>
          <span class="metric-value">{{ systemMetrics.memory }}%</span>
        </div>
        <el-progress 
          :percentage="systemMetrics.memory" 
          :stroke-width="6"
          :color="getProgressColor(systemMetrics.memory)"
          :show-text="false"
        />
      </div>
      
      <div class="metric-item">
        <div class="metric-header">
          <span class="metric-label">磁盘使用率</span>
          <span class="metric-value">{{ systemMetrics.disk }}%</span>
        </div>
        <el-progress 
          :percentage="systemMetrics.disk" 
          :stroke-width="6"
          :color="getProgressColor(systemMetrics.disk)"
          :show-text="false"
        />
      </div>
    </div>
    
    <div class="status-details">
      <div class="detail-grid">
        <div class="detail-item">
          <i class="iconfont-sys">&#xe650;</i>
          <div class="detail-content">
            <div class="detail-label">运行时间</div>
            <div class="detail-value">{{ systemMetrics.uptime }}</div>
          </div>
        </div>

        <div class="detail-item">
          <i class="iconfont-sys">&#xe651;</i>
          <div class="detail-content">
            <div class="detail-label">在线用户</div>
            <div class="detail-value">{{ systemMetrics.onlineUsers }}</div>
          </div>
        </div>

        <div class="detail-item">
          <i class="iconfont-sys">&#xe64f;</i>
          <div class="detail-content">
            <div class="detail-label">今日请求</div>
            <div class="detail-value">{{ systemMetrics.requests }}</div>
          </div>
        </div>

        <div class="detail-item">
          <i class="iconfont-sys">&#xe652;</i>
          <div class="detail-content">
            <div class="detail-label">错误数量</div>
            <div class="detail-value">{{ systemMetrics.errors }}</div>
          </div>
        </div>
      </div>
    </div>
    
    <div class="widget-footer">
      <el-button type="primary" size="small" @click="refreshStatus" :loading="refreshing" link>
        <i class="iconfont-sys">&#xe643;</i>
        刷新状态
      </el-button>
      <el-button type="primary" size="small" @click="viewDetails" link>
        <i class="iconfont-sys">&#xe653;</i>
        查看详情
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'

interface SystemMetrics {
  cpu: number
  memory: number
  disk: number
  uptime: string
  onlineUsers: number
  requests: number
  errors: number
}

const refreshing = ref(false)
const autoRefresh = ref<NodeJS.Timeout>()

const systemMetrics = ref<SystemMetrics>({
  cpu: 45,
  memory: 62,
  disk: 38,
  uptime: '7天 12小时 30分钟',
  onlineUsers: 24,
  requests: 1520,
  errors: 3
})

const overallStatus = computed(() => {
  const maxUsage = Math.max(systemMetrics.value.cpu, systemMetrics.value.memory, systemMetrics.value.disk)
  if (maxUsage > 80) return '警告'
  if (maxUsage > 60) return '注意'
  return '正常'
})

const getProgressColor = (percentage: number) => {
  if (percentage > 80) return '#ff4757'
  if (percentage > 60) return '#ffa726'
  return '#5dade2'
}

const refreshStatus = async () => {
  refreshing.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 模拟随机数据更新
    systemMetrics.value = {
      cpu: Math.floor(Math.random() * 80) + 10,
      memory: Math.floor(Math.random() * 70) + 20,
      disk: Math.floor(Math.random() * 60) + 30,
      uptime: '7天 12小时 31分钟',
      onlineUsers: Math.floor(Math.random() * 50) + 10,
      requests: Math.floor(Math.random() * 1000) + 1000,
      errors: Math.floor(Math.random() * 10)
    }
    
    ElMessage.success('状态已刷新')
  } catch (error) {
    ElMessage.error('刷新失败')
  } finally {
    refreshing.value = false
  }
}

const viewDetails = () => {
  ElMessage.info('详细监控功能开发中...')
}

const startAutoRefresh = () => {
  // 取消自动刷新功能
  // autoRefresh.value = setInterval(() => {
  //   refreshStatus()
  // }, 30000) // 30秒自动刷新
}

const stopAutoRefresh = () => {
  if (autoRefresh.value) {
    clearInterval(autoRefresh.value)
  }
}

onMounted(() => {
  // 不启动自动刷新
  // startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped lang="scss">
.system-status-widget {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  height: 100%;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  }
  
  .widget-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    
    .widget-title {
      font-size: 20px;
      font-weight: 700;
      color: #1a1a1a;
      margin: 0;
      display: flex;
      align-items: center;
      gap: 10px;
      background: linear-gradient(45deg, #4facfe, #00f2fe);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
      
      i {
        font-size: 22px;
        color: #4facfe;
      }
    }
    
    .status-indicator {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 13px;
      font-weight: 600;
      padding: 8px 16px;
      border-radius: 20px;
      
      .status-dot {
        width: 10px;
        height: 10px;
        border-radius: 50%;
      }
      
      &.正常 {
        background: rgba(87, 214, 141, 0.15);
        color: #57d68d;
        border: 1px solid rgba(87, 214, 141, 0.3);
        
        .status-dot {
          background: #57d68d;
          box-shadow: 0 0 8px rgba(87, 214, 141, 0.5);
        }
      }
      
      &.注意 {
        background: rgba(255, 167, 38, 0.15);
        color: #ffa726;
        border: 1px solid rgba(255, 167, 38, 0.3);
        
        .status-dot {
          background: #ffa726;
          box-shadow: 0 0 8px rgba(255, 167, 38, 0.5);
        }
      }
      
      &.警告 {
        background: rgba(255, 71, 87, 0.15);
        color: #ff4757;
        border: 1px solid rgba(255, 71, 87, 0.3);
        
        .status-dot {
          background: #ff4757;
          animation: pulse 2s infinite;
          box-shadow: 0 0 8px rgba(255, 71, 87, 0.5);
        }
      }
    }
  }
  
  .status-metrics {
    margin-bottom: 24px;
    
    .metric-item {
      margin-bottom: 20px;
      background: rgba(255, 255, 255, 0.6);
      padding: 16px;
      border-radius: 12px;
      border: 1px solid rgba(255, 255, 255, 0.3);
      transition: all 0.3s ease;
      
      &:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
        background: rgba(255, 255, 255, 0.9);
      }
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .metric-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;
        
        .metric-label {
          font-size: 15px;
          color: #666;
          font-weight: 600;
        }
        
        .metric-value {
          font-size: 16px;
          font-weight: 700;
          color: #333;
          background: linear-gradient(45deg, #4facfe, #00f2fe);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
        }
      }
      
      :deep(.el-progress-bar__outer) {
        background: #f0f0f0;
        border-radius: 8px;
        box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.1);
      }
      
      :deep(.el-progress-bar__inner) {
        border-radius: 8px;
        position: relative;
        
        &::after {
          content: '';
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
          animation: shimmer 2s infinite;
        }
      }
    }
  }
  
  .status-details {
    flex: 1;
    margin-bottom: 24px;
    
    .detail-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 16px;
      
      .detail-item {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 16px;
        background: rgba(255, 255, 255, 0.6);
        border-radius: 12px;
        border: 1px solid rgba(255, 255, 255, 0.3);
        transition: all 0.3s ease;
        
        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
          background: rgba(255, 255, 255, 0.9);
        }
        
        i {
          font-size: 18px;
          color: #4facfe;
          flex-shrink: 0;
          width: 24px;
          height: 24px;
          display: flex;
          align-items: center;
          justify-content: center;
          background: rgba(79, 172, 254, 0.1);
          border-radius: 50%;
        }
        
        .detail-content {
          flex: 1;
          min-width: 0;
          
          .detail-label {
            font-size: 13px;
            color: #999;
            margin-bottom: 4px;
            font-weight: 500;
          }
          
          .detail-value {
            font-size: 15px;
            font-weight: 700;
            color: #333;
          }
        }
      }
    }
  }
  
  .widget-footer {
    display: flex;
    justify-content: space-between;
    padding-top: 20px;
    border-top: 1px solid rgba(255, 255, 255, 0.3);
    gap: 12px;
    
    .el-button {
      flex: 1;
      background: linear-gradient(45deg, #4facfe, #00f2fe);
      border: none;
      color: white;
      padding: 8px 16px;
      border-radius: 8px;
      font-weight: 600;
      font-size: 13px;
      transition: all 0.3s ease;
      
      &:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(79, 172, 254, 0.4);
      }
    }
  }
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(255, 71, 87, 0.7);
    transform: scale(1);
  }
  70% {
    box-shadow: 0 0 0 6px rgba(255, 71, 87, 0);
    transform: scale(1.1);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(255, 71, 87, 0);
    transform: scale(1);
  }
}

@keyframes shimmer {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}
</style> 