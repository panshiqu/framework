/*
Package network .
1.暂时仅支持单处理器，可以随时按订阅模式扩展成多处理器，带类型注册处理器进而实现消息分发
2.不管主动停止Stop还是被动停止Accept error，继续接收的消息都应该记录后因为GetBind==nil而返回错误（除非登陆、注册等等）
*/
package network

import (
	"encoding/binary"
	"log"
	"net"

	"io"
	"sync"
)

// Processor 处理器
type Processor interface {
	OnMessage(net.Conn, uint16, uint16, []byte)
}

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
		size := make([]byte, 2)

		if _, err := io.ReadFull(conn, size); err != nil {
			break
		}

		n := binary.BigEndian.Uint16(size)
		data := make([]byte, n)
		copy(data, size)

		if _, err := io.ReadFull(conn, data[2:]); err != nil {
			break
		}

		s.processor.OnMessage(conn,
			binary.BigEndian.Uint16(data[2:]),
			binary.BigEndian.Uint16(data[4:]),
			data[6:])
	}
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

// SendMessage 发送消息
func SendMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	size := len(data) + 6
	message := make([]byte, size)
	binary.BigEndian.PutUint16(message, uint16(size))
	binary.BigEndian.PutUint16(message[2:], mcmd)
	binary.BigEndian.PutUint16(message[4:], scmd)
	copy(message[6:], data)

	if _, err := conn.Write(message); err != nil {
		return err
	}

	return nil
}
