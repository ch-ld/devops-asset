<template>
  <el-dialog
    v-model="dialogVisible"
    title="添加云账号"
    width="600px"
    :close-on-click-modal="false"
    destroy-on-close
  >
    <div class="provider-modal-content">
      <!-- 步骤指示器 -->
      <el-steps :active="currentStep" align-center class="steps-container">
        <el-step title="选择云厂商" />
        <el-step title="配置认证" />
        <el-step title="完成设置" />
      </el-steps>

      <!-- 第一步：选择云厂商 -->
      <div v-if="currentStep === 0" class="step-content">
        <h3>选择云厂商类型</h3>
        <div class="provider-options">
          <div 
            v-for="option in providerOptions" 
            :key="option.value"
            class="provider-option"
            :class="{ active: formData.type === option.value }"
            @click="selectProvider(option.value)"
          >
            <div class="option-icon">{{ option.icon }}</div>
            <div class="option-name">{{ option.label }}</div>
            <div v-if="formData.type === option.value" class="option-check">✅</div>
          </div>
        </div>
      </div>

      <!-- 第二步：配置认证 -->
      <div v-if="currentStep === 1" class="step-content">
        <h3>配置{{ getCurrentProviderName() }}认证信息</h3>
        <el-form :model="formData" label-width="120px">
          <el-form-item label="账号名称" required>
            <el-input v-model="formData.name" placeholder="请输入账号名称" />
          </el-form-item>
          <el-form-item label="AccessKey ID" required>
            <el-input v-model="formData.access_key_id" placeholder="请输入AccessKey ID" />
          </el-form-item>
          <el-form-item label="AccessKey Secret" required>
            <el-input 
              v-model="formData.access_key_secret" 
              type="password" 
              placeholder="请输入AccessKey Secret" 
              show-password
            />
          </el-form-item>
          <el-form-item label="默认区域">
            <el-select 
              v-model="formData.region" 
              placeholder="请选择默认区域"
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
        </el-form>
      </div>

      <!-- 第三步：完成设置 -->
      <div v-if="currentStep === 2" class="step-content">
        <h3>确认配置信息</h3>
        <div class="config-summary">
          <div class="summary-item">
            <span class="label">云厂商：</span>
            <span class="value">{{ getCurrentProviderName() }}</span>
          </div>
          <div class="summary-item">
            <span class="label">账号名称：</span>
            <span class="value">{{ formData.name }}</span>
          </div>
          <div class="summary-item">
            <span class="label">AccessKey ID：</span>
            <span class="value">{{ formData.access_key_id }}</span>
          </div>
          <div class="summary-item">
            <span class="label">默认区域：</span>
            <span class="value">{{ formData.region || '未设置' }}</span>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button v-if="currentStep > 0" @click="prevStep">上一步</el-button>
        <el-button 
          v-if="currentStep < 2" 
          type="primary" 
          @click="nextStep"
          :disabled="!canProceed"
        >
          下一步
        </el-button>
        <el-button 
          v-if="currentStep === 2" 
          type="primary" 
          @click="handleSubmit"
          :loading="submitting"
        >
          完成
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getProviderRegions } from '@/api/system/host'

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
const submitting = ref(false)
const availableRegions = ref([])

// 表单数据
const formData = reactive({
  name: '',
  type: '',
  access_key_id: '',
  access_key_secret: '',
  region: ''
})

// 云厂商选项
const providerOptions = [
  { value: 'aliyun', label: '阿里云', icon: '☁️' },
  { value: 'tencent', label: '腾讯云', icon: '☁️' },
  { value: 'aws', label: 'AWS', icon: '☁️' }
]

// 计算属性
const canProceed = computed(() => {
  if (currentStep.value === 0) {
    return !!formData.type
  }
  if (currentStep.value === 1) {
    return !!(formData.name && formData.access_key_id && formData.access_key_secret)
  }
  return true
})

// 监听visible变化
watch(() => props.visible, (val) => {
  dialogVisible.value = val
  if (val) {
    resetForm()
    if (props.record) {
      Object.assign(formData, props.record)
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
watch(() => formData.type, (newType) => {
  if (newType) {
    fetchRegions(newType)
  } else {
    availableRegions.value = []
  }
})

// 方法
function resetForm() {
  Object.assign(formData, {
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
  formData.type = type
}

function getCurrentProviderName() {
  const provider = providerOptions.find(p => p.value === formData.type)
  return provider?.label || '云厂商'
}

function nextStep() {
  if (canProceed.value && currentStep.value < 2) {
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
  submitting.value = true
  try {
    // 这里应该调用API提交数据
    await new Promise(resolve => setTimeout(resolve, 1000)) // 模拟API调用
    ElMessage.success('云账号添加成功')
    emit('success')
    dialogVisible.value = false
  } catch (error) {
    ElMessage.error('添加失败：' + error.message)
  } finally {
    submitting.value = false
  }
}

async function fetchRegions(type: string) {
  try {
    const response = await getProviderRegions(type)
    availableRegions.value = response.data || []
  } catch (error) {
    console.error('获取区域失败:', error)
    availableRegions.value = []
  }
}

function handleRegionFocus() {
  if (formData.type && availableRegions.value.length === 0) {
    fetchRegions(formData.type)
  }
}
</script>

<style scoped>
.provider-modal-content {
  padding: 20px 0;
}

.steps-container {
  margin-bottom: 30px;
}

.step-content {
  min-height: 300px;
}

.step-content h3 {
  margin-bottom: 20px;
  color: #333;
}

.provider-options {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 16px;
  margin-top: 20px;
}

.provider-option {
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  padding: 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
}

.provider-option:hover {
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.2);
}

.provider-option.active {
  border-color: #1890ff;
  background-color: #f0f8ff;
}

.option-icon {
  font-size: 32px;
  margin-bottom: 8px;
}

.option-name {
  font-weight: 500;
  color: #333;
}

.option-check {
  position: absolute;
  top: 8px;
  right: 8px;
  font-size: 16px;
}

.config-summary {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 20px;
}

.summary-item {
  display: flex;
  margin-bottom: 12px;
}

.summary-item:last-child {
  margin-bottom: 0;
}

.label {
  width: 120px;
  color: #666;
}

.value {
  color: #333;
  font-weight: 500;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
