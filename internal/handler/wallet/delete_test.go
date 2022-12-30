package wallet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
)

// configureAPIServerDelete configures the API server for the DELETE wallet operation
func configureAPIServerDelete() (*operations.MassaWalletAPI, error) {
	// Load the Swagger specification for the API
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	// Create a new MassaWalletAPI instance
	localAPI := operations.NewMassaWalletAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	// Create a sync.Map to store wallets in memory
	var walletStorage sync.Map

	// Set the handler for the DELETE wallet operation
	localAPI.RestWalletDeleteHandler = NewDelete(&walletStorage)

	// Configure the API server
	server.ConfigureAPI()

	return localAPI, nil
}

// Test_walletDelete_Handle tests the handling of the DELETE wallet operation
func Test_walletDelete_Handle(t *testing.T) {
	// Run the createTestWallet function before running the tests
	// createTestWallet create a wallet called "precondition_wallet" to test the delete function
	api_create, err := configureAPIServeCreate()
	if err != nil {
		panic(err)
	}
	createTestWallet(t, api_create, "precondition_wallet", `{"Nickname": "precondition_wallet", "Password": "1234"}`, 200)

	// Configure the API server for the DELETE wallet operation
	api_Delete, err := configureAPIServerDelete()
	if err != nil {
		panic(err)
	}

	// Define the expected result struct of the test cases
	type want struct {
		header     http.Header
		statusCode int
	}

	// Define the test cases
	// The first test case tests the deletion of the wallet created with the createTestWallet(t) function.
	// The second test case tests the deletion of a wallet that does not exist.
	testsDelete := []struct {
		name string
		body string
		want want
	}{
		{"passing", "precondition_wallet", want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 204}},
		{"failing_wallet_does_not_exist", "wallet_does_not_exist", want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 500}},
	}

	// Iterate through the test cases
	for _, tt := range testsDelete {
		t.Run(tt.name, func(t *testing.T) {
			// Get the handler for the DELETE wallet operation
			handler_delete, exist := api_Delete.HandlerFor("DELETE", "/rest/wallet/{nickname}")
			if !exist {
				t.Fatalf("Endpoint doesn't exist")
			}

			// Create a DELETE request for the wallet with the specified nickname
			httpRequest, err := http.NewRequest("DELETE", fmt.Sprintf("/rest/wallet/%s", tt.body), strings.NewReader(""))
			if err != nil {
				t.Fatalf(err.Error())
			}

			httpRequest.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			handler_delete.ServeHTTP(resp, httpRequest)

			if resp.Result().StatusCode != tt.want.statusCode {
				t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, tt.want.statusCode)
			}
		})
	}
}
