<template>
  <el-dialog
    :model-value="visible"
    title="添加SSL证书"
    width="600px"
    :before-close="handleClose"
    class="add-certificate-modal"
  >
    <el-form 
      ref="formRef"
      :model="form" 
      :rules="rules"
      label-width="120px"
      class="certificate-form"
    >
      <el-form-item label="证书类型" prop="cert_type">
        <el-radio-group v-model="form.cert_type">
          <el-radio label="manual">手动上传</el-radio>
          <el-radio label="letsencrypt">Let's Encrypt自动申请</el-radio>
          <el-radio label="acme">ACME自动申请</el-radio>
        </el-radio-group>
      </el-form-item>

      <!-- 手动上传模式 -->
      <template v-if="form.cert_type === 'manual'">
        <el-form-item label="证书文件" prop="certificate">
          <el-input
            v-model="form.certificate"
            type="textarea"
            :rows="6"
            placeholder="请粘贴证书内容（PEM格式）"
          />
        </el-form-item>
        
        <el-form-item label="私钥文件" prop="private_key">
          <el-input
            v-model="form.private_key"
            type="textarea"
            :rows="6"
            placeholder="请粘贴私钥内容（PEM格式）"
          />
        </el-form-item>
        
        <el-form-item label="证书链" prop="certificate_chain">
          <el-input
            v-model="form.certificate_chain"
            type="textarea"
            :rows="4"
            placeholder="请粘贴证书链内容（可选）"
          />
        </el-form-item>
      </template>

      <!-- 自动申请模式 -->
      <template v-else>
        <el-form-item label="域名列表" prop="domains">
          <el-input
            v-model="form.domains"
            placeholder="请输入域名，多个域名用逗号分隔"
          />
          <div class="form-help">
            示例：example.com,www.example.com,*.example.com
          </div>
        </el-form-item>
        
        <el-form-item label="验证方式" prop="challenge_type">
          <el-select v-model="form.challenge_type" style="width: 100%">
            <el-option label="HTTP验证" value="http-01" />
            <el-option label="DNS验证" value="dns-01" />
            <el-option label="TLS-ALPN验证" value="tls-alpn-01" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="ACME服务器" prop="acme_server" v-if="form.cert_type === 'acme'">
          <el-select v-model="form.acme_server" style="width: 100%">
            <el-option label="Let's Encrypt" value="https://acme-v02.api.letsencrypt.org/directory" />
            <el-option label="Let's Encrypt (测试)" value="https://acme-staging-v02.api.letsencrypt.org/directory" />
            <el-option label="ZeroSSL" value="https://acme.zerossl.com/v2/DV90" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="自定义服务器" prop="custom_acme_server" v-if="form.acme_server === 'custom'">
          <el-input v-model="form.custom_acme_server" placeholder="请输入ACME服务器地址" />
        </el-form-item>
      </template>

      <el-form-item label="自动续费" prop="auto_renew">
        <el-switch v-model="form.auto_renew" />
        <span class="form-help">启用后将在证书过期前自动续费</span>
      </el-form-item>
      
      <el-form-item label="邮箱地址" prop="email">
        <el-input v-model="form.email" placeholder="用于接收证书相关通知" />
      </el-form-item>
    </el-form>

    <!-- 申请进度显示 -->
    <div v-if="applyProgress" class="apply-progress">
      <div class="progress-header">
        <h4>证书申请进度</h4>
        <el-button
          v-if="progressError"
          type="danger"
          size="small"
          @click="handleCancelApply"
        >
          取消申请
        </el-button>
      </div>

      <el-steps :active="currentStep" finish-status="success" process-status="process">
        <el-step
          v-for="(step, index) in progressSteps"
          :key="index"
          :title="step.title"
          :description="step.description"
          :status="step.status"
        />
      </el-steps>

      <div v-if="progressError" class="progress-error">
        <el-alert
          :title="progressError"
          type="error"
          :closable="false"
          show-icon
        />
      </div>

      <div class="progress-info">
        <p v-if="!progressError">
          正在执行: {{ progressSteps[currentStep]?.title }}
        </p>
        <p class="progress-tip">
          请耐心等待，证书申请过程可能需要几分钟时间
        </p>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button
          @click="handleClose"
          class="cancel-btn"
          :disabled="applyProgress && !progressError"
        >
          取消
        </el-button>
        <el-button
          v-if="!applyProgress"
          type="primary"
          @click="handleSubmit"
          :loading="submitLoading"
          class="submit-btn"
        >
          {{ form.cert_type === 'manual' ? '上传证书' : '申请证书' }}
        </el-button>
        <el-button
          v-else-if="progressError"
          type="primary"
          @click="handleAutoApply"
          class="submit-btn"
        >
          重新申请
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, nextTick } from 'vue'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import { certificateApi } from '@/api/dns/certificate'
import type { Domain } from '@/types/dns'

interface Props {
  visible: boolean
  domain?: Domain | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 响应式数据
const formRef = ref<FormInstance>()
const submitLoading = ref(false)

const form = reactive({
  cert_type: 'letsencrypt',
  certificate: '',
  private_key: '',
  certificate_chain: '',
  domains: '',
  challenge_type: 'http-01',
  acme_server: 'https://acme-v02.api.letsencrypt.org/directory',
  custom_acme_server: '',
  auto_renew: true,
  email: ''
})

// 表单验证规则
const rules: FormRules = {
  cert_type: [
    { required: true, message: '请选择证书类型', trigger: 'change' }
  ],
  certificate: [
    { 
      required: true, 
      validator: (rule, value, callback) => {
        if (form.cert_type === 'manual' && !value) {
          callback(new Error('请输入证书内容'))
        } else {
          callback()
        }
      },
      trigger: 'blur' 
    }
  ],
  private_key: [
    { 
      required: true, 
      validator: (rule, value, callback) => {
        if (form.cert_type === 'manual' && !value) {
          callback(new Error('请输入私钥内容'))
        } else {
          callback()
        }
      },
      trigger: 'blur' 
    }
  ],
  domains: [
    { 
      required: true, 
      validator: (rule, value, callback) => {
        if (form.cert_type !== 'manual' && !value) {
          callback(new Error('请输入域名列表'))
        } else {
          callback()
        }
      },
      trigger: 'blur' 
    }
  ],
  challenge_type: [
    { required: true, message: '请选择验证方式', trigger: 'change' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

// 方法
const resetForm = () => {
  Object.assign(form, {
    cert_type: 'letsencrypt',
    certificate: '',
    private_key: '',
    certificate_chain: '',
    domains: props.domain?.name || '',
    challenge_type: 'http-01',
    acme_server: 'https://acme-v02.api.letsencrypt.org/directory',
    custom_acme_server: '',
    auto_renew: true,
    email: ''
  })
  formRef.value?.clearValidate()
}

const handleClose = () => {
  emit('update:visible', false)
}

// 申请进度相关
const applyProgress = ref(false)
const currentStep = ref(0)
const progressSteps = ref([
  { title: '验证域名', status: 'wait', description: '验证域名所有权' },
  { title: '生成密钥', status: 'wait', description: '生成证书私钥和CSR' },
  { title: 'DNS验证', status: 'wait', description: '添加DNS TXT记录' },
  { title: '申请证书', status: 'wait', description: '向CA申请证书' },
  { title: '下载证书', status: 'wait', description: '下载并保存证书' }
])
const progressError = ref('')

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    if (form.cert_type === 'manual') {
      // 手动上传证书
      submitLoading.value = true
      const submitData = { ...form }

      await certificateApi.create({
        domain_id: props.domain?.id,
        ...submitData
      })

      ElMessage.success('证书上传成功')
      emit('success')
    } else {
      // 自动申请证书，显示进度
      await handleAutoApply()
    }
  } catch (error) {
    if (error !== 'validation-failed') {
      ElMessage.error(form.cert_type === 'manual' ? '证书上传失败' : '证书申请失败')
    }
  } finally {
    submitLoading.value = false
  }
}

const handleAutoApply = async () => {
  applyProgress.value = true
  currentStep.value = 0
  progressError.value = ''

  // 重置所有步骤状态
  progressSteps.value.forEach(step => {
    step.status = 'wait'
  })

  try {
    const submitData = { ...form }

    // 处理域名列表
    submitData.domains = form.domains.split(',').map(d => d.trim()).filter(Boolean).join(',')

    // 处理ACME服务器
    if (form.cert_type === 'acme' && form.acme_server === 'custom') {
      submitData.acme_server = form.custom_acme_server
    }

    // 开始申请流程
    await processApplySteps(submitData)

    ElMessage.success('证书申请成功')
    emit('success')
  } catch (error: any) {
    progressError.value = error.message || '申请失败'
    progressSteps.value[currentStep.value].status = 'error'
    ElMessage.error('证书申请失败: ' + progressError.value)
  }
}

const processApplySteps = async (submitData: any) => {
  // 步骤1: 验证域名
  await executeStep(0, async () => {
    await new Promise(resolve => setTimeout(resolve, 1000))
    // 这里应该调用域名验证API
  })

  // 步骤2: 生成密钥
  await executeStep(1, async () => {
    await new Promise(resolve => setTimeout(resolve, 1500))
    // 这里应该调用密钥生成API
  })

  // 步骤3: DNS验证
  await executeStep(2, async () => {
    await new Promise(resolve => setTimeout(resolve, 2000))
    // 这里应该调用DNS验证API
  })

  // 步骤4: 申请证书
  await executeStep(3, async () => {
    await certificateApi.create({
      domain_id: props.domain?.id,
      ...submitData
    })
  })

  // 步骤5: 下载证书
  await executeStep(4, async () => {
    await new Promise(resolve => setTimeout(resolve, 1000))
    // 这里应该调用证书下载API
  })
}

const executeStep = async (stepIndex: number, stepFunction: () => Promise<void>) => {
  currentStep.value = stepIndex
  progressSteps.value[stepIndex].status = 'process'

  try {
    await stepFunction()
    progressSteps.value[stepIndex].status = 'finish'
  } catch (error) {
    progressSteps.value[stepIndex].status = 'error'
    throw error
  }
}

const handleCancelApply = () => {
  applyProgress.value = false
  currentStep.value = 0
  progressError.value = ''
  progressSteps.value.forEach(step => {
    step.status = 'wait'
  })
}

// 监听弹窗显示状态
watch(() => props.visible, (visible) => {
  if (visible) {
    nextTick(() => {
      resetForm()
    })
  }
})
</script>

<style scoped lang="scss">
.add-certificate-modal {
  :deep(.el-dialog) {
    border-radius: 16px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  }
  
  :deep(.el-dialog__header) {
    padding: 24px 24px 0;
    border-bottom: 1px solid #f0f2f5;
  }
  
  :deep(.el-dialog__body) {
    padding: 24px;
    max-height: 70vh;
    overflow-y: auto;
  }
  
  :deep(.el-dialog__footer) {
    padding: 0 24px 24px;
    border-top: 1px solid #f0f2f5;
  }
}

.certificate-form {
  .el-form-item {
    margin-bottom: 24px;
  }
  
  :deep(.el-form-item__label) {
    font-weight: 500;
    color: #374151;
  }
  
  :deep(.el-input__wrapper) {
    border-radius: 8px;
    transition: all 0.3s ease;
  }
  
  :deep(.el-textarea__inner) {
    border-radius: 8px;
    transition: all 0.3s ease;
    font-family: 'Courier New', monospace;
    font-size: 12px;
  }
  
  :deep(.el-select) {
    .el-input__wrapper {
      border-radius: 8px;
    }
  }
  
  .form-help {
    margin-top: 4px;
    font-size: 12px;
    color: #6b7280;
    margin-left: 8px;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 20px;
  
  .cancel-btn {
    padding: 10px 20px;
    border-radius: 8px;
    border: 1px solid #d1d5db;
    background: #f9fafb;
    color: #374151;
    font-weight: 500;
  }
  
  .submit-btn {
    padding: 10px 20px;
    border-radius: 8px;
    background: linear-gradient(135deg, #10b981 0%, #059669 100%);
    border: none;
    font-weight: 600;
  }
}

/* 表单验证错误样式 */
.certificate-form :deep(.el-form-item.is-error) {
  .el-input__wrapper {
    border-color: #ef4444;
    box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.1);
  }
  
  .el-textarea__inner {
    border-color: #ef4444;
    box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.1);
  }
}

/* 表单聚焦样式 */
.certificate-form :deep(.el-input__wrapper:focus-within) {
  border-color: #10b981;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.1);
}

.certificate-form :deep(.el-textarea__inner:focus) {
  border-color: #10b981;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.1);
}

/* 申请进度样式 */
.apply-progress {
  margin: 24px 0;
  padding: 20px;
  background: #f8fafc;
  border-radius: 12px;
  border: 1px solid #e2e8f0;

  .progress-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h4 {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
      color: #1f2937;
    }
  }

  :deep(.el-steps) {
    margin-bottom: 20px;
  }

  .progress-error {
    margin: 16px 0;
  }

  .progress-info {
    margin-top: 16px;

    p {
      margin: 8px 0;
      font-size: 14px;
      color: #6b7280;

      &:first-child {
        font-weight: 500;
        color: #374151;
      }
    }

    .progress-tip {
      font-size: 12px;
      color: #9ca3af;
      font-style: italic;
    }
  }
}

/* 步骤样式优化 */
.apply-progress :deep(.el-step) {
  .el-step__head {
    .el-step__icon {
      border-radius: 50%;
      transition: all 0.3s ease;
    }

    &.is-process .el-step__icon {
      background: #3b82f6;
      border-color: #3b82f6;
      animation: pulse 2s infinite;
    }

    &.is-finish .el-step__icon {
      background: #10b981;
      border-color: #10b981;
    }

    &.is-error .el-step__icon {
      background: #ef4444;
      border-color: #ef4444;
    }
  }

  .el-step__title {
    font-weight: 500;

    &.is-process {
      color: #3b82f6;
    }

    &.is-finish {
      color: #10b981;
    }

    &.is-error {
      color: #ef4444;
    }
  }
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.8;
  }
}
</style>
