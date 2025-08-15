<template>
  <el-dialog
    v-model="dialogVisible"
    :title="modalTitle"
    width="500px"
    :close-on-click-modal="false"
    @close="handleCancel"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="标签名称" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入标签名称"
        />
      </el-form-item>

      <el-form-item label="标签颜色" prop="color">
        <el-color-picker
          v-model="formData.color"
          :predefine="predefineColors"
        />
        <span style="margin-left: 8px; color: #666; font-size: 12px;">
          标签在界面中的显示颜色
        </span>
      </el-form-item>

      <el-form-item label="背景颜色" prop="background_color">
        <el-color-picker
          v-model="formData.background_color"
          :predefine="predefineBackgroundColors"
        />
        <span style="margin-left: 8px; color: #666; font-size: 12px;">
          标签背景颜色（可选）
        </span>
      </el-form-item>

      <el-form-item label="排序" prop="sort">
        <el-input-number
          v-model="formData.sort"
          :min="0"
          :max="999"
          style="width: 200px"
        />
        <span style="margin-left: 8px; color: #666; font-size: 12px;">
          数字越小越靠前
        </span>
      </el-form-item>

      <el-form-item label="描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="请输入标签描述"
        />
      </el-form-item>

      <!-- 预览效果 -->
      <el-form-item label="预览效果">
        <el-tag
          :color="formData.background_color"
          :style="{ 
            color: formData.color,
            backgroundColor: formData.background_color 
          }"
        >
          {{ formData.name || '标签名称' }}
        </el-tag>
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
import { tagApi } from '@/api/dns/tag'
import type { DnsTag } from '@/types/dns'

interface Props {
  visible: boolean
  mode?: 'add' | 'edit' | 'view'
  tag?: DnsTag | null
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'add',
  tag: null
})

const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  name: '',
  color: '#FFFFFF',
  background_color: '#409EFF',
  sort: 0,
  description: ''
})

// 预定义颜色
const predefineColors = [
  '#FFFFFF',
  '#000000',
  '#409EFF',
  '#67C23A',
  '#E6A23C',
  '#F56C6C',
  '#909399',
  '#FF6B6B',
  '#4ECDC4',
  '#45B7D1',
  '#96CEB4',
  '#FFEAA7'
]

const predefineBackgroundColors = [
  '#409EFF',
  '#67C23A',
  '#E6A23C',
  '#F56C6C',
  '#909399',
  '#FF6B6B',
  '#4ECDC4',
  '#45B7D1',
  '#96CEB4',
  '#FFEAA7',
  '#DDA0DD',
  '#98D8C8'
]

// 对话框可见性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 模态框标题
const modalTitle = computed(() => {
  const titles = {
    add: '添加标签',
    edit: '编辑标签',
    view: '查看标签'
  }
  return titles[props.mode]
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入标签名称', trigger: 'blur' },
    { min: 1, max: 20, message: '标签名称长度在 1 到 20 个字符', trigger: 'blur' }
  ],
  color: [
    { required: true, message: '请选择标签颜色', trigger: 'change' }
  ],
  background_color: [
    { required: true, message: '请选择背景颜色', trigger: 'change' }
  ],
  sort: [
    { required: true, message: '请输入排序值', trigger: 'blur' },
    { type: 'number', min: 0, max: 999, message: '排序值必须在 0-999 之间', trigger: 'blur' }
  ]
}

// 事件处理
const handleOk = async () => {
  try {
    await formRef.value?.validate()
    
    loading.value = true
    
    const data = {
      name: formData.name,
      color: formData.color,
      background_color: formData.background_color,
      sort: formData.sort,
      description: formData.description
    }
    
    if (props.mode === 'add') {
      await tagApi.create(data)
      ElMessage.success('标签创建成功')
    } else if (props.tag) {
      await tagApi.update(props.tag.id, data)
      ElMessage.success('标签更新成功')
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
  formData.color = '#FFFFFF'
  formData.background_color = '#409EFF'
  formData.sort = 0
  formData.description = ''
  
  formRef.value?.clearValidate()
}

// 监听props变化
watch(() => props.visible, (visible) => {
  if (visible) {
    if (props.mode === 'add') {
      resetForm()
    } else if (props.tag) {
      // 编辑模式，填充表单数据
      formData.name = props.tag.name
      formData.color = props.tag.color || '#FFFFFF'
      formData.background_color = props.tag.background_color || '#409EFF'
      formData.sort = props.tag.sort || 0
      formData.description = props.tag.description || ''
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
