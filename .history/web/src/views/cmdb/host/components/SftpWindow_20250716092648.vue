<template>
  <div>
    <el-dialog
      v-model="visible"
      :title="title"
      width="90%"
      :destroy-on-close="true"
      class="sftp-modal"
      :close-on-click-modal="false"
      @close="handleClose"
    >
      <div class="sftp-layout">
        <!-- 左侧目录树 -->
        <div class="sftp-sidebar">
          <div class="sftp-header">
            <el-select v-model="currentHost" style="width: 100%" @change="handleHostChange">
              <el-option
                v-for="host in hostList"
                :key="host.id"
                :value="host.id"
                :label="host.name + ' (' + getHostIP(host) + ')'"
              >
              </el-option>
            </el-select>
          </div>

          <div class="sftp-navigation">
            <el-tree
              v-if="directoryTree.length > 0"
              :data="directoryTree"
              :current-node-key="selectedKeys[0]"
              :expand-on-click-node="false"
              :default-expanded-keys="expandedKeys"
              node-key="key"
              @node-click="handleSelect"
            >
              <template #default="{ node, data }">
                <span>
                  <el-icon v-if="data.key.endsWith('/')"><FolderIcon /></el-icon>
                  <el-icon v-else><FileIcon /></el-icon>
                  {{ node.label }}
                </span>
              </template>
            </el-tree>
            <el-empty v-else description="No directories" />
          </div>
        </div>

        <!-- 右侧文件列表 -->
        <div class="sftp-content">
          <div class="sftp-toolbar">
            <!-- 路径导航 -->
            <div class="sftp-breadcrumb">
              <el-breadcrumb>
                <el-breadcrumb-item>
                  <a @click="navigateToPath('/')">Root</a>
                </el-breadcrumb-item>
                <template v-for="(part, index) in pathParts" :key="index">
                  <el-breadcrumb-item v-if="part">
                    <a @click="navigateToPath(getPathUntil(index))">{{ part }}</a>
                  </el-breadcrumb-item>
                </template>
              </el-breadcrumb>
            </div>

            <!-- 操作按钮 -->
            <div class="sftp-actions">
              <el-button-group>
                <el-button @click="refreshCurrentPath">
                  <el-icon><RefreshIcon /></el-icon>
                  刷新
                </el-button>
                <el-button @click="showCreateDirModal">
                  <el-icon><PlusIcon /></el-icon>
                  新建文件夹
                </el-button>
                <el-button @click="showUploadModal">
                  <el-icon><UploadIcon /></el-icon>
                  上传文件
                </el-button>
              </el-button-group>
            </div>
          </div>

          <!-- 文件列表 -->
          <div class="sftp-file-list">
            <el-table :data="fileList" :height="500" v-loading="loading" size="small">
              <el-table-column label="名称" prop="name" min-width="200" sortable>
                <template #default="{ row }">
                  <div class="file-name-cell">
                    <el-icon v-if="row.is_dir" class="file-icon folder"><FolderIcon /></el-icon>
                    <el-icon v-else class="file-icon file"><FileIcon /></el-icon>
                    <a v-if="row.is_dir" @click="navigateTo(row)">{{ row.name }}</a>
                    <span v-else>{{ row.name }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="大小" prop="size" width="120" sortable>
                <template #default="{ row }">
                  <span>{{ formatFileSize(row.size) }}</span>
                </template>
              </el-table-column>
              <el-table-column label="修改时间" prop="mod_time" width="180" sortable />
              <el-table-column label="权限" prop="permissions" width="100" />
              <el-table-column label="操作" width="120" fixed="right">
                <template #default="{ row }">
                  <div class="file-actions">
                    <el-tooltip v-if="row.is_dir" content="打开">
                      <el-button type="primary" link @click="navigateTo(row)">
                        <el-icon><FolderOpenIcon /></el-icon>
                      </el-button>
                    </el-tooltip>
                    <el-tooltip v-else content="下载">
                      <el-button type="primary" link @click="downloadFile(row)">
                        <el-icon><DownloadIcon /></el-icon>
                      </el-button>
                    </el-tooltip>
                    <el-tooltip content="重命名">
                      <el-button type="primary" link @click="showRenameModal(row)">
                        <el-icon><EditIcon /></el-icon>
                      </el-button>
                    </el-tooltip>
                    <el-tooltip content="删除">
                      <el-button type="danger" link @click="confirmDelete(row)">
                        <el-icon><DeleteIcon /></el-icon>
                      </el-button>
                    </el-tooltip>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 创建目录模态框 -->
    <el-dialog v-model="createDirModalVisible" title="新建文件夹" width="400px">
      <el-form :model="dirForm" label-width="100px">
        <el-form-item
          label="文件夹名称"
          prop="name"
          :rules="[{ required: true, message: '请输入文件夹名称', trigger: 'blur' }]"
        >
          <el-input v-model="dirForm.name" placeholder="请输入文件夹名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createDirModalVisible = false">取消</el-button>
          <el-button type="primary" @click="createDirectory">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 重命名模态框 -->
    <el-dialog v-model="renameModalVisible" title="重命名" width="400px">
      <el-form :model="renameForm" label-width="100px">
        <el-form-item
          label="新名称"
          prop="newName"
          :rules="[{ required: true, message: '请输入新名称', trigger: 'blur' }]"
        >
          <el-input v-model="renameForm.newName" placeholder="请输入新名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="renameModalVisible = false">取消</el-button>
          <el-button type="primary" @click="renameFile">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 上传文件模态框 -->
    <el-dialog v-model="uploadModalVisible" title="上传文件" width="400px">
      <el-form :model="uploadForm" label-width="100px">
        <el-form-item label="选择文件" prop="file">
          <el-upload
            :auto-upload="false"
            :limit="1"
            :on-change="beforeUpload"
            :file-list="uploadFileList"
          >
            <el-button type="primary">
              <el-icon><UploadIcon /></el-icon>
              选择文件
            </el-button>
          </el-upload>
        </el-form-item>
        <el-form-item label="上传路径" prop="path">
          <el-input v-model="uploadForm.path" disabled />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="uploadModalVisible = false">取消</el-button>
          <el-button type="primary" @click="uploadFile" :disabled="!uploadForm.file"
            >确定</el-button
          >
        </span>
      </template>
    </el-dialog>
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
  import {
    FolderIcon, FileIcon, RefreshIcon, PlusIcon, UploadIcon,
    FolderOpenIcon, DownloadIcon, EditIcon, DeleteIcon
  } from '@element-plus/icons-vue';

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
      label: 'Root',
      children: []
    });

    // 添加当前目录下的目录
    const currentPathNode = tree[0]; // 根节点
    dirs.forEach(dir => {
      currentPathNode.children.push({
        key: dir.path,
        label: dir.name,
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
  const handleSelect = (node) => {
    if (node.isLeaf) {
      navigateTo(node.data);
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
