<template>
  <div class="import-cert-page">
    <div class="page-header">
      <h1>导入SSL证书</h1>
      <p>导入已有的SSL证书文件到系统中进行管理</p>
    </div>

    <div class="form-container">
      <el-form :model="formData" label-width="120px" :rules="formRules" ref="formRef">
        <el-form-item label="证书名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入证书名称" />
        </el-form-item>

        <el-form-item label="证书内容" prop="certContent">
          <el-input
            v-model="formData.certContent"
            type="textarea"
            :rows="8"
            placeholder="请粘贴证书内容，以 -----BEGIN CERTIFICATE----- 开头，以 -----END CERTIFICATE----- 结尾"
            style="font-family: monospace;"
          />
          <div class="form-tip">
            支持文本格式：.pem、.crt、.PEM、.CRT
          </div>
        </el-form-item>

        <el-form-item label="证书文件">
          <el-upload
            class="upload-demo"
            drag
            action="#"
            :auto-upload="false"
            :on-change="handleCertFileChange"
            accept=".pem,.crt,.PEM,.CRT"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将证书文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                上传文件格式：.pem、.crt、.PEM、.CRT
              </div>
            </template>
          </el-upload>
          <div class="form-tip">
            正确文件格式：证书文件以 -----BEGIN CERTIFICATE----- 开头，以 -----END CERTIFICATE----- 结尾
          </div>
        </el-form-item>

        <el-form-item label="私钥" prop="privateKey">
          <el-input
            v-model="formData.privateKey"
            type="textarea"
            :rows="8"
            placeholder="请粘贴私钥内容，以 -----BEGIN (RSA|EC) PRIVATE KEY----- 开头，以 -----END (RSA|EC) PRIVATE KEY----- 结尾"
            style="font-family: monospace;"
          />
          <div class="form-tip">
            上传文件格式：.key、.KEY
          </div>
        </el-form-item>

        <el-form-item label="证书链">
          <el-input
            v-model="formData.certChain"
            type="textarea"
            :rows="8"
            placeholder="请粘贴证书链内容，以 -----BEGIN CERTIFICATE----- 开头，以 -----END CERTIFICATE----- 结尾"
            style="font-family: monospace;"
          />
          <div class="form-tip">
            上传文件格式：.pem、.PEM
          </div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">导入证书</el-button>
          <el-button @click="goBack">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElForm } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'

const router = useRouter()
const formRef = ref<InstanceType<typeof ElForm>>()
const submitting = ref(false)

const formData = reactive({
  name: '',
  certContent: '',
  certFile: null,
  privateKey: '',
  certChain: ''
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入证书名称', trigger: 'blur' }
  ],
  certContent: [
    { required: true, message: '请输入证书内容或上传证书文件', trigger: 'blur' }
  ]
}

const handleCertFileChange = (file: any) => {
  formData.certFile = file.raw

  // 读取文件内容并填充到文本框
  const reader = new FileReader()
  reader.onload = (e) => {
    const content = e.target?.result as string
    if (content) {
      formData.certContent = content
      ElMessage.success('证书文件读取成功')
    }
  }
  reader.readAsText(file.raw)
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    // 验证证书内容格式
    if (formData.certContent && !formData.certContent.includes('-----BEGIN CERTIFICATE-----')) {
      ElMessage.error('证书内容格式不正确，请确保包含完整的证书内容')
      return
    }

    if (formData.privateKey && !formData.privateKey.includes('-----BEGIN') && !formData.privateKey.includes('PRIVATE KEY-----')) {
      ElMessage.error('私钥内容格式不正确，请确保包含完整的私钥内容')
      return
    }

    console.log('提交数据:', formData)

    // 模拟提交
    await new Promise(resolve => setTimeout(resolve, 2000))

    ElMessage.success('证书导入成功')
    router.push('/dns/certs')

  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    submitting.value = false
  }
}

const goBack = () => {
  router.push('/dns/certs')
}
</script>

<style scoped>
.import-cert-page {
  padding: 20px;
}

.page-header {
  margin-bottom: 30px;
}

.page-header h1 {
  margin: 0 0 10px 0;
  font-size: 24px;
  color: #303133;
}

.page-header p {
  margin: 0;
  color: #606266;
}

.form-container {
  max-width: 800px;
  background: white;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.upload-demo {
  width: 100%;
}

.form-tip {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
  line-height: 1.4;
}

:deep(.el-textarea__inner) {
  font-family: 'Courier New', Consolas, monospace;
  font-size: 13px;
  line-height: 1.4;
}

:deep(.el-upload-dragger) {
  border: 2px dashed #d9d9d9;
  border-radius: 6px;
  width: 100%;
  height: 120px;
  text-align: center;
  position: relative;
  overflow: hidden;
  background-color: #fafafa;
  transition: border-color 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
}

:deep(.el-upload-dragger:hover) {
  border-color: #409eff;
}

:deep(.el-upload__text) {
  color: #606266;
  font-size: 14px;
  text-align: center;
}

:deep(.el-upload__text em) {
  color: #409eff;
  font-style: normal;
}
</style>


