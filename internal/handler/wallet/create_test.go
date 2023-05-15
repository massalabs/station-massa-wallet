package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func Test_walletCreate_Handle(t *testing.T) {
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

			assert.True(t, result.Success, "New password created")

			assertWallet(t, nickname)

			err = cleanupTestData([]string{nickname})
			assert.NoError(t, err)
		}
	}
}

// createTestWallet tests the creation of a wallet.
func createTestWallet(t *testing.T, api *operations.MassaWalletAPI, testName string, nickname string, statusCode int) {
	t.Run(testName, func(t *testing.T) {
		handler, exist := api.HandlerFor("post", "/api/accounts/{nickname}")
		assert.True(t, exist)

		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s", nickname), "")
		assert.NoError(t, err)

		assert.Equal(t, statusCode, resp.Result().StatusCode)

		if resp.Result().StatusCode != http.StatusOK {
			return
		}

		var wallet models.Account
		err = json.Unmarshal(resp.Body.Bytes(), &wallet)

		assert.NoError(t, err)
		assert.Equal(t, models.Nickname(nickname), wallet.Nickname)
	})
}

func assertWallet(t *testing.T, nickname string) {
	wallet, err := wallet.Load(nickname)
	assert.NoError(t, err)

	assert.Equal(t, uint8(0), wallet.Version)
	assert.Equal(t, nickname, wallet.Nickname)

	minAddrLen := 52
	assert.GreaterOrEqual(t, len(wallet.Address), minAddrLen)
}
