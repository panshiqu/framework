package network

import (
	"net"
	"sync"
)

// RPC 远程过程调用
// 暂时没有限制最大连接数
type RPC struct {
	address string
	mutex   sync.Mutex
	conns   []net.Conn
}

// NewRPC 创建RPC
func NewRPC(address string) *RPC {
	return &RPC{
		address: address,
		conns:   make([]net.Conn, 0, 8),
	}
}

func (r *RPC) get() (conn net.Conn, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if size := len(r.conns); size != 0 {
		conn, r.conns = r.conns[size-1], r.conns[:size-1]
		return
	}

	return net.Dial("tcp", r.address)
}

func (r *RPC) put(conn net.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.conns = append(r.conns, conn)
}

// Call 调用
func (r *RPC) Call(mcmd uint16, scmd uint16, data []byte) (dt []byte, err error) {
	conn, err := r.get()
	if err != nil {
		return
	}

	if err = SendMessage(conn, mcmd, scmd, data); err != nil {
		return
	}

	_, _, dt, err = RecvMessage(conn)
	if err != nil {
		return
	}

	r.put(conn)
	return
}
