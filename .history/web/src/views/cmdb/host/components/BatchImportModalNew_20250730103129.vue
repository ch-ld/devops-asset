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
                max-height="500"
                class="preview-table"
                :row-class-name="getRowClassName"
              >
                <el-table-column type="index" label="#" width="50" />
                <el-table-column prop="instance_id" label="实例ID" min-width="140" fixed="left">
                  <template #default="{ row }">
                    <div class="cell-content">
                      <span :class="{ 'error-text': !row.instance_id }">
                        {{ row.instance_id || '缺失' }}
                      </span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="name" label="主机名称" min-width="120" fixed="left">
                  <template #default="{ row }">
                    <div class="cell-content">
                      <span :class="{ 'error-text': !row.name }">
                        {{ row.name || '缺失' }}
                      </span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="group_name" label="主机组" min-width="100" />
                <el-table-column prop="resource_type" label="资源类型" min-width="90" />
                <el-table-column prop="provider_id" label="云厂商ID" min-width="80">
                  <template #default="{ row }">
                    <span>{{ row.provider_id || '-' }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="public_ip" label="公网IP" min-width="140">
                  <template #default="{ row }">
                    <div class="cell-content">
                      <span class="ip-display">{{ formatIPDisplay(row.public_ip) }}</span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="private_ip" label="私网IP" min-width="140">
                  <template #default="{ row }">
                    <div class="cell-content">
                      <span class="ip-display">{{ formatIPDisplay(row.private_ip) }}</span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="username" label="用户名" min-width="100">
                  <template #default="{ row }">
                    <div class="cell-content">
                      <span :class="{ 'error-text': !row.username }">
                        {{ row.username || '缺失' }}
                      </span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="password" label="密码" min-width="80">
                  <template #default="{ row }">
                    <div class="cell-content">
                      <span v-if="row.password" class="password-mask">••••••</span>
                      <span v-else class="password-placeholder">-</span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="os" label="操作系统" min-width="110" />
                <el-table-column prop="region" label="地域" min-width="100" />
                <el-table-column prop="status" label="状态" min-width="80" />
                <el-table-column label="验证状态" width="100" fixed="right">
                  <template #default="{ row }">
                    <div class="validation-cell">
                      <el-tag 
                        :type="row._valid ? 'success' : 'danger'" 
                        size="small"
                        effect="light"
                      >
                        {{ row._valid ? '有效' : '无效' }}
                      </el-tag>
                      <el-tooltip v-if="!row._valid && row._errors" placement="top" width="300">
                        <template #content>
                          <div class="error-tooltip">
                            <div v-for="error in row._errors" :key="error" class="error-item">
                              • {{ error }}
                            </div>
                          </div>
                        </template>
                        <el-icon class="error-icon" color="#F56C6C" size="14">
                          <InfoFilled />
                        </el-icon>
                      </el-tooltip>
                    </div>
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
import { batchImportHosts } from '@/api/system/host'

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

        // 数据验证和标准化 - 更新以支持完整字段
        const processedData = jsonData.map((item, index) => {
          try {
            const processed = {
              // 必填字段
              provider_id: item['云厂商ID'] || item['provider_id'] || '',
              instance_id: item['实例ID'] || item['instance_id'] || '',
              name: item['主机名称'] || item['name'] || '',
              group_name: item['主机组'] || item['group_name'] || '',
              resource_type: item['资源类型'] || item['resource_type'] || 'manual',
              region: item['地域'] || item['region'] || '',
              username: item['用户名'] || item['username'] || '',
              password: item['密码'] || item['password'] || '',
              public_ip: item['公网IP'] || item['public_ip'] || '',
              private_ip: item['私网IP'] || item['private_ip'] || '',
              configuration: item['配置信息'] || item['configuration'] || '',
              os: item['操作系统'] || item['os'] || '',
              status: item['状态'] || item['status'] || 'running',
              expired_at: item['过期时间'] || item['expired_at'] || '',
              remark: item['备注'] || item['remark'] || ''
            }

            // 数据格式化和验证
            const errors: string[] = []

            // 验证必填字段
            if (!processed.instance_id) {
              errors.push('实例ID不能为空')
            }
            if (!processed.name) {
              errors.push('主机名称不能为空')
            }
            if (!processed.username) {
              errors.push('用户名不能为空')
            }

            // 处理云厂商ID
            if (processed.provider_id && processed.provider_id !== '') {
              const providerId = parseInt(processed.provider_id.toString())
              if (isNaN(providerId)) {
                errors.push('云厂商ID必须是数字')
              } else {
                processed.provider_id = providerId
              }
            } else {
              processed.provider_id = null
            }

            // 处理IP地址 - 支持JSON数组格式
            try {
              if (processed.public_ip) {
                if (typeof processed.public_ip === 'string' && processed.public_ip.startsWith('[')) {
                  JSON.parse(processed.public_ip) // 验证JSON格式
                } else if (typeof processed.public_ip === 'string') {
                  // 转换普通IP为JSON数组格式
                  processed.public_ip = JSON.stringify([processed.public_ip])
                }
              }
            } catch (e) {
              errors.push('公网IP格式错误，应为JSON数组格式')
            }

            try {
              if (processed.private_ip) {
                if (typeof processed.private_ip === 'string' && processed.private_ip.startsWith('[')) {
                  JSON.parse(processed.private_ip) // 验证JSON格式
                } else if (typeof processed.private_ip === 'string') {
                  // 转换普通IP为JSON数组格式
                  processed.private_ip = JSON.stringify([processed.private_ip])
                }
              }
            } catch (e) {
              errors.push('私网IP格式错误，应为JSON数组格式')
            }

            // 处理配置信息 - 支持JSON对象格式
            try {
              if (processed.configuration && typeof processed.configuration === 'string') {
                if (processed.configuration.startsWith('{')) {
                  JSON.parse(processed.configuration) // 验证JSON格式
                } else {
                  // 尝试从简单字符串转换为JSON对象
                  const simpleConfig = { description: processed.configuration }
                  processed.configuration = JSON.stringify(simpleConfig)
                }
              }
            } catch (e) {
              errors.push('配置信息格式错误，应为JSON对象格式')
            }

            // 验证时间格式
            if (processed.expired_at && processed.expired_at !== '') {
              const timeRegex = /^\d{4}-\d{2}-\d{2}(\s\d{2}:\d{2}:\d{2})?$/
              if (!timeRegex.test(processed.expired_at)) {
                errors.push('过期时间格式错误，应为 YYYY-MM-DD 或 YYYY-MM-DD HH:mm:ss')
              }
            }

            // 设置验证状态和错误信息
            processed._valid = errors.length === 0
            processed._errors = errors
            processed._rowIndex = index + 2 // Excel行号（从2开始，因为第1行是标题）

            return processed
          } catch (error) {
            return {
              _valid: false,
              _errors: [`第${index + 2}行解析失败: ${error.message}`],
              _rowIndex: index + 2,
              name: `解析失败的行 ${index + 2}`,
              instance_id: '',
              username: '',
              // ... 其他字段设为默认值
            }
          }
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
    // 头部字段 - 与后端期望字段完全对应
    [
      '云厂商ID', '实例ID', '主机名称', '主机组', '资源类型', '地域', 
      '用户名', '密码', '公网IP', '私网IP', '配置信息', '操作系统', 
      '状态', '过期时间', '备注'
    ],
    // 示例数据行1 - 阿里云ECS
    [
      '1', 'i-bp1234567890abcdef', 'web-server-01', '生产环境', 'ecs', 'cn-beijing',
      'root', 'password123', '["1.2.3.4"]', '["10.0.0.1"]', '{"cpu_cores":2,"memory_size":4,"instance_type":"ecs.c6.large"}', 'CentOS 7',
      'running', '2024-12-31 23:59:59', '测试服务器'
    ],
    // 示例数据行2 - 腾讯云CVM
    [
      '2', 'ins-1234567890', 'db-server-01', '测试环境', 'cvm', 'ap-shanghai',
      'ubuntu', 'password456', '["1.2.3.5"]', '["10.0.0.2"]', '{"cpu_cores":4,"memory_size":8,"instance_type":"SA2.MEDIUM4"}', 'Ubuntu 20.04',
      'running', '2024-12-31 23:59:59', '数据库服务器'
    ],
    // 示例数据行3 - 自建主机
    [
      '', 'manual-001', 'app-server-01', '开发环境', 'manual', 'local',
      'admin', 'admin123', '["192.168.1.100"]', '["192.168.1.100"]', '{"cpu_cores":8,"memory_size":16,"disk_size":500}', 'CentOS 8',
      'running', '', '开发应用服务器'
    ]
  ]

  const ws = XLSX.utils.aoa_to_sheet(templateData)
  
  // 设置列宽
  const colWidths = [
    { wch: 10 }, // 云厂商ID
    { wch: 20 }, // 实例ID
    { wch: 15 }, // 主机名称
    { wch: 12 }, // 主机组
    { wch: 12 }, // 资源类型
    { wch: 12 }, // 地域
    { wch: 10 }, // 用户名
    { wch: 10 }, // 密码
    { wch: 15 }, // 公网IP
    { wch: 15 }, // 私网IP
    { wch: 25 }, // 配置信息
    { wch: 12 }, // 操作系统
    { wch: 10 }, // 状态
    { wch: 20 }, // 过期时间
    { wch: 15 }  // 备注
  ]
  ws['!cols'] = colWidths

  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, '主机导入模板')
  
  // 添加说明工作表
  const instructionData = [
    ['字段说明'],
    [''],
    ['必填字段：'],
    ['实例ID', '主机的唯一标识符，必须唯一'],
    ['主机名称', '主机的显示名称'],
    ['主机组', '主机所属分组名称'],
    ['资源类型', 'ecs/cvm/ec2/manual 等'],
    ['用户名', 'SSH登录用户名'],
    [''],
    ['可选字段：'],
    ['云厂商ID', '云账号ID，自建主机可为空'],
    ['密码', 'SSH登录密码'],
    ['地域', '主机所在地域'],
    ['公网IP', 'JSON数组格式，如：["1.2.3.4"]'],
    ['私网IP', 'JSON数组格式，如：["10.0.0.1"]'],
    ['配置信息', 'JSON格式，如：{"cpu_cores":2,"memory_size":4}'],
    ['操作系统', '如：CentOS 7, Ubuntu 20.04'],
    ['状态', 'running/stopped/error'],
    ['过期时间', '格式：YYYY-MM-DD HH:mm:ss'],
    ['备注', '主机备注信息'],
    [''],
    ['注意事项：'],
    ['1. IP地址使用JSON数组格式'],
    ['2. 配置信息使用JSON对象格式'],
    ['3. 时间格式为 YYYY-MM-DD HH:mm:ss'],
    ['4. 云厂商ID对应系统中的云账号']
  ]
  
  const instructionWs = XLSX.utils.aoa_to_sheet(instructionData)
  instructionWs['!cols'] = [{ wch: 15 }, { wch: 40 }]
  XLSX.utils.book_append_sheet(wb, instructionWs, '字段说明')
  
  XLSX.writeFile(wb, '主机批量导入模板.xlsx')
  ElMessage.success('模板下载成功')
}

// 步骤操作
const handleNext = async () => {
  if (currentStep.value === 0 && selectedFile.value) {
    try {
      loading.value = true
      const data = await parseFileData(selectedFile.value)
      previewData.value = data
      currentStep.value = 1
    } catch (error) {
      ElMessage.error('文件解析失败：' + (error as Error).message)
    } finally {
      loading.value = false
    }
  }
}

const handlePrevious = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

const handleImport = async () => {
  try {
    loading.value = true
    const validData = previewData.value.filter(item => item._valid)

    if (validData.length === 0) {
      ElMessage.error('没有有效的数据可以导入')
      return
    }

    const response = await batchImportHosts(validData)
    importResult.value = response.data
    currentStep.value = 2

    ElMessage.success('导入完成')
  } catch (error) {
    ElMessage.error('导入失败：' + (error as Error).message)
  } finally {
    loading.value = false
  }
}

const handleFinish = () => {
  emit('success')
  handleCancel()
}

const handleCancel = () => {
  dialogVisible.value = false
  // 重置状态
  currentStep.value = 0
  selectedFile.value = null
  previewData.value = []
  importResult.value = null
  loading.value = false
}

// 监听对话框关闭
watch(dialogVisible, (newVal) => {
  if (!newVal) {
    handleCancel()
  }
})
</script>

<style scoped>
.batch-import-dialog {
  --el-dialog-border-radius: 12px;
}

.batch-import-dialog :deep(.el-dialog__header) {
  padding: 0;
  margin: 0;
}

.batch-import-dialog :deep(.el-dialog__body) {
  padding: 0 24px 24px;
}

.batch-import-dialog :deep(.el-dialog__footer) {
  padding: 0 24px 24px;
}

/* 对话框头部 */
.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 24px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 12px 12px 0 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-icon {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dialog-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.dialog-subtitle {
  margin: 4px 0 0;
  font-size: 14px;
  opacity: 0.9;
}

.close-btn {
  color: white !important;
  background: rgba(255, 255, 255, 0.1) !important;
  border: none !important;
  border-radius: 8px !important;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.2) !important;
}

/* 进度指示器 */
.progress-indicator {
  margin: 32px 0;
}

.steps-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.step-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  flex: 1;
  max-width: 200px;
}

.step-circle {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #f5f7fa;
  border: 2px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  color: #909399;
  transition: all 0.3s ease;
  z-index: 2;
  position: relative;
}

.step-item.active .step-circle {
  background: #409eff;
  border-color: #409eff;
  color: white;
}

.step-item.completed .step-circle {
  background: #67c23a;
  border-color: #67c23a;
  color: white;
}

.check-icon {
  font-size: 18px;
}

.step-content {
  margin-top: 12px;
  text-align: center;
}

.step-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
}

.step-item.active .step-title {
  color: #409eff;
}

.step-item.completed .step-title {
  color: #67c23a;
}

.step-desc {
  font-size: 12px;
  color: #909399;
}

.step-line {
  position: absolute;
  top: 20px;
  left: 50%;
  right: -50%;
  height: 2px;
  background: #e4e7ed;
  z-index: 1;
}

.step-item.completed + .step-item .step-line {
  background: #67c23a;
}

/* 步骤面板 */
.step-content-wrapper {
  min-height: 400px;
}

.step-panel {
  animation: fadeInUp 0.3s ease;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.panel-header {
  text-align: center;
  margin-bottom: 32px;
}

.panel-header h3 {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.panel-header p {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

/* 上传区域 */
.upload-section {
  margin-bottom: 32px;
}

.upload-dragger {
  width: 100%;
}

.upload-dragger :deep(.el-upload-dragger) {
  width: 100%;
  height: 180px;
  border: 2px dashed #d9d9d9;
  border-radius: 12px;
  background: #fafbfc;
  transition: all 0.3s ease;
}

.upload-dragger :deep(.el-upload-dragger:hover) {
  border-color: #409eff;
  background: #f0f9ff;
}

.upload-dragger.has-file :deep(.el-upload-dragger) {
  border-color: #67c23a;
  background: #f0f9ff;
}

.upload-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 16px;
}

.upload-icon {
  opacity: 0.8;
}

.upload-text {
  text-align: center;
}

.primary-text {
  font-size: 16px;
  color: #303133;
  font-weight: 500;
  margin-bottom: 8px;
}

.secondary-text {
  font-size: 14px;
  color: #909399;
}

/* 文件信息 */
.file-info {
  margin-top: 16px;
}

.file-card {
  display: flex;
  align-items: center;
  padding: 16px;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  gap: 12px;
}

.file-details {
  flex: 1;
}

.file-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.file-size {
  font-size: 12px;
  color: #909399;
}

.remove-file {
  color: #f56c6c !important;
}

/* 模板区域 */
.template-section {
  border-top: 1px solid #ebeef5;
  padding-top: 24px;
}

.template-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  background: #fff7e6;
  border: 1px solid #ffd591;
  border-radius: 8px;
}

.template-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.template-content h4 {
  margin: 0 0 4px;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.template-content p {
  margin: 0;
  font-size: 12px;
  color: #909399;
}

/* 统计卡片 */
.stats-cards {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  flex: 1;
  display: flex;
  align-items: center;
  padding: 20px;
  background: white;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  gap: 16px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-card.total {
  border-color: #409eff;
  background: linear-gradient(135deg, #409eff10, #409eff05);
}

.stat-card.valid {
  border-color: #67c23a;
  background: linear-gradient(135deg, #67c23a10, #67c23a05);
}

.stat-card.invalid {
  border-color: #f56c6c;
  background: linear-gradient(135deg, #f56c6c10, #f56c6c05);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(64, 158, 255, 0.1);
}

.stat-card.valid .stat-icon {
  background: rgba(103, 194, 58, 0.1);
}

.stat-card.invalid .stat-icon {
  background: rgba(245, 108, 108, 0.1);
}

.stat-number {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

/* 预览表格 */
.preview-table-wrapper {
  margin-bottom: 24px;
}

.preview-table {
  border-radius: 8px;
  overflow: hidden;
}

.preview-table :deep(.el-table__header) {
  background: #f8f9fa;
}

.preview-table :deep(.el-table__header th) {
  background: #f8f9fa;
  color: #303133;
  font-weight: 600;
}

.preview-table :deep(.invalid-row) {
  background: #fef0f0;
}

.cell-content {
  display: flex;
  align-items: center;
}

.error-text {
  color: #f56c6c;
  font-style: italic;
}

.password-mask {
  color: #909399;
  font-family: monospace;
}

.table-footer {
  padding: 12px 16px;
  background: #f8f9fa;
  border: 1px solid #ebeef5;
  border-top: none;
  border-radius: 0 0 8px 8px;
  text-align: center;
}

/* 错误摘要 */
.error-summary {
  margin-top: 16px;
}

/* 空数据 */
.empty-data {
  text-align: center;
  padding: 60px 0;
}

/* 结果页面 */
.result-section {
  text-align: center;
  padding: 40px 0;
}

.result-icon {
  margin-bottom: 24px;
}

.result-title {
  margin: 0 0 12px;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.result-desc {
  margin: 0 0 32px;
  font-size: 16px;
  color: #606266;
  line-height: 1.5;
}

.import-stats {
  display: flex;
  justify-content: center;
  gap: 32px;
  margin-bottom: 24px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.stat-item.success .stat-value {
  color: #67c23a;
}

.stat-item.failed .stat-value {
  color: #f56c6c;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.error-messages {
  max-width: 600px;
  margin: 0 auto;
  text-align: left;
}

.error-list {
  margin: 0;
  padding: 0;
  list-style: none;
}

.error-list li {
  padding: 8px 0;
  color: #f56c6c;
  font-size: 14px;
  border-bottom: 1px solid #fde2e2;
}

.error-list li:last-child {
  border-bottom: none;
}

/* 底部操作栏 */
.dialog-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 24px;
  border-top: 1px solid #ebeef5;
}

.footer-left,
.footer-right {
  display: flex;
  gap: 12px;
}

.footer-right {
  margin-left: auto;
}
</style>
