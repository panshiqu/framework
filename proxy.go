package main

import (
	"log"
	"os"
	"os/signal"

	"./network"
	"./proxy"
	"./utils"
)

func handleSignal(server *network.Server, client *network.Client) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	log.Println("Got signal:", s)

	server.Stop()
	client.Stop()
}

func main() {
	//读取命令行参数
	args := utils.GetLoginArgs()

	//读取全局配置文件
	config,err := utils.GetGConfig(args.ConfigPath)

	if err != nil {
		log.Println("ReadJSON ConfigProxy", err)
		return
	}

	server := network.NewServer(config.Proxy.ListenIP)
	client := network.NewClient(config.Manager.ListenIP)
	processor := proxy.NewProcessor(server, client, config)

	server.Register(processor)
	client.Register(processor)

	go handleSignal(server, client)
	go client.Start()

	if err := server.Start(); err != nil {
		log.Println("Start", err)
	}
}
