package network

import (
	"log"
	"net"
	"sync"

	"../define"
)

// Server 服务器
type Server struct {
	listener  net.Listener
	processor Processor

	mutex       sync.Mutex
	waitgroup   sync.WaitGroup
	connections map[net.Conn]interface{}
}

// NewServer 创建服务器
func NewServer(address string) *Server {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	return &Server{
		listener:    listener,
		connections: make(map[net.Conn]interface{}),
	}
}

// Register 注册处理
func (s *Server) Register(processor Processor) {
	s.processor = processor
}

// Start 开始服务
func (s *Server) Start() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.stop()
			return err
		}

		go s.handleConn(conn)
	}
}

func (s *Server) stop() {
	s.mutex.Lock()
	for conn := range s.connections {
		conn.Close()
	}
	s.connections = nil
	s.mutex.Unlock()
	s.waitgroup.Wait()
}

// Stop 停止服务
func (s *Server) Stop() {
	s.listener.Close()
}

func (s *Server) addConn(conn net.Conn) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.connections == nil {
		return false
	}
	s.connections[conn] = nil
	return true
}

func (s *Server) removeConn(conn net.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.connections != nil {
		delete(s.connections, conn)
		conn.Close()
	}
}

func (s *Server) handleConn(conn net.Conn) {
	s.waitgroup.Add(1)
	defer s.waitgroup.Done()
	if !s.addConn(conn) {
		conn.Close()
		return
	}
	defer s.removeConn(conn)

	for {
		mcmd, scmd, data, err := RecvMessage(conn)
		if err != nil {
			break
		}

		if err := s.processor.OnMessage(conn, mcmd, scmd, data); err != nil {
			me, ok := err.(*define.MyError)

			if !ok {
				me = &define.MyError{
					Errno:   define.ErrnoFailure,
					Errdesc: err.Error(),
				}
			}

			if err := SendMessage(conn, mcmd, scmd, []byte(me.Error())); err != nil {
				break
			}
		}
	}

	s.processor.OnClose(conn)
}

// SetBind 设置绑定
func (s *Server) SetBind(conn net.Conn, v interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.connections != nil {
		s.connections[conn] = v
	}
}

// GetBind 获取绑定
func (s *Server) GetBind(conn net.Conn) interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.connections != nil {
		return s.connections[conn]
	}
	return nil
}
