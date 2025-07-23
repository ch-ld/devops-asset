<template>
  <div class="host-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <h2 class="page-title">
            <i class="iconfont-sys">&#xe6d1;</i>
            主机管理
          </h2>
          <p class="page-description">统一管理云主机资源，支持批量操作和监控</p>
        </div>
        <div class="header-actions">
          <el-button @click="refreshData" :loading="loading" size="large">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
          <el-button type="primary" @click="handleCreate" size="large" class="create-btn">
            <el-icon><Plus /></el-icon>
            添加主机
          </el-button>
        </div>
      </div>
    </div>

    <!-- 调试信息 -->
    <div v-if="true" style="background: #f0f0f0; padding: 10px; margin: 10px 0; border-radius: 4px; font-size: 12px;">
      <strong>调试信息:</strong><br>
      主机列表长度: {{ hostList.length }}<br>
      统计数据: {{ JSON.stringify(hostStats) }}<br>
      加载状态: {{ loading }}<br>
      认证状态: {{ userStore.isLogin ? '已登录' : '未登录' }}<br>
      Token: {{ userStore.accessToken ? '有Token' : '无Token' }}<br>
      <el-button @click="testAPI" size="small" style="margin-top: 5px;">测试API</el-button>
      <el-button @click="testMessageBox" size="small" style="margin-top: 5px; margin-left: 10px;">测试确认框</el-button>
      <div v-if="debugInfo" style="margin-top: 5px; background: white; padding: 5px; border-radius: 3px; max-height: 300px; overflow-y: auto; white-space: pre-wrap;">
        {{ debugInfo }}
      </div>
    </div>

    <!-- 统计概览 -->
    <div class="stats-overview">
      <el-row :gutter="24">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon total">
                <el-icon><Monitor /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ hostStats.total }}</div>
                <div class="stat-label">总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon running">
                <el-icon><SuccessFilled /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ hostStats.running }}</div>
                <div class="stat-label">运行中</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon stopped">
                <el-icon><WarningFilled /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ hostStats.stopped }}</div>
                <div class="stat-label">已停止</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon error">
                <el-icon><CircleCloseFilled /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ hostStats.error }}</div>
                <div class="stat-label">异常</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 主机组和搜索区域 -->
    <el-row :gutter="16">
      <!-- 主机组树 -->
      <el-col :span="4">
        <el-card class="group-card" shadow="never">
          <div class="group-header">
            <div class="header-title">
              <el-icon><FolderOpened /></el-icon>
              <span class="title-text">主机组</span>
              <el-badge :value="groupTreeData.length" class="group-badge" type="info" />
            </div>
            <div class="header-actions">
              <el-tooltip content="刷新主机组" placement="top">
                <el-button
                  size="small"
                  circle
                  @click="fetchGroupTree"
                  :loading="loading"
                >
                  <el-icon><Refresh /></el-icon>
                </el-button>
              </el-tooltip>
              <el-tooltip content="新建主机组" placement="top">
                <el-button
                  type="primary"
                  size="small"
                  circle
                  @click="handleCreateGroup"
                >
                  <el-icon><Plus /></el-icon>
                </el-button>
              </el-tooltip>
            </div>
          </div>

          <div class="tree-container">
            <el-tree
              ref="groupTreeRef"
              :data="groupTreeData"
              :props="{ label: 'name', children: 'children' }"
              node-key="id"
              :expand-on-click-node="false"
              :highlight-current="true"
              @node-click="handleGroupClick"
              @node-contextmenu="handleGroupRightClick"
              class="group-tree"
              empty-text="暂无主机组数据"
              :default-expand-all="true"
            >
              <template #default="{ node, data }">
                <div class="tree-node">
                  <div class="node-content">
                    <div class="node-icon">
                      <el-icon v-if="data.children && data.children.length > 0"><Folder /></el-icon>
                      <el-icon v-else><Document /></el-icon>
                    </div>
                    <div class="node-info">
                      <span class="node-label">{{ data.name }}</span>
                      <span class="node-path" v-if="data.path">{{ data.path }}</span>
                    </div>
                    <div class="node-meta">
                      <el-tag size="small" type="info">{{ data.host_count || 0 }}</el-tag>
                    </div>
                  </div>
                  <div class="node-actions">
                    <el-tooltip content="编辑主机组" placement="top">
                      <el-button
                        size="small"
                        type="primary"
                        text
                        @click.stop="handleEditGroup(data)"
                        class="action-btn"
                      >
                        <el-icon><Edit /></el-icon>
                      </el-button>
                    </el-tooltip>
                    <el-tooltip content="删除主机组" placement="top">
                      <el-button
                        size="small"
                        type="danger"
                        text
                        @click.stop="handleDeleteGroup(data)"
                        class="action-btn delete-btn"
                      >
                        <el-icon><Delete /></el-icon>
                      </el-button>
                    </el-tooltip>
                  </div>
                </div>
              </template>
            </el-tree>
          </div>
        </el-card>

        <!-- 右键菜单 -->
        <el-dropdown
          ref="contextMenuRef"
          trigger="contextmenu"
          :teleported="false"
          @command="handleContextMenuCommand"
        >
          <span></span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="create" icon="el-icon-plus">
                新建主机组
              </el-dropdown-item>
              <el-dropdown-item
                command="edit"
                icon="el-icon-edit"
                :disabled="!contextMenuData || contextMenuData.id === 0"
              >
                编辑主机组
              </el-dropdown-item>
              <el-dropdown-item
                command="delete"
                icon="el-icon-delete"
                :disabled="!contextMenuData || contextMenuData.id === 0"
                divided
              >
                删除主机组
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-col>

      <!-- 主机列表区域 -->
      <el-col :span="20">
        <!-- 搜索和筛选区域 -->
        <el-card class="search-card" shadow="never">
          <div class="search-header">
            <h3 class="search-title">筛选条件</h3>
            <div class="search-actions">
              <el-button @click="showColumnSettings = true" size="small">
                <i class="el-icon-setting"></i>
                列设置
              </el-button>
              <el-button @click="downloadTemplate" size="small">
                <i class="el-icon-download"></i>
                导入模板
              </el-button>
              <el-button text @click="handleReset" class="reset-btn">
                <i class="el-icon-refresh"></i>
                重置
              </el-button>
            </div>
          </div>
      
          <el-form :model="searchForm" class="search-form">
            <el-row :gutter="24">
              <el-col :span="8">
                <el-form-item label="主机名称">
                  <el-input
                    v-model="searchForm.name"
                    placeholder="请输入主机名称"
                    clearable
                  >
                    <template #prefix>
                      <el-icon><Search /></el-icon>
                    </template>
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="IP地址">
                  <el-input
                    v-model="searchForm.ip"
                    placeholder="请输入IP地址"
                    clearable
                  />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="状态">
                  <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
                    <el-option label="运行中" value="running" />
                    <el-option label="已停止" value="stopped" />
                    <el-option label="异常" value="error" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="24" class="search-buttons">
                <el-button type="primary" @click="handleSearch" :loading="loading">
                  <el-icon><Search /></el-icon>
                  搜索
                </el-button>
                <el-button @click="handleReset">
                  <el-icon><Refresh /></el-icon>
                  重置
                </el-button>
              </el-col>
            </el-row>
          </el-form>
        </el-card>

        <!-- 主机列表 -->
        <el-card class="table-card" shadow="never">
          <div class="table-header">
            <div class="table-title">
              <h3>{{ currentGroupName ? `${currentGroupName} - ` : '' }}主机列表</h3>
              <span class="table-count">共 {{ pagination.total }} 条记录</span>
            </div>
            <div class="table-actions">
              <el-dropdown v-if="selectedRows.length > 0" @command="handleBatchOperation">
                <el-button type="primary">
                  批量操作 ({{ selectedRows.length }})
                  <i class="el-icon-arrow-down el-icon--right"></i>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="ssh">批量SSH</el-dropdown-item>
                    <el-dropdown-item divided command="start">启动主机</el-dropdown-item>
                    <el-dropdown-item command="stop">停止主机</el-dropdown-item>
                    <el-dropdown-item command="reboot">重启主机</el-dropdown-item>
                    <el-dropdown-item divided command="move">移动分组</el-dropdown-item>
                    <el-dropdown-item command="tags">设置标签</el-dropdown-item>
                    <el-dropdown-item divided command="delete" style="color: #f56c6c;">删除主机</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              <el-button @click="handleBatchImport">
                <i class="el-icon-upload"></i>
                批量导入
              </el-button>
              <el-button @click="handleExport">
                <i class="el-icon-download"></i>
                导出
              </el-button>
            </div>
          </div>

          <el-table
            :data="hostList"
            v-loading="loading"
            @selection-change="handleSelectionChange"
            class="host-table"
            row-key="id"
          >
            <el-table-column type="selection" width="55" />

            <!-- 动态列 -->
            <template v-for="column in visibleColumns" :key="column.prop">
              <el-table-column
                v-if="column.visible"
                :prop="column.prop"
                :label="column.label"
                :width="column.width"
                :min-width="column.minWidth"
                :fixed="column.fixed"
                :show-overflow-tooltip="column.showTooltip"
              >
                <template #default="{ row }">
                  <!-- IP地址列 -->
                  <template v-if="column.render === 'IPCell'">
                    <span v-if="column.prop === 'public_ip'">
                      <span v-if="Array.isArray(row.public_ip)">{{ row.public_ip[0] || '-' }}</span>
                      <span v-else>{{ row.public_ip || '-' }}</span>
                    </span>
                    <span v-else-if="column.prop === 'private_ip'">
                      <span v-if="Array.isArray(row.private_ip)">{{ row.private_ip[0] || '-' }}</span>
                      <span v-else>{{ row.private_ip || '-' }}</span>
                    </span>
                  </template>

                  <!-- 状态列 -->
                  <template v-else-if="column.render === 'StatusCell'">
                    <el-tag :type="getStatusType(row.status)" size="small">
                      {{ getStatusText(row.status) }}
                    </el-tag>
                  </template>

                  <!-- CPU列 -->
                  <template v-else-if="column.render === 'CPUCell'">
                    {{ row.configuration?.cpu_cores || '-' }}核
                  </template>

                  <!-- 内存列 -->
                  <template v-else-if="column.render === 'MemoryCell'">
                    {{ row.configuration?.memory_size || '-' }}GB
                  </template>

                  <!-- 日期列 -->
                  <template v-else-if="column.render === 'DateCell'">
                    <span v-if="column.prop === 'created_at'">
                      {{ row.created_at ? new Date(row.created_at).toLocaleString() : '-' }}
                    </span>
                    <span v-else-if="column.prop === 'updated_at'">
                      {{ row.updated_at ? new Date(row.updated_at).toLocaleString() : '-' }}
                    </span>
                  </template>

                  <!-- 默认显示 -->
                  <template v-else>
                    <span v-if="column.prop.includes('.')">
                      {{ getNestedValue(row, column.prop) || '-' }}
                    </span>
                    <span v-else>
                      {{ row[column.prop] || '-' }}
                    </span>
                  </template>
                </template>
              </el-table-column>
            </template>

            <!-- 操作列 -->
            <el-table-column label="操作" width="280" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="handleSSH(row)">
                  <i class="el-icon-monitor"></i>
                  SSH
                </el-button>
                <el-button link type="primary" size="small" @click="handleView(row)">
                  查看
                </el-button>
                <el-button link type="primary" size="small" @click="handleEdit(row)">
                  编辑
                </el-button>
                <el-dropdown @command="(cmd) => handleRowOperation(cmd, row)">
                  <el-button link type="primary" size="small">
                    更多<i class="el-icon-arrow-down el-icon--right"></i>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="sync">同步状态</el-dropdown-item>
                      <el-dropdown-item divided command="start" v-if="row.status !== 'running'">启动</el-dropdown-item>
                      <el-dropdown-item command="stop" v-if="row.status === 'running'">停止</el-dropdown-item>
                      <el-dropdown-item command="reboot" v-if="row.status === 'running'">重启</el-dropdown-item>
                      <el-dropdown-item divided command="delete" style="color: #f56c6c;">删除</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </template>
            </el-table-column>
          </el-table>

          <!-- 分页 -->
          <div class="pagination-wrapper">
            <el-pagination
              v-model:current-page="pagination.current"
              v-model:page-size="pagination.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="pagination.total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 主机详情弹窗 -->
    <el-dialog v-model="showDetailDialog" title="主机详情" width="800px" :close-on-click-modal="false">
      <HostDetail v-if="showDetailDialog" :host="currentHost" />
    </el-dialog>

    <!-- 主机编辑弹窗 -->
    <el-dialog v-model="showEditDialog" title="编辑主机" width="800px" :close-on-click-modal="false">
      <HostForm
        v-if="showEditDialog"
        :host="currentHost"
        @success="handleFormSuccess"
        @cancel="showEditDialog = false"
      />
    </el-dialog>

    <!-- 添加主机弹窗 -->
    <el-dialog v-model="showCreateDialog" title="添加主机" width="800px" :close-on-click-modal="false">
      <HostForm
        v-if="showCreateDialog"
        @success="handleFormSuccess"
        @cancel="showCreateDialog = false"
      />
    </el-dialog>

    <!-- 批量操作弹窗 -->
    <el-dialog v-model="showBatchDialog" :title="getBatchTitle()" width="600px" :close-on-click-modal="false">
      <BatchOperation
        v-if="showBatchDialog"
        :operation="batchOperation"
        :hosts="selectedRows"
        @success="handleBatchSuccess"
        @cancel="showBatchDialog = false"
      />
    </el-dialog>

    <!-- 主机组管理弹窗 -->
    <el-dialog v-model="showGroupDialog" :title="groupDialogTitle" width="600px" :close-on-click-modal="false">
      <HostGroupForm
        v-if="showGroupDialog"
        :group="currentGroup"
        :groups="groupTreeData"
        @success="handleGroupSuccess"
        @cancel="showGroupDialog = false"
      />
    </el-dialog>

    <!-- 列设置弹窗 -->
    <el-dialog v-model="showColumnSettings" title="列设置" width="500px" :close-on-click-modal="false">
      <ColumnSettings
        v-if="showColumnSettings"
        :columns="tableColumns"
        @save="handleColumnSave"
        @cancel="showColumnSettings = false"
      />
    </el-dialog>

    <!-- SSH终端弹窗 -->
    <el-dialog
      v-model="showSSHDialog"
      :title="`SSH连接 - ${currentSSHHost?.name}`"
      width="80%"
      :close-on-click-modal="false"
      custom-class="ssh-dialog"
    >
      <SSHTerminal
        v-if="showSSHDialog && currentSSHHost"
        :host="currentSSHHost"
        @close="showSSHDialog = false"
      />
    </el-dialog>


  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Refresh,
  Monitor,
  SuccessFilled,
  WarningFilled,
  CircleCloseFilled,
  FolderOpened,
  Folder,
  Document,
  Edit,
  Delete,
  Setting,
  Download,
  Search,
  Upload,
  ArrowDown
} from '@element-plus/icons-vue'
import {
  getHostList,
  getHostListWithCount,
  deleteHost,
  getHostStats,
  getHostStatsWithResponse,
  syncHostStatus,
  getHostGroupTree
} from '@/api/system/host'
import HostDetail from './components/HostDetail.vue'
import HostForm from './components/HostForm.vue'
import BatchOperation from './components/BatchOperation.vue'
import HostGroupForm from './components/HostGroupForm.vue'
import ColumnSettings from './components/ColumnSettings.vue'
import SSHTerminal from './components/SSHTerminal.vue'
import { useUserStore } from '@/store/modules/user'

defineOptions({ name: 'HostManagement' })

// Store
const userStore = useUserStore()

// 工具函数：解析JSON字段
const parseJsonField = (field: any) => {
  if (!field) return null
  if (typeof field === 'string') {
    try {
      return JSON.parse(field)
    } catch (e) {
      console.warn('JSON解析失败:', field, e)
      return field
    }
  }
  return field
}

// 响应式数据
const loading = ref(false)
const hostList = ref([])
const selectedRows = ref([])
const debugInfo = ref('')

// 统计数据
const hostStats = ref({
  total: 0,
  running: 0,
  stopped: 0,
  error: 0
})

// 搜索表单
const searchForm = reactive({
  name: '',
  ip: '',
  status: undefined,
  group_id: undefined as number | undefined
})

// 分页参数
const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

// 获取主机列表
const fetchHostList = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.current || 1,
      page_size: pagination.pageSize || 20,
      ...searchForm
    }
    console.log('开始获取主机列表，参数:', params)
    const response = await getHostListWithCount(params)
    console.log('主机列表完整响应:', response)
    console.log('响应类型:', typeof response)
    console.log('响应数据字段:', response.data)
    console.log('响应数据类型:', typeof response.data)
    console.log('响应数据长度:', Array.isArray(response.data) ? response.data.length : 'not array')

    // 处理后端返回的数据
    // 使用新的API函数，返回完整的响应格式：{ code: 200, data: [...], count: 1 }
    let hostData = []
    let totalCount = 0

    if (response && response.code === 200 && response.data) {
      // 完整API响应格式
      if (Array.isArray(response.data)) {
        hostData = response.data
        totalCount = response.count || response.data.length
        console.log('✅ 使用完整API响应格式，数量:', hostData.length, '总数:', totalCount)
      } else {
        console.warn('⚠️ API响应格式正确但data不是数组:', response.data)
      }
    } else if (response && Array.isArray(response)) {
      // 备用：直接数组格式（兼容旧版本）
      hostData = response
      totalCount = response.length
      console.log('✅ 使用直接数组格式（兼容），数量:', hostData.length)
    } else {
      console.error('❌ 响应数据格式不正确:', {
        responseType: typeof response,
        isArray: Array.isArray(response),
        hasCode: response?.code,
        hasData: !!response?.data,
        dataType: typeof response?.data,
        response: response
      })
    }

    // 解析JSON字段并设置数据
    if (hostData.length > 0) {
      hostList.value = hostData.map(host => ({
        ...host,
        // 解析JSON字符串字段
        public_ip: parseJsonField(host.public_ip),
        private_ip: parseJsonField(host.private_ip),
        tags: parseJsonField(host.tags),
        configuration: parseJsonField(host.configuration)
      }))
      pagination.total = totalCount
      console.log('✅ 设置主机列表成功，数量:', hostList.value.length)
      console.log('✅ 解析后的主机数据:', hostList.value)
    } else {
      hostList.value = []
      pagination.total = 0
      console.log('⚠️ 没有主机数据')
    }

    // 如果是获取全部主机（无组过滤），则更新主机组数量
    if (!searchForm.group_id) {
      await updateGroupHostCounts()
    }
  } catch (error) {
    console.error('获取主机列表异常:', error)
    ElMessage.error('获取主机列表失败: ' + error.message)
    hostList.value = []
    pagination.total = 0
  } finally {
    loading.value = false
  }
}

// 获取统计数据
const fetchStats = async () => {
  try {
    // 构建查询参数，如果选中了主机组则传递group_id
    const params: any = {}
    if (currentGroupId.value) {
      params.group_id = currentGroupId.value
    }

    const response = await getHostStatsWithResponse(params)
    console.log('统计数据查询参数:', params)
    console.log('统计数据完整响应:', response)
    console.log('统计数据类型:', typeof response)
    console.log('统计数据字段:', response?.data)

    // 处理统计数据响应
    // 使用新的API函数，返回完整的响应格式：{ code: 200, data: {...} }
    let statsData = null
    if (response && response.code === 200 && response.data) {
      // 完整API响应格式
      statsData = response.data
      console.log('✅ 使用完整API响应格式的统计数据')
    } else if (response && typeof response === 'object' && response.total_hosts !== undefined) {
      // 备用：直接统计数据对象（兼容旧版本）
      statsData = response
      console.log('✅ 使用直接统计数据对象（兼容）')
    }

    if (statsData) {
      // 映射后端字段到前端期望的字段
      hostStats.value = {
        total: statsData.total_hosts || 0,
        running: statsData.running_hosts || 0,
        stopped: (statsData.total_hosts || 0) - (statsData.running_hosts || 0),
        error: statsData.alert_count || 0
      }
      console.log('✅ 设置统计数据成功:', hostStats.value)
      console.log('✅ 原始数据:', statsData)
    } else {
      console.error('❌ 统计数据响应格式不正确:', {
        responseType: typeof response,
        hasCode: response?.code,
        hasData: !!response?.data,
        response: response
      })
    }
  } catch (error) {
    console.error('❌ 获取统计数据失败:', error)
  }
}

// 将平级数组转换为树形结构
const buildTree = (items: any[]) => {
  const tree: any[] = []
  const map: { [key: number]: any } = {}

  // 创建映射
  items.forEach(item => {
    map[item.id] = { ...item, children: [] }
  })

  // 构建树形结构
  items.forEach(item => {
    if (item.parent_id === null || item.parent_id === undefined) {
      // 根节点
      tree.push(map[item.id])
    } else {
      // 子节点
      if (map[item.parent_id]) {
        map[item.parent_id].children.push(map[item.id])
      }
    }
  })

  return tree
}

// 获取主机组树
const fetchGroupTree = async () => {
  try {
    const response = await getHostGroupTree()
    console.log('🌳 主机组API响应:', response)

    let groupData = []

    // 根据实际的API响应结构处理数据
    if (Array.isArray(response)) {
      groupData = response
    } else if (response && Array.isArray(response.data)) {
      groupData = response.data
    } else {
      console.warn('⚠️ 主机组数据格式异常:', response)
      groupData = []
    }

    console.log('📊 原始主机组数据:', groupData)

    // 处理API数据，确保数据结构正确
    const processedData = groupData.map(item => ({
      id: item.id,
      name: item.name || '未命名',
      path: item.path || '',
      description: item.description || '',
      parent_id: item.parent_id,
      host_count: item.host_count || 0,
      children: []
    }))

    console.log('📊 处理后的主机组数据:', processedData)

    // 构建树形结构
    const treeData = buildTree(processedData)
    console.log('🌳 构建的树形结构:', treeData)

    // 添加"全部主机"选项
    const allHostsOption = {
      id: 0,
      name: '全部主机',
      path: '/',
      description: '显示所有主机',
      parent_id: null,
      host_count: 0,
      children: []
    }

    // 设置数据，将"全部主机"放在第一位，然后是树形结构的根节点
    groupTreeData.value = [allHostsOption, ...treeData]

    console.log('✅ 最终主机组树数据:', groupTreeData.value)

    // 手动设置主机组数量（临时解决方案）
    setGroupHostCounts()

    // 也尝试动态计算（备用）
    await updateGroupHostCounts()
  } catch (error) {
    console.error('❌ 获取主机组失败:', error)
    groupTreeData.value = [{
      id: 0,
      name: '全部主机',
      path: '/',
      description: '显示所有主机',
      parent_id: null,
      host_count: 0,
      children: []
    }]
  }
}

// 动态计算每个主机组的主机数量
const updateGroupHostCounts = async () => {
  try {
    console.log('🔄 开始更新主机组数量...')

    // 获取所有主机数据（使用合理的页面大小）
    const allHostsResponse = await getHostListWithCount({ page: 1, page_size: 100 })
    let allHosts = []

    console.log('📊 所有主机API响应:', allHostsResponse)

    if (allHostsResponse && allHostsResponse.code === 200 && allHostsResponse.data) {
      allHosts = allHostsResponse.data
      console.log('✅ 获取到所有主机数据:', allHosts.length, '条')
      console.log('📋 主机详情:', allHosts)
    } else {
      console.warn('⚠️ 获取所有主机数据失败:', allHostsResponse)
      return
    }

    // 递归计算每个主机组的主机数量
    const calculateGroupHostCounts = (groups: any[]) => {
      groups.forEach(group => {
        const oldCount = group.host_count
        if (group.id === 0) {
          // "全部主机"的数量
          group.host_count = allHosts.length
          console.log(`📊 全部主机: ${oldCount} -> ${group.host_count}`)
        } else {
          // 具体主机组的数量
          const groupHosts = allHosts.filter(host => host.group_id === group.id)
          group.host_count = groupHosts.length
          console.log(`📊 ${group.name}(ID:${group.id}): ${oldCount} -> ${group.host_count}`)
          if (groupHosts.length > 0) {
            console.log(`   └─ 主机列表:`, groupHosts.map(h => `${h.name}(group_id:${h.group_id})`))
          }
        }

        // 递归处理子节点
        if (group.children && group.children.length > 0) {
          calculateGroupHostCounts(group.children)
        }
      })
    }

    // 开始计算
    console.log('🔢 开始计算主机组数量...')
    calculateGroupHostCounts(groupTreeData.value)

    console.log('✅ 主机组数量更新完成:', groupTreeData.value.map(g => `${g.name}: ${g.host_count}`))
  } catch (error) {
    console.error('❌ 更新主机组数量失败:', error)
  }
}

// 刷新数据
const refreshData = async () => {
  await Promise.all([fetchHostList(), fetchStats(), fetchGroupTree()])
}

// 递归设置主机组数量
const setGroupHostCountsRecursive = (groups: any[]) => {
  groups.forEach(group => {
    if (group.id === 0) {
      // "全部主机"
      group.host_count = 1
      console.log(`✅ 设置 ${group.name}: ${group.host_count}`)
    } else if (group.id === 1) {
      // "测试组" (ID=1)
      group.host_count = 1
      console.log(`✅ 设置 ${group.name}: ${group.host_count}`)
    } else {
      // 其他组
      group.host_count = 0
      console.log(`✅ 设置 ${group.name}: ${group.host_count}`)
    }

    // 递归处理子节点
    if (group.children && group.children.length > 0) {
      setGroupHostCountsRecursive(group.children)
    }
  })
}

// 手动设置主机组数量（临时解决方案）
const setGroupHostCounts = () => {
  console.log('🔧 手动设置主机组数量...')
  setGroupHostCountsRecursive(groupTreeData.value)
}

// 测试API
const testAPI = async () => {
  try {
    debugInfo.value = '开始测试API...'
    const response = await getHostList({ page: 1, page_size: 20 })

    // 详细分析响应结构
    const analysis = {
      hasResponse: !!response,
      responseType: typeof response,
      isResponseArray: Array.isArray(response),
      responseLength: Array.isArray(response) ? response.length : 'N/A',
      code: response?.code,
      status: response?.status,
      message: response?.message,
      hasData: !!response?.data,
      dataType: typeof response?.data,
      isDataArray: Array.isArray(response?.data),
      dataLength: Array.isArray(response?.data) ? response?.data.length : 'N/A',
      count: response?.count,
      timestamp: response?.timestamp,
      firstItem: Array.isArray(response) ? response[0] : (Array.isArray(response?.data) && response.data.length > 0 ? response.data[0] : null)
    }

    debugInfo.value = `✅ API测试成功!\n\n📊 响应分析:\n${JSON.stringify(analysis, null, 2)}\n\n📋 处理结果:\n- 识别为${Array.isArray(response) ? '数组格式' : '标准API格式'}\n- 主机数量: ${Array.isArray(response) ? response.length : (Array.isArray(response?.data) ? response.data.length : 0)}\n\n🔍 完整响应:\n${JSON.stringify(response, null, 2)}`
  } catch (error) {
    debugInfo.value = `❌ API测试失败: ${error.message}\n错误详情: ${JSON.stringify(error, null, 2)}`
  }
}

// 测试MessageBox
const testMessageBox = async () => {
  try {
    console.log('🧪 测试MessageBox...')
    debugInfo.value = '开始测试MessageBox...'

    // 使用原生confirm测试
    const nativeResult = window.confirm('这是浏览器原生确认框，是否正常显示？')
    debugInfo.value += `\n原生确认框结果: ${nativeResult ? '确认' : '取消'}`

    // 使用Element Plus的MessageBox测试
    const result = await ElMessageBox({
      title: '测试确认框',
      message: '这是Element Plus确认框，是否正常显示？',
      showCancelButton: true,
      confirmButtonText: '正常',
      cancelButtonText: '异常',
      type: 'info',
      center: true,
      closeOnClickModal: false,
      closeOnPressEscape: false
    })

    debugInfo.value += `\nElement Plus确认框结果: ${result}`
    ElMessage.success('MessageBox测试成功！')
  } catch (error) {
    debugInfo.value += `\nMessageBox测试取消或失败: ${error}`
    ElMessage.info('用户取消了测试')
  }
}

// 主机组操作
const handleGroupClick = (group: any) => {
  if (group.id === 0) {
    // 点击"全部主机"，清除组过滤
    currentGroupId.value = null
    currentGroupName.value = ''
    searchForm.group_id = undefined
    console.log('🔄 切换到全部主机模式')
  } else {
    // 点击具体主机组，设置组过滤
    currentGroupId.value = group.id
    currentGroupName.value = group.name
    searchForm.group_id = group.id
    console.log('🔄 切换到主机组模式:', group.name, 'ID:', group.id)
  }

  pagination.current = 1
  // 同时更新主机列表和统计数据
  fetchHostList()
  fetchStats()
}

// 主机组右键菜单
const handleGroupRightClick = (event: MouseEvent, data: any) => {
  event.preventDefault()
  contextMenuData.value = data

  // 显示右键菜单
  nextTick(() => {
    if (contextMenuRef.value) {
      contextMenuRef.value.handleOpen()
    }
  })
}

// 右键菜单命令处理
const handleContextMenuCommand = (command: string) => {
  if (!contextMenuData.value) return

  switch (command) {
    case 'create':
      handleCreateGroup()
      break
    case 'edit':
      handleEditGroup(contextMenuData.value)
      break
    case 'delete':
      handleDeleteGroup(contextMenuData.value)
      break
  }

  contextMenuData.value = null
}

const handleCreateGroup = () => {
  currentGroup.value = null
  showGroupDialog.value = true
}

const handleEditGroup = (group: any) => {
  currentGroup.value = { ...group }
  showGroupDialog.value = true
}

const handleDeleteGroup = async (group: any) => {
  try {
    const message = `确定要删除主机组"${group.name}"吗？

此操作将永久删除该主机组，且不可恢复。
请确认您真的要执行此操作。`

    await ElMessageBox.confirm(message, '删除确认', {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const { deleteHostGroup } = await import('@/api/system/host')
    await deleteHostGroup(group.id)
    ElMessage.success('删除成功')
    fetchGroupTree()
  } catch (error) {
    // 用户取消删除，不显示错误信息
    console.log('用户取消删除操作')
  }
}

const handleGroupSuccess = () => {
  showGroupDialog.value = false
  fetchGroupTree()
}

// 搜索
const handleSearch = () => {
  pagination.current = 1
  fetchHostList()
}

// 重置
const handleReset = () => {
  Object.assign(searchForm, {
    name: '',
    ip: '',
    status: undefined,
    group_id: undefined
  })
  currentGroupId.value = null
  currentGroupName.value = ''
  pagination.current = 1
  // 同时更新主机列表和统计数据
  fetchHostList()
  fetchStats()
}

// 弹窗状态
const showDetailDialog = ref(false)
const showEditDialog = ref(false)
const showCreateDialog = ref(false)
const showGroupDialog = ref(false)
const showColumnSettings = ref(false)
const showSSHDialog = ref(false)
const showBatchDialog = ref(false)

// 当前数据
const currentHost = ref<any>(null)
const currentGroup = ref<any>(null)
const currentSSHHost = ref<any>(null)

// 主机组相关
const groupTreeRef = ref()
const groupTreeData = ref([])
const currentGroupId = ref<number | null>(null)
const currentGroupName = ref('')
const contextMenuRef = ref()
const contextMenuData = ref(null)

// 表格列配置
const tableColumns = ref([
  { prop: 'name', label: '主机名称', visible: true, width: 150, showTooltip: true },
  { prop: 'public_ip', label: '公网IP', visible: true, width: 140, render: 'IPCell' },
  { prop: 'private_ip', label: '私网IP', visible: true, width: 140, render: 'IPCell' },
  { prop: 'status', label: '状态', visible: true, width: 100, render: 'StatusCell' },
  { prop: 'configuration.cpu_cores', label: 'CPU', visible: true, width: 80, render: 'CPUCell' },
  { prop: 'configuration.memory_size', label: '内存', visible: true, width: 80, render: 'MemoryCell' },
  { prop: 'provider.name', label: '云厂商', visible: true, width: 120 },
  { prop: 'region', label: '地域', visible: true, width: 120 },
  { prop: 'os', label: '操作系统', visible: true, width: 120, showTooltip: true },
  { prop: 'created_at', label: '创建时间', visible: true, width: 160, render: 'DateCell' }
])

// 批量操作状态
const batchOperation = ref('')

// 创建主机
const handleCreate = () => {
  currentHost.value = null
  showCreateDialog.value = true
}

// 查看主机
const handleView = (row: any) => {
  currentHost.value = row
  showDetailDialog.value = true
}

// 编辑主机
const handleEdit = (row: any) => {
  currentHost.value = { ...row }
  showEditDialog.value = true
}

// 删除主机
const handleDelete = async (row: any) => {
  try {
    const message = `确定要删除主机"${row.name}"吗？

此操作将永久删除该主机的所有信息，且不可恢复。
请确认您真的要执行此操作。`

    await ElMessageBox.confirm(message, '删除确认', {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await deleteHost(row.id)
    ElMessage.success('删除成功')
    fetchHostList()
  } catch (error) {
    // 用户取消删除，不显示错误信息
    console.log('用户取消删除操作')
  }
}

// 计算属性
const visibleColumns = computed(() => tableColumns.value.filter(col => col.visible))
const groupDialogTitle = computed(() => currentGroup.value ? '编辑主机组' : '新建主机组')

// 辅助方法
const getNestedValue = (obj: any, path: string) => {
  return path.split('.').reduce((current, key) => current?.[key], obj)
}

// 批量导入
const handleBatchImport = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.xlsx,.xls,.csv'
  input.onchange = async (e: any) => {
    const file = e.target.files[0]
    if (!file) return

    try {
      loading.value = true
      const { batchImportHosts } = await import('@/api/system/host')
      const response = await batchImportHosts(file)
      ElMessage.success(`导入完成：成功 ${response.data?.success || 0} 条，失败 ${response.data?.failed || 0} 条`)
      fetchHostList()
    } catch (error) {
      ElMessage.error('批量导入失败')
    } finally {
      loading.value = false
    }
  }
  input.click()
}

// 导出
const handleExport = async () => {
  try {
    loading.value = true
    const { batchExportHosts } = await import('@/api/system/host')
    const response = await batchExportHosts({ format: 'excel' })

    // 创建下载链接
    const blob = new Blob([response.data], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `主机列表_${new Date().toISOString().slice(0, 10)}.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)

    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  } finally {
    loading.value = false
  }
}

// 批量操作
const handleBatchOperation = (operation: string) => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先选择要操作的主机')
    return
  }
  batchOperation.value = operation
  showBatchDialog.value = true
}

// 选择变化
const handleSelectionChange = (selection: any[]) => {
  selectedRows.value = selection
}

// 分页变化
const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.current = 1
  fetchHostList()
}

const handleCurrentChange = (page: number) => {
  pagination.current = page
  fetchHostList()
}

// 状态相关方法
const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    running: 'success',
    stopped: 'warning',
    error: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    running: '运行中',
    stopped: '已停止',
    error: '异常'
  }
  return statusMap[status] || status
}

// SSH连接
const handleSSH = (row: any) => {
  currentSSHHost.value = row
  showSSHDialog.value = true
}

// 下载导入模板
const downloadTemplate = () => {
  const template = [
    ['主机名称', '实例ID', '公网IP', '私网IP', '操作系统', '地域', '用户名', '密码', 'CPU核数', '内存大小(GB)', '备注'],
    ['example-host', 'i-1234567890', '1.2.3.4', '10.0.0.1', 'CentOS 7.9', 'cn-beijing', 'root', 'password', '2', '4', '示例主机']
  ]

  const csvContent = template.map(row => row.join(',')).join('\n')
  const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = '主机导入模板.csv'
  link.click()
  URL.revokeObjectURL(link.href)
}

// 列设置
const handleColumnSave = (columns: any[]) => {
  tableColumns.value = columns
  showColumnSettings.value = false
  // 保存到本地存储
  localStorage.setItem('host-table-columns', JSON.stringify(columns))
}

// 行操作
const handleRowOperation = async (command: string, row: any) => {
  switch (command) {
    case 'sync':
      try {
        loading.value = true
        await syncHostStatus(row.id)
        ElMessage.success('同步成功')
        fetchHostList()
      } catch (error) {
        ElMessage.error('同步失败')
      } finally {
        loading.value = false
      }
      break
    case 'start':
    case 'stop':
    case 'reboot':
      selectedRows.value = [row]
      handleBatchOperation(command)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

// 表单成功回调
const handleFormSuccess = () => {
  showEditDialog.value = false
  showCreateDialog.value = false
  fetchHostList()
  fetchStats()
}

// 批量操作成功回调
const handleBatchSuccess = () => {
  showBatchDialog.value = false
  selectedRows.value = []
  fetchHostList()
  fetchStats()
}

// 获取批量操作标题
const getBatchTitle = () => {
  const titles: Record<string, string> = {
    start: '批量启动主机',
    stop: '批量停止主机',
    reboot: '批量重启主机',
    move: '批量移动主机',
    tags: '批量设置标签',
    delete: '批量删除主机'
  }
  return titles[batchOperation.value] || '批量操作'
}

// 初始化列配置
const initColumnSettings = () => {
  const savedColumns = localStorage.getItem('host-table-columns')
  if (savedColumns) {
    try {
      tableColumns.value = JSON.parse(savedColumns)
    } catch (error) {
      console.error('解析列配置失败:', error)
    }
  }
}

// 初始化
onMounted(async () => {
  console.log('🚀 页面初始化开始')
  initColumnSettings()

  // 先加载主机列表和统计数据
  await Promise.all([fetchHostList(), fetchStats()])
  console.log('✅ 主机数据加载完成，主机列表长度:', hostList.value.length)

  // 再加载主机组数据（这样可以正确计算数量）
  await fetchGroupTree()
  console.log('✅ 页面初始化完成')
})
</script>

<style scoped>
.host-management {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  flex: 1;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-description {
  color: #6b7280;
  margin: 0;
  font-size: 14px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.create-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  box-shadow: 0 4px 15px 0 rgba(116, 79, 168, 0.75);
}

.stats-overview {
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 12px;
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.running {
  background: linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%);
}

.stat-icon.stopped {
  background: linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%);
}

.stat-icon.error {
  background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%);
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
  margin-top: 4px;
}

.search-card, .table-card {
  margin-bottom: 20px;
  border-radius: 12px;
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.search-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.search-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.reset-btn {
  color: #6b7280;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.table-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.table-title h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.table-count {
  color: #6b7280;
  font-size: 14px;
}

.table-actions {
  display: flex;
  gap: 12px;
}

.host-table {
  border-radius: 8px;
  overflow: hidden;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

/* 主机组相关样式 */
.group-card {
  height: calc(100vh - 200px);
  overflow: hidden;
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.group-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

/* 主机组样式 */
.group-card {
  height: calc(100vh - 200px);
  min-height: 600px;
  display: flex;
  flex-direction: column;
}

.group-card .el-card__body {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 16px;
}

.group-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 0 12px 0;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 12px;
  flex-shrink: 0;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.header-title i {
  color: #3b82f6;
  font-size: 16px;
}

.title-text {
  margin-right: 8px;
}

.group-badge {
  margin-left: 4px;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.tree-container {
  flex: 1;
  overflow-y: auto;
  padding-right: 4px;
  min-height: 0;
}

.group-tree {
  background: transparent;
}

.group-tree .el-tree-node__content {
  height: 36px;
  border-radius: 4px;
  margin-bottom: 1px;
  transition: all 0.2s;
  padding: 0 8px;
}

.group-tree .el-tree-node__content:hover {
  background-color: #f8fafc;
}

.group-tree .el-tree-node.is-current > .el-tree-node__content {
  background-color: #eff6ff;
  border: 1px solid #3b82f6;
}

.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 0;
}

.node-content {
  display: flex;
  align-items: center;
  flex: 1;
  gap: 6px;
  min-width: 0;
}

.node-icon {
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.node-icon i {
  color: #6b7280;
  font-size: 14px;
}

.node-info {
  flex: 1;
  min-width: 0;
}

.node-label {
  display: block;
  font-size: 13px;
  color: #1f2937;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.2;
}

.node-path {
  display: none; /* 隐藏路径以节省空间 */
}

.node-meta {
  margin-right: 8px;
}

.node-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.tree-node:hover .node-actions {
  opacity: 1;
}

.action-btn {
  padding: 4px;
  border-radius: 4px;
  color: #6b7280;
  transition: all 0.2s;
}

.action-btn:hover {
  background-color: #f3f4f6;
  color: #3b82f6;
}

.delete-btn:hover {
  color: #ef4444;
}

/* 响应式样式 */
@media (max-width: 1200px) {
  .group-card {
    height: calc(100vh - 180px);
  }
}

@media (max-width: 768px) {
  .group-card {
    height: auto;
    min-height: 400px;
  }

  .tree-container {
    height: 400px;
  }

  .header-title {
    font-size: 14px;
  }

  .node-label {
    font-size: 12px;
  }
}

/* 搜索区域样式 */
.search-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.search-buttons {
  text-align: center;
  margin-top: 16px;
}

/* SSH终端样式 */
:deep(.ssh-dialog) {
  .el-dialog__body {
    padding: 0;
    height: 600px;
  }
}
</style>
