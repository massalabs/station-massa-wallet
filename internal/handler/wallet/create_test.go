package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
)

func Test_walletCreate_validation(t *testing.T) {
	api, prompterApp, resChan, err := MockAPI()
	if err != nil {
		panic(err)
	}

	nickname := "GoodNickname"

	tests := []struct {
		name       string
		body       string
		password   string
		statusCode int
	}{
		{"create success", fmt.Sprintf(`{"Nickname": "%s"}`, nickname), "GoodPassword", http.StatusOK},
		{"invalid Nickname", `{"Nickname": "  "}`, "", http.StatusBadRequest},
		{"Nickname Missing", `{}`, "", http.StatusUnprocessableEntity},
		{"Prompt action canceled by user", fmt.Sprintf(`{"Nickname": "%s"}`, nickname), "cancel", http.StatusInternalServerError},
	}
	for _, test := range tests {
		testResult := make(chan walletapp.EventData)

		if len(test.password) > 0 {
			// Send password to prompter app and wait for result
			go func(res chan walletapp.EventData) {
				if test.password == "cancel" {
					// Send cancel to prompter app to unlock the handler
					prompterApp.App().CtrlChan <- walletapp.Cancel
				} else {
					prompterApp.App().PasswordChan <- test.password
					// forward test result to test goroutine
				}
				res <- (<-resChan)
			}(testResult)
		}

		createTestWallet(t, api, test.name, test.body, test.statusCode)

		if len(test.password) > 0 && test.statusCode == http.StatusOK {
			result := <-testResult
			if !result.Success {
				t.Fatalf("Expected success, got error")
			}
			msg := "New password created"
			if result.Data != msg {
				t.Fatalf(fmt.Sprintf("Expected error message to be %s, got %s", msg, result.Data))
			}

			err = cleanupTestData([]string{nickname})
			if err != nil {
				t.Fatalf("while cleaning up TestData: %s", err)
			}

		}
	}
}

// createTestWallet tests the creation of a wallet.
func createTestWallet(t *testing.T, api *operations.MassaWalletAPI, name string, inputBody string, statusCode int) {
	t.Run(name, func(t *testing.T) {
		resp, err := processHTTPRequest(api, "POST", "/api/accounts", inputBody)
		if err != nil {
			t.Fatalf("while serving HTTP request: %s", err)
		}

		if resp.Result().StatusCode != statusCode {
			t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, statusCode)
		}

		if resp.Result().StatusCode != http.StatusOK {
			return
		}

		var wallet models.Account
		err = json.Unmarshal(resp.Body.Bytes(), &wallet)
		if err != nil {
			t.Fatalf("impossible to hydrate models.Account: %s", err)
		}

		var body operations.RestCreateAccountBody
		err = json.Unmarshal([]byte(inputBody), &body)
		if err != nil {
			t.Fatalf("impossible to hydrate operations.RestCreateAccountBody: %s", err)
		}

		if wallet.Nickname != body.Nickname {
			t.Fatalf("the wallet nickname was: %s, want %s", wallet.Nickname, `toto`)
		}
	})
}
