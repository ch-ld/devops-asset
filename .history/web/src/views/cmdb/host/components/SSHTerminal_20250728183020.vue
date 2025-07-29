<template>
  <div class="ssh-terminal">
    <div class="terminal-header">
      <div class="terminal-info">
        <span class="host-info">{{ host.name }} ({{ getHostIP() }})</span>
        <el-tag :type="getStatusType(connectionStatus)" size="small">
          {{ getStatusText(connectionStatus) }}
        </el-tag>
      </div>
      <div class="terminal-actions">
        <el-button size="small" @click="handleReconnect" :loading="connecting">
          <el-icon><Refresh /></el-icon>
          重连
        </el-button>
        <el-button size="small" @click="handleClear">
          <el-icon><Delete /></el-icon>
          清屏
        </el-button>
        <el-button size="small" @click="openNewWindow">
          <el-icon><FullScreen /></el-icon>
          新窗口打开
        </el-button>
      </div>
    </div>

    <div class="terminal-container" ref="terminalContainer">
      <div class="terminal-output" ref="terminalOutput" v-html="terminalContent"></div>
      <div class="terminal-input" v-if="connectionStatus === 'connected'">
        <span class="prompt">{{ currentPrompt }}</span>
        <input
          ref="inputRef"
          v-model="currentInput"
          @keydown="handleKeydown"
          @focus="inputFocused = true"
          @blur="inputFocused = false"
          class="input-field"
        />
        <span class="cursor" :class="{ blink: inputFocused }">|</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { ElMessage, ElTag, ElButton, ElIcon } from 'element-plus'
import { Refresh, Delete, FullScreen } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/modules/user'

const props = defineProps({
  host: {
    type: Object,
    required: true
  },
  visible: {
    type: Boolean,
    default: false
  }
})

const userStore = useUserStore()

// 响应式数据
const terminalContainer = ref(null)
const terminalOutput = ref(null)
const inputRef = ref(null)
const terminalContent = ref('')
const currentInput = ref('')
const currentPrompt = ref('$ ')
const connectionStatus = ref('disconnected') // disconnected, connecting, connected, error
const connecting = ref(false)
const inputFocused = ref(false)

let websocket = null
let pingInterval = null

// 监听visible变化
watch(() => props.visible, (newVal) => {
  if (newVal) {
    nextTick(() => {
      connectWebSocket()
      focusInput()
    })
  } else {
    disconnectWebSocket()
  }
})

// 组件挂载
onMounted(() => {
  if (props.visible) {
    connectWebSocket()
  }
})

// 组件卸载
onUnmounted(() => {
  disconnectWebSocket()
})

// 获取主机IP
const getHostIP = () => {
  return props.host.publicIP || props.host.privateIP || props.host.managementIP || 'Unknown'
}

// 获取连接状态类型
const getStatusType = (status) => {
  const statusMap = {
    'disconnected': 'info',
    'connecting': 'warning',
    'connected': 'success',
    'error': 'danger'
  }
  return statusMap[status] || 'info'
}

// 获取连接状态文本
const getStatusText = (status) => {
  const statusMap = {
    'disconnected': '未连接',
    'connecting': '连接中',
    'connected': '已连接',
    'error': '连接错误'
  }
  return statusMap[status] || '未知'
}

// 连接WebSocket
const connectWebSocket = () => {
  if (connecting.value || connectionStatus.value === 'connected') return

  connecting.value = true
  connectionStatus.value = 'connecting'

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const token = userStore.token
  const wsUrl = `${protocol}//${window.location.host}/api/v1/cmdb/ws/ssh?host_id=${props.host.id}&token=${token}`

  websocket = new WebSocket(wsUrl)

  websocket.onopen = () => {
    connecting.value = false
    connectionStatus.value = 'connected'
    ElMessage.success('SSH连接成功')
    startPing()
    focusInput()
  }

  websocket.onmessage = (event) => {
    const data = event.data
    appendToTerminal(data)
  }

  websocket.onclose = () => {
    connecting.value = false
    connectionStatus.value = 'disconnected'
    stopPing()
    appendToTerminal('\r\n[连接已断开]\r\n')
  }

  websocket.onerror = (error) => {
    connecting.value = false
    connectionStatus.value = 'error'
    stopPing()
    ElMessage.error('SSH连接失败')
    console.error('WebSocket error:', error)
  }
}
