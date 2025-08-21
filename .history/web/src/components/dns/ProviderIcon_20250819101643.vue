<template>
  <div class="provider-icon" :class="`provider-${type}`">
    <!-- 阿里云 -->
    <svg v-if="type === 'aliyun'" viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
    </svg>
    
    <!-- 腾讯云 -->
    <svg v-else-if="type === 'tencent'" viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
    </svg>
    
    <!-- AWS -->
    <svg v-else-if="type === 'aws'" viewBox="0 0 24 24" fill="currentColor">
      <path d="M9 11H7l1.5-4.5L10 11h-1zm6.5 0H14l1.5-4.5L17 11h-1.5zM12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z"/>
    </svg>
    
    <!-- GoDaddy -->
    <svg v-else-if="type === 'godaddy'" viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm5 9h-4v4h-2v-4H7V9h4V7h2v2h4v2z"/>
    </svg>
    
    <!-- Cloudflare -->
    <svg v-else-if="type === 'cloudflare'" viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
    </svg>
    
    <!-- 默认图标 -->
    <svg v-else viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
    </svg>
  </div>
</template>

<script setup lang="ts">
interface Props {
  type: string
  size?: string
}

withDefaults(defineProps<Props>(), {
  size: '24px'
})
</script>

<style scoped lang="scss">
.provider-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: v-bind(size);
  height: v-bind(size);
  border-radius: 4px;
  
  svg {
    width: 100%;
    height: 100%;
  }
  
  &.provider-aliyun {
    background: linear-gradient(135deg, #ff6a00, #ff8f00);
    color: white;
  }
  
  &.provider-tencent {
    background: linear-gradient(135deg, #006eff, #0084ff);
    color: white;
  }
  
  &.provider-aws {
    background: linear-gradient(135deg, #ff9900, #ffb84d);
    color: white;
  }
  
  &.provider-godaddy {
    background: linear-gradient(135deg, #1baa6e, #2bc77a);
    color: white;
  }
  
  &.provider-cloudflare {
    background: linear-gradient(135deg, #f38020, #f5a623);
    color: white;
  }
  
  // 默认样式
  &:not([class*="provider-"]) {
    background: linear-gradient(135deg, #909399, #b3b7bc);
    color: white;
  }
}
</style>
