package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	//“_” 操作引用包是无法通过包名来调用包中的导出函数，而是只是为了简单的调用其 init() 函数。
	_ "github.com/go-sql-driver/mysql"
	"./db"
	"./define"
	"./network"
	"./utils"
)

func handleSignal(server *network.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	log.Println("Got signal:", s)

	db.GameEngine.Close()
	db.LogEngine.Close()
	server.Stop()
}

func main() {
	//解析命令行参数
	args := utils.GetDBArgs()
	fmt.Println(args)

	config := &define.ConfigDB{}
	if err := utils.ReadJSON(args.ConfigPath, config); err != nil {
		log.Println("ReadJSON ConfigDB", err)
		return
	}

	server := network.NewServer(config.ListenIP)
	processor := db.NewProcessor(server, config)

	if processor == nil {
		return
	}

	server.Register(processor)
	go handleSignal(server)

	if err := server.Start(); err != nil {
		log.Println("Start", err)
	}
}
