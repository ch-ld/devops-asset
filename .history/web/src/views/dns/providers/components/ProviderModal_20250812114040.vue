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
    <!-- 顶部提供商选择卡片 -->
    <div class="provider-selection-card">
      <h3 class="card-title">
        <el-icon><Setting /></el-icon>
        基本信息
      </h3>
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="140px"
        class="provider-form"
      >
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="提供商名称" prop="name">
              <el-input
                v-model="formData.name"
                placeholder="请输入提供商名称"
                :disabled="mode === 'view'"
                size="large"
                :prefix-icon="Edit"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="提供商类型" prop="type">
              <el-select
                v-model="formData.type"
                placeholder="请选择提供商类型"
                :disabled="mode === 'view'"
                @change="handleTypeChange"
                size="large"
                style="width: 100%"
              >
                <el-option label="阿里云" value="aliyun">
                  <div class="provider-option">
                    <el-icon><Cloudy /></el-icon>
                    <span>阿里云</span>
                  </div>
                </el-option>
                <el-option label="腾讯云" value="tencent">
                  <div class="provider-option">
                    <el-icon><Cloudy /></el-icon>
                    <span>腾讯云</span>
                  </div>
                </el-option>
                <el-option label="Cloudflare" value="cloudflare">
                  <div class="provider-option">
                    <el-icon><Lightning /></el-icon>
                    <span>Cloudflare</span>
                  </div>
                </el-option>
                <el-option label="DNSPod" value="dnspod">
                  <div class="provider-option">
                    <el-icon><Monitor /></el-icon>
                    <span>DNSPod</span>
                  </div>
                </el-option>
                <el-option label="GoDaddy" value="godaddy">
                  <div class="provider-option">
                    <el-icon><Setting /></el-icon>
                    <span>GoDaddy</span>
                  </div>
                </el-option>
                <el-option label="AWS Route53" value="route53">
                  <div class="provider-option">
                    <el-icon><Connection /></el-icon>
                    <span>AWS Route53</span>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="备注描述" prop="remark">
          <el-input
            v-model="formData.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注描述信息"
            :disabled="mode === 'view'"
            show-word-limit
            maxlength="200"
          />
        </el-form-item>
      </el-form>
    </div>

    <!-- API配置卡片 -->
    <div class="api-config-card" v-if="formData.type">
      <h3 class="card-title">
        <el-icon><Key /></el-icon>
        API配置 - {{ getProviderDisplayName(formData.type) }}
      </h3>

      <div class="provider-info-banner">
        <el-alert
          :title="`请确保您的 ${getProviderDisplayName(formData.type)} API凭证具有DNS管理权限`"
          type="info"
          show-icon
          :closable="false"
        />
      </div>

      <template v-if="formData.type === 'aliyun'">
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="Access Key ID" prop="accessKeyId">
              <el-input
                v-model="formData.accessKeyId"
                placeholder="请输入阿里云Access Key ID"
                :disabled="mode === 'view'"
                show-password
                size="large"
                :prefix-icon="Key"
              />
              <div class="field-hint">通常以LTAI开头，16-24位字符</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Access Key Secret" prop="accessKeySecret">
              <el-input
                v-model="formData.accessKeySecret"
                placeholder="请输入阿里云Access Key Secret"
                :disabled="mode === 'view'"
                show-password
                size="large"
                :prefix-icon="Lock"
              />
              <div class="field-hint">30位字符的密钥</div>
            </el-form-item>
          </el-col>
        </el-row>
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
    </div>

    <!-- 其他设置卡片 -->
    <div class="settings-card">
      <h3 class="card-title">
        <el-icon><Setting /></el-icon>
        其他设置
      </h3>

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
import { 
  Edit, Setting, Key, Lock, Cloudy, Lightning, Monitor, Connection,
  CircleCheck, Warning
} from '@element-plus/icons-vue'
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
// 获取提供商显示名称
const getProviderDisplayName = (type: string) => {
  const nameMap: Record<string, string> = {
    aliyun: '阿里云',
    tencent: '腾讯云',
    cloudflare: 'Cloudflare',
    dnspod: 'DNSPod',
    godaddy: 'GoDaddy',
    route53: 'AWS Route53'
  }
  return nameMap[type] || type
}
</script>

<style scoped>
/* 提供商模态框整体样式 */
:deep(.provider-modal .el-dialog) {
  border-radius: 20px;
  background: white;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.8);
}

:deep(.provider-modal .el-dialog__header) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 24px 32px;
  border-radius: 20px 20px 0 0;
  border-bottom: none;
}

:deep(.provider-modal .el-dialog__title) {
  color: white;
  font-size: 20px;
  font-weight: 700;
}

:deep(.provider-modal .el-dialog__headerbtn .el-dialog__close) {
  color: white;
  font-size: 18px;
}

:deep(.provider-modal .el-dialog__body) {
  padding: 32px;
  background: #f8f9fa;
}

/* 卡片样式 */
.provider-selection-card,
.api-config-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  border: 1px solid #e9ecef;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 600;
  color: #343a40;
  padding-bottom: 12px;
  border-bottom: 2px solid #e9ecef;
}

.card-title .el-icon {
  color: #667eea;
  font-size: 20px;
}

/* 提供商选项样式 */
.provider-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.provider-option .el-icon {
  color: #667eea;
}

/* 信息横幅 */
.provider-info-banner {
  margin-bottom: 20px;
}

.provider-info-banner :deep(.el-alert) {
  border-radius: 12px;
  border: none;
  background: linear-gradient(135deg, #e3f2fd 0%, #f3e5f5 100%);
}

/* 表单样式优化 */
.provider-form :deep(.el-form-item__label) {
  font-weight: 600;
  color: #495057;
}

.provider-form :deep(.el-input__wrapper) {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  border: 1px solid #e9ecef;
  transition: all 0.3s ease;
}

.provider-form :deep(.el-input__wrapper:hover) {
  border-color: #667eea;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
}

.provider-form :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.25);
}

.provider-form :deep(.el-select .el-input__wrapper) {
  border-radius: 12px;
}

.provider-form :deep(.el-textarea__inner) {
  border-radius: 12px;
  border: 1px solid #e9ecef;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  resize: vertical;
}

.provider-form :deep(.el-textarea__inner:hover) {
  border-color: #667eea;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
}

.provider-form :deep(.el-textarea__inner:focus) {
  border-color: #667eea;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.25);
}

/* 字段提示样式 */
.field-hint {
  font-size: 12px;
  color: #6c757d;
  margin-top: 4px;
  padding-left: 12px;
}

/* 底部按钮样式 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 24px 32px;
  background: white;
  border-radius: 0 0 20px 20px;
  border-top: 1px solid #e9ecef;
}

.dialog-footer .el-button {
  min-width: 100px;
  padding: 12px 24px;
  border-radius: 10px;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.dialog-footer .el-button--primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.dialog-footer .el-button--primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.6);
}

.dialog-footer .el-button--success {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(79, 172, 254, 0.4);
}

.dialog-footer .el-button--success:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(79, 172, 254, 0.6);
}

.dialog-footer .el-button:not(.el-button--primary):not(.el-button--success) {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  color: #495057;
  border: 1px solid #dee2e6;
}

.dialog-footer .el-button:not(.el-button--primary):not(.el-button--success):hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* 下拉框选项样式 */
:deep(.el-select-dropdown__item) {
  padding: 12px 16px;
}

:deep(.el-select-dropdown__item:hover) {
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  color: #667eea;
}

/* 响应式设计 */
@media (max-width: 768px) {
  :deep(.provider-modal .el-dialog) {
    width: 95% !important;
    margin: 5vh auto;
  }
  
  :deep(.provider-modal .el-dialog__body) {
    padding: 20px;
  }
  
  .provider-selection-card,
  .api-config-card {
    padding: 16px;
  }
  
  .dialog-footer {
    padding: 16px 20px;
    flex-direction: column;
  }
  
  .dialog-footer .el-button {
    width: 100%;
  }
}

/* 动画效果 */
@keyframes slideInDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.provider-selection-card,
.api-config-card {
  animation: slideInDown 0.3s ease-out;
}
</style>
