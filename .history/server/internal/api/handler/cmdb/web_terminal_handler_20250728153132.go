package cmdb

import (
	"api-server/internal/config"
	"api-server/internal/middleware/middleware"
	model "api-server/internal/model/cmdb"
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
	HostID          uint               `json:"host_id"`           // 主机ID
	Width           int                `json:"width"`             // 终端宽度
	Height          int                `json:"height"`            // 终端高度
	Preferences     TerminalPreference `json:"preferences"`       // 偏好设置
	AuthUsername    string             `json:"username"`          // 覆盖主机默认用户名
	AuthPassword    string             `json:"password"`          // 覆盖主机默认密码
	PrivateKey      string             `json:"private_key"`       // 使用私钥认证
	PreferredIPType string             `json:"preferred_ip_type"` // IP类型偏好：public, private, auto
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

	// 智能选择IP地址
	targetIP, err := selectTargetIP(host, settings.PreferredIPType)
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	// 获取认证信息
	username := host.Username
	var password string
	var privateKey string
	port := host.Port
	if port == 0 {
		port = 22 // 默认SSH端口
	}

	// 如果提供了覆盖的认证信息，则使用它
	if settings.AuthUsername != "" {
		username = settings.AuthUsername
	}

	if settings.AuthPassword != "" {
		password = settings.AuthPassword
	} else if host.Password != "" {
		// 解密主机保存的密码
		key := []byte(config.GlobalConfig.App.AesKey)
		var err error
		password, err = encryption.DecryptAES(host.Password, key)
		if err != nil {
			wsConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: Failed to decrypt host credentials: %v", err)))
			return
		}
	}

	// 获取私钥（如果存在）
	if settings.PrivateKey != "" {
		privateKey = settings.PrivateKey
	} else if host.PrivateKey != "" {
		// 解密主机保存的私钥
		key := []byte(config.GlobalConfig.App.AesKey)
		var err error
		privateKey, err = encryption.DecryptAES(host.PrivateKey, key)
		if err != nil {
			wsConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: Failed to decrypt host private key: %v", err)))
			return
		}
	}

	// 设置SSH客户端配置
	sshConfig := &gosshtool.ClientConfig{
		User:            username,
		Auth:            []gosshtool.AuthMethod{},
		HostKeyCallback: gosshtool.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	// 添加认证方法
	if privateKey != "" {
		// 优先使用私钥认证
		signer, err := gosshtool.ParsePrivateKey([]byte(privateKey))
		if err == nil {
			sshConfig.Auth = append(sshConfig.Auth, gosshtool.PublicKeys(signer))
		}
	}

	if password != "" {
		// 添加密码认证
		sshConfig.Auth = append(sshConfig.Auth, gosshtool.Password(password))
	}

	// 如果没有任何认证方法，返回错误
	if len(sshConfig.Auth) == 0 {
		wsConn.WriteMessage(websocket.TextMessage, []byte("Error: No valid authentication method available"))
		return
	}

	// 连接SSH服务器，使用动态端口
	client, err := gosshtool.Dial("tcp", fmt.Sprintf("%s:%d", targetIP, port), sshConfig)
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

// selectTargetIP 智能选择目标IP地址
func selectTargetIP(host *model.Host, preferredType string) (string, error) {
	// 解析公网IP
	var publicIPs []string
	if len(host.PublicIP) > 2 { // 检查是否为非空JSON数组
		if err := json.Unmarshal(host.PublicIP, &publicIPs); err != nil {
			publicIPs = []string{}
		}
	}

	// 解析私网IP
	var privateIPs []string
	if len(host.PrivateIP) > 2 { // 检查是否为非空JSON数组
		if err := json.Unmarshal(host.PrivateIP, &privateIPs); err != nil {
			privateIPs = []string{}
		}
	}

	// 过滤空字符串
	publicIPs = filterEmptyStrings(publicIPs)
	privateIPs = filterEmptyStrings(privateIPs)

	// 根据偏好类型选择IP
	switch preferredType {
	case "public":
		if len(publicIPs) > 0 {
			return publicIPs[0], nil
		}
		return "", fmt.Errorf("no public IP available for this host")
	case "private":
		if len(privateIPs) > 0 {
			return privateIPs[0], nil
		}
		return "", fmt.Errorf("no private IP available for this host")
	default: // "auto" 或其他值
		// 优先使用公网IP，如果没有则使用私网IP
		if len(publicIPs) > 0 {
			return publicIPs[0], nil
		}
		if len(privateIPs) > 0 {
			return privateIPs[0], nil
		}
		return "", fmt.Errorf("no IP address available for this host")
	}
}

// filterEmptyStrings 过滤空字符串
func filterEmptyStrings(strs []string) []string {
	var result []string
	for _, str := range strs {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}
