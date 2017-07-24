package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/game"
	"github.com/panshiqu/framework/network"
	"github.com/panshiqu/framework/utils"
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
	config := &define.ConfigGame{}
	if err := utils.ReadJSON("./config/game.json", config); err != nil {
		log.Println("ReadJSON ConfigGame", err)
		return
	}

	server := network.NewServer(config.ListenIP)
	client := network.NewClient(config.DialIP)
	processor := game.NewProcessor(server, client, config)

	server.Register(processor)
	client.Register(processor)

	go handleSignal(server, client)
	go client.Start()

	go func() {
		http.HandleFunc("/", processor.Monitor)
		log.Println(http.ListenAndServe(config.PprofIP, nil))
	}()

	if err := server.Start(); err != nil {
		log.Println("Start", err)
	}
}
