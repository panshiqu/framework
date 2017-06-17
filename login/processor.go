package login

import (
	"encoding/json"
	"fmt"
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

	return &define.MyError{Errno: 1, Errdesc: fmt.Sprint("unknown main cmd ", mcmd)}
}

// OnMainCommon 通用主命令
func (p *Processor) OnMainCommon(conn net.Conn, scmd uint16, data []byte) error {
	switch scmd {
	case define.LoginFastRegister:
		return p.OnSubFastRegister(conn, data)
	}

	return &define.MyError{Errno: 1, Errdesc: fmt.Sprint("unknown sub cmd ", scmd)}
}

// OnSubFastRegister 快速注册子命令
func (p *Processor) OnSubFastRegister(conn net.Conn, data []byte) error {
	fastRegister := &define.FastRegister{}

	if err := json.Unmarshal(data, fastRegister); err != nil {
		return err
	}

	// 获取客户端地址
	fastRegister.IP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())

	log.Println(fastRegister)

	return nil
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
