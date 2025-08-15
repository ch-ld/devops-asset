<template>
  <div class="dns-domain-detail">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="breadcrumb-section">
          <el-button 
            type="primary" 
            link 
            @click="handleGoBack"
            class="back-button"
          >
            <el-icon><ArrowLeft /></el-icon>
            {{ domainInfo.name }}
          </el-button>
          <div class="domain-status">
            <el-tag 
              :type="getDomainStatusType(domainInfo.status)"
              size="small"
            >
              {{ getDomainStatusText(domainInfo.status) }}
            </el-tag>
          </div>
        </div>
        
        <!-- 同步状态提示 -->
        <div v-if="syncStatus.is_syncing" class="sync-notice">
          <el-alert
            title="域名解析记录正在同步中，请稍后查看最新状态"
            type="success"
            :closable="false"
            show-icon
          />
        </div>
      </div>
    </div>

    <!-- 操作按钮区域 -->
    <div class="action-section">
      <div class="action-content">
        <div class="action-left">
          <el-button type="primary" @click="handleAddRecord">
            <el-icon><Plus /></el-icon>
            添加记录
          </el-button>
          <el-button @click="handleBatchImport">
            <el-icon><Upload /></el-icon>
            批量导入
          </el-button>
          <el-button @click="handleBatchExport">
            <el-icon><Download /></el-icon>
            批量导出
          </el-button>
          <el-button @click="handleSyncRecords" :loading="syncLoading">
            <el-icon><Refresh /></el-icon>
            同步记录
          </el-button>
        </div>
        
        <div class="action-right">
          <el-input
            v-model="searchKeyword"
            placeholder="请输入主机记录"
            style="width: 200px"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #suffix>
              <el-icon class="search-icon" @click="handleSearch"><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </div>
    </div>

    <!-- 筛选区域 -->
    <div class="filter-section">
      <div class="filter-content">
        <div class="filter-tabs">
          <el-radio-group v-model="filters.recordType" @change="handleFilterChange">
            <el-radio-button label="">全部记录</el-radio-button>
            <el-radio-button label="A">A</el-radio-button>
            <el-radio-button label="CNAME">CNAME</el-radio-button>
            <el-radio-button label="MX">MX</el-radio-button>
            <el-radio-button label="TXT">TXT</el-radio-button>
            <el-radio-button label="NS">NS</el-radio-button>
            <el-radio-button label="AAAA">AAAA</el-radio-button>
          </el-radio-group>
        </div>
        
        <div class="filter-selects">
          <el-select v-model="filters.line" placeholder="全部线路" clearable style="width: 120px">
            <el-option label="全部线路" value="" />
            <el-option label="默认" value="default" />
            <el-option label="电信" value="telecom" />
            <el-option label="联通" value="unicom" />
            <el-option label="移动" value="mobile" />
            <el-option label="海外" value="overseas" />
          </el-select>
          
          <el-select v-model="filters.status" placeholder="解析状态" clearable style="width: 120px">
            <el-option label="全部状态" value="" />
            <el-option label="正常" value="enabled" />
            <el-option label="暂停" value="disabled" />
          </el-select>
        </div>
      </div>
    </div>

    <!-- 解析记录表格 -->
    <div class="table-section">
      <el-table
        :data="recordList"
        :loading="loading"
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="name" label="主机记录" min-width="120">
          <template #default="{ row }">
            <div class="record-name">
              <span class="name-text">{{ row.name || '@' }}</span>
              <el-tooltip v-if="row.name === '@'" content="@表示直接解析主域名">
                <el-icon class="info-icon"><QuestionFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="type" label="记录类型" width="80" align="center">
          <template #default="{ row }">
            <el-tag size="small" :type="getRecordTypeColor(row.type)">
              {{ row.type }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="解析请求来源" width="120" align="center">
          <template #default="{ row }">
            <span>{{ getLineText(row.line) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="value" label="记录值" min-width="200">
          <template #default="{ row }">
            <div class="record-value">
              <span class="value-text">{{ row.value }}</span>
              <el-button 
                type="primary" 
                link 
                size="small"
                @click="handleCopyValue(row.value)"
              >
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="mx_priority" label="MX优先级" width="100" align="center">
          <template #default="{ row }">
            <span>{{ row.type === 'MX' ? (row.mx_priority || '-') : '-' }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="ttl" label="TTL" width="80" align="center">
          <template #default="{ row }">
            <span>{{ row.ttl || 600 }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="weight" label="权重" width="80" align="center">
          <template #default="{ row }">
            <span>{{ row.weight || '-' }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.enabled"
              :loading="row.statusLoading"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        
        <el-table-column label="更新时间(UTC+8)" width="160" align="center">
          <template #default="{ row }">
            <span class="update-time">{{ formatTime(row.updated_at) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button 
                type="primary" 
                link 
                size="small"
                @click="handleEditRecord(row)"
              >
                修改
              </el-button>
              <el-button 
                type="primary" 
                link 
                size="small"
                @click="handleDeleteRecord(row)"
              >
                删除
              </el-button>
              <el-dropdown @command="(command) => handleMoreAction(command, row)">
                <el-button type="primary" link size="small">
                  更多
                  <el-icon><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="copy">复制记录</el-dropdown-item>
                    <el-dropdown-item command="test">解析测试</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 批量操作和分页 -->
      <div class="table-footer">
        <div class="batch-actions">
          <el-button 
            :disabled="selectedRecords.length === 0"
            @click="handleBatchEnable"
          >
            批量启用
          </el-button>
          <el-button 
            :disabled="selectedRecords.length === 0"
            @click="handleBatchDisable"
          >
            批量暂停
          </el-button>
          <el-button 
            :disabled="selectedRecords.length === 0"
            @click="handleBatchDelete"
          >
            批量删除
          </el-button>
        </div>
        
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 添加/编辑记录模态框 -->
    <AliRecordModal
      v-model:visible="recordModalVisible"
      :domain-id="domainId"
      :domain-name="domainInfo.name"
      :record="currentRecord"
      @success="handleRecordSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft,
  Plus,
  Upload,
  Download,
  Refresh,
  Search,
  QuestionFilled,
  CopyDocument,
  ArrowDown
} from '@element-plus/icons-vue'
import { getRecords, deleteRecord, batchSyncRecords } from '@/api/dns/record'
import { getDomains } from '@/api/dns'
import AliRecordModal from './components/AliRecordModal.vue'

const route = useRoute()
const router = useRouter()

// 获取域名ID
const domainId = computed(() => parseInt(route.params.id as string))

// 响应式数据
const loading = ref(false)
const syncLoading = ref(false)
const recordList = ref([])
const selectedRecords = ref([])
const searchKeyword = ref('')
const recordModalVisible = ref(false)
const currentRecord = ref(null)

// 域名信息
const domainInfo = ref({
  name: '',
  status: 'active'
})

// 同步状态
const syncStatus = ref({
  is_syncing: false
})

// 筛选条件
const filters = reactive({
  recordType: '',
  line: '',
  status: ''
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 获取域名信息
const fetchDomainInfo = async () => {
  try {
    const response = await getDomains({ 
      page: 1, 
      page_size: 1,
      keyword: domainId.value.toString() 
    })
    if (response.items && response.items.length > 0) {
      domainInfo.value = response.items[0]
    }
  } catch (error) {
    ElMessage.error('获取域名信息失败')
  }
}

// 获取解析记录列表
const fetchRecords = async () => {
  loading.value = true
  try {
    const params = {
      domain_id: domainId.value,
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchKeyword.value,
      type: filters.recordType,
      line: filters.line,
      status: filters.status
    }
    
    const response = await getRecords(params)
    recordList.value = response.items || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('获取解析记录失败')
  } finally {
    loading.value = false
  }
}

// 工具函数
const getDomainStatusType = (status: string) => {
  const types = {
    active: 'success',
    paused: 'warning',
    error: 'danger'
  }
  return types[status as keyof typeof types] || 'success'
}

const getDomainStatusText = (status: string) => {
  const texts = {
    active: '正常',
    paused: '暂停',
    error: '异常'
  }
  return texts[status as keyof typeof texts] || '正常'
}

const getRecordTypeColor = (type: string) => {
  const colors = {
    A: 'success',
    CNAME: 'primary',
    MX: 'warning',
    TXT: 'info',
    NS: 'danger',
    AAAA: 'success'
  }
  return colors[type] || ''
}

const getLineText = (line) => {
  const texts = {
    default: '默认',
    telecom: '电信',
    unicom: '联通',
    mobile: '移动',
    overseas: '海外'
  }
  return texts[line] || '默认'
}

const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}

// 事件处理
const handleGoBack = () => {
  // 返回DNS记录页面，并传递域名ID参数以便恢复状态
  router.push({
    path: '/dns/records',
    query: { 
      domainId: domainId.value.toString(),
      viewMode: 'list' // 强制切换到列表模式
    }
  })
}

const handleSearch = () => {
  pagination.page = 1
  fetchRecords()
}

const handleFilterChange = () => {
  pagination.page = 1
  fetchRecords()
}

const handleSelectionChange = (selection) => {
  selectedRecords.value = selection
}

const handleAddRecord = () => {
  currentRecord.value = null
  recordModalVisible.value = true
}

const handleEditRecord = (record) => {
  currentRecord.value = record
  recordModalVisible.value = true
}

const handleDeleteRecord = async (record) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除解析记录 ${record.name}.${domainInfo.value.name} 吗？`,
      '确认删除',
      { type: 'warning' }
    )
    
    await deleteRecord(record.id)
    ElMessage.success('删除成功')
    fetchRecords()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleStatusChange = async (record) => {
  record.statusLoading = true
  try {
    // 调用状态更新API
    ElMessage.success(record.enabled ? '启用成功' : '暂停成功')
  } catch (error) {
    record.enabled = !record.enabled // 回滚状态
    ElMessage.error('状态更新失败')
  } finally {
    record.statusLoading = false
  }
}

const handleCopyValue = async (value) => {
  try {
    await navigator.clipboard.writeText(value)
    ElMessage.success('复制成功')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

const handleSyncRecords = async () => {
  syncLoading.value = true
  try {
    await batchSyncRecords({
      domain_id: domainId.value,
      provider_id: domainInfo.value.registrar_id,
      dry_run: false
    })
    ElMessage.success('同步成功')
    fetchRecords()
  } catch (error) {
    ElMessage.error('同步失败')
  } finally {
    syncLoading.value = false
  }
}

const handleRecordSuccess = () => {
  recordModalVisible.value = false
  fetchRecords()
}

const handleMoreAction = (command, record) => {
  switch (command) {
    case 'copy':
      ElMessage.info('复制记录功能开发中...')
      break
    case 'test':
      ElMessage.info('解析测试功能开发中...')
      break
  }
}

const handleBatchImport = () => {
  ElMessage.info('批量导入功能开发中...')
}

const handleBatchExport = () => {
  ElMessage.info('批量导出功能开发中...')
}

const handleBatchEnable = () => {
  ElMessage.info('批量启用功能开发中...')
}

const handleBatchDisable = () => {
  ElMessage.info('批量暂停功能开发中...')
}

const handleBatchDelete = () => {
  ElMessage.info('批量删除功能开发中...')
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchRecords()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  fetchRecords()
}

// 初始化
onMounted(() => {
  fetchDomainInfo()
  fetchRecords()
})
</script>

<style scoped lang="scss">
.dns-domain-detail {
  background: #f5f5f5;
  min-height: 100vh;

  .page-header {
    background: white;
    padding: 16px 24px;
    border-bottom: 1px solid #e8e8e8;

    .header-content {
      max-width: 1200px;
      margin: 0 auto;

      .breadcrumb-section {
        display: flex;
        align-items: center;
        gap: 12px;
        margin-bottom: 16px;

        .back-button {
          font-size: 16px;
          font-weight: 500;
          padding: 0;

          .el-icon {
            margin-right: 4px;
          }
        }

        .domain-status {
          margin-left: auto;
        }
      }

      .sync-notice {
        margin-top: 12px;
      }
    }
  }

  .action-section {
    background: white;
    padding: 16px 24px;
    border-bottom: 1px solid #e8e8e8;

    .action-content {
      max-width: 1200px;
      margin: 0 auto;
      display: flex;
      justify-content: space-between;
      align-items: center;

      .action-left {
        display: flex;
        gap: 8px;
      }

      .action-right {
        .search-icon {
          cursor: pointer;
          color: #8c8c8c;

          &:hover {
            color: #1890ff;
          }
        }
      }
    }
  }

  .filter-section {
    background: white;
    padding: 12px 24px;
    border-bottom: 1px solid #e8e8e8;

    .filter-content {
      max-width: 1200px;
      margin: 0 auto;
      display: flex;
      justify-content: space-between;
      align-items: center;

      .filter-tabs {
        .el-radio-group {
          .el-radio-button {
            margin-right: 0;
          }
        }
      }

      .filter-selects {
        display: flex;
        gap: 12px;
      }
    }
  }

  .table-section {
    background: white;
    padding: 24px;

    .el-table {
      max-width: 1200px;
      margin: 0 auto;
    }

    .record-name {
      display: flex;
      align-items: center;
      gap: 4px;

      .name-text {
        font-weight: 500;
      }

      .info-icon {
        color: #8c8c8c;
        font-size: 12px;
      }
    }

    .record-value {
      display: flex;
      align-items: center;
      gap: 8px;

      .value-text {
        flex: 1;
        word-break: break-all;
      }
    }

    .update-time {
      font-size: 12px;
      color: #8c8c8c;
    }

    .action-buttons {
      display: flex;
      gap: 4px;
    }

    .table-footer {
      margin-top: 24px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      max-width: 1200px;
      margin-left: auto;
      margin-right: auto;

      .batch-actions {
        display: flex;
        gap: 8px;
      }
    }
  }
}

// 响应式设计
@media (max-width: 1200px) {
  .dns-domain-detail {
    .page-header .header-content,
    .action-section .action-content,
    .filter-section .filter-content,
    .table-section .el-table,
    .table-section .table-footer {
      max-width: 100%;
      padding: 0 16px;
    }
  }
}

@media (max-width: 768px) {
  .dns-domain-detail {
    .page-header {
      padding: 12px 16px;
    }

    .action-section {
      padding: 12px 16px;

      .action-content {
        flex-direction: column;
        gap: 12px;
        align-items: stretch;

        .action-left {
          justify-content: center;
          flex-wrap: wrap;
        }
      }
    }

    .filter-section {
      padding: 12px 16px;

      .filter-content {
        flex-direction: column;
        gap: 12px;
        align-items: stretch;

        .filter-tabs {
          .el-radio-group {
            display: flex;
            flex-wrap: wrap;
            gap: 4px;
          }
        }

        .filter-selects {
          justify-content: space-between;
        }
      }
    }

    .table-section {
      padding: 16px;

      .table-footer {
        flex-direction: column;
        gap: 16px;
        align-items: stretch;

        .batch-actions {
          justify-content: center;
          flex-wrap: wrap;
        }
      }
    }
  }
}
</style>
