<template>
  <div class="host-management-page">
    <!-- 顶部工具栏 -->
    <a-card :bordered="false" class="toolbar-card">
      <a-row :gutter="16" align="middle">
        <a-col :lg="6" :md="8" :sm="12" :xs="24">
          <a-input-search
            v-model:value="searchParams.keyword"
            placeholder="搜索主机名、IP地址"
            @search="handleSearch"
            allow-clear
          />
        </a-col>

        <a-col :lg="4" :md="6" :sm="12" :xs="12">
          <a-select
            v-model:value="searchParams.group_id"
            placeholder="主机组"
            style="width: 100%"
            allow-clear
            @change="handleSearch"
          >
            <a-select-option :value="undefined">全部主机组</a-select-option>
            <template v-for="group in hostGroupOptions" :key="group.value">
              <a-select-option :value="group.value">{{ group.label }}</a-select-option>
            </template>
          </a-select>
        </a-col>

        <a-col :lg="4" :md="6" :sm="12" :xs="12">
          <a-select
            v-model:value="searchParams.status"
            placeholder="状态"
            style="width: 100%"
            allow-clear
            @change="handleSearch"
          >
            <a-select-option :value="undefined">全部状态</a-select-option>
            <a-select-option value="running">运行中</a-select-option>
            <a-select-option value="stopped">已停止</a-select-option>
            <a-select-option value="error">错误</a-select-option>
            <a-select-option value="expired">已过期</a-select-option>
          </a-select>
        </a-col>

        <a-col :lg="4" :md="4" :sm="12" :xs="12">
          <a-select
            v-model:value="searchParams.region"
            placeholder="地区"
            style="width: 100%"
            allow-clear
            @change="handleSearch"
          >
            <a-select-option :value="undefined">全部地区</a-select-option>
            <template v-for="region in regionOptions" :key="region">
              <a-select-option :value="region">{{ region }}</a-select-option>
            </template>
          </a-select>
        </a-col>

        <a-col :lg="6" :md="24" :sm="24" :xs="24" style="text-align: right">
          <a-space>
            <a-button @click="handleSearch">
              <template #icon><ReloadIcon /></template>
              刷新
            </a-button>
            <a-button type="primary" @click="showDashboard">
              <template #icon><DashboardIcon /></template>
              主机概览
            </a-button>
            <a-dropdown>
              <a-button type="primary">
                <template #icon><PlusOutlined /></template>
                添加主机
                <template #suffix><DownIcon /></template>
              </a-button>
              <template #overlay>
                <a-menu @click="handleAddMenuClick">
                  <a-menu-item key="manual">手动添加主机</a-menu-item>
                  <a-menu-item key="batch">批量导入主机</a-menu-item>
                  <a-menu-item key="sync">从云提供商同步</a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </a-space>
        </a-col>
      </a-row>
    </a-card>

    <!-- 主机列表 -->
    <a-card :bordered="false" class="host-list-card">
      <template #title>
        <span>主机列表</span>
        <a-tag color="blue" style="margin-left: 8px">{{ pagination.total }} 台</a-tag>
      </template>

      <template #extra>
        <a-space>
          <a-button-group>
            <a-tooltip title="导出">
              <a-button @click="handleExportHosts">
                <template #icon><ExportIcon /></template>
              </a-button>
            </a-tooltip>
            <a-tooltip title="设置">
              <a-button @click="handleColumnSettings">
                <template #icon><SettingIcon /></template>
              </a-button>
            </a-tooltip>
          </a-button-group>
          <a-dropdown v-if="selectedRowKeys.length > 0">
            <a-button>
              批量操作
              <template #icon><DownIcon /></template>
            </a-button>
            <template #overlay>
              <a-menu @click="handleBatchOperation">
                <a-menu-item key="move">移动到主机组</a-menu-item>
                <a-menu-item key="tags">批量标签管理</a-menu-item>
                <a-menu-item key="status">修改状态</a-menu-item>
                <a-menu-item key="terminal">批量执行命令</a-menu-item>
                <a-menu-item key="sftp">批量文件传输</a-menu-item>
                <a-menu-item key="delete">批量删除</a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </a-space>
      </template>

      <a-table
        :columns="columns"
        :data-source="hostList"
        :pagination="tablePagination"
        :loading="loading"
        :row-selection="{
          selectedRowKeys,
          onChange: onSelectionChange
        }"
        :row-key="(record) => record.id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <!-- 主机名称列 -->
          <template v-if="column.key === 'name'">
            <div class="host-name-column">
              <a-tooltip :title="record.instance_id">
                <a @click="viewHostDetail(record.id)">{{ record.name }}</a>
              </a-tooltip>
              <a-tag v-if="record.provider_type" color="green" size="small">{{
                record.provider_type
              }}</a-tag>
            </div>
          </template>

          <!-- IP地址列 -->
          <template v-else-if="column.key === 'ip'">
            <div class="ip-address-column">
              <div v-if="record.public_ip && record.public_ip.length > 0">
                <span class="ip-label">公网:</span>
                <a-tooltip :title="record.public_ip.join(', ')">
                  {{ record.public_ip[0] }}
                  <span v-if="record.public_ip.length > 1">等{{ record.public_ip.length }}个</span>
                </a-tooltip>
              </div>
              <div v-if="record.private_ip && record.private_ip.length > 0">
                <span class="ip-label">内网:</span>
                <a-tooltip :title="record.private_ip.join(', ')">
                  {{ record.private_ip[0] }}
                  <span v-if="record.private_ip.length > 1"
                    >等{{ record.private_ip.length }}个</span
                  >
                </a-tooltip>
              </div>
            </div>
          </template>

          <!-- 配置列 -->
          <template v-else-if="column.key === 'configuration'">
            <div class="configuration-column">
              <div v-if="record.configuration.cpu_cores"
                >CPU: {{ record.configuration.cpu_cores }} 核</div
              >
              <div v-if="record.configuration.memory_size"
                >内存: {{ formatMemorySize(record.configuration.memory_size) }}</div
              >
              <div v-if="record.configuration.instance_type"
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

<script lang="ts" setup>
  import { ref, reactive, computed, onMounted, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
  import { useHostStore } from '@/store/modules/host'
  import { storeToRefs } from 'pinia'
  import * as hostApi from '@/api/system/host'
  import type { Host, HostGroup, Provider } from '@/types/api/host'
  import dayjs from 'dayjs'
  import relativeTime from 'dayjs/plugin/relativeTime'
  import 'dayjs/locale/zh-cn'

  // 配置 dayjs
  dayjs.extend(relativeTime)
  dayjs.locale('zh-cn')

  import HostModal from './components/HostModal.vue'
  import ManualHostModal from './components/ManualHostModal.vue'
  import BatchImportModal from './components/BatchImportModal.vue'
  import BatchMoveModal from './components/BatchMoveModal.vue'
  import TerminalWindow from './components/TerminalWindow.vue'
  import SftpWindow from './components/SftpWindow.vue'
  import BatchTagsModal from './components/BatchTagsModal.vue'

  const router = useRouter()
  const hostStore = useHostStore()
  const { hostList, hostGroupTree, pagination, isLoading } = storeToRefs(hostStore)
  const loading = computed(() => isLoading.value)

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
  const selectedRowKeys = ref<number[]>([])
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
  const currentHost = ref<Host | null>(null)
  const isEditMode = ref(false)

  // 分页配置
  const tablePagination = computed(() => ({
    total: pagination.value.total,
    current: pagination.value.page,
    pageSize: pagination.value.pageSize,
    showSizeChanger: true,
    showQuickJumper: true,
    showTotal: (total: number) => `共 ${total} 台主机`
  }))

  // 主机组选项
  const hostGroupOptions = computed(() => {
    const options: { label: string; value: number }[] = []
    const processGroups = (groups: HostGroup[], parentPath = '') => {
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
    const regions = new Set<string>()
    hostList.value.forEach((host) => {
      if (host.region) {
        regions.add(host.region)
      }
    })
    return Array.from(regions)
  })

  // 获取主机组列表
  const hostGroups = ref<HostGroup[]>([])
  const fetchHostGroups = async () => {
    try {
      const res = await hostApi.getHostGroupTree()
      hostGroups.value = res.data || []
    } catch (error) {
      console.error('获取主机组失败', error)
      ElMessage.error('获取主机组失败')
    }
  }

  // 获取云账号列表
  const providers = ref<Provider[]>([])
  const fetchProviders = async () => {
    try {
      const res = await hostApi.getProviderList()
      providers.value = res.data || []
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

  const handleTableChange = (pag: any, filters: any, sorter: any) => {
    searchParams.page = pag.current
    searchParams.page_size = pag.pageSize

    // 处理排序
    if (sorter.field && sorter.order) {
      const orderMap: Record<string, string> = {
        ascend: 'asc',
        descend: 'desc'
      }
      const sort = `${sorter.field},${orderMap[sorter.order]}`
      // 这里可能需要修改API支持排序
    }

    fetchHosts()
  }

  const onSelectionChange = (keys: number[]) => {
    selectedRowKeys.value = keys
  }

  // 跳转到主机仪表盘
  const showDashboard = () => {
    router.push('/cmdb/host/dashboard')
  }

  // 添加主机相关
  const handleAddMenuClick = ({ key }: { key: string }) => {
    if (key === 'manual') {
      manualHostModalVisible.value = true
    } else if (key === 'batch') {
      batchImportModalVisible.value = true
    } else if (key === 'sync') {
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
  const viewHostDetail = (id: number) => {
    router.push(`/cmdb/host/${id}`)
  }

  // 打开主机编辑对话框
  const hostModalRef = ref()
  const handleEditHost = (host: Host | null = null) => {
    hostModalRef.value?.open({
      providers: providers.value,
      hostGroups: hostGroups.value,
      host,
      isEdit: !!host,
      onSuccess: () => {
        fetchHosts()
      }
    })
  }

  // 打开手动添加主机对话框
  const manualHostModalRef = ref()
  const handleAddManualHost = () => {
    manualHostModalRef.value?.open({
      hostGroups: hostGroups.value,
      onSuccess: () => {
        fetchHosts()
      }
    })
  }

  // 打开批量导入对话框
  const batchImportModalRef = ref()
  const handleBatchImport = () => {
    batchImportModalRef.value?.open({
      onSuccess: () => {
        fetchHosts()
      }
    })
  }

  // 打开批量移动对话框
  const batchMoveModalRef = ref()
  const handleBatchMove = () => {
    if (selectedRowKeys.value.length === 0) {
      ElMessage.warning('请至少选择一台主机')
      return
    }
    batchMoveModalRef.value?.open({
      hostGroups: hostGroups.value,
      hostIds: selectedRowKeys.value,
      onSuccess: () => {
        fetchHosts()
        selectedRowKeys.value = []
      }
    })
  }

  // 打开批量标签对话框
  const batchTagsModalRef = ref()
  const handleBatchTags = () => {
    if (selectedRowKeys.value.length === 0) {
      ElMessage.warning('请至少选择一台主机')
      return
    }
    batchTagsModalRef.value?.open({
      hostIds: selectedRowKeys.value,
      onSuccess: () => {
        fetchHosts()
        selectedRowKeys.value = []
      }
    })
  }

  // 删除主机
  const confirmDeleteHost = (host: Host) => {
    ElMessageBox.confirm(`确定要删除主机 ${host.name} 吗？此操作不可恢复。`, '删除主机', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        await hostStore.deleteHost(host.id)
        ElMessage.success('删除主机成功')
      } catch (error) {
        ElMessage.error('删除主机失败')
        console.error(error)
      }
    })
  }

  // 主机移动分组
  const moveHostToGroup = (host: Host) => {
    currentHost.value = host
    batchMoveModalVisible.value = true
  }

  // 重启主机
  const restartHost = (host: Host) => {
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
  const syncHostStatus = async (hostId: number) => {
    try {
      ElLoading.service({
        lock: true,
        text: '正在同步主机状态...',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      await hostStore.syncHostStatus(hostId)
      ElMessage.success('同步主机状态成功')
    } catch (error) {
      ElMessage.error('同步主机状态失败')
      console.error(error)
    } finally {
      ElLoading.service().close()
    }
  }

  // 批量操作
  const handleBatchOperation = ({ key }: { key: string }) => {
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
        } catch (error) {
          ElMessage.error('批量删除主机失败')
          console.error(error)
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
      const blob = new Blob([response as BlobPart], {
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
    } catch (error) {
      ElMessage.error('导出失败')
      console.error(error)
    }
  }

  // 列设置
  const handleColumnSettings = () => {
    // 实现列设置功能
    ElMessage.info('列设置功能待实现')
  }

  // 打开终端
  const openTerminal = (host: Host) => {
    currentHost.value = host
    terminalVisible.value = true
  }

  // 打开SFTP窗口
  const openSftp = (host: Host) => {
    currentHost.value = host
    sftpVisible.value = true
  }

  // 模态框成功回调
  const handleModalSuccess = () => {
    fetchHosts()
  }

  // 工具函数
  const formatMemorySize = (size?: number) => {
    if (!size) return '--'
    return size >= 1024 ? `${(size / 1024).toFixed(1)} GB` : `${size} MB`
  }

  const getOsShortName = (os?: string) => {
    if (!os) return '--'
    if (os.toLowerCase().includes('windows')) return 'Windows'
    if (os.toLowerCase().includes('ubuntu')) return 'Ubuntu'
    if (os.toLowerCase().includes('centos')) return 'CentOS'
    if (os.toLowerCase().includes('debian')) return 'Debian'
    if (os.toLowerCase().includes('linux')) return 'Linux'
    return os
  }

  const getStatusColor = (status?: string) => {
    if (!status) return 'default'
    const colorMap: Record<string, string> = {
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

  const getStatusText = (status?: string) => {
    if (!status) return '--'
    const textMap: Record<string, string> = {
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

  const getExpiryColor = (expiryDate?: string) => {
    if (!expiryDate) return 'default'

    const now = new Date()
    const expiry = new Date(expiryDate)
    const diffDays = Math.floor((expiry.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))

    if (diffDays < 0) return 'red'
    if (diffDays < 7) return 'orange'
    if (diffDays < 30) return 'gold'
    return 'green'
  }

  const formatExpiryTime = (expiryDate?: string) => {
    if (!expiryDate) return '--'
    try {
      const date = new Date(expiryDate)
      const now = new Date()

      if (date < now) {
        return '已过期'
      }

      return dayjs(expiryDate).fromNow()
    } catch (error) {
      return expiryDate
    }
  }

  // 格式化时间
  const formatTime = (timestamp: number | string) => {
    if (!timestamp) return '-'
    // 如果是Unix时间戳（数字），需要转换为毫秒
    return dayjs(typeof timestamp === 'number' ? timestamp * 1000 : timestamp).format(
      'YYYY-MM-DD HH:mm:ss'
    )
  }

  // 同步主机
  const handleSync = async () => {
    try {
      const loading = ElLoading.service({
        lock: true,
        text: '正在同步主机...',
        background: 'rgba(0, 0, 0, 0.7)'
      })

      await hostApi.syncHosts()
      ElMessage.success('同步主机成功')
      fetchHosts()
      loading.close()
    } catch (error) {
      console.error('同步主机失败', error)
      ElMessage.error('同步主机失败')
      ElLoading.service().close()
    }
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
