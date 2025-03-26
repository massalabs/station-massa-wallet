package wallet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_exportFileWallet_handler(t *testing.T) {
	api, _, err := MockAPI()

	assert.NoError(t, err)

	handler, exist := api.HandlerFor("get", "/api/accounts/{nickname}/exportFile")
	assert.True(t, exist)
	// test empty configuration first.
	t.Run("Export file of unknown wallet", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s/exportFile", "nobody"), "")
		assert.NoError(t, err)

		verifyStatusCode(t, resp, http.StatusNotFound)
	})

	// test with one wallet configuration.
	t.Run("Export file of created wallet", func(t *testing.T) {
		nickname := "trololol"
		password := "zePassword"
		createAccount(password, nickname, t, prompterAppMock)

		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s/exportFile", nickname), "")
		assert.NoError(t, err)

		verifyStatusCode(t, resp, http.StatusOK)
		verifyHeader(t, resp, "Content-Type", "application/octet-stream")
		verifyHeader(t, resp, "Content-Disposition", fmt.Sprintf("attachment; filename=%q", "wallet_trololol.yaml"))
		verifyBodyWalletBackup(t, resp, nickname)
	})
}

func verifyHeader(t *testing.T, resp *httptest.ResponseRecorder, headerName, headerValue string) {
	assert.Equal(t, resp.Header().Get(headerName), headerValue)
}

func verifyBodyWalletBackup(t *testing.T, resp *httptest.ResponseRecorder, nickname string) {
	body := resp.Body.String()
	assert.NotEmpty(t, body)

	// check the first line
	assert.Contains(t, body, "Nickname: "+nickname)
}
