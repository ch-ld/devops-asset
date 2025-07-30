<template>
  <el-dialog
    v-model="dialogVisible"
    :title="null"
    width="720px"
    :close-on-click-modal="false"
    @close="handleCancel"
    class="modern-provider-modal"
    destroy-on-close
    :show-close="false"
  >
    <!-- 自定义头部 -->
    <template #header>
      <div class="modal-header">
        <div class="header-left">
          <div class="header-icon">
            <el-icon size="28" color="#1890ff">☁️</el-icon>
          </div>
          <div class="header-content">
            <h2 class="header-title">{{ isEdit ? '编辑云账号' : '添加云账号' }}</h2>
            <p class="header-subtitle">配置云厂商账号信息，用于资源管理和监控</p>
          </div>
        </div>
        <div class="header-right">
          <el-button 
                            type="text" 
            @click="handleCancel"
            class="close-btn"
          >
            <el-icon size="20">✕</el-icon>
          </el-button>
        </div>
      </div>
    </template>

    <div class="modal-body">
      <el-form 
        ref="formRef" 
        :model="formState" 
        :rules="rules" 
        class="modern-form"
        :label-position="'top'"
      >
        <!-- 步骤指示器 -->
        <div class="steps-container">
          <el-steps :active="currentStep" align-center>
            <el-step title="选择云厂商" />
            <el-step title="配置认证" />
            <el-step title="完成设置" />
          </el-steps>
        </div>

        <!-- 步骤1: 选择云厂商 -->
        <div v-show="currentStep === 0" class="step-content">
          <div class="step-title">
            <h3>选择云厂商类型</h3>
            <p>请选择您要添加的云厂商类型</p>
          </div>
          
          <div class="provider-grid">
            <div 
              v-for="provider in providerOptions" 
              :key="provider.value"
              class="provider-card"
              :class="{ active: formState.type === provider.value }"
              @click="selectProvider(provider.value)"
            >
              <div class="provider-icon">
                <img :src="provider.icon" :alt="provider.label" />
              </div>
              <div class="provider-info">
                <h4>{{ provider.label }}</h4>
                <p>{{ provider.description }}</p>
              </div>
              <div class="provider-check">
                <el-icon v-if="formState.type === provider.value" color="#1890ff">
                  ✅
                </el-icon>
              </div>
            </div>
          </div>
        </div>

        <!-- 步骤2: 配置认证 -->
        <div v-show="currentStep === 1" class="step-content">
          <div class="step-title">
            <h3>配置认证信息</h3>
            <p>请输入{{ getCurrentProviderName() }}的访问密钥</p>
          </div>

          <div class="auth-form">
            <el-form-item label="账号名称" prop="name">
              <el-input
                v-model="formState.name"
                placeholder="请输入云账号名称，用于标识此账号"
                size="large"
                clearable
              >
                <template #prefix>
                  <el-icon>{{ iconMap.User }}</el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="AccessKey ID" prop="access_key">
              <el-input
                v-model="formState.access_key"
                placeholder="请输入 AccessKey ID"
                size="large"
                clearable
              >
                <template #prefix>
                  <el-icon>{{ iconMap.Key }}</el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="AccessKey Secret" prop="secret_key">
              <el-input
                v-model="formState.secret_key"
                placeholder="请输入 AccessKey Secret"
                type="password"
                size="large"
                show-password
                clearable
              >
                <template #prefix>
                  <el-icon>{{ iconMap.Lock }}</el-icon>
                </template>
              </el-input>
            </el-form-item>

            <!-- 验证AccessKey按钮 -->
            <el-form-item>
              <el-button
                type="primary"
                :loading="validating"
                :disabled="!canValidateCredentials"
                @click="validateAccessKey"
                style="width: 100%; margin-top: 10px;"
                size="large"
              >
                <el-icon v-if="!validating">{{ iconMap.Check }}</el-icon>
                <span>{{ validating ? '验证中...' : '验证AccessKey' }}</span>
              </el-button>

              <!-- 验证结果显示 -->
              <div v-if="accessKeyValidationResult !== null" class="validation-result" :class="accessKeyValidationResult ? 'success' : 'error'">
                <el-icon v-if="accessKeyValidationResult">{{ iconMap.CircleCheck }}</el-icon>
                <el-icon v-else>{{ iconMap.CircleClose }}</el-icon>
                <span>{{ accessKeyValidationMessage }}</span>
              </div>
            </el-form-item>

            <!-- 区域选择 -->
            <el-form-item
              v-if="needsRegion"
              label="默认区域"
              prop="region"
            >
              <el-select
                v-model="formState.region"
                placeholder="请选择默认区域"
                size="large"
                style="width: 100%"
                :loading="loadingRegions"
                filterable
              >
                <el-option
                  v-for="region in availableRegions"
                  :key="region.value"
                  :label="region.label"
                  :value="region.value"
                />
              </el-select>
            </el-form-item>

            <!-- 安全提示 -->
            <el-alert
              title="安全提示"
              type="info"
              :closable="false"
              show-icon
            >
              <template #default>
                <p>• 请确保提供的密钥具有足够的权限</p>
                <p>• 建议使用只读权限的子账号密钥</p>
                <p>• 密钥信息将被加密存储</p>
              </template>
            </el-alert>
          </div>
        </div>

        <!-- 步骤3: 完成设置 -->
        <div v-show="currentStep === 2" class="step-content">
          <div class="step-title">
            <h3>验证连接</h3>
            <p>正在验证您的云账号配置...</p>
          </div>

          <div class="validation-content">
            <div v-if="validationResult === null" class="validation-loading">
              <el-icon class="loading-icon" size="48">⏳</el-icon>
              <p>正在验证连接...</p>
            </div>

            <div v-else-if="validationResult === true" class="validation-success">
              <el-icon class="success-icon" size="48" color="#52c41a">✅</el-icon>
              <h4>验证成功！</h4>
              <p>云账号配置正确，可以正常访问云资源</p>
            </div>

            <div v-else class="validation-error">
              <el-icon class="error-icon" size="48" color="#ff4d4f">❌</el-icon>
              <h4>验证失败</h4>
              <p>{{ validationError || '无法连接到云厂商，请检查配置信息' }}</p>
              <el-button type="primary" @click="currentStep = 1">返回修改</el-button>
            </div>
          </div>
        </div>
      </el-form>
    </div>

    <!-- 自定义底部 -->
    <template #footer>
      <div class="modal-footer">
        <div class="footer-left">
          <el-button v-if="currentStep > 0" @click="prevStep">
            <el-icon>←</el-icon>
            上一步
          </el-button>
        </div>
        <div class="footer-right">
          <el-button @click="handleCancel">取消</el-button>
          <el-button 
            v-if="currentStep < 2" 
            type="primary" 
            @click="nextStep"
            :disabled="!canNextStep"
          >
            下一步
            <el-icon>→</el-icon>
          </el-button>
          <el-button 
            v-else
            type="primary" 
            @click="handleSubmit"
            :loading="loading"
            :disabled="validationResult !== true"
          >
            {{ isEdit ? '更新' : '创建' }}
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, PropType } from 'vue'
import { ElMessage } from 'element-plus'
// 使用字符串图标替代 Element Plus 图标
const iconMap = {
  Close: '✕',
  Check: '✅',
  User: '👤',
  Key: '🔑',
  Lock: '🔒',
  Loading: '⏳',
  CircleCheck: '✅',
  CircleClose: '❌',
  ArrowLeft: '←',
  ArrowRight: '→'
}
import {
  createProvider,
  updateProvider,
  validateProviderCredentials as validateProvider,
  getProviderRegions
} from '@/api/system/host'

interface Provider {
  id?: number
  name: string
  type: string
  access_key: string
  secret_key: string
  region?: string
}

// Props
const props = defineProps({
  visible: Boolean,
  record: Object as PropType<Provider | null>
})

// Emits
const emit = defineEmits(['success', 'cancel'])

// 响应式数据
const dialogVisible = ref(false)
const currentStep = ref(0)
const loading = ref(false)
const loadingRegions = ref(false)
const validationResult = ref<boolean | null>(null)
const validationError = ref('')
const formRef = ref()

// AccessKey验证相关数据
const validating = ref(false)
const accessKeyValidationResult = ref<boolean | null>(null)
const accessKeyValidationMessage = ref('')

// 表单数据
const formState = ref({
  name: '',
  type: '',
  access_key: '',
  secret_key: '',
  region: ''
})

// 是否编辑模式
const isEdit = computed(() => !!props.record)

// 云厂商选项
const providerOptions = ref([
  {
    value: 'aliyun',
    label: '阿里云',
    description: '阿里巴巴云计算服务',
    icon: 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzIiIGhlaWdodD0iMzIiIHZpZXdCb3g9IjAgMCAzMiAzMiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3Qgd2lkdGg9IjMyIiBoZWlnaHQ9IjMyIiByeD0iNCIgZmlsbD0iI0ZGNkEwMCIvPgo8cGF0aCBkPSJNOCAxNkMxMiAxMiAyMCAxMiAyNCAxNkMyMCAyMCAxMiAyMCA4IDE2WiIgZmlsbD0id2hpdGUiLz4KPC9zdmc+'
  },
  {
    value: 'tencent',
    label: '腾讯云',
    description: '腾讯云计算服务',
    icon: 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzIiIGhlaWdodD0iMzIiIHZpZXdCb3g9IjAgMCAzMiAzMiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3Qgd2lkdGg9IjMyIiBoZWlnaHQ9IjMyIiByeD0iNCIgZmlsbD0iIzAwNkVGRiIvPgo8cGF0aCBkPSJNMTYgOEMxMiAxMiA4IDE2IDE2IDE2QzIwIDEyIDI0IDggMTYgOFoiIGZpbGw9IndoaXRlIi8+CjxwYXRoIGQ9Ik0xNiAyNEMyMCAyMCAyNCAxNiAxNiAxNkMxMiAyMCA4IDI0IDE2IDI0WiIgZmlsbD0id2hpdGUiLz4KPC9zdmc+'
  },
  {
    value: 'aws',
    label: 'AWS',
    description: '亚马逊云计算服务',
    icon: 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzIiIGhlaWdodD0iMzIiIHZpZXdCb3g9IjAgMCAzMiAzMiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3Qgd2lkdGg9IjMyIiBoZWlnaHQ9IjMyIiByeD0iNCIgZmlsbD0iI0ZGOTkwMCIvPgo8cGF0aCBkPSJNOCAxMkgxNkwyNCAxMkwyMCAyMEgxMkw4IDIwVjEyWiIgZmlsbD0id2hpdGUiLz4KPC9zdmc+'
  }
])

// 可用区域
const availableRegions = ref([])

// 是否需要区域选择
const needsRegion = computed(() => {
  return ['aliyun', 'tencent', 'aws'].includes(formState.value.type)
})

// 能否进入下一步
const canNextStep = computed(() => {
  if (currentStep.value === 0) {
    return !!formState.value.type
  }
  if (currentStep.value === 1) {
    return formState.value.name &&
           formState.value.access_key &&
           formState.value.secret_key &&
           (!needsRegion.value || formState.value.region)
  }
  return false
})

// 是否可以验证AccessKey
const canValidateCredentials = computed(() => {
  return formState.value.type &&
         formState.value.access_key &&
         formState.value.secret_key &&
         (!needsRegion.value || formState.value.region)
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入账号名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择云厂商类型', trigger: 'change' }
  ],
  access_key: [
    { required: true, message: '请输入AccessKey ID', trigger: 'blur' }
  ],
  secret_key: [
    { required: true, message: '请输入AccessKey Secret', trigger: 'blur' }
  ],
  region: [
    { required: true, message: '请选择区域', trigger: 'change' }
  ]
}

// 监听props变化
watch(
  () => props.visible,
  (val) => {
    dialogVisible.value = val
    if (val) {
      resetForm()
      if (props.record) {
        Object.assign(formState.value, props.record)
        currentStep.value = 1 // 编辑时直接跳到配置步骤
      }
    }
  }
)

// 监听dialogVisible变化，发出事件
watch(
  () => dialogVisible.value,
  (val) => {
    if (!val) {
      emit('cancel')
    }
  }
)

// 监听厂商类型变化，自动获取区域列表
watch(
  () => formState.value.type,
  (newType) => {
    if (newType) {
      fetchRegions(newType)
    } else {
      availableRegions.value = []
    }
  }
)

// 方法
function selectProvider(type: string) {
  formState.value.type = type
  if (needsRegion.value) {
    fetchRegions(type)
  }
}

function getCurrentProviderName() {
  const provider = providerOptions.value.find(p => p.value === formState.value.type)
  return provider?.label || '云厂商'
}

async function fetchRegions(type: string) {
  if (!type) return

  loadingRegions.value = true
  try {
    const response = await getProviderRegions(type)
    availableRegions.value = response.data || []
  } catch (error) {
    console.error('获取区域失败:', error)
    ElMessage.error('获取区域列表失败')
  } finally {
    loadingRegions.value = false
  }
}

function nextStep() {
  if (currentStep.value === 1) {
    // 进入验证步骤
    currentStep.value = 2
    validateCredentials()
  } else {
    currentStep.value++
  }
}

function prevStep() {
  if (currentStep.value > 0) {
    currentStep.value--
    validationResult.value = null
  }
}

// 验证AccessKey
async function validateAccessKey() {
  if (!canValidateCredentials.value) {
    ElMessage.warning('请先填写完整的AccessKey信息')
    return
  }

  validating.value = true
  accessKeyValidationResult.value = null
  accessKeyValidationMessage.value = ''

  try {
    const response = await validateProvider({
      type: formState.value.type,
      access_key: formState.value.access_key,
      secret_key: formState.value.secret_key,
      region: formState.value.region
    })

    if (response.code === 200) {
      accessKeyValidationResult.value = response.data.valid
      accessKeyValidationMessage.value = response.data.message || (response.data.valid ? 'AccessKey验证成功' : 'AccessKey验证失败')

      if (response.data.valid) {
        ElMessage.success('AccessKey验证成功！')
      } else {
        ElMessage.error(accessKeyValidationMessage.value)
      }
    } else {
      accessKeyValidationResult.value = false
      accessKeyValidationMessage.value = response.message || '验证请求失败'
      ElMessage.error(accessKeyValidationMessage.value)
    }
  } catch (error: any) {
    console.error('验证AccessKey失败:', error)
    accessKeyValidationResult.value = false
    accessKeyValidationMessage.value = error.message || '验证AccessKey时发生错误'
    ElMessage.error(accessKeyValidationMessage.value)
  } finally {
    validating.value = false
  }
}

async function validateCredentials() {
  validationResult.value = null
  validationError.value = ''

  try {
    await validateProvider({
      type: formState.value.type,
      access_key: formState.value.access_key,
      secret_key: formState.value.secret_key,
      region: formState.value.region
    })
    validationResult.value = true
  } catch (error: any) {
    validationResult.value = false
    validationError.value = error.message || '验证失败'
  }
}

async function handleSubmit() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true

    if (isEdit.value) {
      await updateProvider(props.record!.id!, formState.value)
      ElMessage.success('更新成功')
    } else {
      await createProvider(formState.value)
      ElMessage.success('创建成功')
    }

    emit('success')
    dialogVisible.value = false
  } catch (error: any) {
    console.error(error)
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}

function handleCancel() {
  dialogVisible.value = false
  emit('cancel')
}

function resetForm() {
  currentStep.value = 0
  validationResult.value = null
  validationError.value = ''
  formState.value = {
    name: '',
    type: '',
    access_key: '',
    secret_key: '',
    region: ''
  }
  availableRegions.value = []

  // 重置AccessKey验证状态
  validating.value = false
  accessKeyValidationResult.value = null
  accessKeyValidationMessage.value = ''
}

// 暴露方法给父组件
function open(record?: Provider) {
  dialogVisible.value = true
  if (record) {
    Object.assign(formState.value, record)
    currentStep.value = 1
  }
}

defineExpose({
  open
})
</script>

<style scoped>
/* 模态框样式 */
.modern-provider-modal {
  --primary-color: #1890ff;
  --success-color: #52c41a;
  --error-color: #ff4d4f;
  --border-radius: 8px;
  --box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

:deep(.el-dialog) {
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
}

:deep(.el-dialog__header) {
  padding: 0;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.el-dialog__body) {
  padding: 0;
}

:deep(.el-dialog__footer) {
  padding: 0;
  border-top: 1px solid #f0f0f0;
}

/* 头部样式 */
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 32px;
  background: linear-gradient(135deg, #f6f9fc 0%, #ffffff 100%);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--primary-color), #40a9ff);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
}

.header-content {
  flex: 1;
}

.header-title {
  margin: 0 0 4px 0;
  font-size: 20px;
  font-weight: 600;
  color: #262626;
}

.header-subtitle {
  margin: 0;
  font-size: 14px;
  color: #8c8c8c;
}

.close-btn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.close-btn:hover {
  background: #f5f5f5;
}

/* 主体样式 */
.modal-body {
  padding: 32px;
  min-height: 400px;
}

.modern-form {
  height: 100%;
}

/* 步骤指示器 */
.steps-container {
  margin-bottom: 32px;
}

:deep(.el-steps) {
  margin-bottom: 0;
}

:deep(.el-step__title) {
  font-size: 14px;
  font-weight: 500;
}

/* 步骤内容 */
.step-content {
  min-height: 320px;
}

.step-title {
  text-align: center;
  margin-bottom: 32px;
}

.step-title h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
  font-weight: 600;
  color: #262626;
}

.step-title p {
  margin: 0;
  font-size: 14px;
  color: #8c8c8c;
}

/* 云厂商选择网格 */
.provider-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
  max-width: 480px;
  margin: 0 auto;
}

.provider-card {
  display: flex;
  align-items: center;
  padding: 20px;
  border: 2px solid #f0f0f0;
  border-radius: var(--border-radius);
  cursor: pointer;
  transition: all 0.3s ease;
  background: #ffffff;
}

.provider-card:hover {
  border-color: var(--primary-color);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.15);
  transform: translateY(-2px);
}

.provider-card.active {
  border-color: var(--primary-color);
  background: linear-gradient(135deg, #e6f7ff 0%, #f6ffed 100%);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.2);
}

.provider-icon {
  width: 48px;
  height: 48px;
  margin-right: 16px;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8f9fa;
}

.provider-icon img {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

.provider-info {
  flex: 1;
}

.provider-info h4 {
  margin: 0 0 4px 0;
  font-size: 16px;
  font-weight: 600;
  color: #262626;
}

.provider-info p {
  margin: 0;
  font-size: 14px;
  color: #8c8c8c;
}

.provider-check {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 认证表单 */
.auth-form {
  max-width: 480px;
  margin: 0 auto;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #262626;
}

:deep(.el-input__wrapper) {
  border-radius: 6px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.02);
  transition: all 0.2s;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

:deep(.el-input.is-focus .el-input__wrapper) {
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

/* 验证内容 */
.validation-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 200px;
  text-align: center;
}

.validation-loading,
.validation-success,
.validation-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.loading-icon {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.success-icon,
.error-icon {
  margin-bottom: 8px;
}

.validation-success h4,
.validation-error h4 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.validation-success h4 {
  color: var(--success-color);
}

.validation-error h4 {
  color: var(--error-color);
}

.validation-success p,
.validation-error p {
  margin: 0 0 16px 0;
  font-size: 14px;
  color: #8c8c8c;
}

/* 底部样式 */
.modal-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 32px;
  background: #fafafa;
}

.footer-left,
.footer-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* 响应式 */
@media (max-width: 768px) {
  .modal-header {
    padding: 20px;
  }

  .modal-body {
    padding: 20px;
  }

  .modal-footer {
    padding: 16px 20px;
  }

  .provider-grid {
    grid-template-columns: 1fr;
  }
}

/* Element Plus 组件样式覆盖 */
:deep(.el-select-dropdown) {
  z-index: 99999 !important;
  background-color: #ffffff !important;
  border: 1px solid #dcdfe6 !important;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1) !important;
}

:deep(.el-select-dropdown .el-select-dropdown__item) {
  color: #606266 !important;
  background-color: #ffffff !important;
}

:deep(.el-select-dropdown .el-select-dropdown__item:hover) {
  background-color: #f5f7fa !important;
  color: #409eff !important;
}

:deep(.el-select-dropdown .el-select-dropdown__item.selected) {
  background-color: #409eff !important;
  color: #ffffff !important;
}

:deep(.el-popper) {
  z-index: 99999 !important;
}

:deep(.el-select__popper) {
  z-index: 99999 !important;
}

:deep(.el-tooltip__popper) {
  z-index: 99999 !important;
}

/* 特别针对区域选择下拉框的样式修复 */
:deep(.el-form-item:has([placeholder*="区域"]) .el-select-dropdown),
:deep(.el-form-item:has([placeholder*="region"]) .el-select-dropdown) {
  background-color: #ffffff !important;
  border: 1px solid #dcdfe6 !important;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12) !important;
  z-index: 99999 !important;
}

:deep(.el-form-item:has([placeholder*="区域"]) .el-option),
:deep(.el-form-item:has([placeholder*="region"]) .el-option) {
  color: #606266 !important;
  background-color: #ffffff !important;
  padding: 8px 12px !important;
}

:deep(.el-form-item:has([placeholder*="区域"]) .el-option:hover),
:deep(.el-form-item:has([placeholder*="region"]) .el-option:hover) {
  background-color: #f5f7fa !important;
  color: #409eff !important;
}

:deep(.el-form-item:has([placeholder*="区域"]) .el-option.selected),
:deep(.el-form-item:has([placeholder*="region"]) .el-option.selected) {
  background-color: #409eff !important;
  color: #ffffff !important;
}

/* AccessKey验证结果样式 */
.validation-result {
  margin-top: 10px;
  padding: 8px 12px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  font-size: 13px;
}

.validation-result.success {
  background: #f0f9f0;
  border: 1px solid #b3e19d;
  color: #67c23a;
}

.validation-result.error {
  background: #fef0f0;
  border: 1px solid #fbc4c4;
  color: #f56c6c;
}

.validation-result .el-icon {
  margin-right: 6px;
  font-size: 14px;
}
</style>
