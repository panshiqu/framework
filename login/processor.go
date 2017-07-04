package login

import (
	"encoding/json"
	"log"
	"net"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
)

// Processor 处理器
type Processor struct {
	rpc    *network.RPC        // 数据库
	server *network.Server     // 服务器
	client *network.Client     // 客户端
	config *define.ConfigLogin // 配置
}

// OnMessage 收到消息
func (p *Processor) OnMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	log.Println("OnMessage", mcmd, scmd, string(data))

	switch mcmd {
	case define.LoginCommon:
		return p.OnMainCommon(conn, scmd, data)
	}

	return define.ErrUnknownMainCmd
}

// OnMainCommon 通用主命令
func (p *Processor) OnMainCommon(conn net.Conn, scmd uint16, data []byte) error {
	switch scmd {
	case define.LoginFastRegister:
		return p.OnSubFastRegister(conn, data)
	}

	return define.ErrUnknownSubCmd
}

// OnSubFastRegister 快速注册子命令
func (p *Processor) OnSubFastRegister(conn net.Conn, data []byte) error {
	fastRegister := &define.FastRegister{}
	replyFastRegister := &define.ReplyFastRegister{}

	if err := json.Unmarshal(data, fastRegister); err != nil {
		return err
	}

	// 数据库请求
	if err := p.rpc.JSONCall(define.DBCommon, define.DBFastRegister, fastRegister, replyFastRegister); err != nil {
		return err
	}

	// 只更新不查询字段
	replyFastRegister.UserName = fastRegister.Name
	replyFastRegister.UserIcon = fastRegister.Icon
	replyFastRegister.UserGender = fastRegister.Gender

	// 回复客户端
	return network.SendJSONMessage(conn, define.LoginCommon, define.LoginFastRegister, replyFastRegister)
}

// OnClose 连接关闭
func (p *Processor) OnClose(conn net.Conn) {

}

// OnClientMessage 客户端收到消息
func (p *Processor) OnClientMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) {
	log.Println("OnClientMessage", mcmd, scmd, string(data))
}

// OnClientConnect 客户端连接成功
func (p *Processor) OnClientConnect(conn net.Conn) {
	// 构造服务
	service := &define.Service{
		ID:          p.config.ID,
		IP:          p.config.ListenIP,
		ServiceType: define.ServiceLogin,
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
func NewProcessor(server *network.Server, client *network.Client, config *define.ConfigLogin) *Processor {
	return &Processor{
		rpc:    network.NewRPC(config.DBIP),
		server: server,
		client: client,
		config: config,
	}
}
