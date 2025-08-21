<template>
  <el-dialog
    v-model="visibleLocal"
    width="1000px"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    @close="handleClose"
    class="import-records-modal"
    align-center
  >
    <template #header>
      <div class="modal-header">
        <div class="header-icon">
          <el-icon size="28"><Upload /></el-icon>
        </div>
        <div class="header-content">
          <h3 class="modal-title">导入DNS解析记录</h3>
          <p class="modal-subtitle">批量导入CSV或Excel格式的DNS记录</p>
        </div>
      </div>
    </template>
    <div class="import-body">
      <!-- 步骤指示器 -->
      <div class="steps-container">
        <el-steps :active="currentStep" finish-status="success" class="import-steps" align-center>
          <el-step>
            <template #icon>
              <el-icon><UploadFilled /></el-icon>
            </template>
            <template #title>
              <span class="step-title">上传文件</span>
            </template>
            <template #description>
              <span class="step-desc">选择要导入的文件</span>
            </template>
          </el-step>
          <el-step>
            <template #icon>
              <el-icon><View /></el-icon>
            </template>
            <template #title>
              <span class="step-title">数据预览</span>
            </template>
            <template #description>
              <span class="step-desc">预览并确认导入数据</span>
            </template>
          </el-step>
          <el-step>
            <template #icon>
              <el-icon><CircleCheck /></el-icon>
            </template>
            <template #title>
              <span class="step-title">导入完成</span>
            </template>
            <template #description>
              <span class="step-desc">查看导入结果</span>
            </template>
          </el-step>
        </el-steps>
      </div>

      <!-- 步骤1: 上传文件 -->
      <div v-if="currentStep === 0" class="step-content upload-step">
        <div class="info-section">
          <el-alert
            type="info"
            :closable="false"
            show-icon
            class="info-alert"
          >
            <template #title>
              <div class="alert-content">
                <span class="alert-text">支持 CSV 与 Excel(xlsx) 格式</span>
                <span class="alert-subtext">建议先下载模板，按模板填写后上传</span>
              </div>
            </template>
          </el-alert>
        </div>

        <div class="template-section">
          <div class="section-title">
            <el-icon><Download /></el-icon>
            <span>下载模板</span>
          </div>
          <div class="template-buttons">
            <el-button 
              @click="downloadTemplate('csv')" 
              :loading="downloading.csv"
              :icon="DocumentCopy"
              type="primary"
              plain
            >
              下载CSV模板
            </el-button>
            <el-button 
              @click="downloadTemplate('excel')" 
              :loading="downloading.excel"
              :icon="Document"
              type="success"
              plain
            >
              下载Excel模板
            </el-button>
          </div>
        </div>

        <div class="upload-section">
          <div class="section-title">
            <el-icon><UploadFilled /></el-icon>
            <span>上传文件</span>
          </div>
          <el-upload
            drag
            :auto-upload="false"
            :on-change="onFileChange"
            :limit="1"
            accept=".csv,.xlsx"
            class="upload-area"
          >
            <div class="upload-content">
              <div class="upload-icon">
                <el-icon size="48"><UploadFilled /></el-icon>
              </div>
              <div class="upload-text">
                <h4>将文件拖到此处</h4>
                <p>或 <em>点击上传</em></p>
              </div>
              <div class="upload-tip">
                <el-icon><InfoFilled /></el-icon>
                <span>仅支持 CSV 或 Excel(xlsx) 格式</span>
              </div>
            </div>
          </el-upload>
        </div>
      </div>

      <!-- 步骤2: 数据预览 -->
      <div v-if="currentStep === 1" class="step-content">
        <div class="step-header">
          <h3>
            <el-icon><View /></el-icon>
            数据预览
          </h3>
          <p>请确认以下数据无误后，点击导入</p>
        </div>

        <div class="preview-summary">
          <el-tag size="large" type="primary">总计：{{ previewData.length }} 条记录</el-tag>
          <el-tag size="large" type="success">有效：{{ validCount }} 条</el-tag>
          <el-tag size="large" type="danger" v-if="invalidCount > 0">无效：{{ invalidCount }} 条</el-tag>
        </div>

        <el-table :data="previewData" border max-height="400" class="preview-table">
          <el-table-column type="index" label="#" width="50" />
          <el-table-column prop="name" label="主机记录" width="120">
            <template #default="{ row }">
              <div class="cell-content">
                <span>{{ row.name || '@' }}</span>
                <el-icon v-if="row._error" class="error-icon"><WarningFilled /></el-icon>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="type" label="类型" width="80">
            <template #default="{ row }">
              <el-tag :type="getRecordTypeColor(row.type)" size="small">{{ row.type }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="value" label="记录值" min-width="150" show-overflow-tooltip />
          <el-table-column prop="ttl" label="TTL" width="80">
            <template #default="{ row }">
              {{ row.ttl || 600 }}
            </template>
          </el-table-column>
          <el-table-column prop="priority" label="优先级" width="80">
            <template #default="{ row }">
              {{ row.priority || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="weight" label="权重" width="80">
            <template #default="{ row }">
              {{ row.weight || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" min-width="100" show-overflow-tooltip />
          <el-table-column v-if="invalidCount > 0" prop="_error" label="错误信息" width="150">
            <template #default="{ row }">
              <el-text v-if="row._error" type="danger" size="small">{{ row._error }}</el-text>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 步骤3: 导入结果 -->
      <div v-if="currentStep === 2" class="step-content">
        <el-result :icon="resultIcon" :title="resultTitle" :sub-title="resultSub">
          <template #extra>
            <el-button type="primary" @click="emit('success')" v-if="result?.successCount">刷新记录</el-button>
          </template>
        </el-result>
        <el-table v-if="result?.failedItems?.length" :data="result.failedItems" size="small" style="margin-top: 8px">
          <el-table-column prop="row" label="#" width="60" />
          <el-table-column prop="reason" label="失败原因" />
        </el-table>
      </div>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button v-if="currentStep > 0" @click="prevStep">上一步</el-button>
        <el-button v-if="currentStep === 0" type="primary" :disabled="!file" @click="nextStep">下一步</el-button>
        <el-button v-if="currentStep === 1" type="primary" :loading="loading" :disabled="invalidCount > 0" @click="handleImport">开始导入</el-button>
        <el-button v-if="currentStep === 2" type="primary" @click="handleClose">完成</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { recordApi } from '@/api/dns/record'
import api from '@/api/client'
import { UploadFilled, View, WarningFilled, Download, DocumentCopy, Document, InfoFilled, CircleCheck, Upload } from '@element-plus/icons-vue'

const props = defineProps<{ visible: boolean; domainId: number | null }>()
const emit = defineEmits(['update:visible', 'success'])

const visibleLocal = ref(props.visible)
const file = ref<File | null>(null)
const loading = ref(false)
const downloading = ref({ csv: false, excel: false })
const result = ref<{ successCount: number; failedCount: number; total: number; failedItems: any[] } | null>(null)
const currentStep = ref(0)
const previewData = ref<any[]>([])

watch(() => props.visible, v => {
  visibleLocal.value = v
  if (v) {
    resetModal()
  }
})
watch(visibleLocal, v => emit('update:visible', v))

const onFileChange = (f: any) => {
  file.value = f?.raw || null
  result.value = null
}

// 计算属性
const validCount = computed(() => previewData.value.filter(item => !item._error).length)
const invalidCount = computed(() => previewData.value.filter(item => item._error).length)

const resultIcon = computed(() => (result.value && result.value.failedCount === 0 ? 'success' : 'warning'))
const resultTitle = computed(() => {
  if (!result.value) return '等待导入'
  const { successCount, failedCount } = result.value
  return failedCount === 0 ? `全部成功：成功 ${successCount} 条` : `部分成功：成功 ${successCount} 条，失败 ${failedCount} 条`
})
const resultSub = computed(() => (result.value ? `共 ${result.value.total} 条` : ''))

// 步骤控制
const nextStep = async () => {
  if (currentStep.value === 0 && file.value) {
    try {
      loading.value = true
      const rows = await parseFile(file.value)
      previewData.value = rows.map((row, index) => {
        const item = { ...row, _error: '' }

        // 验证必填字段
        if (!item.name && item.name !== '@') {
          item._error = '主机记录不能为空'
        } else if (!item.type) {
          item._error = '记录类型不能为空'
        } else if (!item.value) {
          item._error = '记录值不能为空'
        } else if (!['A', 'AAAA', 'CNAME', 'MX', 'TXT', 'NS', 'SRV', 'PTR'].includes(item.type.toUpperCase())) {
          item._error = '不支持的记录类型'
        }

        return item
      })

      currentStep.value = 1
      ElMessage.success(`解析成功，共 ${previewData.value.length} 条记录`)
    } catch (error) {
      ElMessage.error('文件解析失败，请检查文件格式')
    } finally {
      loading.value = false
    }
  }
}

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

const handleImport = async () => {
  if (!file.value || invalidCount.value > 0) return
  if (!props.domainId) {
    ElMessage.warning('请先选择域名')
    return
  }
  loading.value = true
  try {
    const validData = previewData.value.filter(item => !item._error)
    const payload = {
      format: 'csv',
      data: validData.map(r => ({
        domain_id: props.domainId,
        name: r.name,
        type: r.type,
        value: r.value,
        ttl: toInt(r.ttl, 600),
        priority: toInt(r.priority),
        weight: toInt(r.weight),
        port: toInt(r.port),
        remark: r.remark || ''
      }))
    }
    const resp = await recordApi.import(payload as any)
    result.value = {
      successCount: (resp as any).data?.success || (resp as any).success,
      failedCount: (resp as any).data?.failed || (resp as any).failed,
      total: (resp as any).data?.total || (resp as any).total,
      failedItems: (resp as any).data?.failed_items || (resp as any).failed_items || []
    }
    currentStep.value = 2
    ElMessage.success('导入完成')
  } catch (e: any) {
    console.error(e)
    ElMessage.error(e?.message || '导入失败')
  } finally {
    loading.value = false
  }
}

const toInt = (v: any, def?: number) => {
  const n = Number(v)
  return Number.isFinite(n) ? n : def
}

const parseFile = async (f: File): Promise<any[]> => {
  const text = await f.text()
  // 简易 CSV 解析（以逗号分隔，首行为表头）
  const lines = text.split(/\r?\n/).filter(Boolean)
  if (lines.length === 0) return []
  const headers = lines[0].split(',').map(s => s.trim().toLowerCase())
  const idx = (h: string) => headers.indexOf(h)
  const out: any[] = []
  for (let i = 1; i < lines.length; i++) {
    const cols = splitCSVLine(lines[i])
    out.push({
      name: cols[idx('name')] || '',
      type: cols[idx('type')] || 'A',
      value: cols[idx('value')] || '',
      ttl: cols[idx('ttl')] || '',
      priority: cols[idx('priority')] || '',
      weight: cols[idx('weight')] || '',
      port: cols[idx('port')] || '',
      remark: cols[idx('remark')] || ''
    })
  }
  return out
}

const splitCSVLine = (line: string): string[] => {
  const res: string[] = []
  let cur = ''
  let inQuotes = false
  for (let i = 0; i < line.length; i++) {
    const ch = line[i]
    if (ch === '"') {
      if (inQuotes && line[i + 1] === '"') { cur += '"'; i++; continue }
      inQuotes = !inQuotes
      continue
    }
    if (ch === ',' && !inQuotes) { res.push(cur); cur = ''; continue }
    cur += ch
  }
  res.push(cur)
  return res.map(s => s.trim())
}

const downloadTemplate = async (type: 'csv' | 'excel') => {
  downloading.value[type] = true
  try {
    // 优先走后端模板接口
    const url = type === 'csv' ? '/api/v1/dns/records/template/csv' : '/api/v1/dns/records/template/excel'
    const blob = await api.get({ url, responseType: 'blob' }) as unknown as Blob
    const link = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = link
    a.download = type === 'csv' ? 'dns-records-template.csv' : 'dns-records-template.csv'
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(link)
  } catch (e) {
    // 兜底：前端生成CSV
    const content = 'name,type,value,ttl,priority,weight,port,remark\n'
    const blob = new Blob([content], { type: 'text/csv;charset=utf-8;' })
    const link = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = link
    a.download = 'dns-records-template.csv'
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(link)
  } finally {
    downloading.value[type] = false
  }
}

const handleClose = () => {
  visibleLocal.value = false
}

const resetModal = () => {
  currentStep.value = 0
  file.value = null
  previewData.value = []
  result.value = null
}

// 工具函数
const getRecordTypeColor = (type: string) => {
  const colorMap: Record<string, string> = {
    A: 'success',
    AAAA: 'success',
    CNAME: 'warning',
    MX: 'danger',
    TXT: 'info',
    NS: 'primary',
    SRV: 'info'
  }
  return colorMap[type] || 'info'
}
</script>

<style scoped lang="scss">
.import-records-modal {
  :deep(.el-dialog) {
    border-radius: 8px;
    overflow: hidden;
    background: white;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    border: 1px solid #e1e4e8;
  }

  :deep(.el-dialog__header) {
    padding: 0;
    border-bottom: 1px solid rgba(226, 232, 240, 0.3);
  }

  :deep(.el-dialog__body) {
    padding: 0;
    background: transparent;
  }

  :deep(.el-dialog__footer) {
    padding: 16px 24px;
    border-top: 1px solid #e1e4e8;
    background: #f6f8fa;
  }

  .modal-header {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 24px;
    background: #0366d6;
    color: white;
    border-bottom: 1px solid #e1e4e8;

    .header-icon {
      width: 40px;
      height: 40px;
      border-radius: 6px;
      background: rgba(255, 255, 255, 0.15);
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 20px;
    }

    .header-content {
      flex: 1;

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
        font-weight: 400;
      }
    }
  }

  .import-body {
    padding: 32px;
    
    .steps-container {
      margin-bottom: 32px;
      padding: 24px;
      background: rgba(248, 250, 252, 0.6);
      border-radius: 16px;
      border: 1px solid rgba(226, 232, 240, 0.5);
      
      .import-steps {
        :deep(.el-step__title) {
          .step-title {
            font-weight: 700;
            color: #1e293b;
            font-size: 16px;
          }
        }
        
        :deep(.el-step__description) {
          .step-desc {
            color: #64748b;
            font-weight: 500;
          }
        }
        
        :deep(.el-step__icon) {
          width: 48px;
          height: 48px;
          border-radius: 50%;
          border: 2px solid #e2e8f0;
          background: rgba(255, 255, 255, 0.9);
          
          &.is-process {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-color: #667eea;
            color: white;
          }
          
          &.is-finish {
            background: linear-gradient(135deg, #10b981 0%, #059669 100%);
            border-color: #10b981;
            color: white;
          }
        }
        
        :deep(.el-step__line) {
          background: linear-gradient(135deg, #e2e8f0 0%, #cbd5e1 100%);
          
          &.is-finish {
            background: linear-gradient(135deg, #10b981 0%, #059669 100%);
          }
        }
      }
    }
    
    .step-content {
      min-height: 400px;
      
      &.upload-step {
        display: flex;
        flex-direction: column;
        gap: 24px;
        
        .info-section {
          .info-alert {
            border-radius: 12px;
            border: 1px solid rgba(59, 130, 246, 0.2);
            background: rgba(239, 246, 255, 0.8);
            backdrop-filter: blur(10px);
            
            .alert-content {
              display: flex;
              flex-direction: column;
              gap: 4px;
              
              .alert-text {
                font-size: 16px;
                font-weight: 600;
                color: #1e40af;
              }
              
              .alert-subtext {
                font-size: 14px;
                color: #3730a3;
                opacity: 0.8;
              }
            }
          }
        }
        
        .template-section,
        .upload-section {
          .section-title {
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 18px;
            font-weight: 700;
            color: #1e293b;
            margin-bottom: 16px;
            
            .el-icon {
              color: #667eea;
            }
          }
        }
        
        .template-buttons {
          display: flex;
          gap: 16px;
          
          :deep(.el-button) {
            border-radius: 12px;
            font-weight: 600;
            padding: 12px 24px;
            transition: all 0.3s ease;
            
            &:hover {
              transform: translateY(-2px);
            }
            
            &.el-button--primary.is-plain {
              border-color: #667eea;
              color: #667eea;
              
              &:hover {
                background: #667eea;
                color: white;
              }
            }
            
            &.el-button--success.is-plain {
              border-color: #10b981;
              color: #10b981;
              
              &:hover {
                background: #10b981;
                color: white;
              }
            }
          }
        }
        
        .upload-area {
          width: 100%;
          
          :deep(.el-upload) {
            width: 100%;
          }
          
          :deep(.el-upload-dragger) {
            width: 100%;
            height: 240px;
            border: 2px dashed #e2e8f0;
            border-radius: 16px;
            background: rgba(248, 250, 252, 0.6);
            backdrop-filter: blur(10px);
            transition: all 0.3s ease;
            display: flex;
            align-items: center;
            justify-content: center;
            
            &:hover {
              border-color: #667eea;
              background: rgba(102, 126, 234, 0.05);
              transform: translateY(-2px);
              box-shadow: 0 8px 32px rgba(102, 126, 234, 0.15);
            }
          }
          
          .upload-content {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 16px;
            
            .upload-icon {
              color: #667eea;
              opacity: 0.7;
            }
            
            .upload-text {
              text-align: center;
              
              h4 {
                margin: 0 0 8px 0;
                font-size: 20px;
                font-weight: 600;
                color: #1e293b;
              }
              
              p {
                margin: 0;
                font-size: 16px;
                color: #64748b;
                
                em {
                  color: #667eea;
                  font-weight: 600;
                  text-decoration: underline;
                }
              }
            }
            
            .upload-tip {
              display: flex;
              align-items: center;
              gap: 6px;
              font-size: 14px;
              color: #94a3b8;
              font-weight: 500;
              
              .el-icon {
                color: #667eea;
              }
            }
          }
        }
      }
    }
    
    .step-header {
      margin-bottom: 24px;
      padding: 20px;
      background: rgba(248, 250, 252, 0.6);
      border-radius: 12px;
      border: 1px solid rgba(226, 232, 240, 0.5);
      
      h3 {
        display: flex;
        align-items: center;
        gap: 10px;
        margin: 0 0 8px 0;
        font-size: 20px;
        font-weight: 700;
        color: #1e293b;
      }
      
      p {
        margin: 0;
        color: #64748b;
        font-size: 15px;
        font-weight: 500;
      }
    }
    
    .preview-summary {
      margin-bottom: 20px;
      display: flex;
      gap: 16px;
      
      :deep(.el-tag) {
        padding: 8px 16px;
        border-radius: 10px;
        font-weight: 600;
        font-size: 14px;
        border: none;
        
        &.el-tag--primary {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          color: white;
        }
        
        &.el-tag--success {
          background: linear-gradient(135deg, #10b981 0%, #059669 100%);
          color: white;
        }
        
        &.el-tag--danger {
          background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
          color: white;
        }
      }
    }
    
    .preview-table {
      border-radius: 16px;
      overflow: hidden;
      box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
      border: 1px solid rgba(226, 232, 240, 0.3);
      
      :deep(.el-table__header) {
        background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
        
        th {
          background: transparent;
          color: #475569;
          font-weight: 700;
          border-bottom: 2px solid #e2e8f0;
        }
      }
      
      :deep(.el-table__row) {
        background: rgba(255, 255, 255, 0.8);
        backdrop-filter: blur(5px);
        
        &:hover {
          background: rgba(255, 255, 255, 0.95);
          transform: translateY(-1px);
          box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
        }
      }
      
      .cell-content {
        display: flex;
        align-items: center;
        gap: 8px;
        
        .error-icon {
          color: #ef4444;
          font-size: 16px;
        }
      }
    }
  }
  
  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 16px;

    :deep(.el-button) {
      border-radius: 12px;
      font-weight: 600;
      padding: 14px 28px;
      font-size: 15px;
      transition: all 0.3s ease;
      border: 2px solid transparent;

      &.el-button--default {
        background: rgba(255, 255, 255, 0.9);
        backdrop-filter: blur(10px);
        border-color: #e2e8f0;
        color: #64748b;
        
        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
          border-color: #94a3b8;
          color: #475569;
        }
      }

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
</style>

