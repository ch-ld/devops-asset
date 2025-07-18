<template>
  <div class="host-dashboard-container">
    <el-card shadow="never" class="dashboard-header">
      <template #header>
        <div class="card-header">
          <h3>主机概览仪表盘</h3>
        </div>
      </template>
      <div class="dashboard-content">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card class="stat-card">
              <template #header>
                <div class="card-header">
                  <span>主机总数</span>
                </div>
              </template>
              <div class="stat-value">{{ hostCount }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stat-card">
              <template #header>
                <div class="card-header">
                  <span>在线主机</span>
                </div>
              </template>
              <div class="stat-value">{{ onlineCount }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stat-card">
              <template #header>
                <div class="card-header">
                  <span>异常主机</span>
                </div>
              </template>
              <div class="stat-value">{{ errorCount }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stat-card">
              <template #header>
                <div class="card-header">
                  <span>即将到期</span>
                </div>
              </template>
              <div class="stat-value">{{ expiringCount }}</div>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-card>

    <el-card shadow="never" class="chart-card">
      <template #header>
        <div class="card-header">
          <h3>主机状态分布</h3>
        </div>
      </template>
      <div class="chart-container">
        <div style="height: 300px; display: flex; align-items: center; justify-content: center">
          <el-empty description="暂无数据" />
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
  import { ref, onMounted } from 'vue'
  import { useHostStore } from '@/store/modules/host'
  import { storeToRefs } from 'pinia'

  // 使用host store
  const hostStore = useHostStore()
  const { hostList } = storeToRefs(hostStore)

  // 统计数据
  const hostCount = ref(0)
  const onlineCount = ref(0)
  const errorCount = ref(0)
  const expiringCount = ref(0)

  // 初始化数据
  onMounted(async () => {
    try {
      // 加载主机数据
      await hostStore.fetchHosts()

      // 更新统计数据
      hostCount.value = hostList.value.length
      onlineCount.value = hostList.value.filter((h) => h.status === 'running').length
      errorCount.value = hostList.value.filter((h) => h.status === 'error').length

      // 计算30天内到期的主机数量
      const now = new Date()
      const thirtyDaysLater = new Date()
      thirtyDaysLater.setDate(now.getDate() + 30)

      expiringCount.value = hostList.value.filter((h) => {
        if (!h.expired_at) return false
        const expireDate = new Date(h.expired_at)
        return expireDate > now && expireDate < thirtyDaysLater
      }).length
    } catch (error) {
      console.error('加载主机数据失败:', error)
    }
  })
</script>

<style scoped>
  .host-dashboard-container {
    padding: 20px;
  }

  .dashboard-header {
    margin-bottom: 20px;
  }

  .chart-card {
    margin-bottom: 20px;
  }

  .stat-card {
    text-align: center;
  }

  .stat-value {
    font-size: 24px;
    font-weight: bold;
    color: #409eff;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .chart-container {
    height: 400px;
  }
</style>
