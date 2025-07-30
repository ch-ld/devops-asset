<template>
  <div class="announcement-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <h2 class="page-title">
            <i class="iconfont-sys">&#xe6d1;</i>
            公告管理
          </h2>
          <p class="page-description">管理系统公告信息，支持发布、编辑和删除操作</p>
        </div>
        <div class="header-actions">
          <el-button type="primary" size="large" @click="handleCreate" class="create-btn">
            <i class="iconfont-sys">&#xe63f;</i>
            新增公告
          </el-button>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选区域 -->
    <el-card class="search-card" shadow="never">
      <div class="search-header">
        <h3 class="search-title">筛选条件</h3>
        <el-button text @click="handleReset" class="reset-btn">
          <i class="iconfont-sys">&#xe643;</i>
          重置
        </el-button>
      </div>
      
      <el-form :model="searchForm" class="search-form">
        <el-row :gutter="24">
          <el-col :span="6">
            <el-form-item label="公告标题">
              <el-input
                v-model="searchForm.title"
                placeholder="请输入公告标题"
                clearable
              >
                <template #prefix>
                  <i class="iconfont-sys">&#xe654;</i>
                </template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="公告类型">
              <el-select v-model="searchForm.type" placeholder="请选择类型" clearable>
                <el-option label="通知" :value="1" />
                <el-option label="公告" :value="2" />
                <el-option label="系统" :value="3" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="状态">
              <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
                <el-option label="启用" :value="1" />
                <el-option label="禁用" :value="2" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item>
              <el-button type="primary" @click="handleSearch" class="search-btn">
                <i class="iconfont-sys">&#xe654;</i>
                搜索
              </el-button>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 数据统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon total">
          <i class="iconfont-sys">&#xe6d1;</i>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ pagination.total }}</div>
          <div class="stat-label">总公告数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon active">
          <i class="iconfont-sys">&#xe61f;</i>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ activeCount }}</div>
          <div class="stat-label">启用公告</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon pinned">
          <i class="iconfont-sys">&#xe6d2;</i>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ pinnedCount }}</div>
          <div class="stat-label">置顶公告</div>
        </div>
      </div>
    </div>

    <!-- 公告列表 -->
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="table-header">
          <h3 class="table-title">公告列表</h3>
          <div class="table-actions">
            <el-button-group>
              <el-button :type="viewMode === 'table' ? 'primary' : ''" @click="viewMode = 'table'">
                <i class="iconfont-sys">&#xe655;</i>
                表格视图
              </el-button>
              <el-button :type="viewMode === 'card' ? 'primary' : ''" @click="viewMode = 'card'">
                <i class="iconfont-sys">&#xe656;</i>
                卡片视图
              </el-button>
            </el-button-group>
          </div>
        </div>
      </template>

      <!-- 表格视图 -->
      <el-table 
        v-if="viewMode === 'table'"
        :data="announcementList" 
        v-loading="loading"
        class="modern-table"
        stripe
        :header-cell-style="{ backgroundColor: '#fafafa', color: '#606266' }"
      >
        <el-table-column prop="title" label="标题" min-width="200">
          <template #default="{ row }">
            <div class="announcement-title-cell">
              <div class="title-content">
                <span class="title-text">{{ row.title }}</span>
                <div class="title-badges">
                  <el-tag v-if="row.is_pinned === 1" type="warning" size="small" effect="light">
                    <i class="el-icon-top"></i>
                    置顶
                  </el-tag>
                  <el-tag :type="getTypeTagType(row.type)" size="small" effect="light">
                    {{ getTypeLabel(row.type) }}
                  </el-tag>
                </div>
              </div>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" effect="light">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="priority" label="优先级" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.priority > 0" type="warning" effect="plain">
              {{ row.priority }}
            </el-tag>
            <span v-else class="text-muted">{{ row.priority }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="read_count" label="阅读数" width="100" align="center">
          <template #default="{ row }">
            <div class="read-count">
              <i class="el-icon-view"></i>
              {{ row.read_count || 0 }}
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="publisher" label="发布者" width="120">
          <template #default="{ row }">
            <div class="publisher-info">
              <el-avatar :size="24" :src="row.publisher?.avatar">
                {{ row.publisher?.name?.charAt(0) }}
              </el-avatar>
              <span class="publisher-name">{{ row.publisher?.name || '-' }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            <div class="time-info">
              <i class="el-icon-time"></i>
              {{ formatTime(row.created_at) }}
            </div>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button 
                type="primary" 
                size="small" 
                @click="handleView(row)"
                link
              >
                <i class="el-icon-view"></i>
                查看
              </el-button>
              <el-button 
                type="primary" 
                size="small" 
                @click="handleEdit(row)"
                link
              >
                <i class="el-icon-edit"></i>
                编辑
              </el-button>
              <el-button 
                type="danger" 
                size="small" 
                @click="handleDelete(row)"
                link
              >
                <i class="el-icon-delete"></i>
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 卡片视图 -->
      <div v-else class="card-view" v-loading="loading">
        <div class="announcement-cards">
          <div 
            v-for="item in announcementList" 
            :key="item.id"
            class="announcement-card"
            @click="handleView(item)"
          >
            <div class="card-header">
              <div class="card-badges">
                <el-tag :type="getTypeTagType(item.type)" size="small" effect="light">
                  {{ getTypeLabel(item.type) }}
                </el-tag>
                <el-tag v-if="item.is_pinned === 1" type="warning" size="small" effect="light">
                  <i class="el-icon-top"></i>
                  置顶
                </el-tag>
                <el-tag :type="item.status === 1 ? 'success' : 'danger'" size="small" effect="light">
                  {{ item.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </div>
              <div class="card-actions">
                <el-dropdown trigger="click">
                  <el-button type="text" size="small">
                    <i class="el-icon-more"></i>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item @click.stop="handleEdit(item)">
                        <i class="el-icon-edit"></i>
                        编辑
                      </el-dropdown-item>
                      <el-dropdown-item @click.stop="handleDelete(item)" divided>
                        <i class="el-icon-delete"></i>
                        删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>
            
            <div class="card-content">
              <h4 class="card-title">{{ item.title }}</h4>
              <p class="card-description">{{ item.content }}</p>
            </div>
            
            <div class="card-footer">
              <div class="card-meta">
                <div class="meta-item">
                  <el-avatar :size="20" :src="item.publisher?.avatar">
                    {{ item.publisher?.name?.charAt(0) }}
                  </el-avatar>
                  <span>{{ item.publisher?.name }}</span>
                </div>
                <div class="meta-item">
                  <i class="el-icon-time"></i>
                  <span>{{ formatTime(item.created_at) }}</span>
                </div>
                <div class="meta-item">
                  <i class="el-icon-view"></i>
                  <span>{{ item.read_count || 0 }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
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

    <!-- 创建/编辑对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="isEdit ? '编辑公告' : '创建公告'"
      width="900px"
      @closed="handleDialogClosed"
      class="modern-dialog"
    >
      <el-form 
        :model="form" 
        :rules="rules" 
        ref="formRef" 
        label-width="100px"
        class="modern-form"
      >
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="标题" prop="title">
              <el-input v-model="form.title" placeholder="请输入公告标题" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="类型" prop="type">
              <el-select v-model="form.type" placeholder="请选择类型" style="width: 100%">
                <el-option label="通知" :value="1" />
                <el-option label="公告" :value="2" />
                <el-option label="系统" :value="3" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="24">
          <el-col :span="8">
            <el-form-item label="状态" prop="status">
              <el-select v-model="form.status" placeholder="请选择状态" style="width: 100%">
                <el-option label="启用" :value="1" />
                <el-option label="禁用" :value="2" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="是否置顶" prop="is_pinned">
              <el-switch 
                v-model="form.is_pinned" 
                :active-value="1" 
                :inactive-value="2"
                active-text="置顶" 
                inactive-text="不置顶" 
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="优先级" prop="priority">
              <el-input-number 
                v-model="form.priority" 
                :min="0" 
                :max="999" 
                placeholder="数字越大优先级越高"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="发布时间" prop="publish_time">
              <el-date-picker
                v-model="form.publish_time"
                type="datetime"
                placeholder="选择发布时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="过期时间" prop="expire_time">
              <el-date-picker
                v-model="form.expire_time"
                type="datetime"
                placeholder="选择过期时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="内容" prop="content">
          <el-input 
            v-model="form.content" 
            type="textarea" 
            :rows="6"
            placeholder="请输入公告内容"
            show-word-limit
            maxlength="1000"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
            {{ isEdit ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 查看对话框 -->
    <el-dialog 
      v-model="viewDialogVisible" 
      :title="currentAnnouncement?.title"
      width="700px"
      class="view-dialog"
    >
      <div v-if="currentAnnouncement" class="announcement-detail">
        <div class="announcement-meta">
          <div class="meta-badges">
            <el-tag :type="getTypeTagType(currentAnnouncement.type)" effect="light">
              {{ getTypeLabel(currentAnnouncement.type) }}
            </el-tag>
            <el-tag v-if="currentAnnouncement.is_pinned === 1" type="warning" effect="light">
              <i class="el-icon-top"></i>
              置顶
            </el-tag>
            <el-tag :type="currentAnnouncement.status === 1 ? 'success' : 'danger'" effect="light">
              {{ currentAnnouncement.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </div>
          <div class="meta-info">
            <div class="info-item">
              <i class="el-icon-time"></i>
              <span>发布时间：{{ formatTime(currentAnnouncement.created_at) }}</span>
            </div>
            <div v-if="currentAnnouncement.publisher" class="info-item">
              <i class="el-icon-user"></i>
              <span>发布者：{{ currentAnnouncement.publisher.name }}</span>
            </div>
            <div class="info-item">
              <i class="el-icon-view"></i>
              <span>阅读数：{{ currentAnnouncement.read_count || 0 }}</span>
            </div>
            <div class="info-item">
              <i class="el-icon-sort"></i>
              <span>优先级：{{ currentAnnouncement.priority || 0 }}</span>
            </div>
          </div>
        </div>
        <div class="announcement-content">
          {{ currentAnnouncement.content }}
        </div>
      </div>
      
      <template #footer>
        <el-button type="primary" @click="viewDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, FormInstance } from 'element-plus'
import { 
  getAnnouncementList, 
  createAnnouncement, 
  updateAnnouncement, 
  deleteAnnouncement,
  getAnnouncementDetail 
} from '@/api/system/announcement'

defineOptions({ name: 'AnnouncementManagement' })

// 响应式数据
const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const viewDialogVisible = ref(false)
const isEdit = ref(false)
const announcementList = ref([])
const currentAnnouncement = ref(null)
const viewMode = ref('table')

// 计算统计数据
const activeCount = computed(() => {
  return announcementList.value.filter(item => item.status === 1).length
})

const pinnedCount = computed(() => {
  return announcementList.value.filter(item => item.is_pinned === 1).length
})

// 搜索表单
const searchForm = reactive({
  title: '',
  type: undefined,
  status: undefined
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 表单数据
const form = reactive({
  title: '',
  content: '',
  type: 1,
  status: 1,
  is_pinned: 2,
  priority: 0,
  publish_time: '',
  expire_time: ''
})

// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入公告标题', trigger: 'blur' },
    { min: 1, max: 100, message: '标题长度应在1-100字符之间', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入公告内容', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择公告类型', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择公告状态', trigger: 'change' }
  ]
}

const formRef = ref<FormInstance>()

// 获取公告列表
const fetchAnnouncementList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...searchForm
    }
    
    const response = await getAnnouncementList(params)
    if (response && response.code === 200) {
      announcementList.value = response.data || []
      pagination.total = response.count || 0
    } else {
      ElMessage.error('获取公告列表失败')
    }
  } catch (error) {
    console.error('获取公告列表失败:', error)
    ElMessage.error('获取公告列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchAnnouncementList()
}

// 重置搜索
const handleReset = () => {
  searchForm.title = ''
  searchForm.type = undefined
  searchForm.status = undefined
  pagination.page = 1
  fetchAnnouncementList()
}

// 分页相关
const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  fetchAnnouncementList()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchAnnouncementList()
}

// 创建公告
const handleCreate = () => {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

// 编辑公告
const handleEdit = (row: any) => {
  isEdit.value = true
  Object.assign(form, {
    id: row.id,
    title: row.title,
    content: row.content,
    type: row.type,
    status: row.status,
    is_pinned: row.is_pinned,
    priority: row.priority,
    publish_time: row.publish_time ? (typeof row.publish_time === 'number' ? new Date(row.publish_time * 1000).toISOString().slice(0, 19) : row.publish_time) : '',
    expire_time: row.expire_time ? (typeof row.expire_time === 'number' ? new Date(row.expire_time * 1000).toISOString().slice(0, 19) : row.expire_time) : ''
  })
  dialogVisible.value = true
}

// 查看公告
const handleView = async (row: any) => {
  try {
    const response = await getAnnouncementDetail(row.id)
    if (response && response.code === 200) {
      currentAnnouncement.value = response.data
      viewDialogVisible.value = true
    } else {
      ElMessage.error('获取公告详情失败')
    }
  } catch (error) {
    console.error('获取公告详情失败:', error)
    ElMessage.error('获取公告详情失败')
  }
}

// 删除公告
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除公告"${row.title}"吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const response = await deleteAnnouncement(row.id)
    if (response && response.code === 200) {
      ElMessage.success('删除成功')
      fetchAnnouncementList()
    } else {
      ElMessage.error(response?.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除公告失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        const data = { ...form }
        
        if (isEdit.value) {
          const response = await updateAnnouncement(data.id, data)
          if (response && response.code === 200) {
            ElMessage.success('更新成功')
            dialogVisible.value = false
            fetchAnnouncementList()
          } else {
            ElMessage.error(response?.message || '更新失败')
          }
        } else {
          const response = await createAnnouncement(data)
          if (response && response.code === 200) {
            ElMessage.success('创建成功')
            dialogVisible.value = false
            fetchAnnouncementList()
          } else {
            ElMessage.error(response?.message || '创建失败')
          }
        }
      } catch (error) {
        console.error('提交失败:', error)
        ElMessage.error('提交失败')
      } finally {
        submitLoading.value = false
      }
    }
  })
}

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    title: '',
    content: '',
    type: 1,
    status: 1,
    is_pinned: 2,
    priority: 0,
    publish_time: '',
    expire_time: ''
  })
}

// 对话框关闭事件
const handleDialogClosed = () => {
  resetForm()
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

// 获取类型标签
const getTypeLabel = (type: number) => {
  const typeMap = {
    1: '通知',
    2: '公告',
    3: '系统'
  }
  return typeMap[type] || '未知'
}

// 获取类型标签类型
const getTypeTagType = (type: number) => {
  const typeMap = {
    1: 'info',
    2: 'warning',
    3: 'danger'
  }
  return typeMap[type] || 'info'
}

// 格式化时间
const formatTime = (timestamp: number | string) => {
  if (!timestamp) return '-'
  // 如果是Unix时间戳（数字），需要转换为毫秒
  const date = typeof timestamp === 'number' ? new Date(timestamp * 1000) : new Date(timestamp)
  return date.toLocaleString('zh-CN')
}

onMounted(() => {
  fetchAnnouncementList()
})
</script>

<style scoped lang="scss">
.announcement-management {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 24px;
  position: relative;
  overflow: hidden;
  
  // 添加背景装饰元素
  &::before {
    content: '';
    position: absolute;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background: radial-gradient(circle, rgba(255, 255, 255, 0.05) 1px, transparent 1px);
    background-size: 50px 50px;
    animation: float 20s infinite linear;
    pointer-events: none;
  }
  
  @keyframes float {
    0% { transform: translate(0, 0) rotate(0deg); }
    100% { transform: translate(-50px, -50px) rotate(360deg); }
  }
  
  // 页面头部
  .page-header {
    margin-bottom: 32px;
    position: relative;
    z-index: 1;
    
    .header-content {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      background: rgba(255, 255, 255, 0.95);
      backdrop-filter: blur(20px);
      border-radius: 20px;
      padding: 32px 40px;
      box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1), 0 8px 25px rgba(102, 126, 234, 0.15);
      border: 1px solid rgba(255, 255, 255, 0.3);
      position: relative;
      overflow: hidden;
      
      &::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 2px;
        background: linear-gradient(90deg, #667eea, #764ba2, #f093fb, #f5576c);
        background-size: 200% 100%;
        animation: gradient-flow 3s ease-in-out infinite;
      }
      
      @keyframes gradient-flow {
        0%, 100% { background-position: 0% 50%; }
        50% { background-position: 100% 50%; }
      }
    }
    
    .header-left {
      .page-title {
        font-size: 32px;
        font-weight: 800;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        margin: 0 0 12px 0;
        display: flex;
        align-items: center;
        gap: 16px;
        letter-spacing: -0.5px;
        
        .iconfont-sys {
          font-size: 36px;
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
          animation: icon-pulse 2s ease-in-out infinite;
        }
        
        @keyframes icon-pulse {
          0%, 100% { transform: scale(1); }
          50% { transform: scale(1.05); }
        }
      }
      
      .page-description {
        color: #555;
        font-size: 16px;
        margin: 0;
        font-weight: 500;
        letter-spacing: 0.2px;
      }
    }
    
    .header-actions {
      .create-btn {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        border: none;
        border-radius: 16px;
        padding: 16px 32px;
        font-weight: 700;
        font-size: 16px;
        letter-spacing: 0.5px;
        box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3), 0 4px 15px rgba(0, 0, 0, 0.1);
        transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
        position: relative;
        overflow: hidden;
        
        &::before {
          content: '';
          position: absolute;
          top: 0;
          left: -100%;
          width: 100%;
          height: 100%;
          background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
          transition: left 0.5s;
        }
        
        &:hover {
          transform: translateY(-4px) scale(1.02);
          box-shadow: 0 15px 40px rgba(102, 126, 234, 0.4), 0 8px 25px rgba(0, 0, 0, 0.15);
          
          &::before {
            left: 100%;
          }
        }
        
        &:active {
          transform: translateY(-2px) scale(0.98);
        }
        
        i {
          margin-right: 8px;
          font-size: 18px;
        }
      }
    }
  }
  
  // 搜索卡片
  .search-card {
    margin-bottom: 32px;
    border-radius: 20px;
    border: none;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.08), 0 4px 20px rgba(102, 126, 234, 0.1);
    overflow: hidden;
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(10px);
    position: relative;
    z-index: 1;
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 15px 50px rgba(0, 0, 0, 0.12), 0 8px 30px rgba(102, 126, 234, 0.15);
    }
    
    :deep(.el-card__body) {
      padding: 32px;
    }
    
    .search-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;
      
      .search-title {
        font-size: 18px;
        font-weight: 600;
        color: #1a1a1a;
        margin: 0;
      }
      
      .reset-btn {
        color: #667eea;
        font-weight: 500;
        
        &:hover {
          background: rgba(102, 126, 234, 0.1);
        }
      }
    }
    
    .search-form {
      :deep(.el-form-item__label) {
        font-weight: 500;
        color: #333;
      }
      
      :deep(.el-input__wrapper) {
        border-radius: 8px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
        transition: all 0.3s ease;
        
        &:hover, &.is-focus {
          box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
        }
      }
      
      :deep(.el-select) {
        width: 100%;
      }
      
      .search-btn {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        border: none;
        border-radius: 8px;
        font-weight: 500;
      }
    }
  }
  
  // 统计卡片网格
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 24px;
    margin-bottom: 32px;
    position: relative;
    z-index: 1;
    
    .stat-card {
      background: rgba(255, 255, 255, 0.98);
      backdrop-filter: blur(15px);
      border-radius: 20px;
      padding: 32px 28px;
      display: flex;
      align-items: center;
      gap: 20px;
      box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08), 0 4px 16px rgba(102, 126, 234, 0.1);
      border: 1px solid rgba(255, 255, 255, 0.3);
      transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
      position: relative;
      overflow: hidden;
      
      &::before {
        content: '';
        position: absolute;
        top: 0;
        left: -100%;
        width: 100%;
        height: 2px;
        transition: left 0.5s ease;
      }
      
      &:hover {
        transform: translateY(-8px) scale(1.02);
        box-shadow: 0 20px 60px rgba(0, 0, 0, 0.12), 0 8px 30px rgba(102, 126, 234, 0.2);
        
        &::before {
          left: 100%;
        }
        
        .stat-icon {
          transform: scale(1.1) rotate(5deg);
        }
        
        .stat-value {
          transform: scale(1.05);
        }
      }
      
      .stat-icon {
        width: 56px;
        height: 56px;
        border-radius: 16px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 28px;
        color: white;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
        transition: all 0.3s ease;
        
        &.total {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          
          &:hover {
            box-shadow: 0 8px 30px rgba(102, 126, 234, 0.4);
          }
          
          &::before {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          }
        }
        
        &.active {
          background: linear-gradient(135deg, #56CC9D 0%, #6FCEB7 100%);
          
          &:hover {
            box-shadow: 0 8px 30px rgba(86, 204, 157, 0.4);
          }
          
          &::before {
            background: linear-gradient(135deg, #56CC9D 0%, #6FCEB7 100%);
          }
        }
        
        &.pinned {
          background: linear-gradient(135deg, #FFB74D 0%, #FF8A65 100%);
          
          &:hover {
            box-shadow: 0 8px 30px rgba(255, 183, 77, 0.4);
          }
          
          &::before {
            background: linear-gradient(135deg, #FFB74D 0%, #FF8A65 100%);
          }
        }
      }
      
      .stat-content {
        flex: 1;
        
        .stat-value {
          font-size: 32px;
          font-weight: 800;
          background: linear-gradient(135deg, #1a1a1a 0%, #333 100%);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
          line-height: 1;
          margin-bottom: 8px;
          transition: all 0.3s ease;
          letter-spacing: -1px;
        }
        
        .stat-label {
          font-size: 15px;
          color: #555;
          font-weight: 600;
          text-transform: uppercase;
          letter-spacing: 0.5px;
        }
      }
    }
  }
  
  // 表格卡片
  .table-card {
    border-radius: 20px;
    border: none;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.08), 0 4px 20px rgba(102, 126, 234, 0.1);
    overflow: hidden;
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(10px);
    position: relative;
    z-index: 1;
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 15px 50px rgba(0, 0, 0, 0.12), 0 8px 30px rgba(102, 126, 234, 0.15);
    }
    
    :deep(.el-card__header) {
      background: linear-gradient(135deg, #f8f9ff 0%, #f0f4ff 100%);
      border-bottom: 1px solid rgba(102, 126, 234, 0.1);
      padding: 24px 32px;
      position: relative;
      
      &::before {
        content: '';
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        height: 2px;
        background: linear-gradient(90deg, #667eea, #764ba2);
        opacity: 0.3;
      }
    }
    
    :deep(.el-card__body) {
      padding: 0;
    }
    
    .table-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .table-title {
        font-size: 20px;
        font-weight: 700;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        margin: 0;
        letter-spacing: -0.3px;
      }
      
      .table-actions {
        :deep(.el-button-group) {
          .el-button {
            border-radius: 8px;
            
            &:first-child {
              border-top-right-radius: 0;
              border-bottom-right-radius: 0;
            }
            
            &:last-child {
              border-top-left-radius: 0;
              border-bottom-left-radius: 0;
            }
          }
        }
      }
    }
  }
  
  // 现代表格样式
  .modern-table {
    :deep(.el-table__header) {
      th {
        border-bottom: 2px solid rgba(102, 126, 234, 0.1);
        font-weight: 700;
        background: linear-gradient(135deg, #f8f9ff 0%, #f0f4ff 100%) !important;
        
        .cell {
          padding: 20px 16px;
          color: #333;
          font-size: 14px;
          letter-spacing: 0.3px;
        }
      }
    }
    
    :deep(.el-table__body) {
      tr {
        transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
        
        &:hover {
          background: linear-gradient(135deg, #f8f9ff 0%, #f0f4ff 50%) !important;
          transform: translateX(4px);
          box-shadow: 4px 0 12px rgba(102, 126, 234, 0.1);
        }
        
        td {
          border-bottom: 1px solid rgba(0, 0, 0, 0.05);
          
          .cell {
            padding: 20px 16px;
          }
        }
      }
    }
    
    .announcement-title-cell {
      .title-content {
        .title-text {
          font-weight: 500;
          color: #333;
          margin-bottom: 4px;
          display: block;
        }
        
        .title-badges {
          display: flex;
          gap: 4px;
          flex-wrap: wrap;
        }
      }
    }
    
    .publisher-info {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .publisher-name {
        font-size: 13px;
        color: #666;
      }
    }
    
    .time-info {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 13px;
      color: #666;
      
      i {
        font-size: 12px;
        opacity: 0.7;
      }
    }
    
    .read-count {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 13px;
      color: #666;
      
      i {
        font-size: 12px;
        opacity: 0.7;
      }
    }
    
    .action-buttons {
      display: flex;
      gap: 8px;
      
      .el-button {
        &.is-link {
          padding: 4px 8px;
          font-size: 12px;
        }
      }
    }
  }
  
  // 卡片视图
  .card-view {
    padding: 24px;
    
    .announcement-cards {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
      gap: 20px;
      
      .announcement-card {
        background: white;
        border-radius: 12px;
        overflow: hidden;
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
        transition: all 0.3s ease;
        cursor: pointer;
        
        &:hover {
          transform: translateY(-4px);
          box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
        }
        
        .card-header {
          display: flex;
          justify-content: space-between;
          align-items: flex-start;
          padding: 16px 20px 12px;
          
          .card-badges {
            display: flex;
            gap: 6px;
            flex-wrap: wrap;
          }
        }
        
        .card-content {
          padding: 0 20px 16px;
          
          .card-title {
            font-size: 16px;
            font-weight: 600;
            color: #333;
            margin: 0 0 8px 0;
            line-height: 1.4;
          }
          
          .card-description {
            font-size: 14px;
            color: #666;
            line-height: 1.5;
            margin: 0;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
            overflow: hidden;
          }
        }
        
        .card-footer {
          padding: 12px 20px 16px;
          border-top: 1px solid #f0f0f0;
          
          .card-meta {
            display: flex;
            align-items: center;
            gap: 16px;
            font-size: 12px;
            color: #999;
            
            .meta-item {
              display: flex;
              align-items: center;
              gap: 4px;
              
              span {
                white-space: nowrap;
              }
            }
          }
        }
      }
    }
  }
  
  // 分页
  .pagination-wrapper {
    margin-top: 32px;
    padding: 24px 32px;
    display: flex;
    justify-content: center;
    background: linear-gradient(135deg, #f8f9ff 0%, #f0f4ff 100%);
    border-top: 1px solid rgba(102, 126, 234, 0.1);
    
    :deep(.el-pagination) {
      .el-pagination__goto,
      .el-pagination__classifier,
      .el-pagination__total {
        color: #555;
        font-weight: 500;
      }
      
      .btn-prev,
      .btn-next,
      .el-pager li {
        background: rgba(255, 255, 255, 0.8);
        border: 1px solid rgba(102, 126, 234, 0.1);
        color: #333;
        font-weight: 500;
        transition: all 0.3s ease;
        
        &:hover {
          background: rgba(102, 126, 234, 0.1);
          border-color: rgba(102, 126, 234, 0.2);
          transform: translateY(-1px);
        }
        
        &.is-active {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          border-color: transparent;
          color: white;
          box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
        }
      }
    }
  }
  
  // 对话框样式
  :deep(.modern-dialog) {
    .el-dialog {
      border-radius: 20px;
      overflow: hidden;
      box-shadow: 0 25px 80px rgba(0, 0, 0, 0.15), 0 10px 40px rgba(102, 126, 234, 0.1);
      border: 1px solid rgba(255, 255, 255, 0.2);
    }
    
    .el-dialog__header {
      background: linear-gradient(135deg, #f8f9ff 0%, #f0f4ff 100%);
      padding: 24px 32px;
      border-bottom: 1px solid rgba(102, 126, 234, 0.1);
      position: relative;
      
      &::before {
        content: '';
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        height: 2px;
        background: linear-gradient(90deg, #667eea, #764ba2);
        opacity: 0.3;
      }
      
      .el-dialog__title {
        font-weight: 700;
        font-size: 20px;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        letter-spacing: -0.3px;
      }
    }
    
    .el-dialog__body {
      padding: 32px;
      background: rgba(255, 255, 255, 0.95);
    }
    
    .el-dialog__footer {
      padding: 20px 32px;
      border-top: 1px solid rgba(102, 126, 234, 0.1);
      background: linear-gradient(135deg, #f8f9ff 0%, #f0f4ff 100%);
    }
  }
  
  .modern-form {
    :deep(.el-form-item__label) {
      font-weight: 600;
      color: #333;
      font-size: 14px;
      letter-spacing: 0.3px;
    }
    
    :deep(.el-input__wrapper),
    :deep(.el-textarea__inner) {
      border-radius: 12px;
      border: 1px solid rgba(102, 126, 234, 0.2);
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
      transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
      
      &:hover {
        border-color: rgba(102, 126, 234, 0.4);
        box-shadow: 0 4px 15px rgba(102, 126, 234, 0.1);
        transform: translateY(-1px);
      }
      
      &.is-focus {
        border-color: #667eea;
        box-shadow: 0 6px 20px rgba(102, 126, 234, 0.2);
        transform: translateY(-2px);
      }
    }
    
    :deep(.el-select) {
      width: 100%;
      
      .el-select__wrapper {
        border-radius: 12px;
        border: 1px solid rgba(102, 126, 234, 0.2);
        transition: all 0.3s ease;
        
        &:hover {
          border-color: rgba(102, 126, 234, 0.4);
          box-shadow: 0 4px 15px rgba(102, 126, 234, 0.1);
        }
      }
    }
    
    :deep(.el-switch) {
      .el-switch__core {
        border-radius: 20px;
        
        &::after {
          border-radius: 50%;
        }
      }
      
      &.is-checked .el-switch__core {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      }
    }
    
    :deep(.el-date-editor) {
      .el-input__wrapper {
        border-radius: 12px;
      }
    }
  }
  
  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 16px;
    
    .el-button {
      border-radius: 12px;
      padding: 12px 24px;
      font-weight: 600;
      letter-spacing: 0.3px;
      transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
      
      &.el-button--primary {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        border: none;
        box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
        
        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
        }
        
        &:active {
          transform: translateY(0);
        }
      }
      
      &.el-button--default {
        border: 1px solid rgba(102, 126, 234, 0.2);
        color: #555;
        background: rgba(255, 255, 255, 0.8);
        
        &:hover {
          background: rgba(102, 126, 234, 0.05);
          border-color: rgba(102, 126, 234, 0.3);
          transform: translateY(-1px);
        }
      }
    }
  }
  
  // 查看对话框
  :deep(.view-dialog) {
    .el-dialog {
      border-radius: 16px;
    }
  }
  
  .announcement-detail {
    .announcement-meta {
      margin-bottom: 20px;
      
      .meta-badges {
        display: flex;
        gap: 8px;
        margin-bottom: 16px;
        flex-wrap: wrap;
      }
      
      .meta-info {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 12px;
        
        .info-item {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 14px;
          color: #666;
          
          i {
            font-size: 14px;
            opacity: 0.7;
          }
        }
      }
    }
    
    .announcement-content {
      background: #f8f9fa;
      border-radius: 8px;
      padding: 20px;
      line-height: 1.6;
      color: #333;
      white-space: pre-wrap;
      border-left: 4px solid #667eea;
    }
  }
  
  // 响应式设计
  @media (max-width: 768px) {
    padding: 16px;
    
    .page-header .header-content {
      flex-direction: column;
      gap: 16px;
      align-items: stretch;
      text-align: center;
    }
    
    .stats-grid {
      grid-template-columns: 1fr;
    }
    
    .search-form {
      :deep(.el-row) {
        flex-direction: column;
        
        .el-col {
          width: 100% !important;
          max-width: 100% !important;
        }
      }
    }
    
    .card-view .announcement-cards {
      grid-template-columns: 1fr;
    }
  }
  
  // 辅助类
  .text-muted {
    color: #999 !important;
  }
}
</style> 
