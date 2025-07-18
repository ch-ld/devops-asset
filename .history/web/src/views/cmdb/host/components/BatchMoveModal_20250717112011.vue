<template>
  <a-modal
    title="批量移动主机"
    :open="visible"
    :width="500"
    :confirm-loading="loading"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <div class="move-container">
      <a-alert
        :message="`已选择 ${selectedHostIds.length} 台主机`"
        type="info"
        show-icon
        style="margin-bottom: 16px"
      />

      <a-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
      >
        <a-form-item label="目标主机组" name="group_id">
          <a-tree-select
            v-model:value="formData.group_id"
            :tree-data="hostGroupOptions"
            placeholder="选择目标主机组"
            allow-clear
            tree-default-expand-all
          />
        </a-form-item>

        <a-form-item label="操作说明">
          <a-alert
            message="移动说明"
            description="选择目标主机组后，所选主机将被移动到该组下。如果不选择主机组，则将主机移出所有分组。"
            type="warning"
            show-icon
          />
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
  import { ref, reactive, computed } from 'vue'
  import { ElMessage } from 'element-plus'
  // Remove FormInstance import and just use any type
  import type { HostGroup } from '@/types/api/host'

  interface Props {
    hostGroups: HostGroup[]
  }

  const props = defineProps<Props>()

  const emit = defineEmits<{
    success: []
  }>()

  const visible = ref(false)
  const loading = ref(false)
  const selectedHostIds = ref<number[]>([])
  const formRef = ref<any>() // Use any instead of FormInstance

  const formData = reactive({
    group_id: undefined as number | undefined
  })

  const rules = {
    // 不设置必填规则，允许移出分组
  }

  const hostGroupOptions = computed(() => {
    const buildTreeData = (groups: HostGroup[], parentId: number | null = null): any[] => {
      if (!groups || !Array.isArray(groups)) {
        return []
      }
      return groups
        .filter((group) => group.parent_id === parentId)
        .map((group) => ({
          title: group.name,
          value: group.id,
          key: group.id,
          children: buildTreeData(groups, group.id)
        }))
    }
    return buildTreeData(props.hostGroups || [])
  })

  const resetForm = () => {
    formData.group_id = undefined
  }

  const handleSubmit = async () => {
    try {
      await formRef.value?.validate()
      loading.value = true

      // 这里应该调用批量移动API
      console.log('Move hosts:', {
        hostIds: selectedHostIds.value,
        groupId: formData.group_id
      })

      // 模拟API调用
      await new Promise((resolve) => setTimeout(resolve, 1000))

      ElMessage.success('移动成功')
      visible.value = false
      emit('success')
    } catch (error) {
      console.error('Move error:', error)
      ElMessage.error('移动失败')
    } finally {
      loading.value = false
    }
  }

  const handleCancel = () => {
    visible.value = false
    resetForm()
  }

  const open = (hostIds: number[]) => {
    selectedHostIds.value = hostIds
    visible.value = true
    resetForm()
  }

  defineExpose({
    open
  })
</script>

<style scoped>
  .move-container {
    .ant-form-item {
      margin-bottom: 16px;
    }
  }
</style>
