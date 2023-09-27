package openapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomResponder_WriteResponse(t *testing.T) {
	// Test case 1: Successful response
	responder := NewCustomResponder([]byte("Hello, World!"), map[string]string{"Content-Type": "text/plain"}, http.StatusOK)

	recorder := httptest.NewRecorder()
	responder.WriteResponse(recorder, nil)

	assert.Equal(t, "Hello, World!", recorder.Body.String())
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "text/plain", recorder.Header().Get("Content-Type"))

	// Test case 2: Error while writing response body
	responderWithError := NewCustomResponder([]byte("Hello, World!"), map[string]string{"Content-Type": "text/plain"}, http.StatusOK)

	recorderWithError := &mockResponseWriter{errOnWrite: true}
	responderWithError.WriteResponse(recorderWithError, nil)
}

func TestNewNotFoundResponder(t *testing.T) {
	responder := NewNotFoundResponder()

	assert.Equal(t, "Page not found", string(responder.Body))
	assert.Equal(t, http.StatusNotFound, responder.StatusCode)
	assert.Equal(t, "text/html", responder.Header["Content-Type"])
}

func TestNewInternalServerErrorResponder(t *testing.T) {
	err := DummyFunction()
	responder := NewInternalServerErrorResponder(err)

	assert.Equal(t, err.Error(), string(responder.Body))
	assert.Equal(t, http.StatusInternalServerError, responder.StatusCode)
	assert.Equal(t, "text/html", responder.Header["Content-Type"])
}

// Mock response writer for testing
type mockResponseWriter struct {
	httptest.ResponseRecorder
	errOnWrite bool
}

func (m *mockResponseWriter) Write(p []byte) (int, error) {
	if m.errOnWrite {
		return 0, assert.AnError
	}
	return m.ResponseRecorder.Write(p)
}

func DummyFunction() error {
	return fmt.Errorf("Some error occurred")
}
