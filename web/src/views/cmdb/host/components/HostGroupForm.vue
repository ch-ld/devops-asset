<template>
  <div class="host-group-form">
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="100px"
      @submit.prevent
    >
      <el-form-item label="主机组名称" prop="name">
        <el-input v-model="formData.name" placeholder="请输入主机组名称" />
      </el-form-item>

      <el-form-item label="父级主机组" prop="parent_id">
        <div class="custom-select-wrapper">
          <select
            v-model="formData.parent_id"
            class="custom-select"
            @change="handleParentGroupChange"
          >
            <option value="">请选择父级主机组（可选）</option>
            <option
              v-for="group in parentGroupOptions"
              :key="group.value"
              :value="group.value"
            >
              {{ group.label }}
            </option>
          </select>
          <div class="select-arrow">
            <el-icon><ArrowDown /></el-icon>
          </div>
        </div>
        <div style="margin-top: 5px; font-size: 12px; color: #999;">
          调试: 父级选项数量 {{ parentGroupOptions.length }}
        </div>
      </el-form-item>

      <el-form-item label="排序" prop="sort">
        <el-input-number
          v-model="formData.sort"
          :min="0"
          :max="999"
          placeholder="排序值"
        />
      </el-form-item>

      <el-form-item label="描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="请输入主机组描述"
        />
      </el-form-item>

      <!-- 操作按钮 -->
      <div class="form-actions">
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import { 
  createHostGroup, 
  updateHostGroup
} from '@/api/system/host'
import type { HostGroup } from '@/types/api/host'

interface Props {
  group?: HostGroup | null
  groups: HostGroup[]
}

const props = defineProps<Props>()
const emit = defineEmits<{
  success: []
  cancel: []
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const isEdit = computed(() => !!props.group)

// 表单数据
const formData = reactive({
  name: '',
  parent_id: undefined as number | undefined,
  sort: 0,
  description: ''
})

// 表单验证规则
const formRules: FormRules = {
  name: [{ required: true, message: '请输入主机组名称', trigger: 'blur' }]
}

// 分组选项（排除当前编辑的分组）
const groupOptions = computed(() => {
  console.log('🌳 HostGroupForm props.groups:', props.groups)

  if (!props.groups || !Array.isArray(props.groups)) {
    return []
  }

  // 构建树形结构的函数
  const buildTreeData = (groups: HostGroup[], parentId: number | null = null): any[] => {
    return groups
      .filter(group => {
        // 排除当前编辑的分组
        if (props.group && group.id === props.group.id) {
          return false
        }
        // 过滤父级关系
        return group.parent_id === parentId
      })
      .map(group => ({
        value: group.id,
        label: group.name,
        children: buildTreeData(groups, group.id)
      }))
  }

  // 过滤掉"全部主机"选项（id为0）
  const filteredGroups = props.groups.filter(group => group.id !== 0)
  const treeData = buildTreeData(filteredGroups)

  console.log('🌳 HostGroupForm最终树形数据:', treeData)

  // 验证数据结构
  treeData.forEach((item, index) => {
    console.log(`🔍 HostGroupForm选项 ${index}:`, {
      value: item.value,
      label: item.label,
      hasChildren: item.children && item.children.length > 0
    })
  })

  return treeData
})

// 扁平化主机组选项
const parentGroupOptions = computed(() => {
  const flattenOptions = (options: any[], prefix = ''): any[] => {
    const result: any[] = []

    options.forEach(option => {
      const label = prefix ? `${prefix} / ${option.label}` : option.label
      result.push({
        value: option.value,
        label: label
      })

      if (option.children && option.children.length > 0) {
        result.push(...flattenOptions(option.children, label))
      }
    })

    return result
  }

  const flattened = flattenOptions(groupOptions.value)
  console.log('🔍 HostGroupForm扁平化选项:', flattened)
  return flattened
})

// 初始化表单数据
const initFormData = () => {
  if (props.group) {
    Object.assign(formData, {
      name: props.group.name,
      parent_id: props.group.parent_id,
      sort: props.group.sort || 0,
      description: props.group.description || ''
    })
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    if (isEdit.value && props.group) {
      await updateHostGroup(props.group.id, formData)
      ElMessage.success('更新成功')
    } else {
      await createHostGroup(formData)
      ElMessage.success('创建成功')
    }

    emit('success')
  } catch (error) {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  } finally {
    submitting.value = false
  }
}

// 处理父级主机组选择变化
const handleParentGroupChange = (event: Event) => {
  const target = event.target as HTMLSelectElement
  const value = target.value
  console.log('🔄 父级主机组选择变化:', value)

  // 确保数值类型正确
  if (value === '') {
    formData.parent_id = undefined
  } else {
    formData.parent_id = parseInt(value)
  }
}

// 取消
const handleCancel = () => {
  emit('cancel')
}

// 初始化
onMounted(() => {
  initFormData()
})
</script>

<style scoped>
.host-group-form {
  padding: 20px 0;
}

.form-actions {
  margin-top: 24px;
  text-align: center;
}

.form-actions .el-button {
  margin: 0 8px;
}

/* 自定义Select样式 */
.custom-select-wrapper {
  position: relative;
  width: 100%;
}

.custom-select {
  width: 100%;
  height: 40px;
  padding: 8px 32px 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  background-color: #ffffff;
  font-size: 14px;
  color: #606266;
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  cursor: pointer;
  transition: all 0.3s ease;
  outline: none;
}

.custom-select:hover {
  border-color: #c0c4cc;
}

.custom-select:focus {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.custom-select option {
  padding: 8px 12px;
  color: #606266;
  background-color: #ffffff;
}

.custom-select option:hover {
  background-color: #f5f7fa;
}

.select-arrow {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  color: #c0c4cc;
  transition: transform 0.3s ease;
}

.custom-select:focus + .select-arrow {
  color: #409eff;
  transform: translateY(-50%) rotate(180deg);
}
</style>
