package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/awnumar/memguard"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/stretchr/testify/assert"
)

func TestAddAssetHandler(t *testing.T) {
	api, _, err := MockAPI()
	assert.NoError(t, err)

	nickname := "GoodNickname"
	password := "zePassword"

	_, err = account.Generate(memguard.NewBufferFromBytes([]byte(password)), nickname)
	assert.Nil(t, err)

	type testCase struct {
		Name            string
		InvalidAddress  string
		ExpectedCode    int
		ExpectedMessage string
	}

	testCases := []testCase{
		{
			Name:            "InvalidAddress1",
			InvalidAddress:  "InvalidAddress1",
			ExpectedCode:    http.StatusUnprocessableEntity,
			ExpectedMessage: "assetAddress in query should match '^AS[0-9a-zA-Z]+$'",
		},
		{
			Name:            "InvalidAddress2",
			InvalidAddress:  "AS12GwD3UEk2BP1zMx2zSdvKov97z8gs1MtsoN4u4C9emLBbhYa3U",
			ExpectedCode:    http.StatusNotFound,
			ExpectedMessage: "Asset with the provided address not found in the network.",
		},
		// Add more test cases here as needed
	}

	// Get the handler for the AddAsset endpoint
	handler, exist := api.HandlerFor("post", "/api/accounts/{nickname}/assets")
	assert.True(t, exist)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Send the POST request with the current invalid address
			resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/assets?assetAddress=%s", nickname, tc.InvalidAddress), "")
			assert.NoError(t, err)

			// Assert that the response status code matches the expected code
			assert.Equal(t, tc.ExpectedCode, resp.Result().StatusCode)

			// Parse the error response from the body
			var errorResponse models.Error
			_ = json.Unmarshal(resp.Body.Bytes(), &errorResponse)

			// Assert that the error message matches the expected message
			assert.Equal(t, tc.ExpectedMessage, errorResponse.Message)
		})
	}

	t.Run("AssetAlreadyExists", func(t *testing.T) {
		existingAddress := "AS17gQyPvtwGQ2rfvE6L91J3N7ebvnvsSuh44vADVrPSFVW3vw96"

		// Get the handler for the AddAsset endpoint
		handler, exist := api.HandlerFor("post", "/api/accounts/{nickname}/assets")
		assert.True(t, exist)

		// Add the asset with the existing address for the first time
		_ = addAssetTest(t, api, nickname, existingAddress)

		// Attempt to add the asset with the existing address again
		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/assets?assetAddress=%s", nickname, existingAddress), "")
		assert.NoError(t, err)

		// Assert that the response status code is 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)

		// Parse the error response from the body
		var errorResponse models.Error
		err = json.Unmarshal(resp.Body.Bytes(), &errorResponse)
		assert.NoError(t, err)

		// Assert that the error message matches the expected value
		assert.Equal(t, "Asset with the provided address already exists.", errorResponse.Message)
	})
}

func addAssetTest(t *testing.T, api *operations.MassaWalletAPI, nickname, assetAddress string) *models.AssetInfo {
	t.Helper()

	handler, exist := api.HandlerFor("post", "/api/accounts/{nickname}/assets")
	assert.True(t, exist)

	// Send the POST request with a valid address to add a new asset
	resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/assets?assetAddress=%s", nickname, assetAddress), "")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.Result().StatusCode, "response is %s", resp.Body.String())

	// Parse the response body to get the added asset
	var addedAsset models.AssetInfo
	err = json.Unmarshal(resp.Body.Bytes(), &addedAsset)
	assert.NoError(t, err)

	return &addedAsset
}
