package proxy

import (
	"encoding/json"
	"net"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	"../define"
	"../network"
	"../utils"
)

// Session 会话
type Session struct {
	client, login, game net.Conn

	userid uint32   // 用户编号
	status int32 // 会话状态
	close  chan bool
	log *log.Logger
}

// OnMessage 收到消息
func (s *Session) OnMessage(mcmd uint16, scmd uint16, data []byte) (err error) {
	defer utils.Trace("Session OnMessage", mcmd, scmd)()
	utils.LogMessage(s.log, "Session OnMessage", mcmd,scmd,data)
	
	atomic.StoreInt32(&s.status, define.KeepAliveSafe)

	switch mcmd {
	case define.LoginCommon:
		if mcmd == define.LoginCommon && scmd == define.LoginFastRegister {
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
		}else if mcmd == define.LoginCommon && scmd == define.LoginRegisterCheck {
			s.closeLogin()

			if s.login, err = sins.Dial(define.ServiceLogin,
				define.GameUnknown, define.LevelUnknown); err != nil {
				return err
			}
			go s.RecvMessage(s.login)
		}

		if s.login == nil {
			s.client.Close()
			return nil
		}

		return network.SendMessage(s.login, mcmd, scmd, data)

	case define.GameCommon, define.GameTable:
		if mcmd == define.GameCommon {
			if scmd == define.GameFastLogin {
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
					UserID:    fastLogin.UserID,
					Timestamp: time.Now().Unix(),
				}

				newFastLogin.Signature = utils.Signature(newFastLogin.Timestamp)

				s.log.WithFields(log.Fields{
					"NewFastLogin": newFastLogin,
				}).Info("FastLogin NewFastLogin Info")

				if data, err = json.Marshal(newFastLogin); err != nil {
					return err
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

		return network.SendMessage(s.game, mcmd, scmd, data)
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
			break
		}

		if mcmd == define.LoginCommon && scmd == define.LoginFastRegister {
			replyFastRegister := &define.ReplyFastRegister{}
			json.Unmarshal(data, replyFastRegister)
			s.userid = replyFastRegister.UserID
		}

		if err := network.SendMessage(s.client, mcmd, scmd, data); err != nil {
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
				network.SendMessage(s.client, define.GLobalCommon, define.GLobalKeepAlive, nil)

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
		log: GetProxyLogger(),
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
