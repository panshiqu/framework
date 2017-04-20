package network

import (
	"log"
	"net"

	"bufio"
	"encoding/binary"
	"sync"
)

// Processor 处理器
type Processor interface {
	OnMessage(net.Conn, *Message)
}

// Server 服务器
type Server struct {
	listener  net.Listener
	processor Processor

	mutex       sync.Mutex
	waitgroup   sync.WaitGroup
	connections map[net.Conn]bool
}

// NewServer 创建服务器
func NewServer(address string, processor Processor) *Server {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	return &Server{
		listener:    listener,
		processor:   processor,
		connections: make(map[net.Conn]bool),
	}
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
	for conn := range s.connections {
		conn.Close()
	}

	s.waitgroup.Wait()
}

// Stop 停止服务
func (s *Server) Stop() {
	s.listener.Close()
}

func (s *Server) addConn(conn net.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.connections[conn] = true
}

func (s *Server) removeConn(conn net.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.connections, conn)
}

func (s *Server) handleConn(conn net.Conn) {
	s.addConn(conn)
	defer s.removeConn(conn)

	s.waitgroup.Add(1)
	defer s.waitgroup.Done()

	rd := bufio.NewReader(conn)

	for {
		buf, err := rd.Peek(2)
		if err != nil {
			break
		}

		size := int(binary.BigEndian.Uint16(buf))
		if _, err = rd.Peek(size); err != nil {
			break
		}

		message := make([]byte, size)
		if _, err = rd.Read(message); err != nil {
			break
		}

		s.processor.OnMessage(conn, NewRecvMessage(message))
	}
}
