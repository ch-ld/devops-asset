<template>
  <div class="p-4">
    <el-card shadow="never">
      <template #header>
        <div class="flex justify-between items-center">
          <span>厂商管理</span>
          <el-button type="primary" @click="handleAdd">添加厂商</el-button>
        </div>
      </template>

      <div class="flex justify-between mb-3">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索厂商名称..."
          style="width: 240px"
          clearable
          @keyup.enter="onSearch"
        >
          <template #append>
            <el-button icon="Search" @click="onSearch" />
          </template>
        </el-input>
        <el-button @click="fetchProviders" :loading="loading">刷新</el-button>
      </div>

      <el-table :data="providerList" border stripe style="width: 100%" :loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="厂商名称" />
        <el-table-column prop="type" label="类型" width="120" align="center">
          <template #default="{ row }">
            <el-tag>{{ row.type.toUpperCase() }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status ? 'success' : 'danger'">{{
              row.status ? '有效' : '无效'
            }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="220" align="center">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              size="small"
              @click="handleSync(row.id)"
              :loading="syncing[row.id]"
              >同步</el-button
            >
            <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除此厂商吗？" @confirm="handleDelete(row.id)">
              <el-button link type="danger" size="small">删除</el-button>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="mt-3 text-right">
        <el-pagination
          background
          layout="prev, pager, next"
          :page-size="pagination.pageSize"
          :current-page="pagination.current"
          :total="pagination.total"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <ProviderModal ref="providerModalRef" @success="fetchProviders" />
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import ProviderModal from './ProviderModal.vue'
  import * as hostApi from '@/api/system/host'
  import type { Provider } from '@/types/api/host'

  const providerList = ref<Provider[]>([])
  const loading = ref(false)
  const searchKeyword = ref('')
  const providerModalRef = ref<InstanceType<typeof ProviderModal>>()
  const syncing = reactive<Record<number, boolean>>({})

  const pagination = reactive({
    current: 1,
    pageSize: 10,
    total: 0
  })

  const fetchProviders = async () => {
    loading.value = true
    try {
      const res = await hostApi.getProviderList()
      // 暂无分页接口，简单前端分页
      const filtered = res.filter((p: Provider) => p.name.includes(searchKeyword.value))
      pagination.total = filtered.length
      const start = (pagination.current - 1) * pagination.pageSize
      providerList.value = filtered.slice(start, start + pagination.pageSize)
    } catch (err) {
      ElMessage.error('获取厂商失败')
    } finally {
      loading.value = false
    }
  }

  const handlePageChange = (page: number) => {
    pagination.current = page
    fetchProviders()
  }

  const onSearch = () => {
    pagination.current = 1
    fetchProviders()
  }

  const handleAdd = () => {
    providerModalRef.value?.open()
  }

  const handleEdit = (record: Provider) => {
    providerModalRef.value?.open(record)
  }

  const handleDelete = async (id: number) => {
    try {
      await hostApi.deleteProvider(id)
      ElMessage.success('删除成功')
      fetchProviders()
    } catch (err) {
      ElMessage.error('删除失败')
    }
  }

  const handleSync = async (id: number) => {
    syncing[id] = true
    try {
      await hostApi.syncResources(id)
      ElMessage.success('同步指令已发送')
    } catch (err) {
      ElMessage.error('同步失败')
    } finally {
      syncing[id] = false
    }
  }

  const formatTime = (ts: string) => {
    if (!ts) return '--'
    return new Date(ts).toLocaleString()
  }

  onMounted(fetchProviders)
</script>

<style scoped></style>
