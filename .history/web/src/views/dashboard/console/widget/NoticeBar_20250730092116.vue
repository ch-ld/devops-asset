<template>
  <div class="notice-bar art-custom-card">
    <div class="notice-title">
      <i class="iconfont-sys">&#xe6d1;</i>
      <span>公告栏</span>
      <el-button 
        type="link" 
        size="small" 
        @click="refreshNotices"
        :loading="loading"
        style="margin-left: auto"
      >
        刷新
      </el-button>
    </div>
    
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-wrapper">
      <el-skeleton :rows="3" animated />
    </div>
    
    <!-- 公告列表 -->
    <ul v-else-if="notices.length > 0" class="notice-list">
      <li 
        v-for="(item, idx) in notices" 
        :key="item.id || idx"
        @click="handleNoticeClick(item)"
        class="notice-item"
      >
        <div class="notice-left">
          <span class="notice-type" :class="getTypeClass(item.type)">
            {{ getTypeLabel(item.type) }}
          </span>
          <span class="notice-content">{{ item.title }}</span>
        </div>
        <div class="notice-right">
          <span class="notice-time">{{ formatTime(item.created_at) }}</span>
          <i v-if="item.is_pinned === 1" class="pin-icon iconfont-sys">&#xe6d2;</i>
        </div>
      </li>
    </ul>
    
    <!-- 空状态 -->
    <div v-else class="empty-state">
      <i class="iconfont-sys">&#xe8d7;</i>
      <p>暂无公告</p>
    </div>
    
    <!-- 公告详情弹窗 -->
    <el-dialog 
      v-model="showDetailDialog" 
      :title="selectedNotice?.title" 
      width="600px"
      class="notice-detail-dialog"
    >
      <div v-if="selectedNotice" class="notice-detail">
        <div class="notice-meta">
          <span class="notice-type" :class="getTypeClass(selectedNotice.type)">
            {{ getTypeLabel(selectedNotice.type) }}
          </span>
          <span class="publish-time">发布时间：{{ formatTime(selectedNotice.created_at) }}</span>
          <span v-if="selectedNotice.publisher" class="publisher">
            发布者：{{ selectedNotice.publisher.name }}
          </span>
        </div>
        <div class="notice-content" v-html="selectedNotice.content"></div>
        <div v-if="selectedNotice.read_count" class="read-count">
          阅读数：{{ selectedNotice.read_count }}
        </div>
      </div>
      
      <template #footer>
        <el-button @click="showDetailDialog = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getActiveAnnouncements, getAnnouncementDetail } from '@/api/system/announcement'

interface Notice {
  id: number
  title: string
  content: string
  type: number
  is_pinned: number
  created_at: number
  read_count: number
  publisher?: {
    id: number
    name: string
  }
}

const notices = ref<Notice[]>([])
const loading = ref(false)
const showDetailDialog = ref(false)
const selectedNotice = ref<Notice | null>(null)

// 获取公告列表
const fetchNotices = async () => {
  loading.value = true
  try {
    const response = await getActiveAnnouncements(5)
    if (response && response.code === 200) {
      notices.value = response.data || []
    } else {
      ElMessage.error('获取公告失败')
    }
  } catch (error) {
    console.error('获取公告失败:', error)
    ElMessage.error('获取公告失败')
  } finally {
    loading.value = false
  }
}

// 刷新公告
const refreshNotices = () => {
  fetchNotices()
}

// 处理公告点击
const handleNoticeClick = async (notice: Notice) => {
  try {
    const response = await getAnnouncementDetail(notice.id)
    if (response && response.code === 200) {
      selectedNotice.value = response.data
      showDetailDialog.value = true
    } else {
      ElMessage.error('获取公告详情失败')
    }
  } catch (error) {
    console.error('获取公告详情失败:', error)
    ElMessage.error('获取公告详情失败')
  }
}

// 获取类型标签
const getTypeLabel = (type: number) => {
  const typeMap = {
    1: '通知',
    2: '公告',
    3: '系统'
  }
  return typeMap[type] || '未知'
}

// 获取类型样式类
const getTypeClass = (type: number) => {
  const classMap = {
    1: 'type-notice',
    2: 'type-announcement',
    3: 'type-system'
  }
  return classMap[type] || 'type-default'
}

// 格式化时间
const formatTime = (timestamp: number | string) => {
  if (!timestamp) return ''
  try {
    // 如果是Unix时间戳（数字），需要转换为毫秒
    const date = typeof timestamp === 'number' ? new Date(timestamp * 1000) : new Date(timestamp)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    
    // 计算时间差
    const seconds = Math.floor(diff / 1000)
    const minutes = Math.floor(seconds / 60)
    const hours = Math.floor(minutes / 60)
    const days = Math.floor(hours / 24)
    
    if (days > 7) {
      // 超过7天显示具体日期
      return date.toLocaleDateString('zh-CN')
    } else if (days > 0) {
      return `${days}天前`
    } else if (hours > 0) {
      return `${hours}小时前`
    } else if (minutes > 0) {
      return `${minutes}分钟前`
    } else if (seconds > 0) {
      return `${seconds}秒前`
    } else {
      return '刚刚'
    }
  } catch (error) {
    return String(timestamp)
  }
}

onMounted(() => {
  fetchNotices()
})
</script>

<style lang="scss" scoped>
.notice-bar {
  padding: 18px 24px;
  margin-bottom: 18px;
  
  .notice-title {
    display: flex;
    align-items: center;
    font-size: 18px;
    font-weight: 600;
    color: var(--el-color-primary);
    margin-bottom: 10px;
    
    .iconfont-sys {
      font-size: 20px;
      margin-right: 8px;
    }
  }
  
  .loading-wrapper {
    padding: 10px 0;
  }
  
  .notice-list {
    list-style: none;
    padding: 0;
    margin: 0;
    
    .notice-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 0;
      border-bottom: 1px dashed var(--el-border-color-lighter);
      cursor: pointer;
      transition: background-color 0.3s;
      
      &:hover {
        background-color: var(--el-color-primary-light-9);
        border-radius: 4px;
        padding-left: 8px;
        padding-right: 8px;
      }
      
      &:last-child {
        border-bottom: none;
      }
      
      .notice-left {
        display: flex;
        align-items: center;
        flex: 1;
        min-width: 0;
        
        .notice-type {
          font-size: 12px;
          padding: 2px 6px;
          border-radius: 4px;
          margin-right: 8px;
          white-space: nowrap;
          
          &.type-notice {
            background-color: var(--el-color-info-light-8);
            color: var(--el-color-info);
          }
          
          &.type-announcement {
            background-color: var(--el-color-warning-light-8);
            color: var(--el-color-warning);
          }
          
          &.type-system {
            background-color: var(--el-color-danger-light-8);
            color: var(--el-color-danger);
          }
        }
        
        .notice-content {
          flex: 1;
          @include ellipsis;
          font-size: 15px;
          color: var(--art-text-gray-900);
        }
      }
      
      .notice-right {
        display: flex;
        align-items: center;
        margin-left: 16px;
        
        .notice-time {
          color: var(--art-text-gray-500);
          font-size: 13px;
          white-space: nowrap;
        }
        
        .pin-icon {
          color: var(--el-color-warning);
          margin-left: 4px;
          font-size: 14px;
        }
      }
    }
  }
  
  .empty-state {
    text-align: center;
    padding: 40px 0;
    color: var(--art-text-gray-500);
    
    .iconfont-sys {
      font-size: 48px;
      margin-bottom: 16px;
      display: block;
    }
    
    p {
      margin: 0;
      font-size: 14px;
    }
  }
}

.notice-detail-dialog {
  .notice-detail {
    .notice-meta {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 16px;
      padding-bottom: 12px;
      border-bottom: 1px solid var(--el-border-color-lighter);
      
      .notice-type {
        font-size: 12px;
        padding: 2px 6px;
        border-radius: 4px;
        
        &.type-notice {
          background-color: var(--el-color-info-light-8);
          color: var(--el-color-info);
        }
        
        &.type-announcement {
          background-color: var(--el-color-warning-light-8);
          color: var(--el-color-warning);
        }
        
        &.type-system {
          background-color: var(--el-color-danger-light-8);
          color: var(--el-color-danger);
        }
      }
      
      .publish-time,
      .publisher {
        font-size: 12px;
        color: var(--art-text-gray-500);
      }
    }
    
    .notice-content {
      line-height: 1.6;
      color: var(--art-text-gray-900);
      margin-bottom: 16px;
      
      :deep(p) {
        margin-bottom: 8px;
      }
      
      :deep(ul),
      :deep(ol) {
        padding-left: 20px;
      }
    }
    
    .read-count {
      font-size: 12px;
      color: var(--art-text-gray-500);
      text-align: right;
    }
  }
}
</style> 
