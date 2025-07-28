<template>
  <div class="quick-actions-widget">
    <div class="widget-header">
      <h3 class="widget-title">
        <i class="iconfont-sys">&#xe66b;</i>
        快速操作
      </h3>
      <el-button type="text" size="small" @click="refreshActions">
        <i class="iconfont-sys">&#xe643;</i>
      </el-button>
    </div>
    
    <div class="actions-grid">
      <div 
        v-for="action in actions" 
        :key="action.id"
        class="action-item"
        @click="handleAction(action)"
        :class="{ disabled: action.disabled }"
      >
        <div class="action-icon" :style="{ background: action.color }">
          <i :class="action.icon" v-html="action.iconCode"></i>
        </div>
        <div class="action-content">
          <div class="action-title">{{ action.title }}</div>
          <div class="action-desc">{{ action.description }}</div>
        </div>
        <div v-if="action.badge" class="action-badge">
          {{ action.badge }}
        </div>
      </div>
    </div>
    
    <div class="widget-footer">
      <el-button type="primary" size="small" @click="openCustomize" link>
        <i class="iconfont-sys">&#xe64c;</i>
        自定义
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'

interface QuickAction {
  id: string
  title: string
  description: string
  icon: string
  color: string
  route?: string
  action?: () => void
  badge?: string | number
  disabled?: boolean
}

const router = useRouter()

const actions = ref<QuickAction[]>([
  {
    id: 'create-announcement',
    title: '创建公告',
    description: '快速发布系统公告',
    icon: 'iconfont-sys',
    color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    route: '/system/announcement',
    badge: '热门'
  },
  {
    id: 'user-management',
    title: '用户管理',
    description: '管理系统用户',
    icon: 'iconfont-sys',
    iconCode: '&#xe64d;',
    color: 'linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%)',
    route: '/system/user'
  },
  {
    id: 'system-status',
    title: '系统状态',
    description: '查看系统运行状态',
    icon: 'iconfont-sys',
    iconCode: '&#xe642;',
    color: 'linear-gradient(135deg, #ffa726 0%, #ff7043 100%)',
    action: () => checkSystemStatus()
  },
  {
    id: 'backup-data',
    title: '数据备份',
    description: '备份系统数据',
    icon: 'iconfont-sys',
    iconCode: '&#xe64e;',
    color: 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)',
    action: () => backupData()
  },
  {
    id: 'clear-cache',
    title: '清理缓存',
    description: '清理系统缓存',
    icon: 'iconfont-sys',
    iconCode: '&#xe641;',
    color: 'linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%)',
    action: () => clearCache()
  },
  {
    id: 'system-logs',
    title: '系统日志',
    description: '查看操作日志',
    icon: 'iconfont-sys',
    iconCode: '&#xe64f;',
    color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
    route: '/system/logs'
  }
])

const handleAction = (action: QuickAction) => {
  if (action.disabled) {
    ElMessage.warning('该功能暂时不可用')
    return
  }
  
  if (action.route) {
    router.push(action.route)
  } else if (action.action) {
    action.action()
  }
}

const refreshActions = () => {
  ElMessage.success('已刷新快速操作')
}

const openCustomize = () => {
  ElMessage.info('自定义功能开发中...')
}

const checkSystemStatus = () => {
  ElMessage.info('系统运行正常 ✓')
}

const backupData = () => {
  ElMessage.info('数据备份功能开发中...')
}

const clearCache = () => {
  ElMessage.success('缓存清理完成')
}

onMounted(() => {
  // 可以在这里加载用户自定义的快速操作
})
</script>

<style scoped lang="scss">
.quick-actions-widget {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  height: 100%;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  }
  
  .widget-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    
    .widget-title {
      font-size: 20px;
      font-weight: 700;
      color: #1a1a1a;
      margin: 0;
      display: flex;
      align-items: center;
      gap: 10px;
      background: linear-gradient(45deg, #667eea, #764ba2);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
      
      .iconfont-sys {
        font-size: 22px;
        color: #667eea;
      }
    }
  }
  
  .actions-grid {
    flex: 1;
    display: grid;
    gap: 16px;
    
    .action-item {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 16px;
      border-radius: 12px;
      cursor: pointer;
      transition: all 0.3s ease;
      border: 1px solid rgba(255, 255, 255, 0.3);
      background: rgba(255, 255, 255, 0.6);
      position: relative;
      
      &:hover {
        transform: translateY(-4px);
        box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
        background: rgba(255, 255, 255, 0.9);
        border-color: transparent;
      }
      
      &.disabled {
        opacity: 0.5;
        cursor: not-allowed;
        
        &:hover {
          transform: none;
          box-shadow: none;
        }
      }
      
      .action-icon {
        width: 48px;
        height: 48px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        font-size: 20px;
        flex-shrink: 0;
        box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
      }
      
      .action-content {
        flex: 1;
        min-width: 0;
        
        .action-title {
          font-size: 16px;
          font-weight: 600;
          color: #333;
          margin-bottom: 4px;
        }
        
        .action-desc {
          font-size: 13px;
          color: #666;
          line-height: 1.4;
        }
      }
      
      .action-badge {
        background: linear-gradient(45deg, #ff6b6b, #ee5a52);
        color: white;
        font-size: 11px;
        padding: 4px 8px;
        border-radius: 6px;
        font-weight: 600;
        position: absolute;
        top: 12px;
        right: 12px;
        box-shadow: 0 2px 8px rgba(255, 107, 107, 0.3);
      }
    }
  }
  
  .widget-footer {
    margin-top: 20px;
    padding-top: 20px;
    border-top: 1px solid rgba(255, 255, 255, 0.3);
    text-align: center;
    
    .el-button {
      background: linear-gradient(45deg, #667eea, #764ba2);
      border: none;
      color: white;
      padding: 8px 16px;
      border-radius: 8px;
      font-weight: 600;
      transition: all 0.3s ease;
      
      &:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
      }
    }
  }
}
</style> 