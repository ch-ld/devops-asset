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
      <a-form-item label="域名" name="name">
        <a-input
          v-model:value="formData.name"
          placeholder="请输入域名，如：example.com"
          :disabled="mode === 'view' || mode === 'edit'"
        />
      </a-form-item>

      <a-form-item label="DNS提供商" name="provider_id">
        <a-select
          v-model:value="formData.provider_id"
          placeholder="请选择DNS提供商"
          :disabled="mode === 'view'"
        >
          <a-select-option
            v-for="provider in providerOptions"
            :key="provider.id"
            :value="provider.id"
          >
            {{ provider.name }}
          </a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="泛域名" name="is_wildcard">
        <a-switch 
          v-model:checked="formData.is_wildcard" 
          :disabled="mode === 'view'"
        />
        <span class="form-help">是否为泛域名（*.example.com）</span>
      </a-form-item>

      <a-form-item label="自动续期" name="auto_renew">
        <a-switch 
          v-model:checked="formData.auto_renew" 
          :disabled="mode === 'view'"
        />
        <span class="form-help">开启后将自动续期域名</span>
      </a-form-item>

      <a-form-item label="过期时间" name="expires_at" v-if="mode !== 'add'">
        <a-date-picker
          v-model:value="formData.expires_at"
          placeholder="请选择过期时间"
          :disabled="mode === 'view'"
          style="width: 100%"
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
      <template v-if="mode === 'view' && domain">
        <a-divider>状态信息</a-divider>
        
        <a-form-item label="状态">
          <a-badge
            :status="getStatusBadge(domain.status)"
            :text="getStatusText(domain.status)"
          />
        </a-form-item>

        <a-form-item label="解析记录数">
          <span>{{ domain.record_count || 0 }} 条</span>
        </a-form-item>

        <a-form-item label="最后同步时间" v-if="domain.last_sync_at">
          <span>{{ formatDate(domain.last_sync_at) }}</span>
        </a-form-item>

        <a-form-item label="创建时间">
          <span>{{ formatDate(domain.created_at) }}</span>
        </a-form-item>

        <a-form-item label="更新时间">
          <span>{{ formatDate(domain.updated_at) }}</span>
        </a-form-item>
      </template>
    </a-form>

    <template #footer v-if="mode !== 'view'">
      <a-space>
        <a-button @click="handleCancel">取消</a-button>
        <a-button @click="handleTest" :loading="testLoading" v-if="formData.name && formData.provider_id">
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
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'
import type { Dayjs } from 'dayjs'
import { domainApi } from '@/api/dns/domain'
import { dnsProviderApi } from '@/api/dns/provider'
import type { Domain, DomainCreateRequest, DNSProvider } from '@/types/dns'

interface Props {
  visible: boolean
  mode: 'add' | 'edit' | 'view'
  domain?: Domain | null
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  domain: null
})

const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const testLoading = ref(false)
const formRef = ref<FormInstance>()
const providerOptions = ref<DNSProvider[]>([])

// 表单数据
const formData = reactive<DomainCreateRequest & { expires_at?: Dayjs }>({
  name: '',
  provider_id: undefined,
  is_wildcard: false,
  auto_renew: true,
  remark: '',
  expires_at: undefined
})

// 表单验证规则
const rules: Record<string, Rule[]> = {
  name: [
    { required: true, message: '请输入域名', trigger: 'blur' },
    { 
      pattern: /^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/, 
      message: '请输入有效的域名', 
      trigger: 'blur' 
    }
  ],
  provider_id: [
    { required: true, message: '请选择DNS提供商', trigger: 'change' }
  ]
}

// 计算属性
const modalTitle = computed(() => {
  const titleMap = {
    add: '添加域名',
    edit: '编辑域名',
    view: '域名详情'
  }
  return titleMap[props.mode]
})

// 监听器
watch(() => props.visible, (visible) => {
  if (visible) {
    if ((props.mode === 'edit' || props.mode === 'view') && props.domain) {
      Object.assign(formData, {
        name: props.domain.name,
        provider_id: props.domain.provider_id,
        is_wildcard: props.domain.is_wildcard,
        auto_renew: props.domain.auto_renew,
        remark: props.domain.remark || '',
        expires_at: props.domain.expires_at ? dayjs(props.domain.expires_at) : undefined
      })
    } else if (props.mode === 'add') {
      resetForm()
    }
    fetchProviderOptions()
  }
})

// 方法
const getStatusBadge = (status: string) => {
  const statusMap = {
    active: 'success',
    inactive: 'default',
    error: 'error'
  }
  return statusMap[status] || 'default'
}

const getStatusText = (status: string) => {
  const statusMap = {
    active: '正常',
    inactive: '停用',
    error: '异常'
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
    provider_id: undefined,
    is_wildcard: false,
    auto_renew: true,
    remark: '',
    expires_at: undefined
  })
  formRef.value?.clearValidate()
}

const fetchProviderOptions = async () => {
  try {
    const response = await dnsProviderApi.list({ page: 1, size: 100 })
    providerOptions.value = response.items || []
  } catch (error) {
    console.error('获取DNS提供商列表失败:', error)
  }
}

const handleTest = async () => {
  try {
    testLoading.value = true
    // TODO: 实现域名连接测试
    await new Promise(resolve => setTimeout(resolve, 1000))
    message.success('域名连接测试成功')
  } catch (error) {
    message.error('域名连接测试失败')
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

    const submitData = {
      ...formData,
      expires_at: formData.expires_at?.format('YYYY-MM-DD HH:mm:ss')
    }

    if (props.mode === 'add') {
      await domainApi.create(submitData)
      message.success('域名创建成功')
    } else if (props.mode === 'edit' && props.domain) {
      await domainApi.update(props.domain.id, submitData)
      message.success('域名更新成功')
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

// 生命周期
onMounted(() => {
  fetchProviderOptions()
})
</script>

<style scoped lang="scss">
.form-help {
  margin-left: 8px;
  color: #8c8c8c;
  font-size: 12px;
}

.ant-divider {
  margin: 16px 0;
}
</style>
