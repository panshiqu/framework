package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
)

// Processor 处理器
type Processor struct {
	client *network.Client
}

// OnMessage 收到消息
func (p *Processor) OnMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) error {
	return nil
}

// OnClose 连接关闭
func (p *Processor) OnClose(conn net.Conn) {

}

// OnClientMessage 客户端收到消息
func (p *Processor) OnClientMessage(conn net.Conn, mcmd uint16, scmd uint16, data []byte) {

}

// OnClientConnect 客户端连接成功
func (p *Processor) OnClientConnect(conn net.Conn) {
	// 构造服务
	fastRegister := &define.FastRegister{
		Account:  "panshiqu",
		Password: "111111",
		Machine:  "panshiqu",
		Name:     "panshiqu",
		Icon:     0,
		Gender:   0,
	}

	// 发送快速注册消息
	if err := p.client.SendJSONMessage(define.LoginCommon, define.LoginFastRegister, fastRegister); err != nil {
		log.Println("OnClientConnect SendJSONMessage", err)
		return
	}

	log.Println("OnClientConnect", fastRegister)
}

func handleSignal(client *network.Client) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	log.Println("Got signal:", s)

	client.Stop()
}

func main() {
	client := network.NewClient("127.0.0.1:8081")
	client.Register(&Processor{client: client})
	go handleSignal(client)
	client.Start()
}
