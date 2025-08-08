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
    >
      <el-form-item label="域名" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入域名，如：example.com"
        />
      </el-form-item>

      <el-form-item label="注册商类型" prop="registrar_type">
        <el-select
          v-model="formData.registrar_type"
          placeholder="请选择注册商类型"
          style="width: 100%"
        >
          <el-option label="阿里云" value="aliyun" />
          <el-option label="腾讯云" value="tencent" />
          <el-option label="AWS" value="aws" />
          <el-option label="GoDaddy" value="godaddy" />
          <el-option label="其他" value="other" />
        </el-select>
      </el-form-item>

      <el-form-item label="域名分组" prop="group_id">
        <el-select
          v-model="formData.group_id"
          placeholder="请选择域名分组"
          style="width: 100%"
          clearable
        >
          <el-option
            v-for="group in groups"
            :key="group.id"
            :label="group.name"
            :value="group.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="自动续期">
        <el-switch 
          v-model="formData.auto_renew"
        />
        <span style="margin-left: 8px; color: #666; font-size: 12px;">
          域名即将过期时自动续期
        </span>
      </el-form-item>

      <el-form-item label="过期时间" prop="expires_at">
        <el-date-picker
          v-model="formData.expires_at"
          type="date"
          placeholder="选择过期时间"
          style="width: 100%"
        />
      </el-form-item>

      <el-form-item label="备注" prop="remark">
        <el-input
          v-model="formData.remark"
          type="textarea"
          :rows="3"
          placeholder="请输入备注信息"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button 
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
import { domainApi } from '@/api/dns/domain'
import type { Domain, DomainGroup } from '@/types/dns'

interface Props {
  visible: boolean
  mode?: 'add' | 'edit' | 'view'
  domain?: Domain | null
  groups?: DomainGroup[]
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'add',
  domain: null,
  groups: () => []
})

const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  name: '',
  registrar_type: '',
  group_id: undefined as number | undefined,
  auto_renew: false,
  expires_at: null as Date | null,
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
    add: '添加域名',
    edit: '编辑域名',
    view: '查看域名'
  }
  return titles[props.mode]
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入域名', trigger: 'blur' },
    { 
      pattern: /^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/,
      message: '请输入有效的域名格式',
      trigger: 'blur'
    }
  ],
  registrar_type: [
    { required: true, message: '请选择注册商类型', trigger: 'change' }
  ]
}

// 事件处理
const handleOk = async () => {
  try {
    await formRef.value?.validate()
    
    loading.value = true
    
    const data = {
      name: formData.name,
      registrar_type: formData.registrar_type,
      group_id: formData.group_id,
      auto_renew: formData.auto_renew,
      expires_at: formData.expires_at?.toISOString().split('T')[0],
      remark: formData.remark
    }
    
    if (props.mode === 'add') {
      await domainApi.create(data)
      ElMessage.success('域名创建成功')
    } else if (props.domain) {
      await domainApi.update(props.domain.id, data)
      ElMessage.success('域名更新成功')
    }
    
    emit('success')
    emit('update:visible', false)
  } catch (error) {
    console.error('操作失败:', error)
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
  formData.name = ''
  formData.registrar_type = ''
  formData.group_id = undefined
  formData.auto_renew = false
  formData.expires_at = null
  formData.remark = ''
  
  formRef.value?.clearValidate()
}

// 监听props变化
watch(() => props.visible, (visible) => {
  if (visible) {
    if (props.mode === 'add') {
      resetForm()
    } else if (props.domain) {
      // 编辑模式，填充表单数据
      formData.name = props.domain.name
      formData.registrar_type = props.domain.registrar_type
      formData.group_id = props.domain.group_id
      formData.auto_renew = props.domain.auto_renew || false
      formData.expires_at = props.domain.expires_at ? new Date(props.domain.expires_at) : null
      formData.remark = props.domain.remark || ''
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
