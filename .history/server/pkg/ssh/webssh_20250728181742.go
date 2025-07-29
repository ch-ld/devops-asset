package ssh

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

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
	}

	sess.Stdout = turn
	sess.Stderr = turn

	// 设置伪终端
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echo
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := sess.RequestPty("xterm", 150, 30, modes); err != nil {
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

	writer, err := t.WsConn.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, err
	}
	defer writer.Close()

	return writer.Write(p)
}

// Close 关闭会话
func (t *Turn) Close() error {
	if t.Session != nil {
		t.Session.Close()
	}
	return t.WsConn.Close()
}

// Read 从WebSocket读取数据
func (t *Turn) Read(p []byte) (n int, err error) {
	for {
		msgType, reader, err := t.WsConn.NextReader()
		if err != nil {
			return 0, err
		}
		if msgType != websocket.TextMessage {
			continue
		}
		return reader.Read(p)
	}
}

// LoopRead 循环读取WebSocket数据并写入SSH
func (t *Turn) LoopRead(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return errors.New("LoopRead exit")
		default:
			_, wsData, err := t.WsConn.ReadMessage()
			if err != nil {
				return fmt.Errorf("reading webSocket message err:%s", err)
			}

			// 处理数据，如果是base64编码则解码
			body := wsData
			if len(wsData) > 1 && wsData[0] == 'b' {
				// 尝试base64解码
				if decoded, err := base64.StdEncoding.DecodeString(string(wsData[1:])); err == nil {
					body = decoded
				}
			}

			if _, err := t.StdinPipe.Write(body); err != nil {
				return fmt.Errorf("StdinPipe write err:%s", err)
			}
		}
	}
}

// SessionWait 等待SSH会话结束
func (t *Turn) SessionWait() error {
	if err := t.Session.Wait(); err != nil {
		return err
	}
	return nil
}
