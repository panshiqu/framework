package proxy

import (
	"net"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
	"github.com/panshiqu/framework/utils"
)

// Session 会话
type Session struct {
	client, login, game net.Conn
}

// OnMessage 收到消息
func (s *Session) OnMessage(mcmd uint16, scmd uint16, data []byte) (err error) {
	defer utils.Trace("Session OnMessage", mcmd, scmd)()

	if mcmd == define.LoginCommon && scmd == define.LoginFastRegister {
		if s.login != nil {
			s.login.Close()
			s.login = nil
		}

		if s.login, err = net.Dial("tcp", "127.0.0.1:8081"); err != nil {
			return err
		}

		go s.RecvMessage(s.login)
	}

	return network.SendMessage(s.login, mcmd, scmd, data)
}

// OnClose 连接关闭
func (s *Session) OnClose() {
	defer utils.Trace("Session OnClose")()

	if s.login != nil {
		s.login.Close()
		s.login = nil
	}

	if s.game != nil {
		s.game.Close()
		s.game = nil
	}
}

// RecvMessage 收到消息
func (s *Session) RecvMessage(conn net.Conn) {
	defer utils.Trace("Session RecvMessage")()

	for {
		mcmd, scmd, data, err := network.RecvMessage(conn)
		if err != nil {
			break
		}

		network.SendMessage(s.client, mcmd, scmd, data)
	}

	s.client.Close()
}

// NewSession 创建会话
func NewSession(client net.Conn) *Session {
	return &Session{
		client: client,
	}
}
