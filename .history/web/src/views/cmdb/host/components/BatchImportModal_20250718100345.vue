<template>
  <el-dialog
    v-model="dialogVisible"
    title="批量导入主机"
    width="600px"
    :close-on-click-modal="false"
  >
    <div class="import-container">
      <el-steps :active="currentStep" simple finish-status="success" style="margin-bottom: 24px">
        <el-step title="选择文件" />
        <el-step title="导入预览" />
        <el-step title="导入完成" />
      </el-steps>

      <!-- 步骤 0 -->
      <div v-if="currentStep === 0">
        <el-alert
          title="导入说明"
          type="info"
          show-icon
          :closable="false"
          description="请选择 Excel(.xlsx) 或 CSV(.csv) 文件进行导入"
          class="mb-3"
        />

        <el-upload
          class="upload-area"
          drag
          :auto-upload="false"
          :before-upload="handleFileSelect"
          :show-file-list="false"
          accept=".xlsx,.csv"
        >
          <el-icon class="el-icon--upload"><Upload /></el-icon>
          <div class="el-upload__text">点击或拖拽文件到此处上传</div>
          <template #tip>
            <div class="el-upload__tip">支持 .xlsx / .csv，大小 &lt; 10MB</div>
          </template>
        </el-upload>

        <el-alert
          v-if="selectedFile"
          :title="`已选择文件: ${selectedFile.name}`"
          type="success"
          show-icon
          closable
          class="mt-2"
          @close="selectedFile = null"
        />

        <el-divider />

        <div class="template-section">
          <h4>模板下载</h4>
          <p class="text-secondary">如无模板，请先下载并按格式填写数据。</p>
          <el-button type="primary" link @click="downloadTemplate">
            <el-icon><Download /></el-icon>
            下载 Excel 模板
          </el-button>
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

        <div v-if="importResult && importResult.failed_msg.length">
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
  import * as hostApi from '@/api/system/host'
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

  const handleFileSelect = (file: File): boolean => {
    const isValidType =
      file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' ||
      file.type === 'text/csv' ||
      file.name.endsWith('.xlsx') ||
      file.name.endsWith('.csv')

    if (!isValidType) {
      ElMessage.error('只支持Excel(.xlsx)或CSV(.csv)格式文件')
      return false
    }

    const isLt10M = file.size / 1024 / 1024 < 10
    if (!isLt10M) {
      ElMessage.error('文件大小不能超过10MB')
      return false
    }

    selectedFile.value = file
    return false // 阻止默认上传行为
  }

  const handleNext = async () => {
    if (!selectedFile.value) return

    loading.value = true
    try {
      const result = await hostApi.batchImportHosts(selectedFile.value)
      // API 返回的是 { data: BatchOperationResult } 结构
      importResult.value = result
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
      const result = await hostApi.batchImportHosts(selectedFile.value)
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
    // 创建模板数据
    const templateData = [
      [
        'provider_id',
        'instance_id',
        'name',
        'resource_type',
        'region',
        'username',
        'password',
        'public_ip',
        'private_ip',
        'configuration',
        'os',
        'status',
        'expired_at',
        'remark'
      ],
      [
        '1',
        'i-example123',
        '示例主机',
        'ecs',
        'cn-hangzhou',
        'root',
        'password123',
        '["1.2.3.4"]',
        '["192.168.1.10"]',
        '{"cpu_cores":2,"memory_size":4}',
        'CentOS 7',
        'running',
        '2024-12-31 23:59:59',
        '示例备注'
      ]
    ]

    // 创建CSV内容
    const csvContent = templateData.map((row) => row.join(',')).join('\n')

    // 创建下载链接
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = 'host_import_template.csv'
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
