package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
)

var uid, cid int

// Processor 处理器
type Processor struct {
	client *network.Client
}

// OnMessage 收到消息
func (p *Processor) OnMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	return nil
}

// OnClose 连接关闭
func (p *Processor) OnClose(conn net.Conn) {

}

// OnClientMessage 客户端收到消息
func (p *Processor) OnClientMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) {
	log.Println("OnClientMessage", mcmd, scmd, string(data))

	if mcmd == define.GLobalCommon && scmd == define.GLobalKeepAlive {
		p.client.SendMessage(mcmd, scmd, nil)
	}

	if mcmd == define.LoginCommon && scmd == define.LoginFastRegister {
		replyFastRegister := &define.ReplyFastRegister{}

		if err := json.Unmarshal(data, replyFastRegister); err != nil {
			log.Println("json.Unmarshal replyFastRegister", err)
			return
		}

		// 记录用户编号
		uid = replyFastRegister.UserID

		// 快速登陆
		fastLogin := &define.FastLogin{
			GameType:  define.GameLandlords,
			GameLevel: define.LevelOne,
		}

		// 发送快速登陆消息
		if err := p.client.SendJSONMessage(define.GameCommon, define.GameFastLogin, fastLogin); err != nil {
			log.Println("GameFastLogin SendJSONMessage", err)
			return
		}

		log.Println("GameFastLogin", fastLogin)
	}

	if mcmd == define.GameCommon {
		switch scmd {
		case define.GameFastLogin:

		case define.GameNotifySitDown:
			notifySitDown := &define.NotifySitDown{}

			if err := json.Unmarshal(data, notifySitDown); err != nil {
				log.Println("json.Unmarshal NotifySitDown", err)
				return
			}

			// 自己坐下发送准备
			if notifySitDown.UserID == uid {
				p.client.SendMessage(define.GameCommon, define.GameReady, nil)

				// 记录椅子编号
				cid = notifySitDown.ChairID
			}

		case define.GameNotifyStatus:
			notifyStatus := &define.NotifyStatus{}

			if err := json.Unmarshal(data, notifyStatus); err != nil {
				log.Println("json.Unmarshal NotifyStatus", err)
				return
			}

			// 自己空闲状态再次发送准备
			if notifyStatus.ChairID == cid && notifyStatus.UserStatus == define.UserStatusFree {
				p.client.SendMessage(define.GameCommon, define.GameReady, nil)
			}
		}
	}
}

// OnClientConnect 客户端连接成功
func (p *Processor) OnClientConnect(conn net.Conn) {
	// 快速注册
	fastRegister := &define.FastRegister{
		Account:  *account,
		Password: "111111",
		Machine:  *account,
		Name:     *account,
		Icon:     1,
		Gender:   define.GenderFemale,
	}

	// 发送快速注册消息
	if err := p.client.SendJSONMessage(define.LoginCommon, define.LoginFastRegister, fastRegister); err != nil {
		log.Println("OnClientConnect SendJSONMessage", err)
		return
	}

	log.Println("OnClientConnect", fastRegister)
}

func handleSignal(client *network.Client) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	log.Println("Got signal:", s)

	client.Stop()
}

var account = flag.String("account", "panshiqu", "account")

func main() {
	flag.Parse()
	client := network.NewClient("127.0.0.1:8083")
	client.Register(&Processor{client: client})
	go handleSignal(client)
	client.Start()
}
