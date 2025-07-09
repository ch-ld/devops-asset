<template>
  <div class="user-stat art-custom-card">
    <div class="stat-title">
      <i class="anticon anticon-usergroup"></i>
      <span>用户统计</span>
    </div>
    <div class="stat-list">
      <div class="stat-item">
        <i class="anticon anticon-user"></i>
        <span class="stat-label">用户总数</span>
        <span class="stat-value">1024</span>
      </div>
      <div class="stat-item">
        <i class="anticon anticon-user-active"></i>
        <span class="stat-label">活跃用户</span>
        <span class="stat-value">256</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 预留后期动态数据
</script>

<style lang="scss" scoped>
.user-stat {
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
.anticon-usergroup:before { content: '\1F465'; }
.anticon-user:before { content: '\1F464'; }
.anticon-user-active:before { content: '\1F91D'; }
</style> 