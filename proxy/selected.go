package proxy

import (
	"net"
	"sync"

	"../define"
)

var sins Selected

// Selected 已选服务
type Selected struct {
	mutex    sync.RWMutex
	selected map[int]*define.Service
}

// Dial 连接
func (s *Selected) Dial(st, gt, gl int) (net.Conn, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, v := range s.selected {
		if v.ServiceType == st && v.GameType == gt && v.GameLevel == gl {
			return net.Dial("tcp", v.IP)
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

// Del 减少
func (s *Selected) Del(v *define.Service) {
	s.mutex.Lock()
	delete(s.selected, v.ID)
	s.mutex.Unlock()
}
