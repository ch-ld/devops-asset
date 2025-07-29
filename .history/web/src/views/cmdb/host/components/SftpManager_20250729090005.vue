<template>
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
const uploadFileList = ref([])
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
    if (response.code === 200) {
      fileList.value = response.data || []
    } else {
      ElMessage.error(response.message || '加载文件列表失败')
    }
  } catch (error) {
    ElMessage.error('加载文件列表失败')
    console.error('Load files error:', error)
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
  const segments = pathSegments.value.slice(0, index + 1)
  currentPath.value = '/' + segments.join('/')
  loadFiles()
}

// 处理行双击事件
const handleRowDoubleClick = (row: FileInfo) => {
  if (row.isDir) {
    currentPath.value = currentPath.value.endsWith('/') 
      ? currentPath.value + row.name 
      : currentPath.value + '/' + row.name
    loadFiles()
  }
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
const formatTime = (timeStr: string) => {
  return new Date(timeStr).toLocaleString()
}

// 处理文件选择
const handleFileChange = (file: any) => {
  // 文件选择处理逻辑
}

// 处理上传
const handleUpload = async () => {
  const files = uploadRef.value?.uploadFiles || []
  if (files.length === 0) {
    ElMessage.warning('请选择要上传的文件')
    return
  }

  uploading.value = true
  try {
    for (const fileItem of files) {
      await sftpApi.uploadFile(
        props.host.id.toString(), 
        currentPath.value, 
        fileItem.raw
      )
    }
    ElMessage.success('文件上传成功')
    showUploadDialog.value = false
    uploadFileList.value = []
    loadFiles()
  } catch (error) {
    ElMessage.error('文件上传失败')
    console.error('Upload error:', error)
  } finally {
    uploading.value = false
  }
}

// 处理创建文件夹
const handleCreateDir = async () => {
  if (!createDirForm.value.name.trim()) {
    ElMessage.warning('请输入文件夹名称')
    return
  }

  try {
    const dirPath = currentPath.value.endsWith('/') 
      ? currentPath.value + createDirForm.value.name 
      : currentPath.value + '/' + createDirForm.value.name

    await sftpApi.createDirectory(props.host.id.toString(), dirPath)
    ElMessage.success('文件夹创建成功')
    showCreateDirDialog.value = false
    createDirForm.value.name = ''
    loadFiles()
  } catch (error) {
    ElMessage.error('文件夹创建失败')
    console.error('Create directory error:', error)
  }
}

// 显示重命名对话框
const handleShowRenameDialog = (row: FileInfo) => {
  renameForm.value.oldName = row.name
  renameForm.value.newName = row.name
  showRenameDialog.value = true
}

// 处理重命名
const handleRename = async () => {
  if (!renameForm.value.newName.trim()) {
    ElMessage.warning('请输入新名称')
    return
  }

  try {
    const oldPath = currentPath.value.endsWith('/') 
      ? currentPath.value + renameForm.value.oldName 
      : currentPath.value + '/' + renameForm.value.oldName
    
    const newPath = currentPath.value.endsWith('/') 
      ? currentPath.value + renameForm.value.newName 
      : currentPath.value + '/' + renameForm.value.newName

    await sftpApi.renameFile(props.host.id.toString(), oldPath, newPath)
    ElMessage.success('重命名成功')
    showRenameDialog.value = false
    loadFiles()
  } catch (error) {
    ElMessage.error('重命名失败')
    console.error('Rename error:', error)
  }
}

// 下载文件
const downloadFile = async (row: FileInfo) => {
  try {
    const filePath = currentPath.value.endsWith('/') 
      ? currentPath.value + row.name 
      : currentPath.value + '/' + row.name

    const blob = await sftpApi.downloadFile(props.host.id.toString(), filePath)
    
    // 创建下载链接
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = row.name
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('文件下载成功')
  } catch (error) {
    ElMessage.error('文件下载失败')
    console.error('Download error:', error)
  }
}

// 删除文件
const deleteFile = async (row: FileInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 ${row.name} 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const filePath = currentPath.value.endsWith('/') 
      ? currentPath.value + row.name 
      : currentPath.value + '/' + row.name

    await sftpApi.deleteFile(props.host.id.toString(), filePath)
    ElMessage.success('删除成功')
    loadFiles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
      console.error('Delete error:', error)
    }
  }
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
