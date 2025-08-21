<template>
  <div class="dns-cert-detail-container">
    <div class="page-header">
      <el-button type="default" @click="goBack" :icon="ArrowLeft" text>返回</el-button>
      <h2 style="margin: 12px 0 0 0">证书详情</h2>
    </div>

    <el-skeleton :rows="6" animated v-if="loading" />

    <div v-else>
      <CertificateDetail
        v-if="certificate"
        :certificate="certificate"
        @download="handleDownload"
        @renew="handleRenew"
        @revoke="handleRevoke"
      />
      <el-empty v-else description="未找到证书" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import CertificateDetail from '@/views/dns/certs/components/CertificateDetail.vue'
import { certificateApi } from '@/api/dns/certificate'
import type { Certificate } from '@/types/dns'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const certificate = ref<Certificate | null>(null)

const id = route.params.id as string

const fetchDetail = async () => {
  loading.value = true
  try {
    const res = await certificateApi.get(Number(id))
    certificate.value = res?.data as Certificate
  } catch (err) {
    console.error('加载证书详情失败:', err)
    ElMessage.error('加载证书详情失败')
  } finally {
    loading.value = false
  }
}

onMounted(fetchDetail)

const goBack = () => router.back()

const handleDownload = async () => {
  if (!id) return
  try {
    const blob = await certificateApi.download(Number(id), 'pem')
    const url = window.URL.createObjectURL(blob as Blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${certificate.value?.common_name || 'certificate'}.pem`
    a.click()
    window.URL.revokeObjectURL(url)
  } catch (e) {
    ElMessage.error('下载失败')
  }
}

const handleRenew = async () => {
  if (!id) return
  try {
    await ElMessageBox.confirm('确定要续期该证书吗？', '确认续期', { type: 'warning' })
    await certificateApi.renew(Number(id))
    ElMessage.success('续期请求已提交')
    await fetchDetail()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('续期失败')
  }
}

const handleRevoke = async () => {
  if (!id) return
  try {
    await ElMessageBox.confirm('确定要吊销该证书吗？此操作不可逆！', '确认吊销', { type: 'warning' })
    await certificateApi.revoke(Number(id))
    ElMessage.success('证书已吊销')
    await fetchDetail()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('吊销失败')
  }
}
</script>

<style scoped>
.dns-cert-detail-container {
  padding: 16px;
}
.page-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}
</style>

