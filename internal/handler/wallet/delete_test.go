package wallet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
)

func deleteWallet(t *testing.T, api *operations.MassaWalletAPI, nickname string) *httptest.ResponseRecorder {
	handler, exist := api.HandlerFor("DELETE", "/api/accounts/{nickname}")
	assert.True(t, exist, "Endpoint doesn't exist")

	resp, err := handleHTTPRequest(handler, "DELETE", fmt.Sprintf("/api/accounts/%s", nickname), "")
	assert.NoError(t, err)
	return resp
}

func Test_walletDelete_Handle(t *testing.T) {
	api, prompterApp, resChan, err := MockAPI()
	assert.NoError(t, err)

	testResult := make(chan walletapp.EventData)

	nickname := "walletToDelete"
	password := "zePassword"
	_, errGenerate := wallet.Generate(nickname, password)
	assert.Nil(t, errGenerate)

	t.Run("invalid nickname", func(t *testing.T) {
		resp := deleteWallet(t, api, "toto")
		verifyStatusCode(t, resp, http.StatusNotFound)
	})

	t.Run("invalid password", func(t *testing.T) {
		// Send password to prompter app and wait for result
		go func() {
			prompterApp.App().PromptInput <- "invalid password"
			// forward test result to test goroutine
			failRes := <-resChan

			checkResultChannel(t, failRes, false, utils.WrongPassword)

			// Send cancel to prompter app to unlock the handler
			prompterApp.App().CtrlChan <- walletapp.Cancel
		}()

		resp := deleteWallet(t, api, nickname)
		verifyStatusCode(t, resp, http.StatusUnauthorized)
	})

	t.Run("canceled by user", func(t *testing.T) {
		go func() {
			// Send wrong password to prompter app and wait for result
			prompterApp.App().PromptInput <- "this is not the password"
			// forward test result to test goroutine
			failRes := <-resChan

			checkResultChannel(t, failRes, false, utils.WrongPassword)

			// Send cancel to prompter app to unlock the handler
			prompterApp.App().CtrlChan <- walletapp.Cancel
		}()

		resp := deleteWallet(t, api, nickname)

		verifyStatusCode(t, resp, http.StatusUnauthorized)

		_, err = wallet.Load(nickname)
		assert.NoError(t, err)
	})

	t.Run("delete success", func(t *testing.T) {
		// Send password to prompter app and wait for result
		go func() {
			prompterApp.App().PromptInput <- password
			// forward test result to test goroutine
			testResult <- (<-resChan)
		}()

		resp := deleteWallet(t, api, nickname)

		verifyStatusCode(t, resp, http.StatusNoContent)

		result := <-testResult

		checkResultChannel(t, result, true, utils.MsgAccountDeleted)

		_, err = wallet.Load(nickname)
		assert.Error(t, err, "Wallet should have been deleted")
	})
}
