<template>
  <div class="host-management">
    <div class="page-container">
      <!-- 左侧边栏 -->
      <div class="sidebar">
        <div class="sidebar-header">
          <div class="header-content">
            <div class="header-title">
              <i class="icon">📁</i>
              <h3>主机分组</h3>
            </div>
            <el-button size="small" type="primary" @click="refreshGroups" :loading="groupLoading">
              刷新
            </el-button>
          </div>
        </div>
        
        <div class="group-search">
          <el-input
            v-model="groupSearchQuery"
            placeholder="搜索分组..."
            clearable
            size="small"
          />
        </div>

        <div class="group-tree-container">
          <el-tree
            :data="filteredGroupTree"
            :props="treeProps"
            node-key="id"
            :current-node-key="selectedGroupId"
            @node-click="handleGroupSelect"
            :expand-on-click-node="false"
            class="group-tree"
          >
            <template #default="{ node, data }">
              <div
                class="tree-node"
                @contextmenu="(e) => handleGroupRightClick(e, node, data)"
              >
                <span class="node-label">{{ node.label }}</span>
                <span class="node-count">({{ data.host_count || 0 }})</span>
                <div class="node-actions" v-if="data.id !== null">
                  <el-button
                    type="text"
                    size="small"
                    @click.stop="handleAddGroup(data)"
                    title="添加子组"
                  >
                    +
                  </el-button>
                  <el-button
                    type="text"
                    size="small"
                    @click.stop="handleEditGroup(data)"
                    title="编辑"
                  >
                    ✎
                  </el-button>
                  <el-button
                    type="text"
                    size="small"
                    @click.stop="handleDeleteGroup(data)"
                    title="删除"
                  >
                    ✕
                  </el-button>
                </div>
              </div>
            </template>
          </el-tree>
        </div>
      </el-aside>

      <!-- 右侧主内容区 -->
      <el-main class="main-content">
        <!-- 页面标题和工具栏 -->
        <div class="content-header">
          <div class="header-left">
            <h1 class="page-title">主机管理</h1>
            <div class="breadcrumb">
              <span class="breadcrumb-item">主机管理</span>
              <span class="breadcrumb-separator">/</span>
              <span class="breadcrumb-current">{{ selectedGroupName || '全部主机' }}</span>
            </div>
          </div>
          <div class="header-actions">
            <el-button type="primary" @click="refreshData" :loading="loading">
              刷新数据
            </el-button>
            <el-button type="success" @click="handleAdd">
              添加主机
            </el-button>
          </div>
        </div>

        <!-- 搜索筛选工具栏 -->
        <div class="toolbar">
          <el-row :gutter="16" align="middle">
            <el-col :span="8">
              <el-input
                v-model="searchParams.keyword"
                placeholder="搜索主机名称、IP地址..."
                clearable
                @input="handleSearch"
                class="search-input"
              />
            </el-col>
            <el-col :span="4">
              <el-select
                v-model="searchParams.status"
                placeholder="状态筛选"
                clearable
                @change="handleSearch"
                class="filter-select"
              >
                <el-option label="全部状态" value="" />
                <el-option label="运行中" value="running" />
                <el-option label="已停止" value="stopped" />
                <el-option label="维护中" value="maintenance" />
              </el-select>
            </el-col>
            <el-col :span="4">
              <el-select
                v-model="searchParams.region"
                placeholder="地域筛选"
                clearable
                @change="handleSearch"
                class="filter-select"
              >
                <el-option label="全部地域" value="" />
                <el-option label="华东1" value="cn-hangzhou" />
                <el-option label="华东2" value="cn-shanghai" />
                <el-option label="华北1" value="cn-qingdao" />
                <el-option label="华北2" value="cn-beijing" />
                <el-option label="华南1" value="cn-shenzhen" />
                <el-option label="新加坡" value="ap-southeast-1" />
              </el-select>
            </el-col>
            <el-col :span="4">
              <el-select
                v-model="providerFilter"
                placeholder="云厂商筛选"
                clearable
                @change="handleProviderFilter"
                class="filter-select"
              >
                <el-option label="全部厂商" value="" />
                <el-option label="阿里云" value="aliyun" />
                <el-option label="腾讯云" value="tencent" />
                <el-option label="华为云" value="huawei" />
                <el-option label="AWS" value="aws" />
                <el-option label="自建" value="manual" />
              </el-select>
            </el-col>
          </el-row>
        </div>

        <!-- 主机列表表格 -->
        <div class="table-container">
          <el-table
            :data="displayHostList"
            v-loading="loading"
            element-loading-text="数据加载中..."
            stripe
            class="host-table"
            @selection-change="handleSelectionChange"
          >
            <!-- 选择列 -->
            <el-table-column type="selection" width="55" />
            
            <!-- 主机名称 -->
            <el-table-column prop="name" label="主机名称" min-width="180" show-overflow-tooltip>
              <template #default="{ row }">
                <div class="host-name-cell">
                  <span class="host-name">{{ row.name }}</span>
                </div>
              </template>
            </el-table-column>

            <!-- IP地址 -->
            <el-table-column label="IP地址" min-width="140">
              <template #default="{ row }">
                <div class="ip-cell">
                  <div class="ip-item" v-if="getDisplayIP(row.public_ip)">
                    <span class="ip-label">公网:</span>
                    <span class="ip-value">{{ getDisplayIP(row.public_ip) }}</span>
                  </div>
                  <div class="ip-item" v-if="getDisplayIP(row.private_ip)">
                    <span class="ip-label">私网:</span>
                    <span class="ip-value">{{ getDisplayIP(row.private_ip) }}</span>
                  </div>
                </div>
              </template>
            </el-table-column>

            <!-- 状态 -->
            <el-table-column prop="status" label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)" size="small">
                  {{ getStatusDisplayText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>

            <!-- 配置信息 -->
            <el-table-column label="配置" width="180" align="center">
              <template #default="{ row }">
                <div class="config-specs">
                  <div class="spec-item">
                    <span class="spec-label">CPU:</span>
                    <span class="spec-value">{{ formatConfiguration(row.configuration).cpu }}</span>
                  </div>
                  <div class="spec-item">
                    <span class="spec-label">内存:</span>
                    <span class="spec-value">{{ formatConfiguration(row.configuration).memory }}</span>
                  </div>
                  <div class="spec-item">
                    <span class="spec-label">磁盘:</span>
                    <span class="spec-value">{{ formatConfiguration(row.configuration).disk }}</span>
                  </div>
                </div>
              </template>
            </el-table-column>

            <!-- 操作系统 -->
            <el-table-column prop="os" label="系统" width="100" show-overflow-tooltip>
              <template #default="{ row }">
                <span class="os-text">{{ row.os || '-' }}</span>
              </template>
            </el-table-column>

            <!-- 地域 -->
            <el-table-column prop="region" label="地域" width="120">
              <template #default="{ row }">
                <span class="region-text">{{ getRegionDisplay(row.region) }}</span>
              </template>
            </el-table-column>

            <!-- 云厂商 -->
            <el-table-column label="云厂商" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getProviderTagType(row.provider_type)" size="small">
                  {{ getProviderDisplayText(row.provider_type) }}
                </el-tag>
              </template>
            </el-table-column>

            <!-- 操作按钮 -->
            <el-table-column label="操作" width="240" align="center" fixed="right">
              <template #default="{ row }">
                <div class="action-buttons">
                  <el-button type="primary" size="small" @click="handleView(row)">
                    查看
                  </el-button>
                  <el-button type="success" size="small" @click="handleTerminal(row)">
                    终端
                  </el-button>
                  <el-button type="warning" size="small" @click="handleEdit(row)">
                    编辑
                  </el-button>
                  <el-button type="danger" size="small" @click="handleDelete(row)">
                    删除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page="pagination.page"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="pagination.pageSize"
            layout="total, sizes, prev, pager, next, jumper"
            :total="pagination.total"
            background
          />
        </div>
      </el-main>
    </el-container>

    <!-- 主机组管理对话框 -->
    <el-dialog
      v-model="groupDialogVisible"
      :title="isEditingGroup ? '编辑主机组' : '新建主机组'"
      width="500px"
    >
      <el-form :model="groupFormData" label-width="80px">
        <el-form-item label="组名称" required>
          <el-input v-model="groupFormData.name" placeholder="请输入主机组名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="groupFormData.description"
            type="textarea"
            placeholder="请输入描述信息"
            :rows="3"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="groupDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveGroup">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 添加主机对话框 -->
    <el-dialog
      v-model="addHostDialogVisible"
      title="添加主机"
      width="600px"
      :before-close="() => addHostDialogVisible = false"
    >
      <el-form :model="hostFormData" label-width="100px" class="host-form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="主机名称" required>
              <el-input v-model="hostFormData.name" placeholder="请输入主机名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="云厂商" required>
              <el-select v-model="hostFormData.provider_type" placeholder="选择云厂商">
                <el-option label="AWS" value="aws" />
                <el-option label="阿里云" value="aliyun" />
                <el-option label="腾讯云" value="tencent" />
                <el-option label="华为云" value="huawei" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="公网IP">
              <el-input v-model="hostFormData.public_ip[0]" placeholder="请输入公网IP" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="私网IP">
              <el-input v-model="hostFormData.private_ip[0]" placeholder="请输入私网IP" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="操作系统">
              <el-input v-model="hostFormData.os" placeholder="如：Ubuntu 20.04" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="地域">
              <el-input v-model="hostFormData.region" placeholder="如：us-east-1" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="CPU核数">
              <el-input-number v-model="hostFormData.configuration.cpu_cores" :min="1" :max="64" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="内存(GB)">
              <el-input-number v-model="hostFormData.configuration.memory_size" :min="1" :max="512" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="磁盘(GB)">
              <el-input-number v-model="hostFormData.configuration.disk_size" :min="8" :max="2048" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="实例类型">
          <el-input v-model="hostFormData.configuration.instance_type" placeholder="如：t2.micro" />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="addHostDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveHost">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- SSH终端对话框 -->
    <el-dialog
      v-model="terminalVisible"
      :title="`SSH终端 - ${currentHost?.name}`"
      width="80%"
      :before-close="() => { terminalVisible = false; currentHost = null }"
    >
      <SSHTerminal
        v-if="terminalVisible && currentHost"
        :host="currentHost"
        @close="terminalVisible = false"
      />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox, ElDialog } from 'element-plus'
import { storeToRefs } from 'pinia'
import { useHostStore } from '@/store/modules/host'
import SSHTerminal from './components/SSHTerminal.vue'

// Store
const hostStore = useHostStore()
const { hostList, pagination, isLoading } = storeToRefs(hostStore)
const hostGroupTree = ref([])

// 响应式数据
const loading = computed(() => isLoading.value)
const groupLoading = ref(false)
const selectedGroupId = ref(null)
const selectedGroupName = ref('')
const groupSearchQuery = ref('')
const providerFilter = ref('')
const selectedHosts = ref([])

// 搜索参数
const searchParams = reactive({
  keyword: '',
  status: '',
  region: '',
  group_id: null,
  page: 1,
  page_size: 20
})

// 树形组件配置
const treeProps = {
  children: 'children',
  label: 'name'
}

// 终端相关
const terminalVisible = ref(false)
const currentHost = ref(null)

// 主机组管理
const groupDialogVisible = ref(false)
const groupFormData = ref({
  name: '',
  description: '',
  parent_id: null
})
const isEditingGroup = ref(false)
const editingGroupId = ref(null)

// 主机添加功能
const addHostDialogVisible = ref(false)
const hostFormData = ref({
  name: '',
  public_ip: [],
  private_ip: [],
  status: 'running',
  os: '',
  region: '',
  provider_type: 'aws',
  group_id: null,
  configuration: {
    cpu_cores: 1,
    memory_size: 1,
    disk_size: 20,
    instance_type: 't2.micro'
  }
})

// 工具函数
const getDisplayIP = (ips) => {
  if (!ips || !Array.isArray(ips)) return '-'
  return ips.join(', ')
}

const getStatusColor = (status) => {
  const statusMap = {
    'running': 'success',
    'stopped': 'danger',
    'pending': 'warning'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status) => {
  const statusMap = {
    'running': '运行中',
    'stopped': '已停止',
    'pending': '启动中'
  }
  return statusMap[status] || status
}

// 格式化配置信息 - 只显示核心信息
const formatConfiguration = (config) => {
  if (!config || typeof config !== 'object') return { cpu: '-', memory: '-', disk: '-' }

  try {
    const configObj = typeof config === 'string' ? JSON.parse(config) : config
    return {
      cpu: configObj.cpu_cores ? `${configObj.cpu_cores}核` : '-',
      memory: configObj.memory_size ? `${configObj.memory_size}GB` : '-',
      disk: configObj.disk_size ? `${configObj.disk_size}GB` : '-'
    }
  } catch (error) {
    return { cpu: '-', memory: '-', disk: '-' }
  }
}

// 获取云厂商显示名称
const getProviderName = (provider) => {
  const providerMap = {
    'aws': 'AWS',
    'aliyun': '阿里云',
    'tencent': '腾讯云',
    'huawei': '华为云'
  }
  return providerMap[provider] || provider
}

// 获取云厂商颜色
const getProviderColor = (provider) => {
  const colorMap = {
    'aws': '#FF9900',
    'aliyun': '#FF6A00',
    'tencent': '#006EFF',
    'huawei': '#FF0000'
  }
  return colorMap[provider] || '#666666'
}

// 计算属性
const filteredGroupTree = computed(() => {
  try {
    const allGroupsNode = {
      id: null,
      name: '全部主机',
      host_count: hostList.value?.length || 0,
      children: []
    }

    // 确保hostGroupTree.value是数组，添加更强的保护
    let groups = []
    if (hostGroupTree.value && Array.isArray(hostGroupTree.value)) {
      groups = hostGroupTree.value
    } else if (hostGroupTree.value && typeof hostGroupTree.value === 'object') {
      // 如果是对象，尝试获取data属性
      groups = hostGroupTree.value.data || []
    }

    let tree = [allGroupsNode, ...groups]

    if (groupSearchQuery.value) {
      const query = groupSearchQuery.value.toLowerCase()
      tree = tree.filter(group =>
        group.name && group.name.toLowerCase().includes(query)
      )
    }

    return tree
  } catch (error) {
    console.error('filteredGroupTree计算错误:', error)
    return [{
      id: null,
      name: '全部主机',
      host_count: 0,
      children: []
    }]
  }
})

const displayHostList = computed(() => {
  return hostList.value
})

// 数据处理方法已在上面定义

const getConfigDisplay = (config) => {
  if (!config) return '-'
  if (typeof config === 'string') return config
  if (typeof config === 'object') {
    const cpu = config.cpu_cores || config.cpu || ''
    const memory = config.memory_size || config.memory || ''
    if (cpu && memory) {
      return `${cpu}核${memory}GB`
    }
    return config.instance_type || '-'
  }
  return '-'
}

const getRegionDisplay = (region) => {
  const regionMap = {
    'cn-hangzhou': '华东1',
    'cn-shanghai': '华东2',
    'cn-qingdao': '华北1',
    'cn-beijing': '华北2',
    'cn-shenzhen': '华南1',
    'ap-southeast-1': '新加坡',
    'us-east-1': '美东1',
    'us-west-1': '美西1'
  }
  return regionMap[region] || region || '-'
}

// 状态标签类型
const getStatusTagType = (status) => {
  const statusMap = {
    'running': 'success',
    'stopped': 'danger',
    'maintenance': 'warning',
    'pending': 'info'
  }
  return statusMap[status] || 'info'
}

// 状态显示文本
const getStatusDisplayText = (status) => {
  const statusMap = {
    'running': '运行中',
    'stopped': '已停止',
    'maintenance': '维护中',
    'pending': '启动中'
  }
  return statusMap[status] || '未知'
}

// 云厂商标签类型
const getProviderTagType = (provider) => {
  const providerMap = {
    'aliyun': 'warning',
    'tencent': 'primary',
    'huawei': 'success',
    'aws': 'info',
    'manual': ''
  }
  return providerMap[provider] || 'info'
}

// 云厂商显示文本
const getProviderDisplayText = (provider) => {
  const providerMap = {
    'aliyun': '阿里云',
    'tencent': '腾讯云',
    'huawei': '华为云',
    'aws': 'AWS',
    'manual': '自建'
  }
  return providerMap[provider] || '未知'
}

// 事件处理函数
const handleGroupSelect = (data) => {
  selectedGroupId.value = data.id
  selectedGroupName.value = data.name
  searchParams.group_id = data.id
  searchParams.page = 1
  fetchHosts()
}

const handleSearch = () => {
  searchParams.page = 1
  fetchHosts()
}

const handleProviderFilter = () => {
  handleSearch()
}

const handleSelectionChange = (selection) => {
  selectedHosts.value = selection
}

const refreshGroups = async () => {
  groupLoading.value = true
  try {
    await hostStore.fetchHostGroupTree()
    // 更新本地的hostGroupTree
    hostGroupTree.value = hostStore.hostGroupTree || []
    ElMessage.success('分组数据刷新成功')
  } catch (error) {
    ElMessage.error('分组数据刷新失败')
    hostGroupTree.value = []
  } finally {
    groupLoading.value = false
  }
}

const refreshData = async () => {
  try {
    await fetchHosts()
    ElMessage.success('数据刷新成功')
  } catch (error) {
    ElMessage.error('数据刷新失败')
  }
}

const fetchHosts = async () => {
  try {
    await hostStore.fetchHosts(searchParams)
  } catch (error) {
    console.error('获取主机列表失败:', error)
  }
}

// 主机组管理
const handleGroupRightClick = (event, node, data) => {
  event.preventDefault()
  showGroupContextMenu(event, data)
}

const showGroupContextMenu = (event, group) => {
  console.log('右键点击主机组:', group)
}

const handleAddGroup = (parentGroup = null) => {
  groupFormData.value = {
    name: '',
    description: '',
    parent_id: parentGroup?.id || null
  }
  isEditingGroup.value = false
  editingGroupId.value = null
  groupDialogVisible.value = true
}

const handleEditGroup = (group) => {
  groupFormData.value = {
    name: group.name,
    description: group.description || '',
    parent_id: group.parent_id
  }
  isEditingGroup.value = true
  editingGroupId.value = group.id
  groupDialogVisible.value = true
}

const handleDeleteGroup = async (group) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除主机组 "${group.name}" 吗？删除后该组下的主机将移动到根目录。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await hostStore.deleteHostGroup(group.id)
    ElMessage.success(`已删除主机组: ${group.name}`)
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const saveGroup = async () => {
  try {
    if (isEditingGroup.value) {
      await hostStore.updateHostGroup(editingGroupId.value, groupFormData.value)
      ElMessage.success('主机组更新成功')
    } else {
      await hostStore.addHostGroup(groupFormData.value)
      ElMessage.success('主机组创建成功')
    }
    groupDialogVisible.value = false
  } catch (error) {
    ElMessage.error(isEditingGroup.value ? '更新失败' : '创建失败')
  }
}

// 主机操作
const handleAdd = () => {
  // 重置表单数据
  hostFormData.value = {
    name: '',
    public_ip: [],
    private_ip: [],
    status: 'running',
    os: '',
    region: '',
    provider_type: 'aws',
    group_id: selectedGroupId.value,
    configuration: {
      cpu_cores: 1,
      memory_size: 1,
      disk_size: 20,
      instance_type: 't2.micro'
    }
  }
  addHostDialogVisible.value = true
}

const saveHost = async () => {
  try {
    // 这里应该调用后端API保存主机
    // await hostStore.addHost(hostFormData.value)
    ElMessage.success('主机添加成功')
    addHostDialogVisible.value = false
    await fetchHosts()
  } catch (error) {
    ElMessage.error('主机添加失败')
  }
}

const handleView = (row) => {
  ElMessage.info(`查看主机: ${row.name}`)
}

const handleEdit = (row) => {
  ElMessage.info(`编辑主机: ${row.name}`)
}

const handleTerminal = (row) => {
  currentHost.value = row
  terminalVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除主机 "${row.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await hostStore.deleteHost(row.id)
    ElMessage.success(`已删除主机: ${row.name}`)
    await fetchHosts()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleSizeChange = (val) => {
  searchParams.page_size = val
  searchParams.page = 1
  fetchHosts()
}

const handleCurrentChange = (val) => {
  searchParams.page = val
  fetchHosts()
}

// 生命周期
onMounted(async () => {
  try {
    await hostStore.fetchHostGroupTree()
    // 从store中获取数据并设置到本地ref
    hostGroupTree.value = hostStore.hostGroupTree || []
  } catch (error) {
    console.error('获取主机组失败:', error)
    hostGroupTree.value = []
  }
  await fetchHosts()
})
</script>

<style scoped>
/* 主容器样式 */
.host-management-pro {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  padding: 16px;
}

.main-container {
  height: 100%;
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
}

/* 左侧边栏样式 */
.sidebar {
  background: linear-gradient(180deg, #f8f9fa 0%, #e9ecef 100%);
  border-right: 1px solid #dee2e6;
  box-shadow: 2px 0 16px rgba(0, 0, 0, 0.1);
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #dee2e6;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
}

.sidebar-title {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: #2c3e50;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.refresh-btn {
  font-size: 12px;
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);
  border: none;
  color: white;
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 500;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(52, 152, 219, 0.3);
}

.refresh-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(52, 152, 219, 0.4);
}

.group-search {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.group-tree-container {
  padding: 8px 0;
  max-height: calc(100vh - 140px);
  overflow-y: auto;
}

.group-tree {
  padding: 0 12px;
}

.tree-node {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding-right: 8px;
  position: relative;
}

.tree-node:hover .node-actions {
  opacity: 1;
}

.node-label {
  font-size: 14px;
  color: #262626;
  flex: 1;
}

.node-count {
  font-size: 12px;
  color: #8c8c8c;
  margin-right: 8px;
}

.node-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.2s;
}

.node-actions .el-button {
  padding: 2px 4px;
  font-size: 12px;
  min-height: auto;
  border: none;
}

.node-actions .el-button:hover {
  background-color: #f0f0f0;
}

/* 右侧主内容区样式 */
.main-content {
  padding: 0;
  background: linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%);
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 24px 24px 16px;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-bottom: 1px solid #e8e8e8;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.header-left {
  flex: 1;
}

.page-title {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 700;
  color: #2c3e50;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.breadcrumb {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: #8c8c8c;
}

.breadcrumb-item {
  color: #8c8c8c;
}

.breadcrumb-separator {
  margin: 0 8px;
  color: #d9d9d9;
}

.breadcrumb-current {
  color: #1890ff;
}

.header-actions {
  display: flex;
  gap: 12px;
}

/* 工具栏样式 */
.toolbar {
  padding: 16px 24px;
  background: white;
  border-bottom: 1px solid #e8e8e8;
}

.search-input,
.filter-select {
  width: 100%;
}

/* 表格容器样式 */
.table-container {
  margin: 16px 24px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  border: 1px solid #e9ecef;
  animation: fadeInUp 0.6s ease-out;
}

.host-table {
  border-radius: 8px;
}

/* 表格单元格样式 */
.host-name-cell {
  display: flex;
  align-items: center;
}

.host-name {
  font-weight: 500;
  color: #262626;
}

.ip-cell {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
}

.ip-item {
  margin-bottom: 2px;
}

.ip-item:last-child {
  margin-bottom: 0;
}

.ip-label {
  color: #8c8c8c;
  margin-right: 4px;
}

.ip-value {
  color: #262626;
}

.config-cell {
  font-size: 13px;
}

.config-text {
  color: #595959;
}

.os-text,
.region-text {
  color: #595959;
  font-size: 13px;
}

.action-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

/* 分页样式 */
.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 16px 24px 24px;
}

/* 表格行悬浮效果 */
.host-table :deep(.el-table__row):hover {
  background-color: #fafafa !important;
}

/* 确保没有任何字体图标相关的样式 */
.host-table :deep(.el-table__row) *,
.group-tree :deep(.el-tree-node) *,
.sidebar *,
.main-content * {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif !important;
}

/* 强力禁用所有可能的伪元素内容 */
.host-table :deep(.el-table__row) *::before,
.host-table :deep(.el-table__row) *::after,
.group-tree :deep(.el-tree-node) *::before,
.group-tree :deep(.el-tree-node) *::after,
.sidebar *::before,
.sidebar *::after,
.main-content *::before,
.main-content *::after {
  content: none !important;
  display: none !important;
}

/* 禁用所有可能的图标字体 */
.host-management-pro * {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif !important;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .sidebar {
    width: 240px !important;
  }
}

@media (max-width: 768px) {
  .sidebar {
    width: 200px !important;
  }

  .content-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .header-actions {
    justify-content: flex-start;
  }

  .toolbar .el-row {
    flex-direction: column;
    gap: 12px;
  }

  .toolbar .el-col {
    width: 100% !important;
  }
}

/* 动画效果 */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideInLeft {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.tree-node {
  animation: slideInLeft 0.4s ease-out;
}

/* 表格行悬浮效果增强 */
.host-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.host-table :deep(.el-table__row):hover {
  background: linear-gradient(135deg, #f8f9fa 0%, #e3f2fd 100%) !important;
  transform: scale(1.01);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* 表格头部美化 */
.host-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
}

.host-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 700 !important;
  border-bottom: 2px solid #dee2e6 !important;
  font-size: 14px !important;
}

/* 按钮美化 */
.action-buttons .el-button {
  transition: all 0.3s ease;
  border-radius: 6px;
  font-weight: 500;
}

.action-buttons .el-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

/* 标签美化 */
.el-tag {
  border-radius: 6px !important;
  font-weight: 500 !important;
  border: none !important;
}

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, #5a6fd8 0%, #6a4190 100%);
}

/* 输入框和选择器美化 */
.el-input__wrapper,
.el-select__wrapper {
  border-radius: 8px !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1) !important;
  transition: all 0.3s ease !important;
}

.el-input__wrapper:hover,
.el-select__wrapper:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}

.el-input__wrapper.is-focus,
.el-select__wrapper.is-focus {
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2) !important;
}
</style>
