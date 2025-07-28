<template>
  <el-dialog
    v-model="dialogVisible"
    title=""
    width="900px"
    :close-on-click-modal="false"
    class="batch-import-dialog"
    :show-close="false"
  >
    <!-- 自定义头部 -->
    <template #header>
      <div class="dialog-header">
        <div class="header-left">
          <div class="header-icon">
            <el-icon size="24"><Upload /></el-icon>
          </div>
          <div class="header-content">
            <h2 class="dialog-title">批量导入主机</h2>
            <p class="dialog-subtitle">快速导入多台主机到系统中</p>
          </div>
        </div>
        <el-button 
          type="text" 
          size="large" 
          @click="handleCancel"
          class="close-btn"
        >
          <el-icon size="20"><Close /></el-icon>
        </el-button>
      </div>
    </template>

    <div class="import-container">
      <!-- 进度指示器 -->
      <div class="progress-indicator">
        <div class="steps-wrapper">
          <div 
            v-for="(step, index) in steps" 
            :key="index"
            :class="['step-item', { 
              active: currentStep === index, 
              completed: currentStep > index 
            }]"
          >
            <div class="step-circle">
              <el-icon v-if="currentStep > index" class="check-icon">
                <Check />
              </el-icon>
              <span v-else>{{ index + 1 }}</span>
            </div>
            <div class="step-content">
              <div class="step-title">{{ step.title }}</div>
              <div class="step-desc">{{ step.desc }}</div>
            </div>
            <div v-if="index < steps.length - 1" class="step-line"></div>
          </div>
        </div>
      </div>

      <!-- 步骤内容 -->
      <div class="step-content-wrapper">
        <!-- 步骤 0: 选择文件 -->
        <div v-if="currentStep === 0" class="step-panel">
          <div class="panel-header">
            <h3>选择导入文件</h3>
            <p>支持 Excel (.xlsx) 和 CSV (.csv) 格式文件</p>
          </div>

          <div class="upload-section">
            <el-upload
              class="upload-dragger"
              drag
              :auto-upload="false"
              :on-change="handleFileChange"
              :show-file-list="false"
              accept=".xlsx,.csv"
              :class="{ 'has-file': selectedFile }"
            >
              <div class="upload-content">
                <div class="upload-icon">
                  <el-icon size="48" color="#409EFF">
                    <UploadFilled />
                  </el-icon>
                </div>
                <div class="upload-text">
                  <div class="primary-text">点击或拖拽文件到此处</div>
                  <div class="secondary-text">支持 .xlsx / .csv 格式，文件大小不超过 10MB</div>
                </div>
              </div>
            </el-upload>

            <!-- 文件信息 -->
            <div v-if="selectedFile" class="file-info">
              <div class="file-card">
                <div class="file-icon">
                  <el-icon size="24" color="#67C23A">
                    <Document />
                  </el-icon>
                </div>
                <div class="file-details">
                  <div class="file-name">{{ selectedFile.name }}</div>
                  <div class="file-size">{{ formatFileSize(selectedFile.size) }}</div>
                </div>
                <el-button 
                  type="text" 
                  @click="selectedFile = null"
                  class="remove-file"
                >
                  <el-icon><Close /></el-icon>
                </el-button>
              </div>
            </div>
          </div>

          <!-- 模板下载区域 -->
          <div class="template-section">
            <div class="template-card">
              <div class="template-header">
                <div class="template-icon">
                  <el-icon size="20" color="#E6A23C">
                    <Download />
                  </el-icon>
                </div>
                <div class="template-content">
                  <h4>需要导入模板？</h4>
                  <p>下载标准模板，按照格式填写数据后导入</p>
                </div>
              </div>
              <el-button type="primary" plain @click="downloadTemplate">
                下载模板
              </el-button>
            </div>
          </div>
        </div>

        <!-- 步骤 1: 数据预览 -->
        <div v-if="currentStep === 1" class="step-panel">
          <div class="panel-header">
            <h3>数据预览</h3>
            <p>请检查解析的数据是否正确，确认无误后继续导入</p>
          </div>

          <div v-if="previewData.length > 0" class="preview-section">
            <!-- 统计信息 -->
            <div class="stats-cards">
              <div class="stat-card total">
                <div class="stat-icon">
                  <el-icon size="24"><Document /></el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-number">{{ previewData.length }}</div>
                  <div class="stat-label">总记录数</div>
                </div>
              </div>
              <div class="stat-card valid">
                <div class="stat-icon">
                  <el-icon size="24" color="#67C23A"><CircleCheck /></el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-number">{{ validCount }}</div>
                  <div class="stat-label">有效记录</div>
                </div>
              </div>
              <div class="stat-card invalid">
                <div class="stat-icon">
                  <el-icon size="24" color="#F56C6C"><CircleClose /></el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-number">{{ invalidCount }}</div>
                  <div class="stat-label">无效记录</div>
                </div>
              </div>
            </div>

            <!-- 数据表格 -->
            <div class="preview-table-wrapper">
              <el-table
                :data="previewData.slice(0, 10)"
                border
                size="default"
                max-height="400"
                class="preview-table"
                :row-class-name="getRowClassName"
              >
                <el-table-column type="index" label="#" width="50" />
                <el-table-column prop="name" label="主机名称" min-width="120">
                  <template #default="{ row }">
                    <div class="cell-content">
                      <span :class="{ 'error-text': !row.name }">
                        {{ row.name || '缺失' }}
                      </span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="public_ip" label="公网IP" min-width="120">
                  <template #default="{ row }">
                    <div class="cell-content">
                      <span :class="{ 'error-text': !row.public_ip }">
                        {{ row.public_ip || '缺失' }}
                      </span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="private_ip" label="私网IP" min-width="120" />
                <el-table-column prop="username" label="用户名" min-width="100">
                  <template #default="{ row }">
                    <div class="cell-content">
                      <span :class="{ 'error-text': !row.username }">
                        {{ row.username || '缺失' }}
                      </span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="password" label="密码" min-width="100">
                  <template #default="{ row }">
                    <div class="cell-content">
                      <span v-if="row.password" class="password-mask">••••••</span>
                      <span v-else class="error-text">缺失</span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="os" label="操作系统" min-width="100" />
                <el-table-column prop="provider_type" label="提供商" min-width="100" />
                <el-table-column label="状态" width="80" fixed="right">
                  <template #default="{ row }">
                    <el-tag 
                      :type="row._valid ? 'success' : 'danger'" 
                      size="small"
                      effect="light"
                    >
                      {{ row._valid ? '有效' : '无效' }}
                    </el-tag>
                  </template>
                </el-table-column>
              </el-table>

              <div v-if="previewData.length > 10" class="table-footer">
                <el-text type="info">
                  <el-icon><InfoFilled /></el-icon>
                  仅显示前 10 条记录，共 {{ previewData.length }} 条数据
                </el-text>
              </div>
            </div>

            <!-- 错误提示 -->
            <div v-if="invalidCount > 0" class="error-summary">
              <el-alert
                title="数据验证警告"
                type="warning"
                show-icon
                :closable="false"
              >
                <template #default>
                  <p>发现 {{ invalidCount }} 条无效记录，这些记录将被跳过。</p>
                  <p>请确保必填字段（主机名称、公网IP、用户名、密码）不为空。</p>
                </template>
              </el-alert>
            </div>
          </div>

          <div v-else class="empty-data">
            <el-empty 
              description="没有解析到有效数据"
              :image-size="120"
            >
              <el-button type="primary" @click="currentStep = 0">
                重新选择文件
              </el-button>
            </el-empty>
          </div>
        </div>

        <!-- 步骤 2: 导入结果 -->
        <div v-if="currentStep === 2" class="step-panel">
          <div class="result-section">
            <div class="result-icon">
              <el-icon size="64" color="#67C23A">
                <CircleCheck />
              </el-icon>
            </div>
            <h3 class="result-title">导入完成</h3>
            <p class="result-desc">
              成功导入 {{ importResult?.success || 0 }} 条主机记录
              <span v-if="importResult?.failed > 0">
                ，跳过 {{ importResult.failed }} 条无效记录
              </span>
            </p>
            
            <!-- 导入统计 -->
            <div class="import-stats">
              <div class="stat-item">
                <span class="stat-label">总计：</span>
                <span class="stat-value">{{ importResult?.total || 0 }}</span>
              </div>
              <div class="stat-item success">
                <span class="stat-label">成功：</span>
                <span class="stat-value">{{ importResult?.success || 0 }}</span>
              </div>
              <div v-if="importResult?.failed > 0" class="stat-item failed">
                <span class="stat-label">失败：</span>
                <span class="stat-value">{{ importResult.failed }}</span>
              </div>
            </div>

            <!-- 错误信息 -->
            <div v-if="importResult?.failed_msg?.length" class="error-messages">
              <el-collapse>
                <el-collapse-item title="查看错误详情" name="errors">
                  <ul class="error-list">
                    <li v-for="(error, index) in importResult.failed_msg" :key="index">
                      {{ error }}
                    </li>
                  </ul>
                </el-collapse-item>
              </el-collapse>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部操作栏 -->
    <template #footer>
      <div class="dialog-footer">
        <div class="footer-left">
          <el-button 
            v-if="currentStep > 0 && currentStep < 2" 
            @click="handlePrevious"
            :icon="ArrowLeft"
          >
            上一步
          </el-button>
        </div>
        <div class="footer-right">
          <el-button @click="handleCancel">
            {{ currentStep === 2 ? '关闭' : '取消' }}
          </el-button>
          <el-button
            v-if="currentStep === 0"
            type="primary"
            :disabled="!selectedFile"
            @click="handleNext"
            :icon="ArrowRight"
          >
            下一步
          </el-button>
          <el-button 
            v-if="currentStep === 1" 
            type="primary" 
            :loading="loading" 
            @click="handleImport"
            :icon="Upload"
          >
            {{ loading ? '导入中...' : '确认导入' }}
          </el-button>
          <el-button 
            v-if="currentStep === 2" 
            type="primary" 
            @click="handleFinish"
            :icon="Check"
          >
            完成
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Upload,
  Close,
  Check,
  UploadFilled,
  Document,
  Download,
  CircleCheck,
  CircleClose,
  InfoFilled,
  ArrowLeft,
  ArrowRight
} from '@element-plus/icons-vue'
import * as XLSX from 'xlsx'
import { batchImportHosts } from '@/api/cmdb/host'

interface Props {
  modelValue: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 步骤配置
const steps = [
  { title: '选择文件', desc: '上传Excel或CSV文件' },
  { title: '数据预览', desc: '检查数据格式和内容' },
  { title: '导入完成', desc: '查看导入结果' }
]

// 状态管理
const currentStep = ref(0)
const loading = ref(false)
const selectedFile = ref<File | null>(null)
const previewData = ref<any[]>([])
const importResult = ref<any>(null)

// 计算属性
const validCount = computed(() => {
  return previewData.value.filter(item => item._valid).length
})

const invalidCount = computed(() => {
  return previewData.value.filter(item => !item._valid).length
})

// 文件处理
const handleFileChange = (file: any) => {
  const rawFile = file.raw

  // 文件大小检查
  if (rawFile.size > 10 * 1024 * 1024) {
    ElMessage.error('文件大小不能超过 10MB')
    return
  }

  // 文件类型检查
  const isValidType =
    rawFile.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' ||
    rawFile.type === 'text/csv' ||
    rawFile.type === 'application/csv' ||
    rawFile.name.endsWith('.xlsx') ||
    rawFile.name.endsWith('.csv')

  if (!isValidType) {
    ElMessage.error('只支持 .xlsx 和 .csv 格式的文件')
    return
  }

  selectedFile.value = rawFile
}

// 文件大小格式化
const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 解析文件数据
const parseFileData = async (file: File): Promise<any[]> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()

    reader.onload = (e) => {
      try {
        const data = e.target?.result
        let jsonData: any[] = []

        if (file.name.endsWith('.csv')) {
          // 解析CSV
          const text = data as string
          const lines = text.split('\n').filter(line => line.trim())
          if (lines.length < 2) {
            reject(new Error('CSV文件格式不正确'))
            return
          }

          const headers = lines[0].split(',').map(h => h.trim().replace(/"/g, ''))
          jsonData = lines.slice(1).map(line => {
            const values = line.split(',').map(v => v.trim().replace(/"/g, ''))
            const obj: any = {}
            headers.forEach((header, index) => {
              obj[header] = values[index] || ''
            })
            return obj
          })
        } else {
          // 解析Excel
          const workbook = XLSX.read(data, { type: 'binary' })
          const sheetName = workbook.SheetNames[0]
          const worksheet = workbook.Sheets[sheetName]
          jsonData = XLSX.utils.sheet_to_json(worksheet)
        }

        // 数据验证和标准化
        const processedData = jsonData.map(item => {
          const processed = {
            name: item['主机名称'] || item['name'] || '',
            public_ip: item['公网IP'] || item['public_ip'] || '',
            private_ip: item['私网IP'] || item['private_ip'] || '',
            username: item['用户名'] || item['username'] || '',
            password: item['密码'] || item['password'] || '',
            os: item['操作系统'] || item['os'] || '',
            provider_type: item['提供商类型'] || item['provider_type'] || '',
            region: item['区域'] || item['region'] || '',
            configuration: item['配置规格'] || item['configuration'] || '',
            remark: item['备注'] || item['remark'] || ''
          }

          // 验证必填字段
          processed._valid = !!(processed.name && processed.public_ip && processed.username && processed.password)

          return processed
        })

        resolve(processedData)
      } catch (error) {
        reject(error)
      }
    }

    reader.onerror = () => reject(new Error('文件读取失败'))

    if (file.name.endsWith('.csv')) {
      reader.readAsText(file, 'UTF-8')
    } else {
      reader.readAsBinaryString(file)
    }
  })
}

// 表格行样式
const getRowClassName = ({ row }: { row: any }) => {
  return row._valid ? '' : 'invalid-row'
}

// 下载模板
const downloadTemplate = () => {
  const templateData = [
    ['主机名称', '公网IP', '私网IP', '用户名', '密码', '操作系统', '提供商类型', '区域', '配置规格', '备注'],
    ['web-server-01', '1.2.3.4', '10.0.0.1', 'root', 'password123', 'CentOS 7', 'aliyun', 'cn-beijing', '2核4G', '测试服务器'],
    ['db-server-01', '1.2.3.5', '10.0.0.2', 'root', 'password456', 'Ubuntu 20.04', 'tencent', 'ap-shanghai', '4核8G', '数据库服务器']
  ]

  const ws = XLSX.utils.aoa_to_sheet(templateData)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, '主机导入模板')
  XLSX.writeFile(wb, '主机导入模板.xlsx')
}
