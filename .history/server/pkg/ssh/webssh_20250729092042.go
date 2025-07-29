package ssh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// TerminalMessage 终端消息结构
type TerminalMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// ResizeMessage 调整终端大小消息
type ResizeMessage struct {
	Cols int `json:"cols"`
	Rows int `json:"rows"`
}

// NewSSHClient 创建SSH客户端
func NewSSHClient(conf *SSHClientConfig) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		Timeout:         conf.Timeout,
		User:            conf.UserName,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略know_hosts检查
	}

	switch strings.ToUpper(conf.AuthModel) {
	case "PASSWORD":
		config.Auth = []ssh.AuthMethod{ssh.Password(conf.Password)}
	case "PUBLICKEY":
		signer, err := ssh.ParsePrivateKey([]byte(conf.PublicKey))
		if err != nil {
			return nil, err
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	default:
		return nil, fmt.Errorf("AuthModel %s is not supported", conf.AuthModel)
	}

	c, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", conf.IP, conf.Port), config)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Turn WebSSH会话管理
type Turn struct {
	StdinPipe io.WriteCloser
	Session   *ssh.Session
	WsConn    *websocket.Conn
	mu        sync.Mutex
	cols      int
	rows      int
}

// NewTurn 创建新的WebSSH会话
func NewTurn(wsConn *websocket.Conn, sshClient *ssh.Client) (*Turn, error) {
	sess, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	stdinPipe, err := sess.StdinPipe()
	if err != nil {
		return nil, err
	}

	turn := &Turn{
		StdinPipe: stdinPipe,
		Session:   sess,
		WsConn:    wsConn,
		cols:      80,  // 默认终端宽度
		rows:      24,  // 默认终端高度
	}

	sess.Stdout = turn
	sess.Stderr = turn

	// 设置伪终端
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echo
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := sess.RequestPty("xterm-256color", turn.rows, turn.cols, modes); err != nil {
		return nil, err
	}

	if err := sess.Shell(); err != nil {
		return nil, err
	}

	return turn, nil
}

// Write 实现io.Writer接口，将SSH输出写入WebSocket
func (t *Turn) Write(p []byte) (n int, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 发送二进制数据而不是文本消息，避免编码问题
	err = t.WsConn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

// Close 关闭会话
func (t *Turn) Close() error {
	if t.Session != nil {
		t.Session.Close()
	}
	return t.WsConn.Close()
}

// ResizeTerminal 调整终端大小
func (t *Turn) ResizeTerminal(cols, rows int) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	t.cols = cols
	t.rows = rows
	
	return t.Session.WindowChange(rows, cols)
}

// LoopRead 循环读取WebSocket数据并写入SSH
func (t *Turn) LoopRead(ctx context.Context) error {
	// 设置心跳检测
	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return errors.New("LoopRead exit")
		case <-pingTicker.C:
			// 发送ping消息保持连接
			if err := t.WsConn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return fmt.Errorf("ping failed: %s", err)
			}
		default:
			// 设置读取超时
			t.WsConn.SetReadDeadline(time.Now().Add(60 * time.Second))
			
			msgType, wsData, err := t.WsConn.ReadMessage()
			if err != nil {
				return fmt.Errorf("reading webSocket message err:%s", err)
			}

			switch msgType {
			case websocket.TextMessage:
				// 处理控制消息（如调整终端大小）
				if err := t.handleControlMessage(wsData); err != nil {
					// 如果不是控制消息，则作为输入处理
					if _, err := t.StdinPipe.Write(wsData); err != nil {
						return fmt.Errorf("StdinPipe write err:%s", err)
					}
				}
			case websocket.BinaryMessage:
				// 直接写入二进制数据
				if _, err := t.StdinPipe.Write(wsData); err != nil {
					return fmt.Errorf("StdinPipe write err:%s", err)
				}
			case websocket.PongMessage:
				// 处理pong消息，重置心跳
				continue
			}
		}
	}
}

// handleControlMessage 处理控制消息
func (t *Turn) handleControlMessage(data []byte) error {
	var msg TerminalMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return err // 不是JSON格式，返回错误让上层作为普通输入处理
	}

	switch msg.Type {
	case "resize":
		resizeData, ok := msg.Data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid resize data")
		}
		
		cols, _ := resizeData["cols"].(float64)
		rows, _ := resizeData["rows"].(float64)
		
		if cols > 0 && rows > 0 {
			return t.ResizeTerminal(int(cols), int(rows))
		}
	case "ping":
		// 响应ping消息
		response := TerminalMessage{
			Type: "pong",
			Data: time.Now().Unix(),
		}
		responseData, _ := json.Marshal(response)
		return t.WsConn.WriteMessage(websocket.TextMessage, responseData)
	}
	
	return nil
}

// SessionWait 等待SSH会话结束
func (t *Turn) SessionWait() error {
	if err := t.Session.Wait(); err != nil {
		return err
	}
	return nil
}
