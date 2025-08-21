<template>
  <el-dialog
    v-model="visibleLocal"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    @close="handleClose"
    class="bulk-records-modal"
    width="1000px"
    align-center
  >
    <template #header>
      <div class="modal-header">
        <div class="header-icon">
          <el-icon size="28"><Grid /></el-icon>
        </div>
        <div class="header-content">
          <h3 class="modal-title">批量添加DNS记录</h3>
          <p class="modal-subtitle">快速批量创建多个DNS解析记录</p>
        </div>
      </div>
    </template>
    <div class="bulk-body">
      <div class="info-section">
        <el-alert 
          type="info" 
          :closable="false" 
          show-icon
          class="info-alert"
        >
          <template #title>
            <div class="alert-content">
              <span class="alert-text">支持从Excel/CSV粘贴多行数据，或手动逐行添加</span>
              <div class="supported-types">
                <span class="types-label">支持类型：</span>
                <el-tag size="small" v-for="type in ['A', 'AAAA', 'CNAME', 'MX', 'TXT', 'SRV']" :key="type" class="type-tag">
                  {{ type }}
                </el-tag>
              </div>
            </div>
          </template>
        </el-alert>
      </div>
      
      <div class="toolbar">
        <div class="toolbar-left">
          <el-button @click="addRow" :icon="Plus" type="primary" plain>
            添加一行
          </el-button>
          <el-button @click="pasteRows" :icon="DocumentCopy" type="success" plain>
            从剪贴板粘贴
          </el-button>
        </div>
        <div class="toolbar-right">
          <el-checkbox v-model="syncAfterCreate" class="sync-checkbox">
            <span class="checkbox-text">创建后立即同步到云厂商</span>
          </el-checkbox>
        </div>
      </div>
      
      <div class="table-container">
        <el-table :data="rows" class="records-table" size="default" stripe>
          <el-table-column label="序号" width="60" type="index" align="center" />
          <el-table-column prop="name" label="主机记录" width="140">
            <template #default="{ row }">
              <el-input v-model="row.name" placeholder="@ 或 www" size="small" class="table-input" />
            </template>
          </el-table-column>
          <el-table-column prop="type" label="类型" width="100">
            <template #default="{ row }">
              <el-select v-model="row.type" size="small" class="table-select">
                <el-option label="A" value="A" />
                <el-option label="AAAA" value="AAAA" />
                <el-option label="CNAME" value="CNAME" />
                <el-option label="MX" value="MX" />
                <el-option label="TXT" value="TXT" />
                <el-option label="SRV" value="SRV" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column prop="value" label="记录值" min-width="200">
            <template #default="{ row }">
              <el-input v-model="row.value" placeholder="1.1.1.1 或 cname.example.com" size="small" class="table-input" />
            </template>
          </el-table-column>
          <el-table-column prop="ttl" label="TTL" width="100">
            <template #default="{ row }">
              <el-input-number v-model="row.ttl" :min="10" :max="86400" :step="10" size="small" class="table-number" />
            </template>
          </el-table-column>
          <el-table-column prop="priority" label="优先级" width="90">
            <template #default="{ row }">
              <el-input-number v-model="row.priority" :min="0" :max="65535" size="small" class="table-number" />
            </template>
          </el-table-column>
          <el-table-column prop="weight" label="权重" width="80">
            <template #default="{ row }">
              <el-input-number v-model="row.weight" :min="0" :max="65535" size="small" class="table-number" />
            </template>
          </el-table-column>
          <el-table-column prop="port" label="端口" width="80">
            <template #default="{ row }">
              <el-input-number v-model="row.port" :min="1" :max="65535" size="small" class="table-number" />
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" min-width="120">
            <template #default="{ row }">
              <el-input v-model="row.remark" placeholder="可选" size="small" class="table-input" />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="80" align="center">
            <template #default="{ $index }">
              <el-button size="small" type="danger" link @click="removeRow($index)" :icon="Delete">
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      
      <div class="status-section" v-if="errorMessage || result">
        <el-alert v-if="errorMessage" type="error" :title="errorMessage" show-icon class="status-alert" />
        <el-result 
          v-if="result" 
          :icon="result.failedCount ? 'warning' : 'success'" 
          :title="resultTitle" 
          :sub-title="`共处理 ${result.total} 条记录`" 
          class="result-display"
        >
          <template #extra>
            <div class="result-stats">
              <div class="stat-item success">
                <span class="stat-label">成功</span>
                <span class="stat-value">{{ result.successCount }}</span>
              </div>
              <div class="stat-item failed" v-if="result.failedCount > 0">
                <span class="stat-label">失败</span>
                <span class="stat-value">{{ result.failedCount }}</span>
              </div>
            </div>
          </template>
        </el-result>
      </div>
    </div>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose" size="large">
          取消
        </el-button>
        <el-button 
          type="primary" 
          :loading="loading" 
          :disabled="rows.length === 0" 
          @click="handleSubmit"
          size="large"
          :icon="loading ? Loading : Check"
        >
          {{ loading ? '提交中...' : `提交 (${rows.length} 条)` }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, DocumentCopy, Delete, Grid, Check, Loading } from '@element-plus/icons-vue'
import { recordApi } from '@/api/dns/record'

const props = defineProps<{ visible: boolean; domainId: number | null; providerId?: number | null }>()
const emit = defineEmits(['update:visible','success'])

const visibleLocal = ref(props.visible)
watch(() => props.visible, v => {
  visibleLocal.value = v
  if (v) {
    resetModal()
  }
})
watch(visibleLocal, v => emit('update:visible', v))

const rows = ref<any[]>([])
const loading = ref(false)
const errorMessage = ref('')
const result = ref<{successCount:number; failedCount:number; total:number}|null>(null)
const syncAfterCreate = ref(false)

const addRow = () => rows.value.push({ name:'', type:'A', value:'', ttl:600, priority:null, weight:null, port:null, remark:'' })
const removeRow = (i:number) => rows.value.splice(i,1)

const pasteRows = async () => {
  try {
    const text = await navigator.clipboard.readText()
    if (!text) return
    const lines = text.split(/\r?\n/).filter(Boolean)
    for (const line of lines) {
      const cols = line.split(/\t|,/) // 支持tab或逗号
      rows.value.push({
        name: cols[0] || '',
        type: cols[1] || 'A',
        value: cols[2] || '',
        ttl: Number(cols[3] || 600),
        priority: toNum(cols[4]),
        weight: toNum(cols[5]),
        port: toNum(cols[6]),
        remark: cols[7] || ''
      })
    }
  } catch {}
}

const toNum = (v:any) => {
  const n = Number(v); return Number.isFinite(n) ? n : null
}

const resultTitle = computed(() => {
  if (!result.value) return ''
  return result.value.failedCount ? `部分成功：成功 ${result.value.successCount}，失败 ${result.value.failedCount}` : `全部成功：成功 ${result.value.successCount}`
})

const handleSubmit = async () => {
  if (!props.domainId) { ElMessage.warning('请先选择域名'); return }
  if (rows.value.length === 0) return
  loading.value = true; errorMessage.value=''; result.value=null
  try {
    // 优先使用推荐的 batchCreate 接口
    const resp = await recordApi.batchCreate({
      domain_id: props.domainId,
      records: rows.value.map(r => ({ ...r, domain_id: props.domainId })),
      sync: syncAfterCreate.value,
      provider_id: props.providerId || undefined
    } as any)
    result.value = {
      successCount: (resp as any).data?.success || (resp as any).success,
      failedCount: (resp as any).data?.failed || (resp as any).failed,
      total: (resp as any).data?.total || (resp as any).total
    }
    ElMessage.success('提交完成')
    emit('success')
  } catch (e:any) {
    errorMessage.value = e?.message || '提交失败'
  } finally {
    loading.value = false
  }
}

const handleClose = () => {
  visibleLocal.value = false
}

const resetModal = () => {
  rows.value = []
  loading.value = false
  errorMessage.value = ''
  result.value = null
  syncAfterCreate.value = false
}
</script>

<style scoped lang="scss">
.bulk-records-modal {
  :deep(.el-dialog) {
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  :deep(.el-dialog__body) {
    padding: 24px;
  }

  :deep(.el-dialog__footer) {
    padding: 16px 24px;
    border-top: 1px solid #e1e4e8;
    background: #f6f8fa;
  }

  .modal-header {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 20px 24px;
    background: #0366d6;
    color: white;

    .header-icon {
      width: 40px;
      height: 40px;
      border-radius: 6px;
      background: rgba(255, 255, 255, 0.2);
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
    }

    .header-content {
      .modal-title {
        margin: 0 0 4px 0;
        font-size: 18px;
        font-weight: 600;
        color: white;
      }

      .modal-subtitle {
        margin: 0;
        font-size: 14px;
        color: rgba(255, 255, 255, 0.8);
      }
    }
  }

  .bulk-body {
    display: flex;
    flex-direction: column;
    gap: 16px;
    
    .info-section {
      .info-alert {
        :deep(.el-alert__content) {
          padding: 0;
        }
        
        .alert-content {
          .alert-text {
            font-size: 14px;
            margin-bottom: 8px;
            display: block;
          }
          
          .supported-types {
            display: flex;
            align-items: center;
            gap: 8px;
            
            .types-label {
              font-size: 12px;
              color: #586069;
            }
            
            .type-tag {
              background: #0366d6;
              color: white;
              border: none;
              font-size: 11px;
              padding: 2px 6px;
            }
          }
        }
      }
    }
    
    .toolbar {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px;
      background: #f6f8fa;
      border: 1px solid #e1e4e8;
      border-radius: 6px;
      
      .toolbar-left {
        display: flex;
        gap: 8px;
      }
      
      .toolbar-right {
        .sync-checkbox {
          :deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
            background: #0366d6;
            border-color: #0366d6;
          }
        }
      }
    }
    
    .table-container {
      border: 1px solid #e1e4e8;
      border-radius: 6px;
      overflow: hidden;
      
      .records-table {
        :deep(.el-table__header) {
          background: #f6f8fa;
          
          th {
            color: #24292e;
            font-weight: 600;
            border-bottom: 1px solid #e1e4e8;
          }
        }
        
        :deep(.el-table__row) {
          &:hover {
            background: #f6f8fa;
          }
        }
        
        .table-input,
        .table-select,
        .table-number {
          :deep(.el-input__wrapper),
          :deep(.el-select__wrapper) {
            border: 1px solid #d0d7de;
            border-radius: 4px;
            
            &:hover {
              border-color: #0366d6;
            }
            
            &.is-focus,
            &.is-focused {
              border-color: #0366d6;
              box-shadow: 0 0 0 3px rgba(3, 102, 214, 0.1);
            }
          }
        }
      }
    }
    
    .status-section {
      .status-alert {
        border-radius: 6px;
      }
      
      .result-display {
        border-radius: 6px;
        border: 1px solid #e1e4e8;
        
        .result-stats {
          display: flex;
          gap: 16px;
          justify-content: center;
          
          .stat-item {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 4px;
            padding: 12px 16px;
            border-radius: 6px;
            border: 1px solid #e1e4e8;
            background: white;
            
            &.success {
              border-color: #28a745;
              
              .stat-label,
              .stat-value {
                color: #28a745;
              }
            }
            
            &.failed {
              border-color: #d73a49;
              
              .stat-label,
              .stat-value {
                color: #d73a49;
              }
            }
          }
        }
      }
    }
  }
  
  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 8px;

    :deep(.el-button) {
      border-radius: 6px;
      font-weight: 500;
      padding: 6px 16px;

      &.el-button--primary {
        background: #0366d6;
        border-color: #0366d6;

        &:hover {
          background: #0256cc;
          border-color: #0256cc;
        }
      }
    }
  }
}

