<template>
  <div class="column-settings">
    <div class="settings-header">
      <h4>自定义表格列</h4>
      <p>拖拽调整列的顺序，勾选控制列的显示</p>
    </div>

    <div class="column-list">
      <div v-for="(element, index) in localColumns" :key="element.prop" class="column-item">
        <div class="column-info">
          <i class="el-icon-rank drag-handle"></i>
          <el-checkbox v-model="element.visible">
            {{ element.label }}
          </el-checkbox>
        </div>
        <div class="column-actions">
          <el-input-number
            v-model="element.width"
            :min="80"
            :max="500"
            size="small"
            controls-position="right"
            style="width: 100px"
          />
          <span class="width-label">px</span>
        </div>
      </div>
    </div>

    <div class="settings-actions">
      <el-button @click="handleReset">重置默认</el-button>
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" @click="handleSave">保存设置</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
// import { VueDraggable as draggable } from 'vue-draggable-plus'

interface Column {
  prop: string
  label: string
  visible: boolean
  width?: number
  minWidth?: number
  fixed?: string
  showTooltip?: boolean
  render?: string
}

interface Props {
  columns: Column[]
}

const props = defineProps<Props>()
const emit = defineEmits<{
  save: [columns: Column[]]
  cancel: []
}>()

const localColumns = ref<Column[]>([])

// 默认列配置
const defaultColumns: Column[] = [
  { prop: 'name', label: '主机名称', visible: true, width: 150, showTooltip: true },
  { prop: 'public_ip', label: '公网IP', visible: true, width: 140, render: 'IPCell' },
  { prop: 'private_ip', label: '私网IP', visible: true, width: 140, render: 'IPCell' },
  { prop: 'status', label: '状态', visible: true, width: 100, render: 'StatusCell' },
  { prop: 'configuration.cpu_cores', label: 'CPU', visible: true, width: 80, render: 'CPUCell' },
  { prop: 'configuration.memory_size', label: '内存', visible: true, width: 80, render: 'MemoryCell' },
  { prop: 'provider.name', label: '云厂商', visible: true, width: 120 },
  { prop: 'region', label: '地域', visible: true, width: 120 },
  { prop: 'os', label: '操作系统', visible: true, width: 120, showTooltip: true },
  { prop: 'created_at', label: '创建时间', visible: true, width: 160, render: 'DateCell' }
]

// 初始化
const initColumns = () => {
  localColumns.value = JSON.parse(JSON.stringify(props.columns))
}

// 拖拽结束 - 暂时禁用
// const handleDragEnd = () => {
//   // 拖拽结束后的处理逻辑
// }

// 重置默认
const handleReset = () => {
  localColumns.value = JSON.parse(JSON.stringify(defaultColumns))
  ElMessage.success('已重置为默认设置')
}

// 保存设置
const handleSave = () => {
  // 验证至少有一列可见
  const visibleColumns = localColumns.value.filter(col => col.visible)
  if (visibleColumns.length === 0) {
    ElMessage.warning('至少需要显示一列')
    return
  }

  emit('save', localColumns.value)
  ElMessage.success('列设置已保存')
}

// 取消
const handleCancel = () => {
  emit('cancel')
}

onMounted(() => {
  initColumns()
})
</script>

<style scoped>
.column-settings {
  padding: 20px 0;
}

.settings-header {
  margin-bottom: 20px;
}

.settings-header h4 {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.settings-header p {
  margin: 0;
  font-size: 14px;
  color: #6b7280;
}

.column-list {
  max-height: 400px;
  overflow-y: auto;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  padding: 8px;
}

.column-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  border-radius: 6px;
  margin-bottom: 8px;
  background: #f9fafb;
  transition: all 0.2s;
}

.column-item:hover {
  background: #f3f4f6;
}

.column-item:last-child {
  margin-bottom: 0;
}

.column-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.drag-handle {
  cursor: move;
  color: #9ca3af;
  font-size: 16px;
}

.drag-handle:hover {
  color: #6b7280;
}

.column-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.width-label {
  font-size: 12px;
  color: #6b7280;
}

.settings-actions {
  margin-top: 24px;
  text-align: center;
  border-top: 1px solid #e5e7eb;
  padding-top: 20px;
}

.settings-actions .el-button {
  margin: 0 8px;
}
</style>
