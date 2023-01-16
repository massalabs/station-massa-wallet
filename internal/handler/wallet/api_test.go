package wallet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
)

// Prompt struct will be used to drive the password prompter externally
// directly from the tests we will set the value to be returned
// hence we will be able to:
// - return the right, wrong or an empty password
// - return an error
type Prompt struct {
	Password string
	Err      error
}

// testPrompter implements the password.Asker interface for test purpose.
type testPrompter struct {
	mockPasswordEntry chan Prompt
}

// Ask simulates a password entry by returning the content given through the mockPasswordEntry channel.
func (t *testPrompter) Ask(name string) (string, error) {
	result := <-t.mockPasswordEntry
	return result.Password, result.Err
}

// MockAPI mocks the wallet API.
// All the wallet endpoints are mocked. You can use the Prompt channel to drive the password entry expected values.
func MockAPI() (*operations.MassaWalletAPI, chan Prompt, error) {
	mockChan := make(chan Prompt, 2) // buffered channel

	// Load the Swagger specification
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, nil, err
	}

	// Create a new MassaWalletAPI instance
	massaWalletAPI := operations.NewMassaWalletAPI(swaggerSpec)
	var app *fyne.App

	// Set wallet API endpoints
	AppendEndpoints(massaWalletAPI, app, &testPrompter{mockPasswordEntry: mockChan})

	// instantiates the server configure its API.
	server := restapi.NewServer(massaWalletAPI)
	server.ConfigureAPI()

	return massaWalletAPI, mockChan, err
}

// processHTTPRequest simulates the processing of an HTTP request on the given API.
func processHTTPRequest(api *operations.MassaWalletAPI, httpMethod string, endpoint string, body string) (*httptest.ResponseRecorder, error) {
	handler, exist := api.HandlerFor(httpMethod, endpoint)
	if !exist {
		return nil, fmt.Errorf("Endpoint doesn't exist.")
	}

	return handleHTTPRequest(handler, httpMethod, endpoint, body)
}

// handleHTTPRequest handles the processing of an HTTP request on the given API.
func handleHTTPRequest(handler http.Handler, httpMethod string, endpoint string, body string) (*httptest.ResponseRecorder, error) {
	httpRequest, err := http.NewRequest(httpMethod, endpoint, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, httpRequest)

	return resp, nil
}
