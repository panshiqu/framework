package main

import (
	"log"
	"net"
	"os"
	"os/signal"

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

}

func handleSignal(client *network.Client) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	log.Println("Got signal:", s)

	client.Stop()
}

func main() {
	client := network.NewClient("127.0.0.1:8080")
	client.Register(&Processor{client: client})
	go handleSignal(client)
	client.Start()
}
