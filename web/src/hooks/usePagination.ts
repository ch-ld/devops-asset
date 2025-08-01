import { reactive } from 'vue'
import type { TableProps } from 'ant-design-vue'

/**
 * 分页Hook
 */
export function usePagination(initialPageSize = 10) {
  const pagination = reactive({
    current: 1,
    pageSize: initialPageSize,
    total: 0,
    showSizeChanger: true,
    showQuickJumper: true,
    showTotal: (total: number, range: [number, number]) => 
      `共 ${total} 条记录，显示第 ${range[0]}-${range[1]} 条`,
    pageSizeOptions: ['10', '20', '50', '100']
  })

  const handleTableChange: TableProps['onChange'] = (pag, filters, sorter) => {
    if (pag) {
      pagination.current = pag.current || 1
      pagination.pageSize = pag.pageSize || initialPageSize
    }
  }

  const resetPagination = () => {
    pagination.current = 1
    pagination.total = 0
  }

  return {
    pagination,
    handleTableChange,
    resetPagination
  }
}
