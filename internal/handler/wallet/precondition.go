package wallet

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
)

func configureAPIServerPrecondition() (*operations.MassaWalletAPI, error) {
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
func Test_walletPrecondition_Handle() {
	api_create, err := configureAPIServerPrecondition()
	if err != nil {
		panic(err)
	}
	type want struct {
		header     http.Header
		statusCode int
	}

	testPrecondition := struct {
		name string
		body string
		want want
	}{"precondition_wallet", `{"Nickname": "precondition_wallet", "Password": "1234"}`, want{header: http.Header{"Content-Type": {"application/json"}}, statusCode: 200}}
	handler_create, exist := api_create.HandlerFor("post", "/rest/wallet")
	if !exist {
		log.Fatalln("Endpoint doesn't exist")
	}

	httpRequest, err := http.NewRequest("POST", "/rest/wallet", strings.NewReader(testPrecondition.body))
	if err != nil {
		log.Fatalln(err.Error())
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	handler_create.ServeHTTP(resp, httpRequest)

	if resp.Result().StatusCode != testPrecondition.want.statusCode {
		log.Fatalln("Endpoint doesn't exist")
		// log.Fatalln("the status code was: %d, want %d", resp.Result().StatusCode, testPrecondition.want.statusCode)
	}

}
