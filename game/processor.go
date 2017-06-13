package game

import (
	"encoding/json"
	"log"
	"net"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
)

// Processor 处理器
type Processor struct {
	server *network.Server    // 服务器
	client *network.Client    // 客户端
	config *define.ConfigGame // 配置
}

// OnMessage 收到消息
func (p *Processor) OnMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	log.Println("OnMessage", mcmd, scmd, string(data))
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
		ServiceType: define.ServiceGame,
		IsServe:     true,
	}

	data, err := json.Marshal(service)
	if err != nil {
		log.Println("OnClientConnect Marshal", err)
		return
	}

	// 发送注册服务消息
	if err := p.client.SendMessage(define.ManagerCommon,
		define.ManagerRegisterService, data); err != nil {
		log.Println("OnClientConnect SendMessage", err)
		return
	}

	log.Println("OnClientConnect", service)
}

// NewProcessor 创建处理器
func NewProcessor(server *network.Server, client *network.Client, config *define.ConfigGame) *Processor {
	return &Processor{
		server: server,
		client: client,
		config: config,
	}
}
