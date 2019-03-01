package main

import (
	"log"
	"os"
	"os/signal"

	"./define"
	"./login"
	"./network"
	"./utils"
)

//信号处理函数
func handleSignal(server *network.Server, client *network.Client) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	log.Println("Got signal:", s)

	server.Stop()
	client.Stop()
}

func main() {
	//读取login配置文件
	config := &define.ConfigLogin{}
	if err := utils.ReadJSON("./config/login.json", config); err != nil {
		log.Println("ReadJSON ConfigLogin", err)
		return
	}

	server := network.NewServer(config.ListenIP)
	client := network.NewClient(config.DialIP)
	processor := login.NewProcessor(server, client, config)

	server.Register(processor)
	client.Register(processor)

	go handleSignal(server, client)
	go client.Start()

	if err := server.Start(); err != nil {
		log.Println("Start", err)
	}
}
