<template>
  <div class="auth-container">
    <!-- 动态背景 -->
    <div class="background-wrapper">
      <!-- 渐变背景 -->
      <div class="gradient-bg"></div>

      <!-- 浮动粒子效果 -->
      <div class="particles">
        <div class="particle" v-for="i in 50" :key="i" :style="getParticleStyle(i)"></div>
      </div>

      <!-- 几何图形装饰 -->
      <div class="geometric-shapes">
        <div class="shape shape-1"></div>
        <div class="shape shape-2"></div>
        <div class="shape shape-3"></div>
        <div class="shape shape-4"></div>
      </div>

      <!-- 科技网格 -->
      <div class="tech-grid">
        <svg width="100%" height="100%" class="grid-svg">
          <defs>
            <pattern id="grid" width="60" height="60" patternUnits="userSpaceOnUse">
              <path d="M 60 0 L 0 0 0 60" fill="none" stroke="rgba(99, 102, 241, 0.1)" stroke-width="1"/>
            </pattern>
          </defs>
          <rect width="100%" height="100%" fill="url(#grid)" />
        </svg>
      </div>
    </div>
    <!-- 主登录卡片 -->
    <div class="auth-card">
      <!-- 卡片头部 -->
      <div class="card-header">
        <div class="brand-section">
          <div class="brand-icon">
            <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
              <rect width="48" height="48" rx="12" fill="url(#brandGradient)"/>
              <path d="M16 20L24 14L32 20V32L24 38L16 32V20Z" stroke="white" stroke-width="2" fill="none"/>
              <circle cx="24" cy="24" r="3" fill="white"/>
              <defs>
                <linearGradient id="brandGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" style="stop-color:#0f52ba"/>
                  <stop offset="100%" style="stop-color:#00bcd4"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <div class="brand-text">
            <h1 class="brand-title">DevOps Asset</h1>
            <p class="brand-subtitle">深色科技蓝 DevOps 平台 · Cloud · DNS · Cert · K8s</p>
          </div>
        </div>

        <div class="welcome-section">
          <h2 class="welcome-title">欢迎回来</h2>
          <p class="welcome-subtitle">请登录您的账户以继续</p>
        </div>
      </div>

      <!-- 登录表单 -->
      <div class="card-body">
        <ElForm
          ref="formRef"
          :model="formData"
          :rules="rules"
          @keyup.enter="handleSubmit"
          class="auth-form"
          size="large"
        >
          <ElFormItem prop="username">
            <div class="input-group">
              <div class="input-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                  <path d="M12 12C14.7614 12 17 9.76142 17 7C17 4.23858 14.7614 2 12 2C9.23858 2 7 4.23858 7 7C7 9.76142 9.23858 12 12 12Z" fill="currentColor"/>
                  <path d="M12 14C7.58172 14 4 17.5817 4 22H20C20 17.5817 16.4183 14 12 14Z" fill="currentColor"/>
                </svg>
              </div>
              <ElInput
                placeholder="请输入用户名"
                v-model.trim="formData.username"
                class="modern-input"
              />
            </div>
          </ElFormItem>

          <ElFormItem prop="password">
            <div class="input-group">
              <div class="input-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                  <path d="M6 10V8C6 5.79086 7.79086 4 10 4H14C16.2091 4 18 5.79086 18 8V10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                  <rect x="4" y="10" width="16" height="10" rx="2" fill="currentColor"/>
                </svg>
              </div>
              <ElInput
                placeholder="请输入密码"
                type="password"
                show-password
                v-model.trim="formData.password"
                class="modern-input"
              />
            </div>
          </ElFormItem>

          <ElFormItem prop="captcha">
            <div class="captcha-row">
              <div class="input-group captcha-input">
                <div class="input-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <rect x="3" y="5" width="18" height="14" rx="2" stroke="currentColor" stroke-width="2"/>
                    <path d="M9 12L11 14L15 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </div>
                <ElInput
                  placeholder="请输入验证码"
                  v-model.trim="formData.captcha"
                  class="modern-input"
                />
              </div>
              <div class="captcha-wrapper">
                <div class="captcha-image-container">
                  <img
                    :src="captchaImageUrl"
                    @click="refreshCaptcha"
                    class="captcha-image"
                    alt="验证码"
                    @error="handleCaptchaImageError"
                    @load="() => {}"
                  />
                </div>
                <div class="captcha-refresh" @click="refreshCaptcha" title="刷新验证码">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                    <path d="M1 4V10H7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M23 20V14H17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M20.49 9A9 9 0 0 0 5.64 5.64L1 10M23 14L18.36 18.36A9 9 0 0 1 3.51 15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </div>
              </div>
            </div>
          </ElFormItem>

          <div class="form-options">
            <ElCheckbox v-model="formData.rememberPassword" class="remember-checkbox">
              记住密码
            </ElCheckbox>
            <router-link to="/auth/forget-password" class="forgot-link">忘记密码？</router-link>
          </div>

          <ElFormItem>
            <ElButton
              type="primary"
              size="large"
              @click="handleSubmit"
              class="auth-button"
              :loading="loading"
            >
              {{ loading ? '登录中...' : '立即登录' }}
            </ElButton>
          </ElFormItem>
        </ElForm>
      </div>

      <!-- 卡片底部 -->
      <div class="card-footer">
        <div class="register-link">
          <span>还没有账号？</span>
          <RouterLink to="/auth/register" class="link">立即注册</RouterLink>
        </div>

        <div class="footer-links">
          <a href="https://your-docs.com/devops-asset" target="_blank">文档</a>
          <span>·</span>
          <a href="https://github.com/your-repo/devops-asset" target="_blank">Github</a>
          <span>·</span>
          <a href="https://your-homepage.com" target="_blank">官网</a>
        </div>

        <div class="copyright">
          Copyright © 2024 DevOps Asset
        </div>
      </div>
    </div>
  </div>

</template>

<script setup lang="ts">
  import AppConfig from '@/config'
  import { ElForm, ElMessage, ElNotification } from 'element-plus'
  import { useUserStore } from '@/store/modules/user'
  // ApiStatus已不再使用，直接比较状态码
  // import { ApiStatus } from '@/utils/http/status'
  import { SystemThemeEnum } from '@/enums/appEnum'
  import { useSettingStore } from '@/store/modules/setting'
  import type { FormInstance, FormRules } from 'element-plus'
  import { onMounted, ref, reactive, computed } from 'vue'
  import { getCaptcha, userLogin, getUserInfo } from '@/api/system/api'
  import defaultAvatar from '@/assets/img/user/avatar.png'
  import { storeToRefs } from 'pinia'

  defineOptions({ name: 'Login' })
  const settingStore = useSettingStore()
  const { isDark, systemThemeType } = storeToRefs(settingStore)

  const userStore = useUserStore()
  const router = useRouter()

  // 定义一些未使用的变量为_开头，避免lint警告
  const _systemName = AppConfig.systemInfo.name
  const formRef = ref<FormInstance>()
  const formData = reactive({
    username: import.meta.env.DEV ? 'admin' : '',
    password: import.meta.env.DEV ? '123456' : '',
    rememberPassword: true,
    captcha: ''
  })

  const rules = computed<FormRules>(() => ({
    username: [{ required: true, message: '请输入账号', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
    captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }]
  }))

  const loading = ref(false)
  const captchaImageUrl = ref('') // 验证码图片的URL
  const captchaImageID = ref('') // 验证码图片的ID

  // 粒子动画样式生成
  const getParticleStyle = (index: number) => {
    const size = Math.random() * 4 + 2
    const duration = Math.random() * 20 + 10
    const delay = Math.random() * 5
    const x = Math.random() * 100
    const y = Math.random() * 100

    return {
      width: `${size}px`,
      height: `${size}px`,
      left: `${x}%`,
      top: `${y}%`,
      animationDuration: `${duration}s`,
      animationDelay: `${delay}s`
    }
  }

  const handleSubmit = async () => {
    if (!formRef.value) return

    await formRef.value.validate(async (valid) => {
      if (valid) {
        loading.value = true
        // 延时辅助函数
        const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))
        try {
          const res = await userLogin({
            username: formData.username,
            password: formData.password,
            captcha: formData.captcha,
            captcha_id: captchaImageID.value
          })
          // 调试查看登录响应
          console.log('Login response:', res)
          // 修改状态码判断为200，而非ApiStatus.success
          if (res.code === 200 && res.data) {
            // 设置 token
            userStore.setToken(res.data.access_token)
            // 获取用户信息
            const userRes = await getUserInfo()
            console.log('User info response:', userRes)
            // 修改状态码判断
            if (userRes.code === 200 && userRes.data) {
              console.log('获取用户信息成功:', userRes.data)

              // 转换为UserInfo类型
              const userInfo: Api.User.UserInfo = {
                userId: userRes.data.id,
                userName: userRes.data.username,
                roles: [userRes.data.role?.code || ''],
                buttons: [],
                avatar: defaultAvatar,
                email: userRes.data.email,
                phone: userRes.data.phone
              }

              userStore.setUserInfo(userInfo)
            } else {
              ElMessage.error(userRes.message)
              console.error('获取用户信息失败:', userRes.message)
            }
            // 设置登录状态
            userStore.setLoginStatus(true)
            // 延时辅助函数
            await delay(1000)
            // 登录成功提示
            showLoginSuccessNotice()
            // 跳转首页
            console.log('登录成功，跳转首页')
            router.push('/')
          } else {
            ElMessage.error(res.message)
            refreshCaptcha()
          }
        } finally {
          await delay(1000)
          loading.value = false
        }
      }
    })
  }

  // 登录成功提示
  const showLoginSuccessNotice = () => {
    setTimeout(() => {
      ElNotification({
        title: '登录成功',
        type: 'success',
        showClose: true,
        duration: 2500,
        zIndex: 10000,
        message: `欢迎回来!`
      })
    }, 300)
  }

  // 切换主题
  import { useTheme } from '@/composables/useTheme'

  // 定义但不使用，添加_前缀
  const _toggleTheme = () => {
    let { LIGHT, DARK } = SystemThemeEnum
    useTheme().switchThemeStyles(systemThemeType.value === LIGHT ? DARK : LIGHT)
  }

  const refreshCaptcha = async () => {
    try {
      const captchaRes = await getCaptcha(80, 240)
      // 调试检查响应格式
      console.log('Captcha response:', captchaRes)

      // 全面检查各种可能的响应格式
      if (captchaRes) {
        if (captchaRes.code === 200 && captchaRes.data) {
          // 正常格式
          if (captchaRes.data.image && captchaRes.data.id) {
            captchaImageUrl.value = captchaRes.data.image
            captchaImageID.value = captchaRes.data.id
            console.log('验证码加载成功(标准格式):', {
              id: captchaImageID.value,
              imageLength: captchaImageUrl.value?.length
            })
            return
          }

          // 处理响应中直接包含image和id的情况
          if (typeof captchaRes.data === 'object') {
            const { id, image } = captchaRes.data as { id?: string; image?: string }
            if (id && image) {
              captchaImageUrl.value = image
              captchaImageID.value = id
              console.log('验证码加载成功(直接包含):', { id, imageLength: image.length })
              return
            }
          }
        }

        // 尝试在外层查找id和image
        // 使用类型断言避免类型错误
        const anyResponse = captchaRes as any
        if (anyResponse.image && anyResponse.id) {
          captchaImageUrl.value = anyResponse.image
          captchaImageID.value = anyResponse.id
          console.log('验证码加载成功(外层数据):', { id: anyResponse.id })
          return
        }
      }

      // 如果所有格式都不匹配
      console.error('验证码响应格式无法识别:', captchaRes)
      ElMessage.error('验证码格式错误')
    } catch (error) {
      console.error('Error refreshing captcha:', error)
      ElMessage.error('验证码获取失败')
    }
  }

  // 验证码图片加载错误处理
  const handleCaptchaImageError = () => {
    console.error('验证码图片加载失败，尝试重新获取')
    // 静默重新获取验证码，不显示错误提示
    setTimeout(() => {
      refreshCaptcha()
    }, 500)
  }

  onMounted(() => {
    refreshCaptcha() // 页面加载时获取验证码
  })
</script>

<style lang="scss" scoped>
.auth-container {
  position: relative;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.background-wrapper {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
}

.gradient-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg,
    #0a0f1c 0%,
    #0f1e3a 25%,
    #0b2a6f 50%,
    #0d3b8c 75%,
    #0f52ba 100%
  );
  background-size: 400% 400%;
  animation: gradientShift 20s ease infinite;
}

@keyframes gradientShift {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.particles {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.particle {
  position: absolute;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 50%;
  animation: float linear infinite;
}

@keyframes float {
  0% {
    transform: translateY(100vh) rotate(0deg);
    opacity: 0;
  }
  10% {
    opacity: 1;
  }
  90% {
    opacity: 1;
  }
  100% {
    transform: translateY(-100px) rotate(360deg);
    opacity: 0;
  }
}

.geometric-shapes {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  animation: shapeFloat 20s ease-in-out infinite;
}

.shape-1 {
  width: 200px;
  height: 200px;
  top: 10%;
  left: 10%;
  animation-delay: 0s;
}

.shape-2 {
  width: 150px;
  height: 150px;
  top: 60%;
  right: 15%;
  animation-delay: 5s;
}

.shape-3 {
  width: 100px;
  height: 100px;
  bottom: 20%;
  left: 20%;
  animation-delay: 10s;
}

.shape-4 {
  width: 120px;
  height: 120px;
  top: 30%;
  right: 30%;
  animation-delay: 15s;
}

@keyframes shapeFloat {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(180deg); }
}

.tech-grid {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0.4;
  background-image:
    linear-gradient(rgba(99, 102, 241, 0.1) 1px, transparent 1px),
    linear-gradient(90deg, rgba(99, 102, 241, 0.1) 1px, transparent 1px);
  background-size: 50px 50px;
  animation: gridMove 30s linear infinite;
}

@keyframes gridMove {
  0% { transform: translate(0, 0); }
  100% { transform: translate(50px, 50px); }
}

.auth-card {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 420px;
  margin: 0 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.1),
    0 0 0 1px rgba(255, 255, 255, 0.2);
  overflow: hidden;
  animation: cardSlideIn 0.8s ease-out;
}

@keyframes cardSlideIn {
  0% {
    opacity: 0;
    transform: translateY(30px) scale(0.95);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.card-header {
  padding: 40px 40px 20px;
  text-align: center;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.brand-section {
  margin-bottom: 30px;
}

.brand-icon {
  margin-bottom: 16px;
  display: flex;
  justify-content: center;
}

.brand-text {
  .brand-title {
    font-size: 28px;
    font-weight: 700;
    color: #1a1a1a;
    margin: 0 0 8px 0;
    background: linear-gradient(135deg, #007cf0, #00dfd8);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .brand-subtitle {
    font-size: 14px;
    color: #6b7280;
    margin: 0;
  }
}

.welcome-section {
  .welcome-title {
    font-size: 24px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 8px 0;
  }

  .welcome-subtitle {
    font-size: 14px;
    color: #6b7280;
    margin: 0;
  }
}

.card-body {
  padding: 20px 40px 40px;
}

.auth-form {
  .el-form-item {
    margin-bottom: 24px;

    &:last-child {
      margin-bottom: 0;
    }
  }
}

.input-group {
  width: 100%; // 保证输入框组占满行宽，左右对齐
  position: relative;
  display: flex;
  align-items: center;
  gap: 12px;
}

.input-group .input-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 2;
  color: #9ca3af;
  transition: color 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.input-group .modern-input {
  :deep(.el-input__wrapper) {
    padding-left: 48px;
    border-radius: 12px;
    border: 2px solid #e5e7eb;
    background: #ffffff;
    box-shadow: none;
    transition: all 0.3s ease;

    &:hover {
      border-color: #d1d5db;
    }

    &.is-focus {
      border-color: #6366f1;
      box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
    }
  }

  :deep(.el-input__inner) {
    color: #1f2937;

    &::placeholder {
      color: #9ca3af;
    }
  }
}

.input-group:focus-within .input-icon {
  color: #6366f1;
}

.captcha-row {
  display: grid; // 采用网格便于对齐
  grid-template-columns: minmax(160px, 1fr) auto; // 取消刷新按钮列
  gap: 12px;
  align-items: center; // 垂直居中
}

.captcha-row .captcha-input {
  width: 100%;
  min-width: 160px; // 最小宽度
}
.captcha-row .captcha-input .modern-input {
  width: 100%;
}

.captcha-wrapper {
  display: flex;
  align-items: center;
  gap: 0; // 去除图片与按钮间距
}

.captcha-wrapper .captcha-image-container {
  position: relative;
  width: 180px; // 加宽验证码图片，实际宽度可根据图片自适应
  height: 48px;
  border-radius: 12px;
  overflow: hidden;
  border: 2px solid #e5e7eb;
  background: #ffffff;
  transition: border-color 0.3s ease;

  &:hover {
    border-color: #6366f1;
  }
}

.captcha-wrapper .captcha-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  cursor: pointer;
  display: block;
  background: #ffffff;
}

.captcha-wrapper .captcha-refresh {
  display: none !important;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;

  .remember-checkbox {
    :deep(.el-checkbox__label) {
      font-size: 14px;
      color: #6b7280;
    }

    :deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
      background-color: #6366f1;
      border-color: #6366f1;
    }
  }

  .forgot-link {
    font-size: 14px;
    color: #6366f1;
    text-decoration: none;
    transition: color 0.3s ease;

    &:hover {
      color: #4f46e5;
      text-decoration: underline;
    }
  }
}

.auth-button {
  width: 100%;
  height: 52px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #007cf0, #00dfd8);
  border: none;
  transition: all 0.3s ease;

  &:hover {
    background: linear-gradient(135deg, #005bb5, #00c8c0);
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(99, 102, 241, 0.3);
  }

  &:active {
    transform: translateY(0);
  }

  &.is-loading {
    background: linear-gradient(135deg, #9ca3af, #6b7280);
  }
}

.card-footer {
  padding: 20px 40px 40px;
  text-align: center;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
}

.register-link {
  margin-bottom: 24px;
  font-size: 14px;
  color: #6b7280;

  .link {
    color: #6366f1;
    text-decoration: none;
    font-weight: 500;
    margin-left: 4px;
    transition: color 0.3s ease;

    &:hover {
      color: #4f46e5;
      text-decoration: underline;
    }
  }
}

.footer-links {
  margin-bottom: 16px;
  font-size: 13px;
  color: #9ca3af;

  a {
    color: #6b7280;
    text-decoration: none;
    transition: color 0.3s ease;

    &:hover {
      color: #6366f1;
    }
  }

  span {
    margin: 0 8px;
  }
}

.copyright {
  font-size: 12px;
  color: #9ca3af;
}

// 响应式设计
@media (max-width: 768px) {
  .auth-card {
    margin: 0 16px;
    border-radius: 16px;
  }

  .card-header {
    padding: 32px 24px 16px;
  }

  .card-body {
    padding: 16px 24px 32px;
  }

  .card-footer {
    padding: 16px 24px 32px;
  }

  .brand-text .brand-title {
    font-size: 24px;
  }

  .welcome-section .welcome-title {
    font-size: 20px;
  }
}

// 暗色主题支持
@media (prefers-color-scheme: dark) {
  .auth-card {
    background: rgba(17, 24, 39, 0.95);
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.3),
      0 0 0 1px rgba(255, 255, 255, 0.1);
  }

  .brand-text .brand-title {
    color: #f9fafb;
  }

  .welcome-section .welcome-title {
    color: #f9fafb;
  }

  .input-group .modern-input :deep(.el-input__wrapper) {
    background: rgba(31, 41, 55, 0.8);
    border-color: #374151;

    &:hover {
      border-color: #4b5563;
    }

    &.is-focus {
      border-color: #6366f1;
    }
  }

  .input-group .modern-input :deep(.el-input__inner) {
    color: #f9fafb;

    &::placeholder {
      color: #9ca3af;
    }
  }

  .captcha-wrapper {
    border-color: #374151;

    &:hover {
      border-color: #4b5563;
    }
  }
}
</style>
