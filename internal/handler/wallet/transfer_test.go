package wallet

import (
	"fmt"
	"net/http"
	"testing"

	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/cache"
	"github.com/massalabs/station-massa-wallet/pkg/config"
	"github.com/stretchr/testify/assert"
)

func Test_transfer_handler(t *testing.T) {
	api, resChan, err := MockAPI()
	assert.NoError(t, err)

	handler, exist := api.HandlerFor("post", "/api/accounts/{nickname}/transfer")
	assert.True(t, exist, "Endpoint doesn't exist")

	nickname := "wallet1"
	password := "password"
	createAccount(password, nickname, t, prompterAppMock)

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

	t.Run("Transfer with invalid fee", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/transfer", nickname), `{
			"fee": "none",
			"amount": "5",
			"recipientAddress": "AU1eQkRhZZBa5VNc24fCejxgFDpe1FHChpwiUksQB9StNb3rWm6i"
		}`)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
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
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    password,
				Fees:        "1000",
			}
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
		verifyStatusCode(t, resp, http.StatusOK)
	})

	// Regression: when the account has an enabled sign rule, the handler caches
	// the private key after sending the operation. The password buffer used for
	// signing is destroyed by acc.Sign, so the handler must back up the password
	// before signing. Without the backup, this case returned HTTP 500 even though
	// the operation was submitted, surfacing as a misleading error toast in the UI.
	t.Run("Transfer valid with enabled sign rule caches private key", func(t *testing.T) {
		account, err := prompterAppMock.App().Wallet.GetAccount(nickname)
		assert.NoError(t, err)

		cfg := config.Get()
		ruleID, err := cfg.AddSignRule(nickname, config.SignRule{
			Name:     "Disable password prompt",
			Contract: "AS12UMSUxgpRBB6ArZDJ19arHoxNkkpdfofQGekAiAJqsuE6PEFJy",
			RuleType: config.RuleTypeDisablePasswordPrompt,
			Enabled:  true,
		})
		assert.NoError(t, err)
		assert.True(t, cfg.HasEnabledRule(nickname))

		t.Cleanup(func() {
			_ = cfg.DeleteSignRule(nickname, ruleID)
			address, addrErr := account.Address.String()
			assert.NoError(t, addrErr)
			cache.Init().Remove(cache.KeyHash([]byte("pkey" + address)))
		})

		testResult := make(chan walletapp.EventData)

		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.SignPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Password:    password,
				Fees:        "1000",
			}

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
		verifyStatusCode(t, resp, http.StatusOK)

		// The whole point of the cache block after transfer is to populate the cache.
		// If the password backup is missing, CachePrivateKeyFromPassword fails and
		// the handler returns 500 before reaching this assertion.
		pkey, err := cache.PrivateKeyFromCache(account)
		assert.NoError(t, err)
		assert.NotNil(t, pkey)
	})
}
