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

      <div class="card-extra" style="margin-bottom: 16px; text-align: right">
        <div class="el-space--horizontal" style="gap: 8px; display: inline-flex">
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
        <el-table-column label="主机名称" prop="name" min-width="200" sortable="custom">
          <template #default="{ row }">
            <div class="host-name-column">
              <el-tooltip :content="row.instance_id">
                <a @click="viewHostDetail(row.id)">{{ row.name }}</a>
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
              <div v-if="row.configuration && row.configuration.cpu_cores"
                >CPU: {{ row.configuration.cpu_cores }} 核</div
              >
              <div v-if="row.configuration && row.configuration.memory_size"
                >内存: {{ formatMemorySize(row.configuration.memory_size) }}</div
              >
              <div v-if="row.configuration && row.configuration.instance_type"
                >实例类型: {{ row.configuration.instance_type }}</div
              >
            </div>
          </template>
        </el-table-column>

        <el-table-column label="操作系统" min-width="150">
          <template #default="{ row }">
            <div class="os-column">
              <el-tooltip :content="row.os">
                <span>{{ getOsShortName(row.os) }}</span>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100" sortable="custom">
          <template #default="{ row }">
            <el-tag :type="getStatusColor(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="到期时间" width="150" sortable="custom">
          <template #default="{ row }">
            <div class="expired-at-column">
              <el-tag v-if="row.expired_at" :type="getExpiryColor(row.expired_at)">
                {{ formatExpiryTime(row.expired_at) }}
              </el-tag>
              <span v-else>--</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="创建时间" width="150" sortable="custom">
          <template #default="{ row }">
            <div class="time-info">
              {{ formatTime(row.created_at) }}
            </div>
          </template>
        </el-table-column>

        <el-table-column label="所属分组" width="150">
          <template #default="{ row }">
            <div class="group-column">
              <el-tooltip :content="row.group ? row.group.name : '未分组'">
                <span>{{ row.group ? row.group.name : '未分组' }}</span>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <div class="action-column">
              <div class="el-space--horizontal" style="gap: 4px; display: inline-flex">
                <el-tooltip content="SSH 终端">
                  <el-button link @click="openTerminal(row)">
                    <el-icon><Monitor /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="SFTP 文件管理">
                  <el-button link @click="openSftp(row)">
                    <el-icon><Folder /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-dropdown @command="(command) => handleOperationCommand(command, row)">
                  <el-button link>
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
                      <el-divider />
                      <el-dropdown-item command="restart">
                        <el-icon><Refresh /></el-icon> 重启主机
                      </el-dropdown-item>
                      <el-dropdown-item command="sync">
                        <el-icon><Refresh /></el-icon> 同步状态
                      </el-dropdown-item>
                      <el-divider />
                      <el-dropdown-item command="delete">
                        <el-icon><Delete /></el-icon> 删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

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
    Delete,
    Search
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

  // 主机组选项
  const hostGroupOptions = computed(() => {
    const options = []
    const processGroups = (groups, prefix = '') => {
      groups.forEach((group) => {
        const label = prefix ? `${prefix} / ${group.name}` : group.name
        options.push({ label, value: group.id })

        // 查找子组
        const children = hostGroupTree.value.filter((g) => g.parent_id === group.id)
        if (children.length > 0) {
          processGroups(children, label)
        }
      })
    }

    // 获取所有一级节点
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

  // 生命周期钩子
  onMounted(async () => {
    await Promise.all([hostStore.fetchHosts(searchParams), hostStore.fetchHostGroupTree()])
  })

  // 查询相关
  const handleSearch = () => {
    searchParams.page = 1 // 重置到第一页
    hostStore.fetchHosts(searchParams)
  }

  // 表格事件处理
  const onSelectionChange = (selection) => {
    selectedRowKeys.value = selection.map((item) => item.id)
  }

  const handleTableChange = ({ prop, order }) => {
    if (prop && order) {
      // 处理排序
      const sortOrder = order === 'ascending' ? 'asc' : 'desc'
      console.log(`排序列: ${prop}, 排序方向: ${sortOrder}`)
      // 这里可以发送带排序的API请求
    }
  }

  // 跳转到详情页
  const viewHostDetail = (id) => {
    router.push(`/cmdb/host/detail/${id}`)
  }

  // 添加主机相关
  const handleAddMenuClick = (command) => {
    if (command === 'manual') {
      manualHostModalVisible.value = true
    } else if (command === 'batch') {
      batchImportModalVisible.value = true
    } else if (command === 'sync') {
      ElMessageBox.confirm('确定要从云提供商同步主机信息吗？这可能需要一些时间。', '同步主机', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async () => {
          try {
            const loading = ElLoading.service({
              lock: true,
              text: '正在同步主机数据...',
              background: 'rgba(0, 0, 0, 0.7)'
            })
            await hostApi.syncHosts()
            ElMessage.success('主机同步成功')
            await hostStore.fetchHosts(searchParams)
            loading.close()
          } catch (error) {
            ElMessage.error('同步失败: ' + error.message)
          }
        })
        .catch(() => {
          // 用户取消
        })
    }
  }

  // 处理下拉菜单操作
  const handleOperationCommand = (command, row) => {
    switch (command) {
      case 'view':
        viewHostDetail(row.id)
        break
      case 'edit':
        handleEditHost(row)
        break
      case 'move':
        moveHostToGroup(row)
        break
      case 'restart':
        restartHost(row)
        break
      case 'sync':
        syncHostStatus(row.id)
        break
      case 'delete':
        confirmDeleteHost(row)
        break
    }
  }

  // 批量操作
  const handleBatchOperation = (command) => {
    if (selectedRowKeys.value.length === 0) {
      ElMessage.warning('请先选择要操作的主机')
      return
    }

    switch (command) {
      case 'move':
        batchMoveModalVisible.value = true
        break
      case 'tags':
        batchTagsVisible.value = true
        break
      case 'status':
        ElMessage.info('批量修改状态功能待实现')
        break
      case 'terminal':
        ElMessage.info('批量执行命令功能待实现')
        break
      case 'sftp':
        ElMessage.info('批量文件传输功能待实现')
        break
      case 'delete':
        confirmBatchDelete()
        break
    }
  }

  // 显示仪表盘
  const showDashboard = () => {
    router.push('/cmdb/host/dashboard')
  }

  // 编辑主机
  const handleEditHost = (host) => {
    currentHost.value = host
    isEditMode.value = true
    hostModalVisible.value = true
  }

  // 移动主机到分组
  const moveHostToGroup = (host) => {
    selectedRowKeys.value = [host.id]
    batchMoveModalVisible.value = true
  }

  // 重启主机
  const restartHost = (host) => {
    ElMessageBox.confirm(`确定要重启主机 ${host.name} 吗？`, '重启主机', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
      .then(() => {
        ElMessage.success('已发送重启命令')
      })
      .catch(() => {})
  }

  // 同步主机状态
  const syncHostStatus = async (id) => {
    try {
      const loading = ElLoading.service({
        lock: true,
        text: '正在同步主机状态...',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      await hostApi.syncHostStatus(id)
      ElMessage.success('同步成功')
      await hostStore.fetchHosts(searchParams)
      loading.close()
    } catch (error) {
      ElMessage.error('同步失败: ' + error.message)
    }
  }

  // 确认删除主机
  const confirmDeleteHost = (host) => {
    ElMessageBox.confirm(`确定要删除主机 ${host.name} 吗？此操作不可恢复！`, '删除主机', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'error'
    })
      .then(async () => {
        try {
          await hostApi.deleteHost(host.id)
          ElMessage.success('删除成功')
          await hostStore.fetchHosts(searchParams)
        } catch (error) {
          ElMessage.error('删除失败: ' + error.message)
        }
      })
      .catch(() => {})
  }

  // 确认批量删除
  const confirmBatchDelete = () => {
    ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRowKeys.value.length} 台主机吗？此操作不可恢复！`,
      '批量删除',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'error'
      }
    )
      .then(async () => {
        try {
          await hostStore.batchDeleteHosts(selectedRowKeys.value)
          ElMessage.success('批量删除成功')
          selectedRowKeys.value = []
          await hostStore.fetchHosts(searchParams)
        } catch (error) {
          ElMessage.error('批量删除失败: ' + error.message)
        }
      })
      .catch(() => {})
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
    hostStore.fetchHosts(searchParams)
  }

  // 重试方法
  const handleRetry = () => {
    hostStore.fetchHosts(searchParams)
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
      running: 'success',
      stopped: 'info',
      error: 'danger',
      expired: 'danger',
      starting: 'primary',
      stopping: 'info',
      rebooting: 'primary'
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

    if (diffDays < 0) return 'danger'
    if (diffDays < 7) return 'warning'
    if (diffDays < 30) return 'info'
    return 'success'
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
      .card-header {
        display: flex;
        align-items: center;
        height: 32px;
      }

      .host-name-column {
        display: flex;
        flex-direction: column;

        .el-tag {
          margin-top: 4px;
        }
      }

      .ip-address-column {
        .ip-label {
          color: var(--el-text-color-secondary);
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
  :deep(.el-card) {
    --el-card-bg-color: var(--el-bg-color);
  }

  html.dark {
    .host-management-page {
      .ip-address-column {
        .ip-label {
          color: var(--el-text-color-secondary);
        }
      }
    }
  }
</style>
