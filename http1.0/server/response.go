package server

import "fmt"

type HttpResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

func NewResponse() *HttpResponse {
	return &HttpResponse{
		StatusCode: 200,
		Headers:    map[string]string{},
		Body:       "",
	}
}

func (res *HttpResponse) ToBytes() []byte {
	if len(res.Body) > 0 {
		res.Headers["Content-Length"] = fmt.Sprintf("%d", len(res.Body))
	}

	response := fmt.Sprintf("HTTP/1.0 %d \n", res.StatusCode)

	for k, v := range res.Headers {
		response += fmt.Sprintf("%s: %s\n", k, v)
	}

	response += "\n" + res.Body + "\n"
	return []byte(response)
}
