package wallet

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
)

func configureAPIServerImport() (*operations.MassaWalletAPI, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	localAPI := operations.NewMassaWalletAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	var walletStorage sync.Map

	localAPI.RestWalletImportHandler = NewImport(&walletStorage)

	server.ConfigureAPI()

	return localAPI, nil
}

func Test_walletImport_Handle(t *testing.T) {
	api_Import, err := configureAPIServerImport()
	if err != nil {
		panic(err)
	}

	type want struct {
		header     http.Header
		statusCode int
	}
	testsImport := struct {
		name string
		body string
		want want
	}{"create_passing", `{
		"address": "A12TFdPyw8Sg9qouzgTWwW5yo5PBDu5C3BWEGPjB9vRx9s3b42qv",
		"keyPairs": [
			{
				"nonce": "86zrpLuzBXBtePQiC5b4d1",
				"privateKey": "HvAH6XuMNamRCuCuAaGsUKrCSjFibwyZ35aHZ4zBd5iNM5x2YM74vLUUH9KhAxDKGxWJ4V3YWNvGGiziGjC4yA1J72NKmcVMitHZM23eW44FAHay4iA",
				"publicKey": "ub4aTM9RSBydGJCbkxe8v7GqWpZNNXuh7uGQgthBpaWhocvA1",
				"salt": "4B28WQKc6jaYN6ymx6xoX8GzwHqF"
			}
		],
		"nickname": "imported"
	}`, want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 204}}

	t.Run(testsImport.name, func(t *testing.T) {

		handler_create, exist := api_Import.HandlerFor("put", "/rest/wallet")
		if !exist {
			t.Fatalf("Endpoint doesn't exist")
		}

		httpRequest, err := http.NewRequest("PUT", "/rest/wallet", strings.NewReader(testsImport.body))
		if err != nil {
			t.Fatalf(err.Error())
		}

		httpRequest.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		handler_create.ServeHTTP(resp, httpRequest)

		if resp.Result().StatusCode != testsImport.want.statusCode {
			t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, testsImport.want.statusCode)
		}
		// Run the cleanupTestData function after running the tests
		// cleanupTestData Clean up test data by listing all created wallets with tests and deleting them
		t.Run("cleanupTestData", func(t *testing.T) { cleanupTestData(t) })
	})
}
