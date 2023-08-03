package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	utils "github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/stretchr/testify/assert"
)

func TestAddAssetHandler_InvalidAddress(t *testing.T) {
	api, _, _, _, err := MockAPI()
	assert.NoError(t, err)

	nickname := "GoodNickname"
	invalidAddress := "InvalidAddress"
	password := "zePassword"


	_, errGenerate := wallet.Generate(nickname, password)
	assert.Nil(t, errGenerate)

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
	// Assert that the error message match the expected values
	assert.Equal(t, "assetAddress in query should match '^AS[0-9a-zA-Z]+$'", errorResponse.Message)

	// Remove the json file created
	err = utils.RemoveJSONFile()
	assert.NoError(t, err)

	err = cleanupTestData([]string{nickname})
	assert.NoError(t, err)
}

func TestAddAssetHandler_AssetAlreadyExists(t *testing.T) {
	api, _, _, _, err := MockAPI()
	assert.NoError(t, err)

	nickname := "GoodNickname"
	existingAddress := "AS17gQyPvtwGQ2rfvE6L91J3N7ebvnvsSuh44vADVrPSFVW3vw96"

	password := "zePassword"


	_, errGenerate := wallet.Generate(nickname, password)
	assert.Nil(t, errGenerate)

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
	err = utils.RemoveJSONFile()
	assert.NoError(t, err)

	err = cleanupTestData([]string{nickname})
	assert.NoError(t, err)
}
