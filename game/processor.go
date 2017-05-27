package game

import (
	"log"
	"net"

	"github.com/panshiqu/framework/network"
)

// Processor 处理器
type Processor struct {
	server *network.Server
	client *network.Client
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

}

// OnClientConnect 客户端连接成功
func (p *Processor) OnClientConnect(conn net.Conn) {

}

// NewProcessor 创建处理器
func NewProcessor(server *network.Server, client *network.Client) *Processor {
	return &Processor{
		server: server,
		client: client,
	}
}
