package network

import (
	"net"
	"sync"
	"time"

	"github.com/panshiqu/framework/utils"
)

// Client 客户端
type Client struct {
	stop      bool
	address   string
	processor Processor
	delay     time.Duration
	mutex     sync.RWMutex
	conn      net.Conn
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
	for !c.stop {
		c.mutex.Lock() // 因为重连将会重新设置conn所以增加写锁

		for {
			var err error

			if c.conn != nil {
				c.conn.Close() // 重连前关闭上次连接
			}

			if c.conn, err = net.Dial("tcp", c.address); err == nil {
				break
			}

			utils.Sleep(&c.delay)
		}

		c.mutex.Unlock()

		c.mutex.RLock()

		c.processor.OnClientConnect(c.conn)

		for {
			mcmd, scmd, data, err := RecvMessage(c.conn)
			if err != nil {
				break
			}

			c.processor.OnClientMessage(c.conn, mcmd, scmd, data)
		}

		c.mutex.RUnlock()

		c.delay = 0
	}
}

// Stop 停止服务
func (c *Client) Stop() {
	c.mutex.RLock()
	c.stop = true
	c.conn.Close()
	c.mutex.RUnlock()
}

// SendMessage 发送消息
func (c *Client) SendMessage(mcmd uint16, scmd uint16, data []byte) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return SendMessage(c.conn, mcmd, scmd, data)
}
