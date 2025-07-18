<template>
  <div>
    <a-modal
      v-model:visible="visible"
      :title="title"
      width="90%"
      :footer="null"
      :destroyOnClose="true"
      class="sftp-modal"
      :maskClosable="false"
      @cancel="handleClose"
    >
      <div class="sftp-layout">
        <!-- 左侧目录树 -->
        <div class="sftp-sidebar">
          <div class="sftp-header">
            <a-select v-model:value="currentHost" style="width: 100%" @change="handleHostChange">
              <a-select-option v-for="host in hostList" :key="host.id" :value="host.id">
                {{ host.name }} ({{ getHostIP(host) }})
              </a-select-option>
            </a-select>
          </div>

          <div class="sftp-navigation">
            <a-tree
              v-if="directoryTree.length > 0"
              :tree-data="directoryTree"
              :selectedKeys="selectedKeys"
              :expandedKeys="expandedKeys"
              :auto-expand-parent="true"
              @select="handleSelect"
              @expand="handleExpand"
            >
              <template #title="{ title, key }">
                <span>
                  <folder-outlined v-if="key.endsWith('/')" />
                  <file-outlined v-else />
                  {{ title }}
                </span>
              </template>
            </a-tree>
            <a-empty v-else description="No directories" />
          </div>
        </div>

        <!-- 右侧文件列表 -->
        <div class="sftp-content">
          <div class="sftp-toolbar">
            <!-- 路径导航 -->
            <div class="sftp-breadcrumb">
              <a-breadcrumb>
                <a-breadcrumb-item>
                  <a @click="navigateToPath('/')">Root</a>
                </a-breadcrumb-item>
                <template v-for="(part, index) in pathParts" :key="index">
                  <a-breadcrumb-item v-if="part">
                    <a @click="navigateToPath(getPathUntil(index))">{{ part }}</a>
                  </a-breadcrumb-item>
                </template>
              </a-breadcrumb>
            </div>

            <!-- 操作按钮 -->
            <div class="sftp-actions">
              <a-button-group>
                <a-button @click="refreshCurrentPath">
                  <template #icon><ReloadOutlined /></template>
                  刷新
                </a-button>
                <a-button @click="showCreateDirModal">
                  <template #icon><FolderAddOutlined /></template>
                  新建文件夹
                </a-button>
                <a-button @click="showUploadModal">
                  <template #icon><UploadOutlined /></template>
                  上传文件
                </a-button>
              </a-button-group>
            </div>
          </div>

          <!-- 文件列表 -->
          <div class="sftp-file-list">
            <a-table
              :columns="columns"
              :data-source="fileList"
              :pagination="false"
              :loading="loading"
              size="small"
              :scroll="{ y: 500 }"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'name'">
                  <div class="file-name-cell">
                    <folder-filled v-if="record.is_dir" class="file-icon folder" />
                    <file-filled v-else class="file-icon file" />
                    <a v-if="record.is_dir" @click="navigateTo(record)">{{ record.name }}</a>
                    <span v-else>{{ record.name }}</span>
                  </div>
                </template>
                <template v-else-if="column.key === 'size'">
                  <span>{{ formatFileSize(record.size) }}</span>
                </template>
                <template v-else-if="column.key === 'actions'">
                  <div class="file-actions">
                    <a-tooltip v-if="record.is_dir" title="打开">
                      <a @click="navigateTo(record)"><folder-open-outlined /></a>
                    </a-tooltip>
                    <a-tooltip v-else title="下载">
                      <a @click="downloadFile(record)"><download-outlined /></a>
                    </a-tooltip>
                    <a-tooltip title="重命名">
                      <a @click="showRenameModal(record)"><edit-outlined /></a>
                    </a-tooltip>
                    <a-tooltip title="删除">
                      <a @click="confirmDelete(record)"><delete-outlined /></a>
                    </a-tooltip>
                  </div>
                </template>
              </template>
            </a-table>
          </div>
        </div>
      </div>
    </a-modal>

    <!-- 创建目录模态框 -->
    <a-modal
      v-model:visible="createDirModalVisible"
      title="新建文件夹"
      @ok="createDirectory"
      @cancel="createDirModalVisible = false"
    >
      <a-form :model="dirForm" layout="vertical">
        <a-form-item
          label="文件夹名称"
          name="name"
          :rules="[{ required: true, message: '请输入文件夹名称' }]"
        >
          <a-input v-model:value="dirForm.name" placeholder="请输入文件夹名称" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 重命名模态框 -->
    <a-modal
      v-model:visible="renameModalVisible"
      title="重命名"
      @ok="renameFile"
      @cancel="renameModalVisible = false"
    >
      <a-form :model="renameForm" layout="vertical">
        <a-form-item
          label="新名称"
          name="newName"
          :rules="[{ required: true, message: '请输入新名称' }]"
        >
          <a-input v-model:value="renameForm.newName" placeholder="请输入新名称" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 上传文件模态框 -->
    <a-modal
      v-model:visible="uploadModalVisible"
      title="上传文件"
      @ok="uploadFile"
      @cancel="uploadModalVisible = false"
      :ok-button-props="{ disabled: !uploadForm.file }"
    >
      <a-form :model="uploadForm" layout="vertical">
        <a-form-item label="选择文件" name="file">
          <a-upload
            v-model:file-list="uploadFileList"
            :before-upload="beforeUpload"
            :multiple="false"
            :max-count="1"
          >
            <a-button>
              <template #icon><upload-outlined /></template>
              选择文件
            </a-button>
          </a-upload>
        </a-form-item>
        <a-form-item label="上传路径" name="path">
          <a-input v-model:value="uploadForm.path" disabled />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
  import { ref, computed, watch, onMounted } from 'vue';
  import { ElMessage, ElMessageBox } from 'element-plus';
  import { useHostStore } from '@/store/modules/host';
  import {
    FolderOutlined, FileOutlined, FolderFilled, FileFilled,
    ReloadOutlined, FolderAddOutlined, UploadOutlined,
    DownloadOutlined, EditOutlined, DeleteOutlined,
    FolderOpenOutlined
  } from '@ant-design/icons-vue';
  import type { TableColumnType } from 'ant-design-vue';

  interface FileItem {
    key: string;
    name: string;
    path: string;
    size: number;
    is_dir: boolean;
    mode: string;
    mod_time: string;
    permissions: string;
  }

  // 状态
  const visible = ref(false);
  const loading = ref(false);
  const currentHost = ref(null);
  const currentPath = ref('/home');
  const fileList = ref([]);
  const hostStore = useHostStore();
  const hostList = computed(() => hostStore.hostList);

  // 目录树状态
  const directoryTree = ref([]);
  const selectedKeys = ref([]);
  const expandedKeys = ref([]);

  // 模态框状态
  const createDirModalVisible = ref(false);
  const renameModalVisible = ref(false);
  const uploadModalVisible = ref(false);

  // 表单状态
  const dirForm = ref({ name: '' });
  const renameForm = ref({ name: '', newName: '', path: '', isDir: false });
  const uploadForm = ref({ file: null, path: currentPath.value });
  const uploadFileList = ref([]);

  // 表格列定义
  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      sorter: (a, b) => a.name.localeCompare(b.name),
      width: '40%',
      ellipsis: true
    },
    {
      title: '大小',
      dataIndex: 'size',
      key: 'size',
      sorter: (a, b) => a.size - b.size,
      width: '15%'
    },
    {
      title: '修改时间',
      dataIndex: 'mod_time',
      key: 'mod_time',
      sorter: (a, b) => new Date(a.mod_time).getTime() - new Date(b.mod_time).getTime(),
      width: '25%'
    },
    {
      title: '权限',
      dataIndex: 'permissions',
      key: 'permissions',
      width: '10%'
    },
    {
      title: '操作',
      key: 'actions',
      width: '10%'
    }
  ];

  // 计算属性
  const title = computed(() => {
    const host = hostList.value.find(h => h.id === currentHost.value);
    return host ? `文件管理 - ${host.name}` : '文件管理';
  });

  const pathParts = computed(() => {
    return currentPath.value.split('/').filter(Boolean);
  });

  // 方法 - 基础功能
  const open = (host) => {
    visible.value = true;

    if (host) {
      currentHost.value = host.id;
      listFiles();
    } else if (hostList.value.length > 0 && !currentHost.value) {
      currentHost.value = hostList.value[0].id;
      listFiles();
    }
  };

  const handleClose = () => {
    visible.value = false;
  };

  const getHostIP = (host) => {
    try {
      const ips = typeof host.public_ip === 'string' ? JSON.parse(host.public_ip) : host.public_ip;
      return ips && ips.length > 0 ? ips[0] : '无IP';
    } catch (error) {
      return '无IP';
    }
  };

  const handleHostChange = (hostId) => {
    currentHost.value = hostId;
    currentPath.value = '/home';
    listFiles();
  };

  const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  // 方法 - 文件操作
  const listFiles = async () => {
    if (!currentHost.value) return;

    loading.value = true;
    try {
      const response = await fetch(`/api/v1/cmdb/sftp/list?host_id=${currentHost.value}&path=${encodeURIComponent(currentPath.value)}`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      });

      const data = await response.json();
      if (data.code === 0) {
        fileList.value = data.data.map((item) => ({
          ...item,
          key: item.path
        }));
        updateDirectoryTree();
      } else {
        ElMessage.error(data.message || '获取文件列表失败');
      }
    } catch (error) {
      console.error('Failed to list files:', error);
      ElMessage.error('获取文件列表失败');
    } finally {
      loading.value = false;
    }
  };

  const navigateTo = (file) => {
    if (file.is_dir) {
      currentPath.value = file.path;
      listFiles();
    }
  };

  const navigateToPath = (path) => {
    currentPath.value = path || '/';
    listFiles();
  };

  const getPathUntil = (index) => {
    return '/' + pathParts.value.slice(0, index + 1).join('/');
  };

  const refreshCurrentPath = () => {
    listFiles();
  };

  const updateDirectoryTree = () => {
    // 基于文件列表构建目录树
    const tree = [];
    const dirs = fileList.value.filter(item => item.is_dir);

    // 添加根目录
    tree.push({
      key: '/',
      title: 'Root',
      children: []
    });

    // 添加当前目录下的目录
    const currentPathNode = tree[0]; // 根节点
    dirs.forEach(dir => {
      currentPathNode.children.push({
        key: dir.path,
        title: dir.name,
        isLeaf: false
      });
    });

    directoryTree.value = tree;
    selectedKeys.value = [currentPath.value];
    expandedKeys.value = ['/'];
  };

  // 文件操作 - 目录创建
  const showCreateDirModal = () => {
    dirForm.value.name = '';
    createDirModalVisible.value = true;
  };

  const createDirectory = async () => {
    if (!dirForm.value.name || !currentHost.value) {
      return;
    }

    try {
      const formData = new FormData();
      formData.append('host_id', String(currentHost.value));
      formData.append('path', `${currentPath.value}/${dirForm.value.name}`);

      const response = await fetch('/api/v1/cmdb/sftp/mkdir', {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        },
        body: formData
      });

      const data = await response.json();
      if (data.code === 0) {
        ElMessage.success('创建目录成功');
        createDirModalVisible.value = false;
        listFiles();
      } else {
        ElMessage.error(data.message || '创建目录失败');
      }
    } catch (error) {
      console.error('Failed to create directory:', error);
      ElMessage.error('创建目录失败');
    }
  };

  // 文件操作 - 重命名
  const showRenameModal = (file) => {
    renameForm.value = {
      name: file.name,
      newName: file.name,
      path: file.path,
      isDir: file.is_dir
    };
    renameModalVisible.value = true;
  };

  const renameFile = async () => {
    if (!renameForm.value.newName || !currentHost.value) {
      return;
    }

    try {
      const formData = new FormData();
      formData.append('host_id', String(currentHost.value));
      formData.append('old_path', renameForm.value.path);

      // 构建新路径
      const pathParts = renameForm.value.path.split('/');
      pathParts[pathParts.length - 1] = renameForm.value.newName;
      const newPath = pathParts.join('/');
      formData.append('new_path', newPath);

      const response = await fetch('/api/v1/cmdb/sftp/rename', {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        },
        body: formData
      });

      const data = await response.json();
      if (data.code === 0) {
        ElMessage.success('重命名成功');
        renameModalVisible.value = false;
        listFiles();
      } else {
        ElMessage.error(data.message || '重命名失败');
      }
    } catch (error) {
      console.error('Failed to rename:', error);
      ElMessage.error('重命名失败');
    }
  };

  // 文件操作 - 删除
  const confirmDelete = (file) => {
    const fileType = file.is_dir ? '目录' : '文件';
    ElMessageBox.confirm(
      `确定要删除${fileType} "${file.name}" 吗？`,
      `确认删除${fileType}`,
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    ).then(() => {
      deleteFile(file);
    }).catch(() => {
      // 取消操作
    });
  };

  const deleteFile = async (file) => {
    if (!currentHost.value) return;

    try {
      const response = await fetch(`/api/v1/cmdb/sftp/delete?host_id=${currentHost.value}&path=${encodeURIComponent(file.path)}`, {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      });

      const data = await response.json();
      if (data.code === 0) {
        ElMessage.success('删除成功');
        listFiles();
      } else {
        ElMessage.error(data.message || '删除失败');
      }
    } catch (error) {
      console.error('Failed to delete:', error);
      ElMessage.error('删除失败');
    }
  };

  // 文件操作 - 下载
  const downloadFile = async (file) => {
    if (!currentHost.value) return;

    try {
      const a = document.createElement('a');
      a.href = `/api/v1/cmdb/sftp/download?host_id=${currentHost.value}&path=${encodeURIComponent(file.path)}`;
      a.download = file.name;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
    } catch (error) {
      console.error('Failed to download:', error);
      ElMessage.error('下载失败');
    }
  };

  // 文件操作 - 上传
  const showUploadModal = () => {
    uploadForm.value = { file: null, path: currentPath.value };
    uploadFileList.value = [];
    uploadModalVisible.value = true;
  };

  const beforeUpload = (file) => {
    uploadForm.value.file = file;
    return false; // 阻止自动上传
  };

  const uploadFile = async () => {
    if (!uploadForm.value.file || !currentHost.value) {
      return;
    }

    try {
      const formData = new FormData();
      formData.append('host_id', String(currentHost.value));
      formData.append('path', uploadForm.value.path);
      formData.append('file', uploadForm.value.file);

      const response = await fetch('/api/v1/cmdb/sftp/upload', {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        },
        body: formData
      });

      const data = await response.json();
      if (data.code === 0) {
        ElMessage.success('上传成功');
        uploadModalVisible.value = false;
        listFiles();
      } else {
        ElMessage.error(data.message || '上传失败');
      }
    } catch (error) {
      console.error('Failed to upload:', error);
      ElMessage.error('上传失败');
    }
  };

  // 监听树操作
  const handleSelect = (keys) => {
    if (keys.length > 0) {
      selectedKeys.value = keys;
      navigateToPath(keys[0]);
    }
  };

  const handleExpand = (keys) => {
    expandedKeys.value = keys;
  };

  // 监听主机列表变化
  watch(() => hostStore.hostList, () => {
    if (hostList.value.length > 0 && !currentHost.value) {
      currentHost.value = hostList.value[0].id;
    }
  });

  // 生命周期钩子
  onMounted(() => {
    hostStore.fetchHosts();
  });

  defineExpose({
    open
  });
</script>

<style lang="scss">
  .sftp-modal {
    .ant-modal-body {
      padding: 0;
      height: calc(80vh - 55px);
    }

    .sftp-layout {
      display: flex;
      height: 100%;

      .sftp-sidebar {
        width: 250px;
        border-right: 1px solid #e8e8e8;
        display: flex;
        flex-direction: column;

        .sftp-header {
          padding: 12px;
          border-bottom: 1px solid #e8e8e8;
        }

        .sftp-navigation {
          flex: 1;
          padding: 12px;
          overflow: auto;
        }
      }

      .sftp-content {
        flex: 1;
        display: flex;
        flex-direction: column;

        .sftp-toolbar {
          display: flex;
          justify-content: space-between;
          padding: 12px;
          border-bottom: 1px solid #e8e8e8;

          .sftp-breadcrumb {
            flex: 1;
            display: flex;
            align-items: center;
          }

          .sftp-actions {
            display: flex;
            gap: 8px;
          }
        }

        .sftp-file-list {
          flex: 1;
          padding: 0 12px;
          overflow: auto;

          .file-name-cell {
            display: flex;
            align-items: center;

            .file-icon {
              margin-right: 8px;

              &.folder {
                color: #1890ff;
              }

              &.file {
                color: #8c8c8c;
              }
            }
          }

          .file-actions {
            display: flex;
            gap: 8px;

            a {
              font-size: 16px;
            }
          }
        }
      }
    }
  }
</style>
