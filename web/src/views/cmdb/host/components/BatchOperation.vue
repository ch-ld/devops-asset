<template>
  <div class="batch-operation">
    <!-- 启动/停止/重启操作 -->
    <div v-if="['start', 'stop', 'reboot'].includes(operation)">
      <el-alert
        :title="getOperationTitle()"
        :description="getOperationDescription()"
        type="warning"
        show-icon
        :closable="false"
      />
      
      <div class="host-list">
        <h4>将要操作的主机 ({{ hosts.length }}台):</h4>
        <el-table :data="hosts" size="small" max-height="300">
          <el-table-column prop="name" label="主机名称" />
          <el-table-column prop="public_ip" label="公网IP" width="140">
            <template #default="{ row }">
              {{ Array.isArray(row.public_ip) ? row.public_ip[0] : row.public_ip }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="当前状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <!-- 移动分组操作 -->
    <div v-else-if="operation === 'move'">
      <el-alert
        title="批量移动主机分组"
        description="选择目标分组，将选中的主机移动到指定分组"
        type="info"
        show-icon
        :closable="false"
      />
      
      <el-form :model="moveForm" label-width="100px" class="operation-form">
        <el-form-item label="目标分组" required>
          <el-tree-select
            v-model="moveForm.group_id"
            :data="groupTree"
            :props="{ label: 'name', value: 'id' }"
            placeholder="请选择目标分组"
            clearable
            check-strictly
          />
        </el-form-item>
      </el-form>

      <div class="host-list">
        <h4>将要移动的主机 ({{ hosts.length }}台):</h4>
        <el-table :data="hosts" size="small" max-height="200">
          <el-table-column prop="name" label="主机名称" />
          <el-table-column prop="group.name" label="当前分组" />
        </el-table>
      </div>
    </div>

    <!-- 设置标签操作 -->
    <div v-else-if="operation === 'tags'">
      <el-alert
        title="批量设置标签"
        description="为选中的主机批量添加、移除或替换标签"
        type="info"
        show-icon
        :closable="false"
      />
      
      <el-form :model="tagsForm" label-width="100px" class="operation-form">
        <el-form-item label="操作类型" required>
          <el-radio-group v-model="tagsForm.action">
            <el-radio label="add">添加标签</el-radio>
            <el-radio label="remove">移除标签</el-radio>
            <el-radio label="replace">替换标签</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="标签" required>
          <div class="tags-input">
            <el-tag
              v-for="tag in tagsForm.tags"
              :key="tag"
              closable
              @close="removeTag(tag)"
              class="tag-item"
            >
              {{ tag }}
            </el-tag>
            <el-input
              v-if="inputVisible"
              ref="inputRef"
              v-model="inputValue"
              size="small"
              @keyup.enter="handleInputConfirm"
              @blur="handleInputConfirm"
              class="tag-input"
            />
            <el-button v-else size="small" @click="showInput">+ 添加标签</el-button>
          </div>
        </el-form-item>
      </el-form>

      <div class="host-list">
        <h4>将要操作的主机 ({{ hosts.length }}台):</h4>
        <el-table :data="hosts" size="small" max-height="200">
          <el-table-column prop="name" label="主机名称" />
          <el-table-column label="当前标签">
            <template #default="{ row }">
              <el-tag v-for="tag in row.tags" :key="tag" size="small" class="current-tag">
                {{ tag }}
              </el-tag>
              <span v-if="!row.tags || row.tags.length === 0">-</span>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <!-- 删除操作 -->
    <div v-else-if="operation === 'delete'">
      <el-alert
        title="批量删除主机"
        description="此操作将永久删除选中的主机，请谨慎操作！"
        type="error"
        show-icon
        :closable="false"
      />
      
      <div class="host-list">
        <h4>将要删除的主机 ({{ hosts.length }}台):</h4>
        <el-table :data="hosts" size="small" max-height="300">
          <el-table-column prop="name" label="主机名称" />
          <el-table-column prop="public_ip" label="公网IP" width="140">
            <template #default="{ row }">
              {{ Array.isArray(row.public_ip) ? row.public_ip[0] : row.public_ip }}
            </template>
          </el-table-column>
          <el-table-column prop="provider.name" label="云厂商" width="120" />
        </el-table>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="operation-actions">
      <el-button @click="handleCancel">取消</el-button>
      <el-button 
        type="primary" 
        @click="handleConfirm" 
        :loading="operating"
        :disabled="!canConfirm"
      >
        确认{{ getOperationText() }}
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  batchDeleteHosts,
  batchMoveHosts,
  batchUpdateTags,
  batchStart,
  batchStop,
  batchReboot,
  getHostGroupTree
} from '@/api/system/host'
import type { Host, HostGroup } from '@/types/api/host'

interface Props {
  operation: string
  hosts: Host[]
}

const props = defineProps<Props>()
const emit = defineEmits<{
  success: []
  cancel: []
}>()

const operating = ref(false)
const groupTree = ref<HostGroup[]>([])

// 移动分组表单
const moveForm = reactive({
  group_id: undefined as number | undefined
})

// 标签表单
const tagsForm = reactive({
  action: 'add' as 'add' | 'remove' | 'replace',
  tags: [] as string[]
})

// 标签输入相关
const inputVisible = ref(false)
const inputValue = ref('')
const inputRef = ref()

// 计算属性
const canConfirm = computed(() => {
  switch (props.operation) {
    case 'move':
      return moveForm.group_id !== undefined
    case 'tags':
      return tagsForm.tags.length > 0
    default:
      return true
  }
})

// 获取操作标题
const getOperationTitle = () => {
  const titles: Record<string, string> = {
    start: '批量启动主机',
    stop: '批量停止主机',
    reboot: '批量重启主机'
  }
  return titles[props.operation] || ''
}

// 获取操作描述
const getOperationDescription = () => {
  const descriptions: Record<string, string> = {
    start: '此操作将启动选中的所有主机，请确认操作',
    stop: '此操作将停止选中的所有主机，请确认操作',
    reboot: '此操作将重启选中的所有主机，请确认操作'
  }
  return descriptions[props.operation] || ''
}

// 获取操作文本
const getOperationText = () => {
  const texts: Record<string, string> = {
    start: '启动',
    stop: '停止',
    reboot: '重启',
    move: '移动',
    tags: '设置标签',
    delete: '删除'
  }
  return texts[props.operation] || '操作'
}

// 状态相关方法
const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    running: 'success',
    stopped: 'warning',
    error: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    running: '运行中',
    stopped: '已停止',
    error: '异常'
  }
  return statusMap[status] || status
}

// 标签操作
const removeTag = (tag: string) => {
  tagsForm.tags.splice(tagsForm.tags.indexOf(tag), 1)
}

const showInput = () => {
  inputVisible.value = true
  nextTick(() => {
    inputRef.value?.focus()
  })
}

const handleInputConfirm = () => {
  if (inputValue.value && !tagsForm.tags.includes(inputValue.value)) {
    tagsForm.tags.push(inputValue.value)
  }
  inputVisible.value = false
  inputValue.value = ''
}

// 获取主机分组树
const fetchGroupTree = async () => {
  try {
    const response = await getHostGroupTree()
    groupTree.value = response.data || []
  } catch (error) {
    console.error('获取主机分组失败:', error)
  }
}

// 确认操作
const handleConfirm = async () => {
  try {
    operating.value = true
    const hostIds = props.hosts.map(host => host.id)

    switch (props.operation) {
      case 'start':
        await batchStart({ ids: hostIds })
        ElMessage.success('批量启动操作已提交')
        break
      case 'stop':
        await batchStop({ ids: hostIds })
        ElMessage.success('批量停止操作已提交')
        break
      case 'reboot':
        await batchReboot({ ids: hostIds })
        ElMessage.success('批量重启操作已提交')
        break
      case 'move':
        await batchMoveHosts({ ids: hostIds, group_id: moveForm.group_id })
        ElMessage.success('批量移动成功')
        break
      case 'tags':
        await batchUpdateTags({ 
          ids: hostIds, 
          tags: tagsForm.tags, 
          action: tagsForm.action 
        })
        ElMessage.success('批量设置标签成功')
        break
      case 'delete':
        await ElMessageBox.confirm(
          `确定要删除这 ${props.hosts.length} 台主机吗？此操作不可恢复！`,
          '确认删除',
          { type: 'error' }
        )
        await batchDeleteHosts(hostIds)
        ElMessage.success('批量删除成功')
        break
    }

    emit('success')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`${getOperationText()}失败`)
    }
  } finally {
    operating.value = false
  }
}

// 取消操作
const handleCancel = () => {
  emit('cancel')
}

// 初始化
onMounted(() => {
  if (props.operation === 'move') {
    fetchGroupTree()
  }
})
</script>

<style scoped>
.batch-operation {
  padding: 20px 0;
}

.operation-form {
  margin: 20px 0;
}

.host-list {
  margin: 20px 0;
}

.host-list h4 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.operation-actions {
  margin-top: 24px;
  text-align: center;
}

.operation-actions .el-button {
  margin: 0 8px;
}

.tags-input {
  min-height: 32px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 4px 8px;
}

.tag-item,
.current-tag {
  margin-right: 8px;
  margin-bottom: 4px;
}

.tag-input {
  width: 90px;
  margin-left: 8px;
  vertical-align: bottom;
}
</style>
