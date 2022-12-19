package wallet

import (
	"encoding/json"
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
	api_create, err := configureAPIServer()
	if err != nil {
		panic(err)
	}
	api_get, err := configureAPIServerGet()
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

	testsGet := []struct {
		name string
		body string
		want want
	}{
		{"Passed_list_empty", ``, want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 200}},
		{"Passed_list_with_wallets", `{"Nickname": "toto"}`, want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 200}},
	}

	t.Run(testsGet[0].name, func(t *testing.T) {
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

		if resp.Result().StatusCode != testsGet[0].want.statusCode {
			t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, testCreate.want.statusCode)
		}
	})

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

	t.Run(testsGet[1].name, func(t *testing.T) {
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

		if resp.Result().StatusCode != testsGet[1].want.statusCode {
			t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, testCreate.want.statusCode)
		}

		if resp.Result().StatusCode != 200 {
			return
		}
		var wallets []models.Wallet
		err = json.Unmarshal(resp.Body.Bytes(), &wallets)
		if err != nil {
			t.Fatalf("impossible to hydrate models.Wallet: %s", err)
		}
		var body operations.RestWalletCreateBody
		err = json.Unmarshal([]byte(testsGet[1].body), &body)
		if err != nil {
			t.Fatalf("impossible to hydrate operations.RestWalletCreateBody: %s", err)
		}

		if *wallets[0].Nickname != *body.Nickname {
			t.Fatalf("the wallet nickname was: %s, want %s", *wallets[0].Nickname, `toto`)
		}
	})
}
