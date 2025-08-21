<template>
  <el-dialog
    :model-value="visible"
    :title="isEdit ? '编辑分组' : '新建分组'"
    width="600px"
    :before-close="handleClose"
    class="domain-group-modal"
  >
    <el-form 
      ref="formRef"
      :model="form" 
      :rules="rules"
      label-width="100px" 
      class="group-form"
    >
      <el-form-item label="分组名称" prop="name">
        <el-input 
          v-model="form.name" 
          placeholder="请输入分组名称"
          maxlength="50"
          show-word-limit
        />
      </el-form-item>
      
      <el-form-item label="描述信息" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="请输入分组描述信息"
          maxlength="200"
          show-word-limit
        />
      </el-form-item>
      
      <el-form-item label="父分组" prop="parent_id">
        <el-tree-select
          v-model="form.parent_id"
          :data="treeData"
          :props="treeProps"
          placeholder="选择父分组（可选）"
          check-strictly
          clearable
          style="width: 100%"
          :render-after-expand="false"
        />
      </el-form-item>
      
      <el-form-item label="排序权重" prop="sort">
        <el-input-number 
          v-model="form.sort" 
          :min="0" 
          :max="999" 
          placeholder="数值越大排序越靠前"
          style="width: 200px"
        />
        <span class="form-help">数值越大排序越靠前</span>
      </el-form-item>
      
      <el-form-item label="分组颜色" prop="color">
        <el-color-picker v-model="form.color" :predefine="predefineColors" />
        <span class="form-help">用于区分不同分组的颜色标识</span>
      </el-form-item>
      
      <el-form-item label="分组状态" prop="status">
        <el-radio-group v-model="form.status">
          <el-radio label="active">启用</el-radio>
          <el-radio label="inactive">禁用</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose" class="cancel-btn">
          取消
        </el-button>
        <el-button 
          type="primary" 
          @click="handleSubmit"
          :loading="submitLoading"
          class="submit-btn"
        >
          {{ isEdit ? '保存修改' : '创建分组' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, nextTick } from 'vue'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import { domainGroupApi } from '@/api/dns/domainGroup'
import type { DomainGroup } from '@/types/dns'

interface Props {
  visible: boolean
  group?: DomainGroup | null
  parentGroups: DomainGroup[]
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 响应式数据
const formRef = ref<FormInstance>()
const submitLoading = ref(false)

const form = reactive({
  name: '',
  description: '',
  parent_id: null as number | null,
  sort: 0,
  color: '#3b82f6',
  status: 'active' as 'active' | 'inactive'
})

// 预定义颜色
const predefineColors = [
  '#3b82f6', '#10b981', '#f59e0b', '#ff4d4f', '#8b5cf6',
  '#06b6d4', '#84cc16', '#f97316', '#ec4899', '#6366f1'
]

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: '请输入分组名称', trigger: 'blur' },
    { min: 2, max: 50, message: '分组名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 200, message: '描述信息长度不能超过 200 个字符', trigger: 'blur' }
  ],
  sort: [
    { type: 'number', min: 0, max: 999, message: '排序权重必须在 0-999 之间', trigger: 'blur' }
  ]
}

// 计算属性
const isEdit = computed(() => !!(props.group && props.group.id))

const treeData = computed(() => {
  // 构建树形数据，排除当前编辑的分组及其子分组
  const buildTree = (groups: DomainGroup[]): any[] => {
    return groups
      .filter(group => {
        // 如果是编辑模式，排除当前分组及其子分组
        if (isEdit.value && props.group) {
          return group.id !== props.group.id && !isDescendantOf(group, props.group)
        }
        return true
      })
      .map(group => ({
        value: group.id,
        label: group.name,
        children: group.children ? buildTree(group.children) : []
      }))
  }
  
  return buildTree(props.parentGroups)
})

const treeProps = {
  children: 'children',
  label: 'label',
  value: 'value'
}

// 方法
const isDescendantOf = (node: DomainGroup, ancestor: DomainGroup): boolean => {
  if (!ancestor.children) return false
  return ancestor.children.some(child => 
    child.id === node.id || isDescendantOf(node, child)
  )
}

const resetForm = () => {
  Object.assign(form, {
    name: '',
    description: '',
    parent_id: null,
    sort: 0,
    color: '#3b82f6',
    status: 'active'
  })
  formRef.value?.clearValidate()
}

const loadGroupData = () => {
  if (props.group && props.group.id) {
    Object.assign(form, {
      name: props.group.name || '',
      description: props.group.description || '',
      parent_id: props.group.parent_id || null,
      sort: props.group.sort || 0,
      color: props.group.color || '#3b82f6',
      status: props.group.status || 'active'
    })
  } else {
    resetForm()
  }
}

const handleClose = () => {
  emit('update:visible', false)
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitLoading.value = true
    
    const submitData = { ...form }
    
    if (isEdit.value && props.group?.id) {
      await domainGroupApi.update(props.group.id, submitData)
      ElMessage.success('修改成功')
    } else {
      await domainGroupApi.create(submitData)
      ElMessage.success('创建成功')
    }
    
    emit('success')
  } catch (error) {
    if (error !== 'validation-failed') {
      ElMessage.error(isEdit.value ? '修改失败' : '创建失败')
    }
  } finally {
    submitLoading.value = false
  }
}

// 监听弹窗显示状态
watch(() => props.visible, (visible) => {
  if (visible) {
    nextTick(() => {
      loadGroupData()
    })
  }
})

// 监听分组数据变化
watch(() => props.group, () => {
  if (props.visible) {
    loadGroupData()
  }
}, { deep: true })
</script>

<style scoped lang="scss">
.domain-group-modal {
  :deep(.el-dialog) {
    border-radius: 16px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  }
  
  :deep(.el-dialog__header) {
    padding: 24px 24px 0;
    border-bottom: 1px solid #f0f2f5;
    
    .el-dialog__title {
      font-size: 18px;
      font-weight: 600;
      color: #1f2937;
    }
  }
  
  :deep(.el-dialog__body) {
    padding: 24px;
  }
  
  :deep(.el-dialog__footer) {
    padding: 0 24px 24px;
    border-top: 1px solid #f0f2f5;
  }
}

.group-form {
  .el-form-item {
    margin-bottom: 24px;
  }
  
  :deep(.el-form-item__label) {
    font-weight: 500;
    color: #374151;
  }
  
  :deep(.el-input__wrapper) {
    border-radius: 8px;
    transition: all 0.3s ease;
  }
  
  :deep(.el-textarea__inner) {
    border-radius: 8px;
    transition: all 0.3s ease;
  }
  
  :deep(.el-input-number) {
    .el-input__wrapper {
      border-radius: 8px;
    }
  }
  
  :deep(.el-tree-select) {
    .el-input__wrapper {
      border-radius: 8px;
    }
  }
  
  .form-help {
    margin-left: 12px;
    font-size: 12px;
    color: #6b7280;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 20px;
  
  .cancel-btn {
    padding: 10px 20px;
    border-radius: 8px;
    border: 1px solid #d1d5db;
    background: #f9fafb;
    color: #374151;
    font-weight: 500;
    transition: all 0.3s ease;
    
    &:hover {
      background: #f3f4f6;
      border-color: #9ca3af;
    }
  }
  
  .submit-btn {
    padding: 10px 20px;
    border-radius: 8px;
    background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
    border: none;
    font-weight: 600;
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 8px 20px rgba(59, 130, 246, 0.4);
    }
  }
}

/* 表单验证错误样式 */
.group-form :deep(.el-form-item.is-error) {
  .el-input__wrapper {
    border-color: #ef4444;
    box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.1);
  }
  
  .el-textarea__inner {
    border-color: #ef4444;
    box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.1);
  }
}

/* 表单聚焦样式 */
.group-form :deep(.el-input__wrapper:focus-within) {
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1);
}

.group-form :deep(.el-textarea__inner:focus) {
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1);
}
</style>
