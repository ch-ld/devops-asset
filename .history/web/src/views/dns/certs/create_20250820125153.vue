<!-- 这个文件已被临时禁用，请使用 import-cert.vue -->
<template>
  <div>
    <h1>此页面已被禁用</h1>
    <p>请使用新的证书申请页面</p>
  </div>
</template>

<script setup lang="ts">
// 临时禁用此文件
</script>

<style scoped>
/* 临时禁用样式 */
</style>


            <!-- 证书类型 -->
            <el-form-item label="证书类型" prop="certType">
              <div class="cert-type-cards">
                <div 
                  class="cert-type-card"
                  :class="{ active: formData.certType === 'single' }"
                  @click="selectCertType('single')"
                >
                  <div class="type-icon">
                    <el-icon><Document /></el-icon>
                  </div>
                  <div class="type-content">
                    <div class="type-title">单域名证书</div>
                    <div class="type-desc">仅保护选择的域名</div>
                    <div class="type-example">例如：example.com</div>
                  </div>
                  <div class="type-radio">
                    <el-radio v-model="formData.certType" label="single"></el-radio>
                  </div>
                </div>

                <div 
                  class="cert-type-card"
                  :class="{ active: formData.certType === 'wildcard' }"
                  @click="selectCertType('wildcard')"
                >
                  <div class="type-icon">
                    <el-icon><Connection /></el-icon>
                  </div>
                  <div class="type-content">
                    <div class="type-title">通配符证书</div>
                    <div class="type-desc">保护域名及其所有子域名</div>
                    <div class="type-example">例如：*.example.com</div>
                  </div>
                  <div class="type-radio">
                    <el-radio v-model="formData.certType" label="wildcard"></el-radio>
                  </div>
                </div>

                <div 
                  class="cert-type-card"
                  :class="{ active: formData.certType === 'multi' }"
                  @click="selectCertType('multi')"
                >
                  <div class="type-icon">
                    <el-icon><Plus /></el-icon>
                  </div>
                  <div class="type-content">
                    <div class="type-title">多域名证书</div>
                    <div class="type-desc">保护多个不同的域名</div>
                    <div class="type-example">例如：多个域名</div>
                  </div>
                  <div class="type-radio">
                    <el-radio v-model="formData.certType" label="multi"></el-radio>
                  </div>
                </div>
              </div>
            </el-form-item>

            <!-- 邮箱地址 -->
            <el-form-item label="邮箱地址" prop="email" required>
              <el-input 
                v-model="formData.email" 
                placeholder="请输入邮箱地址，用于接收证书相关通知"
                size="large"
              >
                <template #prefix>
                  <el-icon><Message /></el-icon>
                </template>
              </el-input>
              <div class="form-tip">
                <el-icon><InfoFilled /></el-icon>
                邮箱用于接收证书到期提醒和重要通知
              </div>
            </el-form-item>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
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
const route = useRoute()

// 表单引用
const formRef = ref()

// 当前步骤
const currentStep = ref(0)

// 表单数据
const formData = reactive({
  applyType: 'auto', // auto: 自动申请, upload: 导入证书, csr: 自定义CSR
  domainId: null,
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

// 响应式数据
const applying = ref(false)
const domainOptions = ref([])
const customDomainInput = ref('')
const domainExamples = ref([
  'example.com',
  '*.example.com',
  'www.example.com',
  'api.example.com'
])

// 计算属性
const canSubmit = computed(() => {
  if (formData.applyType === 'auto') {
    return formData.domainId && formData.email && formData.certType
  }
  return false
})

// 生命周期
onMounted(() => {
  // 检查URL参数，设置默认申请类型
  const type = route.query.type as string
  if (type && ['auto', 'upload', 'csr'].includes(type)) {
    formData.applyType = type
  }

  loadDomains()
})

// 方法
const goBack = () => {
  router.back()
}

// 选择申请方式
const selectApplyType = (type: string) => {
  formData.applyType = type
  currentStep.value = 1
}

// 选择证书类型
const selectCertType = (type: string) => {
  formData.certType = type
}

// 处理域名变化
const handleDomainChange = (domainId: number) => {
  const domain = domainOptions.value.find(d => d.id === domainId)
  if (domain) {
    formData.domains = [domain.name]
    // 清空自定义输入
    customDomainInput.value = ''
  }
}

// 验证域名格式
const validateDomain = (domain: string) => {
  const domainRegex = /^(\*\.)?[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/
  return domainRegex.test(domain)
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

  // 清空域名选择
  formData.domainId = null
  formData.domains = [domain]

  // 根据域名类型自动设置证书类型
  if (domain.startsWith('*.')) {
    formData.certType = 'wildcard'
  } else {
    formData.certType = 'single'
  }

  ElMessage.success(`已设置域名：${domain}`)
}

// 清除域名选择
const clearDomainSelection = () => {
  formData.domainId = null
  formData.domains = []
  customDomainInput.value = ''
  ElMessage.info('已清除域名选择，请重新选择')
}

// 使用示例域名
const useExampleDomain = (example: string) => {
  customDomainInput.value = example
  addCustomDomain()
}

// 数据加载
const loadDomains = async () => {
  try {
    const response = await domainApi.list({
      page: 1,
      page_size: 1000,
      status: 'active' // 只获取活跃的域名
    })

    // 只显示已配置DNS Provider的域名
    domainOptions.value = (response.data?.items || []).filter(domain =>
      domain.provider_id && domain.provider
    )
  } catch (error) {
    console.error('加载域名列表失败:', error)
    ElMessage.error('加载域名列表失败')
  }
}

// 证书申请
const handleApply = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    applying.value = true

    // 获取选中的域名信息
    const selectedDomain = domainOptions.value.find(d => d.id === formData.domainId)
    if (!selectedDomain) {
      throw new Error('请选择有效的域名')
    }

    // 构建申请数据
    const requestData = {
      domain_id: formData.domainId,
      domains: formData.certType === 'wildcard'
        ? [`*.${selectedDomain.name}`, selectedDomain.name] // 通配符证书包含主域名
        : [selectedDomain.name],
      email: formData.email,
      provider_id: selectedDomain.provider_id!,
      key_type: 'RSA2048',
      valid_days: 90,
      auto_renew: true,
      remark: formData.remark || `${formData.certType === 'wildcard' ? '通配符' : '单域名'}SSL证书`
    }

    // 使用更长的超时时间进行证书申请
    await Promise.race([
      certificateApi.create(requestData),
      new Promise((_, reject) =>
        setTimeout(() => reject(new Error('证书申请超时，但后台仍在处理中，请稍后查看证书列表')), 60000)
      )
    ])

    ElMessage.success('证书申请成功，正在后台处理中')
    router.push('/dns/certs')
  } catch (error: any) {
    console.error('证书申请失败:', error)
    if (error.message?.includes('超时')) {
      ElMessage.warning('证书申请可能需要较长时间，请稍后查看证书列表')
      // 即使超时也跳转到列表页面，因为后台可能仍在处理
      setTimeout(() => {
        router.push('/dns/certs')
      }, 2000)
    } else {
      ElMessage.error(error.message || '证书申请失败')
    }
  } finally {
    applying.value = false
  }
}
</script>

<style scoped>
/* 基础布局 */
.cert-create {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.page-header {
  background: white;
  border-bottom: 1px solid #e4e7ed;
  padding: 20px 24px;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 1200px;
  margin: 0 auto;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.back-btn {
  color: #606266;
  font-size: 16px;
  padding: 8px 16px;
}

.back-btn:hover {
  color: #409eff;
  background-color: #f0f9ff;
}

.header-title h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.header-title p {
  margin: 6px 0 0 0;
  font-size: 14px;
  color: #909399;
}

.header-actions {
  display: flex;
  gap: 16px;
}

.main-content {
  max-width: 1000px;
  margin: 32px auto;
  padding: 0 24px;
}

.steps-container {
  background: white;
  border-radius: 12px;
  padding: 32px;
  margin-bottom: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.form-container {
  background: white;
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.cert-form {
  max-width: none;
}

.step-section {
  margin-bottom: 40px;
}

.section-header {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 2px solid #f0f2f5;
}

.section-header h3 {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-header p {
  margin: 0;
  font-size: 14px;
  color: #909399;
}

/* 申请方式卡片样式 */
.apply-type-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 16px;
}

.apply-type-card {
  position: relative;
  border: 2px solid #e4e7ed;
  border-radius: 12px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  background: white;
}

.apply-type-card:hover {
  border-color: #409eff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
  transform: translateY(-2px);
}

.apply-type-card.active {
  border-color: #409eff;
  background: linear-gradient(135deg, #f0f9ff 0%, #e6f7ff 100%);
  box-shadow: 0 4px 20px rgba(64, 158, 255, 0.2);
}

.card-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: linear-gradient(135deg, #409eff 0%, #67c23a 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}

.card-icon .el-icon {
  font-size: 24px;
  color: white;
}

.card-content {
  flex: 1;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-desc {
  font-size: 14px;
  color: #606266;
  margin-bottom: 12px;
  line-height: 1.5;
}

.card-features {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.card-features span {
  font-size: 12px;
  color: #909399;
}

.card-radio {
  position: absolute;
  top: 16px;
  right: 16px;
}

.card-radio .el-radio {
  margin: 0;
}

/* 证书类型卡片样式 */
.cert-type-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.cert-type-card {
  position: relative;
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  background: white;
  display: flex;
  align-items: center;
  gap: 16px;
}

.cert-type-card:hover {
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.15);
}

.cert-type-card.active {
  border-color: #409eff;
  background: #f0f9ff;
}

.type-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: #409eff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.type-icon .el-icon {
  font-size: 20px;
  color: white;
}

.type-content {
  flex: 1;
}

.type-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
}

.type-desc {
  font-size: 14px;
  color: #606266;
  margin-bottom: 4px;
}

.type-example {
  font-size: 12px;
  color: #909399;
}

.type-radio {
  position: absolute;
  top: 12px;
  right: 12px;
}

// 自定义域名输入样式
.custom-domain-section {
  margin-bottom: 12px;

  .domain-examples {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 8px;
    flex-wrap: wrap;

    .example-label {
      color: #666;
      font-size: 14px;
      white-space: nowrap;
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

/* 域名选择优化 */
.domain-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 8px 0;
}

.domain-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.domain-name {
  font-weight: 500;
  color: #303133;
}

.domain-status {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #67c23a;
}

.form-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 8px;
  padding: 8px 12px;
  background: #f0f9ff;
  border-radius: 6px;
  font-size: 13px;
  color: #606266;
}

.form-tip .el-icon {
  color: #409eff;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .main-content {
    margin: 16px auto;
    padding: 0 16px;
  }

  .steps-container,
  .form-container {
    padding: 20px;
  }

  .apply-type-cards {
    grid-template-columns: 1fr;
  }

  .cert-type-cards {
    grid-template-columns: 1fr;
  }

  .header-content {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }

  .header-actions {
    width: 100%;
    justify-content: flex-end;
  }
}

/* 已选择域名显示样式 */
.selected-domain-display {
  margin-bottom: 12px;

  .el-tag {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 8px 12px;
    font-size: 14px;
    font-weight: 500;
  }

  .domain-tip {
    margin-left: 8px;
    font-size: 12px;
    color: #909399;
  }
}
</style>
