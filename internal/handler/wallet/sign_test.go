package wallet

import (
	"encoding/base64"
	"encoding/binary"
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

const callSCString = "AKT4CASAzuTNAqCNBgEAXBwUw39NBQYix8Ovph0TUiJuDDEnlFYUPgsbeMbrA4cLZm9yd2FyZEJ1cm7FAQDgfY7fLW7qpwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACoAAAAweGZDRERBRTI1MTAwNjIxYTViQzg4MTlkQzlEMzg0MjUzNEQ3QmY0NzYAAAAANQAAAEFTMTJUUm9TY01kd0xLOFlwdDZOQkFwcHl6Q0Z3N1FlRzVlM3hGdnhwQ0FuQW5ZTGZ1TVVUKgAAADB4NTM4NDRGOTU3N0MyMzM0ZTU0MUFlYzdEZjcxNzRFQ2U1ZEYxZkNmMKc2qgAAAAAA"

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

	transactionData := fmt.Sprintf(`{"operation":"%s"}`, callSCString)
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

	t.Run("sign a plain text message", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().PromptInput <- password
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		message := "a message"

		// Encode the operation type as binary Uvarint
		opIDBuf := make([]byte, binary.MaxVarintLen64)
		nbBytes := binary.PutUvarint(opIDBuf, uint64(SignPlainTextMessage))

		// Encode the operation ID in base64
		encodedOpID := base64.StdEncoding.EncodeToString(opIDBuf[:nbBytes])

		// Encode the message in base64
		encodedMessage := base64.StdEncoding.EncodeToString([]byte(message))

		// Combine the encoded operation ID and the encoded message
		finalMessage := fmt.Sprintf(`{"operation":"%s %s"}`, encodedOpID, encodedMessage)

		transactionData := finalMessage

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
		transactionDataBatch := fmt.Sprintf(`{"operation":"%s","batch":true}`, callSCString)
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

		transactionDataBatch = fmt.Sprintf(`{"operation":"%s","correlationId":"%s"}`, callSCString, correlationId)
		// Send new transaction without password prompt
		resp = signTransaction(t, api, nickname, transactionDataBatch)
		verifyStatusCode(t, resp, http.StatusOK)
	})

	err = cleanupTestData([]string{nickname})
	assert.NoError(t, err)
}
