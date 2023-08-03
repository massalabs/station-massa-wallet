package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	utils "github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/stretchr/testify/assert"
)

func TestGetAllAssetsHandler(t *testing.T) {
	api, _, _, _, err := MockAPI()
	assert.NoError(t, err)

	nickname := "GoodNickname"
	password := "zePassword"


	_, errGenerate := wallet.Generate(nickname, password)
	assert.Nil(t, errGenerate)

	// Get the handler for the GetAllAssets endpoint
	handler, exist := api.HandlerFor("get", "/api/accounts/{nickname}/assets")
	assert.True(t, exist)

	// Send the GET request to retrieve all assets for the created wallet
	resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s/assets", nickname), "")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

	// Parse the response body to get the assetsWithBalance
	var assetsWithBalance []*models.AssetInfoWithBalance
	err = json.Unmarshal(resp.Body.Bytes(), &assetsWithBalance)
	assert.NoError(t, err)

	// Assert that assetsWithBalance contains the expected data
	assert.Len(t, assetsWithBalance, 1, "the assets list should have 1 item")

	assert.Equal(t, "1000000", assetsWithBalance[0].Balance)
	assert.Equal(t, "Massa", assetsWithBalance[0].AssetInfo.Name)
	assert.Equal(t, "XMA", assetsWithBalance[0].AssetInfo.Symbol)

	// Remove the json file created
	err = utils.RemoveJSONFile()
	assert.NoError(t, err)

	err = cleanupTestData([]string{nickname})
	assert.NoError(t, err)
}

func TestAddGetDeleteAsset(t *testing.T) {
	api, _, _, _, err := MockAPI()
	assert.NoError(t, err)

	nickname := "GoodNickname"
	password := "zePassword"


	_, errGenerate := wallet.Generate(nickname, password)
	assert.Nil(t, errGenerate)

	// Create the test wallet first

	// Step 1: Get the assets before adding
	assetsBeforeAdd := getAssets(t, api, nickname)

	// Step 2: Add a new asset to the wallet
	validAddress := "AS17gQyPvtwGQ2rfvE6L91J3N7ebvnvsSuh44vADVrPSFVW3vw96"
	addedAsset := addAssetTest(t, api, nickname, validAddress)

	// Step 3: Get the assets again after adding
	assetsAfterAdd := getAssets(t, api, nickname)

	// Assert that the number of assets increased by 1 after adding
	assert.Len(t, assetsAfterAdd, len(assetsBeforeAdd)+1, "the assets list should have increased by 1 after adding")

	// Assert the first asset remains unchanged
	assertAssetInfoWithBalanceEqual(t, assetsAfterAdd[0], assetsBeforeAdd[0])

	// Assert the newly added asset
	assert.Equal(t, "MassaToken", addedAsset.Name)
	assert.Equal(t, "MST", addedAsset.Symbol)
	assert.Equal(t, int64(18), *addedAsset.Decimals)

	// Assert the balance of the newly added asset
	assert.Equal(t, "1000000", assetsAfterAdd[0].Balance)

	// Step 4: Delete the imported asset
	deleteAssetTest(t, api, nickname, validAddress)

	// Step 5: Get the assets again after deletion
	assetsAfterDelete := getAssets(t, api, nickname)

	// Assert that the number of assets is back to the original count after deletion
	assert.Len(t, assetsAfterDelete, len(assetsBeforeAdd), "the assets list should be back to the original count after deletion")

	// Assert the first asset remains unchanged after deletion
	assertAssetInfoWithBalanceEqual(t, assetsAfterDelete[0], assetsBeforeAdd[0])

	// Remove the json file created
	err = utils.RemoveJSONFile()
	assert.NoError(t, err)

	err = cleanupTestData([]string{nickname})
	assert.NoError(t, err)
}

func getAssets(t *testing.T, api *operations.MassaWalletAPI, nickname string) []*models.AssetInfoWithBalance {
	t.Helper()

	handler, exist := api.HandlerFor("get", "/api/accounts/{nickname}/assets")
	assert.True(t, exist)

	// Send the GET request to retrieve all assets for the wallet
	resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s/assets", nickname), "")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)

	// Parse the response body to get the assetsWithBalance
	var assetsWithBalance []*models.AssetInfoWithBalance
	err = json.Unmarshal(resp.Body.Bytes(), &assetsWithBalance)
	assert.NoError(t, err)

	return assetsWithBalance
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

func deleteAssetTest(t *testing.T, api *operations.MassaWalletAPI, nickname, assetAddress string) {
	t.Helper()

	handler, exist := api.HandlerFor("delete", "/api/accounts/{nickname}/assets")
	assert.True(t, exist)

	// Send the DELETE request to delete the asset
	resp, err := handleHTTPRequest(handler, "DELETE", fmt.Sprintf("/api/accounts/%s/assets?assetAddress=%s", nickname, assetAddress), "")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)
}

func assertAssetInfoWithBalanceEqual(t *testing.T, actual, expected *models.AssetInfoWithBalance) {
	assert.Equal(t, expected.Balance, actual.Balance)
	assert.Equal(t, expected.AssetInfo.Name, actual.AssetInfo.Name)
	assert.Equal(t, expected.AssetInfo.Symbol, actual.AssetInfo.Symbol)
	assert.Equal(t, expected.AssetInfo.Decimals, actual.AssetInfo.Decimals)
}
