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
