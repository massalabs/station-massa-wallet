package wallet

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
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

func signMessage(t *testing.T, api *operations.MassaWalletAPI, nickname string, body string) *httptest.ResponseRecorder {
	handler, exist := api.HandlerFor("post", "/api/accounts/{nickname}/signMessage")
	assert.True(t, exist)

	resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/signMessage", nickname), body)
	assert.NoError(t, err)

	return resp
}

func TestPrepareOperation(t *testing.T) {
	acc := account.NewAccount(t)
	fees := uint64(1000)
	operationB64 := "AKT4CASAzuTNAqCNBgEAXBwUw39NBQYix8Ovph0TUiJuDDEnlFYUPgsbeMbrA4cLZm9yd2FyZEJ1cm7FAQDgfY7fLW7qpwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACoAAAAweGZDRERBRTI1MTAwNjIxYTViQzg4MTlkQzlEMzg0MjUzNEQ3QmY0NzYAAAAANQAAAEFTMTJUUm9TY01kd0xLOFlwdDZOQkFwcHl6Q0Z3N1FlRzVlM3hGdnhwQ0FuQW5ZTGZ1TVVUKgAAADB4NTM4NDRGOTU3N0MyMzM0ZTU0MUFlYzdEZjcxNzRFQ2U1ZEYxZkNmMKc2qgAAAAAA"

	operation, msgToSign, err := prepareOperation(acc, fees, operationB64, ChainIDUnitTests)
	assert.NoError(t, err)
	expected := []byte{0xe8, 0x7, 0xa4, 0xf8, 0x8, 0x4, 0x80, 0xce, 0xe4, 0xcd, 0x2, 0xa0, 0x8d, 0x6, 0x1, 0x0, 0x5c, 0x1c, 0x14, 0xc3, 0x7f, 0x4d, 0x5, 0x6, 0x22, 0xc7, 0xc3, 0xaf, 0xa6, 0x1d, 0x13, 0x52, 0x22, 0x6e, 0xc, 0x31, 0x27, 0x94, 0x56, 0x14, 0x3e, 0xb, 0x1b, 0x78, 0xc6, 0xeb, 0x3, 0x87, 0xb, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x42, 0x75, 0x72, 0x6e, 0xc5, 0x1, 0x0, 0xe0, 0x7d, 0x8e, 0xdf, 0x2d, 0x6e, 0xea, 0xa7, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2a, 0x0, 0x0, 0x0, 0x30, 0x78, 0x66, 0x43, 0x44, 0x44, 0x41, 0x45, 0x32, 0x35, 0x31, 0x30, 0x30, 0x36, 0x32, 0x31, 0x61, 0x35, 0x62, 0x43, 0x38, 0x38, 0x31, 0x39, 0x64, 0x43, 0x39, 0x44, 0x33, 0x38, 0x34, 0x32, 0x35, 0x33, 0x34, 0x44, 0x37, 0x42, 0x66, 0x34, 0x37, 0x36, 0x0, 0x0, 0x0, 0x0, 0x35, 0x0, 0x0, 0x0, 0x41, 0x53, 0x31, 0x32, 0x54, 0x52, 0x6f, 0x53, 0x63, 0x4d, 0x64, 0x77, 0x4c, 0x4b, 0x38, 0x59, 0x70, 0x74, 0x36, 0x4e, 0x42, 0x41, 0x70, 0x70, 0x79, 0x7a, 0x43, 0x46, 0x77, 0x37, 0x51, 0x65, 0x47, 0x35, 0x65, 0x33, 0x78, 0x46, 0x76, 0x78, 0x70, 0x43, 0x41, 0x6e, 0x41, 0x6e, 0x59, 0x4c, 0x66, 0x75, 0x4d, 0x55, 0x54, 0x2a, 0x0, 0x0, 0x0, 0x30, 0x78, 0x35, 0x33, 0x38, 0x34, 0x34, 0x46, 0x39, 0x35, 0x37, 0x37, 0x43, 0x32, 0x33, 0x33, 0x34, 0x65, 0x35, 0x34, 0x31, 0x41, 0x65, 0x63, 0x37, 0x44, 0x66, 0x37, 0x31, 0x37, 0x34, 0x45, 0x43, 0x65, 0x35, 0x64, 0x46, 0x31, 0x66, 0x43, 0x66, 0x30, 0xa7, 0x36, 0xaa, 0x0, 0x0, 0x0, 0x0, 0x0}
	assert.Equal(t, expected, operation)
	expectedMsg := []byte{0x0, 0x0, 0x0, 0x0, 0x4, 0xa0, 0xf8, 0xfe, 0x0, 0x2d, 0x96, 0xbc, 0xda, 0xcb, 0xbe, 0x41, 0x38, 0x2c, 0xa2, 0x3e, 0x52, 0xe3, 0xd2, 0x19, 0x6c, 0xba, 0x65, 0xe7, 0xa1, 0xac, 0xd2, 0x9, 0xdf, 0xc9, 0x5c, 0x6b, 0x32, 0xb6, 0xa1, 0x8a, 0x93, 0xe8, 0x7, 0xa4, 0xf8, 0x8, 0x4, 0x80, 0xce, 0xe4, 0xcd, 0x2, 0xa0, 0x8d, 0x6, 0x1, 0x0, 0x5c, 0x1c, 0x14, 0xc3, 0x7f, 0x4d, 0x5, 0x6, 0x22, 0xc7, 0xc3, 0xaf, 0xa6, 0x1d, 0x13, 0x52, 0x22, 0x6e, 0xc, 0x31, 0x27, 0x94, 0x56, 0x14, 0x3e, 0xb, 0x1b, 0x78, 0xc6, 0xeb, 0x3, 0x87, 0xb, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x42, 0x75, 0x72, 0x6e, 0xc5, 0x1, 0x0, 0xe0, 0x7d, 0x8e, 0xdf, 0x2d, 0x6e, 0xea, 0xa7, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2a, 0x0, 0x0, 0x0, 0x30, 0x78, 0x66, 0x43, 0x44, 0x44, 0x41, 0x45, 0x32, 0x35, 0x31, 0x30, 0x30, 0x36, 0x32, 0x31, 0x61, 0x35, 0x62, 0x43, 0x38, 0x38, 0x31, 0x39, 0x64, 0x43, 0x39, 0x44, 0x33, 0x38, 0x34, 0x32, 0x35, 0x33, 0x34, 0x44, 0x37, 0x42, 0x66, 0x34, 0x37, 0x36, 0x0, 0x0, 0x0, 0x0, 0x35, 0x0, 0x0, 0x0, 0x41, 0x53, 0x31, 0x32, 0x54, 0x52, 0x6f, 0x53, 0x63, 0x4d, 0x64, 0x77, 0x4c, 0x4b, 0x38, 0x59, 0x70, 0x74, 0x36, 0x4e, 0x42, 0x41, 0x70, 0x70, 0x79, 0x7a, 0x43, 0x46, 0x77, 0x37, 0x51, 0x65, 0x47, 0x35, 0x65, 0x33, 0x78, 0x46, 0x76, 0x78, 0x70, 0x43, 0x41, 0x6e, 0x41, 0x6e, 0x59, 0x4c, 0x66, 0x75, 0x4d, 0x55, 0x54, 0x2a, 0x0, 0x0, 0x0, 0x30, 0x78, 0x35, 0x33, 0x38, 0x34, 0x34, 0x46, 0x39, 0x35, 0x37, 0x37, 0x43, 0x32, 0x33, 0x33, 0x34, 0x65, 0x35, 0x34, 0x31, 0x41, 0x65, 0x63, 0x37, 0x44, 0x66, 0x37, 0x31, 0x37, 0x34, 0x45, 0x43, 0x65, 0x35, 0x64, 0x46, 0x31, 0x66, 0x43, 0x66, 0x30, 0xa7, 0x36, 0xaa, 0x0, 0x0, 0x0, 0x0, 0x0}
	assert.Equal(t, expectedMsg, msgToSign)
}

func Test_walletSign_Handle(t *testing.T) {
	api, prompterApp, _, resChan, err := MockAPI()
	assert.NoError(t, err)

	transactionData := fmt.Sprintf(`{"chainId": `+strconv.FormatUint(ChainIDUnitTests, 10)+`, "operation":"%s"}`, callSCString)
	nickname := "walletToDelete"
	password := "zePassword"
	createAccount(password, nickname, t, prompterApp)

	t.Run("invalid nickname", func(t *testing.T) {
		resp := signTransaction(t, api, "Johnny", transactionData)
		verifyStatusCode(t, resp, http.StatusNotFound)
	})

	t.Run("sign transaction OK", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{CorrelationID: PromptCorrelationTestId},
				Password:    password,
				Fees:        "14400",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "")
	})

	t.Run("sign a plain text message", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{CorrelationID: PromptCorrelationTestId},
				Password:    password,
				Fees:        "1000",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		message := "Test"
		requestBody := fmt.Sprintf(`{"message":"%s"}`, message)

		resp := signMessage(t, api, nickname, requestBody)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "")
	})

	// The handler will not return until a the good password is sent or the action is canceled
	t.Run("invalid password try, then valid password", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		//nolint:staticcheck
		go func(res chan walletapp.EventData) {
			// Send wrong password to prompter app and wait for result
			prompterApp.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{CorrelationID: PromptCorrelationTestId},
				Password:    "this is not the password",
				Fees:        "1000",
			}
			// forward test result to test goroutine
			failRes := <-resChan

			checkResultChannel(t, failRes, false, utils.WrongPassword)

			// Send password to prompter app to unlock the handler
			prompterApp.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{CorrelationID: PromptCorrelationTestId},
				Password:    password,
				Fees:        "1000",
			}

			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "")
	})

	t.Run("invalid password try, then action canceled by user", func(t *testing.T) {
		//nolint:staticcheck
		go func() {
			// Send wrong password to prompter app and wait for result
			prompterApp.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{CorrelationID: PromptCorrelationTestId},
				Password:    "this is not the password",
				Fees:        "1000",
			}
			// forward test result to test goroutine
			failRes := <-resChan

			checkResultChannel(t, failRes, false, utils.WrongPassword)

			// Send cancel to prompter app to unlock the handler
			prompterApp.App().CtrlChan <- walletapp.Cancel
		}()

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, http.StatusUnauthorized)
	})

	t.Run("sign transaction batch", func(t *testing.T) {
		transactionDataBatch := fmt.Sprintf(`{"chainId": `+strconv.FormatUint(ChainIDUnitTests, 10)+`, "operation":"%s","batch":true}`, callSCString)
		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{CorrelationID: PromptCorrelationTestId},
				Password:    password,
				Fees:        "1000",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionDataBatch)
		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "")

		var body models.SignResponse
		err = json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)

		correlationId := base64.StdEncoding.EncodeToString(body.CorrelationID)

		transactionDataBatch = fmt.Sprintf(`{"chainId": `+strconv.FormatUint(ChainIDUnitTests, 10)+`, "operation":"%s","correlationId":"%s"}`, callSCString, correlationId)
		// Send new transaction without password prompt
		// ---
		// No, we disable batch signing, so send password prompt:
		go func(res chan walletapp.EventData) {
			prompterApp.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{CorrelationID: PromptCorrelationTestId},
				Password:    password,
				Fees:        "1000",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)
		// ---
		resp = signTransaction(t, api, nickname, transactionDataBatch)
		verifyStatusCode(t, resp, http.StatusOK)

		// Send new transaction with incorrect correlation id
		correlationId = base64.StdEncoding.EncodeToString([]byte("wrong correlation id"))
		transactionDataBatch = fmt.Sprintf(`{"chainId": `+strconv.FormatUint(ChainIDUnitTests, 10)+`, "operation":"%s","correlationId":"%s"}`, callSCString, correlationId)

		// We disable batch signing, so send password prompt:
		go func(res chan walletapp.EventData) {
			prompterApp.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{CorrelationID: PromptCorrelationTestId},
				Password:    password,
				Fees:        "1000",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp = signTransaction(t, api, nickname, transactionDataBatch)

		verifyStatusCode(t, resp, http.StatusOK) // remove this line when enabling batch signing
		// Uncomment the code below when enabling batch signing:
		// var bodyError operations.SignInternalServerError
		// err = json.Unmarshal(resp.Body.Bytes(), &bodyError)
		// assert.NoError(t, err)
		// verifyStatusCode(t, resp, http.StatusNotFound)
	})
}
