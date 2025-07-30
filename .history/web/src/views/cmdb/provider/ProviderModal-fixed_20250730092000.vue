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
                            type="link" 
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

        <!-- 第一步：选择云厂商 -->
        <div v-if="currentStep === 0" class="step-content">
          <h3 class="step-title">选择云厂商类型</h3>
          <div class="provider-grid">
            <div 
              v-for="provider in providerOptions" 
              :key="provider.value"
              class="provider-card"
              :class="{ active: formState.type === provider.value }"
              @click="selectProvider(provider.value)"
            >
              <div class="provider-icon">{{ provider.icon }}</div>
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

        <!-- 第二步：配置认证 -->
        <div v-if="currentStep === 1" class="step-content">
          <h3 class="step-title">配置{{ getCurrentProviderName() }}认证信息</h3>
          
          <el-form-item label="云账号名称" prop="name">
            <el-input
              v-model="formState.name"
              placeholder="请输入云账号名称，用于标识此账号"
              size="large"
              clearable
            >
              <template #prefix>
                <el-icon>👤</el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="AccessKey ID" prop="access_key_id">
            <el-input
              v-model="formState.access_key_id"
              placeholder="请输入 AccessKey ID"
              size="large"
              clearable
            >
              <template #prefix>
                <el-icon>🔑</el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="AccessKey Secret" prop="access_key_secret">
            <el-input
              v-model="formState.access_key_secret"
              type="password"
              placeholder="请输入 AccessKey Secret"
              size="large"
              show-password
              clearable
            >
              <template #prefix>
                <el-icon>🔒</el-icon>
              </template>
            </el-input>
          </el-form-item>

          <!-- 区域选择 -->
          <el-form-item label="默认区域" prop="region">
            <el-select
              v-model="formState.region"
              placeholder="请选择默认区域"
              size="large"
              clearable
              filterable
              :loading="loadingRegions"
              style="width: 100%"
              @focus="handleRegionFocus"
            >
              <el-option
                v-for="region in availableRegions"
                :key="region.value"
                :label="region.label"
                :value="region.value"
              />
            </el-select>
          </el-form-item>
        </div>

        <!-- 第三步：完成设置 -->
        <div v-if="currentStep === 2" class="step-content">
          <h3 class="step-title">确认配置信息</h3>
          <div class="config-summary">
            <div class="summary-item">
              <span class="label">云厂商：</span>
              <span class="value">{{ getCurrentProviderName() }}</span>
            </div>
            <div class="summary-item">
              <span class="label">账号名称：</span>
              <span class="value">{{ formState.name }}</span>
            </div>
            <div class="summary-item">
              <span class="label">AccessKey ID：</span>
              <span class="value">{{ formState.access_key_id }}</span>
            </div>
            <div class="summary-item">
              <span class="label">默认区域：</span>
              <span class="value">{{ formState.region || '未设置' }}</span>
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
          >
            {{ isEdit ? '更新' : '创建' }}
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  createProvider, 
  updateProvider, 
  getProviderRegions 
} from '@/api/system/host'

// Props
const props = defineProps({
  visible: Boolean,
  record: Object
})

// Emits
const emit = defineEmits(['success', 'cancel'])

// 响应式数据
const dialogVisible = ref(false)
const currentStep = ref(0)
const loading = ref(false)
const loadingRegions = ref(false)
const formRef = ref()
const availableRegions = ref([])

// 表单数据
const formState = reactive({
  name: '',
  type: '',
  access_key_id: '',
  access_key_secret: '',
  region: ''
})

// 是否编辑模式
const isEdit = computed(() => !!props.record)

// 云厂商选项
const providerOptions = [
  { 
    value: 'aliyun', 
    label: '阿里云', 
    icon: '☁️',
    description: '阿里巴巴云计算服务'
  },
  { 
    value: 'tencent', 
    label: '腾讯云', 
    icon: '☁️',
    description: '腾讯云计算服务'
  },
  { 
    value: 'aws', 
    label: 'AWS', 
    icon: '☁️',
    description: '亚马逊云计算服务'
  }
]

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入云账号名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择云厂商类型', trigger: 'change' }
  ],
  access_key_id: [
    { required: true, message: '请输入AccessKey ID', trigger: 'blur' }
  ],
  access_key_secret: [
    { required: true, message: '请输入AccessKey Secret', trigger: 'blur' }
  ]
}

// 计算属性
const canNextStep = computed(() => {
  if (currentStep.value === 0) {
    return !!formState.type
  }
  if (currentStep.value === 1) {
    return !!(formState.name && formState.access_key_id && formState.access_key_secret)
  }
  return true
})

// 监听visible变化
watch(() => props.visible, (val) => {
  dialogVisible.value = val
  if (val) {
    resetForm()
    if (props.record) {
      Object.assign(formState, props.record)
      currentStep.value = 1
    }
  }
})

// 监听dialogVisible变化
watch(dialogVisible, (val) => {
  if (!val) {
    emit('cancel')
  }
})

// 监听厂商类型变化
watch(() => formState.type, (newType) => {
  if (newType) {
    fetchRegions(newType)
  } else {
    availableRegions.value = []
  }
})

// 方法
function resetForm() {
  Object.assign(formState, {
    name: '',
    type: '',
    access_key_id: '',
    access_key_secret: '',
    region: ''
  })
  currentStep.value = 0
  availableRegions.value = []
}

function selectProvider(type: string) {
  formState.type = type
}

function getCurrentProviderName() {
  const provider = providerOptions.find(p => p.value === formState.type)
  return provider?.label || '云厂商'
}

function nextStep() {
  if (canNextStep.value && currentStep.value < 2) {
    currentStep.value++
  }
}

function prevStep() {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

function handleCancel() {
  dialogVisible.value = false
}

async function handleSubmit() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true

    if (isEdit.value) {
      await updateProvider(props.record.id, formState)
    } else {
      await createProvider(formState)
    }

    // 先关闭模态框
    dialogVisible.value = false
    emit('success')
    
    // 延迟显示成功消息，确保在模态框关闭后显示在页面顶部
    setTimeout(() => {
      ElMessage({
        message: `云账号"${formState.name}"${isEdit.value ? '更新' : '添加'}成功`,
        type: 'success',
        duration: 3000,
        showClose: true,
        customClass: 'success-message-top'
      })
    }, 100)
    
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}

async function fetchRegions(type: string) {
  if (!type) return

  loadingRegions.value = true
  try {
    const response = await getProviderRegions(type)
    console.log('获取区域响应:', response) // 调试日志
    // 确保数据格式正确
    if (Array.isArray(response)) {
      availableRegions.value = response
    } else if (response && Array.isArray(response.data)) {
      availableRegions.value = response.data
    } else {
      availableRegions.value = []
    }
  } catch (error) {
    console.error('获取区域失败:', error)
    ElMessage.error('获取区域列表失败')
    availableRegions.value = []
  } finally {
    loadingRegions.value = false
  }
}

function handleRegionFocus() {
  if (formState.type && availableRegions.value.length === 0) {
    fetchRegions(formState.type)
  }
}
</script>

<style scoped>
.modern-provider-modal {
  border-radius: 12px;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-icon {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.header-content h2 {
  margin: 0 0 4px 0;
  font-size: 20px;
  font-weight: 600;
}

.header-content p {
  margin: 0;
  font-size: 14px;
  opacity: 0.9;
}

.close-btn {
  color: white !important;
  border: none !important;
  background: rgba(255, 255, 255, 0.1) !important;
  border-radius: 8px !important;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.2) !important;
}

.modal-body {
  padding: 32px;
  min-height: 400px;
}

.steps-container {
  margin-bottom: 32px;
}

.step-content {
  min-height: 300px;
}

.step-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 24px;
  text-align: center;
}

.provider-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
  max-width: 500px;
  margin: 0 auto;
}

.provider-card {
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  gap: 16px;
  position: relative;
}

.provider-card:hover {
  border-color: #1890ff;
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.15);
}

.provider-card.active {
  border-color: #1890ff;
  background-color: #f0f8ff;
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.15);
}

.provider-icon {
  font-size: 32px;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8fafc;
  border-radius: 8px;
}

.provider-info h4 {
  margin: 0 0 4px 0;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.provider-info p {
  margin: 0;
  font-size: 14px;
  color: #666;
}

.provider-check {
  position: absolute;
  top: 12px;
  right: 12px;
  font-size: 18px;
}

.config-summary {
  background: #f8fafc;
  border-radius: 8px;
  padding: 24px;
  max-width: 500px;
  margin: 0 auto;
}

.summary-item {
  display: flex;
  margin-bottom: 16px;
  align-items: center;
}

.summary-item:last-child {
  margin-bottom: 0;
}

.summary-item .label {
  width: 120px;
  color: #666;
  font-weight: 500;
}

.summary-item .value {
  color: #333;
  font-weight: 600;
}

.modal-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 32px;
  background: #f8fafc;
  border-top: 1px solid #e5e7eb;
}

.footer-left, .footer-right {
  display: flex;
  gap: 12px;
}

/* Element Plus 组件样式覆盖 */
:deep(.el-select-dropdown) {
  z-index: 99999 !important;
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

/* 成功消息样式 - 确保显示在最顶层 */
:global(.success-message-top) {
  z-index: 999999 !important;
  position: fixed !important;
  top: 20px !important;
  left: 50% !important;
  transform: translateX(-50%) !important;
}

/* 表单样式优化 */
:deep(.el-form-item__label) {
  font-weight: 600;
  color: #333;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
}

:deep(.el-select .el-input__wrapper) {
  border-radius: 8px;
}

/* 步骤条样式优化 */
:deep(.el-steps) {
  margin-bottom: 32px;
}

:deep(.el-step__title) {
  font-weight: 600;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .modal-header {
    padding: 20px 24px;
  }

  .modal-body {
    padding: 24px;
  }

  .modal-footer {
    padding: 16px 24px;
  }

  .provider-grid {
    grid-template-columns: 1fr;
  }
}
</style>
