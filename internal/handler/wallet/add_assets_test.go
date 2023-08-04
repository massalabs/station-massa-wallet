package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	utils "github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAddAssetHandler(t *testing.T) {
	api, _, _, _, err := MockAPI()
	assert.NoError(t, err)

	nickname := "GoodNickname"
	password := "zePassword"

	_, errGenerate := wallet.Generate(nickname, password)
	assert.Nil(t, errGenerate)

	t.Run("InvalidAddress", func(t *testing.T) {
		invalidAddress := "InvalidAddress"
		// Get the handler for the AddAsset endpoint
		handler, exist := api.HandlerFor("post", "/api/accounts/{nickname}/assets")
		assert.True(t, exist)

		// Send the POST request with an invalid address
		resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/assets?assetAddress=%s", nickname, invalidAddress), "")
		assert.NoError(t, err)

		// Assert that the response status code is 422 Unprocessable Entity
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Result().StatusCode)

		// Parse the error response from the body
		var errorResponse models.Error
		_ = json.Unmarshal(resp.Body.Bytes(), &errorResponse)
		// Assert that the error message matches the expected values
		assert.Equal(t, "assetAddress in query should match '^AS[0-9a-zA-Z]+$'", errorResponse.Message)

		// Remove the json file created
		err = RemoveJSONFile()
		assert.NoError(t, err)
	})

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

		// Remove the json file created
		err = RemoveJSONFile()
		assert.NoError(t, err)

		err = cleanupTestData([]string{nickname})
		assert.NoError(t, err)
	})
}

func addAssetTest(t *testing.T, api *operations.MassaWalletAPI, nickname, assetAddress string) *models.AssetInfo {
	t.Helper()

	handler, exist := api.HandlerFor("post", "/api/accounts/{nickname}/assets")
	assert.True(t, exist)

	// Send the POST request with a valid address to add a new asset
	resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/api/accounts/%s/assets?assetAddress=%s", nickname, assetAddress), "")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)

	// Parse the response body to get the added asset
	var addedAsset models.AssetInfo
	err = json.Unmarshal(resp.Body.Bytes(), &addedAsset)
	assert.NoError(t, err)

	return &addedAsset
}

// RemoveJSONFile removes the assets.json file if it exists.
func RemoveJSONFile() error {
	assetsJSONPath, err := utils.GetAssetsJSONPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(assetsJSONPath); os.IsNotExist(err) {
		// The file does not exist, so there's nothing to remove.
		return nil
	}

	// The file exists, so let's attempt to remove it.
	if err := os.Remove(assetsJSONPath); err != nil {
		return errors.Wrap(err, "failed to remove assets JSON file")
	}

	return nil
}
