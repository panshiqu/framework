package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"

	"./utils"
	"./define"
	"./game"
	"./network"
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
		log.Println("ReadJSON Config", err)
		return
	}
	define.CG = config.Game

	server := network.NewServer(config.Game.ListenIP)
	client := network.NewClient(config.Manager.ListenIP)
	processor := game.NewProcessor(server, client, config)

	server.Register(processor)
	client.Register(processor)

	go handleSignal(server, client)
	go client.Start()

	go func() {
		http.HandleFunc("/", processor.Monitor)
		log.Println(http.ListenAndServe(define.CG.PprofIP, nil))
	}()

	if err := server.Start(); err != nil {
		log.Println("Start", err)
	}
}
