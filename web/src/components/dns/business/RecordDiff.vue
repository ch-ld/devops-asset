<template>
  <div class="record-diff">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>DNS记录差异对比</span>
          <div>
            <el-button 
              type="warning" 
              size="small"
              @click="$emit('sync-all')"
              :loading="syncLoading"
              v-if="hasDifferences"
            >
              同步所有差异
            </el-button>
            <el-button 
              type="default" 
              size="small"
              @click="$emit('refresh')"
              :loading="refreshLoading"
            >
              刷新对比
            </el-button>
          </div>
        </div>
      </template>

      <div v-if="!localRecords.length && !remoteRecords.length" class="empty-state">
        <el-empty description="暂无DNS记录" />
      </div>

      <div v-else>
        <!-- 统计信息 -->
        <el-row :gutter="16" style="margin-bottom: 16px;">
          <el-col :span="6">
            <el-statistic title="本地记录" :value="localRecords.length" />
          </el-col>
          <el-col :span="6">
            <el-statistic title="远程记录" :value="remoteRecords.length" />
          </el-col>
          <el-col :span="6">
            <el-statistic title="需要同步" :value="pendingRecords.length" />
          </el-col>
          <el-col :span="6">
            <el-statistic title="冲突记录" :value="conflictRecords.length" />
          </el-col>
        </el-row>

        <!-- 差异详情 -->
        <el-table 
          :data="diffRecords" 
          v-loading="loading"
          row-key="key"
          :expand-row-keys="expandedRows"
          @expand-change="handleExpandChange"
        >
          <el-table-column type="expand">
            <template #default="{ row }">
              <div class="record-detail">
                <el-row :gutter="16">
                  <el-col :span="12">
                    <h4>本地记录</h4>
                    <div v-if="row.local" class="record-info">
                      <p><strong>记录值:</strong> {{ row.local.value }}</p>
                      <p><strong>TTL:</strong> {{ row.local.ttl }}</p>
                      <p v-if="row.local.priority"><strong>优先级:</strong> {{ row.local.priority }}</p>
                      <p><strong>最后同步:</strong> {{ formatDate(row.local.last_sync_at) }}</p>
                    </div>
                    <div v-else class="no-record">
                      <el-tag type="warning">本地无此记录</el-tag>
                    </div>
                  </el-col>
                  <el-col :span="12">
                    <h4>远程记录</h4>
                    <div v-if="row.remote" class="record-info">
                      <p><strong>记录值:</strong> {{ row.remote.value }}</p>
                      <p><strong>TTL:</strong> {{ row.remote.ttl }}</p>
                      <p v-if="row.remote.priority"><strong>优先级:</strong> {{ row.remote.priority }}</p>
                      <p><strong>记录ID:</strong> {{ row.remote.remote_id }}</p>
                    </div>
                    <div v-else class="no-record">
                      <el-tag type="warning">远程无此记录</el-tag>
                    </div>
                  </el-col>
                </el-row>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="name" label="记录名" width="200" />
          <el-table-column prop="type" label="类型" width="80" />
          
          <el-table-column label="差异状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getDiffStatusType(row.status)" size="small">
                {{ getDiffStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column label="值差异" show-overflow-tooltip>
            <template #default="{ row }">
              <div v-if="row.status === 'conflict'" class="value-diff">
                <div class="local-value">
                  <el-tag type="info" size="small">本地:</el-tag>
                  {{ row.local?.value || '-' }}
                </div>
                <div class="remote-value" style="margin-top: 4px;">
                  <el-tag type="warning" size="small">远程:</el-tag>
                  {{ row.remote?.value || '-' }}
                </div>
              </div>
              <div v-else>
                {{ (row.local || row.remote)?.value || '-' }}
              </div>
            </template>
          </el-table-column>

          <el-table-column label="操作" width="200">
            <template #default="{ row }">
              <el-button-group size="small">
                <el-button 
                  type="primary"
                  @click="handleSync(row, 'push')"
                  v-if="row.status === 'local-only' || row.status === 'conflict'"
                  :loading="row.syncing"
                >
                  推送到远程
                </el-button>
                <el-button 
                  type="success"
                  @click="handleSync(row, 'pull')"
                  v-if="row.status === 'remote-only' || row.status === 'conflict'"
                  :loading="row.syncing"
                >
                  拉取到本地
                </el-button>
                <el-button 
                  type="danger"
                  @click="handleDelete(row)"
                  v-if="row.status === 'local-only'"
                  :loading="row.deleting"
                >
                  删除本地
                </el-button>
              </el-button-group>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { DnsRecord } from '@/types/dns'

interface DiffRecord {
  key: string
  name: string
  type: string
  status: 'synced' | 'local-only' | 'remote-only' | 'conflict'
  local?: DnsRecord
  remote?: any
  syncing?: boolean
  deleting?: boolean
}

interface Props {
  localRecords: DnsRecord[]
  remoteRecords: any[]
  loading?: boolean
  syncLoading?: boolean
  refreshLoading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  localRecords: () => [],
  remoteRecords: () => [],
  loading: false,
  syncLoading: false,
  refreshLoading: false
})

const emit = defineEmits<{
  'sync-record': [record: DiffRecord, action: 'push' | 'pull']
  'delete-record': [record: DiffRecord]
  'sync-all': []
  refresh: []
}>()

// 响应式数据
const expandedRows = ref<string[]>([])

// 计算属性
const diffRecords = computed(() => {
  const records = new Map<string, DiffRecord>()
  
  // 处理本地记录
  props.localRecords.forEach(local => {
    const key = `${local.name}.${local.type}`
    records.set(key, {
      key,
      name: local.name,
      type: local.type,
      status: 'local-only',
      local
    })
  })
  
  // 处理远程记录
  props.remoteRecords.forEach(remote => {
    const key = `${remote.name}.${remote.type}`
    const existing = records.get(key)
    
    if (existing) {
      // 检查是否有差异
      const hasValueDiff = existing.local!.value !== remote.value
      const hasTtlDiff = existing.local!.ttl !== remote.ttl
      const hasPriorityDiff = existing.local!.priority !== remote.priority
      
      if (hasValueDiff || hasTtlDiff || hasPriorityDiff) {
        existing.status = 'conflict'
      } else {
        existing.status = 'synced'
      }
      existing.remote = remote
    } else {
      records.set(key, {
        key,
        name: remote.name,
        type: remote.type,
        status: 'remote-only',
        remote
      })
    }
  })
  
  return Array.from(records.values()).sort((a, b) => {
    // 先按状态排序（冲突优先），再按名称排序
    const statusOrder = { conflict: 0, 'local-only': 1, 'remote-only': 2, synced: 3 }
    const statusDiff = statusOrder[a.status] - statusOrder[b.status]
    if (statusDiff !== 0) return statusDiff
    return a.name.localeCompare(b.name)
  })
})

const pendingRecords = computed(() => 
  diffRecords.value.filter(r => r.status !== 'synced')
)

const conflictRecords = computed(() => 
  diffRecords.value.filter(r => r.status === 'conflict')
)

const hasDifferences = computed(() => pendingRecords.value.length > 0)

// 工具方法
const getDiffStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    synced: 'success',
    'local-only': 'info',
    'remote-only': 'warning',
    conflict: 'danger'
  }
  return statusMap[status] || 'info'
}

const getDiffStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    synced: '已同步',
    'local-only': '仅本地',
    'remote-only': '仅远程',
    conflict: '冲突'
  }
  return statusMap[status] || status
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 事件处理
const handleExpandChange = (row: DiffRecord, expanded: boolean) => {
  if (expanded) {
    expandedRows.value.push(row.key)
  } else {
    const index = expandedRows.value.indexOf(row.key)
    if (index > -1) {
      expandedRows.value.splice(index, 1)
    }
  }
}

const handleSync = async (record: DiffRecord, action: 'push' | 'pull') => {
  try {
    record.syncing = true
    emit('sync-record', record, action)
  } catch (error) {
    console.error('同步失败:', error)
    ElMessage.error('同步失败')
  } finally {
    record.syncing = false
  }
}

const handleDelete = async (record: DiffRecord) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除本地记录 "${record.name}.${record.type}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    record.deleting = true
    emit('delete-record', record)
  } catch (error) {
    // 用户取消删除
  } finally {
    record.deleting = false
  }
}

// 监听数据变化，自动展开冲突记录
watch(() => conflictRecords.value, (newConflicts) => {
  if (newConflicts.length > 0 && expandedRows.value.length === 0) {
    // 自动展开第一个冲突记录
    expandedRows.value = [newConflicts[0].key]
  }
}, { immediate: true })
</script>

<style scoped lang="scss">
.record-diff {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .empty-state {
    text-align: center;
    padding: 40px 0;
  }

  .record-detail {
    padding: 16px;
    background-color: #f8f9fa;
    border-radius: 4px;
    margin: 8px 0;

    h4 {
      margin: 0 0 12px 0;
      color: #409eff;
    }

    .record-info {
      p {
        margin: 4px 0;
        font-size: 14px;
      }
    }

    .no-record {
      padding: 20px;
      text-align: center;
    }
  }

  .value-diff {
    .local-value,
    .remote-value {
      font-size: 12px;
      word-break: break-all;
    }
  }

  :deep(.el-table__expand-icon) {
    font-size: 16px;
  }
}
</style>
