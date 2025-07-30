<template>
  <el-dialog
    v-model="dialogVisible"
    title=""
    width="800px"
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
                  type="link"
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

      <!-- 步骤 1 -->
      <div v-if="currentStep === 1">
        <el-alert
          v-if="importResult"
          :title="`解析完成：共 ${importResult.total} 条，成功 ${importResult.success} 条，失败 ${importResult.failed} 条`"
          :type="importResult.failed > 0 ? 'warning' : 'success'"
          show-icon
          class="mb-2"
        />

        <div v-if="importResult && importResult.failed_msg && importResult.failed_msg.length">
          <h4>错误信息</h4>
          <ul class="error-list">
            <li v-for="(errMsg, idx) in importResult.failed_msg" :key="idx" class="text-danger">
              {{ errMsg }}
            </li>
          </ul>
        </div>
      </div>

      <!-- 步骤 2 -->
      <div v-if="currentStep === 2">
        <el-result
          icon="success"
          title="导入完成"
          :sub-title="`成功导入 ${importResult?.success || 0} 条主机记录`"
        />
      </div>
    </div>

    <template #footer>
      <el-space>
        <el-button v-if="currentStep > 0" @click="handlePrevious">上一步</el-button>
        <el-button @click="handleCancel">取消</el-button>
        <el-button
          v-if="currentStep === 0"
          type="primary"
          :disabled="!selectedFile"
          @click="handleNext"
        >
          下一步
        </el-button>
        <el-button v-if="currentStep === 1" type="primary" :loading="loading" @click="handleImport">
          确认导入
        </el-button>
        <el-button v-if="currentStep === 2" type="primary" @click="handleFinish">完成</el-button>
      </el-space>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
  import { ref, watch } from 'vue'
  import { ElMessage } from 'element-plus'
  import { Upload, Download } from '@element-plus/icons-vue'
  import { batchImportHosts, batchExportHosts } from '@/api/system/host'
  import type { BatchOperationResult } from '@/types/api/host'

  const props = defineProps<{
    visible: boolean
  }>()

  const emit = defineEmits(['success', 'update:visible'])

  const dialogVisible = ref(props.visible)
  const loading = ref(false)
  const currentStep = ref(0)
  const selectedFile = ref<File | null>(null)
  const importResult = ref<BatchOperationResult | null>(null)

  watch(
    () => props.visible,
    (val) => {
      dialogVisible.value = val
    }
  )
  watch(dialogVisible, (val) => emit('update:visible', val))

  const handleFileChange = (uploadFile: any) => {
    // 获取原始文件对象
    const file = uploadFile.raw || uploadFile

    console.log('文件对象:', file)
    console.log('文件类型:', file.type)
    console.log('文件名:', file.name)

    const isValidType =
      file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' ||
      file.type === 'text/csv' ||
      file.type === 'application/csv' ||
      file.name.endsWith('.xlsx') ||
      file.name.endsWith('.csv')

    if (!isValidType) {
      ElMessage.error('只支持Excel(.xlsx)或CSV(.csv)格式文件')
      return
    }

    const isLt10M = file.size / 1024 / 1024 < 10
    if (!isLt10M) {
      ElMessage.error('文件大小不能超过10MB')
      return
    }

    selectedFile.value = file
    console.log('文件已选择:', file.name, '文件大小:', file.size, '字节')
  }

  const handleNext = async () => {
    if (!selectedFile.value) {
      console.error('没有选择文件')
      return
    }

    console.log('开始上传文件:', selectedFile.value.name)
    console.log('文件对象详情:', selectedFile.value)

    loading.value = true
    try {
      const result = await batchImportHosts(selectedFile.value)
      // API 返回的是 { data: BatchOperationResult } 结构
      console.log('API响应:', result)

      // 确保正确提取数据
      if (result && typeof result === 'object') {
        // 如果result有data字段，使用data；否则直接使用result
        const data = result.data || result
        importResult.value = {
          total: data.total || 0,
          success: data.success || 0,
          failed: data.failed || 0,
          failed_msg: data.failed_msg || []
        }
      } else {
        importResult.value = {
          total: 0,
          success: 0,
          failed: 0,
          failed_msg: ['API响应格式错误']
        }
      }

      currentStep.value = 1
    } catch (error) {
      ElMessage.error('文件解析失败')
      console.error('Import preview error:', error)
    } finally {
      loading.value = false
    }
  }

  const handlePrevious = () => {
    currentStep.value = Math.max(0, currentStep.value - 1)
  }

  const handleImport = async () => {
    if (!selectedFile.value) return

    loading.value = true
    try {
      const result = await batchImportHosts(selectedFile.value)
      // API 返回的是 { data: BatchOperationResult } 结构
      importResult.value = result
      currentStep.value = 2
      ElMessage.success(`导入完成，成功${result.success}条，失败${result.failed}条`)
      emit('success')
    } catch (error) {
      ElMessage.error('导入失败')
      console.error('Import error:', error)
    } finally {
      loading.value = false
    }
  }

  const handleFinish = () => {
    dialogVisible.value = false
    emit('success')
  }

  const handleCancel = () => {
    dialogVisible.value = false
    resetModal()
  }

  const resetModal = () => {
    currentStep.value = 0
    selectedFile.value = null
    importResult.value = null
    loading.value = false
  }

  const downloadTemplate = () => {
    // 创建模板数据 - 使用中文字段名，便于用户理解
    const templateData = [
      [
        '云厂商ID',
        '实例ID',
        '主机名称',
        '主机组',
        '资源类型',
        '地域',
        '用户名',
        '密码',
        '公网IP',
        '私网IP',
        '配置信息',
        '操作系统',
        '状态',
        '过期时间',
        '备注'
      ],
      [
        '1',
        'i-example123',
        'web-server-01',
        '生产环境',
        'ecs',
        'ap-southeast-1a',
        'root',
        'your_password',
        '["13.229.230.216"]',
        '["172.31.23.13"]',
        '{"cpu_cores":2,"memory_size":4}',
        'Amazon Linux',
        'running',
        '2024-12-31 23:59:59',
        '生产环境Web服务器'
      ],
      [
        '1',
        'i-example456',
        'db-server-01',
        '数据库组',
        'ecs',
        'ap-southeast-1b',
        'admin',
        'your_password',
        '["13.212.189.93"]',
        '["172.31.38.63"]',
        '{"cpu_cores":4,"memory_size":8}',
        'Ubuntu 20.04',
        'stopped',
        '2024-12-31 23:59:59',
        '数据库服务器'
      ]
    ]

    // 创建CSV内容，使用UTF-8 BOM确保中文正确显示
    // 正确处理包含逗号的字段，用双引号包围
    const csvContent = '\uFEFF' + templateData.map((row) =>
      row.map(field => {
        // 如果字段包含逗号、双引号或换行符，需要用双引号包围
        if (field.includes(',') || field.includes('"') || field.includes('\n')) {
          // 转义内部的双引号
          return '"' + field.replace(/"/g, '""') + '"'
        }
        return field
      }).join(',')
    ).join('\n')

    // 创建下载链接
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = '主机导入模板.csv'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    ElMessage.success('模板下载完成')
  }

  const open = () => {
    dialogVisible.value = true
    resetModal()
  }

  defineExpose({
    open
  })
</script>

<style scoped>
  .import-container {
    .template-section {
      h4 {
        margin-bottom: 8px;
        color: #333;
      }

      p {
        margin-bottom: 12px;
        color: #666;
      }
    }
  }
</style>
