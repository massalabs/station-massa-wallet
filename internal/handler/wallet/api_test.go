package wallet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/bluele/gcache"
	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
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

// testPrompter implements the password.PasswordAsker interface for test purpose.
type testPrompterPassword struct {
	mockPasswordEntry chan PasswordPrompt
}

// Ask simulates a password entry by returning the content given through the mockPasswordEntry channel.
func (t *testPrompterPassword) Ask(name string) (string, error) {
	passwordPrompter := <-t.mockPasswordEntry
	return passwordPrompter.Password, passwordPrompter.Err
}

type PrivateKeyPrompt struct {
	PrivateKey string
	Err        error
}

type testPrompterPrivatekey struct {
	mockPrivateKeyEntry chan PrivateKeyPrompt
}

type testConfirmDelete struct {
	mockPasswordEntry chan PasswordPrompt
}

// Ask simulates a private key entry by returning the content given through the mockPrivateKeyEntry channel.
func (t *testPrompterPrivatekey) Ask() (string, error) {
	PrivateKeyPrompter := <-t.mockPrivateKeyEntry
	return PrivateKeyPrompter.PrivateKey, PrivateKeyPrompter.Err
}

// Confirm simulates a password entry by returning the content given through the mockPrivateKeyEntry channel.
func (t *testConfirmDelete) Confirm(walletName string) (string, error) {
	confirmDeletePrompter := <-t.mockPasswordEntry
	return confirmDeletePrompter.Password, confirmDeletePrompter.Err
}

// MockAPI mocks the wallet API.
// All the wallet endpoints are mocked. You can use the Prompt channel to drive the password entry expected values.
func MockAPI() (*operations.MassaWalletAPI, wallet.WalletPrompterInterface, chan walletapp.EventData, error) {
	// Load the Swagger specification
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, nil, nil, err
	}
	// Create a new MassaWalletAPI instance
	massaWalletAPI := operations.NewMassaWalletAPI(swaggerSpec)

	resultChannel := make(chan walletapp.EventData)

	prompterApp := NewWalletPrompterMock(walletapp.NewWalletApp(), resultChannel)

	// Set wallet API endpoints
	AppendEndpoints(massaWalletAPI, prompterApp, gcache.New(20).LRU().Build())

	// instantiates the server configure its API.
	server := restapi.NewServer(massaWalletAPI)
	server.ConfigureAPI()

	return massaWalletAPI, prompterApp, resultChannel, err
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
