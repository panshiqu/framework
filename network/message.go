package network

import (
	"encoding/binary"
)

// Message 消息
type Message struct {
	mcmd uint16 // 主命令
	scmd uint16 // 子命令
	body []byte // 消息体
}

// Mcmd 主命令
func (m *Message) Mcmd() uint16 {
	return m.mcmd
}

// Scmd 子命令
func (m *Message) Scmd() uint16 {
	return m.scmd
}

// Body 消息体
func (m *Message) Body() []byte {
	return m.body
}

// NewRecvMessage 创建接收消息
func NewRecvMessage(message []byte) *Message {
	return &Message{
		mcmd: binary.BigEndian.Uint16(message[2:]),
		scmd: binary.BigEndian.Uint16(message[4:]),
		body: message[6:],
	}
}

// NewSendMessage 创建发送消息
func NewSendMessage(mcmd uint16, scmd uint16, body []byte) []byte {
	size := len(body) + 6
	message := make([]byte, size)
	binary.BigEndian.PutUint16(message, uint16(size))
	binary.BigEndian.PutUint16(message[2:], mcmd)
	binary.BigEndian.PutUint16(message[4:], scmd)
	copy(message[6:], body)
	return message
}
