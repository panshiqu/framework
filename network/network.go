/*
Package network server and client

1.暂时仅支持单处理器，可以随时按订阅模式扩展成多处理器，带类型注册处理器进而实现消息分发

2.不管主动停止Stop还是被动停止Accept error，继续接收的消息都应该记录后因为GetBind==nil而返回错误（除非登陆、注册等等）

3.OnMessage返回值可以扩展成any类型，通过断言error, MyError, others进而区别对待，将可以实现快捷回复消息（现仅在db中实现）

4.RPC发送接收必须匹配才能实现同步调用，则必然需要特殊处理return nil，但若nil等同于ErrSuccess则需要与真正的回复互斥，所以捎带实现第3点
*/
package network

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/utils"
)

const MessageHead = 0x5277

// Processor 处理器
type Processor interface {
	OnMessage(net.Conn, uint16, uint16, []byte) error
	OnClose(net.Conn)

	OnClientMessage(net.Conn, uint16, uint16, []byte)
	OnClientConnect(net.Conn)
}

// RecvMessage 接收消息
func RecvMessage(conn net.Conn) (uint16, uint16, []byte, error) {
	size := make([]byte, 4)

	if _, err := io.ReadFull(conn, size); err != nil {
		return 0, 0, nil, utils.Wrap(err)
	}

	if binary.BigEndian.Uint16(size) != MessageHead {
		return 0, 0, nil, utils.Wrap(define.ErrMessageHead)
	}

	n := binary.BigEndian.Uint16(size[2:])
	if n > define.LengthLimit {
		return 0, 0, nil, utils.Wrap(define.ErrLengthLimit)
	}

	message := make([]byte, n)
	copy(message, size)

	if _, err := io.ReadFull(conn, message[4:]); err != nil {
		return 0, 0, nil, utils.Wrap(err)
	}

	return binary.BigEndian.Uint16(message[4:]), binary.BigEndian.Uint16(message[6:]), message[8:], nil
}

// SendMessage 发送消息
func SendMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	size := len(data) + 8
	if size > define.LengthLimit {
		return utils.Wrap(define.ErrLengthLimit)
	}

	message := make([]byte, size)
	binary.BigEndian.PutUint16(message, MessageHead)
	binary.BigEndian.PutUint16(message[2:], uint16(size))
	binary.BigEndian.PutUint16(message[4:], mcmd)
	binary.BigEndian.PutUint16(message[6:], scmd)
	copy(message[8:], data)

	return utils.Wrap(utils.Error(conn.Write(message)))
}

// SendJSONMessage 发送消息
func SendJSONMessage(conn net.Conn, mcmd uint16, scmd uint16, js any) error {
	data, err := json.Marshal(js)
	if err != nil {
		return utils.Wrap(err)
	}

	return utils.Wrap(SendMessage(conn, mcmd, scmd, data))
}
