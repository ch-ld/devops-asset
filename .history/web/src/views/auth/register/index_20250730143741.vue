<template>
  <div class="login-bg-animate">
    <svg
      class="bg-anim-svg"
      width="100vw"
      height="100vh"
      viewBox="0 0 1920 1080"
      preserveAspectRatio="none"
    >
      <circle cx="300" cy="300" r="40" fill="#5D87FF22">
        <animate attributeName="r" values="40;60;40" dur="4s" repeatCount="indefinite" />
      </circle>
      <circle cx="1600" cy="200" r="30" fill="#38C0FC22">
        <animate attributeName="r" values="30;50;30" dur="5s" repeatCount="indefinite" />
      </circle>
      <circle cx="500" cy="800" r="25" fill="#5D87FF11">
        <animate attributeName="r" values="25;40;25" dur="6s" repeatCount="indefinite" />
      </circle>
      <line x1="300" y1="300" x2="1600" y2="200" stroke="#5D87FF33" stroke-width="2">
        <animate attributeName="x2" values="1600;1400;1600" dur="6s" repeatCount="indefinite" />
      </line>
      <line x1="500" y1="800" x2="300" y2="300" stroke="#38C0FC33" stroke-width="2">
        <animate attributeName="x2" values="300;400;300" dur="5s" repeatCount="indefinite" />
      </line>
      <circle cx="1200" cy="700" r="18" fill="#38C0FC22">
        <animate attributeName="r" values="18;32;18" dur="7s" repeatCount="indefinite" />
      </circle>
      <line x1="1200" y1="700" x2="1600" y2="200" stroke="#5D87FF22" stroke-width="2">
        <animate attributeName="y2" values="200;300;200" dur="7s" repeatCount="indefinite" />
      </line>
    </svg>
    <div class="login-center-card register-card">
      <div class="logo-row">
        <ArtLogo class="logo-svg" :size="48" />
        <span class="logo-title">DevOps Asset</span>
      </div>
      <div class="slogan">欢迎注册 DevOps 资产管理平台</div>
      <ElForm
        ref="formRef"
        :model="formData"
        :rules="rules"
        @keyup.enter="handleSubmit"
        class="login-form register-form"
      >
        <ElFormItem prop="username">
          <ElInput placeholder="请输入用户名" size="large" v-model.trim="formData.username" />
        </ElFormItem>
        <ElFormItem prop="password">
          <ElInput
            placeholder="请输入密码"
            size="large"
            type="password"
            show-password
            v-model.trim="formData.password"
          />
        </ElFormItem>
        <ElFormItem prop="confirmPassword">
          <ElInput
            placeholder="请确认密码"
            size="large"
            type="password"
            show-password
            v-model.trim="formData.confirmPassword"
          />
        </ElFormItem>
        <ElFormItem prop="name">
          <ElInput placeholder="请输入真实姓名" size="large" v-model.trim="formData.name" />
        </ElFormItem>
        <ElFormItem prop="phone">
          <ElInput placeholder="请输入手机号" size="large" v-model.trim="formData.phone" />
        </ElFormItem>
        <ElFormItem prop="gender">
          <ElSelect placeholder="请选择性别" size="large" v-model="formData.gender" style="width: 100%">
            <ElOption label="男" :value="1" />
            <ElOption label="女" :value="2" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem>
          <ElButton
            type="primary"
            size="large"
            :loading="loading"
            @click="handleSubmit"
            class="login-btn"
          >
            注册
          </ElButton>
        </ElFormItem>
      </ElForm>
      <div class="footer-links">
        <RouterLink to="/auth/login" class="link">已有账号？立即登录</RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage, type FormRules } from 'element-plus'
import { useRouter } from 'vue-router'
import ArtLogo from '@/components/core/base/art-logo/index.vue'
import { registerAPI } from '@/api/system/api'

const router = useRouter()
const formRef = ref()
const loading = ref(false)

// 表单数据
const formData = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  name: '',
  phone: '',
  gender: 1
})

// 验证规则
const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度在 3 到 50 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== formData.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  name: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' },
    { min: 2, max: 50, message: '姓名长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式', trigger: 'blur' }
  ],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }]
}

// 提交注册
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        await registerAPI({
          username: formData.username,
          password: formData.password,
          name: formData.name,
          phone: formData.phone,
          gender: formData.gender
        })
        
        ElMessage.success('注册成功！请登录')
        router.push('/auth/login')
      } catch (error: any) {
        ElMessage.error(error.message || '注册失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped lang="scss">
@import './style.scss';

.register-card {
  min-height: 600px;
  width: 420px;
}

.register-form {
  .el-form-item {
    margin-bottom: 20px;
  }
}

.footer-links {
  text-align: center;
  margin-top: 20px;
  
  .link {
    color: var(--el-color-primary);
    text-decoration: none;
    font-size: 14px;
    
    &:hover {
      text-decoration: underline;
    }
  }
}
</style>
