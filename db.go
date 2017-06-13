package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/panshiqu/framework/db"
	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
	"github.com/panshiqu/framework/utils"
)

func handleSignal(server *network.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	log.Println("Got signal:", s)

	server.Stop()
}

func main() {
	config := &define.ConfigDB{}
	if err := utils.ReadJSON("./config/db.json", config); err != nil {
		log.Println("ReadJSON ConfigDB", err)
		return
	}

	server := network.NewServer(config.ListenIP)
	processor := db.NewProcessor(server)
	server.Register(processor)
	go handleSignal(server)

	if err := server.Start(); err != nil {
		log.Println("Start", err)
	}
}
