<template>
  <HostDashboard />
</template>

<script setup>
  import HostDashboard from './dashboard/index.vue'
</script>

<style lang="scss" scoped>
  .host-dashboard-container {
    padding: 0;

    .stat-card {
      height: 100%;

      .stat-content {
        display: flex;
        align-items: center;

        .stat-icon {
          display: flex;
          align-items: center;
          justify-content: center;
          width: 64px;
          height: 64px;
          border-radius: 8px;
          font-size: 32px;
          color: #fff;
        }

        .stat-info {
          margin-left: 16px;

          .stat-title {
            color: rgba(0, 0, 0, 0.45);
            font-size: 14px;
          }

          .stat-value {
            font-size: 24px;
            font-weight: 600;
            margin-top: 4px;
          }
        }
      }
    }
  }
</style>
