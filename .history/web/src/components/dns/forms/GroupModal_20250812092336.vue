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
      label-width="100px"
    >
      <el-form-item label="分组名称" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入分组名称"
        />
      </el-form-item>

      <el-form-item label="父分组" prop="parent_id">
        <el-tree-select
          v-model="formData.parent_id"
          :data="groupTreeData"
          placeholder="请选择父分组（可选）"
          clearable
          check-strictly
          :render-after-expand="false"
          style="width: 100%"
        />
      </el-form-item>

      <el-form-item label="分组颜色" prop="color">
        <el-color-picker
          v-model="formData.color"
          show-alpha
          :predefine="predefineColors"
        />
        <span style="margin-left: 8px; color: #666; font-size: 12px;">
          用于在界面中标识该分组
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
          placeholder="请输入分组描述"
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
import { groupApi } from '@/api/dns/group'
import type { DomainGroup } from '@/types/dns'

interface Props {
  visible: boolean
  mode?: 'add' | 'edit' | 'view'
  group?: DomainGroup | null
  groups?: DomainGroup[]
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'add',
  group: null,
  groups: () => []
})

const emit = defineEmits<Emits>()

// 响应式数据
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  name: '',
  parent_id: undefined as number | undefined,
  color: '#409EFF',
  sort: 0,
  description: ''
})

// 预定义颜色
const predefineColors = [
  '#ff4500',
  '#ff8c00',
  '#ffd700',
  '#90ee90',
  '#00ced1',
  '#1e90ff',
  '#c71585',
  '#409EFF',
  '#67C23A',
  '#E6A23C',
  '#F56C6C',
  '#909399'
]

// 对话框可见性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 模态框标题
const modalTitle = computed(() => {
  const titles = {
    add: '添加分组',
    edit: '编辑分组',
    view: '查看分组'
  }
  return titles[props.mode]
})

// 构建分组树数据
const groupTreeData = computed(() => {
  const buildTree = (groups: DomainGroup[], parentId: number | null = null): any[] => {
    return groups
      .filter(group => group.parent_id === parentId)
      .map(group => ({
        value: group.id,
        label: group.name,
        children: buildTree(groups, group.id)
      }))
  }
  
  return buildTree(props.groups.filter(g => g.id !== props.group?.id))
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入分组名称', trigger: 'blur' },
    { min: 1, max: 50, message: '分组名称长度在 1 到 50 个字符', trigger: 'blur' }
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
      parent_id: formData.parent_id,
      color: formData.color,
      sort: formData.sort,
      description: formData.description
    }
    
    if (props.mode === 'add') {
      await groupApi.create(data)
      ElMessage.success('分组创建成功')
    } else if (props.group) {
      await groupApi.update(props.group.id, data)
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
  formData.name = ''
  formData.parent_id = undefined
  formData.color = '#409EFF'
  formData.sort = 0
  formData.description = ''
  
  formRef.value?.clearValidate()
}

// 监听props变化
watch(() => props.visible, (visible) => {
  if (visible) {
    if (props.mode === 'add') {
      resetForm()
    } else if (props.group) {
      // 编辑模式，填充表单数据
      formData.name = props.group.name
      formData.parent_id = props.group.parent_id
      formData.color = props.group.color || '#409EFF'
      formData.sort = props.group.sort || 0
      formData.description = props.group.description || ''
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
