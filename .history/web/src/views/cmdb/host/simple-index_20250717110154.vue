<template>
  <div class="host-management-page">
    <!-- 错误提示组件 -->
    <el-alert
      v-if="loadError"
      type="error"
      title="数据加载失败"
      description="主机数据加载失败，请稍后重试或联系管理员"
      show-icon
      :closable="true"
      @close="loadError = false"
      class="error-alert"
    >
      <template #default>
        <el-button @click="fetchHostsData" type="primary" size="small">
          <el-icon><ReloadIcon /></el-icon>
          重试
        </el-button>
      </template>
    </el-alert>
    
    <!-- 顶部工具栏 -->
    <el-card shadow="never" class="toolbar-card">
      <el-row :gutter="16">
        <el-col :lg="6" :md="8" :sm="12" :xs="24">
          <el-input
            v-model="searchParams.keyword"
            placeholder="搜索主机名、IP地址"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #suffix>
              <el-icon @click="handleSearch"><Search /></el-icon>
            </template>
          </el-input>
        </el-col>

        <el-col :lg="4" :md="6" :sm="12" :xs="12">
          <el-select
            v-model="searchParams.group_id"
            placeholder="主机组"
            style="width: 100%"
            clearable
            @change="handleSearch"
          >
            <el-option :value="undefined" label="全部主机组" />
            <el-option
              v-for="group in hostGroupOptions"
              :key="group.value"
              :value="group.value"
              :label="group.label"
            />
          </el-select>
        </el-col>
        </el-row>
    </el-card>

    <!-- 主机列表 -->
    <el-card shadow="never" class="host-list-card">
      <template #header>
        <div class="card-header">
          <span>主机列表</span>
          <el-tag type="info" style="margin-left: 8px">{{ pagination.total }} 台</el-tag>
          <div class="header-operations">
            <el-button-group>
              <el-tooltip content="导出">
                <el-button @click="handleExportHosts">
                  <el-icon><Download /></el-icon>
                </el-button>
              </el-tooltip>
              <el-tooltip content="设置">
                <el-button @click="handleColumnSettings">
                  <el-icon><Setting /></el-icon>
                </el-button>
              </el-tooltip>
            </el-button-group>
          </div>
        </div>
      </template>

      <!-- 使用骨架屏提升加载体验 -->
      <el-skeleton :loading="loading" animated :count="3" :throttle="500">
        <template #default>
          <el-table
            :data="hostList || []"
            v-loading="loading"
            @selection-change="onSelectionChange"
            style="width: 100%"
          >
            <!-- 表格列内容 -->
          </el-table>
        </template>
      </el-skeleton>

      <!-- 分页组件 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="searchParams.page"
          v-model:page-size="searchParams.page_size"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="pagination.total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
  import { ref, reactive, computed, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
  import { useHostStore } from '@/store/modules/host'
  import { storeToRefs } from 'pinia'
  import { Search } from '@element-plus/icons-vue'

  const router = useRouter()
  const hostStore = useHostStore()
  const { hostList, hostGroupTree, pagination, isLoading } = storeToRefs(hostStore)
  const loading = computed(() => isLoading.value)
  const loadError = ref(false)

  // 搜索参数
  const searchParams = reactive({
    keyword: '',
    group_id: undefined,
    status: undefined,
    region: undefined,
    page: 1,
    page_size: 10
  })

  // 主机组选项
  const hostGroupOptions = computed(() => {
    try {
      const options = []
      const processGroups = (groups, parentPath = '') => {
        if (!groups || !Array.isArray(groups)) return;
        
        groups.forEach((group) => {
          if (!group) return; // 跳过空组
          
          const label = parentPath ? `${parentPath} / ${group.name}` : group.name
          options.push({ label, value: group.id })

          // 查找子组
          const children = (hostGroupTree.value || []).filter((g) => g && g.parent_id === group.id)
          if (children.length > 0) {
            processGroups(children, label)
          }
        })
      }

      // 从根节点开始处理
      const rootGroups = (hostGroupTree.value || []).filter((g) => g && !g.parent_id)
      processGroups(rootGroups)

      return options
    } catch (error) {
      console.error('生成主机组选项时出错:', error)
      return [] // 发生错误时返回空数组
    }
  })

  // 初始化数据
  onMounted(() => {
    fetchHostsData()
  })

  // 方法
  const fetchHostsData = async () => {
    try {
      loadError.value = false
      // 同时获取主机组和主机列表
      await Promise.all([
        hostStore.fetchHostGroupTree(),
        fetchHosts()
      ])
    } catch (error) {
      loadError.value = true
      ElMessage.error('初始化数据失败')
      console.error(error)
    }
  }

  const fetchHosts = async () => {
    try {
      await hostStore.fetchHosts({
        ...searchParams,
        page: searchParams.page,
        page_size: searchParams.page_size
      })
    } catch (error) {
      loadError.value = true
      ElMessage.error('获取主机列表失败')
      console.error(error)
    }
  }

  // 表格相关
  const selectedRowKeys = ref([])
  const selectedRows = ref([])

  const handleSearch = () => {
    searchParams.page = 1
    fetchHosts()
  }

  const handleSizeChange = (size) => {
    searchParams.page_size = size
    fetchHosts()
  }

  const handleCurrentChange = (page) => {
    searchParams.page = page
    fetchHosts()
  }

  const onSelectionChange = (rows) => {
    try {
      selectedRows.value = rows || []
      selectedRowKeys.value = (rows || []).map(row => row?.id).filter(Boolean)
    } catch (error) {
      console.error('选择行变化处理错误:', error)
      selectedRows.value = []
      selectedRowKeys.value = []
    }
  }

  // 导出主机
  const handleExportHosts = () => {
    ElMessage.info('导出功能待实现')
  }

  // 列设置
  const handleColumnSettings = () => {
    ElMessage.info('列设置功能待实现')
  }
</script>

<style lang="scss" scoped>
  .host-management-page {
    .error-alert {
      margin-bottom: 16px;
    }
    
    .toolbar-card {
      margin-bottom: 16px;
    }

    .host-list-card {
      .card-header {
        display: flex;
        align-items: center;

        .header-operations {
          margin-left: auto;
          display: flex;
          gap: 8px;
        }
      }

      .pagination-container {
        margin-top: 20px;
        display: flex;
        justify-content: flex-end;
      }
    }
  }
</style>
