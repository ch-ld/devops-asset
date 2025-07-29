<template>
  <el-dialog
    v-model="dialogVisible"
    :before-close="handleClose"
    destroy-on-close
    class="modern-sftp-dialog"
    width="90%"
    top="5vh"
  >
    <template #header>
      <div class="dialog-header">
        <div class="header-left">
          <div class="header-icon">
            <el-icon size="24"><FolderOpened /></el-icon>
          </div>
          <div class="header-info">
            <h3 class="dialog-title">文件管理器</h3>
            <p class="host-info">{{ props.host?.name || '未知主机' }} ({{ props.host?.publicIP || props.host?.privateIP }})</p>
          </div>
        </div>
        <div class="header-actions">
          <el-button circle size="small" @click="refreshFiles" :loading="loading" class="action-btn">
            <el-icon><Refresh /></el-icon>
          </el-button>
          <el-button circle size="small" @click="showUploadDialog = true" type="primary" class="action-btn">
            <el-icon><Upload /></el-icon>
          </el-button>
          <el-button circle size="small" @click="showCreateDirDialog = true" type="success" class="action-btn">
            <el-icon><FolderAdd /></el-icon>
          </el-button>
        </div>
      </div>
    </template>

    <div class="modern-sftp-manager">
      <!-- 路径导航栏 -->
      <div class="path-bar">
        <div class="path-navigation">
          <el-icon class="path-icon"><Location /></el-icon>
          <div class="breadcrumb-container">
            <span class="path-segment root" @click="navigateToPath(-1)">
              <el-icon><HomeFilled /></el-icon>
              根目录
            </span>
            <template v-for="(segment, index) in pathSegments" :key="index">
              <el-icon class="separator"><ArrowRight /></el-icon>
              <span
                class="path-segment"
                @click="navigateToPath(index)"
                :class="{ active: index === pathSegments.length - 1 }"
              >
                {{ segment }}
              </span>
            </template>
          </div>
        </div>
        <div class="view-controls">
          <el-button-group>
            <el-button
              :type="viewMode === 'grid' ? 'primary' : ''"
              size="small"
              @click="viewMode = 'grid'"
            >
              <el-icon><Grid /></el-icon>
            </el-button>
            <el-button
              :type="viewMode === 'list' ? 'primary' : ''"
              size="small"
              @click="viewMode = 'list'"
            >
              <el-icon><List /></el-icon>
            </el-button>
          </el-button-group>
        </div>
      </div>

      <!-- 文件列表区域 -->
      <div class="file-content" v-loading="loading" element-loading-text="加载中...">
        <!-- 网格视图 -->
        <div v-if="viewMode === 'grid'" class="grid-view">
          <div
            v-for="file in fileList"
            :key="file.name"
            class="file-card"
            @dblclick="handleRowDoubleClick(file)"
            @contextmenu.prevent="showContextMenu($event, file)"
          >
            <div class="file-icon">
              <el-icon v-if="file.isDir" size="48" class="folder-icon">
                <Folder />
              </el-icon>
              <el-icon v-else size="48" :class="getFileIconClass(file.name)">
                <component :is="getFileIcon(file.name)" />
              </el-icon>
            </div>
            <div class="file-info">
              <div class="file-name" :title="file.name">{{ file.name }}</div>
              <div class="file-meta">
                <span v-if="!file.isDir" class="file-size">{{ formatFileSize(file.size) }}</span>
                <span class="file-time">{{ formatTime(file.modTime) }}</span>
              </div>
            </div>
            <div class="file-actions">
              <el-dropdown trigger="click" @command="handleFileAction">
                <el-button circle size="small" class="more-btn">
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item v-if="!file.isDir" :command="{action: 'download', file}">
                      <el-icon><Download /></el-icon> 下载
                    </el-dropdown-item>
                    <el-dropdown-item :command="{action: 'rename', file}">
                      <el-icon><Edit /></el-icon> 重命名
                    </el-dropdown-item>
                    <el-dropdown-item :command="{action: 'delete', file}" class="danger-item">
                      <el-icon><Delete /></el-icon> 删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </div>

        <!-- 列表视图 -->
        <div v-else class="list-view">
          <el-table
            :data="fileList"
            @row-dblclick="handleRowDoubleClick"
            class="modern-table"
            :header-cell-style="{ background: '#f8fafc', color: '#374151' }"
          >
            <el-table-column width="60" align="center">
              <template #default="{ row }">
                <el-icon v-if="row.isDir" size="24" class="folder-icon">
                  <Folder />
                </el-icon>
                <el-icon v-else size="24" :class="getFileIconClass(row.name)">
                  <component :is="getFileIcon(row.name)" />
                </el-icon>
              </template>
            </el-table-column>
            <el-table-column prop="name" label="文件名" min-width="200">
              <template #default="{ row }">
                <div class="file-name-cell">
                  <span class="name">{{ row.name }}</span>
                  <span v-if="row.isDir" class="type-badge">文件夹</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="size" label="大小" width="120" align="right">
              <template #default="{ row }">
                <span class="file-size">{{ row.isDir ? '-' : formatFileSize(row.size) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="mode" label="权限" width="100" align="center">
              <template #default="{ row }">
                <el-tag size="small" type="info">{{ row.mode }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="modTime" label="修改时间" width="180">
              <template #default="{ row }">
                <span class="file-time">{{ formatTime(row.modTime) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" align="center">
              <template #default="{ row }">
                <div class="table-actions">
                  <el-button
                    v-if="!row.isDir"
                    size="small"
                    type="primary"
                    @click="downloadFile(row)"
                    class="action-btn"
                  >
                    <el-icon><Download /></el-icon>
                  </el-button>
                  <el-button
                    size="small"
                    @click="handleShowRenameDialog(row)"
                    class="action-btn"
                  >
                    <el-icon><Edit /></el-icon>
                  </el-button>
                  <el-button
                    size="small"
                    type="danger"
                    @click="deleteFile(row)"
                    class="action-btn"
                  >
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 空状态 -->
        <div v-if="!loading && fileList.length === 0" class="empty-state">
          <el-icon size="64" class="empty-icon"><FolderOpened /></el-icon>
          <p class="empty-text">此文件夹为空</p>
          <el-button type="primary" @click="showUploadDialog = true">
            <el-icon><Upload /></el-icon>
            上传文件
          </el-button>
        </div>
      </div>
    </div>

    <!-- 上传对话框 -->
    <el-dialog v-model="showUploadDialog" title="上传文件" width="500px">
      <el-upload
        ref="uploadRef"
        :auto-upload="false"
        :on-change="handleFileChange"
        :file-list="uploadFileList"
        drag
        multiple
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
      </el-upload>
      <template #footer>
        <el-button @click="showUploadDialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpload" :loading="uploading">
          上传
        </el-button>
      </template>
    </el-dialog>

    <!-- 创建文件夹对话框 -->
    <el-dialog v-model="showCreateDirDialog" title="新建文件夹" width="400px">
      <el-form :model="createDirForm" label-width="80px">
        <el-form-item label="文件夹名">
          <el-input v-model="createDirForm.name" placeholder="请输入文件夹名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDirDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreateDir">确定</el-button>
      </template>
    </el-dialog>

    <!-- 重命名对话框 -->
    <el-dialog v-model="showRenameDialog" title="重命名" width="400px">
      <el-form :model="renameForm" label-width="80px">
        <el-form-item label="新名称">
          <el-input v-model="renameForm.newName" placeholder="请输入新名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRenameDialog = false">取消</el-button>
        <el-button type="primary" @click="handleRename">确定</el-button>
      </template>
    </el-dialog>
  </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Refresh, Upload, FolderAdd, Folder, Document, UploadFilled,
  FolderOpened, Location, HomeFilled, ArrowRight, Grid, List,
  MoreFilled, Download, Edit, Delete, Picture, VideoPlay,
  DocumentCopy, Files, Cpu, Setting
} from '@element-plus/icons-vue'
import { sftpApi } from '@/api/cmdb/sftp'

// 定义文件信息接口
interface FileInfo {
  name: string
  size: number
  mode: string
  modTime: string
  isDir: boolean
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

// 响应式数据
const dialogVisible = ref(false)
const loading = ref(false)
const uploading = ref(false)
const currentPath = ref('/')
const fileList = ref<FileInfo[]>([])
const showUploadDialog = ref(false)
const showCreateDirDialog = ref(false)
const showRenameDialog = ref(false)
const uploadFileList = ref<any[]>([])
const uploadRef = ref()
const viewMode = ref<'grid' | 'list'>('grid')

// 表单数据
const createDirForm = ref({
  name: ''
})

const renameForm = ref({
  oldName: '',
  newName: ''
})

// 计算属性
const pathSegments = computed(() => {
  return currentPath.value.split('/').filter(segment => segment !== '')
})

// 监听visible变化
watch(() => props.visible, (newVal) => {
  dialogVisible.value = newVal
  if (newVal) {
    loadFiles()
  }
})

// 监听dialogVisible变化
watch(dialogVisible, (newVal) => {
  emit('update:visible', newVal)
})

// 关闭弹窗
const handleClose = () => {
  dialogVisible.value = false
}

// 组件挂载
onMounted(() => {
  if (props.visible) {
    dialogVisible.value = true
    loadFiles()
  }
})

// 加载文件列表
const loadFiles = async () => {
  if (!props.host?.id) return

  loading.value = true
  try {
    const response = await sftpApi.listFiles(props.host.id.toString(), currentPath.value)
    
    if (response.code === 0 || response.code === 200) {
      fileList.value = response.data || []
    } else {
      ElMessage.error(response.message || '加载文件列表失败')
      fileList.value = []
    }
  } catch (error: any) {
    console.error('Load files error:', error)
    ElMessage.error('加载文件列表失败：' + (error.message || '网络错误'))
    fileList.value = []
  } finally {
    loading.value = false
  }
}

// 刷新文件列表
const refreshFiles = () => {
  loadFiles()
}

// 导航到指定路径
const navigateToPath = (index: number) => {
  if (index === -1) {
    // 点击root，回到根目录
    currentPath.value = '/'
  } else {
    const segments = pathSegments.value.slice(0, index + 1)
    currentPath.value = '/' + segments.join('/')
  }
  loadFiles()
}

// 处理行双击事件
const handleRowDoubleClick = (row: FileInfo) => {
  if (row.isDir) {
    // 进入目录
    currentPath.value = currentPath.value.endsWith('/') 
      ? currentPath.value + row.name
      : currentPath.value + '/' + row.name
    loadFiles()
  } else {
    // 下载文件
    downloadFile(row)
  }
}

// 下载文件
const downloadFile = async (file: FileInfo) => {
  if (!props.host?.id) return

  try {
    const filePath = currentPath.value.endsWith('/') 
      ? currentPath.value + file.name
      : currentPath.value + '/' + file.name
      
    const blob = await sftpApi.downloadFile(props.host.id.toString(), filePath)
    
    // 创建下载链接
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = file.name
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('文件下载成功')
  } catch (error: any) {
    console.error('Download file error:', error)
    ElMessage.error('文件下载失败：' + (error.message || '网络错误'))
  }
}

// 文件选择处理
const handleFileChange = (file: any) => {
  uploadFileList.value.push(file)
}

// 处理文件上传
const handleUpload = async () => {
  if (!props.host?.id || uploadFileList.value.length === 0) {
    ElMessage.warning('请选择要上传的文件')
    return
  }

  uploading.value = true
  try {
    for (const fileWrapper of uploadFileList.value) {
      const file = fileWrapper.raw || fileWrapper
      await sftpApi.uploadFile(props.host.id.toString(), currentPath.value, file)
    }
    
    ElMessage.success('文件上传成功')
    showUploadDialog.value = false
    uploadFileList.value = []
    uploadRef.value?.clearFiles()
    loadFiles() // 刷新文件列表
  } catch (error: any) {
    console.error('Upload file error:', error)
    ElMessage.error('文件上传失败：' + (error.message || '网络错误'))
  } finally {
    uploading.value = false
  }
}

// 创建目录
const handleCreateDir = async () => {
  if (!props.host?.id || !createDirForm.value.name.trim()) {
    ElMessage.warning('请输入文件夹名称')
    return
  }

  try {
    const dirPath = currentPath.value.endsWith('/') 
      ? currentPath.value + createDirForm.value.name
      : currentPath.value + '/' + createDirForm.value.name

    const response = await sftpApi.createDirectory(props.host.id.toString(), dirPath)
    
    if (response.code === 0 || response.code === 200) {
      ElMessage.success('文件夹创建成功')
      showCreateDirDialog.value = false
      createDirForm.value.name = ''
      loadFiles() // 刷新文件列表
    } else {
      ElMessage.error(response.message || '文件夹创建失败')
    }
  } catch (error: any) {
    console.error('Create directory error:', error)
    ElMessage.error('文件夹创建失败：' + (error.message || '网络错误'))
  }
}

// 显示重命名对话框
const handleShowRenameDialog = (file: FileInfo) => {
  renameForm.value.oldName = file.name
  renameForm.value.newName = file.name
  showRenameDialog.value = true
}

// 重命名文件
const handleRename = async () => {
  if (!props.host?.id || !renameForm.value.newName.trim()) {
    ElMessage.warning('请输入新名称')
    return
  }

  if (renameForm.value.oldName === renameForm.value.newName) {
    ElMessage.warning('新名称与原名称相同')
    return
  }

  try {
    const oldPath = currentPath.value.endsWith('/') 
      ? currentPath.value + renameForm.value.oldName
      : currentPath.value + '/' + renameForm.value.oldName
      
    const newPath = currentPath.value.endsWith('/') 
      ? currentPath.value + renameForm.value.newName
      : currentPath.value + '/' + renameForm.value.newName

    const response = await sftpApi.renameFile(props.host.id.toString(), oldPath, newPath)
    
    if (response.code === 0 || response.code === 200) {
      ElMessage.success('重命名成功')
      showRenameDialog.value = false
      renameForm.value.oldName = ''
      renameForm.value.newName = ''
      loadFiles() // 刷新文件列表
    } else {
      ElMessage.error(response.message || '重命名失败')
    }
  } catch (error: any) {
    console.error('Rename file error:', error)
    ElMessage.error('重命名失败：' + (error.message || '网络错误'))
  }
}

// 删除文件
const deleteFile = async (file: FileInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 "${file.name}" 吗？${file.isDir ? '删除文件夹将会删除其中的所有内容。' : ''}`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const filePath = currentPath.value.endsWith('/') 
      ? currentPath.value + file.name
      : currentPath.value + '/' + file.name

    const response = await sftpApi.deleteFile(props.host.id.toString(), filePath)
    
    if (response.code === 0 || response.code === 200) {
      ElMessage.success('删除成功')
      loadFiles() // 刷新文件列表
    } else {
      ElMessage.error(response.message || '删除失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Delete file error:', error)
      ElMessage.error('删除失败：' + (error.message || '网络错误'))
    }
  }
}

// 格式化文件大小
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化时间
const formatTime = (timeStr: string): string => {
  try {
    const date = new Date(timeStr)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    return timeStr
  }
}

// 返回上级目录
const goToParentDirectory = () => {
  if (currentPath.value === '/') return

  const segments = currentPath.value.split('/').filter(s => s !== '')
  segments.pop()
  currentPath.value = segments.length > 0 ? '/' + segments.join('/') : '/'
  loadFiles()
}

// 获取文件图标
const getFileIcon = (fileName: string) => {
  const ext = fileName.split('.').pop()?.toLowerCase()

  switch (ext) {
    case 'jpg':
    case 'jpeg':
    case 'png':
    case 'gif':
    case 'bmp':
    case 'svg':
      return Picture
    case 'mp4':
    case 'avi':
    case 'mov':
    case 'wmv':
    case 'flv':
      return VideoPlay
    case 'txt':
    case 'md':
    case 'log':
      return DocumentCopy
    case 'zip':
    case 'rar':
    case '7z':
    case 'tar':
    case 'gz':
      return Files
    case 'exe':
    case 'msi':
    case 'deb':
    case 'rpm':
      return Cpu
    case 'conf':
    case 'config':
    case 'ini':
    case 'yaml':
    case 'yml':
    case 'json':
      return Setting
    default:
      return Document
  }
}

// 获取文件图标样式类
const getFileIconClass = (fileName: string) => {
  const ext = fileName.split('.').pop()?.toLowerCase()

  switch (ext) {
    case 'jpg':
    case 'jpeg':
    case 'png':
    case 'gif':
    case 'bmp':
    case 'svg':
      return 'image-icon'
    case 'mp4':
    case 'avi':
    case 'mov':
    case 'wmv':
    case 'flv':
      return 'video-icon'
    case 'txt':
    case 'md':
    case 'log':
      return 'text-icon'
    case 'zip':
    case 'rar':
    case '7z':
    case 'tar':
    case 'gz':
      return 'archive-icon'
    case 'exe':
    case 'msi':
    case 'deb':
    case 'rpm':
      return 'executable-icon'
    case 'conf':
    case 'config':
    case 'ini':
    case 'yaml':
    case 'yml':
    case 'json':
      return 'config-icon'
    default:
      return 'document-icon'
  }
}

// 处理文件操作
const handleFileAction = (command: any) => {
  const { action, file } = command

  switch (action) {
    case 'download':
      downloadFile(file)
      break
    case 'rename':
      handleShowRenameDialog(file)
      break
    case 'delete':
      deleteFile(file)
      break
  }
}

// 显示右键菜单
const showContextMenu = (event: MouseEvent, file: FileInfo) => {
  // 这里可以实现右键菜单功能
  console.log('Right click on file:', file.name)
}
</script>

<style scoped>
/* 现代化对话框样式 */
:deep(.modern-sftp-dialog) {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.modern-sftp-dialog .el-dialog__header) {
  padding: 0;
  margin: 0;
  border-bottom: 1px solid #e5e7eb;
}

:deep(.modern-sftp-dialog .el-dialog__body) {
  padding: 0;
  height: 70vh;
  overflow: hidden;
}

/* 对话框头部 */
.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-icon {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-info h3 {
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
  gap: 8px;
}

.action-btn {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.5);
}

/* 主容器 */
.modern-sftp-manager {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #f8fafc;
}

/* 路径导航栏 */
.path-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: white;
  border-bottom: 1px solid #e5e7eb;
}

.path-navigation {
  display: flex;
  align-items: center;
  gap: 8px;
}

.path-icon {
  color: #6b7280;
}

.breadcrumb-container {
  display: flex;
  align-items: center;
  gap: 4px;
}

.path-segment {
  padding: 6px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
  color: #6b7280;
  display: flex;
  align-items: center;
  gap: 4px;
}

.path-segment:hover {
  background: #f3f4f6;
  color: #374151;
}

.path-segment.root {
  color: #3b82f6;
  font-weight: 500;
}

.path-segment.active {
  background: #dbeafe;
  color: #1d4ed8;
  font-weight: 500;
}

.separator {
  color: #d1d5db;
  font-size: 12px;
}

.view-controls {
  display: flex;
  gap: 8px;
}

/* 文件内容区域 */
.file-content {
  flex: 1;
  overflow: auto;
  padding: 24px;
}

/* 网格视图 */
.grid-view {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
}

.file-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
  position: relative;
  border: 2px solid transparent;
}

.file-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
  border-color: #e5e7eb;
}

.file-icon {
  text-align: center;
  margin-bottom: 12px;
}

.folder-icon {
  color: #3b82f6;
}

.image-icon {
  color: #10b981;
}

.video-icon {
  color: #f59e0b;
}

.text-icon {
  color: #6b7280;
}

.archive-icon {
  color: #8b5cf6;
}

.executable-icon {
  color: #ef4444;
}

.config-icon {
  color: #06b6d4;
}

.document-icon {
  color: #6b7280;
}

.file-info {
  text-align: center;
}

.file-name {
  font-weight: 500;
  color: #374151;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-meta {
  font-size: 12px;
  color: #6b7280;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.file-actions {
  position: absolute;
  top: 12px;
  right: 12px;
  opacity: 0;
  transition: opacity 0.2s;
}

.file-card:hover .file-actions {
  opacity: 1;
}

.more-btn {
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid #e5e7eb;
}

/* 列表视图 */
.list-view {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

:deep(.modern-table) {
  border: none;
}

:deep(.modern-table .el-table__header) {
  background: #f8fafc;
}

:deep(.modern-table .el-table__row:hover > .el-table__cell) {
  background: #f8fafc;
}

.file-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.name {
  font-weight: 500;
  color: #374151;
}

.type-badge {
  background: #dbeafe;
  color: #1d4ed8;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
}

.file-size {
  color: #6b7280;
  font-size: 13px;
}

.file-time {
  color: #6b7280;
  font-size: 13px;
}

.table-actions {
  display: flex;
  gap: 4px;
  justify-content: center;
}

.table-actions .action-btn {
  background: transparent;
  border: 1px solid #e5e7eb;
  color: #6b7280;
  width: 32px;
  height: 32px;
  padding: 0;
}

.table-actions .action-btn:hover {
  background: #f3f4f6;
  border-color: #d1d5db;
  color: #374151;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #6b7280;
}

.empty-icon {
  color: #d1d5db;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  margin-bottom: 24px;
}

/* 下拉菜单样式 */
:deep(.el-dropdown-menu__item.danger-item) {
  color: #ef4444;
}

:deep(.el-dropdown-menu__item.danger-item:hover) {
  background: #fef2f2;
  color: #dc2626;
}

/* 上传对话框样式 */
.el-upload__text {
  margin-top: 12px;
  color: #6b7280;
}

.el-icon--upload {
  font-size: 67px;
  color: #d1d5db;
  margin-bottom: 16px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .grid-view {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 16px;
  }

  .file-card {
    padding: 16px;
  }

  .path-bar {
    padding: 12px 16px;
  }

  .file-content {
    padding: 16px;
  }
}
</style>
