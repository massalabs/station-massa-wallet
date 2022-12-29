package wallet

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
)

func configureAPIServer() (*operations.MassaWalletAPI, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	localAPI := operations.NewMassaWalletAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	var walletStorage sync.Map

	localAPI.RestWalletCreateHandler = NewCreate(&walletStorage)

	server.ConfigureAPI()

	return localAPI, nil
}

func Test_walletCreate_Handle(t *testing.T) {
	api, err := configureAPIServer()
	if err != nil {
		panic(err)
	}

	type want struct {
		header     http.Header
		statusCode int
	}

	tests := []struct {
		name string
		body string
		want want
	}{
		{"passing", `{"Nickname": "toto", "Password": "1234"}`, want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 200}},
		{"without Password", `{"Nickname": "toto"}`, want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 422}},
		{"without Nickname", `{"Password": "1234"}`, want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 422}},
		{"without Password and Nickname", `{}`, want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 422}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, exist := api.HandlerFor("post", "/rest/wallet")
			if !exist {
				t.Fatalf("Endpoint doesn't exist")
			}

			httpRequest, err := http.NewRequest("POST", "/rest/wallet", strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf(err.Error())
			}

			httpRequest.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, httpRequest)

			if resp.Result().StatusCode != tt.want.statusCode {
				t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, tt.want.statusCode)
			}

			if !reflect.DeepEqual(resp.Header(), tt.want.header) {
				t.Fatalf("the response header was: %v, want: %v", resp.Header(), tt.want.header)
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
			err = json.Unmarshal([]byte(tt.body), &body)
			if err != nil {
				t.Fatalf("impossible to hydrate operations.RestWalletCreateBody: %s", err)
			}

			if *wallet.Nickname != *body.Nickname {
				t.Fatalf("the wallet nickname was: %s, want %s", *wallet.Nickname, `toto`)
			}
		})
	}
	// Run the cleanupTestData function after running the tests
	// createTestWallet Clean up test data by listing all created wallets with tests and deleting them
	t.Run("cleanupTestData", func(t *testing.T) { cleanupTestData(t) })
}
