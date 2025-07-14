<template>
  <div class="host-management">
    <a-card :bordered="false" class="search-card">
      <a-form layout="inline" :model="searchForm" @submit="handleSearch">
        <a-form-item label="搜索关键字">
          <a-input
            v-model:value="searchForm.keyword"
            placeholder="主机名称、实例ID、IP地址..."
            style="width: 240px"
            @press-enter="handleSearch"
          />
        </a-form-item>
        <a-form-item label="状态">
          <a-select
            v-model:value="searchForm.status"
            placeholder="选择状态"
            style="width: 120px"
            allow-clear
          >
            <a-select-option value="running">运行中</a-select-option>
            <a-select-option value="stopped">已停止</a-select-option>
            <a-select-option value="starting">启动中</a-select-option>
            <a-select-option value="stopping">停止中</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="主机组">
          <a-tree-select
            v-model:value="searchForm.group_id"
            :tree-data="hostGroupOptions"
            placeholder="选择主机组"
            style="width: 200px"
            allow-clear
            tree-default-expand-all
          />
        </a-form-item>
        <a-form-item label="区域">
          <a-input
            v-model:value="searchForm.region"
            placeholder="区域"
            style="width: 120px"
          />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">
            <template #icon><SearchOutlined /></template>
            搜索
          </a-button>
          <a-button @click="handleReset" style="margin-left: 8px">
            <template #icon><ReloadOutlined /></template>
            重置
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card :bordered="false" class="table-card">
      <template #title>
        <div class="table-header">
          <span>主机列表</span>
          <div class="table-actions">
            <a-button type="primary" @click="openHostModal()">
              <template #icon><PlusOutlined /></template>
              添加主机
            </a-button>
            <a-button @click="openManualHostModal()">
              <template #icon><CloudServerOutlined /></template>
              自建主机
            </a-button>
            <a-dropdown>
              <a-button>
                <template #icon><DownloadOutlined /></template>
                批量导入
                <DownOutlined />
              </a-button>
              <template #overlay>
                <a-menu @click="handleBatchAction">
                  <a-menu-item key="import">
                    <template #icon><UploadOutlined /></template>
                    导入主机
                  </a-menu-item>
                  <a-menu-item key="export">
                    <template #icon><DownloadOutlined /></template>
                    导出主机
                  </a-menu-item>
                  <a-menu-item key="template">
                    <template #icon><FileExcelOutlined /></template>
                    下载模板
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
            <a-button @click="syncAllHosts">
              <template #icon><SyncOutlined /></template>
              同步主机
            </a-button>
          </div>
        </div>
      </template>

      <div class="table-alert" v-if="selectedRowKeys.length > 0">
        <a-alert
          :message="`已选择 ${selectedRowKeys.length} 项`"
          type="info"
          show-icon
          closable
          @close="selectedRowKeys = []"
        >
          <template #action>
            <a-button size="small" @click="handleBatchDelete">批量删除</a-button>
            <a-button size="small" @click="handleBatchStatusChange">批量状态变更</a-button>
            <a-button size="small" @click="handleBatchMove">批量移动</a-button>
          </template>
        </a-alert>
      </div>

      <a-table
        :columns="columns"
        :data-source="hostList"
        :loading="isLoading"
        :pagination="paginationConfig"
        :row-selection="rowSelection"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <a @click="viewHostDetail(record)">{{ record.name }}</a>
          </template>
          
          <template v-if="column.key === 'status'">
            <a-tag :color="getStatusColor(record.status)">
              {{ getStatusText(record.status) }}
            </a-tag>
          </template>

          <template v-if="column.key === 'public_ip'">
            <div v-if="record.public_ip && record.public_ip.length > 0">
              <a-tag v-for="ip in record.public_ip" :key="ip">{{ ip }}</a-tag>
            </div>
            <span v-else>-</span>
          </template>

          <template v-if="column.key === 'private_ip'">
            <div v-if="record.private_ip && record.private_ip.length > 0">
              <a-tag v-for="ip in record.private_ip" :key="ip" color="blue">{{ ip }}</a-tag>
            </div>
            <span v-else>-</span>
          </template>

          <template v-if="column.key === 'provider'">
            <a-tag v-if="record.provider" :color="getProviderColor(record.provider.type)">
              {{ record.provider.name }}
            </a-tag>
            <a-tag v-else color="default">自建</a-tag>
          </template>

          <template v-if="column.key === 'group'">
            <span v-if="record.group">{{ record.group.name }}</span>
            <span v-else class="text-gray-400">未分组</span>
          </template>

          <template v-if="column.key === 'expired_at'">
            <span v-if="record.expired_at">
              {{ formatDate(record.expired_at) }}
            </span>
            <span v-else>-</span>
          </template>

          <template v-if="column.key === 'action'">
            <a-space>
              <a-button type="link" size="small" @click="openWebSsh(record)">
                <template #icon><CodeOutlined /></template>
                SSH
              </a-button>
              <a-button type="link" size="small" @click="openSftpWindow(record)">
                <template #icon><FolderOutlined /></template>
                SFTP
              </a-button>
              <a-dropdown>
                <a-button type="link" size="small">
                  更多 <DownOutlined />
                </a-button>
                <template #overlay>
                  <a-menu @click="({ key }) => handleMoreAction(key, record)">
                    <a-menu-item key="edit">
                      <template #icon><EditOutlined /></template>
                      编辑
                    </a-menu-item>
                    <a-menu-item key="sync">
                      <template #icon><SyncOutlined /></template>
                      同步状态
                    </a-menu-item>
                    <a-menu-item key="move">
                      <template #icon><DragOutlined /></template>
                      移动分组
                    </a-menu-item>
                    <a-menu-divider />
                    <a-menu-item key="delete" class="text-red-500">
                      <template #icon><DeleteOutlined /></template>
                      删除
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 主机表单弹窗 -->
    <HostModal 
      ref="hostModalRef" 
      :providers="providerList"
      :host-groups="hostGroupTree"
      @success="handleModalSuccess" 
    />

    <!-- 自建主机表单弹窗 -->
    <ManualHostModal 
      ref="manualHostModalRef" 
      :host-groups="hostGroupTree"
      @success="handleModalSuccess" 
    />

    <!-- SFTP文件管理窗口 -->
    <SftpWindow ref="sftpWindowRef" />

    <!-- 批量导入弹窗 -->
    <BatchImportModal 
      ref="batchImportModalRef" 
      @success="handleModalSuccess" 
    />

    <!-- 批量移动弹窗 -->
    <BatchMoveModal 
      ref="batchMoveModalRef" 
      :host-groups="hostGroupTree"
      @success="handleModalSuccess" 
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { useHostStore } from '@/store/modules/host'
import { storeToRefs } from 'pinia'
import type { Host } from '@/types/api/host'
import { 
  SearchOutlined, 
  ReloadOutlined, 
  PlusOutlined, 
  CloudServerOutlined,
  DownloadOutlined,
  UploadOutlined,
  FileExcelOutlined,
  SyncOutlined,
  CodeOutlined,
  FolderOutlined,
  EditOutlined,
  DeleteOutlined,
  DragOutlined,
  DownOutlined
} from '@ant-design/icons-vue'

// 导入组件
import HostModal from './components/HostModal.vue'
import ManualHostModal from './components/ManualHostModal.vue'
import SftpWindow from './components/SftpWindow.vue'
import BatchImportModal from './components/BatchImportModal.vue'
import BatchMoveModal from './components/BatchMoveModal.vue'

const router = useRouter()
const hostStore = useHostStore()
const { hostList, providerList, hostGroupTree, isLoading, pagination } = storeToRefs(hostStore)

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: undefined,
  group_id: undefined,
  region: ''
})

// 选中的行
const selectedRowKeys = ref<number[]>([])

// 表格行选择配置
const rowSelection = {
  selectedRowKeys: selectedRowKeys,
  onChange: (keys: number[]) => {
    selectedRowKeys.value = keys
  }
}

// 分页配置
const paginationConfig = computed(() => ({
  current: pagination.value.page,
  pageSize: pagination.value.pageSize,
  total: pagination.value.total,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) => 
    `第 ${range[0]}-${range[1]} 项，共 ${total} 项`
}))

// 主机组选项
const hostGroupOptions = computed(() => {
  const buildTreeData = (groups: any[], parentId: number | null = null): any[] => {
    return groups
      .filter(group => group.parent_id === parentId)
      .map(group => ({
        title: group.name,
        value: group.id,
        key: group.id,
        children: buildTreeData(groups, group.id)
      }))
  }
  return buildTreeData(hostGroupTree.value)
})

// 表格列配置
const columns = [
  {
    title: '主机名称',
    dataIndex: 'name',
    key: 'name',
    width: 200,
    ellipsis: true
  },
  {
    title: '实例ID',
    dataIndex: 'instance_id',
    key: 'instance_id',
    width: 150,
    ellipsis: true
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 100
  },
  {
    title: '公网IP',
    dataIndex: 'public_ip',
    key: 'public_ip',
    width: 150
  },
  {
    title: '私网IP',
    dataIndex: 'private_ip',
    key: 'private_ip',
    width: 150
  },
  {
    title: '云账号',
    dataIndex: 'provider',
    key: 'provider',
    width: 120
  },
  {
    title: '主机组',
    dataIndex: 'group',
    key: 'group',
    width: 120
  },
  {
    title: '区域',
    dataIndex: 'region',
    key: 'region',
    width: 100
  },
  {
    title: '操作系统',
    dataIndex: 'os',
    key: 'os',
    width: 120,
    ellipsis: true
  },
  {
    title: '过期时间',
    dataIndex: 'expired_at',
    key: 'expired_at',
    width: 120
  },
  {
    title: '操作',
    key: 'action',
    width: 200,
    fixed: 'right'
  }
]

// 组件引用
const hostModalRef = ref()
const manualHostModalRef = ref()
const sftpWindowRef = ref()
const batchImportModalRef = ref()
const batchMoveModalRef = ref()

// 工具函数
const getStatusColor = (status: string) => {
  const colorMap: Record<string, string> = {
    running: 'green',
    stopped: 'red',
    starting: 'blue',
    stopping: 'orange',
    unknown: 'default'
  }
  return colorMap[status] || 'default'
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    running: '运行中',
    stopped: '已停止',
    starting: '启动中',
    stopping: '停止中',
    unknown: '未知'
  }
  return textMap[status] || status
}

const getProviderColor = (type: string) => {
  const colorMap: Record<string, string> = {
    aliyun: 'orange',
    tencent: 'blue',
    aws: 'yellow',
    manual: 'default'
  }
  return colorMap[type] || 'default'
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString()
}

// 事件处理
const handleSearch = () => {
  fetchHostsWithSearch()
}

const handleReset = () => {
  Object.assign(searchForm, {
    keyword: '',
    status: undefined,
    group_id: undefined,
    region: ''
  })
  fetchHostsWithSearch()
}

const handleTableChange = (pag: any) => {
  pagination.value.page = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchHostsWithSearch()
}

const fetchHostsWithSearch = async () => {
  try {
    await hostStore.fetchHosts({
      page: pagination.value.page,
      page_size: pagination.value.pageSize,
      ...searchForm
    })
  } catch (error) {
    message.error('获取主机列表失败')
  }
}

const handleBatchAction = ({ key }: { key: string }) => {
  switch (key) {
    case 'import':
      batchImportModalRef.value?.open()
      break
    case 'export':
      handleExportHosts()
      break
    case 'template':
      handleDownloadTemplate()
      break
  }
}

const handleExportHosts = async () => {
  try {
    // 实现导出逻辑
    message.success('导出成功')
  } catch (error) {
    message.error('导出失败')
  }
}

const handleDownloadTemplate = () => {
  // 实现下载模板逻辑
  message.info('模板下载功能待实现')
}

const syncAllHosts = async () => {
  try {
    await hostStore.syncHosts()
    message.success('同步成功')
  } catch (error) {
    message.error('同步失败')
  }
}

const handleBatchDelete = () => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 台主机吗？`,
    onOk: async () => {
      try {
        await hostStore.batchDeleteHosts(selectedRowKeys.value)
        selectedRowKeys.value = []
        message.success('删除成功')
      } catch (error) {
        message.error('删除失败')
      }
    }
  })
}

const handleBatchStatusChange = () => {
  // 实现批量状态变更逻辑
  message.info('批量状态变更功能待实现')
}

const handleBatchMove = () => {
  batchMoveModalRef.value?.open(selectedRowKeys.value)
}

const handleMoreAction = (key: string, record: Host) => {
  switch (key) {
    case 'edit':
      openHostModal(record)
      break
    case 'sync':
      syncHostStatus(record.id)
      break
    case 'move':
      // 实现移动逻辑
      message.info('移动功能待实现')
      break
    case 'delete':
      deleteHost(record.id)
      break
  }
}

const openHostModal = (record?: Host) => {
  hostModalRef.value?.open(record)
}

const openManualHostModal = () => {
  manualHostModalRef.value?.open()
}

const openWebSsh = (record: Host) => {
  const routeUrl = router.resolve({
    path: '/webssh',
    query: { host_id: record.id }
  })
  window.open(routeUrl.href, '_blank')
}

const openSftpWindow = (record: Host) => {
  sftpWindowRef.value?.open(record.id)
}

const viewHostDetail = (record: Host) => {
  router.push(`/cmdb/host/${record.id}`)
}

const syncHostStatus = async (hostId: number) => {
  try {
    await hostStore.syncHostStatus(hostId)
    message.success('同步成功')
  } catch (error) {
    message.error('同步失败')
  }
}

const deleteHost = (hostId: number) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除此主机吗？',
    onOk: async () => {
      try {
        await hostStore.deleteHost(hostId)
        message.success('删除成功')
      } catch (error) {
        message.error('删除失败')
      }
    }
  })
}

const handleModalSuccess = () => {
  fetchHostsWithSearch()
  selectedRowKeys.value = []
}

// 初始化
onMounted(async () => {
  try {
    await Promise.all([
      hostStore.fetchProviders(),
      hostStore.fetchHostGroupTree(),
      fetchHostsWithSearch()
    ])
  } catch (error) {
    message.error('初始化失败')
  }
})
</script>

<style scoped lang="scss">
.host-management {
  .search-card {
    margin-bottom: 16px;
  }

  .table-card {
    .table-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .table-actions {
        display: flex;
        gap: 8px;
      }
    }

    .table-alert {
      margin-bottom: 16px;
    }
  }

  .text-gray-400 {
    color: #9ca3af;
  }

  .text-red-500 {
    color: #ef4444;
  }
}
</style> 