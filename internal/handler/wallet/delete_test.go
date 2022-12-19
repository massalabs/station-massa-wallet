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

func configureAPIServerDelete() (*operations.MassaWalletAPI, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	localAPI := operations.NewMassaWalletAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	var walletStorage sync.Map

	localAPI.RestWalletDeleteHandler = NewDelete(&walletStorage)

	server.ConfigureAPI()

	return localAPI, nil
}

func Test_walletDelete_Handle(t *testing.T) {
	api_create, err := configureAPIServer()
	if err != nil {
		panic(err)
	}
	api_Delete, err := configureAPIServerDelete()
	if err != nil {
		panic(err)
	}

	type want struct {
		header     http.Header
		statusCode int
	}

	testCreate := struct {
		name string
		body string
		want want
	}{"create_passing", `{"Nickname": "toto", "Password": "1234"}`, want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 200}}

	testsDelete := []struct {
		name string
		body string
		want want
	}{
		{"passing", "toto", want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 204}},
		{"without Nickname", "titi", want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 500}},
	}
	t.Run(testCreate.name, func(t *testing.T) {

		handler_create, exist := api_create.HandlerFor("post", "/rest/wallet")
		if !exist {
			t.Fatalf("Endpoint doesn't exist")
		}

		httpRequest, err := http.NewRequest("POST", "/rest/wallet", strings.NewReader(testCreate.body))
		if err != nil {
			t.Fatalf(err.Error())
		}

		httpRequest.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		handler_create.ServeHTTP(resp, httpRequest)

		if resp.Result().StatusCode != testCreate.want.statusCode {
			t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, testCreate.want.statusCode)
		}
	})
	for _, tt := range testsDelete {
		t.Run(tt.name, func(t *testing.T) {
			handler_delete, exist := api_Delete.HandlerFor("DELETE", "/rest/wallet/{nickname}")
			if !exist {
				t.Fatalf("Endpoint doesn't exist")
			}

			httpRequest, err := http.NewRequest("DELETE", fmt.Sprintf("/rest/wallet/%s", tt.body), strings.NewReader(""))
			if err != nil {
				t.Fatalf(err.Error())
			}

			httpRequest.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			handler_delete.ServeHTTP(resp, httpRequest)

			if resp.Result().StatusCode != tt.want.statusCode {
				t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, testCreate.want.statusCode)
			}
		})
	}
}
