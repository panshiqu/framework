package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/manager"
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
	config := &define.ConfigManager{}
	if err := utils.ReadJSON("./config/manager.json", config); err != nil {
		log.Println("ReadJSON ConfigManager", err)
		return
	}

	server := network.NewServer(config.ListenIP)
	processor := manager.NewProcessor(server)
	server.Register(processor)
	go handleSignal(server)

	go func() {
		http.HandleFunc("/", processor.Monitor)
		log.Println(http.ListenAndServe(config.PprofIP, nil))
	}()

	if err := server.Start(); err != nil {
		log.Println("Start", err)
	}
}
