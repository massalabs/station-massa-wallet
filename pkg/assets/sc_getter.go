package assets

import (
	"github.com/go-openapi/swag"
	"github.com/massalabs/station-massa-wallet/api/server/models"
)

// getAssetInfoFromSC retrieves the asset information for a given contract address by making a smart contract call.
func (s *AssetsStore) getAssetInfoFromSC(contractAddress string) (*models.AssetInfo, error) {
	// mocking for now implementation to call the smart contract and retrieve asset information goes here.
	// todo : Replace the following dummy assetInfo with the actual asset information retrieved from the smart contract.

	assetInfoFromSC := &models.AssetInfo{
		Name:     "Token C",
		Symbol:   "TKC",
		Decimals: swag.Int64(6),
	}

	return assetInfoFromSC, nil
}
