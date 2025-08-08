<template>
  <div class="notification-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title">
          <h1>通知管理</h1>
          <p>管理系统通知消息，包括证书过期提醒、监控告警等</p>
        </div>
        <div class="page-actions">
          <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="item">
            <el-button @click="handleMarkAllRead" :disabled="unreadCount === 0">
              <el-icon><Check /></el-icon>
              全部已读
            </el-button>
          </el-badge>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            创建通知
          </el-button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-container">
      <el-row :gutter="16">
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="总通知数"
              :value="statistics.total"
              :value-style="{ color: '#409eff' }"
            >
              <template #suffix>
                <el-icon><Bell /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="未读通知"
              :value="statistics.unread"
              :value-style="{ color: '#f56c6c' }"
            >
              <template #suffix>
                <el-icon><Warning /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="今日通知"
              :value="statistics.today"
              :value-style="{ color: '#67c23a' }"
            >
              <template #suffix>
                <el-icon><Calendar /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <el-statistic
              title="告警通知"
              :value="statistics.alert"
              :value-style="{ color: '#e6a23c' }"
            >
              <template #suffix>
                <el-icon><AlarmClock /></el-icon>
              </template>
            </el-statistic>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 快速筛选 -->
    <el-card class="filter-card">
      <el-row :gutter="16">
        <el-col :span="18">
          <el-space wrap>
            <el-button-group>
              <el-button 
                :type="activeFilter === 'all' ? 'primary' : 'default'"
                @click="setFilter('all')"
              >
                全部 ({{ statistics.total }})
              </el-button>
              <el-button 
                :type="activeFilter === 'unread' ? 'primary' : 'default'"
                @click="setFilter('unread')"
              >
                未读 ({{ statistics.unread }})
              </el-button>
              <el-button 
                :type="activeFilter === 'alert' ? 'primary' : 'default'"
                @click="setFilter('alert')"
              >
                告警 ({{ statistics.alert }})
              </el-button>
              <el-button 
                :type="activeFilter === 'cert' ? 'primary' : 'default'"
                @click="setFilter('cert')"
              >
                证书 ({{ statistics.cert }})
              </el-button>
            </el-button-group>

            <el-select 
              v-model="searchForm.level" 
              placeholder="通知级别"
              clearable
              style="width: 120px"
            >
              <el-option label="信息" value="info" />
              <el-option label="警告" value="warning" />
              <el-option label="错误" value="error" />
              <el-option label="成功" value="success" />
            </el-select>

            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              style="width: 240px"
            />
          </el-space>
        </el-col>
        <el-col :span="6">
          <div class="search-input">
            <el-input
              v-model="searchForm.keyword"
              placeholder="搜索通知内容"
              clearable
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 通知列表 -->
    <el-card class="notification-list-card">
      <template #header>
        <div class="list-header">
          <span>通知列表</span>
          <div class="list-actions">
            <el-button @click="handleRefresh">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button 
              type="danger" 
              @click="handleBatchDelete" 
              :disabled="selectedRowKeys.length === 0"
            >
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="notificationList"
        @selection-change="handleSelectionChange"
        stripe
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-icon 
              :color="row.isRead ? '#909399' : '#409eff'" 
              :size="16"
            >
              <CircleClose v-if="!row.isRead" />
              <CircleCheck v-else />
            </el-icon>
          </template>
        </el-table-column>

        <el-table-column label="级别" width="80">
          <template #default="{ row }">
            <el-tag :type="getLevelType(row.level)" size="small">
              {{ getLevelText(row.level) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="title" label="标题" min-width="200">
          <template #default="{ row }">
            <div class="notification-title" @click="handleView(row)">
              <el-link type="primary" :underline="false">{{ row.title }}</el-link>
              <el-tag v-if="row.category" type="info" size="small" class="category-tag">
                {{ getCategoryText(row.category) }}
              </el-tag>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="content" label="内容" show-overflow-tooltip min-width="300">
          <template #default="{ row }">
            <div class="notification-content">
              {{ row.content }}
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="createdAt" label="时间" width="160">
          <template #default="{ row }">
            <div class="notification-time">
              {{ formatDateTime(row.createdAt) }}
            </div>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button text @click="handleView(row)">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button 
              v-if="!row.isRead" 
              text 
              @click="handleMarkRead(row)"
            >
              <el-icon><Check /></el-icon>
              已读
            </el-button>
            <el-button text @click="handleDelete(row)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handlePageSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 通知详情弹窗 -->
    <el-dialog
      v-model="detailModalVisible"
      title="通知详情"
      width="600px"
    >
      <div v-if="currentNotification" class="notification-detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="标题">
            {{ currentNotification.title }}
          </el-descriptions-item>
          <el-descriptions-item label="级别">
            <el-tag :type="getLevelType(currentNotification.level)">
              {{ getLevelText(currentNotification.level) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="类别">
            {{ getCategoryText(currentNotification.category) }}
          </el-descriptions-item>
          <el-descriptions-item label="内容">
            <div class="notification-content-detail">
              {{ currentNotification.content }}
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatDateTime(currentNotification.createdAt) }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="currentNotification.isRead ? 'success' : 'warning'">
              {{ currentNotification.isRead ? '已读' : '未读' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="detailModalVisible = false">关闭</el-button>
          <el-button 
            v-if="currentNotification && !currentNotification.isRead" 
            type="primary" 
            @click="handleMarkRead(currentNotification)"
          >
            标记为已读
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 创建通知弹窗 -->
    <NotificationModal
      v-model:visible="notificationModalVisible"
      @success="handleNotificationSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, 
  Search, 
  Refresh, 
  Delete,
  Check,
  View,
  Bell,
  Warning,
  Calendar,
  AlarmClock,
  CircleCheck,
  CircleClose
} from '@element-plus/icons-vue'
import NotificationModal from './components/NotificationModal.vue'

// 响应式数据
const loading = ref(false)
const notificationList = ref<any[]>([])
const selectedRowKeys = ref<number[]>([])
const detailModalVisible = ref(false)
const notificationModalVisible = ref(false)
const currentNotification = ref<any>(null)
const activeFilter = ref('all')
const dateRange = ref<string[]>([])

// 搜索表单
const searchForm = reactive({
  keyword: '',
  level: '',
  category: '',
  isRead: null as boolean | null
})

// 分页
const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0
})

// 统计数据
const statistics = reactive({
  total: 0,
  unread: 0,
  today: 0,
  alert: 0,
  cert: 0
})

// 计算未读数量
const unreadCount = computed(() => statistics.unread)

// 获取通知列表
const fetchNotificationList = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取通知列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 模拟数据
    notificationList.value = [
      {
        id: 1,
        title: 'SSL证书即将过期',
        content: '域名 example.com 的SSL证书将在7天后过期，请及时续期',
        level: 'warning',
        category: 'cert',
        isRead: false,
        createdAt: new Date().toISOString()
      },
      {
        id: 2,
        title: 'HTTPS监控异常',
        content: 'https://api.example.com 监控检测到连接超时',
        level: 'error',
        category: 'monitor',
        isRead: true,
        createdAt: new Date(Date.now() - 3600000).toISOString()
      },
      {
        id: 3,
        title: 'DNS解析更新成功',
        content: '域名 test.com 的DNS记录更新已完成',
        level: 'success',
        category: 'dns',
        isRead: false,
        createdAt: new Date(Date.now() - 7200000).toISOString()
      }
    ]
    
    pagination.total = notificationList.value.length
    updateStatistics()
  } catch (error) {
    ElMessage.error('获取通知列表失败')
  } finally {
    loading.value = false
  }
}

// 更新统计数据
const updateStatistics = () => {
  statistics.total = notificationList.value.length
  statistics.unread = notificationList.value.filter(item => !item.isRead).length
  statistics.today = notificationList.value.filter(item => {
    const today = new Date().toDateString()
    return new Date(item.createdAt).toDateString() === today
  }).length
  statistics.alert = notificationList.value.filter(item => 
    item.level === 'error' || item.level === 'warning'
  ).length
  statistics.cert = notificationList.value.filter(item => 
    item.category === 'cert'
  ).length
}

// 处理选择变化
const handleSelectionChange = (selection: any[]) => {
  selectedRowKeys.value = selection.map(item => item.id)
}

// 设置筛选
const setFilter = (filter: string) => {
  activeFilter.value = filter
  
  // 重置搜索条件
  Object.assign(searchForm, {
    keyword: '',
    level: '',
    category: '',
    isRead: null
  })

  // 根据筛选条件设置
  switch (filter) {
    case 'unread':
      searchForm.isRead = false
      break
    case 'alert':
      searchForm.level = 'error'
      break
    case 'cert':
      searchForm.category = 'cert'
      break
  }
  
  fetchNotificationList()
}

// 处理查看
const handleView = (record: any) => {
  currentNotification.value = record
  detailModalVisible.value = true
  
  // 如果是未读通知，标记为已读
  if (!record.isRead) {
    handleMarkRead(record, false)
  }
}

// 处理标记已读
const handleMarkRead = async (record: any, showMessage = true) => {
  try {
    // TODO: 调用API标记已读
    record.isRead = true
    updateStatistics()
    if (showMessage) {
      ElMessage.success('标记为已读')
    }
  } catch (error) {
    ElMessage.error('标记失败')
  }
}

// 处理全部已读
const handleMarkAllRead = async () => {
  try {
    // TODO: 调用API全部标记已读
    notificationList.value.forEach(item => {
      item.isRead = true
    })
    updateStatistics()
    ElMessage.success('全部标记为已读')
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 处理删除
const handleDelete = async (record: any) => {
  try {
    await ElMessageBox.confirm('确定要删除此通知吗？', '确认删除', {
      type: 'warning'
    })
    // TODO: 调用API删除通知
    ElMessage.success('删除成功')
    fetchNotificationList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 处理批量删除
const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(`确定要删除选中的${selectedRowKeys.value.length}条通知吗？`, '确认删除', {
      type: 'warning'
    })
    // TODO: 调用API批量删除
    ElMessage.success('批量删除成功')
    selectedRowKeys.value = []
    fetchNotificationList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 处理刷新
const handleRefresh = () => {
  fetchNotificationList()
}

// 处理添加
const handleAdd = () => {
  notificationModalVisible.value = true
}

// 处理分页
const handlePageChange = (page: number) => {
  pagination.current = page
  fetchNotificationList()
}

const handlePageSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.current = 1
  fetchNotificationList()
}

// 处理通知成功
const handleNotificationSuccess = () => {
  fetchNotificationList()
}

// 工具函数
const getLevelType = (level: string) => {
  const levelMap: Record<string, any> = {
    'info': '',
    'success': 'success',
    'warning': 'warning',
    'error': 'danger'
  }
  return levelMap[level] || ''
}

const getLevelText = (level: string) => {
  const levelMap: Record<string, string> = {
    'info': '信息',
    'success': '成功',
    'warning': '警告',
    'error': '错误'
  }
  return levelMap[level] || level
}

const getCategoryText = (category: string) => {
  const categoryMap: Record<string, string> = {
    'cert': '证书',
    'monitor': '监控',
    'dns': 'DNS',
    'system': '系统',
    'security': '安全'
  }
  return categoryMap[category] || category
}

const formatDateTime = (dateTime: string) => {
  return new Date(dateTime).toLocaleString()
}

onMounted(() => {
  fetchNotificationList()
})
</script>

<style scoped>
.notification-container {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.page-title h1 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
}

.page-title p {
  margin: 0;
  color: var(--el-text-color-regular);
}

.page-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.stats-container {
  margin-bottom: 20px;
}

.filter-card {
  margin-bottom: 20px;
}

.search-input {
  text-align: right;
}

.notification-list-card {
  margin-bottom: 20px;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.list-actions {
  display: flex;
  gap: 8px;
}

.notification-title {
  cursor: pointer;
}

.category-tag {
  margin-left: 8px;
}

.notification-content {
  color: var(--el-text-color-regular);
}

.notification-time {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.notification-content-detail {
  white-space: pre-wrap;
  line-height: 1.6;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
