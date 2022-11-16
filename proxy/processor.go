package proxy

import (
	"encoding/json"
	"log"
	"net"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
	"github.com/panshiqu/framework/utils"
)

// Processor 处理器
type Processor struct {
	server *network.Server     // 服务器
	client *network.Client     // 客户端
	config *define.ConfigProxy // 配置
}

// OnMessage 收到消息
func (p *Processor) OnMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	log.Println("OnMessage", mcmd, scmd, string(data))

	session, ok := p.server.GetBind(conn).(*Session)
	if !ok {
		log.Println("NewSession")
		session = NewSession(conn)
		p.server.SetBind(conn, session)
	}

	return utils.Wrap(session.OnMessage(mcmd, scmd, data))
}

// OnClose 连接关闭
func (p *Processor) OnClose(conn net.Conn) {
	if session, ok := p.server.GetBind(conn).(*Session); ok {
		log.Println("CloseSession")
		session.OnClose()
	}
}

// OnClientMessage 客户端收到消息
func (p *Processor) OnClientMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) {
	log.Println("OnClientMessage", mcmd, scmd, string(data))

	if mcmd != define.ManagerCommon {
		return
	}

	switch scmd {
	// 通知已选服务、所有服务
	case define.ManagerNotifyCurService, define.ManagerNotifyAllService:
		var selected map[int]*define.Service

		if err := json.Unmarshal(data, &selected); err != nil {
			return
		}

		if scmd == define.ManagerNotifyAllService {
			sins.InitAll(selected)
		} else {
			sins.Init(selected)
		}

	// 增加删除服务
	case define.ManagerNotifyAddService, define.ManagerNotifyDelService,
		define.ManagerNotifyIncrService, define.ManagerNotifyDecrService:
		service := &define.Service{}

		if err := json.Unmarshal(data, service); err != nil {
			return
		}

		switch scmd {
		case define.ManagerNotifyAddService:
			sins.Add(service)
		case define.ManagerNotifyDelService:
			sins.Del(service)
		case define.ManagerNotifyIncrService:
			sins.Incr(service)
		case define.ManagerNotifyDecrService:
			sins.Decr(service)
		}

	// 改变已选服务
	case define.ManagerNotifyChangeService:
		var services []*define.Service

		if err := json.Unmarshal(data, &services); err != nil {
			return
		}

		sins.Change(services)
	}
}

// OnClientConnect 客户端连接成功
func (p *Processor) OnClientConnect(conn net.Conn) {
	// 构造服务
	service := &define.Service{
		ID:          p.config.ID,
		IP:          p.config.ListenIP,
		ServiceType: define.ServiceProxy,
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
func NewProcessor(server *network.Server, client *network.Client, config *define.ConfigProxy) *Processor {
	return &Processor{
		server: server,
		client: client,
		config: config,
	}
}
