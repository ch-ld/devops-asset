package ssh

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

// ReaderFunc is an adapter to allow the use of ordinary functions as io.Reader.
// This is added for compatibility in case the Go version does not have it built-in.
type ReaderFunc func(p []byte) (n int, err error)

// Read calls f(p).
func (f ReaderFunc) Read(p []byte) (n int, err error) {
	return f(p)
}

// SshTerminal represents a single SSH session over WebSocket.
type SshTerminal struct {
	mu          sync.Mutex
	Session     *ssh.Session
	wsConn      *websocket.Conn
	resizeEvent chan termSize
}

// termSize defines the structure for terminal resize events.
type termSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// WebSocketMessage defines the structure for WebSocket messages
type WebSocketMessage struct {
	Type      string `json:"type"`
	Data      string `json:"data,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// NewSshTerminal creates a new SshTerminal instance.
func NewSshTerminal(session *ssh.Session, wsConn *websocket.Conn) *SshTerminal {
	return &SshTerminal{
		Session:     session,
		wsConn:      wsConn,
		resizeEvent: make(chan termSize),
	}
}

// Connect establishes the full-duplex communication between WebSocket and SSH session.
func (st *SshTerminal) Connect() error {
	// Setup I/O piping
	wsReader, wsWriter, err := st.getWsReadWriter()
	if err != nil {
		return fmt.Errorf("failed to get ws reader/writer: %w", err)
	}

	sshIn, err := st.Session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get ssh stdin pipe: %w", err)
	}
	sshOut, err := st.Session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get ssh stdout pipe: %w", err)
	}
	sshErr, err := st.Session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get ssh stderr pipe: %w", err)
	}

	// Start piping data
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, _ = io.Copy(sshIn, wsReader)
	}()
	go func() {
		defer wg.Done()
		_, _ = io.Copy(wsWriter, sshOut)
		_, _ = io.Copy(wsWriter, sshErr)
	}()

	// Handle terminal resizing
	go st.handleResize()

	if err := st.Session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %w", err)
	}

	// Wait for the session to finish
	wg.Wait()
	return st.Session.Wait()
}

// Close closes the WebSocket and SSH session.
func (st *SshTerminal) Close() {
	if st.Session != nil {
		st.Session.Close()
	}
	if st.wsConn != nil {
		st.wsConn.Close()
	}
}

// handleResize listens for resize events from the WebSocket connection.
func (st *SshTerminal) handleResize() {
	for {
		select {
		case size, ok := <-st.resizeEvent:
			if !ok {
				return
			}
			st.Session.WindowChange(size.Height, size.Width)
		}
	}
}

// getWsReadWriter returns a reader and writer for the WebSocket connection.
func (st *SshTerminal) getWsReadWriter() (io.Reader, io.Writer, error) {
	// Custom WebSocket reader to handle different message types
	wsReader := func(p []byte) (n int, err error) {
		for {
			msgType, reader, err := st.wsConn.NextReader()
			if err != nil {
				return 0, err
			}
			if msgType == websocket.TextMessage {
				// Handle resize commands
				var size termSize
				if err := websocket.ReadJSON(st.wsConn, &size); err == nil {
					st.resizeEvent <- size
					continue
				}
			}
			if msgType == websocket.BinaryMessage || msgType == websocket.TextMessage {
				return reader.Read(p)
			}
		}
	}

	// WebSocket writer
	wsWriter, err := st.wsConn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return nil, nil, err
	}

	return ReaderFunc(wsReader), wsWriter, nil
}

// Upgrader is a websocket upgrader with a larger buffer.
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections
		return true
	},
	HandshakeTimeout: 5 * time.Second,
}
