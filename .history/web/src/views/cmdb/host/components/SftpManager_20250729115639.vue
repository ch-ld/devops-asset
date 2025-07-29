<template>
  <el-dialog
    title="文件管理"
    v-model="dialogVisible"
    width="1200px"
    :before-close="handleClose"
    destroy-on-close
    class="sftp-dialog"
  >
    <div class="sftp-manager">
      <div class="sftp-header">
        <div class="path-navigation">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item
              v-for="(segment, index) in pathSegments"
              :key="index"
              @click="navigateToPath(index)"
              class="breadcrumb-item"
            >
              {{ segment || 'root' }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="sftp-actions">
          <el-button size="small" @click="refreshFiles" :loading="loading">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
          <el-button size="small" @click="showUploadDialog = true">
            <el-icon><Upload /></el-icon>
            上传
          </el-button>
          <el-button size="small" @click="showCreateDirDialog = true">
            <el-icon><FolderAdd /></el-icon>
            新建文件夹
          </el-button>
        </div>
      </div>

    <div class="file-list" v-loading="loading">
      <el-table :data="fileList" @row-dblclick="handleRowDoubleClick" height="400">
        <el-table-column width="40">
          <template #default="{ row }">
            <el-icon v-if="row.isDir" color="#409EFF">
              <Folder />
            </el-icon>
            <el-icon v-else color="#67C23A">
              <Document />
            </el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="文件名" min-width="200" />
        <el-table-column prop="size" label="大小" width="100">
          <template #default="{ row }">
            {{ row.isDir ? '-' : formatFileSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column prop="mode" label="权限" width="100" />
        <el-table-column prop="modTime" label="修改时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.modTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button 
              v-if="!row.isDir" 
              size="small" 
              type="primary" 
              @click="downloadFile(row)"
            >
              下载
            </el-button>
            <el-button
              size="small"
              @click="handleShowRenameDialog(row)"
            >
              重命名
            </el-button>
            <el-button 
              size="small" 
              type="danger" 
              @click="deleteFile(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
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
  Refresh, Upload, FolderAdd, Folder, Document, UploadFilled 
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

// 响应式数据
const loading = ref(false)
const uploading = ref(false)
const currentPath = ref('/')
const fileList = ref<FileInfo[]>([])
const showUploadDialog = ref(false)
const showCreateDirDialog = ref(false)
const showRenameDialog = ref(false)
const uploadFileList = ref<any[]>([])
const uploadRef = ref()

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
  if (newVal) {
    loadFiles()
  }
})

// 组件挂载
onMounted(() => {
  if (props.visible) {
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
</script>

<style scoped>
.sftp-manager {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.sftp-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border-bottom: 1px solid #e4e7ed;
  background: #f5f7fa;
}

.path-navigation {
  flex: 1;
}

.breadcrumb-item {
  cursor: pointer;
}

.breadcrumb-item:hover {
  color: #409eff;
}

.sftp-actions {
  display: flex;
  gap: 8px;
}

.file-list {
  flex: 1;
  padding: 12px;
}

.el-upload__text {
  margin-top: 12px;
}

.el-icon--upload {
  font-size: 67px;
  color: #c0c4cc;
  margin-bottom: 16px;
}

/* 表格样式优化 */
:deep(.el-table) {
  border: 1px solid #ebeef5;
}

:deep(.el-table th) {
  background-color: #fafafa;
}

:deep(.el-table tr:hover > td) {
  background-color: #f5f7fa;
}

/* 对话框样式 */
:deep(.el-dialog__header) {
  padding: 20px 20px 10px;
}

:deep(.el-dialog__body) {
  padding: 10px 20px;
}

:deep(.el-dialog__footer) {
  padding: 10px 20px 20px;
}
</style>
