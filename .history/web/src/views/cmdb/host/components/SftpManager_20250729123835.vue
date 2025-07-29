<template>
  <el-dialog
    v-model="dialogVisible"
    :before-close="handleClose"
    destroy-on-close
    class="modern-sftp-dialog"
    width="95%"
    top="2vh"
    :close-on-click-modal="false"
  >
    <template #header>
      <div class="dialog-header">
        <div class="header-left">
          <div class="header-icon">
            <el-icon size="28"><FolderOpened /></el-icon>
          </div>
          <div class="header-info">
            <h3 class="dialog-title">高级文件管理器</h3>
            <p class="host-info">{{ props.host?.name || '未知主机' }} ({{ getHostIP() }})</p>
            <div class="connection-status">
              <span class="status-dot connected"></span>
              <span class="status-text">SFTP 已连接</span>
            </div>
          </div>
        </div>
        <div class="header-actions">
          <el-tooltip content="刷新文件列表" placement="bottom">
            <el-button circle size="small" @click="refreshFiles" :loading="loading" class="action-btn">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip content="上传文件" placement="bottom">
            <el-button circle size="small" @click="showUploadDialog = true" type="primary" class="action-btn">
              <el-icon><Upload /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip content="新建文件夹" placement="bottom">
            <el-button circle size="small" @click="showCreateDirDialog = true" type="success" class="action-btn">
              <el-icon><FolderAdd /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip content="批量操作" placement="bottom">
            <el-button circle size="small" @click="toggleBatchMode" :type="batchMode ? 'warning' : ''" class="action-btn">
              <el-icon><Operation /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip content="搜索文件" placement="bottom">
            <el-button circle size="small" @click="toggleSearch" class="action-btn">
              <el-icon><Search /></el-icon>
            </el-button>
          </el-tooltip>
        </div>
      </div>
    </template>

    <div class="modern-sftp-manager">
      <!-- 高级工具栏 -->
      <div class="advanced-toolbar">
        <!-- 路径导航 -->
        <div class="path-navigation">
          <el-icon class="path-icon"><Location /></el-icon>
          <div class="breadcrumb-container">
            <span class="path-segment root" @click="navigateToPath(-1)">
              <el-icon><HomeFilled /></el-icon>
              根目录
            </span>
            <template v-for="(segment, index) in pathSegments" :key="index">
              <el-icon class="separator"><ArrowRight /></el-icon>
              <span
                class="path-segment"
                @click="navigateToPath(index)"
                :class="{ active: index === pathSegments.length - 1 }"
              >
                {{ segment }}
              </span>
            </template>
          </div>
        </div>

        <!-- 搜索框 -->
        <div class="search-container" v-show="searchVisible">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索文件..."
            size="small"
            @input="applyFiltersAndSort"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>

        <!-- 工具按钮组 -->
        <div class="toolbar-actions">
          <el-button-group size="small">
            <el-tooltip content="网格视图" placement="bottom">
              <el-button
                :type="viewMode === 'grid' ? 'primary' : ''"
                @click="viewMode = 'grid'"
              >
                <el-icon><Grid /></el-icon>
              </el-button>
            </el-tooltip>
            <el-tooltip content="列表视图" placement="bottom">
              <el-button
                :type="viewMode === 'list' ? 'primary' : ''"
                @click="viewMode = 'list'"
              >
                <el-icon><List /></el-icon>
              </el-button>
            </el-tooltip>
          </el-button-group>

          <el-button-group size="small">
            <el-tooltip content="按名称排序" placement="bottom">
              <el-button
                :type="sortBy === 'name' ? 'primary' : ''"
                @click="sortFiles('name')"
              >
                <el-icon><Sort /></el-icon>
                名称
                <el-icon v-if="sortBy === 'name'">
                  <component :is="sortOrder === 'asc' ? Top : Bottom" />
                </el-icon>
              </el-button>
            </el-tooltip>
            <el-tooltip content="按大小排序" placement="bottom">
              <el-button
                :type="sortBy === 'size' ? 'primary' : ''"
                @click="sortFiles('size')"
              >
                大小
                <el-icon v-if="sortBy === 'size'">
                  <component :is="sortOrder === 'asc' ? Top : Bottom" />
                </el-icon>
              </el-button>
            </el-tooltip>
            <el-tooltip content="按时间排序" placement="bottom">
              <el-button
                :type="sortBy === 'modTime' ? 'primary' : ''"
                @click="sortFiles('modTime')"
              >
                时间
                <el-icon v-if="sortBy === 'modTime'">
                  <component :is="sortOrder === 'asc' ? Top : Bottom" />
                </el-icon>
              </el-button>
            </el-tooltip>
          </el-button-group>

          <el-button-group size="small">
            <el-tooltip content="显示隐藏文件" placement="bottom">
              <el-button
                :type="showHiddenFiles ? 'primary' : ''"
                @click="showHiddenFiles = !showHiddenFiles; applyFiltersAndSort()"
              >
                <el-icon><View /></el-icon>
              </el-button>
            </el-tooltip>
            <el-tooltip content="全屏模式" placement="bottom">
              <el-button @click="toggleFullscreen">
                <el-icon><FullScreen /></el-icon>
              </el-button>
            </el-tooltip>
          </el-button-group>
        </div>
      </div>

      <!-- 批量操作栏 -->
      <div class="batch-toolbar" v-show="batchMode">
        <div class="batch-info">
          <el-checkbox
            :model-value="selectedFiles.size === filteredFileList.length && filteredFileList.length > 0"
            :indeterminate="selectedFiles.size > 0 && selectedFiles.size < filteredFileList.length"
            @change="selectAllFiles"
          >
            已选择 {{ selectedFiles.size }} 个文件
          </el-checkbox>
        </div>
        <div class="batch-actions">
          <el-button size="small" @click="batchDownload" :disabled="selectedFiles.size === 0">
            <el-icon><Download /></el-icon>
            批量下载
          </el-button>
          <el-button size="small" type="danger" @click="batchDelete" :disabled="selectedFiles.size === 0">
            <el-icon><Delete /></el-icon>
            批量删除
          </el-button>
          <el-button size="small" @click="batchMode = false; selectedFiles.clear()">
            <el-icon><CircleClose /></el-icon>
            取消选择
          </el-button>
        </div>
      </div>

      <!-- 文件列表区域 -->
      <div class="file-content" v-loading="loading" element-loading-text="加载中...">
        <!-- 网格视图 -->
        <div v-if="viewMode === 'grid'" class="grid-view">
          <div
            v-for="file in filteredFileList"
            :key="file.name"
            class="file-card"
            :class="{
              selected: selectedFiles.has(file.name),
              favorite: favoriteFiles.has(file.name)
            }"
            @click="batchMode && handleFileSelect(file.name, !selectedFiles.has(file.name))"
            @dblclick="!batchMode && handleRowDoubleClick(file)"
            @contextmenu.prevent="showContextMenu($event, file)"
          >
            <!-- 批量选择复选框 -->
            <div class="file-checkbox" v-show="batchMode">
              <el-checkbox
                :model-value="selectedFiles.has(file.name)"
                @change="(val) => handleFileSelect(file.name, val)"
                @click.stop
              />
            </div>

            <!-- 收藏标记 -->
            <div class="favorite-mark" v-show="favoriteFiles.has(file.name)">
              <el-icon class="favorite-icon"><StarFilled /></el-icon>
            </div>

            <div class="file-icon">
              <el-icon v-if="file.isDir" size="48" class="folder-icon">
                <Folder />
              </el-icon>
              <el-icon v-else size="48" :class="getFileIconClass(file.name)">
                <component :is="getFileIcon(file.name)" />
              </el-icon>
            </div>

            <div class="file-info">
              <div class="file-name" :title="file.name">{{ file.name }}</div>
              <div class="file-meta">
                <span v-if="!file.isDir" class="file-size">{{ formatFileSize(file.size) }}</span>
                <span class="file-time">{{ formatTime(file.modTime) }}</span>
                <span class="file-permissions">{{ file.mode }}</span>
              </div>
            </div>

            <div class="file-actions" v-show="!batchMode">
              <el-button-group size="small">
                <el-tooltip content="预览" placement="top" v-if="!file.isDir">
                  <el-button circle @click.stop="handlePreviewFile(file)">
                    <el-icon><View /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="收藏" placement="top">
                  <el-button
                    circle
                    :type="favoriteFiles.has(file.name) ? 'warning' : ''"
                    @click.stop="toggleFavorite(file.name)"
                  >
                    <el-icon>
                      <component :is="favoriteFiles.has(file.name) ? StarFilled : Star" />
                    </el-icon>
                  </el-button>
                </el-tooltip>
                <el-dropdown trigger="click" @command="handleFileAction" @click.stop>
                  <el-button circle>
                    <el-icon><MoreFilled /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item v-if="!file.isDir" :command="{action: 'download', file}">
                        <el-icon><Download /></el-icon> 下载
                      </el-dropdown-item>
                      <el-dropdown-item v-if="!file.isDir" :command="{action: 'preview', file}">
                        <el-icon><View /></el-icon> 预览
                      </el-dropdown-item>
                      <el-dropdown-item :command="{action: 'rename', file}">
                        <el-icon><Edit /></el-icon> 重命名
                      </el-dropdown-item>
                      <el-dropdown-item :command="{action: 'copy', file}">
                        <el-icon><CopyDocument /></el-icon> 复制
                      </el-dropdown-item>
                      <el-dropdown-item :command="{action: 'properties', file}">
                        <el-icon><InfoFilled /></el-icon> 属性
                      </el-dropdown-item>
                      <el-dropdown-item divided :command="{action: 'delete', file}" class="danger-item">
                        <el-icon><Delete /></el-icon> 删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </el-button-group>
            </div>
          </div>
        </div>

        <!-- 列表视图 -->
        <div v-else class="list-view">
          <el-table
            :data="filteredFileList"
            @row-dblclick="!batchMode && handleRowDoubleClick"
            @row-click="batchMode && ((row) => handleFileSelect(row.name, !selectedFiles.has(row.name)))"
            class="modern-table"
            :header-cell-style="{ background: '#f8fafc', color: '#374151' }"
            :row-class-name="({ row }) => selectedFiles.has(row.name) ? 'selected-row' : ''"
          >
            <!-- 批量选择列 -->
            <el-table-column width="50" align="center" v-if="batchMode">
              <template #header>
                <el-checkbox
                  :model-value="selectedFiles.size === filteredFileList.length && filteredFileList.length > 0"
                  :indeterminate="selectedFiles.size > 0 && selectedFiles.size < filteredFileList.length"
                  @change="selectAllFiles"
                />
              </template>
              <template #default="{ row }">
                <el-checkbox
                  :model-value="selectedFiles.has(row.name)"
                  @change="(val) => handleFileSelect(row.name, val)"
                />
              </template>
            </el-table-column>

            <!-- 文件图标列 -->
            <el-table-column width="60" align="center">
              <template #default="{ row }">
                <div class="file-icon-container">
                  <el-icon v-if="row.isDir" size="24" class="folder-icon">
                    <Folder />
                  </el-icon>
                  <el-icon v-else size="24" :class="getFileIconClass(row.name)">
                    <component :is="getFileIcon(row.name)" />
                  </el-icon>
                  <el-icon
                    v-if="favoriteFiles.has(row.name)"
                    class="favorite-overlay"
                    size="12"
                  >
                    <StarFilled />
                  </el-icon>
                </div>
              </template>
            </el-table-column>

            <!-- 文件名列 -->
            <el-table-column prop="name" label="文件名" min-width="250" sortable="custom">
              <template #default="{ row }">
                <div class="file-name-cell">
                  <span class="name" :class="{ 'hidden-file': row.name.startsWith('.') }">
                    {{ row.name }}
                  </span>
                  <div class="file-badges">
                    <el-tag v-if="row.isDir" size="small" type="primary" class="type-badge">
                      文件夹
                    </el-tag>
                    <el-tag v-else size="small" type="info" class="type-badge">
                      {{ getFileExtension(row.name) || '文件' }}
                    </el-tag>
                    <el-tag v-if="row.name.startsWith('.')" size="small" type="warning" class="hidden-badge">
                      隐藏
                    </el-tag>
                  </div>
                </div>
              </template>
            </el-table-column>

            <!-- 文件大小列 -->
            <el-table-column prop="size" label="大小" width="120" align="right" sortable="custom">
              <template #default="{ row }">
                <span class="file-size">
                  {{ row.isDir ? '-' : formatFileSize(row.size) }}
                </span>
              </template>
            </el-table-column>

            <!-- 权限列 -->
            <el-table-column prop="mode" label="权限" width="100" align="center">
              <template #default="{ row }">
                <el-tooltip :content="getPermissionDescription(row.mode)" placement="top">
                  <el-tag size="small" :type="getPermissionType(row.mode)">
                    {{ row.mode }}
                  </el-tag>
                </el-tooltip>
              </template>
            </el-table-column>

            <!-- 修改时间列 -->
            <el-table-column prop="modTime" label="修改时间" width="180" sortable="custom">
              <template #default="{ row }">
                <div class="time-cell">
                  <span class="file-time">{{ formatTime(row.modTime) }}</span>
                  <span class="relative-time">{{ getRelativeTime(row.modTime) }}</span>
                </div>
              </template>
            </el-table-column>

            <!-- 操作列 -->
            <el-table-column label="操作" width="220" align="center" v-if="!batchMode">
              <template #default="{ row }">
                <div class="table-actions">
                  <el-tooltip content="预览" placement="top" v-if="!row.isDir">
                    <el-button size="small" @click="handlePreviewFile(row)" class="action-btn">
                      <el-icon><View /></el-icon>
                    </el-button>
                  </el-tooltip>
                  <el-tooltip content="下载" placement="top" v-if="!row.isDir">
                    <el-button size="small" type="primary" @click="downloadFile(row)" class="action-btn">
                      <el-icon><Download /></el-icon>
                    </el-button>
                  </el-tooltip>
                  <el-tooltip content="收藏" placement="top">
                    <el-button
                      size="small"
                      :type="favoriteFiles.has(row.name) ? 'warning' : ''"
                      @click="toggleFavorite(row.name)"
                      class="action-btn"
                    >
                      <el-icon>
                        <component :is="favoriteFiles.has(row.name) ? StarFilled : Star" />
                      </el-icon>
                    </el-button>
                  </el-tooltip>
                  <el-tooltip content="重命名" placement="top">
                    <el-button size="small" @click="handleShowRenameDialog(row)" class="action-btn">
                      <el-icon><Edit /></el-icon>
                    </el-button>
                  </el-tooltip>
                  <el-tooltip content="属性" placement="top">
                    <el-button size="small" @click="showFileProperties(row)" class="action-btn">
                      <el-icon><InfoFilled /></el-icon>
                    </el-button>
                  </el-tooltip>
                  <el-tooltip content="删除" placement="top">
                    <el-button size="small" type="danger" @click="deleteFile(row)" class="action-btn">
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </el-tooltip>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 空状态 -->
        <div v-if="!loading && filteredFileList.length === 0" class="empty-state">
          <el-icon size="64" class="empty-icon"><FolderOpened /></el-icon>
          <p class="empty-text">
            {{ searchKeyword ? '没有找到匹配的文件' : '此文件夹为空' }}
          </p>
          <div class="empty-actions" v-if="!searchKeyword">
            <el-button type="primary" @click="showUploadDialog = true">
              <el-icon><Upload /></el-icon>
              上传文件
            </el-button>
            <el-button @click="showCreateDirDialog = true">
              <el-icon><FolderAdd /></el-icon>
              新建文件夹
            </el-button>
          </div>
        </div>

        <!-- 上传进度显示 -->
        <div class="upload-progress-panel" v-if="uploadQueue.length > 0">
          <div class="progress-header">
            <h4>上传进度</h4>
            <el-button size="small" @click="clearCompletedUploads">清除已完成</el-button>
          </div>
          <div class="progress-list">
            <div
              v-for="item in uploadQueue"
              :key="item.file.name"
              class="progress-item"
            >
              <div class="progress-info">
                <span class="file-name">{{ item.file.name }}</span>
                <span class="file-size">{{ formatFileSize(item.file.size) }}</span>
              </div>
              <el-progress
                :percentage="item.progress"
                :status="item.status === 'error' ? 'exception' : item.status === 'success' ? 'success' : ''"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- 状态栏 -->
      <div class="status-bar">
        <div class="status-left">
          <span class="file-count">
            {{ filteredFileList.length }} 个项目
            <template v-if="selectedFiles.size > 0">
              (已选择 {{ selectedFiles.size }} 个)
            </template>
          </span>
          <span class="folder-size" v-if="currentFolderSize > 0">
            总大小: {{ formatFileSize(currentFolderSize) }}
          </span>
        </div>
        <div class="status-right">
          <span class="current-path">{{ currentPath }}</span>
          <el-button size="small" @click="showOperationHistory = true">
            <el-icon><Timer /></el-icon>
            操作历史
          </el-button>
        </div>
      </div>
    </div>

    <!-- 上传对话框 -->
    <el-dialog v-model="showUploadDialog" title="上传文件" width="500px">
      <el-upload
        ref="uploadRef"
        :auto-upload="false"
        :on-change="handleFileChange"
        :file-list="uploadFileList"
        drag
        multiple
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
      </el-upload>
      <template #footer>
        <el-button @click="showUploadDialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpload" :loading="uploading">
          上传
        </el-button>
      </template>
    </el-dialog>

    <!-- 创建文件夹对话框 -->
    <el-dialog v-model="showCreateDirDialog" title="新建文件夹" width="400px">
      <el-form :model="createDirForm" label-width="80px">
        <el-form-item label="文件夹名">
          <el-input v-model="createDirForm.name" placeholder="请输入文件夹名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDirDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreateDir">确定</el-button>
      </template>
    </el-dialog>

    <!-- 重命名对话框 -->
    <el-dialog v-model="showRenameDialog" title="重命名" width="400px">
      <el-form :model="renameForm" label-width="80px">
        <el-form-item label="新名称">
          <el-input v-model="renameForm.newName" placeholder="请输入新名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRenameDialog = false">取消</el-button>
        <el-button type="primary" @click="handleRename">确定</el-button>
      </template>
    <!-- 高级上传对话框 -->
    <el-dialog v-model="showUploadDialog" title="高级文件上传" width="600px" destroy-on-close>
      <div class="advanced-upload">
        <el-upload
          ref="uploadRef"
          :action="uploadUrl"
          :headers="uploadHeaders"
          :data="uploadData"
          :before-upload="beforeUpload"
          :on-progress="onUploadProgress"
          :on-success="onUploadSuccess"
          :on-error="onUploadError"
          :file-list="uploadFileList"
          :auto-upload="false"
          multiple
          drag
          class="upload-dragger"
        >
          <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">
            将文件拖到此处，或<em>点击上传</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
              支持多文件上传，单个文件不超过 500MB
            </div>
          </template>
        </el-upload>

        <div class="upload-options">
          <el-checkbox v-model="uploadOptions.overwrite">覆盖同名文件</el-checkbox>
          <el-checkbox v-model="uploadOptions.createPath">自动创建目录</el-checkbox>
          <el-checkbox v-model="uploadOptions.preserveTimestamp">保持时间戳</el-checkbox>
        </div>
      </div>

      <template #footer>
        <el-button @click="showUploadDialog = false">取消</el-button>
        <el-button @click="clearUploadList">清空列表</el-button>
        <el-button type="primary" @click="submitUpload" :loading="uploading">
          开始上传
        </el-button>
      </template>
    </el-dialog>

    <!-- 新建文件夹对话框 -->
    <el-dialog v-model="showCreateDirDialog" title="新建文件夹" width="400px">
      <el-form :model="createDirForm" :rules="createDirRules" ref="createDirFormRef">
        <el-form-item label="文件夹名称" prop="name">
          <el-input
            v-model="createDirForm.name"
            placeholder="请输入文件夹名称"
            @keyup.enter="createDirectory"
          />
        </el-form-item>
        <el-form-item label="权限设置" prop="permissions">
          <el-select v-model="createDirForm.permissions" placeholder="选择权限">
            <el-option label="755 (rwxr-xr-x)" value="755" />
            <el-option label="750 (rwxr-x---)" value="750" />
            <el-option label="700 (rwx------)" value="700" />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showCreateDirDialog = false">取消</el-button>
        <el-button type="primary" @click="createDirectory">创建</el-button>
      </template>
    </el-dialog>

    <!-- 重命名对话框 -->
    <el-dialog v-model="showRenameDialog" title="重命名" width="400px">
      <el-form :model="renameForm" :rules="renameRules" ref="renameFormRef">
        <el-form-item label="当前名称">
          <el-input :value="renameForm.oldName" disabled />
        </el-form-item>
        <el-form-item label="新名称" prop="newName">
          <el-input
            v-model="renameForm.newName"
            placeholder="请输入新名称"
            @keyup.enter="renameFile"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showRenameDialog = false">取消</el-button>
        <el-button type="primary" @click="renameFile">确定</el-button>
      </template>
    </el-dialog>

    <!-- 文件预览对话框 -->
    <el-dialog v-model="previewVisible" :title="`预览 - ${previewFile?.name}`" width="80%" top="5vh">
      <div class="file-preview">
        <div class="preview-toolbar">
          <el-button-group size="small">
            <el-button @click="downloadFile(previewFile!)">
              <el-icon><Download /></el-icon>
              下载
            </el-button>
            <el-button @click="editFile(previewFile!)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button @click="copyFileContent">
              <el-icon><CopyDocument /></el-icon>
              复制内容
            </el-button>
          </el-button-group>
        </div>
        <div class="preview-content">
          <pre><code>{{ previewContent }}</code></pre>
        </div>
      </div>
    </el-dialog>

    <!-- 文件属性对话框 -->
    <el-dialog v-model="showPropertiesDialog" :title="`属性 - ${propertiesFile?.name}`" width="500px">
      <div class="file-properties" v-if="propertiesFile">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="文件名">
            {{ propertiesFile.name }}
          </el-descriptions-item>
          <el-descriptions-item label="类型">
            {{ propertiesFile.isDir ? '文件夹' : '文件' }}
          </el-descriptions-item>
          <el-descriptions-item label="大小" v-if="!propertiesFile.isDir">
            {{ formatFileSize(propertiesFile.size) }}
          </el-descriptions-item>
          <el-descriptions-item label="权限">
            {{ propertiesFile.mode }}
          </el-descriptions-item>
          <el-descriptions-item label="修改时间">
            {{ formatTime(propertiesFile.modTime) }}
          </el-descriptions-item>
          <el-descriptions-item label="完整路径">
            {{ currentPath === '/' ? '/' + propertiesFile.name : currentPath + '/' + propertiesFile.name }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    <!-- 操作历史对话框 -->
    <el-dialog v-model="showOperationHistory" title="操作历史" width="700px">
      <el-table :data="operationHistory" height="400">
        <el-table-column prop="type" label="操作类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getOperationTypeColor(row.type)">
              {{ getOperationTypeName(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="file" label="文件" min-width="200" />
        <el-table-column prop="timestamp" label="时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.timestamp) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <el-button @click="clearOperationHistory">清空历史</el-button>
        <el-button type="primary" @click="showOperationHistory = false">关闭</el-button>
      </template>
    </el-dialog>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Refresh, Upload, FolderAdd, Folder, Document, UploadFilled,
  FolderOpened, Location, HomeFilled, ArrowRight, Grid, List,
  MoreFilled, Download, Edit, Delete, Picture, VideoPlay,
  DocumentCopy, Files, Cpu, Setting, Operation, Search,
  View, Sort, Filter, Share, Lock, Unlock, Star, StarFilled,
  Timer, Calendar, User, Warning, InfoFilled, SuccessFilled,
  CircleCheck, CircleClose, Plus, Minus, CopyDocument,
  ScaleToOriginal, FullScreen, Back, Right, Top, Bottom
} from '@element-plus/icons-vue'
import { sftpApi } from '@/api/cmdb/sftp'

// 定义文件信息接口
interface FileInfo {
  name: string
  size: number
  mode: string
  modTime: string
  isDir: boolean
}

const props = defineProps({
  host: {
    type: Object,
    required: true
  },
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible'])

// 响应式数据
const dialogVisible = ref(false)
const loading = ref(false)
const uploading = ref(false)
const currentPath = ref('/')
const fileList = ref<FileInfo[]>([])
const filteredFileList = ref<FileInfo[]>([])
const showUploadDialog = ref(false)
const showCreateDirDialog = ref(false)
const showRenameDialog = ref(false)
const uploadFileList = ref<any[]>([])
const uploadRef = ref()
const viewMode = ref<'grid' | 'list'>('grid')

// 高级功能状态
const batchMode = ref(false)
const selectedFiles = ref<Set<string>>(new Set())
const searchVisible = ref(false)
const searchKeyword = ref('')
const sortBy = ref<'name' | 'size' | 'modTime' | 'type'>('name')
const sortOrder = ref<'asc' | 'desc'>('asc')
const showHiddenFiles = ref(false)
const favoriteFiles = ref<Set<string>>(new Set())

// 文件操作历史
const operationHistory = ref<Array<{
  type: 'upload' | 'download' | 'delete' | 'rename' | 'create'
  file: string
  timestamp: Date
  status: 'success' | 'error'
}>>([])

// 上传进度
const uploadProgress = ref<Map<string, number>>(new Map())
const uploadQueue = ref<Array<{
  file: File
  path: string
  status: 'waiting' | 'uploading' | 'success' | 'error'
  progress: number
}>>([])

// 预览相关
const previewVisible = ref(false)
const previewFile = ref<FileInfo | null>(null)
const previewContent = ref('')

// 文件属性对话框
const showPropertiesDialog = ref(false)
const propertiesFile = ref<FileInfo | null>(null)

// 表单数据
const createDirForm = ref({
  name: '',
  permissions: '755'
})

const renameForm = ref({
  oldName: '',
  newName: ''
})

const uploadOptions = ref({
  overwrite: false,
  createPath: true,
  preserveTimestamp: false
})

// 其他状态
const showOperationHistory = ref(false)
const currentFolderSize = ref(0)

// 表单验证规则
const createDirRules = {
  name: [
    { required: true, message: '请输入文件夹名称', trigger: 'blur' },
    { pattern: /^[^\/\\:*?"<>|]+$/, message: '文件夹名称包含非法字符', trigger: 'blur' }
  ]
}

const renameRules = {
  newName: [
    { required: true, message: '请输入新名称', trigger: 'blur' },
    { pattern: /^[^\/\\:*?"<>|]+$/, message: '文件名包含非法字符', trigger: 'blur' }
  ]
}

// 表单引用
const createDirFormRef = ref()
const renameFormRef = ref()

// 计算属性
const pathSegments = computed(() => {
  return currentPath.value.split('/').filter(segment => segment !== '')
})

// 监听visible变化
watch(() => props.visible, (newVal) => {
  dialogVisible.value = newVal
  if (newVal) {
    loadFiles()
  }
})

// 监听dialogVisible变化
watch(dialogVisible, (newVal) => {
  emit('update:visible', newVal)
})

// 监听文件列表变化，自动应用过滤和排序
watch(fileList, () => {
  applyFiltersAndSort()
  // 计算当前文件夹大小
  currentFolderSize.value = fileList.value
    .filter(file => !file.isDir)
    .reduce((total, file) => total + (file.size || 0), 0)
}, { deep: true })

// 监听搜索关键词变化
watch(searchKeyword, () => {
  applyFiltersAndSort()
})

// 监听隐藏文件显示设置变化
watch(showHiddenFiles, () => {
  applyFiltersAndSort()
})

// 关闭弹窗
const handleClose = () => {
  dialogVisible.value = false
}

// 组件挂载
onMounted(() => {
  if (props.visible) {
    dialogVisible.value = true
    loadFiles()
  }
})

// 加载文件列表
const loadFiles = async () => {
  if (!props.host?.id) return

  loading.value = true
  try {
    const response = await sftpApi.listFiles(props.host.id.toString(), currentPath.value)
    
    if (response.code === 0 || response.code === 200) {
      fileList.value = response.data || []
    } else {
      ElMessage.error(response.message || '加载文件列表失败')
      fileList.value = []
    }
  } catch (error: any) {
    console.error('Load files error:', error)
    ElMessage.error('加载文件列表失败：' + (error.message || '网络错误'))
    fileList.value = []
  } finally {
    loading.value = false
  }
}

// 刷新文件列表
const refreshFiles = () => {
  loadFiles()
}

// 导航到指定路径
const navigateToPath = (index: number) => {
  if (index === -1) {
    // 点击root，回到根目录
    currentPath.value = '/'
  } else {
    const segments = pathSegments.value.slice(0, index + 1)
    currentPath.value = '/' + segments.join('/')
  }
  loadFiles()
}

// 处理行双击事件
const handleRowDoubleClick = (row: FileInfo) => {
  if (row.isDir) {
    // 进入目录
    currentPath.value = currentPath.value.endsWith('/') 
      ? currentPath.value + row.name
      : currentPath.value + '/' + row.name
    loadFiles()
  } else {
    // 下载文件
    downloadFile(row)
  }
}

// 下载文件
const downloadFile = async (file: FileInfo) => {
  if (!props.host?.id) return

  try {
    const filePath = currentPath.value.endsWith('/') 
      ? currentPath.value + file.name
      : currentPath.value + '/' + file.name
      
    const blob = await sftpApi.downloadFile(props.host.id.toString(), filePath)
    
    // 创建下载链接
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = file.name
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('文件下载成功')
  } catch (error: any) {
    console.error('Download file error:', error)
    ElMessage.error('文件下载失败：' + (error.message || '网络错误'))
  }
}

// 文件选择处理
const handleFileChange = (file: any) => {
  uploadFileList.value.push(file)
}

// 处理文件上传
const handleUpload = async () => {
  if (!props.host?.id || uploadFileList.value.length === 0) {
    ElMessage.warning('请选择要上传的文件')
    return
  }

  uploading.value = true
  try {
    for (const fileWrapper of uploadFileList.value) {
      const file = fileWrapper.raw || fileWrapper
      await sftpApi.uploadFile(props.host.id.toString(), currentPath.value, file)
    }
    
    ElMessage.success('文件上传成功')
    showUploadDialog.value = false
    uploadFileList.value = []
    uploadRef.value?.clearFiles()
    loadFiles() // 刷新文件列表
  } catch (error: any) {
    console.error('Upload file error:', error)
    ElMessage.error('文件上传失败：' + (error.message || '网络错误'))
  } finally {
    uploading.value = false
  }
}

// 创建目录
const handleCreateDir = async () => {
  if (!props.host?.id || !createDirForm.value.name.trim()) {
    ElMessage.warning('请输入文件夹名称')
    return
  }

  try {
    const dirPath = currentPath.value.endsWith('/') 
      ? currentPath.value + createDirForm.value.name
      : currentPath.value + '/' + createDirForm.value.name

    const response = await sftpApi.createDirectory(props.host.id.toString(), dirPath)
    
    if (response.code === 0 || response.code === 200) {
      ElMessage.success('文件夹创建成功')
      showCreateDirDialog.value = false
      createDirForm.value.name = ''
      loadFiles() // 刷新文件列表
    } else {
      ElMessage.error(response.message || '文件夹创建失败')
    }
  } catch (error: any) {
    console.error('Create directory error:', error)
    ElMessage.error('文件夹创建失败：' + (error.message || '网络错误'))
  }
}

// 显示重命名对话框
const handleShowRenameDialog = (file: FileInfo) => {
  renameForm.value.oldName = file.name
  renameForm.value.newName = file.name
  showRenameDialog.value = true
}

// 重命名文件
const handleRename = async () => {
  if (!props.host?.id || !renameForm.value.newName.trim()) {
    ElMessage.warning('请输入新名称')
    return
  }

  if (renameForm.value.oldName === renameForm.value.newName) {
    ElMessage.warning('新名称与原名称相同')
    return
  }

  try {
    const oldPath = currentPath.value.endsWith('/') 
      ? currentPath.value + renameForm.value.oldName
      : currentPath.value + '/' + renameForm.value.oldName
      
    const newPath = currentPath.value.endsWith('/') 
      ? currentPath.value + renameForm.value.newName
      : currentPath.value + '/' + renameForm.value.newName

    const response = await sftpApi.renameFile(props.host.id.toString(), oldPath, newPath)
    
    if (response.code === 0 || response.code === 200) {
      ElMessage.success('重命名成功')
      showRenameDialog.value = false
      renameForm.value.oldName = ''
      renameForm.value.newName = ''
      loadFiles() // 刷新文件列表
    } else {
      ElMessage.error(response.message || '重命名失败')
    }
  } catch (error: any) {
    console.error('Rename file error:', error)
    ElMessage.error('重命名失败：' + (error.message || '网络错误'))
  }
}

// 删除文件
const deleteFile = async (file: FileInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 "${file.name}" 吗？${file.isDir ? '删除文件夹将会删除其中的所有内容。' : ''}`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const filePath = currentPath.value.endsWith('/') 
      ? currentPath.value + file.name
      : currentPath.value + '/' + file.name

    const response = await sftpApi.deleteFile(props.host.id.toString(), filePath)
    
    if (response.code === 0 || response.code === 200) {
      ElMessage.success('删除成功')
      loadFiles() // 刷新文件列表
    } else {
      ElMessage.error(response.message || '删除失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Delete file error:', error)
      ElMessage.error('删除失败：' + (error.message || '网络错误'))
    }
  }
}

// 格式化文件大小
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化时间
const formatTime = (timeStr: string): string => {
  try {
    const date = new Date(timeStr)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    return timeStr
  }
}

// 返回上级目录
const goToParentDirectory = () => {
  if (currentPath.value === '/') return

  const segments = currentPath.value.split('/').filter(s => s !== '')
  segments.pop()
  currentPath.value = segments.length > 0 ? '/' + segments.join('/') : '/'
  loadFiles()
}

// 获取文件图标
const getFileIcon = (fileName: string) => {
  const ext = fileName.split('.').pop()?.toLowerCase()

  switch (ext) {
    case 'jpg':
    case 'jpeg':
    case 'png':
    case 'gif':
    case 'bmp':
    case 'svg':
      return Picture
    case 'mp4':
    case 'avi':
    case 'mov':
    case 'wmv':
    case 'flv':
      return VideoPlay
    case 'txt':
    case 'md':
    case 'log':
      return DocumentCopy
    case 'zip':
    case 'rar':
    case '7z':
    case 'tar':
    case 'gz':
      return Files
    case 'exe':
    case 'msi':
    case 'deb':
    case 'rpm':
      return Cpu
    case 'conf':
    case 'config':
    case 'ini':
    case 'yaml':
    case 'yml':
    case 'json':
      return Setting
    default:
      return Document
  }
}

// 获取文件图标样式类
const getFileIconClass = (fileName: string) => {
  const ext = fileName.split('.').pop()?.toLowerCase()

  switch (ext) {
    case 'jpg':
    case 'jpeg':
    case 'png':
    case 'gif':
    case 'bmp':
    case 'svg':
      return 'image-icon'
    case 'mp4':
    case 'avi':
    case 'mov':
    case 'wmv':
    case 'flv':
      return 'video-icon'
    case 'txt':
    case 'md':
    case 'log':
      return 'text-icon'
    case 'zip':
    case 'rar':
    case '7z':
    case 'tar':
    case 'gz':
      return 'archive-icon'
    case 'exe':
    case 'msi':
    case 'deb':
    case 'rpm':
      return 'executable-icon'
    case 'conf':
    case 'config':
    case 'ini':
    case 'yaml':
    case 'yml':
    case 'json':
      return 'config-icon'
    default:
      return 'document-icon'
  }
}

// 处理文件操作
const handleFileAction = (command: any) => {
  const { action, file } = command

  switch (action) {
    case 'download':
      downloadFile(file)
      break
    case 'preview':
      previewFile(file)
      break
    case 'rename':
      handleShowRenameDialog(file)
      break
    case 'copy':
      copyFile(file)
      break
    case 'properties':
      showFileProperties(file)
      break
    case 'delete':
      deleteFile(file)
      break
  }
}

const copyFile = (file: FileInfo) => {
  // 实现文件复制功能
  ElMessage.info('文件复制功能开发中...')
}

// 显示右键菜单
const showContextMenu = (event: MouseEvent, file: FileInfo) => {
  // 这里可以实现右键菜单功能
  console.log('Right click on file:', file.name)
}

// 高级功能方法
const getHostIP = () => {
  return props.host?.publicIP || props.host?.privateIP || '未知IP'
}

const toggleBatchMode = () => {
  batchMode.value = !batchMode.value
  if (!batchMode.value) {
    selectedFiles.value.clear()
  }
}

const toggleSearch = () => {
  searchVisible.value = !searchVisible.value
  if (!searchVisible.value) {
    searchKeyword.value = ''
    filteredFileList.value = fileList.value
  }
}

const handleFileSelect = (fileName: string, selected: boolean) => {
  if (selected) {
    selectedFiles.value.add(fileName)
  } else {
    selectedFiles.value.delete(fileName)
  }
}

const selectAllFiles = () => {
  if (selectedFiles.value.size === filteredFileList.value.length) {
    selectedFiles.value.clear()
  } else {
    filteredFileList.value.forEach(file => {
      selectedFiles.value.add(file.name)
    })
  }
}

const batchDelete = async () => {
  if (selectedFiles.value.size === 0) {
    ElMessage.warning('请先选择要删除的文件')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedFiles.value.size} 个文件吗？`,
      '批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const deletePromises = Array.from(selectedFiles.value).map(fileName => {
      return sftpApi.deleteFile({
        host_id: props.host.id,
        file_path: currentPath.value === '/' ? fileName : `${currentPath.value}/${fileName}`
      })
    })

    await Promise.all(deletePromises)
    ElMessage.success('批量删除成功')
    selectedFiles.value.clear()
    await loadFiles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

const batchDownload = async () => {
  if (selectedFiles.value.size === 0) {
    ElMessage.warning('请先选择要下载的文件')
    return
  }

  for (const fileName of selectedFiles.value) {
    const file = fileList.value.find(f => f.name === fileName)
    if (file && !file.isDir) {
      await downloadFile(file)
    }
  }
}

const sortFiles = (by: typeof sortBy.value) => {
  if (sortBy.value === by) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortBy.value = by
    sortOrder.value = 'asc'
  }
  applyFiltersAndSort()
}

const applyFiltersAndSort = () => {
  let filtered = [...fileList.value]

  // 应用搜索过滤
  if (searchKeyword.value) {
    filtered = filtered.filter(file =>
      file.name.toLowerCase().includes(searchKeyword.value.toLowerCase())
    )
  }

  // 应用隐藏文件过滤
  if (!showHiddenFiles.value) {
    filtered = filtered.filter(file => !file.name.startsWith('.'))
  }

  // 应用排序
  filtered.sort((a, b) => {
    let comparison = 0

    // 文件夹优先
    if (a.isDir && !b.isDir) return -1
    if (!a.isDir && b.isDir) return 1

    switch (sortBy.value) {
      case 'name':
        comparison = a.name.localeCompare(b.name)
        break
      case 'size':
        comparison = (a.size || 0) - (b.size || 0)
        break
      case 'modTime':
        comparison = new Date(a.modTime).getTime() - new Date(b.modTime).getTime()
        break
      case 'type':
        const aExt = a.name.split('.').pop() || ''
        const bExt = b.name.split('.').pop() || ''
        comparison = aExt.localeCompare(bExt)
        break
    }

    return sortOrder.value === 'asc' ? comparison : -comparison
  })

  filteredFileList.value = filtered
}

const toggleFavorite = (fileName: string) => {
  if (favoriteFiles.value.has(fileName)) {
    favoriteFiles.value.delete(fileName)
  } else {
    favoriteFiles.value.add(fileName)
  }
  // 这里可以保存到本地存储或服务器
}

const handlePreviewFile = async (file: FileInfo) => {
  if (file.isDir) return

  const ext = file.name.split('.').pop()?.toLowerCase()
  const textExtensions = ['txt', 'md', 'json', 'xml', 'html', 'css', 'js', 'ts', 'vue', 'py', 'java', 'cpp', 'c', 'h']

  if (textExtensions.includes(ext || '')) {
    try {
      // 这里需要实现文件内容获取API
      previewContent.value = '文件预览功能开发中...'
      previewFile.value = file
      previewVisible.value = true
    } catch (error) {
      ElMessage.error('预览失败')
    }
  } else {
    ElMessage.info('该文件类型不支持预览')
  }
}

const showFileProperties = (file: FileInfo) => {
  propertiesFile.value = file
  showPropertiesDialog.value = true
}

const addToOperationHistory = (operation: any) => {
  operationHistory.value.unshift({
    ...operation,
    timestamp: new Date()
  })
  // 保持历史记录在100条以内
  if (operationHistory.value.length > 100) {
    operationHistory.value = operationHistory.value.slice(0, 100)
  }
}

// 更多辅助方法
const getFileExtension = (fileName: string) => {
  const ext = fileName.split('.').pop()
  return ext ? ext.toUpperCase() : ''
}

const getPermissionDescription = (mode: string) => {
  // 简化的权限描述
  const descriptions: Record<string, string> = {
    '755': '所有者：读写执行，组：读执行，其他：读执行',
    '644': '所有者：读写，组：读，其他：读',
    '600': '所有者：读写，组：无，其他：无',
    '777': '所有者：读写执行，组：读写执行，其他：读写执行'
  }
  return descriptions[mode] || '自定义权限'
}

const getPermissionType = (mode: string) => {
  if (mode.includes('7')) return 'danger'
  if (mode.includes('6')) return 'warning'
  return 'info'
}

const getRelativeTime = (timestamp: string) => {
  const now = new Date()
  const time = new Date(timestamp)
  const diff = now.getTime() - time.getTime()

  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 30) return `${days}天前`
  return '很久以前'
}

const getOperationTypeName = (type: string) => {
  const names: Record<string, string> = {
    'upload': '上传',
    'download': '下载',
    'delete': '删除',
    'rename': '重命名',
    'create': '创建'
  }
  return names[type] || type
}

const getOperationTypeColor = (type: string) => {
  const colors: Record<string, string> = {
    'upload': 'success',
    'download': 'primary',
    'delete': 'danger',
    'rename': 'warning',
    'create': 'info'
  }
  return colors[type] || ''
}

const clearOperationHistory = () => {
  operationHistory.value = []
  ElMessage.success('操作历史已清空')
}

const clearCompletedUploads = () => {
  uploadQueue.value = uploadQueue.value.filter(item =>
    item.status === 'waiting' || item.status === 'uploading'
  )
}

const toggleFullscreen = () => {
  // 实现全屏切换
  ElMessage.info('全屏功能开发中...')
}

const copyFileContent = () => {
  if (previewContent.value) {
    navigator.clipboard.writeText(previewContent.value)
    ElMessage.success('内容已复制到剪贴板')
  }
}

const editFile = (file: FileInfo) => {
  ElMessage.info('文件编辑功能开发中...')
}

// 上传相关方法
const uploadUrl = computed(() => {
  return `/api/v1/cmdb/sftp/upload?host_id=${props.host.id}&path=${encodeURIComponent(currentPath.value)}`
})

const uploadHeaders = computed(() => {
  return {
    'Authorization': `Bearer ${userStore.token}`
  }
})

const uploadData = computed(() => {
  return {
    overwrite: uploadOptions.value.overwrite,
    createPath: uploadOptions.value.createPath,
    preserveTimestamp: uploadOptions.value.preserveTimestamp
  }
})

const beforeUpload = (file: File) => {
  const isLt500M = file.size / 1024 / 1024 < 500
  if (!isLt500M) {
    ElMessage.error('文件大小不能超过 500MB!')
    return false
  }
  return true
}

const onUploadProgress = (event: any, file: any) => {
  const progress = Math.round((event.loaded * 100) / event.total)
  uploadProgress.value.set(file.name, progress)
}

const onUploadSuccess = (response: any, file: any) => {
  ElMessage.success(`${file.name} 上传成功`)
  addToOperationHistory({
    type: 'upload',
    file: file.name,
    status: 'success'
  })
  loadFiles()
}

const onUploadError = (error: any, file: any) => {
  ElMessage.error(`${file.name} 上传失败`)
  addToOperationHistory({
    type: 'upload',
    file: file.name,
    status: 'error'
  })
}

const submitUpload = () => {
  uploadRef.value?.submit()
}

const clearUploadList = () => {
  uploadRef.value?.clearFiles()
  uploadFileList.value = []
}

const createDirectory = async () => {
  try {
    await createDirFormRef.value?.validate()

    const dirPath = currentPath.value === '/'
      ? createDirForm.value.name
      : `${currentPath.value}/${createDirForm.value.name}`

    await sftpApi.createDirectory({
      host_id: props.host.id,
      dir_path: dirPath,
      permissions: createDirForm.value.permissions
    })

    ElMessage.success('文件夹创建成功')
    showCreateDirDialog.value = false
    createDirForm.value.name = ''
    addToOperationHistory({
      type: 'create',
      file: createDirForm.value.name,
      status: 'success'
    })
    await loadFiles()
  } catch (error) {
    ElMessage.error('文件夹创建失败')
  }
}

const renameFile = async () => {
  try {
    await renameFormRef.value?.validate()

    const oldPath = currentPath.value === '/'
      ? renameForm.value.oldName
      : `${currentPath.value}/${renameForm.value.oldName}`

    const newPath = currentPath.value === '/'
      ? renameForm.value.newName
      : `${currentPath.value}/${renameForm.value.newName}`

    await sftpApi.renameFile({
      host_id: props.host.id,
      old_path: oldPath,
      new_path: newPath
    })

    ElMessage.success('重命名成功')
    showRenameDialog.value = false
    addToOperationHistory({
      type: 'rename',
      file: `${renameForm.value.oldName} -> ${renameForm.value.newName}`,
      status: 'success'
    })
    await loadFiles()
  } catch (error) {
    ElMessage.error('重命名失败')
  }
}
</script>

<style scoped>
/* 现代化对话框样式 */
:deep(.modern-sftp-dialog) {
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
}

:deep(.modern-sftp-dialog .el-dialog__header) {
  padding: 0;
  margin: 0;
  border-bottom: none;
}

:deep(.modern-sftp-dialog .el-dialog__body) {
  padding: 0;
  height: 80vh;
  overflow: hidden;
  background: #f8fafc;
}

/* 对话框头部 */
.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  position: relative;
  overflow: hidden;
}

.dialog-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grid" width="10" height="10" patternUnits="userSpaceOnUse"><path d="M 10 0 L 0 0 0 10" fill="none" stroke="rgba(255,255,255,0.1)" stroke-width="0.5"/></pattern></defs><rect width="100" height="100" fill="url(%23grid)"/></svg>');
  opacity: 0.3;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
  position: relative;
  z-index: 1;
}

.header-icon {
  width: 56px;
  height: 56px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.header-info h3 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.host-info {
  margin: 6px 0 0 0;
  font-size: 14px;
  opacity: 0.9;
  font-weight: 500;
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  padding: 4px 12px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  font-size: 12px;
  backdrop-filter: blur(10px);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #10b981;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.5);
}

.status-dot.connected {
  background: #10b981;
  animation: pulse-green 2s infinite;
}

@keyframes pulse-green {
  0%, 100% {
    box-shadow: 0 0 8px rgba(16, 185, 129, 0.5);
  }
  50% {
    box-shadow: 0 0 16px rgba(16, 185, 129, 0.8);
  }
}

.header-actions {
  display: flex;
  gap: 12px;
  position: relative;
  z-index: 1;
}

.action-btn {
  background: rgba(255, 255, 255, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: white;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.25);
  border-color: rgba(255, 255, 255, 0.4);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

/* 主容器 */
.modern-sftp-manager {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #f8fafc;
}

/* 高级工具栏 */
.advanced-toolbar {
  display: flex;
  flex-direction: column;
  background: white;
  border-bottom: 1px solid #e5e7eb;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.advanced-toolbar > div {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
}

.advanced-toolbar > div:not(:last-child) {
  border-bottom: 1px solid #f3f4f6;
}

.search-container {
  flex: 1;
  max-width: 300px;
  margin: 0 20px;
}

.toolbar-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.toolbar-actions .el-button-group {
  margin-right: 8px;
}

/* 批量操作栏 */
.batch-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 24px;
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  border-bottom: 1px solid #f59e0b;
}

.batch-info {
  display: flex;
  align-items: center;
  gap: 12px;
  font-weight: 500;
  color: #92400e;
}

.batch-actions {
  display: flex;
  gap: 8px;
}

.path-navigation {
  display: flex;
  align-items: center;
  gap: 8px;
}

.path-icon {
  color: #6b7280;
}

.breadcrumb-container {
  display: flex;
  align-items: center;
  gap: 4px;
}

.path-segment {
  padding: 6px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
  color: #6b7280;
  display: flex;
  align-items: center;
  gap: 4px;
}

.path-segment:hover {
  background: #f3f4f6;
  color: #374151;
}

.path-segment.root {
  color: #3b82f6;
  font-weight: 500;
}

.path-segment.active {
  background: #dbeafe;
  color: #1d4ed8;
  font-weight: 500;
}

.separator {
  color: #d1d5db;
  font-size: 12px;
}

.view-controls {
  display: flex;
  gap: 8px;
}

/* 文件内容区域 */
.file-content {
  flex: 1;
  overflow: auto;
  padding: 24px;
}

/* 网格视图 */
.grid-view {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 24px;
  padding: 24px;
}

.file-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  position: relative;
  border: 2px solid transparent;
  overflow: hidden;
}

.file-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #667eea, #764ba2);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.file-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
  border-color: #e5e7eb;
}

.file-card:hover::before {
  opacity: 1;
}

.file-card.selected {
  border-color: #3b82f6;
  background: #eff6ff;
}

.file-card.selected::before {
  opacity: 1;
  background: #3b82f6;
}

.file-card.favorite {
  position: relative;
}

.file-checkbox {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 2;
}

.favorite-mark {
  position: absolute;
  top: 12px;
  right: 12px;
  z-index: 2;
}

.favorite-icon {
  color: #f59e0b;
  filter: drop-shadow(0 2px 4px rgba(245, 158, 11, 0.3));
}

.file-icon {
  text-align: center;
  margin-bottom: 12px;
}

.folder-icon {
  color: #3b82f6;
}

.image-icon {
  color: #10b981;
}

.video-icon {
  color: #f59e0b;
}

.text-icon {
  color: #6b7280;
}

.archive-icon {
  color: #8b5cf6;
}

.executable-icon {
  color: #ef4444;
}

.config-icon {
  color: #06b6d4;
}

.document-icon {
  color: #6b7280;
}

.file-info {
  text-align: center;
}

.file-name {
  font-weight: 500;
  color: #374151;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-meta {
  font-size: 12px;
  color: #6b7280;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.file-actions {
  position: absolute;
  top: 12px;
  right: 12px;
  opacity: 0;
  transition: opacity 0.2s;
}

.file-card:hover .file-actions {
  opacity: 1;
}

.more-btn {
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid #e5e7eb;
}

/* 列表视图 */
.list-view {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

:deep(.modern-table) {
  border: none;
}

:deep(.modern-table .el-table__header) {
  background: #f8fafc;
}

:deep(.modern-table .el-table__row:hover > .el-table__cell) {
  background: #f8fafc;
}

.file-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.name {
  font-weight: 500;
  color: #374151;
}

.type-badge {
  background: #dbeafe;
  color: #1d4ed8;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
}

.file-size {
  color: #6b7280;
  font-size: 13px;
}

.file-time {
  color: #6b7280;
  font-size: 13px;
}

.table-actions {
  display: flex;
  gap: 4px;
  justify-content: center;
}

.table-actions .action-btn {
  background: transparent;
  border: 1px solid #e5e7eb;
  color: #6b7280;
  width: 32px;
  height: 32px;
  padding: 0;
}

.table-actions .action-btn:hover {
  background: #f3f4f6;
  border-color: #d1d5db;
  color: #374151;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #6b7280;
}

.empty-icon {
  color: #d1d5db;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  margin-bottom: 24px;
}

/* 下拉菜单样式 */
:deep(.el-dropdown-menu__item.danger-item) {
  color: #ef4444;
}

:deep(.el-dropdown-menu__item.danger-item:hover) {
  background: #fef2f2;
  color: #dc2626;
}

/* 上传对话框样式 */
.el-upload__text {
  margin-top: 12px;
  color: #6b7280;
}

.el-icon--upload {
  font-size: 67px;
  color: #d1d5db;
  margin-bottom: 16px;
}

/* 状态栏 */
.status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 24px;
  background: white;
  border-top: 1px solid #e5e7eb;
  font-size: 13px;
  color: #6b7280;
}

.status-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.file-count {
  font-weight: 500;
}

.folder-size {
  color: #9ca3af;
}

.status-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.current-path {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  background: #f3f4f6;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
}

/* 上传进度面板 */
.upload-progress-panel {
  position: fixed;
  bottom: 20px;
  right: 20px;
  width: 350px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  border: 1px solid #e5e7eb;
  z-index: 1000;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f3f4f6;
}

.progress-header h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #374151;
}

.progress-list {
  max-height: 300px;
  overflow-y: auto;
  padding: 16px 20px;
}

.progress-item {
  margin-bottom: 16px;
}

.progress-item:last-child {
  margin-bottom: 0;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.progress-info .file-name {
  font-weight: 500;
  color: #374151;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  margin-right: 12px;
}

.progress-info .file-size {
  font-size: 12px;
  color: #6b7280;
}

/* 高级上传对话框 */
.advanced-upload {
  padding: 20px 0;
}

.upload-dragger {
  margin-bottom: 20px;
}

.upload-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  margin-top: 16px;
}

/* 文件预览 */
.file-preview {
  height: 60vh;
  display: flex;
  flex-direction: column;
}

.preview-toolbar {
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
  background: #f8fafc;
}

.preview-content {
  flex: 1;
  overflow: auto;
  padding: 20px;
  background: #1a1a1a;
  color: #ffffff;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
}

.preview-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .grid-view {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 16px;
    padding: 16px;
  }

  .file-card {
    padding: 16px;
  }

  .advanced-toolbar > div {
    padding: 12px 16px;
  }

  .file-content {
    padding: 16px;
  }

  .status-bar {
    padding: 8px 16px;
    flex-direction: column;
    gap: 8px;
    align-items: flex-start;
  }

  .upload-progress-panel {
    width: calc(100vw - 40px);
    right: 20px;
    left: 20px;
  }

  .dialog-header {
    padding: 20px 16px;
  }

  .header-left {
    gap: 12px;
  }

  .header-icon {
    width: 48px;
    height: 48px;
  }

  .header-info h3 {
    font-size: 20px;
  }
}
</style>
