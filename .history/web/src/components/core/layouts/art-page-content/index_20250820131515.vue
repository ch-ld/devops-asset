<!-- 布局内容 -->
<template>
  <div class="layout-content" :class="{ 'no-basic-layout': isFullPage }" :style="containerStyle">
    <!-- 节日滚动 -->
    <ArtFestivalTextScroll v-if="!isFullPage" />

    <RouterView v-if="isRefresh" v-slot="{ Component, route }" :style="contentStyle">
      <!-- 路由信息调试 -->
      <div v-if="isOpenRouteInfo === 'true'" class="route-info">
        router meta：{{ route.meta }}
      </div>

      <!-- 使用错误边界包裹路由组件，防止DOM渲染错误 -->
      <ErrorBoundary>
        <!-- 路由动画 -->
        <Transition :name="showTransitionMask ? '' : actualTransition" mode="out-in" appear>
          <KeepAlive :max="10" :exclude="keepAliveExclude" v-if="route.meta.keepAlive">
            <component :is="Component" :key="route.fullPath" />
          </KeepAlive>
          <component :is="Component" :key="route.fullPath" v-else />
        </Transition>
      </ErrorBoundary>
    </RouterView>

    <!-- 全屏页面切换过渡遮罩（用于提升页面切换视觉体验） -->
    <Teleport to="body">
      <div v-show="showTransitionMask" class="full-page-mask" />
    </Teleport>
  </div>
</template>
<script setup lang="ts">
  import type { CSSProperties } from 'vue'
  import { useRoute } from 'vue-router'
  import { useCommon } from '@/composables/useCommon'
  import { useSettingStore } from '@/store/modules/setting'
  import { useWorktabStore } from '@/store/modules/worktab'
  import { defineComponent, h } from 'vue'

  // 定义错误边界组件
  const ErrorBoundary = defineComponent({
    name: 'ErrorBoundary',
    setup(_, { slots }) {
      const error = ref(false)
      const errorMsg = ref('')

      onErrorCaptured((err) => {
        error.value = true
        errorMsg.value = err.message
        console.error('路由组件渲染错误:', err)
        return false // 阻止错误继续传播
      })

      return () => {
        if (error.value) {
          return h('div', { class: 'error-container' }, [
            h('h3', { class: 'error-title' }, '页面渲染错误'),
            h('p', { class: 'error-message' }, errorMsg.value),
            h(
              'button',
              {
                class: 'error-retry-btn',
                onClick: () => {
                  error.value = false
                  window.location.reload() // 重新加载页面
                }
              },
              '刷新页面'
            )
          ])
        }

        return slots.default?.()
      }
    }
  })

  const route = useRoute()
  const { containerMinHeight } = useCommon()
  const { pageTransition, containerWidth, refresh } = storeToRefs(useSettingStore())
  const { keepAliveExclude } = storeToRefs(useWorktabStore())

  const isRefresh = shallowRef(true)
  const isOpenRouteInfo = import.meta.env.VITE_OPEN_ROUTE_INFO
  const showTransitionMask = ref(false)

  // 检查当前路由是否需要使用无基础布局模式
  const isFullPage = computed(() => route.matched.some((r) => r.meta?.isFullPage))
  const prevIsFullPage = ref(isFullPage.value)

  // 切换动画名称：从全屏返回时不使用动画
  const actualTransition = computed(() =>
    prevIsFullPage.value && !isFullPage.value ? '' : pageTransition.value
  )

  // 监听全屏状态变化，显示过渡遮罩
  watch(isFullPage, (val, oldVal) => {
    if (val !== oldVal) {
      showTransitionMask.value = true
      // 延迟隐藏遮罩，给足时间让页面完成切换
      setTimeout(() => {
        showTransitionMask.value = false
      }, 50)
    }

    nextTick(() => {
      prevIsFullPage.value = val
    })
  })

  const containerStyle = computed(
    (): CSSProperties =>
      isFullPage.value
        ? {
            position: 'fixed',
            top: 0,
            left: 0,
            width: '100%',
            height: '100vh',
            zIndex: 2000,
            background: 'var(--art-bg-color)'
          }
        : {
            maxWidth: containerWidth.value
          }
  )

  const contentStyle = computed(
    (): CSSProperties => ({
      minHeight: containerMinHeight.value
    })
  )

  const reload = () => {
    isRefresh.value = false
    nextTick(() => {
      isRefresh.value = true
    })
  }

  watch(refresh, reload, { flush: 'post' })
</script>

<style lang="scss" scoped>
  .layout-content {
    &.no-basic-layout {
      overflow: auto;
    }
  }

  .route-info {
    padding: 6px 8px;
    margin-bottom: 12px;
    font-size: 14px;
    color: var(--art-gray-600);
    background: var(--art-gray-200);
    border: 1px solid var(--art-border-dashed-color);
    border-radius: 6px;
  }

  .full-page-mask {
    position: fixed;
    top: 0;
    left: 0;
    z-index: 2000;
    width: 100vw;
    height: 100vh;
    pointer-events: none;
    background-color: var(--art-main-bg-color);
  }

  /* 错误边界样式 */
  :deep(.error-container) {
    padding: 24px;
    margin: 20px auto;
    max-width: 500px;
    text-align: center;
    background-color: var(--art-bg-color);
    border: 1px solid var(--art-border-color);
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  :deep(.error-title) {
    margin-bottom: 16px;
    font-size: 20px;
    color: var(--art-text-color);
  }

  :deep(.error-message) {
    margin-bottom: 20px;
    font-size: 14px;
    color: var(--art-text-secondary-color);
    word-break: break-word;
  }

  :deep(.error-retry-btn) {
    padding: 8px 16px;
    font-size: 14px;
    color: #fff;
    background-color: var(--art-primary-color);
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.3s;

    &:hover {
      background-color: var(--art-primary-color-hover);
    }
  }
</style>
