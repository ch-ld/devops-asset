<template>
  <el-dialog
    v-model="dialogVisible"
    title="创建通知"
    width="600px"
    :close-on-click-modal="false"
    @close="handleCancel"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="通知标题" prop="title">
        <el-input
          v-model="formData.title"
          placeholder="请输入通知标题"
          maxlength="100"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="通知级别" prop="level">
        <el-select v-model="formData.level" placeholder="请选择通知级别">
          <el-option label="信息" value="info" />
          <el-option label="成功" value="success" />
          <el-option label="警告" value="warning" />
          <el-option label="错误" value="error" />
        </el-select>
      </el-form-item>

      <el-form-item label="通知类别" prop="category">
        <el-select v-model="formData.category" placeholder="请选择通知类别">
          <el-option label="证书" value="cert" />
          <el-option label="监控" value="monitor" />
          <el-option label="DNS" value="dns" />
          <el-option label="系统" value="system" />
          <el-option label="安全" value="security" />
        </el-select>
      </el-form-item>

      <el-form-item label="通知内容" prop="content">
        <el-input
          v-model="formData.content"
          type="textarea"
          :rows="4"
          placeholder="请输入通知内容"
          maxlength="500"
          show-word-limit
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleOk">
          创建
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElForm } from 'element-plus'

interface Props {
  visible: boolean
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  visible: false
})

const emit = defineEmits<Emits>()

const formRef = ref<InstanceType<typeof ElForm>>()
const loading = ref(false)

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 表单数据
const formData = reactive({
  title: '',
  level: 'info',
  category: 'system',
  content: ''
})

// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入通知标题', trigger: 'blur' },
    { min: 2, max: 100, message: '标题长度为2-100个字符', trigger: 'blur' }
  ],
  level: [
    { required: true, message: '请选择通知级别', trigger: 'change' }
  ],
  category: [
    { required: true, message: '请选择通知类别', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请输入通知内容', trigger: 'blur' },
    { min: 5, max: 500, message: '内容长度为5-500个字符', trigger: 'blur' }
  ]
}

// 重置表单
const resetForm = () => {
  Object.assign(formData, {
    title: '',
    level: 'info',
    category: 'system',
    content: ''
  })
  formRef.value?.resetFields()
}

// 处理确认
const handleOk = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true

    // TODO: 调用API创建通知
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    ElMessage.success('通知创建成功')
    emit('success')
    handleCancel()
  } catch (error: any) {
    if (error.fields) {
      // 表单验证错误
      return
    }
    ElMessage.error('通知创建失败')
  } finally {
    loading.value = false
  }
}

// 处理取消
const handleCancel = () => {
  dialogVisible.value = false
  resetForm()
}
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
