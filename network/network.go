/*
Package network server and client

1.暂时仅支持单处理器，可以随时按订阅模式扩展成多处理器，带类型注册处理器进而实现消息分发

2.不管主动停止Stop还是被动停止Accept error，继续接收的消息都应该记录后因为GetBind==nil而返回错误（除非登陆、注册等等）

3.OnMessage返回error请以如下格式创建，请自行校验Json数据的合法性，该数据将直接回复给客户端

	var ErrSuccess = errors.New(`{"errno":0,"errdesc":"success"}`)
*/
package network

import (
	"encoding/binary"
	"io"
	"net"
)

// Processor 处理器
type Processor interface {
	OnMessage(net.Conn, uint16, uint16, []byte) error
	OnClose(net.Conn)

	OnClientMessage(net.Conn, uint16, uint16, []byte)
	OnClientConnect(net.Conn)
}

// RecvMessage 接收消息
func RecvMessage(conn net.Conn) (uint16, uint16, []byte, error) {
	size := make([]byte, 2)

	if _, err := io.ReadFull(conn, size); err != nil {
		return 0, 0, nil, err
	}

	n := binary.BigEndian.Uint16(size)
	data := make([]byte, n)
	copy(data, size)

	if _, err := io.ReadFull(conn, data[2:]); err != nil {
		return 0, 0, nil, err
	}

	return binary.BigEndian.Uint16(data[2:]), binary.BigEndian.Uint16(data[4:]), data[6:], nil
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
