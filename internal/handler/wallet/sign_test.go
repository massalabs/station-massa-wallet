package wallet

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/stretchr/testify/assert"
)

func signTransaction(t *testing.T, api *operations.MassaWalletAPI, nickname string, body string) *httptest.ResponseRecorder {
	handler, exist := api.HandlerFor("post", "/api/accounts/{nickname}/sign")
	assert.True(t, exist)

	resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/sign", nickname), body)
	assert.NoError(t, err)
	return resp
}

func Test_walletSign_Handle(t *testing.T) {
	api, prompterApp, _, resChan, err := MockAPI()
	assert.NoError(t, err)

	transactionData := `{"operation":"MjIzM3QyNHQ="}`
	nickname := "walletToDelete"
	password := "zePassword"
	_, errGenerate := wallet.Generate(nickname, password)
	assert.Nil(t, errGenerate)

	t.Run("invalid nickname", func(t *testing.T) {
		resp := signTransaction(t, api, "Johnny", transactionData)
		verifyStatusCode(t, resp, http.StatusNotFound)
	})

	t.Run("sign transaction OK", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().PromptInput <- password
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, utils.MsgAccountUnprotected)
	})

	// The handler will not return until a the good password is sent or the action is canceled
	t.Run("invalid password try, then valid password", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		//nolint:staticcheck
		go func(res chan walletapp.EventData) {
			// Send wrong password to prompter app and wait for result
			prompterApp.App().PromptInput <- "this is not the password"
			// forward test result to test goroutine
			failRes := <-resChan

			checkResultChannel(t, failRes, false, utils.WrongPassword)

			// Send password to prompter app to unlock the handler
			prompterApp.App().PromptInput <- password

			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, utils.MsgAccountUnprotected)
	})

	t.Run("invalid password try, then action canceled by user", func(t *testing.T) {
		//nolint:staticcheck
		go func() {
			// Send wrong password to prompter app and wait for result
			prompterApp.App().PromptInput <- "this is not the password"
			// forward test result to test goroutine
			failRes := <-resChan

			checkResultChannel(t, failRes, false, utils.WrongPassword)

			// Send cancel to prompter app to unlock the handler
			prompterApp.App().CtrlChan <- walletapp.Cancel
		}()

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusUnauthorized)
	})

	t.Run("sign transaction batch OK", func(t *testing.T) {
		transactionDataBatch := `{"operation":"MjIzM3QyNHQ=","batch":true}`
		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().PromptInput <- password
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionDataBatch)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "Unprotect Success")

		var body models.SignResponse
		err = json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)

		correlationId := base64.StdEncoding.EncodeToString(body.CorrelationID)

		transactionDataBatch = fmt.Sprintf(`{"operation":"MjIzM3QyNHQ=","correlationId":"%s"}`, correlationId)
		// Send new transaction without password prompt
		resp = signTransaction(t, api, nickname, transactionDataBatch)
		verifyStatusCode(t, resp, http.StatusOK)
	})

	err = cleanupTestData([]string{nickname})
	assert.NoError(t, err)
}
