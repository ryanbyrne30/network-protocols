package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Server struct {
	host    string
	port    string
	maxConn int
	conns   int
}

type Client struct {
	conn net.Conn
}

type Config struct {
	Host    string
	Port    string
	MaxConn int
}

func New(config *Config) *Server {
	return &Server{
		host:    config.Host,
		port:    config.Port,
		maxConn: config.MaxConn,
		conns:   0,
	}
}

func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		if server.conns >= server.maxConn {
			continue
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Adding client: %s\n", conn.RemoteAddr())
		server.conns += 1
		client := &Client{
			conn: conn,
		}
		go client.handleRequest(func() {
			log.Printf("Removing client: %s\n", conn.RemoteAddr())
			server.conns -= 1
		})
	}
}

func (client *Client) handleRequest(callback func()) {
	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			client.conn.Close()
			callback()
			return
		}
		fmt.Printf("[%s] Message incoming: %s", client.conn.RemoteAddr(), string(message))
		msg := fmt.Sprintf("[%s] Message received.\n", client.conn.LocalAddr())
		client.conn.Write([]byte(msg))
	}
}
