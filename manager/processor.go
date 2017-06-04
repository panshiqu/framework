package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
	"github.com/panshiqu/framework/utils"
)

// Processor 处理器
type Processor struct {
	mutex    sync.Mutex              // 锁
	server   *network.Server         // 服务器
	services map[int]*define.Service // 服务表（所有已开启的服务）
	selected map[int]*define.Service // 已选表（负载均衡策略后选择的服务）
}

// OnMessage 收到消息
func (p *Processor) OnMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	log.Println("OnMessage", mcmd, scmd, string(data))

	switch mcmd {
	case define.ManagerCommon:
		return p.OnMainCommon(conn, scmd, data)
	}

	return define.NewError(fmt.Sprint("unknown main cmd ", mcmd))
}

// OnMainCommon 通用主命令
func (p *Processor) OnMainCommon(conn net.Conn, scmd uint16, data []byte) error {
	switch scmd {
	case define.ManagerRegisterService:
		return p.OnSubRegisterService(conn, data)
	case define.ManagerUpdateCount:
		return p.OnSubUpdateCount(conn, data)
	case define.ManagerOpenService:
		return p.OnSubOpenService(conn, data)
	case define.ManagerShutService:
		return p.OnSubShutService(conn, data)
	}

	return define.NewError(fmt.Sprint("unknown sub cmd ", scmd))
}

// OnSubRegisterService 注册服务子命令
func (p *Processor) OnSubRegisterService(conn net.Conn, data []byte) error {
	service := &define.Service{}

	if err := json.Unmarshal(data, service); err != nil {
		return define.NewError(err.Error())
	}

	// 加锁
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 重复注册服务
	if _, ok := p.services[service.ID]; ok {
		return define.NewError("repeat register service")
	}

	// 设置网络连接
	service.Conn = conn

	// 服务表增加
	p.services[service.ID] = service

	// 是否存在类似已选
	if !p.isExistSimilarSelected(service) {
		// 增加已选服务
		p.addSelectedService(service)
	}

	// 仅通知代理
	if service.ServiceType != define.ServiceProxy {
		return nil
	}

	// 通知已选服务
	if err := network.SendJSONMessage(conn, define.ManagerCommon, define.ManagerNotifyCurService, p.selected); err != nil {
		log.Println("OnSubRegisterService SendJSONMessage", err)
	}

	return nil
}

// OnSubUnRegisterService 注销服务子命令
func (p *Processor) OnSubUnRegisterService(conn net.Conn, data []byte) error {
	// 加锁
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, v := range p.services {
		if v.Conn == conn {
			// 服务表删除
			delete(p.services, v.ID)

			// 改变已选服务
			p.changeSelectedService(v.ID)

			break
		}
	}

	return nil
}

// OnSubUpdateCount 更新计数子命令
func (p *Processor) OnSubUpdateCount(conn net.Conn, data []byte) error {
	updateCount := &define.Service{}

	if err := json.Unmarshal(data, updateCount); err != nil {
		return define.NewError(err.Error())
	}

	// 加锁
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 获取服务
	service, ok := p.services[updateCount.ID]
	if !ok {
		return define.NewError("not exist service")
	}

	// 更新计数
	service.Count = updateCount.Count

	// 计数小于对应容量
	if service.Count < p.getServiceCapacity(service.ServiceType) {
		return nil
	}

	// 改变已选服务
	p.changeSelectedService(updateCount.ID)

	return nil
}

// OnSubOpenService 开启服务
func (p *Processor) OnSubOpenService(conn net.Conn, data []byte) error {
	openService := &define.Service{}

	if err := json.Unmarshal(data, openService); err != nil {
		return define.NewError(err.Error())
	}

	// 加锁
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 获取服务
	service, ok := p.services[openService.ID]
	if !ok {
		return define.NewError("not exist service")
	}

	// 服务已开启
	if service.IsServe {
		return define.NewError("service already open")
	}

	// 开启服务
	service.IsServe = true

	// 是否存在类似已选
	if !p.isExistSimilarSelected(service) {
		// 增加已选服务
		p.addSelectedService(service)
	}

	return nil
}

// OnSubShutService 关闭服务
func (p *Processor) OnSubShutService(conn net.Conn, data []byte) error {
	shutService := &define.Service{}

	if err := json.Unmarshal(data, shutService); err != nil {
		return define.NewError(err.Error())
	}

	// 加锁
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 获取服务
	service, ok := p.services[shutService.ID]
	if !ok {
		return define.NewError("not exist service")
	}

	// 服务已关闭
	if !service.IsServe {
		return define.NewError("service already shut")
	}

	// 关闭服务
	service.IsServe = false

	// 改变已选服务
	p.changeSelectedService(shutService.ID)

	return nil
}

// OnClose 连接关闭
func (p *Processor) OnClose(conn net.Conn) {
	defer utils.Trace("Processor OnClose")()

	p.OnSubUnRegisterService(conn, nil)
}

// OnClientMessage 客户端收到消息
func (p *Processor) OnClientMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) {
	// nothing to do
}

// OnClientConnect 客户端连接成功
func (p *Processor) OnClientConnect(conn net.Conn) {
	// nothing to do
}

// isSimilarService 是否类似服务
func (p *Processor) isSimilarService(l, r *define.Service) bool {
	return l.ServiceType == r.ServiceType && l.GameType == r.GameType && l.GameLevel == r.GameLevel
}

// isExistSimilarSelected 是否存在类似已选
func (p *Processor) isExistSimilarSelected(service *define.Service) bool {
	defer utils.Trace("Processor isExistSimilarSelected", service.ID)()

	for _, v := range p.selected {
		if p.isSimilarService(v, service) {
			return true
		}
	}

	return false
}

// getSimilarService 获取类似服务
func (p *Processor) getSimilarService(service *define.Service) *define.Service {
	defer utils.Trace("Processor getSimilarService", service.ID)()

	var min, max *define.Service
	for _, v := range p.services {
		if !p.isSimilarService(v, service) ||
			!v.IsServe {
			continue
		}

		// 总是记录计数最少的服务
		if min == nil || v.Count < min.Count {
			min = v
		}

		// 只关心游戏服务且计数小于游戏容量
		if service.ServiceType != define.ServiceGame ||
			v.Count > define.CapacityGame {
			continue
		}

		// 总是记录计数较多的游戏服务（尽可能在相同游戏服务中玩，便于快速组桌开始游戏）
		if max == nil || v.Count > max.Count {
			max = v
		}
	}

	if service.ServiceType == define.ServiceGame && max != nil {
		return max
	}

	return min
}

// addSelectedService 增加已选服务
func (p *Processor) addSelectedService(service *define.Service) {
	defer utils.Trace("Processor addSelectedService", service.ID)()

	// 已选表增加
	p.selected[service.ID] = service

	// 通知增加已选服务
	p.notifySelectedService(define.ManagerNotifyAddService, service)
}

// delSelectedService 删除已选服务
func (p *Processor) delSelectedService(service *define.Service) {
	defer utils.Trace("Processor delSelectedService", service.ID)()

	// 已选表删除
	delete(p.selected, service.ID)

	// 通知删除已选服务
	p.notifySelectedService(define.ManagerNotifyDelService, service)
}

// notifySelectedService 通知已选服务
func (p *Processor) notifySelectedService(scmd uint16, service *define.Service) {
	defer utils.Trace("Processor notifySelectedService", scmd, service.ID)()

	data, err := json.Marshal(service)
	if err != nil {
		log.Println("notifySelectedService Marshal", err)
		return
	}

	for _, v := range p.services {
		// 不通知自己
		if v.ID == service.ID {
			continue
		}

		// 仅通知代理
		if v.ServiceType != define.ServiceProxy {
			continue
		}

		// 通知已选服务
		if err := network.SendMessage(v.Conn, define.ManagerCommon, scmd, data); err != nil {
			log.Println("notifySelectedService SendMessage", err)
		}
	}
}

// changeSelectedService 改变已选服务
func (p *Processor) changeSelectedService(id int) {
	defer utils.Trace("Processor changeSelectedService", id)()

	// 是否存在
	oldService, ok := p.selected[id]
	if !ok {
		return
	}

	// 获取类似服务
	newService := p.getSimilarService(oldService)
	if oldService == newService {
		return
	}

	// 删除已选服务
	p.delSelectedService(oldService)

	if newService == nil {
		return
	}

	// 增加已选服务
	p.addSelectedService(newService)
}

// getServiceCapacity 获取服务容量
func (p *Processor) getServiceCapacity(tp int) int {
	switch tp {
	case define.ServiceProxy:
		return define.CapacityProxy
	case define.ServiceLogin:
		return define.CapacityLogin
	case define.ServiceGame:
		return define.CapacityGame
	}

	return 0
}

// NewProcessor 创建处理器
func NewProcessor(server *network.Server) *Processor {
	return &Processor{
		server:   server,
		services: make(map[int]*define.Service),
		selected: make(map[int]*define.Service),
	}
}

// Monitor 监视器
func (p *Processor) Monitor(w http.ResponseWriter, r *http.Request) {
	// 加锁
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 打印服务表
	fmt.Fprintln(w, "services:")
	for _, v := range p.services {
		fmt.Fprintln(w, v)
	}

	// 打印已选表
	fmt.Fprintln(w, "selected:")
	for _, v := range p.selected {
		fmt.Fprintln(w, v)
	}
}
