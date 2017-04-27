package network

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

// ErrDisconnect 断开连接
var ErrDisconnect = errors.New(`{"errno":1,"errdesc":"disconnect"}`)

// Client 客户端
type Client struct {
	stop      chan bool
	address   string
	processor Processor
	delay     time.Duration
	mutex     sync.RWMutex
	conn      net.Conn
}

// NewClient 创建客户端
func NewClient(address string) *Client {
	return &Client{
		stop:    make(chan bool),
		address: address,
	}
}

// Register 注册处理
func (c *Client) Register(processor Processor) {
	c.processor = processor
}

// Start 开始服务
func (c *Client) Start() {
	for {
		c.delay = 0

		c.mutex.Lock() // 因为重连将会重新设置conn所以增加写锁

		for {
			var err error

			if c.conn != nil { // 重连前关闭上次连接
				c.conn.Close()
			}

			log.Println("Dial", c.address, "delay", c.delay)
			if c.conn, err = net.Dial("tcp", c.address); err == nil {
				break
			}

			if c.delay != 0 {
				c.delay *= 2
			} else {
				c.delay = 5 * time.Millisecond
			}
			if max := 1 * time.Second; c.delay > max {
				c.delay = max
			}

			select {
			case <-c.stop:
				c.mutex.Unlock() // 防止goroutine正阻塞在mutex.RLock
				return
			case <-time.After(c.delay):
			}
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

		select {
		case <-c.stop:
			return
		default:
		}
	}
}

// Stop 停止服务
func (c *Client) Stop() {
	close(c.stop)
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.conn != nil { // 担心正在重连时停止conn==nil
		c.conn.Close()
	}
}

// SendMessage 发送消息
func (c *Client) SendMessage(mcmd uint16, scmd uint16, data []byte) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.conn != nil { // 担心正在重连时发送conn==nil
		return SendMessage(c.conn, mcmd, scmd, data)
	}
	return ErrDisconnect
}
