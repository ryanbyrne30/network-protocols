package main

import (
	"flag"

	"github.com/ryanbyrne30/network-protocols/http1.0/server"
)

var host string
var port string

func main() {
	flag.StringVar(&host, "h", "127.0.0.1", "Host to run the server on")
	flag.StringVar(&port, "p", "5555", "Port to run the server on")
	flag.Parse()

	server := server.New(&server.Config{
		Host: host,
		Port: port,
	})

	server.Run()
}
