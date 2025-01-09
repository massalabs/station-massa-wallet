package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/stretchr/testify/assert"
)

func Test_getWallets_handler(t *testing.T) {
	api, _, err := MockAPI()
	assert.NoError(t, err)

	wallet.ClearAccounts(t, prompterAppMock.App().Wallet.WalletPath)

	// test empty configuration first.
	t.Run("Get empty list", func(t *testing.T) {
		resp, err := processHTTPRequest(api, "GET", "/api/accounts", "")
		assert.NoError(t, err)

		assert.Equal(t, resp.Result().StatusCode, http.StatusOK)

		var accounts []models.Account
		err = json.Unmarshal(resp.Body.Bytes(), &accounts)
		assert.NoError(t, err)

		assert.Len(t, accounts, 0, "the account list should be empty")
	})

	// Create accounts
	password := "zePassword"

	nicknames := []string{"account1", "account2", "account3"}
	for _, nickname := range nicknames {
		createAccount(password, nickname, t, prompterAppMock)
	}

	t.Run("Get multiple accounts", func(t *testing.T) {
		resp, err := processHTTPRequest(api, "GET", "/api/accounts", "")
		assert.NoError(t, err)

		assert.Equal(t, resp.Result().StatusCode, http.StatusOK)

		var accounts []models.Account
		err = json.Unmarshal(resp.Body.Bytes(), &accounts)
		assert.NoError(t, err)

		for idx, nickname := range nicknames {
			assert.Equal(t, accounts[idx].Nickname, models.Nickname(nickname))
		}

		assertAccountsBody(t, resp, true)
	})

	wallet.ClearAccounts(t, prompterAppMock.App().Wallet.WalletPath)
}

func Test_getWallet_handler(t *testing.T) {
	nickname := "trololol"
	password := "zePassword"

	api, resChan, err := MockAPI()
	assert.NoError(t, err)

	createAccount(password, nickname, t, prompterAppMock)

	handler, exist := api.HandlerFor("get", "/api/accounts/{nickname}")
	assert.True(t, exist, "Endpoint doesn't exist")

	// test empty configuration first.
	t.Run("Get unknown account", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s", "nobody"), "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})

	// test with one wallet configuration.
	t.Run("Passed_get_ciphered", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s", nickname), "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		assertAccountBody(t, resp, nickname, true)

		wallet.ClearAccounts(t, prompterAppMock.App().Wallet.WalletPath)
	})

	// test with un-ciphered data.
	t.Run("Passed_get_un-ciphered", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)
		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterAppMock.App().PromptInput <- &walletapp.StringPromptInput{
				BaseMessage: walletapp.BaseMessage{},
				Message:     password,
			}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s?ciphered=false", nickname), "")
		assert.NoError(t, err)

		verifyStatusCode(t, resp, http.StatusOK)
		assertAccountBody(t, resp, nickname, false)

		result := <-testResult

		checkResultChannel(t, result, true, "")
	})
}

func assertAccountsBody(t *testing.T, resp *httptest.ResponseRecorder, cyphered bool) {
	var accounts []models.Account
	err := json.Unmarshal(resp.Body.Bytes(), &accounts)
	assert.NoError(t, err)

	for _, acc := range accounts {
		assertAccountContent(t, acc, cyphered)
	}
}

func assertAccountBody(t *testing.T, resp *httptest.ResponseRecorder, nickname string, cyphered bool) {
	var acc models.Account
	err := json.Unmarshal(resp.Body.Bytes(), &acc)
	assert.NoError(t, err)
	assert.Equal(t, models.Nickname(nickname), acc.Nickname)
	assertAccountContent(t, acc, cyphered)
}

func assertAccountContent(t *testing.T, acc models.Account, cyphered bool) {
	assert.NotEmpty(t, acc.Address)
	assert.NotEmpty(t, acc.Nickname)
	assert.NotEmpty(t, acc.CandidateBalance)
	assert.NotEmpty(t, acc.Balance)

	if cyphered {
		assert.Empty(t, acc.KeyPair)
	} else {
		assert.NotEmpty(t, acc.KeyPair)
		assert.NotEmpty(t, acc.KeyPair.Nonce)
		assert.NotEmpty(t, acc.KeyPair.PrivateKey)
		assert.NotEmpty(t, acc.KeyPair.PublicKey)
		assert.NotEmpty(t, acc.KeyPair.Salt)
	}
}
