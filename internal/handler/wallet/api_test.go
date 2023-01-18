package wallet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
)

// Prompt struct will be used to drive the password prompter externally
// directly from the tests we will set the value to be returned
// hence we will be able to:
// - return the right, wrong or an empty password
// - return an error
type PasswordPrompt struct {
	Password string
	Err      error
}

type PrivateKeyPrompt struct {
	PrivateKey string
	Err        error
}

// testPrompter implements the password.PasswordAsker interface for test purpose.
type testPrompterPassword struct {
	mockPasswordEntry chan PasswordPrompt
}

type testPrompterPrivatekey struct {
	mockPrivateKeyEntry chan PrivateKeyPrompt
}

// Ask simulates a password entry by returning the content given through the mockPasswordEntry channel.
func (t *testPrompterPassword) Ask(name string) (string, error) {
	passwordPrompter := <-t.mockPasswordEntry
	return passwordPrompter.Password, passwordPrompter.Err
}

// Ask simulates a password entry by returning the content given through the mockPasswordEntry channel.
func (t *testPrompterPrivatekey) Ask() (string, error) {
	PrivateKeyPrompter := <-t.mockPrivateKeyEntry
	return PrivateKeyPrompter.PrivateKey, PrivateKeyPrompter.Err
}

// MockAPI mocks the wallet API.
// All the wallet endpoints are mocked. You can use the Prompt channel to drive the password entry expected values.
func MockAPI() (*operations.MassaWalletAPI, chan PasswordPrompt, chan PrivateKeyPrompt, error) {
	mockChanPassword := make(chan PasswordPrompt, 2) // buffered channel
	mockChanPrivateKey := make(chan PrivateKeyPrompt, 2)
	// Load the Swagger specification
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, nil, nil, err
	}

	// Create a new MassaWalletAPI instance
	massaWalletAPI := operations.NewMassaWalletAPI(swaggerSpec)

	// Set wallet API endpoints
	AppendEndpoints(massaWalletAPI, &testPrompterPrivatekey{mockPrivateKeyEntry: mockChanPrivateKey}, &testPrompterPassword{mockPasswordEntry: mockChanPassword})

	// instantiates the server configure its API.
	server := restapi.NewServer(massaWalletAPI)
	server.ConfigureAPI()

	return massaWalletAPI, mockChanPassword, mockChanPrivateKey, err
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
