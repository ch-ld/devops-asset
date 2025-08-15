<template>
  <div class="dns-page-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title">
          <h1>{{ title }}</h1>
          <p v-if="description">{{ description }}</p>
        </div>
        <div class="page-actions">
          <slot name="header-actions" />
        </div>
      </div>
    </div>

    <!-- 统计卡片区域 -->
    <div v-if="showStats" class="stats-container">
      <slot name="stats" />
    </div>

    <!-- 搜索表单区域 -->
    <div v-if="showSearch" class="search-container">
      <slot name="search" />
    </div>

    <!-- 主要内容区域 -->
    <div class="content-container">
      <slot />
    </div>

    <!-- 浮动操作区域 -->
    <div v-if="$slots['floating-actions']" class="floating-actions">
      <slot name="floating-actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
export interface Props {
  /** 页面标题 */
  title: string
  /** 页面描述 */
  description?: string
  /** 是否显示统计卡片区域 */
  showStats?: boolean
  /** 是否显示搜索区域 */
  showSearch?: boolean
}

withDefaults(defineProps<Props>(), {
  showStats: true,
  showSearch: true
})
</script>

<style scoped>
.dns-page-container {
  padding: 0;
}

.page-header {
  margin-bottom: 24px;
  
  .page-header-content {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    
    .page-title {
      h1 {
        margin: 0 0 8px 0;
        font-size: 24px;
        font-weight: 600;
        color: #1f2937;
        line-height: 1.2;
      }
      
      p {
        margin: 0;
        color: #6b7280;
        font-size: 14px;
        line-height: 1.4;
      }
    }
    
    .page-actions {
      display: flex;
      gap: 8px;
      align-items: center;
      flex-shrink: 0;
    }
  }
}

.stats-container {
  margin-bottom: 24px;
}

.search-container {
  margin-bottom: 24px;
}

.content-container {
  position: relative;
}

.floating-actions {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

@media (max-width: 768px) {
  .page-header .page-header-content {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
    
    .page-actions {
      justify-content: flex-end;
    }
  }
  
  .floating-actions {
    bottom: 16px;
    right: 16px;
  }
}
</style>
