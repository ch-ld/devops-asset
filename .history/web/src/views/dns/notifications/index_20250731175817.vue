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
          <a-badge :count="unreadCount" :offset="[10, 0]">
            <a-button @click="handleMarkAllRead" :disabled="unreadCount === 0">
              <template #icon>
                <CheckOutlined />
              </template>
              全部已读
            </a-button>
          </a-badge>
          <a-button type="primary" @click="handleAdd">
            <template #icon>
              <PlusOutlined />
            </template>
            创建通知
          </a-button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-container">
      <a-row :gutter="16">
        <a-col :span="6">
          <a-card>
            <a-statistic
              title="总通知数"
              :value="statistics.total"
              :value-style="{ color: '#1890ff' }"
            >
              <template #suffix>
                <BellOutlined />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic
              title="未读通知"
              :value="statistics.unread"
              :value-style="{ color: '#ff4d4f' }"
            >
              <template #suffix>
                <ExclamationCircleOutlined />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic
              title="今日通知"
              :value="statistics.today"
              :value-style="{ color: '#52c41a' }"
            >
              <template #suffix>
                <CalendarOutlined />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic
              title="紧急通知"
              :value="statistics.critical"
              :value-style="{ color: '#fa8c16' }"
            >
              <template #suffix>
                <WarningOutlined />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-container">
      <a-card :bordered="false">
        <a-form
          ref="searchFormRef"
          :model="searchForm"
          layout="inline"
          class="search-form"
        >
          <a-form-item label="关键词" name="keyword">
            <a-input
              v-model:value="searchForm.keyword"
              placeholder="请输入标题或内容关键词"
              allow-clear
              style="width: 200px"
            />
          </a-form-item>
          <a-form-item label="类型" name="type">
            <a-select
              v-model:value="searchForm.type"
              placeholder="请选择通知类型"
              allow-clear
              style="width: 150px"
            >
              <a-select-option value="certificate_expiring">证书即将过期</a-select-option>
              <a-select-option value="certificate_expired">证书已过期</a-select-option>
              <a-select-option value="monitor_down">监控离线</a-select-option>
              <a-select-option value="monitor_up">监控恢复</a-select-option>
              <a-select-option value="system">系统通知</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="严重程度" name="severity">
            <a-select
              v-model:value="searchForm.severity"
              placeholder="请选择严重程度"
              allow-clear
              style="width: 120px"
            >
              <a-select-option value="info">信息</a-select-option>
              <a-select-option value="warning">警告</a-select-option>
              <a-select-option value="error">错误</a-select-option>
              <a-select-option value="critical">紧急</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="状态" name="status">
            <a-select
              v-model:value="searchForm.status"
              placeholder="请选择状态"
              allow-clear
              style="width: 120px"
            >
              <a-select-option value="unread">未读</a-select-option>
              <a-select-option value="read">已读</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="时间范围" name="time_range">
            <a-range-picker
              v-model:value="searchForm.time_range"
              style="width: 240px"
            />
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="handleSearch">
              <template #icon>
                <SearchOutlined />
              </template>
              搜索
            </a-button>
            <a-button @click="handleReset" style="margin-left: 8px">
              重置
            </a-button>
          </a-form-item>
        </a-form>
      </a-card>
    </div>

    <!-- 通知列表 -->
    <div class="table-container">
      <a-card :bordered="false">
        <template #title>
          <div class="table-title">
            <span>通知列表</span>
            <div class="table-actions">
              <a-button @click="handleRefresh">
                <template #icon>
                  <ReloadOutlined />
                </template>
                刷新
              </a-button>
              <a-button @click="handleBatchRead" :disabled="!hasSelected">
                <template #icon>
                  <CheckOutlined />
                </template>
                批量已读
              </a-button>
              <a-button @click="handleBatchDelete" :disabled="!hasSelected">
                <template #icon>
                  <DeleteOutlined />
                </template>
                批量删除
              </a-button>
            </div>
          </div>
        </template>

        <a-table
          :columns="columns"
          :data-source="tableData"
          :loading="loading"
          :pagination="pagination"
          :row-selection="rowSelection"
          row-key="id"
          @change="handleTableChange"
          :row-class-name="getRowClassName"
        >
          <!-- 通知内容列 -->
          <template #content="{ record }">
            <div class="notification-content">
              <div class="notification-header">
                <span class="notification-title" :class="{ 'unread': record.status === 'unread' }">
                  {{ record.title }}
                </span>
                <div class="notification-meta">
                  <a-tag :color="getTypeColor(record.type)" size="small">
                    {{ getTypeText(record.type) }}
                  </a-tag>
                  <a-tag :color="getSeverityColor(record.severity)" size="small">
                    {{ getSeverityText(record.severity) }}
                  </a-tag>
                </div>
              </div>
              <div class="notification-message">
                {{ truncateText(record.message, 100) }}
              </div>
              <div class="notification-time">
                <ClockCircleOutlined />
                {{ formatRelativeTime(record.created_at) }}
                <span v-if="record.read_at" class="read-time">
                  · 已读于 {{ formatRelativeTime(record.read_at) }}
                </span>
              </div>
            </div>
          </template>

          <!-- 资源信息列 -->
          <template #resource="{ record }">
            <div class="resource-info" v-if="record.resource_type && record.resource_id">
              <div class="resource-type">{{ getResourceTypeText(record.resource_type) }}</div>
              <div class="resource-id">#{{ record.resource_id }}</div>
            </div>
            <span v-else class="no-data">-</span>
          </template>

          <!-- 通知渠道列 -->
          <template #channels="{ record }">
            <div class="notification-channels">
              <a-tag
                v-for="channel in record.channels.slice(0, 3)"
                :key="channel"
                :color="getChannelColor(channel)"
                size="small"
              >
                {{ getChannelName(channel) }}
              </a-tag>
              <a-tag v-if="record.channels.length > 3" size="small">
                +{{ record.channels.length - 3 }}
              </a-tag>
            </div>
          </template>

          <!-- 发送状态列 -->
          <template #sent_status="{ record }">
            <div class="sent-status">
              <a-badge
                v-if="record.sent_at"
                status="success"
                text="已发送"
              />
              <a-badge
                v-else
                status="processing"
                text="待发送"
              />
              <div v-if="record.sent_at" class="sent-time">
                {{ formatRelativeTime(record.sent_at) }}
              </div>
            </div>
          </template>

          <!-- 操作列 -->
          <template #action="{ record }">
            <a-space>
              <a @click="handleView(record)">
                <EyeOutlined />
                查看
              </a>
              <a 
                v-if="record.status === 'unread'" 
                @click="handleMarkRead(record)"
              >
                <CheckOutlined />
                标记已读
              </a>
              <a-dropdown>
                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="handleResend(record)" v-if="record.sent_at">
                      <template #icon>
                        <SendOutlined />
                      </template>
                      重新发送
                    </a-menu-item>
                    <a-menu-item @click="handleViewDetails(record)">
                      <template #icon>
                        <InfoCircleOutlined />
                      </template>
                      查看详情
                    </a-menu-item>
                    <a-menu-divider />
                    <a-menu-item @click="handleDelete(record)" class="danger">
                      <template #icon>
                        <DeleteOutlined />
                      </template>
                      删除
                    </a-menu-item>
                  </a-menu>
                </template>
                <a>
                  更多
                  <DownOutlined />
                </a>
              </a-dropdown>
            </a-space>
          </template>
        </a-table>
      </a-card>
    </div>

    <!-- 通知表单弹窗 -->
    <NotificationModal
      v-model:visible="modalVisible"
      :mode="modalMode"
      :notification="currentNotification"
      @success="handleModalSuccess"
    />

    <!-- 通知详情弹窗 -->
    <NotificationDetailModal
      v-model:visible="detailModalVisible"
      :notification="currentNotification"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  ReloadOutlined,
  DeleteOutlined,
  CheckOutlined,
  BellOutlined,
  ExclamationCircleOutlined,
  CalendarOutlined,
  WarningOutlined,
  ClockCircleOutlined,
  EyeOutlined,
  SendOutlined,
  InfoCircleOutlined,
  DownOutlined
} from '@ant-design/icons-vue'
import type { TableColumnsType } from 'ant-design-vue'
import type { Dayjs } from 'dayjs'
import { usePagination } from '@/hooks/usePagination'
import { useSelection } from '@/hooks/useSelection'
import NotificationModal from './components/NotificationModal.vue'
import NotificationDetailModal from './components/NotificationDetailModal.vue'
import { notificationApi } from '@/api/dns/notification'
import type { Notification } from '@/types/dns'

// 响应式数据
const loading = ref(false)
const modalVisible = ref(false)
const detailModalVisible = ref(false)
const modalMode = ref<'add' | 'edit' | 'view'>('add')
const currentNotification = ref<Notification | null>(null)
const tableData = ref<Notification[]>([])
const unreadCount = ref(0)
const statistics = ref({
  total: 0,
  unread: 0,
  today: 0,
  critical: 0
})

// 搜索表单
const searchFormRef = ref()
const searchForm = reactive({
  keyword: '',
  type: undefined,
  severity: undefined,
  status: undefined,
  time_range: null as [Dayjs, Dayjs] | null
})

// 分页和选择
const { pagination, handleTableChange } = usePagination()
const { selectedRowKeys, rowSelection, hasSelected } = useSelection()

// 表格列定义
const columns: TableColumnsType = [
  {
    title: '通知内容',
    dataIndex: 'content',
    key: 'content',
    slots: { customRender: 'content' },
    width: 400
  },
  {
    title: '资源信息',
    dataIndex: 'resource',
    key: 'resource',
    slots: { customRender: 'resource' },
    width: 120
  },
  {
    title: '通知渠道',
    dataIndex: 'channels',
    key: 'channels',
    slots: { customRender: 'channels' },
    width: 150
  },
  {
    title: '发送状态',
    dataIndex: 'sent_status',
    key: 'sent_status',
    slots: { customRender: 'sent_status' },
    width: 120
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'action' },
    width: 150,
    fixed: 'right'
  }
]

// 工具方法
const getTypeColor = (type: string) => {
  const colorMap = {
    certificate_expiring: 'orange',
    certificate_expired: 'red',
    monitor_down: 'red',
    monitor_up: 'green',
    system: 'blue'
  }
  return colorMap[type] || 'default'
}

const getTypeText = (type: string) => {
  const textMap = {
    certificate_expiring: '证书即将过期',
    certificate_expired: '证书已过期',
    monitor_down: '监控离线',
    monitor_up: '监控恢复',
    system: '系统通知'
  }
  return textMap[type] || type
}

const getSeverityColor = (severity: string) => {
  const colorMap = {
    info: 'blue',
    warning: 'orange',
    error: 'red',
    critical: 'magenta'
  }
  return colorMap[severity] || 'default'
}

const getSeverityText = (severity: string) => {
  const textMap = {
    info: '信息',
    warning: '警告',
    error: '错误',
    critical: '紧急'
  }
  return textMap[severity] || severity
}

const getResourceTypeText = (resourceType: string) => {
  const textMap = {
    certificate: '证书',
    monitor: '监控',
    domain: '域名',
    provider: '提供商'
  }
  return textMap[resourceType] || resourceType
}

const getChannelColor = (channel: string) => {
  const colorMap = {
    email: 'blue',
    sms: 'green',
    webhook: 'purple',
    dingtalk: 'orange',
    wechat: 'cyan'
  }
  return colorMap[channel] || 'default'
}

const getChannelName = (channel: string) => {
  const nameMap = {
    email: '邮件',
    sms: '短信',
    webhook: 'Webhook',
    dingtalk: '钉钉',
    wechat: '微信'
  }
  return nameMap[channel] || channel
}

const truncateText = (text: string, maxLength: number) => {
  if (!text) return ''
  return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
}

const formatRelativeTime = (date: string) => {
  if (!date) return ''
  const diff = Date.now() - new Date(date).getTime()
  const minutes = Math.floor(diff / (1000 * 60))

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  if (days < 7) return `${days}天前`
  return new Date(date).toLocaleDateString('zh-CN')
}

const getRowClassName = (record: Notification) => {
  return record.status === 'unread' ? 'unread-row' : ''
}

// 事件处理
const handleAdd = () => {
  modalMode.value = 'add'
  currentNotification.value = null
  modalVisible.value = true
}

const handleView = (record: Notification) => {
  modalMode.value = 'view'
  currentNotification.value = record
  modalVisible.value = true

  // 如果是未读通知，标记为已读
  if (record.status === 'unread') {
    handleMarkRead(record, false)
  }
}

const handleViewDetails = (record: Notification) => {
  currentNotification.value = record
  detailModalVisible.value = true

  // 如果是未读通知，标记为已读
  if (record.status === 'unread') {
    handleMarkRead(record, false)
  }
}

const handleMarkRead = async (record: Notification, showMessage = true) => {
  try {
    await notificationApi.markAsRead(record.id)
    if (showMessage) {
      message.success('已标记为已读')
    }
    await fetchData()
    await fetchUnreadCount()
  } catch (error) {
    if (showMessage) {
      message.error('操作失败')
    }
  }
}

const handleMarkAllRead = async () => {
  try {
    const unreadNotifications = tableData.value.filter(item => item.status === 'unread')
    if (unreadNotifications.length === 0) {
      message.info('没有未读通知')
      return
    }

    const unreadIds = unreadNotifications.map(item => item.id)
    await notificationApi.batchMarkAsRead(unreadIds)
    message.success('所有通知已标记为已读')
    await fetchData()
    await fetchUnreadCount()
  } catch (error) {
    message.error('操作失败')
  }
}

const handleBatchRead = async () => {
  if (!hasSelected.value) return

  try {
    await notificationApi.batchMarkAsRead(selectedRowKeys.value)
    message.success('批量标记已读成功')
    selectedRowKeys.value = []
    await fetchData()
    await fetchUnreadCount()
  } catch (error) {
    message.error('批量标记已读失败')
  }
}

const handleResend = async (record: Notification) => {
  try {
    // TODO: 实现重新发送功能
    message.success('通知重新发送成功')
  } catch (error) {
    message.error('重新发送失败')
  }
}

const handleDelete = (record: Notification) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除通知 "${record.title}" 吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await notificationApi.delete(record.id)
        message.success('删除成功')
        await fetchData()
        await fetchStatistics()
      } catch (error) {
        message.error('删除失败')
      }
    }
  })
}

const handleBatchDelete = () => {
  if (!hasSelected.value) return

  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 个通知吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await notificationApi.batchDelete(selectedRowKeys.value)
        message.success('批量删除成功')
        selectedRowKeys.value = []
        await fetchData()
        await fetchStatistics()
      } catch (error) {
        message.error('批量删除失败')
      }
    }
  })
}

const handleSearch = () => {
  pagination.current = 1
  fetchData()
}

const handleReset = () => {
  searchFormRef.value?.resetFields()
  pagination.current = 1
  fetchData()
}

const handleRefresh = () => {
  fetchData()
  fetchStatistics()
  fetchUnreadCount()
}

const handleModalSuccess = () => {
  modalVisible.value = false
  fetchData()
  fetchStatistics()
}

// 数据获取
const fetchData = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.current,
      size: pagination.pageSize,
      ...searchForm,
      start_time: searchForm.time_range?.[0]?.format('YYYY-MM-DD HH:mm:ss'),
      end_time: searchForm.time_range?.[1]?.format('YYYY-MM-DD HH:mm:ss')
    }
    delete params.time_range

    const response = await notificationApi.list(params)
    tableData.value = response.items || []
    pagination.total = response.total || 0
  } catch (error) {
    message.error('获取通知列表失败')
  } finally {
    loading.value = false
  }
}

const fetchStatistics = async () => {
  try {
    // TODO: 实现统计数据获取
    // 暂时使用模拟数据
    const total = tableData.value.length
    const unread = tableData.value.filter(item => item.status === 'unread').length
    const today = tableData.value.filter(item => {
      const today = new Date().toDateString()
      return new Date(item.created_at).toDateString() === today
    }).length
    const critical = tableData.value.filter(item => item.severity === 'critical').length

    statistics.value = { total, unread, today, critical }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const fetchUnreadCount = async () => {
  try {
    const response = await notificationApi.getUnreadCount()
    unreadCount.value = response.count || 0
  } catch (error) {
    console.error('获取未读数量失败:', error)
  }
}

// 生命周期
onMounted(() => {
  fetchData()
  fetchStatistics()
  fetchUnreadCount()
})
</script>

<style scoped lang="scss">
.notification-container {
  padding: 24px;
  background: #f5f5f5;
  min-height: 100vh;
}

.page-header {
  margin-bottom: 24px;

  .page-header-content {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;

    .page-title {
      h1 {
        margin: 0 0 8px 0;
        font-size: 24px;
        font-weight: 600;
        color: #262626;
      }

      p {
        margin: 0;
        color: #8c8c8c;
        font-size: 14px;
      }
    }

    .page-actions {
      display: flex;
      gap: 8px;
    }
  }
}

.stats-container {
  margin-bottom: 24px;
}

.search-container {
  margin-bottom: 24px;

  .search-form {
    .ant-form-item {
      margin-bottom: 16px;
    }
  }
}

.table-container {
  .table-title {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .table-actions {
      .ant-btn {
        margin-left: 8px;
      }
    }
  }
}

.notification-content {
  .notification-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 8px;

    .notification-title {
      font-weight: 500;
      color: #262626;

      &.unread {
        font-weight: 600;
        color: #1890ff;
      }
    }

    .notification-meta {
      display: flex;
      gap: 4px;
    }
  }

  .notification-message {
    color: #595959;
    font-size: 13px;
    line-height: 1.4;
    margin-bottom: 8px;
  }

  .notification-time {
    display: flex;
    align-items: center;
    font-size: 12px;
    color: #8c8c8c;

    .anticon {
      margin-right: 4px;
    }

    .read-time {
      color: #52c41a;
    }
  }
}

.resource-info {
  .resource-type {
    font-weight: 500;
    color: #262626;
  }

  .resource-id {
    font-size: 12px;
    color: #8c8c8c;
    margin-top: 2px;
  }
}

.notification-channels {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.sent-status {
  .sent-time {
    font-size: 12px;
    color: #8c8c8c;
    margin-top: 2px;
  }
}

.no-data {
  color: #d9d9d9;
  font-style: italic;
}

.danger {
  color: #ff4d4f !important;
}

// 未读行样式
:deep(.unread-row) {
  background-color: #f6ffed;

  &:hover {
    background-color: #f6ffed !important;
  }
}
</style>
