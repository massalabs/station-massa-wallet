package wallet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/stretchr/testify/assert"
)

func backupWallet(t *testing.T, api *operations.MassaWalletAPI, nickname string) *httptest.ResponseRecorder {
	handler, exist := api.HandlerFor("POST", "/api/accounts/{nickname}/backup")
	assert.True(t, exist, "Endpoint doesn't exist")

	resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/backup", nickname), "")
	assert.NoError(t, err)
	return resp
}

func Test_walletBackupAccount_Handle(t *testing.T) {
	api, prompterApp, resChan, err := MockAPI()
	assert.NoError(t, err)

	nickname := "walletToBackup"
	password := "zePassword"
	wlt, walletError := wallet.Generate(nickname, password)
	assert.Nil(t, walletError)

	walletError = wlt.Unprotect(password)
	assert.Nil(t, walletError)

	t.Run("invalid nickname", func(t *testing.T) {
		resp := backupWallet(t, api, "toto")
		verifyStatusCode(t, resp, http.StatusNotFound)
	})

	t.Run("invalid backup method", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		go func() {
			prompterApp.App().PromptInput <- "unknown method"
			testResult <- (<-resChan)
		}()

		resp := backupWallet(t, api, nickname)
		verifyStatusCode(t, resp, http.StatusUnauthorized)

		result := <-testResult

		checkResultChannel(t, result, false, utils.ErrPromptInputType)
	})

	t.Run("export canceled by user", func(t *testing.T) {
		go func() {
			prompterApp.App().CtrlChan <- walletapp.Cancel
		}()

		resp := backupWallet(t, api, nickname)
		verifyStatusCode(t, resp, http.StatusUnauthorized)
	})

	t.Run("export yaml file", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		go func() {
			prompterApp.App().PromptInput <- string(prompt.YamlFileBackup)
			testResult <- (<-resChan)
		}()

		resp := backupWallet(t, api, nickname)
		verifyStatusCode(t, resp, http.StatusNoContent)

		result := <-testResult

		checkResultChannel(t, result, true, "Backup Success")
	})

	t.Run("chose private backup then cancel", func(t *testing.T) {
		go func() {
			// send backup method
			prompterApp.App().PromptInput <- string(prompt.PrivateKeyBackup)
			prompterApp.App().CtrlChan <- walletapp.Cancel
		}()

		resp := backupWallet(t, api, nickname)
		verifyStatusCode(t, resp, http.StatusUnauthorized)
	})

	t.Run("backup private key, wrong password and cancel", func(t *testing.T) {
		go func() {
			// send backup method
			prompterApp.App().PromptInput <- string(prompt.PrivateKeyBackup)
			// send password
			prompterApp.App().PromptInput <- "wrong password"

			result := <-resChan

			checkResultChannel(t, result, false, utils.WrongPassword)

			prompterApp.App().CtrlChan <- walletapp.Cancel
		}()

		resp := backupWallet(t, api, nickname)
		verifyStatusCode(t, resp, http.StatusUnauthorized)
	})

	t.Run("backup private key success", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)

		go func() {
			// send backup method
			prompterApp.App().PromptInput <- string(prompt.PrivateKeyBackup)
			// send password
			prompterApp.App().PromptInput <- password

			testResult <- (<-resChan)
		}()

		resp := backupWallet(t, api, nickname)
		verifyStatusCode(t, resp, http.StatusNoContent)

		result := <-testResult

		data, ok := result.Data.(KeyPair)
		assert.True(t, ok, "Data is not of expected type")

		assert.Equal(t, wlt.GetPrivKey(), data.PrivateKey)
		assert.Equal(t, wlt.GetPupKey(), data.PublicKey)

		checkResultChannel(t, result, true, "Backup Success")
	})

	err = cleanupTestData([]string{nickname})
	assert.NoError(t, err)

	err = os.Remove(WalletBackupFilepath)
	assert.NoError(t, err)
}
