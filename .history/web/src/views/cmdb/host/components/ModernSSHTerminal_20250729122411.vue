<template>
  <el-dialog
    v-model="dialogVisible"
    :before-close="handleClose"
    destroy-on-close
    class="modern-terminal-dialog"
    width="90%"
    top="5vh"
  >
    <template #header>
      <div class="terminal-dialog-header">
        <div class="header-left">
          <div class="terminal-icon">
            <el-icon size="24"><Monitor /></el-icon>
          </div>
          <div class="terminal-info">
            <h3 class="dialog-title">SSH 终端</h3>
            <p class="host-info">{{ props.host?.name || '未知主机' }} ({{ getHostIP() }})</p>
          </div>
        </div>
        <div class="header-actions">
          <div class="connection-status">
            <span class="status-dot" :class="connectionStatus"></span>
            <span class="status-text">{{ getStatusText(connectionStatus) }}</span>
          </div>
          <el-button circle size="small" @click="handleReconnect" :loading="connecting" class="action-btn">
            <el-icon><Refresh /></el-icon>
          </el-button>
          <el-button circle size="small" @click="handleClear" class="action-btn">
            <el-icon><Delete /></el-icon>
          </el-button>
          <el-button circle size="small" @click="openNewWindow" class="action-btn">
            <el-icon><FullScreen /></el-icon>
          </el-button>
        </div>
      </div>
    </template>

    <div class="modern-terminal">
      <!-- 终端标签栏 -->
      <div class="terminal-tabs">
        <div class="tab-list">
          <div 
            v-for="(tab, index) in terminalTabs" 
            :key="tab.id"
            class="terminal-tab"
            :class="{ active: activeTabId === tab.id }"
            @click="switchTab(tab.id)"
          >
            <span class="tab-title">{{ tab.title }}</span>
            <el-icon class="tab-close" @click.stop="closeTab(tab.id)" v-if="terminalTabs.length > 1">
              <Close />
            </el-icon>
          </div>
          <el-button size="small" @click="addNewTab" class="add-tab-btn">
            <el-icon><Plus /></el-icon>
          </el-button>
        </div>
        <div class="terminal-controls">
          <el-button-group size="small">
            <el-button @click="increaseFontSize">
              <el-icon><ZoomIn /></el-icon>
            </el-button>
            <el-button @click="decreaseFontSize">
              <el-icon><ZoomOut /></el-icon>
            </el-button>
            <el-button @click="resetFontSize">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </el-button-group>
        </div>
      </div>

      <!-- 终端内容区域 -->
      <div class="terminal-content" ref="terminalContainer">
        <div 
          v-for="tab in terminalTabs" 
          :key="tab.id"
          v-show="activeTabId === tab.id"
          class="terminal-pane"
          :style="{ fontSize: fontSize + 'px' }"
        >
          <div class="terminal-output" ref="terminalOutput" v-html="getTabContent(tab.id)"></div>
          <div class="terminal-input-line" v-if="connectionStatus === 'connected'">
            <span class="prompt">{{ currentPrompt }}</span>
            <input
              ref="inputRef"
              v-model="currentInput"
              @keydown="handleKeydown"
              @focus="inputFocused = true"
              @blur="inputFocused = false"
              class="terminal-input"
              autocomplete="off"
              spellcheck="false"
            />
            <span class="cursor" :class="{ blink: inputFocused }">█</span>
          </div>
        </div>

        <!-- 连接状态提示 -->
        <div v-if="connectionStatus === 'disconnected'" class="connection-prompt">
          <div class="prompt-content">
            <el-icon size="48" class="prompt-icon"><Monitor /></el-icon>
            <h3>SSH 终端未连接</h3>
            <p>点击连接按钮开始SSH会话</p>
            <el-button type="primary" @click="handleReconnect" :loading="connecting">
              <el-icon><Refresh /></el-icon>
              {{ connecting ? '连接中...' : '立即连接' }}
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Monitor, Refresh, Delete, FullScreen, Close, Plus, 
  ZoomIn, ZoomOut 
} from '@element-plus/icons-vue'
import { useUserStore } from '@/store/modules/user'

// 定义接口
interface TerminalTab {
  id: string
  title: string
  content: string
}

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

const emit = defineEmits(['update:visible'])

const userStore = useUserStore()

// 响应式数据
const dialogVisible = ref(false)
const terminalContainer = ref<HTMLElement>()
const terminalOutput = ref<HTMLElement>()
const inputRef = ref<HTMLInputElement>()
const currentInput = ref('')
const currentPrompt = ref('$ ')
const connectionStatus = ref<'disconnected' | 'connecting' | 'connected' | 'error'>('disconnected')
const connecting = ref(false)
const inputFocused = ref(false)
const fontSize = ref(14)

// 终端标签相关
const terminalTabs = ref<TerminalTab[]>([
  { id: '1', title: 'Terminal 1', content: '' }
])
const activeTabId = ref('1')
const tabCounter = ref(1)

let websocket: WebSocket | null = null
let pingInterval: NodeJS.Timeout | null = null

// 监听visible变化
watch(() => props.visible, (newVal) => {
  dialogVisible.value = newVal
  if (newVal) {
    nextTick(() => {
      connectWebSocket()
      focusInput()
    })
  } else {
    disconnectWebSocket()
  }
})

// 监听dialogVisible变化
watch(dialogVisible, (newVal) => {
  emit('update:visible', newVal)
})

// 计算属性
const getHostIP = () => {
  return props.host?.publicIP || props.host?.privateIP || '未知IP'
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'connected': return '已连接'
    case 'connecting': return '连接中...'
    case 'error': return '连接错误'
    default: return '未连接'
  }
}

// 方法
const handleClose = () => {
  dialogVisible.value = false
}

const getTabContent = (tabId: string) => {
  const tab = terminalTabs.value.find(t => t.id === tabId)
  return tab?.content || ''
}

const switchTab = (tabId: string) => {
  activeTabId.value = tabId
}

const addNewTab = () => {
  tabCounter.value++
  const newTab: TerminalTab = {
    id: tabCounter.value.toString(),
    title: `Terminal ${tabCounter.value}`,
    content: ''
  }
  terminalTabs.value.push(newTab)
  activeTabId.value = newTab.id
}

const closeTab = (tabId: string) => {
  if (terminalTabs.value.length <= 1) return
  
  const index = terminalTabs.value.findIndex(t => t.id === tabId)
  if (index > -1) {
    terminalTabs.value.splice(index, 1)
    if (activeTabId.value === tabId) {
      activeTabId.value = terminalTabs.value[Math.max(0, index - 1)].id
    }
  }
}

const increaseFontSize = () => {
  if (fontSize.value < 24) fontSize.value++
}

const decreaseFontSize = () => {
  if (fontSize.value > 10) fontSize.value--
}

const resetFontSize = () => {
  fontSize.value = 14
}

const handleReconnect = () => {
  connectWebSocket()
}

const handleClear = () => {
  const activeTab = terminalTabs.value.find(t => t.id === activeTabId.value)
  if (activeTab) {
    activeTab.content = ''
  }
}

const openNewWindow = () => {
  // 实现新窗口打开功能
  ElMessage.info('新窗口功能开发中...')
}

const focusInput = () => {
  nextTick(() => {
    inputRef.value?.focus()
  })
}

const handleKeydown = (event: KeyboardEvent) => {
  // 处理键盘事件
  if (event.key === 'Enter') {
    sendCommand()
  }
}

const sendCommand = () => {
  if (!currentInput.value.trim()) return
  
  // 发送命令到WebSocket
  if (websocket && websocket.readyState === WebSocket.OPEN) {
    websocket.send(currentInput.value + '\n')
    currentInput.value = ''
  }
}

const connectWebSocket = () => {
  if (connecting.value) return
  
  connecting.value = true
  connectionStatus.value = 'connecting'
  
  // 这里需要根据实际的WebSocket地址进行连接
  // websocket = new WebSocket(`ws://localhost:8080/api/v1/cmdb/ssh/ws?host_id=${props.host.id}`)
  
  // 模拟连接过程
  setTimeout(() => {
    connecting.value = false
    connectionStatus.value = 'connected'
    ElMessage.success('SSH连接成功')
  }, 2000)
}

const disconnectWebSocket = () => {
  if (websocket) {
    websocket.close()
    websocket = null
  }
  if (pingInterval) {
    clearInterval(pingInterval)
    pingInterval = null
  }
  connectionStatus.value = 'disconnected'
}

// 组件挂载
onMounted(() => {
  if (props.visible) {
    dialogVisible.value = true
    connectWebSocket()
  }
})

// 组件卸载
onUnmounted(() => {
  disconnectWebSocket()
})
</script>

<style scoped>
/* 现代化终端对话框样式 */
:deep(.modern-terminal-dialog) {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.modern-terminal-dialog .el-dialog__header) {
  padding: 0;
  margin: 0;
  border-bottom: none;
}

:deep(.modern-terminal-dialog .el-dialog__body) {
  padding: 0;
  height: 75vh;
  overflow: hidden;
  background: #1a1a1a;
}

/* 终端对话框头部 */
.terminal-dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
  color: white;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.terminal-icon {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.terminal-info h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.host-info {
  margin: 4px 0 0 0;
  font-size: 14px;
  opacity: 0.9;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  font-size: 13px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #ef4444;
}

.status-dot.connected {
  background: #10b981;
}

.status-dot.connecting {
  background: #f59e0b;
  animation: pulse 1.5s infinite;
}

.status-dot.error {
  background: #ef4444;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.action-btn {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: white;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.3);
}

/* 现代化终端主体 */
.modern-terminal {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #1a1a1a;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

/* 终端标签栏 */
.terminal-tabs {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #2d2d2d;
  border-bottom: 1px solid #404040;
  padding: 8px 16px;
}

.tab-list {
  display: flex;
  align-items: center;
  gap: 4px;
}

.terminal-tab {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #404040;
  border-radius: 8px 8px 0 0;
  cursor: pointer;
  transition: all 0.2s;
  color: #a0a0a0;
  font-size: 13px;
}

.terminal-tab:hover {
  background: #4a4a4a;
  color: #ffffff;
}

.terminal-tab.active {
  background: #1a1a1a;
  color: #ffffff;
  border-bottom: 2px solid #10b981;
}

.tab-title {
  font-weight: 500;
}

.tab-close {
  font-size: 12px;
  opacity: 0.6;
  transition: opacity 0.2s;
}

.tab-close:hover {
  opacity: 1;
  color: #ef4444;
}

.add-tab-btn {
  background: transparent;
  border: 1px solid #404040;
  color: #a0a0a0;
  padding: 6px;
}

.add-tab-btn:hover {
  background: #404040;
  border-color: #4a4a4a;
  color: #ffffff;
}

.terminal-controls {
  display: flex;
  gap: 8px;
}

/* 终端内容区域 */
.terminal-content {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.terminal-pane {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #1a1a1a;
  color: #ffffff;
}

.terminal-output {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  line-height: 1.4;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.terminal-input-line {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: #1a1a1a;
  border-top: 1px solid #2d2d2d;
}

.prompt {
  color: #10b981;
  font-weight: bold;
  margin-right: 8px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

.terminal-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: #ffffff;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: inherit;
}

.cursor {
  color: #10b981;
  font-weight: bold;
  margin-left: 2px;
}

.cursor.blink {
  animation: blink 1s infinite;
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}

/* 连接状态提示 */
.connection-prompt {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #1a1a1a;
}

.prompt-content {
  text-align: center;
  color: #a0a0a0;
}

.prompt-icon {
  color: #4a4a4a;
  margin-bottom: 16px;
}

.prompt-content h3 {
  margin: 0 0 8px 0;
  color: #ffffff;
  font-size: 18px;
  font-weight: 500;
}

.prompt-content p {
  margin: 0 0 24px 0;
  font-size: 14px;
}

/* 滚动条样式 */
.terminal-output::-webkit-scrollbar {
  width: 8px;
}

.terminal-output::-webkit-scrollbar-track {
  background: #2d2d2d;
}

.terminal-output::-webkit-scrollbar-thumb {
  background: #4a4a4a;
  border-radius: 4px;
}

.terminal-output::-webkit-scrollbar-thumb:hover {
  background: #5a5a5a;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .terminal-tabs {
    padding: 6px 12px;
  }

  .terminal-tab {
    padding: 6px 12px;
    font-size: 12px;
  }

  .terminal-output {
    padding: 12px;
    font-size: 13px;
  }

  .terminal-input-line {
    padding: 8px 12px;
  }
}
</style>
