package main

import (
	"flag"
	"api"
)

func main() {
	flag.Parse()

	api.NewServer()
	api.InitApi()
	api.StartServer()
}
