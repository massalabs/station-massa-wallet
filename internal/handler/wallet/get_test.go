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

func configureAPIServerGet() (*operations.MassaWalletAPI, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	localAPI := operations.NewMassaWalletAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	var walletStorage sync.Map

	localAPI.RestWalletListHandler = NewList(&walletStorage)

	server.ConfigureAPI()

	return localAPI, nil
}

func Test_walletGet_Handle(t *testing.T) {
	api_get, err := configureAPIServerGet()
	if err != nil {
		panic(err)
	}
	t.Run("Passed_list_empty", func(t *testing.T) {
		handler_get, exist := api_get.HandlerFor("get", "/rest/wallet")
		if !exist {
			t.Fatalf("Endpoint doesn't exist")
		}

		httpRequest, err := http.NewRequest("Get", "/rest/wallet", strings.NewReader(""))
		if err != nil {
			t.Fatalf(err.Error())
		}

		httpRequest.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		handler_get.ServeHTTP(resp, httpRequest)

		if resp.Result().StatusCode != 200 {
			t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, 200)
		}

		var wallets []models.Wallet
		err = json.Unmarshal(resp.Body.Bytes(), &wallets)
		if err != nil {
			t.Fatalf("impossible to hydrate models.Wallet: %s", err)
		}

		if len(wallets) != 0 {
			t.Fatalf("the wallets list should be empty")
		}
	})

	t.Run("Passed_list_with_wallets", func(t *testing.T) {
		// createTestWallet create a wallet called "precondition_wallet" to test the get function
		api_create, err := configureAPIServeCreate()
		if err != nil {
			panic(err)
		}
		createTestWallet(t, api_create, "precondition_wallet", `{"Nickname": "precondition_wallet", "Password": "1234"}`, 200)
		handler_get, exist := api_get.HandlerFor("get", "/rest/wallet")
		if !exist {
			t.Fatalf("Endpoint doesn't exist")
		}

		httpRequest, err := http.NewRequest("Get", "/rest/wallet", strings.NewReader(""))
		if err != nil {
			t.Fatalf(err.Error())
		}

		httpRequest.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		handler_get.ServeHTTP(resp, httpRequest)

		if resp.Result().StatusCode != 200 {
			t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, 200)
		}

		if resp.Result().StatusCode != 200 {
			return
		}
		var wallets []models.Wallet
		err = json.Unmarshal(resp.Body.Bytes(), &wallets)
		if err != nil {
			t.Fatalf("impossible to hydrate models.Wallet: %s", err)
		}

		if *wallets[0].Nickname != "precondition_wallet" {
			t.Fatalf("the wallet nickname was: %s, want %s", *wallets[0].Nickname, "precondition_wallet")
		}
	})

	// Run the cleanupTestData function after running the tests
	// createTestWallet Clean up test data by deleting the created wallet.
	err = cleanupTestData([]string{"precondition_wallet"})
	if err != nil {
		log.Printf("Error while cleaning up TestData ")
	}
}
