<!-- 通知组件 -->
<template>
  <div
    class="notice"
    v-show="visible"
    :style="{
      transform: show ? 'scaleY(1)' : 'scaleY(0.9)',
      opacity: show ? 1 : 0
    }"
    @click.stop=""
  >
    <div class="header">
      <span class="text">{{ $t('notice.title') }}</span>
      <span class="btn">{{ $t('notice.btnRead') }}</span>
    </div>

    <ul class="bar">
      <li
        v-for="(item, index) in barList"
        :key="index"
        :class="{ active: barActiveIndex === index }"
        @click="changeBar(index)"
      >
        {{ item.name }} ({{ item.num }})
      </li>
    </ul>

    <div class="content">
      <div class="scroll">
        <!-- 通知 -->
        <ul class="notice-list" v-show="barActiveIndex === 0">
          <li v-for="(item, index) in noticeList" :key="index">
            <div
              class="icon"
              :style="{ background: getNoticeStyle(item.type).backgroundColor + '!important' }"
            >
              <i
                class="iconfont-sys"
                :style="{ color: getNoticeStyle(item.type).iconColor + '!important' }"
                v-html="getNoticeStyle(item.type).icon"
              >
              </i>
            </div>
            <div class="text">
              <h4>{{ item.title }}</h4>
              <p>{{ item.time }}</p>
            </div>
          </li>
        </ul>

        <!-- 消息 -->
        <ul class="user-list" v-show="barActiveIndex === 1">
          <li v-for="(item, index) in msgList" :key="index">
            <div class="avatar">
              <img :src="item.avatar" />
            </div>
            <div class="text">
              <h4>{{ item.title }}</h4>
              <p>{{ item.time }}</p>
            </div>
          </li>
        </ul>

        <!-- 待办 -->
        <ul class="base" v-show="barActiveIndex === 2">
          <li v-for="(item, index) in pendingList" :key="index">
            <h4>{{ item.title }}</h4>
            <p>{{ item.time }}</p>
          </li>
        </ul>

        <!-- 空状态 -->
        <div class="empty-tips" v-show="currentTabIsEmpty">
          <i class="iconfont-sys">&#xe8d7;</i>
          <p>{{ $t('notice.text[0]') }}{{ barList[barActiveIndex].name }}</p>
        </div>
      </div>

      <div class="btn-wrapper">
        <el-button class="view-all" @click="handleViewAll" v-ripple>
          {{ $t('notice.viewAll') }}
        </el-button>
      </div>
    </div>

    <div style="height: 100px"></div>
  </div>
</template>

<script setup lang="ts">
  import { computed, ref, watch, onMounted, type Ref, type ComputedRef } from 'vue'
  import { useI18n } from 'vue-i18n'
  import AppConfig from '@/config'
  import { getActiveAnnouncements } from '@/api/system/announcement'

  defineOptions({ name: 'ArtNotification' })

  interface NoticeItem {
    /** 标题 */
    title: string
    /** 时间 */
    time: string
    /** 类型 */
    type: NoticeType
  }

  interface MessageItem {
    /** 标题 */
    title: string
    /** 时间 */
    time: string
    /** 头像 */
    avatar: string
  }

  interface PendingItem {
    /** 标题 */
    title: string
    /** 时间 */
    time: string
  }

  interface BarItem {
    /** 名称 */
    name: ComputedRef<string>
    /** 数量 */
    num: number
  }

  interface NoticeStyle {
    /** 图标 */
    icon: string
    /** 图标颜色 */
    iconColor: string
    /** 背景颜色 */
    backgroundColor: string
  }

  type NoticeType = 'email' | 'message' | 'collection' | 'user' | 'notice'

  const { t } = useI18n()

  const props = defineProps<{
    value: boolean
  }>()

  const show = ref(false)
  const visible = ref(false)
  const barActiveIndex = ref(0)

  const useNotificationData = () => {
    // 通知数据（包含系统公告）
    const noticeList = ref<NoticeItem[]>([])
    
    // 格式化时间
    const formatTime = (timestamp: number | string) => {
      if (!timestamp) return ''
      const date = typeof timestamp === 'number' ? new Date(timestamp * 1000) : new Date(timestamp)
      return date.toLocaleString('zh-CN')
    }

    // 获取系统公告
    const loadAnnouncements = async () => {
      try {
        // 实际调用获取活跃公告的API
        const response = await getActiveAnnouncements(10)
        if (response && response.code === 200) {
          const announcements = response.data.map(item => ({
            title: `【${getTypeLabel(item.type)}】${item.title}`,
            time: formatTime(item.created_at),
            type: getAnnouncementNoticeType(item.type)
          }))
          
          // 合并公告和其他通知
          const otherNotices = [
            {
              title: '冷月呆呆给你发了一条消息',
              time: '2024-4-21 8:05',
              type: 'message' as NoticeType
            },
            {
              title: '小肥猪关注了你',
              time: '2020-3-17 21:12',
              type: 'collection' as NoticeType
            },
            {
              title: '小肥猪给你发了一封邮件',
              time: '2024-1-20 0:15',
              type: 'email' as NoticeType
            }
          ]
          
          noticeList.value = [...announcements, ...otherNotices]
          return
        }
        
        // 如果API调用失败，使用模拟数据
        const mockNotices = [
          {
            title: '【系统公告】网站维护通知',
            time: '2024-6-13 10:30',
            type: 'notice' as NoticeType
          },
          {
            title: '【重要通知】功能更新说明',
            time: '2024-6-12 15:45',
            type: 'notice' as NoticeType
          },
          {
            title: '冷月呆呆给你发了一条消息',
            time: '2024-4-21 8:05',
            type: 'message' as NoticeType
          },
          {
            title: '小肥猪关注了你',
            time: '2020-3-17 21:12',
            type: 'collection' as NoticeType
          },
          {
            title: '【系统提醒】定期数据备份',
            time: '2024-6-10 9:00',
            type: 'notice' as NoticeType
          },
          {
            title: '小肥猪给你发了一封邮件',
            time: '2024-1-20 0:15',
            type: 'email' as NoticeType
          }
        ]
        noticeList.value = mockNotices
      } catch (error) {
        console.error('加载公告失败:', error)
        // 使用默认数据作为后备
        noticeList.value = [
          {
            title: '系统公告加载失败，请稍后重试',
            time: new Date().toLocaleString(),
            type: 'notice' as NoticeType
          }
        ]
      }
    }

    // 获取公告类型标签
    const getTypeLabel = (type: number) => {
      const typeMap = {
        1: '通知',
        2: '公告', 
        3: '系统'
      }
      return typeMap[type] || '通知'
    }

    // 格式化时间
    const formatTime = (timestamp: number | string) => {
      if (!timestamp) return ''
      try {
        // 如果是Unix时间戳（数字），需要转换为毫秒
        const date = typeof timestamp === 'number' ? new Date(timestamp * 1000) : new Date(timestamp)
        return date.toLocaleString('zh-CN', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit'
        })
      } catch (error) {
        return String(timestamp)
      }
    }

    // 将公告类型映射到通知类型
    const getAnnouncementNoticeType = (announcementType: number): NoticeType => {
      // 1: 通知, 2: 公告, 3: 系统
      switch (announcementType) {
        case 1:
          return 'notice'
        case 2:
          return 'notice'
        case 3:
          return 'notice'
        default:
          return 'notice'
      }
    }

    // 消息数据
    const msgList = ref<MessageItem[]>([
      {
        title: '池不胖 关注了你',
        time: '2021-2-26 23:50',
        avatar: '/src/assets/img/avatar/avatar1.webp'
      },
      {
        title: '唐不苦 关注了你',
        time: '2021-2-21 8:05',
        avatar: '/src/assets/img/avatar/avatar2.webp'
      },
      {
        title: '中小鱼 关注了你',
        time: '2020-1-17 21:12',
        avatar: '/src/assets/img/avatar/avatar3.webp'
      },
      {
        title: '何小荷 关注了你',
        time: '2021-01-14 0:20',
        avatar: '/src/assets/img/avatar/avatar4.webp'
      },
      {
        title: '誶誶淰 关注了你',
        time: '2020-12-20 0:15',
        avatar: '/src/assets/img/avatar/avatar5.webp'
      },
      {
        title: '冷月呆呆 关注了你',
        time: '2020-12-17 22:06',
        avatar: '/src/assets/img/avatar/avatar6.webp'
      }
    ])

    // 待办数据
    const pendingList = ref<PendingItem[]>([])

    // 标签栏数据
    const barList = computed<BarItem[]>(() => [
      {
        name: computed(() => t('notice.bar[0]')),
        num: noticeList.value.length
      },
      {
        name: computed(() => t('notice.bar[1]')),
        num: msgList.value.length
      },
      {
        name: computed(() => t('notice.bar[2]')),
        num: pendingList.value.length
      }
    ])

    return {
      noticeList,
      msgList,
      pendingList,
      barList,
      loadAnnouncements
    }
  }

  // 样式管理
  const useNotificationStyles = () => {
    const noticeStyleMap: Record<NoticeType, NoticeStyle> = {
      email: {
        icon: '&#xe72e;',
        iconColor: 'rgb(var(--art-warning))',
        backgroundColor: 'rgb(var(--art-bg-warning))'
      },
      message: {
        icon: '&#xe747;',
        iconColor: 'rgb(var(--art-success))',
        backgroundColor: 'rgb(var(--art-bg-success))'
      },
      collection: {
        icon: '&#xe714;',
        iconColor: 'rgb(var(--art-danger))',
        backgroundColor: 'rgb(var(--art-bg-danger))'
      },
      user: {
        icon: '&#xe608;',
        iconColor: 'rgb(var(--art-info))',
        backgroundColor: 'rgb(var(--art-bg-info))'
      },
      notice: {
        icon: '&#xe6c2;',
        iconColor: 'rgb(var(--art-primary))',
        backgroundColor: 'rgb(var(--art-bg-primary))'
      }
    }

    const getRandomColor = (): string => {
      const index = Math.floor(Math.random() * AppConfig.systemMainColor.length)
      return AppConfig.systemMainColor[index]
    }

    const getNoticeStyle = (type: NoticeType): NoticeStyle => {
      const defaultStyle: NoticeStyle = {
        icon: '&#xe747;',
        iconColor: '#FFFFFF',
        backgroundColor: getRandomColor()
      }

      return noticeStyleMap[type] || defaultStyle
    }

    return {
      getNoticeStyle
    }
  }

  // 动画管理
  const useNotificationAnimation = () => {
    const showNotice = (open: boolean) => {
      if (open) {
        visible.value = open
        setTimeout(() => {
          show.value = open
        }, 5)
      } else {
        show.value = open
        setTimeout(() => {
          visible.value = open
        }, 350)
      }
    }

    return {
      showNotice
    }
  }

  // 标签页管理
  const useTabManagement = (
    noticeList: Ref<NoticeItem[]>,
    msgList: Ref<MessageItem[]>,
    pendingList: Ref<PendingItem[]>,
    businessHandlers: {
      handleNoticeAll: () => void
      handleMsgAll: () => void
      handlePendingAll: () => void
    }
  ) => {
    const changeBar = (index: number) => {
      barActiveIndex.value = index
    }

    // 检查当前标签页是否为空
    const currentTabIsEmpty = computed(() => {
      const tabDataMap = [noticeList.value, msgList.value, pendingList.value]

      const currentData = tabDataMap[barActiveIndex.value]
      return currentData && currentData.length === 0
    })

    const handleViewAll = () => {
      // 查看全部处理器映射
      const viewAllHandlers: Record<number, () => void> = {
        0: businessHandlers.handleNoticeAll,
        1: businessHandlers.handleMsgAll,
        2: businessHandlers.handlePendingAll
      }

      const handler = viewAllHandlers[barActiveIndex.value]
      handler?.()
    }

    return {
      changeBar,
      currentTabIsEmpty,
      handleViewAll
    }
  }

  // 业务逻辑处理
  const useBusinessLogic = () => {
    const handleNoticeAll = () => {
      // 处理查看全部通知
      console.log('查看全部通知')
    }

    const handleMsgAll = () => {
      // 处理查看全部消息
      console.log('查看全部消息')
    }

    const handlePendingAll = () => {
      // 处理查看全部待办
      console.log('查看全部待办')
    }

    return {
      handleNoticeAll,
      handleMsgAll,
      handlePendingAll
    }
  }

  // 组合所有逻辑
  const { noticeList, msgList, pendingList, barList, loadAnnouncements } = useNotificationData()
  const { getNoticeStyle } = useNotificationStyles()
  const { showNotice } = useNotificationAnimation()
  const { handleNoticeAll, handleMsgAll, handlePendingAll } = useBusinessLogic()
  const { changeBar, currentTabIsEmpty, handleViewAll } = useTabManagement(
    noticeList,
    msgList,
    pendingList,
    { handleNoticeAll, handleMsgAll, handlePendingAll }
  )

  // 监听属性变化
  watch(
    () => props.value,
    (newValue) => {
      showNotice(newValue)
    }
  )

  // 组件挂载时加载公告数据
  onMounted(() => {
    loadAnnouncements()
  })
</script>

<style lang="scss" scoped>
  @use './style';
</style>
