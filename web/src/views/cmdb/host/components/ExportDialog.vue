<template>
  <el-dialog
    title="导出主机数据"
    v-model="visible"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px">
      <!-- 导出范围 -->
      <el-form-item label="导出范围" prop="scope">
        <el-radio-group v-model="formData.scope" @change="handleScopeChange">
          <el-radio value="all">全部主机</el-radio>
          <el-radio value="group">指定主机组</el-radio>
          <el-radio value="current">当前筛选结果</el-radio>
        </el-radio-group>
      </el-form-item>

      <!-- 主机组选择 -->
      <el-form-item 
        v-if="formData.scope === 'group'" 
        label="选择主机组" 
        prop="group_id"
      >
        <el-tree-select
          v-model="formData.group_id"
          :data="hostGroupOptions"
          placeholder="请选择要导出的主机组"
          clearable
          :default-expand-all="true"
          :render-after-expand="false"
          node-key="id"
          :props="{ label: 'name', value: 'id' }"
        />
      </el-form-item>

      <!-- 导出格式 -->
      <el-form-item label="导出格式" prop="format">
        <el-radio-group v-model="formData.format">
          <el-radio value="excel">Excel (.xlsx)</el-radio>
          <el-radio value="csv">CSV (.csv)</el-radio>
        </el-radio-group>
      </el-form-item>

      <!-- 导出字段选择 -->
      <el-form-item label="导出字段" prop="fields">
        <div class="field-selection">
          <div class="field-group">
            <div class="field-group-header">
              <el-checkbox 
                v-model="selectAllBasic" 
                @change="handleSelectAllBasic"
                :indeterminate="basicIndeterminate"
              >
                基础信息
              </el-checkbox>
            </div>
            <div class="field-list">
              <el-checkbox-group v-model="formData.fields">
                <el-checkbox value="name">主机名称</el-checkbox>
                <el-checkbox value="instance_id">实例ID</el-checkbox>
                <el-checkbox value="status">状态</el-checkbox>
                <el-checkbox value="public_ip">公网IP</el-checkbox>
                <el-checkbox value="private_ip">私网IP</el-checkbox>
                <el-checkbox value="os">操作系统</el-checkbox>
                <el-checkbox value="region">区域</el-checkbox>
              </el-checkbox-group>
            </div>
          </div>

          <div class="field-group">
            <div class="field-group-header">
              <el-checkbox 
                v-model="selectAllConfig" 
                @change="handleSelectAllConfig"
                :indeterminate="configIndeterminate"
              >
                配置信息
              </el-checkbox>
            </div>
            <div class="field-list">
              <el-checkbox-group v-model="formData.fields">
                <el-checkbox value="configuration">配置规格</el-checkbox>
                <el-checkbox value="username">用户名</el-checkbox>
                <el-checkbox value="provider_type">提供商类型</el-checkbox>
                <el-checkbox value="resource_type">资源类型</el-checkbox>
                <el-checkbox value="group_name">主机组</el-checkbox>
                <el-checkbox value="provider_name">云账号</el-checkbox>
              </el-checkbox-group>
            </div>
          </div>

          <div class="field-group">
            <div class="field-group-header">
              <el-checkbox 
                v-model="selectAllExtra" 
                @change="handleSelectAllExtra"
                :indeterminate="extraIndeterminate"
              >
                扩展信息
              </el-checkbox>
            </div>
            <div class="field-list">
              <el-checkbox-group v-model="formData.fields">
                <el-checkbox value="tags">标签</el-checkbox>
                <el-checkbox value="expired_at">过期时间</el-checkbox>
                <el-checkbox value="remark">备注</el-checkbox>
                <el-checkbox value="created_at">创建时间</el-checkbox>
                <el-checkbox value="updated_at">更新时间</el-checkbox>
              </el-checkbox-group>
            </div>
          </div>
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleExport" :loading="exporting">
          {{ exporting ? '导出中...' : '确认导出' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getHostGroupTree } from '@/api/cmdb/hostGroup'
import type { HostGroup } from '@/api/cmdb/hostGroup'

interface ExportFormData {
  scope: 'all' | 'group' | 'current'
  group_id?: number
  format: 'excel' | 'csv'
  fields: string[]
}

interface Props {
  modelValue: boolean
  currentFilters?: Record<string, any>
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'export', data: ExportFormData): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const formRef = ref<FormInstance>()
const exporting = ref(false)
const hostGroupOptions = ref<HostGroup[]>([])

// 表单数据
const formData = reactive<ExportFormData>({
  scope: 'all',
  format: 'excel',
  fields: ['name', 'instance_id', 'status', 'public_ip', 'private_ip', 'os', 'region']
})

// 表单验证规则
const rules: FormRules = {
  scope: [{ required: true, message: '请选择导出范围', trigger: 'change' }],
  group_id: [
    { 
      required: true, 
      message: '请选择主机组', 
      trigger: 'change',
      validator: (rule, value, callback) => {
        if (formData.scope === 'group' && !value) {
          callback(new Error('请选择主机组'))
        } else {
          callback()
        }
      }
    }
  ],
  format: [{ required: true, message: '请选择导出格式', trigger: 'change' }],
  fields: [
    { 
      required: true, 
      message: '请至少选择一个导出字段', 
      trigger: 'change',
      validator: (rule, value, callback) => {
        if (!value || value.length === 0) {
          callback(new Error('请至少选择一个导出字段'))
        } else {
          callback()
        }
      }
    }
  ]
}

// 字段分组
const basicFields = ['name', 'instance_id', 'status', 'public_ip', 'private_ip', 'os', 'region']
const configFields = ['configuration', 'username', 'provider_type', 'resource_type', 'group_name', 'provider_name']
const extraFields = ['tags', 'expired_at', 'remark', 'created_at', 'updated_at']

// 全选状态计算
const selectAllBasic = computed({
  get: () => basicFields.every(field => formData.fields.includes(field)),
  set: (value) => {
    if (value) {
      basicFields.forEach(field => {
        if (!formData.fields.includes(field)) {
          formData.fields.push(field)
        }
      })
    } else {
      formData.fields = formData.fields.filter(field => !basicFields.includes(field))
    }
  }
})

const selectAllConfig = computed({
  get: () => configFields.every(field => formData.fields.includes(field)),
  set: (value) => {
    if (value) {
      configFields.forEach(field => {
        if (!formData.fields.includes(field)) {
          formData.fields.push(field)
        }
      })
    } else {
      formData.fields = formData.fields.filter(field => !configFields.includes(field))
    }
  }
})

const selectAllExtra = computed({
  get: () => extraFields.every(field => formData.fields.includes(field)),
  set: (value) => {
    if (value) {
      extraFields.forEach(field => {
        if (!formData.fields.includes(field)) {
          formData.fields.push(field)
        }
      })
    } else {
      formData.fields = formData.fields.filter(field => !extraFields.includes(field))
    }
  }
})

// 半选状态计算
const basicIndeterminate = computed(() => {
  const selectedBasic = basicFields.filter(field => formData.fields.includes(field))
  return selectedBasic.length > 0 && selectedBasic.length < basicFields.length
})

const configIndeterminate = computed(() => {
  const selectedConfig = configFields.filter(field => formData.fields.includes(field))
  return selectedConfig.length > 0 && selectedConfig.length < configFields.length
})

const extraIndeterminate = computed(() => {
  const selectedExtra = extraFields.filter(field => formData.fields.includes(field))
  return selectedExtra.length > 0 && selectedExtra.length < extraFields.length
})

// 处理范围变化
const handleScopeChange = () => {
  if (formData.scope !== 'group') {
    formData.group_id = undefined
  }
}

// 处理全选
const handleSelectAllBasic = (value: boolean) => {
  selectAllBasic.value = value
}

const handleSelectAllConfig = (value: boolean) => {
  selectAllConfig.value = value
}

const handleSelectAllExtra = (value: boolean) => {
  selectAllExtra.value = value
}

// 获取主机组列表
const fetchHostGroups = async () => {
  try {
    const response = await getHostGroupTree()
    hostGroupOptions.value = response.data || []
  } catch (error) {
    console.error('获取主机组列表失败:', error)
    ElMessage.error('获取主机组列表失败')
  }
}

// 处理导出
const handleExport = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    
    exporting.value = true
    emit('export', { ...formData })
    
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

// 处理关闭
const handleClose = () => {
  visible.value = false
  exporting.value = false
}

// 监听对话框打开
watch(visible, (newVal) => {
  if (newVal) {
    fetchHostGroups()
  }
})

// 暴露方法给父组件
defineExpose({
  setExporting: (value: boolean) => {
    exporting.value = value
  }
})
</script>

<style scoped>
.field-selection {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 16px;
  background-color: #fafafa;
}

.field-group {
  margin-bottom: 16px;
}

.field-group:last-child {
  margin-bottom: 0;
}

.field-group-header {
  margin-bottom: 8px;
  font-weight: 500;
  color: #303133;
}

.field-list {
  margin-left: 20px;
}

.field-list .el-checkbox-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 16px;
}

.field-list .el-checkbox {
  margin-right: 0;
  white-space: nowrap;
}

.dialog-footer {
  text-align: right;
}
</style>
