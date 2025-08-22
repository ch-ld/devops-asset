<template>
  <div class="modern-page-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <div class="title-icon">
              <el-icon><Lock /></el-icon>
            </div>
            SSL证书管理
          </h1>
          <p class="page-description">管理您的SSL证书，监控证书状态，确保网站安全访问</p>
        </div>
        <div class="header-actions">
          <el-button class="modern-btn secondary" @click="handleRefresh" :icon="Refresh">
            刷新
          </el-button>
          <el-dropdown @command="handleBulkAction">
            <el-button class="modern-btn secondary">
              批量操作
              <el-icon><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="deploy">批量部署</el-dropdown-item>
                <el-dropdown-item command="renew">批量续期</el-dropdown-item>
                <el-dropdown-item command="export">批量导出</el-dropdown-item>
                <el-dropdown-item command="delete" divided>批量删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-button class="modern-btn primary" @click="handleAddCert" :icon="Plus">
            申请证书
          </el-button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="modern-stats-grid">
      <div class="stat-card">
        <div class="stat-header">
          <div class="stat-icon primary">
            <el-icon><Lock /></el-icon>
          </div>
          <div class="stat-trend up">+12%</div>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ statistics.total || 0 }}</div>
          <div class="stat-label">总证书数</div>
          <div class="stat-description">所有SSL/TLS证书</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-header">
          <div class="stat-icon success">
            <el-icon><Check /></el-icon>
          </div>
          <div class="stat-trend up">+5%</div>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ statistics.valid || 0 }}</div>
          <div class="stat-label">有效证书</div>
          <div class="stat-description">正常使用中</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-header">
          <div class="stat-icon warning">
            <el-icon><Clock /></el-icon>
          </div>
          <div class="stat-trend down">-2</div>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ statistics.expiring || 0 }}</div>
          <div class="stat-label">即将过期</div>
          <div class="stat-description">30天内过期</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-header">
          <div class="stat-icon error">
            <el-icon><Warning /></el-icon>
          </div>
          <div class="stat-trend down">-1</div>
        </div>
        <div class="stat-content">
          <div class="stat-number">{{ statistics.expired || 0 }}</div>
          <div class="stat-label">已过期</div>
          <div class="stat-description">需要立即处理</div>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="modern-search-section">
      <div class="search-content">
        <el-input
          v-model="searchForm.keyword"
          placeholder="搜索证书域名..."
          size="large"
          clearable
          class="search-input"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <div class="search-filters">
          <el-select
            v-model="searchForm.status"
            placeholder="证书状态"
            size="large"
            clearable
            style="width: 150px"
          >
            <el-option label="全部" value="" />
            <el-option label="申请中" value="pending" />
            <el-option label="处理中" value="processing" />
            <el-option label="验证中" value="validating" />
            <el-option label="已签发" value="issued" />
            <el-option label="已过期" value="expired" />
            <el-option label="申请失败" value="failed" />
            <el-option label="已吊销" value="revoked" />
          </el-select>
          <el-select
            v-model="searchForm.ca_type"
            placeholder="CA类型"
            size="large"
            clearable
            style="width: 150px"
          >
            <el-option label="全部" value="" />
            <el-option label="Let's Encrypt" value="letsencrypt" />
            <el-option label="ZeroSSL" value="zerossl" />
            <el-option label="自定义" value="custom" />
          </el-select>
          <el-button class="modern-btn primary" @click="handleSearch" :icon="Search">
            搜索
          </el-button>
          <el-button class="modern-btn secondary" @click="handleReset" :icon="Refresh">
            重置
          </el-button>
        </div>
      </div>
    </div>

    <!-- 证书列表 -->
    <div class="modern-content-card">
      <div class="card-header">
        <div class="header-content">
          <div class="header-left">
            <h3 class="card-title">证书列表</h3>
            <p class="card-subtitle">{{ certificates.length }} 个证书</p>
          </div>
          <div class="header-actions">
            <el-button class="modern-btn secondary" @click="handleRefresh" :icon="Refresh">
              刷新
            </el-button>
            <el-button class="modern-btn secondary" @click="handleImportCert" :icon="Upload">
              导入证书
            </el-button>
            <el-dropdown @command="handleBatchCommand" :disabled="!hasSelected">
              <el-button class="modern-btn warning" :disabled="!hasSelected">
                批量操作 ({{ selectedRows.length }})
                <el-icon><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="renew">批量续期</el-dropdown-item>
                  <el-dropdown-item command="download">批量下载</el-dropdown-item>
                  <el-dropdown-item command="deploy">批量部署</el-dropdown-item>
                  <el-dropdown-item command="export">导出报告</el-dropdown-item>
                  <el-dropdown-item command="delete" divided>批量删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </div>
      <div class="card-content">
        <div class="modern-table">
          <el-table
            ref="tableRef"
            :data="certificates"
            :loading="loading"
            @selection-change="handleSelectionChange"
            row-key="id"
          >
            <el-table-column type="selection" width="55" />
            <el-table-column label="域名" min-width="200">
              <template #default="{ row }">
                <div class="domain-cell">
                  <div class="domain-name">{{ row.common_name || row.domain_name || '-' }}</div>
                  <div class="cert-type" v-if="row.subject_alt_names && row.subject_alt_names.length > 0">
                    <el-tag
                      v-for="san in row.subject_alt_names.slice(0, 2)"
                      :key="san"
                      type="info"
                      size="small"
                      style="margin-right: 4px;"
                    >
                      {{ san }}
                    </el-tag>
                    <el-tag
                      v-if="row.subject_alt_names.length > 2"
                      type="info"
                      size="small"
                    >
                      +{{ row.subject_alt_names.length - 2 }}
                    </el-tag>
                  </div>
                  <div class="cert-type" v-if="row.common_name && row.common_name.startsWith('*.')">
                    <el-tag type="success" size="small">泛域名</el-tag>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="120">
              <template #default="{ row }">
                <div class="status-cell">
                  <el-tag
                    :type="getStatusTagType(row.status)"
                    size="small"
                  >
                    {{ getStatusText(row.status) }}
                  </el-tag>
                  <div v-if="row.status === 'failed' && row.error_message" class="error-message">
                    <el-tooltip :content="row.error_message" placement="top">
                      <el-text type="danger" size="small">
                        {{ row.error_message }}
                      </el-text>
                    </el-tooltip>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="ca_type" label="CA类型" width="120">
              <template #default="{ row }">
                <el-tag type="info" size="small">
                  {{ getCATypeName(row.ca_type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="issued_at" label="签发时间" width="140">
              <template #default="{ row }">
                {{ formatDate(row.issued_at) }}
              </template>
            </el-table-column>
            <el-table-column prop="expires_at" label="过期时间" width="180">
              <template #default="{ row }">
                <div v-if="row.expires_at" class="expires-cell">
                  <div class="expires-date">{{ formatDate(row.expires_at) }}</div>
                  <div class="expires-days" :class="getExpiresClass(row.expires_at)">
                    {{ getExpiringDays(row.expires_at) }}
                  </div>
                </div>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column prop="auto_renew" label="自动续期" width="100">
              <template #default="{ row }">
                <el-switch
                  :model-value="row.auto_renew"
                  @change="val => handleAutoRenewChange(row, Boolean(val))"
                />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="250" fixed="right">
              <template #default="{ row }">
                <el-button
                  size="small"
                  @click="handleViewCert(row)"
                >
                  查看
                </el-button>
                <el-button
                  v-if="row.status === 'issued'"
                  size="small"
                  type="primary"
                  @click="handleDeployCert(row)"
                >
                  部署
                </el-button>
                <el-button
                  v-if="row.status === 'issued'"
                  size="small"
                  @click="handleDownloadCert(row)"
                >
                  下载
                </el-button>
                <el-dropdown @command="cmd => handleCertAction(cmd, row)" trigger="click">
                  <el-button size="small">
                    更多
                    <el-icon><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="renew">手动续期</el-dropdown-item>
                      <el-dropdown-item command="revoke" v-if="row.status === 'issued'">
                        吊销
                      </el-dropdown-item>
                      <el-dropdown-item command="export">导出配置</el-dropdown-item>
                      <el-dropdown-item command="delete" divided>删除</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </template>
            </el-table-column>
          </el-table>

          <!-- 分页 -->
          <div class="pagination-container">
            <el-pagination
              v-model:current-page="pagination.page"
              v-model:page-size="pagination.pageSize"
              :total="pagination.total"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 证书申请弹窗 -->
    <CertificateModal
      v-model:visible="modalVisible"
      :certificate="currentCertificate"
      @success="handleModalSuccess"
    />

    <!-- 证书详情抽屉 -->
    <CertificateDetail
      v-if="currentCertificate"
      v-model:visible="drawerVisible"
      :certificate="currentCertificate"
    />

    <!-- 证书部署弹窗 -->
    <CertificateDeployDialog
      v-model:visible="showDeployDialog"
      :certificate="selectedCertificate"
    />

    <!-- 证书导出弹窗 -->
    <CertificateExportDialog
      v-model:visible="showExportDialog"
      :certificate="selectedCertificate"
    />
  </div>
</template>
