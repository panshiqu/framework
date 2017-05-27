package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/panshiqu/framework/game"
	"github.com/panshiqu/framework/network"
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
	server := network.NewServer("127.0.0.1:8082")
	client := network.NewClient("127.0.0.1:8080")
	processor := game.NewProcessor(server, client)

	server.Register(processor)
	client.Register(processor)

	go handleSignal(server, client)
	go client.Start()

	if err := server.Start(); err != nil {
		log.Println("Start", err)
	}
}
