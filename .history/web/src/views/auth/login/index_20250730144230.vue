<template>
  <div class="login-bg-animate">
    <!-- DevOps 风格的代码背景 -->
    <div class="code-background">
      <div class="code-line">$ kubectl get pods --all-namespaces</div>
      <div class="code-line">$ docker ps -a</div>
      <div class="code-line">$ terraform plan -var-file="production.tfvars"</div>
      <div class="code-line">$ ansible-playbook deploy.yml</div>
      <div class="code-line">$ helm upgrade devops-asset ./charts/devops-asset</div>
      <div class="code-line">$ git push origin feature/new-deployment</div>
      <div class="code-line">$ jenkins build deploy-pipeline</div>
      <div class="code-line">$ prometheus query up{job="devops-asset"}</div>
    </div>
    
    <!-- 科技风格的动画元素 -->
    <svg
      class="bg-anim-svg"
      width="100vw"
      height="100vh"
      viewBox="0 0 1920 1080"
      preserveAspectRatio="none"
    >
      <!-- DevOps 图标节点 -->
      <g class="devops-node" transform="translate(200,200)">
        <circle cx="0" cy="0" r="25" fill="#00D4AA33" stroke="#00D4AA" stroke-width="2">
          <animate attributeName="r" values="25;35;25" dur="3s" repeatCount="indefinite" />
        </circle>
        <text x="0" y="5" text-anchor="middle" fill="#00D4AA" font-size="10" font-family="monospace">K8s</text>
      </g>
      
      <g class="devops-node" transform="translate(1500,300)">
        <circle cx="0" cy="0" r="20" fill="#2496ED33" stroke="#2496ED" stroke-width="2">
          <animate attributeName="r" values="20;30;20" dur="4s" repeatCount="indefinite" />
        </circle>
        <text x="0" y="5" text-anchor="middle" fill="#2496ED" font-size="9" font-family="monospace">Docker</text>
      </g>
      
      <g class="devops-node" transform="translate(400,700)">
        <circle cx="0" cy="0" r="22" fill="#623CE433" stroke="#623CE4" stroke-width="2">
          <animate attributeName="r" values="22;32;22" dur="5s" repeatCount="indefinite" />
        </circle>
        <text x="0" y="5" text-anchor="middle" fill="#623CE4" font-size="9" font-family="monospace">Terraform</text>
      </g>
      
      <g class="devops-node" transform="translate(1200,600)">
        <circle cx="0" cy="0" r="18" fill="#EE0000333" stroke="#EE0000" stroke-width="2">
          <animate attributeName="r" values="18;28;18" dur="6s" repeatCount="indefinite" />
        </circle>
        <text x="0" y="5" text-anchor="middle" fill="#EE0000" font-size="8" font-family="monospace">Ansible</text>
      </g>
      
      <!-- 连接线 -->
      <path d="M 200,200 Q 850,350 1500,300" stroke="#00D4AA44" stroke-width="2" fill="none">
        <animate attributeName="stroke-dasharray" values="0,1000;1000,0" dur="8s" repeatCount="indefinite" />
      </path>
      
      <path d="M 400,700 Q 800,450 1200,600" stroke="#623CE444" stroke-width="2" fill="none">
        <animate attributeName="stroke-dasharray" values="0,800;800,0" dur="10s" repeatCount="indefinite" />
      </path>
      
      <!-- 数据流动效果 -->
      <circle r="4" fill="#00D4AA">
        <animateMotion dur="8s" repeatCount="indefinite">
          <path d="M 200,200 Q 850,350 1500,300" />
        </animateMotion>
      </circle>
      
      <circle r="4" fill="#623CE4">
        <animateMotion dur="10s" repeatCount="indefinite">
          <path d="M 400,700 Q 800,450 1200,600" />
        </animateMotion>
      </circle>
    </svg>
    <div class="login-center-card">
      <div class="logo-row">
        <div class="devops-logo">
          <div class="logo-icon">
            <i class="iconsys-container"></i>
          </div>
          <div class="logo-text">
            <span class="logo-title">DevOps Asset</span>
            <span class="logo-version">v2.0</span>
          </div>
        </div>
      </div>
      <div class="devops-slogan">
        <div class="slogan-main">统一的 DevOps 资产管理平台</div>
        <div class="slogan-sub">Kubernetes • Docker • CI/CD • Monitoring</div>
        <div class="status-indicator">
          <span class="status-dot"></span>
          <span class="status-text">System Online</span>
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
