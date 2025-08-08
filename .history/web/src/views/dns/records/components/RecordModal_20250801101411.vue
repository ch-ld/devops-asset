<template>
  <el-dialog
    v-model="dialogVisible"
    :title="modalTitle"
    width="800px"
    :close-on-click-modal="false"
    @close="handleCancel"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="120px"
      :disabled="mode === 'view'"
    >
      <el-form-item label="域名" prop="domain_id">
        <el-select
          v-model="formData.domain_id"
          placeholder="请选择域名"
          filterable
          style="width: 100%"
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
          placeholder="请输入记录名，如：www、@、*"
        />
      </el-form-item>

      <el-form-item label="记录类型" prop="type">
        <el-select
          v-model="formData.type"
          placeholder="请选择记录类型"
          @change="handleTypeChange"
          style="width: 100%"
        >
          <el-option label="A" value="A" />
          <el-option label="AAAA" value="AAAA" />
          <el-option label="CNAME" value="CNAME" />
          <el-option label="MX" value="MX" />
          <el-option label="TXT" value="TXT" />
          <el-option label="NS" value="NS" />
          <el-option label="SRV" value="SRV" />
        </el-select>
      </el-form-item>

      <el-form-item label="记录值" prop="value">
        <el-input
          v-model="formData.value"
          type="textarea"
          :rows="getValueRows()"
          :placeholder="getValuePlaceholder()"
        />
      </el-form-item>

      <el-form-item label="TTL" prop="ttl">
        <el-input-number
          v-model="formData.ttl"
          :min="60"
          :max="86400"
          :step="60"
          style="width: 200px"
        />
        <span style="margin-left: 8px; color: #666;">秒</span>
      </el-form-item>

      <el-form-item 
        v-if="showPriority" 
        label="优先级" 
        prop="priority"
      >
        <el-input-number
          v-model="formData.priority"
          :min="0"
          :max="65535"
          style="width: 200px"
        />
      </el-form-item>

      <el-form-item label="备注" prop="remark">
        <el-input
          v-model="formData.remark"
          type="textarea"
          :rows="2"
          placeholder="请输入备注信息"
        />
      </el-form-item>

      <!-- 查看模式下显示额外信息 -->
      <template v-if="mode === 'view' && record">
        <el-form-item label="状态">
          <el-tag :type="getStatusType(record.status)">
            {{ getStatusText(record.status) }}
          </el-tag>
        </el-form-item>

        <el-form-item label="同步状态">
          <el-tag :type="getSyncStatusType(record.sync_status)" size="small">
            {{ getSyncStatusText(record.sync_status) }}
          </el-tag>
        </el-form-item>

        <el-form-item label="最后同步时间" v-if="record.last_sync_at">
          {{ formatDate(record.last_sync_at) }}
        </el-form-item>

        <el-form-item label="创建时间">
          {{ formatDate(record.created_at) }}
        </el-form-item>

        <el-form-item label="更新时间">
          {{ formatDate(record.updated_at) }}
        </el-form-item>
      </template>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button 
          v-if="mode !== 'view'" 
          type="primary" 
          :loading="loading"
          @click="handleOk"
        >
          {{ mode === 'add' ? '创建' : '更新' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { recordApi } from '@/api/dns/record'
import type { DnsRecord, Domain } from '@/types/dns'

interface Props {
  visible: boolean
  mode?: 'add' | 'edit' | 'view'
  record?: DnsRecord | null
  domains?: Domain[]
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'add',
  record: null,
  domains: () => []
})

const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  domain_id: undefined as number | undefined,
  name: '',
  type: '',
  value: '',
  ttl: 300,
  priority: undefined as number | undefined,
  remark: ''
})

// 对话框可见性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 模态框标题
const modalTitle = computed(() => {
  const titles = {
    add: '添加DNS记录',
    edit: '编辑DNS记录',
    view: '查看DNS记录'
  }
  return titles[props.mode]
})

// 是否显示优先级字段
const showPriority = computed(() => {
  return ['MX', 'SRV'].includes(formData.type)
})

// 表单验证规则
const rules = {
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
    { required: true, message: '请输入记录值', trigger: 'blur' }
  ],
  ttl: [
    { required: true, message: '请输入TTL值', trigger: 'blur' },
    { type: 'number', min: 60, max: 86400, message: 'TTL值必须在60-86400之间', trigger: 'blur' }
  ]
}

// 工具方法
const getStatusType = (status: string) => {
  const statusMap: Record<string, any> = {
    active: 'success',
    inactive: 'info',
    syncing: 'warning',
    error: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '正常',
    inactive: '停用',
    syncing: '同步中',
    error: '异常'
  }
  return statusMap[status] || status
}

const getSyncStatusType = (syncStatus: string) => {
  const statusMap: Record<string, any> = {
    synced: 'success',
    pending: 'warning',
    failed: 'danger'
  }
  return statusMap[syncStatus] || 'info'
}

const getSyncStatusText = (syncStatus: string) => {
  const statusMap: Record<string, string> = {
    synced: '已同步',
    pending: '待同步',
    failed: '同步失败'
  }
  return statusMap[syncStatus] || syncStatus
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const getValueRows = () => {
  return formData.type === 'TXT' ? 3 : 1
}

const getValuePlaceholder = () => {
  const placeholders: Record<string, string> = {
    A: '请输入IPv4地址，如：192.168.1.1',
    AAAA: '请输入IPv6地址，如：2001:db8::1',
    CNAME: '请输入目标域名，如：example.com',
    MX: '请输入邮件服务器地址，如：mail.example.com',
    TXT: '请输入文本内容',
    NS: '请输入名称服务器，如：ns1.example.com',
    SRV: '请输入服务器信息，格式：权重 端口 目标'
  }
  return placeholders[formData.type] || '请输入记录值'
}

// 事件处理
const handleTypeChange = () => {
  // 清空记录值和优先级
  formData.value = ''
  formData.priority = undefined
  
  // 根据类型设置默认优先级
  if (formData.type === 'MX') {
    formData.priority = 10
  }
}

const handleOk = async () => {
  try {
    await formRef.value?.validate()
    
    loading.value = true
    
    const data = {
      domain_id: formData.domain_id!,
      name: formData.name,
      type: formData.type,
      value: formData.value,
      ttl: formData.ttl,
      priority: formData.priority,
      remark: formData.remark
    }
    
    if (props.mode === 'add') {
      await recordApi.create(data)
      ElMessage.success('DNS记录创建成功')
    } else {
      await recordApi.update(props.record!.id, data)
      ElMessage.success('DNS记录更新成功')
    }
    
    emit('success')
  } catch (error) {
    ElMessage.error(props.mode === 'add' ? '创建失败' : '更新失败')
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  emit('update:visible', false)
}

// 重置表单
const resetForm = () => {
  formData.domain_id = null
  formData.name = ''
  formData.type = ''
  formData.value = ''
  formData.ttl = 300
  formData.priority = undefined
  formData.remark = ''
  
  formRef.value?.clearValidate()
}

// 监听props变化
watch(() => props.visible, (visible) => {
  if (visible) {
    if (props.mode === 'add') {
      resetForm()
    } else if (props.record) {
      // 编辑或查看模式，填充表单数据
      formData.domain_id = props.record.domain_id
      formData.name = props.record.name
      formData.type = props.record.type
      formData.value = props.record.value
      formData.ttl = props.record.ttl
      formData.priority = props.record.priority
      formData.remark = props.record.remark || ''
    }
  }
})
</script>

<style scoped lang="scss">
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
