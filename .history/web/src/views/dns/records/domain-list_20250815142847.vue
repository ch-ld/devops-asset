<template>
  <div class="dns-domain-list">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">公网权威解析</h1>
          <p class="page-description">管理您的域名DNS解析记录，支持多种记录类型和智能解析</p>
        </div>
        
        <!-- 通知横幅 -->
        <div class="notice-banners">
          <div class="notice-banner info">
            <el-icon><InfoFilled /></el-icon>
            <span>【域名实名认证】全面开启域名实名认证，未实名认证域名将被暂停解析</span>
          </div>
          <div class="notice-banner warning">
            <el-icon><WarningFilled /></el-icon>
            <span>【DNS安全增强】新增DNS安全防护功能，有效防护DNS劫持和污染攻击</span>
          </div>
          <div class="notice-banner success">
            <el-icon><SuccessFilled /></el-icon>
            <span>【解析优化】DNS解析全面升级，解析速度提升30%，稳定性大幅提升</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 搜索和操作区域 -->
    <div class="search-section">
      <div class="search-controls">
        <el-button type="primary" @click="handleAddDomain">
          <el-icon><Plus /></el-icon>
          添加域名
        </el-button>
        
        <div class="search-filters">
          <el-select v-model="filters.status" placeholder="解析状态" clearable style="width: 120px">
            <el-option label="全部状态" value="" />
            <el-option label="正常" value="normal" />
            <el-option label="暂停" value="paused" />
            <el-option label="异常" value="error" />
          </el-select>
          
          <el-select v-model="filters.type" placeholder="付费版本" clearable style="width: 120px">
            <el-option label="全部版本" value="" />
            <el-option label="免费版" value="free" />
            <el-option label="个人版" value="personal" />
            <el-option label="企业版" value="enterprise" />
          </el-select>
          
          <el-input
            v-model="filters.keyword"
            placeholder="请输入域名进行搜索"
            style="width: 300px"
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

    <!-- 域名列表表格 -->
    <div class="table-section">
      <el-table
        :data="domainList"
        :loading="loading"
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="name" label="域名" min-width="200">
          <template #default="{ row }">
            <div class="domain-cell">
              <el-button 
                type="primary" 
                link 
                @click="handleViewRecords(row)"
                class="domain-name"
              >
                {{ row.name }}
              </el-button>
              <div class="domain-tags" v-if="row.tags && row.tags.length">
                <el-tag 
                  v-for="tag in row.tags" 
                  :key="tag" 
                  size="small" 
                  class="domain-tag"
                >
                  {{ tag }}
                </el-tag>
              </div>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="record_count" label="记录数" width="80" align="center">
          <template #default="{ row }">
            <span class="record-count">{{ row.record_count || 0 }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="DNS解析数量" width="120" align="center">
          <template #default="{ row }">
            <div class="dns-stats">
              <el-icon class="stats-icon"><TrendCharts /></el-icon>
              <span>{{ formatDnsCount(row.dns_query_count) }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="version" label="付费版本" width="100" align="center">
          <template #default="{ row }">
            <el-tag 
              :type="getVersionTagType(row.version)" 
              size="small"
            >
              {{ getVersionText(row.version) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="DNS解析套餐状态" width="150" align="center">
          <template #default="{ row }">
            <div class="package-status">
              <el-icon 
                :class="['status-icon', getStatusClass(row.status)]"
              >
                <CircleFilled />
              </el-icon>
              <span>{{ getStatusText(row.status) }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="quantity" label="数量" width="80" align="center">
          <template #default="{ row }">
            <span>{{ row.quantity || '-' }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button 
                type="primary" 
                link 
                @click="handleViewRecords(row)"
              >
                解析设置
              </el-button>
              <el-button 
                type="primary" 
                link 
                @click="handleManageDomain(row)"
              >
                域名设置
              </el-button>
              <el-button 
                type="primary" 
                link 
                @click="handleUpgrade(row)"
              >
                升级
              </el-button>
              <el-dropdown @command="(command) => handleMoreAction(command, row)">
                <el-button type="primary" link>
                  更多
                  <el-icon><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="whois">WHOIS查询</el-dropdown-item>
                    <el-dropdown-item command="export">导出记录</el-dropdown-item>
                    <el-dropdown-item command="delete" divided>删除域名</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-section">
        <div class="batch-actions">
          <el-button 
            :disabled="selectedDomains.length === 0"
            @click="handleBatchDelete"
          >
            批量删除
          </el-button>
          <el-button 
            :disabled="selectedDomains.length === 0"
            @click="handleBatchExport"
          >
            批量导出
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Search,
  InfoFilled,
  WarningFilled,
  SuccessFilled,
  TrendCharts,
  CircleFilled,
  ArrowDown
} from '@element-plus/icons-vue'
import { getDomains, recordApi } from '@/api/dns'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const domainList = ref([])
const selectedDomains = ref([])

// 搜索筛选
const filters = reactive({
  keyword: '',
  status: '',
  type: ''
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 获取域名列表
const fetchDomains = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: filters.keyword,
      status: filters.status,
      type: filters.type
    }
    
    const response = await getDomains(params)
    domainList.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('获取域名列表失败')
  } finally {
    loading.value = false
  }
}

// 格式化DNS查询数量
const formatDnsCount = (count) => {
  if (!count) return '0'
  if (count >= 10000) return `${Math.floor(count / 10000)}万`
  if (count >= 1000) return `${Math.floor(count / 1000)}k`
  return count.toString()
}

// 获取版本标签类型
const getVersionTagType = (version) => {
  const types = {
    free: '',
    personal: 'warning',
    enterprise: 'success'
  }
  return types[version] || ''
}

// 获取版本文本
const getVersionText = (version) => {
  const texts = {
    free: '免费版',
    personal: '个人版',
    enterprise: '企业版'
  }
  return texts[version] || '免费版'
}

// 获取状态样式类
const getStatusClass = (status) => {
  const classes = {
    normal: 'status-normal',
    paused: 'status-paused',
    error: 'status-error'
  }
  return classes[status] || 'status-normal'
}

// 获取状态文本
const getStatusText = (status) => {
  const texts = {
    normal: '正常',
    paused: '暂停',
    error: '异常'
  }
  return texts[status] || '正常'
}

// 事件处理
const handleSearch = () => {
  pagination.page = 1
  fetchDomains()
}

const handleSelectionChange = (selection) => {
  selectedDomains.value = selection
}

const handleViewRecords = (domain) => {
  router.push(`/dns/domains/${domain.id}/records`)
}

const handleAddDomain = () => {
  // 添加域名逻辑
  ElMessage.info('添加域名功能开发中...')
}

const handleManageDomain = (domain) => {
  // 域名管理逻辑
  ElMessage.info('域名管理功能开发中...')
}

const handleUpgrade = (domain) => {
  // 升级逻辑
  ElMessage.info('升级功能开发中...')
}

const handleMoreAction = (command, domain) => {
  switch (command) {
    case 'whois':
      ElMessage.info('WHOIS查询功能开发中...')
      break
    case 'export':
      ElMessage.info('导出记录功能开发中...')
      break
    case 'delete':
      handleDeleteDomain(domain)
      break
  }
}

const handleDeleteDomain = async (domain) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除域名 ${domain.name} 吗？此操作不可恢复。`,
      '确认删除',
      { type: 'warning' }
    )
    // 删除逻辑
    ElMessage.success('删除成功')
    fetchDomains()
  } catch (error) {
    // 用户取消
  }
}

const handleBatchDelete = () => {
  ElMessage.info('批量删除功能开发中...')
}

const handleBatchExport = () => {
  ElMessage.info('批量导出功能开发中...')
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchDomains()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  fetchDomains()
}

// 初始化
onMounted(() => {
  fetchDomains()
})
</script>

<style scoped lang="scss">
.dns-domain-list {
  padding: 0;
  background: #f5f5f5;
  min-height: 100vh;

  .page-header {
    background: white;
    padding: 24px;
    margin-bottom: 16px;
    border-bottom: 1px solid #e8e8e8;

    .header-content {
      max-width: 1200px;
      margin: 0 auto;
    }

    .title-section {
      margin-bottom: 20px;

      .page-title {
        font-size: 24px;
        font-weight: 600;
        color: #262626;
        margin: 0 0 8px 0;
      }

      .page-description {
        color: #8c8c8c;
        font-size: 14px;
        margin: 0;
      }
    }

    .notice-banners {
      display: flex;
      gap: 12px;
      flex-wrap: wrap;

      .notice-banner {
        flex: 1;
        min-width: 300px;
        padding: 12px 16px;
        border-radius: 6px;
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 13px;

        &.info {
          background: #e6f7ff;
          border: 1px solid #91d5ff;
          color: #0050b3;
        }

        &.warning {
          background: #fff7e6;
          border: 1px solid #ffd591;
          color: #d46b08;
        }

        &.success {
          background: #f6ffed;
          border: 1px solid #b7eb8f;
          color: #389e0d;
        }
      }
    }
  }

  .search-section {
    background: white;
    padding: 16px 24px;
    margin-bottom: 16px;

    .search-controls {
      max-width: 1200px;
      margin: 0 auto;
      display: flex;
      justify-content: space-between;
      align-items: center;

      .search-filters {
        display: flex;
        gap: 12px;
        align-items: center;

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

  .table-section {
    background: white;
    padding: 0 24px 24px;

    .el-table {
      max-width: 1200px;
      margin: 0 auto;
    }

    .domain-cell {
      .domain-name {
        font-weight: 500;
        font-size: 14px;
      }

      .domain-tags {
        margin-top: 4px;
        display: flex;
        gap: 4px;

        .domain-tag {
          font-size: 12px;
        }
      }
    }

    .record-count {
      font-weight: 500;
      color: #1890ff;
    }

    .dns-stats {
      display: flex;
      align-items: center;
      gap: 4px;
      justify-content: center;

      .stats-icon {
        color: #52c41a;
      }
    }

    .package-status {
      display: flex;
      align-items: center;
      gap: 6px;
      justify-content: center;

      .status-icon {
        font-size: 8px;

        &.status-normal {
          color: #52c41a;
        }

        &.status-paused {
          color: #faad14;
        }

        &.status-error {
          color: #ff4d4f;
        }
      }
    }

    .action-buttons {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
    }

    .pagination-section {
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
  .dns-domain-list {
    .page-header .header-content,
    .search-section .search-controls,
    .table-section .el-table,
    .table-section .pagination-section {
      max-width: 100%;
      padding: 0 16px;
    }
  }
}

@media (max-width: 768px) {
  .dns-domain-list {
    .page-header {
      padding: 16px;

      .notice-banners {
        flex-direction: column;

        .notice-banner {
          min-width: auto;
        }
      }
    }

    .search-section {
      padding: 12px 16px;

      .search-controls {
        flex-direction: column;
        gap: 12px;
        align-items: stretch;

        .search-filters {
          justify-content: space-between;
        }
      }
    }

    .table-section {
      padding: 0 16px 16px;

      .pagination-section {
        flex-direction: column;
        gap: 16px;
        align-items: stretch;

        .batch-actions {
          justify-content: center;
        }
      }
    }
  }
}
</style>
