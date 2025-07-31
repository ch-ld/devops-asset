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
                  <stop offset="0%" style="stop-color:#6366f1"/>
                  <stop offset="100%" style="stop-color:#8b5cf6"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <div class="brand-text">
            <h1 class="brand-title">DevOps Asset</h1>
            <p class="brand-subtitle">统一的 DevOps 资产管理平台</p>
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
                size="large"
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
                size="large"
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
                  size="large"
                  v-model.trim="formData.captcha"
                  class="modern-input"
                />
              </div>
              <div class="captcha-wrapper">
                <img
                  :src="captchaImageUrl"
                  @click="refreshCaptcha"
                  class="captcha-image"
                  alt="点击刷新验证码"
                  @error="handleCaptchaImageError"
                />
                <div class="captcha-refresh" @click="refreshCaptcha">
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
      <ElForm
        ref="formRef"
        :model="formData"
        :rules="rules"
        @keyup.enter="handleSubmit"
        class="login-form"
      >
        <ElFormItem prop="username" class="devops-form-item">
          <ElInput 
            placeholder="Username / 用户名" 
            size="large" 
            v-model.trim="formData.username"
            class="devops-input"
          >
            <template #prefix>
              <i class="iconsys-user" style="color: #00D4AA;"></i>
            </template>
          </ElInput>
        </ElFormItem>
        <ElFormItem prop="password" class="devops-form-item">
          <ElInput
            placeholder="Password / 密码"
            size="large"
            v-model.trim="formData.password"
            type="password"
            autocomplete="off"
            show-password
            class="devops-input"
          >
            <template #prefix>
              <i class="iconsys-lock" style="color: #2496ED;"></i>
            </template>
          </ElInput>
        </ElFormItem>
        <ElFormItem prop="captcha" class="devops-form-item">
          <ElRow :gutter="8">
            <ElCol :span="16">
              <ElInput 
                placeholder="Captcha / 验证码" 
                size="large" 
                v-model.trim="formData.captcha"
                class="devops-input"
              >
                <template #prefix>
                  <i class="iconsys-security" style="color: #623CE4;"></i>
                </template>
              </ElInput>
            </ElCol>
            <ElCol :span="8">
              <div class="captcha-container">
                <img
                  :src="captchaImageUrl"
                  @click="refreshCaptcha"
                  class="captcha-image"
                  alt="点击刷新验证码"
                  @error="handleCaptchaImageError"
                />
                <div class="captcha-overlay">
                  <i class="iconsys-refresh" @click="refreshCaptcha"></i>
                </div>
              </div>
            </ElCol>
          </ElRow>
        </ElFormItem>
        
        <div class="login-options">
          <ElCheckbox v-model="formData.rememberPassword" class="devops-checkbox">
            <span class="checkbox-text">Remember / 记住密码</span>
          </ElCheckbox>
          <router-link to="/auth/forget-password" class="forgot-link">Forgot Password?</router-link>
        </div>
        
        <div class="login-actions">
          <ElButton
            class="devops-login-btn"
            size="large"
            type="primary"
            @click="handleSubmit"
            :loading="loading"
          >
            <template #icon>
              <i class="iconsys-login" v-if="!loading"></i>
            </template>
            {{ loading ? 'AUTHENTICATING...' : 'ACCESS SYSTEM' }}
          </ElButton>
        </div>
      </ElForm>
      <div class="register-link">
        <RouterLink to="/auth/register" class="link">还没有账号？立即注册</RouterLink>
      </div>
      <div class="login-footer">
        <a href="https://your-docs.com/devops-asset" target="_blank">文档</a> ·
        <a href="https://github.com/your-repo/devops-asset" target="_blank">Github</a> ·
        <a href="https://your-homepage.com" target="_blank">官网</a>
        <br />
        Copyright © 2024 DevOps Asset
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
    console.error('验证码图片加载失败')
    ElMessage.error('验证码图片加载失败')
    // 尝试重新获取验证码
    setTimeout(() => {
      refreshCaptcha()
    }, 1000)
  }

  onMounted(() => {
    refreshCaptcha() // 页面加载时获取验证码
  })
</script>

<style lang="scss" scoped>
  @use './index';
</style>
