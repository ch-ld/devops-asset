<template>
  <a-modal
    :visible="visible"
    :title="`通知详情 - ${notification?.title}`"
    :width="900"
    :footer="null"
    @cancel="handleCancel"
  >
    <div class="notification-detail" v-if="notification">
      <!-- 通知头部 -->
      <div class="detail-header">
        <div class="header-left">
          <div class="notification-title">
            <h3>{{ notification.title }}</h3>
            <div class="notification-tags">
              <a-tag :color="getTypeColor(notification.type)">
                {{ getTypeText(notification.type) }}
              </a-tag>
              <a-tag :color="getSeverityColor(notification.severity)">
                {{ getSeverityText(notification.severity) }}
              </a-tag>
              <a-badge
                :status="notification.status === 'read' ? 'success' : 'processing'"
                :text="notification.status === 'read' ? '已读' : '未读'"
              />
            </div>
          </div>
        </div>
        <div class="header-right">
          <a-space>
            <a-button @click="handleMarkRead" v-if="notification.status === 'unread'">
              <CheckOutlined />
              标记已读
            </a-button>
            <a-button @click="handleResend" v-if="notification.sent_at">
              <SendOutlined />
              重新发送
            </a-button>
            <a-button @click="handleEdit">
              <EditOutlined />
              编辑
            </a-button>
          </a-space>
        </div>
      </div>

      <!-- 通知内容 -->
      <div class="detail-content">
        <a-card title="通知内容" size="small">
          <div class="message-content">
            {{ notification.message }}
          </div>
        </a-card>
      </div>

      <!-- 基本信息 -->
      <div class="detail-info">
        <a-row :gutter="24">
          <a-col :span="12">
            <a-card title="基本信息" size="small">
              <a-descriptions :column="1" size="small">
                <a-descriptions-item label="通知ID">
                  {{ notification.id }}
                </a-descriptions-item>
                <a-descriptions-item label="通知类型">
                  <a-tag :color="getTypeColor(notification.type)">
                    {{ getTypeText(notification.type) }}
                  </a-tag>
                </a-descriptions-item>
                <a-descriptions-item label="严重程度">
                  <a-tag :color="getSeverityColor(notification.severity)">
                    {{ getSeverityText(notification.severity) }}
                  </a-tag>
                </a-descriptions-item>
                <a-descriptions-item label="通知状态">
                  <a-badge
                    :status="notification.status === 'read' ? 'success' : 'processing'"
                    :text="notification.status === 'read' ? '已读' : '未读'"
                  />
                </a-descriptions-item>
                <a-descriptions-item label="创建时间">
                  {{ formatDate(notification.created_at) }}
                </a-descriptions-item>
                <a-descriptions-item label="更新时间">
                  {{ formatDate(notification.updated_at) }}
                </a-descriptions-item>
              </a-descriptions>
            </a-card>
          </a-col>
          <a-col :span="12">
            <a-card title="发送信息" size="small">
              <a-descriptions :column="1" size="small">
                <a-descriptions-item label="发送状态">
                  <a-badge
                    v-if="notification.sent_at"
                    status="success"
                    text="已发送"
                  />
                  <a-badge
                    v-else
                    status="processing"
                    text="待发送"
                  />
                </a-descriptions-item>
                <a-descriptions-item label="发送时间" v-if="notification.sent_at">
                  {{ formatDate(notification.sent_at) }}
                </a-descriptions-item>
                <a-descriptions-item label="阅读时间" v-if="notification.read_at">
                  {{ formatDate(notification.read_at) }}
                </a-descriptions-item>
                <a-descriptions-item label="通知渠道">
                  <div class="notification-channels">
                    <a-tag
                      v-for="channel in notification.channels"
                      :key="channel"
                      :color="getChannelColor(channel)"
                      size="small"
                    >
                      {{ getChannelName(channel) }}
                    </a-tag>
                  </div>
                </a-descriptions-item>
              </a-descriptions>
            </a-card>
          </a-col>
        </a-row>
      </div>

      <!-- 关联资源 -->
      <div class="detail-resource" v-if="notification.resource_type && notification.resource_id">
        <a-card title="关联资源" size="small">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item label="资源类型">
              {{ getResourceTypeText(notification.resource_type) }}
            </a-descriptions-item>
            <a-descriptions-item label="资源ID">
              #{{ notification.resource_id }}
            </a-descriptions-item>
          </a-descriptions>
          <div class="resource-actions">
            <a-button @click="handleViewResource" type="link">
              <LinkOutlined />
              查看资源详情
            </a-button>
          </div>
        </a-card>
      </div>

      <!-- 元数据 -->
      <div class="detail-metadata" v-if="notification.metadata && Object.keys(notification.metadata).length">
        <a-card title="扩展信息" size="small">
          <a-descriptions :column="2" size="small">
            <a-descriptions-item
              v-for="[key, value] in Object.entries(notification.metadata)"
              :key="key"
              :label="key"
            >
              {{ value }}
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </div>

      <!-- 发送历史 -->
      <div class="detail-history">
        <a-card title="发送历史" size="small">
          <a-timeline>
            <a-timeline-item color="blue">
              <template #dot>
                <PlusCircleOutlined />
              </template>
              <div class="timeline-content">
                <div class="timeline-title">通知创建</div>
                <div class="timeline-time">{{ formatDate(notification.created_at) }}</div>
              </div>
            </a-timeline-item>
            
            <a-timeline-item 
              v-if="notification.sent_at"
              color="green"
            >
              <template #dot>
                <SendOutlined />
              </template>
              <div class="timeline-content">
                <div class="timeline-title">通知发送</div>
                <div class="timeline-time">{{ formatDate(notification.sent_at) }}</div>
                <div class="timeline-channels">
                  通过
                  <a-tag
                    v-for="channel in notification.channels"
                    :key="channel"
                    :color="getChannelColor(channel)"
                    size="small"
                    style="margin: 0 2px"
                  >
                    {{ getChannelName(channel) }}
                  </a-tag>
                  发送
                </div>
              </div>
            </a-timeline-item>
            
            <a-timeline-item 
              v-if="notification.read_at"
              color="green"
            >
              <template #dot>
                <CheckCircleOutlined />
              </template>
              <div class="timeline-content">
                <div class="timeline-title">通知已读</div>
                <div class="timeline-time">{{ formatDate(notification.read_at) }}</div>
              </div>
            </a-timeline-item>
            
            <a-timeline-item 
              v-if="!notification.sent_at"
              color="gray"
            >
              <template #dot>
                <ClockCircleOutlined />
              </template>
              <div class="timeline-content">
                <div class="timeline-title">等待发送</div>
                <div class="timeline-desc">通知正在队列中等待发送</div>
              </div>
            </a-timeline-item>
          </a-timeline>
        </a-card>
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { message } from 'ant-design-vue'
import {
  CheckOutlined,
  SendOutlined,
  EditOutlined,
  LinkOutlined,
  PlusCircleOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined
} from '@ant-design/icons-vue'
import { notificationApi } from '@/api/dns/notification'
import type { Notification } from '@/types/dns'

interface Props {
  visible: boolean
  notification?: Notification | null
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'edit', notification: Notification): void
  (e: 'refresh'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 方法
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

const getResourceTypeText = (resourceType: string) => {
  const textMap = {
    certificate: '证书',
    monitor: '监控',
    domain: '域名',
    provider: 'DNS提供商'
  }
  return textMap[resourceType] || resourceType
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 事件处理
const handleMarkRead = async () => {
  if (!props.notification) return
  
  try {
    await notificationApi.markAsRead(props.notification.id)
    message.success('已标记为已读')
    emit('refresh')
  } catch (error) {
    message.error('操作失败')
  }
}

const handleResend = async () => {
  if (!props.notification) return
  
  try {
    // TODO: 实现重新发送功能
    message.success('通知重新发送成功')
    emit('refresh')
  } catch (error) {
    message.error('重新发送失败')
  }
}

const handleEdit = () => {
  if (props.notification) {
    emit('edit', props.notification)
    handleCancel()
  }
}

const handleViewResource = () => {
  if (!props.notification?.resource_type || !props.notification?.resource_id) return
  
  // TODO: 根据资源类型跳转到对应的详情页面
  message.info('跳转到资源详情页面')
}

const handleCancel = () => {
  emit('update:visible', false)
}
</script>

<style scoped lang="scss">
.notification-detail {
  .detail-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 24px;
    padding-bottom: 16px;
    border-bottom: 1px solid #f0f0f0;
    
    .header-left {
      flex: 1;
      
      .notification-title {
        h3 {
          margin: 0 0 8px 0;
          font-size: 18px;
          font-weight: 600;
          color: #262626;
        }
        
        .notification-tags {
          display: flex;
          align-items: center;
          gap: 8px;
        }
      }
    }
    
    .header-right {
      flex-shrink: 0;
    }
  }
  
  .detail-content {
    margin-bottom: 24px;
    
    .message-content {
      background: #fafafa;
      padding: 16px;
      border-radius: 6px;
      line-height: 1.6;
      color: #595959;
      white-space: pre-wrap;
    }
  }
  
  .detail-info {
    margin-bottom: 24px;
  }
  
  .detail-resource {
    margin-bottom: 24px;
    
    .resource-actions {
      margin-top: 12px;
    }
  }
  
  .detail-metadata {
    margin-bottom: 24px;
  }
  
  .detail-history {
    .timeline-content {
      .timeline-title {
        font-weight: 500;
        color: #262626;
        margin-bottom: 4px;
      }
      
      .timeline-time {
        font-size: 12px;
        color: #8c8c8c;
        margin-bottom: 4px;
      }
      
      .timeline-desc {
        font-size: 12px;
        color: #8c8c8c;
      }
      
      .timeline-channels {
        font-size: 12px;
        color: #595959;
        margin-top: 4px;
      }
    }
  }
  
  .notification-channels {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
  }
}
</style>
