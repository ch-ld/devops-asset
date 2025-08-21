<template>
  <div class="domain-detail">
    <el-page-header @back="goBack" :content="domainInfo.name">
      <template #extra>
        <el-space>
          <el-button type="primary" @click="editDomain">
            <el-icon><Edit /></el-icon>
            编辑域名
          </el-button>
          <el-button type="success" @click="addCertificate">
            <el-icon><Plus /></el-icon>
            申请证书
          </el-button>
          <el-dropdown>
            <el-button>
              更多操作<el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="exportRecords">导出DNS记录</el-dropdown-item>
                <el-dropdown-item @click="syncRecords">同步DNS记录</el-dropdown-item>
                <el-dropdown-item divided @click="deleteDomain">删除域名</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </el-space>
      </template>
    </el-page-header>

    <el-row :gutter="24" class="mt-4">
      <!-- 左侧：域名基本信息 -->
      <el-col :span="8">
        <el-card title="基本信息" class="mb-4">
          <template #header>
            <div class="card-header">
              <span>基本信息</span>
              <el-button class="button" text @click="editDomain">编辑</el-button>
            </div>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="域名">{{ domainInfo.name }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(domainInfo.status)">{{ domainInfo.status }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="分组">{{ domainInfo.groupName || '-' }}</el-descriptions-item>
            <el-descriptions-item label="标签">
              <el-space wrap>
                <el-tag v-for="tag in domainInfo.tags" :key="tag.id" :color="tag.color">
                  {{ tag.name }}
                </el-tag>
                <span v-if="!domainInfo.tags || domainInfo.tags.length === 0">-</span>
              </el-space>
            </el-descriptions-item>
            <el-descriptions-item label="到期时间">
              <el-text :type="getExpiryType(domainInfo.expiresAt)">
                {{ formatDate(domainInfo.expiresAt) }}
              </el-text>
            </el-descriptions-item>
            <el-descriptions-item label="注册商">{{ domainInfo.registrar || '-' }}</el-descriptions-item>
            <el-descriptions-item label="DNS提供商">{{ domainInfo.providerName || '-' }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDate(domainInfo.createdAt) }}</el-descriptions-item>
            <el-descriptions-item label="更新时间">{{ formatDate(domainInfo.updatedAt) }}</el-descriptions-item>
            <el-descriptions-item label="备注">{{ domainInfo.remark || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>

        <!-- WHOIS信息 -->
        <el-card title="WHOIS信息">
          <template #header>
            <div class="card-header">
              <span>WHOIS信息</span>
              <el-button class="button" text @click="refreshWhois">刷新</el-button>
            </div>
          </template>
          <el-descriptions :column="1" border v-if="whoisInfo">
            <el-descriptions-item label="注册商">{{ whoisInfo.registrar }}</el-descriptions-item>
            <el-descriptions-item label="注册时间">{{ formatDate(whoisInfo.registeredAt) }}</el-descriptions-item>
            <el-descriptions-item label="到期时间">{{ formatDate(whoisInfo.expiresAt) }}</el-descriptions-item>
            <el-descriptions-item label="名称服务器">
              <div v-for="ns in whoisInfo.nameServers" :key="ns">{{ ns }}</div>
            </el-descriptions-item>
          </el-descriptions>
          <el-empty v-else description="暂无WHOIS信息" />
        </el-card>
      </el-col>

      <!-- 右侧：DNS记录和证书 -->
      <el-col :span="16">
        <el-tabs v-model="activeTab" class="domain-tabs">
          <!-- DNS记录标签页 -->
          <el-tab-pane label="DNS记录" name="records">
            <div class="tab-content">
              <div class="mb-4">
                <el-space>
                  <el-button type="primary" @click="addRecord">
                    <el-icon><Plus /></el-icon>
                    添加记录
                  </el-button>
                  <el-button @click="syncRecords">
                    <el-icon><Refresh /></el-icon>
                    同步记录
                  </el-button>
                  <el-button @click="exportRecords">
                    <el-icon><Download /></el-icon>
                    导出记录
                  </el-button>
                </el-space>
              </div>
              <el-table :data="dnsRecords" v-loading="recordsLoading">
                <el-table-column prop="name" label="名称" />
                <el-table-column prop="type" label="类型" width="80" />
                <el-table-column prop="value" label="值" show-overflow-tooltip />
                <el-table-column prop="ttl" label="TTL" width="80" />
                <el-table-column prop="priority" label="优先级" width="80" />
                <el-table-column label="状态" width="100">
                  <template #default="{ row }">
                    <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
                      {{ row.status }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="150">
                  <template #default="{ row }">
                    <el-button text @click="editRecord(row)">编辑</el-button>
                    <el-button text type="danger" @click="deleteRecord(row)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-tab-pane>

          <!-- 证书标签页 -->
          <el-tab-pane label="SSL证书" name="certificates">
            <div class="tab-content">
              <div class="mb-4">
                <el-space>
                  <el-button type="primary" @click="addCertificate">
                    <el-icon><Plus /></el-icon>
                    申请证书
                  </el-button>
                  <el-button @click="refreshCertificates">
                    <el-icon><Refresh /></el-icon>
                    刷新状态
                  </el-button>
                </el-space>
              </div>
              <el-table :data="certificates" v-loading="certificatesLoading">
                <el-table-column prop="commonName" label="证书名称" />
                <el-table-column prop="caType" label="CA类型" width="120" />
                <el-table-column label="状态" width="100">
                  <template #default="{ row }">
                    <el-tag :type="getCertStatusType(row.status)">{{ row.status }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="issuedAt" label="签发时间" width="120">
                  <template #default="{ row }">
                    {{ formatDate(row.issuedAt) }}
                  </template>
                </el-table-column>
                <el-table-column prop="expiresAt" label="到期时间" width="120">
                  <template #default="{ row }">
                    <el-text :type="getExpiryType(row.expiresAt)">
                      {{ formatDate(row.expiresAt) }}
                    </el-text>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="200">
                  <template #default="{ row }">
                    <el-button text @click="viewCertificate(row)">查看</el-button>
                    <el-button text @click="downloadCertificate(row)">下载</el-button>
                    <el-button text @click="deployCertificate(row)">部署</el-button>
                    <el-button text type="danger" @click="revokeCertificate(row)">吊销</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-tab-pane>

          <!-- 监控历史标签页 -->
          <el-tab-pane label="监控历史" name="monitoring">
            <div class="tab-content">
              <el-table :data="monitoringHistory" v-loading="monitoringLoading">
                <el-table-column prop="url" label="监控URL" />
                <el-table-column label="状态" width="100">
                  <template #default="{ row }">
                    <el-tag :type="row.status === 'up' ? 'success' : 'danger'">
                      {{ row.status }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="responseTime" label="响应时间" width="100">
                  <template #default="{ row }">{{ row.responseTime }}ms</template>
                </el-table-column>
                <el-table-column prop="httpStatus" label="HTTP状态" width="100" />
                <el-table-column prop="certExpiryDays" label="证书剩余天数" width="120">
                  <template #default="{ row }">
                    <el-text :type="row.certExpiryDays < 30 ? 'danger' : 'primary'">
                      {{ row.certExpiryDays }}天
                    </el-text>
                  </template>
                </el-table-column>
                <el-table-column prop="checkedAt" label="检查时间" width="160">
                  <template #default="{ row }">
                    {{ formatDateTime(row.checkedAt) }}
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-tab-pane>
        </el-tabs>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, Plus, ArrowDown, Refresh, Download } from '@element-plus/icons-vue'
import { domainApi } from '@/api/dns/domain'
import { certificateApi } from '@/api/dns/certificate'
import { recordApi } from '@/api/dns/record'

const route = useRoute()
const router = useRouter()

const domainId = ref(route.params.id as string)
const activeTab = ref('records')

// 域名信息
const domainInfo = ref<any>({})
const whoisInfo = ref<any>(null)

// 表格数据
const dnsRecords = ref([])
const certificates = ref([])
const monitoringHistory = ref([])

// 加载状态
const recordsLoading = ref(false)
const certificatesLoading = ref(false)
const monitoringLoading = ref(false)

// 获取域名详情
const fetchDomainDetail = async () => {
  try {
    const response = await domainApi.getDetail(domainId.value)
    domainInfo.value = response.data
  } catch (error) {
    ElMessage.error('获取域名详情失败')
  }
}

// 获取DNS记录
const fetchDnsRecords = async () => {
  recordsLoading.value = true
  try {
    const response = await recordApi.list({ domainId: domainId.value })
    dnsRecords.value = response.data
  } catch (error) {
    ElMessage.error('获取DNS记录失败')
  } finally {
    recordsLoading.value = false
  }
}

// 获取证书列表
const fetchCertificates = async () => {
  certificatesLoading.value = true
  try {
    const response = await certificateApi.list({ domainId: domainId.value })
    certificates.value = response.data
  } catch (error) {
    ElMessage.error('获取证书列表失败')
  } finally {
    certificatesLoading.value = false
  }
}

// 获取监控历史
const fetchMonitoringHistory = async () => {
  monitoringLoading.value = true
  try {
    // TODO: 实现监控历史API
    await new Promise(resolve => setTimeout(resolve, 1000))
    monitoringHistory.value = []
  } catch (error) {
    ElMessage.error('获取监控历史失败')
  } finally {
    monitoringLoading.value = false
  }
}

// 刷新WHOIS信息
const refreshWhois = async () => {
  try {
    // TODO: 实现WHOIS API
    ElMessage.success('WHOIS信息已刷新')
  } catch (error) {
    ElMessage.error('刷新WHOIS信息失败')
  }
}

// 操作方法
const goBack = () => {
  router.back()
}

const editDomain = () => {
  router.push(`/dns/domains/edit/${domainId.value}`)
}

const addRecord = () => {
  // TODO: 打开添加记录弹窗
}

const editRecord = (record: any) => {
  // TODO: 打开编辑记录弹窗
}

const deleteRecord = async (record: any) => {
  try {
    await ElMessageBox.confirm('确定要删除此DNS记录吗？', '确认删除', {
      type: 'warning'
    })
    await recordApi.delete(record.id)
    ElMessage.success('删除成功')
    fetchDnsRecords()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const addCertificate = () => {
  router.push(`/dns/certificates/create?domainId=${domainId.value}`)
}

const viewCertificate = (cert: any) => {
  router.push(`/dns/certs/${cert.id}`)
}

const downloadCertificate = async (cert: any) => {
  try {
    await certificateApi.download(cert.id)
    ElMessage.success('证书下载成功')
  } catch (error) {
    ElMessage.error('证书下载失败')
  }
}

const deployCertificate = (cert: any) => {
  // TODO: 打开部署弹窗
}

const revokeCertificate = async (cert: any) => {
  try {
    await ElMessageBox.confirm('确定要吊销此证书吗？此操作不可逆！', '确认吊销', {
      type: 'warning'
    })
    await certificateApi.revoke(cert.id)
    ElMessage.success('证书已吊销')
    fetchCertificates()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('吊销失败')
    }
  }
}

const syncRecords = async () => {
  try {
    await recordApi.sync(domainId.value)
    ElMessage.success('DNS记录同步成功')
    fetchDnsRecords()
  } catch (error) {
    ElMessage.error('DNS记录同步失败')
  }
}

const exportRecords = async () => {
  try {
    await recordApi.export(domainId.value)
    ElMessage.success('DNS记录导出成功')
  } catch (error) {
    ElMessage.error('DNS记录导出失败')
  }
}

const refreshCertificates = () => {
  fetchCertificates()
}

const deleteDomain = async () => {
  try {
    await ElMessageBox.confirm('确定要删除此域名吗？此操作将同时删除相关的DNS记录和证书！', '确认删除', {
      type: 'warning'
    })
    await domainApi.delete(domainId.value)
    ElMessage.success('域名删除成功')
    router.push('/dns/domains')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('域名删除失败')
    }
  }
}

// 工具函数
const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString()
}

const formatDateTime = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    'active': 'success',
    'expired': 'danger',
    'pending': 'warning',
    'suspended': 'info'
  }
  return statusMap[status] || 'info'
}

const getCertStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    'issued': 'success',
    'pending': 'warning',
    'failed': 'danger',
    'expired': 'danger',
    'revoked': 'info'
  }
  return statusMap[status] || 'info'
}

const getExpiryType = (expiryDate: string) => {
  if (!expiryDate) return 'info'
  const days = Math.ceil((new Date(expiryDate).getTime() - Date.now()) / (1000 * 60 * 60 * 24))
  if (days < 7) return 'danger'
  if (days < 30) return 'warning'
  return 'primary'
}

onMounted(() => {
  fetchDomainDetail()
  fetchDnsRecords()
  fetchCertificates()
  fetchMonitoringHistory()
})
</script>

<style scoped>
.domain-detail {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.domain-tabs {
  min-height: 500px;
}

.tab-content {
  padding: 16px 0;
}

.mt-4 {
  margin-top: 1rem;
}

.mb-4 {
  margin-bottom: 1rem;
}
</style> 
