import { ref, computed } from 'vue'
import type { TableRowSelection } from 'ant-design-vue/es/table/interface'

/**
 * 表格选择Hook
 */
export function useSelection<T = any>() {
  const selectedRowKeys = ref<(string | number)[]>([])
  const selectedRows = ref<T[]>([])

  const hasSelected = computed(() => selectedRowKeys.value.length > 0)

  const rowSelection: TableRowSelection = {
    selectedRowKeys: selectedRowKeys.value,
    onChange: (keys: (string | number)[], rows: T[]) => {
      selectedRowKeys.value = keys
      selectedRows.value = rows
    },
    onSelectAll: (selected: boolean, selectedRows: T[], changeRows: T[]) => {
      // 可以在这里添加全选逻辑
    }
  }

  const clearSelection = () => {
    selectedRowKeys.value = []
    selectedRows.value = []
  }

  return {
    selectedRowKeys,
    selectedRows,
    hasSelected,
    rowSelection,
    clearSelection
  }
}
