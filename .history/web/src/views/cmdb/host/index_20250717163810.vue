<template>
  <div class="host-management-page">
    <!-- 顶部工具栏 -->
    <el-card :body-style="{ padding: '16px' }" shadow="never" class="toolbar-card">
      <el-row :gutter="16" align="middle">
        <el-col :lg="6" :md="8" :sm="12" :xs="24">
          <el-input
            v-model="searchParams.keyword"
            placeholder="搜索主机名、IP地址"
            clearable
            @input="handleSearch"
          >
            <template #append>
              <el-button @click="handleSearch">
                <el-icon><Search /></el-icon>
              </el-button>
            </template>
          </el-input>
        </el-col>

        <el-col :lg="4" :md="6" :sm="12" :xs="12">
          <el-select
            v-model="searchParams.group_id"
            placeholder="主机组"
            style="width: 100%"
            clearable
            @change="handleSearch"
          >
            <el-option :value="undefined" label="全部主机组" />
            <el-option
              v-for="group in hostGroupOptions"
              :key="group.value"
              :value="group.value"
              :label="group.label"
            />
          </el-select>
        </el-col>

        <el-col :lg="4" :md="6" :sm="12" :xs="12">
          <el-select
            v-model="searchParams.status"
            placeholder="状态"
            style="width: 100%"
            clearable
            @change="handleSearch"
          >
            <el-option :value="undefined" label="全部状态" />
            <el-option value="running" label="运行中" />
            <el-option value="stopped" label="已停止" />
            <el-option value="error" label="错误" />
            <el-option value="expired" label="已过期" />
          </el-select>
        </el-col>

        <el-col :lg="4" :md="4" :sm="12" :xs="12">
          <el-select
            v-model="searchParams.region"
            placeholder="地区"
            style="width: 100%"
            clearable
            @change="handleSearch"
          >
            <el-option :value="undefined" label="全部地区" />
            <el-option
              v-for="region in regionOptions"
              :key="region"
              :value="region"
              :label="region"
            />
          </el-select>
        </el-col>

        <el-col :lg="6" :md="24" :sm="24" :xs="24" style="text-align: right">
          <div class="el-space--horizontal" style="gap: 8px; display: flex">
            <el-button @click="handleSearch">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button type="primary" @click="showDashboard">
              <el-icon><DataAnalysis /></el-icon>
              主机概览
            </el-button>
            <el-dropdown @command="handleAddMenuClick">
              <el-button type="primary">
                <el-icon><Plus /></el-icon>
                添加主机
                <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="manual">手动添加主机</el-dropdown-item>
                  <el-dropdown-item command="batch">批量导入主机</el-dropdown-item>
                  <el-dropdown-item command="sync">从云提供商同步</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 主机列表 -->
    <el-card :body-style="{ padding: '16px' }" shadow="never" class="host-list-card">
      <template #header>
        <div class="card-header">
          <span>主机列表</span>
          <el-tag type="info" style="margin-left: 8px">{{ pagination.total }} 台</el-tag>
        </div>
      </template>

      <div class="card-extra" style="margin-bottom: 16px; text-align: right;">
        <div class="el-space--horizontal" style="gap: 8px; display: inline-flex;">
          <el-button-group>
            <el-tooltip content="导出">
              <el-button @click="handleExportHosts">
                <el-icon><Download /></el-icon>
              </el-button>
            </el-tooltip>
            <el-tooltip content="设置">
              <el-button @click="handleColumnSettings">
                <el-icon><Setting /></el-icon>
              </el-button>
            </el-tooltip>
          </el-button-group>
          <el-dropdown v-if="selectedRowKeys.length > 0" @command="handleBatchOperation">
            <el-button>
              批量操作
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="move">移动到主机组</el-dropdown-item>
                <el-dropdown-item command="tags">批量标签管理</el-dropdown-item>
                <el-dropdown-item command="status">修改状态</el-dropdown-item>
                <el-dropdown-item command="terminal">批量执行命令</el-dropdown-item>
                <el-dropdown-item command="sftp">批量文件传输</el-dropdown-item>
                <el-dropdown-item command="delete">批量删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>

      <!-- 错误状态显示 -->
      <el-result
        v-if="hasError"
        status="error"
        title="加载失败"
        sub-title="获取主机列表时发生错误，请检查网络连接或稍后重试。"
      >
        <template #extra>
          <el-button type="primary" @click="handleRetry">
            <el-icon><Refresh /></el-icon>
            重试
          </el-button>
        </template>
      </el-result>

      <!-- 空数据状态 -->
      <el-empty 
        v-else-if="!loading && hostList.length === 0"
        description="暂无主机数据，请添加主机或调整筛选条件"
      >
        <el-button type="primary" @click="handleAddMenuClick('manual')">
          <el-icon><Plus /></el-icon>
          添加主机
        </el-button>
      </el-empty>

      <el-table
        v-else
        :data="hostList"
        :border="false"
        :stripe="true"
        :size="'default'"
        style="width: 100%"
        v-loading="loading"
        @selection-change="onSelectionChange"
        @sort-change="handleTableChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column 
          label="主机名称" 
          prop="name" 
          min-width="200"
          sortable="custom"
        >
          <template #default="{ row }">
            <div class="host-name-column">
              <el-tooltip :content="row.instance_id">
                <a @click="viewHostDetail(row.id)">{{ row.name }}</a>
              </el-tooltip>
              <el-tag v-if="row.provider_type" type="success" size="small">
                {{row.provider_type}}
              </el-tag>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="IP地址" min-width="200">
          <template #default="{ row }">
            <div class="ip-address-column">
              <div v-if="row.public_ip && row.public_ip.length > 0">
                <span class="ip-label">公网:</span>
                <el-tooltip :content="row.public_ip.join(', ')">
                  {{ row.public_ip[0] }}
                  <span v-if="row.public_ip.length > 1">等{{ row.public_ip.length }}个</span>
                </el-tooltip>
              </div>
              <div v-if="row.private_ip && row.private_ip.length > 0">
                <span class="ip-label">内网:</span>
                <el-tooltip :content="row.private_ip.join(', ')">
                  {{ row.private_ip[0] }}
                  <span v-if="row.private_ip.length > 1">等{{ row.private_ip.length }}个</span>
                </el-tooltip>
              </div>
            </div>
          </template>
        </el-table-column>

          <!-- 配置列 -->
          <template v-else-if="column.key === 'configuration'">
            <div class="configuration-column">
              <div v-if="record.configuration && record.configuration.cpu_cores"
                >CPU: {{ record.configuration.cpu_cores }} 核</div
              >
              <div v-if="record.configuration && record.configuration.memory_size"
                >内存: {{ formatMemorySize(record.configuration.memory_size) }}</div
              >
              <div v-if="record.configuration && record.configuration.instance_type"
                >实例类型: {{ record.configuration.instance_type }}</div
              >
            </div>
          </template>

          <!-- 操作系统列 -->
          <template v-else-if="column.key === 'os'">
            <div class="os-column">
              <a-tooltip :title="record.os">
                <span>{{ getOsShortName(record.os) }}</span>
              </a-tooltip>
            </div>
          </template>

          <!-- 状态列 -->
          <template v-else-if="column.key === 'status'">
            <a-tag :color="getStatusColor(record.status)">
              {{ getStatusText(record.status) }}
            </a-tag>
          </template>

          <!-- 到期时间列 -->
          <template v-else-if="column.key === 'expired_at'">
            <div class="expired-at-column">
              <a-tag v-if="record.expired_at" :color="getExpiryColor(record.expired_at)">
                {{ formatExpiryTime(record.expired_at) }}
              </a-tag>
              <span v-else>--</span>
            </div>
          </template>

          <!-- 创建时间列 -->
          <template v-else-if="column.key === 'created_at'">
            <div class="time-info">
              {{ formatTime(record.created_at) }}
            </div>
          </template>

          <!-- 所属分组列 -->
          <template v-else-if="column.key === 'group'">
            <div class="group-column">
              <a-tooltip :title="record.group ? record.group.name : '未分组'">
                <span>{{ record.group ? record.group.name : '未分组' }}</span>
              </a-tooltip>
            </div>
          </template>

          <!-- 操作列 -->
          <template v-else-if="column.key === 'action'">
            <div class="action-column">
              <a-space>
                <a-tooltip title="SSH 终端">
                  <a-button type="link" size="small" @click="openTerminal(record)">
                    <template #icon><CodeIcon /></template>
                  </a-button>
                </a-tooltip>
                <a-tooltip title="SFTP 文件管理">
                  <a-button type="link" size="small" @click="openSftp(record)">
                    <template #icon><FolderIcon /></template>
                  </a-button>
                </a-tooltip>
                <a-dropdown>
                  <a-button type="link" size="small">
                    <template #icon><EllipsisIcon /></template>
                  </a-button>
                  <template #overlay>
                    <a-menu>
                      <a-menu-item key="view" @click="viewHostDetail(record.id)">
                        <EyeIcon /> 查看详情
                      </a-menu-item>
                      <a-menu-item key="edit" @click="handleEditHost(record)">
                        <EditIcon /> 编辑
                      </a-menu-item>
                      <a-menu-item key="move" @click="moveHostToGroup(record)">
                        <SwapIcon /> 移动分组
                      </a-menu-item>
                      <a-divider />
                      <a-menu-item key="restart" @click="restartHost(record)">
                        <ReloadIcon /> 重启主机
                      </a-menu-item>
                      <a-menu-item key="sync" @click="syncHostStatus(record.id)">
                        <SyncIcon /> 同步状态
                      </a-menu-item>
                      <a-divider />
                      <a-menu-item key="delete" @click="confirmDeleteHost(record)">
                        <DeleteIcon /> 删除
                      </a-menu-item>
                    </a-menu>
                  </template>
                </a-dropdown>
              </a-space>
            </div>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 主机模态框组件 -->
    <host-modal
      v-model:visible="hostModalVisible"
      :host="currentHost"
      :is-edit="isEditMode"
      @success="handleModalSuccess"
    />

    <!-- 手动添加主机模态框 -->
    <manual-host-modal v-model:visible="manualHostModalVisible" @success="handleModalSuccess" />

    <!-- 批量导入模态框 -->
    <batch-import-modal v-model:visible="batchImportModalVisible" @success="handleModalSuccess" />

    <!-- 批量移动到主机组模态框 -->
    <batch-move-modal
      v-model:visible="batchMoveModalVisible"
      :host-ids="selectedRowKeys"
      @success="handleModalSuccess"
    />

    <!-- SSH 终端窗口 -->
    <terminal-window v-model:visible="terminalVisible" :host="currentHost" />

    <!-- SFTP 文件管理窗口 -->
    <sftp-window v-model:visible="sftpVisible" :host="currentHost" />

    <!-- 批量标签管理对话框 -->
    <batch-tags-modal
      v-model:visible="batchTagsVisible"
      :selected-host-ids="selectedRowKeys"
      @success="handleSearch"
    />
  </div>
</template>

<script setup>
  import { ref, reactive, computed, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
  import { useHostStore } from '@/store/modules/host'
  import { storeToRefs } from 'pinia'
  import * as hostApi from '@/api/system/host'
  import dayjs from 'dayjs'
  import relativeTime from 'dayjs/plugin/relativeTime'
  import 'dayjs/locale/zh-cn'

  // 配置 dayjs
  dayjs.extend(relativeTime)
  dayjs.locale('zh-cn')

  // 导入 Element Plus 图标
  import {
    Refresh as ReloadIcon,
    DataAnalysis as DashboardIcon,
    Plus as PlusOutlined,
    ArrowDown as DownIcon,
    Download as ExportIcon,
    Setting as SettingIcon,
    Monitor as CodeIcon,
    Folder as FolderIcon,
    MoreFilled as EllipsisIcon,
    View as EyeIcon,
    Edit as EditIcon,
    Sort as SwapIcon,
    Refresh as SyncIcon,
    Delete as DeleteIcon,
    Search,
    ArrowDown,
    Plus
  } from '@element-plus/icons-vue'

  import HostModal from './components/HostModal.vue'
  import ManualHostModal from './components/ManualHostModal.vue'
  import BatchImportModal from './components/BatchImportModal.vue'
  import BatchMoveModal from './components/BatchMoveModal.vue'
  import TerminalWindow from './components/TerminalWindow.vue'
  import SftpWindow from './components/SftpWindow.vue'
  import BatchTagsModal from './components/BatchTagsModal.vue'

  const router = useRouter()
  const hostStore = useHostStore()
  const { hostList, hostGroupTree, pagination, isLoading, error } = storeToRefs(hostStore)
  const loading = computed(() => isLoading.value)
  const hasError = computed(() => !!error.value)

  // 搜索参数
  const searchParams = reactive({
    keyword: '',
    group_id: undefined,
    status: undefined,
    region: undefined,
    page: 1,
    page_size: 10
  })

  // 表格相关
  const selectedRowKeys = ref([])
  const columns = [
    {
      title: '主机名称',
      dataIndex: 'name',
      key: 'name',
      sorter: true,
      width: 200
    },
    {
      title: 'IP地址',
      dataIndex: 'ip',
      key: 'ip',
      width: 200
    },
    {
      title: '配置',
      key: 'configuration',
      width: 180
    },
    {
      title: '系统',
      dataIndex: 'os',
      key: 'os',
      width: 150
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      sorter: true,
      width: 100
    },
    {
      title: '到期时间',
      dataIndex: 'expired_at',
      key: 'expired_at',
      sorter: true,
      width: 150
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      sorter: true,
      width: 150
    },
    {
      title: '所属分组',
      dataIndex: ['group', 'name'],
      key: 'group',
      width: 150
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      fixed: 'right'
    }
  ]

  // 模态框控制
  const hostModalVisible = ref(false)
  const manualHostModalVisible = ref(false)
  const batchImportModalVisible = ref(false)
  const batchMoveModalVisible = ref(false)
  const terminalVisible = ref(false)
  const sftpVisible = ref(false)
  const batchTagsVisible = ref(false)

  // 当前操作的主机
  const currentHost = ref(null)
  const isEditMode = ref(false)

  // 分页配置
  const tablePagination = computed(() => ({
    total: pagination.value.total,
    current: pagination.value.page,
    pageSize: pagination.value.pageSize,
    showSizeChanger: true,
    showQuickJumper: true,
    showTotal: (total) => `共 ${total} 台主机`
  }))

  // 主机组选项
  const hostGroupOptions = computed(() => {
    const options = []
    const processGroups = (groups, parentPath = '') => {
      groups.forEach((group) => {
        const label = parentPath ? `${parentPath} / ${group.name}` : group.name
        options.push({ label, value: group.id })

        // 查找子组
        const children = hostGroupTree.value.filter((g) => g.parent_id === group.id)
        if (children.length > 0) {
          processGroups(children, label)
        }
      })
    }

    // 从根节点开始处理
    const rootGroups = hostGroupTree.value.filter((g) => !g.parent_id)
    processGroups(rootGroups)

    return options
  })

  // 地区选项
  const regionOptions = computed(() => {
    const regions = new Set()
    hostList.value.forEach((host) => {
      if (host.region) {
        regions.add(host.region)
      }
    })
    return Array.from(regions)
  })

  // 获取主机组列表
  const hostGroups = ref([])
  const fetchHostGroups = async () => {
    try {
      hostGroups.value = await hostApi.getHostGroupTree()
    } catch (error) {
      console.error('获取主机组失败', error)
      ElMessage.error('获取主机组失败')
    }
  }

  // 获取云账号列表
  const providers = ref([])
  const fetchProviders = async () => {
    try {
      providers.value = await hostApi.getProviderList()
    } catch (error) {
      console.error('获取云账号失败', error)
      ElMessage.error('获取云账号失败')
    }
  }

  // 初始化数据
  onMounted(() => {
    fetchHosts()
    fetchHostGroups()
    fetchProviders()
  })

  // 方法
  const fetchHosts = async () => {
    try {
      await hostStore.fetchHosts({
        ...searchParams,
        page: searchParams.page,
        page_size: searchParams.page_size
      })
    } catch (error) {
      ElMessage.error('获取主机列表失败')
      console.error(error)
    }
  }

  const handleSearch = () => {
    searchParams.page = 1
    fetchHosts()
  }

  const handleTableChange = (pag, filters, sorter) => {
    searchParams.page = pag.current
    searchParams.page_size = pag.pageSize

    // 处理排序
    if (sorter.field && sorter.order) {
      // 排序逻辑待实现
      // const orderMap = {
      //   ascend: 'asc',
      //   descend: 'desc'
      // }
      // 这里可能需要修改API支持排序
    }

    fetchHosts()
  }

  const onSelectionChange = (keys) => {
    selectedRowKeys.value = keys
  }

  // 跳转到主机仪表盘
  const showDashboard = () => {
    router.push('/cmdb/host/dashboard')
  }

  // 添加主机相关
  const handleAddMenuClick = ({ command }) => {
    if (command === 'manual') {
      manualHostModalVisible.value = true
    } else if (command === 'batch') {
      batchImportModalVisible.value = true
    } else if (command === 'sync') {
      ElMessageBox.confirm('确定要从云提供商同步主机信息吗？这可能需要一些时间。', '同步主机', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }).then(async () => {
        try {
          ElLoading.service({
            lock: true,
            text: '正在同步主机信息...',
            background: 'rgba(0, 0, 0, 0.7)'
          })
          await hostStore.syncHosts()
          ElMessage.success('同步主机信息成功')
          fetchHosts()
        } catch (error) {
          ElMessage.error('同步主机信息失败')
          console.error(error)
        } finally {
          ElLoading.service().close()
        }
      })
    }
  }

  // 查看主机详情
  const viewHostDetail = (id) => {
    router.push(`/cmdb/host/${id}`)
  }

  // 打开主机编辑对话框
  const handleEditHost = (host = null) => {
    currentHost.value = host
    isEditMode.value = !!host
    hostModalVisible.value = true
  }

  // 删除主机
  const confirmDeleteHost = (host) => {
    ElMessageBox.confirm(`确定要删除主机 ${host.name} 吗？此操作不可恢复。`, '删除主机', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        await hostStore.deleteHost(host.id)
        ElMessage.success('删除主机成功')
        fetchHosts() // 刷新列表
      } catch (err) {
        ElMessage.error('删除主机失败')
        console.error(err)
      }
    })
  }

  // 主机移动分组
  const moveHostToGroup = (host) => {
    currentHost.value = host
    batchMoveModalVisible.value = true
  }

  // 重启主机
  const restartHost = (host) => {
    ElMessageBox.confirm(`确定要重启主机 ${host.name} 吗？`, '重启主机', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        // 假设有重启主机的API
        ElMessage.loading('正在重启主机...')
        // await hostApi.restartHost(host.id);
        ElMessage.success('已发送重启命令')
      } catch (error) {
        ElMessage.error('重启主机失败')
        console.error(error)
      }
    })
  }

  // 同步主机状态
  const syncHostStatus = async (hostId) => {
    try {
      ElLoading.service({
        lock: true,
        text: '正在同步主机状态...',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      await hostStore.syncHostStatus(hostId)
      ElMessage.success('同步主机状态成功')
      fetchHosts() // 刷新列表
    } catch (err) {
      ElMessage.error('同步主机状态失败')
      console.error(err)
    } finally {
      ElLoading.service().close()
    }
  }

  // 批量操作
  const handleBatchOperation = ({ key }) => {
    if (selectedRowKeys.value.length === 0) {
      ElMessage.warning('请至少选择一个主机')
      return
    }

    if (key === 'move') {
      batchMoveModalVisible.value = true
    } else if (key === 'tags') {
      batchTagsVisible.value = true
    } else if (key === 'status') {
      // 实现批量修改状态
      // TODO: 实现批量修改状态的对话框
    } else if (key === 'terminal') {
      // 实现批量执行命令
      // TODO: 实现批量执行命令的对话框
    } else if (key === 'sftp') {
      // 实现批量文件传输
      // TODO: 实现批量文件传输的对话框
    } else if (key === 'delete') {
      ElMessageBox.confirm(
        `确定要删除选中的 ${selectedRowKeys.value.length} 台主机吗？此操作不可恢复。`,
        '批量删除主机',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).then(async () => {
        try {
          await hostStore.batchDeleteHosts(selectedRowKeys.value)
          ElMessage.success('批量删除主机成功')
          selectedRowKeys.value = []
          fetchHosts() // 刷新列表
        } catch (err) {
          ElMessage.error('批量删除主机失败')
          console.error(err)
        }
      })
    }
  }

  // 导出主机
  const handleExportHosts = async () => {
    try {
      const loading = ElLoading.service({
        lock: true,
        text: '正在准备导出数据...',
        background: 'rgba(0, 0, 0, 0.7)'
      })

      const response = await hostApi.batchExportHosts({
        format: 'excel'
        // 可以传递筛选条件
      })

      // 创建下载链接
      const blob = new Blob([response], {
        type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
      })
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', `主机列表_${dayjs().format('YYYY-MM-DD')}.xlsx`)
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)

      ElMessage.success('导出成功')
      loading.close()
    } catch (err) {
      ElMessage.error('导出失败')
      console.error(err)
    }
  }

  // 列设置
  const handleColumnSettings = () => {
    // 实现列设置功能
    ElMessage.info('列设置功能待实现')
  }

  // 打开终端
  const openTerminal = (host) => {
    currentHost.value = host
    terminalVisible.value = true
  }

  // 打开SFTP窗口
  const openSftp = (host) => {
    currentHost.value = host
    sftpVisible.value = true
  }

  // 模态框成功回调
  const handleModalSuccess = () => {
    fetchHosts()
  }

  // 重试方法
  const handleRetry = () => {
    fetchHosts()
  }

  // 工具函数
  const formatMemorySize = (size) => {
    if (!size) return '--'
    return size >= 1024 ? `${(size / 1024).toFixed(1)} GB` : `${size} MB`
  }

  const getOsShortName = (os) => {
    if (!os) return '--'
    if (os.toLowerCase().includes('windows')) return 'Windows'
    if (os.toLowerCase().includes('ubuntu')) return 'Ubuntu'
    if (os.toLowerCase().includes('centos')) return 'CentOS'
    if (os.toLowerCase().includes('debian')) return 'Debian'
    if (os.toLowerCase().includes('linux')) return 'Linux'
    return os
  }

  const getStatusColor = (status) => {
    if (!status) return 'default'
    const colorMap = {
      running: 'green',
      stopped: 'orange',
      error: 'red',
      expired: 'red',
      starting: 'blue',
      stopping: 'orange',
      rebooting: 'blue'
    }
    return colorMap[status.toLowerCase()] || 'default'
  }

  const getStatusText = (status) => {
    if (!status) return '--'
    const textMap = {
      running: '运行中',
      stopped: '已停止',
      error: '错误',
      expired: '已过期',
      starting: '启动中',
      stopping: '停止中',
      rebooting: '重启中'
    }
    return textMap[status.toLowerCase()] || status
  }

  const getExpiryColor = (expiryDate) => {
    if (!expiryDate) return 'default'

    const now = new Date()
    const expiry = new Date(expiryDate)
    const diffDays = Math.floor((expiry.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))

    if (diffDays < 0) return 'red'
    if (diffDays < 7) return 'orange'
    if (diffDays < 30) return 'gold'
    return 'green'
  }

  const formatExpiryTime = (expiryDate) => {
    if (!expiryDate) return '--'
    try {
      const date = new Date(expiryDate)
      const now = new Date()

      if (date < now) {
        return '已过期'
      }

      return dayjs(expiryDate).fromNow()
    } catch {
      return expiryDate
    }
  }

  // 格式化时间
  // 该函数在模板中使用
  const formatTime = (timestamp) => {
    if (!timestamp) return '-'
    // 如果是Unix时间戳（数字），需要转换为毫秒
    const date = typeof timestamp === 'number' ? new Date(timestamp * 1000) : new Date(timestamp)
    return date.toLocaleString('zh-CN')
  }
</script>

<style lang="scss" scoped>
  .host-management-page {
    .toolbar-card {
      margin-bottom: 16px;
    }

    .host-list-card {
      .ant-card-head {
        min-height: 48px;
      }

      .host-name-column {
        display: flex;
        flex-direction: column;

        .ant-tag {
          margin-top: 4px;
        }
      }

      .ip-address-column {
        .ip-label {
          color: rgba(0, 0, 0, 0.45);
          margin-right: 4px;
        }
      }

      .configuration-column {
        font-size: 12px;
        line-height: 1.5;
      }

      .action-column {
        white-space: nowrap;
      }
    }
  }

  // 暗色主题适配
  html.dark {
    .host-management-page {
      .ip-address-column {
        .ip-label {
          color: rgba(255, 255, 255, 0.45);
        }
      }
    }
  }
</style>
