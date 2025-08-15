<template>
  <el-dialog
    v-model="visible"
    :title="isEdit ? '编辑DNS记录' : '添加DNS记录'"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="120px"
      @submit.prevent
    >
      <el-form-item label="域名" prop="domain_id">
        <el-select
          v-model="formData.domain_id"
          placeholder="请选择域名"
          style="width: 100%"
          filterable
        >
          <el-option
            v-for="domain in domains"
            :key="domain.id"
            :label="domain.name"
            :value="domain.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="记录名" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入记录名，如：www、mail、@（根域名）"
        />
      </el-form-item>

      <el-form-item label="记录类型" prop="type">
        <el-select
          v-model="formData.type"
          placeholder="请选择记录类型"
          style="width: 100%"
          @change="handleTypeChange"
        >
          <el-option label="A" value="A">
            <span>A - IPv4地址记录</span>
          </el-option>
          <el-option label="AAAA" value="AAAA">
            <span>AAAA - IPv6地址记录</span>
          </el-option>
          <el-option label="CNAME" value="CNAME">
            <span>CNAME - 别名记录</span>
          </el-option>
          <el-option label="MX" value="MX">
            <span>MX - 邮件交换记录</span>
          </el-option>
          <el-option label="TXT" value="TXT">
            <span>TXT - 文本记录</span>
          </el-option>
          <el-option label="NS" value="NS">
            <span>NS - 域名服务器记录</span>
          </el-option>
          <el-option label="SRV" value="SRV">
            <span>SRV - 服务记录</span>
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="记录值" prop="value">
        <el-input
          v-model="formData.value"
          :placeholder="getValuePlaceholder()"
          type="textarea"
          :rows="formData.type === 'TXT' ? 3 : 1"
        />
      </el-form-item>

      <el-row :gutter="16">
        <el-col :span="12">
      <el-form-item label="TTL" prop="ttl">
            <el-select
          v-model="formData.ttl"
              placeholder="请选择TTL"
              style="width: 100%"
            >
              <el-option label="1分钟 (60)" :value="60" />
              <el-option label="5分钟 (300)" :value="300" />
              <el-option label="10分钟 (600)" :value="600" />
              <el-option label="30分钟 (1800)" :value="1800" />
              <el-option label="1小时 (3600)" :value="3600" />
              <el-option label="12小时 (43200)" :value="43200" />
              <el-option label="1天 (86400)" :value="86400" />
            </el-select>
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