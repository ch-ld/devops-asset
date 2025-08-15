<template>
  <el-card shadow="never" class="table-card">
    <!-- 表格头部 -->
    <template #header>
      <div class="table-header">
        <div class="table-title">
          <span>{{ title }}</span>
          <el-tag v-if="total !== undefined" type="info" size="small">
            共 {{ total }} 条
          </el-tag>
        </div>
        <div class="table-actions">
          <slot name="header-actions">
            <el-button @click="handleRefresh" :loading="loading">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button 
              v-if="showBatchDelete"
              @click="handleBatchDelete" 
              :disabled="!hasSelected"
              type="danger"
            >
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
          </slot>
        </div>
      </div>
    </template>

    <!-- 表格内容 -->
    <el-table
      ref="tableRef"
      :data="data"
      :loading="loading"
      @selection-change="handleSelectionChange"
      @sort-change="handleSortChange"
      row-key="id"
      v-bind="tableProps"
    >
      <!-- 选择列 -->
      <el-table-column 
        v-if="showSelection"
        type="selection" 
        width="55" 
        :selectable="selectable"
      />
      
      <!-- 序号列 -->
      <el-table-column 
        v-if="showIndex"
        type="index" 
        label="序号" 
        width="60"
        :index="getIndex"
      />
      
      <!-- 动态列 -->
      <slot />
      
      <!-- 操作列 -->
      <el-table-column 
        v-if="showActions"
        label="操作" 
        :width="actionWidth"
        :fixed="actionFixed"
      >
        <template #default="scope">
          <slot name="actions" :row="scope.row" :index="scope.$index">
            <el-button type="primary" size="small" text @click="handleEdit(scope.row)">
              编辑
            </el-button>
            <el-button type="danger" size="small" text @click="handleDelete(scope.row)">
              删除
            </el-button>
          </slot>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div v-if="showPagination" class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="pageSizes"
        :layout="paginationLayout"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Refresh, Delete } from '@element-plus/icons-vue'
import type { TableInstance } from 'element-plus'

export interface Props {
  /** 表格标题 */
  title?: string
  /** 表格数据 */
  data: any[]
  /** 加载状态 */
  loading?: boolean
  /** 数据总数 */
  total?: number
  /** 当前页码 */
  currentPage?: number
  /** 每页条数 */
  pageSize?: number
  /** 每页条数选项 */
  pageSizes?: number[]
  /** 分页布局 */
  paginationLayout?: string
  /** 是否显示选择框 */
  showSelection?: boolean
  /** 是否显示序号 */
  showIndex?: boolean
  /** 是否显示操作列 */
  showActions?: boolean
  /** 操作列宽度 */
  actionWidth?: string | number
  /** 操作列是否固定 */
  actionFixed?: string | boolean
  /** 是否显示分页 */
  showPagination?: boolean
  /** 是否显示批量删除 */
  showBatchDelete?: boolean
  /** 行选择函数 */
  selectable?: (row: any, index: number) => boolean
  /** 表格额外属性 */
  tableProps?: Record<string, any>
}

const props = withDefaults(defineProps<Props>(), {
  title: '数据列表',
  loading: false,
  currentPage: 1,
  pageSize: 20,
  pageSizes: () => [10, 20, 50, 100],
  paginationLayout: 'total, sizes, prev, pager, next, jumper',
  showSelection: true,
  showIndex: false,
  showActions: true,
  actionWidth: 200,
  actionFixed: 'right',
  showPagination: true,
  showBatchDelete: true,
  tableProps: () => ({})
})

const emit = defineEmits<{
  'refresh': []
  'selection-change': [selection: any[]]
  'sort-change': [sortInfo: { column: any; prop: string; order: string }]
  'size-change': [size: number]
  'current-change': [current: number]
  'edit': [row: any]
  'delete': [row: any]
  'batch-delete': [selection: any[]]
}>()

// 表格引用
const tableRef = ref<TableInstance>()

// 选中的行
const selectedRows = ref<any[]>([])

// 计算属性
const hasSelected = computed(() => selectedRows.value.length > 0)

// 获取序号
const getIndex = (index: number) => {
  return (props.currentPage - 1) * props.pageSize + index + 1
}

// 事件处理
const handleRefresh = () => {
  emit('refresh')
}

const handleSelectionChange = (selection: any[]) => {
  selectedRows.value = selection
  emit('selection-change', selection)
}

const handleSortChange = (sortInfo: any) => {
  emit('sort-change', sortInfo)
}

const handleSizeChange = (size: number) => {
  emit('size-change', size)
}

const handleCurrentChange = (current: number) => {
  emit('current-change', current)
}

const handleEdit = (row: any) => {
  emit('edit', row)
}

const handleDelete = (row: any) => {
  emit('delete', row)
}

const handleBatchDelete = () => {
  emit('batch-delete', selectedRows.value)
}

// 暴露方法
defineExpose({
  tableRef,
  selectedRows,
  clearSelection: () => tableRef.value?.clearSelection(),
  toggleRowSelection: (row: any, selected?: boolean) => {
    tableRef.value?.toggleRowSelection(row, selected)
  }
})
</script>

<style scoped>
.table-card {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .table-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 500;
  }

  .table-actions {
    display: flex;
    gap: 8px;
  }

  .pagination-container {
    display: flex;
    justify-content: flex-end;
    margin-top: 16px;
  }
}
</style>
