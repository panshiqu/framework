package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"net"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
)

// LOG 日志数据库
var LOG *sql.DB

// GAME 游戏数据库
var GAME *sql.DB

// Processor 处理器
type Processor struct {
	server *network.Server // 服务器
}

// OnMessage 收到消息
func (p *Processor) OnMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	ret := p.OnMessageEx(conn, mcmd, scmd, data)

	// 必须回复消息
	if ret == nil {
		return define.ErrSuccess
	}

	// 错误直接回复
	if err, ok := ret.(error); ok {
		return err
	}

	// 实现快捷回复消息
	return network.SendJSONMessage(conn, mcmd, scmd, ret)
}

// OnMessageEx 收到消息
func (p *Processor) OnMessageEx(conn net.Conn, mcmd uint16, scmd uint16, data []byte) interface{} {
	log.Println("OnMessage", mcmd, scmd, string(data))

	switch mcmd {
	case define.DBCommon:
		return p.OnMainCommon(conn, scmd, data)
	}

	return define.ErrUnknownMainCmd
}

// OnMainCommon 通用主命令
func (p *Processor) OnMainCommon(conn net.Conn, scmd uint16, data []byte) interface{} {
	switch scmd {
	case define.DBFastRegister:
		return p.OnSubFastRegister(conn, data)
	}

	return define.ErrUnknownSubCmd
}

// OnSubFastRegister 快速注册子命令
func (p *Processor) OnSubFastRegister(conn net.Conn, data []byte) interface{} {
	fastRegister := &define.FastRegister{}

	if err := json.Unmarshal(data, fastRegister); err != nil {
		return err
	}

	log.Println(fastRegister)

	replyFastRegister := &define.ReplyFastRegister{
		UserID: 10000000,
	}

	return replyFastRegister
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
	var err error

	// todo SetMaxOpenConns, SetMaxIdleConns

	if LOG, err = sql.Open("mysql", "root:@/log"); err != nil {
		log.Println("Open log", err)
		return nil
	}

	if err = LOG.Ping(); err != nil {
		log.Println("Ping log", err)
		return nil
	}

	if GAME, err = sql.Open("mysql", "root:@/game"); err != nil {
		log.Println("Open game", err)
		return nil
	}

	if err = GAME.Ping(); err != nil {
		log.Println("Ping game", err)
		return nil
	}

	return &Processor{
		server: server,
	}
}
