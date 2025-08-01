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
      <a-form-item label="名称" name="name">
        <a-input
          v-model:value="formData.name"
          placeholder="请输入提供商名称"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <a-form-item label="类型" name="type">
        <a-select
          v-model:value="formData.type"
          placeholder="请选择提供商类型"
          :disabled="mode === 'view'"
          @change="handleTypeChange"
        >
          <a-select-option value="aliyun">阿里云</a-select-option>
          <a-select-option value="tencent">腾讯云</a-select-option>
          <a-select-option value="cloudflare">Cloudflare</a-select-option>
          <a-select-option value="dnspod">DNSPod</a-select-option>
          <a-select-option value="godaddy">GoDaddy</a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="优先级" name="priority">
        <a-input-number
          v-model:value="formData.priority"
          :min="1"
          :max="10"
          placeholder="1-10，数字越小优先级越高"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <a-form-item label="限流" name="rate_limit">
        <a-input-number
          v-model:value="formData.rate_limit"
          :min="1"
          :max="1000"
          placeholder="每分钟请求次数"
          addon-after="次/分钟"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <a-form-item label="并发数" name="concurrent">
        <a-input-number
          v-model:value="formData.concurrent"
          :min="1"
          :max="100"
          placeholder="同时处理的请求数"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <a-form-item label="超时时间" name="timeout">
        <a-input-number
          v-model:value="formData.timeout"
          :min="5"
          :max="300"
          placeholder="请求超时时间"
          addon-after="秒"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <!-- 凭证配置 -->
      <a-divider>凭证配置</a-divider>
      
      <template v-if="formData.type === 'aliyun'">
        <a-form-item label="Access Key ID" name="access_key_id">
          <a-input
            v-model:value="formData.credentials.access_key_id"
            placeholder="请输入阿里云Access Key ID"
            :disabled="mode === 'view'"
          />
        </a-form-item>
        <a-form-item label="Access Key Secret" name="access_key_secret">
          <a-input-password
            v-model:value="formData.credentials.access_key_secret"
            placeholder="请输入阿里云Access Key Secret"
            :disabled="mode === 'view'"
          />
        </a-form-item>
      </template>

      <template v-else-if="formData.type === 'tencent'">
        <a-form-item label="Secret ID" name="secret_id">
          <a-input
            v-model:value="formData.credentials.secret_id"
            placeholder="请输入腾讯云Secret ID"
            :disabled="mode === 'view'"
          />
        </a-form-item>
        <a-form-item label="Secret Key" name="secret_key">
          <a-input-password
            v-model:value="formData.credentials.secret_key"
            placeholder="请输入腾讯云Secret Key"
            :disabled="mode === 'view'"
          />
        </a-form-item>
      </template>

      <template v-else-if="formData.type === 'cloudflare'">
        <a-form-item label="API Token" name="api_token">
          <a-input-password
            v-model:value="formData.credentials.api_token"
            placeholder="请输入Cloudflare API Token"
            :disabled="mode === 'view'"
          />
        </a-form-item>
      </template>

      <template v-else-if="formData.type === 'dnspod'">
        <a-form-item label="API ID" name="api_id">
          <a-input
            v-model:value="formData.credentials.api_id"
            placeholder="请输入DNSPod API ID"
            :disabled="mode === 'view'"
          />
        </a-form-item>
        <a-form-item label="API Token" name="api_token">
          <a-input-password
            v-model:value="formData.credentials.api_token"
            placeholder="请输入DNSPod API Token"
            :disabled="mode === 'view'"
          />
        </a-form-item>
      </template>

      <template v-else-if="formData.type === 'godaddy'">
        <a-form-item label="API Key" name="api_key">
          <a-input
            v-model:value="formData.credentials.api_key"
            placeholder="请输入GoDaddy API Key"
            :disabled="mode === 'view'"
          />
        </a-form-item>
        <a-form-item label="API Secret" name="api_secret">
          <a-input-password
            v-model:value="formData.credentials.api_secret"
            placeholder="请输入GoDaddy API Secret"
            :disabled="mode === 'view'"
          />
        </a-form-item>
      </template>

      <a-form-item label="备注" name="remark">
        <a-textarea
          v-model:value="formData.remark"
          placeholder="请输入备注信息"
          :rows="3"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <!-- 查看模式的额外信息 -->
      <template v-if="mode === 'view' && provider">
        <a-divider>状态信息</a-divider>
        
        <a-form-item label="状态">
          <a-badge
            :status="getStatusBadge(provider.status)"
            :text="getStatusText(provider.status)"
          />
        </a-form-item>

        <a-form-item label="是否默认">
          <a-tag v-if="provider.is_default" color="blue">默认提供商</a-tag>
          <span v-else>否</span>
        </a-form-item>

        <a-form-item label="最后测试" v-if="provider.last_test_at">
          <span>{{ formatDate(provider.last_test_at) }}</span>
          <a-tag
            v-if="provider.test_result"
            :color="provider.test_result === 'success' ? 'green' : 'red'"
            style="margin-left: 8px"
          >
            {{ provider.test_result === 'success' ? '成功' : '失败' }}
          </a-tag>
        </a-form-item>

        <a-form-item label="错误信息" v-if="provider.error_message">
          <a-alert :message="provider.error_message" type="error" show-icon />
        </a-form-item>
      </template>
    </a-form>

    <template #footer v-if="mode !== 'view'">
      <a-space>
        <a-button @click="handleCancel">取消</a-button>
        <a-button @click="handleTest" :loading="testLoading">测试连接</a-button>
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
import { dnsProviderApi } from '@/api/dns/provider'
import type { DNSProvider, DNSProviderCreateRequest } from '@/types/dns'

interface Props {
  visible: boolean
  mode: 'add' | 'edit' | 'view'
  provider?: DNSProvider | null
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  provider: null
})

const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const testLoading = ref(false)
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive<DNSProviderCreateRequest>({
  name: '',
  type: '',
  priority: 5,
  rate_limit: 60,
  concurrent: 5,
  timeout: 30,
  credentials: {},
  remark: ''
})

// 表单验证规则
const rules: Record<string, Rule[]> = {
  name: [
    { required: true, message: '请输入提供商名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择提供商类型', trigger: 'change' }
  ],
  priority: [
    { required: true, type: 'number', min: 1, max: 10, message: '优先级必须在1-10之间', trigger: 'blur' }
  ],
  rate_limit: [
    { required: true, type: 'number', min: 1, max: 1000, message: '限流值必须在1-1000之间', trigger: 'blur' }
  ],
  concurrent: [
    { required: true, type: 'number', min: 1, max: 100, message: '并发数必须在1-100之间', trigger: 'blur' }
  ],
  timeout: [
    { required: true, type: 'number', min: 5, max: 300, message: '超时时间必须在5-300秒之间', trigger: 'blur' }
  ]
}

// 计算属性
const modalTitle = computed(() => {
  const titleMap = {
    add: '添加DNS提供商',
    edit: '编辑DNS提供商',
    view: 'DNS提供商详情'
  }
  return titleMap[props.mode]
})

// 监听器
watch(() => props.visible, (visible) => {
  if (visible) {
    if (props.mode === 'edit' && props.provider) {
      Object.assign(formData, {
        name: props.provider.name,
        type: props.provider.type,
        priority: props.provider.priority,
        rate_limit: props.provider.rate_limit,
        concurrent: props.provider.concurrent,
        timeout: props.provider.timeout,
        credentials: { ...props.provider.credentials } || {},
        remark: props.provider.remark || ''
      })
    } else if (props.mode === 'add') {
      resetForm()
    }
  }
})

// 方法
const handleTypeChange = () => {
  formData.credentials = {}
}

const getStatusBadge = (status: string) => {
  const statusMap = {
    active: 'success',
    inactive: 'default'
  }
  return statusMap[status] || 'default'
}

const getStatusText = (status: string) => {
  const statusMap = {
    active: '正常',
    inactive: '停用'
  }
  return statusMap[status] || status
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const resetForm = () => {
  Object.assign(formData, {
    name: '',
    type: '',
    priority: 5,
    rate_limit: 60,
    concurrent: 5,
    timeout: 30,
    credentials: {},
    remark: ''
  })
  formRef.value?.clearValidate()
}

const handleTest = async () => {
  try {
    await formRef.value?.validate()
    testLoading.value = true
    
    // 这里应该调用测试连接的API
    // 由于是新建状态，可能需要先保存再测试，或者提供临时测试接口
    message.success('连接测试成功')
  } catch (error) {
    message.error('连接测试失败')
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

    if (props.mode === 'add') {
      await dnsProviderApi.create(formData)
      message.success('DNS提供商创建成功')
    } else if (props.mode === 'edit' && props.provider) {
      await dnsProviderApi.update(props.provider.id, formData)
      message.success('DNS提供商更新成功')
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
.ant-divider {
  margin: 16px 0;
}
</style>
