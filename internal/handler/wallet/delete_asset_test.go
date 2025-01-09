package wallet

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/stretchr/testify/assert"
)

func TestAddGetDeleteAsset(t *testing.T) {
	api, _, err := MockAPI()
	assert.NoError(t, err)

	nickname := "GoodNickname"
	password := "zePassword"
	createAccount(password, nickname, t, prompterAppMock)

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

	// Assert the newly added asset
	assert.Equal(t, validAddress, addedAsset.Address)
	assert.Equal(t, "TestToken", addedAsset.Name)
	assert.Equal(t, "TST", addedAsset.Symbol)
	assert.Equal(t, int64(9), *addedAsset.Decimals)

	// Assert the balance of the newly added asset
	assert.Equal(t, "10000", assetsAfterAdd[len(assetsAfterAdd)-1].Balance)

	// Assert the address of the newly added asset
	assert.Equal(t, validAddress, assetsAfterAdd[len(assetsAfterAdd)-1].Address)

	// Step 4: Delete the imported asset
	deleteAssetTest(t, api, nickname, validAddress)

	// Step 5: Get the assets again after deletion
	assetsAfterDelete := getAssets(t, api, nickname)

	// Assert that the number of assets is back to the original count after deletion
	assert.Len(t, assetsAfterDelete, len(assetsBeforeAdd), "the assets list should be back to the original count after deletion")

	// Assert the first asset remains unchanged after deletion
	assertAssetInfoWithBalanceEqual(t, assetsAfterDelete[0], assetsBeforeAdd[0])
}

func deleteAssetTest(t *testing.T, api *operations.MassaWalletAPI, nickname, assetAddress string) {
	t.Helper()

	handler, exist := api.HandlerFor("delete", "/api/accounts/{nickname}/assets")
	assert.True(t, exist)

	// Send the DELETE request to delete the asset
	resp, err := handleHTTPRequest(handler, "DELETE", fmt.Sprintf("/api/accounts/%s/assets?assetAddress=%s", nickname, assetAddress), "")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, resp.Result().StatusCode)
}
