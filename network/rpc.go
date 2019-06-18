package network

import (
	"encoding/json"
	"log"
	"net"
	"sync"

	"../define"
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
	for {
		var conn net.Conn
		if conn, err = r.get(); err != nil {
			return
		}

		if err = SendMessage(conn, mcmd, scmd, data); err != nil {
			log.Println("Call SendMessage", err)
			conn.Close()
			continue
		}

		if _, _, dt, err = RecvMessage(conn); err != nil {
			log.Println("Call RecvMessage", err)
			conn.Close()
			continue
		}

		r.put(conn)
		return
	}
}

// JSONCall 调用
func (r *RPC) JSONCall(mcmd uint16, scmd uint16, in interface{}, out interface{}) (err error) {
	var indata, outdata []byte

	if indata, err = json.Marshal(in); err != nil {
		return
	}

	if outdata, err = r.Call(mcmd, scmd, indata); err != nil {
		return
	}

	if err = define.CheckError(outdata); err != nil {
		return
	}

	if out == nil {
		return
	}

	return json.Unmarshal(outdata, out)
}


