<template>
  <div class="announcement-marquee">
    <div class="announcement-content">
      <div class="announcement-icon">
        <i class="el-icon-bell"></i>
      </div>
      <div class="announcement-label">
        系统公告
      </div>
      <div class="announcement-content-area">
        <div v-if="displayedAnnouncements.length > 0" class="announcement-list">
          <div
            v-for="(announcement, index) in displayedAnnouncements"
            :key="announcement.id"
            class="announcement-item"
            @click="handleAnnouncementClick(announcement)"
          >
            <span class="announcement-type" :class="getTypeClass(announcement.type)">
              {{ getTypeLabel(announcement.type) }}
            </span>
            <span class="announcement-text">{{ announcement.title }}</span>
            <i v-if="announcement.is_pinned === 1" class="pin-icon el-icon-star-on"></i>
          </div>
        </div>
        <div v-else class="empty-announcement">
          <span class="empty-text">暂无系统公告</span>
        </div>
      </div>
      <div class="announcement-actions">
        <el-button type="text" size="small" @click="handleMoreClick">
          更多 <i class="el-icon-arrow-right"></i>
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getActiveAnnouncements } from '@/api/system/announcement'
import { useRouter } from 'vue-router'

interface Announcement {
  id: number
  title: string
  content: string
  type: number
  is_pinned: number
  created_at: number
  priority: number
}

const router = useRouter()
const announcements = ref<Announcement[]>([])

// 显示的公告数量限制
const maxDisplayCount = 3

// 计算显示的公告
const displayedAnnouncements = computed(() => {
  return announcements.value.slice(0, maxDisplayCount)
})

// 获取公告列表
const fetchAnnouncements = async () => {
  try {
    const response = await getActiveAnnouncements(10)
    if (response && response.code === 200) {
      announcements.value = response.data || []
    }
  } catch (error) {
    console.error('获取公告失败:', error)
  }
}

// 处理公告点击
const handleAnnouncementClick = (announcement: Announcement) => {
  ElMessage.info(`点击了公告: ${announcement.title}`)
  // 可以在这里添加查看详情的逻辑
}

// 处理更多按钮点击
const handleMoreClick = () => {
  router.push('/system/announcement')
}

// 获取类型标签
const getTypeLabel = (type: number) => {
  const typeMap = {
    1: '通知',
    2: '公告',
    3: '系统'
  }
  return typeMap[type] || '通知'
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

onMounted(() => {
  fetchAnnouncements()
})
</script>

<style lang="scss" scoped>
.announcement-marquee {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  margin-bottom: 24px;
  box-shadow: 0 4px 20px rgba(102, 126, 234, 0.3);
  overflow: hidden;
  position: relative;
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.6), transparent);
  }
  
  .announcement-content {
    display: flex;
    align-items: center;
    padding: 12px 20px;
    color: white;
    
    .announcement-icon {
      font-size: 18px;
      margin-right: 12px;
      color: #ffd700;
      animation: bell-ring 2s infinite;
      
      @keyframes bell-ring {
        0%, 100% { transform: rotate(0deg); }
        10%, 30% { transform: rotate(15deg); }
        20% { transform: rotate(-15deg); }
      }
    }
    
    .announcement-label {
      font-size: 14px;
      font-weight: 600;
      margin-right: 16px;
      white-space: nowrap;
      background: rgba(255, 255, 255, 0.2);
      padding: 4px 12px;
      border-radius: 20px;
      backdrop-filter: blur(10px);
    }
    
    .announcement-content-area {
      flex: 1;
      overflow: hidden;
      
      .announcement-list {
        display: flex;
        align-items: center;
        gap: 24px;
        flex-wrap: wrap;
        
        .announcement-item {
          display: flex;
          align-items: center;
          cursor: pointer;
          transition: all 0.3s ease;
          background: rgba(255, 255, 255, 0.1);
          padding: 8px 16px;
          border-radius: 20px;
          backdrop-filter: blur(5px);
          border: 1px solid rgba(255, 255, 255, 0.2);
          
          &:hover {
            transform: translateY(-2px);
            background: rgba(255, 255, 255, 0.2);
            box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
          }
          
          .announcement-type {
            font-size: 11px;
            padding: 3px 8px;
            border-radius: 12px;
            margin-right: 10px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            
            &.type-notice {
              background: rgba(24, 144, 255, 0.3);
              color: #87ceeb;
            }
            
            &.type-announcement {
              background: rgba(255, 193, 7, 0.3);
              color: #ffd700;
            }
            
            &.type-system {
              background: rgba(255, 77, 79, 0.3);
              color: #ff9999;
            }
          }
          
          .announcement-text {
            font-size: 14px;
            font-weight: 500;
            margin-right: 8px;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            max-width: 300px;
          }
          
          .pin-icon {
            color: #ffd700;
            font-size: 14px;
            margin-left: 4px;
            animation: pin-glow 2s ease-in-out infinite;
          }
          
          @keyframes pin-glow {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.6; }
          }
        }
      }
      
      .empty-announcement {
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 16px;
        
        .empty-text {
          color: rgba(255, 255, 255, 0.7);
          font-size: 14px;
          font-style: italic;
        }
      }
    }
    
    .announcement-actions {
      margin-left: 16px;
      
      .el-button {
        color: white;
        font-size: 13px;
        padding: 6px 12px;
        
        &:hover {
          background: rgba(255, 255, 255, 0.1);
          border-radius: 6px;
        }
        
        i {
          margin-left: 4px;
          transition: transform 0.3s ease;
        }
        
        &:hover i {
          transform: translateX(2px);
        }
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .announcement-marquee {
    .announcement-content {
      padding: 10px 16px;
      
      .announcement-label {
        font-size: 12px;
        margin-right: 12px;
      }
      
      .announcement-content-area {
        .announcement-list {
          gap: 16px;
          
          .announcement-item {
            padding: 6px 12px;
            
            .announcement-text {
              font-size: 13px;
              max-width: 200px;
            }
          }
        }
      }
    }
  }
}
</style> 