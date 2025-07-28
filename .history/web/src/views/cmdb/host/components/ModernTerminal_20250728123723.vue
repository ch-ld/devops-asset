<template>
  <div class="modern-terminal">
    <el-dialog
      :title="terminalTitle"
      v-model="dialogVisible"
      :width="isFullscreen ? '100%' : '1200px'"
      :fullscreen="isFullscreen"
      :before-close="handleClose"
      destroy-on-close
      class="terminal-dialog"
    >
      <!-- 终端工具栏 -->
      <div class="terminal-toolbar">
        <div class="terminal-info">
          <div class="host-info">
            <div class="host-avatar">
              <el-icon class="host-icon"><Monitor /></el-icon>
            </div>
            <div class="host-details">
              <span class="host-name">{{ host?.name || 'Unknown' }}</span>
              <span class="host-ip">{{ primaryIP }}</span>
            </div>
          </div>
          <div class="connection-status">
            <div class="status-indicator">
              <div :class="['status-dot', connectionStatus]"></div>
              <el-tag :type="getStatusType(connectionStatus)" size="small" class="status-tag">
                <el-icon class="status-icon">
                  <Loading v-if="connectionStatus === 'connecting'" />
                  <CircleCheck v-else-if="connectionStatus === 'connected'" />
                  <CircleClose v-else-if="connectionStatus === 'error'" />
                  <Remove v-else />
                </el-icon>
                {{ getStatusText(connectionStatus) }}
              </el-tag>
            </div>
            <div v-if="connectionStatus === 'connected'" class="connection-details">
              <span class="connection-time">{{ connectionTime }}</span>
              <span class="ping-info">延迟: {{ ping }}ms</span>
            </div>
          </div>
        </div>
        
        <div class="terminal-actions">
          <el-button-group class="action-group">
            <el-tooltip content="重新连接" placement="bottom">
              <el-button size="small" @click="handleReconnect" :loading="connecting" class="action-btn">
                <el-icon><Refresh /></el-icon>
                重连
              </el-button>
            </el-tooltip>
            <el-tooltip content="清空终端" placement="bottom">
              <el-button size="small" @click="handleClear" class="action-btn">
                <el-icon><Delete /></el-icon>
                清屏
              </el-button>
            </el-tooltip>
            <el-tooltip content="复制选中内容" placement="bottom">
              <el-button size="small" @click="handleCopy" class="action-btn">
                <el-icon><CopyDocument /></el-icon>
                复制
              </el-button>
            </el-tooltip>
            <el-tooltip content="粘贴" placement="bottom">
              <el-button size="small" @click="handlePaste" class="action-btn">
                <el-icon><DocumentCopy /></el-icon>
                粘贴
              </el-button>
            </el-tooltip>
            <el-tooltip :content="isFullscreen ? '退出全屏' : '全屏模式'" placement="bottom">
              <el-button size="small" @click="toggleFullscreen" class="action-btn">
                <el-icon>
                  <FullScreen v-if="!isFullscreen" />
                  <Close v-else />
                </el-icon>
                {{ isFullscreen ? '退出' : '全屏' }}
              </el-button>
            </el-tooltip>
          </el-button-group>
        </div>
      </div>

      <!-- 终端容器 -->
      <div class="terminal-container" ref="terminalContainer">
        <div v-if="connectionStatus === 'disconnected'" class="terminal-placeholder">
          <el-empty description="点击连接按钮开始SSH会话">
            <el-button type="primary" @click="connectWebSocket" :loading="connecting">
              <el-icon><Connection /></el-icon>
              连接SSH
            </el-button>
          </el-empty>
        </div>
        <div v-else-if="connectionStatus === 'error'" class="terminal-placeholder error">
          <el-empty description="连接失败，请检查主机配置">
            <el-button type="primary" @click="handleReconnect" :loading="connecting">
              <el-icon><Refresh /></el-icon>
              重新连接
            </el-button>
          </el-empty>
        </div>
        <div v-else class="xterm-container" ref="xtermContainer"></div>
      </div>

      <!-- 状态栏 -->
      <div class="terminal-statusbar">
        <div class="status-left">
          <span class="status-item">
            <el-icon><Timer /></el-icon>
            延迟: {{ latency }}ms
          </span>
          <span class="status-item" v-if="terminalSize.cols && terminalSize.rows">
            <el-icon><Grid /></el-icon>
            {{ terminalSize.cols }}x{{ terminalSize.rows }}
          </span>
        </div>
        <div class="status-right">
          <span class="status-item">
            <el-icon><Key /></el-icon>
            {{ host?.username || 'root' }}@{{ primaryIP }}
          </span>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Monitor, Loading, CircleCheck, CircleClose, Remove,
  Refresh, Delete, FullScreen, Close, Connection,
  Timer, Grid, Key, CopyDocument, DocumentCopy
} from '@element-plus/icons-vue'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { WebLinksAddon } from 'xterm-addon-web-links'
import { SearchAddon } from 'xterm-addon-search'
import { Unicode11Addon } from 'xterm-addon-unicode11'
import 'xterm/css/xterm.css'
import { useUserStore } from '@/store/modules/user'

// Props
interface Props {
  visible: boolean
  host: any
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  host: null
})

// Emits
const emit = defineEmits<{
  'update:visible': [value: boolean]
}>()

// 响应式数据
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const terminalContainer = ref<HTMLElement>()
const xtermContainer = ref<HTMLElement>()
const isFullscreen = ref(false)
const connecting = ref(false)
const connectionStatus = ref<'disconnected' | 'connecting' | 'connected' | 'error'>('disconnected')
const connectionTime = ref('')
const ping = ref(0)
const terminalSize = ref({ cols: 80, rows: 24 })

// 终端相关
let terminal: Terminal | null = null
let fitAddon: FitAddon | null = null
let websocket: WebSocket | null = null
let pingInterval: number | null = null

// 计算属性
const terminalTitle = computed(() => {
  if (!props.host) return 'SSH终端'
  return `SSH终端 - ${props.host.name}`
})

const primaryIP = computed(() => {
  if (!props.host) return 'Unknown'
  
  // 处理IP地址数组
  if (Array.isArray(props.host.public_ip) && props.host.public_ip.length > 0) {
    return props.host.public_ip[0]
  }
  if (Array.isArray(props.host.private_ip) && props.host.private_ip.length > 0) {
    return props.host.private_ip[0]
  }
  
  // 处理字符串IP
  return props.host.public_ip || props.host.private_ip || 'Unknown'
})

// 状态相关方法
const getStatusType = (status: string) => {
  switch (status) {
    case 'connected': return 'success'
    case 'connecting': return 'warning'
    case 'error': return 'danger'
    default: return 'info'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'connected': return '已连接'
    case 'connecting': return '连接中'
    case 'error': return '连接失败'
    default: return '未连接'
  }
}

// 初始化终端
const initTerminal = () => {
  if (!xtermContainer.value) return

  // 创建终端实例
  terminal = new Terminal({
    fontSize: 14,
    fontFamily: 'Menlo, Monaco, "Courier New", monospace',
    cursorBlink: true,
    cursorStyle: 'block',
    scrollback: 1000,
    convertEol: true,
    allowTransparency: true,
    theme: {
      background: '#1e1e1e',
      foreground: '#ffffff',
      cursor: '#ffffff',
      selection: '#3a3a3a'
    }
  })

  // 添加插件
  fitAddon = new FitAddon()
  const webLinksAddon = new WebLinksAddon()
  const searchAddon = new SearchAddon()
  const unicode11Addon = new Unicode11Addon()

  terminal.loadAddon(fitAddon)
  terminal.loadAddon(webLinksAddon)
  terminal.loadAddon(searchAddon)
  terminal.loadAddon(unicode11Addon)

  // 打开终端
  terminal.open(xtermContainer.value)

  // 自适应大小
  setTimeout(() => {
    fitAddon?.fit()
    terminalSize.value = {
      cols: terminal?.cols || 80,
      rows: terminal?.rows || 24
    }
  }, 100)

  // 处理窗口大小变化
  const handleResize = () => {
    if (fitAddon && terminal) {
      fitAddon.fit()
      terminalSize.value = {
        cols: terminal.cols,
        rows: terminal.rows
      }
      
      // 通知服务器终端大小变化
      if (websocket && websocket.readyState === WebSocket.OPEN) {
        websocket.send(JSON.stringify({
          type: 'resize',
          cols: terminal.cols,
          rows: terminal.rows
        }))
      }
    }
  }

  window.addEventListener('resize', handleResize)

  // 处理用户输入
  terminal.onData((data) => {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
      websocket.send(data)
    }
  })

  // 右键菜单支持复制粘贴
  terminal.attachCustomKeyEventHandler((event) => {
    // Ctrl+C 复制
    if (event.ctrlKey && event.key === 'c' && terminal?.hasSelection()) {
      handleCopy()
      return false
    }
    // Ctrl+V 粘贴
    if (event.ctrlKey && event.key === 'v') {
      handlePaste()
      return false
    }
    return true
  })
}

// WebSocket连接
const connectWebSocket = () => {
  if (!props.host || !props.host.id) {
    ElMessage.error('无法连接：主机信息不完整')
    return
  }

  if (connecting.value) return

  connecting.value = true
  connectionStatus.value = 'connecting'

  // 关闭现有连接
  if (websocket && websocket.readyState !== WebSocket.CLOSED) {
    websocket.close()
  }

  // 获取认证信息
  const userStore = useUserStore()
  const token = userStore.accessToken

  // 创建WebSocket连接
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  // 修复WebSocket连接地址，使用后端端口
  const backendHost = window.location.hostname + ':8080'
  const wsUrl = `${protocol}//${backendHost}/api/v1/ws/terminal?host_id=${props.host.id}&token=${token}`

  websocket = new WebSocket(wsUrl)

  websocket.onopen = () => {
    connecting.value = false
    connectionStatus.value = 'connected'
    connectionTime.value = new Date().toLocaleTimeString()
    ElMessage.success('SSH连接成功')
    
    // 发送终端大小信息
    if (terminal) {
      websocket?.send(JSON.stringify({
        type: 'resize',
        cols: terminal.cols,
        rows: terminal.rows
      }))
    }
    
    // 开始心跳检测
    startPing()
  }

  websocket.onmessage = (event) => {
    if (terminal) {
      // 直接将数据写入终端
      terminal.write(event.data)
    }
  }

  websocket.onclose = () => {
    connecting.value = false
    connectionStatus.value = 'disconnected'
    stopPing()
    
    if (terminal) {
      terminal.write('\r\n\x1b[31m连接已断开\x1b[0m\r\n')
    }
  }

  websocket.onerror = (error) => {
    connecting.value = false
    connectionStatus.value = 'error'
    ElMessage.error('SSH连接失败')
    console.error('WebSocket错误:', error)
    
    if (terminal) {
      terminal.write('\r\n\x1b[31m连接失败\x1b[0m\r\n')
    }
  }
}

// 心跳检测
const startPing = () => {
  pingInterval = window.setInterval(() => {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
      const pingTime = Date.now()
      websocket.send(JSON.stringify({
        type: 'ping',
        timestamp: pingTime.toString()
      }))

      // 计算延迟（简单实现）
      setTimeout(() => {
        ping.value = Date.now() - pingTime
      }, 100)
    }
  }, 30000) // 30秒心跳
}

const stopPing = () => {
  if (pingInterval) {
    clearInterval(pingInterval)
    pingInterval = null
  }
}

// 操作方法
const handleReconnect = () => {
  if (websocket) {
    websocket.close()
  }

  if (terminal) {
    terminal.clear()
    terminal.write('\x1b[33m正在重新连接...\x1b[0m\r\n')
  }

  setTimeout(() => {
    connectWebSocket()
  }, 1000)
}

const handleClear = () => {
  if (terminal) {
    terminal.clear()
  }
}

// 复制选中内容
const handleCopy = async () => {
  if (terminal && terminal.hasSelection()) {
    const selection = terminal.getSelection()
    try {
      await navigator.clipboard.writeText(selection)
      ElMessage.success('已复制到剪贴板')
    } catch (err) {
      console.error('复制失败:', err)
      ElMessage.error('复制失败')
    }
  } else {
    ElMessage.warning('请先选择要复制的内容')
  }
}

// 粘贴剪贴板内容
const handlePaste = async () => {
  try {
    const text = await navigator.clipboard.readText()
    if (text && websocket && websocket.readyState === WebSocket.OPEN) {
      websocket.send(text)
      ElMessage.success('已粘贴')
    }
  } catch (err) {
    console.error('粘贴失败:', err)
    ElMessage.error('粘贴失败')
  }
}

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value

  // 延迟调整终端大小
  nextTick(() => {
    setTimeout(() => {
      if (fitAddon && terminal) {
        fitAddon.fit()
        terminalSize.value = {
          cols: terminal.cols,
          rows: terminal.rows
        }

        // 通知服务器终端大小变化
        if (websocket && websocket.readyState === WebSocket.OPEN) {
          websocket.send(JSON.stringify({
            type: 'resize',
            cols: terminal.cols,
            rows: terminal.rows
          }))
        }
      }
    }, 300)
  })
}

const handleClose = () => {
  if (websocket) {
    websocket.close()
  }

  stopPing()

  if (terminal) {
    terminal.dispose()
    terminal = null
  }

  dialogVisible.value = false
}

// 监听对话框显示状态
watch(() => props.visible, (newVal) => {
  if (newVal && props.host) {
    nextTick(() => {
      initTerminal()
    })
  }
})

// 生命周期
onMounted(() => {
  if (props.visible && props.host) {
    nextTick(() => {
      initTerminal()
    })
  }
})

onUnmounted(() => {
  handleClose()
})
</script>

<style scoped>
.modern-terminal {
  .terminal-dialog {
    :deep(.el-dialog__body) {
      padding: 0;
    }
  }

  .terminal-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    background: linear-gradient(135deg, #2c3e50 0%, #34495e 100%);
    color: white;
    border-radius: 12px 12px 0 0;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

    .terminal-info {
      display: flex;
      flex-direction: column;
      gap: 8px;

      .host-info {
        display: flex;
        align-items: center;
        gap: 12px;
        font-weight: 600;

        .host-avatar {
          width: 40px;
          height: 40px;
          background: linear-gradient(135deg, #3498db, #2980b9);
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          box-shadow: 0 2px 8px rgba(52, 152, 219, 0.3);

          .host-icon {
            font-size: 18px;
            color: white;
          }
        }

        .host-details {
          display: flex;
          flex-direction: column;
          gap: 2px;

          .host-name {
            font-size: 16px;
            font-weight: 600;
          }

          .host-ip {
            font-size: 12px;
            opacity: 0.8;
            background: rgba(255, 255, 255, 0.15);
            padding: 4px 8px;
            border-radius: 12px;
            font-family: 'Monaco', 'Menlo', monospace;
          }
        }
      }

      .connection-status {
        display: flex;
        flex-direction: column;
        gap: 6px;
        font-size: 12px;

        .status-indicator {
          display: flex;
          align-items: center;
          gap: 8px;

          .status-dot {
            width: 8px;
            height: 8px;
            border-radius: 50%;
            animation: pulse 2s infinite;

            &.connected {
              background: #27ae60;
              box-shadow: 0 0 8px rgba(39, 174, 96, 0.6);
            }

            &.connecting {
              background: #f39c12;
              box-shadow: 0 0 8px rgba(243, 156, 18, 0.6);
            }

            &.error {
              background: #e74c3c;
              box-shadow: 0 0 8px rgba(231, 76, 60, 0.6);
            }

            &.disconnected {
              background: #95a5a6;
            }
          }

          .status-tag {
            background: rgba(255, 255, 255, 0.1);
            border: 1px solid rgba(255, 255, 255, 0.2);

            .status-icon {
              font-size: 12px;
            }
          }
        }

        .connection-details {
          display: flex;
          gap: 12px;
          font-size: 11px;
          opacity: 0.8;

          .connection-time,
          .ping-info {
            background: rgba(255, 255, 255, 0.1);
            padding: 2px 6px;
            border-radius: 8px;
            font-family: 'Monaco', 'Menlo', monospace;
          }
        }
      }
    }

    .terminal-actions {
      .action-group {
        display: flex;
        gap: 8px;

        :deep(.el-button) {
          background: rgba(255, 255, 255, 0.1);
          border: 1px solid rgba(255, 255, 255, 0.2);
          color: white;
          border-radius: 8px;
          padding: 8px 12px;
          font-size: 12px;
          transition: all 0.3s ease;
          backdrop-filter: blur(10px);

          &:hover {
            background: rgba(255, 255, 255, 0.2);
            border-color: rgba(255, 255, 255, 0.4);
            transform: translateY(-1px);
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
          }

          &:active {
            transform: translateY(0);
          }

          .el-icon {
            margin-right: 4px;
          }
        }
      }
    }
  }

  .terminal-container {
    height: 500px;
    background: linear-gradient(135deg, #0c0c0c 0%, #1a1a1a 100%);
    position: relative;
    border-radius: 0 0 12px 12px;
    overflow: hidden;
    box-shadow: inset 0 2px 8px rgba(0, 0, 0, 0.3);

    .terminal-placeholder {
      height: 100%;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      background: linear-gradient(135deg, #0c0c0c 0%, #1a1a1a 100%);
      color: #ffffff;
      gap: 16px;

      &.error {
        background: linear-gradient(135deg, #2d1b1b 0%, #3d2b2b 100%);
      }

      .el-icon {
        font-size: 48px;
        opacity: 0.6;
        animation: float 3s ease-in-out infinite;
      }

      :deep(.el-empty__description) {
        color: #ffffff;
        font-size: 14px;
        opacity: 0.8;
      }
    }

    .xterm-container {
      height: 100%;
      padding: 12px;
      background: transparent;

      :deep(.xterm-viewport) {
        background: transparent !important;
      }

      :deep(.xterm-screen) {
        background: transparent !important;
      }
    }
  }

  .terminal-statusbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 20px;
    background: linear-gradient(135deg, #ecf0f1 0%, #bdc3c7 100%);
    border-top: 1px solid rgba(0, 0, 0, 0.1);
    font-size: 11px;
    color: #2c3e50;
    font-family: 'Monaco', 'Menlo', monospace;

    .status-left,
    .status-right {
      display: flex;
      gap: 20px;
    }

    .status-item {
      display: flex;
      align-items: center;
      gap: 6px;
      background: rgba(255, 255, 255, 0.7);
      padding: 4px 8px;
      border-radius: 8px;
      transition: all 0.3s ease;

      &:hover {
        background: rgba(255, 255, 255, 0.9);
        transform: translateY(-1px);
      }

      .el-icon {
        font-size: 12px;
        opacity: 0.7;
      }
    }
  }
}

/* 动画效果 */
@keyframes pulse {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
  100% {
    opacity: 1;
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-10px);
  }
}

/* 滚动条美化 */
:deep(.xterm-viewport::-webkit-scrollbar) {
  width: 8px;
}

:deep(.xterm-viewport::-webkit-scrollbar-track) {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

:deep(.xterm-viewport::-webkit-scrollbar-thumb) {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 4px;
}

:deep(.xterm-viewport::-webkit-scrollbar-thumb:hover) {
  background: rgba(255, 255, 255, 0.5);
}

/* 全屏模式下的样式调整 */
.modern-terminal :deep(.el-dialog.is-fullscreen) {
  .terminal-container {
    height: calc(100vh - 200px);
  }
}
</style>
