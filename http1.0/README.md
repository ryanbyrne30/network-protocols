# HTTP/1.0 Server

Simple [HTTP/1.0](https://www.w3.org/Protocols/HTTP/1.0/spec.html) server.

## Getting Started

```sh
go run main.go
```

## Sending Requests

You can send HTTP requests however you wish (curl, Postman, netcat, etc).

Using Netcat

```sh
nc 127.0.0.1 5555 < request.txt
```

## Configuration

```sh
go run main.go -h
```