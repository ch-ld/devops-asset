<template>
  <a-modal
    title="SFTP文件管理"
    :open="visible"
    :width="900"
    :footer="null"
    @cancel="handleCancel"
  >
    <div class="sftp-container">
      <div class="sftp-toolbar">
        <a-breadcrumb>
          <a-breadcrumb-item>
            <a @click="navigateToPath('/')">根目录</a>
          </a-breadcrumb-item>
          <a-breadcrumb-item v-for="(part, index) in pathParts" :key="index">
            <a @click="navigateToPath(getPathUpTo(index))">{{ part }}</a>
          </a-breadcrumb-item>
        </a-breadcrumb>
        
        <a-space>
          <a-button @click="refreshFiles">
            <template #icon><ReloadOutlined /></template>
            刷新
          </a-button>
          <a-upload
            :show-upload-list="false"
            :before-upload="handleUpload"
            :multiple="true"
          >
            <a-button type="primary">
              <template #icon><UploadOutlined /></template>
              上传文件
            </a-button>
          </a-upload>
          <a-button @click="showCreateFolder = true">
            <template #icon><FolderAddOutlined /></template>
            新建文件夹
          </a-button>
        </a-space>
      </div>

      <a-table
        :columns="columns"
        :data-source="fileList"
        :loading="loading"
        :pagination="false"
        size="small"
        row-key="name"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <a-space>
              <component
                :is="record.is_dir ? FolderOutlined : FileOutlined"
                :style="{ color: record.is_dir ? '#1890ff' : '#666' }"
              />
              <a
                v-if="record.is_dir"
                @click="navigateToPath(currentPath + '/' + record.name)"
              >
                {{ record.name }}
              </a>
              <span v-else>{{ record.name }}</span>
            </a-space>
          </template>

          <template v-if="column.key === 'size'">
            <span v-if="!record.is_dir">{{ formatFileSize(record.size) }}</span>
            <span v-else>-</span>
          </template>

          <template v-if="column.key === 'mod_time'">
            {{ formatDate(record.mod_time) }}
          </template>

          <template v-if="column.key === 'action'">
            <a-space>
              <a-button
                v-if="!record.is_dir"
                type="link"
                size="small"
                @click="downloadFile(record)"
              >
                <template #icon><DownloadOutlined /></template>
                下载
              </a-button>
              <a-popconfirm
                title="确定要删除吗？"
                @confirm="deleteFile(record)"
              >
                <a-button type="link" size="small" danger>
                  <template #icon><DeleteOutlined /></template>
                  删除
                </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>

      <!-- 创建文件夹对话框 -->
      <a-modal
        v-model:open="showCreateFolder"
        title="新建文件夹"
        @ok="createFolder"
      >
        <a-input
          v-model:value="newFolderName"
          placeholder="请输入文件夹名称"
          @press-enter="createFolder"
        />
      </a-modal>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { message } from 'ant-design-vue'
import { 
  ReloadOutlined, 
  UploadOutlined, 
  FolderAddOutlined,
  FolderOutlined,
  FileOutlined,
  DownloadOutlined,
  DeleteOutlined
} from '@ant-design/icons-vue'
import * as hostApi from '@/api/system/host'

const visible = ref(false)
const loading = ref(false)
const hostId = ref<number>()
const currentPath = ref('/')
const fileList = ref<any[]>([])
const showCreateFolder = ref(false)
const newFolderName = ref('')

const pathParts = computed(() => {
  return currentPath.value.split('/').filter(part => part !== '')
})

const columns = [
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name',
    width: 300
  },
  {
    title: '大小',
    dataIndex: 'size',
    key: 'size',
    width: 120
  },
  {
    title: '修改时间',
    dataIndex: 'mod_time',
    key: 'mod_time',
    width: 180
  },
  {
    title: '操作',
    key: 'action',
    width: 120
  }
]

const formatFileSize = (size: number) => {
  if (size === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(size) / Math.log(k))
  return parseFloat((size / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (timestamp: number) => {
  return new Date(timestamp * 1000).toLocaleString()
}

const getPathUpTo = (index: number) => {
  return '/' + pathParts.value.slice(0, index + 1).join('/')
}

const navigateToPath = (path: string) => {
  currentPath.value = path
  loadFiles()
}

const loadFiles = async () => {
  if (!hostId.value) return
  
  loading.value = true
  try {
    const { data } = await hostApi.listFiles(hostId.value, currentPath.value)
    fileList.value = data || []
  } catch (error) {
    message.error('加载文件列表失败')
    console.error('Load files error:', error)
  } finally {
    loading.value = false
  }
}

const refreshFiles = () => {
  loadFiles()
}

const handleUpload = async (file: File) => {
  if (!hostId.value) return false
  
  loading.value = true
  try {
    await hostApi.uploadSftpFile(hostId.value, currentPath.value, file)
    message.success('上传成功')
    loadFiles()
  } catch (error) {
    message.error('上传失败')
    console.error('Upload error:', error)
  } finally {
    loading.value = false
  }
  
  return false // 阻止默认上传行为
}

const downloadFile = (record: any) => {
  if (!hostId.value) return
  
  const filePath = currentPath.value === '/' ? 
    `/${record.name}` : 
    `${currentPath.value}/${record.name}`
  
  const downloadUrl = hostApi.getSftpDownloadUrl(hostId.value, filePath)
  
  // 创建临时链接下载
  const link = document.createElement('a')
  link.href = downloadUrl
  link.download = record.name
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

const deleteFile = async (record: any) => {
  if (!hostId.value) return
  
  const filePath = currentPath.value === '/' ? 
    `/${record.name}` : 
    `${currentPath.value}/${record.name}`
  
  try {
    await hostApi.deleteSftpFile(hostId.value, filePath)
    message.success('删除成功')
    loadFiles()
  } catch (error) {
    message.error('删除失败')
    console.error('Delete error:', error)
  }
}

const createFolder = async () => {
  if (!newFolderName.value.trim()) {
    message.warning('请输入文件夹名称')
    return
  }
  
  // 这里需要实现创建文件夹的API
  message.info('创建文件夹功能待实现')
  showCreateFolder.value = false
  newFolderName.value = ''
}

const handleCancel = () => {
  visible.value = false
  hostId.value = undefined
  currentPath.value = '/'
  fileList.value = []
}

const open = (id: number) => {
  hostId.value = id
  visible.value = true
  currentPath.value = '/'
  loadFiles()
}

defineExpose({
  open
})
</script>

<style scoped>
.sftp-container {
  .sftp-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    padding: 12px;
    background-color: #fafafa;
    border-radius: 6px;
  }
}
</style> 