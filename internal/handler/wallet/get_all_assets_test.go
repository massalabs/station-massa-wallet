package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
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
	assert.Len(t, assetsWithBalance, 1, "the assets list should have 1 item")

	assert.Equal(t, "1000000", assetsWithBalance[0].Balance)
	assert.Equal(t, "Massa", assetsWithBalance[0].AssetInfo.Name)
	assert.Equal(t, "MAS", assetsWithBalance[0].AssetInfo.Symbol)

	// Remove the json file created
	err = RemoveJSONFile()
	assert.NoError(t, err)
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
}
