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
          <el-button @click="toggleFullscreen">
            <el-icon>
              <full-screen v-if="!isFullscreen" />
              <close v-else />
            </el-icon>
            {{ isFullscreen ? '退出全屏' : '全屏' }}
          </el-button>
        </div>
      </div>

      <div class="terminal-content" v-if="host">
        <ssh-terminal :host="host" :visible="dialogVisible" />
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
  import { defineComponent, ref, computed } from 'vue'
  import type { Host } from '@/types/api/host'
  import { FullScreen, Close } from '@element-plus/icons-vue'

  export default defineComponent({
    name: 'TerminalWindow',
    components: {
      FullScreen,
      Close
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
      const isFullscreen = ref(false)

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
      }

      // 关闭终端窗口
      const handleClose = () => {
        emit('update:visible', false)
      }

      return {
        primaryIP,
        title,
        isFullscreen,
        toggleFullscreen,
        handleClose,
        dialogVisible
      }
    }
  })
</script>

<style lang="scss" scoped>
  .terminal-window {
    :deep(.el-dialog__body) {
      padding: 0;
    }

    :deep(.el-dialog) {
      height: calc(80vh);
      display: flex;
      flex-direction: column;
    }

    :deep(.el-dialog__header) {
      flex: 0 0 auto;
    }

    :deep(.el-dialog__body) {
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
    }

    .terminal-actions {
      display: flex;
      gap: 8px;
    }
  }

  .terminal-placeholder {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #f5f5f5;
    color: #666;
    font-size: 16px;
    padding: 20px;
  }
</style>
