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
	testsImport := []struct {
		name string
		body string
		want want
	}{
		{"passing", `{
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
	}`, want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 204}},

		{"fail_empty_fields", `{}`, want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 422}}}
	for _, tt := range testsImport {
		t.Run(tt.name, func(t *testing.T) {
			handler_create, exist := api_Import.HandlerFor("put", "/rest/wallet")
			if !exist {
				t.Fatalf("Endpoint doesn't exist")
			}

			httpRequest, err := http.NewRequest("PUT", "/rest/wallet", strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf(err.Error())
			}

			httpRequest.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			handler_create.ServeHTTP(resp, httpRequest)

			if resp.Result().StatusCode != tt.want.statusCode {
				t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, tt.want.statusCode)
			}

			// Run the cleanupTestData function after running the tests
			// createTestWallet Clean up test data by deleting the created wallets
			err = cleanupTestData([]string{"imported"})
			if err != nil {
				log.Printf("Error while cleaning up TestData ")
			}
		})
	}
}
