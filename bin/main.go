package main

import (
	"cache/bin/config"
	"cache/internal/controller"
	"cache/internal/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := config.GetConfig()
	s := service.New()
	c := controller.New(s, config.Port)
	go c.Init()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
