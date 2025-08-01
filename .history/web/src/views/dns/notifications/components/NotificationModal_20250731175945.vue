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

// 表单验证规则
const rules: Record<string, Rule[]> = {
  type: [
    { required: true, message: '请选择通知类型', trigger: 'change' }
  ],
  severity: [
    { required: true, message: '请选择严重程度', trigger: 'change' }
  ],
  title: [
    { required: true, message: '请输入通知标题', trigger: 'blur' },
    { min: 2, max: 100, message: '标题长度在2-100个字符', trigger: 'blur' }
  ],
  message: [
    { required: true, message: '请输入通知内容', trigger: 'blur' },
    { min: 5, max: 1000, message: '内容长度在5-1000个字符', trigger: 'blur' }
  ],
  channels: [
    {
      type: 'array',
      min: 1,
      message: '请至少选择一个通知渠道',
      trigger: 'change'
    }
  ]
}

// 计算属性
const modalTitle = computed(() => {
  const titleMap = {
    add: '创建通知',
    edit: '编辑通知',
    view: '通知详情'
  }
  return titleMap[props.mode]
})

// 监听器
watch(() => props.visible, (visible) => {
  if (visible) {
    if ((props.mode === 'edit' || props.mode === 'view') && props.notification) {
      const notification = props.notification
      Object.assign(formData, {
        type: notification.type,
        title: notification.title,
        message: notification.message,
        severity: notification.severity,
        resource_type: notification.resource_type,
        resource_id: notification.resource_id,
        channels: notification.channels || ['email'],
        metadata: notification.metadata || {},
        metadata_list: Object.entries(notification.metadata || {}).map(([key, value]) => ({
          key,
          value: String(value)
        }))
      })

      if (notification.resource_type) {
        fetchResourceOptions(notification.resource_type)
      }
    } else if (props.mode === 'add') {
      resetForm()
    }
  }
})

// 方法
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

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const filterOption = (input: string, option: any) => {
  return option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
}

const addMetadata = () => {
  formData.metadata_list.push({ key: '', value: '' })
}

const removeMetadata = (index: number) => {
  formData.metadata_list.splice(index, 1)
}

const handleResourceTypeChange = (resourceType: string) => {
  formData.resource_id = undefined
  resourceOptions.value = []
  if (resourceType) {
    fetchResourceOptions(resourceType)
  }
}

const fetchResourceOptions = async (resourceType: string) => {
  try {
    resourceLoading.value = true

    // 根据资源类型获取对应的资源列表
    // 这里应该调用相应的API
    // 暂时使用模拟数据
    const mockOptions: ResourceOption[] = []

    switch (resourceType) {
      case 'certificate':
        // 获取证书列表
        for (let i = 1; i <= 10; i++) {
          mockOptions.push({
            id: i,
            name: `example${i}.com`
          })
        }
        break
      case 'monitor':
        // 获取监控列表
        for (let i = 1; i <= 10; i++) {
          mockOptions.push({
            id: i,
            name: `监控${i} - https://example${i}.com`
          })
        }
        break
      case 'domain':
        // 获取域名列表
        for (let i = 1; i <= 10; i++) {
          mockOptions.push({
            id: i,
            name: `example${i}.com`
          })
        }
        break
      case 'provider':
        // 获取DNS提供商列表
        mockOptions.push(
          { id: 1, name: '阿里云DNS' },
          { id: 2, name: '腾讯云DNS' },
          { id: 3, name: 'Cloudflare' }
        )
        break
    }

    resourceOptions.value = mockOptions
  } catch (error) {
    message.error('获取资源列表失败')
  } finally {
    resourceLoading.value = false
  }
}

const resetForm = () => {
  Object.assign(formData, {
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
  resourceOptions.value = []
  formRef.value?.clearValidate()
}

const handlePreview = () => {
  previewVisible.value = true
}

const handleOk = async () => {
  if (props.mode === 'view') {
    handleCancel()
    return
  }

  try {
    await formRef.value?.validate()
    loading.value = true

    // 处理元数据
    const metadata = {}
    formData.metadata_list.forEach(item => {
      if (item.key && item.value) {
        metadata[item.key] = item.value
      }
    })

    const submitData = {
      ...formData,
      metadata
    }
    delete submitData.metadata_list

    if (props.mode === 'add') {
      await notificationApi.create(submitData)
      message.success('通知创建成功')
    } else if (props.mode === 'edit' && props.notification) {
      await notificationApi.update(props.notification.id, submitData)
      message.success('通知更新成功')
    }

    emit('success')
  } catch (error) {
    console.error('操作失败:', error)
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  emit('update:visible', false)
  resetForm()
  previewVisible.value = false
}
</script>

<style scoped lang="scss">
.metadata-config {
  .metadata-item {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
  }
}

.notification-preview {
  .preview-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 16px;

    .preview-title {
      display: flex;
      align-items: center;
      gap: 8px;

      .title-text {
        font-size: 16px;
        font-weight: 500;
        color: #262626;
      }
    }
  }

  .preview-content {
    background: #fafafa;
    padding: 16px;
    border-radius: 6px;
    margin-bottom: 16px;
    line-height: 1.6;
    color: #595959;
  }

  .preview-channels {
    margin-bottom: 16px;

    strong {
      color: #262626;
    }
  }

  .preview-time {
    display: flex;
    align-items: center;
    font-size: 12px;
    color: #8c8c8c;

    .anticon {
      margin-right: 4px;
    }
  }
}

.ant-divider {
  margin: 16px 0;
}
</style>
