<template>
  <div class="cert-create">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <el-button 
            type="text" 
            @click="goBack" 
            class="back-btn"
          >
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <div class="header-title">
            <h1>申请SSL证书</h1>
            <p>为您的域名申请免费的Let's Encrypt SSL证书</p>
          </div>
        </div>
        <div class="header-actions">
          <el-button @click="goBack">取消</el-button>
          <el-button 
            type="primary" 
            :loading="applying" 
            @click="handleApply"
            :disabled="!canSubmit"
          >
            <el-icon><Lock /></el-icon>
            申请证书
          </el-button>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <el-card class="form-card">
        <template #header>
          <div class="card-header">
            <el-icon class="header-icon"><Document /></el-icon>
            <span>证书配置</span>
          </div>
        </template>

        <el-form 
          ref="formRef" 
          :model="formData" 
          :rules="formRules" 
          label-width="120px"
          size="large"
        >
          <!-- 申请方式选择 -->
          <el-form-item label="申请方式" prop="applyType">
            <el-radio-group v-model="formData.applyType" @change="handleApplyTypeChange">
              <el-radio label="auto" class="apply-type-radio">
                <div class="radio-content">
                  <div class="radio-header">
                    <el-icon class="radio-icon"><Lightning /></el-icon>
                    <span class="radio-title">自动申请</span>
                    <el-tag type="success" size="small">推荐</el-tag>
                  </div>
                  <div class="radio-desc">
                    使用Let's Encrypt免费申请SSL证书，支持自动续期
                  </div>
                </div>
              </el-radio>
              <el-radio label="upload" class="apply-type-radio">
                <div class="radio-content">
                  <div class="radio-header">
                    <el-icon class="radio-icon"><FolderAdd /></el-icon>
                    <span class="radio-title">导入证书</span>
                  </div>
                  <div class="radio-desc">
                    导入已有的SSL证书文件
                  </div>
                </div>
              </el-radio>
              <el-radio label="csr" class="apply-type-radio">
                <div class="radio-content">
                  <div class="radio-header">
                    <el-icon class="radio-icon"><Upload /></el-icon>
                    <span class="radio-title">自定义CSR</span>
                  </div>
                  <div class="radio-desc">
                    上传CSR文件申请证书
                  </div>
                </div>
              </el-radio>
            </el-radio-group>
          </el-form-item>

          <!-- 自动申请配置 -->
          <template v-if="formData.applyType === 'auto'">
            <!-- 域名选择 -->
            <el-form-item label="选择域名" prop="domainId" required>
              <el-select 
                v-model="formData.domainId" 
                placeholder="请选择要申请证书的域名"
                filterable
                @change="handleDomainChange"
                style="width: 100%"
              >
                <el-option 
                  v-for="domain in domainOptions" 
                  :key="domain.id"
                  :label="domain.name" 
                  :value="domain.id"
                >
                  <div class="domain-option">
                    <span class="domain-name">{{ domain.name }}</span>
                    <el-tag 
                      v-if="domain.provider" 
                      type="info" 
                      size="small"
                    >
                      {{ domain.provider.name }}
                    </el-tag>
                  </div>
                </el-option>
              </el-select>
              <div class="form-tip">
                <el-icon><InfoFilled /></el-icon>
                只显示已配置DNS Provider的域名
              </div>
            </el-form-item>

            <!-- 证书类型 -->
            <el-form-item label="证书类型" prop="certType">
              <el-radio-group v-model="formData.certType" @change="handleCertTypeChange">
                <el-radio label="single">
                  <div class="cert-type-option">
                    <span class="option-title">单域名证书</span>
                    <span class="option-desc">仅保护选择的域名</span>
                  </div>
                </el-radio>
                <el-radio label="wildcard">
                  <div class="cert-type-option">
                    <span class="option-title">通配符证书</span>
                    <span class="option-desc">保护域名及其所有子域名</span>
                  </div>
                </el-radio>
                <el-radio label="multi">
                  <div class="cert-type-option">
                    <span class="option-title">多域名证书</span>
                    <span class="option-desc">保护多个不同的域名</span>
                  </div>
                </el-radio>
              </el-radio-group>
            </el-form-item>

            <!-- 域名列表预览 -->
            <el-form-item label="证书域名" v-if="formData.domains.length > 0">
              <div class="domain-preview">
                <el-tag 
                  v-for="domain in formData.domains" 
                  :key="domain"
                  type="primary"
                  class="domain-tag"
                >
                  <el-icon><Globe /></el-icon>
                  {{ domain }}
                </el-tag>
              </div>
            </el-form-item>

            <!-- 多域名证书的域名编辑 -->
            <el-form-item 
              label="域名列表" 
              v-if="formData.certType === 'multi'"
            >
              <div class="multi-domain-editor">
                <el-tag
                  v-for="domain in formData.domains"
                  :key="domain"
                  closable
                  @close="removeDomain(domain)"
                  class="domain-tag"
                >
                  {{ domain }}
                </el-tag>
                <el-input
                  v-if="domainInputVisible"
                  ref="domainInputRef"
                  v-model="domainInputValue"
                  size="small"
                  @keyup.enter="confirmDomainInput"
                  @blur="confirmDomainInput"
                  placeholder="输入域名"
                  class="domain-input"
                />
                <el-button 
                  v-else 
                  size="small" 
                  @click="showDomainInput"
                  class="add-domain-btn"
                >
                  <el-icon><Plus /></el-icon>
                  添加域名
                </el-button>
              </div>
            </el-form-item>

            <!-- 验证方式 -->
            <el-form-item label="验证方式" prop="challengeType">
              <el-radio-group v-model="formData.challengeType">
                <el-radio label="dns-01">
                  <div class="challenge-option">
                    <div class="option-header">
                      <span class="option-title">DNS验证</span>
                      <el-tag type="success" size="small">推荐</el-tag>
                    </div>
                    <div class="option-desc">
                      通过DNS TXT记录验证域名所有权，支持通配符证书
                    </div>
                  </div>
                </el-radio>
                <el-radio label="http-01" :disabled="formData.certType === 'wildcard'">
                  <div class="challenge-option">
                    <div class="option-header">
                      <span class="option-title">HTTP验证</span>
                      <el-tag v-if="formData.certType === 'wildcard'" type="info" size="small">不支持通配符</el-tag>
                    </div>
                    <div class="option-desc">
                      通过HTTP文件验证域名所有权，需要域名可访问
                    </div>
                  </div>
                </el-radio>
              </el-radio-group>
            </el-form-item>

            <!-- 密钥类型 -->
            <el-form-item label="密钥类型" prop="keyType">
              <el-select v-model="formData.keyType" style="width: 200px">
                <el-option label="RSA 2048" value="RSA2048" />
                <el-option label="RSA 4096" value="RSA4096" />
                <el-option label="ECDSA P-256" value="ECDSA256" />
                <el-option label="ECDSA P-384" value="ECDSA384" />
              </el-select>
              <div class="form-tip">
                <el-icon><InfoFilled /></el-icon>
                推荐使用RSA 2048，兼容性最好
              </div>
            </el-form-item>

            <!-- 邮箱地址 -->
            <el-form-item label="邮箱地址" prop="email">
              <el-input 
                v-model="formData.email" 
                placeholder="用于接收证书相关通知"
                style="width: 300px"
              />
              <div class="form-tip">
                <el-icon><InfoFilled /></el-icon>
                用于Let's Encrypt账户注册和重要通知
              </div>
            </el-form-item>
          </template>

          <!-- 高级设置 -->
          <el-divider content-position="left">
            <el-icon><Setting /></el-icon>
            高级设置
          </el-divider>

          <!-- 自动续期 -->
          <el-form-item label="自动续期">
            <el-switch 
              v-model="formData.autoRenew" 
              active-text="开启" 
              inactive-text="关闭"
            />
            <div class="form-tip">
              <el-icon><InfoFilled /></el-icon>
              开启后将在证书到期前30天自动续期
            </div>
          </el-form-item>

          <!-- 部署主机 -->
          <el-form-item label="部署主机">
            <el-select 
              v-model="formData.deployHosts" 
              multiple 
              placeholder="选择要部署证书的主机（可选）"
              style="width: 100%"
              clearable
            >
              <el-option 
                v-for="host in hostOptions" 
                :key="host.id"
                :label="`${host.name} (${host.ip})`" 
                :value="host.id"
              >
                <div class="host-option">
                  <span class="host-name">{{ host.name }}</span>
                  <span class="host-ip">{{ host.ip }}</span>
                  <el-tag 
                    :type="host.status === 'online' ? 'success' : 'danger'" 
                    size="small"
                  >
                    {{ host.status === 'online' ? '在线' : '离线' }}
                  </el-tag>
                </div>
              </el-option>
            </el-select>
            <div class="form-tip">
              <el-icon><InfoFilled /></el-icon>
              证书申请成功后将自动部署到选择的主机
            </div>
          </el-form-item>

          <!-- 备注信息 -->
          <el-form-item label="备注信息">
            <el-input 
              v-model="formData.remark" 
              type="textarea" 
              :rows="3"
              placeholder="请输入证书用途或备注信息（可选）"
              maxlength="200"
              show-word-limit
            />
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, nextTick } from 'vue'
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
  Globe,
  Plus,
  Setting
} from '@element-plus/icons-vue'
import { domainApi } from '@/api/dns/domain'
import { certificateApi } from '@/api/dns/certificate'
import { hostApi } from '@/api/host'

const router = useRouter()

// 表单引用
const formRef = ref()

// 表单数据
const formData = reactive({
  applyType: 'auto', // auto: 自动申请, upload: 导入证书, csr: 自定义CSR
  domainId: null,
  certType: 'single', // single: 单域名, wildcard: 通配符, multi: 多域名
  domains: [] as string[],
  challengeType: 'dns-01', // dns-01: DNS验证, http-01: HTTP验证
  keyType: 'RSA2048', // RSA2048, RSA4096, ECDSA256, ECDSA384
  email: '',
  autoRenew: true,
  deployHosts: [] as number[],
  remark: '',
  providerId: null
})

// 表单验证规则
const formRules = {
  domainId: [
    { required: true, message: '请选择域名', trigger: 'change' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

// 状态数据
const applying = ref(false)
const domainOptions = ref([])
const hostOptions = ref([])

// 多域名编辑相关
const domainInputVisible = ref(false)
const domainInputValue = ref('')
const domainInputRef = ref()

// 计算属性
const canSubmit = computed(() => {
  if (formData.applyType === 'auto') {
    return formData.domainId && formData.email && formData.domains.length > 0
  }
  return false
})

// 生命周期
onMounted(() => {
  loadDomains()
  loadHosts()
})

// 方法
const goBack = () => {
  router.back()
}

const handleApplyTypeChange = (type: string) => {
  // 重置表单数据
  formData.domainId = null
  formData.domains = []
}

const handleDomainChange = (domainId: number) => {
  const domain = domainOptions.value.find((d: any) => d.id === domainId)
  if (domain) {
    formData.providerId = domain.provider_id || domain.providerId
    updateDomainsForCertType(domain.name)
  }
}

const handleCertTypeChange = (certType: string) => {
  if (formData.domainId) {
    const domain = domainOptions.value.find((d: any) => d.id === formData.domainId)
    if (domain) {
      updateDomainsForCertType(domain.name)
    }
  }
}

const updateDomainsForCertType = (domainName: string) => {
  if (formData.certType === 'wildcard') {
    formData.domains = [`*.${domainName}`]
  } else if (formData.certType === 'single') {
    formData.domains = [domainName]
  } else if (formData.certType === 'multi') {
    if (formData.domains.length === 0) {
      formData.domains = [domainName]
    }
  }
}

// 多域名编辑
const removeDomain = (domain: string) => {
  const index = formData.domains.indexOf(domain)
  if (index > -1) {
    formData.domains.splice(index, 1)
  }
}

const showDomainInput = () => {
  domainInputVisible.value = true
  nextTick(() => {
    domainInputRef.value?.focus()
  })
}

const confirmDomainInput = () => {
  if (domainInputValue.value && !formData.domains.includes(domainInputValue.value)) {
    formData.domains.push(domainInputValue.value)
  }
  domainInputVisible.value = false
  domainInputValue.value = ''
}

// 数据加载
const loadDomains = async () => {
  try {
    const response = await domainApi.list({ page: 1, size: 1000 })
    domainOptions.value = response.data.items || []
  } catch (error) {
    console.error('加载域名列表失败:', error)
  }
}

const loadHosts = async () => {
  try {
    const response = await hostApi.list({ page: 1, size: 1000 })
    hostOptions.value = response.data.items || []
  } catch (error) {
    console.error('加载主机列表失败:', error)
  }
}

// 证书申请
const handleApply = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    applying.value = true

    // 自动申请证书
    const response = await certificateApi.apply({
      domains: formData.domains,
      provider_id: formData.providerId,
      key_type: formData.keyType,
      email: formData.email,
      auto_renew: formData.autoRenew,
      deploy_hosts: formData.deployHosts,
      remark: formData.remark,
      challenge_type: formData.challengeType
    })

    ElMessage.success('证书申请成功')
    router.push('/dns/certs')
  } catch (error: any) {
    console.error('证书申请失败:', error)
    ElMessage.error(error.message || '证书申请失败')
  } finally {
    applying.value = false
  }
}
