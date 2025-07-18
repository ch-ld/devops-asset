<template>
  <el-dialog v-model="visible" title="批量移动主机" width="500px" :close-on-click-modal="false">
    <div class="move-container">
      <el-alert
        :title="`已选择 ${selectedHostIds.length} 台主机`"
        type="info"
        show-icon
        class="mb-2"
      />

      <el-form ref="formRef" :model="formData" label-width="100px">
        <el-form-item label="目标主机组">
          <el-tree-select
            v-model="formData.group_id"
            :data="hostGroupOptions"
            placeholder="选择目标主机组"
            clearable
            node-key="value"
            style="width: 100%"
          />
        </el-form-item>

        <el-alert
          title="移动说明"
          type="warning"
          show-icon
          description="选择目标主机组后，所选主机将被移动到该组下。如果不选择主机组，则将主机移出所有分组。"
        />
      </el-form>
    </div>

    <template #footer>
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
  import { ref, reactive, computed, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import type { HostGroup } from '@/types/api/host'
  import * as hostApi from '@/api/system/host'

  interface Props {
    visible?: boolean
    hostIds?: number[]
  }

  const props = withDefaults(defineProps<Props>(), {
    visible: false,
    hostIds: () => []
  })

  const emit = defineEmits(['update:visible', 'success'])

  const visible = ref(props.visible)
  const loading = ref(false)
  const selectedHostIds = ref<number[]>([])
  const formRef = ref<any>()
  const hostGroups = ref<HostGroup[]>([])

  const formData = reactive({
    group_id: undefined as number | undefined
  })

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
        .filter((g) => g.parent_id === parentId)
        .map((g) => ({
          label: g.name,
          value: g.id,
          children: buildTreeData(groups, g.id)
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
      // 无验证规则，直接通过
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
    .el-form-item {
      margin-bottom: 16px;
    }
  }
</style>
