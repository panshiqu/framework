package utils

import (
	"sync"
	"time"

	"github.com/panshiqu/framework/network"
)

type event struct {
	timer     *time.Timer
	ticker    *time.Ticker
	endtime   time.Time
	duration  time.Duration
	parameter interface{}
}

func (e *event) expire() bool {
	if e.timer != nil {
		select {
		case <-e.timer.C:
			return true
		default:
		}
	}

	if e.ticker != nil {
		select {
		case <-e.ticker.C:
			return true
		default:
		}
	}

	return false
}

// Schedule 时间表
// 命名也只是不想与标准库Timer重名而已
// 定时时间精确到秒，不精确管理goroutine的退出
// 主要逻辑只不过是对标准库Timer、Ticker封装管理而已
type Schedule struct {
	mutex     sync.Mutex
	events    map[int]*event
	processor network.Processor
}

// NewSchedule 创建时间表
func NewSchedule(processor network.Processor) *Schedule {
	return &Schedule{
		events:    make(map[int]*event),
		processor: processor,
	}
}

// Start 开始
func (s *Schedule) Start() {
	for {
		time.Sleep(time.Second)

		s.mutex.Lock()

		for k, v := range s.events {
			if v.expire() {
				if v.timer != nil {
					delete(s.events, k)
				}
				if v.ticker != nil {
					v.endtime = v.endtime.Add(v.duration)
				}
				go s.processor.OnTimer(k, v.parameter)
			}
		}

		s.mutex.Unlock()
	}
}

// Add 添加
func (s *Schedule) Add(id int, duration time.Duration, parameter interface{}, persistence bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	ev, ok := s.events[id]
	if ok {
		if ev.timer == nil {
			return
		}
		if !ev.timer.Stop() {
			<-ev.timer.C // 必须这样做，不然Reset之后有可能立刻到期
		}
		ev.timer.Reset(duration)
	} else {
		ev = new(event)
		s.events[id] = ev
		if !persistence {
			ev.timer = time.NewTimer(duration)
		} else {
			ev.ticker = time.NewTicker(duration)
		}
	}

	ev.duration = duration
	ev.parameter = parameter
	ev.endtime = time.Now().Add(duration)
}

// Remove 移除
func (s *Schedule) Remove(id int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if ev, ok := s.events[id]; ok {
		if ev.timer != nil && !ev.timer.Stop() {
			<-ev.timer.C // 其实不必大费周章，因为delete(s.events, id)
		}
		if ev.ticker != nil {
			ev.ticker.Stop()
		}
		delete(s.events, id)
	}
}

// Surplus 剩余
func (s *Schedule) Surplus(id int) (duration time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if ev, ok := s.events[id]; ok {
		duration = ev.endtime.Sub(time.Now())
	}

	if duration < 0 {
		duration = 0
	}

	return
}
