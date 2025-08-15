<template>
  <el-dialog
    v-model="visibleLocal"
    title="批量添加DNS记录"
    width="920px"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    @close="handleClose"
  >
    <div class="bulk-body">
      <el-alert type="info" :closable="false" title="可直接从Excel/CSV粘贴多行，或手动添加行。支持 A, AAAA, CNAME, MX, TXT, SRV。" style="margin-bottom: 8px" />
      <div class="toolbar">
        <el-button @click="addRow">添加一行</el-button>
        <el-button @click="pasteRows">从剪贴板粘贴</el-button>
        <el-checkbox v-model="syncAfterCreate">创建后立即同步云厂商</el-checkbox>
      </div>
      <el-table :data="rows" size="small" border>
        <el-table-column label="#" width="48" type="index" />
        <el-table-column prop="name" label="主机记录" width="160">
          <template #default="{ row }">
            <el-input v-model="row.name" placeholder="@ 或 www" />
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-select v-model="row.type" style="width:100%">
              <el-option label="A" value="A" />
              <el-option label="AAAA" value="AAAA" />
              <el-option label="CNAME" value="CNAME" />
              <el-option label="MX" value="MX" />
              <el-option label="TXT" value="TXT" />
              <el-option label="SRV" value="SRV" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column prop="value" label="记录值" min-width="220">
          <template #default="{ row }">
            <el-input v-model="row.value" placeholder="1.1.1.1 或 cname.example.com" />
          </template>
        </el-table-column>
        <el-table-column prop="ttl" label="TTL" width="120">
          <template #default="{ row }">
            <el-input-number v-model="row.ttl" :min="10" :max="86400" :step="10" />
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="Priority" width="120">
          <template #default="{ row }">
            <el-input-number v-model="row.priority" :min="0" :max="65535" />
          </template>
        </el-table-column>
        <el-table-column prop="weight" label="Weight" width="120">
          <template #default="{ row }">
            <el-input-number v-model="row.weight" :min="0" :max="65535" />
          </template>
        </el-table-column>
        <el-table-column prop="port" label="Port" width="120">
          <template #default="{ row }">
            <el-input-number v-model="row.port" :min="1" :max="65535" />
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150">
          <template #default="{ row }">
            <el-input v-model="row.remark" placeholder="可选" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="90">
          <template #default="{ $index }">
            <el-button size="small" type="danger" @click="removeRow($index)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-alert v-if="errorMessage" type="error" :title="errorMessage" show-icon style="margin-top:8px" />
      <el-result v-if="result" :icon="result.failedCount? 'warning':'success'" :title="resultTitle" :sub-title="`共 ${result.total} 条`" style="margin-top:8px" />
    </div>
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="loading" :disabled="rows.length === 0" @click="handleSubmit">提交</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { recordApi } from '@/api/dns/record'

const props = defineProps<{ visible: boolean; domainId: number | null; providerId?: number | null }>()
const emit = defineEmits(['update:visible','success'])

const visibleLocal = ref(props.visible)
watch(() => props.visible, v => visibleLocal.value = v)
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
    result.value = { successCount: (resp as any).success, failedCount: (resp as any).failed, total: (resp as any).total }
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
</script>

<style scoped>
.bulk-body { display:flex; flex-direction:column; gap:8px; }
.toolbar { display:flex; gap:8px; align-items:center; }
</style>

