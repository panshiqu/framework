package proxy

import (
	"net"
)

// Session 会话
type Session struct {
	client, login, game net.Conn
}

// NewSession 创建会话
func NewSession(client net.Conn) *Session {
	return &Session{
		client: client,
	}
}
