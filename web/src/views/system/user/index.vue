<template>
  <ArtTableFullScreen>
    <div class="user-page" id="table-full-screen">
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
            <ElButton @click="showDialog('add')" v-ripple>添加应用</ElButton>
          </template>
        </ArtTableHeader>

        <!-- 表格 -->
        <ArtTable
          :data="tableData"
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
              <!-- 自定义用户名和账号列的渲染 -->
              <template #default="scope" v-if="col.prop === 'User.name'">
                {{ scope.row.User?.name || '--' }}
              </template>

              <template #default="scope" v-else-if="col.prop === 'User.username'">
                {{ scope.row.User?.username || '--' }}
              </template>

              <!-- 自定义手机号列的渲染 -->
              <template #default="scope" v-else-if="col.prop === 'User.phone'">
                {{ scope.row.User?.phone || '--' }}
              </template>

              <!-- 自定义性别列的渲染 -->
              <template #default="scope" v-else-if="col.prop === 'User.gender'">
                <ElTag v-if="scope.row.User?.gender === 1" type="success" effect="light"
                  >男</ElTag
                >
                <ElTag v-else-if="scope.row.User?.gender === 2" type="danger" effect="light"
                  >女</ElTag
                >
                <span v-else>--</span>
              </template>

              <!-- 自定义状态列的渲染 -->
              <template #default="scope" v-else-if="col.prop === 'User.status'">
                <ElTag :type="getTagType(scope.row.User?.status)">
                  {{ buildTagText(scope.row.User?.status) }}
                </ElTag>
              </template>

              <!-- 自定义操作列的渲染 -->
              <template #default="scope" v-else-if="col.prop === 'operation'">
                <div class="operation-column-container">
                  <ArtButtonTable type="edit" @click="showDialog('edit', scope.row)" />
                  <ArtButtonTable type="delete" @click="handleDeleteUser(scope.row)" />
                </div>
              </template>
            </ElTableColumn>
          </template>
        </ArtTable>
      </ElCard>
    </div>

    <ElDialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '添加应用' : '编辑应用'"
      width="800px"
      align-center
      :close-on-click-modal="false"
      @closed="resetForm"
    >
      <ElForm ref="formRef" :model="formData" :rules="computedRules" label-width="85px">
        <ElRow :gutter="20">
          <ElCol :span="12">
            <ElFormItem label="应用名称" prop="name">
              <ElInput v-model="formData.name" placeholder="请输入应用名称" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="应用标识" prop="username">
              <ElInput v-model="formData.username" placeholder="请输入应用标识" />
            </ElFormItem>
          </ElCol>
        </ElRow>
        <ElRow :gutter="20">
          <ElCol :span="12">
            <ElFormItem label="环境" prop="department_id">
              <ElSelect
                v-model="formData.department_id"
                placeholder="请选择环境"
                style="width: 100%"
                clearable
              >
                <ElOption
                  v-for="item in safeDepartmentList"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </ElSelect>
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="部署类型" prop="role_id">
              <ElSelect
                v-model="formData.role_id"
                placeholder="请选择部署类型"
                style="width: 100%"
                clearable
              >
                <ElOption
                  v-for="item in safeRoleList"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </ElSelect>
            </ElFormItem>
          </ElCol>
        </ElRow>
        <ElRow :gutter="20">
          <ElCol :span="12">
            <ElFormItem label="Git 仓库" prop="phone">
              <ElInput v-model="formData.phone" placeholder="请输入 Git 仓库地址" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="分支" prop="gender">
              <ElRadioGroup v-model="formData.gender">
                <ElRadio :label="0">main</ElRadio>
                <ElRadio :label="1">develop</ElRadio>
                <ElRadio :label="2">feature</ElRadio>
              </ElRadioGroup>
            </ElFormItem>
          </ElCol>
        </ElRow>
        <ElRow :gutter="20">
          <ElCol :span="12">
            <ElFormItem label="状态" prop="status">
              <ElSelect v-model="formData.status" placeholder="请选择状态" style="width: 100%">
                <ElOption label="启用" :value="1" />
                <ElOption label="禁用" :value="2" />
              </ElSelect>
            </ElFormItem>
          </ElCol>
          <ElCol :span="12" v-if="dialogType === 'add'">
            <ElFormItem label="部署密钥" prop="password">
              <ElInput
                v-model="formData.password"
                type="password"
                placeholder="请输入部署密钥"
                show-password
              />
            </ElFormItem>
          </ElCol>
        </ElRow>
      </ElForm>

      <template #footer>
        <div class="dialog-footer">
          <ElButton @click="dialogVisible = false">取 消</ElButton>
          <ElButton type="primary" @click="handleSubmit">确 定</ElButton>
        </div>
      </template>
    </ElDialog>
  </ArtTableFullScreen>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted, nextTick, computed } from 'vue'
  import {
    getUserList,
    addUser,
    updateUser,
    deleteUser as apiDeleteUser,
    getDepartmentList,
    getRoleList
  } from '@/api/system/api'
  import { FormInstance } from 'element-plus'
  import { ElMessageBox, ElMessage } from 'element-plus'
  import type { FormRules } from 'element-plus'
  import { ApiStatus } from '@/utils/http/status'
  import { useCheckedColumns } from '@/composables/useCheckedColumns'
  import { SearchFormItem } from '@/types'

  // 状态变量
  const dialogType = ref('add')
  const dialogVisible = ref(false)
  const loading = ref(false)
  const tableData = ref([])
  const pagination = reactive({
    currentPage: 1,
    pageSize: 10,
    total: 0
  })

  // 添加部门列表和角色列表的响应式数据
  const departmentList = ref<any[]>([])
  const roleList = ref<any[]>([])

  // 用户表单数据
  const formData = reactive({
    id: '',
    username: '',
    name: '',
    password: '',
    phone: '',
    gender: undefined,
    status: 1,
    department_id: undefined,
    role_id: undefined
  })

  // 搜索表单
  const searchForm = reactive({
    name: '',
    username: '',
    phone: '',
    department_id: undefined as undefined | number,
    role_id: undefined as undefined | number
  })

  // 搜索表单配置项
  const searchItems: SearchFormItem[] = [
    {
      label: '用户名',
      prop: 'name',
      type: 'input',
      elColSpan: 6, // 从8改为6，缩短显示宽度
      config: {
        clearable: true,
        placeholder: '请输入用户名'
      }
    },
    {
      label: '账号',
      prop: 'username',
      type: 'input',
      elColSpan: 6, // 从8改为6，缩短显示宽度
      config: {
        clearable: true,
        placeholder: '请输入账号'
      }
    },
    {
      label: '手机号',
      prop: 'phone',
      type: 'input',
      elColSpan: 6, // 从8改为6，缩短显示宽度
      config: {
        clearable: true,
        placeholder: '请输入手机号'
      }
    },
    {
      label: '部门',
      prop: 'department_id',
      type: 'select',
      elColSpan: 6, // 从8改为6，缩短显示宽度
      config: {
        clearable: true,
        placeholder: '请选择部门'
      },
      options: () =>
        departmentList.value.map((item) => ({
          label: item.name,
          value: item.id
        }))
    },
    {
      label: '角色',
      prop: 'role_id',
      type: 'select',
      elColSpan: 6, // 从8改为6，缩短显示宽度
      config: {
        clearable: true,
        placeholder: '请选择角色'
      },
      options: () =>
        roleList.value.map((item) => ({
          label: item.name,
          value: item.id
        }))
    }
  ]

  // 列配置选项
  const columnOptions = [
    { label: '用户名', prop: 'User.name' },
    { label: '账号', prop: 'User.username' },
    { label: '手机号', prop: 'User.phone' },
    { label: '性别', prop: 'User.gender' },
    { label: '部门', prop: 'department_name' },
    { label: '角色', prop: 'role_name' },
    { label: '状态', prop: 'User.status' },
    { label: '操作', prop: 'operation' }
  ]

  // 动态列配置
  const { columnChecks, columns } = useCheckedColumns(() => [
    {
      prop: 'User.name',
      label: '用户名',
      align: 'center'
    },
    {
      prop: 'User.username',
      label: '账号',
      align: 'center'
    },
    {
      prop: 'User.phone',
      label: '手机号',
      align: 'center'
    },
    {
      prop: 'User.gender',
      label: '性别',
      align: 'center'
    },
    {
      prop: 'department_name',
      label: '部门',
      align: 'center'
    },
    {
      prop: 'role_name',
      label: '角色',
      align: 'center'
    },
    {
      prop: 'User.status',
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

  // 表单实例引用
  const formRef = ref<FormInstance>()
  const searchFormRef = ref<FormInstance>()

  // 刷新表格数据
  const handleRefresh = () => {
    tableData.value = []
    loading.value = true
    loadUserList()
  }

  // 加载用户列表数据
  const loadUserList = async () => {
    loading.value = true
    try {
      // 构建搜索参数，过滤掉undefined和空字符串
      const params: any = {
        page: pagination.currentPage,
        pageSize: pagination.pageSize
      }

      // 只添加有值的搜索条件
      if (searchForm.name) params.name = searchForm.name
      if (searchForm.username) params.username = searchForm.username
      if (searchForm.phone) params.phone = searchForm.phone
      if (searchForm.department_id) params.department_id = searchForm.department_id
      if (searchForm.role_id) params.role_id = searchForm.role_id

      const res = await getUserList(params)
      if (res.code === 200) {
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
        ElMessage.error(res.message || '获取用户列表失败')
      }
    } catch (err) {
      console.error('获取用户列表出错:', err)
      ElMessage.error('获取用户列表失败')
    } finally {
      loading.value = false
    }
  }

  // 加载部门列表数据
  const loadDepartmentList = async () => {
    try {
      const res = await getDepartmentList()
      if (res.code === ApiStatus.success) {
        departmentList.value = res.data || []
      } else {
        ElMessage.error(res.message || '获取部门列表失败')
      }
    } catch (err) {
      console.error('获取部门列表出错:', err)
      ElMessage.error('获取部门列表失败')
    }
  }

  // 加载角色列表数据
  const loadRoleList = async () => {
    try {
      const res = await getRoleList()
      if (res.code === ApiStatus.success) {
        roleList.value = res.data || []
      } else {
        ElMessage.error(res.message || '获取角色列表失败')
      }
    } catch (err) {
      console.error('获取角色列表出错:', err)
      ElMessage.error('获取角色列表失败')
    }
  }

  // 页码变化处理
  const handleCurrentChange = (page: number) => {
    pagination.currentPage = page
    loadUserList()
  }

  // 每页条数变化处理
  const handleSizeChange = (size: number) => {
    pagination.pageSize = size
    pagination.currentPage = 1
    loadUserList()
  }

  // 搜索处理
  const search = () => {
    pagination.currentPage = 1
    loadUserList()
  }

  // 重置搜索
  const resetSearch = () => {
    searchForm.name = ''
    searchForm.username = ''
    searchForm.phone = ''
    searchForm.department_id = undefined
    searchForm.role_id = undefined
    pagination.currentPage = 1
    loadUserList()
  }

  // 显示对话框
  const showDialog = (type: string, row?: any) => {
    dialogVisible.value = true
    dialogType.value = type

    if (type === 'edit' && row) {
      formData.id = row.User.id
      formData.username = row.User.username || ''
      formData.name = row.User.name
      formData.phone = row.User.phone || ''
      formData.gender = row.User.gender === 0 ? 1 : row.User.gender // 如果性别是未知(0)，则默认设为男(1)
      formData.status = row.User.status
      formData.department_id = row.User.department_id
      formData.role_id = row.User.role_id || 1 // 获取角色ID，如果没有则默认为1
      formData.password = '' // 编辑模式下明确清空密码
    } else {
      // 添加用户时重置表单并确保状态为启用
      formData.id = ''
      formData.username = ''
      formData.name = ''
      formData.password = ''
      formData.phone = ''
      formData.gender = undefined
      formData.status = 1
      formData.department_id = undefined
      formData.role_id = undefined

      // 确保下一个渲染周期状态为启用
      nextTick(() => {
        formData.status = 1
      })
    }

    // 强制重新计算验证规则
    nextTick(() => {
      if (formRef.value) {
        formRef.value.clearValidate()
      }
    })
  }

  // 处理删除用户
  const handleDeleteUser = (row: any) => {
    ElMessageBox.confirm('确定要删除该用户吗？', '删除用户', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'error'
    })
      .then(async () => {
        try {
          // 确保用户ID正确传递
          const userId = row.User?.id
          if (!userId) {
            ElMessage.error('用户ID无效')
            return
          }

          const res = await apiDeleteUser(userId)
          if (res.code === 200) {
            ElMessage.success('删除用户成功')
            loadUserList()
          } else {
            ElMessage.error(res.message || '删除用户失败')
          }
        } catch (err) {
          console.error('删除用户出错:', err)
          ElMessage.error('删除用户失败，请稍后重试')
        }
      })
      .catch(() => {
        // 用户取消删除，不做处理
      })
  }

  const getTagType = (status: number) => {
    switch (status) {
      case 1:
        return 'primary'
      case 2:
        return 'warning'
      default:
        return 'info'
    }
  }

  const buildTagText = (status: number) => {
    if (status === 1) {
      return '启用'
    } else if (status === 2) {
      return '禁用'
    } else {
      return '未知'
    }
  }

  // 定义基本验证规则
  const baseRules = {
    username: [
      { required: true, message: '请输入账号', trigger: 'blur' },
      { min: 4, max: 20, message: '长度在 4 到 20 个字符', trigger: 'blur' }
    ],
    name: [
      { required: true, message: '请输入用户名', trigger: 'blur' },
      { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
    ],
    phone: [
      { required: true, message: '请输入手机号', trigger: 'blur' },
      { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式', trigger: 'blur' }
    ],
    gender: [{ required: true, message: '请选择性别', trigger: 'change' }],
    status: [{ required: true, message: '请选择状态', trigger: 'change' }],
    department_id: [{ required: true, message: '请选择部门', trigger: 'change' }],
    role_id: [{ required: true, message: '请选择角色', trigger: 'change' }]
  }

  // 根据对话框类型动态计算验证规则
  const computedRules = computed(() => {
    // 添加模式下的规则
    if (dialogType.value === 'add') {
      return {
        ...baseRules,
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
        ]
      }
    }
    // 编辑模式下的规则
    else {
      return {
        ...baseRules,
        password: [
          { required: false },
          {
            validator: (rule, value, callback) => {
              if (!value || value === '') {
                callback()
              } else if (value.length < 6 || value.length > 20) {
                callback(new Error('长度在 6 到 20 个字符'))
              } else {
                callback()
              }
            },
            trigger: 'blur'
          }
        ]
      }
    }
  })

  // 提交表单
  const handleSubmit = async () => {
    if (!formRef.value) return

    await formRef.value.validate(async (valid) => {
      if (valid) {
        try {
          const submitData = { ...formData }

          // 如果是编辑模式且密码为空，则删除密码字段
          if (dialogType.value === 'edit' && !submitData.password) {
            delete submitData.password
          }

          let res
          if (dialogType.value === 'add') {
            res = await addUser(submitData)
          } else {
            res = await updateUser(submitData)
          }

          if (res.code === 200) {
            ElMessage.success(dialogType.value === 'add' ? '添加成功' : '更新成功')
            dialogVisible.value = false
            loadUserList()
          } else {
            ElMessage.error(res.message || (dialogType.value === 'add' ? '添加失败' : '更新失败'))
          }
        } catch (err) {
          console.error('提交表单出错:', err)
          ElMessage.error(dialogType.value === 'add' ? '添加失败' : '更新失败')
        }
      }
    })
  }

  // 初始化加载数据
  onMounted(async () => {
    // 并行加载所有数据
    await Promise.all([loadUserList(), loadDepartmentList(), loadRoleList()])
  })

  // 创建安全部门列表和角色列表
  const safeDepartmentList = Array.isArray(departmentList.value) ? departmentList.value : [];
  const safeRoleList = Array.isArray(roleList.value) ? roleList.value : [];
  safeDepartmentList.map((item) => ({ ...item }));
  safeRoleList.map((item) => ({ ...item }));
</script>

<style lang="scss" scoped>
  .user-page {
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

    .user {
      .avatar {
        width: 40px;
        height: 40px;
        border-radius: 8px;
        transition: all 0.3s ease;

        &:hover {
          transform: scale(1.1);
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        }
      }

      > div {
        margin-left: 10px;

        .user-name {
          font-weight: 500;
          color: var(--art-text-gray-800);
        }
      }
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
      .el-input-number,
      .el-select {
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

      .el-select {
        .el-select__wrapper {
          border-radius: 8px;
        }
      }

      .el-radio-group {
        .el-radio {
          margin-right: 16px;
          transition: all 0.2s ease;

          &:hover {
            transform: scale(1.05);
          }
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

  .status-hint {
    margin-left: 8px;
    font-size: 12px;
    color: #909399;
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
    .user-page {
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

      .user {
        flex-direction: column;
        align-items: center;
        text-align: center;

        > div {
          margin-left: 0;
          margin-top: 8px;
        }
      }
    }
  }
</style>
