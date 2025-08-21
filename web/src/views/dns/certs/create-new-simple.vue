<template>
  <div class="cert-create-page">
    <div class="page-header">
      <h1>申请SSL证书</h1>
      <p>为您的域名申请免费的Let's Encrypt SSL证书</p>
    </div>
    
    <div class="form-container">
      <el-form :model="formData" label-width="120px">
        <el-form-item label="申请类型">
          <el-radio-group v-model="formData.applyType">
            <el-radio label="auto">自动申请</el-radio>
            <el-radio label="upload">导入证书</el-radio>
            <el-radio label="csr">自定义CSR</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="域名">
          <el-input v-model="formData.domain" placeholder="请输入域名" />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSubmit">提交申请</el-button>
          <el-button @click="goBack">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const formData = reactive({
  applyType: 'auto',
  domain: ''
})

const handleSubmit = () => {
  console.log('提交申请', formData)
}

const goBack = () => {
  router.push('/dns/certs')
}
</script>

<style scoped>
.cert-create-page {
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
</style>
