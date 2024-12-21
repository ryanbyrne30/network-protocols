package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Server struct {
	host string
	port string
}

type Client struct {
	conn net.Conn
}

type Config struct {
	Host string
	Port string
}

func New(config *Config) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn: conn,
		}
		go client.handleRequest()
	}
}

func (client *Client) handleRequest() {
	log.Printf("[%s] Connection established\n", client.conn.RemoteAddr())
	req, err := readRequest(client.conn)
	if errors.Is(err, io.EOF) {
		client.closeConnection()
		return
	}

	log.Printf("[%s] %s %s %v '%s'", client.conn.RemoteAddr(), req.Method, req.Endpoint, req.Headers, strings.ReplaceAll(req.Body, "\n", "\\n"))

	res := NewResponse()
	if err != nil {
		res.StatusCode = 500
		res.Body = err.Error()
	} else {
		res.Body = "Message received"
	}

	client.conn.Write(res.ToBytes())
	client.closeConnection()
}

func (client *Client) closeConnection() {
	log.Printf("Goodbye %s", client.conn.RemoteAddr())
	client.conn.Close()
}
