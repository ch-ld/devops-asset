<template>
  <a-modal
    :visible="visible"
    :title="modalTitle"
    :width="800"
    :confirm-loading="confirmLoading"
    :mask-closable="false"
    @ok="handleOk"
    @cancel="handleCancel"
  >
    <a-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      :label-col="{ span: 6 }"
      :wrapper-col="{ span: 16 }"
    >
      <a-form-item label="域名" name="name">
        <a-input
          v-model:value="formData.name"
          placeholder="请输入域名，如：example.com"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <a-form-item label="状态" name="status">
        <a-select
          v-model:value="formData.status"
          placeholder="请选择状态"
          :disabled="mode === 'view'"
        >
          <a-select-option value="active">正常</a-select-option>
          <a-select-option value="inactive">停用</a-select-option>
          <a-select-option value="expired">已过期</a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="注册商类型" name="registrar_type">
        <a-select
          v-model:value="formData.registrar_type"
          placeholder="请选择注册商类型"
          :disabled="mode === 'view'"
        >
          <a-select-option value="godaddy">GoDaddy</a-select-option>
          <a-select-option value="aliyun">阿里云</a-select-option>
          <a-select-option value="tencent">腾讯云</a-select-option>
          <a-select-option value="dnspod">DNSPod</a-select-option>
          <a-select-option value="cloudflare">Cloudflare</a-select-option>
          <a-select-option value="other">其他</a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="过期时间" name="expires_at">
        <a-date-picker
          v-model:value="formData.expires_at"
          placeholder="请选择过期时间"
          style="width: 100%"
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <a-form-item label="自动续费" name="auto_renew">
        <a-switch
          v-model:checked="formData.auto_renew"
          :disabled="mode === 'view'"
        />
        <span class="form-help-text">开启后将在到期前自动续费</span>
      </a-form-item>

      <a-form-item label="域名分组" name="group_id">
        <a-tree-select
          v-model:value="formData.group_id"
          :tree-data="groupTreeData"
          placeholder="请选择域名分组"
          allow-clear
          tree-default-expand-all
          :disabled="mode === 'view'"
        />
      </a-form-item>

      <a-form-item label="配置信息" name="configuration">
        <a-textarea
          v-model:value="configurationText"
          placeholder="请输入JSON格式的配置信息，如：{&quot;ns&quot;:[&quot;ns1.example.com&quot;,&quot;ns2.example.com&quot;]}"
          :rows="4"
          :disabled="mode === 'view'"
        />
        <div class="form-help-text">
          配置信息格式示例：
          <br />
          {"ns":["ns1.example.com","ns2.example.com"],"ttl":600}
        </div>
      </a-form-item>

      <a-form-item label="备注" name="remark">
        <a-textarea
          v-model:value="formData.remark"
          placeholder="请输入备注信息"
          :rows="3"
          :disabled="mode === 'view'"
        />
      </a-form-item>
    </a-form>

    <template #footer>
      <div v-if="mode === 'view'">
        <a-button @click="handleCancel">关闭</a-button>
      </div>
      <div v-else>
        <a-button @click="handleCancel">取消</a-button>
        <a-button type="primary" :loading="confirmLoading" @click="handleOk">
          {{ mode === 'add' ? '创建' : '更新' }}
        </a-button>
      </div>
    </template>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { message } from 'ant-design-vue'
import type { FormInstance } from 'ant-design-vue'
import { domainApi } from '@/api/dns/domain'
import type { Domain, DomainGroup } from '@/types/dns'
import dayjs from 'dayjs'

interface Props {
  visible: boolean
  mode: 'add' | 'edit' | 'view'
  domain?: Domain | null
  groupOptions: DomainGroup[]
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  domain: null
})

const emit = defineEmits<Emits>()

// 响应式数据
const formRef = ref<FormInstance>()
const confirmLoading = ref(false)
const configurationText = ref('')

// 表单数据
const formData = reactive({
  name: '',
  status: 'active',
  registrar_type: '',
  registrar_id: undefined,
  expires_at: null,
  auto_renew: false,
  group_id: undefined,
  configuration: {},
  remark: ''
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
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ],
  registrar_type: [
    { required: true, message: '请选择注册商类型', trigger: 'change' }
  ]
}

// 计算属性
const modalTitle = computed(() => {
  const titleMap = {
    add: '添加域名',
    edit: '编辑域名',
    view: '查看域名'
  }
  return titleMap[props.mode]
})

const groupTreeData = computed(() => {
  return buildTreeData(props.groupOptions)
})

// 构建树形数据
const buildTreeData = (groups: DomainGroup[]) => {
  return groups.map(group => ({
    title: group.name,
    value: group.id,
    key: group.id,
    children: group.children ? buildTreeData(group.children) : undefined
  }))
}

// 监听弹窗显示状态
watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      resetForm()
      if (props.domain && props.mode !== 'add') {
        loadDomainData()
      }
    }
  }
)

// 监听配置文本变化
watch(
  () => configurationText.value,
  (value) => {
    try {
      if (value.trim()) {
        formData.configuration = JSON.parse(value)
      } else {
        formData.configuration = {}
      }
    } catch (error) {
      // JSON格式错误时不更新formData.configuration
    }
  }
)

// 重置表单
const resetForm = () => {
  formRef.value?.resetFields()
  Object.assign(formData, {
    name: '',
    status: 'active',
    registrar_type: '',
    registrar_id: undefined,
    expires_at: null,
    auto_renew: false,
    group_id: undefined,
    configuration: {},
    remark: ''
  })
  configurationText.value = ''
}

// 加载域名数据
const loadDomainData = () => {
  if (!props.domain) return

  Object.assign(formData, {
    name: props.domain.name || '',
    status: props.domain.status || 'active',
    registrar_type: props.domain.registrar_type || '',
    registrar_id: props.domain.registrar_id,
    expires_at: props.domain.expires_at ? dayjs(props.domain.expires_at) : null,
    auto_renew: props.domain.auto_renew || false,
    group_id: props.domain.group_id,
    configuration: props.domain.configuration || {},
    remark: props.domain.remark || ''
  })

  // 设置配置文本
  if (props.domain.configuration) {
    try {
      configurationText.value = JSON.stringify(props.domain.configuration, null, 2)
    } catch (error) {
      configurationText.value = ''
    }
  }
}

// 处理确认
const handleOk = async () => {
  if (props.mode === 'view') {
    handleCancel()
    return
  }

  try {
    await formRef.value?.validate()
    
    // 验证配置JSON格式
    if (configurationText.value.trim()) {
      try {
        JSON.parse(configurationText.value)
      } catch (error) {
        message.error('配置信息格式错误，请输入有效的JSON格式')
        return
      }
    }

    confirmLoading.value = true

    const submitData = {
      ...formData,
      expires_at: formData.expires_at ? formData.expires_at.format('YYYY-MM-DD HH:mm:ss') : null
    }

    if (props.mode === 'add') {
      await domainApi.create(submitData)
      message.success('域名创建成功')
    } else {
      await domainApi.update(props.domain!.id, submitData)
      message.success('域名更新成功')
    }

    emit('success')
  } catch (error) {
    if (error.errorFields) {
      // 表单验证错误
      return
    }
    message.error(props.mode === 'add' ? '创建失败' : '更新失败')
  } finally {
    confirmLoading.value = false
  }
}

// 处理取消
const handleCancel = () => {
  emit('update:visible', false)
}
</script>

<style scoped lang="scss">
.form-help-text {
  margin-left: 8px;
  color: #8c8c8c;
  font-size: 12px;
  line-height: 1.5;
}

:deep(.ant-form-item-explain) {
  font-size: 12px;
}
</style>
