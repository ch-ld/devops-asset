<template>
  <el-dialog
    v-model="dialogVisible"
    :title="modalTitle"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
    class="record-modal"
  >
    <div class="modal-body">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
        @submit.prevent
        class="record-form"
      >
        <!-- 记录类型 -->
        <el-form-item label="记录类型" prop="type">
          <el-select
            v-model="formData.type"
            placeholder="选择记录类型"
            style="width: 100%"
            @change="handleTypeChange"
          >
            <el-option label="A - 将域名指向一个IPv4地址" value="A" />
            <el-option label="AAAA - 将域名指向一个IPv6地址" value="AAAA" />
            <el-option label="CNAME - 将域名指向另一个域名" value="CNAME" />
            <el-option label="MX - 将域名指向邮件服务器地址" value="MX" />
            <el-option label="TXT - 为域名设置说明" value="TXT" />
            <el-option label="NS - 将子域名指定其他DNS服务器解析" value="NS" />
            <el-option label="SRV - 记录提供特定的服务的服务器" value="SRV" />
          </el-select>
        </el-form-item>

        <!-- 主机记录 -->
        <el-form-item label="主机记录" prop="name">
          <el-input
            v-model="formData.name"
            placeholder="请输入主机记录"
            style="width: 100%"
          >
            <template #suffix>
              <span class="domain-suffix">.{{ domainName }}</span>
            </template>
          </el-input>
          <div class="form-help">
            {{ getRecordHelp() }}
          </div>
        </el-form-item>

        <!-- 解析请求来源 -->
        <el-form-item label="解析请求来源" prop="line">
          <el-select
            v-model="formData.line"
            placeholder="默认"
            style="width: 100%"
          >
            <el-option label="默认" value="default" />
            <el-option label="境内" value="cn" />
            <el-option label="境外" value="abroad" />
          </el-select>
          <div class="form-help">
            升级至企业版后域名支持按地域、运营商分别设置解析
          </div>
        </el-form-item>

        <!-- 记录值 -->
        <el-form-item label="记录值" prop="value">
          <el-input
            v-model="formData.value"
            :placeholder="getValuePlaceholder()"
            style="width: 100%"
          />
          <div class="form-help">
            {{ getValueHelp() }}
          </div>
        </el-form-item>

        <!-- TTL时间 -->
        <el-form-item label="TTL" prop="ttl">
          <el-select
            v-model="formData.ttl"
            style="width: 100%"
          >
            <el-option label="10 (秒)" :value="10" />
            <el-option label="60 (秒)" :value="60" />
            <el-option label="600 (秒)" :value="600" />
            <el-option label="3600 (秒)" :value="3600" />
            <el-option label="43200 (秒)" :value="43200" />
            <el-option label="86400 (秒)" :value="86400" />
          </el-select>
          <div class="form-help">
            缓存生存时间，数值越小，修改记录各地生效时间越快
          </div>
        </el-form-item>

        <!-- MX记录优先级 -->
        <el-form-item
          v-if="formData.type === 'MX'"
          label="MX优先级"
          prop="priority"
        >
          <el-input-number
            v-model="formData.priority"
            :min="1"
            :max="50"
            style="width: 100%"
          />
        </el-form-item>

        <!-- SRV记录参数 -->
        <template v-if="formData.type === 'SRV'">
          <el-form-item label="优先级" prop="priority">
            <el-input-number
              v-model="formData.priority"
              :min="0"
              :max="65535"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="权重" prop="weight">
            <el-input-number
              v-model="formData.weight"
              :min="0"
              :max="65535"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="端口" prop="port">
            <el-input-number
              v-model="formData.port"
              :min="1"
              :max="65535"
              style="width: 100%"
            />
          </el-form-item>
        </template>
      </el-form>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">
          {{ isEdit ? '更新' : '确定' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { recordApi } from '@/api/dns/record'
import type { DNSRecord } from '@/types/dns'

interface Props {
  visible: boolean
  record?: DNSRecord | null
  domainId?: number
  domainName?: string
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  record: null
})

const emit = defineEmits<Emits>()

const formRef = ref()
const loading = ref(false)

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const isEdit = computed(() => !!props.record)
const modalTitle = computed(() => isEdit.value ? '编辑记录' : '添加记录')

const formData = ref({
  domain_id: props.domainId || 0,
  name: '',
  type: 'A',
  value: '',
  ttl: 600,
  priority: 10,
  weight: 1,
  port: 80,
  line: 'default'
})

const formRules = {
  type: [{ required: true, message: '请选择记录类型', trigger: 'change' }],
  name: [{ required: true, message: '请输入主机记录', trigger: 'blur' }],
  value: [{ required: true, message: '请输入记录值', trigger: 'blur' }],
  ttl: [{ required: true, message: '请选择TTL', trigger: 'change' }]
}

const handleTypeChange = () => {
  // 根据类型设置默认值
  switch (formData.value.type) {
    case 'MX':
      formData.value.priority = 10
      break
    case 'SRV':
      formData.value.priority = 0
      formData.value.weight = 1
      formData.value.port = 80
      break
  }
}

const getRecordHelp = () => {
  switch (formData.value.type) {
    case 'A':
    case 'AAAA':
      return '如需要将域名指向一个IP地址，可以使用此类型。例如：www、mail、@'
    case 'CNAME':
      return '如需要将域名指向另一个域名，再由另一个域名提供IP地址，可以使用此类型'
    case 'MX':
      return '如需要设置邮箱，让邮箱能收到邮件，可以使用此类型'
    case 'TXT':
      return '如需要对域名进行标识和说明，可以使用此类型'
    case 'NS':
      return '如需要把子域名交给其他DNS服务商解析，可以使用此类型'
    case 'SRV':
      return '如需要记录提供特定服务的服务器，可以使用此类型'
    default:
      return ''
  }
}

const getValuePlaceholder = () => {
  switch (formData.value.type) {
    case 'A':
      return '请输入IPv4地址，例如：192.168.1.1'
    case 'AAAA':
      return '请输入IPv6地址，例如：2001:db8::1'
    case 'CNAME':
      return '请输入域名，例如：www.example.com'
    case 'MX':
      return '请输入邮件服务器地址，例如：mail.example.com'
    case 'TXT':
      return '请输入文本内容'
    case 'NS':
      return '请输入DNS服务器地址，例如：ns1.example.com'
    case 'SRV':
      return '请输入目标地址，例如：target.example.com'
    default:
      return '请输入记录值'
  }
}

const getValueHelp = () => {
  switch (formData.value.type) {
    case 'A':
      return '将域名指向一个IPv4地址，例如：223.5.5.5'
    case 'AAAA':
      return '将域名指向一个IPv6地址'
    case 'CNAME':
      return '将域名指向另一个域名，可以通过该域名获取IP地址'
    case 'MX':
      return '将域名指向邮件服务器地址，用于邮件收发'
    case 'TXT':
      return '为域名设置说明文字，常用于域名验证'
    case 'NS':
      return '将子域名指定其他DNS服务器解析'
    case 'SRV':
      return '记录提供特定服务的服务器'
    default:
      return ''
  }
}
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  loading.value = true
  try {
    const data = { ...formData.value }

    // 如果没有传入domainId，使用表单中的domain_id
    if (!props.domainId) {
      data.domain_id = formData.value.domain_id
    } else {
      data.domain_id = props.domainId
    }

    if (isEdit.value && props.record) {
      await recordApi.update(props.record.id, data)
      ElMessage.success('记录更新成功')
    } else {
      await recordApi.create(data)
      ElMessage.success('记录创建成功')
    }

    emit('success')
    handleClose()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}

const handleClose = () => {
  dialogVisible.value = false
  resetForm()
}

const resetForm = () => {
  formData.value = {
    domain_id: props.domainId || 0,
    name: '',
    type: 'A',
    value: '',
    ttl: 600,
    priority: 10,
    weight: 1,
    port: 80,
    line: 'default'
  }
  nextTick(() => {
    formRef.value?.clearValidate()
  })
}

// 监听props变化
watch(() => props.record, (newRecord) => {
  if (newRecord) {
    formData.value = {
      domain_id: newRecord.domain_id,
      name: newRecord.name,
      type: newRecord.type,
      value: newRecord.value,
      ttl: newRecord.ttl,
      priority: newRecord.priority || 10,
      weight: newRecord.weight || 1,
      port: newRecord.port || 80,
      line: newRecord.line || 'default'
    }
  } else {
    resetForm()
  }
}, { immediate: true })

watch(() => props.visible, (visible) => {
  if (!visible) {
    resetForm()
  }
})
</script>

<style scoped>
.record-modal {
  .modal-body {
    padding: 20px 0;
  }

  .record-form {
    .form-help {
      font-size: 12px;
      color: #909399;
      margin-top: 4px;
      line-height: 1.4;
    }

    .domain-suffix {
      color: #909399;
      font-size: 12px;
    }
  }

  .dialog-footer {
    text-align: right;
    padding: 20px 0 0;
    border-top: 1px solid #ebeef5;
  }
}
</style>

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
  Plus, Edit, Connection, Document, Clock, Sort, Setting,
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
    icon: 'Connection',
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

// 添加缺失的方法
const updateStatus = async (recordId: number, status: string) => {
  try {
    await recordApi.updateStatus(recordId, status)
    return true
  } catch (error) {
    console.error('更新记录状态失败:', error)
    return false
  }
}

const sync = async (recordId: number) => {
  try {
    await recordApi.sync(recordId)
    return true
  } catch (error) {
    console.error('同步记录失败:', error)
    return false
  }
}

const test = async (recordId: number) => {
  try {
    const result = await recordApi.testResolution({
      name: formData.name,
      type: formData.type
    })
    return result
  } catch (error) {
    console.error('测试记录失败:', error)
    return { success: false, error_msg: '测试失败' }
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
</script>

<style scoped lang="scss">
.record-modal {
  :deep(.el-dialog) {
    border-radius: 16px;
    overflow: hidden;
  }

  :deep(.el-dialog__header) {
    padding: 0;
    margin: 0;
  }

  :deep(.el-dialog__body) {
    padding: 0;
  }

  :deep(.el-dialog__footer) {
    padding: 0;
  }

  .modal-header {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 24px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;

    .header-icon {
      width: 48px;
      height: 48px;
      background: rgba(255, 255, 255, 0.2);
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .header-content {
      h3 {
        font-size: 20px;
        font-weight: 600;
        margin: 0 0 4px 0;
      }

      p {
        font-size: 14px;
        opacity: 0.9;
        margin: 0;
      }
    }
  }

  .modal-body {
    padding: 24px;
    max-height: 70vh;
    overflow-y: auto;

    .steps-container {
      margin-bottom: 32px;

      :deep(.el-steps) {
        .el-step__title {
          font-size: 14px;
          font-weight: 500;
        }

        .el-step__description {
          font-size: 12px;
        }
      }
    }

    .form-section {
      margin-bottom: 32px;

      .section-header {
        margin-bottom: 20px;
        padding-bottom: 12px;
        border-bottom: 1px solid #f3f4f6;

        h4 {
          font-size: 16px;
          font-weight: 600;
          color: #1f2937;
          margin: 0 0 4px 0;
        }

        p {
          font-size: 14px;
          color: #6b7280;
          margin: 0;
        }
      }

      .form-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 20px;

        .form-item {
          &.full-width {
            grid-column: 1 / -1;
          }

          .form-label {
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 14px;
            font-weight: 500;
            color: #374151;
            margin-bottom: 8px;
          }

          .form-control {
            width: 100%;

            :deep(.el-input__wrapper) {
              border-radius: 8px;
              border: 2px solid #e5e7eb;
              transition: all 0.3s ease;

              &:hover {
                border-color: #d1d5db;
              }

              &.is-focus {
                border-color: #3b82f6;
                box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
              }
            }

            :deep(.el-select__wrapper) {
              border-radius: 8px;
              border: 2px solid #e5e7eb;
              transition: all 0.3s ease;

              &:hover {
                border-color: #d1d5db;
              }

              &.is-focused {
                border-color: #3b82f6;
                box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
              }
            }
          }

          .form-help {
            font-size: 12px;
            color: #6b7280;
            margin-top: 4px;
          }

          .input-suffix {
            color: #9ca3af;
            font-size: 14px;
          }
        }
      }

      .record-type-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 12px;

        .record-type-card {
          display: flex;
          align-items: center;
          gap: 12px;
          padding: 16px;
          border: 2px solid #e5e7eb;
          border-radius: 12px;
          cursor: pointer;
          transition: all 0.3s ease;

          &:hover {
            border-color: #d1d5db;
            transform: translateY(-1px);
          }

          &.active {
            border-color: #3b82f6;
            background: #eff6ff;
            box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
          }

          .type-icon {
            width: 40px;
            height: 40px;
            border-radius: 10px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
          }

          .type-content {
            .type-name {
              font-size: 14px;
              font-weight: 600;
              color: #1f2937;
            }

            .type-desc {
              font-size: 12px;
              color: #6b7280;
              margin-top: 2px;
            }
          }
        }
      }

      .domain-option {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;

        .domain-name {
          font-weight: 500;
        }
      }

      .ttl-option {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;

        .ttl-label {
          font-weight: 500;
        }

        .ttl-desc {
          font-size: 12px;
          color: #6b7280;
        }
      }
    }

    .preview-section {
      background: #f8fafc;
      border-radius: 12px;
      padding: 20px;
      margin-top: 24px;

      .section-header {
        margin-bottom: 16px;
        padding-bottom: 0;
        border-bottom: none;

        h4 {
          color: #1f2937;
        }

        p {
          color: #6b7280;
        }
      }

      .record-preview {
        display: grid;
        gap: 12px;

        .preview-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 8px 0;

          .preview-label {
            font-size: 14px;
            color: #6b7280;
            font-weight: 500;
          }

          .preview-value {
            font-size: 14px;
            color: #1f2937;
            font-family: 'Monaco', 'Menlo', monospace;
          }
        }
      }
    }
  }

  .modal-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 24px;
    background: #f9fafb;
    border-top: 1px solid #e5e7eb;

    .footer-left,
    .footer-right {
      display: flex;
      gap: 12px;
    }

    .submit-btn {
      background: linear-gradient(135deg, #10b981 0%, #059669 100%);
      border: none;
      color: white;
      font-weight: 500;

      &:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 8px rgba(16, 185, 129, 0.25);
      }
    }
  }
}
</style>
