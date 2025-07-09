<template>
  <a-modal v-model:open="visible" title="SFTP 文件浏览器" :footer="null" @cancel="handleCancel" width="80vw">
    <div class="sftp-container">
      <div class="flex items-center gap-4 mb-4">
        <a-button @click="goUp" :disabled="currentPath === '/'">
          <template #icon><ArrowUpOutlined /></template>
          返回上级
        </a-button>
        <a-input :value="currentPath" readonly />
        <a-upload
          name="file"
          :action="uploadUrl"
          :headers="uploadHeaders"
          :data="{ path: currentPath, host_id: hostId }"
          :show-upload-list="false"
          @change="handleUploadChange"
        >
          <a-button type="primary">
            <template #icon><UploadOutlined /></template>
            上传文件
          </a-button>
        </a-upload>
      </div>

      <a-table
        :columns="columns"
        :data-source="files"
        :loading="loading"
        :row-key="(record) => record.name"
        @row-dblclick="handleRowDoubleClick"
        size="small"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <div class="flex items-center gap-2">
              <FolderOpenOutlined v-if="record.is_dir" />
              <FileOutlined v-else />
              <span>{{ record.name }}</span>
            </div>
          </template>
          <template v-if="column.key === 'size'">
            {{ record.is_dir ? '-' : (record.size / 1024).toFixed(2) + ' KB' }}
          </template>
          <template v-if="column.key === 'mod_time'">
            {{ new Date(record.mod_time * 1000).toLocaleString() }}
          </template>
          <template v-if="column.key === 'actions'">
            <a-button
              type="link"
              size="small"
              :href="getDownloadUrl(record)"
              :disabled="record.is_dir"
            >
              下载
            </a-button>
            <a-popconfirm
              title="确定要删除吗?"
              ok-text="确定"
              cancel-text="取消"
              @confirm="handleDelete(record)"
            >
              <a-button type="link" size="small" danger>删除</a-button>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useUserStore } from '@/store/modules/user';
import { listFiles, getSftpDownloadUrl, deleteSftpFile } from '@/api/system/host';
import {
  ArrowUpOutlined,
  UploadOutlined,
  FolderOpenOutlined,
  FileOutlined,
} from '@ant-design/icons-vue';
import { message } from 'ant-design-vue';
import type { UploadChangeParam } from 'ant-design-vue';

interface FileInfo {
  name: string;
  size: number;
  is_dir: boolean;
  mod_time: number;
}

const visible = ref(false);
const hostId = ref<number | null>(null);
const userStore = useUserStore();
const loading = ref(false);
const currentPath = ref('/');
const files = ref<FileInfo[]>([]);

const uploadUrl = computed(() => '/api/v1/cmdb/sftp/upload');
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${userStore.token}`,
}));

const open = (id: number) => {
  hostId.value = id;
  visible.value = true;
  fetchFiles('/');
};

const fetchFiles = async (path: string) => {
  if (!hostId.value) return;
  loading.value = true;
  try {
    const response = await listFiles(hostId.value, path);
    files.value = response.data.sort((a, b) => {
      if (a.is_dir !== b.is_dir) return a.is_dir ? -1 : 1;
      return a.name.localeCompare(b.name);
    });
    currentPath.value = path;
  } catch (err) {
    message.error('加载文件列表失败');
  } finally {
    loading.value = false;
  }
};

const goUp = () => {
  if (currentPath.value === '/') return;
  const parentPath = currentPath.value.substring(0, currentPath.value.lastIndexOf('/')) || '/';
  fetchFiles(parentPath);
};

const handleRowDoubleClick = (record: FileInfo) => {
  if (record.is_dir) {
    const newPath = [currentPath.value, record.name].join('/').replace('//', '/');
    fetchFiles(newPath);
  }
};

const handleUploadChange = (info: UploadChangeParam) => {
  if (info.file.status === 'done') {
    message.success(`${info.file.name} 上传成功`);
    fetchFiles(currentPath.value);
  } else if (info.file.status === 'error') {
    message.error(`${info.file.name} 上传失败.`);
  }
};

const getDownloadUrl = (record: FileInfo) => {
  if (!hostId.value || record.is_dir) return '#';
  const filePath = [currentPath.value, record.name].join('/').replace('//', '/');
  return getSftpDownloadUrl(hostId.value, filePath);
};

const handleDelete = async (record: FileInfo) => {
  if (!hostId.value) return;
  const filePath = [currentPath.value, record.name].join('/').replace('//', '/');
  try {
    await deleteSftpFile(hostId.value, filePath);
    message.success(`删除 ${record.name} 成功!`);
    fetchFiles(currentPath.value);
  } catch (error) {
    message.error(`删除失败: ${error}`);
  }
};

const handleCancel = () => {
  visible.value = false;
  currentPath.value = '/';
  files.value = [];
  hostId.value = null;
};

const columns = [
  { title: '名称', key: 'name', ellipsis: true },
  { title: '大小', key: 'size', width: 120 },
  { title: '修改时间', key: 'mod_time', width: 180 },
  { title: '操作', key: 'actions', width: 150, align: 'center' },
];

defineExpose({
  open,
});
</script>

<style scoped>
.sftp-container {
  height: 65vh;
  display: flex;
  flex-direction: column;
}
.ant-table-wrapper {
  flex-grow: 1;
  overflow-y: auto;
}
</style> 