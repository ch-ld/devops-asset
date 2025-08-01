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
