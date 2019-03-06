package main

import (
	"log"
	"os"
	"os/signal"

	"./db"
	"./network"
	"./utils"
	//“_” 操作引用包是无法通过包名来调用包中的导出函数，而是只是为了简单的调用其 init() 函数。
	_ "github.com/go-sql-driver/mysql"
	"./db/models"
)

func handleSignal(server *network.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	log.Println("Got signal:", s)

	models.CloseEngine()
	server.Stop()
}

func main() {
	//解析命令行参数
	args := utils.GetDBArgs()

	//读取全局配置文件
	config,err := utils.GetGConfig(args.ConfigPath)

	if err != nil {
		log.Println("ReadJSON Config", err)
		return
	}

	server := network.NewServer(config.DB.ListenIP)
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
