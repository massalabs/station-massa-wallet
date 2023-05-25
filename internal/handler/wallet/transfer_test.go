package wallet

import (
	"fmt"
	"net/http"
	"testing"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
	"github.com/stretchr/testify/assert"
)

func Test_transfer_handler(t *testing.T) {
	api, prompterApp, resChan, err := MockAPI()
	assert.NoError(t, err)

	handler, exist := api.HandlerFor("post", "/api/accounts/{nickname}/transfer")
	assert.True(t, exist, "Endpoint doesn't exist")

	nickname := "wallet1"
	password := "password"
	_, errGenerate := wallet.Generate(nickname, password)
	assert.Nil(t, errGenerate)

	t.Run("Transfer with unprocessable entity", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/transfer", "nobody"), "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)
	})

	t.Run("Transfer with unknown account", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/transfer", "nobody"), `{
			"fee": "1",
			"amount": "two",
			"recipientAddress": "AU1eQkRhZZBa5VNc24fCejxgFDpe1FHChpwiUksQB9StNb3rWm6i"
		}`)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})

	t.Run("Transfer with invalid amount", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/transfer", nickname), `{
			"fee": "1",
			"amount": "two",
			"recipientAddress": "AU1eQkRhZZBa5VNc24fCejxgFDpe1FHChpwiUksQB9StNb3rWm6i"
		}`)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
	})

	t.Run("Transfer valid", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().PromptInput <- password
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/transfer", nickname), `{
			"fee": "1",
			"amount": "100",
			"recipientAddress": "AU1eQkRhZZBa5VNc24fCejxgFDpe1FHChpwiUksQB9StNb3rWm6i"
		}`)

		result := <-testResult
		assert.True(t, result.Success)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
	})

	err = cleanupTestData([]string{nickname})
	assert.NoError(t, err, fmt.Sprintf("while cleaning up TestData: %s", err))
}
