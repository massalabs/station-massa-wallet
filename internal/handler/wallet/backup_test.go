package wallet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
	"github.com/stretchr/testify/assert"
)

func backupWallet(t *testing.T, api *operations.MassaWalletAPI, nickname string) *httptest.ResponseRecorder {
	handler, exist := api.HandlerFor("GET", "/api/accounts/{nickname}/backup")
	assert.True(t, exist, "Endpoint doesn't exist")

	resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s/backup", nickname), "")
	assert.NoError(t, err)
	return resp
}

func Test_walletBackupAccount_Handle(t *testing.T) {
	api, prompterApp, resChan, err := MockAPI()
	assert.NoError(t, err)

	nickname := "walletToBackup"
	password := "zePassword"
	_, err = wallet.Generate(nickname, password)
	assert.NoError(t, err)

	t.Run("invalid nickname", func(t *testing.T) {
		resp := backupWallet(t, api, "toto")
		verifyStatusCode(t, resp, http.StatusNotFound)
	})

	t.Run("export success", func(t *testing.T) {
		resp := backupWallet(t, api, nickname)
		verifyStatusCode(t, resp, http.StatusNoContent)

		prompterApp.App().PromptInput <- password

		result := <-resChan

		checkResultChannel(t, result, true, "Backup Success")
	})
	err = cleanupTestData([]string{nickname})
	assert.NoError(t, err)
}
