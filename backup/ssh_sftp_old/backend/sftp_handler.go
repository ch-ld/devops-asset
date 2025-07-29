package cmdb

import (
	"api-server/internal/response/response"
	"api-server/internal/service/cmdb"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SftpHandler SFTP文件管理处理器
type SftpHandler struct {
	sftpService *cmdb.SftpService
}

// NewSftpHandler 创建SFTP处理器实例
func NewSftpHandler(sftpService *cmdb.SftpService) *SftpHandler {
	return &SftpHandler{sftpService: sftpService}
}

// FileInfo 文件信息响应结构
type FileInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	IsDir       bool   `json:"is_dir"`
	Mode        string `json:"mode"`
	ModTime     string `json:"mod_time"`
	Permissions string `json:"permissions"`
}

// List 列出指定路径下的文件和目录
// @Summary 列出主机目录内容
// @Description 列出指定主机指定路径下的文件和目录
// @Tags CMDB-文件管理
// @Accept json
// @Produce json
// @Param host_id query int true "主机ID"
// @Param path query string false "路径，默认为/home"
// @Success 200 {array} FileInfo
// @Failure 400 {object} response.Response "无效的参数"
// @Failure 404 {object} response.Response "主机不存在"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/cmdb/sftp/list [get]
func (h *SftpHandler) List(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Query("host_id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的主机ID")
		return
	}

	path := c.DefaultQuery("path", "/home")

	// 调用服务获取文件列表
	fileInfos, err := h.sftpService.ListFiles(uint(hostID), path)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("列出文件失败: %s", err.Error()))
		return
	}

	// 构建响应
	var result []FileInfo
	for _, file := range fileInfos {
		result = append(result, FileInfo{
			Name:        file.Name(),
			Path:        filepath.Join(path, file.Name()),
			Size:        file.Size(),
			IsDir:       file.IsDir(),
			Mode:        file.Mode().String(),
			ModTime:     file.ModTime().Format("2006-01-02 15:04:05"),
			Permissions: file.Mode().Perm().String(),
		})
	}

	response.ReturnData(c, result)
}

// Upload 上传文件到主机
// @Summary 上传文件到主机
// @Description 上传文件到指定主机的指定路径
// @Tags CMDB-文件管理
// @Accept multipart/form-data
// @Produce json
// @Param host_id formData int true "主机ID"
// @Param path formData string false "目标路径，默认为/home"
// @Param file formData file true "要上传的文件"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response "无效的参数"
// @Failure 404 {object} response.Response "主机不存在"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/cmdb/sftp/upload [post]
func (h *SftpHandler) Upload(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.PostForm("host_id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的主机ID")
		return
	}

	// 获取目标路径
	remotePath := c.DefaultPostForm("path", "/home")

	// 获取上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "上传文件无效")
		return
	}
	defer file.Close()

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "upload-*")
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "创建临时文件失败")
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// 将上传内容写入临时文件
	if _, err := io.Copy(tempFile, file); err != nil {
		response.ReturnError(c, response.INTERNAL, "保存上传文件失败")
		return
	}
	tempFile.Close() // 关闭以确保内容被写入

	// 构建目标路径
	targetPath := filepath.Join(remotePath, header.Filename)

	// 上传文件到远程主机
	if err := h.sftpService.UploadFile(uint(hostID), tempFile.Name(), targetPath); err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("上传文件失败: %s", err.Error()))
		return
	}

	response.ReturnSuccess(c)
}

// Download 从主机下载文件
// @Summary 从主机下载文件
// @Description 从指定主机下载指定路径的文件
// @Tags CMDB-文件管理
// @Produce octet-stream
// @Param host_id query int true "主机ID"
// @Param path query string true "文件路径"
// @Success 200 {file} binary "文件内容"
// @Failure 400 {object} response.Response "无效的参数"
// @Failure 404 {object} response.Response "主机不存在或文件不存在"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/cmdb/sftp/download [get]
func (h *SftpHandler) Download(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Query("host_id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的主机ID")
		return
	}

	// 获取文件路径
	remotePath := c.Query("path")
	if remotePath == "" {
		response.ReturnError(c, response.INVALID_ARGUMENT, "文件路径不能为空")
		return
	}

	// 创建临时文件接收下载内容
	tempFile, err := os.CreateTemp("", "download-*")
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "创建临时文件失败")
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// 下载文件到临时位置
	if err := h.sftpService.DownloadFile(uint(hostID), remotePath, tempFile.Name()); err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("下载文件失败: %s", err.Error()))
		return
	}

	// 获取文件信息用于设置响应头
	fileInfo, err := os.Stat(tempFile.Name())
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "获取文件信息失败")
		return
	}

	// 设置响应头
	filename := filepath.Base(remotePath)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// 发送文件
	c.File(tempFile.Name())
}

// Delete 删除主机上的文件
// @Summary 删除主机上的文件
// @Description 删除指定主机上的指定文件或空目录
// @Tags CMDB-文件管理
// @Accept json
// @Produce json
// @Param host_id query int true "主机ID"
// @Param path query string true "文件路径"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response "无效的参数"
// @Failure 404 {object} response.Response "主机不存在或文件不存在"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/cmdb/sftp/delete [delete]
func (h *SftpHandler) Delete(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Query("host_id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的主机ID")
		return
	}

	// 获取文件路径
	remotePath := c.Query("path")
	if remotePath == "" {
		response.ReturnError(c, response.INVALID_ARGUMENT, "文件路径不能为空")
		return
	}

	// 删除远程文件
	if err := h.sftpService.DeleteFile(uint(hostID), remotePath); err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("删除文件失败: %s", err.Error()))
		return
	}

	response.ReturnSuccess(c)
}

// MakeDir 在主机上创建目录
// @Summary 在主机上创建目录
// @Description 在指定主机上创建指定路径的目录
// @Tags CMDB-文件管理
// @Accept json
// @Produce json
// @Param host_id formData int true "主机ID"
// @Param path formData string true "目录路径"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response "无效的参数"
// @Failure 404 {object} response.Response "主机不存在"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/cmdb/sftp/mkdir [post]
func (h *SftpHandler) MakeDir(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.PostForm("host_id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的主机ID")
		return
	}

	// 获取目录路径
	remotePath := c.PostForm("path")
	if remotePath == "" {
		response.ReturnError(c, response.INVALID_ARGUMENT, "目录路径不能为空")
		return
	}

	// 在远程主机上创建目录
	if err := h.sftpService.MakeDir(uint(hostID), remotePath); err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("创建目录失败: %s", err.Error()))
		return
	}

	response.ReturnSuccess(c)
}

// Rename 重命名主机上的文件或目录
// @Summary 重命名主机上的文件或目录
// @Description 重命名指定主机上的文件或目录
// @Tags CMDB-文件管理
// @Accept json
// @Produce json
// @Param host_id formData int true "主机ID"
// @Param old_path formData string true "原路径"
// @Param new_path formData string true "新路径"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response "无效的参数"
// @Failure 404 {object} response.Response "主机不存在或文件不存在"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/cmdb/sftp/rename [post]
func (h *SftpHandler) Rename(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.PostForm("host_id"), 10, 32)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "无效的主机ID")
		return
	}

	// 获取路径信息
	oldPath := c.PostForm("old_path")
	newPath := c.PostForm("new_path")
	if oldPath == "" || newPath == "" {
		response.ReturnError(c, response.INVALID_ARGUMENT, "路径不能为空")
		return
	}

	// 重命名文件或目录
	if err := h.sftpService.RenameFile(uint(hostID), oldPath, newPath); err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("重命名失败: %s", err.Error()))
		return
	}

	response.ReturnSuccess(c)
}

// RegisterRoutes 注册SFTP相关路由
func (h *SftpHandler) RegisterRoutes(r *gin.RouterGroup) {
	sftpGroup := r.Group("/sftp")
	{
		sftpGroup.GET("/list", h.List)
		sftpGroup.POST("/upload", h.Upload)
		sftpGroup.GET("/download", h.Download)
		sftpGroup.DELETE("/delete", h.Delete)
		sftpGroup.POST("/mkdir", h.MakeDir)
		sftpGroup.POST("/rename", h.Rename)
	}
}
