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
            <el-icon><Connection /></el-icon>
            <span>证书域名</span>
            <span class="required">*</span>
          </div>
          <div class="section-content">
            <el-form-item prop="domainIds">
              <el-select
                v-model="formData.domainIds"
                multiple
                placeholder="请选择要申请证书的域名"
                size="large"
                style="width: 100%"
                @change="handleDomainChange"
              >
                <el-option
                  v-for="domain in availableDomains"
                  :key="domain.id"
                  :label="domain.name"
                  :value="domain.id"
                >
                  <div class="domain-option">
                    <span class="domain-name">{{ domain.name }}</span>
                    <el-tag
                      :type="domain.status === 'active' ? 'success' : 'warning'"
                      size="small"
                    >
                      {{ domain.status === 'active' ? '正常' : '异常' }}
                    </el-tag>
                  </div>
                </el-option>
              </el-select>
              <div class="domain-help">
                <el-text type="info" size="small">
                  支持选择多个域名，系统将为选中的域名申请SAN证书
                </el-text>
              </div>
              <div v-if="formData.domains.length > 0" class="selected-domains">
                <el-tag
                  v-for="domain in formData.domains"
                  :key="domain"
                  closable
                  @close="removeDomainByName(domain)"
                  class="domain-tag"
                >
                  {{ domain }}
                </el-tag>
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

        <!-- 核心配置（两列网格） -->
        <div class="form-grid">
          <!-- 左：域名验证方式 -->
          <div class="grid-item">
            <el-form-item label="域名验证方式" prop="challengeType">
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

          <!-- 右：证书颁发机构 -->
          <div class="grid-item">
            <el-form-item label="证书颁发机构" prop="caType">
              <el-select v-model="formData.caType" placeholder="请选择证书提供商" size="large" class="ca-select">
                <el-option
                  v-for="provider in caProviders"
                  :key="provider.type"
                  :label="provider.name"
                  :value="provider.type"
                >
                  <div class="provider-option">
                    <div class="provider-name">{{ provider.name }}</div>
                    <div class="provider-desc">{{ provider.description }}</div>
                    <el-tag v-if="provider.free" type="success" size="small">免费</el-tag>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
            <ul class="field-tips">
              <li>Let's Encrypt：推荐，免费</li>
              <li>ZeroSSL：需EAB授权</li>
            </ul>
          </div>

          <!-- 左：DNS解析服务商 -->
          <div class="grid-item">
            <el-form-item label="DNS解析服务商" prop="providerId">
              <el-select
                v-model="formData.providerId"
                placeholder="请选择DNS解析服务商"
                filterable
                clearable
                :teleported="false"
                size="large"
                class="provider-select"
              >
                <el-option
                  v-for="provider in dnsProviders"
                  :key="Number(provider.id)"
                  :value="Number(provider.id)"
                  :label="provider.name"
                >
                  <div class="provider-option">
                    <ProviderIcon :type="provider.type" size="24px" />
                    <div class="provider-info">
                      <div class="provider-name">{{ provider.name }}</div>
                      <div class="provider-type">{{ getProviderTypeName(provider.type) }}</div>
                    </div>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </div>

          <!-- 右：DNS解析授权（占位操作） -->
          <div class="grid-item">
            <el-form-item label="DNS解析授权">
              <el-space>
                <el-tag type="warning">未授权</el-tag>
                <el-button size="small" @click="openDnsAuthDialog">选择</el-button>
              </el-space>
              <div class="form-tip">选择对应云厂商的授权以自动添加验证记录</div>
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

        <!-- 自动续期配置 -->
        <div class="form-section">
          <div class="section-header">
            <el-icon><Clock /></el-icon>
            <span>自动续期</span>
          </div>
          <div class="section-content">
            <el-form-item>
              <el-switch
                v-model="formData.autoRenew"
                active-text="启用自动续期"
                inactive-text="禁用自动续期"
                size="large"
              />
              <div class="form-tip">
                启用后，证书将在到期前自动续期
              </div>
            </el-form-item>
          </div>
        </div>

        <!-- 部署配置 -->
        <div class="form-section wide">
          <div class="section-header">
            <el-icon><Bell /></el-icon>
            <span>部署配置</span>
          </div>
          <div class="section-content">
            <el-form-item>
              <el-switch
                v-model="formData.deployEnabled"
                active-text="启用自动部署"
                inactive-text="仅申请证书"
                size="large"
              />
            </el-form-item>

            <div v-if="formData.deployEnabled" class="deploy-config">
              <el-form-item label="部署主机">
                <el-select
                  v-model="formData.deployHosts"
                  multiple
                  placeholder="请选择要部署证书的主机"
                  size="large"
                  style="width: 100%"
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
                  placeholder="证书部署路径，如：/etc/ssl/certs/"
                  size="large"
                />
              </el-form-item>

              <el-form-item label="重启命令">
                <el-input
                  v-model="formData.restartCommand"
                  placeholder="部署后执行的重启命令，如：systemctl reload nginx"
                  size="large"
                />
              </el-form-item>
            </div>
          </div>
        </div>

        <!-- 备注 -->
        <div class="form-section">
          <div class="section-header">
            <el-icon><Document /></el-icon>
            <span>备注信息</span>
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
  Connection,
  Message,
  Star,
  Key,
  Calendar,
  Clock,
  Bell,
  Document,
  Close
} from '@element-plus/icons-vue'
import { domainApi } from '@/api/dns/domain'
import { dnsProviderApi } from '@/api/dns/provider'
import { certificateApi } from '@/api/dns/certificate'
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
    { required: true, message: '请至少选择一个域名', trigger: 'change' }
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
    { required: true, message: '请选择DNS解析服务商', trigger: 'change' }
  ]
}

// 状态管理
const applying = ref(false)
const domainInput = ref('')
const dnsProviders = ref<any[]>([])
const availableDomains = ref<any[]>([])
const availableHosts = ref<any[]>([])
const caProviders = ref([
  {
    type: 'letsencrypt',
    name: "Let's Encrypt",
    description: '免费、自动化的证书颁发机构',
    free: true,
    icon: 'Lock'
  },
  {
    type: 'zerossl',
    name: 'ZeroSSL',
    description: '免费SSL证书提供商',
    free: true,
    icon: 'Lock'
  },
  {
    type: 'buypass',
    name: 'Buypass',
    description: '挪威免费证书颁发机构',
    free: true,
    icon: 'Lock'
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
  // 简单的域名验证
  const domainRegex = /^(\*\.)?[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/
  return domainRegex.test(domain)
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
    const response = await dnsProviderApi.list({
      page: 1,
      size: 100
    })
    const data: any = (response as any).data || response
    dnsProviders.value = data?.items || []
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
    // 这里需要导入主机API，暂时使用空数组
    availableHosts.value = []
  } catch (error) {
    console.error('加载主机列表失败:', error)
  }
}

const handleApply = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    applying.value = true

    // 移除时间处理，不再需要

const openDnsAuthDialog = () => {
  ElMessage.info('DNS授权功能占位：后续接入各云厂商授权流程')
}


    const requestData = {
      domains: formData.domains,
      email: formData.email,
      ca_type: formData.caType,
      challenge_type: formData.challengeType,
      provider_id: formData.providerId,
      key_type: formData.keyType,
      auto_renew: formData.autoRenew,
      renew_days: formData.renewDays,
      deploy_hosts: formData.deployEnabled ? formData.deployHosts : [],
      deploy_path: formData.deployEnabled ? formData.deployPath : '',
      restart_command: formData.deployEnabled ? formData.restartCommand : '',
      remark: formData.remark,
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
  loadDomains()
  loadHosts()
})
</script>

<style scoped lang="scss">
.cert-create-advanced {
  min-height: 100vh;
  background: #f5f7fa;
  padding: 20px;

  .page-header {
    background: white;
    border-radius: 16px;
    padding: 32px;
    margin-bottom: 32px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.8);

    .header-content {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .header-left {
        display: flex;
        align-items: center;
        gap: 20px;

        .back-btn {
          color: #666;
          font-weight: 500;
          transition: all 0.3s ease;

          &:hover {
            color: #409eff;
            transform: translateX(-2px);
          }
        }

        .header-title {
          h1 {
            margin: 0;
            font-size: 28px;
            font-weight: 700;
            color: #1a1a1a;
            display: flex;
            align-items: center;
            gap: 12px;

            .title-icon {
              font-size: 32px;
              color: #409eff;
            }
          }

          p {
            margin: 8px 0 0 0;
            color: #666;
            font-size: 16px;
            font-weight: 400;
          }
        }
      }

      .header-actions {
        display: flex;
        gap: 16px;

        .cancel-btn {
          border: 2px solid #e4e7ed;
          color: #606266;
          font-weight: 500;
          transition: all 0.3s ease;

          &:hover {
            border-color: #409eff;
            color: #409eff;
          }
        }

        .apply-btn {
          background: linear-gradient(135deg, #409eff, #67c23a);
          border: none;
          font-weight: 600;
          box-shadow: 0 4px 16px rgba(64, 158, 255, 0.3);
          transition: all 0.3s ease;

          &:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 24px rgba(64, 158, 255, 0.4);
          }
        }
      }
    }
  }

  .main-content {

    /* 双列表单布局 */
    .form-grid{display:grid;grid-template-columns:1fr 1fr;gap:16px 24px;margin-bottom:24px}
    .grid-item{background:#fff;border:1px solid #ebeef5;border-radius:8px;padding:16px}
    .field-tips{margin:8px 0 0 0;padding-left:18px;color:#909399;font-size:12px;line-height:1.5}

    max-width: 1000px;
    margin: 0 auto;

    .cert-form {
      .form-section {
        background: white;
        border-radius: 16px;
        padding: 32px;
        margin-bottom: 24px;
        box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
        border: 1px solid rgba(255, 255, 255, 0.8);
        transition: all 0.3s ease;

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
        }

        .section-header {
          display: flex;
          align-items: center;
          gap: 12px;
          margin-bottom: 24px;
          font-size: 18px;
          font-weight: 600;
          color: #1a1a1a;

          .el-icon {
            font-size: 20px;
            color: #409eff;
          }

          .required {
            color: #f56c6c;
            margin-left: 4px;
          }
        }

        .section-content {
          .el-form-item {
            margin-bottom: 0;
          }
        }
      }
    }
  }

  // 域名输入相关样式
  .domain-input-container {
    .domain-input {
      margin-bottom: 16px;
    }

    .domain-examples {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 16px;

      .example-label {
        color: #666;
        font-size: 14px;
      }

      .example-tag {
        cursor: pointer;
        transition: all 0.3s ease;

        &:hover {
          transform: scale(1.05);
          box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
        }
      }
    }
  }

  .domain-list {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;

    .domain-item {
      display: flex;
      align-items: center;
      gap: 8px;
      background: linear-gradient(135deg, #f0f9ff, #e0f2fe);
      border: 1px solid #bae6fd;
      border-radius: 8px;
      padding: 8px 12px;
      font-size: 14px;

      .domain-name {
        color: #0369a1;
        font-weight: 500;
      }

      .remove-btn {
        color: #ef4444;
        padding: 0;
        min-width: auto;
        width: 16px;
        height: 16px;

        &:hover {
          background: rgba(239, 68, 68, 0.1);
        }
      }
    }
  }

  // 提供商信息样式
  .provider-info {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px;
    background: linear-gradient(135deg, #f0f9ff, #e0f2fe);
    border: 1px solid #bae6fd;
    border-radius: 12px;

    .provider-logo {
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .provider-details {
      flex: 1;

      .provider-name {
        font-size: 16px;
        font-weight: 600;
        color: #0369a1;
        margin-bottom: 4px;
      }

      .provider-desc {
        font-size: 14px;
        color: #0284c7;
      }
    }
  }

  // 选择器样式
  .challenge-select,
  .provider-select,
  .key-type-select {
    width: 100%;

    .option-content {
      .option-main {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 4px;

        .option-label {
          font-weight: 500;
          color: #303133;
        }
      }

      .option-desc {
        font-size: 12px;
        color: #909399;
      }
    }
  }

  .provider-option {
    display: flex;
    align-items: center;
    gap: 12px;

    .provider-info {
      .provider-name {
        font-weight: 500;
        color: #303133;
        margin-bottom: 2px;
      }

      .provider-type {
        font-size: 12px;
        color: #909399;
      }
    }
  }

  // 其他组件样式
  .renew-days-input {
    width: 200px;
  }

  .form-tip {
    margin-top: 8px;
    font-size: 12px;
    color: #909399;
    line-height: 1.4;
  }

  .schedule-config {
    .schedule-item {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 12px;

      .schedule-label {
        font-size: 16px;
        color: #303133;
      }
    }

    .schedule-note {
      font-size: 12px;
      color: #909399;
    }
  }

  .notification-config {
    .custom-notification {
      margin-top: 16px;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .cert-create-advanced {
    padding: 16px;

    .page-header {
      padding: 24px;

      .header-content {
        flex-direction: column;
        gap: 20px;
        align-items: flex-start;

        .header-left {
          gap: 16px;

          .header-title h1 {
            font-size: 24px;
          }
        }

        .header-actions {
          width: 100%;
          justify-content: flex-end;
        }
      }
    }

    .main-content {
      .cert-form {
        .form-section {
          padding: 24px;
        }
      }
    }

    .domain-list {
      .domain-item {
        width: 100%;
        justify-content: space-between;
      }
    }

    .provider-info {
      flex-direction: column;
      text-align: center;
    }

    .schedule-config {
      .schedule-item {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;
      }
    }
  }

  // 新增样式
  .domain-option {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;

    .domain-name {
      font-weight: 500;
    }
  }

  .domain-help {
    margin-top: 8px;
  }

  .selected-domains {
    margin-top: 16px;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;

    .domain-tag {
      font-size: 14px;
    }
  }

  .ca-provider-group {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 12px;

    .ca-provider-option {
      width: 100%;
      margin: 0 !important;
      margin-right: 0 !important;
      position: relative;
      z-index: 1;

      .provider-card {
        width: 100%;
        border: 1px solid #e5e7eb;
        border-radius: 12px;
        padding: 16px;
        transition: all 0.3s ease;
        background: #fafafa;
        cursor: pointer;
        position: relative;

        &:hover {
          border-color: #409eff;
          box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
        }

        .provider-info {
          display: flex;
          align-items: center;
          gap: 12px;

          .provider-logo {
            width: 40px;
            height: 40px;
            display: flex;
            align-items: center;
            justify-content: center;
            background: #f0f9ff;
            border-radius: 8px;
            color: #409eff;
          }

          .provider-details {
            flex: 1;

            .provider-name {
              font-size: 16px;
              font-weight: 600;
              color: #1f2937;
              margin-bottom: 4px;
            }

            .provider-desc {
              font-size: 14px;
              color: #6b7280;
            }
          }
        }
      }

      &.is-checked .provider-card {
        border-color: #409eff;
        background: linear-gradient(135deg, #f0f9ff, #e0f2fe);
      }
    }
  }

  .deploy-config {
    margin-top: 16px;
    padding: 16px;
    background: #f8fafc;
    border-radius: 8px;
    border: 1px solid #e2e8f0;
  }
}
</style>
