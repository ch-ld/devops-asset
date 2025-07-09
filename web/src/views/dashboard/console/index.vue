<template>
  <div class="console">
    <!-- 滚动公告栏 - 全宽显示 -->
    <div class="announcement-section">
      <AnnouncementMarquee />
    </div>
    
    <!-- 欢迎卡片和统计数据 -->
    <div class="welcome-stats-row">
      <div class="welcome-section">
        <WelcomeBar />
      </div>
      <div class="stats-section">
        <ServerStat />
        <UserStat />
      </div>
    </div>
    
    <!-- 便捷导航区域 -->
    <div class="nav-section">
      <QuickNav />
    </div>
    
    <!-- 小工具区域 -->
    <div class="widgets-row">
      <QuickActionsWidget />
      <SystemStatusWidget />
    </div>
  </div>
</template>

<script setup lang="ts">
import AnnouncementMarquee from '@/components/widgets/AnnouncementMarquee.vue'
import WelcomeBar from './widget/WelcomeBar.vue'
import ServerStat from './widget/ServerStat.vue'
import UserStat from './widget/UserStat.vue'
import QuickNav from './widget/QuickNav.vue'
import QuickActionsWidget from '@/components/widgets/QuickActionsWidget.vue'
import SystemStatusWidget from '@/components/widgets/SystemStatusWidget.vue'

defineOptions({ name: 'ConsoleDashboard' })
</script>

<style lang="scss" scoped>
.console {
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  min-height: 100vh;
  animation: pageSlideIn 0.3s ease-out;
  
  // 滚动公告栏区域
  .announcement-section {
    margin-bottom: 28px;
    position: relative;
    z-index: 1; // 设置较低的层级
  }
  
  // 欢迎和统计区域
  .welcome-stats-row {
    display: flex;
    gap: 24px;
    margin-bottom: 68px; // 增加与便捷导航的间距
    
    .welcome-section {
      flex: 1;
      min-width: 0;
      position: relative;
      z-index: 1; // 设置较低的层级
      
      :deep(.art-custom-card) {
        background: rgba(255, 255, 255, 0.95);
        backdrop-filter: blur(10px);
        border: 1px solid rgba(255, 255, 255, 0.2);
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
        border-radius: 16px;
        height: 100%;
      }
    }
    
    .stats-section {
      display: flex;
      gap: 24px;
      flex: 2;
      position: relative;
      z-index: 1; // 设置较低的层级
      
      > * {
        flex: 1;
        min-width: 0;
      }
      
      :deep(.art-custom-card) {
        background: rgba(255, 255, 255, 0.95);
        backdrop-filter: blur(10px);
        border: 1px solid rgba(255, 255, 255, 0.2);
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
        border-radius: 16px;
        height: 100%;
      }
    }
    
    @media (max-width: 1200px) {
      flex-direction: column;
      
      .stats-section {
        flex-direction: row;
      }
    }
    
    @media (max-width: 768px) {
      .stats-section {
        flex-direction: column;
      }
    }
  }
  
  // 小工具区域
  .widgets-row {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(380px, 1fr));
    gap: 24px;
    margin-bottom: 28px;
    position: relative;
    z-index: 1; // 设置较低的层级
    
    @media (max-width: 768px) {
      grid-template-columns: 1fr;
    }
  }
  
  // 导航区域
  .nav-section {
    margin-bottom: 28px;
    position: relative;
    z-index: 10; // 确保导航在最上层
    clear: both; // 清除浮动
    
    :deep(.art-custom-card) {
      background: rgba(255, 255, 255, 0.95);
      backdrop-filter: blur(10px);
      border: 1px solid rgba(255, 255, 255, 0.2);
      box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
      border-radius: 16px;
    }
    
          :deep(.quick-nav-ant) {
        position: relative;
        z-index: 20;
      }
  }
}

// 全局美化小工具卡片
:deep(.quick-actions-widget),
:deep(.system-status-widget) {
  background: rgba(255, 255, 255, 0.95) !important;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2) !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1) !important;
}
</style>
