package login

import (
	"net"

	"github.com/panshiqu/framework/network"
)

// Processor 处理器
type Processor struct {
	server *network.Server
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

}

// OnClientConnect 客户端连接成功
func (p *Processor) OnClientConnect(conn net.Conn) {

}
