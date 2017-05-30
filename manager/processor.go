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
	case define.ManagerServiceUpdateCount:
		return p.OnSubServiceUpdateCount(conn, data)
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

	// 不存在类似服务
	if !p.isExistSimilar(service) {
		// 已选表增加
		p.selected[service.ID] = service

		// 广播已选服务
	}

	return nil
}

// isExistSimilar 是否存在类似服务
func (p *Processor) isExistSimilar(service *define.Service) bool {
	for _, v := range p.selected {
		if v.ServiceType == service.ServiceType &&
			v.GameType == service.GameType &&
			v.GameLevel == service.GameLevel {
			return true
		}
	}

	return false
}

// getSimilarService 获取类似服务
func (p *Processor) getSimilarService(service *define.Service) *define.Service {
	var min, max *define.Service
	for _, v := range p.services {
		if v.ServiceType != service.ServiceType ||
			v.GameType != service.GameType ||
			v.GameLevel != service.GameLevel {
			continue
		}

		// 服务已关闭
		if !v.IsServe {
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

// OnSubUnRegisterService 注销服务子命令
func (p *Processor) OnSubUnRegisterService(conn net.Conn, data []byte) error {
	// 加锁
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, v := range p.services {
		if v.Conn == conn {
			// 服务表删除
			delete(p.services, v.ID)

			// 注销服务存在已选表中
			if oldService, ok := p.selected[v.ID]; ok {
				// 已选表删除
				delete(p.selected, v.ID)

				// 获取类似服务成功
				if newService := p.getSimilarService(oldService); newService != nil {
					// 已选表增加
					p.selected[newService.ID] = newService

					// 广播已选服务
				}
			}

			break
		}
	}

	return nil
}

// OnSubServiceUpdateCount 服务更新计数子命令
func (p *Processor) OnSubServiceUpdateCount(conn net.Conn, data []byte) error {
	return nil
}

// OnClose 连接关闭
func (p *Processor) OnClose(conn net.Conn) {
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
