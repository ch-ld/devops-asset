<template>
  <div v-if="dialogVisible" class="modern-edit-modal-mask">
    <div class="modern-edit-modal-wrap">
      <div class="modern-edit-modal">
        <!-- 头部 -->
        <div class="modern-edit-modal-header">
          <div class="header-content">
            <div class="modal-title">
              <svg class="title-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path v-if="isEdit" d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                <path v-if="isEdit" d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                <path v-else d="M12 5v14"/>
                <path v-else d="M5 12h14"/>
              </svg>
              <span>{{ dialogTitle }}</span>
            </div>
            <div class="header-description">{{ isEdit ? '编辑导航项目信息' : '创建新的导航项目' }}</div>
          </div>
          <button class="close-btn" @click="dialogVisible = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>

        <!-- 表单内容 -->
        <div class="modern-edit-modal-body">
          <ElForm ref="formRef" :model="form" :rules="rules" class="modern-form">
            <!-- 基本信息 -->
            <div class="form-section">
              <div class="section-title">
                <svg class="section-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="3"/>
                  <path d="M12 1v6m0 6v6m11-7h-6m-6 0H1"/>
                </svg>
                基本信息
              </div>
              <div class="form-grid">
                <div class="form-item-wrapper">
                  <ElFormItem label="导航标题" prop="title" class="modern-form-item">
                    <ElInput 
                      v-model="form.title" 
                      placeholder="请输入导航标题"
                      class="modern-input"
                    >
                      <template #prefix>
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M4 6h16M4 12h16M4 18h7"/>
                        </svg>
                      </template>
                    </ElInput>
                  </ElFormItem>
                </div>
                <div class="form-item-wrapper">
                  <ElFormItem label="分组名称" prop="group" class="modern-form-item">
                    <el-select 
                      v-model="form.group" 
                      filterable 
                      allow-create 
                      default-first-option 
                      placeholder="请选择或输入分组"
                      class="modern-select"
                    >
                      <el-option v-for="g in groupOptions" :key="g" :label="g" :value="g" />
                    </el-select>
                  </ElFormItem>
                </div>
              </div>
              
              <div class="form-item-wrapper full-width">
                <ElFormItem label="链接地址" prop="url" class="modern-form-item">
                  <ElInput 
                    v-model="form.url" 
                    placeholder="请输入完整的链接地址 (https://www.example.com)"
                    class="modern-input"
                  >
                    <template #prefix>
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/>
                        <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>
                      </svg>
                    </template>
                  </ElInput>
                </ElFormItem>
              </div>
            </div>

            <!-- 外观设置 -->
            <div class="form-section">
              <div class="section-title">
                <svg class="section-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <polygon points="10,8 16,12 10,16"/>
                </svg>
                外观设置
              </div>
              <div class="form-grid">
                <div class="form-item-wrapper">
                  <ElFormItem label="图标" prop="icon" class="modern-form-item">
                    <div class="icon-selector-wrapper">
                      <ArtIconSelector
                        v-model="form.icon"
                        :iconType="iconType"
                        :defaultIcon="form.icon"
                        width="100%"
                        class="modern-icon-selector"
                      />
                      <div class="icon-tip">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <circle cx="12" cy="12" r="10"/>
                          <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/>
                          <line x1="12" y1="17" x2="12.01" y2="17"/>
                        </svg>
                        未填写时将自动抓取网站 favicon，抓取失败用默认logo
                      </div>
                    </div>
                  </ElFormItem>
                </div>
                <div class="form-item-wrapper">
                  <ElFormItem label="排序" prop="sort" class="modern-form-item">
                    <ElInputNumber
                      v-model="form.sort"
                      :min="1"
                      controls-position="right"
                      class="modern-input-number"
                    >
                      <template #decrease-icon>
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <line x1="5" y1="12" x2="19" y2="12"/>
                        </svg>
                      </template>
                      <template #increase-icon>
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <line x1="12" y1="5" x2="12" y2="19"/>
                          <line x1="5" y1="12" x2="19" y2="12"/>
                        </svg>
                      </template>
                    </ElInputNumber>
                  </ElFormItem>
                </div>
              </div>
            </div>

            <!-- 功能配置 -->
            <div class="form-section">
              <div class="section-title">
                <svg class="section-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="3"/>
                  <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1 1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
                </svg>
                功能配置
              </div>
              <div class="form-grid switch-grid">
                <div class="switch-item">
                  <div class="switch-label">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M9 12l2 2 4-4"/>
                      <path d="M21 4H8l-7 8 7 8h13a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2z"/>
                    </svg>
                    <div>
                      <div class="switch-title">启用状态</div>
                      <div class="switch-desc">控制导航是否显示在页面中</div>
                    </div>
                  </div>
                  <ElSwitch 
                    v-model="form.isEnable" 
                    class="modern-switch"
                    active-color="var(--el-color-primary)"
                    inactive-color="var(--el-border-color)"
                  />
                </div>
                <div class="switch-item">
                  <div class="switch-label">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
                      <polyline points="15,3 21,3 21,9"/>
                      <line x1="10" y1="14" x2="21" y2="3"/>
                    </svg>
                    <div>
                      <div class="switch-title">新窗口打开</div>
                      <div class="switch-desc">在新标签页中打开链接</div>
                    </div>
                  </div>
                  <ElSwitch 
                    v-model="form.openInNewTab" 
                    class="modern-switch"
                    active-color="var(--el-color-primary)"
                    inactive-color="var(--el-border-color)"
                  />
                </div>
              </div>
            </div>

            <!-- 详细描述 -->
            <div class="form-section">
              <div class="section-title">
                <svg class="section-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                  <polyline points="14,2 14,8 20,8"/>
                  <line x1="16" y1="13" x2="8" y2="13"/>
                  <line x1="16" y1="17" x2="8" y2="17"/>
                  <polyline points="10,9 9,9 8,9"/>
                </svg>
                详细描述
              </div>
              <div class="form-item-wrapper full-width">
                <ElFormItem label="描述信息" prop="description" class="modern-form-item">
                  <ElInput
                    v-model="form.description"
                    type="textarea"
                    :rows="4"
                    placeholder="请输入导航描述，帮助用户了解这个导航的用途..."
                    class="modern-textarea"
                    maxlength="200"
                    show-word-limit
                  />
                </ElFormItem>
              </div>
            </div>
          </ElForm>
        </div>

        <!-- 底部操作栏 -->
        <div class="modern-edit-modal-footer">
          <button class="footer-btn cancel-btn" @click="dialogVisible = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
            取消
          </button>
          <button class="footer-btn submit-btn" @click="submitForm">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20,6 9,17 4,12"/>
            </svg>
            {{ isEdit ? '保存修改' : '创建导航' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, nextTick } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { IconTypeEnum } from '@/enums/appEnum'
import { addNav, updateNav } from '@/api/system/nav'
import { ApiStatus } from '@/utils/http/status'

const props = defineProps({
  groupOptions: {
    type: Array,
    default: () => []
  }
})

const dialogVisible = ref(false)
const form = reactive({
  id: 0,
  title: '',
  url: '',
  icon: '',
  group: '',
  sort: 1,
  isEnable: true,
  openInNewTab: true,
  description: ''
})
const iconType = ref(IconTypeEnum.UNICODE)
const isEdit = ref(false)
const formRef = ref<FormInstance>()

const rules = reactive<FormRules>({
  title: [
    { required: true, message: '请输入导航标题', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  url: [
    { required: true, message: '请输入链接地址', trigger: 'blur' },
    { 
      pattern: /^https?:\/\/.+/, 
      message: '请输入有效的链接地址', 
      trigger: 'blur' 
    }
  ],
  group: [
    { required: true, message: '请选择或输入分组', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ]
})

const dialogTitle = computed(() => {
  return isEdit.value ? '编辑导航' : '新建导航'
})

const showModal = (type: string, row?: any) => {
  dialogVisible.value = true
  isEdit.value = type === 'edit'
  resetForm()
  
  if (row && type === 'edit') {
    nextTick(() => {
      // 数据回显，确保字段正确映射
      form.id = row.id
      form.title = row.title
      form.url = row.url
      form.icon = row.icon
      form.group = row.group
      form.sort = row.sort || row.order_num  // 兼容不同字段名
      form.isEnable = row.isEnable || row.status === 1  // 兼容不同字段名
      form.openInNewTab = row.openInNewTab || row.open_in_new_tab  // 兼容不同字段名
      form.description = row.description
    })
  }
}

const resetForm = () => {
  formRef.value?.resetFields()
  Object.assign(form, {
    id: 0,
    title: '',
    url: '',
    icon: '',
    group: '',
    sort: 1,
    isEnable: true,
    openInNewTab: true,
    description: ''
  })
}

const submitForm = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    // 字段映射：将前端字段名转换为后端期望的字段名
    const formData = {
      name: form.title,                    // 前端的 title -> 后端的 name
      links: form.url,                     // 前端的 url -> 后端的 links
      group_name: form.group,              // 前端的 group -> 后端的 group_name
      icon_url: form.icon || '',           // 前端的 icon -> 后端的 icon_url
      order_num: form.sort || 1,           // 前端的 sort -> 后端的 order_num
      status: form.isEnable ? 1 : 2,       // 前端的 isEnable -> 后端的 status (1:启用, 2:禁用)
      open_in_new_tab: form.openInNewTab,  // 前端的 openInNewTab -> 后端的 open_in_new_tab
      description: form.description || ''
    }
    
    // 编辑时需要ID在URL中，新增时不需要ID
    if (isEdit.value) {
      // 编辑时，将ID添加到formData中用于API URL构建
      (formData as any).id = form.id
      const res = await updateNav(formData)
      if (res.code === ApiStatus.success) {
        ElMessage.success('编辑成功')
        dialogVisible.value = false
        emit('refresh')
      } else {
        ElMessage.error(`编辑失败: ${res.message}`)
      }
    } else {
      const res = await addNav(formData)
      if (res.code === ApiStatus.success) {
        ElMessage.success('新增成功')
        dialogVisible.value = false
        emit('refresh')
      } else {
        ElMessage.error(`新增失败: ${res.message}`)
      }
    }
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 对外暴露方法
defineExpose({
  showModal
})

// 定义事件
const emit = defineEmits(['refresh'])
</script>

<style lang="scss" scoped>
/* 现代化弹窗遮罩层 */
.modern-edit-modal-mask {
  position: fixed;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(8px);
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modern-edit-modal-wrap {
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

/* 现代化弹窗容器 */
.modern-edit-modal {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  box-shadow: 
    0 20px 60px rgba(0, 0, 0, 0.12),
    0 8px 32px rgba(0, 0, 0, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
  width: 800px;
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
.modern-edit-modal-header {
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

/* 主体内容 */
.modern-edit-modal-body {
  flex: 1;
  padding: 24px 32px;
  overflow-y: auto;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.8), rgba(246, 248, 250, 0.5));
}

/* 现代化表单 */
.modern-form {
  .form-section {
    margin-bottom: 32px;
    
    &:last-child {
      margin-bottom: 0;
    }
  }
  
  .section-title {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 18px;
    font-weight: 600;
    color: #1a1a1a;
    margin-bottom: 20px;
    padding-bottom: 12px;
    border-bottom: 2px solid rgba(22, 119, 255, 0.1);
    
    .section-icon {
      width: 20px;
      height: 20px;
      stroke: #1677ff;
      stroke-width: 2;
    }
  }
  
  .form-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
    
    &.switch-grid {
      gap: 16px;
    }
  }
  
  .form-item-wrapper {
    &.full-width {
      grid-column: 1 / -1;
    }
  }
}

/* 表单项样式 */
:deep(.modern-form-item) {
  .el-form-item__label {
    font-weight: 500;
    color: #333;
    line-height: 1.6;
    margin-bottom: 8px;
  }
  
  .el-form-item__content {
    line-height: 1.6;
  }
  
  .el-form-item__error {
    font-size: 12px;
    color: #ff4d4f;
    margin-top: 4px;
  }
}

/* 输入框样式 */
:deep(.modern-input) {
  .el-input__wrapper {
    border-radius: 12px;
    border: 1px solid rgba(0, 0, 0, 0.1);
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(10px);
    transition: all 0.3s ease;
    padding: 12px 16px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    
    &:hover {
      border-color: #1677ff;
      box-shadow: 0 4px 16px rgba(22, 119, 255, 0.12);
    }
    
    &.is-focus {
      border-color: #1677ff;
      box-shadow: 0 0 0 3px rgba(22, 119, 255, 0.1);
    }
  }
  
  .el-input__prefix {
    svg {
      width: 16px;
      height: 16px;
      stroke: #666;
    }
  }
  
  .el-input__inner {
    font-size: 14px;
    color: #333;
    
    &::placeholder {
      color: #999;
    }
  }
}

/* 选择器样式 */
:deep(.modern-select) {
  .el-select__wrapper {
    border-radius: 12px;
    border: 1px solid rgba(0, 0, 0, 0.1);
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(10px);
    transition: all 0.3s ease;
    padding: 12px 16px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    
    &:hover {
      border-color: #1677ff;
      box-shadow: 0 4px 16px rgba(22, 119, 255, 0.12);
    }
    
    &.is-focus {
      border-color: #1677ff;
      box-shadow: 0 0 0 3px rgba(22, 119, 255, 0.1);
    }
  }
}

/* 数字输入框样式 */
:deep(.modern-input-number) {
  width: 100%;
  
  .el-input-number__decrease,
  .el-input-number__increase {
    border: none;
    background: rgba(22, 119, 255, 0.1);
    border-radius: 8px;
    
    &:hover {
      background: rgba(22, 119, 255, 0.2);
    }
    
    svg {
      width: 14px;
      height: 14px;
      stroke: #1677ff;
    }
  }
  
  .el-input__wrapper {
    border-radius: 12px;
    border: 1px solid rgba(0, 0, 0, 0.1);
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(10px);
    transition: all 0.3s ease;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    
    &:hover {
      border-color: #1677ff;
      box-shadow: 0 4px 16px rgba(22, 119, 255, 0.12);
    }
    
    &.is-focus {
      border-color: #1677ff;
      box-shadow: 0 0 0 3px rgba(22, 119, 255, 0.1);
    }
  }
}

/* 文本域样式 */
:deep(.modern-textarea) {
  .el-textarea__inner {
    border-radius: 12px;
    border: 1px solid rgba(0, 0, 0, 0.1);
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(10px);
    transition: all 0.3s ease;
    padding: 16px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    resize: none;
    font-family: inherit;
    line-height: 1.6;
    
    &:hover {
      border-color: #1677ff;
      box-shadow: 0 4px 16px rgba(22, 119, 255, 0.12);
    }
    
    &:focus {
      border-color: #1677ff;
      box-shadow: 0 0 0 3px rgba(22, 119, 255, 0.1);
    }
    
    &::placeholder {
      color: #999;
    }
  }
}

/* 图标选择器样式 */
.icon-selector-wrapper {
  .modern-icon-selector {
    border-radius: 12px;
    overflow: hidden;
  }
  
  .icon-tip {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 8px;
    padding: 8px 12px;
    background: rgba(22, 119, 255, 0.05);
    border-radius: 8px;
    font-size: 12px;
    color: #666;
    line-height: 1.4;
    
    svg {
      width: 14px;
      height: 14px;
      stroke: #1677ff;
      flex-shrink: 0;
    }
  }
}

/* 开关样式 */
.switch-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 12px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  
  &:hover {
    background: rgba(255, 255, 255, 0.8);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  }
}

.switch-label {
  display: flex;
  align-items: center;
  gap: 12px;
  
  svg {
    width: 20px;
    height: 20px;
    stroke: #1677ff;
    stroke-width: 2;
    flex-shrink: 0;
  }
}

.switch-title {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 2px;
}

.switch-desc {
  font-size: 12px;
  color: #666;
  line-height: 1.4;
}

:deep(.modern-switch) {
  .el-switch__core {
    border-radius: 12px;
    height: 24px;
    
    &::after {
      width: 20px;
      height: 20px;
    }
  }
}

/* 底部操作栏 */
.modern-edit-modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 32px;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
  background: rgba(246, 248, 250, 0.5);
}

.footer-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  
  svg {
    width: 16px;
    height: 16px;
  }
  
  &.cancel-btn {
    background: rgba(0, 0, 0, 0.06);
    color: #666;
    
    &:hover {
      background: rgba(0, 0, 0, 0.1);
      transform: translateY(-1px);
    }
  }
  
  &.submit-btn {
    background: linear-gradient(135deg, #1677ff, #69b1ff);
    color: #fff;
    box-shadow: 0 4px 16px rgba(22, 119, 255, 0.3);
    
    &:hover {
      box-shadow: 0 6px 24px rgba(22, 119, 255, 0.4);
      transform: translateY(-2px);
    }
    
    &:active {
      transform: translateY(0);
    }
  }
}

/* 滚动条样式 */
.modern-edit-modal-body::-webkit-scrollbar {
  width: 6px;
}

.modern-edit-modal-body::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.04);
  border-radius: 3px;
}

.modern-edit-modal-body::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
  
  &:hover {
    background: rgba(0, 0, 0, 0.3);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .modern-edit-modal {
    width: 95vw;
    max-height: 95vh;
  }
  
  .modern-edit-modal-header {
    padding: 20px 24px 16px;
  }
  
  .modern-edit-modal-body {
    padding: 20px 24px;
  }
  
  .modern-edit-modal-footer {
    padding: 16px 24px;
    flex-direction: column-reverse;
    gap: 8px;
    
    .footer-btn {
      width: 100%;
      justify-content: center;
    }
  }
  
  .form-grid {
    grid-template-columns: 1fr !important;
    gap: 16px;
  }
  
  .switch-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
}
</style> 