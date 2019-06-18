package proxy

import (
	"encoding/json"
	"net"

	"../define"
	"../network"
	"../utils"
	log "github.com/sirupsen/logrus"

)

var logger *log.Logger

// Processor 处理器
type Processor struct {
	server *network.Server     // 服务器
	client *network.Client     // 客户端
	config *define.GConfig // 配置
}

// OnMessage 收到消息
func (p *Processor) OnMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	logger.WithFields(log.Fields{
		"mcmd": mcmd,
		"scmd": scmd,
		"data": string(data),
	}).Info("Proxy processor OnMessage")

	session, ok := p.server.GetBind(conn).(*Session)
	if !ok {
		log.Println("NewSession")
		session = NewSession(conn)
		p.server.SetBind(conn, session)
	}

	return session.OnMessage(mcmd, scmd, data)
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
	// 通知已选服务
	case define.ManagerNotifyCurService:
		var selected map[int]*define.Service

		if err := json.Unmarshal(data, &selected); err != nil {
			return
		}

		sins.Init(selected)

	// 增加已选服务
	case define.ManagerNotifyAddService:
		service := &define.Service{}

		if err := json.Unmarshal(data, service); err != nil {
			return
		}

		sins.Add(service)

	// 删除已选服务
	case define.ManagerNotifyDelService:
		service := &define.Service{}

		if err := json.Unmarshal(data, service); err != nil {
			return
		}

		sins.Del(service)
	}
}

// OnClientConnect 客户端连接成功
func (p *Processor) OnClientConnect(conn net.Conn) {
	// 构造服务
	service := &define.Service{
		ID:          p.config.Proxy.ID,
		IP:          p.config.Proxy.ListenIP,
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

//返回proxy logger
func GetProxyLogger()*log.Logger  {
	return logger
}

// NewProcessor 创建处理器
func NewProcessor(server *network.Server, client *network.Client, config *define.GConfig) *Processor {
	if logger == nil {
		logger = utils.GetLogger("proxy")
	}
	return &Processor{
		server: server,
		client: client,
		config: config,
	}
}
