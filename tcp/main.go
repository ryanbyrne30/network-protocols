package main

import (
	"flag"

	"github.com/ryanbyrne30/network-protocols/tcp/server"
)

var host string
var port string
var conns int

func main() {
	flag.StringVar(&host, "h", "127.0.0.1", "Host to run the server on")
	flag.StringVar(&port, "p", "5555", "Port to run the server on")
	flag.IntVar(&conns, "c", 2, "Max connections server can process at a time")
	flag.Parse()

	server := server.New(&server.Config{
		Host:    host,
		Port:    port,
		MaxConn: conns,
	})

	server.Run()
}
