<template>
  <el-card class="stat-card" :class="[`stat-card-${variant}`]" shadow="never">
    <div class="stat-content">
      <div class="stat-icon">
        <el-icon :size="iconSize">
          <component :is="iconComponent" />
        </el-icon>
      </div>
      <div class="stat-info">
        <div class="stat-value" :style="valueStyle">
          <el-statistic 
            :value="value" 
            :precision="precision"
            :value-style="valueStyle"
          />
        </div>
        <div class="stat-label">{{ title }}</div>
        <div v-if="subtitle" class="stat-subtitle">{{ subtitle }}</div>
      </div>
      <div v-if="trend !== undefined" class="stat-trend">
        <el-icon :class="trendClass">
          <component :is="trendIcon" />
        </el-icon>
        <span :class="trendClass">{{ Math.abs(trend) }}%</span>
      </div>
    </div>
    <div v-if="loading" class="stat-loading">
      <el-skeleton :rows="1" animated />
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { 
  ArrowUp, 
  ArrowDown,
  Operation,
  CircleCheck,
  Warning,
  Clock,
  Lock,
  Setting,
  Monitor,
  Document
} from '@element-plus/icons-vue'

export interface Props {
  /** 标题 */
  title: string
  /** 数值 */
  value: number
  /** 副标题 */
  subtitle?: string
  /** 图标名称 */
  icon?: string
  /** 图标大小 */
  iconSize?: number
  /** 卡片变体样式 */
  variant?: 'primary' | 'success' | 'warning' | 'danger' | 'info'
  /** 数值精度 */
  precision?: number
  /** 趋势百分比 */
  trend?: number
  /** 加载状态 */
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  iconSize: 24,
  variant: 'primary',
  precision: 0,
  loading: false
})

// 图标映射
const iconMap = {
  operation: Operation,
  check: CircleCheck,
  warning: Warning,
  clock: Clock,
  lock: Lock,
  setting: Setting,
  monitor: Monitor,
  document: Document
}

// 计算属性
const iconComponent = computed(() => {
  if (!props.icon) return Operation
  return iconMap[props.icon as keyof typeof iconMap] || Operation
})

const valueStyle = computed(() => {
  const colorMap = {
    primary: '#409eff',
    success: '#67c23a', 
    warning: '#e6a23c',
    danger: '#f56c6c',
    info: '#909399'
  }
  return { color: colorMap[props.variant] }
})

const trendIcon = computed(() => {
  return props.trend && props.trend > 0 ? TrendChartUp : TrendChartDown
})

const trendClass = computed(() => {
  if (props.trend === undefined) return ''
  return props.trend > 0 ? 'trend-up' : 'trend-down'
})
</script>

<style scoped>
.stat-card {
  position: relative;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
  position: relative;
}

.stat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 8px;
  background: rgba(64, 158, 255, 0.1);
  color: #409eff;
}

.stat-card-primary .stat-icon {
  background: rgba(64, 158, 255, 0.1);
  color: #409eff;
}

.stat-card-success .stat-icon {
  background: rgba(103, 194, 58, 0.1);
  color: #67c23a;
}

.stat-card-warning .stat-icon {
  background: rgba(230, 162, 60, 0.1);
  color: #e6a23c;
}

.stat-card-danger .stat-icon {
  background: rgba(245, 108, 108, 0.1);
  color: #f56c6c;
}

.stat-card-info .stat-icon {
  background: rgba(144, 147, 153, 0.1);
  color: #909399;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #606266;
  line-height: 1;
}

.stat-subtitle {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.stat-trend {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  font-size: 12px;
}

.trend-up {
  color: #67c23a;
}

.trend-down {
  color: #f56c6c;
}

.stat-loading {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 8px;
  display: flex;
  align-items: center;
  padding: 16px;
}

:deep(.el-statistic__content) {
  display: flex;
  align-items: baseline;
}
</style>
