package db

import (
	"encoding/json"
	"log"
	"net"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
)

// Processor 处理器
type Processor struct {
	server *network.Server // 服务器
}

// OnMessage 收到消息
func (p *Processor) OnMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	log.Println("OnMessage", mcmd, scmd, string(data))

	switch mcmd {
	case define.DBCommon:
		return p.OnMainCommon(conn, scmd, data)
	}

	return define.ErrUnknownMainCmd
}

// OnMainCommon 通用主命令
func (p *Processor) OnMainCommon(conn net.Conn, scmd uint16, data []byte) error {
	switch scmd {
	case define.DBFastRegister:
		return p.OnSubFastRegister(conn, data)
	}

	return define.ErrUnknownSubCmd
}

// OnSubFastRegister 快速注册子命令
func (p *Processor) OnSubFastRegister(conn net.Conn, data []byte) error {
	fastRegister := &define.FastRegister{}

	if err := json.Unmarshal(data, fastRegister); err != nil {
		return err
	}

	log.Println(fastRegister)

	replyFastRegister := &define.ReplyFastRegister{
		UserID: 10000000,
	}

	if err := network.SendJSONMessage(conn, define.DBCommon, define.DBFastRegister, replyFastRegister); err != nil {
		return err
	}

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
