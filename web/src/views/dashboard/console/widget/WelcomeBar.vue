<template>
  <div class="welcome-bar art-custom-card">
    <div class="welcome-info">
      <img class="avatar" src="/src/assets/img/user/avatar.png" alt="avatar" />
      <div class="text">
        <div class="greet">Hi，欢迎回来！</div>
        <div class="username">超级管理员</div>
      </div>
    </div>
    <div class="weather">
      <i class="iconfont-sys">&#xe6c6;</i>
      <span>上海</span>
      <span class="weather-desc">多云 28℃</span>
    </div>
  </div>
</template>

<script setup lang="ts">
// 后续可接入真实天气API
</script>

<style lang="scss" scoped>
.welcome-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 24px;
  min-width: 0;
  .welcome-info {
    display: flex;
    align-items: center;
    .avatar {
      width: 48px;
      height: 48px;
      border-radius: 50%;
      margin-right: 16px;
      object-fit: cover;
      border: 2px solid var(--el-color-primary-light-7);
    }
    .text {
      .greet {
        font-size: 16px;
        color: var(--el-color-primary);
        font-weight: 600;
      }
      .username {
        font-size: 15px;
        color: var(--art-text-gray-900);
        margin-top: 2px;
      }
    }
  }
  .weather {
    display: flex;
    align-items: center;
    font-size: 15px;
    color: var(--el-color-primary);
    .iconfont-sys {
      font-size: 20px;
      margin-right: 6px;
    }
    .weather-desc {
      color: var(--art-text-gray-600);
      margin-left: 10px;
      font-size: 14px;
    }
  }
}
</style> 