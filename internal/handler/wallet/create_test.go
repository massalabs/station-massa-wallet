package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/stretchr/testify/assert"
)

func Test_walletCreate_Handle(t *testing.T) {
	api, resChan, err := MockAPI()
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
		{"nickname already exists", nickname, "GoodPassword", http.StatusBadRequest},
		{"Prompt action canceled by user", nickname, "cancel", http.StatusUnauthorized},
	}

	for _, test := range tests {
		testResult := make(chan walletapp.EventData)

		if len(test.password) > 0 {
			// Send password to prompter app and wait for result
			go func(res chan walletapp.EventData) {
				if test.password == "cancel" {
					// Send cancel to prompter app to unlock the handler
					prompterAppMock.App().CtrlChan <- walletapp.Cancel
				} else {
					prompterAppMock.App().PromptInput <- &walletapp.StringPromptInput{
						BaseMessage: walletapp.BaseMessage{},
						Message:     test.password,
					}
					// forward test result to test goroutine
				}
				res <- (<-resChan)
			}(testResult)
		}

		createTestWallet(t, api, test.name, test.nickname, test.statusCode)

		if len(test.password) > 0 && test.statusCode == http.StatusOK {
			result := <-testResult

			checkResultChannel(t, result, true, "")

			assertWallet(t, prompterAppMock.App().Wallet, nickname)
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

func assertWallet(t *testing.T, wallet *wallet.Wallet, nickname string) {
	acc, err := wallet.GetAccount(nickname)
	assert.NoError(t, err)

	assert.Equal(t, uint8(1), *acc.Version)
	assert.Equal(t, nickname, acc.Nickname)

	addressByteLen := 32
	assert.Len(t, acc.Address.Data, addressByteLen)
}
