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
	api, _, _, err := MockAPI()
	assert.NoError(t, err)

	// test empty configuration first.
	t.Run("Get empty list", func(t *testing.T) {
		resp, err := processHTTPRequest(api, "GET", "/api/accounts", "")
		assert.NoError(t, err)

		assert.Equal(t, resp.Result().StatusCode, http.StatusOK)

		var wallets []models.Account
		err = json.Unmarshal(resp.Body.Bytes(), &wallets)
		assert.NoError(t, err)

		assert.Len(t, wallets, 0, "the wallets list should be empty")
	})

	// Create wallets
	nicknames := []string{"wallet1", "wallet2", "wallet3"}
	for _, nickname := range nicknames {
		password := "zePassword"
		_, errGenerate := wallet.Generate(nickname, password)
		assert.Nil(t, errGenerate)
	}

	t.Run("Get multiple wallets", func(t *testing.T) {
		resp, err := processHTTPRequest(api, "GET", "/api/accounts", "")
		assert.NoError(t, err)

		assert.Equal(t, resp.Result().StatusCode, http.StatusOK)

		var wallets []models.Account
		err = json.Unmarshal(resp.Body.Bytes(), &wallets)
		assert.NoError(t, err)

		for idx, nickname := range nicknames {
			assert.Equal(t, wallets[idx].Nickname, models.Nickname(nickname))
		}

		assertWalletsBody(t, resp, true)
	})

	err = cleanupTestData(nicknames)
	assert.NoError(t, err, fmt.Sprintf("while cleaning up TestData: %s", err))
}

func Test_getWallet_handler(t *testing.T) {
	api, prompterApp, resChan, err := MockAPI()
	assert.NoError(t, err)

	handler, exist := api.HandlerFor("get", "/api/accounts/{nickname}")
	assert.True(t, exist, "Endpoint doesn't exist")

	// test empty configuration first.
	t.Run("Get unknown wallet", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s", "nobody"), "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})

	// test with one wallet configuration.
	t.Run("Passed_get_ciphered", func(t *testing.T) {
		nickname := "trololol"
		password := "zePassword"
		_, errGenerate := wallet.Generate(nickname, password)
		assert.Nil(t, errGenerate)

		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s", nickname), "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		assertWalletBody(t, resp, nickname, true)

		err = cleanupTestData([]string{nickname})
		assert.NoError(t, err, fmt.Sprintf("while cleaning up TestData: %s", err))
	})

	// test with un-ciphered data.
	t.Run("Passed_get_un-ciphered", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)
		nickname := "trololol"
		password := "zePassword"
		_, errGenerate := wallet.Generate(nickname, password)
		assert.Nil(t, errGenerate)

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().PromptInput <- password
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s?ciphered=false", nickname), "")
		assert.NoError(t, err)

		verifyStatusCode(t, resp, http.StatusOK)
		assertWalletBody(t, resp, nickname, false)

		result := <-testResult

		checkResultChannel(t, result, true, "Unprotect Success")

		err = cleanupTestData([]string{nickname})
		assert.NoError(t, err, fmt.Sprintf("while cleaning up TestData: %s", err))
	})
}

func assertWalletsBody(t *testing.T, resp *httptest.ResponseRecorder, cyphered bool) {
	var wallets []models.Account
	err := json.Unmarshal(resp.Body.Bytes(), &wallets)
	assert.NoError(t, err)

	for _, wallet := range wallets {
		assertWalletContent(t, wallet, cyphered)
	}
}

func assertWalletBody(t *testing.T, resp *httptest.ResponseRecorder, nickname string, cyphered bool) {
	var wallet models.Account
	err := json.Unmarshal(resp.Body.Bytes(), &wallet)
	assert.NoError(t, err)
	assert.Equal(t, models.Nickname(nickname), wallet.Nickname)
	assertWalletContent(t, wallet, cyphered)
}

func assertWalletContent(t *testing.T, wallet models.Account, cyphered bool) {
	assert.NotEmpty(t, wallet.Address)
	assert.NotEmpty(t, wallet.Nickname)
	assert.NotEmpty(t, wallet.CandidateBalance)
	assert.NotEmpty(t, wallet.Balance)
	if cyphered {
		assert.Empty(t, wallet.KeyPair)
	} else {
		assert.NotEmpty(t, wallet.KeyPair)
		assert.NotEmpty(t, wallet.KeyPair.Nonce)
		assert.NotEmpty(t, wallet.KeyPair.PrivateKey)
		assert.NotEmpty(t, wallet.KeyPair.PublicKey)
		assert.NotEmpty(t, wallet.KeyPair.Salt)
	}
}
