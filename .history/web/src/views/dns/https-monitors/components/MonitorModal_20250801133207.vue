<template>
  <el-dialog
    v-model="dialogVisible"
    :title="modalTitle"
    width="900px"
    :close-on-click-modal="false"
    @close="handleCancel"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="120px"
      class="monitor-form"
    >
      <!-- 基本信息 -->
      <el-divider content-position="left">基本信息</el-divider>
      
      <el-form-item label="监控名称" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入监控名称"
          :disabled="mode === 'view'"
        />
      </el-form-item>

      <el-form-item label="监控URL" prop="url">
        <el-input
          v-model="formData.url"
          placeholder="请输入完整的HTTPS URL，如：https://example.com"
          :disabled="mode === 'view'"
        >
          <template #prepend>
            <el-select v-model="urlProtocol" style="width: 80px" :disabled="mode === 'view'">
              <el-option label="HTTPS" value="https" />
              <el-option label="HTTP" value="http" />
            </el-select>
          </template>
          <template #append>
            <el-button 
              v-if="formData.url && mode !== 'view'" 
              text
              :loading="testLoading"
              @click="testUrl"
            >
              测试
            </el-button>
          </template>
        </el-input>
      </el-form-item>

      <!-- 监控配置 -->
      <el-divider content-position="left">监控配置</el-divider>
      
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="检查间隔" prop="checkInterval">
            <el-input-number
              v-model="formData.checkInterval"
              :min="1"
              :max="1440"
              placeholder="分钟"
              :disabled="mode === 'view'"
              style="width: 100%"
            />
            <span class="input-suffix">分钟</span>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="超时时间" prop="timeout">
            <el-input-number
              v-model="formData.timeout"
              :min="5"
              :max="300"
              placeholder="秒"
              :disabled="mode === 'view'"
              style="width: 100%"
            />
            <span class="input-suffix">秒</span>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="告警阈值" prop="alertThreshold">
            <el-input-number
              v-model="formData.alertThreshold"
              :min="1000"
              :max="30000"
              placeholder="毫秒"
              :disabled="mode === 'view'"
              style="width: 100%"
            />
            <span class="input-suffix">ms</span>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="重试次数" prop="retryCount">
            <el-input-number
              v-model="formData.retryCount"
              :min="0"
              :max="10"
              :disabled="mode === 'view'"
              style="width: 100%"
            />
            <span class="input-suffix">次</span>
          </el-form-item>
        </el-col>
      </el-row>

      <!-- SSL证书监控 -->
      <el-divider content-position="left">SSL证书监控</el-divider>
      
      <el-form-item label="证书监控">
        <el-switch
          v-model="formData.enableSslCheck"
          active-text="开启"
          inactive-text="关闭"
          :disabled="mode === 'view'"
        />
      </el-form-item>

      <div v-if="formData.enableSslCheck">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="过期提醒" prop="sslExpiryDays">
              <el-input-number
                v-model="formData.sslExpiryDays"
                :min="1"
                :max="365"
                :disabled="mode === 'view'"
                style="width: 100%"
              />
              <span class="input-suffix">天前提醒</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="证书验证">
              <el-switch
                v-model="formData.verifySsl"
                active-text="验证"
                inactive-text="忽略"
                :disabled="mode === 'view'"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <!-- 通知配置 -->
      <el-divider content-position="left">通知配置</el-divider>
      
      <el-form-item label="启用通知">
        <el-switch
          v-model="formData.enableNotification"
          active-text="开启"
          inactive-text="关闭"
          :disabled="mode === 'view'"
        />
      </el-form-item>

      <div v-if="formData.enableNotification">
        <el-form-item label="通知方式" prop="notificationChannels">
          <el-checkbox-group v-model="formData.notificationChannels" :disabled="mode === 'view'">
            <el-checkbox label="email">邮件</el-checkbox>
            <el-checkbox label="sms">短信</el-checkbox>
            <el-checkbox label="webhook">Webhook</el-checkbox>
            <el-checkbox label="dingtalk">钉钉</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="通知邮箱" prop="notificationEmail" v-if="formData.notificationChannels.includes('email')">
          <el-input
            v-model="formData.notificationEmail"
            placeholder="请输入邮箱地址，多个邮箱用逗号分隔"
            :disabled="mode === 'view'"
          />
        </el-form-item>

        <el-form-item label="Webhook URL" prop="webhookUrl" v-if="formData.notificationChannels.includes('webhook')">
          <el-input
            v-model="formData.webhookUrl"
            placeholder="请输入Webhook URL"
            :disabled="mode === 'view'"
          />
        </el-form-item>
      </div>

      <!-- 其他设置 -->
      <el-divider content-position="left">其他设置</el-divider>
      
      <el-form-item label="状态">
        <el-switch
          v-model="formData.enabled"
          active-text="启用"
          inactive-text="禁用"
          :disabled="mode === 'view'"
        />
      </el-form-item>

      <el-form-item label="备注">
        <el-input
          v-model="formData.remark"
          type="textarea"
          :rows="3"
          placeholder="请输入备注信息"
          :disabled="mode === 'view'"
          maxlength="500"
          show-word-limit
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleCancel">{{ mode === 'view' ? '关闭' : '取消' }}</el-button>
        <el-button 
          v-if="mode !== 'view'" 
          type="primary" 
          :loading="loading" 
          @click="handleOk"
        >
          {{ mode === 'add' ? '创建' : '保存' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, ElForm } from 'element-plus'

interface Props {
  visible: boolean
  mode: 'add' | 'edit' | 'view'
  monitor?: any
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  mode: 'add',
  monitor: null
})

const emit = defineEmits<Emits>()

const formRef = ref<InstanceType<typeof ElForm>>()
const loading = ref(false)
const testLoading = ref(false)
const urlProtocol = ref('https')

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const modalTitle = computed(() => {
  const titleMap = {
    add: '添加监控',
    edit: '编辑监控',
    view: '查看监控'
  }
  return titleMap[props.mode]
})

// 表单数据
const formData = reactive({
  name: '',
  url: '',
  checkInterval: 5,
  timeout: 30,
  alertThreshold: 5000,
  retryCount: 3,
  enableSslCheck: true,
  sslExpiryDays: 30,
  verifySsl: true,
  enableNotification: true,
  notificationChannels: ['email'] as string[],
  notificationEmail: '',
  webhookUrl: '',
  enabled: true,
  remark: ''
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入监控名称', trigger: 'blur' },
    { min: 2, max: 50, message: '监控名称长度为2-50个字符', trigger: 'blur' }
  ],
  url: [
    { required: true, message: '请输入监控URL', trigger: 'blur' },
    { 
      pattern: /^https?:\/\/.+/, 
      message: '请输入有效的URL格式', 
      trigger: 'blur' 
    }
  ],
  checkInterval: [
    { required: true, message: '请设置检查间隔', trigger: 'blur' }
  ],
  timeout: [
    { required: true, message: '请设置超时时间', trigger: 'blur' }
  ],
  notificationEmail: [
    { 
      pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/, 
      message: '请输入有效的邮箱地址', 
      trigger: 'blur' 
    }
  ],
  webhookUrl: [
    { 
      pattern: /^https?:\/\/.+/, 
      message: '请输入有效的Webhook URL', 
      trigger: 'blur' 
    }
  ]
}

// 测试URL连通性
const testUrl = async () => {
  if (!formData.url) {
    ElMessage.warning('请先输入URL')
    return
  }

  try {
    testLoading.value = true
    // TODO: 调用API测试URL
    await new Promise(resolve => setTimeout(resolve, 2000))
    ElMessage.success('URL连通性测试成功')
  } catch (error) {
    ElMessage.error('URL连通性测试失败')
  } finally {
    testLoading.value = false
  }
}

// 重置表单
const resetForm = () => {
  Object.assign(formData, {
    name: '',
    url: '',
    checkInterval: 5,
    timeout: 30,
    alertThreshold: 5000,
    retryCount: 3,
    enableSslCheck: true,
    sslExpiryDays: 30,
    verifySsl: true,
    enableNotification: true,
    notificationChannels: ['email'],
    notificationEmail: '',
    webhookUrl: '',
    enabled: true,
    remark: ''
  })
  formRef.value?.resetFields()
}

// 处理确认
const handleOk = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true

    // TODO: 调用API保存数据
    if (props.mode === 'add') {
      // await monitorApi.create(formData)
      ElMessage.success('监控创建成功')
    } else {
      // await monitorApi.update(props.monitor.id, formData)
      ElMessage.success('监控更新成功')
    }

    emit('success')
    handleCancel()
  } catch (error: any) {
    if (error.fields) {
      // 表单验证错误
      return
    }
    ElMessage.error(props.mode === 'add' ? '监控创建失败' : '监控更新失败')
  } finally {
    loading.value = false
  }
}

// 处理取消
const handleCancel = () => {
  dialogVisible.value = false
  resetForm()
}

// 监听编辑数据变化
watch(
  () => props.monitor,
  (newVal) => {
    if (newVal && props.mode !== 'add') {
      Object.assign(formData, {
        name: newVal.name || '',
        url: newVal.url || '',
        checkInterval: newVal.checkInterval || 5,
        timeout: newVal.timeout || 30,
        alertThreshold: newVal.alertThreshold || 5000,
        retryCount: newVal.retryCount || 3,
        enableSslCheck: newVal.enableSslCheck ?? true,
        sslExpiryDays: newVal.sslExpiryDays || 30,
        verifySsl: newVal.verifySsl ?? true,
        enableNotification: newVal.enableNotification ?? true,
        notificationChannels: newVal.notificationChannels || ['email'],
        notificationEmail: newVal.notificationEmail || '',
        webhookUrl: newVal.webhookUrl || '',
        enabled: newVal.enabled ?? true,
        remark: newVal.remark || ''
      })
      
      // 从URL中提取协议
      if (newVal.url) {
        urlProtocol.value = newVal.url.startsWith('https://') ? 'https' : 'http'
      }
    }
  },
  { immediate: true }
)

// 监听弹窗显示状态
watch(
  () => props.visible,
  (visible) => {
    if (visible && props.mode === 'add') {
      resetForm()
    }
  }
)
</script>

<style scoped>
.monitor-form {
  max-height: 600px;
  overflow-y: auto;
}

.input-suffix {
  margin-left: 8px;
  color: var(--el-text-color-regular);
  font-size: 14px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
