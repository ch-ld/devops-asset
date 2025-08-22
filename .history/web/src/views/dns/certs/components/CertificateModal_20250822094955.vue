<template>
  <el-dialog
    v-model="dialogVisible"
    :title="modalTitle"
    width="900px"
    :close-on-click-modal="false"
    @close="handleCancel"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="120px"
      :disabled="mode === 'view'"
    >
      <el-form-item label="域名" prop="domain_id" v-if="mode === 'add'">
        <el-select
          v-model="formData.domain_id"
          placeholder="请选择域名"
          filterable
          style="width: 100%"
        >
          <el-option
            v-for="domain in domainOptions"
            :key="domain.id"
            :label="domain.name"
            :value="domain.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="主域名" prop="common_name" v-if="mode === 'add'">
        <el-input
          v-model="formData.common_name"
          placeholder="请输入主域名，如：example.com"
        />
      </el-form-item>

      <el-form-item label="备用域名" prop="subject_alt_names" v-if="mode === 'add'">
        <el-select
          v-model="formData.subject_alt_names"
          multiple
          allow-create
          placeholder="请输入备用域名，如：www.example.com"
          style="width: 100%"
        >
          <el-option
            v-for="item in subjectAltNameOptions"
            :key="item"
            :label="item"
            :value="item"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="邮箱地址" prop="email" v-if="mode === 'add'">
        <el-input
          v-model="formData.email"
          placeholder="请输入用于证书申请的邮箱地址"
        />
      </el-form-item>

      <el-form-item label="DNS提供商" prop="provider_id" v-if="mode === 'add'">
        <el-select
          v-model="formData.provider_id"
          placeholder="请选择DNS提供商"
          style="width: 100%"
        >
          <el-option
            v-for="provider in providerOptions"
            :key="provider.id"
            :label="provider.name"
            :value="provider.id"
          />
        </el-select>
      </el-form-item>
            :value="item"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="邮箱地址" prop="email" v-if="mode === 'add'">
        <el-input
          v-model="formData.email"
          placeholder="请输入用于证书申请的邮箱地址"
        />
      </el-form-item>

      <el-form-item label="DNS提供商" prop="provider_id" v-if="mode === 'add'">
        <el-select
          v-model="formData.provider_id"
          placeholder="请选择DNS提供商"
          style="width: 100%"
        >
          <el-option
            v-for="provider in providerOptions"
            :key="provider.id"
            :label="provider.name"
            :value="provider.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="密钥类型" prop="key_type" v-if="mode === 'add'">
        <el-select
          v-model="formData.key_type"
          placeholder="请选择密钥类型"
          style="width: 100%"
        >
          <el-option label="RSA2048" value="RSA2048" />
          <el-option label="RSA4096" value="RSA4096" />
          <el-option label="EC256" value="EC256" />
          <el-option label="EC384" value="EC384" />
        </el-select>
      </el-form-item>

      <el-form-item label="有效期" prop="valid_days" v-if="mode === 'add'">
        <el-select
          v-model="formData.valid_days"
          placeholder="请选择证书有效期"
          style="width: 100%"
        >
          <el-option label="90天" :value="90" />
          <el-option label="180天" :value="180" />
          <el-option label="365天" :value="365" />
        </el-select>
      </el-form-item>

      <el-form-item label="自动续期" prop="auto_renew" v-if="mode === 'add'">
        <el-switch
          v-model="formData.auto_renew"
          active-text="开启"
          inactive-text="关闭"
        />
      </el-form-item>

      <el-form-item label="部署主机" prop="deploy_hosts" v-if="mode === 'add'">
        <el-select
          v-model="formData.deploy_hosts"
          multiple
          placeholder="请选择要部署证书的主机"
          style="width: 100%"
        >
          <el-option
            v-for="host in hostOptions"
            :key="host.id"
            :label="`${host.name} (${host.ip})`"
            :value="host.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="备注" prop="remark">
        <el-input
          v-model="formData.remark"
          type="textarea"
          :rows="3"
          placeholder="请输入备注信息"
        />
      </el-form-item>

      <!-- 详情模式显示的字段 -->
      <template v-if="mode === 'view'">
        <el-form-item label="证书状态">
          <el-tag :type="getCertStatusType(formData.status)">
            {{ getCertStatusText(formData.status) }}
          </el-tag>
        </el-form-item>

        <el-form-item label="序列号" v-if="formData.serial_number">
          <span>{{ formData.serial_number }}</span>
        </el-form-item>

        <el-form-item label="有效期" v-if="formData.not_before && formData.not_after">
          <span>{{ formatDate(formData.not_before) }} 至 {{ formatDate(formData.not_after) }}</span>
        </el-form-item>

        <el-form-item label="剩余天数" v-if="formData.not_after">
          <span :class="getDaysLeftClass(formData.not_after)">
            {{ getDaysLeft(formData.not_after) }} 天
          </span>
        </el-form-item>
      </template>
    </el-form>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button
          v-if="mode === 'add'"
          type="primary"
          :loading="loading"
          @click="handleOk"
        >
          {{ loading ? '申请中...' : '申请证书' }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, ElForm } from 'element-plus'
import type { Certificate, Domain } from '@/types/dns'

interface Provider {
  id: number
  name: string
  type: string
}

interface Host {
  id: number
  name: string
  ip: string
}

interface Props {
  visible: boolean
  certificate?: Certificate | null
  domainOptions?: Domain[]
  providerOptions?: Provider[]
  hostOptions?: Host[]
  mode?: 'add' | 'view'
}

const props = withDefaults(defineProps<Props>(), {
  certificate: null,
  domainOptions: () => [],
  providerOptions: () => [],
  hostOptions: () => [],
  mode: 'add'
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
  success: [certificate: Certificate]
}>()

// 响应式数据
const formRef = ref<InstanceType<typeof ElForm>>()
const loading = ref(false)
const subjectAltNameOptions = ref<string[]>([])

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const modalTitle = computed(() => {
  if (props.mode === 'view') {
    return '证书详情'
  }
  return props.certificate ? '编辑证书' : '申请证书'
})

// 表单数据
const formData = reactive({
  domain_id: undefined as number | undefined,
  common_name: '',
  subject_alt_names: [] as string[],
  email: '',
  provider_id: undefined as number | undefined,
  key_type: 'RSA2048',
  valid_days: 90,
  auto_renew: true,
  deploy_hosts: [] as number[],
  remark: '',
  status: '',
  serial_number: '',
  not_before: '',
  not_after: ''
})

// 表单验证规则
const rules = reactive({
  domain_id: [
    { required: true, message: '请选择域名', trigger: 'change' }
  ],
  common_name: [
    { required: true, message: '请输入主域名', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/, message: '请输入有效的域名格式', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  provider_id: [
    { required: true, message: '请选择DNS提供商', trigger: 'change' }
  ],
  key_type: [
    { required: true, message: '请选择密钥类型', trigger: 'change' }
  ],
  valid_days: [
    { required: true, message: '请选择有效期', trigger: 'change' }
  ]
})

// 工具方法
const getCertStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: 'warning',
    processing: 'info',
    issued: 'success',
    expired: 'danger',
    failed: 'danger',
    revoked: 'info'
  }
  return statusMap[status] || 'info'
}

const getCertStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待申请',
    processing: '申请中',
    issued: '已签发',
    expired: '已过期',
    failed: '申请失败',
    revoked: '已吊销'
  }
  return statusMap[status] || status
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const getDaysLeft = (expiresAt: string) => {
  if (!expiresAt) return 0
  const expires = new Date(expiresAt)
  const now = new Date()
  const diffTime = expires.getTime() - now.getTime()
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
}

const getDaysLeftClass = (expiresAt: string) => {
  const days = getDaysLeft(expiresAt)
  if (days <= 7) return 'text-red-500'
  if (days <= 30) return 'text-orange-500'
  return 'text-green-500'
}

// 处理方法
const handleOk = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true

    // 构建提交数据
    const submitData = {
      domain_id: formData.domain_id,
      domains: [formData.common_name, ...formData.subject_alt_names].filter(Boolean),
      email: formData.email,
      provider_id: formData.provider_id,
      key_type: formData.key_type,
      valid_days: formData.valid_days,
      auto_renew: formData.auto_renew,
      deploy_hosts: formData.deploy_hosts,
      remark: formData.remark
    }

    // 这里应该调用API申请证书
    console.log('申请证书数据:', submitData)
    
    // 模拟成功
    setTimeout(() => {
      loading.value = false
      ElMessage.success('证书申请成功')
      emit('success', submitData as any)
      dialogVisible.value = false
    }, 2000)

  } catch (error) {
    console.error('表单验证失败:', error)
    loading.value = false
  }
}

const handleCancel = () => {
  dialogVisible.value = false
  resetForm()
}

const resetForm = () => {
  Object.assign(formData, {
    domain_id: undefined,
    common_name: '',
    subject_alt_names: [],
    email: '',
    provider_id: undefined,
    key_type: 'RSA2048',
    valid_days: 90,
    auto_renew: true,
    deploy_hosts: [],
    remark: '',
    status: '',
    serial_number: '',
    not_before: '',
    not_after: ''
  })
  formRef.value?.clearValidate()
}

// 监听数据变化
watch(() => props.certificate, (newCert) => {
  if (newCert) {
    Object.assign(formData, {
      domain_id: newCert.domain_id,
      common_name: newCert.common_name,
      subject_alt_names: newCert.domains?.filter(d => d !== newCert.common_name) || [],
      email: newCert.email,
      provider_id: newCert.provider_id,
      key_type: newCert.key_type,
      valid_days: newCert.valid_days,
      auto_renew: newCert.auto_renew,
      deploy_hosts: newCert.deploy_hosts || [],
      remark: newCert.remark,
      status: newCert.status,
      serial_number: newCert.serial_number,
      not_before: newCert.not_before,
      not_after: newCert.not_after
    })
  }
}, { immediate: true })

watch(() => props.visible, (visible) => {
  if (!visible) {
    resetForm()
  }
})
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.text-red-500 {
  color: #f56565;
}

.text-orange-500 {
  color: #ed8936;
}

.text-green-500 {
  color: #48bb78;
}
</style>
