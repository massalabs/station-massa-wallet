package wallet

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/bluele/gcache"
	"github.com/go-openapi/loads"
	"github.com/massalabs/station-massa-wallet/api/server/restapi"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/stretchr/testify/assert"
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

// MockAPI mocks the wallet API.
// All the wallet endpoints are mocked. You can use the Prompt channel to drive the password entry expected values.
func MockAPI() (*operations.MassaWalletAPI, prompt.WalletPrompterInterface, *assets.AssetsStore, chan walletapp.EventData, error) {
	os.Setenv("STANDALONE", "1")

	// Load the Swagger specification
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, nil, nil, nil, err
	}
	// Create a new MassaWalletAPI instance
	massaWalletAPI := operations.NewMassaWalletAPI(swaggerSpec)

	resultChannel := make(chan walletapp.EventData)

	walletPath, err := os.MkdirTemp(os.TempDir(), "*-wallet-dir")
	if err != nil {
		log.Fatalf("while creating temporary wallet directory: %s", err.Error())
	}

	wallet, err := wallet.New(walletPath)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	prompterApp := NewWalletPrompterMock(walletapp.NewWalletApp(wallet), resultChannel)

	nodeFetcher := network.NewNodeFetcher()

	AssetsStore, err := assets.NewAssetsStore(walletPath, nodeFetcher)
	if err != nil {
		log.Fatalf("Failed to create AssetsStore: %v", err)
	}

	massaNodeMock := NewNodeFetcherMock()

	// Set wallet API endpoints
	AppendEndpoints(massaWalletAPI, prompterApp, massaNodeMock, AssetsStore, gcache.New(20).LRU().Build())

	// instantiates the server configure its API.
	server := restapi.NewServer(massaWalletAPI)
	server.ConfigureAPI()

	return massaWalletAPI, prompterApp, AssetsStore, resultChannel, err
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

func checkResultChannel(t *testing.T, result walletapp.EventData, success bool, expectedCodeMsg string) {
	assert.Equal(t, success, result.Success)
	assert.Equal(t, expectedCodeMsg, result.CodeMessage)
}
