<template>
  <div class="cert-create">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <el-button
            type="text"
            @click="goBack"
            class="back-btn"
          >
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <div class="header-title">
            <h1>申请SSL证书</h1>
            <p>为您的域名申请免费的Let's Encrypt SSL证书</p>
          </div>
        </div>
        <div class="header-actions">
          <el-button @click="goBack">取消</el-button>
          <el-button
            type="primary"
            :loading="applying"
            @click="handleApply"
            :disabled="!canSubmit"
          >
            <el-icon><Lock /></el-icon>
            申请证书
          </el-button>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <el-card class="form-card">
        <template #header>
          <div class="card-header">
            <el-icon class="header-icon"><Document /></el-icon>
            <span>证书配置</span>
          </div>
        </template>

        <el-form
          ref="formRef"
          :model="formData"
          :rules="formRules"
          label-width="120px"
          size="large"
        >
          <!-- 申请方式选择 -->
          <el-form-item label="申请方式" prop="applyType">
            <el-radio-group v-model="formData.applyType" @change="handleApplyTypeChange">
              <el-radio label="auto" class="apply-type-radio">
                <div class="radio-content">
                  <div class="radio-header">
                    <el-icon class="radio-icon"><Lightning /></el-icon>
                    <span class="radio-title">自动申请</span>
                    <el-tag type="success" size="small">推荐</el-tag>
                  </div>
                  <div class="radio-desc">
                    使用Let's Encrypt免费申请SSL证书，支持自动续期
                  </div>
                </div>
              </el-radio>
              <el-radio label="upload" class="apply-type-radio">
                <div class="radio-content">
                  <div class="radio-header">
                    <el-icon class="radio-icon"><FolderAdd /></el-icon>
                    <span class="radio-title">导入证书</span>
                  </div>
                  <div class="radio-desc">
                    导入已有的SSL证书文件
                  </div>
                </div>
              </el-radio>
              <el-radio label="csr" class="apply-type-radio">
                <div class="radio-content">
                  <div class="radio-header">
                    <el-icon class="radio-icon"><Upload /></el-icon>
                    <span class="radio-title">自定义CSR</span>
                  </div>
                  <div class="radio-desc">
                    上传CSR文件申请证书
                  </div>
                </div>
              </el-radio>
            </el-radio-group>
          </el-form-item>

          <!-- 自动申请配置 -->
          <template v-if="formData.applyType === 'auto'">
            <!-- 域名选择 -->
            <el-form-item label="选择域名" prop="domainId" required>
              <el-select
                v-model="formData.domainId"
                placeholder="请选择要申请证书的域名"
                filterable
                @change="handleDomainChange"
                style="width: 100%"
              >
                <el-option
                  v-for="domain in domainOptions"
                  :key="domain.id"
                  :label="domain.name"
                  :value="domain.id"
                >
                  <div class="domain-option">
                    <span class="domain-name">{{ domain.name }}</span>
                    <el-tag
                      v-if="domain.provider"
                      type="info"
                      size="small"
                    >
                      {{ domain.provider.name }}
                    </el-tag>
                  </div>
                </el-option>
              </el-select>
              <div class="form-tip">
                <el-icon><InfoFilled /></el-icon>
                只显示已配置DNS Provider的域名
              </div>
            </el-form-item>

            <!-- 证书类型 -->
            <el-form-item label="证书类型" prop="certType">
              <el-radio-group v-model="formData.certType" @change="handleCertTypeChange">
                <el-radio label="single">
                  <div class="cert-type-option">
                    <span class="option-title">单域名证书</span>
                    <span class="option-desc">仅保护选择的域名</span>
                  </div>
                </el-radio>
                <el-radio label="wildcard">
                  <div class="cert-type-option">
                    <span class="option-title">通配符证书</span>
                    <span class="option-desc">保护域名及其所有子域名</span>
                  </div>
                </el-radio>
                <el-radio label="multi">
                  <div class="cert-type-option">
                    <span class="option-title">多域名证书</span>
                    <span class="option-desc">保护多个不同的域名</span>
                  </div>
                </el-radio>
              </el-radio-group>
            </el-form-item>

            <!-- 域名列表预览 -->
            <el-form-item label="证书域名" v-if="formData.domains.length > 0">
              <div class="domain-preview">
                <el-tag
                  v-for="domain in formData.domains"
                  :key="domain"
                  type="primary"
                  class="domain-tag"
                >
                  <el-icon><Globe /></el-icon>
                  {{ domain }}
                </el-tag>
              </div>
            </el-form-item>

            <!-- 多域名证书的域名编辑 -->
            <el-form-item
              label="域名列表"
              v-if="formData.certType === 'multi'"
            >
              <div class="multi-domain-editor">
                <el-tag
                  v-for="domain in formData.domains"
                  :key="domain"
                  closable
                  @close="removeDomain(domain)"
                  class="domain-tag"
                >
                  {{ domain }}
                </el-tag>
                <el-input
                  v-if="domainInputVisible"
                  ref="domainInputRef"
                  v-model="domainInputValue"
                  size="small"
                  @keyup.enter="confirmDomainInput"
                  @blur="confirmDomainInput"
                  placeholder="输入域名"
                  class="domain-input"
                />
                <el-button
                  v-else
                  size="small"
                  @click="showDomainInput"
                  class="add-domain-btn"
                >
                  <el-icon><Plus /></el-icon>
                  添加域名
                </el-button>
              </div>
            </el-form-item>

            <!-- 验证方式 -->
            <el-form-item label="验证方式" prop="challengeType">
              <el-radio-group v-model="formData.challengeType">
                <el-radio label="dns-01">
                  <div class="challenge-option">
                    <div class="option-header">
                      <span class="option-title">DNS验证</span>
                      <el-tag type="success" size="small">推荐</el-tag>
                    </div>
                    <div class="option-desc">
                      通过DNS TXT记录验证域名所有权，支持通配符证书
                    </div>
                  </div>
                </el-radio>
                <el-radio label="http-01" :disabled="formData.certType === 'wildcard'">
                  <div class="challenge-option">
                    <div class="option-header">
                      <span class="option-title">HTTP验证</span>
                      <el-tag v-if="formData.certType === 'wildcard'" type="info" size="small">不支持通配符</el-tag>
                    </div>
                    <div class="option-desc">
                      通过HTTP文件验证域名所有权，需要域名可访问
                    </div>
                  </div>
                </el-radio>
              </el-radio-group>
            </el-form-item>

            <!-- 密钥类型 -->
            <el-form-item label="密钥类型" prop="keyType">
              <el-select v-model="formData.keyType" style="width: 200px">
                <el-option label="RSA 2048" value="RSA2048" />
                <el-option label="RSA 4096" value="RSA4096" />
                <el-option label="ECDSA P-256" value="ECDSA256" />
                <el-option label="ECDSA P-384" value="ECDSA384" />
              </el-select>
              <div class="form-tip">
                <el-icon><InfoFilled /></el-icon>
                推荐使用RSA 2048，兼容性最好
              </div>
            </el-form-item>

            <!-- 邮箱地址 -->
            <el-form-item label="邮箱地址" prop="email">
              <el-input
                v-model="formData.email"
                placeholder="用于接收证书相关通知"
                style="width: 300px"
              />
              <div class="form-tip">
                <el-icon><InfoFilled /></el-icon>
                用于Let's Encrypt账户注册和重要通知
              </div>
            </el-form-item>
          </template>

          <!-- 高级设置 -->
          <el-divider content-position="left">
            <el-icon><Setting /></el-icon>
            高级设置
          </el-divider>

          <!-- 自动续期 -->
          <el-form-item label="自动续期">
            <el-switch
              v-model="formData.autoRenew"
              active-text="开启"
              inactive-text="关闭"
            />
            <div class="form-tip">
              <el-icon><InfoFilled /></el-icon>
              开启后将在证书到期前30天自动续期
            </div>
          </el-form-item>

          <!-- 部署主机 -->
          <el-form-item label="部署主机">
            <el-select
              v-model="formData.deployHosts"
              multiple
              placeholder="选择要部署证书的主机（可选）"
              style="width: 100%"
              clearable
            >
              <el-option
                v-for="host in hostOptions"
                :key="host.id"
                :label="`${host.name} (${host.ip})`"
                :value="host.id"
              >
                <div class="host-option">
                  <span class="host-name">{{ host.name }}</span>
                  <span class="host-ip">{{ host.ip }}</span>
                  <el-tag
                    :type="host.status === 'online' ? 'success' : 'danger'"
                    size="small"
                  >
                    {{ host.status === 'online' ? '在线' : '离线' }}
                  </el-tag>
                </div>
              </el-option>
            </el-select>
            <div class="form-tip">
              <el-icon><InfoFilled /></el-icon>
              证书申请成功后将自动部署到选择的主机
            </div>
          </el-form-item>

          <!-- 备注信息 -->
          <el-form-item label="备注信息">
            <el-input
              v-model="formData.remark"
              type="textarea"
              :rows="3"
              placeholder="请输入证书用途或备注信息（可选）"
              maxlength="200"
              show-word-limit
            />
          </el-form-item>

          <!-- 证书上传区域 -->
          <template v-if="formData.applyType === 'upload'">
            <el-form-item label="证书文件" required>
              <el-upload
                ref="certUploadRef"
                :auto-upload="false"
                :show-file-list="false"
                :on-change="handleCertFileChange"
                accept=".crt,.pem,.cer"
                drag
              >
                <el-icon class="el-icon--upload"><upload-filled /></el-icon>
                <div class="el-upload__text">
                  拖拽证书文件到此处或<em>点击上传</em>
                </div>
                <template #tip>
                  <div class="el-upload__tip">
                    支持.crt、.pem、.cer格式的证书文件
                  </div>
                </template>
              </el-upload>
            </el-form-item>

            <el-form-item label="私钥文件" required>
              <el-upload
                ref="keyUploadRef"
                :auto-upload="false"
                :show-file-list="false"
                :on-change="handleKeyFileChange"
                accept=".key,.pem"
                drag
              >
                <el-icon class="el-icon--upload"><upload-filled /></el-icon>
                <div class="el-upload__text">
                  拖拽私钥文件到此处或<em>点击上传</em>
                </div>
                <template #tip>
                  <div class="el-upload__tip">
                    支持.key、.pem格式的私钥文件
                  </div>
                </template>
              </el-upload>
            </el-form-item>

            <el-form-item label="证书链文件">
              <el-upload
                ref="chainUploadRef"
                :auto-upload="false"
                :show-file-list="false"
                :on-change="handleChainFileChange"
                accept=".crt,.pem,.cer"
                drag
              >
                <el-icon class="el-icon--upload"><upload-filled /></el-icon>
                <div class="el-upload__text">
                  拖拽证书链文件到此处或<em>点击上传</em>
                </div>
                <template #tip>
                  <div class="el-upload__tip">
                    可选：支持.crt、.pem、.cer格式的证书链文件
                  </div>
                </template>
              </el-upload>
            </el-form-item>

            <!-- 证书内容预览 -->
            <div v-if="formData.certContent || formData.keyContent" class="cert-preview">
              <el-form-item label="证书内容" v-if="formData.certContent">
                <el-input
                  v-model="formData.certContent"
                  type="textarea"
                  :rows="6"
                  readonly
                  class="cert-content"
                />
              </el-form-item>

              <el-form-item label="私钥内容" v-if="formData.keyContent">
                <el-input
                  v-model="formData.keyContent"
                  type="textarea"
                  :rows="6"
                  readonly
                  class="key-content"
                  show-password
                />
              </el-form-item>

              <el-form-item label="证书链内容" v-if="formData.chainContent">
                <el-input
                  v-model="formData.chainContent"
                  type="textarea"
                  :rows="6"
                  readonly
                  class="chain-content"
                />
              </el-form-item>
            </div>

            <!-- 证书验证结果 -->
            <div v-if="certValidation" class="cert-validation">
              <el-alert
                :type="certValidation.valid ? 'success' : 'error'"
                :title="certValidation.valid ? '证书验证成功' : '证书验证失败'"
                :description="certValidation.message"
                show-icon
                :closable="false"
              />

              <el-descriptions v-if="certValidation.valid && certValidation.details"
                title="证书信息" :column="2" border class="mt-4">
                <el-descriptions-item label="通用名称">
                  {{ certValidation.details.commonName }}
                </el-descriptions-item>
                <el-descriptions-item label="备用名称">
                  {{ certValidation.details.subjectAltNames?.join(', ') || '无' }}
                </el-descriptions-item>
                <el-descriptions-item label="颁发者">
                  {{ certValidation.details.issuer }}
                </el-descriptions-item>
                <el-descriptions-item label="有效期">
                  {{ certValidation.details.validFrom }} 至 {{ certValidation.details.validTo }}
                </el-descriptions-item>
                <el-descriptions-item label="剩余天数">
                  {{ certValidation.details.daysRemaining }} 天
                </el-descriptions-item>
                <el-descriptions-item label="密钥算法">
                  {{ certValidation.details.keyAlgorithm }}
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </div>

          <!-- CSR上传区域 -->
          <div v-if="formData.applyType === 'csr'" class="csr-upload-section">
            <el-form-item label="CSR文件">
              <el-upload
                ref="csrUploadRef"
                :auto-upload="false"
                :show-file-list="false"
                :on-change="handleCSRFileChange"
                accept=".csr,.pem,.txt"
                drag
              >
                <el-icon class="el-icon--upload"><upload-filled /></el-icon>
                <div class="el-upload__text">
                  拖拽CSR文件到此处或<em>点击上传</em>
                </div>
                <template #tip>
                  <div class="el-upload__tip">
                    支持.csr、.pem、.txt格式的文件
                  </div>
                </template>
              </el-upload>
            </el-form-item>

            <el-form-item label="CSR内容" v-if="formData.csrContent">
              <el-input
                v-model="formData.csrContent"
                type="textarea"
                :rows="8"
                readonly
                class="csr-content"
              />
            </el-form-item>

            <!-- CSR验证结果 -->
            <div v-if="csrValidation" class="csr-validation">
              <el-alert 
                :type="csrValidation.valid ? 'success' : 'error'"
                :title="csrValidation.valid ? 'CSR验证成功' : 'CSR验证失败'"
                :description="csrValidation.message"
                show-icon
                :closable="false"
              />
              
              <el-descriptions v-if="csrValidation.valid && csrValidation.details" 
                title="CSR信息" :column="2" border class="mt-4">
                <el-descriptions-item label="通用名称">
                  {{ csrValidation.details.commonName }}
                </el-descriptions-item>
                <el-descriptions-item label="备用名称">
                  {{ csrValidation.details.subjectAltNames?.join(', ') || '-' }}
                </el-descriptions-item>
                <el-descriptions-item label="密钥类型">
                  {{ csrValidation.details.keyType }}
                </el-descriptions-item>
                <el-descriptions-item label="密钥长度">
                  {{ csrValidation.details.keySize }}位
                </el-descriptions-item>
                <el-descriptions-item label="签名算法">
                  {{ csrValidation.details.signatureAlgorithm }}
                </el-descriptions-item>
                <el-descriptions-item label="组织">
                  {{ csrValidation.details.organization || '-' }}
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </div>
        </el-form>
      </div>

      <!-- 步骤2：配置证书信息 -->
      <div v-show="currentStep === 1" class="step-content">
        <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px">
          <el-form-item label="域名" prop="domainId" v-if="formData.applyType === 'auto'">
            <el-select 
              v-model="formData.domainId" 
              placeholder="请选择域名"
              filterable
              @change="handleDomainChange"
            >
              <el-option 
                v-for="domain in domainOptions" 
                :key="domain.id"
                :label="domain.name" 
                :value="domain.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="证书类型" v-if="formData.applyType === 'auto'">
            <el-radio-group v-model="formData.certType" @change="handleCertTypeChange">
              <el-radio label="single">单域名证书</el-radio>
              <el-radio label="wildcard">通配符证书</el-radio>
              <el-radio label="multi">多域名证书</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item 
            label="域名列表" 
            prop="domains" 
            v-if="formData.applyType === 'auto' && formData.certType === 'multi'"
          >
            <el-tag
              v-for="domain in formData.domains"
              :key="domain"
              closable
              @close="removeDomain(domain)"
              class="mr-2 mb-2"
            >
              {{ domain }}
            </el-tag>
            <el-input
              v-if="domainInputVisible"
              ref="domainInputRef"
              v-model="domainInputValue"
              class="domain-input"
              @keyup.enter="handleDomainInputConfirm"
              @blur="handleDomainInputConfirm"
            />
            <el-button v-else @click="showDomainInput" text>
              <el-icon><Plus /></el-icon>
              添加域名
            </el-button>
          </el-form-item>

          <el-form-item label="CA提供商">
            <el-select v-model="formData.caType" placeholder="请选择CA提供商">
              <el-option label="Let's Encrypt" value="letsencrypt" />
              <el-option label="ZeroSSL" value="zerossl" disabled />
              <el-option label="BuyPass" value="buypass" disabled />
            </el-select>
            <div class="form-tip">
              目前仅支持Let's Encrypt，其他CA提供商即将上线
            </div>
          </el-form-item>

          <el-form-item label="DNS提供商" v-if="formData.applyType === 'auto'">
            <el-select v-model="formData.providerId" placeholder="请选择DNS提供商">
              <el-option 
                v-for="provider in providerOptions" 
                :key="provider.id"
                :label="provider.name" 
                :value="provider.id"
              />
            </el-select>
            <div class="form-tip">
              用于DNS-01验证，自动创建TXT记录
            </div>
          </el-form-item>

          <el-form-item label="密钥算法" v-if="formData.applyType === 'auto'">
            <el-select v-model="formData.keyType">
              <el-option label="RSA 2048" value="rsa2048" />
              <el-option label="RSA 4096" value="rsa4096" />
              <el-option label="ECDSA P-256" value="ecdsa256" />
              <el-option label="ECDSA P-384" value="ecdsa384" />
            </el-select>
          </el-form-item>

          <el-form-item label="有效期">
            <el-select v-model="formData.validDays">
              <el-option label="90天（推荐）" :value="90" />
              <el-option label="60天" :value="60" />
              <el-option label="30天" :value="30" />
            </el-select>
            <div class="form-tip">
              Let's Encrypt证书最长有效期为90天
            </div>
          </el-form-item>

          <el-form-item label="自动续期">
            <el-switch 
              v-model="formData.autoRenew" 
              active-text="开启"
              inactive-text="关闭"
            />
            <div class="form-tip">
              在证书到期前30天自动续期
            </div>
          </el-form-item>
        </el-form>
      </div>

      <!-- 步骤3：高级设置 -->
      <div v-show="currentStep === 2" class="step-content">
        <el-form :model="formData" label-width="120px">
          <el-form-item label="自动部署">
            <el-switch 
              v-model="formData.autoDeploy" 
              active-text="开启"
              inactive-text="关闭"
            />
          </el-form-item>

          <div v-if="formData.autoDeploy" class="deploy-settings">
            <el-form-item label="部署主机">
              <el-select 
                v-model="formData.deployHosts" 
                multiple 
                placeholder="请选择部署主机"
                filterable
              >
                <el-option 
                  v-for="host in hostOptions" 
                  :key="host.id"
                  :label="`${host.name} (${host.ip})`" 
                  :value="host.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="部署路径">
              <el-input v-model="formData.deployPath" placeholder="/etc/ssl/certs/" />
            </el-form-item>

            <el-form-item label="重启服务">
              <el-input v-model="formData.restartCommand" placeholder="systemctl reload nginx" />
              <div class="form-tip">
                证书部署后执行的重启命令
              </div>
            </el-form-item>
          </div>

          <el-form-item label="邮件通知">
            <el-switch 
              v-model="formData.emailNotification" 
              active-text="开启"
              inactive-text="关闭"
            />
          </el-form-item>

          <el-form-item label="通知邮箱" v-if="formData.emailNotification">
            <el-input v-model="formData.notificationEmail" placeholder="admin@example.com" />
          </el-form-item>

          <el-form-item label="备注">
            <el-input
              v-model="formData.remark"
              type="textarea"
              :rows="3"
              placeholder="请输入备注信息"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>
        </el-form>
      </div>

      <!-- 步骤4：确认申请 -->
      <div v-show="currentStep === 3" class="step-content">
        <el-alert
          title="请确认证书申请信息"
          type="info"
          :closable="false"
          class="mb-4"
        />
        
        <el-descriptions title="证书信息" :column="2" border>
          <el-descriptions-item label="申请方式">
            {{
              formData.applyType === 'auto' ? '自动申请' :
              formData.applyType === 'upload' ? '导入证书' : '自定义CSR'
            }}
          </el-descriptions-item>
          <el-descriptions-item label="CA提供商">{{ formData.caType }}</el-descriptions-item>
          <el-descriptions-item label="证书类型" v-if="formData.applyType === 'auto'">
            {{ getCertTypeText(formData.certType) }}
          </el-descriptions-item>
          <el-descriptions-item label="域名">
            <div v-if="formData.applyType === 'auto'">
              <el-tag v-for="domain in formData.domains" :key="domain" class="mr-1">
                {{ domain }}
              </el-tag>
            </div>
            <div v-else-if="formData.applyType === 'upload'">
              {{ certValidation?.details?.commonName || '从证书中解析' }}
            </div>
            <div v-else>
              {{ csrValidation?.details?.commonName }}
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="有效期">{{ formData.validDays }}天</el-descriptions-item>
          <el-descriptions-item label="自动续期">
            {{ formData.autoRenew ? '是' : '否' }}
          </el-descriptions-item>
          <el-descriptions-item label="自动部署">
            {{ formData.autoDeploy ? '是' : '否' }}
          </el-descriptions-item>
          <el-descriptions-item label="邮件通知">
            {{ formData.emailNotification ? '是' : '否' }}
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- 步骤控制按钮 -->
      <div class="step-controls">
        <el-button v-if="currentStep > 0" @click="prevStep">上一步</el-button>
        <el-button v-if="currentStep < 3" type="primary" @click="nextStep">下一步</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElForm, ElInput } from 'element-plus'
import { Lightning, Upload, UploadFilled, Plus, FolderAdd } from '@element-plus/icons-vue'
import { certificateApi } from '@/api/dns/certificate'
import { domainApi } from '@/api/dns/domain'
import { providerApi } from '@/api/dns/provider'
import { hostApi } from '@/api/system/host'

const route = useRoute()
const router = useRouter()
const formRef = ref<InstanceType<typeof ElForm>>()
const csrUploadRef = ref()
const domainInputRef = ref<InstanceType<typeof ElInput>>()

const currentStep = ref(0)
const applying = ref(false)

// 表单数据
const formData = reactive({
  applyType: 'auto', // auto | upload | csr
  domainId: null as number | null,
  certType: 'single', // single | wildcard | multi
  domains: [] as string[],
  caType: 'letsencrypt',
  providerId: null as number | null,
  keyType: 'rsa2048',
  validDays: 90,
  autoRenew: true,
  autoDeploy: false,
  deployHosts: [] as number[],
  deployPath: '/etc/ssl/certs/',
  restartCommand: 'systemctl reload nginx',
  emailNotification: false,
  notificationEmail: '',
  remark: '',
  csrContent: '',
  // 证书上传相关
  certContent: '',
  keyContent: '',
  chainContent: ''
})

// CSR验证结果
const csrValidation = ref<any>(null)

// 证书验证结果
const certValidation = ref<any>(null)

// 域名输入
const domainInputVisible = ref(false)
const domainInputValue = ref('')

// 选项数据
const domainOptions = ref([])
const providerOptions = ref([])
const hostOptions = ref([])

// 表单验证规则
const formRules = {
  domainId: [
    { required: true, message: '请选择域名', trigger: 'change' }
  ],
  domains: [
    { 
      validator: (rule: any, value: any, callback: any) => {
        if (formData.certType === 'multi' && (!value || value.length === 0)) {
          callback(new Error('请至少添加一个域名'))
        } else {
          callback()
        }
      }, 
      trigger: 'change' 
    }
  ]
}

// 获取域名选项
const fetchDomainOptions = async () => {
  try {
    const response = await domainApi.list({
      page: 1,
      pageSize: 1000, // 获取所有域名
      status: 'active' // 只获取活跃域名
    })
    domainOptions.value = response.data?.list || []

    // 如果URL中有domainId参数，自动选择
    const domainId = route.query.domainId
    if (domainId) {
      formData.domainId = Number(domainId)
      handleDomainChange(formData.domainId)
    }
  } catch (error) {
    console.error('获取域名列表失败:', error)
    ElMessage.error('获取域名列表失败')
  }
}

// 获取提供商选项
const fetchProviderOptions = async () => {
  try {
    const response = await providerApi.list({
      page: 1,
      pageSize: 1000, // 获取所有提供商
      status: 'active' // 只获取活跃提供商
    })
    providerOptions.value = response.data?.list || []
  } catch (error) {
    console.error('获取提供商列表失败:', error)
    ElMessage.error('获取DNS提供商列表失败')
  }
}

// 获取主机选项
const fetchHostOptions = async () => {
  try {
    const response = await hostApi.list({
      page: 1,
      pageSize: 1000, // 获取所有主机
      status: 'online' // 只获取在线主机
    })
    hostOptions.value = response.data?.list || []
  } catch (error) {
    console.error('获取主机列表失败:', error)
    ElMessage.error('获取主机列表失败')
  }
}

// 处理申请方式变化
const handleApplyTypeChange = (value: string) => {
  if (value === 'auto') {
    formData.csrContent = ''
    csrValidation.value = null
  }
}

// 处理CSR文件变化
const handleCSRFileChange = async (file: any) => {
  try {
    const content = await file.raw.text()
    formData.csrContent = content

    // 验证CSR
    const response = await certificateApi.validateCSR({ csrContent: content })
    csrValidation.value = response.data

    if (response.data.valid) {
      ElMessage.success('CSR文件验证成功')
    }
  } catch (error: any) {
    ElMessage.error('CSR文件验证失败: ' + (error.message || '未知错误'))
    csrValidation.value = { valid: false, message: error.message || '验证失败' }
  }
}

// 处理证书文件变化
const handleCertFileChange = async (file: any) => {
  try {
    const content = await file.raw.text()
    formData.certContent = content

    // 验证证书
    await validateCertificate()

    ElMessage.success('证书文件上传成功')
  } catch (error: any) {
    ElMessage.error('证书文件上传失败: ' + (error.message || '未知错误'))
  }
}

// 处理私钥文件变化
const handleKeyFileChange = async (file: any) => {
  try {
    const content = await file.raw.text()
    formData.keyContent = content

    // 验证证书
    await validateCertificate()

    ElMessage.success('私钥文件上传成功')
  } catch (error: any) {
    ElMessage.error('私钥文件上传失败: ' + (error.message || '未知错误'))
  }
}

// 处理证书链文件变化
const handleChainFileChange = async (file: any) => {
  try {
    const content = await file.raw.text()
    formData.chainContent = content

    // 验证证书
    await validateCertificate()

    ElMessage.success('证书链文件上传成功')
  } catch (error: any) {
    ElMessage.error('证书链文件上传失败: ' + (error.message || '未知错误'))
  }
}

// 验证证书
const validateCertificate = async () => {
  if (!formData.certContent || !formData.keyContent) {
    return
  }

  try {
    const response = await certificateApi.validateCertificate({
      certContent: formData.certContent,
      keyContent: formData.keyContent,
      chainContent: formData.chainContent
    })
    certValidation.value = response.data

    if (response.data.valid) {
      ElMessage.success('证书验证成功')
    } else {
      ElMessage.error('证书验证失败: ' + response.data.message)
    }
  } catch (error: any) {
    ElMessage.error('证书验证失败: ' + (error.message || '未知错误'))
    certValidation.value = { valid: false, message: error.message || '验证失败' }
  }
}

// 处理申请方式变化
const handleApplyTypeChange = (value: string) => {
  // 清空相关数据
  formData.csrContent = ''
  formData.certContent = ''
  formData.keyContent = ''
  formData.chainContent = ''
  csrValidation.value = null
  certValidation.value = null
}

// 处理域名选择变化
const handleDomainChange = (domainId: number) => {
  const domain = domainOptions.value.find((d: any) => d.id === domainId)
  if (domain) {
    formData.providerId = domain.provider_id || domain.providerId
    updateDomainsForCertType(domain.name)
  }
}

// 处理证书类型变化
const handleCertTypeChange = (certType: string) => {
  if (formData.domainId) {
    const domain = domainOptions.value.find((d: any) => d.id === formData.domainId)
    if (domain) {
      updateDomainsForCertType(domain.name)
    }
  }
}

// 根据证书类型更新域名列表
const updateDomainsForCertType = (domainName: string) => {
  if (formData.certType === 'wildcard') {
    formData.domains = [`*.${domainName}`]
  } else if (formData.certType === 'single') {
    formData.domains = [domainName]
  } else if (formData.certType === 'multi') {
    // 多域名证书保持当前域名列表，或初始化为主域名
    if (formData.domains.length === 0) {
      formData.domains = [domainName]
    }
  }
}

// 显示域名输入框
const showDomainInput = () => {
  domainInputVisible.value = true
  nextTick(() => {
    domainInputRef.value?.focus()
  })
}

// 确认域名输入
const handleDomainInputConfirm = () => {
  if (domainInputValue.value) {
    if (!formData.domains.includes(domainInputValue.value)) {
      formData.domains.push(domainInputValue.value)
    }
    domainInputValue.value = ''
  }
  domainInputVisible.value = false
}

// 移除域名
const removeDomain = (domain: string) => {
  const index = formData.domains.indexOf(domain)
  if (index > -1) {
    formData.domains.splice(index, 1)
  }
}

// 获取证书类型文本
const getCertTypeText = (type: string) => {
  const typeMap: Record<string, string> = {
    'single': '单域名证书',
    'wildcard': '通配符证书',
    'multi': '多域名证书'
  }
  return typeMap[type] || type
}

// 步骤控制
const nextStep = async () => {
  if (currentStep.value === 1 && formRef.value) {
    try {
      await formRef.value.validate()
    } catch {
      return
    }
  }
  
  if (currentStep.value < 3) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

// 申请证书
const handleApply = async () => {
  try {
    applying.value = true

    if (formData.applyType === 'upload') {
      // 上传证书
      if (!formData.certContent || !formData.keyContent) {
        ElMessage.error('请上传证书文件和私钥文件')
        return
      }
      if (!certValidation.value?.valid) {
        ElMessage.error('证书验证失败，请检查文件格式')
        return
      }

      await certificateApi.upload({
        cert_content: formData.certContent,
        key_content: formData.keyContent,
        chain_content: formData.chainContent,
        auto_deploy: formData.autoDeploy,
        deploy_hosts: formData.deployHosts,
        deploy_path: formData.deployPath,
        restart_command: formData.restartCommand,
        email_notification: formData.emailNotification,
        notification_email: formData.notificationEmail,
        remark: formData.remark
      })

      ElMessage.success('证书上传成功')
    } else if (formData.applyType === 'csr') {
      // 使用CSR申请
      if (!formData.csrContent) {
        ElMessage.error('请上传CSR文件')
        return
      }
      if (!csrValidation.value?.valid) {
        ElMessage.error('CSR验证失败，请检查文件格式')
        return
      }

      await certificateApi.createWithCSR({
        domains: formData.domains,
        email: formData.notificationEmail || 'admin@example.com',
        provider_id: formData.providerId,
        csr_content: formData.csrContent,
        auto_renew: formData.autoRenew,
        deploy_hosts: formData.deployHosts,
        remark: formData.remark
      })

      ElMessage.success('证书申请已提交，请稍后查看申请状态')
    } else {
      // 自动申请
      await certificateApi.create({
        domain_id: formData.domainId,
        domains: formData.domains,
        email: formData.notificationEmail || 'admin@example.com',
        provider_id: formData.providerId,
        key_type: formData.keyType,
        valid_days: formData.validDays,
        auto_renew: formData.autoRenew,
        deploy_hosts: formData.deployHosts,
        remark: formData.remark
      })

      ElMessage.success('证书申请已提交，请稍后查看申请状态')
    }

    router.push('/dns/certs')
  } catch (error: any) {
    ElMessage.error('操作失败: ' + (error.message || '未知错误'))
  } finally {
    applying.value = false
  }
}

// 返回
const goBack = () => {
  router.back()
}

onMounted(() => {
  fetchDomainOptions()
  fetchProviderOptions()
  fetchHostOptions()

  // 检查URL参数，自动选择申请方式
  const type = route.query.type
  if (type === 'upload') {
    formData.applyType = 'upload'
  } else if (type === 'csr') {
    formData.applyType = 'csr'
  }
})
</script>

<style scoped>
.cert-create {
  padding: 20px;
}

.step-content {
  min-height: 400px;
  padding: 20px 0;
}

.step-controls {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--el-border-color-light);
}

.radio-card {
  display: block;
  width: 100%;
  margin-bottom: 16px;
}

.radio-card :deep(.el-radio__label) {
  width: 100%;
  padding: 16px;
  border: 1px solid var(--el-border-color);
  border-radius: 8px;
  transition: all 0.3s;
}

.radio-card :deep(.el-radio__input.is-checked + .el-radio__label) {
  border-color: var(--el-color-primary);
  background-color: var(--el-color-primary-light-9);
}

.radio-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.radio-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.radio-desc {
  color: var(--el-text-color-regular);
  font-size: 14px;
}

.cert-upload-section {
  margin-top: 24px;
  padding: 24px;
  background-color: var(--el-bg-color-page);
  border-radius: 8px;
}

.cert-preview {
  margin-top: 16px;
}

.cert-content,
.key-content,
.chain-content {
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.cert-validation {
  margin-top: 16px;
}

.csr-upload-section {
  margin-top: 24px;
  padding: 24px;
  background-color: var(--el-bg-color-page);
  border-radius: 8px;
}

.csr-content {
  font-family: 'Courier New', monospace;
}

.csr-validation {
  margin-top: 16px;
}

.deploy-settings {
  margin-left: 24px;
  padding-left: 24px;
  border-left: 2px solid var(--el-color-primary-light-8);
}

.domain-input {
  width: 120px;
  margin-right: 8px;
}

.form-tip {
  font-size: 12px;
  color: var(--el-color-info);
  margin-top: 4px;
}

.mt-4 {
  margin-top: 1rem;
}

.mb-4 {
  margin-bottom: 1rem;
}

.mb-6 {
  margin-bottom: 1.5rem;
}

.mr-1 {
  margin-right: 0.25rem;
}

.mr-2 {
  margin-right: 0.5rem;
}

.mb-2 {
  margin-bottom: 0.5rem;
}
</style> 
