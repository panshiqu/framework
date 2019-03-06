package game

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"../define"
	"../network"
	"../utils"
)

// 数据库
var rpc *network.RPC

// 全局定时器
var sins *utils.Schedule

// Processor 处理器
type Processor struct {
	server *network.Server // 服务器
	client *network.Client // 客户端
	config *define.GConfig
}

// OnTimer 定时器
func (p *Processor) OnTimer(id int, parameter interface{}) {
	if id < define.TimerPerTable {
		return
	}

	if tableFrame := tins.GetTable((id - define.TimerPerTable) / define.TimerPerTable); tableFrame != nil {
		if err := tableFrame.OnTimer((id-define.TimerPerTable)%define.TimerPerTable, parameter); err != nil {
			log.Println("TableFrame OnTimer", err)
		}
	}
}

// OnMessage 收到消息
func (p *Processor) OnMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	log.Println("OnMessage", mcmd, scmd, string(data))

	switch mcmd {
	case define.GameCommon:
		return p.OnMainCommon(conn, scmd, data)
	case define.GameTable:
		return p.OnMainTable(conn, scmd, data)
	}

	return define.ErrUnknownMainCmd
}

// OnMainCommon 通用主命令
func (p *Processor) OnMainCommon(conn net.Conn, scmd uint16, data []byte) error {
	switch scmd {
	case define.GameFastLogin:
		return p.OnSubFastLogin(conn, data)
	case define.GameReady:
		return p.OnSubReady(conn, data)
	}

	return define.ErrUnknownSubCmd
}

// OnMainTable 桌子主命令
func (p *Processor) OnMainTable(conn net.Conn, scmd uint16, data []byte) error {
	// 获取绑定用户
	userItem, ok := p.server.GetBind(conn).(*UserItem)
	if !ok {
		return define.ErrNotExistUser
	}

	// 获取桌子框架
	tableFrame := userItem.TableFrame()
	if tableFrame == nil {
		return define.ErrUserNotSit
	}

	// 通知桌子消息
	return tableFrame.OnMessage(scmd, data, userItem)
}

// OnSubFastLogin 快速登陆子命令
func (p *Processor) OnSubFastLogin(conn net.Conn, data []byte) error {
	defer utils.Trace("Processor OnSubFastLogin")()

	fastLogin := &define.FastLogin{}
	replyFastLogin := &define.ReplyFastLogin{}

	if err := json.Unmarshal(data, fastLogin); err != nil {
		return err
	}

	// 可以判断时间戳是否接近当前时间
	// 即使抓包再封包依然不能模拟登陆
	if utils.Signature(fastLogin.Timestamp) != fastLogin.Signature {
		return define.ErrSignature
	}

	// 查找用户
	if userItem := uins.Search(fastLogin.UserID); userItem != nil {
		// 设置网络连接
		userItem.SetConn(conn)

		// 设置绑定
		p.server.SetBind(conn, userItem)

		// 获取桌子框架
		if tableFrame := userItem.TableFrame(); tableFrame != nil {
			// 发送我的坐下
			userItem.SendJSONMessage(define.GameCommon, define.GameNotifySitDown, userItem.TableUserInfo())

			// 发送同桌玩家信息
			tableFrame.SendTableUserInfo(userItem)

			// 断线重连
			tableFrame.Reconnect(userItem)
		}

		// 设置游戏状态
		userItem.SetUserStatus(define.UserStatusPlaying)

		return nil
	}

	// 数据库请求
	if err := rpc.JSONCall(define.DBCommon, define.DBFastLogin, &fastLogin.UserID, replyFastLogin); err != nil {
		return err
	}

	// 插入用户
	userItem := uins.Insert(conn, replyFastLogin)

	// 设置绑定
	p.server.SetBind(conn, userItem)

	// 用户坐下
	tableFrame := tins.TrySitDown(userItem)

	// 正在游戏设置游戏状态
	if tableFrame.TableStatus() == define.TableStatusGame {
		userItem.SetUserStatus(define.UserStatusPlaying)
	}

	return nil
}

// OnSubReady 准备子命令
func (p *Processor) OnSubReady(conn net.Conn, data []byte) error {
	// 获取绑定用户
	userItem, ok := p.server.GetBind(conn).(*UserItem)
	if !ok {
		return define.ErrNotExistUser
	}

	// 获取桌子框架
	tableFrame := userItem.TableFrame()
	if tableFrame == nil {
		return define.ErrUserNotSit
	}

	// 校验桌子状态
	if tableFrame.TableStatus() == define.TableStatusGame {
		return define.ErrTableStatus
	}

	// 设置准备状态
	userItem.SetUserStatus(define.UserStatusReady)

	// 尝试开始游戏
	tableFrame.StartGame()

	return nil
}

// OnClose 连接关闭
func (p *Processor) OnClose(conn net.Conn) {
	defer utils.Trace("Processor OnClose")()

	// 获取绑定用户
	if userItem, ok := p.server.GetBind(conn).(*UserItem); ok {
		// 获取桌子框架
		if tableFrame := userItem.TableFrame(); tableFrame != nil {
			// 正在游戏设置离线状态
			if tableFrame.TableStatus() == define.TableStatusGame {
				userItem.SetUserStatus(define.UserStatusOffline)
				return
			}

			// 用户站起
			tableFrame.StandUp(userItem)
		}

		// 删除用户
		uins.Delete(userItem.UserID())
	}
}

// OnClientMessage 客户端收到消息
func (p *Processor) OnClientMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) {
	log.Println("OnClientMessage", mcmd, scmd, string(data))
}

// OnClientConnect 客户端连接成功
func (p *Processor) OnClientConnect(conn net.Conn) {
	// 构造服务
	service := &define.Service{
		ID:          p.config.Game.ID,
		IP:          p.config.Game.ListenIP,
		GameType:    define.GameFiveInARow,
		GameLevel:   define.LevelOne,
		ServiceType: define.ServiceGame,
		IsServe:     true,
	}

	// 发送注册服务消息
	if err := p.client.SendJSONMessage(define.ManagerCommon, define.ManagerRegisterService, service); err != nil {
		log.Println("OnClientConnect SendJSONMessage", err)
		return
	}

	log.Println("OnClientConnect", service)
}

// NewProcessor 创建处理器
func NewProcessor(server *network.Server, client *network.Client,config *define.GConfig) *Processor {
	p := &Processor{
		server: server,
		client: client,
		config: config,
	}

	rpc = network.NewRPC(config.DB.ListenIP)

	sins = utils.NewSchedule(p)
	go sins.Start()

	return p
}

// Monitor 监视器
func (p *Processor) Monitor(w http.ResponseWriter, r *http.Request) {
	uins.Monitor(w, r)
	tins.Monitor(w, r)
}

func init() {
	uins.users = make(map[int]*UserItem)
}
