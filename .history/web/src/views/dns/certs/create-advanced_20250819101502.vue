<template>
  <div class="cert-create-advanced">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <el-button text @click="goBack" class="back-btn" size="large">
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <div class="header-title">
            <h1>
              <el-icon class="title-icon"><Lock /></el-icon>
              申请SSL证书
            </h1>
            <p>专业的SSL证书申请配置，支持多种验证方式和加密算法</p>
          </div>
        </div>
        <div class="header-actions">
          <el-button @click="goBack" size="large" class="cancel-btn">
            取消
          </el-button>
          <el-button
            type="primary"
            :loading="applying"
            @click="handleApply"
            :disabled="!canSubmit"
            size="large"
            class="apply-btn"
          >
            <el-icon><Lightning /></el-icon>
            {{ applying ? '申请中...' : '申请证书' }}
          </el-button>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px" class="cert-form">
        
        <!-- 证书域名配置 -->
        <div class="form-section">
          <div class="section-header">
            <el-icon><Globe /></el-icon>
            <span>证书域名</span>
            <span class="required">*</span>
          </div>
          <div class="section-content">
            <el-form-item prop="domains">
              <div class="domain-input-container">
                <el-input
                  v-model="domainInput"
                  placeholder="请输入域名，支持通配符如 *.example.com"
                  @keyup.enter="addDomain"
                  class="domain-input"
                  size="large"
                >
                  <template #append>
                    <el-button @click="addDomain" type="primary">添加</el-button>
                  </template>
                </el-input>
                <div class="domain-examples">
                  <span class="example-label">示例：</span>
                  <el-tag @click="addExampleDomain('example.com')" class="example-tag">example.com</el-tag>
                  <el-tag @click="addExampleDomain('*.example.com')" class="example-tag">*.example.com</el-tag>
                  <el-tag @click="addExampleDomain('sub.example.com')" class="example-tag">sub.example.com</el-tag>
                </div>
              </div>
              <div v-if="formData.domains.length > 0" class="domain-list">
                <div
                  v-for="(domain, index) in formData.domains"
                  :key="index"
                  class="domain-item"
                >
                  <span class="domain-name">{{ domain }}</span>
                  <el-button
                    @click="removeDomain(index)"
                    type="danger"
                    text
                    size="small"
                    class="remove-btn"
                  >
                    <el-icon><Close /></el-icon>
                  </el-button>
                </div>
              </div>
            </el-form-item>
          </div>
        </div>

        <!-- 邮箱配置 -->
        <div class="form-section">
          <div class="section-header">
            <el-icon><Message /></el-icon>
            <span>邮箱</span>
            <span class="required">*</span>
          </div>
          <div class="section-content">
            <el-form-item prop="email">
              <el-input
                v-model="formData.email"
                placeholder="申请人邮箱"
                size="large"
                prefix-icon="Message"
              />
            </el-form-item>
          </div>
        </div>

        <!-- 证书提供商 -->
        <div class="form-section">
          <div class="section-header">
            <el-icon><Star /></el-icon>
            <span>证书提供商</span>
          </div>
          <div class="section-content">
            <div class="provider-info">
              <div class="provider-logo">
                <img src="/images/letsencrypt-logo.png" alt="Let's Encrypt" />
              </div>
              <div class="provider-details">
                <div class="provider-name">Let's Encrypt</div>
                <div class="provider-desc">免费、自动化的证书颁发机构</div>
              </div>
              <el-tag type="success" size="large">免费</el-tag>
            </div>
          </div>
        </div>

        <!-- 域名验证方式 -->
        <div class="form-section">
          <div class="section-header">
            <el-icon><Shield /></el-icon>
            <span>域名验证方式</span>
            <span class="required">*</span>
          </div>
          <div class="section-content">
            <el-form-item prop="challengeType">
              <el-select
                v-model="formData.challengeType"
                placeholder="请选择验证方式"
                size="large"
                class="challenge-select"
              >
                <el-option value="dns" label="DNS直接验证">
                  <div class="option-content">
                    <div class="option-main">
                      <span class="option-label">DNS直接验证</span>
                      <el-tag type="success" size="small">推荐</el-tag>
                    </div>
                    <div class="option-desc">通过DNS记录验证域名所有权</div>
                  </div>
                </el-option>
                <el-option value="cname" label="CNAME代理验证">
                  <div class="option-content">
                    <div class="option-main">
                      <span class="option-label">CNAME代理验证</span>
                    </div>
                    <div class="option-desc">通过CNAME记录代理验证</div>
                  </div>
                </el-option>
                <el-option value="http" label="HTTP文件验证">
                  <div class="option-content">
                    <div class="option-main">
                      <span class="option-label">HTTP文件验证</span>
                    </div>
                    <div class="option-desc">通过HTTP文件验证域名所有权</div>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </div>
        </div>

        <!-- DNS解析服务商 -->
        <div class="form-section">
          <div class="section-header">
            <el-icon><Connection /></el-icon>
            <span>DNS解析服务商</span>
            <span class="required">*</span>
          </div>
          <div class="section-content">
            <el-form-item prop="providerId">
              <el-select
                v-model="formData.providerId"
                placeholder="请选择DNS解析服务商"
                filterable
                size="large"
                class="provider-select"
              >
                <el-option
                  v-for="provider in dnsProviders"
                  :key="provider.id"
                  :value="provider.id"
                  :label="provider.name"
                >
                  <div class="provider-option">
                    <div class="provider-icon">
                      <img :src="getProviderIcon(provider.type)" :alt="provider.name" />
                    </div>
                    <div class="provider-info">
                      <div class="provider-name">{{ provider.name }}</div>
                      <div class="provider-type">{{ getProviderTypeName(provider.type) }}</div>
                    </div>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </div>
        </div>

        <!-- 加密算法 -->
        <div class="form-section">
          <div class="section-header">
            <el-icon><Key /></el-icon>
            <span>加密算法</span>
          </div>
          <div class="section-content">
            <el-form-item prop="keyType">
              <el-select
                v-model="formData.keyType"
                placeholder="请选择加密算法"
                size="large"
                class="key-type-select"
              >
                <el-option value="RSA2048" label="RSA 2048">
                  <div class="option-content">
                    <div class="option-main">
                      <span class="option-label">RSA 2048</span>
                      <el-tag type="primary" size="small">推荐</el-tag>
                    </div>
                    <div class="option-desc">兼容性好，安全性高</div>
                  </div>
                </el-option>
                <el-option value="RSA1024" label="RSA 1024">
                  <div class="option-content">
                    <div class="option-main">
                      <span class="option-label">RSA 1024</span>
                    </div>
                    <div class="option-desc">较低安全性，不推荐</div>
                  </div>
                </el-option>
                <el-option value="RSA3072" label="RSA 3072">
                  <div class="option-content">
                    <div class="option-main">
                      <span class="option-label">RSA 3072</span>
                    </div>
                    <div class="option-desc">高安全性，文件较大</div>
                  </div>
                </el-option>
                <el-option value="RSA4096" label="RSA 4096">
                  <div class="option-content">
                    <div class="option-main">
                      <span class="option-label">RSA 4096</span>
                    </div>
                    <div class="option-desc">最高安全性，文件最大</div>
                  </div>
                </el-option>
                <el-option value="EC256" label="EC 256">
                  <div class="option-content">
                    <div class="option-main">
                      <span class="option-label">EC 256</span>
                    </div>
                    <div class="option-desc">椭圆曲线，高效安全</div>
                  </div>
                </el-option>
                <el-option value="EC384" label="EC 384">
                  <div class="option-content">
                    <div class="option-main">
                      <span class="option-label">EC 384</span>
                    </div>
                    <div class="option-desc">椭圆曲线，超高安全</div>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </div>
        </div>

        <!-- 更新天数 -->
        <div class="form-section">
          <div class="section-header">
            <el-icon><Calendar /></el-icon>
            <span>更新天数</span>
          </div>
          <div class="section-content">
            <el-form-item prop="renewDays">
              <el-input-number
                v-model="formData.renewDays"
                :min="1"
                :max="90"
                size="large"
                class="renew-days-input"
              />
              <div class="form-tip">
                证书到期前多少天自动续期，注意：流水线多少天不后续新证书，请谨慎填写
              </div>
            </el-form-item>
          </div>
        </div>

        <!-- 定时触发 -->
        <div class="form-section">
          <div class="section-header">
            <el-icon><Clock /></el-icon>
            <span>定时触发</span>
          </div>
          <div class="section-content">
            <div class="schedule-config">
              <div class="schedule-item">
                <span class="schedule-label">每天的</span>
                <el-time-picker
                  v-model="formData.scheduleTime"
                  format="HH:mm"
                  placeholder="选择时间"
                  size="large"
                />
              </div>
              <div class="schedule-note">
                建议设置为凌晨时间，避免影响业务
              </div>
            </div>
          </div>
        </div>

        <!-- 失败通知 -->
        <div class="form-section">
          <div class="section-header">
            <el-icon><Bell /></el-icon>
            <span>失败通知</span>
          </div>
          <div class="section-content">
            <div class="notification-config">
              <el-radio-group v-model="formData.notificationType" size="large">
                <el-radio value="default">使用默认通知</el-radio>
                <el-radio value="custom">自定义通知</el-radio>
              </el-radio-group>
              <div v-if="formData.notificationType === 'custom'" class="custom-notification">
                <el-input
                  v-model="formData.notificationEmail"
                  placeholder="请输入通知邮箱"
                  size="large"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- 备注 -->
        <div class="form-section">
          <div class="section-header">
            <el-icon><Document /></el-icon>
            <span>流水线备注</span>
          </div>
          <div class="section-content">
            <el-form-item prop="remark">
              <el-input
                v-model="formData.remark"
                type="textarea"
                :rows="3"
                placeholder="请输入备注信息（可选）"
                maxlength="200"
                show-word-limit
                size="large"
              />
            </el-form-item>
          </div>
        </div>

      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft,
  Lock,
  Lightning,
  Globe,
  Message,
  Star,
  Shield,
  Connection,
  Key,
  Calendar,
  Clock,
  Bell,
  Document,
  Close
} from '@element-plus/icons-vue'
import { domainApi } from '@/api/dns/domain'
import { providerApi } from '@/api/dns/provider'
import { certificateApi } from '@/api/dns/certificate'

const router = useRouter()

// 表单数据
const formData = reactive({
  domains: [] as string[],
  email: '',
  challengeType: 'dns',
  providerId: null as number | null,
  keyType: 'RSA2048',
  renewDays: 35,
  scheduleTime: null,
  notificationType: 'default',
  notificationEmail: '',
  remark: ''
})

// 表单验证规则
const formRules = {
  domains: [
    { required: true, message: '请至少添加一个域名', trigger: 'change' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  challengeType: [
    { required: true, message: '请选择验证方式', trigger: 'change' }
  ],
  providerId: [
    { required: true, message: '请选择DNS解析服务商', trigger: 'change' }
  ]
}

// 状态管理
const applying = ref(false)
const domainInput = ref('')
const dnsProviders = ref<any[]>([])
const formRef = ref()

// 计算属性
const canSubmit = computed(() => {
  return formData.domains.length > 0 && 
         formData.email && 
         formData.challengeType && 
         formData.providerId && 
         !applying.value
})

// 方法
const goBack = () => {
  router.push('/dns/certs')
}

const addDomain = () => {
  const domain = domainInput.value.trim()
  if (domain && !formData.domains.includes(domain)) {
    if (validateDomain(domain)) {
      formData.domains.push(domain)
      domainInput.value = ''
    } else {
      ElMessage.error('请输入有效的域名格式')
    }
  }
}

const addExampleDomain = (domain: string) => {
  if (!formData.domains.includes(domain)) {
    formData.domains.push(domain)
  }
}

const removeDomain = (index: number) => {
  formData.domains.splice(index, 1)
}

const validateDomain = (domain: string) => {
  // 简单的域名验证
  const domainRegex = /^(\*\.)?[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/
  return domainRegex.test(domain)
}

const getProviderIcon = (type: string) => {
  const icons = {
    aliyun: '/images/providers/aliyun.png',
    tencent: '/images/providers/tencent.png',
    aws: '/images/providers/aws.png',
    godaddy: '/images/providers/godaddy.png',
    cloudflare: '/images/providers/cloudflare.png'
  }
  return icons[type] || '/images/providers/default.png'
}

const getProviderTypeName = (type: string) => {
  const names = {
    aliyun: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS',
    godaddy: 'GoDaddy',
    cloudflare: 'Cloudflare'
  }
  return names[type] || type
}

const loadDnsProviders = async () => {
  try {
    const response = await providerApi.list({
      page: 1,
      page_size: 100,
      status: 'active'
    })
    dnsProviders.value = response.data?.items || []
  } catch (error) {
    console.error('加载DNS提供商失败:', error)
    ElMessage.error('加载DNS提供商失败')
  }
}

const handleApply = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    applying.value = true

    const requestData = {
      domains: formData.domains,
      email: formData.email,
      challenge_type: formData.challengeType,
      provider_id: formData.providerId,
      key_type: formData.keyType,
      renew_days: formData.renewDays,
      schedule_time: formData.scheduleTime,
      notification_type: formData.notificationType,
      notification_email: formData.notificationEmail,
      remark: formData.remark,
      auto_renew: true,
      valid_days: 90
    }

    const response = await certificateApi.create(requestData)
    ElMessage.success('证书申请成功')
    router.push('/dns/certs')
  } catch (error: any) {
    console.error('证书申请失败:', error)
    ElMessage.error(error.message || '证书申请失败')
  } finally {
    applying.value = false
  }
}

// 生命周期
onMounted(() => {
  loadDnsProviders()
})
</script>
