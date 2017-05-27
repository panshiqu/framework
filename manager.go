package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/panshiqu/framework/manager"
	"github.com/panshiqu/framework/network"
)

func handleSignal(server *network.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	log.Println("Got signal:", s)

	server.Stop()
}

func main() {
	server := network.NewServer("127.0.0.1:8080")
	processor := manager.NewProcessor(server)
	server.Register(processor)
	go handleSignal(server)

	go func() {
		http.HandleFunc("/", processor.Monitor)
		log.Println(http.ListenAndServe("127.0.0.1:9090", nil))
	}()

	if err := server.Start(); err != nil {
		log.Println("Start", err)
	}
}
