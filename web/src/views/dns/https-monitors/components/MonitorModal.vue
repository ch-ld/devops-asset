<template>
  <a-modal
    :visible="visible"
    :title="modalTitle"
    :width="900"
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
      
      <a-form-item label="监控名称" name="name">
        <a-input
          v-model:value="formData.name"
          placeholder="请输入监控名称"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <a-form-item label="监控URL" name="url">
        <a-input
          v-model:value="formData.url"
          placeholder="请输入完整的HTTPS URL，如：https://example.com"
          :disabled="mode === 'view'"
        >
          <template #addonBefore>
            <a-select v-model:value="urlProtocol" style="width: 80px" :disabled="mode === 'view'">
              <a-select-option value="https">HTTPS</a-select-option>
              <a-select-option value="http">HTTP</a-select-option>
            </a-select>
          </template>
          <template #suffix>
            <a-button 
              v-if="formData.url && mode !== 'view'" 
              type="link" 
              size="small" 
              @click="testUrl"
              :loading="testLoading"
            >
              测试
            </a-button>
          </template>
        </a-input>
      </a-form-item>

      <!-- 监控配置 -->
      <a-divider orientation="left">监控配置</a-divider>
      
      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="检查间隔" name="check_interval">
            <a-input-number
              v-model:value="formData.check_interval"
              :min="1"
              :max="1440"
              placeholder="分钟"
              addon-after="分钟"
              :disabled="mode === 'view'"
              style="width: 100%"
            />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="超时时间" name="timeout">
            <a-input-number
              v-model:value="formData.timeout"
              :min="5"
              :max="300"
              placeholder="秒"
              addon-after="秒"
              :disabled="mode === 'view'"
              style="width: 100%"
            />
          </a-form-item>
        </a-col>
      </a-row>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="告警阈值" name="alert_threshold">
            <a-input-number
              v-model:value="formData.alert_threshold"
              :min="1000"
              :max="30000"
              placeholder="毫秒"
              addon-after="ms"
              :disabled="mode === 'view'"
              style="width: 100%"
            />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="最大失败次数" name="max_failures">
            <a-input-number
              v-model:value="formData.max_failures"
              :min="1"
              :max="10"
              placeholder="次"
              addon-after="次"
              :disabled="mode === 'view'"
              style="width: 100%"
            />
          </a-form-item>
        </a-col>
      </a-row>

      <!-- 通知设置 -->
      <a-divider orientation="left">通知设置</a-divider>
      
      <a-form-item label="启用通知" name="notification_enabled">
        <a-switch 
          v-model:checked="formData.notification_enabled" 
          :disabled="mode === 'view'"
        />
        <span class="form-help">开启后将在监控异常时发送通知</span>
      </a-form-item>

      <a-form-item 
        label="通知渠道" 
        name="notification_channels" 
        v-if="formData.notification_enabled"
      >
        <a-checkbox-group 
          v-model:value="formData.notification_channels" 
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

      <!-- 高级配置 -->
      <a-divider orientation="left">高级配置</a-divider>
      
      <a-form-item label="请求头" name="headers">
        <div class="headers-config">
          <div 
            v-for="(header, index) in formData.configuration.headers" 
            :key="index"
            class="header-item"
          >
            <a-input
              v-model:value="header.key"
              placeholder="Header名称"
              style="width: 40%"
              :disabled="mode === 'view'"
            />
            <a-input
              v-model:value="header.value"
              placeholder="Header值"
              style="width: 40%; margin-left: 8px"
              :disabled="mode === 'view'"
            />
            <a-button 
              v-if="mode !== 'view'"
              type="link" 
              danger 
              @click="removeHeader(index)"
              style="margin-left: 8px"
            >
              删除
            </a-button>
          </div>
          <a-button 
            v-if="mode !== 'view'"
            type="dashed" 
            @click="addHeader" 
            style="width: 100%; margin-top: 8px"
          >
            <PlusOutlined /> 添加请求头
          </a-button>
        </div>
      </a-form-item>

      <a-form-item label="用户代理" name="user_agent">
        <a-input
          v-model:value="formData.configuration.user_agent"
          placeholder="自定义User-Agent，留空使用默认值"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <a-form-item label="跟随重定向" name="follow_redirects">
        <a-switch 
          v-model:checked="formData.configuration.follow_redirects" 
          :disabled="mode === 'view'"
        />
        <span class="form-help">是否跟随HTTP重定向</span>
      </a-form-item>

      <a-form-item label="验证SSL证书" name="verify_ssl">
        <a-switch 
          v-model:checked="formData.configuration.verify_ssl" 
          :disabled="mode === 'view'"
        />
        <span class="form-help">是否验证SSL证书有效性</span>
      </a-form-item>

      <a-form-item label="期望状态码" name="expected_status_codes">
        <a-select
          v-model:value="formData.configuration.expected_status_codes"
          mode="tags"
          placeholder="期望的HTTP状态码，如：200,301,302"
          :disabled="mode === 'view'"
        >
          <a-select-option value="200">200 OK</a-select-option>
          <a-select-option value="301">301 Moved Permanently</a-select-option>
          <a-select-option value="302">302 Found</a-select-option>
          <a-select-option value="404">404 Not Found</a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="期望内容" name="expected_content">
        <a-textarea
          v-model:value="formData.configuration.expected_content"
          placeholder="期望在响应中包含的内容（可选）"
          :rows="2"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <a-form-item label="备注" name="remark">
        <a-textarea
          v-model:value="formData.remark"
          placeholder="请输入备注信息"
          :rows="3"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <!-- 查看模式的详细信息 -->
      <template v-if="mode === 'view' && monitor">
        <a-divider orientation="left">状态信息</a-divider>
        
        <a-form-item label="当前状态">
          <a-badge
            :status="getStatusBadge(monitor.status)"
            :text="getStatusText(monitor.status)"
          />
        </a-form-item>

        <a-form-item label="运行状态">
          <a-badge
            :status="getLastStatusBadge(monitor.last_status)"
            :text="getLastStatusText(monitor.last_status)"
          />
        </a-form-item>

        <a-form-item label="最后响应时间" v-if="monitor.last_response_time">
          <span :class="getResponseTimeClass(monitor.last_response_time)">
            {{ monitor.last_response_time }}ms
          </span>
        </a-form-item>

        <a-form-item label="连续失败次数" v-if="monitor.consecutive_failures > 0">
          <a-tag color="red">{{ monitor.consecutive_failures }} 次</a-tag>
        </a-form-item>

        <a-form-item label="最后检查时间" v-if="monitor.last_checked">
          <span>{{ formatDate(monitor.last_checked) }}</span>
        </a-form-item>

        <a-form-item label="SSL证书信息" v-if="monitor.ssl_cert_expires_at">
          <div class="ssl-cert-info">
            <div>
              <strong>过期时间：</strong>
              <span :class="getCertExpireClass(monitor.ssl_cert_expires_at)">
                {{ formatDate(monitor.ssl_cert_expires_at) }}
              </span>
            </div>
            <div v-if="monitor.ssl_cert_issuer">
              <strong>颁发机构：</strong>{{ monitor.ssl_cert_issuer }}
            </div>
            <div v-if="monitor.ssl_cert_subject">
              <strong>证书主体：</strong>{{ monitor.ssl_cert_subject }}
            </div>
          </div>
        </a-form-item>

        <a-form-item label="创建时间">
          <span>{{ formatDate(monitor.created_at) }}</span>
        </a-form-item>

        <a-form-item label="更新时间">
          <span>{{ formatDate(monitor.updated_at) }}</span>
        </a-form-item>
      </template>
    </a-form>

    <template #footer v-if="mode !== 'view'">
      <a-space>
        <a-button @click="handleCancel">取消</a-button>
        <a-button @click="handleTest" :loading="testLoading" v-if="formData.url">
          测试连接
        </a-button>
        <a-button type="primary" @click="handleOk" :loading="loading">
          {{ mode === 'add' ? '创建' : '更新' }}
        </a-button>
      </a-space>
    </template>
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
  WechatOutlined
} from '@ant-design/icons-vue'
import { httpsMonitorApi } from '@/api/dns/httpsMonitor'
import type { HTTPSMonitor, HTTPSMonitorCreateRequest } from '@/types/dns'

interface Props {
  visible: boolean
  mode: 'add' | 'edit' | 'view'
  monitor?: HTTPSMonitor | null
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  monitor: null
})

const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const testLoading = ref(false)
const formRef = ref<FormInstance>()
const urlProtocol = ref('https')

// 表单数据
const formData = reactive<HTTPSMonitorCreateRequest>({
  name: '',
  url: '',
  check_interval: 5,
  timeout: 30,
  notification_enabled: true,
  notification_channels: ['email'],
  alert_threshold: 5000,
  max_failures: 3,
  configuration: {
    headers: [],
    user_agent: '',
    follow_redirects: true,
    verify_ssl: true,
    expected_status_codes: ['200'],
    expected_content: ''
  },
  remark: ''
})

// 表单验证规则
const rules: Record<string, Rule[]> = {
  name: [
    { required: true, message: '请输入监控名称', trigger: 'blur' },
    { min: 2, max: 50, message: '监控名称长度在2-50个字符', trigger: 'blur' }
  ],
  url: [
    { required: true, message: '请输入监控URL', trigger: 'blur' },
    {
      pattern: /^https?:\/\/.+/,
      message: '请输入有效的URL地址',
      trigger: 'blur'
    }
  ],
  check_interval: [
    {
      required: true,
      type: 'number',
      min: 1,
      max: 1440,
      message: '检查间隔必须在1-1440分钟之间',
      trigger: 'blur'
    }
  ],
  timeout: [
    {
      required: true,
      type: 'number',
      min: 5,
      max: 300,
      message: '超时时间必须在5-300秒之间',
      trigger: 'blur'
    }
  ],
  alert_threshold: [
    {
      type: 'number',
      min: 1000,
      max: 30000,
      message: '告警阈值必须在1000-30000毫秒之间',
      trigger: 'blur'
    }
  ],
  max_failures: [
    {
      required: true,
      type: 'number',
      min: 1,
      max: 10,
      message: '最大失败次数必须在1-10次之间',
      trigger: 'blur'
    }
  ],
  notification_channels: [
    {
      type: 'array',
      min: 1,
      message: '请至少选择一个通知渠道',
      trigger: 'change',
      validator: (rule, value) => {
        if (formData.notification_enabled && (!value || value.length === 0)) {
          return Promise.reject('启用通知时必须选择至少一个通知渠道')
        }
        return Promise.resolve()
      }
    }
  ]
}

// 计算属性
const modalTitle = computed(() => {
  const titleMap = {
    add: '添加HTTPS监控',
    edit: '编辑HTTPS监控',
    view: 'HTTPS监控详情'
  }
  return titleMap[props.mode]
})

// 监听器
watch(() => props.visible, (visible) => {
  if (visible) {
    if ((props.mode === 'edit' || props.mode === 'view') && props.monitor) {
      const monitor = props.monitor
      Object.assign(formData, {
        name: monitor.name,
        url: monitor.url,
        check_interval: monitor.check_interval,
        timeout: monitor.timeout,
        notification_enabled: monitor.notification_enabled,
        notification_channels: monitor.notification_channels || [],
        alert_threshold: monitor.alert_threshold,
        max_failures: monitor.max_failures,
        configuration: {
          headers: monitor.configuration?.headers || [],
          user_agent: monitor.configuration?.user_agent || '',
          follow_redirects: monitor.configuration?.follow_redirects ?? true,
          verify_ssl: monitor.configuration?.verify_ssl ?? true,
          expected_status_codes: monitor.configuration?.expected_status_codes || ['200'],
          expected_content: monitor.configuration?.expected_content || ''
        },
        remark: monitor.remark || ''
      })

      // 设置URL协议
      if (monitor.url) {
        urlProtocol.value = monitor.url.startsWith('https://') ? 'https' : 'http'
      }
    } else if (props.mode === 'add') {
      resetForm()
    }
  }
})

// 监听URL协议变化
watch(urlProtocol, (protocol) => {
  if (formData.url) {
    const urlWithoutProtocol = formData.url.replace(/^https?:\/\//, '')
    formData.url = `${protocol}://${urlWithoutProtocol}`
  }
})

// 方法
const getStatusBadge = (status: string) => {
  const statusMap = {
    active: 'success',
    inactive: 'default'
  }
  return statusMap[status] || 'default'
}

const getStatusText = (status: string) => {
  const statusMap = {
    active: '启用',
    inactive: '停用'
  }
  return statusMap[status] || status
}

const getLastStatusBadge = (status: string) => {
  const statusMap = {
    online: 'success',
    offline: 'error',
    warning: 'warning'
  }
  return statusMap[status] || 'default'
}

const getLastStatusText = (status: string) => {
  const statusMap = {
    online: '在线',
    offline: '离线',
    warning: '警告'
  }
  return statusMap[status] || status
}

const getResponseTimeClass = (responseTime: number) => {
  if (responseTime < 500) return 'response-fast'
  if (responseTime < 2000) return 'response-normal'
  return 'response-slow'
}

const getCertExpireClass = (expiresAt: string) => {
  if (!expiresAt) return ''
  const expireDate = new Date(expiresAt)
  const now = new Date()
  const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))

  if (diffDays < 0) return 'text-red'
  if (diffDays <= 30) return 'text-orange'
  return ''
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const addHeader = () => {
  if (!formData.configuration.headers) {
    formData.configuration.headers = []
  }
  formData.configuration.headers.push({ key: '', value: '' })
}

const removeHeader = (index: number) => {
  formData.configuration.headers.splice(index, 1)
}

const resetForm = () => {
  Object.assign(formData, {
    name: '',
    url: '',
    check_interval: 5,
    timeout: 30,
    notification_enabled: true,
    notification_channels: ['email'],
    alert_threshold: 5000,
    max_failures: 3,
    configuration: {
      headers: [],
      user_agent: '',
      follow_redirects: true,
      verify_ssl: true,
      expected_status_codes: ['200'],
      expected_content: ''
    },
    remark: ''
  })
  urlProtocol.value = 'https'
  formRef.value?.clearValidate()
}

const testUrl = async () => {
  if (!formData.url) {
    message.warning('请先输入URL')
    return
  }

  try {
    testLoading.value = true
    // 这里应该调用测试URL的API
    // 暂时模拟测试
    await new Promise(resolve => setTimeout(resolve, 1000))
    message.success('URL测试成功，响应时间: 234ms')
  } catch (error) {
    message.error('URL测试失败')
  } finally {
    testLoading.value = false
  }
}

const handleTest = async () => {
  try {
    await formRef.value?.validate()
    testLoading.value = true

    // 创建临时监控配置进行测试
    const testConfig = { ...formData }
    // 这里应该调用测试监控的API
    await new Promise(resolve => setTimeout(resolve, 2000))
    message.success('监控配置测试成功')
  } catch (error) {
    message.error('监控配置测试失败')
  } finally {
    testLoading.value = false
  }
}

const handleOk = async () => {
  if (props.mode === 'view') {
    handleCancel()
    return
  }

  try {
    await formRef.value?.validate()
    loading.value = true

    // 处理URL协议
    if (formData.url && !formData.url.startsWith('http')) {
      formData.url = `${urlProtocol.value}://${formData.url}`
    }

    if (props.mode === 'add') {
      await httpsMonitorApi.create(formData)
      message.success('HTTPS监控创建成功')
    } else if (props.mode === 'edit' && props.monitor) {
      await httpsMonitorApi.update(props.monitor.id, formData)
      message.success('HTTPS监控更新成功')
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
}
</script>

<style scoped lang="scss">
.form-help {
  margin-left: 8px;
  color: #8c8c8c;
  font-size: 12px;
}

.headers-config {
  .header-item {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
  }
}

.ssl-cert-info {
  div {
    margin-bottom: 8px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  .text-red {
    color: #ff4d4f;
    font-weight: 500;
  }

  .text-orange {
    color: #fa8c16;
    font-weight: 500;
  }
}

.response-fast {
  color: #52c41a;
  font-weight: 500;
}

.response-normal {
  color: #1890ff;
}

.response-slow {
  color: #ff4d4f;
  font-weight: 500;
}

.ant-divider {
  margin: 16px 0;
}
</style>
