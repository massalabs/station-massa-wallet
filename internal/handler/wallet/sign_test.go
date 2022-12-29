package wallet

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
)

// this prompt struct will be used to drive the password prompter externally
// directly from the tests we will set the value to be returned
// hence we will be able to:
// - return the right, wrong or an empty password
// - return an error
type prompt struct {
	password string
	err      error
}

func configureAPIServerSign(prompt chan prompt) (*operations.MassaWalletAPI, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	localAPI := operations.NewMassaWalletAPI(swaggerSpec)
	server := restapi.NewServer(localAPI)

	// Create a mock password prompt function using the given channel
	pwdPrompt := func(string) (string, error) {
		result := <-prompt
		return result.password, result.err
	}

	localAPI.RestWalletSignOperationHandler = NewSign(pwdPrompt)
	server.ConfigureAPI()

	return localAPI, nil
}

func Test_walletSign_Handle(t *testing.T) {
	// Run the createTestWallet function before running the tests
	// createTestWallet create a wallet called "precondition_wallet" to test the delete function
	createTestWallet(t)

	channel := make(chan prompt, 1) // buffered channel
	api_Sign, err := configureAPIServerSign(channel)
	if err != nil {
		panic(err)
	}

	type want struct {
		statusCode int
	}

	testsSign := []struct {
		name         string
		nickname     string
		body         string
		promptResult prompt
		want         want
	}{
		{"passing", "precondition_wallet", `{"operation":"MjIzM3QyNHQ="}`, prompt{password: "1234", err: nil}, want{statusCode: 200}},
		{"wrong password", "precondition_wallet", `{"operation":"MjIzM3QyNHQ="}`, prompt{password: "4321", err: nil}, want{statusCode: 500}},
		{"wrong nickname", "titi", `{"operation":"MjIzM3QyNHQ="}`, prompt{password: "1234", err: nil}, want{statusCode: 500}},
		// to debug why this test never finish
		// {"prompt error", "titi", `{"operation":"MjIzM3QyNHQ="}`, prompt{password: "1234", err: errors.New("Error while getting password prompt")}, want{statusCode: 500}},
	}
	for _, tt := range testsSign {
		t.Run(tt.name, func(t *testing.T) {
			handler_sign, exist := api_Sign.HandlerFor("post", "/rest/wallet/{nickname}/signOperation")
			if !exist {
				t.Fatalf("Endpoint doesn't exist")
			}

			httpRequest, err := http.NewRequest("POST", fmt.Sprintf("/rest/wallet/%s/signOperation", tt.nickname), strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf(err.Error())
			}

			channel <- tt.promptResult // non blocking call as channel is buffered

			httpRequest.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			handler_sign.ServeHTTP(resp, httpRequest)

			if resp.Result().StatusCode != tt.want.statusCode {
				bodyBytes := new(strings.Builder)
				_, _ = io.Copy(bodyBytes, resp.Result().Body)

				body := strings.TrimSpace(bodyBytes.String())
				t.Logf("the returned body is: %s", body)
				t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, tt.want.statusCode)
			}
		})
	}
	// Run the cleanupTestData function after running the tests
	// cleanupTestData Clean up test data by listing all created wallets with tests and deleting them
	t.Run("cleanupTestData", func(t *testing.T) { cleanupTestData(t) })
}
