package proxy

import (
	"sync"

	"github.com/panshiqu/framework/define"
)

// Selected 已选服务
type Selected struct {
	mutex    sync.RWMutex
	selected map[int]*define.Service
}

// Init 初始化
func (s *Selected) Init(v map[int]*define.Service) {

}

// Add 增加
func (s *Selected) Add(v *define.Service) {

}

// Del 减少
func (s *Selected) Del(v *define.Service) {

}
