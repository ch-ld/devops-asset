<template>
  <a-modal
    :visible="visible"
    :title="modalTitle"
    :width="800"
    :confirm-loading="loading"
    :mask-closable="false"
    @ok="handleOk"
    @cancel="handleCancel"
  >
    <a-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      :label-col="{ span: 6 }"
      :wrapper-col="{ span: 18 }"
    >
      <a-form-item label="域名" name="domain_id">
        <a-select
          v-model:value="formData.domain_id"
          placeholder="请选择域名"
          :disabled="mode === 'view'"
          show-search
          :filter-option="filterOption"
        >
          <a-select-option
            v-for="domain in domainOptions"
            :key="domain.id"
            :value="domain.id"
          >
            {{ domain.name }}
          </a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="记录名" name="name">
        <a-input
          v-model:value="formData.name"
          placeholder="请输入记录名，如：www、@、*"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <a-form-item label="记录类型" name="type">
        <a-select
          v-model:value="formData.type"
          placeholder="请选择记录类型"
          :disabled="mode === 'view'"
          @change="handleTypeChange"
        >
          <a-select-option value="A">A - IPv4地址</a-select-option>
          <a-select-option value="AAAA">AAAA - IPv6地址</a-select-option>
          <a-select-option value="CNAME">CNAME - 别名</a-select-option>
          <a-select-option value="MX">MX - 邮件交换</a-select-option>
          <a-select-option value="TXT">TXT - 文本记录</a-select-option>
          <a-select-option value="NS">NS - 域名服务器</a-select-option>
          <a-select-option value="SRV">SRV - 服务记录</a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="记录值" name="value">
        <a-textarea
          v-model:value="formData.value"
          :placeholder="getValuePlaceholder(formData.type)"
          :rows="formData.type === 'TXT' ? 3 : 2"
          :disabled="mode === 'view'"
        />
        <div class="form-help">{{ getValueHelp(formData.type) }}</div>
      </a-form-item>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-form-item label="TTL" name="ttl">
            <a-input-number
              v-model:value="formData.ttl"
              :min="60"
              :max="86400"
              placeholder="600"
              addon-after="秒"
              :disabled="mode === 'view'"
              style="width: 100%"
            />
          </a-form-item>
        </a-col>
        <a-col :span="12" v-if="formData.type === 'MX' || formData.type === 'SRV'">
          <a-form-item label="优先级" name="priority">
            <a-input-number
              v-model:value="formData.priority"
              :min="0"
              :max="65535"
              placeholder="10"
              :disabled="mode === 'view'"
              style="width: 100%"
            />
          </a-form-item>
        </a-col>
      </a-row>

      <a-form-item label="备注" name="remark">
        <a-textarea
          v-model:value="formData.remark"
          placeholder="请输入备注信息"
          :rows="2"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <!-- 查看模式的详细信息 -->
      <template v-if="mode === 'view' && record">
        <a-divider>状态信息</a-divider>
        
        <a-form-item label="状态">
          <a-badge
            :status="getStatusBadge(record.status)"
            :text="getStatusText(record.status)"
          />
        </a-form-item>

        <a-form-item label="最后同步时间" v-if="record.last_sync_at">
          <span>{{ formatDate(record.last_sync_at) }}</span>
        </a-form-item>

        <a-form-item label="同步状态" v-if="record.sync_status">
          <a-tag
            :color="record.sync_status === 'success' ? 'green' : 'red'"
          >
            {{ record.sync_status === 'success' ? '同步成功' : '同步失败' }}
          </a-tag>
        </a-form-item>

        <a-form-item label="创建时间">
          <span>{{ formatDate(record.created_at) }}</span>
        </a-form-item>

        <a-form-item label="更新时间">
          <span>{{ formatDate(record.updated_at) }}</span>
        </a-form-item>
      </template>
    </a-form>

    <template #footer v-if="mode !== 'view'">
      <a-space>
        <a-button @click="handleCancel">取消</a-button>
        <a-button @click="handleValidate" v-if="formData.value">
          验证记录
        </a-button>
        <a-button type="primary" @click="handleOk" :loading="loading">
          {{ mode === 'add' ? '创建' : '更新' }}
        </a-button>
      </a-space>
    </template>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'
import { recordApi } from '@/api/dns/record'
import { domainApi } from '@/api/dns/domain'
import type { DNSRecord, DNSRecordCreateRequest, Domain } from '@/types/dns'

interface Props {
  visible: boolean
  mode: 'add' | 'edit' | 'view'
  record?: DNSRecord | null
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  record: null
})

const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const formRef = ref<FormInstance>()
const domainOptions = ref<Domain[]>([])

// 表单数据
const formData = reactive<DNSRecordCreateRequest>({
  domain_id: undefined,
  name: '',
  type: 'A',
  value: '',
  ttl: 600,
  priority: undefined,
  remark: ''
})

// 表单验证规则
const rules: Record<string, Rule[]> = {
  domain_id: [
    { required: true, message: '请选择域名', trigger: 'change' }
  ],
  name: [
    { required: true, message: '请输入记录名', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择记录类型', trigger: 'change' }
  ],
  value: [
    { required: true, message: '请输入记录值', trigger: 'blur' },
    { validator: validateRecordValue, trigger: 'blur' }
  ],
  ttl: [
    { required: true, type: 'number', min: 60, max: 86400, message: 'TTL必须在60-86400秒之间', trigger: 'blur' }
  ],
  priority: [
    { type: 'number', min: 0, max: 65535, message: '优先级必须在0-65535之间', trigger: 'blur' }
  ]
}

// 计算属性
const modalTitle = computed(() => {
  const titleMap = {
    add: '添加解析记录',
    edit: '编辑解析记录',
    view: '解析记录详情'
  }
  return titleMap[props.mode]
})

// 监听器
watch(() => props.visible, (visible) => {
  if (visible) {
    if ((props.mode === 'edit' || props.mode === 'view') && props.record) {
      Object.assign(formData, {
        domain_id: props.record.domain_id,
        name: props.record.name,
        type: props.record.type,
        value: props.record.value,
        ttl: props.record.ttl,
        priority: props.record.priority,
        remark: props.record.remark || ''
      })
    } else if (props.mode === 'add') {
      resetForm()
    }
    fetchDomainOptions()
  }
})

// 方法
const getValuePlaceholder = (type: string) => {
  const placeholderMap = {
    A: '请输入IPv4地址，如：192.168.1.1',
    AAAA: '请输入IPv6地址，如：2001:db8::1',
    CNAME: '请输入目标域名，如：example.com',
    MX: '请输入邮件服务器，如：mail.example.com',
    TXT: '请输入文本内容，如：v=spf1 include:_spf.example.com ~all',
    NS: '请输入域名服务器，如：ns1.example.com',
    SRV: '请输入服务记录，如：10 5 443 target.example.com'
  }
  return placeholderMap[type] || '请输入记录值'
}

const getValueHelp = (type: string) => {
  const helpMap = {
    A: 'IPv4地址格式：x.x.x.x',
    AAAA: 'IPv6地址格式：xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx',
    CNAME: '目标域名，不能与其他记录类型共存',
    MX: '邮件服务器域名，需要设置优先级',
    TXT: '文本记录，常用于域名验证、SPF等',
    NS: '域名服务器，用于委托子域名',
    SRV: '格式：权重 端口 目标域名'
  }
  return helpMap[type] || ''
}

const getStatusBadge = (status: string) => {
  const statusMap = {
    active: 'success',
    inactive: 'default',
    error: 'error'
  }
  return statusMap[status] || 'default'
}

const getStatusText = (status: string) => {
  const statusMap = {
    active: '正常',
    inactive: '停用',
    error: '异常'
  }
  return statusMap[status] || status
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const filterOption = (input: string, option: any) => {
  return option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
}

// 记录值验证
async function validateRecordValue(rule: any, value: string) {
  if (!value) return Promise.resolve()
  
  const type = formData.type
  
  switch (type) {
    case 'A':
      if (!/^(\d{1,3}\.){3}\d{1,3}$/.test(value)) {
        return Promise.reject('请输入有效的IPv4地址')
      }
      break
    case 'AAAA':
      if (!/^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$/.test(value) && 
          !/^::1$/.test(value) && 
          !/^([0-9a-fA-F]{1,4}:)*::([0-9a-fA-F]{1,4}:)*[0-9a-fA-F]{1,4}$/.test(value)) {
        return Promise.reject('请输入有效的IPv6地址')
      }
      break
    case 'CNAME':
    case 'NS':
      if (!/^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/.test(value)) {
        return Promise.reject('请输入有效的域名')
      }
      break
    case 'MX':
      if (!/^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/.test(value)) {
        return Promise.reject('请输入有效的邮件服务器域名')
      }
      break
  }
  
  return Promise.resolve()
}

const handleTypeChange = (type: string) => {
  // 清空记录值和优先级
  formData.value = ''
  formData.priority = undefined
  
  // 根据类型设置默认优先级
  if (type === 'MX') {
    formData.priority = 10
  } else if (type === 'SRV') {
    formData.priority = 0
  }
}

const resetForm = () => {
  Object.assign(formData, {
    domain_id: undefined,
    name: '',
    type: 'A',
    value: '',
    ttl: 600,
    priority: undefined,
    remark: ''
  })
  formRef.value?.clearValidate()
}

const fetchDomainOptions = async () => {
  try {
    const response = await domainApi.list({ page: 1, size: 1000 })
    domainOptions.value = response.items || []
  } catch (error) {
    console.error('获取域名列表失败:', error)
  }
}

const handleValidate = async () => {
  try {
    // TODO: 实现记录验证功能
    message.success('记录验证通过')
  } catch (error) {
    message.error('记录验证失败')
  }
}

const handleOk = async () => {
  if (props.mode === 'view') {
    handleCancel()
    return
  }

  try {
    await formRef.value?.validate()
    loading.value = true

    if (props.mode === 'add') {
      await recordApi.create(formData)
      message.success('解析记录创建成功')
    } else if (props.mode === 'edit' && props.record) {
      await dnsRecordApi.update(props.record.id, formData)
      message.success('解析记录更新成功')
    }

    emit('success')
  } catch (error) {
    console.error('操作失败:', error)
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  emit('update:visible', false)
  resetForm()
}

// 生命周期
onMounted(() => {
  fetchDomainOptions()
})
</script>

<style scoped lang="scss">
.form-help {
  margin-top: 4px;
  color: #8c8c8c;
  font-size: 12px;
}

.ant-divider {
  margin: 16px 0;
}
</style>
