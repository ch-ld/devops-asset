<template>
  <el-dialog
    v-model="dialogVisible"
    :title="modalTitle"
    width="800px"
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
      <el-form-item label="域名" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入域名，如：example.com"
          :disabled="mode === 'view' || mode === 'edit'"
        />
      </el-form-item>

      <el-form-item label="DNS提供商" prop="provider_id">
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

      <el-form-item label="域名分组" prop="group_id">
        <el-select
          v-model="formData.group_id"
          placeholder="请选择域名分组"
          style="width: 100%"
          clearable
        >
          <el-option
            v-for="group in groups"
            :key="group.id"
            :label="group.name"
            :value="group.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="注册商类型" prop="registrar_type">
        <el-select
          v-model="formData.registrar_type"
          placeholder="请选择注册商类型"
          style="width: 100%"
        >
          <el-option label="阿里云" value="aliyun" />
          <el-option label="腾讯云" value="tencent" />
          <el-option label="AWS" value="aws" />
          <el-option label="GoDaddy" value="godaddy" />
          <el-option label="其他" value="other" />
        </el-select>
      </el-form-item>

      <el-form-item label="泛域名">
        <el-switch 
          v-model="formData.is_wildcard"
        />
        <span style="margin-left: 8px; color: #666; font-size: 12px;">
          是否为泛域名（*.example.com）
        </span>
      </el-form-item>

      <el-form-item label="自动续期">
        <el-switch 
          v-model="formData.auto_renew"
        />
        <span style="margin-left: 8px; color: #666; font-size: 12px;">
          域名即将过期时自动续期
        </span>
      </el-form-item>

      <el-form-item label="过期时间" prop="expired_at" v-if="formData.expired_at">
        <el-date-picker
          v-model="formData.expired_at"
          type="date"
          placeholder="选择过期时间"
          style="width: 100%"
        />
      </el-form-item>

      <el-form-item label="备注" prop="remark">
        <el-input
          v-model="formData.remark"
          type="textarea"
          :rows="3"
          placeholder="请输入备注信息"
        />
      </el-form-item>

      <!-- 查看模式下显示额外信息 -->
      <template v-if="mode === 'view' && domain">
        <el-form-item label="状态">
          <el-tag :type="getStatusType(domain.status)">
            {{ getStatusText(domain.status) }}
          </el-tag>
        </el-form-item>

        <el-form-item label="创建时间" v-if="domain.created_at">
          {{ formatDate(domain.created_at) }}
        </el-form-item>

        <el-form-item label="更新时间" v-if="domain.updated_at">
          {{ formatDate(domain.updated_at) }}
        </el-form-item>
      </template>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button 
          v-if="mode !== 'view'" 
          type="primary" 
          :loading="loading"
          @click="handleOk"
        >
          {{ mode === 'add' ? '创建' : '更新' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { domainApi } from '@/api/dns/domain'
import { dnsProviderApi } from '@/api/dns/provider'
import type { Domain, DomainGroup, DNSProvider } from '@/types/dns'

interface Props {
  visible: boolean
  mode?: 'add' | 'edit' | 'view'
  domain?: Domain | null
  groups?: DomainGroup[]
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'add',
  domain: null,
  groups: () => []
})

const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const formRef = ref()
const providerOptions = ref<DNSProvider[]>([])

const formData = reactive({
  name: '',
  registrar_id: undefined as number | undefined,
  group_id: undefined as number | undefined,
  registrar_type: '',
  auto_renew: false,
  expires_at: null as Date | null,
  remark: ''
})

// 对话框可见性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 模态框标题
const modalTitle = computed(() => {
  const titles = {
    add: '添加域名',
    edit: '编辑域名',
    view: '查看域名'
  }
  return titles[props.mode]
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入域名', trigger: 'blur' },
    { 
      pattern: /^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/,
      message: '请输入有效的域名格式',
      trigger: 'blur'
    }
  ],
  provider_id: [
    { required: true, message: '请选择DNS提供商', trigger: 'change' }
  ],
  registrar_type: [
    { required: true, message: '请选择注册商类型', trigger: 'change' }
  ]
}

// 工具方法
const getStatusType = (status: string) => {
  const statusMap: Record<string, any> = {
    active: 'success',
    inactive: 'info',
    expired: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '正常',
    inactive: '停用',
    expired: '已过期'
  }
  return statusMap[status] || status
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

// 事件处理
const handleOk = async () => {
  try {
    await formRef.value?.validate()
    
    loading.value = true
    
    const data = {
      name: formData.name,
      provider_id: formData.provider_id!,
      group_id: formData.group_id,
      registrar_type: formData.registrar_type,
      is_wildcard: formData.is_wildcard,
      auto_renew: formData.auto_renew,
      expired_at: formData.expired_at,
      remark: formData.remark
    }
    
    if (props.mode === 'add') {
      await domainApi.create(data)
      ElMessage.success('域名创建成功')
    } else {
      await domainApi.update(props.domain!.id, data)
      ElMessage.success('域名更新成功')
    }
    
    emit('success')
  } catch (error) {
    ElMessage.error(props.mode === 'add' ? '创建失败' : '更新失败')
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  emit('update:visible', false)
}

// 重置表单
const resetForm = () => {
  formData.name = ''
  formData.provider_id = undefined
  formData.group_id = undefined
  formData.registrar_type = ''
  formData.is_wildcard = false
  formData.auto_renew = false
  formData.expired_at = null
  formData.remark = ''
  
  formRef.value?.clearValidate()
}

// 获取DNS提供商列表
const fetchProviders = async () => {
  try {
    const response = await providerApi.list({ page: 1, size: 100 })
    providerOptions.value = response.list || []
  } catch (error) {
    console.error('获取DNS提供商列表失败:', error)
  }
}

// 监听props变化
watch(() => props.visible, (visible) => {
  if (visible) {
    fetchProviders()
    if (props.mode === 'add') {
      resetForm()
    } else if (props.domain) {
      // 编辑或查看模式，填充表单数据
      formData.name = props.domain.name
      formData.provider_id = props.domain.provider_id
      formData.group_id = props.domain.group_id
      formData.registrar_type = props.domain.registrar_type
      formData.is_wildcard = props.domain.is_wildcard || false
      formData.auto_renew = props.domain.auto_renew || false
      formData.expired_at = props.domain.expired_at ? new Date(props.domain.expired_at) : null
      formData.remark = props.domain.remark || ''
    }
  }
})

// 初始化
onMounted(() => {
  fetchProviders()
})
</script>

<style scoped lang="scss">
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
