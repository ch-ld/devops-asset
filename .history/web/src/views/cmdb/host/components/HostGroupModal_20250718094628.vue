<template>
  <el-dialog
    :title="isEdit ? '编辑主机组' : '添加主机组'"
    v-model="visible"
    width="600px"
    :close-on-click-modal="false"
  >
    <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
      <el-form-item label="组名称" prop="name">
        <el-input v-model="formData.name" placeholder="请输入主机组名称" />
      </el-form-item>

      <el-form-item label="父级分组" prop="parent_id">
        <el-tree-select
          v-model="formData.parent_id"
          :data="hostGroupOptions"
          placeholder="选择父级分组（可选）"
          clearable
          :default-expand-all="true"
          :render-after-expand="false"
        />
      </el-form-item>

      <el-form-item label="描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="请输入描述信息"
        />
      </el-form-item>

      <el-form-item label="排序" prop="sort">
        <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
      </el-form-item>

      <el-alert
        v-if="parentGroup"
        :title="`将在 ${parentGroup.name} 下创建子分组`"
        type="info"
        show-icon
        class="mb-3"
      />
    </el-form>

    <template #footer>
      <el-button @click="handleCancel">取 消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSubmit">确 定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
  import { ref, reactive, computed } from 'vue'
  import { ElMessage } from 'element-plus'
  import type { FormInstance } from 'element-plus'
  import type { HostGroup } from '@/types/api/host'
  import { useHostStore } from '@/store/modules/host'

  interface Props {
    hostGroups: HostGroup[]
  }

  const props = defineProps<Props>()

  const emit = defineEmits<{
    success: []
  }>()

  const hostStore = useHostStore()

  const visible = ref(false)
  const loading = ref(false)
  const isEdit = ref(false)
  const parentGroup = ref<HostGroup | null>(null)
  const formRef = ref<FormInstance>()

  const formData = reactive({
    id: undefined as number | undefined,
    name: '',
    parent_id: undefined as number | undefined,
    description: '',
    sort: 0
  })

  const rules = {
    name: [
      { required: true, message: '请输入主机组名称', trigger: 'blur' },
      { min: 1, max: 50, message: '名称长度应在1-50个字符之间', trigger: 'blur' }
    ]
  }

  const hostGroupOptions = computed(() => {
    const buildTreeData = (
      groups: HostGroup[],
      parentId: number | null = null,
      excludeId?: number
    ): any[] => {
      return groups
        .filter((group) => group.parent_id === parentId && group.id !== excludeId)
        .map((group) => ({
          label: group.name,
          value: group.id,
          key: group.id,
          disabled: isEdit.value && group.id === formData.id, // 禁止选择自己作为父级
          children: buildTreeData(groups, group.id, excludeId)
        }))
    }
    return buildTreeData(props.hostGroups, null, formData.id)
  })

  const resetForm = () => {
    Object.assign(formData, {
      id: undefined,
      name: '',
      parent_id: undefined,
      description: '',
      sort: 0
    })
    parentGroup.value = null
  }

  const handleSubmit = async () => {
    try {
      await formRef.value?.validate()
      loading.value = true

      const submitData = {
        ...formData,
        parent_id: parentGroup.value?.id || formData.parent_id
      }

      if (isEdit.value) {
        await hostStore.updateHostGroup(formData.id!, submitData)
        ElMessage.success('更新成功')
      } else {
        await hostStore.addHostGroup(submitData)
        ElMessage.success('创建成功')
      }

      visible.value = false
      emit('success')
    } catch (error) {
      console.error('Submit error:', error)
      ElMessage.error('操作失败')
    } finally {
      loading.value = false
    }
  }

  const handleCancel = () => {
    visible.value = false
    resetForm()
  }

  const open = (record?: HostGroup, parent?: HostGroup) => {
    visible.value = true
    isEdit.value = !!record
    parentGroup.value = parent || null

    if (record) {
      Object.assign(formData, {
        id: record.id,
        name: record.name,
        parent_id: record.parent_id,
        description: record.description,
        sort: record.sort
      })
    } else {
      resetForm()
      if (parent) {
        formData.parent_id = parent.id
      }
    }
  }

  defineExpose({
    open
  })
</script>

<style scoped>
  .ant-form-item {
    margin-bottom: 16px;
  }
</style>
