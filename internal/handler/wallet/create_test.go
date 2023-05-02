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
		nickname   string
		password   string
		statusCode int
	}{
		{"create success", nickname, "GoodPassword", http.StatusOK},
		{"invalid Nickname", " ", "", http.StatusBadRequest},
		{"Prompt action canceled by user", nickname, "cancel", http.StatusUnauthorized},
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

		createTestWallet(t, api, test.name, test.nickname, test.statusCode)

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
func createTestWallet(t *testing.T, api *operations.MassaWalletAPI, testName string, nickname string, statusCode int) {
	t.Run(testName, func(t *testing.T) {
		handler, exist := api.HandlerFor("post", "/api/accounts/{nickname}")
		if !exist {
			panic("Endpoint doesn't exist")
		}
		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s", nickname), "")
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

		if wallet.Nickname != models.Nickname(nickname) {
			t.Fatalf("the wallet nickname was: %s, want %s", wallet.Nickname, `toto`)
		}
	})
}
