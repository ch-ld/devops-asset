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
            延迟: {{ ping }}ms
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

    <!-- IP选择对话框 -->
    <IPSelectionDialog
      ref="ipSelectionDialog"
      :host="host"
      connection-type="ssh"
      @confirm="handleIPSelection"
      @cancel="handleIPSelectionCancel"
    />
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
import IPSelectionDialog from './IPSelectionDialog.vue'
import { analyzeHostIPs, autoSelectConnectionIP, type ConnectionParams } from '@/utils/ipSelection'

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

// IP选择相关
const ipSelectionDialog = ref()
const selectedConnectionParams = ref<ConnectionParams | null>(null)

// 终端相关
let terminal: Terminal | null = null
let fitAddon: FitAddon | null = null
let websocket: WebSocket | null = null
let pingInterval: number | null = null
let pingMeasurementInterval: number | null = null
let lastPingTime = ref<number | null>(null)
let connectionTimeout: number | null = null // 连接超时定时器

// 计算属性
const terminalTitle = computed(() => {
  if (!props.host) return 'SSH终端'
  return `SSH终端 - ${props.host.name}`
})

const primaryIP = computed(() => {
  if (!props.host) return 'Unknown'
  
  // 尝试解析公网IP
  try {
    const publicIPs = Array.isArray(props.host.public_ip) 
      ? props.host.public_ip 
      : JSON.parse(props.host.public_ip || '[]')
    // 过滤空字符串，找到第一个有效的公网IP
    const validPublicIP = publicIPs.find((ip: string) => ip && ip.trim() !== '')
    if (validPublicIP) {
      return validPublicIP
    }
  } catch (e) {
    // 如果解析失败，尝试直接使用（但也要检查是否为空）
    if (props.host.public_ip && props.host.public_ip.trim() !== '') {
      return props.host.public_ip
    }
  }
  
  // 尝试解析私网IP
  try {
    const privateIPs = Array.isArray(props.host.private_ip) 
      ? props.host.private_ip 
      : JSON.parse(props.host.private_ip || '[]')
    // 过滤空字符串，找到第一个有效的私网IP
    const validPrivateIP = privateIPs.find((ip: string) => ip && ip.trim() !== '')
    if (validPrivateIP) {
      return validPrivateIP
    }
  } catch (e) {
    // 如果解析失败，尝试直接使用（但也要检查是否为空）
    if (props.host.private_ip && props.host.private_ip.trim() !== '') {
      return props.host.private_ip
    }
  }
  
  return 'Unknown'
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

// 连接WebSocket
const connectWebSocket = async () => {
  console.log('🔍 [Frontend] Starting WebSocket connection...')
  console.log('🔍 [Frontend] Current connection status:', connectionStatus.value)
  console.log('🔍 [Frontend] Connecting flag:', connecting.value)
  
  if (connecting.value || connectionStatus.value === 'connected') {
    console.log('⚠️ [Frontend] Already connecting or connected, skipping...')
    return
  }

  try {
    connecting.value = true
    connectionStatus.value = 'connecting'
    console.log('🔍 [Frontend] Set connection status to connecting')

    // 如果没有选择连接参数，先选择IP
    if (!selectedConnectionParams.value) {
      console.log('🔍 [Frontend] No connection params, analyzing host IPs...')
      const analysis = analyzeHostIPs(props.host)
      console.log('🔍 [Frontend] IP analysis result:', analysis)
      
      if (analysis.needsSelection) {
        console.log('🔍 [Frontend] Multiple IPs found, opening selection dialog...')
        ipSelectionDialog.value?.open()
        return
      } else {
        selectedConnectionParams.value = autoSelectConnectionIP(props.host)
        console.log('🔍 [Frontend] Auto-selected connection params:', selectedConnectionParams.value)
      }
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const token = useUserStore().accessToken
    
    // 获取正确的API主机地址
    // 在开发环境中，前端可能运行在不同端口，需要连接到正确的后端端口
    let apiHost = window.location.host
    
    // 如果当前是开发环境（通过端口判断），使用正确的后端端口
    if (window.location.port === '3006' || window.location.port === '5173') {
      // 开发环境，连接到后端服务器
      apiHost = `${window.location.hostname}:8080`
    }
    
    const wsUrl = `${protocol}//${apiHost}/api/v1/cmdb/ws/ssh?host_id=${props.host.id}&token=${token}`
    console.log('🔍 [Frontend] WebSocket URL:', wsUrl)
    console.log('🔍 [Frontend] Host ID:', props.host.id)
    console.log('🔍 [Frontend] Token (first 20 chars):', token ? token.substring(0, 20) + '...' : 'null')

    websocket = new WebSocket(wsUrl)
    websocket.binaryType = 'arraybuffer' // 设置接收二进制数据
    console.log('🔍 [Frontend] WebSocket created, waiting for connection...')

    websocket.onopen = () => {
      console.log('✅ [Frontend] WebSocket connection opened')
      connecting.value = false
      connectionStatus.value = 'connecting' // WebSocket连接成功，但SSH还在连接中
      ElMessage.info('正在连接SSH服务器...')
      console.log('🔍 [Frontend] Set status to connecting, waiting for SSH connection...')
      
      // 设置连接超时（10秒后如果还没有收到数据就认为连接失败）
      connectionTimeout = window.setTimeout(() => {
        if (connectionStatus.value === 'connecting') {
          console.log('⏰ [Frontend] Connection timeout reached')
          connectionStatus.value = 'error'
          ElMessage.error('SSH连接超时，请检查主机配置和网络连通性')
          if (websocket) {
            websocket.close()
          }
        }
      }, 10000)
      
      // 初始化终端（但不显示为已连接）
      initTerminal()
    }

    websocket.onmessage = (event) => {
      console.log('📨 [Frontend] Received WebSocket message:', {
        type: typeof event.data,
        isArrayBuffer: event.data instanceof ArrayBuffer,
        size: event.data instanceof ArrayBuffer ? event.data.byteLength : event.data.length,
        preview: event.data instanceof ArrayBuffer 
          ? `[Binary: ${event.data.byteLength} bytes]` 
          : event.data.substring(0, 100) + (event.data.length > 100 ? '...' : '')
      })
      
      if (event.data instanceof ArrayBuffer) {
        // 处理二进制数据
        const uint8Array = new Uint8Array(event.data)
        if (terminal) {
          terminal.write(uint8Array)
        }
        
        // 如果收到终端数据，说明SSH连接已成功
        if (connectionStatus.value === 'connecting') {
          console.log('✅ [Frontend] Received binary data, SSH connection successful!')
          // 清除连接超时定时器
          if (connectionTimeout) {
            clearTimeout(connectionTimeout)
            connectionTimeout = null
          }
          connectionStatus.value = 'connected'
          connectionTime.value = new Date().toLocaleTimeString()
          ElMessage.success('SSH连接成功')
          startPing()
          startPingMeasurement()
        }
      } else if (typeof event.data === 'string') {
        console.log('📝 [Frontend] Received text message:', event.data)
        try {
          const message = JSON.parse(event.data)
          console.log('🔧 [Frontend] Parsed as JSON control message:', message)
          handleControlMessage(message)
        } catch (e) {
          console.log('📝 [Frontend] Not JSON, treating as terminal output or error message')
          // 不是JSON，可能是错误消息或终端输出
          if (terminal) {
            terminal.write(event.data)
          }
          
          // 检查是否是错误消息
          if (event.data.includes('连接失败') || 
              event.data.includes('authentication failed') || 
              event.data.includes('Connection refused') ||
              event.data.includes('Host is unreachable') ||
              event.data.includes('connection timeout') ||
              event.data.includes('用户名或密码错误') ||
              event.data.includes('目标主机拒绝连接') ||
              event.data.includes('连接超时') ||
              event.data.includes('网络不可达')) {
            console.log('❌ [Frontend] Detected error message:', event.data)
            // 清除连接超时定时器
            if (connectionTimeout) {
              clearTimeout(connectionTimeout)
              connectionTimeout = null
            }
            connectionStatus.value = 'error'
            ElMessage.error(event.data)
            return
          }
          
          // 如果收到非错误的文本数据，说明连接成功
          if (connectionStatus.value === 'connecting') {
            console.log('✅ [Frontend] Received non-error text, SSH connection successful!')
            // 清除连接超时定时器
            if (connectionTimeout) {
              clearTimeout(connectionTimeout)
              connectionTimeout = null
            }
            connectionStatus.value = 'connected'
            connectionTime.value = new Date().toLocaleTimeString()
            ElMessage.success('SSH连接成功')
            startPing()
            startPingMeasurement()
          }
        }
      }
    }

    websocket.onclose = (event) => {
      console.log('🔌 [Frontend] WebSocket connection closed:', {
        code: event.code,
        reason: event.reason,
        wasClean: event.wasClean
      })
      connecting.value = false
      connectionStatus.value = 'error'
      stopPing()
      stopPingMeasurement()
      
      // 清除连接超时定时器
      if (connectionTimeout) {
        clearTimeout(connectionTimeout)
        connectionTimeout = null
      }
      
      // 根据关闭代码给出不同的错误提示
      if (event.code === 1006) {
        console.log('❌ [Frontend] Connection failed - server not reachable')
        ElMessage.error('无法连接到后端服务器，请检查服务器是否启动')
      } else if (event.code !== 1000) {
        console.log('⚠️ [Frontend] Connection closed unexpectedly')
        ElMessage.warning('SSH连接已断开')
      }
    }

    websocket.onerror = (error) => {
      console.log('❌ [Frontend] WebSocket error:', error)
      connecting.value = false
      connectionStatus.value = 'error'
      stopPing()
      stopPingMeasurement()
      
      // 清除连接超时定时器
      if (connectionTimeout) {
        clearTimeout(connectionTimeout)
        connectionTimeout = null
      }
      
      ElMessage.error('无法连接到后端服务器，请检查网络连接和服务器状态')
      console.error('WebSocket error:', error)
    }

  } catch (error: any) {
    console.log('❌ [Frontend] Exception in connectWebSocket:', error)
    connecting.value = false
    connectionStatus.value = 'error'
    ElMessage.error('连接失败: ' + error.message)
  }
}

// 处理控制消息
const handleControlMessage = (message: any) => {
  switch (message.type) {
    case 'pong':
      // 更新延迟
      const currentTime = Date.now()
      if (lastPingTime.value !== null) {
        ping.value = currentTime - lastPingTime.value
      }
      break
    default:
      console.log('Unknown control message:', message)
  }
}

// 初始化终端
const initTerminal = () => {
  if (!xtermContainer.value) return

  // 销毁现有终端
  if (terminal) {
    terminal.dispose()
  }

  // 创建新终端实例
  terminal = new Terminal({
    fontSize: 14,
    fontFamily: 'Consolas, Monaco, "Courier New", monospace',
    allowProposedApi: true, // 允许使用proposed API
    theme: {
      background: '#1e1e1e',
      foreground: '#ffffff',
      cursor: '#ffffff',
      cursorAccent: '#000000',
      selectionBackground: '#3399ff33',
      black: '#000000',
      red: '#cd3131',
      green: '#0dbc79',
      yellow: '#e5e510',
      blue: '#2472c8',
      magenta: '#bc3fbc',
      cyan: '#11a8cd',
      white: '#e5e5e5',
      brightBlack: '#666666',
      brightRed: '#f14c4c',
      brightGreen: '#23d18b',
      brightYellow: '#f5f543',
      brightBlue: '#3b8eea',
      brightMagenta: '#d670d6',
      brightCyan: '#29b8db',
      brightWhite: '#ffffff'
    },
    allowTransparency: true,
    cursorBlink: true,
    cursorStyle: 'block',
    scrollback: 10000,
    tabStopWidth: 4
  })

  // 添加插件
  fitAddon = new FitAddon()
  terminal.loadAddon(fitAddon)
  terminal.loadAddon(new WebLinksAddon())
  terminal.loadAddon(new SearchAddon())
  terminal.loadAddon(new Unicode11Addon())
  terminal.unicode.activeVersion = '11'

  // 挂载到DOM
  terminal.open(xtermContainer.value)
  
  // 适配终端大小
  nextTick(() => {
    if (fitAddon && terminal) {
      fitAddon.fit()
      terminalSize.value = {
        cols: terminal.cols,
        rows: terminal.rows
      }
      // 发送终端大小到后端
      sendTerminalResize()
    }
  })

  // 监听终端输入
  terminal.onData((data) => {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
      // 发送二进制数据
      const encoder = new TextEncoder()
      const uint8Array = encoder.encode(data)
      websocket.send(uint8Array.buffer)
    }
  })

  // 监听终端大小变化
  terminal.onResize((size) => {
    terminalSize.value = { cols: size.cols, rows: size.rows }
    sendTerminalResize()
  })

  // 监听窗口大小变化
  const resizeObserver = new ResizeObserver(() => {
    if (fitAddon && terminal) {
      fitAddon.fit()
      const newSize = { cols: terminal.cols, rows: terminal.rows }
      if (newSize.cols !== terminalSize.value.cols || newSize.rows !== terminalSize.value.rows) {
        terminalSize.value = newSize
        sendTerminalResize()
      }
    }
  })
  
  if (xtermContainer.value) {
    resizeObserver.observe(xtermContainer.value)
  }

  // 保存 observer 以便清理 (使用 any 类型避免 TypeScript 错误)
  ;(terminal as any)._resizeObserver = resizeObserver
}

// 发送终端大小调整消息
const sendTerminalResize = () => {
  if (websocket && websocket.readyState === WebSocket.OPEN) {
    const resizeMessage = {
      type: 'resize',
      data: {
        cols: terminalSize.value.cols,
        rows: terminalSize.value.rows
      }
    }
    websocket.send(JSON.stringify(resizeMessage))
  }
}

// 开始心跳检测
const startPing = () => {
  pingInterval = window.setInterval(() => {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
      const pingMessage = {
        type: 'ping',
        data: Date.now()
      }
      lastPingTime.value = Date.now()
      websocket.send(JSON.stringify(pingMessage))
    }
  }, 30000) // 30秒发送一次心跳
}

// 停止心跳检测
const stopPing = () => {
  if (pingInterval) {
    clearInterval(pingInterval)
    pingInterval = null
  }
}

// 开始延迟测量
const startPingMeasurement = () => {
  pingMeasurementInterval = window.setInterval(() => {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
      lastPingTime.value = Date.now()
      // 发送WebSocket ping帧
      try {
        // 创建一个ping帧（如果浏览器支持）
        websocket.send(JSON.stringify({ type: 'ping', data: Date.now() }))
      } catch (e) {
        console.warn('Ping measurement failed:', e)
      }
    }
  }, 5000) // 5秒测量一次延迟
}

// 停止延迟测量
const stopPingMeasurement = () => {
  if (pingMeasurementInterval) {
    clearInterval(pingMeasurementInterval)
    pingMeasurementInterval = null
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
          const windowSize = {width: terminal.cols, height: terminal.rows}
          const blob = new Blob([JSON.stringify(windowSize)], {type : 'application/json'})
          websocket.send(blob)
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
  stopPingMeasurement()

  if (terminal) {
    terminal.dispose()
    terminal = null
  }

  dialogVisible.value = false
}

// IP选择处理方法
const handleIPSelection = (ipType: 'public' | 'private', ip: string) => {
  selectedConnectionParams.value = {
    ip,
    ipType,
    preferredIPType: ipType
  }

  // 选择IP后立即连接
  connectWebSocket()
}

const handleIPSelectionCancel = () => {
  // 用户取消IP选择，不进行连接
  connecting.value = false
  connectionStatus.value = 'disconnected'
}

// 监听对话框显示状态
watch(() => props.visible, (newVal) => {
  if (newVal && props.host) {
    nextTick(() => {
      initTerminal()
    })
  }
})

// 监听主机变化，重置连接参数
watch(() => props.host, (newHost, oldHost) => {
  // 当主机发生变化时，重置连接参数
  if (newHost && oldHost && newHost.id !== oldHost.id) {
    selectedConnectionParams.value = null
    // 如果终端正在连接或已连接，需要断开
    if (connectionStatus.value === 'connecting' || connectionStatus.value === 'connected') {
      handleClose()
    }
  }
}, { deep: true })

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
