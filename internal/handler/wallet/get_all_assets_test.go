package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/stretchr/testify/assert"
)

func TestGetAllAssetsHandler(t *testing.T) {
	api, prompterApp, _, _, err := MockAPI()
	assert.NoError(t, err)

	nickname := "GoodNickname"
	password := "zePassword"
	createAccount(password, nickname, t, prompterApp)

	// Get the assetsWithBalance
	assetsWithBalance := getAssets(t, api, nickname)

	// Assert that assetsWithBalance contains the expected data
	counter := getExpectedAssetsCount(t)
	assert.Len(t, assetsWithBalance, counter+1) // +1 for native MAS

	assert.Equal(t, "1000000", assetsWithBalance[0].Balance)
	assert.Equal(t, "Massa", assetsWithBalance[0].AssetInfo.Name)
	assert.Equal(t, "MAS", assetsWithBalance[0].AssetInfo.Symbol)
}

func getAssets(t *testing.T, api *operations.MassaWalletAPI, nickname string) []*models.AssetInfoWithBalance {
	t.Helper()

	handler, exist := api.HandlerFor("get", "/api/accounts/{nickname}/assets")
	assert.True(t, exist)

	// Send the GET request to retrieve all assets for the wallet
	resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s/assets", nickname), "")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.Result().StatusCode, "response is %s", resp.Body.String())

	// Parse the response body to get the assetsWithBalance
	var assetsWithBalance []*models.AssetInfoWithBalance
	err = json.Unmarshal(resp.Body.Bytes(), &assetsWithBalance)
	assert.NoError(t, err)

	return assetsWithBalance
}

func assertAssetInfoWithBalanceEqual(t *testing.T, actual, expected *models.AssetInfoWithBalance) {
	assert.Equal(t, expected.Balance, actual.Balance)
	assert.Equal(t, expected.AssetInfo.Name, actual.AssetInfo.Name)
	assert.Equal(t, expected.AssetInfo.Symbol, actual.AssetInfo.Symbol)
	assert.Equal(t, expected.AssetInfo.Decimals, actual.AssetInfo.Decimals)
	assert.Equal(t, expected.AssetInfo.ChainID, actual.AssetInfo.ChainID)
}

func getExpectedAssetsCount(t *testing.T) int {
	tempDir, err := os.MkdirTemp(os.TempDir(), "*-wallet-dir")
	assert.NoError(t, err)
	nodeFetcher := network.NewNodeFetcher()
	store, err := assets.NewAssetsStore(tempDir, nodeFetcher)
	assert.NoError(t, err)
	defaultAssets, err := store.Default()
	assert.NoError(t, err)
	counter := 0

	for _, asset := range defaultAssets {
		if asset.ChainID == 77658377 {
			counter += 1
		}
	}

	return counter
}
