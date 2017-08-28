package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
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
	if err := utils.ReadJSON("./config/game.json", &define.CG); err != nil {
		log.Println("ReadJSON ConfigGame", err)
		return
	}

	server := network.NewServer(define.CG.ListenIP)
	client := network.NewClient(define.CG.DialIP)
	processor := game.NewProcessor(server, client)

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
