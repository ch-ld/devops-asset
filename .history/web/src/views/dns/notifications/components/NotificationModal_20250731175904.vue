<template>
  <a-modal
    :visible="visible"
    :title="modalTitle"
    :width="800"
    :confirm-loading="loading"
    :mask-closable="false"
    @ok="handleOk"
    @cancel="handleCancel"
  >
    <a-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      :label-col="{ span: 6 }"
      :wrapper-col="{ span: 18 }"
    >
      <!-- 基本信息 -->
      <a-divider orientation="left">基本信息</a-divider>
      
      <a-form-item label="通知类型" name="type">
        <a-select
          v-model:value="formData.type"
          placeholder="请选择通知类型"
          :disabled="mode === 'view'"
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
          v-model:value="formData.severity"
          placeholder="请选择严重程度"
          :disabled="mode === 'view'"
        >
          <a-select-option value="info">
            <a-tag color="blue">信息</a-tag>
          </a-select-option>
          <a-select-option value="warning">
            <a-tag color="orange">警告</a-tag>
          </a-select-option>
          <a-select-option value="error">
            <a-tag color="red">错误</a-tag>
          </a-select-option>
          <a-select-option value="critical">
            <a-tag color="magenta">紧急</a-tag>
          </a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="通知标题" name="title">
        <a-input
          v-model:value="formData.title"
          placeholder="请输入通知标题"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <a-form-item label="通知内容" name="message">
        <a-textarea
          v-model:value="formData.message"
          placeholder="请输入通知内容"
          :rows="4"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <!-- 关联资源 -->
      <a-divider orientation="left">关联资源</a-divider>
      
      <a-form-item label="资源类型" name="resource_type">
        <a-select
          v-model:value="formData.resource_type"
          placeholder="请选择资源类型（可选）"
          allow-clear
          :disabled="mode === 'view'"
          @change="handleResourceTypeChange"
        >
          <a-select-option value="certificate">证书</a-select-option>
          <a-select-option value="monitor">监控</a-select-option>
          <a-select-option value="domain">域名</a-select-option>
          <a-select-option value="provider">DNS提供商</a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item 
        label="资源ID" 
        name="resource_id" 
        v-if="formData.resource_type"
      >
        <a-select
          v-model:value="formData.resource_id"
          placeholder="请选择关联的资源"
          show-search
          :filter-option="filterOption"
          :disabled="mode === 'view'"
          :loading="resourceLoading"
        >
          <a-select-option
            v-for="resource in resourceOptions"
            :key="resource.id"
            :value="resource.id"
          >
            {{ resource.name }}
          </a-select-option>
        </a-select>
      </a-form-item>

      <!-- 通知渠道 -->
      <a-divider orientation="left">通知渠道</a-divider>
      
      <a-form-item label="通知渠道" name="channels">
        <a-checkbox-group 
          v-model:value="formData.channels" 
          :disabled="mode === 'view'"
        >
          <a-row>
            <a-col :span="8">
              <a-checkbox value="email">
                <MailOutlined /> 邮件通知
              </a-checkbox>
            </a-col>
            <a-col :span="8">
              <a-checkbox value="sms">
                <MessageOutlined /> 短信通知
              </a-checkbox>
            </a-col>
            <a-col :span="8">
              <a-checkbox value="webhook">
                <ApiOutlined /> Webhook
              </a-checkbox>
            </a-col>
            <a-col :span="8">
              <a-checkbox value="dingtalk">
                <DingtalkOutlined /> 钉钉通知
              </a-checkbox>
            </a-col>
            <a-col :span="8">
              <a-checkbox value="wechat">
                <WechatOutlined /> 微信通知
              </a-checkbox>
            </a-col>
          </a-row>
        </a-checkbox-group>
      </a-form-item>

      <!-- 扩展信息 -->
      <a-divider orientation="left">扩展信息</a-divider>
      
      <a-form-item label="元数据" name="metadata">
        <div class="metadata-config">
          <div 
            v-for="(meta, index) in formData.metadata_list" 
            :key="index"
            class="metadata-item"
          >
            <a-input
              v-model:value="meta.key"
              placeholder="键名"
              style="width: 40%"
              :disabled="mode === 'view'"
            />
            <a-input
              v-model:value="meta.value"
              placeholder="键值"
              style="width: 40%; margin-left: 8px"
              :disabled="mode === 'view'"
            />
            <a-button 
              v-if="mode !== 'view'"
              type="link" 
              danger 
              @click="removeMetadata(index)"
              style="margin-left: 8px"
            >
              删除
            </a-button>
          </div>
          <a-button 
            v-if="mode !== 'view'"
            type="dashed" 
            @click="addMetadata" 
            style="width: 100%; margin-top: 8px"
          >
            <PlusOutlined /> 添加元数据
          </a-button>
        </div>
      </a-form-item>

      <!-- 查看模式的详细信息 -->
      <template v-if="mode === 'view' && notification">
        <a-divider orientation="left">状态信息</a-divider>
        
        <a-form-item label="通知状态">
          <a-badge
            :status="notification.status === 'read' ? 'success' : 'processing'"
            :text="notification.status === 'read' ? '已读' : '未读'"
          />
        </a-form-item>

        <a-form-item label="发送状态">
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
        </a-form-item>

        <a-form-item label="发送时间" v-if="notification.sent_at">
          <span>{{ formatDate(notification.sent_at) }}</span>
        </a-form-item>

        <a-form-item label="阅读时间" v-if="notification.read_at">
          <span>{{ formatDate(notification.read_at) }}</span>
        </a-form-item>

        <a-form-item label="创建时间">
          <span>{{ formatDate(notification.created_at) }}</span>
        </a-form-item>

        <a-form-item label="更新时间">
          <span>{{ formatDate(notification.updated_at) }}</span>
        </a-form-item>
      </template>
    </a-form>

    <template #footer v-if="mode !== 'view'">
      <a-space>
        <a-button @click="handleCancel">取消</a-button>
        <a-button @click="handlePreview" v-if="formData.title && formData.message">
          预览
        </a-button>
        <a-button type="primary" @click="handleOk" :loading="loading">
          {{ mode === 'add' ? '创建' : '更新' }}
        </a-button>
      </a-space>
    </template>
  </a-modal>

  <!-- 预览弹窗 -->
  <a-modal
    v-model:visible="previewVisible"
    title="通知预览"
    :width="600"
    :footer="null"
  >
    <div class="notification-preview">
      <div class="preview-header">
        <div class="preview-title">
          <a-tag :color="getSeverityColor(formData.severity)">
            {{ getSeverityText(formData.severity) }}
          </a-tag>
          <span class="title-text">{{ formData.title }}</span>
        </div>
        <div class="preview-type">
          <a-tag :color="getTypeColor(formData.type)">
            {{ getTypeText(formData.type) }}
          </a-tag>
        </div>
      </div>
      <div class="preview-content">
        {{ formData.message }}
      </div>
      <div class="preview-channels" v-if="formData.channels.length">
        <strong>通知渠道：</strong>
        <a-tag
          v-for="channel in formData.channels"
          :key="channel"
          :color="getChannelColor(channel)"
          style="margin-left: 4px"
        >
          {{ getChannelName(channel) }}
        </a-tag>
      </div>
      <div class="preview-time">
        <ClockCircleOutlined />
        {{ new Date().toLocaleString('zh-CN') }}
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { message } from 'ant-design-vue'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'
import {
  PlusOutlined,
  MailOutlined,
  MessageOutlined,
  ApiOutlined,
  DingtalkOutlined,
  WechatOutlined,
  ClockCircleOutlined
} from '@ant-design/icons-vue'
import { notificationApi } from '@/api/dns/notification'
import type { Notification, NotificationCreateRequest } from '@/types/dns'

interface Props {
  visible: boolean
  mode: 'add' | 'edit' | 'view'
  notification?: Notification | null
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

interface ResourceOption {
  id: number
  name: string
}

const props = withDefaults(defineProps<Props>(), {
  notification: null
})

const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const resourceLoading = ref(false)
const previewVisible = ref(false)
const formRef = ref<FormInstance>()
const resourceOptions = ref<ResourceOption[]>([])

// 表单数据
const formData = reactive<NotificationCreateRequest & { metadata_list: Array<{ key: string; value: string }> }>({
  type: '',
  title: '',
  message: '',
  severity: 'info',
  resource_type: undefined,
  resource_id: undefined,
  channels: ['email'],
  metadata: {},
  metadata_list: []
})
