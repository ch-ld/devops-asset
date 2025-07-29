package cmdb

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"api-server/internal/config"
	repo "api-server/internal/repository/cmdb"
	"api-server/pkg/crypto/encryption"
	"api-server/pkg/ssh"

	"github.com/gin-gonic/gin"
	gosshtool "golang.org/x/crypto/ssh"
)

// WebSSH管理相关接口处理器
// 负责主机WebSSH连接与操作的HTTP/WebSocket请求处理
type WebSshHandler struct {
	hostRepo *repo.HostRepository
}

// NewWebSSHHandler 创建WebSSH处理器实例
func NewWebSshHandler(hostRepo *repo.HostRepository) *WebSshHandler {
	return &WebSshHandler{hostRepo: hostRepo}
}

// ConnectWebSSH 建立WebSSH连接接口
func (h *WebSshHandler) HandleSSH(c *gin.Context) {
	// 验证token参数（WebSocket连接无法使用标准的JWT中间件）
	token := c.Query("token")
	if token == "" {
		c.JSON(401, gin.H{"error": "Missing token parameter"})
		return
	}

	// TODO: 这里应该验证token的有效性
	// 暂时跳过token验证，但在生产环境中必须实现

	wsConn, err := ssh.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}
	defer wsConn.Close()

	hostIDStr := c.Query("host_id")
	hostID, err := strconv.Atoi(hostIDStr)
	if err != nil {
		wsConn.WriteMessage(1, []byte("Error: Invalid host_id query parameter"))
		return
	}

	host, err := h.hostRepo.FindByID(uint(hostID))
	if err != nil {
		wsConn.WriteMessage(1, []byte("Error: Host not found"))
		return
	}

	// 智能选择IP地址 (默认自动选择)
	targetIP, err := selectTargetIP(host, "auto")
	if err != nil {
		wsConn.WriteMessage(1, []byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	if host.Username == "" || host.Password == "" {
		wsConn.WriteMessage(1, []byte("Error: Host SSH credentials are not configured"))
		return
	}

	key := []byte(config.GlobalConfig.App.AesKey)
	decryptedPassword, err := encryption.DecryptAES(host.Password, key)
	if err != nil {
		wsConn.WriteMessage(1, []byte(fmt.Sprintf("Error: Failed to decrypt host credentials: %v", err)))
		return
	}

	sshClient, err := sshConnect(host.Username, decryptedPassword, targetIP, 22)
	if err != nil {
		wsConn.WriteMessage(1, []byte(fmt.Sprintf("Error: Failed to connect via SSH: %v", err)))
		return
	}
	defer sshClient.Close()

	sshSession, err := sshClient.NewSession()
	if err != nil {
		wsConn.WriteMessage(1, []byte(fmt.Sprintf("Error: Failed to create SSH session: %v", err)))
		return
	}

	terminal := ssh.NewSshTerminal(sshSession, wsConn)
	defer terminal.Close()

	if err := terminal.Connect(); err != nil {
		log.Printf("SSH terminal connection ended with error: %v", err)
	} else {
		log.Println("SSH terminal connection ended gracefully.")
	}
}

// sshConnect helper function to establish an SSH connection.
func sshConnect(user, password, host string, port int) (*gosshtool.Client, error) {
	config := &gosshtool.ClientConfig{
		User: user,
		Auth: []gosshtool.AuthMethod{
			gosshtool.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key gosshtool.PublicKey) error {
			return nil // In a real-world scenario, you should verify the host key.
		},
		Timeout: 10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	return gosshtool.Dial("tcp", addr, config)
}
