package manager

import (
	"fmt"
	"log"
	"net"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
)

// Processor 处理器
type Processor struct {
	server *network.Server
}

// OnMessage 收到消息
func (p *Processor) OnMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	log.Println("OnMessage", mcmd, scmd, string(data))

	switch mcmd {
	case define.ManagerCommon:
		return p.OnMainCommon(conn, scmd, data)
	}

	return define.NewError(fmt.Sprint("unknown main cmd", mcmd))
}

// OnMainCommon 通用主命令
func (p *Processor) OnMainCommon(conn net.Conn, scmd uint16, data []byte) error {
	switch scmd {
	case define.ManagerRegisterService:
		return p.OnSubRegisterService(conn, data)
	}

	return define.NewError(fmt.Sprint("unknown sub cmd", scmd))
}

// OnSubRegisterService 注册服务子命令
func (p *Processor) OnSubRegisterService(conn net.Conn, data []byte) error {
	log.Println(string(data))
	return nil
}

// OnClose 连接关闭
func (p *Processor) OnClose(conn net.Conn) {

}

// OnClientMessage 客户端收到消息
func (p *Processor) OnClientMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) {
	// nothing to do
}

// OnClientConnect 客户端连接成功
func (p *Processor) OnClientConnect(conn net.Conn) {
	// nothing to do
}

// NewProcessor 创建处理器
func NewProcessor(server *network.Server) *Processor {
	return &Processor{
		server: server,
	}
}
