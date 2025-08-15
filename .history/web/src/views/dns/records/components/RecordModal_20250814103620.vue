<template>
  <el-dialog
    v-model="dialogVisible"
    :title="modalTitle"
    width="800px"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    @close="handleClose"
    class="record-modal"
  >
    <template #header>
      <div class="modal-header">
        <div class="header-icon">
          <el-icon size="24" :color="isEdit ? '#f59e0b' : '#10b981'">
            <component :is="isEdit ? 'Edit' : 'Plus'" />
          </el-icon>
        </div>
        <div class="header-content">
          <h3>{{ modalTitle }}</h3>
          <p>{{ modalSubtitle }}</p>
        </div>
      </div>
    </template>

    <div class="modal-body">
      <!-- 步骤指示器 -->
      <div class="steps-container" v-if="!isEdit">
        <el-steps :active="currentStep" align-center>
          <el-step title="基本信息" description="设置记录名称和类型" />
          <el-step title="记录配置" description="配置记录值和参数" />
          <el-step title="高级设置" description="TTL和其他选项" />
        </el-steps>
      </div>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="0"
        @submit.prevent
        class="record-form"
      >
        <!-- 步骤1: 基本信息 -->
        <div v-show="currentStep === 0 || isEdit" class="form-section">
          <div class="section-header">
            <h4>基本信息</h4>
            <p>设置DNS记录的基本信息</p>
          </div>

          <div class="form-grid">
            <div class="form-item" v-if="!domainId">
              <label class="form-label">
                <el-icon><Globe /></el-icon>
                域名
              </label>
              <el-select
                v-model="formData.domain_id"
                placeholder="请选择域名"
                size="large"
                filterable
                class="form-control"
              >
                <el-option
                  v-for="domain in domains"
                  :key="domain.id"
                  :label="domain.name"
                  :value="domain.id"
                >
                  <div class="domain-option">
                    <span class="domain-name">{{ domain.name }}</span>
                    <el-tag size="small" type="info">{{ domain.provider?.name }}</el-tag>
                  </div>
                </el-option>
              </el-select>
            </div>

            <div class="form-item">
              <label class="form-label">
                <el-icon><Edit /></el-icon>
                记录名
              </label>
              <el-input
                v-model="formData.name"
                placeholder="www, mail, @ (根域名)"
                size="large"
                class="form-control"
              >
                <template #suffix>
                  <span class="input-suffix">.{{ selectedDomainName }}</span>
                </template>
              </el-input>
              <div class="form-help">
                完整域名: {{ fullDomainName }}
              </div>
            </div>

            <div class="form-item">
              <label class="form-label">
                <el-icon><Document /></el-icon>
                记录类型
              </label>
              <div class="record-type-grid">
                <div
                  v-for="type in recordTypes"
                  :key="type.value"
                  :class="['record-type-card', { active: formData.type === type.value }]"
                  @click="selectRecordType(type.value)"
                >
                  <div class="type-icon" :style="{ background: type.color }">
                    <el-icon size="20">
                      <component :is="type.icon" />
                    </el-icon>
                  </div>
                  <div class="type-content">
                    <div class="type-name">{{ type.value }}</div>
                    <div class="type-desc">{{ type.description }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 步骤2: 记录配置 -->
        <div v-show="currentStep === 1 || isEdit" class="form-section">
          <div class="section-header">
            <h4>记录配置</h4>
            <p>配置{{ formData.type }}记录的具体值</p>
          </div>

          <div class="form-grid">
            <div class="form-item full-width">
              <label class="form-label">
                <el-icon><Connection /></el-icon>
                记录值
              </label>
              <el-input
                v-model="formData.value"
                :placeholder="getValuePlaceholder()"
                :type="formData.type === 'TXT' ? 'textarea' : 'text'"
                :rows="formData.type === 'TXT' ? 4 : 1"
                size="large"
                class="form-control"
              />
              <div class="form-help">
                {{ getValueHelp() }}
              </div>
            </div>

            <!-- MX和SRV记录的优先级 -->
            <div class="form-item" v-if="needsPriority">
              <label class="form-label">
                <el-icon><Sort /></el-icon>
                优先级
              </label>
              <el-input-number
                v-model="formData.priority"
                :min="0"
                :max="65535"
                size="large"
                class="form-control"
                placeholder="0-65535"
              />
              <div class="form-help">
                数值越小优先级越高
              </div>
            </div>

            <!-- SRV记录的权重和端口 -->
            <div class="form-item" v-if="formData.type === 'SRV'">
              <label class="form-label">
                <el-icon><Scale /></el-icon>
                权重
              </label>
              <el-input-number
                v-model="formData.weight"
                :min="0"
                :max="65535"
                size="large"
                class="form-control"
                placeholder="0-65535"
              />
            </div>

            <div class="form-item" v-if="formData.type === 'SRV'">
              <label class="form-label">
                <el-icon><Connection /></el-icon>
                端口
              </label>
              <el-input-number
                v-model="formData.port"
                :min="1"
                :max="65535"
                size="large"
                class="form-control"
                placeholder="1-65535"
              />
            </div>
          </div>
        </div>

        <!-- 步骤3: 高级设置 -->
        <div v-show="currentStep === 2 || isEdit" class="form-section">
          <div class="section-header">
            <h4>高级设置</h4>
            <p>配置TTL和其他高级选项</p>
          </div>

          <div class="form-grid">
            <div class="form-item">
              <label class="form-label">
                <el-icon><Clock /></el-icon>
                TTL (生存时间)
              </label>
              <el-select
                v-model="formData.ttl"
                placeholder="选择TTL值"
                size="large"
                class="form-control"
              >
                <el-option
                  v-for="ttl in ttlOptions"
                  :key="ttl.value"
                  :label="ttl.label"
                  :value="ttl.value"
                >
                  <div class="ttl-option">
                    <span class="ttl-label">{{ ttl.label }}</span>
                    <span class="ttl-desc">{{ ttl.description }}</span>
                  </div>
                </el-option>
              </el-select>
              <div class="form-help">
                TTL决定DNS记录在缓存中保存的时间
              </div>
            </div>

            <div class="form-item">
              <label class="form-label">
                <el-icon><Document /></el-icon>
                备注
              </label>
              <el-input
                v-model="formData.remark"
                placeholder="添加备注信息（可选）"
                type="textarea"
                :rows="3"
                size="large"
                class="form-control"
              />
              <div class="form-help">
                备注信息仅用于管理，不会影响DNS解析
              </div>
            </div>
          </div>
        </div>
      </el-form>

      <!-- 预览区域 -->
      <div class="preview-section" v-if="formData.type && formData.value">
        <div class="section-header">
          <h4>记录预览</h4>
          <p>确认您的DNS记录配置</p>
        </div>
        <div class="record-preview">
          <div class="preview-item">
            <span class="preview-label">完整域名:</span>
            <span class="preview-value">{{ fullDomainName }}</span>
          </div>
          <div class="preview-item">
            <span class="preview-label">记录类型:</span>
            <el-tag :type="getRecordTypeTagType(formData.type)">{{ formData.type }}</el-tag>
          </div>
          <div class="preview-item">
            <span class="preview-label">记录值:</span>
            <span class="preview-value">{{ formData.value }}</span>
          </div>
          <div class="preview-item" v-if="formData.ttl">
            <span class="preview-label">TTL:</span>
            <span class="preview-value">{{ formatTTL(formData.ttl) }}</span>
          </div>
          <div class="preview-item" v-if="formData.priority">
            <span class="preview-label">优先级:</span>
            <span class="preview-value">{{ formData.priority }}</span>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="modal-footer">
        <div class="footer-left">
          <el-button v-if="!isEdit && currentStep > 0" @click="prevStep" size="large">
            <el-icon><ArrowLeft /></el-icon>
            上一步
          </el-button>
        </div>
        <div class="footer-right">
          <el-button @click="handleClose" size="large">取消</el-button>
          <el-button
            v-if="!isEdit && currentStep < 2"
            type="primary"
            @click="nextStep"
            size="large"
            :disabled="!canProceedToNext"
          >
            下一步
            <el-icon><ArrowRight /></el-icon>
          </el-button>
          <el-button
            v-else
            type="primary"
            :loading="loading"
            @click="handleSubmit"
            size="large"
            class="submit-btn"
          >
            <el-icon><Check /></el-icon>
            {{ isEdit ? '更新记录' : '创建记录' }}
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Plus, Edit, Globe, Document, Connection, Clock, Sort, Scale,
  ArrowLeft, ArrowRight, Check
} from '@element-plus/icons-vue'
import { recordApi } from '@/api/dns/record'
import type { DnsRecord, Domain } from '@/types/dns'

interface Props {
  visible: boolean
  record?: DnsRecord | null
  domainId?: number | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  record: null,
  domainId: null
})

const emit = defineEmits<Emits>()

// 响应式数据
const formRef = ref()
const loading = ref(false)
const currentStep = ref(0)

const isEdit = computed(() => !!props.record)

// 对话框可见性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 表单数据
const formData = reactive({
  domain_id: props.domainId || 0,
  name: '',
  type: 'A',
  value: '',
  ttl: 600,
  priority: 10,
  weight: 10,
  port: 80,
  remark: ''
})

// 计算属性
const modalTitle = computed(() => {
  return isEdit.value ? '编辑DNS记录' : '添加DNS记录'
})

const modalSubtitle = computed(() => {
  return isEdit.value ? '修改现有DNS记录的配置' : '为域名添加新的DNS解析记录'
})

const selectedDomainName = computed(() => {
  const domain = domains.value.find(d => d.id === formData.domain_id)
  return domain?.name || 'example.com'
})

const fullDomainName = computed(() => {
  if (!formData.name || formData.name === '@') {
    return selectedDomainName.value
  }
  return `${formData.name}.${selectedDomainName.value}`
})

const needsPriority = computed(() => {
  return ['MX', 'SRV'].includes(formData.type)
})

const canProceedToNext = computed(() => {
  switch (currentStep.value) {
    case 0:
      return formData.domain_id && formData.name && formData.type
    case 1:
      return formData.value
    case 2:
      return formData.ttl
    default:
      return false
  }
})

// 记录类型配置
const recordTypes = ref([
  {
    value: 'A',
    description: 'IPv4地址记录',
    icon: 'Connection',
    color: 'linear-gradient(135deg, #10b981 0%, #059669 100%)'
  },
  {
    value: 'AAAA',
    description: 'IPv6地址记录',
    icon: 'Connection',
    color: 'linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%)'
  },
  {
    value: 'CNAME',
    description: '别名记录',
    icon: 'Link',
    color: 'linear-gradient(135deg, #f59e0b 0%, #d97706 100%)'
  },
  {
    value: 'MX',
    description: '邮件交换记录',
    icon: 'Message',
    color: 'linear-gradient(135deg, #ef4444 0%, #dc2626 100%)'
  },
  {
    value: 'TXT',
    description: '文本记录',
    icon: 'Document',
    color: 'linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%)'
  },
  {
    value: 'NS',
    description: '域名服务器记录',
    icon: 'Globe',
    color: 'linear-gradient(135deg, #6b7280 0%, #4b5563 100%)'
  }
])

// TTL选项
const ttlOptions = ref([
  { value: 60, label: '1分钟', description: '快速更新，适合测试' },
  { value: 300, label: '5分钟', description: '较快更新' },
  { value: 600, label: '10分钟', description: '推荐设置' },
  { value: 1800, label: '30分钟', description: '标准设置' },
  { value: 3600, label: '1小时', description: '稳定设置' },
  { value: 43200, label: '12小时', description: '长期稳定' },
  { value: 86400, label: '1天', description: '最大稳定性' }
])

// 域名列表（从父组件获取或API获取）
const domains = ref<Domain[]>([])

// 表单验证规则
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

// 工具方法
const selectRecordType = (type: string) => {
  formData.type = type
  // 清空记录值，让用户重新输入
  formData.value = ''
}

const getValuePlaceholder = () => {
  const placeholders: Record<string, string> = {
    A: '192.168.1.1',
    AAAA: '2001:db8::1',
    CNAME: 'example.com',
    MX: 'mail.example.com',
    TXT: '"v=spf1 include:_spf.example.com ~all"',
    NS: 'ns1.example.com',
    SRV: 'target.example.com'
  }
  return placeholders[formData.type] || '请输入记录值'
}

const getValueHelp = () => {
  const helps: Record<string, string> = {
    A: '输入IPv4地址，如：192.168.1.1',
    AAAA: '输入IPv6地址，如：2001:db8::1',
    CNAME: '输入目标域名，如：example.com',
    MX: '输入邮件服务器域名，如：mail.example.com',
    TXT: '输入文本内容，通常用于验证或配置',
    NS: '输入域名服务器，如：ns1.example.com',
    SRV: '输入目标服务器域名和端口'
  }
  return helps[formData.type] || ''
}

const getRecordTypeTagType = (type: string) => {
  const typeMap: Record<string, any> = {
    A: 'success',
    AAAA: 'primary',
    CNAME: 'warning',
    MX: 'danger',
    TXT: 'info',
    NS: 'info'
  }
  return typeMap[type] || 'info'
}

const formatTTL = (ttl: number) => {
  if (ttl >= 86400) {
    return `${Math.floor(ttl / 86400)}天`
  } else if (ttl >= 3600) {
    return `${Math.floor(ttl / 3600)}小时`
  } else if (ttl >= 60) {
    return `${Math.floor(ttl / 60)}分钟`
  }
  return `${ttl}秒`
}

// 步骤控制
const nextStep = () => {
  if (currentStep.value < 2) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

// 事件处理
const handleClose = () => {
  emit('update:visible', false)
  resetForm()
}

const resetForm = () => {
  currentStep.value = 0
  Object.assign(formData, {
    domain_id: props.domainId || 0,
    name: '',
    type: 'A',
    value: '',
    ttl: 600,
    priority: 10,
    weight: 10,
    port: 80,
    remark: ''
  })
  nextTick(() => {
    formRef.value?.clearValidate()
  })
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    loading.value = true

    const submitData = {
      ...formData,
      // 确保数值类型正确
      ttl: Number(formData.ttl),
      priority: needsPriority.value ? Number(formData.priority) : undefined,
      weight: formData.type === 'SRV' ? Number(formData.weight) : undefined,
      port: formData.type === 'SRV' ? Number(formData.port) : undefined
    }

    if (isEdit.value && props.record) {
      await recordApi.update(props.record.id, submitData)
      ElMessage.success('DNS记录更新成功')
    } else {
      await recordApi.create(submitData)
      ElMessage.success('DNS记录创建成功')
    }

    emit('success')
    handleClose()
  } catch (error) {
    console.error('操作DNS记录失败:', error)
    ElMessage.error(isEdit.value ? 'DNS记录更新失败' : 'DNS记录创建失败')
  } finally {
    loading.value = false
  }
}

const loadRecord = () => {
  if (props.record) {
    Object.assign(formData, {
      domain_id: props.record.domain_id,
      name: props.record.name,
      type: props.record.type,
      value: props.record.value,
      ttl: props.record.ttl,
      priority: props.record.priority || 10,
      weight: props.record.weight || 10,
      port: props.record.port || 80,
      remark: props.record.remark || ''
    })
  } else if (props.domainId) {
    formData.domain_id = props.domainId
  }
}

// 监听器
watch(() => props.visible, (newVal) => {
  if (newVal) {
    loadRecord()
    currentStep.value = 0
  }
})

watch(() => props.domainId, (newVal) => {
  if (newVal) {
    formData.domain_id = newVal
  }
})
    type: 'A',
    value: '',
    ttl: 600,
    priority: 10,
    weight: 10,
    port: 80,
    remark: ''
  })
  
  nextTick(() => {
    formRef.value?.clearValidate()
  })
}

function loadRecord() {
  if (props.record) {
    Object.assign(formData, {
      domain_id: props.record.domain_id,
      name: props.record.name,
      type: props.record.type,
      value: props.record.value,
      ttl: props.record.ttl,
      priority: props.record.priority || 10,
      weight: props.record.weight || 10,
      port: props.record.port || 80,
      remark: props.record.remark || ''
    })
  } else {
    resetForm()
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    loading.value = true
    
    if (isEdit.value) {
      await recordApi.update(props.record!.id, formData)
      ElMessage.success('DNS记录更新成功')
    } else {
      await recordApi.create(formData)
      ElMessage.success('DNS记录创建成功')
    }
    
    emit('success')
  } catch (error) {
    console.error('操作DNS记录失败:', error)
    ElMessage.error(isEdit.value ? 'DNS记录更新失败' : 'DNS记录创建失败')
  } finally {
    loading.value = false
  }
}

function handleClose() {
  emit('update:visible', false)
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    loadRecord()
  }
})
</script>

<style scoped lang="scss">
.dialog-footer {
  text-align: right;
}

:deep(.el-dialog__body) {
  padding: 20px 24px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-select-dropdown__item) {
  height: auto;
  line-height: 1.4;
  padding: 8px 20px;
  
  span {
    font-size: 13px;
    color: #606266;
  }
}
</style>
