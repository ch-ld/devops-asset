<template>
  <el-dialog 
    v-model="visible" 
    :title="isEdit ? '编辑域名' : '添加域名'"
    width="600px"
    class="domain-modal"
    :close-on-click-modal="false"
    destroy-on-close
  >
    <div class="modal-content">
      <el-form 
        ref="formRef" 
        :model="form" 
        :rules="rules" 
        label-width="100px"
        class="domain-form"
      >
        <div class="form-section">
          <h4 class="section-title">
            <el-icon><Monitor /></el-icon>
            基本信息
          </h4>
          
          <el-form-item label="域名" prop="name">
            <el-input 
              v-model="form.name" 
              placeholder="请输入域名，如：example.com"
              :disabled="isEdit"
              clearable
            >
              <template #prepend>
                <el-icon><Monitor /></el-icon>
              </template>
            </el-input>
            <div class="form-tip">域名格式：example.com 或 subdomain.example.com</div>
          </el-form-item>
          
          <el-form-item label="分组" prop="group_id">
            <el-select 
              v-model="form.group_id" 
              placeholder="请选择域名分组"
              clearable
              style="width: 100%"
            >
              <el-option label="无分组" value="" />
              <el-option
                v-for="group in groups"
                :key="group.id"
                :label="group.name"
                :value="group.id"
              />
            </el-select>
            <div class="form-tip">选择域名所属的分组，便于管理</div>
          </el-form-item>
          
          <el-form-item label="注册商" prop="registrar_type">
            <el-select 
              v-model="form.registrar_type" 
              placeholder="请选择注册商"
              style="width: 100%"
            >
              <el-option label="阿里云" value="aliyun">
                <div class="registrar-option">
                  <span>阿里云</span>
                  <el-tag size="small" type="success">推荐</el-tag>
                </div>
              </el-option>
              <el-option label="腾讯云" value="tencent">
                <div class="registrar-option">
                  <span>腾讯云</span>
                </div>
              </el-option>
              <el-option label="AWS Route53" value="route53">
                <div class="registrar-option">
                  <span>AWS Route53</span>
                </div>
              </el-option>
              <el-option label="Cloudflare" value="cloudflare">
                <div class="registrar-option">
                  <span>Cloudflare</span>
                </div>
              </el-option>
              <el-option label="DNSPod" value="dnspod">
                <div class="registrar-option">
                  <span>DNSPod</span>
                </div>
              </el-option>
              <el-option label="GoDaddy" value="godaddy">
                <div class="registrar-option">
                  <span>GoDaddy</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
        </div>
        
        <div class="form-section">
          <h4 class="section-title">
            <el-icon><Setting /></el-icon>
            配置选项
          </h4>
          
          <el-form-item label="状态" prop="status">
            <el-radio-group v-model="form.status" class="status-radio-group">
              <el-radio value="active" class="status-radio active">
                <div class="radio-content">
                  <el-icon><CircleCheck /></el-icon>
                  <span>正常</span>
                </div>
              </el-radio>
              <el-radio value="inactive" class="status-radio inactive">
                <div class="radio-content">
                  <el-icon><CircleClose /></el-icon>
                  <span>禁用</span>
                </div>
              </el-radio>
            </el-radio-group>
          </el-form-item>
          
          <el-form-item label="到期时间" prop="expires_at">
            <el-date-picker
              v-model="form.expires_at"
              type="datetime"
              placeholder="选择到期时间"
              style="width: 100%"
              format="YYYY-MM-DD HH:mm:ss"
              value-format="YYYY-MM-DD HH:mm:ss"
            />
            <div class="form-tip">域名到期时间，用于续费提醒</div>
          </el-form-item>
          
          <el-form-item label="自动续费">
            <el-switch
              v-model="form.auto_renew"
              active-text="启用"
              inactive-text="禁用"
              inline-prompt
            />
            <div class="form-tip">启用后将在到期前自动续费</div>
          </el-form-item>
        </div>
        
        <div class="form-section">
          <h4 class="section-title">
            <el-icon><Document /></el-icon>
            附加信息
          </h4>
          
          <el-form-item label="标签" prop="tags">
            <el-tag
              v-for="tag in form.tags"
              :key="tag"
              closable
              @close="removeTag(tag)"
              class="tag-item"
            >
              {{ tag }}
            </el-tag>
            <el-input
              v-if="inputVisible"
              ref="inputRef"
              v-model="inputValue"
              size="small"
              style="width: 100px"
              @keyup.enter="handleInputConfirm"
              @blur="handleInputConfirm"
            />
            <el-button v-else @click="showInput" size="small" type="primary" text>
              <el-icon><Plus /></el-icon>
              添加标签
            </el-button>
            <div class="form-tip">为域名添加标签，便于分类管理</div>
          </el-form-item>
          
          <el-form-item label="备注" prop="remark">
            <el-input
              v-model="form.remark"
              type="textarea"
              :rows="3"
              placeholder="请输入备注信息..."
            />
          </el-form-item>
        </div>
      </el-form>
    </div>
    
    <template #footer>
      <div class="modal-footer">
        <el-button @click="handleCancel" size="large">取消</el-button>
        <el-button 
          type="primary" 
          @click="handleSubmit" 
          :loading="loading"
          size="large"
          class="submit-btn"
        >
          <el-icon><Check /></el-icon>
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, nextTick } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { 
  Monitor, 
  Setting, 
  Document, 
  CircleCheck, 
  CircleClose,
  Plus,
  Check
} from '@element-plus/icons-vue'
import { domainApi } from '@/api/dns/domain'
import type { Domain } from '@/types/dns'

interface Props {
  visible: boolean
  domain?: Domain | null
  groups?: any[]
}

const props = withDefaults(defineProps<Props>(), {
  domain: null,
  groups: () => []
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'success': []
}>()

// 响应式数据
const formRef = ref<FormInstance>()
const inputRef = ref()
const loading = ref(false)
const inputVisible = ref(false)
const inputValue = ref('')

const visible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const isEdit = computed(() => !!props.domain)

// 表单数据
const form = reactive({
  name: '',
  group_id: '' as string | number,
  registrar_type: '',
  status: 'active' as 'active' | 'inactive' | 'expired',
  expires_at: '',
  auto_renew: false,
  tags: [] as string[],
  remark: ''
})

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: '请输入域名', trigger: 'blur' },
    { 
      pattern: /^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$/, 
      message: '请输入有效的域名格式', 
      trigger: 'blur' 
    }
  ],
  registrar_type: [
    { required: true, message: '请选择注册商', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

// 监听domain变化，初始化表单
watch(() => props.domain, (newDomain) => {
  if (newDomain) {
    Object.assign(form, {
      name: newDomain.name || '',
      group_id: newDomain.group_id || null,
      registrar_type: newDomain.registrar_type || '',
      status: newDomain.status || 'active',
      expires_at: newDomain.expires_at || '',
      auto_renew: newDomain.auto_renew || false,
      tags: newDomain.tags || [],
      remark: newDomain.remark || ''
    })
  } else {
    Object.assign(form, {
      name: '',
      group_id: null,
      registrar_type: '',
      status: 'active',
      expires_at: '',
      auto_renew: false,
      tags: [],
      remark: ''
    })
  }
}, { immediate: true })

// 标签操作
const removeTag = (tag: string) => {
  form.tags.splice(form.tags.indexOf(tag), 1)
}

const showInput = () => {
  inputVisible.value = true
  nextTick(() => {
    inputRef.value?.focus()
  })
}

const handleInputConfirm = () => {
  if (inputValue.value && !form.tags.includes(inputValue.value)) {
    form.tags.push(inputValue.value)
  }
  inputVisible.value = false
  inputValue.value = ''
}

// 事件处理
const handleCancel = () => {
  visible.value = false
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    loading.value = true
    
    const submitData = { ...form }
    
    if (isEdit.value && props.domain) {
      await domainApi.update(props.domain.id, submitData)
      ElMessage.success('域名更新成功')
    } else {
      await domainApi.create(submitData)
      ElMessage.success('域名创建成功')
    }
    
    visible.value = false
    emit('success')
  } catch (error) {
    ElMessage.error(isEdit.value ? '域名更新失败' : '域名创建失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.domain-modal {
  .modal-content {
    max-height: 60vh;
    overflow-y: auto;
  }
  
  .domain-form {
    .form-section {
      margin-bottom: 24px;
      
      .section-title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 600;
        color: #1f2937;
        margin-bottom: 16px;
        padding-bottom: 8px;
        border-bottom: 2px solid #e5e7eb;
        
        .el-icon {
          color: #3b82f6;
        }
      }
      
      .form-tip {
        font-size: 12px;
        color: #6b7280;
        margin-top: 4px;
        line-height: 1.4;
      }
    }
    
    .registrar-option {
      display: flex;
      align-items: center;
      justify-content: space-between;
      width: 100%;
    }
    
    .status-radio-group {
      display: flex;
      gap: 16px;
      
      .status-radio {
        margin-right: 0;
        
        .radio-content {
          display: flex;
          align-items: center;
          gap: 8px;
          padding: 8px 16px;
          border: 2px solid #e5e7eb;
          border-radius: 8px;
          transition: all 0.3s ease;
          
          .el-icon {
            font-size: 18px;
          }
        }
        
        &.active .radio-content {
          .el-icon {
            color: #10b981;
          }
        }
        
        &.inactive .radio-content {
          .el-icon {
            color: #ef4444;
          }
        }
        
        &:deep(.el-radio__input.is-checked) + .radio-content {
          border-color: #3b82f6;
          background: #eff6ff;
        }
      }
    }
    
    .tag-item {
      margin-right: 8px;
      margin-bottom: 8px;
    }
  }
  
  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    
    .submit-btn {
      min-width: 120px;
      
      .el-icon {
        margin-right: 4px;
      }
    }
  }
}

/* 自定义对话框样式 */
:deep(.el-dialog) {
  border-radius: 12px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
}

:deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 12px 12px 0 0;
  padding: 20px 24px;
  
  .el-dialog__title {
    font-size: 18px;
    font-weight: 600;
  }
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-dialog__footer) {
  padding: 20px 24px;
  background: #f9fafb;
  border-radius: 0 0 12px 12px;
}
</style>
