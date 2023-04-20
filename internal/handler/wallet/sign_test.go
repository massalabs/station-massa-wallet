package wallet

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func signTransaction(t *testing.T, api *operations.MassaWalletAPI, nickname string, body string) *httptest.ResponseRecorder {
	handler, exist := api.HandlerFor("post", "/rest/wallet/{nickname}/signOperation")
	if !exist {
		t.Fatalf("Endpoint doesn't exist")
	}

	resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/rest/wallet/%s/signOperation", nickname), body)
	if err != nil {
		t.Fatalf("while serving HTTP request: %s", err)
	}
	return resp
}

func Test_walletSign_Handle(t *testing.T) {
	api, prompterApp, resChan, err := MockAPI()
	if err != nil {
		panic(err)
	}

	transactionData := `{"operation":"MjIzM3QyNHQ="}`
	nickname := "walletToDelete"
	password := "zePassword"
	_, err = wallet.Generate(nickname, password)
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Run("invalid nickname", func(t *testing.T) {
		resp := signTransaction(t, api, "Johnny", transactionData)
		verifyStatusCode(t, resp, 500)
	})

	t.Run("sign transaction OK", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().PasswordChan <- password
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, 200)

		result := <-testResult

		if !result.Success {
			t.Fatalf("Expected success, got error")
		}

		msg := "Unprotect Success"
		if result.Data != msg {
			t.Fatalf(fmt.Sprintf("Expected error message to be %s, got %s", msg, result.Data))
		}
	})

	// The handler will not return until a the good password is sent or the action is canceled
	t.Run("invalid password try, then valid password", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		//nolint:staticcheck
		go func(res chan walletapp.EventData) {
			// Send wrong password to prompter app and wait for result
			prompterApp.App().PasswordChan <- "this is not the password"
			// forward test result to test goroutine
			failRes := <-resChan
			if failRes.Success {
				//nolint:govet
				t.Fatalf("Expected error, got success")
			}

			msg := "error unprotecting wallet:opening the private key seal: cipher: message authentication failed"
			if failRes.Data != msg {
				//nolint:govet
				t.Fatalf(fmt.Sprintf("Expected error message to be %s, got %s", msg, failRes.Data))
			}

			// Send password to prompter app to unlock the handler
			prompterApp.App().PasswordChan <- password

			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, 200)

		result := <-testResult

		if !result.Success {
			t.Fatalf("Expected success, got error")
		}

		msg := "Unprotect Success"
		if result.Data != msg {
			t.Fatalf(fmt.Sprintf("Expected error message to be %s, got %s", msg, result.Data))
		}
	})

	t.Run("invalid password try, then action canceled by user", func(t *testing.T) {
		//nolint:staticcheck
		go func() {
			// Send wrong password to prompter app and wait for result
			prompterApp.App().PasswordChan <- "this is not the password"
			// forward test result to test goroutine
			failRes := <-resChan
			if failRes.Success {
				//nolint:govet
				t.Fatalf("Expected error, got success")
			}

			msg := "error unprotecting wallet:opening the private key seal: cipher: message authentication failed"
			if failRes.Data != msg {
				//nolint:govet
				t.Fatalf(fmt.Sprintf("Expected error message to be %s, got %s", msg, failRes.Data))
			}

			// Send cancel to prompter app to unlock the handler
			prompterApp.App().CtrlChan <- walletapp.Cancel
		}()

		resp := signTransaction(t, api, nickname, transactionData)
		verifyStatusCode(t, resp, 500)
	})

	t.Run("sign transation batch OK", func(t *testing.T) {
		transactionDataBatch := `{"operation":"MjIzM3QyNHQ=","batch":true}`
		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().PasswordChan <- password
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := signTransaction(t, api, nickname, transactionDataBatch)
		verifyStatusCode(t, resp, 200)

		result := <-testResult

		if !result.Success {
			t.Fatalf("Expected success, got error")
		}

		msg := "Unprotect Success"
		if result.Data != msg {
			t.Fatalf(fmt.Sprintf("Expected error message to be %s, got %s", msg, result.Data))
		}

		var body models.Signature
		err = json.Unmarshal(resp.Body.Bytes(), &body)
		if err != nil {
			t.Fatalf("while unmarshalling: %s", err)
		}

		correlationId := base64.StdEncoding.EncodeToString(body.CorrelationID)

		transactionDataBatch = fmt.Sprintf(`{"operation":"MjIzM3QyNHQ=","correlationId":"%s"}`, correlationId)
		// Send new transaction without password prompt
		resp = signTransaction(t, api, nickname, transactionDataBatch)
		verifyStatusCode(t, resp, 200)
	})

	err = cleanupTestData([]string{nickname})
	if err != nil {
		t.Fatalf("while cleaning up TestData: %s", err)
	}
}
