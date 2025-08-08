<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑标签' : '新建标签'"
    width="500px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="80px"
    >
      <el-form-item label="标签名称" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入标签名称"
          maxlength="20"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="标签颜色" prop="color">
        <div class="color-picker-container">
          <el-color-picker 
            v-model="formData.color" 
            show-alpha 
            :predefine="predefineColors"
          />
          <el-tag :color="formData.color" class="preview-tag">
            {{ formData.name || '预览' }}
          </el-tag>
        </div>
      </el-form-item>

      <el-form-item label="标签描述">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="请输入标签描述"
          maxlength="100"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="排序">
        <el-input-number
          v-model="formData.sort"
          :min="0"
          :max="999"
          placeholder="排序值"
        />
        <div class="form-tip">
          数值越小越靠前，默认为0
        </div>
      </el-form-item>

      <el-form-item label="状态">
        <el-switch
          v-model="formData.status"
          :active-value="1"
          :inactive-value="0"
          active-text="启用"
          inactive-text="禁用"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">
          {{ isEdit ? '保存' : '创建' }}
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
  editData?: any
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success', data: any): void
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  editData: null
})

const emit = defineEmits<Emits>()

const formRef = ref<InstanceType<typeof ElForm>>()
const saving = ref(false)

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const isEdit = computed(() => !!props.editData?.id)

// 预定义颜色
const predefineColors = [
  '#ff4500',
  '#ff8c00',
  '#ffd700',
  '#90ee90',
  '#00ced1',
  '#1e90ff',
  '#c71585',
  '#ff69b4',
  '#ba55d3',
  '#9370db',
  '#3cb371',
  '#32cd32',
  '#f0e68c',
  '#dda0dd',
  '#87ceeb'
]

// 表单数据
const formData = reactive({
  name: '',
  color: '#1e90ff',
  description: '',
  sort: 0,
  status: 1
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入标签名称', trigger: 'blur' },
    { min: 2, max: 20, message: '标签名称长度为2-20个字符', trigger: 'blur' }
  ],
  color: [
    { required: true, message: '请选择标签颜色', trigger: 'change' }
  ]
}

// 重置表单
const resetForm = () => {
  Object.assign(formData, {
    name: '',
    color: '#1e90ff',
    description: '',
    sort: 0,
    status: 1
  })
  formRef.value?.resetFields()
}

// 处理关闭
const handleClose = () => {
  dialogVisible.value = false
  resetForm()
}

// 处理保存
const handleSave = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    saving.value = true

    const data = { ...formData }
    
    // TODO: 调用API创建或更新标签
    if (isEdit.value) {
      // await tagApi.update(props.editData.id, data)
      ElMessage.success('标签更新成功')
    } else {
      // const response = await tagApi.create(data)
      ElMessage.success('标签创建成功')
      emit('success', data) // 临时返回表单数据
    }

    handleClose()
  } catch (error: any) {
    if (error.fields) {
      // 表单验证错误
      return
    }
    ElMessage.error(isEdit.value ? '标签更新失败' : '标签创建失败')
  } finally {
    saving.value = false
  }
}

// 监听编辑数据变化
watch(
  () => props.editData,
  (newVal) => {
    if (newVal) {
      Object.assign(formData, {
        name: newVal.name || '',
        color: newVal.color || '#1e90ff',
        description: newVal.description || '',
        sort: newVal.sort || 0,
        status: newVal.status ?? 1
      })
    }
  },
  { immediate: true }
)

// 监听弹窗显示状态
watch(
  () => props.visible,
  (visible) => {
    if (visible && !isEdit.value) {
      resetForm()
    }
  }
)
</script>

<style scoped>
.color-picker-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.preview-tag {
  min-width: 60px;
  text-align: center;
}

.form-tip {
  font-size: 12px;
  color: var(--el-color-info);
  margin-top: 4px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style> 
