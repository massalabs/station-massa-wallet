package openapi

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

type CustomResponder struct {
	Body       []byte
	Header     map[string]string
	StatusCode int
}

func (c *CustomResponder) WriteResponse(writer http.ResponseWriter, producer runtime.Producer) {
	for k, v := range c.Header {
		writer.Header().Set(k, v)
	}

	writer.WriteHeader(c.StatusCode)

	_, err := writer.Write(c.Body)
	if err != nil {
		panic(err)
	}
}

func NewCustomResponder(body []byte, header map[string]string, statusCode int) *CustomResponder {
	return &CustomResponder{Body: body, Header: header, StatusCode: statusCode}
}

func NewNotFoundResponder() *CustomResponder {
	return NewCustomResponder(
		[]byte("Page not found"),
		map[string]string{"Content-Type": "text/html"},
		http.StatusNotFound)
}

func NewInternalServerErrorResponder(err error) *CustomResponder {
	return NewCustomResponder(
		[]byte(err.Error()),
		map[string]string{"Content-Type": "text/html"},
		http.StatusInternalServerError)
}
