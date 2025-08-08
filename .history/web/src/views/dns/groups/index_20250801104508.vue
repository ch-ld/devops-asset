<template>
  <div class="domain-group-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title">
          <h1>域名分组</h1>
          <p>管理域名分组，用于组织和管理域名</p>
        </div>
        <div class="page-actions">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加分组
          </el-button>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-container">
      <el-card shadow="never">
        <el-form
          ref="searchFormRef"
          :model="searchForm"
          :inline="true"
          class="search-form"
        >
          <el-form-item label="分组名称" prop="keyword">
            <el-input
              v-model="searchForm.keyword"
              placeholder="请输入分组名称"
              clearable
              style="width: 200px"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button @click="handleReset">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>

    <!-- 分组列表 -->
    <div class="table-container">
      <el-card shadow="never">
        <div class="table-header">
          <div class="table-title">
            <h3>分组列表</h3>
            <span>共 {{ pagination.total }} 个分组</span>
          </div>
          <div class="table-actions">
            <el-button 
              type="danger" 
              :disabled="!hasSelected"
              @click="handleBatchDelete"
            >
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
          </div>
        </div>

        <el-table
          v-loading="loading"
          :data="groups"
          row-key="id"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="50" />
          <el-table-column prop="name" label="分组名称" min-width="200">
            <template #default="{ row }">
              <div class="group-name">
                <el-icon><Folder /></el-icon>
                <span>{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="描述" min-width="300" />
          <el-table-column prop="domain_count" label="域名数量" width="120">
            <template #default="{ row }">
              <el-tag type="info">{{ row.domain_count || 0 }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="sort" label="排序" width="100" />
          <el-table-column prop="created_at" label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleEdit(row)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button type="danger" link @click="handleDelete(row)">
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="table-pagination">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="fetchGroups"
            @current-change="fetchGroups"
          />
        </div>
      </el-card>
    </div>

    <!-- 分组表单弹窗 -->
    <GroupModal
      v-model:visible="modalVisible"
      :group="currentGroup"
      :mode="modalMode"
      @success="handleModalSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, Search, Refresh, Delete, Folder, Edit 
} from '@element-plus/icons-vue'
import GroupModal from './components/GroupModal.vue'
import { domainGroupApi } from '@/api/dns/domainGroup'
import type { DomainGroup } from '@/types/dns'

// 响应式数据
const loading = ref(false)
const modalVisible = ref(false)
const modalMode = ref<'add' | 'edit'>('add')
const currentGroup = ref<DomainGroup | null>(null)
const groups = ref<DomainGroup[]>([])

// 搜索表单
const searchFormRef = ref()
const searchForm = reactive({
  keyword: ''
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 选择
const selectedGroups = ref<DomainGroup[]>([])
const hasSelected = computed(() => selectedGroups.value.length > 0)

// 获取分组列表
const fetchGroups = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      size: pagination.pageSize,
      keyword: searchForm.keyword
    }
    
    const response = await domainGroupApi.list(params)
    groups.value = response.list || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取分组列表失败:', error)
    ElMessage.error('获取分组列表失败')
  } finally {
    loading.value = false
  }
}

// 事件处理
const handleAdd = () => {
  modalMode.value = 'add'
  currentGroup.value = null
  modalVisible.value = true
}

const handleEdit = (row: DomainGroup) => {
  modalMode.value = 'edit'
  currentGroup.value = row
  modalVisible.value = true
}

const handleDelete = async (row: DomainGroup) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除分组"${row.name}"吗？删除后该分组下的域名将移动到默认分组。`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await domainGroupApi.delete(row.id)
    ElMessage.success('删除成功')
    await fetchGroups()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

const handleBatchDelete = async () => {
  if (!hasSelected.value) return
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedGroups.value.length} 个分组吗？`,
      '批量删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const ids = selectedGroups.value.map(item => item.id)
    await domainGroupApi.batchDelete(ids)
    ElMessage.success('批量删除成功')
    await fetchGroups()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
      ElMessage.error('批量删除失败')
    }
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchGroups()
}

const handleReset = () => {
  searchFormRef.value?.resetFields()
  pagination.page = 1
  fetchGroups()
}

const handleSelectionChange = (selection: DomainGroup[]) => {
  selectedGroups.value = selection
}

const handleModalSuccess = () => {
  modalVisible.value = false
  fetchGroups()
}

// 工具方法
const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 初始化
onMounted(() => {
  fetchGroups()
})
</script>

<style scoped lang="scss">
.domain-group-container {
  padding: 16px;
}

.page-header {
  background: #fff;
  padding: 16px 24px;
  margin-bottom: 16px;
  border-radius: 8px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03);

  .page-header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .page-title {
      h1 {
        margin: 0 0 8px 0;
        font-size: 20px;
        font-weight: 600;
        color: #1f2937;
      }

      p {
        margin: 0;
        color: #6b7280;
        font-size: 14px;
      }
    }
  }
}

.search-container, .table-container {
  margin-bottom: 16px;

  .search-form {
    margin-bottom: 0;
  }
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  .table-title {
    h3 {
      margin: 0 8px 0 0;
      font-size: 16px;
      font-weight: 600;
      color: #1f2937;
    }

    span {
      color: #6b7280;
      font-size: 14px;
    }
  }

  .table-actions {
    display: flex;
    gap: 8px;
  }
}

.table-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.group-name {
  display: flex;
  align-items: center;
  gap: 8px;

  .el-icon {
    color: #6b7280;
  }
}
</style> 
