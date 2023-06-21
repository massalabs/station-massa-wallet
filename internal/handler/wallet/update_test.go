package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"

	"github.com/stretchr/testify/assert"
)

func Test_ModifyWallets_handler(t *testing.T) {
	api, _, _, err := MockAPI()
	assert.NoError(t, err)

	// Create account
	nickname := "trololol-old"
	password := "zePassword"
	_, errGenerate := wallet.Generate(nickname, password)
	assert.Nil(t, errGenerate)

	handler, exist := api.HandlerFor("put", "/api/accounts/{nickname}")
	assert.True(t, exist, "Endpoint doesn't exist")

	newNickname := "toby-new"
	payload := fmt.Sprintf(`{
		"nickname": "%s"
	}`, newNickname)

	t.Run("update invalid payload", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "PUT", fmt.Sprintf("/api/accounts/%s", nickname), "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)
	})

	t.Run("update unknown account", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "PUT", fmt.Sprintf("/api/accounts/%s", "bhabhabha"), payload)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
	})

	t.Run("update account", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "PUT", fmt.Sprintf("/api/accounts/%s", nickname), payload)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

		assertWalletBody(t, resp, newNickname, true)

		// Check that the old nickname doesn't exist anymore
		resp, err = processHTTPRequest(api, "GET", "/api/accounts", "")
		assert.NoError(t, err)
		assert.Equal(t, resp.Result().StatusCode, http.StatusOK)

		var wallets []models.Account
		err = json.Unmarshal(resp.Body.Bytes(), &wallets)
		assert.NoError(t, err)

		assert.Len(t, wallets, 1, "there should be only one wallet")
		assert.Equal(t, wallets[0].Nickname, models.Nickname(newNickname), "the nickname should have been updated")

		// Clean account file
		err = cleanupTestData([]string{newNickname})
		assert.NoError(t, err, fmt.Sprintf("while cleaning up TestData: %s", err))
	})
}
