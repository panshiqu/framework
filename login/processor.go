package login

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net"

	"../define"
	"../network"
	"../utils"
)

var logger *log.Logger

// Processor 处理器
type Processor struct {
	rpc    *network.RPC    // 数据库
	server *network.Server // 服务器
	client *network.Client // 客户端
	config *define.GConfig // 配置
}

// OnMessage 收到消息
func (p *Processor) OnMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	logger.WithFields(log.Fields{
		"mcmd": mcmd,
		"scmd": scmd,
		"data": string(data),
	}).Info("OnMessage receive message")

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
	case define.LoginRegisterCheck:
		return p.onSubFasterRegisterCheck(conn,data)
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

// 注册检查
func (p* Processor) onSubFasterRegisterCheck(conn net.Conn,data []byte)error {
	registerCheck := &define.FastRegisterCheck{}
	ret := &define.MyError{}

	if err := json.Unmarshal(data, registerCheck); err != nil {
		return err
	}
	// 数据库请求
	if err := p.rpc.JSONCall(define.DBCommon, define.DBRegisterCheck, registerCheck, ret); err != nil {
		return err
	}

	return network.SendJSONMessage(conn,define.LoginCommon, define.LoginRegisterCheck, ret)
}

// OnClose 连接关闭
func (p *Processor) OnClose(conn net.Conn) {

}

// OnClientMessage 客户端收到消息
func (p *Processor) OnClientMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) {
	logger.WithFields(log.Fields{
		"mcmd": mcmd,
		"scmd": scmd,
		"data": string(data),
	}).Info("OnClientMessage receive message")
}

// OnClientConnect 客户端连接成功
func (p *Processor) OnClientConnect(conn net.Conn) {
	// 构造服务
	service := &define.Service{
		ID:          p.config.Login.ID,
		IP:          p.config.Login.ListenIP,
		ServiceType: define.ServiceLogin,
		IsServe:     true,
	}

	// 发送注册服务消息
	if err := p.client.SendJSONMessage(define.ManagerCommon, define.ManagerRegisterService, service); err != nil {
		logger.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("OnClientConnect SendJSONMessage Error")
		return
	}

	logger.WithFields(log.Fields{
		"service": service,
	}).Info("OnClientConnect register service")
}

// NewProcessor 创建处理器
func NewProcessor(server *network.Server, client *network.Client, config *define.GConfig) *Processor {
	if logger == nil {
		logger = utils.GetLogger("login")
	}
	return &Processor{
		rpc:    network.NewRPC(config.DB.ListenIP),
		server: server,
		client: client,
		config: config,
	}
}
