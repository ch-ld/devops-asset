<template>
  <div class="dns-record-create">
    <!-- 面包屑导航 -->
    <div class="breadcrumb-nav">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item>DNS管理</el-breadcrumb-item>
        <el-breadcrumb-item>解析记录</el-breadcrumb-item>
        <el-breadcrumb-item>添加记录</el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <div class="icon-wrapper">
            <el-icon size="32"><Plus /></el-icon>
          </div>
          <div class="title-content">
            <h1>添加DNS解析记录</h1>
            <p>为您的域名添加A、CNAME、MX等类型的DNS解析记录</p>
          </div>
        </div>
        <div class="action-section">
          <el-button @click="handleCancel">
            <el-icon><ArrowLeft /></el-icon>
            返回列表
          </el-button>
        </div>
      </div>
    </div>

    <!-- 创建表单 -->
    <div class="form-container">
      <el-card class="form-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span class="card-title">
              <el-icon><Edit /></el-icon>
              记录信息
            </span>
          </div>
        </template>

        <el-form
          ref="formRef"
          :model="formData"
          :rules="formRules"
          label-width="140px"
          class="record-form"
          @submit.prevent
        >
          <!-- 域名选择 -->
          <el-form-item label="目标域名" prop="domain_id" class="form-item-highlight">
            <el-select
              v-model="formData.domain_id"
              placeholder="请选择要添加记录的域名"
              style="width: 100%"
              filterable
              size="large"
            >
              <el-option
                v-for="domain in domains"
                :key="domain.id"
                :label="domain.name"
                :value="domain.id"
              >
                <div class="domain-option">
                  <span class="domain-name">{{ domain.name }}</span>
                  <el-tag
                    v-if="domain.status"
                    :type="domain.status === 'active' ? 'success' : 'danger'"
                    size="small"
                  >
                    {{ domain.status === 'active' ? '正常' : '异常' }}
                  </el-tag>
                </div>
              </el-option>
            </el-select>
          </el-form-item>

          <el-divider content-position="left">
            <span class="divider-text">记录配置</span>
          </el-divider>

          <!-- 记录名和类型 -->
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="记录名" prop="name">
                <el-input
                  v-model="formData.name"
                  placeholder="如：www、mail、@（根域名）"
                  size="large"
                >
                  <template #prefix>
                    <el-icon><Document /></el-icon>
                  </template>
                </el-input>
                <div class="form-tip">
                  <el-icon><InfoFilled /></el-icon>
                  使用 @ 表示根域名，www 表示 www.域名.com
                </div>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="记录类型" prop="type">
                <el-select
                  v-model="formData.type"
                  placeholder="请选择记录类型"
                  style="width: 100%"
                  size="large"
                  @change="handleTypeChange"
                >
                  <el-option-group label="常用类型">
                    <el-option label="A - IPv4地址" value="A">
                      <div class="type-option">
                        <div class="type-main">
                          <span class="type-name">A</span>
                          <span class="type-desc">IPv4地址记录</span>
                        </div>
                        <span class="type-example">如：192.168.1.1</span>
                      </div>
                    </el-option>
                    <el-option label="CNAME - 别名" value="CNAME">
                      <div class="type-option">
                        <div class="type-main">
                          <span class="type-name">CNAME</span>
                          <span class="type-desc">别名记录</span>
                        </div>
                        <span class="type-example">如：www.example.com</span>
                      </div>
                    </el-option>
                    <el-option label="MX - 邮件" value="MX">
                      <div class="type-option">
                        <div class="type-main">
                          <span class="type-name">MX</span>
                          <span class="type-desc">邮件交换记录</span>
                        </div>
                        <span class="type-example">如：mail.example.com</span>
                      </div>
                    </el-option>
                  </el-option-group>
                  <el-option-group label="其他类型">
                    <el-option label="AAAA - IPv6地址" value="AAAA" />
                    <el-option label="TXT - 文本记录" value="TXT" />
                    <el-option label="NS - 域名服务器" value="NS" />
                    <el-option label="SRV - 服务记录" value="SRV" />
                  </el-option-group>
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <!-- 记录值 -->
          <el-form-item label="记录值" prop="value">
            <el-input
              v-model="formData.value"
              :placeholder="getValuePlaceholder()"
              :type="formData.type === 'TXT' ? 'textarea' : 'text'"
              :rows="formData.type === 'TXT' ? 4 : 1"
              size="large"
            >
              <template v-if="formData.type !== 'TXT'" #prefix>
                <el-icon><Link /></el-icon>
              </template>
            </el-input>
            <div class="form-tip">
              <el-icon><InfoFilled /></el-icon>
              {{ getValueTip() }}
            </div>
          </el-form-item>

          <el-divider content-position="left">
            <span class="divider-text">高级配置</span>
          </el-divider>

          <!-- TTL和优先级 -->
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="TTL（生存时间）" prop="ttl">
                <el-select
                  v-model="formData.ttl"
                  placeholder="请选择TTL"
                  style="width: 100%"
                  size="large"
                >
                  <el-option label="1分钟 (60秒)" :value="60" />
                  <el-option label="5分钟 (300秒)" :value="300" />
                  <el-option label="10分钟 (600秒)" :value="600" />
                  <el-option label="30分钟 (1800秒)" :value="1800" />
                  <el-option label="1小时 (3600秒)" :value="3600" />
                  <el-option label="12小时 (43200秒)" :value="43200" />
                  <el-option label="1天 (86400秒)" :value="86400" />
                </el-select>
                <div class="form-tip">
                  <el-icon><InfoFilled /></el-icon>
                  TTL越小解析生效越快，但查询频率会增加
                </div>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item
                v-if="formData.type === 'MX' || formData.type === 'SRV'"
                label="优先级"
                prop="priority"
              >
                <el-input-number
                  v-model="formData.priority"
                  :min="0"
                  :max="65535"
                  size="large"
                  style="width: 100%"
                />
                <div class="form-tip">
                  <el-icon><InfoFilled /></el-icon>
                  数值越小优先级越高
                </div>
              </el-form-item>
            </el-col>
          </el-row>

          <!-- SRV记录专用字段 -->
          <el-row v-if="formData.type === 'SRV'" :gutter="24">
            <el-col :span="12">
              <el-form-item label="权重" prop="weight">
                <el-input-number
                  v-model="formData.weight"
                  :min="0"
                  :max="65535"
                  size="large"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="端口" prop="port">
                <el-input-number
                  v-model="formData.port"
                  :min="1"
                  :max="65535"
                  size="large"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <!-- 备注 -->
          <el-form-item label="备注信息" prop="remark">
            <el-input
              v-model="formData.remark"
              placeholder="请输入备注信息（可选）"
              type="textarea"
              :rows="3"
            />
          </el-form-item>

          <!-- 操作按钮 -->
          <el-form-item class="form-actions">
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              @click="handleSubmit"
              class="submit-btn"
            >
              <el-icon><Check /></el-icon>
              创建记录
            </el-button>
            <el-button size="large" @click="handleReset" class="reset-btn">
              <el-icon><Refresh /></el-icon>
              重置表单
            </el-button>
            <el-button size="large" @click="handleCancel" class="cancel-btn">
              <el-icon><Close /></el-icon>
              取消
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Plus,
  ArrowLeft,
  Edit,
  Document,
  Link,
  InfoFilled,
  Check,
  Refresh,
  Close
} from '@element-plus/icons-vue'
import { recordApi } from '@/api/dns/record'
import { domainApi } from '@/api/dns/domain'
import type { Domain } from '@/types/dns'

const router = useRouter()
const formRef = ref()
const loading = ref(false)
const domains = ref<Domain[]>([])

const formData = reactive({
  domain_id: 0,
  name: '',
  type: 'A',
  value: '',
  ttl: 600,
  priority: 10,
  weight: 10,
  port: 80,
  remark: ''
})

const formRules = {
  domain_id: [
    { required: true, message: '请选择域名', trigger: 'change' }
  ],
  name: [
    { required: true, message: '请输入记录名', trigger: 'blur' },
    { max: 255, message: '记录名长度不能超过255个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择记录类型', trigger: 'change' }
  ],
  value: [
    { required: true, message: '请输入记录值', trigger: 'blur' }
  ],
  ttl: [
    { required: true, message: '请选择TTL', trigger: 'change' }
  ]
}

function getValuePlaceholder() {
  switch (formData.type) {
    case 'A':
      return '请输入IPv4地址，如：192.168.1.1'
    case 'AAAA':
      return '请输入IPv6地址，如：2001:db8::1'
    case 'CNAME':
      return '请输入目标域名，如：www.example.com'
    case 'MX':
      return '请输入邮件服务器，如：mail.example.com'
    case 'TXT':
      return '请输入文本内容，如：v=spf1 include:_spf.example.com ~all'
    case 'NS':
      return '请输入域名服务器，如：ns1.example.com'
    case 'SRV':
      return '请输入目标主机，如：target.example.com'
    default:
      return '请输入记录值'
  }
}

function getValueTip() {
  switch (formData.type) {
    case 'A':
      return '填写服务器的IPv4地址'
    case 'AAAA':
      return '填写服务器的IPv6地址'
    case 'CNAME':
      return '填写指向的目标域名，不能指向IP地址'
    case 'MX':
      return '填写邮件服务器的域名'
    case 'TXT':
      return '填写文本信息，常用于SPF、DKIM等验证'
    case 'NS':
      return '填写域名服务器的域名'
    case 'SRV':
      return '填写提供服务的目标主机'
    default:
      return '请根据记录类型填写相应的值'
  }
}

function handleTypeChange() {
  // 根据记录类型调整默认值
  switch (formData.type) {
    case 'MX':
      if (!formData.priority) formData.priority = 10
      break
    case 'SRV':
      if (!formData.priority) formData.priority = 10
      if (!formData.weight) formData.weight = 10
      if (!formData.port) formData.port = 80
      break
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    loading.value = true
    
    await recordApi.create(formData)
    ElMessage.success('DNS记录创建成功')
    
    // 返回列表页
    router.push('/dns/records')
  } catch (error) {
    console.error('创建DNS记录失败:', error)
    ElMessage.error('DNS记录创建失败')
  } finally {
    loading.value = false
  }
}

function handleReset() {
  Object.assign(formData, {
    domain_id: 0,
    name: '',
    type: 'A',
    value: '',
    ttl: 600,
    priority: 10,
    weight: 10,
    port: 80,
    remark: ''
  })
  
  formRef.value?.clearValidate()
}

function handleCancel() {
  router.push('/dns/records')
}

async function fetchDomains() {
  try {
    const response = await domainApi.list({ page: 1, page_size: 100 })
    const data = (response as any).data || response
    domains.value = data.items || data.list || []
  } catch (error) {
    console.error('获取域名列表失败:', error)
    ElMessage.error('获取域名列表失败')
  }
}

onMounted(() => {
  fetchDomains()
})
</script>

<style scoped lang="scss">
.dns-record-create {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 24px;
}

.breadcrumb-nav {
  margin-bottom: 20px;
  
  :deep(.el-breadcrumb__inner) {
    color: rgba(255, 255, 255, 0.8);
    
    &.is-link {
      color: white;
      font-weight: 500;
      
      &:hover {
        color: #ffd700;
      }
    }
  }
  
  :deep(.el-breadcrumb__separator) {
    color: rgba(255, 255, 255, 0.6);
  }
}

.page-header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 32px;
  margin-bottom: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-section {
  display: flex;
  align-items: center;
  gap: 20px;
}

.icon-wrapper {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.3);
}

.title-content {
  h1 {
    margin: 0 0 8px 0;
    font-size: 32px;
    font-weight: 700;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
  
  p {
    margin: 0;
    color: #64748b;
    font-size: 16px;
  }
}

.form-container {
  max-width: 1000px;
  margin: 0 auto;
}

.form-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  
  :deep(.el-card__header) {
    background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
    border-bottom: 1px solid rgba(226, 232, 240, 0.5);
    border-radius: 16px 16px 0 0;
  }
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
}

.record-form {
  padding: 8px 0;
}

.form-item-highlight {
  :deep(.el-form-item__label) {
    font-weight: 600;
    color: #1e293b;
  }
  
  :deep(.el-select) {
    border-radius: 12px;
    
    .el-input__wrapper {
      border-radius: 12px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    }
  }
}

.domain-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  .domain-name {
    font-weight: 500;
  }
}

.divider-text {
  font-weight: 600;
  color: #475569;
  font-size: 14px;
}

.form-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 6px;
  font-size: 12px;
  color: #64748b;
  
  .el-icon {
    color: #94a3b8;
  }
}

.type-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  .type-main {
    display: flex;
    flex-direction: column;
    
    .type-name {
      font-weight: 600;
      color: #1e293b;
    }
    
    .type-desc {
      font-size: 12px;
      color: #64748b;
    }
  }
  
  .type-example {
    font-size: 11px;
    color: #94a3b8;
    font-family: monospace;
  }
}

.form-actions {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #e2e8f0;
  
  :deep(.el-form-item__content) {
    display: flex;
    gap: 16px;
  }
}

.submit-btn {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  border: none;
  padding: 12px 32px;
  border-radius: 12px;
  font-weight: 600;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 20px rgba(16, 185, 129, 0.4);
  }
}

.reset-btn {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 12px;
  font-weight: 600;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 20px rgba(245, 158, 11, 0.4);
  }
}

.cancel-btn {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #475569;
  padding: 12px 24px;
  border-radius: 12px;
  font-weight: 600;
  
  &:hover {
    background: #f1f5f9;
    border-color: #cbd5e1;
    transform: translateY(-1px);
  }
}

// Element Plus 样式覆盖
:deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  transition: all 0.2s ease;
  
  &:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
  }
  
  &.is-focus {
    box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
  }
}

:deep(.el-select-dropdown) {
  border-radius: 12px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(0, 0, 0, 0.08);
}

:deep(.el-textarea__inner) {
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  
  &:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
  }
  
  &:focus {
    box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
  }
}

:deep(.el-input-number) {
  width: 100%;
  
  .el-input__wrapper {
    border-radius: 8px;
  }
}
</style>
