package network

import (
	"net"
	"time"

	"github.com/panshiqu/framework/utils"
)

// Client 客户端
type Client struct {
	address   string
	processor Processor

	conn  net.Conn
	delay time.Duration
}

// NewClient 创建客户端
func NewClient(address string) *Client {
	return &Client{
		address: address,
	}
}

// Register 注册处理
func (c *Client) Register(processor Processor) {
	c.processor = processor
}

// Start 开始服务
func (c *Client) Start() {
	defer func() {
		utils.Sleep(&c.delay)
		go c.Start()
	}()

	var err error

	if c.conn, err = net.Dial("tcp", c.address); err != nil {
		return
	}

	c.processor.OnClientConnect(c.conn)

	for {
		mcmd, scmd, data, err := RecvMessage(c.conn)
		if err != nil {
			break
		}

		c.processor.OnClientMessage(c.conn, mcmd, scmd, data)
	}

	c.conn = nil
	c.delay = 0
}

// SendMessage 发送消息
func (c *Client) SendMessage(mcmd uint16, scmd uint16, data []byte) error {
	if c.conn != nil {
		return SendMessage(c.conn, mcmd, scmd, data)
	}

	return nil
}
