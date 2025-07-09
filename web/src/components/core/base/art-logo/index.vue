<!-- 系统logo -->
<template>
  <div class="art-logo">
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
      <defs>
        <linearGradient id="devops-blue" x1="0" y1="0" x2="40" y2="40" gradientUnits="userSpaceOnUse">
          <stop stop-color="#5D87FF"/>
          <stop offset="1" stop-color="#38C0FC"/>
        </linearGradient>
      </defs>
      <polygon points="20,4 36,36 4,36" fill="url(#devops-blue)" opacity="0.9"/>
      <ellipse cx="20" cy="28" rx="10" ry="4" fill="#fff" opacity="0.2"/>
      <path d="M10 30 Q20 20 30 30" stroke="#fff" stroke-width="2" fill="none" opacity="0.7"/>
    </svg>
  </div>
</template>

<script setup lang="ts">
  defineOptions({ name: 'ArtLogo' })
</script>

<style lang="scss" scoped>
  .art-logo {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
  }
</style>
