package main

import (
	"flag"
	"api"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()

	api.NewServer()
	api.InitApi()

	api.StartServer()

	api.AuthorizedAtHeyService()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	api.StopServer()
}
