<template>
  <div class="terminal-window">
    <el-dialog
      :title="title"
      v-model="dialogVisible"
      :width="isFullscreen ? '100%' : '800px'"
      :fullscreen="isFullscreen"
      :before-close="handleClose"
      destroy-on-close
    >
      <div class="terminal-toolbar">
        <div class="terminal-info">
          <span v-if="host">{{ host.name }} - {{ primaryIP }}</span>
        </div>
        <div class="terminal-actions">
          <el-tooltip content="字体大小">
            <el-button-group>
              <el-button @click="increaseFontSize">
                <el-icon><font-size-outlined /></el-icon>
                +
              </el-button>
              <el-button @click="decreaseFontSize">
                <el-icon><font-size-outlined /></el-icon>
                -
              </el-button>
            </el-button-group>
          </el-tooltip>

          <el-dropdown @command="changeTheme">
            <el-button>
              <el-icon><bg-colors-outlined /></el-icon>
              主题
              <el-icon class="el-icon--right"><down-outlined /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="dark">深色主题</el-dropdown-item>
                <el-dropdown-item command="light">浅色主题</el-dropdown-item>
                <el-dropdown-item command="solarized">Solarized</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

          <el-dropdown @command="changeCursor">
            <el-button>
              <el-icon><border-outlined /></el-icon>
              光标样式
              <el-icon class="el-icon--right"><down-outlined /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="block">块状</el-dropdown-item>
                <el-dropdown-item command="bar">竖线</el-dropdown-item>
                <el-dropdown-item command="underline">下划线</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

          <el-button @click="toggleFullscreen">
            <el-icon>
              <fullscreen-outlined v-if="!isFullscreen" />
              <fullscreen-exit-outlined v-else />
            </el-icon>
            {{ isFullscreen ? '退出全屏' : '全屏' }}
          </el-button>
        </div>
      </div>

      <div class="terminal-container" ref="terminalContainer"></div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
  import { defineComponent, ref, onMounted, watch, computed, nextTick } from 'vue'
  import { ElMessage } from 'element-plus'
  import { Terminal } from '@xterm/xterm'
  import { FitAddon } from '@xterm/addon-fit'
  import { WebLinksAddon } from '@xterm/addon-web-links'
  import { SearchAddon } from '@xterm/addon-search'
  import { Unicode11Addon } from '@xterm/addon-unicode11'
  import '@xterm/xterm/css/xterm.css'
  import type { Host } from '@/types/api/host'
  import {
    Edit as FontSizeOutlined,
    Brush as BgColorsOutlined,
    Grid as BorderOutlined,
    FullScreen as FullscreenOutlined,
    Close as FullscreenExitOutlined,
    ArrowDown as DownOutlined
  } from '@element-plus/icons-vue'
  import { useUserStore } from '@/store/modules/user'

  export default defineComponent({
    name: 'TerminalWindow',
    components: {
      FontSizeOutlined,
      BgColorsOutlined,
      BorderOutlined,
      FullscreenOutlined,
      FullscreenExitOutlined,
      DownOutlined
    },
    props: {
      visible: {
        type: Boolean,
        default: false
      },
      host: {
        type: Object as () => Host | null,
        default: null
      }
    },
    emits: ['update:visible'],
    setup(props, { emit }) {
      const terminalContainer = ref<HTMLElement | null>(null)
      const terminal = ref<Terminal | null>(null)
      const fitAddon = ref<FitAddon | null>(null)
      const ws = ref<WebSocket | null>(null)
      const isFullscreen = ref(false)
      const terminalPrefs = ref({
        font_size: 14,
        theme: 'dark',
        fullscreen: false,
        cursor_style: 'block'
      })

      // 计算主机IP
      const primaryIP = computed(() => {
        if (!props.host) return ''
        try {
          const ipList =
            typeof props.host.public_ip === 'string'
              ? JSON.parse(props.host.public_ip)
              : props.host.public_ip || []
          return ipList.length > 0 ? ipList[0] : ''
        } catch {
          return ''
        }
      })

      // 计算窗口标题
      const title = computed(() => {
        if (!props.host) return 'Terminal'
        return `Terminal: ${props.host.name} (${primaryIP.value})`
      })

      // 处理对话框可见性
      const dialogVisible = computed({
        get: () => props.visible,
        set: (value) => emit('update:visible', value)
      })

      // 切换全屏模式
      const toggleFullscreen = () => {
        isFullscreen.value = !isFullscreen.value
        terminalPrefs.value.fullscreen = isFullscreen.value

        // 调整terminal大小以适应容器
        nextTick(() => {
          fitAddon.value?.fit()
        })

        // 保存偏好设置
        savePreferences()
      }

      // 改变主题
      const changeTheme = ({ key }: { key: string }) => {
        terminalPrefs.value.theme = key
        applyTheme()
        savePreferences()
      }

      // 应用主题
      const applyTheme = () => {
        if (!terminal.value) return

        // 清除现有主题类
        const terminalEl = terminalContainer.value?.querySelector('.xterm')
        if (terminalEl) {
          terminalEl.classList.remove('theme-dark', 'theme-light', 'theme-solarized')
          terminalEl.classList.add(`theme-${terminalPrefs.value.theme}`)
        }

        // 设置终端颜色
        switch (terminalPrefs.value.theme) {
          case 'dark':
            terminal.value.options.theme = {
              background: '#1e1e1e',
              foreground: '#f0f0f0'
            }
            break
          case 'light':
            terminal.value.options.theme = {
              background: '#ffffff',
              foreground: '#000000'
            }
            break
          case 'solarized':
            terminal.value.options.theme = {
              background: '#002b36',
              foreground: '#839496',
              selectionBackground: '#073642',
              black: '#073642',
              red: '#dc322f',
              green: '#859900',
              yellow: '#b58900',
              blue: '#268bd2',
              magenta: '#d33682',
              cyan: '#2aa198',
              white: '#eee8d5'
            }
            break
        }
      }

      // 改变光标样式
      const changeCursor = ({ key }: { key: string }) => {
        if (!terminal.value) return
        terminalPrefs.value.cursor_style = key

        // 设置光标样式
        terminal.value.options.cursorStyle = key as 'block' | 'underline' | 'bar'
        savePreferences()
      }

      // 增加字体大小
      const increaseFontSize = () => {
        if (terminalPrefs.value.font_size < 24) {
          terminalPrefs.value.font_size += 2
          updateFontSize()
          savePreferences()
        }
      }

      // 减小字体大小
      const decreaseFontSize = () => {
        if (terminalPrefs.value.font_size > 10) {
          terminalPrefs.value.font_size -= 2
          updateFontSize()
          savePreferences()
        }
      }

      // 更新字体大小
      const updateFontSize = () => {
        if (!terminal.value) return
        terminal.value.options.fontSize = terminalPrefs.value.font_size
        // 调整terminal大小以适应容器
        nextTick(() => {
          fitAddon.value?.fit()
        })
      }

      // 保存偏好设置
      const savePreferences = () => {
        localStorage.setItem('terminal_preferences', JSON.stringify(terminalPrefs.value))
      }

      // 加载偏好设置
      const loadPreferences = () => {
        try {
          const savedPrefs = localStorage.getItem('terminal_preferences')
          if (savedPrefs) {
            const parsedPrefs = JSON.parse(savedPrefs)
            Object.assign(terminalPrefs.value, parsedPrefs)
          }
        } catch (e) {
          console.error('Failed to load terminal preferences:', e)
        }
      }

      // 初始化终端
      const initTerminal = () => {
        if (!terminalContainer.value) return

        // 加载用户偏好设置
        loadPreferences()

        // 创建终端实例
        terminal.value = new Terminal({
          fontSize: terminalPrefs.value.font_size,
          fontFamily: 'Menlo, Monaco, "Courier New", monospace',
          cursorBlink: true,
          cursorStyle: terminalPrefs.value.cursor_style as 'block' | 'underline' | 'bar',
          scrollback: 1000,
          convertEol: true,
          allowTransparency: true
        })

        // 添加插件
        fitAddon.value = new FitAddon()
        const webLinksAddon = new WebLinksAddon()
        const searchAddon = new SearchAddon()
        const unicode11Addon = new Unicode11Addon()

        terminal.value.loadAddon(fitAddon.value)
        terminal.value.loadAddon(webLinksAddon)
        terminal.value.loadAddon(searchAddon)
        terminal.value.loadAddon(unicode11Addon)

        // 打开终端并应用主题
        terminal.value.open(terminalContainer.value)
        applyTheme()

        // 自适应容器大小
        setTimeout(() => {
          fitAddon.value?.fit()
        }, 100)

        // 处理窗口大小变化
        window.addEventListener('resize', () => {
          fitAddon.value?.fit()
        })

        // 如果已经可见且有主机信息，则自动连接
        if (props.visible && props.host) {
          connectWebSocket()
        }
      }

      // 连接WebSocket
      const connectWebSocket = () => {
        if (!props.host || !props.host.id) {
          ElMessage.error('无法连接：主机信息不完整')
          return
        }

        // 关闭现有连接
        if (ws.value && ws.value.readyState !== WebSocket.CLOSED) {
          ws.value.close()
        }

        // 获取认证信息
        const userStore = useUserStore()
        const token = userStore.accessToken

        // 创建WebSocket连接
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const wsUrl = `${protocol}//${window.location.host}/api/v1/cmdb/terminal/connect?host_id=${props.host.id}&token=${token}`

        try {
          ws.value = new WebSocket(wsUrl)

          // 连接成功
          ws.value.onopen = () => {
            ElMessage.success('终端连接成功')

            // 清空终端
            if (terminal.value) {
              terminal.value.clear()
              terminal.value.write('\x1b[1;32mConnected to SSH terminal\x1b[0m\r\n')
            }
          }

          // 接收消息
          ws.value.onmessage = (event) => {
            if (terminal.value) {
              terminal.value.write(event.data)
            }
          }

          // 连接关闭
          ws.value.onclose = () => {
            if (terminal.value) {
              terminal.value.write('\r\n\x1b[1;31mConnection closed\x1b[0m\r\n')
            }
          }

          // 连接错误
          ws.value.onerror = (error) => {
            console.error('WebSocket error:', error)
            ElMessage.error('终端连接发生错误')
            if (terminal.value) {
              terminal.value.write('\r\n\x1b[1;31mConnection error\x1b[0m\r\n')
            }
          }

          // 监听终端输入
          if (terminal.value) {
            terminal.value.onData((data) => {
              if (ws.value && ws.value.readyState === WebSocket.OPEN) {
                ws.value.send(data)
              }
            })
          }
        } catch (error) {
          console.error('Failed to connect WebSocket:', error)
          ElMessage.error('无法建立终端连接')
        }
      }

      // 关闭终端窗口
      const handleClose = () => {
        // 关闭WebSocket连接
        if (ws.value) {
          ws.value.close()
          ws.value = null
        }

        emit('update:visible', false)
      }

      // 监听可见性变化
      watch(
        () => props.visible,
        (newVal) => {
          if (newVal) {
            // 延迟初始化，等待DOM渲染完成
            nextTick(() => {
              if (!terminal.value) {
                initTerminal()
              } else {
                // 已有终端实例，重新连接WebSocket
                connectWebSocket()
              }
            })
          } else {
            // 关闭WebSocket连接
            if (ws.value) {
              ws.value.close()
            }
          }
        }
      )

      // 监听主机变化
      watch(
        () => props.host,
        (newHost) => {
          if (props.visible && newHost && terminal.value) {
            // 重新连接WebSocket
            connectWebSocket()
          }
        }
      )

      // 组件挂载时初始化
      onMounted(() => {
        if (props.visible) {
          initTerminal()
        }
      })

      return {
        terminalContainer,
        primaryIP,
        title,
        isFullscreen,
        terminalPrefs,
        toggleFullscreen,
        changeTheme,
        changeCursor,
        increaseFontSize,
        decreaseFontSize,
        handleClose,
        dialogVisible
      }
    }
  })
</script>

<style lang="scss" scoped>
  .terminal-modal {
    :deep(.ant-modal-body) {
      padding: 0;
    }

    :deep(.ant-modal-content) {
      height: calc(80vh);
      display: flex;
      flex-direction: column;
    }

    :deep(.ant-modal-header) {
      flex: 0 0 auto;
    }

    :deep(.ant-modal-body) {
      flex: 1 1 auto;
      display: flex;
      flex-direction: column;
      overflow: hidden;
    }
  }

  .terminal-toolbar {
    display: flex;
    justify-content: space-between;
    padding: 8px 16px;
    border-bottom: 1px solid #e8e8e8;
    background-color: #fafafa;
    flex: 0 0 auto;

    .terminal-info {
      display: flex;
      align-items: center;

      .host-name {
        font-weight: 500;
        margin-right: 16px;
      }

      .host-ip {
        color: #999;
      }
    }

    .terminal-actions {
      display: flex;
      gap: 8px;
    }
  }

  .terminal-container {
    flex: 1;
    overflow: hidden;
    background-color: #1e1e1e;
  }

  :deep(.xterm) {
    height: 100%;

    &.theme-dark {
      background-color: #1e1e1e;
      .xterm-viewport {
        background-color: #1e1e1e;
      }
    }

    &.theme-light {
      background-color: #ffffff;
      .xterm-viewport {
        background-color: #ffffff;
      }
    }

    &.theme-solarized {
      background-color: #002b36;
      .xterm-viewport {
        background-color: #002b36;
      }
    }
  }
</style>
