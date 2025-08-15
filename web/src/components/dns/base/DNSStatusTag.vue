<template>
  <el-tag :type="tagType" :effect="effect" :size="size">
    <el-icon v-if="showIcon" class="status-icon">
      <component :is="iconComponent" />
    </el-icon>
    {{ statusText }}
  </el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { 
  CircleCheck, 
  Warning, 
  Clock, 
  CircleClose, 
  Loading,
  Lock,
  Unlock
} from '@element-plus/icons-vue'

export interface Props {
  /** 状态值 */
  status: string
  /** 状态类型：domain, record, certificate, provider, monitor */
  type?: 'domain' | 'record' | 'certificate' | 'provider' | 'monitor'
  /** 标签大小 */
  size?: 'small' | 'default' | 'large'
  /** 标签效果 */
  effect?: 'dark' | 'light' | 'plain'
  /** 是否显示图标 */
  showIcon?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'domain',
  size: 'small',
  effect: 'light',
  showIcon: true
})

// 状态映射配置
const statusConfig = {
  domain: {
    active: { text: '正常', type: 'success', icon: CircleCheck },
    inactive: { text: '未激活', type: 'info', icon: Clock },
    expired: { text: '已过期', type: 'danger', icon: Warning },
    expiring: { text: '即将过期', type: 'warning', icon: Clock },
    error: { text: '异常', type: 'danger', icon: CircleClose }
  },
  record: {
    active: { text: '正常', type: 'success', icon: CircleCheck },
    inactive: { text: '未激活', type: 'info', icon: Clock },
    syncing: { text: '同步中', type: 'warning', icon: Loading },
    error: { text: '同步失败', type: 'danger', icon: CircleClose },
    pending: { text: '待同步', type: 'info', icon: Clock }
  },
  certificate: {
    issued: { text: '已签发', type: 'success', icon: Lock },
    pending: { text: '申请中', type: 'warning', icon: Loading },
    expired: { text: '已过期', type: 'danger', icon: Unlock },
    revoked: { text: '已吊销', type: 'danger', icon: CircleClose },
    expiring: { text: '即将过期', type: 'warning', icon: Warning }
  },
  provider: {
    active: { text: '正常', type: 'success', icon: CircleCheck },
    inactive: { text: '未激活', type: 'info', icon: Clock },
    error: { text: '连接失败', type: 'danger', icon: CircleClose },
    testing: { text: '测试中', type: 'warning', icon: Loading }
  },
  monitor: {
    online: { text: '在线', type: 'success', icon: CircleCheck },
    offline: { text: '离线', type: 'danger', icon: CircleClose },
    warning: { text: '告警', type: 'warning', icon: Warning },
    checking: { text: '检测中', type: 'warning', icon: Loading }
  }
} as const

// 计算属性
const config = computed(() => {
  const typeConfig = statusConfig[props.type] || statusConfig.domain
  return typeConfig[props.status as keyof typeof typeConfig] || { 
    text: props.status, 
    type: 'info', 
    icon: Clock 
  }
})

const statusText = computed(() => config.value.text)
const tagType = computed(() => config.value.type as any)
const iconComponent = computed(() => config.value.icon)
</script>

<style scoped>
.status-icon {
  margin-right: 4px;
}
</style>
