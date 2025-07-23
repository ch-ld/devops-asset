<template>
  <div class="provider-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <div class="page-title">
            <h1>云账号管理</h1>
            <p>管理和监控您的云厂商账号配置</p>
          </div>
        </div>
        <div class="header-right">
          <el-button type="primary" size="large" @click="handleCreate">
            <el-icon>{{ iconMap.Plus }}</el-icon>
            添加云账号
          </el-button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-icon aliyun">
          <div class="icon-placeholder">阿里云</div>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ getProviderCount('aliyun') }}</div>
          <div class="stat-label">阿里云账号</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon tencent">
          <div class="icon-placeholder">腾讯云</div>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ getProviderCount('tencent') }}</div>
          <div class="stat-label">腾讯云账号</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon aws">
          <div class="icon-placeholder">AWS</div>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ getProviderCount('aws') }}</div>
          <div class="stat-label">AWS账号</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon total">
          <el-icon size="24"><CloudServer /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ providerList.length }}</div>
          <div class="stat-label">总计账号</div>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-section">
      <el-card class="search-card" shadow="never">
        <el-form :model="searchForm" inline class="search-form">
          <el-form-item label="账号名称">
            <el-input
              v-model="searchForm.name"
              placeholder="请输入账号名称"
              clearable
              style="width: 200px"
              @keyup.enter="handleSearch"
            >
              <template #prefix>
                <el-icon>{{ iconMap.Search }}</el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item label="云厂商">
            <el-select
              v-model="searchForm.provider_type"
              placeholder="请选择云厂商"
              clearable
              style="width: 150px"
            >
              <el-option label="阿里云" value="aliyun" />
              <el-option label="腾讯云" value="tencent" />
              <el-option label="AWS" value="aws" />
            </el-select>
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" @click="handleSearch">
              <el-icon>{{ iconMap.Search }}</el-icon>
              搜索
            </el-button>
            <el-button @click="handleReset">
              <el-icon>{{ iconMap.Refresh }}</el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>

    <!-- 云账号列表 -->
    <div class="provider-list">
      <el-card v-if="providerList.length === 0 && !loading" class="empty-card" shadow="never">
        <el-empty description="暂无云账号">
          <el-button type="primary" @click="handleCreate">添加第一个云账号</el-button>
        </el-empty>
      </el-card>

      <div v-else class="provider-grid">
        <div
          v-for="provider in providerList"
          :key="provider.id"
          class="provider-card"
          :class="`provider-${provider.type}`"
        >
          <!-- 卡片头部 -->
          <div class="card-header">
            <div class="provider-info">
              <div class="provider-avatar" :class="`avatar-${provider.type}`">
                <div class="avatar-text">{{ getProviderLabel(provider.type).charAt(0) }}</div>
              </div>
              <div class="provider-details">
                <h3 class="provider-name">{{ provider.name }}</h3>
                <p class="provider-type">{{ getProviderLabel(provider.type) }}</p>
              </div>
            </div>
            <div class="card-actions">
              <el-dropdown @command="handleCommand">
                <el-button type="text" class="more-btn">
                  <el-icon>{{ iconMap.MoreFilled }}</el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item :command="`edit-${provider.id}`">
                      <el-icon>{{ iconMap.Edit }}</el-icon>
                      编辑
                    </el-dropdown-item>
                    <el-dropdown-item :command="`sync-${provider.id}`">
                      <el-icon>{{ iconMap.Refresh }}</el-icon>
                      同步资源
                    </el-dropdown-item>
                    <el-dropdown-item :command="`delete-${provider.id}`" divided>
                      <el-icon>{{ iconMap.Delete }}</el-icon>
                      删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>

          <!-- 卡片内容 -->
          <div class="card-content">
            <div class="info-item">
              <span class="info-label">AccessKey ID:</span>
              <span class="info-value">{{ maskAccessKey(provider.access_key_id) }}</span>
            </div>
            <div class="info-item" v-if="provider.region">
              <span class="info-label">默认区域:</span>
              <span class="info-value">{{ provider.region }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">创建时间:</span>
              <span class="info-value">{{ formatDate(provider.created_at) }}</span>
            </div>
          </div>

          <!-- 卡片底部 -->
          <div class="card-footer">
            <div class="status-info">
              <el-tag :type="getStatusType(provider.status)" size="small">
                {{ getStatusText(provider.status) }}
              </el-tag>
            </div>
            <div class="quick-actions">
              <el-button size="small" type="primary" @click="handleEdit(provider)">
                编辑
              </el-button>
              <el-button size="small" @click="handleSync(provider)">
                同步
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="3" animated />
    </div>

    <!-- 云账号弹窗 -->
    <ProviderModal
      :visible="showModal"
      :record="currentProvider"
      @success="handleModalSuccess"
      @cancel="handleModalCancel"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
// 使用字符串图标替代 Element Plus 图标
const iconMap = {
  Plus: '➕',
  Search: '🔍',
  Refresh: '🔄',
  CloudServer: '☁️',
  Edit: '✏️',
  Delete: '🗑️',
  MoreFilled: '⋯'
}
import {
  getProviderList,
  deleteProvider,
  syncResources
} from '@/api/system/host'
import ProviderModal from './ProviderModal.vue'

defineOptions({ name: 'ProviderManagement' })

// 响应式数据
const loading = ref(false)
const providerList = ref([])
const showModal = ref(false)
const currentProvider = ref(null)
// const providerModalRef = ref() // 不再需要ref

// 搜索表单
const searchForm = reactive({
  name: '',
  provider_type: undefined,
})

// 方法
async function fetchProviderList() {
  loading.value = true
  try {
    const response = await getProviderList(searchForm)
    providerList.value = response.data || []
  } catch (error) {
    console.error('获取云账号列表失败:', error)
    ElMessage.error('获取云账号列表失败')
  } finally {
    loading.value = false
  }
}

function getProviderCount(type: string) {
  return providerList.value.filter((p: any) => p.type === type).length
}



function getProviderLabel(type: string) {
  const labels = {
    aliyun: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS'
  }
  return labels[type] || type
}

function maskAccessKey(accessKey: string) {
  if (!accessKey) return ''
  return accessKey.substring(0, 8) + '****' + accessKey.substring(accessKey.length - 4)
}

function getStatusType(status: string) {
  const types = {
    active: 'success',
    inactive: 'danger',
    pending: 'warning'
  }
  return types[status] || 'info'
}

function getStatusText(status: string) {
  const texts = {
    active: '正常',
    inactive: '异常',
    pending: '待验证'
  }
  return texts[status] || '未知'
}

function handleSearch() {
  fetchProviderList()
}

function handleReset() {
  searchForm.name = ''
  searchForm.provider_type = undefined
  fetchProviderList()
}

function handleCreate() {
  currentProvider.value = null
  showModal.value = true
}

function handleEdit(row: any) {
  currentProvider.value = row
  showModal.value = true
}

async function handleSync(row: any) {
  try {
    await syncResources(row.id)
    ElMessage.success('同步任务已启动')
    fetchProviderList()
  } catch (error) {
    console.error('同步失败:', error)
    ElMessage.error('同步失败')
  }
}

async function handleDelete(row: any) {
  try {
    await ElMessageBox.confirm(
      `确定要删除云账号 "${row.name}" 吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await deleteProvider(row.id)
    ElMessage.success('删除成功')
    fetchProviderList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

function handleCommand(command: string) {
  const [action, id] = command.split('-')
  const provider = providerList.value.find((p: any) => p.id === parseInt(id))
  
  if (!provider) return
  
  switch (action) {
    case 'edit':
      handleEdit(provider)
      break
    case 'sync':
      handleSync(provider)
      break
    case 'delete':
      handleDelete(provider)
      break
  }
}

function handleModalSuccess() {
  showModal.value = false
  fetchProviderList()
}

function handleModalCancel() {
  showModal.value = false
}

function formatDate(dateStr: string) {
  if (!dateStr) return '--'
  try {
    return new Date(dateStr).toLocaleString('zh-CN')
  } catch (error) {
    return dateStr
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchProviderList()
})
</script>

<style scoped>
.provider-management {
  padding: 24px;
  background: #f5f7fa;
  min-height: 100vh;
}

/* 页面头部 */
.page-header {
  margin-bottom: 24px;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.3);
}

.header-left {
  flex: 1;
}

.page-title h1 {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 700;
  color: white;
}

.page-title p {
  margin: 0;
  font-size: 16px;
  color: rgba(255, 255, 255, 0.8);
}

.header-right {
  display: flex;
  gap: 12px;
}

/* 统计卡片 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  border: 1px solid #f0f0f0;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.stat-icon.aliyun {
  background: linear-gradient(135deg, #ff6a00, #ff8f00);
}

.stat-icon.tencent {
  background: linear-gradient(135deg, #006eff, #0099ff);
}

.stat-icon.aws {
  background: linear-gradient(135deg, #ff9900, #ffb84d);
}

.stat-icon.total {
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  color: white;
}

.icon-placeholder {
  color: white;
  font-size: 12px;
  font-weight: 600;
  text-align: center;
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  color: #262626;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #8c8c8c;
  font-weight: 500;
}

/* 搜索区域 */
.search-section {
  margin-bottom: 24px;
}

.search-card {
  border-radius: 12px;
  border: 1px solid #f0f0f0;
}

:deep(.el-card__body) {
  padding: 20px 24px;
}

.search-form {
  margin: 0;
}

:deep(.el-form-item) {
  margin-bottom: 0;
  margin-right: 24px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #262626;
}

/* 云账号列表 */
.provider-list {
  margin-bottom: 24px;
}

.empty-card {
  border-radius: 12px;
  border: 1px solid #f0f0f0;
  text-align: center;
  padding: 40px;
}

.provider-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 20px;
}

.provider-card {
  background: white;
  border-radius: 12px;
  border: 1px solid #f0f0f0;
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.provider-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  border-color: #1890ff;
}

.provider-card.provider-aliyun:hover {
  border-color: #ff6a00;
}

.provider-card.provider-tencent:hover {
  border-color: #006eff;
}

.provider-card.provider-aws:hover {
  border-color: #ff9900;
}

/* 卡片头部 */
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px 16px;
  border-bottom: 1px solid #f5f5f5;
}

.provider-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.provider-avatar {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8f9fa;
}

.avatar-text {
  color: white;
  font-size: 16px;
  font-weight: 700;
}

.avatar-aliyun {
  background: linear-gradient(135deg, #ff6a00, #ff8f00);
}

.avatar-tencent {
  background: linear-gradient(135deg, #006eff, #0099ff);
}

.avatar-aws {
  background: linear-gradient(135deg, #ff9900, #ffb84d);
}

.provider-details {
  flex: 1;
}

.provider-name {
  margin: 0 0 4px 0;
  font-size: 16px;
  font-weight: 600;
  color: #262626;
}

.provider-type {
  margin: 0;
  font-size: 14px;
  color: #8c8c8c;
}

.card-actions {
  display: flex;
  align-items: center;
}

.more-btn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 卡片内容 */
.card-content {
  padding: 16px 24px;
}

.info-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.info-item:last-child {
  margin-bottom: 0;
}

.info-label {
  font-size: 14px;
  color: #8c8c8c;
  font-weight: 500;
}

.info-value {
  font-size: 14px;
  color: #262626;
  font-family: 'Monaco', 'Menlo', monospace;
}

/* 卡片底部 */
.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px 20px;
  border-top: 1px solid #f5f5f5;
  background: #fafafa;
}

.status-info {
  display: flex;
  align-items: center;
}

.quick-actions {
  display: flex;
  gap: 8px;
}

/* 加载状态 */
.loading-container {
  background: white;
  border-radius: 12px;
  padding: 24px;
  border: 1px solid #f0f0f0;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .provider-grid {
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  }
}

@media (max-width: 768px) {
  .provider-management {
    padding: 16px;
  }

  .header-content {
    flex-direction: column;
    gap: 20px;
    text-align: center;
    padding: 24px;
  }

  .stats-cards {
    grid-template-columns: 1fr;
  }

  .provider-grid {
    grid-template-columns: 1fr;
  }

  .search-form {
    flex-direction: column;
  }

  :deep(.el-form-item) {
    margin-right: 0;
    margin-bottom: 16px;
  }
}

/* Element Plus 样式覆盖 */
:deep(.el-button--large) {
  padding: 12px 24px;
  font-size: 16px;
  border-radius: 8px;
}

:deep(.el-card) {
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

:deep(.el-empty) {
  padding: 40px 0;
}

:deep(.el-tag) {
  border-radius: 4px;
  font-weight: 500;
}
</style>
