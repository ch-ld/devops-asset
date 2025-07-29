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

// 断开WebSocket连接
const disconnectWebSocket = () => {
  if (websocket) {
    websocket.close()
    websocket = null
  }
  stopPing()
}

// 开始心跳检测
const startPing = () => {
  pingInterval = setInterval(() => {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
      websocket.send('ping')
    }
  }, 30000)
}

// 停止心跳检测
const stopPing = () => {
  if (pingInterval) {
    clearInterval(pingInterval)
    pingInterval = null
  }
}

// 添加内容到终端
const appendToTerminal = (data) => {
  // 处理ANSI转义序列和特殊字符
  const processedData = data
    .replace(/\r\n/g, '<br>')
    .replace(/\n/g, '<br>')
    .replace(/\r/g, '<br>')
    .replace(/ /g, '&nbsp;')

  terminalContent.value += processedData
  scrollToBottom()
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (terminalOutput.value) {
      terminalOutput.value.scrollTop = terminalOutput.value.scrollHeight
    }
  })
}

// 处理键盘输入
const handleKeydown = (event) => {
  if (connectionStatus.value !== 'connected') return

  if (event.key === 'Enter') {
    event.preventDefault()
    sendCommand(currentInput.value + '\r')
    currentInput.value = ''
  } else if (event.key === 'Tab') {
    event.preventDefault()
    sendCommand('\t')
  } else if (event.ctrlKey && event.key === 'c') {
    event.preventDefault()
    sendCommand('\x03') // Ctrl+C
  } else if (event.ctrlKey && event.key === 'd') {
    event.preventDefault()
    sendCommand('\x04') // Ctrl+D
  }
}

// 发送命令
const sendCommand = (command) => {
  if (websocket && websocket.readyState === WebSocket.OPEN) {
    websocket.send(command)
  }
}

// 重连
const handleReconnect = () => {
  disconnectWebSocket()
  setTimeout(() => {
    connectWebSocket()
  }, 1000)
}

// 清屏
const handleClear = () => {
  terminalContent.value = ''
}

// 新窗口打开
const openNewWindow = () => {
  const url = `/ssh-terminal?host_id=${props.host.id}`
  window.open(url, '_blank', 'width=1200,height=800')
}

// 聚焦输入框
const focusInput = () => {
  nextTick(() => {
    if (inputRef.value) {
      inputRef.value.focus()
    }
  })
}
</script>

<style scoped>
.ssh-terminal {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #1e1e1e;
  color: #ffffff;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

.terminal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #2d2d2d;
  border-bottom: 1px solid #404040;
}

.terminal-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.host-info {
  font-weight: bold;
  color: #ffffff;
}

.terminal-actions {
  display: flex;
  gap: 8px;
}

.terminal-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.terminal-output {
  flex: 1;
  padding: 12px;
  overflow-y: auto;
  font-size: 14px;
  line-height: 1.4;
  white-space: pre-wrap;
  word-break: break-all;
}

.terminal-input {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background: #1e1e1e;
  border-top: 1px solid #404040;
}

.prompt {
  color: #00ff00;
  margin-right: 4px;
}

.input-field {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: #ffffff;
  font-family: inherit;
  font-size: 14px;
}

.cursor {
  color: #ffffff;
  margin-left: 2px;
}

.cursor.blink {
  animation: blink 1s infinite;
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}

/* 滚动条样式 */
.terminal-output::-webkit-scrollbar {
  width: 8px;
}

.terminal-output::-webkit-scrollbar-track {
  background: #2d2d2d;
}

.terminal-output::-webkit-scrollbar-thumb {
  background: #555;
  border-radius: 4px;
}

.terminal-output::-webkit-scrollbar-thumb:hover {
  background: #777;
}
</style>
