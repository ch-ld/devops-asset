<template>
  <a-modal
    title="批量导入主机"
    :open="visible"
    :width="600"
    :confirm-loading="loading"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <div class="import-container">
      <a-steps :current="currentStep" size="small" style="margin-bottom: 24px">
        <a-step title="选择文件" />
        <a-step title="导入预览" />
        <a-step title="导入完成" />
      </a-steps>

      <!-- 第一步：选择文件 -->
      <div v-if="currentStep === 0">
        <a-alert
          message="导入说明"
          description="请选择Excel或CSV格式的主机文件进行批量导入。系统支持.xlsx和.csv格式文件。"
          type="info"
          show-icon
          style="margin-bottom: 16px"
        />

        <a-upload-dragger
          :before-upload="handleFileSelect"
          :show-upload-list="false"
          accept=".xlsx,.csv"
        >
          <p class="ant-upload-drag-icon">
            <InboxOutlined />
          </p>
          <p class="ant-upload-text">点击或拖拽文件到此区域上传</p>
          <p class="ant-upload-hint">支持.xlsx和.csv格式文件</p>
        </a-upload-dragger>

        <div v-if="selectedFile" style="margin-top: 16px">
          <a-alert
            :message="`已选择文件: ${selectedFile.name}`"
            type="success"
            show-icon
            closable
            @close="selectedFile = null"
          />
        </div>

        <a-divider />

        <div class="template-section">
          <h4>模板下载</h4>
          <p>如果您还没有导入模板，请先下载模板文件并按照格式填写数据。</p>
          <a-button type="link" @click="downloadTemplate">
            <template #icon><DownloadOutlined /></template>
            下载Excel模板
          </a-button>
        </div>
      </div>

      <!-- 第二步：导入预览 -->
      <div v-if="currentStep === 1">
        <a-alert
          v-if="importResult"
          :message="`解析完成，共${importResult.total}条记录，成功${importResult.success}条，失败${importResult.failed}条`"
          :type="importResult.failed > 0 ? 'warning' : 'success'"
          show-icon
          style="margin-bottom: 16px"
        />

        <div v-if="importResult && importResult.failed_msg.length > 0">
          <h4>错误信息：</h4>
          <a-list size="small" bordered>
            <a-list-item v-for="(error, index) in importResult.failed_msg" :key="index">
              <a-typography-text type="danger">{{ error }}</a-typography-text>
            </a-list-item>
          </a-list>
        </div>
      </div>

      <!-- 第三步：导入完成 -->
      <div v-if="currentStep === 2">
        <a-result
          status="success"
          title="导入完成"
          :sub-title="`成功导入${importResult?.success || 0}条主机记录`"
        />
      </div>
    </div>

    <template #footer>
      <a-space>
        <a-button v-if="currentStep > 0" @click="handlePrevious">上一步</a-button>
        <a-button @click="handleCancel">取消</a-button>
        <a-button
          v-if="currentStep === 0"
          type="primary"
          :disabled="!selectedFile"
          @click="handleNext"
        >
          下一步
        </a-button>
        <a-button v-if="currentStep === 1" type="primary" :loading="loading" @click="handleImport">
          确认导入
        </a-button>
        <a-button v-if="currentStep === 2" type="primary" @click="handleFinish"> 完成 </a-button>
      </a-space>
    </template>
  </a-modal>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import { Upload as InboxOutlined, Download as DownloadOutlined } from '@element-plus/icons-vue'
  import * as hostApi from '@/api/system/host'
  import type { BatchOperationResult } from '@/types/api/host'

  const emit = defineEmits<{
    success: []
  }>()

  const visible = ref(false)
  const loading = ref(false)
  const currentStep = ref(0)
  const selectedFile = ref<File | null>(null)
  const importResult = ref<BatchOperationResult | null>(null)

  const handleFileSelect = (file: File) => {
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
      importResult.value = result.data
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
      importResult.value = result.data
      currentStep.value = 2
      ElMessage.success('导入完成')
    } catch (error) {
      ElMessage.error('导入失败')
      console.error('Import error:', error)
    } finally {
      loading.value = false
    }
  }

  const handleFinish = () => {
    visible.value = false
    emit('success')
  }

  const handleCancel = () => {
    visible.value = false
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
    visible.value = true
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
