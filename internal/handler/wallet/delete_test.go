package wallet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func deleteWallet(t *testing.T, api *operations.MassaWalletAPI, nickname string) *httptest.ResponseRecorder {
	handler, exist := api.HandlerFor("DELETE", "/api/accounts/{nickname}")
	assert.True(t, exist, "Endpoint doesn't exist")

	resp, err := handleHTTPRequest(handler, "DELETE", fmt.Sprintf("/api/accounts/%s", nickname), "")
	assert.NoError(t, err)

	return resp
}

func Test_walletDelete_Handle(t *testing.T) {
	api, prompterApp, _, resChan, err := MockAPI()
	assert.NoError(t, err)

	testResult := make(chan walletapp.EventData)

	nickname := "walletToDelete"
	password := "zePassword"
	createAccount(password, nickname, t, prompterApp)

	t.Run("invalid nickname", func(t *testing.T) {
		resp := deleteWallet(t, api, "toto")
		verifyStatusCode(t, resp, http.StatusNotFound)
	})

	t.Run("invalid password", func(t *testing.T) {
		// Send password to prompter app and wait for result
		go func() {
			prompterApp.App().PromptInput <- &walletapp.StringPromptInput{
				BaseMessage: walletapp.BaseMessage{CorrelationID: PromptCorrelationTestId},
				Message:     "invalid password",
			}
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
			prompterApp.App().PromptInput <- &walletapp.StringPromptInput{
				BaseMessage: walletapp.BaseMessage{CorrelationID: PromptCorrelationTestId},
				Message:     "this is not the password",
			}
			// forward test result to test goroutine
			failRes := <-resChan

			checkResultChannel(t, failRes, false, utils.WrongPassword)

			// Send cancel to prompter app to unlock the handler
			prompterApp.App().CtrlChan <- walletapp.Cancel
		}()

		resp := deleteWallet(t, api, nickname)

		verifyStatusCode(t, resp, http.StatusUnauthorized)

		_, err = prompterApp.App().Wallet.GetAccount(nickname)
		assert.NoError(t, err)
	})

	t.Run("invalid prompt correlation id", func(t *testing.T) {
		// first send a invalid id and then the correct one
		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().PromptInput <- &walletapp.StringPromptInput{
				BaseMessage: walletapp.BaseMessage{CorrelationID: "666"},
				Message:     password,
			}
			// forward test result to test goroutine
			failRes := <-resChan

			checkResultChannel(t, failRes, false, utils.ErrWrongPromptCorrelationIdMsg)

			prompterApp.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{CorrelationID: PromptCorrelationTestId},
				Password:    password,
				Fees:        "1000",
			}

			// forward test result to test goroutine
			failRes = <-resChan
		}(testResult)

		resp := deleteWallet(t, api, nickname)
		verifyStatusCode(t, resp, http.StatusUnauthorized)
	})

	t.Run("delete success", func(t *testing.T) {
		// Send password to prompter app and wait for result
		go func() {
			prompterApp.App().PromptInput <- &walletapp.StringPromptInput{
				BaseMessage: walletapp.BaseMessage{CorrelationID: PromptCorrelationTestId},
				Message:     password,
			}
			// forward test result to test goroutine
			testResult <- (<-resChan)
		}()

		resp := deleteWallet(t, api, nickname)

		verifyStatusCode(t, resp, http.StatusNoContent)

		result := <-testResult

		checkResultChannel(t, result, true, "")

		_, err = prompterApp.App().Wallet.GetAccount(nickname)
		assert.Error(t, err, "Wallet should have been deleted")
	})
}
