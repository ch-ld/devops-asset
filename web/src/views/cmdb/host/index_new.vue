<template>
  <div class="host-management">
    <div class="page-container">
      <!-- 左侧边栏 -->
      <div class="sidebar">
        <div class="sidebar-header">
          <div class="header-content">
            <h3 class="title">
              <i class="icon">📁</i>
              主机分组
            </h3>
            <el-button size="small" type="primary" @click="refreshGroups" :loading="groupLoading">
              刷新
            </el-button>
          </div>
        </div>
        
        <!-- 搜索框 -->
        <div class="search-section">
          <el-input
            v-model="groupSearchQuery"
            placeholder="搜索分组..."
            size="small"
            clearable
            class="search-input"
          >
            <template #prefix>
              <i class="search-icon">🔍</i>
            </template>
          </el-input>
        </div>

        <!-- 分组列表 -->
        <div class="group-list">
          <div 
            v-for="group in filteredGroups" 
            :key="group.id"
            :class="['group-item', { active: selectedGroupId === group.id }]"
            @click="handleGroupSelect(group)"
          >
            <div class="group-content">
              <div class="group-info">
                <i class="group-icon">{{ group.id === null ? '🏠' : '📂' }}</i>
                <span class="group-name">{{ group.name }}</span>
              </div>
              <div class="group-count">{{ group.host_count || 0 }}</div>
            </div>
            <div class="group-actions" v-if="group.id !== null">
              <el-button size="small" text @click.stop="handleEditGroup(group)" title="编辑">
                ✏️
              </el-button>
              <el-button size="small" text @click.stop="handleDeleteGroup(group)" title="删除">
                🗑️
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 主内容区 -->
      <div class="main-content">
        <!-- 页面头部 -->
        <div class="content-header">
          <div class="header-left">
            <h1 class="page-title">主机管理</h1>
            <div class="breadcrumb">
              <span>主机管理</span>
              <span class="separator">/</span>
              <span class="current">{{ currentGroupName }}</span>
            </div>
          </div>
          <div class="header-actions">
            <el-button @click="refreshHosts" :loading="loading">
              刷新数据
            </el-button>
            <el-button type="primary" @click="handleAdd">
              添加主机
            </el-button>
          </div>
        </div>

        <!-- 搜索工具栏 -->
        <div class="toolbar">
          <div class="toolbar-left">
            <el-input
              v-model="searchQuery"
              placeholder="搜索主机名称、IP地址..."
              size="default"
              clearable
              class="search-input"
            >
              <template #prefix>
                <i>🔍</i>
              </template>
            </el-input>
          </div>
          <div class="toolbar-right">
            <el-select v-model="statusFilter" placeholder="状态筛选" clearable size="default">
              <el-option label="运行中" value="running" />
              <el-option label="已停止" value="stopped" />
              <el-option label="重启中" value="restarting" />
            </el-select>
            <el-select v-model="regionFilter" placeholder="地域筛选" clearable size="default">
              <el-option v-for="region in regions" :key="region" :label="region" :value="region" />
            </el-select>
            <el-select v-model="providerFilter" placeholder="云厂商筛选" clearable size="default">
              <el-option label="AWS" value="aws" />
              <el-option label="阿里云" value="aliyun" />
              <el-option label="腾讯云" value="tencent" />
              <el-option label="华为云" value="huawei" />
            </el-select>
          </div>
        </div>

        <!-- 主机列表 -->
        <div class="table-container">
          <el-table
            :data="filteredHosts"
            v-loading="loading"
            class="host-table"
            stripe
            @selection-change="handleSelectionChange"
          >
            <!-- 选择框 -->
            <el-table-column type="selection" width="50" />
            
            <!-- 主机名称 -->
            <el-table-column prop="name" label="主机名称" min-width="150" show-overflow-tooltip>
              <template #default="{ row }">
                <div class="host-name">
                  <i class="host-icon">🖥️</i>
                  <span>{{ row.name }}</span>
                </div>
              </template>
            </el-table-column>

            <!-- IP地址 -->
            <el-table-column label="IP地址" min-width="160">
              <template #default="{ row }">
                <div class="ip-info">
                  <div class="ip-item" v-if="getDisplayIP(row.public_ip)">
                    <span class="ip-label">公网</span>
                    <span class="ip-value">{{ getDisplayIP(row.public_ip) }}</span>
                  </div>
                  <div class="ip-item" v-if="getDisplayIP(row.private_ip)">
                    <span class="ip-label">私网</span>
                    <span class="ip-value">{{ getDisplayIP(row.private_ip) }}</span>
                  </div>
                  <div class="ip-item" v-if="!getDisplayIP(row.public_ip) && !getDisplayIP(row.private_ip)">
                    <span class="ip-empty">暂无IP</span>
                  </div>
                </div>
              </template>
            </el-table-column>

            <!-- 状态 -->
            <el-table-column label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>

            <!-- 配置 -->
            <el-table-column label="配置" width="180" align="center">
              <template #default="{ row }">
                <div class="config-info">
                  <div class="config-item">
                    <span class="config-label">CPU</span>
                    <span class="config-value">{{ formatConfiguration(row.configuration).cpu }}</span>
                  </div>
                  <div class="config-item">
                    <span class="config-label">内存</span>
                    <span class="config-value">{{ formatConfiguration(row.configuration).memory }}</span>
                  </div>
                  <div class="config-item">
                    <span class="config-label">磁盘</span>
                    <span class="config-value">{{ formatConfiguration(row.configuration).disk }}</span>
                  </div>
                </div>
              </template>
            </el-table-column>

            <!-- 系统 -->
            <el-table-column prop="os" label="系统" width="120" show-overflow-tooltip>
              <template #default="{ row }">
                <span class="os-info">{{ row.os || '-' }}</span>
              </template>
            </el-table-column>

            <!-- 地域 -->
            <el-table-column prop="region" label="地域" width="120">
              <template #default="{ row }">
                <span class="region-info">{{ row.region || '-' }}</span>
              </template>
            </el-table-column>

            <!-- 云厂商 -->
            <el-table-column label="云厂商" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getProviderType(row.provider_type)" size="small">
                  {{ getProviderText(row.provider_type) }}
                </el-tag>
              </template>
            </el-table-column>

            <!-- 操作 -->
            <el-table-column label="操作" width="200" align="center" fixed="right">
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
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </div>

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
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// 响应式数据
const loading = ref(false)
const groupLoading = ref(false)
const selectedGroupId = ref(null)
const searchQuery = ref('')
const groupSearchQuery = ref('')
const statusFilter = ref('')
const regionFilter = ref('')
const providerFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const addHostDialogVisible = ref(false)

// 主机数据
const hosts = ref([])
const groups = ref([
  { id: null, name: '全部主机', host_count: 0 },
  { id: 1, name: '测试组', host_count: 0 },
  { id: 2, name: '生产组', host_count: 0 },
  { id: 3, name: '开发组', host_count: 0 }
])

// 表单数据
const hostFormData = reactive({
  name: '',
  provider_type: 'aws',
  public_ip: [''],
  private_ip: [''],
  os: '',
  region: '',
  configuration: {
    cpu_cores: 1,
    memory_size: 1,
    disk_size: 20,
    instance_type: 't2.micro'
  }
})

// 计算属性
const currentGroupName = computed(() => {
  const group = groups.value.find(g => g.id === selectedGroupId.value)
  return group ? group.name : '全部主机'
})

const filteredGroups = computed(() => {
  if (!groupSearchQuery.value) return groups.value
  return groups.value.filter(group =>
    group.name.toLowerCase().includes(groupSearchQuery.value.toLowerCase())
  )
})

const filteredHosts = computed(() => {
  let result = hosts.value

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(host =>
      host.name.toLowerCase().includes(query) ||
      (host.public_ip && host.public_ip.some(ip => ip.includes(query))) ||
      (host.private_ip && host.private_ip.some(ip => ip.includes(query)))
    )
  }

  // 状态过滤
  if (statusFilter.value) {
    result = result.filter(host => host.status === statusFilter.value)
  }

  // 地域过滤
  if (regionFilter.value) {
    result = result.filter(host => host.region === regionFilter.value)
  }

  // 云厂商过滤
  if (providerFilter.value) {
    result = result.filter(host => host.provider_type === providerFilter.value)
  }

  return result
})

const regions = computed(() => {
  const regionSet = new Set()
  hosts.value.forEach(host => {
    if (host.region) regionSet.add(host.region)
  })
  return Array.from(regionSet)
})

// 工具函数
const getDisplayIP = (ip) => {
  if (!ip) return ''
  if (Array.isArray(ip)) return ip[0] || ''
  return ip
}

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

const getStatusType = (status) => {
  const statusMap = {
    'running': 'success',
    'stopped': 'danger',
    'restarting': 'warning',
    'pending': 'info'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status) => {
  const statusMap = {
    'running': '运行中',
    'stopped': '已停止',
    'restarting': '重启中',
    'pending': '待启动'
  }
  return statusMap[status] || status
}

const getProviderType = (provider) => {
  const providerMap = {
    'aws': 'warning',
    'aliyun': 'success',
    'tencent': 'primary',
    'huawei': 'info'
  }
  return providerMap[provider] || 'info'
}

const getProviderText = (provider) => {
  const providerMap = {
    'aws': 'AWS',
    'aliyun': '阿里云',
    'tencent': '腾讯云',
    'huawei': '华为云'
  }
  return providerMap[provider] || provider
}

// 事件处理
const refreshGroups = async () => {
  groupLoading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success('分组列表刷新成功')
  } catch (error) {
    ElMessage.error('刷新失败')
  } finally {
    groupLoading.value = false
  }
}

const refreshHosts = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success('主机列表刷新成功')
  } catch (error) {
    ElMessage.error('刷新失败')
  } finally {
    loading.value = false
  }
}

const handleGroupSelect = (group) => {
  selectedGroupId.value = group.id
  // 根据分组筛选主机
  refreshHosts()
}

const handleAdd = () => {
  // 重置表单
  Object.assign(hostFormData, {
    name: '',
    provider_type: 'aws',
    public_ip: [''],
    private_ip: [''],
    os: '',
    region: '',
    configuration: {
      cpu_cores: 1,
      memory_size: 1,
      disk_size: 20,
      instance_type: 't2.micro'
    }
  })
  addHostDialogVisible.value = true
}

const saveHost = async () => {
  try {
    // 验证表单
    if (!hostFormData.name) {
      ElMessage.error('请输入主机名称')
      return
    }

    // 模拟保存
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success('主机添加成功')
    addHostDialogVisible.value = false
    refreshHosts()
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

const handleView = (row) => {
  ElMessage.info(`查看主机: ${row.name}`)
}

const handleTerminal = (row) => {
  ElMessage.info(`连接终端: ${row.name}`)
}

const handleEdit = (row) => {
  ElMessage.info(`编辑主机: ${row.name}`)
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除主机 "${row.name}" 吗？`, '确认删除', {
      type: 'warning'
    })
    ElMessage.success('删除成功')
    refreshHosts()
  } catch {
    // 用户取消
  }
}

const handleEditGroup = (group) => {
  ElMessage.info(`编辑分组: ${group.name}`)
}

const handleDeleteGroup = async (group) => {
  try {
    await ElMessageBox.confirm(`确定要删除分组 "${group.name}" 吗？`, '确认删除', {
      type: 'warning'
    })
    ElMessage.success('删除成功')
    refreshGroups()
  } catch {
    // 用户取消
  }
}

const handleSelectionChange = (selection) => {
  console.log('选中的主机:', selection)
}

const handleSizeChange = (size) => {
  pageSize.value = size
  refreshHosts()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  refreshHosts()
}

// 初始化
onMounted(() => {
  selectedGroupId.value = null
  refreshHosts()

  // 模拟主机数据
  hosts.value = [
    {
      id: 1,
      name: 'web-server-01',
      public_ip: ['54.123.45.67'],
      private_ip: ['10.0.1.10'],
      status: 'running',
      os: 'Ubuntu 20.04',
      region: 'us-east-1',
      provider_type: 'aws',
      configuration: {
        cpu_cores: 2,
        memory_size: 4,
        disk_size: 40,
        instance_type: 't3.small'
      }
    },
    {
      id: 2,
      name: 'db-server-01',
      public_ip: [],
      private_ip: ['10.0.1.20'],
      status: 'running',
      os: 'CentOS 7',
      region: 'us-east-1',
      provider_type: 'aws',
      configuration: {
        cpu_cores: 4,
        memory_size: 16,
        disk_size: 100,
        instance_type: 't3.large'
      }
    }
  ]

  // 更新分组主机数量
  groups.value[0].host_count = hosts.value.length
})
</script>

<style scoped>
/* 主容器 */
.host-management {
  height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.page-container {
  display: flex;
  height: 100%;
  max-width: 1400px;
  margin: 0 auto;
  background: white;
  box-shadow: 0 0 30px rgba(0, 0, 0, 0.1);
}

/* 左侧边栏 */
.sidebar {
  width: 280px;
  background: linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%);
  border-right: 1px solid #e9ecef;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #e9ecef;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
}

.icon {
  font-size: 20px;
}

.search-section {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.search-input {
  width: 100%;
}

.search-icon {
  color: #999;
  font-size: 14px;
}

/* 分组列表 */
.group-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.group-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  margin: 2px 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
}

.group-item:hover {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  transform: translateX(4px);
}

.group-item.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.group-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.group-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.group-icon {
  font-size: 16px;
}

.group-name {
  font-weight: 500;
  font-size: 14px;
}

.group-count {
  background: rgba(0, 0, 0, 0.1);
  color: inherit;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.group-item.active .group-count {
  background: rgba(255, 255, 255, 0.2);
}

.group-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.group-item:hover .group-actions {
  opacity: 1;
}

/* 主内容区 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #fafbfc;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 24px 32px;
  background: white;
  border-bottom: 1px solid #e9ecef;
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6c757d;
  font-size: 14px;
}

.separator {
  color: #adb5bd;
}

.current {
  color: #495057;
  font-weight: 500;
}

.header-actions {
  display: flex;
  gap: 12px;
}

/* 工具栏 */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 32px;
  background: white;
  border-bottom: 1px solid #f0f0f0;
}

.toolbar-left {
  flex: 1;
  max-width: 400px;
}

.toolbar-right {
  display: flex;
  gap: 12px;
}

/* 表格容器 */
.table-container {
  flex: 1;
  margin: 20px 32px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.host-table {
  width: 100%;
}

/* 主机名称 */
.host-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.host-icon {
  font-size: 16px;
}

/* IP信息 */
.ip-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.ip-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}

.ip-label {
  background: #f8f9fa;
  color: #495057;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
  min-width: 32px;
  text-align: center;
}

.ip-value {
  font-family: 'Monaco', 'Menlo', monospace;
  color: #2c3e50;
  font-weight: 500;
}

.ip-empty {
  color: #adb5bd;
  font-style: italic;
}

/* 配置信息 */
.config-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.config-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.config-label {
  color: #6c757d;
  font-weight: 500;
  min-width: 32px;
}

.config-value {
  color: #2c3e50;
  font-weight: 600;
  font-family: 'Monaco', 'Menlo', monospace;
}

/* 系统和地域信息 */
.os-info, .region-info {
  color: #495057;
  font-size: 13px;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: 4px;
}

.action-buttons .el-button {
  padding: 4px 8px;
  font-size: 12px;
}

/* 分页 */
.pagination-container {
  display: flex;
  justify-content: center;
  padding: 20px 32px;
  background: white;
  border-top: 1px solid #f0f0f0;
}

/* 对话框 */
.host-form {
  padding: 20px 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 表格样式增强 */
.host-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
}

.host-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 600 !important;
  border-bottom: 2px solid #dee2e6 !important;
  font-size: 13px !important;
}

.host-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.host-table :deep(.el-table__row):hover {
  background: linear-gradient(135deg, #f8f9fa 0%, #e3f2fd 100%) !important;
}

/* 标签美化 */
.el-tag {
  border: none !important;
  font-weight: 500 !important;
  border-radius: 6px !important;
}

/* 按钮美化 */
.el-button {
  border-radius: 6px !important;
  font-weight: 500 !important;
  transition: all 0.3s ease !important;
}

.el-button:hover {
  transform: translateY(-1px) !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}

/* 输入框美化 */
.el-input :deep(.el-input__wrapper) {
  border-radius: 8px !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05) !important;
  transition: all 0.3s ease !important;
}

.el-input :deep(.el-input__wrapper):hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}

.el-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2) !important;
}

/* 选择器美化 */
.el-select :deep(.el-select__wrapper) {
  border-radius: 8px !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05) !important;
}

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, #5a6fd8 0%, #6a4190 100%);
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .page-container {
    max-width: 100%;
  }

  .sidebar {
    width: 240px;
  }
}

@media (max-width: 768px) {
  .sidebar {
    display: none;
  }

  .content-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .toolbar {
    flex-direction: column;
    gap: 16px;
  }

  .toolbar-right {
    flex-wrap: wrap;
  }

  .table-container {
    margin: 16px;
  }

  .action-buttons {
    flex-direction: column;
    gap: 2px;
  }
}
</style>
