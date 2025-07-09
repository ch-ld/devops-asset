<template>
  <div class="server-stat art-custom-card">
    <div class="stat-title">
      <i class="anticon anticon-cloud-server"></i>
      <span>服务器统计</span>
    </div>
    <div class="stat-list">
      <div class="stat-item">
        <i class="anticon anticon-database"></i>
        <span class="stat-label">服务器总数</span>
        <span class="stat-value">12</span>
      </div>
      <div class="stat-item">
        <i class="anticon anticon-cloud"></i>
        <span class="stat-label">在线服务器</span>
        <span class="stat-value">10</span>
      </div>
      <div class="stat-item">
        <i class="anticon anticon-disconnect"></i>
        <span class="stat-label">离线服务器</span>
        <span class="stat-value">2</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 预留后期动态数据
</script>

<style lang="scss" scoped>
.server-stat {
  padding: 18px 24px;
  margin-bottom: 18px;
  .stat-title {
    display: flex;
    align-items: center;
    font-size: 18px;
    font-weight: 600;
    color: var(--el-color-primary);
    margin-bottom: 10px;
    .anticon {
      font-size: 22px;
      margin-right: 8px;
      color: #1677ff;
    }
  }
  .stat-list {
    display: flex;
    gap: 32px;
    .stat-item {
      display: flex;
      flex-direction: column;
      align-items: flex-start;
      position: relative;
      .anticon {
        font-size: 20px;
        margin-bottom: 2px;
        color: #69b1ff;
      }
      .stat-label {
        font-size: 14px;
        color: var(--art-text-gray-600);
      }
      .stat-value {
        font-size: 22px;
        font-weight: 600;
        color: var(--art-text-gray-900);
        margin-top: 2px;
      }
    }
  }
}
// Ant Design 图标
.anticon {
  display: inline-block;
  font-style: normal;
  vertical-align: middle;
}
.anticon-cloud-server:before { content: '\1F5A5'; }
.anticon-database:before { content: '\1F4BE'; }
.anticon-cloud:before { content: '\2601'; }
.anticon-disconnect:before { content: '\26D4'; }
</style> 