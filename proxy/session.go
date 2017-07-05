package proxy

import (
	"encoding/json"
	"net"
	"time"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
	"github.com/panshiqu/framework/utils"
)

// Session 会话
type Session struct {
	client, login, game net.Conn

	userid int // 用户编号
}

// OnMessage 收到消息
func (s *Session) OnMessage(mcmd uint16, scmd uint16, data []byte) (err error) {
	defer utils.Trace("Session OnMessage", mcmd, scmd)()

	switch mcmd {
	case define.LoginCommon:
		if scmd == define.LoginFastRegister {
			s.closeLogin()

			if s.login, err = sins.Dial(define.ServiceLogin,
				define.GameUnknown, define.LevelUnknown); err != nil {
				return err
			}

			go s.RecvMessage(s.login)

			// 填充客户端地址
			fastRegister := &define.FastRegister{}

			if err = json.Unmarshal(data, fastRegister); err != nil {
				return err
			}

			fastRegister.IP, _, _ = net.SplitHostPort(s.client.RemoteAddr().String())

			if data, err = json.Marshal(fastRegister); err != nil {
				return err
			}
		}

		if s.login == nil {
			s.client.Close()
			return nil
		}

		return network.SendMessage(s.login, mcmd, scmd, data)

	case define.GameCommon:
		switch scmd {
		case define.GameFastLogin:
			s.closeGame()

			fastLogin := &define.FastLogin{}

			if err = json.Unmarshal(data, fastLogin); err != nil {
				return err
			}

			if s.game, err = sins.Dial(define.ServiceGame,
				fastLogin.GameType, fastLogin.GameLevel); err != nil {
				return err
			}

			go s.RecvMessage(s.game)

			newFastLogin := &define.FastLogin{
				UserID:    s.userid,
				Timestamp: time.Now().Unix(),
			}

			newFastLogin.Signature = utils.Signature(newFastLogin.Timestamp)

			if data, err = json.Marshal(newFastLogin); err != nil {
				return err
			}

		case define.GameLogout:
			s.closeGame()
			return nil
		}

		if s.game == nil {
			s.client.Close()
			return nil
		}

		return network.SendMessage(s.game, mcmd, scmd, data)
	}

	return nil
}

// OnClose 连接关闭
func (s *Session) OnClose() {
	s.closeLogin()
	s.closeGame()
}

// RecvMessage 收到消息
func (s *Session) RecvMessage(conn net.Conn) {
	defer utils.Trace("Session RecvMessage")()

	for {
		mcmd, scmd, data, err := network.RecvMessage(conn)
		if err != nil {
			break
		}

		if mcmd == define.LoginCommon && scmd == define.LoginFastRegister {
			replyFastRegister := &define.ReplyFastRegister{}
			json.Unmarshal(data, replyFastRegister)
			s.userid = replyFastRegister.UserID
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

func (s *Session) closeLogin() {
	if s.login != nil {
		s.login.Close()
		s.login = nil
	}
}

func (s *Session) closeGame() {
	if s.game != nil {
		s.game.Close()
		s.game = nil
	}
}
