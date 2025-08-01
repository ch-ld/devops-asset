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
      <a-form-item label="域名" name="domain_id" v-if="mode === 'add'">
        <a-select
          v-model:value="formData.domain_id"
          placeholder="请选择域名"
          show-search
          :filter-option="filterOption"
        >
          <a-select-option
            v-for="domain in domainOptions"
            :key="domain.id"
            :value="domain.id"
          >
            {{ domain.name }}
          </a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="主域名" name="common_name" v-if="mode === 'add'">
        <a-input
          v-model:value="formData.common_name"
          placeholder="请输入主域名，如：example.com"
        />
      </a-form-item>

      <a-form-item label="备用域名" name="subject_alt_names" v-if="mode === 'add'">
        <a-select
          v-model:value="formData.subject_alt_names"
          mode="tags"
          placeholder="请输入备用域名，如：www.example.com"
          :token-separators="[',', ' ']"
        />
      </a-form-item>

      <a-form-item label="CA类型" name="ca_type" v-if="mode === 'add'">
        <a-select v-model:value="formData.ca_type" placeholder="请选择CA类型">
          <a-select-option value="letsencrypt">Let's Encrypt</a-select-option>
          <a-select-option value="zerossl">ZeroSSL</a-select-option>
          <a-select-option value="buypass">Buypass</a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="自动续期" name="auto_renew" v-if="mode !== 'view'">
        <a-switch v-model:checked="formData.auto_renew" />
        <span class="form-help">开启后将在证书过期前自动续期</span>
      </a-form-item>

      <a-form-item label="续期天数" name="renew_days" v-if="formData.auto_renew && mode !== 'view'">
        <a-input-number
          v-model:value="formData.renew_days"
          :min="1"
          :max="90"
          placeholder="30"
          addon-after="天"
        />
        <span class="form-help">在证书过期前多少天开始续期</span>
      </a-form-item>

      <a-form-item label="备注" name="remark" v-if="mode !== 'view'">
        <a-textarea
          v-model:value="formData.remark"
          placeholder="请输入备注信息"
          :rows="3"
        />
      </a-form-item>

      <!-- 查看模式的详细信息 -->
      <template v-if="mode === 'view' && certificate">
        <a-form-item label="状态">
          <a-badge
            :status="getStatusBadge(certificate.status)"
            :text="getStatusText(certificate.status)"
          />
        </a-form-item>

        <a-form-item label="序列号" v-if="certificate.serial_number">
          <a-typography-text copyable>{{ certificate.serial_number }}</a-typography-text>
        </a-form-item>

        <a-form-item label="指纹" v-if="certificate.fingerprint">
          <a-typography-text copyable>{{ certificate.fingerprint }}</a-typography-text>
        </a-form-item>

        <a-form-item label="签发时间" v-if="certificate.issued_at">
          <span>{{ formatDate(certificate.issued_at) }}</span>
        </a-form-item>

        <a-form-item label="过期时间" v-if="certificate.expires_at">
          <span :class="getExpireClass(certificate.expires_at)">
            {{ formatDate(certificate.expires_at) }}
          </span>
        </a-form-item>

        <a-form-item label="最后续期" v-if="certificate.last_renew_at">
          <span>{{ formatDate(certificate.last_renew_at) }}</span>
        </a-form-item>

        <a-form-item label="自动续期">
          <a-switch :checked="certificate.auto_renew" disabled />
        </a-form-item>

        <a-form-item label="续期天数" v-if="certificate.auto_renew">
          <span>{{ certificate.renew_days }}天</span>
        </a-form-item>

        <a-form-item label="备注" v-if="certificate.remark">
          <span>{{ certificate.remark }}</span>
        </a-form-item>
      </template>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'
import { certificateApi } from '@/api/dns/certificate'
import { domainApi } from '@/api/dns/domain'
import type { Certificate, CertificateCreateRequest, Domain } from '@/types/dns'

interface Props {
  visible: boolean
  mode: 'add' | 'edit' | 'view'
  certificate?: Certificate | null
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  certificate: null
})

const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const formRef = ref<FormInstance>()
const domainOptions = ref<Domain[]>([])

// 表单数据
const formData = reactive<CertificateCreateRequest>({
  domain_id: undefined,
  common_name: '',
  subject_alt_names: [],
  ca_type: 'letsencrypt',
  auto_renew: true,
  renew_days: 30,
  remark: ''
})

// 表单验证规则
const rules: Record<string, Rule[]> = {
  domain_id: [
    { required: true, message: '请选择域名', trigger: 'change' }
  ],
  common_name: [
    { required: true, message: '请输入主域名', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/, message: '请输入有效的域名', trigger: 'blur' }
  ],
  ca_type: [
    { required: true, message: '请选择CA类型', trigger: 'change' }
  ],
  renew_days: [
    { type: 'number', min: 1, max: 90, message: '续期天数必须在1-90之间', trigger: 'blur' }
  ]
}

// 计算属性
const modalTitle = computed(() => {
  const titleMap = {
    add: '申请证书',
    edit: '编辑证书',
    view: '证书详情'
  }
  return titleMap[props.mode]
})

// 监听器
watch(() => props.visible, (visible) => {
  if (visible) {
    if (props.mode === 'edit' && props.certificate) {
      Object.assign(formData, {
        domain_id: props.certificate.domain_id,
        common_name: props.certificate.common_name,
        subject_alt_names: props.certificate.subject_alt_names || [],
        ca_type: props.certificate.ca_type,
        auto_renew: props.certificate.auto_renew,
        renew_days: props.certificate.renew_days,
        remark: props.certificate.remark || ''
      })
    } else if (props.mode === 'add') {
      resetForm()
    }
    fetchDomainOptions()
  }
})

// 方法
const filterOption = (input: string, option: any) => {
  return option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
}

const getStatusBadge = (status: string) => {
  const statusMap = {
    pending: 'processing',
    issued: 'success',
    expired: 'error',
    revoked: 'default'
  }
  return statusMap[status] || 'default'
}

const getStatusText = (status: string) => {
  const statusMap = {
    pending: '申请中',
    issued: '已签发',
    expired: '已过期',
    revoked: '已吊销'
  }
  return statusMap[status] || status
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const getExpireClass = (expiresAt: string) => {
  if (!expiresAt) return ''
  const expireDate = new Date(expiresAt)
  const now = new Date()
  const diffDays = Math.ceil((expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  
  if (diffDays < 0) return 'text-red'
  if (diffDays <= 30) return 'text-orange'
  return ''
}

const resetForm = () => {
  Object.assign(formData, {
    domain_id: undefined,
    common_name: '',
    subject_alt_names: [],
    ca_type: 'letsencrypt',
    auto_renew: true,
    renew_days: 30,
    remark: ''
  })
  formRef.value?.clearValidate()
}

const fetchDomainOptions = async () => {
  try {
    const response = await domainApi.list({ page: 1, size: 1000 })
    domainOptions.value = response.items || []
  } catch (error) {
    console.error('获取域名选项失败:', error)
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
      await certificateApi.issue(formData)
      message.success('证书申请已提交')
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
  fetchDomainOptions()
})
</script>

<style scoped lang="scss">
.form-help {
  margin-left: 8px;
  color: #8c8c8c;
  font-size: 12px;
}

.text-red {
  color: #ff4d4f;
}

.text-orange {
  color: #fa8c16;
}
</style>
