package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type HttpRequest struct {
	Method   string
	Endpoint string
	Headers  map[string]string
	Body     string
}

func readRequest(conn net.Conn) (*HttpRequest, error) {
	request := &HttpRequest{
		Method:   "",
		Endpoint: "",
		Headers:  map[string]string{},
		Body:     "",
	}

	reader := bufio.NewReader(conn)
	bodyLength := 0

	// parse head
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return request, err
		}

		isFirstLine := request.Method == ""
		isBodySeperator := strings.TrimSpace(line) == ""

		if isFirstLine {
			method, endpoint, err := parseRequestLine(line)
			request.Method = method
			request.Endpoint = endpoint
			if err != nil {
				return request, err
			}
		} else if isBodySeperator {
			break
		} else {
			key, val, err := parseHeader(line)
			if err != nil {
				return request, err
			}
			key = strings.ToLower(key)
			request.Headers[key] = val
			if key == "content-length" {
				length, err := strconv.Atoi(val)
				if err != nil {
					return request, fmt.Errorf("header Content-Length has invalid value. Expected integer. Was %s", val)
				}
				bodyLength = length
			}
		}
	}

	if bodyLength == 0 {
		return request, nil
	}

	// parse body
	buf := make([]byte, bodyLength)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return request, err
	}

	request.Body = string(buf)
	return request, nil
}

func parseRequestLine(line string) (method string, endpoint string, err error) {
	parts := strings.Split(strings.TrimSpace(line), " ")
	if len(parts) != 3 {
		return "", "", fmt.Errorf("invalid request line: %s", line)
	}

	method = strings.ToUpper(parts[0])
	endpoint = parts[1]
	version := parts[2]

	if strings.Split(version, ".")[0] != "HTTP/1" {
		return method, endpoint, fmt.Errorf("invalid version. Expected HTTP/1.x. Was %s", version)
	}

	validMethods := []string{
		"GET", "POST", "PUT", "PATCH", "DELETE", "OPTION", "HEAD",
	}

	for _, m := range validMethods {
		if m == method {
			return method, endpoint, nil
		}
	}

	return method, endpoint, fmt.Errorf("invalid method: %s", method)
}

func parseHeader(line string) (key string, value string, err error) {
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid header. Required at least 2 parts. Was %d. %s", len(parts), line)
	}
	key = parts[0]
	value = strings.TrimSpace(strings.Join(parts[1:], ":"))
	return key, value, nil
}
