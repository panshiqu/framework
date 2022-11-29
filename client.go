package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/game/fiveinarow"
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

	if mcmd == define.GlobalCommon && scmd == define.GlobalKeepAlive {
		p.client.SendMessage(mcmd, scmd, nil)
	}

	if mcmd == define.LoginCommon {
		switch scmd {
		case define.LoginFastRegister:
			replyFastRegister := &define.ReplyFastRegister{}

			if err := json.Unmarshal(data, replyFastRegister); err != nil {
				log.Println("json.Unmarshal replyFastRegister", err)
				return
			}

			// 记录用户编号
			uid = replyFastRegister.UserID

			// 快速登陆
			fastLogin := &define.FastLogin{}
			fastLogin.GameType = define.GameFiveInARow
			fastLogin.GameLevel = define.LevelOne

			// 发送快速登陆消息
			if err := p.client.SendJSONMessage(define.GameCommon, define.GameFastLogin, fastLogin); err != nil {
				log.Println("GameFastLogin SendJSONMessage", err)
				return
			}

			log.Println("GameFastLogin", fastLogin)

			// 发送签到天数
			if err := p.client.SendMessage(define.LoginCommon, define.LoginSignInDays, nil); err != nil {
				log.Println("LoginSignInDays SendMessage", err)
				return
			}

		case define.LoginSignInDays:
			replySignInDays := &define.ReplySignInDays{}

			if err := json.Unmarshal(data, replySignInDays); err != nil {
				log.Println("json.Unmarshal replySignInDays", err)
				return
			}

			// 发送签到
			if err := p.client.SendMessage(define.LoginCommon, define.LoginSignIn, nil); err != nil {
				log.Println("LoginSignIn SendMessage", err)
				return
			}

		case define.LoginSignIn:
			replySignIn := &define.ReplySignIn{}

			if err := json.Unmarshal(data, replySignIn); err != nil {
				log.Println("json.Unmarshal replySignIn", err)
				return
			}
		}
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

	// 五子棋
	if mcmd == define.GameTable {
		switch scmd {
		// 广播开始
		case fiveinarow.GameBroadcastStart:
			broadcastStart := &fiveinarow.BroadcastStart{}

			if err := json.Unmarshal(data, broadcastStart); err != nil {
				log.Println("json.Unmarshal BroadcastStart", err)
				return
			}

			if broadcastStart.ChairID == cid {
				go p.onUserInput()
			}

		// 广播落子
		case fiveinarow.GameBroadcastPlaceStone:
			broadcastPlaceStone := &fiveinarow.BroadcastPlaceStone{}

			if err := json.Unmarshal(data, broadcastPlaceStone); err != nil {
				log.Println("json.Unmarshal BroadcastPlaceStone", err)
				return
			}

			if !broadcastPlaceStone.IsWin && broadcastPlaceStone.ChairID != cid {
				go p.onUserInput()
			}
		}
	}
}

func (p *Processor) onUserInput() {
	fmt.Println("onUserInput")
	placeStone := &fiveinarow.PlaceStone{}
	fmt.Scan(&placeStone.PositionX, &placeStone.PositionY)
	p.client.SendJSONMessage(define.GameTable, fiveinarow.GamePlaceStone, placeStone)
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

func httpStart(client *network.Client) {
	http.HandleFunc("/sendmessage", func(w http.ResponseWriter, r *http.Request) {
		mcmd, merr := strconv.Atoi(r.FormValue("mcmd"))
		if merr != nil {
			fmt.Fprintln(w, merr)
			log.Println(merr)
			return
		}

		scmd, serr := strconv.Atoi(r.FormValue("scmd"))
		if serr != nil {
			fmt.Fprintln(w, serr)
			log.Println(serr)
			return
		}

		if err := client.SendMessage(uint16(mcmd), uint16(scmd), []byte(r.FormValue("data"))); err != nil {
			fmt.Fprintln(w, "SendMessage", err)
			log.Println("SendMessage", err)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

var account = flag.String("account", "panshiqu", "account")

func main() {
	flag.Parse()
	client := network.NewClient("127.0.0.1:8888")
	client.Register(&Processor{client: client})
	go handleSignal(client)
	go httpStart(client)
	client.Start()
}
