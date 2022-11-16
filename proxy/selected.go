package proxy

import (
	"net"
	"sync"

	"github.com/panshiqu/framework/define"
)

var sins Selected

// Selected 已选服务、所有服务
// 支持通过游戏编号精准重连
type Selected struct {
	mutex    sync.RWMutex
	services map[int]*define.Service
	selected map[int]*define.Service
}

// Get 获取已选代理
// 已意识到已选代理维护在这里当前无意义
// 但未来将有能力为已连上代理的客户端快速回复或推送已选代理
func (s *Selected) Get() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, v := range s.selected {
		if v.ServiceType == define.ServiceProxy {
			return v.IP
		}
	}

	return ""
}

// Dial 连接
func (s *Selected) Dial(st, gi, gt, gl int) (net.Conn, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if gi != 0 {
		for _, v := range s.services {
			if v.ServiceType == st && v.ID == gi {
				return net.Dial("tcp", v.IP)
			}
		}
	} else {
		for _, v := range s.selected {
			if v.ServiceType == st && v.GameType == gt && v.GameLevel == gl {
				return net.Dial("tcp", v.IP)
			}
		}
	}

	return nil, define.ErrNotExistService
}

// Init 初始化
func (s *Selected) Init(v map[int]*define.Service) {
	s.mutex.Lock()
	s.selected = v
	s.mutex.Unlock()
}

// Add 增加
func (s *Selected) Add(v *define.Service) {
	s.mutex.Lock()
	s.selected[v.ID] = v
	s.mutex.Unlock()
}

// Del 删除
func (s *Selected) Del(v *define.Service) {
	s.mutex.Lock()
	delete(s.selected, v.ID)
	s.mutex.Unlock()
}

// Change 改变
func (s *Selected) Change(v []*define.Service) {
	s.mutex.Lock()
	s.selected[v[0].ID] = v[0]
	delete(s.selected, v[1].ID)
	s.mutex.Unlock()
}

// InitAll 初始所有
func (s *Selected) InitAll(v map[int]*define.Service) {
	s.mutex.Lock()
	s.services = v
	s.mutex.Unlock()
}

// Incr 增加
func (s *Selected) Incr(v *define.Service) {
	s.mutex.Lock()
	s.services[v.ID] = v
	s.mutex.Unlock()
}

// Decr 删除
func (s *Selected) Decr(v *define.Service) {
	s.mutex.Lock()
	delete(s.services, v.ID)
	s.mutex.Unlock()
}
