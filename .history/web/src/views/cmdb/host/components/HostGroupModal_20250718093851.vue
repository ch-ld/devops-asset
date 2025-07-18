<template>
  <a-modal
    :title="isEdit ? '编辑主机组' : '添加主机组'"
    :open="visible"
    :width="600"
    :confirm-loading="loading"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <a-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      :label-col="{ span: 6 }"
      :wrapper-col="{ span: 18 }"
    >
      <a-form-item label="组名称" name="name">
        <a-input v-model:value="formData.name" placeholder="请输入主机组名称" />
      </a-form-item>

      <a-form-item label="父级分组" name="parent_id">
        <a-tree-select
          v-model:value="formData.parent_id"
          :tree-data="hostGroupOptions"
          placeholder="选择父级分组（可选）"
          allow-clear
          tree-default-expand-all
        />
      </a-form-item>

      <a-form-item label="描述" name="description">
        <a-textarea v-model:value="formData.description" placeholder="请输入描述信息" :rows="3" />
      </a-form-item>

      <a-form-item label="排序" name="sort">
        <a-input-number
          v-model:value="formData.sort"
          placeholder="排序值"
          :min="0"
          style="width: 100%"
        />
      </a-form-item>

      <a-alert
        v-if="parentGroup"
        :message="`将在 ${parentGroup.name} 下创建子分组`"
        type="info"
        show-icon
        style="margin-bottom: 16px"
      />
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
  import { ref, reactive, computed } from 'vue'
  import { message } from 'ant-design-vue'
  import type { FormInstance } from 'ant-design-vue'
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
          title: group.name,
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
        message.success('更新成功')
      } else {
        await hostStore.addHostGroup(submitData)
        message.success('创建成功')
      }

      visible.value = false
      emit('success')
    } catch (error) {
      console.error('Submit error:', error)
      message.error('操作失败')
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
