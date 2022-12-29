package wallet

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
)

// configureAPIServerPrecondition sets up the API server with the required endpoints
// and handler functions for the precondition tests
func configureAPIServerPrecondition() (*operations.MassaWalletAPI, error) {
	// Load the Swagger specification
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	// Create a new MassaWalletAPI instance
	localAPI := operations.NewMassaWalletAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	// Create a sync.Map to store the wallets
	var walletStorage sync.Map

	// Set the handler function for the create wallet endpoint
	localAPI.RestWalletCreateHandler = NewCreate(&walletStorage)

	// Configure the API endpoints
	server.ConfigureAPI()

	return localAPI, nil
}

// createTestWallet creates a new wallet which will be used as a preconditon to run other tests
func createTestWallet(t *testing.T) {
	t.Log("\n................... createTestWallet start ...................\n")
	// Configure the API server
	api, err := configureAPIServerPrecondition()
	if err != nil {
		// Panic if there was an error configuring the API server
		panic(err)
	}

	// Set up the request body and expected response
	requestBody := `{"Nickname": "precondition_wallet", "Password": "1234"}`
	expectedHeader := http.Header{"Content-Type": {"application/json"}}
	expectedStatusCode := 200
	testPrecondition := struct {
		name string
		body string
		want struct {
			header     http.Header
			statusCode int
		}
	}{
		"precondition_wallet",
		requestBody,
		struct {
			header     http.Header
			statusCode int
		}{
			expectedHeader,
			expectedStatusCode,
		},
	}

	// Get the handler function for the create wallet endpoint
	handler, exist := api.HandlerFor("post", "/rest/wallet")
	if !exist {
		// Log an error and exit if the endpoint does not exist
		log.Fatalln("Endpoint '/rest/wallet' does not exist")
	}

	// Create a new HTTP request to the create wallet
	t.Logf("\n................... create wallet %v ...................\n", strings.NewReader(testPrecondition.body))
	httpRequest, err := http.NewRequest("POST", "/rest/wallet", strings.NewReader(testPrecondition.body))
	if err != nil {
		// Log an error and exit if there was an error creating the request
		log.Fatalln(err.Error())
	}

	// Set the request header to specify that the body of the request is in JSON format
	httpRequest.Header.Set("Content-Type", "application/json")

	// Create a new HTTP response recorder to record the response from the handler
	resp := httptest.NewRecorder()

	// Serve the HTTP request to the handler function
	handler.ServeHTTP(resp, httpRequest)

	// Check if the status code of the response is the expected value
	if resp.Result().StatusCode != testPrecondition.want.statusCode {
		// Log an error and exit if the status code is not the expected value
		log.Fatalln("Unexpected status code: got", resp.Result().StatusCode, "want", testPrecondition.want.statusCode)
	}
	t.Log("\n................... createTestWallet complete ...................\n")

}
