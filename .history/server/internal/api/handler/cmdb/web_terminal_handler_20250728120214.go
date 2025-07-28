package cmdb

import (
	"api-server/internal/config"
	"api-server/internal/middleware/middleware"
	repo "api-server/internal/repository/cmdb"
	"api-server/pkg/crypto/encryption"
	"api-server/pkg/ssh"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	gosshtool "golang.org/x/crypto/ssh"
)

// WebTerminalHandler WebSSH终端处理器
// 负责主机WebSSH终端连接与操作的HTTP/WebSocket请求处理
type WebTerminalHandler struct {
	hostRepo *repo.HostRepository
}

// TerminalPreference 终端偏好设置
type TerminalPreference struct {
	FontSize    int    `json:"font_size"`    // 字体大小，默认14
	Theme       string `json:"theme"`        // 主题：dark, light, solarized
	Fullscreen  bool   `json:"fullscreen"`   // 是否全屏
	CursorStyle string `json:"cursor_style"` // 光标样式：block, bar, underline
}

// TerminalSettings 终端连接配置
type TerminalSettings struct {
	HostID       uint               `json:"host_id"`     // 主机ID
	Width        int                `json:"width"`       // 终端宽度
	Height       int                `json:"height"`      // 终端高度
	Preferences  TerminalPreference `json:"preferences"` // 偏好设置
	AuthUsername string             `json:"username"`    // 覆盖主机默认用户名
	AuthPassword string             `json:"password"`    // 覆盖主机默认密码
	PrivateKey   string             `json:"private_key"` // 使用私钥认证
}

// NewWebTerminalHandler 创建WebTerminal处理器实例
func NewWebTerminalHandler(hostRepo *repo.HostRepository) *WebTerminalHandler {
	return &WebTerminalHandler{hostRepo: hostRepo}
}

// HandleSSHConnection 处理WebSocket SSH连接
// @Summary 建立WebSocket终端连接
// @Description 与远程主机建立SSH终端连接
// @Tags CMDB-终端管理
// @Accept json
// @Produce json
// @Param host_id query int true "主机ID"
// @Param width query int false "终端宽度(列数)"
// @Param height query int false "终端高度(行数)"
// @Success 101 {string} string "Switching Protocols to websocket"
// @Failure 400 {string} string "无效的参数"
// @Failure 404 {string} string "主机不存在"
// @Failure 500 {string} string "内部服务器错误"
// @Router /api/v1/ws/terminal [get]
func (h *WebTerminalHandler) HandleSSHConnection(c *gin.Context) {
	// 升级到WebSocket连接
	wsConn, err := ssh.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}
	defer wsConn.Close()

	// 接收终端设置
	var settings TerminalSettings
	if err := wsConn.ReadJSON(&settings); err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: Invalid terminal settings: %v", err)))
		return
	}

	// 获取主机信息
	host, err := h.hostRepo.FindByID(settings.HostID)
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("Error: Host not found"))
		return
	}

	// 提取IP地址 (假设公网IP是JSON数组格式)
	var ipList []string
	if err := json.Unmarshal(host.PublicIP, &ipList); err != nil || len(ipList) == 0 {
		wsConn.WriteMessage(websocket.TextMessage, []byte("Error: No public IP configured for the host"))
		return
	}
	targetIP := ipList[0]

	// 获取认证信息
	username := host.Username
	var password string

	// 如果提供了覆盖的认证信息，则使用它
	if settings.AuthUsername != "" {
		username = settings.AuthUsername
	}

	if settings.AuthPassword != "" {
		password = settings.AuthPassword
	} else {
		// 解密主机保存的密码
		key := []byte(config.GlobalConfig.App.AesKey)
		var err error
		password, err = encryption.DecryptAES(host.Password, key)
		if err != nil {
			wsConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: Failed to decrypt host credentials: %v", err)))
			return
		}
	}

	// 设置SSH客户端配置
	sshConfig := &gosshtool.ClientConfig{
		User: username,
		Auth: []gosshtool.AuthMethod{
			gosshtool.Password(password),
		},
		HostKeyCallback: gosshtool.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	// 如果提供了私钥，添加私钥认证
	if settings.PrivateKey != "" {
		signer, err := gosshtool.ParsePrivateKey([]byte(settings.PrivateKey))
		if err == nil {
			sshConfig.Auth = append(sshConfig.Auth, gosshtool.PublicKeys(signer))
		}
	}

	// 连接SSH服务器
	client, err := gosshtool.Dial("tcp", fmt.Sprintf("%s:22", targetIP), sshConfig)
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: SSH connection failed: %v", err)))
		return
	}
	defer client.Close()

	// 创建新会话
	session, err := client.NewSession()
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: Failed to create session: %v", err)))
		return
	}
	defer session.Close()

	// 请求伪终端
	termWidth := 80
	termHeight := 24
	if settings.Width > 0 {
		termWidth = settings.Width
	}
	if settings.Height > 0 {
		termHeight = settings.Height
	}

	err = session.RequestPty("xterm-256color", termHeight, termWidth, gosshtool.TerminalModes{
		gosshtool.ECHO:          1,
		gosshtool.TTY_OP_ISPEED: 14400,
		gosshtool.TTY_OP_OSPEED: 14400,
	})
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: Failed to request PTY: %v", err)))
		return
	}

	// 创建终端实例并启动连接
	terminal := ssh.NewSshTerminal(session, wsConn)
	if err := terminal.Connect(); err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
	}
}

// GetTerminalPreferences 获取终端偏好设置
// @Summary 获取终端偏好设置
// @Description 获取当前用户的终端偏好设置
// @Tags CMDB-终端管理
// @Produce json
// @Success 200 {object} TerminalPreference
// @Router /api/v1/cmdb/terminal/preferences [get]
func (h *WebTerminalHandler) GetTerminalPreferences(c *gin.Context) {
	// 默认偏好设置
	defaultPrefs := TerminalPreference{
		FontSize:    14,
		Theme:       "dark",
		Fullscreen:  false,
		CursorStyle: "block",
	}

	// TODO: 实际实现应从数据库读取用户偏好设置
	// 此处简单返回默认值
	c.JSON(200, defaultPrefs)
}

// SaveTerminalPreferences 保存终端偏好设置
// @Summary 保存终端偏好设置
// @Description 保存当前用户的终端偏好设置
// @Tags CMDB-终端管理
// @Accept json
// @Produce json
// @Param preferences body TerminalPreference true "终端偏好设置"
// @Success 200 {object} response.Response
// @Router /api/v1/cmdb/terminal/preferences [put]
func (h *WebTerminalHandler) SaveTerminalPreferences(c *gin.Context) {
	var prefs TerminalPreference
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(400, gin.H{"error": "Invalid preferences format"})
		return
	}

	// TODO: 实际实现应将设置保存到数据库

	c.JSON(200, gin.H{"message": "Preferences saved successfully"})
}

// 辅助函数：从上下文获取用户ID
func getUserIDFromContext(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	id, ok := userID.(uint)
	if !ok {
		return 0
	}
	return id
}

// RegisterRoutes 注册终端相关路由
func (h *WebTerminalHandler) RegisterRoutes(r *gin.RouterGroup) {
	// WebSocket终端连接路由（需要JWT认证）
	wsGroup := r.Group("/ws")
	wsGroup.Use(middleware.JWTAuth())
	{
		wsGroup.GET("/terminal", h.HandleSSHConnection)
	}

	// 终端偏好设置路由（需要JWT认证）
	terminalGroup := r.Group("/cmdb/terminal")
	terminalGroup.Use(middleware.JWTAuth())
	{
		terminalGroup.GET("/preferences", h.GetTerminalPreferences)
		terminalGroup.PUT("/preferences", h.SaveTerminalPreferences)
	}
}
