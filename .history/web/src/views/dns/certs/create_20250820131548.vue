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

        <el-form-item label="私钥文件">
          <el-upload
            class="upload-demo"
            drag
            action="#"
            :auto-upload="false"
            :on-change="handleKeyFileChange"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将私钥文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持 .key, .pem 格式的私钥文件
              </div>
            </template>
          </el-upload>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit">导入证书</el-button>
          <el-button @click="goBack">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { UploadFilled } from '@element-plus/icons-vue'

const router = useRouter()

const formData = reactive({
  name: '',
  domain: '',
  certFile: null,
  keyFile: null
})

const handleCertFileChange = (file: any) => {
  formData.certFile = file.raw
  console.log('证书文件:', file.name)
}

const handleKeyFileChange = (file: any) => {
  formData.keyFile = file.raw
  console.log('私钥文件:', file.name)
}

const handleSubmit = () => {
  console.log('导入证书', formData)
  // 这里可以添加实际的导入逻辑
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
  max-width: 600px;
  background: white;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.upload-demo {
  width: 100%;
}
</style>


