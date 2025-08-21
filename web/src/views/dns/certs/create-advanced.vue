<template>
  <div class="cert-create-page">
    <!-- 简化后的紧凑头部 -->
    <div class="page-header">
      <div class="header-left">
        <el-button text @click="goBack" class="back-btn" size="small">
          <el-icon><ArrowLeft /></el-icon>
          返回证书管理
        </el-button>
        <div class="divider"></div>
        <h1>
          <span class="icon">🔒</span>
          申请SSL证书
        </h1>
      </div>
      <div class="header-actions">
        <el-button @click="goBack" class="cancel-btn">
          取消
        </el-button>
        <el-button
          type="primary"
          :loading="applying"
          @click="handleApply"
          :disabled="!canSubmit"
          class="apply-btn"
        >
          <el-icon><Lightning /></el-icon>
          {{ applying ? '申请中...' : '申请证书' }}
        </el-button>
      </div>
    </div>

    <!-- 主要内容：左右分栏布局 -->
    <div class="main-container">
      <!-- 左侧：配置表单 -->
      <div class="config-panel">
        <el-form ref="formRef" :model="formData" :rules="formRules" class="cert-form">
          
          <!-- 域名配置卡片 -->
          <div class="config-card">
            <div class="card-header">
              <span class="icon">🌐</span>
              <span class="title">域名配置</span>
              <span class="required">*</span>
            </div>
            <div class="card-content">
              <!-- 域名选择 -->
              <div class="form-row">
                <el-select
                  v-model="formData.domainIds"
                  multiple
                  placeholder="选择已有域名"
                  class="domain-select"
                  @change="handleDomainChange"
                >
                  <el-option
                    v-for="domain in availableDomains"
                    :key="domain.id"
                    :label="domain.name"
                    :value="domain.id"
                  >
                    <div class="domain-option">
                      <span>{{ domain.name }}</span>
                      <el-tag :type="domain.status === 'active' ? 'success' : 'warning'" size="small">
                        {{ domain.status === 'active' ? '正常' : '异常' }}
                      </el-tag>
                    </div>
                  </el-option>
                </el-select>
              </div>
              
              <!-- 自定义域名输入 -->
              <div class="form-row">
                <el-input
                  v-model="customDomainInput"
                  placeholder="或输入自定义域名，如：example.com 或 *.example.com"
                  @keyup.enter="addCustomDomain"
                >
                  <template #append>
                    <el-button @click="addCustomDomain" :disabled="!customDomainInput.trim()">
                      添加
                    </el-button>
                  </template>
                </el-input>
              </div>
              
              <!-- 已选域名显示 -->
              <div v-if="formData.domains.length > 0" class="selected-domains">
                <el-tag
                  v-for="domain in formData.domains"
                  :key="domain"
                  closable
                  @close="removeDomainByName(domain)"
                  :type="domain.startsWith('*') ? 'warning' : 'primary'"
                >
                  {{ domain }}
                </el-tag>
              </div>
              
              <!-- 示例域名 -->
              <div class="quick-add">
                <span class="label">快速添加：</span>
                <el-tag
                  v-for="example in domainExamples.slice(0, 3)"
                  :key="example"
                  size="small"
                  class="example-tag"
                  @click="addExampleDomain(example)"
                >
                  {{ example }}
                </el-tag>
              </div>
            </div>
          </div>

          <!-- 基础信息卡片 -->
          <div class="config-card">
            <div class="card-header">
              <span class="icon">📧</span>
              <span class="title">基础信息</span>
            </div>
            <div class="card-content">
              <div class="form-grid">
                <el-form-item label="申请人邮箱" prop="email">
                  <el-input
                    v-model="formData.email"
                    placeholder="用于接收证书通知"
                    prefix-icon="Message"
                  />
                </el-form-item>
                <el-form-item label="证书提供商" prop="caType">
                  <el-select v-model="formData.caType">
                    <el-option
                      v-for="provider in caProviders"
                      :key="provider.type"
                      :label="provider.name"
                      :value="provider.type"
                    >
                      <div class="provider-info">
                        <span>{{ provider.name }}</span>
                        <el-tag v-if="provider.free" type="success" size="small">免费</el-tag>
                      </div>
                    </el-option>
                  </el-select>
                </el-form-item>
              </div>
            </div>
          </div>

          <!-- 验证方式卡片 -->
          <div class="config-card">
            <div class="card-header">
              <span class="icon">🔍</span>
              <span class="title">验证方式</span>
            </div>
            <div class="card-content">
              <div class="verification-tabs">
                <div 
                  v-for="method in verificationMethods"
                  :key="method.value"
                  class="tab-item"
                  :class="{ active: formData.challengeType === method.value }"
                  @click="formData.challengeType = method.value"
                >
                  <div class="tab-icon">{{ method.icon }}</div>
                  <div class="tab-content">
                    <div class="tab-title">{{ method.label }}</div>
                    <div class="tab-desc">{{ method.desc }}</div>
                  </div>
                  <el-tag v-if="method.recommended" type="success" size="small">推荐</el-tag>
                </div>
              </div>
              
              <!-- DNS提供商选择 -->
              <div v-if="formData.challengeType === 'dns' || formData.challengeType === 'cname'" class="provider-section">
                <el-form-item label="DNS解析服务商" prop="providerId">
                  <el-select
                    v-model="formData.providerId"
                    placeholder="选择DNS服务商"
                    filterable
                  >
                    <el-option
                      v-for="provider in dnsProviders"
                      :key="Number(provider.id)"
                      :value="Number(provider.id)"
                      :label="provider.name"
                    >
                      <div class="provider-option">
                        <ProviderIcon :type="provider.type" size="20px" />
                        <span>{{ provider.name }}</span>
                        <span class="provider-type">{{ getProviderTypeName(provider.type) }}</span>
                      </div>
                    </el-option>
                  </el-select>
                </el-form-item>
              </div>
              
              <!-- HTTP验证说明 -->
              <div v-if="formData.challengeType === 'http'" class="http-notice">
                <el-alert type="info" :closable="false">
                  <p>HTTP验证需要在网站根目录创建验证文件，请确保域名可正常访问</p>
                </el-alert>
              </div>
            </div>
          </div>

          <!-- 高级选项（可折叠） -->
          <div class="config-card">
            <div class="card-header" @click="showAdvanced = !showAdvanced">
              <span class="icon">⚙️</span>
              <span class="title">高级选项</span>
              <el-icon class="expand-icon" :class="{ expanded: showAdvanced }">
                <ArrowDown />
              </el-icon>
            </div>
            <el-collapse-transition>
              <div v-show="showAdvanced" class="card-content">
                <div class="form-grid">
                  <el-form-item label="加密算法">
                    <el-select v-model="formData.keyType">
                      <el-option value="RSA2048" label="RSA 2048（推荐）" />
                      <el-option value="RSA4096" label="RSA 4096" />
                      <el-option value="EC256" label="EC 256" />
                    </el-select>
                  </el-form-item>
                  <el-form-item label="续期天数">
                    <el-input-number
                      v-model="formData.renewDays"
                      :min="1"
                      :max="90"
                    />
                  </el-form-item>
                </div>
                
                <div class="toggle-options">
                  <div class="toggle-item">
                    <el-switch v-model="formData.autoRenew" />
                    <span class="toggle-label">启用自动续期</span>
                  </div>
                  <div class="toggle-item">
                    <el-switch v-model="formData.deployEnabled" />
                    <span class="toggle-label">启用自动部署</span>
                  </div>
                </div>
                
                <!-- 部署配置 -->
                <div v-if="formData.deployEnabled" class="deploy-section">
                  <div class="section-title">部署配置</div>
                  <el-form-item label="目标主机">
                    <el-select
                      v-model="formData.deployHosts"
                      multiple
                      placeholder="选择部署主机"
                      :loading="hostsLoading"
                    >
                      <el-option
                        v-for="host in availableHosts"
                        :key="host.id"
                        :label="`${host.name} (${host.ip})`"
                        :value="host.id"
                      />
                    </el-select>
                  </el-form-item>
                  <el-form-item label="部署路径">
                    <el-input
                      v-model="formData.deployPath"
                      placeholder="/etc/ssl/certs/"
                    />
                  </el-form-item>
                </div>
                
                <!-- 通知配置 -->
                <div class="notification-section">
                  <div class="section-title">通知配置</div>
                  <el-form-item label="通知方式">
                    <el-radio-group v-model="formData.notificationType">
                      <el-radio value="default">使用申请邮箱</el-radio>
                      <el-radio value="custom">自定义邮箱</el-radio>
                      <el-radio value="none">不接收通知</el-radio>
                    </el-radio-group>
                  </el-form-item>
                  <el-form-item v-if="formData.notificationType === 'custom'" label="通知邮箱">
                    <el-input
                      v-model="formData.notificationEmail"
                      placeholder="custom@example.com"
                    />
                  </el-form-item>
                </div>
              </div>
            </el-collapse-transition>
          </div>

        </el-form>
      </div>

      <!-- 右侧：预览和帮助 -->
      <div class="info-panel">
        <!-- 配置预览 -->
        <div class="preview-card">
          <div class="card-title">
            <span class="icon">👀</span>
            配置预览
          </div>
          <div class="preview-content">
            <div class="preview-item">
              <span class="label">域名数量：</span>
              <span class="value">{{ formData.domains.length }} 个</span>
            </div>
            <div class="preview-item">
              <span class="label">验证方式：</span>
              <span class="value">{{ getVerificationName(formData.challengeType) }}</span>
            </div>
            <div class="preview-item">
              <span class="label">证书提供商：</span>
              <span class="value">{{ getCAName(formData.caType) }}</span>
            </div>
            <div class="preview-item">
              <span class="label">加密算法：</span>
              <span class="value">{{ formData.keyType }}</span>
            </div>
            <div class="preview-item">
              <span class="label">自动续期：</span>
              <span class="value">{{ formData.autoRenew ? '已启用' : '已禁用' }}</span>
            </div>
            <div class="preview-item">
              <span class="label">自动部署：</span>
              <span class="value">{{ formData.deployEnabled ? '已启用' : '已禁用' }}</span>
            </div>
          </div>
        </div>

        <!-- 域名列表 -->
        <div v-if="formData.domains.length > 0" class="domains-card">
          <div class="card-title">
            <span class="icon">📋</span>
            选择的域名
          </div>
          <div class="domains-list">
            <div
              v-for="domain in formData.domains"
              :key="domain"
              class="domain-item"
            >
              <span class="domain-name">{{ domain }}</span>
              <el-tag v-if="domain.startsWith('*')" type="warning" size="small">通配符</el-tag>
              <el-button
                text
                type="danger"
                size="small"
                @click="removeDomainByName(domain)"
              >
                <el-icon><Close /></el-icon>
              </el-button>
            </div>
          </div>
        </div>

        <!-- 帮助信息 -->
        <div class="help-card">
          <div class="card-title">
            <span class="icon">💡</span>
            使用提示
          </div>
          <div class="help-content">
            <div class="help-item">
              <strong>域名格式：</strong>
              <p>支持单域名（example.com）和通配符域名（*.example.com）</p>
            </div>
            <div class="help-item">
              <strong>验证方式：</strong>
              <p>DNS验证适合大多数场景，HTTP验证需要网站可访问</p>
            </div>
            <div class="help-item">
              <strong>自动续期：</strong>
              <p>推荐开启，避免证书过期导致的服务中断</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft,
  Lightning,
  ArrowDown,
  Close
} from '@element-plus/icons-vue'
import { domainApi } from '@/api/dns/domain'
import { dnsProviderApi } from '@/api/dns/provider'
import { certificateApi } from '@/api/dns/certificate'
import { getHostListWithCount } from '@/api/system/host'
import ProviderIcon from '@/components/dns/ProviderIcon.vue'

const router = useRouter()

// 表单数据
const formData = reactive({
  domainIds: [] as number[],
  domains: [] as string[],
  email: '',
  caType: 'letsencrypt',
  challengeType: 'dns',
  providerId: null as number | null, // 绑定为数字类型，配合选项Number(id)
  keyType: 'RSA2048',
  autoRenew: true,
  renewDays: 35,
  notificationType: 'default',
  notificationEmail: '',
  deployEnabled: false,
  deployHosts: [] as number[],
  deployPath: '/etc/ssl/certs/',
  restartCommand: '',
  remark: ''
})

// 表单验证规则
const formRules = {
  domainIds: [
    {
      validator: (rule: any, value: any, callback: any) => {
        // 检查是否选择了域名或者手动输入了域名
        if ((!value || value.length === 0) && (!formData.domains || formData.domains.length === 0)) {
          callback(new Error('请至少选择一个域名或手动输入域名'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  caType: [
    { required: true, message: '请选择证书提供商', trigger: 'change' }
  ],
  challengeType: [
    { required: true, message: '请选择验证方式', trigger: 'change' }
  ],
  providerId: [
    {
      validator: (rule: any, value: any, callback: any) => {
        // 只有选择DNS验证方式时才需要DNS解析服务商
        if (formData.challengeType === 'dns' || formData.challengeType === 'cname') {
          if (!value) {
            callback(new Error('请选择DNS解析服务商'))
          } else {
            callback()
          }
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ],
  notificationEmail: [
    {
      validator: (rule: any, value: any, callback: any) => {
        // 只有选择自定义通知时才需要验证邮箱
        if (formData.notificationType === 'custom') {
          if (!value || !value.trim()) {
            callback(new Error('请输入通知邮箱'))
          } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
            callback(new Error('请输入正确的邮箱格式'))
          } else {
            callback()
          }
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  deployHosts: [
    {
      validator: (rule: any, value: any, callback: any) => {
        // 只有启用自动部署时才需要选择主机
        if (formData.deployEnabled) {
          if (!value || value.length === 0) {
            callback(new Error('请选择要部署的主机'))
          } else {
            callback()
          }
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ],
  deployPath: [
    {
      validator: (rule: any, value: any, callback: any) => {
        // 只有启用自动部署时才需要部署路径
        if (formData.deployEnabled) {
          if (!value || !value.trim()) {
            callback(new Error('请输入部署路径'))
          } else {
            callback()
          }
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 状态管理
const applying = ref(false)
const hostsLoading = ref(false)
const customDomainInput = ref('')
const showAdvanced = ref(false)
const dnsProviders = ref<any[]>([])
const availableDomains = ref<any[]>([])
const availableHosts = ref<any[]>([])
const domainExamples = ref([
  'example.com',
  '*.example.com',
  'www.example.com'
])
const caProviders = ref([
  {
    type: 'letsencrypt',
    name: "Let's Encrypt",
    description: '免费、自动化的证书颁发机构',
    free: true
  },
  {
    type: 'zerossl',
    name: 'ZeroSSL',
    description: '免费SSL证书提供商',
    free: true
  }
])
const verificationMethods = ref([
  {
    value: 'dns',
    label: 'DNS验证',
    desc: '自动创建DNS记录验证',
    icon: '🌐',
    recommended: true
  },
  {
    value: 'http',
    label: 'HTTP验证',
    desc: '网站根目录文件验证',
    icon: '📁',
    recommended: false
  },
  {
    value: 'cname',
    label: 'CNAME验证',
    desc: '通过CNAME记录代理验证',
    icon: '🔗',
    recommended: false
  }
])
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

const handleDomainChange = (domainIds: number[]) => {
  // 根据选中的域名ID更新域名列表
  formData.domains = domainIds.map(id => {
    const domain = availableDomains.value.find(d => d.id === id)
    return domain ? domain.name : ''
  }).filter(Boolean)
}

const removeDomainByName = (domainName: string) => {
  const index = formData.domains.indexOf(domainName)
  if (index > -1) {
    formData.domains.splice(index, 1)
    // 同时移除对应的ID
    const domain = availableDomains.value.find(d => d.name === domainName)
    if (domain) {
      const idIndex = formData.domainIds.indexOf(domain.id)
      if (idIndex > -1) {
        formData.domainIds.splice(idIndex, 1)
      }
    }
  }
}

const validateDomain = (domain: string) => {
  // 增强的域名验证，支持通配符和多级域名
  const domainRegex = /^(\*\.)?[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$/
  
  // 检查基本格式
  if (!domainRegex.test(domain)) {
    return false
  }
  
  // 检查域名长度不超过253个字符
  if (domain.length > 253) {
    return false
  }
  
  // 检查每个标签不超过63个字符
  const labels = domain.split('.')
  for (const label of labels) {
    if (label.length > 63 || label.length === 0) {
      return false
    }
    
    // 标签不能以-开头或结尾
    if (label.startsWith('-') || label.endsWith('-')) {
      return false
    }
  }
  
  // 对于通配符域名，确保通配符只在最前面
  if (domain.includes('*') && !domain.startsWith('*.')) {
    return false
  }
  
  return true
}

// 添加自定义域名
const addCustomDomain = () => {
  const domain = customDomainInput.value.trim()
  if (!domain) {
    ElMessage.warning('请输入域名')
    return
  }

  if (!validateDomain(domain)) {
    ElMessage.error('域名格式不正确')
    return
  }

  if (formData.domains.includes(domain)) {
    ElMessage.warning('域名已存在')
    return
  }

  formData.domains.push(domain)
  customDomainInput.value = ''
  ElMessage.success('域名添加成功')
}

// 添加示例域名
const addExampleDomain = (example: string) => {
  if (!formData.domains.includes(example)) {
    formData.domains.push(example)
    ElMessage.success(`已添加域名：${example}`)
  } else {
    ElMessage.warning('域名已存在')
  }
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

const getVerificationName = (type: string) => {
  const method = verificationMethods.value.find(m => m.value === type)
  return method ? method.label : type
}

const getCAName = (type: string) => {
  const provider = caProviders.value.find(p => p.type === type)
  return provider ? provider.name : type
}

const loadDnsProviders = async () => {
  try {
    const response = await dnsProviderApi.list({
      page: 1,
      size: 100
    })
    const data: any = (response as any).data || response
    dnsProviders.value = (data?.items || []).map((p: any) => ({ ...p, id: Number(p.id) }))
  } catch (error) {
    console.error('加载DNS提供商失败:', error)
    ElMessage.error('加载DNS提供商失败')
  }
}

const loadDomains = async () => {
  try {
    const response = await domainApi.list({
      page: 1,
      size: 100,
      status: 'active'
    })
    const data: any = (response as any).data || response
    availableDomains.value = data?.items || []
  } catch (error) {
    console.error('加载域名列表失败:', error)
    ElMessage.error('加载域名列表失败')
  }
}

const loadHosts = async () => {
  try {
    hostsLoading.value = true
    
    // 使用正确的CMDB主机查询API
    const response = await getHostListWithCount({
      page: 1,
      page_size: 100
      // 移除status过滤，显示所有主机
    })

    console.log('证书创建 - CMDB主机API响应:', response)
    
    // 根据API响应结构解析数据
    // getHostListWithCount返回格式: {code: 200, data: [...], count: 10}
    let hostList = []
    
    if (response && response.data && Array.isArray(response.data)) {
      // 标准API响应格式
      hostList = response.data
    } else if (response && Array.isArray(response)) {
      // 直接数组格式
      hostList = response
    } else {
      console.warn('证书创建 - 未识别的主机API响应格式:', response)
      hostList = []
    }
    
    // 映射主机数据，确保字段完整
    availableHosts.value = hostList.map((host: any) => ({
      id: host.id,
      name: host.name || host.hostname || `主机-${host.id}`,
      ip: host.ip || host.private_ip || host.public_ip || '未知IP',
      status: host.status || 'unknown',
      provider: host.provider || host.provider_name || '未知',
      region: host.region || '未知'
    }))
    
    console.log(`证书创建 - 成功加载 ${availableHosts.value.length} 台主机:`, availableHosts.value)
  } catch (error) {
    console.error('证书创建 - 加载主机列表失败:', error)
    ElMessage.error('加载主机列表失败，请确保CMDB主机管理模块正常运行')
  } finally {
    hostsLoading.value = false
  }
}

const handleApply = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    applying.value = true

    const requestData = {
      domain_id: Array.isArray(formData.domainIds) && formData.domainIds.length > 0 ? Number(formData.domainIds[0]) : undefined,
      domains: formData.domains,
      email: formData.email,
      ca_type: formData.caType,
      challenge_type: formData.challengeType,
      provider_id: Number(formData.providerId),
      key_type: formData.keyType,
      auto_renew: formData.autoRenew,
      renew_days: formData.renewDays,
      notification_type: formData.notificationType,
      notification_email: formData.notificationEmail || undefined,
      deploy_hosts: formData.deployEnabled ? formData.deployHosts : [],
      deploy_path: formData.deployEnabled ? formData.deployPath : '',
      restart_command: formData.deployEnabled ? formData.restartCommand : '',
      remark: formData.remark,
      valid_days: 90
    }

    // 使用异步提交模式，立即返回成功并提示用户
    ElMessage.success('证书申请已提交，正在后台处理中...')
    
    // 跳转到证书列表页面
    router.push('/dns/certs')
    
    // 在后台异步处理申请
    certificateApi.create(requestData).then(() => {
      console.log('证书申请处理完成')
    }).catch(error => {
      console.error('证书申请失败:', error)
      // 可以在这里添加一些后台错误处理逻辑
    }).finally(() => {
      applying.value = false
    })
    
  } catch (error: any) {
    console.error('证书申请失败:', error)
    if (error.message?.includes('验证') || error.errors) {
      ElMessage.error('请检查输入信息是否正确')
    } else {
      ElMessage.error(error.message || '证书申请失败，请稍后重试')
    }
    applying.value = false
  }
}

// 生命周期

onMounted(() => {
  loadDnsProviders()
  loadDomains()
  loadHosts()
})
</script>

<style scoped lang="scss">
.cert-create-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

// 页面头部
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  padding: 16px 0;
  
  .header-left {
    display: flex;
    align-items: center;
    gap: 16px;
    
    .back-btn {
      color: rgba(255, 255, 255, 0.8);
      font-weight: 500;
      transition: all 0.3s ease;
      
      &:hover {
        color: white;
        transform: translateX(-2px);
      }
    }
    
    .divider {
      width: 1px;
      height: 20px;
      background: rgba(255, 255, 255, 0.3);
    }
    
    h1 {
      margin: 0;
      font-size: 24px;
      font-weight: 700;
      color: white;
      display: flex;
      align-items: center;
      gap: 8px;
      
      .icon {
        font-size: 24px;
      }
    }
  }
  
  .header-actions {
    display: flex;
    gap: 12px;
    
    .cancel-btn {
      background: rgba(255, 255, 255, 0.1);
      border: 1px solid rgba(255, 255, 255, 0.3);
      color: white;
      
      &:hover {
        background: rgba(255, 255, 255, 0.2);
      }
    }
    
    .apply-btn {
      background: rgba(255, 255, 255, 0.9);
      color: #667eea;
      border: none;
      font-weight: 600;
      
      &:hover {
        background: white;
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
      }
    }
  }
}

// 主要容器：左右分栏
.main-container {
  display: grid;
  grid-template-columns: 1fr 350px;
  gap: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

// 配置面板
.config-panel {
  .cert-form {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
}

// 配置卡片
.config-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  
  .card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 16px 20px;
    background: #f8fafc;
    border-bottom: 1px solid #e2e8f0;
    font-weight: 600;
    color: #2d3748;
    cursor: pointer;
    transition: all 0.3s ease;
    
    .icon {
      font-size: 18px;
    }
    
    .title {
      flex: 1;
    }
    
    .required {
      color: #f56c6c;
    }
    
    .expand-icon {
      transition: transform 0.3s ease;
      
      &.expanded {
        transform: rotate(180deg);
      }
    }
    
    &:hover {
      background: #edf2f7;
    }
  }
  
  .card-content {
    padding: 20px;
  }
}

// 表单网格
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-row {
  margin-bottom: 12px;
}

// 域名选择相关
.domain-select {
  width: 100%;
}

.domain-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.selected-domains {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.quick-add {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 12px;
  
  .label {
    color: #666;
    font-size: 14px;
    white-space: nowrap;
  }
  
  .example-tag {
    cursor: pointer;
    transition: all 0.3s ease;
    
    &:hover {
      transform: scale(1.05);
    }
  }
}

// 验证方式选项卡
.verification-tabs {
  display: grid;
  gap: 12px;
  margin-bottom: 16px;
  
  .tab-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    border: 2px solid #e2e8f0;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.3s ease;
    
    .tab-icon {
      font-size: 20px;
    }
    
    .tab-content {
      flex: 1;
      
      .tab-title {
        font-weight: 600;
        color: #2d3748;
        margin-bottom: 2px;
      }
      
      .tab-desc {
        font-size: 12px;
        color: #718096;
      }
    }
    
    &.active {
      border-color: #667eea;
      background: linear-gradient(135deg, #f0f4ff, #e6f2ff);
    }
    
    &:hover:not(.active) {
      border-color: #cbd5e0;
      background: #f7fafc;
    }
  }
}

// 提供商选项
.provider-option {
  display: flex;
  align-items: center;
  gap: 8px;
  
  .provider-type {
    font-size: 12px;
    color: #718096;
    margin-left: auto;
  }
}

.provider-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

// 提供商部分
.provider-section {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e2e8f0;
}

// HTTP通知
.http-notice {
  margin-top: 16px;
  
  .el-alert {
    border-radius: 8px;
  }
}

// 高级选项
.toggle-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin: 16px 0;
  
  .toggle-item {
    display: flex;
    align-items: center;
    gap: 12px;
    
    .toggle-label {
      font-weight: 500;
      color: #2d3748;
    }
  }
}

.deploy-section,
.notification-section {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e2e8f0;
  
  .section-title {
    font-size: 14px;
    font-weight: 600;
    color: #4a5568;
    margin-bottom: 12px;
  }
}

// 信息面板
.info-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: sticky;
  top: 20px;
  height: fit-content;
}

// 预览卡片
.preview-card,
.domains-card,
.help-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  
  .card-title {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 16px;
    font-weight: 600;
    color: #2d3748;
    
    .icon {
      font-size: 16px;
    }
  }
}

.preview-content {
  .preview-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 0;
    border-bottom: 1px solid #f1f5f9;
    
    &:last-child {
      border-bottom: none;
    }
    
    .label {
      color: #718096;
      font-size: 14px;
    }
    
    .value {
      color: #2d3748;
      font-weight: 500;
      font-size: 14px;
    }
  }
}

.domains-list {
  .domain-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 12px;
    margin: 4px 0;
    background: #f8fafc;
    border-radius: 8px;
    
    .domain-name {
      flex: 1;
      color: #2d3748;
      font-weight: 500;
    }
  }
}

.help-content {
  .help-item {
    margin-bottom: 16px;
    
    &:last-child {
      margin-bottom: 0;
    }
    
    strong {
      color: #2d3748;
      font-size: 14px;
    }
    
    p {
      margin: 4px 0 0 0;
      color: #718096;
      font-size: 13px;
      line-height: 1.5;
    }
  }
}

// 响应式设计
@media (max-width: 1024px) {
  .main-container {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .info-panel {
    position: static;
  }
}

@media (max-width: 768px) {
  .cert-create-page {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
    
    .header-actions {
      width: 100%;
      justify-content: flex-end;
    }
  }
  
  .form-grid {
    grid-template-columns: 1fr;
  }
  
  .verification-tabs {
    .tab-item {
      .tab-content {
        .tab-desc {
          display: none;
        }
      }
    }
  }
}
</style>
