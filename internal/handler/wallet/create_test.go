package wallet

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
)

func configureAPIServeCreate() (*operations.MassaWalletAPI, error) {
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

	server.ConfigureAPI()

	return localAPI, nil
}

func Test_walletCreate_Handle(t *testing.T) {
	api_create, err := configureAPIServeCreate()
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name       string
		body       string
		statusCode int
	}{
		{"passing", `{"Nickname": "toto", "Password": "1234"}`, 200},
		{"without Password", `{"Nickname": "toto"}`, 422},
		{"without Nickname", `{"Password": "1234"}`, 422},
		{"without Password and Nickname", `{}`, 422},
	}
	for _, tt := range tests {
		createTestWallet(t, api_create, tt.name, tt.body, tt.statusCode)
	}
	// Run the cleanupTestData function after running the tests
	// createTestWallet Clean up test data by deleting the created wallets
	err = cleanupTestData([]string{"toto"})
	if err != nil {
		log.Printf("Error while cleaning up TestData ")
	}
}

func createTestWallet(t *testing.T, api *operations.MassaWalletAPI, name string, inputBody string, statusCode int) {
	t.Run(name, func(t *testing.T) {
		handler, exist := api.HandlerFor("post", "/rest/wallet")
		if !exist {
			t.Fatalf("Endpoint doesn't exist")
		}

		httpRequest, err := http.NewRequest("POST", "/rest/wallet", strings.NewReader(inputBody))
		if err != nil {
			t.Fatalf(err.Error())
		}

		httpRequest.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, httpRequest)

		if resp.Result().StatusCode != statusCode {
			t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, statusCode)
		}

		if resp.Result().StatusCode != 200 {
			return
		}

		var wallet models.Wallet
		err = json.Unmarshal(resp.Body.Bytes(), &wallet)
		if err != nil {
			t.Fatalf("impossible to hydrate models.Wallet: %s", err)
		}

		var body operations.RestWalletCreateBody
		err = json.Unmarshal([]byte(inputBody), &body)
		if err != nil {
			t.Fatalf("impossible to hydrate operations.RestWalletCreateBody: %s", err)
		}

		if *wallet.Nickname != *body.Nickname {
			t.Fatalf("the wallet nickname was: %s, want %s", *wallet.Nickname, `toto`)
		}
	})
}
