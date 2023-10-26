package openapi

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-openapi/runtime"
)

// CustomResponder is a go-swagger compatible struct that allows the caller to fully customize the HTTP response.
// The caller can define the response body, header, and status code.
type CustomResponder struct {
	Body       []byte
	Header     map[string]string
	StatusCode int
}

// WriteResponse writes the defined HTTP response to the provided ResponseWriter.
//
// It sets the response headers, status code, and body. If an error occurs while writing the body, the error is logged to Stderr.
// Note: runtime.Producer is ignored.
func (c *CustomResponder) WriteResponse(writer http.ResponseWriter, producer runtime.Producer) {
	for k, v := range c.Header {
		writer.Header().Set(k, v)
	}

	writer.WriteHeader(c.StatusCode)

	_, err := writer.Write(c.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "while writing HTTP response: %s\n", err.Error())
	}
}

// NewCustomResponder creates a new CustomResponder with the given body, header, and status code.
func NewCustomResponder(body []byte, header map[string]string, statusCode int) *CustomResponder {
	return &CustomResponder{Body: body, Header: header, StatusCode: statusCode}
}

type PayloadResponder struct {
	StatusCode int
	Payload    interface{}
}

func (p *PayloadResponder) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(p.StatusCode)

	if p.Payload != nil {
		payload := p.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

func NewPayloadResponder(statusCode int, payload interface{}) *PayloadResponder {
	return &PayloadResponder{StatusCode: statusCode, Payload: payload}
}
