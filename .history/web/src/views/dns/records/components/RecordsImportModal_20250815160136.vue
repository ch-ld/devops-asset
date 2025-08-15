<template>
  <el-dialog
    v-model="visibleLocal"
    title="导入DNS解析记录"
    width="920px"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    @close="handleClose"
  >
    <div class="import-body">
      <!-- 步骤指示器 -->
      <el-steps :active="currentStep" finish-status="success" class="import-steps" style="margin-bottom: 20px">
        <el-step title="上传文件" description="选择要导入的文件" />
        <el-step title="数据预览" description="预览并确认导入数据" />
        <el-step title="导入完成" description="查看导入结果" />
      </el-steps>

      <!-- 步骤1: 上传文件 -->
      <div v-if="currentStep === 0" class="step-content">
        <el-alert
          type="info"
          :closable="false"
          style="margin-bottom: 12px"
          title="支持 CSV 与 Excel(xlsx)。建议先下载模板，按模板填写后上传。"
        />

        <div class="tool-row">
          <el-button @click="downloadTemplate('csv')" :loading="downloading.csv">下载CSV模板</el-button>
          <el-button @click="downloadTemplate('excel')" :loading="downloading.excel">下载Excel模板</el-button>
        </div>

        <el-upload
          drag
          :auto-upload="false"
          :on-change="onFileChange"
          :limit="1"
          accept=".csv,.xlsx"
          style="width: 100%"
        >
          <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">
            将文件拖到此处，或 <em>点击上传</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">仅支持 CSV 或 Excel(xlsx)</div>
          </template>
        </el-upload>
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
import { UploadFilled, View, WarningFilled } from '@element-plus/icons-vue'

const props = defineProps<{ visible: boolean; domainId: number | null }>()
const emit = defineEmits(['update:visible', 'success'])

const visibleLocal = ref(props.visible)
const file = ref<File | null>(null)
const loading = ref(false)
const downloading = ref({ csv: false, excel: false })
const result = ref<{ successCount: number; failedCount: number; total: number; failedItems: any[] } | null>(null)
const currentStep = ref(0)
const previewData = ref<any[]>([])

watch(() => props.visible, v => (visibleLocal.value = v))
watch(visibleLocal, v => emit('update:visible', v))

const onFileChange = (f: any) => {
  file.value = f?.raw || null
  result.value = null
}

const resultIcon = computed(() => (result.value && result.value.failedCount === 0 ? 'success' : 'warning'))
const resultTitle = computed(() => {
  if (!result.value) return '等待导入'
  const { successCount, failedCount } = result.value
  return failedCount === 0 ? `导入成功 ${successCount} 条` : `部分成功：成功 ${successCount} 条，失败 ${failedCount} 条`
})
const resultSub = computed(() => (result.value ? `共 ${result.value.total} 条` : ''))

const handleImport = async () => {
  if (!file.value) return
  if (!props.domainId) {
    ElMessage.warning('请先选择域名')
    return
  }
  loading.value = true
  try {
    // 读取文件并解析为 FormData -> 后端接收原始文件或数据列表二选一
    // 为简化先走文件直传通道：前端读取并解析CSV/Excel为 JSON 列表再提交
    const rows = await parseFile(file.value)
    const payload = {
      format: 'csv',
      data: rows.map(r => ({
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
      successCount: (resp as any).success,
      failedCount: (resp as any).failed,
      total: (resp as any).total,
      failedItems: (resp as any).failed_items || []
    }
    ElMessage.success('导入完成')
    // 导入完成后可自动刷新由外部处理，这里抛 success 交给父组件
    emit('success')
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
</script>

<style scoped>
.tool-row { margin-bottom: 12px; display: flex; gap: 8px; }
.result { margin-top: 12px; }
</style>

