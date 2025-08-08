<template>
  <el-dialog
    v-model="dialogVisible"
    :title="modalTitle"
    width="600px"
    :close-on-click-modal="false"
    @close="handleCancel"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item label="分组名称" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入分组名称"
          maxlength="50"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="上级分组" prop="parent_id">
        <el-select
          v-model="formData.parent_id"
          placeholder="请选择上级分组"
          style="width: 100%"
          clearable
        >
          <el-option label="无上级分组" :value="null" />
          <el-option
            v-for="group in parentGroups"
            :key="group.id"
            :label="group.name"
            :value="group.id"
            :disabled="group.id === formData.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="描述信息" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="请输入分组描述"
          maxlength="500"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="排序" prop="sort">
        <el-input-number
          v-model="formData.sort"
          :min="0"
          :max="9999"
          style="width: 200px"
        />
        <span style="margin-left: 8px; color: #666; font-size: 12px;">
          数字越小排序越靠前
        </span>
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
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { domainGroupApi } from '@/api/dns/domainGroup'
import type { DomainGroup } from '@/types/dns'

interface Props {
  visible: boolean
  mode?: 'add' | 'edit'
  group?: DomainGroup | null
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'add',
  group: null
})

const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const formRef = ref()
const parentGroups = ref<DomainGroup[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  parent_id: null as number | null,
  description: '',
  sort: 0
})

// 对话框可见性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 模态框标题
const modalTitle = computed(() => {
  const titles = {
    add: '添加分组',
    edit: '编辑分组'
  }
  return titles[props.mode]
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入分组名称', trigger: 'blur' },
    { min: 2, max: 50, message: '分组名称长度为2-50个字符', trigger: 'blur' }
  ]
}

// 获取父级分组列表
const fetchParentGroups = async () => {
  try {
    const response = await domainGroupApi.list({ page: 1, size: 1000 })
    parentGroups.value = response.list || []
  } catch (error) {
    console.error('获取父级分组列表失败:', error)
  }
}

// 事件处理
const handleOk = async () => {
  try {
    await formRef.value?.validate()
    
    loading.value = true
    
    const data = {
      name: formData.name,
      parent_id: formData.parent_id,
      description: formData.description,
      sort: formData.sort
    }
    
    if (props.mode === 'add') {
      await domainGroupApi.create(data)
      ElMessage.success('分组创建成功')
    } else if (props.group) {
      await domainGroupApi.update(props.group.id, data)
      ElMessage.success('分组更新成功')
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
  formData.id = undefined
  formData.name = ''
  formData.parent_id = null
  formData.description = ''
  formData.sort = 0
  
  formRef.value?.clearValidate()
}

// 监听props变化
watch(() => props.visible, (visible) => {
  if (visible) {
    fetchParentGroups()
    if (props.mode === 'add') {
      resetForm()
    } else if (props.group) {
      // 编辑模式，填充表单数据
      formData.id = props.group.id
      formData.name = props.group.name
      formData.parent_id = props.group.parent_id || null
      formData.description = props.group.description || ''
      formData.sort = props.group.sort || 0
    }
  }
})

// 初始化
onMounted(() => {
  fetchParentGroups()
})
</script>

<style scoped lang="scss">
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style> 
