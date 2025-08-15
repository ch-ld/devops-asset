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
      >
        <el-input-number
          v-model="formData.priority"
          :min="0"
          :max="65535"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row v-if="formData.type === 'SRV'" :gutter="16">
        <el-col :span="12">
          <el-form-item label="权重" prop="weight">
            <el-input-number
              v-model="formData.weight"
              :min="0"
              :max="65535"
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
              style="width: 100%"
        />
      </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="备注" prop="remark">
        <el-input
          v-model="formData.remark"
          placeholder="请输入备注信息（可选）"
          type="textarea"
          :rows="2"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { recordApi } from '@/api/dns/record'
import type { DnsRecord, Domain } from '@/types/dns'

interface Props {
  visible: boolean
  record?: DnsRecord | null
  domains: Domain[]
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  record: null
})

const emit = defineEmits<Emits>()

const formRef = ref()
const loading = ref(false)

const isEdit = computed(() => !!props.record)

// 对话框可见性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

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
    { required: true, message: '请输入记录值', trigger: 'blur' },
    { validator: validateRecordValue, trigger: 'blur' }
  ],
  ttl: [
    { required: true, message: '请选择TTL', trigger: 'change' }
  ],
  priority: [
    { type: 'number', min: 0, max: 65535, message: '优先级必须在0-65535之间', trigger: 'blur' }
  ],
  weight: [
    { type: 'number', min: 0, max: 65535, message: '权重必须在0-65535之间', trigger: 'blur' }
  ],
  port: [
    { type: 'number', min: 1, max: 65535, message: '端口必须在1-65535之间', trigger: 'blur' }
  ]
}

function validateRecordValue(rule: any, value: string, callback: Function) {
  if (!value) {
    callback(new Error('请输入记录值'))
    return
  }

  const type = formData.type
  
  switch (type) {
    case 'A':
      if (!/^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/.test(value)) {
        callback(new Error('请输入有效的IPv4地址'))
        return
      }
      break
    case 'AAAA':
      if (!/^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$/.test(value) && 
          !/^::$/.test(value) && 
          !/^::1$/.test(value)) {
        callback(new Error('请输入有效的IPv6地址'))
        return
      }
      break
    case 'CNAME':
    case 'NS':
      if (!/^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*\.?$/.test(value)) {
        callback(new Error('请输入有效的域名'))
        return
      }
      break
  }
  
  callback()
}

function getValuePlaceholder() {
  switch (formData.type) {
    case 'A':
      return '如：192.168.1.1'
    case 'AAAA':
      return '如：2001:0db8:85a3:0000:0000:8a2e:0370:7334'
    case 'CNAME':
      return '如：www.example.com'
    case 'MX':
      return '如：mail.example.com'
    case 'TXT':
      return '如：v=spf1 include:_spf.example.com ~all'
    case 'NS':
      return '如：ns1.example.com'
    case 'SRV':
      return '如：target.example.com'
    default:
      return '请输入记录值'
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

function resetForm() {
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
