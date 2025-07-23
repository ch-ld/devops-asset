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
          <i class="el-icon-refresh"></i>
          重连
        </el-button>
        <el-button size="small" @click="handleClear">
          <i class="el-icon-delete"></i>
          清屏
        </el-button>
        <el-button size="small" @click="openNewWindow">
          <i class="el-icon-full-screen"></i>
          新窗口打开
        </el-button>
      </div>
    </div>

    <div class="terminal-container" ref="terminalContainer">
      <div class="terminal-output" ref="terminalOutput">
        <div v-for="(line, index) in outputLines" :key="index" class="output-line">
          {{ line }}
        </div>
      </div>
      <div class="terminal-input">
        <span class="prompt">{{ currentPrompt }}</span>
        <input
          ref="inputRef"
          v-model="currentInput"
          @keydown="handleKeydown"
          @focus="inputFocused = true"
          @blur="inputFocused = false"
          class="input-field"
          :disabled="connectionStatus !== 'connected'"
        />
        <span class="cursor" :class="{ blink: inputFocused }">|</span>
      </div>
    </div>

    <div class="terminal-footer">
      <div class="connection-info">
        <span>连接时间: {{ connectionTime }}</span>
        <span v-if="connectionStatus === 'connected'">
          | 延迟: {{ latency }}ms
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import type { Host } from '@/types/api/host'

interface Props {
  host: Host
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
}>()

const terminalOutput = ref<HTMLElement>()
const inputRef = ref<HTMLInputElement>()
const connecting = ref(false)
const connectionStatus = ref<'disconnected' | 'connecting' | 'connected' | 'error'>('disconnected')
const connectionTime = ref('')
const latency = ref(0)
const inputFocused = ref(false)
const currentInput = ref('')
const currentPrompt = ref('$ ')
const outputLines = ref<string[]>(['正在连接SSH...'])

let websocket: WebSocket | null = null
let pingInterval: number | null = null
let commandHistory: string[] = []
let historyIndex = -1

// 获取主机IP
const getHostIP = () => {
  if (Array.isArray(props.host.public_ip)) {
    return props.host.public_ip[0] || props.host.private_ip[0] || 'Unknown'
  }
  return props.host.public_ip || props.host.private_ip || 'Unknown'
}

// 状态相关
const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    connected: 'success',
    connecting: 'warning',
    disconnected: 'info',
    error: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    connected: '已连接',
    connecting: '连接中',
    disconnected: '未连接',
    error: '连接错误'
  }
  return statusMap[status] || status
}

// 处理键盘输入
const handleKeydown = (event: KeyboardEvent) => {
  if (connectionStatus.value !== 'connected') return

  switch (event.key) {
    case 'Enter':
      event.preventDefault()
      handleCommand()
      break
    case 'ArrowUp':
      event.preventDefault()
      navigateHistory(-1)
      break
    case 'ArrowDown':
      event.preventDefault()
      navigateHistory(1)
      break
    case 'Tab':
      event.preventDefault()
      // TODO: 实现命令补全
      break
  }
}

// 处理命令
const handleCommand = () => {
  const command = currentInput.value.trim()
  if (!command) return

  // 添加到历史记录
  commandHistory.unshift(command)
  if (commandHistory.length > 100) {
    commandHistory.pop()
  }
  historyIndex = -1

  // 显示命令
  outputLines.value.push(`${currentPrompt.value}${command}`)

  // 发送命令到服务器
  if (websocket && websocket.readyState === WebSocket.OPEN) {
    websocket.send(JSON.stringify({
      type: 'input',
      data: command + '\n'
    }))
  }

  currentInput.value = ''
  scrollToBottom()
}

// 历史记录导航
const navigateHistory = (direction: number) => {
  const newIndex = historyIndex + direction
  if (newIndex >= -1 && newIndex < commandHistory.length) {
    historyIndex = newIndex
    currentInput.value = historyIndex === -1 ? '' : commandHistory[historyIndex]
  }
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (terminalOutput.value) {
      terminalOutput.value.scrollTop = terminalOutput.value.scrollHeight
    }
  })
}

// 连接WebSocket
const connectWebSocket = () => {
  if (connecting.value) return

  connecting.value = true
  connectionStatus.value = 'connecting'

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${protocol}//${window.location.host}/api/v1/ws/ssh?host_id=${props.host.id}`

  websocket = new WebSocket(wsUrl)

  websocket.onopen = () => {
    connecting.value = false
    connectionStatus.value = 'connected'
    connectionTime.value = new Date().toLocaleTimeString()
    ElMessage.success('SSH连接成功')
    
    // 开始心跳检测
    startPing()
  }

  websocket.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data)

      if (message.type === 'output') {
        // 处理输出数据
        const lines = message.data.split('\n')
        lines.forEach((line: string, index: number) => {
          if (index === lines.length - 1 && line === '') return
          outputLines.value.push(line)
        })
        scrollToBottom()
      } else if (message.type === 'prompt') {
        currentPrompt.value = message.data
      } else if (message.type === 'pong') {
        const now = Date.now()
        latency.value = now - parseInt(message.timestamp)
      } else if (message.type === 'error') {
        ElMessage.error(`SSH错误: ${message.message}`)
        connectionStatus.value = 'error'
        outputLines.value.push(`错误: ${message.message}`)
      }
    } catch (error) {
      console.error('解析WebSocket消息失败:', error)
    }
  }

  websocket.onclose = () => {
    connecting.value = false
    connectionStatus.value = 'disconnected'
    stopPing()
    outputLines.value.push('连接已断开')
  }

  websocket.onerror = (error) => {
    connecting.value = false
    connectionStatus.value = 'error'
    ElMessage.error('SSH连接失败')
    console.error('WebSocket错误:', error)
  }
}

// 心跳检测
const startPing = () => {
  pingInterval = window.setInterval(() => {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
      websocket.send(JSON.stringify({
        type: 'ping',
        timestamp: Date.now().toString()
      }))
    }
  }, 30000) // 30秒心跳
}

const stopPing = () => {
  if (pingInterval) {
    clearInterval(pingInterval)
    pingInterval = null
  }
}

// 重连
const handleReconnect = () => {
  if (websocket) {
    websocket.close()
  }
  outputLines.value = ['正在重新连接...']
  setTimeout(() => {
    connectWebSocket()
  }, 1000)
}

// 清屏
const handleClear = () => {
  outputLines.value = []
}

// 新窗口打开
const openNewWindow = () => {
  const url = `/terminal?host=${props.host.id}`
  window.open(url, '_blank', 'width=1200,height=800')
}

// 初始化
onMounted(async () => {
  await nextTick()
  connectWebSocket()

  // 聚焦输入框
  if (inputRef.value) {
    inputRef.value.focus()
  }
})

// 清理
onUnmounted(() => {
  if (websocket) {
    websocket.close()
  }

  stopPing()
})
</script>

<style scoped>
.ssh-terminal {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #1e1e1e;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.terminal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #2d2d2d;
  border-bottom: 1px solid #404040;
}

.terminal-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.host-info {
  color: #d4d4d4;
  font-weight: 500;
}

.terminal-actions {
  display: flex;
  gap: 8px;
}

.terminal-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 16px;
  overflow: hidden;
}

.terminal-output {
  flex: 1;
  overflow-y: auto;
  margin-bottom: 8px;
  padding: 8px;
  background: #1e1e1e;
  border-radius: 4px;
}

.output-line {
  color: #d4d4d4;
  font-size: 14px;
  line-height: 1.4;
  white-space: pre-wrap;
  word-break: break-all;
  margin-bottom: 2px;
}

.terminal-input {
  display: flex;
  align-items: center;
  background: #2d2d2d;
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid #404040;
}

.prompt {
  color: #00ff00;
  font-weight: bold;
  margin-right: 8px;
  white-space: nowrap;
}

.input-field {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: #d4d4d4;
  font-family: inherit;
  font-size: 14px;
}

.input-field:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.cursor {
  color: #d4d4d4;
  margin-left: 2px;
}

.cursor.blink {
  animation: blink 1s infinite;
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}

.terminal-footer {
  padding: 8px 16px;
  background: #2d2d2d;
  border-top: 1px solid #404040;
}

.connection-info {
  font-size: 12px;
  color: #9ca3af;
}

.connection-info span {
  margin-right: 16px;
}

/* 滚动条样式 */
.terminal-output::-webkit-scrollbar {
  width: 8px;
}

.terminal-output::-webkit-scrollbar-track {
  background: #1e1e1e;
}

.terminal-output::-webkit-scrollbar-thumb {
  background: #404040;
  border-radius: 4px;
}

.terminal-output::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>
