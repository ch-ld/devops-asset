<template>
  <div class="host-management-page">
    <!-- 顶部工具栏 -->
    <el-card shadow="hover" class="toolbar-card">
      <el-row :gutter="16" align="middle">
        <el-col :lg="6" :md="8" :sm="12" :xs="24">
          <el-input
            v-model="searchParams.keyword"
            placeholder="搜索主机名、IP地址"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
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
          <el-space>
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
          </el-space>
        </el-col>
      </el-row>
    </el-card>

    <!-- 主机列表 -->
    <el-card shadow="hover" class="host-list-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">主机列表</span>
          <el-tag type="info" effect="plain" class="host-count-tag"
            >{{ pagination.total }} 台</el-tag
          >
          <div class="header-operations">
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
              <el-button type="primary" plain>
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
      </template>

      <!-- 加载中状态 -->
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="5" animated />
      </div>

      <!-- 错误状态显示 -->
      <el-result
        v-else-if="hasError"
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
        :image-size="200"
      >
        <el-button type="primary" @click="handleAddMenuClick('manual')">
          <el-icon><Plus /></el-icon>
          添加主机
        </el-button>
      </el-empty>

      <!-- 主机表格 -->
      <el-table
        v-else
        :data="hostList"
        @selection-change="onSelectionChange"
        style="width: 100%"
        border
        stripe
        highlight-current-row
      >
        <el-table-column type="selection" width="55" />

        <el-table-column prop="name" label="主机名称" min-width="200" sortable>
          <template #default="{ row }">
            <div class="host-name-column">
              <el-tooltip :content="row.instance_id">
                <el-link type="primary" @click="viewHostDetail(row.id)">{{ row.name }}</el-link>
              </el-tooltip>
              <el-tag v-if="row.provider_type" type="success" size="small" effect="light">
                {{ row.provider_type }}
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
                  <span
                    >{{ row.public_ip[0] }}
                    <span v-if="row.public_ip.length > 1">等{{ row.public_ip.length }}个</span>
                  </span>
                </el-tooltip>
              </div>
              <div v-if="row.private_ip && row.private_ip.length > 0">
                <span class="ip-label">内网:</span>
                <el-tooltip :content="row.private_ip.join(', ')">
                  <span
                    >{{ row.private_ip[0] }}
                    <span v-if="row.private_ip.length > 1">等{{ row.private_ip.length }}个</span>
                  </span>
                </el-tooltip>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="配置" min-width="180">
          <template #default="{ row }">
            <div class="configuration-column">
              <div v-if="row.configuration && row.configuration.cpu_cores">
                <el-icon><Cpu /></el-icon> {{ row.configuration.cpu_cores }} 核
              </div>
              <div v-if="row.configuration && row.configuration.memory_size">
                <el-icon><Coin /></el-icon> {{ formatMemorySize(row.configuration.memory_size) }}
              </div>
              <div v-if="row.configuration && row.configuration.instance_type">
                <el-icon><Monitor /></el-icon> {{ row.configuration.instance_type }}
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag
              :type="getStatusType(row.status)"
              effect="light"
              size="small"
              class="status-tag"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="主机组" width="150">
          <template #default="{ row }">
            <span v-if="row.group">{{ row.group.name }}</span>
            <span v-else class="text-muted">未分组</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-dropdown @command="(cmd) => handleRowCommand(cmd, row)" trigger="click">
              <el-button type="primary" link>
                操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="view">
                    <el-icon><View /></el-icon> 查看详情
                  </el-dropdown-item>
                  <el-dropdown-item command="edit">
                    <el-icon><Edit /></el-icon> 编辑
                  </el-dropdown-item>
                  <el-dropdown-item command="terminal">
                    <el-icon><Monitor /></el-icon> SSH终端
                  </el-dropdown-item>
                  <el-dropdown-item command="sftp">
                    <el-icon><FolderOpened /></el-icon> 文件传输
                  </el-dropdown-item>
                  <el-dropdown-item command="sync">
                    <el-icon><RefreshRight /></el-icon> 同步状态
                  </el-dropdown-item>
                  <el-dropdown-item command="restart">
                    <el-icon><Refresh /></el-icon> 重启
                  </el-dropdown-item>
                  <el-dropdown-item command="move">
                    <el-icon><Position /></el-icon> 移动到组
                  </el-dropdown-item>
                  <el-dropdown-item command="delete" divided>
                    <el-icon><Delete /></el-icon> 删除
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>
    </el-card>

    <!-- 主机模态框 -->
    <manual-host-modal
      v-model:visible="hostModalVisible"
      :host="currentHost"
      :is-edit="isEditMode"
      @success="handleModalSuccess"
    />

    <!-- 终端模态框 -->
    <terminal-window v-if="terminalVisible" v-model:visible="terminalVisible" :host="currentHost" />

    <!-- SFTP模态框 -->
    <sftp-window v-if="sftpVisible" v-model:visible="sftpVisible" :host="currentHost" />

    <!-- 批量导入 -->
    <batch-import-modal
      v-if="batchImportModalVisible"
      :visible="batchImportModalVisible"
      @success="handleBatchSuccess"
      @update:visible="(val) => (batchImportModalVisible.value = val)"
    />
  </div>
</template>

<script setup>
  import { ref, reactive, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { useHostStore } from '@/store/modules/host'
  import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
  import {
    Search,
    Refresh,
    DataAnalysis,
    Plus,
    ArrowDown,
    Download,
    Setting,
    View,
    Edit,
    FolderOpened,
    RefreshRight,
    Position,
    Delete,
    Cpu,
    Coin,
    Monitor
  } from '@element-plus/icons-vue'
  import ManualHostModal from './components/ManualHostModal.vue'
  import TerminalWindow from './components/TerminalWindow.vue'
  import SftpWindow from './components/SftpWindow.vue'
  import BatchImportModal from './components/BatchImportModal.vue'

  const router = useRouter()
  const hostStore = useHostStore()

  // 搜索参数
  const searchParams = reactive({
    keyword: '',
    group_id: undefined,
    status: undefined,
    region: undefined
  })

  // 分页参数
  const pagination = reactive({
    current: 1,
    pageSize: 10,
    total: 0
  })

  // 数据状态
  const loading = ref(false)
  const hasError = ref(false)
  const hostList = ref([])
  const hostGroupOptions = ref([])
  const regionOptions = ref([])
  const selectedRowKeys = ref([])

  // 模态框状态
  const hostModalVisible = ref(false)
  const terminalVisible = ref(false)
  const sftpVisible = ref(false)
  const isEditMode = ref(false)
  const currentHost = ref({})
  const batchImportModalVisible = ref(false)

  // 获取主机列表
  const fetchHosts = async () => {
    loading.value = true
    hasError.value = false

    try {
      const result = await hostStore.getHostList({
        page: pagination.current,
        page_size: pagination.pageSize,
        keyword: searchParams.keyword,
        status: searchParams.status,
        group_id: searchParams.group_id,
        region: searchParams.region
      })

      hostList.value = result.data || []
      pagination.total = result.total || 0

      // 提取地区选项
      const regions = new Set()
      hostList.value.forEach((host) => {
        if (host.region) regions.add(host.region)
      })
      regionOptions.value = Array.from(regions)
    } catch (error) {
      console.error('获取主机列表失败:', error)
      hasError.value = true
      hostList.value = []
    } finally {
      loading.value = false
    }
  }

  // 获取主机组选项
  const fetchHostGroups = async () => {
    try {
      const groups = await hostStore.getHostGroupTree()
      hostGroupOptions.value = formatHostGroupOptions(groups)
    } catch (error) {
      console.error('获取主机组失败:', error)
      hostGroupOptions.value = []
    }
  }

  // 格式化主机组选项
  const formatHostGroupOptions = (groups, prefix = '') => {
    let options = []
    groups.forEach((group) => {
      const label = prefix ? `${prefix} / ${group.name}` : group.name
      options.push({
        label,
        value: group.id
      })

      if (group.children && group.children.length > 0) {
        options = options.concat(formatHostGroupOptions(group.children, label))
      }
    })
    return options
  }

  // 获取状态类型
  const getStatusType = (status) => {
    switch (status) {
      case 'running':
        return 'success'
      case 'stopped':
        return 'info'
      case 'error':
        return 'danger'
      case 'expired':
        return 'warning'
      default:
        return 'info'
    }
  }

  // 获取状态文本
  const getStatusText = (status) => {
    switch (status) {
      case 'running':
        return '运行中'
      case 'stopped':
        return '已停止'
      case 'error':
        return '错误'
      case 'expired':
        return '已过期'
      default:
        return status || '未知'
    }
  }

  // 监听页面加载
  onMounted(() => {
    fetchHosts()
    fetchHostGroups()
  })

  // 搜索处理
  const handleSearch = () => {
    pagination.current = 1
    fetchHosts()
  }

  // 重试加载
  const handleRetry = () => {
    fetchHosts()
  }

  // 分页大小改变
  const handleSizeChange = (size) => {
    pagination.pageSize = size
    fetchHosts()
  }

  // 当前页改变
  const handleCurrentChange = (page) => {
    pagination.current = page
    fetchHosts()
  }

  // 表格选择改变
  const onSelectionChange = (selection) => {
    selectedRowKeys.value = selection.map((item) => item.id)
  }

  // 添加主机菜单点击
  const handleAddMenuClick = (command) => {
    switch (command) {
      case 'manual':
        currentHost.value = {}
        isEditMode.value = false
        hostModalVisible.value = true
        break
      case 'batch':
        batchImportModalVisible.value = true
        break
      case 'sync':
        ElMessageBox.confirm('确定要从云提供商同步主机信息吗？此操作可能耗时较长。', '同步主机', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(async () => {
          const loading = ElLoading.service({ text: '正在同步主机...' })
          try {
            await hostStore.syncHosts() // Assuming hostStore has a syncHosts method
            ElMessage.success('同步完成')
            fetchHosts()
          } catch (e) {
            console.error(e)
            ElMessage.error('同步失败')
          } finally {
            loading.close()
          }
        })
        break
    }
  }

  // 显示主机概览
  const showDashboard = () => {
    router.push('/cmdb/host-dashboard')
  }

  // 查看主机详情
  const viewHostDetail = (id) => {
    router.push(`/cmdb/host/${id}`)
  }

  // 打开主机编辑对话框
  const handleEditHost = (host) => {
    currentHost.value = host
    isEditMode.value = true
    hostModalVisible.value = true
  }

  // 行操作命令处理
  const handleRowCommand = (command, row) => {
    switch (command) {
      case 'view':
        viewHostDetail(row.id)
        break
      case 'edit':
        handleEditHost(row)
        break
      case 'delete':
        confirmDeleteHost(row)
        break
      case 'sync':
        syncHostStatus(row.id)
        break
      case 'restart':
        restartHost(row)
        break
      case 'move':
        moveHostToGroup()
        break
      case 'terminal':
        openTerminal(row)
        break
      case 'sftp':
        openSftp(row)
        break
    }
  }

  // 批量操作
  const handleBatchOperation = (command) => {
    if (selectedRowKeys.value.length === 0) {
      ElMessage.warning('请至少选择一个主机')
      return
    }

    switch (command) {
      case 'delete':
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
            fetchHosts()
          } catch (error) {
            ElMessage.error('批量删除主机失败')
            console.error(error)
          }
        })
        break
      default:
        ElMessage.info(`${command} 功能待实现`)
    }
  }

  // 导出主机
  const handleExportHosts = () => {
    ElMessage.info('导出功能待实现')
  }

  // 列设置
  const handleColumnSettings = () => {
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
        fetchHosts()
      } catch (error) {
        ElMessage.error('删除主机失败')
        console.error(error)
      }
    })
  }

  // 同步主机状态
  const syncHostStatus = async (hostId) => {
    try {
      const loading = ElLoading.service({
        lock: true,
        text: '正在同步主机状态...',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      await hostStore.syncHostStatus(hostId)
      ElMessage.success('同步主机状态成功')
      fetchHosts()
      loading.close()
    } catch (error) {
      ElMessage.error('同步主机状态失败')
      console.error(error)
    }
  }

  // 重启主机
  const restartHost = (host) => {
    ElMessageBox.confirm(`确定要重启主机 ${host.name} 吗？`, '重启主机', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      ElMessage.success('已发送重启命令')
    })
  }

  // 移动主机到分组
  const moveHostToGroup = () => {
    ElMessage.info('移动主机功能待实现')
  }

  // 批量导入成功回调
  const handleBatchSuccess = () => {
    batchImportModalVisible.value = false
    fetchHosts()
  }

  // 工具函数
  const formatMemorySize = (size) => {
    if (!size) return '--'
    return size >= 1024 ? `${(size / 1024).toFixed(1)} GB` : `${size} MB`
  }
</script>

<style lang="scss" scoped>
  .host-management-page {
    padding: 16px;

    .toolbar-card {
      margin-bottom: 16px;

      :deep(.el-card__body) {
        padding: 16px;
      }
    }

    .host-list-card {
      :deep(.el-card__header) {
        padding: 12px 20px;
      }

      .card-header {
        display: flex;
        align-items: center;

        .card-title {
          font-size: 16px;
          font-weight: 500;
        }

        .host-count-tag {
          margin-left: 8px;
        }

        .header-operations {
          margin-left: auto;
          display: flex;
          gap: 8px;
        }
      }

      .loading-container {
        padding: 20px 0;
      }

      .host-name-column {
        display: flex;
        flex-direction: column;
        gap: 4px;
      }

      .ip-address-column {
        .ip-label {
          color: var(--el-text-color-secondary);
          margin-right: 4px;
          font-size: 13px;
        }
      }

      .configuration-column {
        font-size: 13px;
        line-height: 1.8;

        .el-icon {
          margin-right: 4px;
          vertical-align: middle;
        }
      }

      .status-tag {
        text-align: center;
        width: 70px;
      }

      .text-muted {
        color: var(--el-text-color-secondary);
        font-size: 13px;
      }

      .pagination-container {
        margin-top: 20px;
        display: flex;
        justify-content: flex-end;
      }
    }
  }
</style>
