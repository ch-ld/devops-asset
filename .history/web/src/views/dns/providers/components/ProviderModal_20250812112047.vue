<template>
  <el-dialog
    v-model="dialogVisible"
    :title="modalTitle"
    width="900px"
    :close-on-click-modal="false"
    @close="handleCancel"
    class="provider-modal"
    :append-to-body="true"
    destroy-on-close
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
          :disabled="mode === 'view'"
        />
      </el-form-item>

      <el-form-item label="提供商类型" prop="type">
        <el-select
          v-model="formData.type"
          placeholder="请选择提供商类型"
          :disabled="mode === 'view'"
          @change="handleTypeChange"
        >
          <el-option label="阿里云" value="aliyun" />
          <el-option label="腾讯云" value="tencent" />
          <el-option label="Cloudflare" value="cloudflare" />
          <el-option label="DNSPod" value="dnspod" />
          <el-option label="GoDaddy" value="godaddy" />
          <el-option label="AWS Route53" value="route53" />
        </el-select>
      </el-form-item>

      <el-form-item label="描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="2"
          placeholder="请输入描述信息"
          :disabled="mode === 'view'"
        />
      </el-form-item>

      <!-- API配置 -->
      <el-divider content-position="left">API配置</el-divider>

      <template v-if="formData.type === 'aliyun'">
        <el-form-item label="Access Key ID" prop="accessKeyId">
          <el-input
            v-model="formData.accessKeyId"
            placeholder="请输入阿里云Access Key ID"
            :disabled="mode === 'view'"
            show-password
          />
        </el-form-item>
        <el-form-item label="Access Key Secret" prop="accessKeySecret">
          <el-input
            v-model="formData.accessKeySecret"
            placeholder="请输入阿里云Access Key Secret"
            :disabled="mode === 'view'"
            show-password
          />
        </el-form-item>
      </template>

      <template v-else-if="formData.type === 'cloudflare'">
        <el-form-item label="API Token" prop="apiToken">
          <el-input
            v-model="formData.apiToken"
            placeholder="请输入Cloudflare API Token"
            :disabled="mode === 'view'"
            show-password
          />
        </el-form-item>
        <el-form-item label="Zone ID" prop="zoneId">
          <el-input
            v-model="formData.zoneId"
            placeholder="请输入Zone ID（可选）"
            :disabled="mode === 'view'"
          />
        </el-form-item>
      </template>

      <template v-else-if="formData.type === 'route53'">
        <el-form-item label="Access Key ID" prop="accessKeyId">
          <el-input
            v-model="formData.accessKeyId"
            placeholder="请输入AWS Access Key ID"
            :disabled="mode === 'view'"
            show-password
          />
        </el-form-item>
        <el-form-item label="Secret Access Key" prop="secretAccessKey">
          <el-input
            v-model="formData.secretAccessKey"
            placeholder="请输入AWS Secret Access Key"
            :disabled="mode === 'view'"
            show-password
          />
        </el-form-item>
        <el-form-item label="Region" prop="region">
          <el-select v-model="formData.region" placeholder="请选择区域">
            <el-option label="us-east-1" value="us-east-1" />
            <el-option label="us-west-2" value="us-west-2" />
            <el-option label="eu-west-1" value="eu-west-1" />
            <el-option label="ap-southeast-1" value="ap-southeast-1" />
          </el-select>
        </el-form-item>
      </template>

      <template v-else-if="formData.type === 'tencent'">
        <el-form-item label="Secret ID" prop="accessKeyId">
          <el-input
            v-model="formData.accessKeyId"
            placeholder="请输入腾讯云Secret ID"
            :disabled="mode === 'view'"
            show-password
          />
        </el-form-item>
        <el-form-item label="Secret Key" prop="accessKeySecret">
          <el-input
            v-model="formData.accessKeySecret"
            placeholder="请输入腾讯云Secret Key"
            :disabled="mode === 'view'"
            show-password
          />
        </el-form-item>
        <el-form-item label="区域" prop="region">
          <el-select v-model="formData.region" placeholder="请选择区域">
            <el-option label="ap-beijing" value="ap-beijing" />
            <el-option label="ap-shanghai" value="ap-shanghai" />
            <el-option label="ap-guangzhou" value="ap-guangzhou" />
            <el-option label="ap-hongkong" value="ap-hongkong" />
          </el-select>
        </el-form-item>
      </template>

      <template v-else-if="formData.type === 'dnspod'">
        <el-form-item label="API Token" prop="apiToken">
          <el-input
            v-model="formData.apiToken"
            placeholder="请输入DNSPod API Token"
            :disabled="mode === 'view'"
            show-password
          />
        </el-form-item>
      </template>

      <template v-else-if="formData.type === 'godaddy'">
        <el-form-item label="API Key" prop="apiKey">
          <el-input
            v-model="formData.apiKey"
            placeholder="请输入GoDaddy API Key"
            :disabled="mode === 'view'"
            show-password
          />
        </el-form-item>
        <el-form-item label="API Secret" prop="apiSecret">
          <el-input
            v-model="formData.apiSecret"
            placeholder="请输入GoDaddy API Secret"
            :disabled="mode === 'view'"
            show-password
          />
        </el-form-item>
      </template>

      <template v-else-if="formData.type">
        <el-form-item label="API Key" prop="apiKey">
          <el-input
            v-model="formData.apiKey"
            placeholder="请输入API Key"
            :disabled="mode === 'view'"
            show-password
          />
        </el-form-item>
        <el-form-item label="API Secret" prop="apiSecret">
          <el-input
            v-model="formData.apiSecret"
            placeholder="请输入API Secret"
            :disabled="mode === 'view'"
            show-password
          />
        </el-form-item>
      </template>

      <!-- 其他设置 -->
      <el-divider content-position="left">其他设置</el-divider>

      <el-form-item label="状态">
        <el-switch
          v-model="formData.enabled"
          active-text="启用"
          inactive-text="禁用"
          :disabled="mode === 'view'"
        />
      </el-form-item>

      <el-form-item label="默认提供商">
        <el-switch
          v-model="formData.isDefault"
          active-text="是"
          inactive-text="否"
          :disabled="mode === 'view'"
        />
      </el-form-item>

      <el-form-item label="备注">
        <el-input
          v-model="formData.remark"
          type="textarea"
          :rows="3"
          placeholder="请输入备注信息"
          :disabled="mode === 'view'"
          maxlength="500"
          show-word-limit
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleCancel">{{ mode === 'view' ? '关闭' : '取消' }}</el-button>
        <el-button 
          v-if="mode !== 'view'" 
          @click="handleTest" 
          :loading="testLoading"
        >
          测试连接
        </el-button>
        <el-button 
          v-if="mode !== 'view'" 
          type="primary" 
          :loading="loading" 
          @click="handleOk"
        >
          {{ mode === 'add' ? '创建' : '保存' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, ElForm } from 'element-plus'
import { dnsProviderApi } from '@/api/dns/provider'

interface Props {
  visible: boolean
  mode: 'add' | 'edit' | 'view'
  provider?: any
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  mode: 'add',
  provider: null
})

const emit = defineEmits<Emits>()

const formRef = ref<InstanceType<typeof ElForm>>()
const loading = ref(false)
const testLoading = ref(false)

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const modalTitle = computed(() => {
  const titleMap = {
    add: '添加DNS提供商',
    edit: '编辑DNS提供商',
    view: 'DNS提供商详情'
  }
  return titleMap[props.mode]
})

// 表单数据
const formData = reactive({
  name: '',
  type: '',
  description: '',
  accessKeyId: '',
  accessKeySecret: '',
  apiToken: '',
  zoneId: '',
  secretAccessKey: '',
  region: '',
  apiKey: '',
  apiSecret: '',
  enabled: true,
  isDefault: false,
  remark: ''
})

// 表单验证规则
const rules = reactive<Record<string, any>>({
  name: [
    { required: true, message: '请输入提供商名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度为2-50个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择提供商类型', trigger: 'change' }
  ]
})

// 处理类型变化
const handleTypeChange = (type: string) => {
  // 清空API配置
  Object.assign(formData, {
    accessKeyId: '',
    accessKeySecret: '',
    apiToken: '',
    zoneId: '',
    secretAccessKey: '',
    region: '',
    apiKey: '',
    apiSecret: ''
  })

  // 根据类型设置验证规则
  updateRules(type)
}

// 更新验证规则
const updateRules = (type: string) => {
  // 重置规则
  delete rules.accessKeyId
  delete rules.accessKeySecret
  delete rules.apiToken
  delete rules.secretAccessKey
  delete rules.apiKey
  delete rules.apiSecret

  switch (type) {
    case 'aliyun':
      rules.accessKeyId = [
        { required: true, message: '请输入Access Key ID', trigger: 'blur' }
      ]
      rules.accessKeySecret = [
        { required: true, message: '请输入Access Key Secret', trigger: 'blur' }
      ]
      break
    case 'cloudflare':
      rules.apiToken = [
        { required: true, message: '请输入API Token', trigger: 'blur' }
      ]
      break
    case 'route53':
      rules.accessKeyId = [
        { required: true, message: '请输入Access Key ID', trigger: 'blur' }
      ]
      rules.secretAccessKey = [
        { required: true, message: '请输入Secret Access Key', trigger: 'blur' }
      ]
      break
    default:
      rules.apiKey = [
        { required: true, message: '请输入API Key', trigger: 'blur' }
      ]
      rules.apiSecret = [
        { required: true, message: '请输入API Secret', trigger: 'blur' }
      ]
  }
}

// 测试连接
const handleTest = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    testLoading.value = true

    // 如果是编辑模式且有ID，直接测试现有提供商
    if (props.mode === 'edit' && props.provider?.id) {
      await dnsProviderApi.test(props.provider.id)
    } else {
      // 新建模式，调用临时测试接口
      const testData = {
        type: formData.type,
        credentials: buildCredentials()
      }

      // 调用临时测试接口
      await dnsProviderApi.testConnection(testData)
    }

    ElMessage.success('连接测试成功')
  } catch (error: any) {
    if (error.fields) {
      ElMessage.warning('请先完善配置信息')
      return
    }
    console.error('连接测试失败:', error)
    ElMessage.error(error.message || '连接测试失败')
  } finally {
    testLoading.value = false
  }
}

// 重置表单
const resetForm = () => {
  Object.assign(formData, {
    name: '',
    type: '',
    description: '',
    accessKeyId: '',
    accessKeySecret: '',
    apiToken: '',
    zoneId: '',
    secretAccessKey: '',
    region: '',
    apiKey: '',
    apiSecret: '',
    enabled: true,
    isDefault: false,
    remark: ''
  })
  formRef.value?.resetFields()
}

// 处理确认
const handleOk = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true

    // 构建提交数据
    const submitData = {
      name: formData.name,
      type: formData.type,
      credentials: buildCredentials(),
      remark: formData.remark
    }

    if (props.mode === 'add') {
      await dnsProviderApi.create(submitData)
      ElMessage.success('提供商创建成功')
    } else if (props.provider?.id) {
      await dnsProviderApi.update(props.provider.id, submitData)
      ElMessage.success('提供商更新成功')
    }

    emit('success')
    handleCancel()
  } catch (error: any) {
    if (error.fields) {
      // 表单验证错误
      return
    }
    console.error('保存提供商失败:', error)
    ElMessage.error(props.mode === 'add' ? '提供商创建失败' : '提供商更新失败')
  } finally {
    loading.value = false
  }
}

// 构建凭证信息
const buildCredentials = () => {
  const credentials: Record<string, string> = {}

  switch (formData.type) {
    case 'aliyun':
      credentials.access_key_id = formData.accessKeyId
      credentials.access_key_secret = formData.accessKeySecret
      credentials.region = formData.region || 'cn-hangzhou'
      break
    case 'cloudflare':
      credentials.api_token = formData.apiToken
      if (formData.zoneId) {
        credentials.zone_id = formData.zoneId
      }
      break
    case 'route53':
      credentials.access_key_id = formData.accessKeyId
      credentials.secret_access_key = formData.secretAccessKey
      credentials.region = formData.region || 'us-east-1'
      break
    case 'tencent':
      credentials.secret_id = formData.accessKeyId
      credentials.secret_key = formData.accessKeySecret
      credentials.region = formData.region || 'ap-beijing'
      break
    case 'dnspod':
      credentials.api_token = formData.apiToken
      break
    case 'godaddy':
      credentials.api_key = formData.apiKey
      credentials.api_secret = formData.apiSecret
      break
    default:
      credentials.api_key = formData.apiKey
      credentials.api_secret = formData.apiSecret
  }

  return credentials
}

// 处理取消
const handleCancel = () => {
  dialogVisible.value = false
  resetForm()
}

// 监听编辑数据变化
watch(
  () => props.provider,
  (newVal) => {
    if (newVal && props.mode !== 'add') {
      Object.assign(formData, {
        name: newVal.name || '',
        type: newVal.type || '',
        description: newVal.description || '',
        accessKeyId: newVal.accessKeyId || '',
        accessKeySecret: newVal.accessKeySecret || '',
        apiToken: newVal.apiToken || '',
        zoneId: newVal.zoneId || '',
        secretAccessKey: newVal.secretAccessKey || '',
        region: newVal.region || '',
        apiKey: newVal.apiKey || '',
        apiSecret: newVal.apiSecret || '',
        enabled: newVal.enabled ?? true,
        isDefault: newVal.isDefault ?? false,
        remark: newVal.remark || ''
      })
      
      if (newVal.type) {
        updateRules(newVal.type)
      }
    }
  },
  { immediate: true }
)

// 监听弹窗显示状态
watch(
  () => props.visible,
  (visible) => {
    if (visible && props.mode === 'add') {
      resetForm()
    }
  }
)
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
