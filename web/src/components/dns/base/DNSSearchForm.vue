<template>
  <el-card shadow="never" class="search-form-card">
    <el-form
      ref="formRef"
      :model="searchForm"
      inline
      class="search-form"
      @submit.prevent="handleSearch"
    >
      <!-- 动态渲染搜索字段 -->
      <el-form-item
        v-for="field in fields"
        :key="field.key"
        :label="field.label"
        :prop="field.key"
      >
        <!-- 输入框 -->
        <el-input
          v-if="field.type === 'input'"
          v-model="searchForm[field.key]"
          :placeholder="field.placeholder"
          clearable
          :style="{ width: field.width || '200px' }"
        />
        
        <!-- 选择框 -->
        <el-select
          v-else-if="field.type === 'select'"
          v-model="searchForm[field.key]"
          :placeholder="field.placeholder"
          clearable
          :style="{ width: field.width || '150px' }"
        >
          <el-option
            v-for="option in field.options"
            :key="option.value"
            :label="option.label"
            :value="option.value"
          />
        </el-select>
        
        <!-- 日期选择器 -->
        <el-date-picker
          v-else-if="field.type === 'date'"
          v-model="searchForm[field.key]"
          type="date"
          :placeholder="field.placeholder"
          clearable
          :style="{ width: field.width || '150px' }"
        />
        
        <!-- 日期范围选择器 -->
        <el-date-picker
          v-else-if="field.type === 'daterange'"
          v-model="searchForm[field.key]"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          clearable
          :style="{ width: field.width || '240px' }"
        />
      </el-form-item>
      
      <!-- 操作按钮 -->
      <el-form-item>
        <el-button type="primary" @click="handleSearch" :loading="loading">
          <el-icon><Search /></el-icon>
          搜索
        </el-button>
        <el-button @click="handleReset">
          <el-icon><Refresh /></el-icon>
          重置
        </el-button>
        <el-button 
          v-if="showExport" 
          @click="handleExport"
          :loading="exportLoading"
        >
          <el-icon><Download /></el-icon>
          导出
        </el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { Search, Refresh, Download } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'

export interface SearchField {
  /** 字段键名 */
  key: string
  /** 字段标签 */
  label: string
  /** 字段类型 */
  type: 'input' | 'select' | 'date' | 'daterange'
  /** 占位符 */
  placeholder?: string
  /** 组件宽度 */
  width?: string
  /** 选择框选项 */
  options?: Array<{ label: string; value: any }>
}

export interface Props {
  /** 搜索字段配置 */
  fields: SearchField[]
  /** 初始值 */
  modelValue?: Record<string, any>
  /** 是否显示导出按钮 */
  showExport?: boolean
  /** 搜索加载状态 */
  loading?: boolean
  /** 导出加载状态 */
  exportLoading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => ({}),
  showExport: false,
  loading: false,
  exportLoading: false
})

const emit = defineEmits<{
  'update:modelValue': [value: Record<string, any>]
  'search': [params: Record<string, any>]
  'reset': []
  'export': [params: Record<string, any>]
}>()

// 表单引用
const formRef = ref<FormInstance>()

// 搜索表单数据
const searchForm = reactive<Record<string, any>>({})

// 初始化表单数据
const initFormData = () => {
  props.fields.forEach(field => {
    if (props.modelValue[field.key] !== undefined) {
      searchForm[field.key] = props.modelValue[field.key]
    } else {
      // 设置默认值
      switch (field.type) {
        case 'input':
          searchForm[field.key] = ''
          break
        case 'select':
          searchForm[field.key] = ''
          break
        case 'date':
          searchForm[field.key] = ''
          break
        case 'daterange':
          searchForm[field.key] = []
          break
        default:
          searchForm[field.key] = ''
      }
    }
  })
}

// 监听modelValue变化
watch(() => props.modelValue, (newVal) => {
  Object.assign(searchForm, newVal)
}, { immediate: true })

// 监听searchForm变化，同步到父组件
watch(searchForm, (newVal) => {
  emit('update:modelValue', { ...newVal })
}, { deep: true })

// 搜索处理
const handleSearch = () => {
  emit('search', { ...searchForm })
}

// 重置处理
const handleReset = () => {
  formRef.value?.resetFields()
  initFormData()
  emit('reset')
}

// 导出处理
const handleExport = () => {
  emit('export', { ...searchForm })
}

// 初始化
initFormData()
</script>

<style scoped>
.search-form-card {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: -18px;
}

.search-form .el-form-item {
  margin-bottom: 18px;
}

.search-form .el-form-item:last-child {
  margin-right: 0;
}
</style>
