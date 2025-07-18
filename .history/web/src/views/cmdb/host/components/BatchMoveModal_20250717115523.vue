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
  import { ref, reactive, computed, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import type { HostGroup } from '@/types/api/host'
  import * as hostApi from '@/api/system/host'
  import { useHostStore } from '@/store/modules/host'

  interface Props {
    visible?: boolean;
    hostIds?: number[];
  }

  const props = withDefaults(defineProps<Props>(), {
    visible: false,
    hostIds: () => []
  })

  const emit = defineEmits(['update:visible', 'success'])

  const hostStore = useHostStore()
  const visible = ref(props.visible)
  const loading = ref(false)
  const selectedHostIds = ref<number[]>([])
  const formRef = ref<any>()
  const hostGroups = ref<HostGroup[]>([])

  const formData = reactive({
    group_id: undefined as number | undefined
  })

  const rules = {
    // 不设置必填规则，允许移出分组
  }

  const fetchHostGroups = async () => {
    try {
      loading.value = true
      const response = await hostApi.getHostGroupTree()
      hostGroups.value = response as HostGroup[]
    } catch (error) {
      console.error('Failed to fetch host groups:', error)
      ElMessage.error('获取主机组失败，请刷新重试')
      hostGroups.value = []
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    fetchHostGroups()
  })

  const hostGroupOptions = computed(() => {
    if (!hostGroups.value || !Array.isArray(hostGroups.value)) {
      return []
    }
    
    const buildTreeData = (groups: HostGroup[], parentId: number | null = null): any[] => {
      return groups
        .filter((group) => group.parent_id === parentId)
        .map((group) => ({
          title: group.name,
          value: group.id,
          key: group.id,
          children: buildTreeData(groups, group.id)
        }))
    }
    return buildTreeData(hostGroups.value)
  })

  const resetForm = () => {
    formData.group_id = undefined
  }

  const handleSubmit = async () => {
    if (!selectedHostIds.value.length) {
      ElMessage.warning('请选择要移动的主机')
      return
    }
    
    try {
      await formRef.value?.validate()
      loading.value = true

      try {
        // 调用批量移动API
        // 直接调用API而不是通过store
        await hostApi.batchMoveHosts({
          ids: selectedHostIds.value,
          group_id: formData.group_id
        })
        
        ElMessage.success('移动成功')
        visible.value = false
        emit('success')
      } catch (error: any) {
        console.error('Move error:', error)
        ElMessage.error(`移动失败: ${error.message || '未知错误'}`)
      }
    } catch (error) {
      console.error('Validation error:', error)
      ElMessage.error('表单验证失败')
    } finally {
      loading.value = false
    }
  }

  const handleCancel = () => {
    visible.value = false
    resetForm()
    emit('update:visible', false)
  }

  const open = (hostIds: number[]) => {
    selectedHostIds.value = hostIds || []
    visible.value = true
    resetForm()
    
    // 确保有主机组数据
    if (!hostGroups.value.length) {
      fetchHostGroups()
    }
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
