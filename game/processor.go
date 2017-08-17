package game

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
	"github.com/panshiqu/framework/utils"
)

// Processor 处理器
type Processor struct {
	rpc    *network.RPC    // 数据库
	server *network.Server // 服务器
	client *network.Client // 客户端
}

// OnMessage 收到消息
func (p *Processor) OnMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	log.Println("OnMessage", mcmd, scmd, string(data))

	switch mcmd {
	case define.GameCommon:
		return p.OnMainCommon(conn, scmd, data)
	}

	return define.ErrUnknownMainCmd
}

// OnMainCommon 通用主命令
func (p *Processor) OnMainCommon(conn net.Conn, scmd uint16, data []byte) error {
	switch scmd {
	case define.GameFastLogin:
		return p.OnSubFastLogin(conn, data)
	}

	return define.ErrUnknownSubCmd
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
		// 设置绑定
		p.server.SetBind(conn, userItem)

		// 回复数据
		replyFastLogin.UserID = userItem.UserID()
		replyFastLogin.UserName = userItem.UserName()
		replyFastLogin.UserIcon = userItem.UserIcon()
		replyFastLogin.UserLevel = userItem.UserLevel()
		replyFastLogin.UserGender = userItem.UserGender()
		replyFastLogin.BindPhone = userItem.BindPhone()
		replyFastLogin.UserScore = userItem.UserScore()
		replyFastLogin.UserDiamond = userItem.UserDiamond()

		// 回复客户端
		return network.SendJSONMessage(conn, define.GameCommon, define.GameFastLogin, replyFastLogin)
	}

	// 数据库请求
	if err := p.rpc.JSONCall(define.DBCommon, define.DBFastLogin, &fastLogin.UserID, replyFastLogin); err != nil {
		return err
	}

	// 插入用户
	userItem := uins.Insert(conn, replyFastLogin)

	// 设置绑定
	p.server.SetBind(conn, userItem)

	// 回复客户端
	if err := network.SendJSONMessage(conn, define.GameCommon, define.GameFastLogin, replyFastLogin); err != nil {
		return err
	}

	// 用户坐下
	tins.TrySitDown(userItem)

	return nil
}

// OnClose 连接关闭
func (p *Processor) OnClose(conn net.Conn) {
	defer utils.Trace("Processor OnClose")()

	// 获取绑定用户
	if userItem, ok := p.server.GetBind(conn).(*UserItem); ok {
		if tableFrame := userItem.TableFrame(); tableFrame != nil {
			tableFrame.StandUp(userItem)
		}
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
		ID:          define.CG.ID,
		IP:          define.CG.ListenIP,
		GameType:    define.GameLandlords,
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
func NewProcessor(server *network.Server, client *network.Client) *Processor {
	return &Processor{
		rpc:    network.NewRPC(define.CG.DBIP),
		server: server,
		client: client,
	}
}

// Monitor 监视器
func (p *Processor) Monitor(w http.ResponseWriter, r *http.Request) {
	uins.Monitor(w, r)
	tins.Monitor(w, r)
}

func init() {
	uins.users = make(map[int]*UserItem)
}
