package wallet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
)

func configureAPIServerSign() (*operations.MassaWalletAPI, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	localAPI := operations.NewMassaWalletAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	// Create a mock password prompt function
	pwdPrompt := func(string) (string, error) {
		return "1234", nil
	}

	localAPI.RestWalletSignOperationHandler = NewSign(pwdPrompt)
	// params := localAPI.RestWalletSignOperationHandler
	server.ConfigureAPI()

	return localAPI, nil
}

func Test_walletSign_Handle(t *testing.T) {
	api_Sign, err := configureAPIServerSign()
	if err != nil {
		panic(err)
	}

	type want struct {
		header     http.Header
		statusCode int
	}
	testsSign := []struct {
		name string
		body string
		want want
	}{
		{"passing", "toto", want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 200}},
		{"failed", "titi", want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 200}},
	}
	for _, tt := range testsSign {
		t.Run(tt.name, func(t *testing.T) {
			handler_sign, exist := api_Sign.HandlerFor("post", "/rest/wallet/{nickname}/signOperation")
			if !exist {
				t.Fatalf("Endpoint doesn't exist")
			}
			httpRequest, err := http.NewRequest("POST", fmt.Sprintf("/rest/wallet/%s/signOperation", tt.body), strings.NewReader(""))
			if err != nil {
				t.Fatalf(err.Error())
			}

			httpRequest.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			handler_sign.ServeHTTP(resp, httpRequest)

			if resp.Result().StatusCode != tt.want.statusCode {
				t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, tt.want.statusCode)
			}
		})
	}
}
