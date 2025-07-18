<template>
  <div class="host-management-page">
    <!-- 顶部工具栏 -->
    <el-card shadow="never" class="toolbar-card">
      <el-row :gutter="16">
        <el-col :lg="6" :md="8" :sm="12" :xs="24">
          <el-input
            v-model="searchParams.keyword"
            placeholder="搜索主机名、IP地址"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #suffix>
              <el-icon @click="handleSearch"><Search /></el-icon>
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
    <el-card shadow="never" class="host-list-card">
      <template #header>
        <div class="card-header">
          <span>主机列表</span>
          <el-tag type="info" style="margin-left: 8px">{{ pagination.total }} 台</el-tag>
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
      </template>

      <el-table
        :data="hostList"
        v-loading="loading"
        @selection-change="onSelectionChange"
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="主机名称" min-width="200" sortable>
          <template #default="{ row }">
            <div class="host-name-column">
              <el-tooltip :content="row.instance_id">
                <el-link @click="viewHostDetail(row.id)">{{ row.name }}</el-link>
              </el-tooltip>
              <el-tag v-if="row.provider_type" type="success" size="small">
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

        <el-table-column label="配置" min-width="180">
          <template #default="{ row }">
            <div class="configuration-column">
              <div v-if="row.configuration && row.configuration.cpu_cores">
                CPU: {{ row.configuration.cpu_cores }} 核
              </div>
              <div v-if="row.configuration && row.configuration.memory_size">
                内存: {{ formatMemorySize(row.configuration.memory_size) }}
              </div>
              <div v-if="row.configuration && row.configuration.instance_type">
                实例类型: {{ row.configuration.instance_type }}
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-space>
              <el-tooltip content="SSH 终端">
                <el-button type="primary" link @click="openTerminal(row)">
                  <el-icon><Monitor /></el-icon>
                </el-button>
              </el-tooltip>
              <el-tooltip content="SFTP 文件管理">
                <el-button type="primary" link @click="openSftp(row)">
                  <el-icon><Folder /></el-icon>
                </el-button>
              </el-tooltip>
              <el-dropdown @command="(cmd) => handleRowCommand(cmd, row)">
                <el-button type="primary" link>
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="view">
                      <el-icon><View /></el-icon> 查看详情
                    </el-dropdown-item>
                    <el-dropdown-item command="edit">
                      <el-icon><Edit /></el-icon> 编辑
                    </el-dropdown-item>
                    <el-dropdown-item command="move">
                      <el-icon><Sort /></el-icon> 移动分组
                    </el-dropdown-item>
                    <el-dropdown-item command="restart">
                      <el-icon><Refresh /></el-icon> 重启主机
                    </el-dropdown-item>
                    <el-dropdown-item command="sync">
                      <el-icon><Refresh /></el-icon> 同步状态
                    </el-dropdown-item>
                    <el-dropdown-item command="delete" divided>
                      <el-icon><Delete /></el-icon> 删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </el-space>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 组件引用 -->
    <host-modal
      v-model:visible="hostModalVisible"
      :host="currentHost"
      :is-edit="isEditMode"
      @success="handleModalSuccess"
    />

    <terminal-window v-model:visible="terminalVisible" :host="currentHost" />
    <sftp-window v-model:visible="sftpVisible" :host="currentHost" />
  </div>
</template>

<script setup>
  import { ref, reactive, computed, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
  import { useHostStore } from '@/store/modules/host'
  import { storeToRefs } from 'pinia'
  import * as hostApi from '@/api/system/host'

  // 导入 Element Plus 图标
  import {
    Search,
    Refresh,
    DataAnalysis,
    Plus,
    ArrowDown,
    Download,
    Setting,
    Monitor,
    Folder,
    MoreFilled,
    View,
    Edit,
    Sort,
    Delete
  } from '@element-plus/icons-vue'

  import HostModal from './components/HostModal.vue'
  import TerminalWindow from './components/TerminalWindow.vue'
  import SftpWindow from './components/SftpWindow.vue'

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
  const selectedRowKeys = ref([])
  const selectedRows = ref([])

  // 模态框控制
  const hostModalVisible = ref(false)
  const terminalVisible = ref(false)
  const sftpVisible = ref(false)

  // 当前操作的主机
  const currentHost = ref(null)
  const isEditMode = ref(false)

  // 主机组选项
  const hostGroupOptions = computed(() => {
    const options = []
    const processGroups = (groups, parentPath = '') => {
      if (!groups || !Array.isArray(groups)) return

      groups.forEach((group) => {
        const label = parentPath ? `${parentPath} / ${group.name}` : group.name
        options.push({ label, value: group.id })

        // 查找子组
        const children = (hostGroupTree.value || []).filter((g) => g.parent_id === group.id)
        if (children.length > 0) {
          processGroups(children, label)
        }
      })
    }

    // 从根节点开始处理
    const rootGroups = (hostGroupTree.value || []).filter((g) => !g.parent_id)
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

  // 初始化数据
  onMounted(() => {
    fetchHosts()
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

  const handleSizeChange = (size) => {
    searchParams.page_size = size
    fetchHosts()
  }

  const handleCurrentChange = (page) => {
    searchParams.page = page
    fetchHosts()
  }

  const onSelectionChange = (rows) => {
    selectedRows.value = rows
    selectedRowKeys.value = rows.map((row) => row.id)
  }

  // 跳转到主机仪表盘
  const showDashboard = () => {
    router.push('/cmdb/host/dashboard')
  }

  // 添加主机相关
  const handleAddMenuClick = (command) => {
    if (command === 'manual') {
      hostModalVisible.value = true
      currentHost.value = null
      isEditMode.value = false
    } else if (command === 'batch') {
      ElMessage.info('批量导入功能待实现')
    } else if (command === 'sync') {
      ElMessageBox.confirm('确定要从云提供商同步主机信息吗？这可能需要一些时间。', '同步主机', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }).then(async () => {
        try {
          const loading = ElLoading.service({
            lock: true,
            text: '正在同步主机信息...',
            background: 'rgba(0, 0, 0, 0.7)'
          })
          await hostStore.syncHosts()
          ElMessage.success('同步主机信息成功')
          fetchHosts()
          loading.close()
        } catch (error) {
          ElMessage.error('同步主机信息失败')
          console.error(error)
        }
      })
    }
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
        moveHostToGroup(row)
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
  const moveHostToGroup = (host) => {
    ElMessage.info('移动主机功能待实现')
  }

  // 工具函数
  const formatMemorySize = (size) => {
    if (!size) return '--'
    return size >= 1024 ? `${(size / 1024).toFixed(1)} GB` : `${size} MB`
  }
</script>

<style lang="scss" scoped>
  .host-management-page {
    .toolbar-card {
      margin-bottom: 16px;
    }

    .host-list-card {
      .card-header {
        display: flex;
        align-items: center;

        .header-operations {
          margin-left: auto;
          display: flex;
          gap: 8px;
        }
      }

      .host-name-column {
        display: flex;
        flex-direction: column;
        gap: 4px;
      }

      .ip-address-column {
        .ip-label {
          color: #909399;
          margin-right: 4px;
        }
      }

      .configuration-column {
        font-size: 12px;
        line-height: 1.5;
      }

      .pagination-container {
        margin-top: 20px;
        display: flex;
        justify-content: flex-end;
      }
    }
  }
</style>
