<template>
  <div class="provider-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <div class="title-section">
            <h1>云账号管理</h1>
            <p>管理和监控您的云厂商账号配置</p>
          </div>
        </div>
        <div class="header-right">
          <el-button type="primary" size="large" @click="handleCreate">
            添加云账号
          </el-button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-section">
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon">☁️</div>
          <div class="stat-content">
            <div class="stat-number">0</div>
            <div class="stat-label">阿里云账号</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">☁️</div>
          <div class="stat-content">
            <div class="stat-number">0</div>
            <div class="stat-label">腾讯云账号</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">☁️</div>
          <div class="stat-content">
            <div class="stat-number">0</div>
            <div class="stat-label">AWS账号</div>
          </div>
        </div>
        <div class="stat-card total">
          <div class="stat-number">0</div>
          <div class="stat-label">总计账号</div>
        </div>
      </div>
    </div>

    <!-- 搜索表单 -->
    <div class="search-section">
      <el-form :model="searchForm" inline>
        <el-form-item label="账号名称">
          <el-input
            v-model="searchForm.name"
            placeholder="请输入账号名称"
            clearable
          />
        </el-form-item>
        <el-form-item label="云厂商">
          <el-select
            v-model="searchForm.provider_type"
            placeholder="请选择云厂商"
            clearable
          >
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            搜索
          </el-button>
          <el-button @click="handleReset">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 云账号列表 -->
    <div class="provider-list">
      <el-card v-if="providerList.length === 0 && !loading" class="empty-card" shadow="never">
        <el-empty description="暂无云账号">
          <el-button type="primary" @click="handleCreate">添加第一个云账号</el-button>
        </el-empty>
      </el-card>

      <div v-else class="provider-grid">
        <!-- 这里将来显示云账号列表 -->
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="3" animated />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

// 响应式数据
const loading = ref(false)
const providerList = ref([])

// 搜索表单
const searchForm = reactive({
  name: '',
  provider_type: undefined,
})

// 获取厂商列表
async function fetchProviderList() {
  loading.value = true
  try {
    // 这里将来调用API获取数据
    // const response = await getProviderList(searchForm)
    // providerList.value = response.data || []
    providerList.value = []
  } catch (error) {
    console.error('获取厂商列表失败:', error)
    ElMessage.error('获取厂商列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
function handleSearch() {
  fetchProviderList()
}

// 重置
function handleReset() {
  searchForm.name = ''
  searchForm.provider_type = undefined
  fetchProviderList()
}

// 创建厂商
function handleCreate() {
  ElMessage.info('添加云账号功能正在开发中...')
}

// 页面加载时获取数据
onMounted(() => {
  fetchProviderList()
})
</script>

<style scoped>
.provider-page {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-section h1 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
}

.title-section p {
  margin: 0;
  color: #6b7280;
  font-size: 14px;
}

.stats-section {
  margin-bottom: 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.stat-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-card.total {
  background: #f3f4f6;
}

.stat-icon {
  font-size: 32px;
}

.stat-number {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
}

.search-section {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 24px;
}

.provider-list {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 20px;
}

.empty-card {
  border: none;
  box-shadow: none;
}

.loading-container {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 20px;
}
</style>
