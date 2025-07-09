<template>
  <div v-if="visible" class="modern-modal-mask">
    <div class="modern-modal-wrap">
      <div class="modern-modal">
        <!-- 头部 -->
        <div class="modern-modal-header">
          <div class="header-content">
            <div class="modal-title">
              <svg class="title-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 12l2 2 4-4"/>
                <path d="M21 12c-1 0-3-1-3-3s2-3 3-3 3 1 3 3-2 3-3 3"/>
                <path d="M3 12c1 0 3-1 3-3s-2-3-3-3-3 1-3 3 2 3 3 3"/>
                <path d="M12 3c0 1-1 3-3 3s-3-2-3-3 1-3 3-3 3 2 3 3"/>
                <path d="M12 21c0-1-1-3-3-3s-3 2-3 3 1 3 3 3 3-2 3-3"/>
              </svg>
              <span>导航管理</span>
            </div>
            <div class="header-description">管理您的快速导航项目</div>
          </div>
          <button class="close-btn" @click="close">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>

        <!-- 工具栏 -->
        <div class="modern-toolbar">
          <div class="toolbar-left">
            <div class="select-wrapper">
              <label class="modern-checkbox">
                <input type="checkbox" v-model="allChecked" @change="toggleAll" />
                <span class="checkbox-mark"></span>
                <span class="checkbox-text">全选</span>
              </label>
            </div>
            <div class="selected-count" v-if="selectedIds.length">
              已选择 {{ selectedIds.length }} 项
            </div>
          </div>
          <div class="toolbar-right">
            <button 
              class="action-btn danger-btn" 
              :disabled="!selectedIds.length" 
              @click="confirmBatchDelete"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3,6 5,6 21,6"/>
                <path d="M19,6v14a2,2 0 0,1-2,2H7a2,2 0 0,1-2-2V6m3,0V4a2,2 0 0,1 2-2h4a2,2 0 0,1 2,2v2"/>
                <line x1="10" y1="11" x2="10" y2="17"/>
                <line x1="14" y1="11" x2="14" y2="17"/>
              </svg>
              批量删除
            </button>
            <button class="action-btn primary-btn" @click="onAdd">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              新增导航
            </button>
          </div>
        </div>

        <!-- 内容区域 -->
        <div class="modern-modal-body">
          <!-- 空状态 -->
          <div v-if="!props.navList.length" class="empty-state">
            <div class="empty-illustration">
              <svg viewBox="0 0 200 200" fill="none">
                <circle cx="100" cy="100" r="80" fill="url(#gradient1)" opacity="0.1"/>
                <circle cx="100" cy="100" r="60" fill="url(#gradient2)" opacity="0.2"/>
                <circle cx="100" cy="100" r="40" fill="url(#gradient3)" opacity="0.3"/>
                <defs>
                  <linearGradient id="gradient1" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" style="stop-color:#1677ff;stop-opacity:1" />
                    <stop offset="100%" style="stop-color:#69b1ff;stop-opacity:1" />
                  </linearGradient>
                  <linearGradient id="gradient2" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" style="stop-color:#ff6b6b;stop-opacity:1" />
                    <stop offset="100%" style="stop-color:#ffd93d;stop-opacity:1" />
                  </linearGradient>
                  <linearGradient id="gradient3" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" style="stop-color:#6c5ce7;stop-opacity:1" />
                    <stop offset="100%" style="stop-color:#a29bfe;stop-opacity:1" />
                  </linearGradient>
                </defs>
              </svg>
            </div>
            <div class="empty-title">暂无导航项目</div>
            <div class="empty-description">创建您的第一个导航项目，开始管理您的快捷链接</div>
            <button class="empty-action-btn" @click="onAdd">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              立即创建
            </button>
          </div>

          <!-- 导航列表 -->
          <div v-else class="nav-grid">
            <div 
              class="nav-card" 
              v-for="(nav, idx) in props.navList" 
              :key="nav.id"
              :class="{ 'selected': selectedIds.includes(nav.id) }"
            >
              <div class="card-header">
                <label class="card-checkbox" @click.stop>
                  <input type="checkbox" v-model="selectedIds" :value="nav.id" />
                  <span class="checkbox-mark"></span>
                </label>
                <div class="card-actions">
                  <button 
                    class="card-action-btn move-up-btn" 
                    @click.stop="moveUp(idx)" 
                    :disabled="idx === 0"
                    :title="idx === 0 ? '已在最顶部，无法上移' : '上移'"

                  >
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="18,15 12,9 6,15"/>
                    </svg>
                  </button>
                  <button 
                    class="card-action-btn move-down-btn" 
                    @click.stop="moveDown(idx)" 
                    :disabled="idx >= props.navList.length - 1"
                    :title="idx >= props.navList.length - 1 ? '已在最底部，无法下移' : '下移'"

                  >
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="6,9 12,15 18,9"/>
                    </svg>
                  </button>
                  <button 
                    class="card-action-btn edit-btn" 
                    @click.stop="onEdit(nav)"
                    title="编辑"

                  >
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                      <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                  </button>
                  <button 
                    class="card-action-btn delete-btn" 
                    @click.stop="onDelete(nav.id)"
                    title="删除"

                  >
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="3,6 5,6 21,6"/>
                      <path d="M19,6v14a2,2 0 0,1-2,2H7a2,2 0 0,1-2-2V6m3,0V4a2,2 0 0,1 2-2h4a2,2 0 0,1 2,2v2"/>
                      <line x1="10" y1="11" x2="10" y2="17"/>
                      <line x1="14" y1="11" x2="14" y2="17"/>
                    </svg>
                  </button>
                </div>
              </div>
              <div class="card-content">
                <div class="nav-logo-wrapper">
              <img class="nav-logo" :src="getNavLogo(nav)" :alt="nav.title" />
                  <div class="logo-overlay"></div>
                </div>
              <div class="nav-info">
                <div class="nav-title">{{ nav.title }}</div>
                  <div class="nav-description">{{ nav.description }}</div>
              </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 确认删除弹窗 -->
    <div v-if="showConfirm" class="confirm-modal-mask">
      <div class="confirm-modal">
        <div class="confirm-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="15" y1="9" x2="9" y2="15"/>
            <line x1="9" y1="9" x2="15" y2="15"/>
          </svg>
        </div>
        <div class="confirm-title">确认删除</div>
        <div class="confirm-description">
          您确定要删除选中的 {{ selectedIds.length }} 个导航项目吗？此操作无法撤销。
        </div>
        <div class="confirm-actions">
          <button class="confirm-btn cancel-btn" @click="showConfirm = false">
            取消
          </button>
          <button class="confirm-btn danger-btn" @click="doBatchDelete">
            确认删除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, defineExpose } from 'vue'
const props = defineProps<{ navList: any[], getNavLogo: (nav:any)=>string }>()
const emit = defineEmits(['close', 'edit', 'delete', 'add', 'move', 'batchDelete'])
const visible = ref(false)
const selectedIds = ref<number[]>([])
const allChecked = ref(false)
const showConfirm = ref(false)

function show() { 
  visible.value = true 
}
function close() { visible.value = false; emit('close') }
function onEdit(nav) { 
  emit('edit', nav) 
}
function onDelete(id) { 
  emit('delete', id) 
}
function onAdd() { emit('add') }
function moveUp(idx) { 
  emit('move', {from: idx, to: idx-1}) 
}
function moveDown(idx) { 
  emit('move', {from: idx, to: idx+1}) 
}
function toggleAll() {
  if (allChecked.value) selectedIds.value = props.navList.map(n=>n.id)
  else selectedIds.value = []
}
function toggleSelect(id) {
  const i = selectedIds.value.indexOf(id)
  if (i === -1) selectedIds.value.push(id)
  else selectedIds.value.splice(i, 1)
}
function confirmBatchDelete() { showConfirm.value = true }
function doBatchDelete() {
  emit('batchDelete', selectedIds.value)
  selectedIds.value = []
  allChecked.value = false
  showConfirm.value = false
}
watch(() => props.navList, () => { selectedIds.value = []; allChecked.value = false }, {deep:true})
defineExpose({ show, close })
</script>

<style lang="scss" scoped>
/* 现代化弹窗遮罩层 */
.modern-modal-mask {
  position: fixed;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(8px);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modern-modal-wrap {
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

/* 现代化弹窗容器 */
.modern-modal {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  box-shadow: 
    0 20px 60px rgba(0, 0, 0, 0.12),
    0 8px 32px rgba(0, 0, 0, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
  width: 920px;
  max-width: 95vw;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.2);
  animation: slideUp 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(40px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* 头部样式 */
.modern-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28px 32px 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(255, 255, 255, 0.7));
}

.header-content {
  flex: 1;
  
  .modal-title {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 24px;
    font-weight: 600;
    color: #1a1a1a;
    margin-bottom: 4px;
    
    .title-icon {
      width: 28px;
      height: 28px;
      stroke: #1677ff;
      stroke-width: 2;
    }
  }
  
  .header-description {
    font-size: 14px;
    color: #666;
    margin-left: 40px;
  }
}

.close-btn {
  width: 44px;
  height: 44px;
  border: none;
  background: rgba(0, 0, 0, 0.04);
  border-radius: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  
  svg {
    width: 20px;
    height: 20px;
    stroke: #666;
  }
  
  &:hover {
    background: rgba(255, 77, 79, 0.1);
    transform: scale(1.05);
    
    svg {
      stroke: #ff4d4f;
    }
  }
}

/* 工具栏样式 */
.modern-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 32px;
  background: rgba(246, 248, 250, 0.5);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.modern-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  
  input[type="checkbox"] {
    position: absolute;
    opacity: 0;
    width: 0;
    height: 0;
  }
  
  .checkbox-mark {
    width: 20px;
    height: 20px;
    border: 2px solid #d9d9d9;
    border-radius: 6px;
    background: #fff;
    position: relative;
    transition: all 0.3s ease;
    
    &::after {
      content: '';
      position: absolute;
      left: 6px;
      top: 2px;
      width: 4px;
      height: 8px;
      border: solid #fff;
      border-width: 0 2px 2px 0;
      transform: rotate(45deg);
      opacity: 0;
      transition: opacity 0.3s ease;
    }
  }
  
  input:checked + .checkbox-mark {
    background: #1677ff;
    border-color: #1677ff;
    
    &::after {
      opacity: 1;
    }
  }
  
  .checkbox-text {
    font-size: 14px;
    color: #333;
    font-weight: 500;
  }
}

.selected-count {
  font-size: 14px;
  color: #1677ff;
  font-weight: 500;
  padding: 6px 12px;
  background: rgba(22, 119, 255, 0.1);
  border-radius: 20px;
}

.toolbar-right {
  display: flex;
  gap: 12px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  
  svg {
    width: 18px;
    height: 18px;
  }
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

.primary-btn {
  background: linear-gradient(135deg, #1677ff, #69b1ff);
  color: #fff;
  box-shadow: 0 4px 16px rgba(22, 119, 255, 0.3);
  
  &:hover:not(:disabled) {
    box-shadow: 0 6px 24px rgba(22, 119, 255, 0.4);
    transform: translateY(-2px);
  }
  
  &:active {
    transform: translateY(0);
  }
}

.danger-btn {
  background: rgba(255, 77, 79, 0.1);
  color: #ff4d4f;
  border: 1px solid rgba(255, 77, 79, 0.3);
  
  &:hover:not(:disabled) {
    background: rgba(255, 77, 79, 0.15);
    border-color: #ff4d4f;
    transform: translateY(-1px);
  }
}

/* 主体内容 */
.modern-modal-body {
  flex: 1;
  padding: 24px 32px 32px;
  overflow-y: auto;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.8), rgba(246, 248, 250, 0.5));
}

/* 空状态样式 */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  
  .empty-illustration {
    width: 200px;
    height: 200px;
    margin: 0 auto 24px;
    opacity: 0.8;
  }
  
  .empty-title {
    font-size: 20px;
    font-weight: 600;
    color: #333;
    margin-bottom: 8px;
  }
  
  .empty-description {
    font-size: 14px;
    color: #666;
    margin-bottom: 32px;
    line-height: 1.6;
  }
  
  .empty-action-btn {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 12px 24px;
    background: linear-gradient(135deg, #1677ff, #69b1ff);
    color: #fff;
    border: none;
    border-radius: 12px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s ease;
    box-shadow: 0 4px 16px rgba(22, 119, 255, 0.3);
    
    svg {
      width: 18px;
      height: 18px;
    }
    
    &:hover {
      box-shadow: 0 6px 24px rgba(22, 119, 255, 0.4);
      transform: translateY(-2px);
    }
  }
}

/* 网格布局 */
.nav-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

/* 导航卡片 */
.nav-card {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  cursor: default;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, rgba(22, 119, 255, 0.05), rgba(105, 177, 255, 0.05));
    opacity: 0;
    transition: opacity 0.3s ease;
  }
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
    
    &::before {
      opacity: 1;
    }
  }
  
  &.selected {
    border-color: #1677ff;
    background: rgba(22, 119, 255, 0.05);
    
    &::before {
      opacity: 1;
    }
  }
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.card-checkbox {
  display: flex;
  align-items: center;
  cursor: pointer;
  
    input[type="checkbox"] {
    position: absolute;
    opacity: 0;
    width: 0;
    height: 0;
  }
  
  .checkbox-mark {
    width: 18px;
    height: 18px;
    border: 2px solid #d9d9d9;
    border-radius: 4px;
    background: #fff;
    position: relative;
    transition: all 0.3s ease;
    
    &::after {
      content: '';
      position: absolute;
      left: 4px;
      top: 1px;
      width: 3px;
      height: 6px;
      border: solid #fff;
      border-width: 0 2px 2px 0;
      transform: rotate(45deg);
      opacity: 0;
      transition: opacity 0.3s ease;
    }
  }
  
  input:checked + .checkbox-mark {
    background: #1677ff;
    border-color: #1677ff;
    
    &::after {
      opacity: 1;
    }
  }
}

.card-actions {
  display: flex;
  gap: 6px;
  opacity: 0.9;
  transition: opacity 0.3s ease;
  
  .nav-card:hover & {
    opacity: 1;
  }
}

.card-action-btn {
  width: 36px;
  height: 36px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.9);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  pointer-events: auto;
  z-index: 10;
  
  svg {
    width: 18px;
    height: 18px;
    stroke: #666;
  }
  
  &:hover:not(:disabled) {
    background: rgba(255, 255, 255, 1);
    transform: scale(1.05);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
  }
  
  &:disabled {
    opacity: 0.3;
    cursor: not-allowed;
    background: rgba(0, 0, 0, 0.05);
    transform: none !important;
    
    svg {
      stroke: #ccc;
    }
    
    &:hover {
      transform: none !important;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
    }
  }
  
  &.edit-btn:hover:not(:disabled) {
    background: rgba(22, 119, 255, 0.1);
    border-color: #1677ff;
    svg { stroke: #1677ff; }
  }
  
  &.delete-btn:hover:not(:disabled) {
    background: rgba(255, 77, 79, 0.1);
    border-color: #ff4d4f;
    svg { stroke: #ff4d4f; }
  }
  
  // 为移动按钮添加特殊样式
  &.move-up-btn {
    &:hover:not(:disabled) {
      background: rgba(52, 211, 153, 0.1);
      border-color: #34d399;
      svg { stroke: #34d399; }
    }
    
    &:disabled {
      svg { stroke: #ccc; }
    }
  }
  
  &.move-down-btn {
    &:hover:not(:disabled) {
      background: rgba(52, 211, 153, 0.1);
      border-color: #34d399;
      svg { stroke: #34d399; }
    }
    
    &:disabled {
      svg { stroke: #ccc; }
    }
  }
}

.card-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.nav-logo-wrapper {
  position: relative;
  width: 48px;
  height: 48px;
  flex-shrink: 0;
  
  .nav-logo {
    width: 100%;
    height: 100%;
    border-radius: 12px;
    object-fit: cover;
    background: rgba(246, 248, 250, 0.8);
    border: 1px solid rgba(0, 0, 0, 0.06);
    transition: transform 0.3s ease;
  }
  
  .logo-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, rgba(22, 119, 255, 0.1), rgba(105, 177, 255, 0.1));
    border-radius: 12px;
    opacity: 0;
    transition: opacity 0.3s ease;
  }
  
  .nav-card:hover & {
    .nav-logo {
      transform: scale(1.05);
    }
    
    .logo-overlay {
      opacity: 1;
    }
  }
}

.nav-info {
  flex: 1;
  min-width: 0;
  
  .nav-title {
    font-size: 16px;
    font-weight: 600;
    color: #1a1a1a;
    margin-bottom: 4px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .nav-description {
    font-size: 12px;
    color: #666;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}

/* 确认删除弹窗 */
.confirm-modal-mask {
  position: fixed;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(4px);
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: fadeIn 0.2s ease-out;
}

.confirm-modal {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 32px;
  width: 420px;
  max-width: 90vw;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.2);
  animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.confirm-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 20px;
  background: rgba(255, 77, 79, 0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  
  svg {
    width: 28px;
    height: 28px;
    stroke: #ff4d4f;
  }
}

.confirm-title {
  font-size: 20px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.confirm-description {
  font-size: 14px;
  color: #666;
  margin-bottom: 28px;
  line-height: 1.6;
}

.confirm-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.confirm-btn {
  padding: 10px 24px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  min-width: 80px;
  
  &.cancel-btn {
    background: rgba(0, 0, 0, 0.06);
    color: #666;
    
    &:hover {
      background: rgba(0, 0, 0, 0.1);
    }
  }
  
  &.danger-btn {
    background: #ff4d4f;
    color: #fff;
    
    &:hover {
      background: #ff7875;
      transform: translateY(-1px);
    }
  }
}

/* 滚动条样式 */
.modern-modal-body::-webkit-scrollbar {
  width: 6px;
}

.modern-modal-body::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.04);
  border-radius: 3px;
}

.modern-modal-body::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
  
  &:hover {
    background: rgba(0, 0, 0, 0.3);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .modern-modal {
    width: 95vw;
    max-height: 95vh;
  }
  
  .modern-modal-header {
    padding: 20px 24px 16px;
  }
  
  .modern-toolbar {
    padding: 16px 24px;
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .modern-modal-body {
    padding: 20px 24px 24px;
  }
  
  .nav-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .nav-card {
    padding: 16px;
  }
}
</style>