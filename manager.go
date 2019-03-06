package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"./manager"
	"./network"
	"./utils"
)

func handleSignal(server *network.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	log.Println("Got signal:", s)

	server.Stop()
}

func main() {
	//读取命令行参数
	args := utils.GetLoginArgs()
	fmt.Println(args.ConfigPath)

	//读取全局配置文件
	config,err := utils.GetGConfig(args.ConfigPath)

	if  err != nil {
		log.Println("ReadJSON Config", err)
		return
	}

	server := network.NewServer(config.Manager.ListenIP)
	processor := manager.NewProcessor(server)
	server.Register(processor)

	go handleSignal(server)

	go func() {
		http.HandleFunc("/", processor.Monitor)
		log.Println(http.ListenAndServe(config.Manager.PprofIP, nil))
	}()

	if err := server.Start(); err != nil {
		log.Println("Start", err)
	}
}
