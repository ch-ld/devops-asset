<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑分组' : '新建分组'"
    width="500px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="80px"
    >
      <el-form-item label="分组名称" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入分组名称"
          maxlength="50"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="父分组" prop="parentId">
        <el-cascader
          v-model="formData.parentId"
          :options="groupOptions"
          :props="cascaderProps"
          placeholder="请选择父分组（可选）"
          clearable
          filterable
        />
      </el-form-item>

      <el-form-item label="分组描述">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="请输入分组描述"
          maxlength="200"
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
import { domainGroupApi } from '@/api/dns/domainGroup'

interface Props {
  visible: boolean
  editData?: any
  excludeIds?: number[]
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success', data: any): void
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  editData: null,
  excludeIds: () => []
})

const emit = defineEmits<Emits>()

const formRef = ref<InstanceType<typeof ElForm>>()
const saving = ref(false)
const groupOptions = ref([])

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const isEdit = computed(() => !!props.editData?.id)

// 表单数据
const formData = reactive({
  name: '',
  parentId: null as number | null,
  description: '',
  sort: 0,
  status: 1
})

// 级联选择器配置
const cascaderProps = {
  value: 'id',
  label: 'name',
  children: 'children',
  checkStrictly: true,
  emitPath: false
}

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入分组名称', trigger: 'blur' },
    { min: 2, max: 50, message: '分组名称长度为2-50个字符', trigger: 'blur' }
  ]
}

// 获取分组选项
const fetchGroupOptions = async () => {
  try {
    const response = await domainGroupApi.list()
    groupOptions.value = buildTreeOptions(response.data)
  } catch (error) {
    console.error('获取分组列表失败:', error)
  }
}

// 构建树形选项
const buildTreeOptions = (items: any[], parentId = 0): any[] => {
  const result: any[] = []
  items.forEach(item => {
    // 排除当前编辑的分组和其子分组
    if (props.excludeIds.includes(item.id)) {
      return
    }
    
    if (item.parentId === parentId) {
      const children = buildTreeOptions(items, item.id)
      const option: any = {
        id: item.id,
        name: item.name,
        value: item.id,
        label: item.name
      }
      if (children.length > 0) {
        option.children = children
      }
      result.push(option)
    }
  })
  return result
}

// 重置表单
const resetForm = () => {
  Object.assign(formData, {
    name: '',
    parentId: null,
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
    
    if (isEdit.value) {
      await domainGroupApi.update(props.editData.id, data)
      ElMessage.success('分组更新成功')
    } else {
      const response = await domainGroupApi.create(data)
      ElMessage.success('分组创建成功')
      emit('success', response.data)
    }

    handleClose()
  } catch (error: any) {
    if (error.fields) {
      // 表单验证错误
      return
    }
    ElMessage.error(isEdit.value ? '分组更新失败' : '分组创建失败')
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
        parentId: newVal.parentId || null,
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
    if (visible) {
      fetchGroupOptions()
      if (!isEdit.value) {
        resetForm()
      }
    }
  }
)
</script>

<style scoped>
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
