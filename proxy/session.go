package proxy

import (
	"net"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
)

// Session 会话
type Session struct {
	client, login, game net.Conn
}

// OnMessage 收到消息
func (s *Session) OnMessage(mcmd uint16, scmd uint16, data []byte) error {
	if mcmd == define.LoginCommon && scmd == define.LoginFastRegister {
		if s.login != nil {
			s.login.Close()
		}

		conn, err := net.Dial("tcp", "127.0.0.1:8081")
		if err != nil {
			return err
		}

		go s.RecvMessage(conn)

		s.login = conn
	}

	return network.SendMessage(s.login, mcmd, scmd, data)
}

// OnClose 连接关闭
func (s *Session) OnClose() {
	if s.login != nil {
		s.login.Close()
	}

	if s.game != nil {
		s.game.Close()
	}
}

// RecvMessage 收到消息
func (s *Session) RecvMessage(conn net.Conn) {
	for {
		mcmd, scmd, data, err := network.RecvMessage(conn)
		if err != nil {
			break
		}

		network.SendMessage(s.client, mcmd, scmd, data)
	}
}

// NewSession 创建会话
func NewSession(client net.Conn) *Session {
	return &Session{
		client: client,
	}
}
