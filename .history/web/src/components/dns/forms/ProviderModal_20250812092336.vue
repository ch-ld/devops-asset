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
    >
      <el-form-item label="提供商名称" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入提供商名称"
        />
      </el-form-item>

      <el-form-item label="提供商类型" prop="type">
        <el-select
          v-model="formData.type"
          placeholder="请选择提供商类型"
          style="width: 100%"
          @change="handleTypeChange"
        >
          <el-option label="阿里云DNS" value="aliyun" />
          <el-option label="AWS Route53" value="aws" />
          <el-option label="Cloudflare" value="cloudflare" />
          <el-option label="腾讯云DNS" value="tencent" />
          <el-option label="DNSPod" value="dnspod" />
          <el-option label="GoDaddy" value="godaddy" />
        </el-select>
      </el-form-item>

      <!-- 阿里云配置 -->
      <template v-if="formData.type === 'aliyun'">
        <el-form-item label="Access Key ID" prop="access_key_id">
          <el-input
            v-model="formData.access_key_id"
            placeholder="请输入Access Key ID"
            show-password
          />
        </el-form-item>
        <el-form-item label="Access Key Secret" prop="access_key_secret">
          <el-input
            v-model="formData.access_key_secret"
            type="password"
            placeholder="请输入Access Key Secret"
            show-password
          />
        </el-form-item>
        <el-form-item label="区域" prop="region">
          <el-select
            v-model="formData.region"
            placeholder="请选择区域"
            style="width: 100%"
          >
            <el-option label="cn-hangzhou" value="cn-hangzhou" />
            <el-option label="cn-beijing" value="cn-beijing" />
            <el-option label="cn-shanghai" value="cn-shanghai" />
            <el-option label="cn-shenzhen" value="cn-shenzhen" />
          </el-select>
        </el-form-item>
      </template>

      <!-- AWS Route53配置 -->
      <template v-if="formData.type === 'aws'">
        <el-form-item label="Access Key ID" prop="access_key_id">
          <el-input
            v-model="formData.access_key_id"
            placeholder="请输入AWS Access Key ID"
            show-password
          />
        </el-form-item>
        <el-form-item label="Secret Access Key" prop="secret_access_key">
          <el-input
            v-model="formData.secret_access_key"
            type="password"
            placeholder="请输入AWS Secret Access Key"
            show-password
          />
        </el-form-item>
        <el-form-item label="区域" prop="region">
          <el-select
            v-model="formData.region"
            placeholder="请选择AWS区域"
            style="width: 100%"
          >
            <el-option label="us-east-1" value="us-east-1" />
            <el-option label="us-west-2" value="us-west-2" />
            <el-option label="ap-southeast-1" value="ap-southeast-1" />
            <el-option label="eu-west-1" value="eu-west-1" />
          </el-select>
        </el-form-item>
      </template>

      <!-- Cloudflare配置 -->
      <template v-if="formData.type === 'cloudflare'">
        <el-form-item label="API Token" prop="api_token">
          <el-input
            v-model="formData.api_token"
            type="password"
            placeholder="请输入Cloudflare API Token"
            show-password
          />
        </el-form-item>
        <el-form-item label="Zone ID" prop="zone_id">
          <el-input
            v-model="formData.zone_id"
            placeholder="请输入Zone ID（可选）"
          />
        </el-form-item>
      </template>

      <!-- 腾讯云配置 -->
      <template v-if="formData.type === 'tencent'">
        <el-form-item label="Secret ID" prop="secret_id">
          <el-input
            v-model="formData.secret_id"
            placeholder="请输入Secret ID"
            show-password
          />
        </el-form-item>
        <el-form-item label="Secret Key" prop="secret_key">
          <el-input
            v-model="formData.secret_key"
            type="password"
            placeholder="请输入Secret Key"
            show-password
          />
        </el-form-item>
        <el-form-item label="区域" prop="region">
          <el-select
            v-model="formData.region"
            placeholder="请选择区域"
            style="width: 100%"
          >
            <el-option label="ap-beijing" value="ap-beijing" />
            <el-option label="ap-shanghai" value="ap-shanghai" />
            <el-option label="ap-guangzhou" value="ap-guangzhou" />
          </el-select>
        </el-form-item>
      </template>

      <!-- DNSPod配置 -->
      <template v-if="formData.type === 'dnspod'">
        <el-form-item label="API Token" prop="api_token">
          <el-input
            v-model="formData.api_token"
            type="password"
            placeholder="请输入DNSPod API Token"
            show-password
          />
        </el-form-item>
      </template>

      <!-- GoDaddy配置 -->
      <template v-if="formData.type === 'godaddy'">
        <el-form-item label="API Key" prop="api_key">
          <el-input
            v-model="formData.api_key"
            placeholder="请输入GoDaddy API Key"
            show-password
          />
        </el-form-item>
        <el-form-item label="API Secret" prop="api_secret">
          <el-input
            v-model="formData.api_secret"
            type="password"
            placeholder="请输入GoDaddy API Secret"
            show-password
          />
        </el-form-item>
      </template>

      <el-form-item label="备注" prop="remark">
        <el-input
          v-model="formData.remark"
          type="textarea"
          :rows="3"
          placeholder="请输入备注信息"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button 
          type="default" 
          :loading="testLoading"
          @click="handleTestConnection"
        >
          测试连接
        </el-button>
        <el-button 
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
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { providerApi } from '@/api/dns/provider'
import type { DnsProvider } from '@/types/dns'

interface Props {
  visible: boolean
  mode?: 'add' | 'edit' | 'view'
  provider?: DnsProvider | null
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'add',
  provider: null
})

const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const testLoading = ref(false)
const formRef = ref()

const formData = reactive({
  name: '',
  type: '',
  access_key_id: '',
  access_key_secret: '',
  secret_access_key: '',
  api_token: '',
  zone_id: '',
  secret_id: '',
  secret_key: '',
  api_key: '',
  api_secret: '',
  region: '',
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
    add: '添加DNS提供商',
    edit: '编辑DNS提供商',
    view: '查看DNS提供商'
  }
  return titles[props.mode]
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入提供商名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择提供商类型', trigger: 'change' }
  ],
  access_key_id: [
    { required: true, message: '请输入Access Key ID', trigger: 'blur' }
  ],
  access_key_secret: [
    { required: true, message: '请输入Access Key Secret', trigger: 'blur' }
  ],
  secret_access_key: [
    { required: true, message: '请输入Secret Access Key', trigger: 'blur' }
  ],
  api_token: [
    { required: true, message: '请输入API Token', trigger: 'blur' }
  ],
  secret_id: [
    { required: true, message: '请输入Secret ID', trigger: 'blur' }
  ],
  secret_key: [
    { required: true, message: '请输入Secret Key', trigger: 'blur' }
  ],
  api_key: [
    { required: true, message: '请输入API Key', trigger: 'blur' }
  ],
  api_secret: [
    { required: true, message: '请输入API Secret', trigger: 'blur' }
  ]
}

// 事件处理
const handleTypeChange = () => {
  // 清空所有凭证信息
  formData.access_key_id = ''
  formData.access_key_secret = ''
  formData.secret_access_key = ''
  formData.api_token = ''
  formData.zone_id = ''
  formData.secret_id = ''
  formData.secret_key = ''
  formData.api_key = ''
  formData.api_secret = ''
  formData.region = ''
}

const handleTestConnection = async () => {
  try {
    await formRef.value?.validate()
    
    testLoading.value = true
    
    const credentials = getCredentials()
    
    await providerApi.testConnection({
      type: formData.type,
      credentials
    })
    
    ElMessage.success('连接测试成功')
  } catch (error) {
    console.error('连接测试失败:', error)
    ElMessage.error('连接测试失败')
  } finally {
    testLoading.value = false
  }
}

const handleOk = async () => {
  try {
    await formRef.value?.validate()
    
    loading.value = true
    
    const data = {
      name: formData.name,
      type: formData.type,
      credentials: getCredentials(),
      remark: formData.remark
    }
    
    if (props.mode === 'add') {
      await providerApi.create(data)
      ElMessage.success('DNS提供商创建成功')
    } else if (props.provider) {
      await providerApi.update(props.provider.id, data)
      ElMessage.success('DNS提供商更新成功')
    }
    
    emit('success')
    emit('update:visible', false)
  } catch (error) {
    console.error('操作失败:', error)
    ElMessage.error(props.mode === 'add' ? '创建失败' : '更新失败')
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  emit('update:visible', false)
}

// 工具方法
const getCredentials = () => {
  const credentials: any = {}
  
  switch (formData.type) {
    case 'aliyun':
      credentials.access_key_id = formData.access_key_id
      credentials.access_key_secret = formData.access_key_secret
      credentials.region = formData.region
      break
    case 'aws':
      credentials.access_key_id = formData.access_key_id
      credentials.secret_access_key = formData.secret_access_key
      credentials.region = formData.region
      break
    case 'cloudflare':
      credentials.api_token = formData.api_token
      if (formData.zone_id) {
        credentials.zone_id = formData.zone_id
      }
      break
    case 'tencent':
      credentials.secret_id = formData.secret_id
      credentials.secret_key = formData.secret_key
      credentials.region = formData.region
      break
    case 'dnspod':
      credentials.api_token = formData.api_token
      break
    case 'godaddy':
      credentials.api_key = formData.api_key
      credentials.api_secret = formData.api_secret
      break
  }
  
  return credentials
}

// 重置表单
const resetForm = () => {
  formData.name = ''
  formData.type = ''
  formData.access_key_id = ''
  formData.access_key_secret = ''
  formData.secret_access_key = ''
  formData.api_token = ''
  formData.zone_id = ''
  formData.secret_id = ''
  formData.secret_key = ''
  formData.api_key = ''
  formData.api_secret = ''
  formData.region = ''
  formData.remark = ''
  
  formRef.value?.clearValidate()
}

// 监听props变化
watch(() => props.visible, (visible) => {
  if (visible) {
    if (props.mode === 'add') {
      resetForm()
    } else if (props.provider) {
      // 编辑模式，填充表单数据
      formData.name = props.provider.name
      formData.type = props.provider.type
      formData.remark = props.provider.remark || ''
      
      // 注意：出于安全考虑，凭证信息不回显
    }
  }
})
</script>

<style scoped lang="scss">
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
