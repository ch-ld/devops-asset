<template>
  <el-dialog 
    v-model="visible" 
    title="批量导入域名"
    width="700px"
    class="import-modal"
    :close-on-click-modal="false"
    destroy-on-close
  >
    <div class="import-content">
      <!-- 步骤指示器 -->
      <el-steps :active="currentStep" finish-status="success" class="import-steps">
        <el-step title="下载模板" description="下载导入模板文件" />
        <el-step title="上传文件" description="上传填写好的数据文件" />
        <el-step title="数据预览" description="预览并确认导入数据" />
        <el-step title="导入完成" description="查看导入结果" />
      </el-steps>

      <!-- 步骤1: 下载模板 -->
      <div v-if="currentStep === 0" class="step-content">
        <div class="step-header">
          <h3>
            <el-icon><Download /></el-icon>
            下载导入模板
          </h3>
          <p>请先下载模板文件，按照格式填写域名数据</p>
        </div>
        
        <div class="template-options">
          <div class="template-card" @click="downloadTemplate('xlsx')">
            <div class="template-icon">
              <el-icon size="48"><Document /></el-icon>
            </div>
            <div class="template-info">
              <h4>Excel 模板</h4>
              <p>支持 .xlsx 格式，推荐使用</p>
              <el-tag type="success" size="small">推荐</el-tag>
            </div>
          </div>
          
          <div class="template-card" @click="downloadTemplate('csv')">
            <div class="template-icon">
              <el-icon size="48"><Document /></el-icon>
            </div>
            <div class="template-info">
              <h4>CSV 模板</h4>
              <p>支持 .csv 格式，通用性好</p>
            </div>
          </div>
        </div>
        
        <!-- 模板说明 -->
        <el-card class="template-guide">
          <template #header>
            <h4>
              <el-icon><InfoFilled /></el-icon>
              填写说明
            </h4>
          </template>
          <div class="guide-content">
            <div class="guide-item">
              <strong>域名 (name)*</strong>
              <span>必填，格式如：example.com</span>
            </div>
            <div class="guide-item">
              <strong>注册商 (registrar_type)*</strong>
              <span>必填，支持：aliyun、tencent、route53、cloudflare、dnspod、godaddy</span>
            </div>
            <div class="guide-item">
              <strong>状态 (status)</strong>
              <span>可选，active（正常）、inactive（禁用），默认：active</span>
            </div>
            <div class="guide-item">
              <strong>分组ID (group_id)</strong>
              <span>可选，域名分组ID，不填表示无分组</span>
            </div>
            <div class="guide-item">
              <strong>到期时间 (expires_at)</strong>
              <span>可选，格式：2024-12-31 23:59:59</span>
            </div>
            <div class="guide-item">
              <strong>自动续费 (auto_renew)</strong>
              <span>可选，true（启用）、false（禁用），默认：false</span>
            </div>
            <div class="guide-item">
              <strong>标签 (tags)</strong>
              <span>可选，多个标签用逗号分隔，如：生产,重要</span>
            </div>
            <div class="guide-item">
              <strong>备注 (remark)</strong>
              <span>可选，域名备注信息</span>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 步骤2: 上传文件 -->
      <div v-if="currentStep === 1" class="step-content">
        <div class="step-header">
          <h3>
            <el-icon><Upload /></el-icon>
            上传数据文件
          </h3>
          <p>请上传填写完成的域名数据文件</p>
        </div>
        
        <el-upload
          ref="uploadRef"
          class="upload-dragger"
          drag
          :auto-upload="false"
          :limit="1"
          :accept="'.xlsx,.xls,.csv'"
          :on-change="handleFileChange"
          :on-exceed="handleExceed"
        >
          <div class="upload-content">
            <el-icon class="upload-icon"><UploadFilled /></el-icon>
            <div class="upload-text">
              <p>将文件拖到此处，或<em>点击上传</em></p>
              <p class="upload-tip">支持 .xlsx、.xls、.csv 格式文件</p>
            </div>
          </div>
        </el-upload>
        
        <div v-if="uploadFile" class="file-info">
          <div class="file-item">
            <el-icon><Document /></el-icon>
            <span class="file-name">{{ uploadFile.name }}</span>
            <span class="file-size">{{ formatFileSize(uploadFile.size) }}</span>
            <el-button type="danger" text @click="removeFile">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
      </div>

      <!-- 步骤3: 数据预览 -->
      <div v-if="currentStep === 2" class="step-content">
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
          <el-table-column prop="name" label="域名" width="180">
            <template #default="{ row }">
              <div class="cell-content">
                <span>{{ row.name }}</span>
                <el-icon v-if="row._error" class="error-icon"><WarningFilled /></el-icon>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="registrar_type" label="注册商" width="120" />
          <el-table-column prop="status" label="状态" width="80">
            <template #default="{ row }">
              <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
                {{ row.status === 'active' ? '正常' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="expires_at" label="到期时间" width="160" />
          <el-table-column prop="auto_renew" label="自动续费" width="80">
            <template #default="{ row }">
              <el-tag :type="row.auto_renew ? 'success' : 'info'" size="small">
                {{ row.auto_renew ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="tags" label="标签" />
          <el-table-column prop="_error" label="错误信息" width="200">
            <template #default="{ row }">
              <span v-if="row._error" class="error-text">{{ row._error }}</span>
              <el-tag v-else type="success" size="small">有效</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 步骤4: 导入结果 -->
      <div v-if="currentStep === 3" class="step-content">
        <div class="step-header">
          <h3>
            <el-icon><CircleCheck /></el-icon>
            导入完成
          </h3>
        </div>
        
        <div class="result-summary">
          <div class="result-card success">
            <div class="result-icon">
              <el-icon><CircleCheck /></el-icon>
            </div>
            <div class="result-info">
              <div class="result-number">{{ importResult.success }}</div>
              <div class="result-label">成功导入</div>
            </div>
          </div>
          
          <div class="result-card total">
            <div class="result-icon">
              <el-icon><Document /></el-icon>
            </div>
            <div class="result-info">
              <div class="result-number">{{ importResult.total }}</div>
              <div class="result-label">总计记录</div>
            </div>
          </div>
          
          <div v-if="importResult.failed > 0" class="result-card failed">
            <div class="result-icon">
              <el-icon><CircleClose /></el-icon>
            </div>
            <div class="result-info">
              <div class="result-number">{{ importResult.failed }}</div>
              <div class="result-label">导入失败</div>
            </div>
          </div>
        </div>
        
        <div v-if="importResult.failed_items && importResult.failed_items.length > 0" class="failed-list">
          <h4>失败详情：</h4>
          <el-scrollbar max-height="200px">
            <div class="failed-item" v-for="(item, index) in importResult.failed_items" :key="index">
              {{ item }}
            </div>
          </el-scrollbar>
        </div>
      </div>
    </div>
    
    <template #footer>
      <div class="modal-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button v-if="currentStep > 0 && currentStep < 3" @click="prevStep">上一步</el-button>
        <el-button 
          v-if="currentStep < 2" 
          type="primary" 
          @click="nextStep"
          :disabled="!canNext"
        >
          下一步
        </el-button>
        <el-button 
          v-if="currentStep === 2" 
          type="primary" 
          @click="handleImport"
          :loading="importing"
          :disabled="validCount === 0"
        >
          开始导入
        </el-button>
        <el-button 
          v-if="currentStep === 3" 
          type="primary" 
          @click="handleFinish"
        >
          完成
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Download, 
  Upload, 
  Document, 
  InfoFilled,
  UploadFilled,
  Delete,
  View,
  CircleCheck,
  CircleClose,
  WarningFilled
} from '@element-plus/icons-vue'
import { domainApi } from '@/api/dns/domain'
import * as XLSX from 'xlsx'

interface Props {
  visible: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:visible': [value: boolean]
  'success': []
}>()

// 响应式数据
const currentStep = ref(0)
const uploadRef = ref()
const uploadFile = ref<File | null>(null)
const previewData = ref<any[]>([])
const importing = ref(false)
const importResult = ref({
  success: 0,
  failed: 0,
  total: 0,
  failed_items: [] as string[]
})

const visible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const canNext = computed(() => {
  if (currentStep.value === 0) return true
  if (currentStep.value === 1) return !!uploadFile.value
  return false
})

const validCount = computed(() => {
  return previewData.value.filter(item => !item._error).length
})

const invalidCount = computed(() => {
  return previewData.value.filter(item => item._error).length
})

// 下载模板
const downloadTemplate = async (format: 'xlsx' | 'csv') => {
  try {
    const templateData = [
      {
        name: 'example.com',
        registrar_type: 'aliyun',
        status: 'active',
        group_id: '',
        expires_at: '2024-12-31 23:59:59',
        auto_renew: false,
        tags: '生产,重要',
        remark: '示例域名'
      }
    ]
    
    const ws = XLSX.utils.json_to_sheet(templateData)
    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, '域名数据')
    
    const fileName = `域名导入模板.${format}`
    XLSX.writeFile(wb, fileName)
    
    ElMessage.success(`${format.toUpperCase()} 模板下载成功`)
  } catch (error) {
    ElMessage.error('模板下载失败')
  }
}

// 文件上传处理
const handleFileChange = (file: any) => {
  uploadFile.value = file.raw
  parseFile(file.raw)
}

const handleExceed = () => {
  ElMessage.warning('只能上传一个文件')
}

const removeFile = () => {
  uploadFile.value = null
  previewData.value = []
  uploadRef.value?.clearFiles()
}

const formatFileSize = (size: number) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return Math.round(size / 1024) + ' KB'
  return Math.round(size / (1024 * 1024)) + ' MB'
}

// 解析文件
const parseFile = async (file: File) => {
  try {
    const buffer = await file.arrayBuffer()
    const workbook = XLSX.read(buffer, { type: 'buffer' })
    const worksheet = workbook.Sheets[workbook.SheetNames[0]]
    const jsonData = XLSX.utils.sheet_to_json(worksheet)
    
    // 验证和转换数据
    previewData.value = jsonData.map((row: any, index: number) => {
      const item = {
        name: row.name || '',
        registrar_type: row.registrar_type || '',
        status: row.status || 'active',
        group_id: row.group_id || null,
        expires_at: row.expires_at || '',
        auto_renew: row.auto_renew === true || row.auto_renew === 'true',
        tags: row.tags ? row.tags.split(',').map((t: string) => t.trim()) : [],
        remark: row.remark || '',
        _error: ''
      }
      
      // 验证必填字段
      if (!item.name) {
        item._error = '域名不能为空'
      } else if (!/^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$/.test(item.name)) {
        item._error = '域名格式不正确'
      } else if (!item.registrar_type) {
        item._error = '注册商不能为空'
      } else if (!['aliyun', 'tencent', 'route53', 'cloudflare', 'dnspod', 'godaddy'].includes(item.registrar_type)) {
        item._error = '注册商类型不支持'
      }
      
      return item
    })
    
    ElMessage.success(`解析成功，共 ${previewData.value.length} 条记录`)
  } catch (error) {
    ElMessage.error('文件解析失败，请检查文件格式')
  }
}

// 步骤控制
const nextStep = () => {
  if (currentStep.value < 3) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

// 导入数据
const handleImport = async () => {
  try {
    importing.value = true
    
    const validData = previewData.value.filter(item => !item._error)
    
    const result = await domainApi.import({
      data: validData,
      format: 'json'
    })
    
    importResult.value = result
    currentStep.value = 3
    
    ElMessage.success('导入完成')
  } catch (error) {
    ElMessage.error('导入失败')
  } finally {
    importing.value = false
  }
}

const handleCancel = () => {
  visible.value = false
  resetModal()
}

const handleFinish = () => {
  visible.value = false
  emit('success')
  resetModal()
}

const resetModal = () => {
  currentStep.value = 0
  uploadFile.value = null
  previewData.value = []
  importing.value = false
  importResult.value = {
    success: 0,
    failed: 0,
    total: 0,
    failed_items: []
  }
  uploadRef.value?.clearFiles()
}
</script>

<style scoped>
.import-modal {
  .import-content {
    .import-steps {
      margin-bottom: 32px;
    }
    
    .step-content {
      .step-header {
        text-align: center;
        margin-bottom: 24px;
        
        h3 {
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 8px;
          font-size: 18px;
          font-weight: 600;
          color: #1f2937;
          margin-bottom: 8px;
        }
        
        p {
          color: #6b7280;
          margin: 0;
        }
      }
      
      .template-options {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 16px;
        margin-bottom: 24px;
        
        .template-card {
          display: flex;
          align-items: center;
          gap: 16px;
          padding: 20px;
          border: 2px solid #e5e7eb;
          border-radius: 12px;
          cursor: pointer;
          transition: all 0.3s ease;
          
          &:hover {
            border-color: #3b82f6;
            box-shadow: 0 4px 12px rgba(59, 130, 246, 0.15);
          }
          
          .template-icon {
            color: #3b82f6;
          }
          
          .template-info {
            flex: 1;
            
            h4 {
              margin: 0 0 4px 0;
              font-size: 16px;
              font-weight: 600;
            }
            
            p {
              margin: 0 0 8px 0;
              color: #6b7280;
              font-size: 14px;
            }
          }
        }
      }
      
      .template-guide {
        .guide-content {
          .guide-item {
            display: flex;
            margin-bottom: 12px;
            
            strong {
              min-width: 160px;
              color: #374151;
            }
            
            span {
              color: #6b7280;
              font-size: 14px;
            }
          }
        }
      }
      
      .upload-dragger {
        .upload-content {
          text-align: center;
          padding: 40px;
          
          .upload-icon {
            font-size: 48px;
            color: #3b82f6;
            margin-bottom: 16px;
          }
          
          .upload-text {
            p {
              margin: 8px 0;
              
              &.upload-tip {
                font-size: 12px;
                color: #9ca3af;
              }
              
              em {
                color: #3b82f6;
                font-style: normal;
              }
            }
          }
        }
      }
      
      .file-info {
        margin-top: 16px;
        
        .file-item {
          display: flex;
          align-items: center;
          gap: 12px;
          padding: 12px 16px;
          background: #f3f4f6;
          border-radius: 8px;
          
          .file-name {
            flex: 1;
            font-weight: 500;
          }
          
          .file-size {
            color: #6b7280;
            font-size: 14px;
          }
        }
      }
      
      .preview-summary {
        display: flex;
        gap: 12px;
        margin-bottom: 16px;
        
        .el-tag {
          font-weight: 600;
        }
      }
      
      .preview-table {
        .cell-content {
          display: flex;
          align-items: center;
          gap: 8px;
          
          .error-icon {
            color: #ef4444;
          }
        }
        
        .error-text {
          color: #ef4444;
          font-size: 12px;
        }
      }
      
      .result-summary {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 16px;
        margin-bottom: 24px;
        
        .result-card {
          display: flex;
          align-items: center;
          gap: 16px;
          padding: 20px;
          border-radius: 12px;
          
          &.success {
            background: linear-gradient(135deg, #10b981 0%, #059669 100%);
            color: white;
          }
          
          &.total {
            background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
            color: white;
          }
          
          &.failed {
            background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
            color: white;
          }
          
          .result-icon {
            font-size: 32px;
          }
          
          .result-info {
            .result-number {
              font-size: 24px;
              font-weight: 700;
              margin-bottom: 4px;
            }
            
            .result-label {
              font-size: 14px;
              opacity: 0.9;
            }
          }
        }
      }
      
      .failed-list {
        h4 {
          margin-bottom: 12px;
          color: #ef4444;
        }
        
        .failed-item {
          padding: 8px 12px;
          background: #fef2f2;
          border: 1px solid #fecaca;
          border-radius: 6px;
          margin-bottom: 8px;
          font-size: 14px;
          color: #991b1b;
        }
      }
    }
  }
  
  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }
}

/* Element Plus 样式覆盖 */
:deep(.el-upload-dragger) {
  border: 2px dashed #d1d5db;
  border-radius: 12px;
  background: #f9fafb;
  transition: all 0.3s ease;
  
  &:hover {
    border-color: #3b82f6;
    background: #eff6ff;
  }
}

:deep(.el-dialog) {
  border-radius: 12px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
}

:deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 12px 12px 0 0;
  padding: 20px 24px;
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-dialog__footer) {
  padding: 20px 24px;
  background: #f9fafb;
  border-radius: 0 0 12px 12px;
}
</style>
