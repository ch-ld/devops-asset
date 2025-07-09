<template>
  <ArtTableFullScreen>
    <div class="department-page" id="table-full-screen">
      <!-- 搜索栏 -->
      <ArtSearchBar
        v-model:filter="searchForm"
        :items="searchItems"
        @reset="resetSearch"
        @search="search"
      ></ArtSearchBar>

      <ElCard shadow="never" class="art-table-card">
        <!-- 表格头部 -->
        <ArtTableHeader
          :columnList="columnOptions"
          v-model:columns="columnChecks"
          @refresh="handleRefresh"
        >
          <template #left>
            <ElButton @click="showDialog('add')" v-ripple>添加环境</ElButton>
          </template>
        </ArtTableHeader>

        <!-- 表格 -->
        <ArtTable
          :data="safeTableData"
          :currentPage="pagination.currentPage"
          :pageSize="pagination.pageSize"
          :total="pagination.total"
          :loading="loading"
          :hideOnSinglePage="false"
          :marginTop="10"
          height="100%"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        >
          <template #default>
            <ElTableColumn
              v-for="col in filteredColumns"
              :key="col.prop || col.type"
              v-bind="col"
            >
              <!-- 自定义人数列的渲染 -->
              <template #default="scope" v-if="col.prop === 'users'">
                {{ Array.isArray(scope.row.users) ? scope.row.users.length : 0 }}
              </template>

              <!-- 自定义状态列的渲染 -->
              <template #default="scope" v-else-if="col.prop === 'status'">
                <ElTag :type="scope.row.status === 1 ? 'primary' : 'warning'">
                  {{ scope.row.status === 1 ? '启用' : '禁用' }}
                </ElTag>
              </template>

              <!-- 自定义操作列的渲染 -->
              <template #default="scope" v-else-if="col.prop === 'operation'">
                <div class="operation-column-container">
                  <ArtButtonTable type="edit" @click="showDialog('edit', scope.row)" />
                  <ArtButtonTable type="delete" @click="deleteDepartment(scope.row.id)" />
                </div>
              </template>
            </ElTableColumn>
          </template>
        </ArtTable>
      </ElCard>
    </div>

    <ElDialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '添加环境' : '编辑环境'"
      width="600px"
      align-center
      :close-on-click-modal="false"
      @closed="resetForm"
    >
      <ElForm ref="formRef" :model="formData" :rules="rules" label-width="85px">
        <ElRow :gutter="20">
          <ElCol :span="12">
            <ElFormItem label="环境名称" prop="name">
              <ElInput v-model="formData.name" placeholder="请输入环境名称" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="排序" prop="sort">
              <ElInputNumber
                v-model="formData.sort"
                style="width: 100%"
                :min="1"
                controls-position="right"
                placeholder="请输入排序号"
              />
            </ElFormItem>
          </ElCol>
        </ElRow>
        <ElRow :gutter="20">
          <ElCol :span="12">
            <ElFormItem label="启用" prop="status">
              <ElSwitch v-model="formData.status" />
            </ElFormItem>
          </ElCol>
        </ElRow>
      </ElForm>
      <template #footer>
        <div class="dialog-footer">
          <ElButton @click="dialogVisible = false">取消</ElButton>
          <ElButton type="primary" @click="submitForm" :loading="submitLoading">提交</ElButton>
        </div>
      </template>
    </ElDialog>
  </ArtTableFullScreen>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted, computed } from 'vue'
  import { ElMessageBox, ElMessage } from 'element-plus'
  import type { FormInstance, FormRules } from 'element-plus'
  import {
    getDepartmentList,
    addDepartment,
    updateDepartment,
    deleteDepartment as apiDeleteDepartment
  } from '@/api/system/api'
  import { ApiStatus } from '@/utils/http/status'
  import { useCheckedColumns } from '@/composables/useCheckedColumns'
  import { SearchFormItem } from '@/types'

  // 状态变量
  const dialogType = ref('add')
  const dialogVisible = ref(false)
  const loading = ref(false)
  const submitLoading = ref(false)
  const currentId = ref<number | null>(null)
  const tableData = ref<any[]>([])

  // 表单数据
  const formData = reactive({
    name: '',
    sort: '1',
    status: true
  })

  // 分页配置
  const pagination = reactive({
    currentPage: 1,
    pageSize: 10,
    total: 0
  })

  // 搜索表单
  const searchForm = reactive({
    name: '',
    status: null as null | number
  })

  // 搜索表单配置项
  const searchItems: SearchFormItem[] = [
    {
      label: '环境名称',
      prop: 'name',
      type: 'input',
      elColSpan: 6, // 从8改为6，缩短显示宽度
      config: {
        clearable: true,
        placeholder: '请输入环境名称'
      }
    },
    {
      label: '状态',
      prop: 'status',
      type: 'select',
      elColSpan: 6, // 从8改为6，缩短显示宽度
      config: {
        clearable: true,
        placeholder: '请选择状态'
      },
      options: [
        { label: '启用', value: 1 },
        { label: '禁用', value: 2 }
      ]
    }
  ]

  // 列配置选项
  const columnOptions = [
    { label: '环境名称', prop: 'name' },
    { label: '排序', prop: 'sort' },
    { label: '应用数量', prop: 'users' },
    { label: '状态', prop: 'status' },
    { label: '操作', prop: 'operation' }
  ]

  // 动态列配置
  const { columnChecks, columns } = useCheckedColumns(() => [
    {
      prop: 'name',
      label: '环境名称',
      align: 'center'
    },
    {
      prop: 'sort',
      label: '排序',
      sortable: true,
      align: 'center'
    },
    {
      prop: 'users',
      label: '应用数量',
      align: 'center'
    },
    {
      prop: 'status',
      label: '状态',
      align: 'center'
    },
    {
      prop: 'operation',
      label: '操作',
      align: 'center'
    }
  ])

  // 根据列选中状态筛选得到最终显示的列
  const filteredColumns = computed(() => {
    return columns.value
  })

  // 表单验证规则
  const rules = reactive<FormRules>({
    name: [
      { required: true, message: '请输入部门名称', trigger: 'blur' },
      { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
    ],
    sort: [
      { required: true, message: '请输入排序号', trigger: 'blur' },
      { pattern: /^[0-9]*$/, message: '请输入数字', trigger: 'blur' }
    ]
  })

  // 表单实例引用
  const formRef = ref<FormInstance>()
  const searchFormRef = ref<FormInstance>()

  // 在组件挂载时获取部门列表
  onMounted(async () => {
    await refreshDepartmentList()
  })

  // 刷新表格数据
  const handleRefresh = () => {
    tableData.value = []
    loading.value = true
    refreshDepartmentList()
  }

  // 刷新部门列表
  const refreshDepartmentList = async () => {
    loading.value = true
    try {
      const params = {
        page: pagination.currentPage,
        pageSize: pagination.pageSize,
        ...searchForm
      }
      const res = await getDepartmentList(params)
      if (res.code === ApiStatus.success) {
        tableData.value = res.data || []

        // 使用返回值中的count字段作为总数
        if (res.count !== undefined) {
          pagination.total = res.count
        } else if (res.meta && res.meta.count) {
          pagination.total = res.meta.count
        } else if (res.meta && res.meta.total) {
          pagination.total = res.meta.total
        } else {
          pagination.total = res.data?.length || 0
        }
      } else {
        ElMessage.error(`获取部门列表失败: ${res.message}`)
      }
    } catch (err) {
      console.error('获取部门列表出错:', err)
      ElMessage.error('获取部门列表失败，请检查网络连接')
    } finally {
      loading.value = false
    }
  }

  // 搜索部门
  const search = () => {
    pagination.currentPage = 1
    refreshDepartmentList() // 调用刷新方法获取数据
  }

  // 重置搜索条件并刷新列表
  const resetSearch = () => {
    searchForm.name = ''
    searchForm.status = null
    pagination.currentPage = 1
    refreshDepartmentList()
  }

  // 重置表单
  const resetForm = () => {
    formData.name = ''
    formData.sort = '1'
    formData.status = true
    currentId.value = null
  }

  // 显示对话框
  const showDialog = (type: string, row?: any) => {
    dialogType.value = type
    dialogVisible.value = true

    if (type === 'edit' && row) {
      currentId.value = row.id
      formData.name = row.name
      formData.sort = String(row.sort)
      formData.status = row.status === 1
    } else {
      resetForm()
    }
  }

  // 删除部门
  const deleteDepartment = (id: number) => {
    ElMessageBox.confirm('确定要删除该部门吗？', '删除部门', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
      .then(async () => {
        loading.value = true
        try {
          const res = await apiDeleteDepartment(id)
          if (res.code === ApiStatus.success) {
            ElMessage.success('删除部门成功')
            await refreshDepartmentList()
          } else {
            ElMessage.error(`删除部门失败: ${res.message}`)
          }
        } catch (err) {
          console.error('删除部门出错:', err)
          ElMessage.error('删除部门失败')
        } finally {
          loading.value = false
        }
      })
      .catch(() => {
        // 用户取消删除，不做任何操作
      })
  }

  // 提交表单
  const submitForm = async () => {
    if (!formRef.value) return

    await formRef.value.validate(async (valid) => {
      if (valid) {
        submitLoading.value = true
        try {
          const params = {
            name: formData.name,
            sort: parseInt(formData.sort),
            status: formData.status ? 1 : 0
          }

          if (dialogType.value === 'edit') {
            // 编辑部门
            if (!currentId.value) {
              ElMessage.error('部门ID无效')
              return
            }

            const res = await updateDepartment({
              ...params,
              id: currentId.value
            })

            if (res.code === ApiStatus.success) {
              ElMessage.success('修改部门成功')
              dialogVisible.value = false
              await refreshDepartmentList()
            } else {
              ElMessage.error(`修改部门失败: ${res.message}`)
            }
          } else {
            // 添加部门
            const res = await addDepartment(params)

            if (res.code === ApiStatus.success) {
              ElMessage.success('添加部门成功')
              dialogVisible.value = false
              await refreshDepartmentList()
            } else {
              ElMessage.error(`添加部门失败: ${res.message}`)
            }
          }
        } catch (err) {
          console.error('提交表单出错:', err)
          ElMessage.error(`${dialogType.value === 'add' ? '添加' : '修改'}部门失败`)
        } finally {
          submitLoading.value = false
        }
      }
    })
  }

  // 页码变化处理
  const handleCurrentChange = (page: number) => {
    pagination.currentPage = page
    refreshDepartmentList() // 调用刷新方法获取数据
  }

  // 每页条数变化处理
  const handleSizeChange = (size: number) => {
    pagination.pageSize = size
    pagination.currentPage = 1
    refreshDepartmentList() // 调用刷新方法获取数据
  }

  // 使用后端数据进行表格渲染，不再进行前端过滤和分页
  const safeTableData = Array.isArray(tableData) ? tableData : []
</script>

<style lang="scss" scoped>
  .department-page {
    .table-container {
      flex: 1;
      min-height: 0;
      padding: 16px;
    }

    .search-container {
      display: flex;
      justify-content: space-between;
      margin-bottom: 16px;

      .el-input {
        width: 240px;
        margin-right: 16px;
      }
    }

    .operation-column-container {
      display: flex;
      justify-content: center;
      align-items: center;
      gap: 8px;

      .art-button-table {
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

        &:hover {
          transform: scale(1.1);
        }
      }
    }

    .svg-icon {
      width: 1.8em;
      height: 1.8em;
      overflow: hidden;
      vertical-align: -8px;
      fill: currentcolor;
    }

    // 优化卡片样式
    :deep(.art-table-card) {
      border-radius: 12px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
      transition: all 0.3s ease;

      &:hover {
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
      }
    }

    // 优化表格样式
    :deep(.el-table) {
      border-radius: 8px;
      overflow: hidden;

      .el-table__header {
        background: var(--art-gray-50);
      }

      .el-table__row {
        transition: all 0.2s ease;

        &:hover {
          background: var(--art-gray-50);
          transform: translateY(-1px);
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
        }
      }
    }

    // 优化按钮样式
    :deep(.el-button) {
      border-radius: 8px;
      font-weight: 500;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

      &:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      }

      &.el-button--primary {
        background: var(--main-color);
        border-color: var(--main-color);

        &:hover {
          background: var(--main-color-hover);
          border-color: var(--main-color-hover);
        }
      }
    }

    // 优化标签样式
    :deep(.el-tag) {
      border-radius: 6px;
      font-weight: 500;
      transition: all 0.2s ease;

      &:hover {
        transform: scale(1.05);
      }
    }

    // 优化对话框样式
    :deep(.el-dialog) {
      border-radius: 12px;
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);

      .el-dialog__header {
        border-bottom: 1px solid var(--art-card-border);
        padding: 20px 24px;

        .el-dialog__title {
          font-weight: 600;
          color: var(--art-text-color-primary);
        }
      }

      .el-dialog__body {
        padding: 24px;
      }

      .el-dialog__footer {
        border-top: 1px solid var(--art-card-border);
        padding: 16px 24px;
      }
    }

    // 优化表单样式
    :deep(.el-form-item) {
      margin-bottom: 20px;

      .el-form-item__label {
        font-weight: 500;
        color: var(--art-text-color-primary);
      }

      .el-input,
      .el-input-number {
        .el-input__wrapper {
          border-radius: 8px;
          transition: all 0.3s ease;

          &:hover {
            box-shadow: 0 0 0 1px var(--main-color) inset;
          }

          &.is-focus {
            box-shadow: 0 0 0 1px var(--main-color) inset;
          }
        }
      }

      .el-switch {
        .el-switch__core {
          border-radius: 12px;
        }
      }
    }

    // 优化搜索栏样式
    :deep(.art-search-bar) {
      background: var(--art-bg-color);
      border-radius: 12px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
      margin-bottom: 16px;
      padding: 16px;
      transition: all 0.3s ease;

      &:hover {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      }
    }

    // 优化表格头部样式
    :deep(.art-table-header) {
      background: var(--art-bg-color);
      border-radius: 8px;
      margin-bottom: 16px;
      padding: 16px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    }

    // 添加页面加载动画
    animation: pageSlideIn 0.3s ease-out;
  }

  // 页面加载动画
  @keyframes pageSlideIn {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  // 响应式优化
  @media (max-width: 768px) {
    .department-page {
      .search-container {
        flex-direction: column;
        gap: 12px;

        .el-input {
          width: 100%;
          margin-right: 0;
        }
      }

      .operation-column-container {
        flex-direction: column;
        gap: 4px;
      }
    }
  }
</style>
