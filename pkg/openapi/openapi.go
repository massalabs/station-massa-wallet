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

// NewNotFoundResponder creates a not found reponse.
//
// A CustomResponder with a "Page not found" body, a "Content-Type" header set to "text/html"
// and a status code of 404 (Not Found) is instanciate.
func NewNotFoundResponder() *CustomResponder {
	return NewCustomResponder(
		[]byte("Page not found"),
		map[string]string{"Content-Type": "text/html"},
		http.StatusNotFound)
}

// NewInternalServerResponder creates an internal server error reponse.
//
// A CustomResponder with a body with the error content, a "Content-Type" header set to "text/html"
// and a status code of 505 (Internal Server Error) is instanciate.
func NewInternalServerErrorResponder(err error) *CustomResponder {
	return NewCustomResponder(
		[]byte(err.Error()),
		map[string]string{"Content-Type": "text/html"},
		http.StatusInternalServerError)
}
