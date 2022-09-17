package proxy

import (
	"encoding/json"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
	"github.com/panshiqu/framework/utils"
)

// Session 会话
type Session struct {
	client, login, game net.Conn

	userid int   // 用户编号
	status int32 // 会话状态
	close  chan bool
}

// OnMessage 收到消息
func (s *Session) OnMessage(mcmd uint16, scmd uint16, data []byte) (err error) {
	defer utils.Trace("Session OnMessage", mcmd, scmd)()

	atomic.StoreInt32(&s.status, define.KeepAliveSafe)

	switch mcmd {
	case define.LoginCommon:
		if scmd == define.LoginFastRegister {
			s.closeLogin()

			if s.login, err = sins.Dial(define.ServiceLogin,
				define.GameUnknown, define.LevelUnknown); err != nil {
				return utils.Wrap(err)
			}

			go s.RecvMessage(s.login)

			// 填充客户端地址
			fastRegister := &define.FastRegister{}

			if err = json.Unmarshal(data, fastRegister); err != nil {
				return utils.Wrap(err)
			}

			fastRegister.IP, _, _ = net.SplitHostPort(s.client.RemoteAddr().String())

			if data, err = json.Marshal(fastRegister); err != nil {
				return utils.Wrap(err)
			}
		}

		if s.login == nil {
			s.client.Close()
			return nil
		}

		return utils.Wrap(network.SendMessage(s.login, mcmd, scmd, data))

	case define.GameCommon, define.GameTable:
		if mcmd == define.GameCommon {
			if scmd == define.GameFastLogin {
				s.closeGame()

				fastLogin := &define.FastLogin{}

				if err = json.Unmarshal(data, fastLogin); err != nil {
					return utils.Wrap(err)
				}

				if s.game, err = sins.Dial(define.ServiceGame,
					fastLogin.GameType, fastLogin.GameLevel); err != nil {
					return utils.Wrap(err)
				}

				go s.RecvMessage(s.game)

				newFastLogin := &define.FastLogin{
					UserID:    s.userid,
					Timestamp: time.Now().Unix(),
				}

				newFastLogin.Signature = utils.Signature(newFastLogin.Timestamp)

				if data, err = json.Marshal(newFastLogin); err != nil {
					return utils.Wrap(err)
				}
			} else if scmd == define.GameLogout {
				s.closeGame()
				return nil
			}
		}

		if s.game == nil {
			s.client.Close()
			return nil
		}

		return utils.Wrap(network.SendMessage(s.game, mcmd, scmd, data))
	}

	return nil
}

// OnClose 连接关闭
func (s *Session) OnClose() {
	close(s.close)
	s.closeLogin()
	s.closeGame()
}

// RecvMessage 收到消息
func (s *Session) RecvMessage(conn net.Conn) {
	defer utils.Trace("Session RecvMessage")()

	for {
		mcmd, scmd, data, err := network.RecvMessage(conn)
		if err != nil {
			log.Println(utils.Wrap(err))
			break
		}

		if mcmd == define.LoginCommon && scmd == define.LoginFastRegister {
			replyFastRegister := &define.ReplyFastRegister{}
			json.Unmarshal(data, replyFastRegister)
			s.userid = replyFastRegister.UserID
		}

		if err := network.SendMessage(s.client, mcmd, scmd, data); err != nil {
			log.Println(mcmd, scmd, utils.Wrap(err))
			break
		}
	}

	s.client.Close()
}

// KeepAlive 保活
func (s *Session) KeepAlive() {
	defer utils.Trace("Session KeepAlive")()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.close:
			return

		case <-ticker.C:
			switch atomic.LoadInt32(&s.status) {
			case define.KeepAliveSafe:
				atomic.StoreInt32(&s.status, define.KeepAliveWarn)

			case define.KeepAliveWarn:
				atomic.StoreInt32(&s.status, define.KeepAliveDead)
				network.SendMessage(s.client, define.GlobalCommon, define.GlobalKeepAlive, nil)

			case define.KeepAliveDead:
				s.client.Close()
			}

		default:
			time.Sleep(time.Second)
		}
	}
}

// NewSession 创建会话
func NewSession(client net.Conn) *Session {
	ses := &Session{
		client: client,
		status: define.KeepAliveSafe,
		close:  make(chan bool),
	}

	go ses.KeepAlive()

	return ses
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
