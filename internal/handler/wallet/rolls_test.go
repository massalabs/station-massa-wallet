package wallet

import (
	"fmt"
	"net/http"
	"testing"

	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/stretchr/testify/assert"
)

func Test_traderolls_handler(t *testing.T) {
	api, resChan, err := MockAPI()
	assert.NoError(t, err)

	handler, exist := api.HandlerFor("post", "/api/accounts/{nickname}/rolls")
	assert.True(t, exist, "Endpoint doesn't exist")

	nickname := "wallet1"
	password := "password"
	createAccount(password, nickname, t, prompterAppMock)

	t.Run("Trade rolls with unprocessable entity", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/rolls", nickname), "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)
	})

	t.Run("Trade rolls with unknown account", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/rolls", "not-a-nickname"), `{
			"fee": "1",
			"amount": "2",
			"side": "buy"
		}`)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})

	t.Run("Trade rolls with invalid trade side", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/rolls", nickname), `{
			"fee": "1",
			"amount": "2",
			"side": "lolmao"
		}`)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)
	})

	t.Run("Trade rolls with invalid amount", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/rolls", nickname), `{
			"fee": "1",
			"amount": "two",
			"side": "buy"
		}`)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})

	t.Run("Trade rolls with invalid decimal amount", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/rolls", nickname), `{
			"fee": "1",
			"amount": "6.666",
			"side": "buy"
		}`)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})

	t.Run("Trade rolls with invalid fee", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/rolls", nickname), `{
			"fee": "yo",
			"amount": "2",
			"side": "buy"
		}`)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})

	t.Run("Buy rolls valid", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    password,
				Fees:        "1000",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/rolls", nickname), `{
			"fee": "1",
			"amount": "100",
			"side": "buy"
		}`)

		result := <-testResult

		assert.True(t, result.Success)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
	})

	t.Run("Sell rolls valid", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    password,
				Fees:        "1000",
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/rolls", nickname), `{
			"fee": "1234",
			"amount": "123456789",
			"side": "sell"
		}`)

		result := <-testResult
		assert.True(t, result.Success)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
	})
}
