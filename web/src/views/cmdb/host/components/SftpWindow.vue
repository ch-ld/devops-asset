<template>
  <div class="sftp-window">
    <el-dialog
      :title="title"
      v-model="dialogVisible"
      :width="isFullscreen ? '100%' : '1200px'"
      :fullscreen="isFullscreen"
      :before-close="handleClose"
      destroy-on-close
    >
      <div class="sftp-toolbar">
        <div class="sftp-info">
          <span v-if="host">{{ host.name }} - SFTP文件管理</span>
        </div>
        <div class="sftp-actions">
          <el-button @click="toggleFullscreen">
            <el-icon>
              <full-screen v-if="!isFullscreen" />
              <close v-else />
            </el-icon>
            {{ isFullscreen ? '退出全屏' : '全屏' }}
          </el-button>
        </div>
      </div>

      <div class="sftp-content" v-if="host">
        <sftp-manager :host="host" :visible="dialogVisible" />
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
  import { defineComponent, ref, computed } from 'vue'
  import type { Host } from '@/types/api/host'
  import { FullScreen, Close } from '@element-plus/icons-vue'
  import SftpManager from './SftpManager.vue'

  export default defineComponent({
    name: 'SftpWindow',
    components: {
      FullScreen,
      Close,
      SftpManager
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

      const dialogVisible = computed({
        get: () => props.visible,
        set: (value) => emit('update:visible', value)
      })

      const title = computed(() => {
        if (props.host) {
          return `SFTP文件管理 - ${props.host.name}`
        }
        return 'SFTP文件管理'
      })

      const toggleFullscreen = () => {
        isFullscreen.value = !isFullscreen.value
      }

      const handleClose = () => {
        dialogVisible.value = false
      }

      return {
        dialogVisible,
        title,
        isFullscreen,
        toggleFullscreen,
        handleClose
      }
    }
  })
</script>

<style scoped>
  .sftp-window {
    height: 100%;
  }

  .sftp-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 0 16px 0;
    border-bottom: 1px solid #e4e7ed;
    margin-bottom: 16px;
  }

  .sftp-info {
    font-size: 16px;
    font-weight: 500;
    color: #303133;
  }

  .sftp-actions {
    display: flex;
    gap: 8px;
  }

  .sftp-content {
    height: 600px;
    display: flex;
    flex-direction: column;
  }

  /* 全屏模式下的样式调整 */
  :deep(.el-dialog.is-fullscreen) .sftp-content {
    height: calc(100vh - 200px);
  }

  /* 对话框样式优化 */
  :deep(.el-dialog__header) {
    padding: 20px 20px 0;
  }

  :deep(.el-dialog__body) {
    padding: 20px;
    height: calc(100% - 60px);
  }

  :deep(.el-dialog__footer) {
    padding: 0 20px 20px;
  }
</style>
