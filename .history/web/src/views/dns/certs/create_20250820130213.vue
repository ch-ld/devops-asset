<template>
  <div class="import-cert-page">
    <div class="page-header">
      <h1>导入SSL证书</h1>
      <p>导入已有的SSL证书文件到系统中进行管理</p>
    </div>

    <div class="form-container">
      <el-form :model="formData" label-width="120px">
        <el-form-item label="证书名称">
          <el-input v-model="formData.name" placeholder="请输入证书名称" />
        </el-form-item>

        <el-form-item label="域名">
          <el-input v-model="formData.domain" placeholder="请输入域名" />
        </el-form-item>

        <el-form-item label="证书文件">
          <el-upload
            class="upload-demo"
            drag
            action="#"
            :auto-upload="false"
            :on-change="handleCertFileChange"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将证书文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持 .crt, .pem, .cer 格式的证书文件
              </div>
            </template>
          </el-upload>
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


