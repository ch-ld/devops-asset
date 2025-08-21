<template>
  <div class="cert-create-page">
    <!-- 背景装饰 -->
    <div class="page-background">
      <div class="bg-pattern"></div>
      <div class="bg-gradient"></div>
    </div>

    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <el-button
            text
            @click="goBack"
            class="back-btn"
            size="large"
            :disabled="applying"
          >
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <div class="header-title">
            <h1>
              <el-icon class="title-icon"><Lock /></el-icon>
              申请SSL证书
            </h1>
            <p>为您的域名申请免费的Let's Encrypt SSL证书，保护网站安全</p>
          </div>
        </div>
        <div class="header-actions">
          <el-button @click="goBack" size="large" :disabled="applying" class="cancel-btn">
            取消
          </el-button>
          <el-button
            type="primary"
            :loading="applying"
            @click="handleApply"
            :disabled="!canSubmit || applying"
            size="large"
            class="apply-btn"
          >
            <el-icon><Lightning /></el-icon>
            {{ applying ? '申请中...' : '申请证书' }}
          </el-button>
        </div>
      </div>
    </div>

    <!-- 申请进度 -->
    <div v-if="showProgress" class="progress-section">
      <div class="progress-container">
        <div class="progress-card">
          <div class="progress-header">
            <div class="header-left">
              <div class="progress-icon">
                <el-icon><Lightning /></el-icon>
              </div>
              <div class="header-info">
                <h3>证书申请进度</h3>
                <p>正在为您申请SSL证书，请稍候...</p>
              </div>
            </div>
            <el-tag :type="getProgressType(currentStep)" size="large" class="status-tag">
              {{ getProgressText(currentStep) }}
            </el-tag>
          </div>

          <div class="progress-steps-container">
            <el-steps :active="currentStep" :status="stepStatus" direction="horizontal" class="progress-steps">
              <el-step title="提交申请" description="验证参数并创建申请">
                <template #icon>
                  <el-icon><Document /></el-icon>
                </template>
              </el-step>
              <el-step title="域名验证" description="验证域名所有权">
                <template #icon>
                  <el-icon><Connection /></el-icon>
                </template>
              </el-step>
              <el-step title="证书签发" description="CA机构签发证书">
                <template #icon>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-step>
              <el-step title="完成部署" description="证书安装完成">
                <template #icon>
                  <el-icon><Check /></el-icon>
                </template>
              </el-step>
            </el-steps>
          </div>

          <div class="progress-details">
            <div class="details-grid">
              <div class="detail-item">
                <div class="detail-icon">
                  <el-icon><Globe /></el-icon>
                </div>
                <div class="detail-content">
                  <span class="label">申请域名</span>
                  <span class="value">{{ progressInfo.domain }}</span>
                </div>
              </div>
              <div class="detail-item">
                <div class="detail-icon">
                  <el-icon><Lock /></el-icon>
                </div>
                <div class="detail-content">
                  <span class="label">证书类型</span>
                  <span class="value">{{ progressInfo.certType }}</span>
                </div>
              </div>
              <div class="detail-item">
                <div class="detail-icon">
                  <el-icon><Clock /></el-icon>
                </div>
                <div class="detail-content">
                  <span class="label">开始时间</span>
                  <span class="value">{{ progressInfo.startTime }}</span>
                </div>
              </div>
              <div v-if="progressInfo.estimatedTime" class="detail-item">
                <div class="detail-icon">
                  <el-icon><Timer /></el-icon>
                </div>
                <div class="detail-content">
                  <span class="label">预计完成</span>
                  <span class="value">{{ progressInfo.estimatedTime }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 实时日志 -->
          <div v-if="logs.length > 0" class="logs-section">
            <div class="logs-header">
              <h4>
                <el-icon><Document /></el-icon>
                申请日志
              </h4>
              <el-button text @click="clearLogs" size="small">
                <el-icon><Delete /></el-icon>
                清空
              </el-button>
            </div>
            <div class="logs-container">
              <div
                v-for="(log, index) in logs"
                :key="index"
                :class="['log-item', `log-${log.level}`]"
              >
                <div class="log-indicator"></div>
                <span class="log-time">{{ log.time }}</span>
                <span class="log-message">{{ log.message }}</span>
              </div>
            </div>
          </div>

          <!-- 错误信息 -->
          <div v-if="errorInfo" class="error-section">
            <div class="error-card">
              <div class="error-header">
                <el-icon class="error-icon"><Warning /></el-icon>
                <div class="error-title">{{ errorInfo.title }}</div>
              </div>
              <div class="error-message">{{ errorInfo.message }}</div>
              <div class="error-actions">
                <el-button @click="retryApply" type="primary" size="large">
                  <el-icon><Refresh /></el-icon>
                  重试申请
                </el-button>
                <el-button @click="resetForm" size="large">
                  <el-icon><Edit /></el-icon>
                  重新配置
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 申请表单 -->
    <div v-if="!showProgress" class="form-section">
      <div class="form-container">
        <!-- 步骤指示器 -->
        <div class="steps-indicator">
          <div class="steps-wrapper">
            <div class="step-item active">
              <div class="step-circle">
                <el-icon><Lightning /></el-icon>
              </div>
              <div class="step-info">
                <div class="step-title">选择申请方式</div>
                <div class="step-desc">选择证书申请类型</div>
              </div>
            </div>
            <div class="step-connector"></div>
            <div class="step-item">
              <div class="step-circle">
                <el-icon><Setting /></el-icon>
              </div>
              <div class="step-info">
                <div class="step-title">配置证书信息</div>
                <div class="step-desc">填写证书相关信息</div>
              </div>
            </div>
            <div class="step-connector"></div>
            <div class="step-item">
              <div class="step-circle">
                <el-icon><Check /></el-icon>
              </div>
              <div class="step-info">
                <div class="step-title">完成申请</div>
                <div class="step-desc">确认信息并提交申请</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 申请方式选择 -->
        <div class="form-card apply-method-card">
          <div class="card-header">
            <div class="header-icon">
              <el-icon><Lightning /></el-icon>
            </div>
            <div class="header-content">
              <h3>选择申请方式</h3>
              <p>请选择您希望的证书申请方式</p>
            </div>
          </div>

          <div class="apply-types">
            <div
              class="type-card recommended"
              :class="{ active: formData.applyType === 'auto' }"
              @click="formData.applyType = 'auto'"
            >
              <div class="type-header">
                <div class="type-icon auto">
                  <el-icon><Lightning /></el-icon>
                </div>
                <div class="type-badge">
                  <el-tag type="success" size="small" effect="dark">推荐</el-tag>
                </div>
              </div>
              <div class="type-content">
                <div class="type-title">自动申请</div>
                <div class="type-desc">使用Let's Encrypt免费申请SSL证书，支持自动续期</div>
                <div class="type-features">
                  <div class="feature">
                    <el-icon><Check /></el-icon>
                    <span>完全免费</span>
                  </div>
                  <div class="feature">
                    <el-icon><Check /></el-icon>
                    <span>自动续期</span>
                  </div>
                  <div class="feature">
                    <el-icon><Check /></el-icon>
                    <span>90天有效期</span>
                  </div>
                </div>
              </div>
              <div class="type-radio">
                <el-radio v-model="formData.applyType" value="auto" size="large" />
              </div>
            </div>

            <div 
              class="type-card" 
              :class="{ active: formData.applyType === 'upload' }"
              @click="formData.applyType = 'upload'"
            >
              <div class="type-icon">
                <el-icon><Upload /></el-icon>
              </div>
              <div class="type-content">
                <div class="type-title">导入证书</div>
                <div class="type-desc">导入已有的SSL证书文件</div>
                <div class="type-features">
                  <div class="feature">• 支持现有证书</div>
                  <div class="feature">• 快速导入</div>
                  <div class="feature">• 灵活管理</div>
                </div>
              </div>
              <div class="type-radio">
                <el-radio v-model="formData.applyType" value="upload">upload</el-radio>
              </div>
            </div>

            <div 
              class="type-card" 
              :class="{ active: formData.applyType === 'csr' }"
              @click="formData.applyType = 'csr'"
            >
              <div class="type-icon">
                <el-icon><Document /></el-icon>
              </div>
              <div class="type-content">
                <div class="type-title">自定义CSR</div>
                <div class="type-desc">上传CSR文件申请证书</div>
                <div class="type-features">
                  <div class="feature">• 自定义配置</div>
                  <div class="feature">• 高级选项</div>
                  <div class="feature">• 专业用户</div>
                </div>
              </div>
              <div class="type-radio">
                <el-radio v-model="formData.applyType" value="csr">csr</el-radio>
              </div>
            </div>
          </div>
        </div>

        <!-- 证书配置 -->
        <div class="form-card">
          <div class="card-header">
            <h3><el-icon><Setting /></el-icon>配置证书信息</h3>
            <p>请选择域名和证书类型</p>
          </div>

          <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px">
            <!-- 域名选择 -->
            <el-form-item label="选择域名" prop="domainId" required>
              <el-select 
                v-model="formData.domainId" 
                placeholder="请选择要申请证书的域名"
                filterable
                @change="handleDomainChange"
                style="width: 100%"
                size="large"
              >
                <el-option 
                  v-for="domain in domainOptions" 
                  :key="domain.id"
                  :label="domain.name" 
                  :value="domain.id"
                >
                  <div class="domain-option">
                    <div class="domain-info">
                      <span class="domain-name">{{ domain.name }}</span>
                      <el-tag 
                        v-if="domain.provider" 
                        type="primary" 
                        size="small"
                      >
                        {{ domain.provider.name }}
                      </el-tag>
                    </div>
                    <div class="domain-status">
                      <el-icon color="#67c23a"><Connection /></el-icon>
                      <span>已配置DNS</span>
                    </div>
                  </div>
                </el-option>
              </el-select>
              <div class="form-tip">
                <el-icon><InfoFilled /></el-icon>
                只显示已配置DNS Provider的域名，确保能够自动验证域名所有权
              </div>
            </el-form-item>

            <!-- 证书类型 -->
            <el-form-item label="证书类型" required>
              <el-radio-group v-model="formData.certType" class="cert-types">
                <div class="cert-type-card" :class="{ active: formData.certType === 'single' }">
                  <div class="cert-icon">
                    <el-icon><Lock /></el-icon>
                  </div>
                  <div class="cert-content">
                    <div class="cert-title">单域名证书</div>
                    <div class="cert-desc">仅保护选择的域名</div>
                    <div class="cert-example">例如：example.com</div>
                  </div>
                  <div class="cert-radio">
                    <el-radio value="single">single</el-radio>
                  </div>
                </div>

                <div class="cert-type-card" :class="{ active: formData.certType === 'wildcard' }">
                  <div class="cert-icon">
                    <el-icon><FolderAdd /></el-icon>
                  </div>
                  <div class="cert-content">
                    <div class="cert-title">通配符证书</div>
                    <div class="cert-desc">保护域名及其所有子域名</div>
                    <div class="cert-example">例如：*.example.com</div>
                  </div>
                  <div class="cert-radio">
                    <el-radio value="wildcard">wildcard</el-radio>
                  </div>
                </div>

                <div class="cert-type-card" :class="{ active: formData.certType === 'multi' }">
                  <div class="cert-icon">
                    <el-icon><Plus /></el-icon>
                  </div>
                  <div class="cert-content">
                    <div class="cert-title">多域名证书</div>
                    <div class="cert-desc">保护多个不同的域名</div>
                    <div class="cert-example">例如：多个域名</div>
                  </div>
                  <div class="cert-radio">
                    <el-radio value="multi">multi</el-radio>
                  </div>
                </div>
              </el-radio-group>
            </el-form-item>

            <!-- 邮箱地址 -->
            <el-form-item label="邮箱地址" prop="email" required>
              <el-input
                v-model="formData.email"
                placeholder="请输入邮箱地址"
                size="large"
                prefix-icon="Message"
              />
              <div class="form-tip">
                <el-icon><InfoFilled /></el-icon>
                邮箱用于接收证书到期提醒和重要通知
              </div>
            </el-form-item>

            <!-- 备注 -->
            <el-form-item label="备注">
              <el-input
                v-model="formData.remark"
                type="textarea"
                :rows="3"
                placeholder="请输入备注信息（可选）"
                maxlength="200"
                show-word-limit
              />
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft,
  Lock,
  Document,
  Lightning,
  FolderAdd,
  Upload,
  InfoFilled,
  Connection,
  Plus,
  Setting,
  Message
} from '@element-plus/icons-vue'
import { domainApi } from '@/api/dns/domain'
import { certificateApi } from '@/api/dns/certificate'

const router = useRouter()

// 表单数据
const formData = reactive({
  applyType: 'auto', // auto: 自动申请, upload: 导入证书, csr: 自定义CSR
  domainId: null as number | null,
  certType: 'single', // single: 单域名, wildcard: 通配符, multi: 多域名
  domains: [] as string[],
  email: '',
  remark: ''
})

// 表单验证规则
const formRules = {
  domainId: [
    { required: true, message: '请选择域名', trigger: 'change' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

// 状态管理
const applying = ref(false)
const showProgress = ref(false)
const currentStep = ref(0)
const stepStatus = ref<'wait' | 'process' | 'finish' | 'error'>('process')
const certificateId = ref<number | null>(null)

// 域名选项
const domainOptions = ref<any[]>([])
const formRef = ref()

// 进度信息
const progressInfo = reactive({
  domain: '',
  certType: '',
  startTime: '',
  estimatedTime: ''
})

// 日志信息
const logs = ref<Array<{
  time: string
  level: 'info' | 'success' | 'warning' | 'error'
  message: string
}>>([])

// 错误信息
const errorInfo = ref<{
  title: string
  message: string
} | null>(null)

// 轮询定时器
let progressTimer: NodeJS.Timeout | null = null

// 计算属性
const canSubmit = computed(() => {
  return formData.domainId && formData.email && !applying.value
})

// 获取进度类型
const getProgressType = (step: number) => {
  if (stepStatus.value === 'error') return 'danger'
  if (step === 3 && stepStatus.value === 'finish') return 'success'
  if (step > 0) return 'primary'
  return 'info'
}

// 获取进度文本
const getProgressText = (step: number) => {
  if (stepStatus.value === 'error') return '申请失败'
  if (step === 3 && stepStatus.value === 'finish') return '申请成功'
  if (step > 0) return '申请中'
  return '准备中'
}

// 数据加载
const loadDomains = async () => {
  try {
    const response = await domainApi.list({
      page: 1,
      page_size: 1000,
      status: 'active'
    })

    domainOptions.value = (response.data?.items || []).filter(domain =>
      domain.provider_id && domain.provider
    )
  } catch (error) {
    console.error('加载域名列表失败:', error)
    ElMessage.error('加载域名列表失败')
  }
}

// 域名变化处理
const handleDomainChange = (domainId: number) => {
  const domain = domainOptions.value.find(d => d.id === domainId)
  if (domain) {
    progressInfo.domain = domain.name
  }
}

// 添加日志
const addLog = (level: 'info' | 'success' | 'warning' | 'error', message: string) => {
  logs.value.push({
    time: new Date().toLocaleTimeString(),
    level,
    message
  })

  // 保持最新的20条日志
  if (logs.value.length > 20) {
    logs.value.shift()
  }
}

// 证书申请
const handleApply = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    applying.value = true
    showProgress.value = true
    currentStep.value = 0
    stepStatus.value = 'process'
    errorInfo.value = null
    logs.value = []

    // 设置进度信息
    const selectedDomain = domainOptions.value.find(d => d.id === formData.domainId)
    if (!selectedDomain) {
      throw new Error('请选择有效的域名')
    }

    progressInfo.domain = selectedDomain.name
    progressInfo.certType = formData.certType === 'wildcard' ? '通配符证书' :
                           formData.certType === 'multi' ? '多域名证书' : '单域名证书'
    progressInfo.startTime = new Date().toLocaleString()
    progressInfo.estimatedTime = new Date(Date.now() + 5 * 60 * 1000).toLocaleString()

    addLog('info', '开始申请SSL证书...')
    addLog('info', `域名: ${selectedDomain.name}`)
    addLog('info', `证书类型: ${progressInfo.certType}`)

    // 第一步：提交申请
    currentStep.value = 1
    addLog('info', '正在提交申请...')

    const requestData = {
      domain_id: formData.domainId!,
      domains: formData.certType === 'wildcard'
        ? [`*.${selectedDomain.name}`, selectedDomain.name]
        : [selectedDomain.name],
      email: formData.email,
      provider_id: selectedDomain.provider_id!,
      key_type: 'RSA2048',
      valid_days: 90,
      auto_renew: true,
      remark: formData.remark || `${progressInfo.certType}`
    }

    const response = await certificateApi.create(requestData)
    certificateId.value = response.data.id

    addLog('success', '申请提交成功')

    // 开始轮询进度
    startProgressPolling()

  } catch (error: any) {
    console.error('证书申请失败:', error)
    stepStatus.value = 'error'
    errorInfo.value = {
      title: '证书申请失败',
      message: error.message || '未知错误，请重试'
    }
    addLog('error', `申请失败: ${error.message || '未知错误'}`)
    applying.value = false
  }
}

// 开始进度轮询
const startProgressPolling = () => {
  if (!certificateId.value) return

  progressTimer = setInterval(async () => {
    try {
      const response = await certificateApi.get(certificateId.value!)
      const cert = response.data

      // 根据证书状态更新进度
      switch (cert.status) {
        case 'pending':
          currentStep.value = 1
          addLog('info', '等待域名验证...')
          break
        case 'processing':
          currentStep.value = 2
          addLog('info', '正在验证域名所有权...')
          break
        case 'issued':
          currentStep.value = 3
          stepStatus.value = 'finish'
          addLog('success', '证书申请成功！')
          addLog('success', `证书已签发，有效期至: ${new Date(cert.valid_to).toLocaleDateString()}`)
          stopProgressPolling()
          applying.value = false

          // 3秒后跳转到证书列表
          setTimeout(() => {
            ElMessage.success('证书申请成功！')
            router.push('/dns/certs')
          }, 3000)
          break
        case 'failed':
          stepStatus.value = 'error'
          errorInfo.value = {
            title: '证书申请失败',
            message: cert.error_message || '证书申请过程中发生错误'
          }
          addLog('error', `申请失败: ${cert.error_message || '未知错误'}`)
          stopProgressPolling()
          applying.value = false
          break
      }
    } catch (error) {
      console.error('获取证书状态失败:', error)
    }
  }, 3000) // 每3秒轮询一次
}

// 停止进度轮询
const stopProgressPolling = () => {
  if (progressTimer) {
    clearInterval(progressTimer)
    progressTimer = null
  }
}

// 重试申请
const retryApply = () => {
  errorInfo.value = null
  logs.value = []
  currentStep.value = 0
  stepStatus.value = 'process'
  handleApply()
}

// 重置表单
const resetForm = () => {
  showProgress.value = false
  errorInfo.value = null
  logs.value = []
  currentStep.value = 0
  stepStatus.value = 'process'
  applying.value = false
  certificateId.value = null
}

// 返回
const goBack = async () => {
  if (applying.value) {
    try {
      await ElMessageBox.confirm(
        '证书申请正在进行中，确定要离开吗？',
        '确认离开',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
    } catch {
      return
    }
  }

  stopProgressPolling()
  router.push('/dns/certs')
}

// 生命周期
onMounted(() => {
  loadDomains()
})

onUnmounted(() => {
  stopProgressPolling()
})
</script>

<style scoped lang="scss">
.cert-create-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;

  .page-header {
    background: white;
    border-radius: 12px;
    padding: 24px;
    margin-bottom: 24px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);

    .header-content {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .header-left {
        display: flex;
        align-items: center;
        gap: 16px;

        .back-btn {
          color: #666;
          &:hover {
            color: #409eff;
          }
        }

        .header-title {
          h1 {
            margin: 0;
            font-size: 24px;
            font-weight: 600;
            color: #303133;
          }

          p {
            margin: 4px 0 0 0;
            color: #909399;
            font-size: 14px;
          }
        }
      }

      .header-actions {
        display: flex;
        gap: 12px;
      }
    }
  }

  .progress-section {
    margin-bottom: 24px;

    .progress-card {
      background: white;
      border-radius: 12px;
      padding: 24px;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);

      .progress-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 24px;

        h3 {
          margin: 0;
          font-size: 18px;
          font-weight: 600;
          color: #303133;
        }
      }

      .progress-steps {
        margin-bottom: 24px;
      }

      .progress-details {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 16px;
        margin-bottom: 24px;
        padding: 16px;
        background: #f8f9fa;
        border-radius: 8px;

        .detail-item {
          display: flex;
          flex-direction: column;
          gap: 4px;

          .label {
            font-size: 12px;
            color: #909399;
            font-weight: 500;
          }

          .value {
            font-size: 14px;
            color: #303133;
            font-weight: 600;
          }
        }
      }

      .logs-section {
        margin-bottom: 24px;

        h4 {
          margin: 0 0 12px 0;
          font-size: 16px;
          font-weight: 600;
          color: #303133;
        }

        .logs-container {
          max-height: 200px;
          overflow-y: auto;
          background: #1e1e1e;
          border-radius: 8px;
          padding: 12px;

          .log-item {
            display: flex;
            gap: 12px;
            margin-bottom: 8px;
            font-family: 'Consolas', 'Monaco', monospace;
            font-size: 12px;

            .log-time {
              color: #888;
              min-width: 80px;
            }

            .log-message {
              flex: 1;
            }

            &.log-info .log-message {
              color: #61dafb;
            }

            &.log-success .log-message {
              color: #98d982;
            }

            &.log-warning .log-message {
              color: #ffb86c;
            }

            &.log-error .log-message {
              color: #ff5555;
            }
          }
        }
      }

      .error-section {
        .error-actions {
          margin-top: 16px;
          display: flex;
          gap: 12px;
        }
      }
    }
  }

  .form-section {
    .form-container {
      max-width: 1000px;
      margin: 0 auto;

      .steps-indicator {
        display: flex;
        justify-content: center;
        margin-bottom: 32px;
        gap: 40px;

        .step-item {
          display: flex;
          align-items: center;
          gap: 12px;
          padding: 16px 24px;
          background: white;
          border-radius: 12px;
          box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
          transition: all 0.3s ease;

          &.active {
            background: linear-gradient(135deg, #409eff, #67c23a);
            color: white;
            transform: translateY(-2px);
            box-shadow: 0 8px 25px rgba(64, 158, 255, 0.3);
          }

          .step-number {
            width: 32px;
            height: 32px;
            border-radius: 50%;
            background: #f0f0f0;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 600;
            font-size: 14px;
          }

          &.active .step-number {
            background: rgba(255, 255, 255, 0.2);
            color: white;
          }

          .step-info {
            .step-title {
              font-weight: 600;
              font-size: 14px;
              margin-bottom: 2px;
            }

            .step-desc {
              font-size: 12px;
              opacity: 0.8;
            }
          }
        }
      }

      .form-card {
        background: white;
        border-radius: 12px;
        padding: 24px;
        margin-bottom: 24px;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);

        .card-header {
          margin-bottom: 24px;

          h3 {
            margin: 0 0 8px 0;
            font-size: 18px;
            font-weight: 600;
            color: #303133;
            display: flex;
            align-items: center;
            gap: 8px;
          }

          p {
            margin: 0;
            color: #909399;
            font-size: 14px;
          }
        }

        .apply-types {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
          gap: 16px;

          .type-card {
            border: 2px solid #e4e7ed;
            border-radius: 12px;
            padding: 20px;
            cursor: pointer;
            transition: all 0.3s ease;
            display: flex;
            align-items: center;
            gap: 16px;

            &:hover {
              border-color: #409eff;
              transform: translateY(-2px);
              box-shadow: 0 8px 25px rgba(64, 158, 255, 0.15);
            }

            &.active {
              border-color: #409eff;
              background: linear-gradient(135deg, rgba(64, 158, 255, 0.1), rgba(103, 194, 58, 0.1));
            }

            .type-icon {
              width: 48px;
              height: 48px;
              border-radius: 12px;
              background: linear-gradient(135deg, #409eff, #67c23a);
              display: flex;
              align-items: center;
              justify-content: center;
              color: white;
              font-size: 20px;
            }

            .type-content {
              flex: 1;

              .type-title {
                font-weight: 600;
                font-size: 16px;
                color: #303133;
                margin-bottom: 4px;
                display: flex;
                align-items: center;
                gap: 8px;
              }

              .type-desc {
                color: #606266;
                font-size: 14px;
                margin-bottom: 8px;
              }

              .type-features {
                display: flex;
                flex-wrap: wrap;
                gap: 8px;

                .feature {
                  font-size: 12px;
                  color: #909399;
                }
              }
            }

            .type-radio {
              .el-radio {
                margin: 0;
              }
            }
          }
        }

        .cert-types {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
          gap: 16px;

          .cert-type-card {
            border: 2px solid #e4e7ed;
            border-radius: 12px;
            padding: 16px;
            cursor: pointer;
            transition: all 0.3s ease;
            display: flex;
            align-items: center;
            gap: 12px;

            &:hover {
              border-color: #409eff;
              transform: translateY(-1px);
              box-shadow: 0 4px 15px rgba(64, 158, 255, 0.15);
            }

            &.active {
              border-color: #409eff;
              background: rgba(64, 158, 255, 0.05);
            }

            .cert-icon {
              width: 40px;
              height: 40px;
              border-radius: 8px;
              background: linear-gradient(135deg, #409eff, #67c23a);
              display: flex;
              align-items: center;
              justify-content: center;
              color: white;
              font-size: 16px;
            }

            .cert-content {
              flex: 1;

              .cert-title {
                font-weight: 600;
                font-size: 14px;
                color: #303133;
                margin-bottom: 2px;
              }

              .cert-desc {
                color: #606266;
                font-size: 12px;
                margin-bottom: 2px;
              }

              .cert-example {
                color: #909399;
                font-size: 11px;
              }
            }

            .cert-radio {
              .el-radio {
                margin: 0;
              }
            }
          }
        }

        .domain-option {
          display: flex;
          justify-content: space-between;
          align-items: center;
          width: 100%;

          .domain-info {
            display: flex;
            align-items: center;
            gap: 8px;

            .domain-name {
              font-weight: 600;
              color: #303133;
            }
          }

          .domain-status {
            display: flex;
            align-items: center;
            gap: 4px;
            font-size: 12px;
            color: #67c23a;
          }
        }

        .form-tip {
          margin-top: 8px;
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 12px;
          color: #909399;
        }
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .cert-create-page {
    padding: 12px;

    .page-header .header-content {
      flex-direction: column;
      gap: 16px;
      align-items: flex-start;
    }

    .form-section .form-container .steps-indicator {
      flex-direction: column;
      gap: 12px;
    }

    .apply-types,
    .cert-types {
      grid-template-columns: 1fr !important;
    }
  }
}
</style>
