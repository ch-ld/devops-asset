<template>
  <div class="certificate-import-wrapper">
    <!-- 步骤条 -->
    <el-steps :active="activeStep" align-center finish-status="success" class="import-steps">
      <el-step title="上传文件" />
      <el-step title="验证文件" />
      <el-step title="完成导入" />
    </el-steps>

    <!-- 步骤内容 -->
    <div class="step-content">
      <!-- Step 1 -->
      <div v-if="activeStep === 0" class="step-upload">
        <el-form :model="form" label-width="120px" ref="uploadFormRef">
          <el-form-item label="证书名称" prop="name" :rules="[{ required: true, message: '请输入名称' }]"><el-input v-model="form.name" placeholder="如：example.com" /></el-form-item>
          <el-form-item label="证书文件" prop="certFile" :rules="[{ required: true, message: '请上传证书' }]">
            <el-upload
              class="upload-block"
              drag
              action="#"
              :before-upload="() => false"
              :on-change="file => (form.certFile = file.raw)"
              accept=".crt,.pem,.cer"
            >
              <upload-filled />
              <div class="el-upload__text">将 .crt / .pem / .cer 文件拖到此处，或点击上传</div>
            </el-upload>
          </el-form-item>
          <el-form-item label="私钥文件" prop="keyFile" :rules="[{ required: true, message: '请上传私钥' }]">
            <el-upload
              class="upload-block"
              drag
              action="#"
              :before-upload="() => false"
              :on-change="file => (form.keyFile = file.raw)"
              accept=".key,.pem"
            >
              <upload-filled />
              <div class="el-upload__text">将 .key / .pem 文件拖到此处，或点击上传</div>
            </el-upload>
          </el-form-item>
          <el-form-item label="证书链 (可选)">
            <el-upload
              class="upload-block"
              drag
              action="#"
              :before-upload="() => false"
              :on-change="file => (form.chainFile = file.raw)"
              accept=".pem,.crt"
            >
              <upload-filled />
              <div class="el-upload__text">将链文件拖到此处，或点击上传</div>
            </el-upload>
          </el-form-item>
        </el-form>
      </div>

      <!-- Step 2 -->
      <div v-if="activeStep === 1" class="step-validate" v-loading="validating">
        <el-result icon="info" title="正在验证文件">
          <template #extra>
            <el-progress :percentage="validateProgress" status="success" style="width: 240px" />
          </template>
        </el-result>
      </div>

      <!-- Step 3 -->
      <div v-if="activeStep === 2" class="step-finish">
        <el-result icon="success" title="导入完成" sub-title="证书已成功导入系统">
          <template #extra>
            <el-button type="primary" @click="goList">返回列表</el-button>
          </template>
        </el-result>
      </div>
    </div>

    <!-- 底部操作栏 -->
    <div class="footer-actions">
      <el-button v-if="activeStep === 0" @click="goList">取消</el-button>
      <el-button v-if="activeStep === 0" type="primary" @click="nextStep">下一步</el-button>
      <el-button v-if="activeStep === 1" @click="prevStep">上一步</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import { certificateApi } from '@/api/dns/certificate'

const router = useRouter()

interface FormState {
  name: string
  certFile: File | null
  keyFile: File | null
  chainFile: File | null
}

const form = reactive<FormState>({ name: '', certFile: null, keyFile: null, chainFile: null })
const activeStep = ref(0)
const validating = ref(false)
const validateProgress = ref(0)
const uploadFormRef = ref()

function goList() {
  router.push('/dns/certs')
}

function prevStep() {
  activeStep.value--
}

async function nextStep() {
  if (activeStep.value === 0) {
    // 校验表单
    await uploadFormRef.value.validate()
    activeStep.value = 1
    await validateFiles()
  }
}

async function validateFiles() {
  try {
    validating.value = true
    validateProgress.value = 20
    // 读取文件文本
    const certContent = await readFileText(form.certFile!)
    const keyContent = await readFileText(form.keyFile!)
    const chainContent = form.chainFile ? await readFileText(form.chainFile) : ''
    validateProgress.value = 60
    // 后端验证
    const { data } = await certificateApi.validateCertificate({ certContent, keyContent, chainContent })
    if (data.code !== 0) throw new Error(data.message || '验证失败')
    validateProgress.value = 100
    // 询问导入
    await ElMessageBox.confirm('文件验证通过，确定导入吗？', '确认导入', { confirmButtonText: '导入', cancelButtonText: '取消' })
    await uploadCertificate(certContent, keyContent, chainContent)
  } catch (e: any) {
    ElMessage.error(e.message || '验证失败')
    activeStep.value = 0
  } finally {
    validating.value = false
  }
}

async function uploadCertificate(certContent: string, keyContent: string, chainContent: string) {
  try {
    const { data } = await certificateApi.upload({ cert_content: certContent, key_content: keyContent, chain_content: chainContent, remark: form.name })
    if (data.code !== 0) throw new Error(data.message || '上传失败')
    activeStep.value = 2
    ElMessage.success('导入成功')
  } catch (e: any) {
    ElMessage.error(e.message || '导入失败')
    activeStep.value = 0
  }
}

function readFileText(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result as string)
    reader.onerror = reject
    reader.readAsText(file)
  })
}
</script>

<style scoped>
.certificate-import-wrapper {
  max-width: 760px;
  margin: 0 auto;
  padding: 24px;
}
.import-steps {
  margin-bottom: 32px;
}
.step-content {
  min-height: 260px;
}
.upload-block {
  width: 100%;
}
.footer-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
}
</style>
