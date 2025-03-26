package wallet

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bluele/gcache"
	"github.com/go-openapi/loads"
	"github.com/massalabs/station-massa-wallet/api/server/restapi"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/cache"
	"github.com/massalabs/station-massa-wallet/pkg/config"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/massalabs/station/pkg/logger"
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

var (
	prompterAppMock *walletPrompterMock
	testAssetStore  *assets.AssetsStore
	testCache       gcache.Cache
)

// MockAPI mocks the wallet API.
// All the wallet endpoints are mocked. You can use the Prompt channel to drive the password entry expected values.
func MockAPI() (*operations.MassaWalletAPI, chan walletapp.EventData, error) {
	os.Setenv("STANDALONE", "1")

	// Load the Swagger specification
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	if err := logger.InitializeGlobal(filepath.Join(walletPath, "unit-test.log")); err != nil {
		log.Fatalf("while initializing global logger: %s", err.Error())
	}

	prompterAppMock = NewWalletPrompterMock(walletapp.NewWalletApp(wallet), resultChannel)

	nodeFetcher := network.NewNodeFetcher()

	assets.SetFileDirOverride(walletPath)

	testAssetStore = assets.InitAssetsStore(nodeFetcher)

	massaNodeMock := NewNodeFetcherMock()

	// Load config file with config file path override
	config.SetConfigFileDirOverride(walletPath)
	config.Load()

	testCache = cache.Init()

	// Set wallet API endpoints
	AppendEndpoints(massaWalletAPI, prompterAppMock, massaNodeMock)

	// instantiates the server configure its API.
	server := restapi.NewServer(massaWalletAPI)
	server.ConfigureAPI()

	return massaWalletAPI, resultChannel, err
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
