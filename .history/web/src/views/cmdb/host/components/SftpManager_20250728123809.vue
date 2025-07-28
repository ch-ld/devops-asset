<template>
  <div class="sftp-manager">
    <el-dialog
      title="SFTP文件管理"
      v-model="dialogVisible"
      :width="isFullscreen ? '100%' : '1200px'"
      :fullscreen="isFullscreen"
      :before-close="handleClose"
      destroy-on-close
      class="sftp-dialog"
    >
      <!-- 工具栏 -->
      <div class="sftp-toolbar">
        <div class="toolbar-left">
          <div class="host-info">
            <div class="host-avatar">
              <el-icon class="host-icon"><FolderOpened /></el-icon>
            </div>
            <div class="host-details">
              <span class="host-name">{{ host?.name || 'Unknown' }}</span>
              <span class="host-ip">{{ primaryIP }}</span>
            </div>
          </div>
          <div class="path-breadcrumb">
            <div class="breadcrumb-wrapper">
              <el-icon class="path-icon"><Folder /></el-icon>
              <el-breadcrumb separator="/" class="custom-breadcrumb">
                <el-breadcrumb-item @click="navigateToPath('/')" class="breadcrumb-item">
                  <el-icon><HomeFilled /></el-icon>
                  根目录
                </el-breadcrumb-item>
                <el-breadcrumb-item
                  v-for="(segment, index) in pathSegments"
                  :key="index"
                  @click="navigateToPath(getPathUpTo(index))"
                  class="breadcrumb-item"
                >
                  {{ segment }}
                </el-breadcrumb-item>
              </el-breadcrumb>
            </div>
          </div>
        </div>
        
        <div class="toolbar-right">
          <el-button-group>
            <el-button size="small" @click="refreshFiles" :loading="loading">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button size="small" @click="showUploadDialog">
              <el-icon><Upload /></el-icon>
              上传
            </el-button>
            <el-button size="small" @click="showCreateFolderDialog">
              <el-icon><FolderAdd /></el-icon>
              新建文件夹
            </el-button>
            <el-button size="small" @click="toggleFullscreen">
              <el-icon>
                <FullScreen v-if="!isFullscreen" />
                <Close v-else />
              </el-icon>
              {{ isFullscreen ? '退出全屏' : '全屏' }}
            </el-button>
          </el-button-group>
        </div>
      </div>

      <!-- 文件列表 -->
      <div class="file-list-container">
        <el-table
          :data="fileList"
          v-loading="loading"
          height="500"
          @row-dblclick="handleRowDoubleClick"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column label="名称" min-width="200">
            <template #default="{ row }">
              <div class="file-item">
                <el-icon class="file-icon" :class="getFileIconClass(row)">
                  <component :is="getFileIcon(row)" />
                </el-icon>
                <span class="file-name">{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="大小" width="120">
            <template #default="{ row }">
              {{ row.isDir ? '-' : formatFileSize(row.size) }}
            </template>
          </el-table-column>
          <el-table-column label="修改时间" width="180">
            <template #default="{ row }">
              {{ formatTime(row.modTime) }}
            </template>
          </el-table-column>
          <el-table-column label="权限" width="100">
            <template #default="{ row }">
              {{ row.mode }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button-group>
                <el-button 
                  v-if="!row.isDir" 
                  size="small" 
                  type="primary" 
                  @click="downloadFile(row)"
                >
                  <el-icon><Download /></el-icon>
                  下载
                </el-button>
                <el-button 
                  size="small" 
                  type="warning" 
                  @click="showRenameDialog(row)"
                >
                  <el-icon><Edit /></el-icon>
                  重命名
                </el-button>
                <el-button 
                  size="small" 
                  type="danger" 
                  @click="deleteFile(row)"
                >
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </el-button-group>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 状态栏 -->
      <div class="sftp-statusbar">
        <div class="status-left">
          <span class="status-item">
            <el-icon><Files /></el-icon>
            {{ fileList.length }} 项
          </span>
          <span class="status-item" v-if="selectedFiles.length > 0">
            <el-icon><Select /></el-icon>
            已选择 {{ selectedFiles.length }} 项
          </span>
        </div>
        <div class="status-right">
          <span class="status-item">
            <el-icon><FolderOpened /></el-icon>
            {{ currentPath }}
          </span>
        </div>
      </div>
    </el-dialog>

    <!-- 上传对话框 -->
    <el-dialog
      title="上传文件"
      v-model="uploadDialogVisible"
      width="500px"
      destroy-on-close
    >
      <el-upload
        ref="uploadRef"
        :action="uploadUrl"
        :headers="uploadHeaders"
        :data="uploadData"
        :on-success="handleUploadSuccess"
        :on-error="handleUploadError"
        :before-upload="beforeUpload"
        multiple
        drag
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持多文件上传，单个文件大小不超过100MB
          </div>
        </template>
      </el-upload>
      
      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="startUpload">开始上传</el-button>
      </template>
    </el-dialog>

    <!-- 新建文件夹对话框 -->
    <el-dialog
      title="新建文件夹"
      v-model="createFolderDialogVisible"
      width="400px"
      destroy-on-close
    >
      <el-form :model="createFolderForm" label-width="80px">
        <el-form-item label="文件夹名">
          <el-input 
            v-model="createFolderForm.name" 
            placeholder="请输入文件夹名称"
            @keyup.enter="createFolder"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="createFolderDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="createFolder">创建</el-button>
      </template>
    </el-dialog>

    <!-- 重命名对话框 -->
    <el-dialog
      title="重命名"
      v-model="renameDialogVisible"
      width="400px"
      destroy-on-close
    >
      <el-form :model="renameForm" label-width="80px">
        <el-form-item label="新名称">
          <el-input 
            v-model="renameForm.newName" 
            placeholder="请输入新名称"
            @keyup.enter="renameFile"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="renameDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="renameFile">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  FolderOpened, Refresh, Upload, FolderAdd, FullScreen, Close,
  Download, Edit, Delete, Files, Select, UploadFilled,
  Document, Folder, Picture, VideoPlay, HomeFilled
} from '@element-plus/icons-vue'
import { useUserStore } from '@/store/modules/user'
import { sftpApi } from '@/api/cmdb/sftp'

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

const isFullscreen = ref(false)
const loading = ref(false)
const currentPath = ref('/')
const fileList = ref<any[]>([])
const selectedFiles = ref<any[]>([])

// 对话框状态
const uploadDialogVisible = ref(false)
const createFolderDialogVisible = ref(false)
const renameDialogVisible = ref(false)

// 表单数据
const createFolderForm = ref({ name: '' })
const renameForm = ref({ oldName: '', newName: '' })

// 上传相关
const uploadRef = ref()
const userStore = useUserStore()

const uploadUrl = computed(() => `/api/v1/cmdb/sftp/upload`)
const uploadHeaders = computed(() => ({
  'Authorization': `Bearer ${userStore.accessToken}`
}))
const uploadData = computed(() => ({
  host_id: props.host?.id,
  path: currentPath.value
}))

// 计算属性
const primaryIP = computed(() => {
  if (!props.host) return 'Unknown'
  
  if (Array.isArray(props.host.public_ip) && props.host.public_ip.length > 0) {
    return props.host.public_ip[0]
  }
  if (Array.isArray(props.host.private_ip) && props.host.private_ip.length > 0) {
    return props.host.private_ip[0]
  }
  
  return props.host.public_ip || props.host.private_ip || 'Unknown'
})

const pathSegments = computed(() => {
  return currentPath.value.split('/').filter(segment => segment !== '')
})

// 方法
const getPathUpTo = (index: number) => {
  const segments = pathSegments.value.slice(0, index + 1)
  return '/' + segments.join('/')
}

const navigateToPath = (path: string) => {
  currentPath.value = path
  loadFiles()
}

const loadFiles = async () => {
  if (!props.host?.id) return
  
  loading.value = true
  try {
    const response = await sftpApi.listFiles(props.host.id, currentPath.value)
    fileList.value = response.data || []
  } catch (error) {
    ElMessage.error('加载文件列表失败')
    console.error('Load files error:', error)
  } finally {
    loading.value = false
  }
}

const refreshFiles = () => {
  loadFiles()
}

const handleRowDoubleClick = (row: any) => {
  if (row.isDir) {
    const newPath = currentPath.value === '/' 
      ? `/${row.name}` 
      : `${currentPath.value}/${row.name}`
    navigateToPath(newPath)
  }
}

const handleSelectionChange = (selection: any[]) => {
  selectedFiles.value = selection
}

// 文件图标
const getFileIcon = (file: any) => {
  if (file.isDir) return Folder

  const ext = file.name.split('.').pop()?.toLowerCase()
  if (['jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg'].includes(ext)) return Picture
  if (['mp4', 'avi', 'mkv', 'mov', 'wmv'].includes(ext)) return VideoPlay
  if (['mp3', 'wav', 'flac', 'aac'].includes(ext)) return Document

  return Document
}

const getFileIconClass = (file: any) => {
  if (file.isDir) return 'folder-icon'

  const ext = file.name.split('.').pop()?.toLowerCase()
  if (['jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg'].includes(ext)) return 'image-icon'
  if (['mp4', 'avi', 'mkv', 'mov', 'wmv'].includes(ext)) return 'video-icon'
  if (['mp3', 'wav', 'flac', 'aac'].includes(ext)) return 'audio-icon'

  return 'file-icon'
}

// 格式化文件大小
const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'

  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化时间
const formatTime = (timestamp: string) => {
  return new Date(timestamp).toLocaleString()
}

// 文件操作
const downloadFile = async (file: any) => {
  if (!props.host?.id) return

  try {
    const filePath = currentPath.value === '/'
      ? `/${file.name}`
      : `${currentPath.value}/${file.name}`

    const response = await sftpApi.downloadFile(props.host.id, filePath)

    // 创建下载链接
    const blob = new Blob([response.data])
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = file.name
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    ElMessage.success('文件下载成功')
  } catch (error) {
    ElMessage.error('文件下载失败')
    console.error('Download file error:', error)
  }
}

const deleteFile = async (file: any) => {
  if (!props.host?.id) return

  try {
    await ElMessageBox.confirm(
      `确定要删除 "${file.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const filePath = currentPath.value === '/'
      ? `/${file.name}`
      : `${currentPath.value}/${file.name}`

    await sftpApi.deleteFile(props.host.id, filePath)
    ElMessage.success('删除成功')
    loadFiles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
      console.error('Delete file error:', error)
    }
  }
}

// 对话框操作
const showUploadDialog = () => {
  uploadDialogVisible.value = true
}

const showCreateFolderDialog = () => {
  createFolderForm.value.name = ''
  createFolderDialogVisible.value = true
}

const showRenameDialog = (file: any) => {
  renameForm.value.oldName = file.name
  renameForm.value.newName = file.name
  renameDialogVisible.value = true
}

const createFolder = async () => {
  if (!props.host?.id || !createFolderForm.value.name.trim()) {
    ElMessage.warning('请输入文件夹名称')
    return
  }

  try {
    const folderPath = currentPath.value === '/'
      ? `/${createFolderForm.value.name}`
      : `${currentPath.value}/${createFolderForm.value.name}`

    await sftpApi.createFolder(props.host.id, folderPath)
    ElMessage.success('文件夹创建成功')
    createFolderDialogVisible.value = false
    loadFiles()
  } catch (error) {
    ElMessage.error('文件夹创建失败')
    console.error('Create folder error:', error)
  }
}

const renameFile = async () => {
  if (!props.host?.id || !renameForm.value.newName.trim()) {
    ElMessage.warning('请输入新名称')
    return
  }

  try {
    const oldPath = currentPath.value === '/'
      ? `/${renameForm.value.oldName}`
      : `${currentPath.value}/${renameForm.value.oldName}`
    const newPath = currentPath.value === '/'
      ? `/${renameForm.value.newName}`
      : `${currentPath.value}/${renameForm.value.newName}`

    await sftpApi.renameFile(props.host.id, oldPath, newPath)
    ElMessage.success('重命名成功')
    renameDialogVisible.value = false
    loadFiles()
  } catch (error) {
    ElMessage.error('重命名失败')
    console.error('Rename file error:', error)
  }
}

// 上传相关
const beforeUpload = (file: File) => {
  const isLt100M = file.size / 1024 / 1024 < 100
  if (!isLt100M) {
    ElMessage.error('文件大小不能超过100MB!')
    return false
  }
  return true
}

const startUpload = () => {
  uploadRef.value?.submit()
}

const handleUploadSuccess = (response: any, file: any) => {
  ElMessage.success(`${file.name} 上传成功`)
  loadFiles()
}

const handleUploadError = (error: any, file: any) => {
  ElMessage.error(`${file.name} 上传失败`)
  console.error('Upload error:', error)
}

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
}

const handleClose = () => {
  dialogVisible.value = false
}

// 生命周期
onMounted(() => {
  if (props.visible && props.host) {
    loadFiles()
  }
})
</script>

<style scoped>
.sftp-manager {
  .sftp-dialog {
    :deep(.el-dialog__body) {
      padding: 0;
    }
  }

  .sftp-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    background: linear-gradient(135deg, #42a5f5 0%, #478ed1 100%);
    color: white;
    border-radius: 8px 8px 0 0;

    .toolbar-left {
      display: flex;
      flex-direction: column;
      gap: 8px;

      .host-info {
        display: flex;
        align-items: center;
        gap: 8px;
        font-weight: 600;

        .host-icon {
          font-size: 16px;
        }

        .host-name {
          font-size: 16px;
        }

        .host-ip {
          font-size: 14px;
          opacity: 0.8;
          background: rgba(255, 255, 255, 0.2);
          padding: 2px 8px;
          border-radius: 12px;
        }
      }

      .path-breadcrumb {
        :deep(.el-breadcrumb) {
          .el-breadcrumb__item {
            .el-breadcrumb__inner {
              color: rgba(255, 255, 255, 0.8);
              cursor: pointer;

              &:hover {
                color: white;
              }
            }
          }

          .el-breadcrumb__separator {
            color: rgba(255, 255, 255, 0.6);
          }
        }
      }
    }

    .toolbar-right {
      :deep(.el-button-group) {
        .el-button {
          background: rgba(255, 255, 255, 0.1);
          border-color: rgba(255, 255, 255, 0.2);
          color: white;

          &:hover {
            background: rgba(255, 255, 255, 0.2);
            border-color: rgba(255, 255, 255, 0.3);
          }
        }
      }
    }
  }

  .file-list-container {
    .file-item {
      display: flex;
      align-items: center;
      gap: 8px;

      .file-icon {
        font-size: 16px;

        &.folder-icon {
          color: #ffa726;
        }

        &.image-icon {
          color: #66bb6a;
        }

        &.video-icon {
          color: #ef5350;
        }

        &.audio-icon {
          color: #ab47bc;
        }

        &.file-icon {
          color: #78909c;
        }
      }

      .file-name {
        font-weight: 500;
      }
    }
  }

  .sftp-statusbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 16px;
    background: #f5f5f5;
    border-top: 1px solid #e0e0e0;
    font-size: 12px;
    color: #666;

    .status-left,
    .status-right {
      display: flex;
      gap: 16px;
    }

    .status-item {
      display: flex;
      align-items: center;
      gap: 4px;

      .el-icon {
        font-size: 12px;
      }
    }
  }
}

/* 全屏模式下的样式调整 */
.sftp-manager :deep(.el-dialog.is-fullscreen) {
  .file-list-container {
    .el-table {
      height: calc(100vh - 250px) !important;
    }
  }
}
</style>
