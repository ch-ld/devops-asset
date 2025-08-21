<template>
  <el-dialog
    v-model="dialogVisible"
    :title="modalTitle"
    width="680px"
    :close-on-click-modal="false"
    @close="handleClose"
    class="record-modal"
    align-center
  >
    <template #header="{ titleId, titleClass }">
      <div class="modal-header">
        <div class="header-icon">
          <el-icon size="24"><Plus /></el-icon>
        </div>
        <div class="header-content">
          <h3 :id="titleId" :class="titleClass" class="modal-title">{{ modalTitle }}</h3>
          <p class="modal-subtitle">为域名 {{ domainName }} 添加DNS解析记录</p>
        </div>
      </div>
    </template>

    <div class="modal-body">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        @submit.prevent
        class="record-form"
        label-position="top"
      >
        <!-- 记录类型 -->
        <el-form-item label="记录类型" prop="type" class="form-item-type">
          <el-select
            v-model="formData.type"
            placeholder="选择记录类型"
            size="large"
            @change="handleTypeChange"
            class="type-select"
          >
            <el-option-group label="常用类型">
              <el-option value="A">
                <div class="option-content">
                  <div class="option-main">
                    <span class="option-type">A</span>
                    <span class="option-desc">将域名指向IPv4地址</span>
                  </div>
                  <span class="option-example">如：192.168.1.1</span>
                </div>
              </el-option>
              <el-option value="CNAME">
                <div class="option-content">
                  <div class="option-main">
                    <span class="option-type">CNAME</span>
                    <span class="option-desc">将域名指向另一个域名</span>
                  </div>
                  <span class="option-example">如：www.example.com</span>
                </div>
              </el-option>
              <el-option value="MX">
                <div class="option-content">
                  <div class="option-main">
                    <span class="option-type">MX</span>
                    <span class="option-desc">邮件服务器记录</span>
                  </div>
                  <span class="option-example">如：mail.example.com</span>
                </div>
              </el-option>
            </el-option-group>
            <el-option-group label="其他类型">
              <el-option value="AAAA" label="AAAA - IPv6地址记录" />
              <el-option value="TXT" label="TXT - 文本记录" />
              <el-option value="NS" label="NS - 域名服务器记录" />
              <el-option value="SRV" label="SRV - 服务记录" />
            </el-option-group>
          </el-select>
        </el-form-item>

        <!-- 主机记录 -->
        <el-form-item label="主机记录" prop="name" class="form-item-name">
          <el-input
            v-model="formData.name"
            placeholder="请输入主机记录"
            size="large"
            class="name-input"
          >
            <template #suffix>
              <span class="domain-suffix">.{{ domainName }}</span>
            </template>
          </el-input>
          <div class="form-help">
            <el-icon><InfoFilled /></el-icon>
            <span>{{ getRecordHelp() }}</span>
          </div>
        </el-form-item>

        <!-- 解析线路 -->
        <el-form-item label="解析线路" prop="line" class="form-item-line">
          <el-select
            v-model="formData.line"
            placeholder="默认"
            size="large"
            class="line-select"
          >
            <el-option label="默认" value="default" />
            <el-option label="境内" value="cn" />
            <el-option label="境外" value="abroad" />
          </el-select>
          <div class="form-help">
            <el-icon><InfoFilled /></el-icon>
            <span>升级至企业版后域名支持按地域、运营商分别设置解析</span>
          </div>
        </el-form-item>

        <!-- 记录值 -->
        <el-form-item label="记录值" prop="value" class="form-item-value">
          <el-input
            v-model="formData.value"
            :placeholder="getValuePlaceholder()"
            size="large"
            :type="formData.type === 'TXT' ? 'textarea' : 'text'"
            :rows="3"
            class="value-input"
          />
          <div class="form-help">
            <el-icon><InfoFilled /></el-icon>
            <span>{{ getValueHelp() }}</span>
          </div>
        </el-form-item>

        <!-- TTL时间 -->
        <el-form-item label="TTL" prop="ttl" class="form-item-ttl">
          <el-select
            v-model="formData.ttl"
            size="large"
            class="ttl-select"
          >
            <el-option label="10秒 (测试用)" :value="10" />
            <el-option label="1分钟" :value="60" />
            <el-option label="10分钟 (推荐)" :value="600" />
            <el-option label="1小时" :value="3600" />
            <el-option label="12小时" :value="43200" />
            <el-option label="1天" :value="86400" />
          </el-select>
          <div class="form-help">
            <el-icon><InfoFilled /></el-icon>
            <span>缓存生存时间，数值越小，修改记录各地生效时间越快</span>
          </div>
        </el-form-item>

        <!-- MX记录优先级 -->
        <el-form-item
          v-if="formData.type === 'MX'"
          label="MX优先级"
          prop="priority"
        >
          <el-input-number
            v-model="formData.priority"
            :min="1"
            :max="50"
            style="width: 100%"
          />
        </el-form-item>

        <!-- SRV记录参数 -->
        <template v-if="formData.type === 'SRV'">
          <el-form-item label="优先级" prop="priority">
            <el-input-number
              v-model="formData.priority"
              :min="0"
              :max="65535"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="权重" prop="weight">
            <el-input-number
              v-model="formData.weight"
              :min="0"
              :max="65535"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="端口" prop="port">
            <el-input-number
              v-model="formData.port"
              :min="1"
              :max="65535"
              style="width: 100%"
            />
          </el-form-item>
        </template>
      </el-form>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">
          {{ isEdit ? '更新' : '确定' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, InfoFilled } from '@element-plus/icons-vue'
import { recordApi } from '@/api/dns/record'
import type { DNSRecord } from '@/types/dns'

interface Props {
  visible: boolean
  record?: DNSRecord | null
  domainId?: number
  domainName?: string
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  record: null
})

const emit = defineEmits<Emits>()

const formRef = ref()
const loading = ref(false)

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const isEdit = computed(() => !!props.record)
const modalTitle = computed(() => isEdit.value ? '编辑记录' : '添加记录')



const formData = ref({
  domain_id: props.domainId || 0,
  name: '',
  type: 'A',
  value: '',
  ttl: 600,
  priority: 10,
  weight: 1,
  port: 80,
  line: 'default'
})

const formRules = {
  type: [{ required: true, message: '请选择记录类型', trigger: 'change' }],
  name: [{ required: true, message: '请输入主机记录', trigger: 'blur' }],
  value: [{ required: true, message: '请输入记录值', trigger: 'blur' }],
  ttl: [{ required: true, message: '请选择TTL', trigger: 'change' }]
}

const handleTypeChange = () => {
  // 根据类型设置默认值
  switch (formData.value.type) {
    case 'MX':
      formData.value.priority = 10
      break
    case 'SRV':
      formData.value.priority = 0
      formData.value.weight = 1
      formData.value.port = 80
      break
  }
}

const getRecordHelp = () => {
  switch (formData.value.type) {
    case 'A':
    case 'AAAA':
      return '如需要将域名指向一个IP地址，可以使用此类型。例如：www、mail、@'
    case 'CNAME':
      return '如需要将域名指向另一个域名，再由另一个域名提供IP地址，可以使用此类型'
    case 'MX':
      return '如需要设置邮箱，让邮箱能收到邮件，可以使用此类型'
    case 'TXT':
      return '如需要对域名进行标识和说明，可以使用此类型'
    case 'NS':
      return '如需要把子域名交给其他DNS服务商解析，可以使用此类型'
    case 'SRV':
      return '如需要记录提供特定服务的服务器，可以使用此类型'
    default:
      return ''
  }
}

const getValuePlaceholder = () => {
  switch (formData.value.type) {
    case 'A':
      return '请输入IPv4地址，例如：192.168.1.1'
    case 'AAAA':
      return '请输入IPv6地址，例如：2001:db8::1'
    case 'CNAME':
      return '请输入域名，例如：www.example.com'
    case 'MX':
      return '请输入邮件服务器地址，例如：mail.example.com'
    case 'TXT':
      return '请输入文本内容'
    case 'NS':
      return '请输入DNS服务器地址，例如：ns1.example.com'
    case 'SRV':
      return '请输入目标地址，例如：target.example.com'
    default:
      return '请输入记录值'
  }
}

const getValueHelp = () => {
  switch (formData.value.type) {
    case 'A':
      return '将域名指向一个IPv4地址，例如：223.5.5.5'
    case 'AAAA':
      return '将域名指向一个IPv6地址'
    case 'CNAME':
      return '将域名指向另一个域名，可以通过该域名获取IP地址'
    case 'MX':
      return '将域名指向邮件服务器地址，用于邮件收发'
    case 'TXT':
      return '为域名设置说明文字，常用于域名验证'
    case 'NS':
      return '将子域名指定其他DNS服务器解析'
    case 'SRV':
      return '记录提供特定服务的服务器'
    default:
      return ''
  }
}
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  loading.value = true
  try {
    const data = { ...formData.value }

    // 如果没有传入domainId，使用表单中的domain_id
    if (!props.domainId) {
      data.domain_id = formData.value.domain_id
    } else {
      data.domain_id = props.domainId
    }

    if (isEdit.value && props.record) {
      await recordApi.update(props.record.id, data)
      ElMessage.success('记录更新成功')
    } else {
      await recordApi.create(data)
      ElMessage.success('记录创建成功')
    }

    emit('success')
    handleClose()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}

const handleClose = () => {
  dialogVisible.value = false
  resetForm()
}

const resetForm = () => {
  formData.value = {
    domain_id: props.domainId || 0,
    name: '',
    type: 'A',
    value: '',
    ttl: 600,
    priority: 10,
    weight: 1,
    port: 80,
    line: 'default'
  }
  nextTick(() => {
    formRef.value?.clearValidate()
  })
}

// 监听props变化
watch(() => props.record, (newRecord) => {
  if (newRecord) {
    formData.value = {
      domain_id: newRecord.domain_id,
      name: newRecord.name,
      type: newRecord.type,
      value: newRecord.value,
      ttl: newRecord.ttl,
      priority: newRecord.priority || 10,
      weight: newRecord.weight || 1,
      port: newRecord.port || 80,
      line: newRecord.line || 'default'
    }
  } else {
    resetForm()
  }
}, { immediate: true })

watch(() => props.visible, (visible) => {
  if (!visible) {
    resetForm()
  }
})
</script>

<style scoped lang="scss">
.record-modal {
  :deep(.el-dialog) {
    border-radius: 8px;
    overflow: hidden;
    background: white;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    border: 1px solid #e1e4e8;
  }

  :deep(.el-dialog__header) {
    padding: 0;
    border-bottom: 1px solid rgba(226, 232, 240, 0.3);
  }

  :deep(.el-dialog__body) {
    padding: 32px;
    background: transparent;
  }

  :deep(.el-dialog__footer) {
    padding: 16px 24px;
    border-top: 1px solid #e1e4e8;
    background: #f6f8fa;
  }

  .modal-header {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 24px;
    background: #0366d6;
    color: white;
    border-bottom: 1px solid #e1e4e8;

    .header-icon {
      width: 40px;
      height: 40px;
      border-radius: 6px;
      background: rgba(255, 255, 255, 0.15);
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 20px;
    }

    .header-content {
      flex: 1;

      .modal-title {
        margin: 0 0 4px 0;
        font-size: 18px;
        font-weight: 600;
        color: white;
      }

      .modal-subtitle {
        margin: 0;
        font-size: 14px;
        color: rgba(255, 255, 255, 0.8);
        font-weight: 400;
      }
    }
  }

  .modal-body {
    .record-form {
      :deep(.el-form-item) {
        margin-bottom: 28px;

        .el-form-item__label {
          font-weight: 700;
          color: #1e293b;
          margin-bottom: 12px;
          line-height: 1.4;
          font-size: 15px;
          background: linear-gradient(135deg, #1e293b 0%, #475569 100%);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
        }

        .el-form-item__content {
          line-height: 1.4;
        }
      }

      .form-item-type {
        .type-select {
          :deep(.el-select__wrapper) {
            border-radius: 12px;
            border: 2px solid #e2e8f0;
            background: rgba(255, 255, 255, 0.9);
            backdrop-filter: blur(10px);
            transition: all 0.3s ease;
            padding: 16px 20px;
            font-size: 15px;

            &:hover {
              border-color: #94a3b8;
              transform: translateY(-1px);
              box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
            }

            &.is-focused {
              border-color: #667eea;
              box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.15);
              transform: translateY(-1px);
            }
          }
        }

        .option-content {
          display: flex;
          justify-content: space-between;
          align-items: center;
          width: 100%;
          padding: 8px 0;

          .option-main {
            display: flex;
            align-items: center;
            gap: 16px;

            .option-type {
              font-weight: 700;
              background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
              -webkit-background-clip: text;
              -webkit-text-fill-color: transparent;
              background-clip: text;
              min-width: 70px;
              font-size: 16px;
            }

            .option-desc {
              color: #64748b;
              font-weight: 500;
            }
          }

          .option-example {
            font-size: 13px;
            color: #94a3b8;
            font-weight: 500;
          }
        }
      }

      .form-item-name,
      .form-item-value {
        .name-input,
        .value-input {
          :deep(.el-input__wrapper) {
            border-radius: 12px;
            border: 2px solid #e2e8f0;
            background: rgba(255, 255, 255, 0.9);
            backdrop-filter: blur(10px);
            transition: all 0.3s ease;
            padding: 16px 20px;
            font-size: 15px;

            &:hover {
              border-color: #94a3b8;
              transform: translateY(-1px);
              box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
            }

            &.is-focus {
              border-color: #667eea;
              box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.15);
              transform: translateY(-1px);
            }
          }

          :deep(.el-textarea__inner) {
            border-radius: 12px;
            border: 2px solid #e2e8f0;
            background: rgba(255, 255, 255, 0.9);
            backdrop-filter: blur(10px);
            transition: all 0.3s ease;
            padding: 16px 20px;
            font-size: 15px;
            line-height: 1.6;

            &:hover {
              border-color: #94a3b8;
              box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
            }

            &:focus {
              border-color: #667eea;
              box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.15);
            }
          }
        }
      }

      .form-help {
        display: flex;
        align-items: flex-start;
        gap: 8px;
        font-size: 14px;
        color: #64748b;
        margin-top: 12px;
        line-height: 1.5;
        padding: 12px 16px;
        background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
        border-radius: 10px;
        border: 1px solid rgba(226, 232, 240, 0.5);

        .el-icon {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
          font-size: 16px;
          margin-top: 2px;
          flex-shrink: 0;
        }
      }

      .domain-suffix {
        color: #64748b;
        font-size: 15px;
        font-weight: 600;
        background: rgba(102, 126, 234, 0.1);
        padding: 4px 8px;
        border-radius: 6px;
      }

      .form-item-line,
      .form-item-ttl,
      .form-item-priority {
        margin-bottom: 24px;

        .line-select,
        .ttl-select,
        .priority-input {
          :deep(.el-select__wrapper),
          :deep(.el-input__wrapper) {
            border-radius: 12px;
            border: 2px solid #e2e8f0;
            background: rgba(255, 255, 255, 0.9);
            backdrop-filter: blur(10px);
            transition: all 0.3s ease;
            padding: 16px 20px;
            font-size: 15px;

            &:hover {
              border-color: #94a3b8;
              transform: translateY(-1px);
              box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
            }

            &.is-focused,
            &.is-focus {
              border-color: #667eea;
              box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.15);
              transform: translateY(-1px);
            }
          }
        }
      }

      :deep(.el-input-number) {
        width: 100%;
        
        .el-input__wrapper {
          border-radius: 12px;
          border: 2px solid #e2e8f0;
          background: rgba(255, 255, 255, 0.9);
          backdrop-filter: blur(10px);
          transition: all 0.3s ease;
          padding: 16px 20px;

          &:hover {
            border-color: #94a3b8;
            transform: translateY(-1px);
            box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
          }

          &.is-focus {
            border-color: #667eea;
            box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.15);
            transform: translateY(-1px);
          }
        }
        
        .el-input-number__decrease,
        .el-input-number__increase {
          background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
          border: none;
          color: #64748b;
          transition: all 0.3s ease;
          
          &:hover {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
          }
        }
      }
    }
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 16px;

    :deep(.el-button) {
      border-radius: 12px;
      font-weight: 600;
      padding: 14px 28px;
      font-size: 15px;
      transition: all 0.3s ease;
      border: 2px solid transparent;

      &.el-button--default {
        background: rgba(255, 255, 255, 0.9);
        backdrop-filter: blur(10px);
        border-color: #e2e8f0;
        color: #64748b;
        
        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
          border-color: #94a3b8;
          color: #475569;
        }
      }

      &.el-button--primary {
        background: #0366d6;
        border-color: #0366d6;

        &:hover {
          background: #0256cc;
          border-color: #0256cc;
        }
      }
    }
  }
}
</style>
