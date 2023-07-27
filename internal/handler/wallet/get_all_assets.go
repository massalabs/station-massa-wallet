package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
)

func NewGetAllAssets(assetsStore *assets.AssetsStore) operations.GetAllAssetsHandler {
	return &getAllAssets{
		assetsStore: assetsStore,
	}
}

type getAllAssets struct {
	assetsStore *assets.AssetsStore
}

func (h *getAllAssets) Handle(params operations.GetAllAssetsParams) middleware.Responder {
	// Retrieve all assets from the asset store
	assets := make([]*models.AssetInfoWithBalance, 0)

	for assetAddress, assetInfo := range h.assetsStore.ContractAssets {
		// For now, mock the balance to 10 for each asset
		assetWithBalance := &models.AssetInfoWithBalance{
			AssetInfo: assetInfo,
			Balance:   getBalance(assetAddress),
		}
		assets = append(assets, assetWithBalance)
	}

	// Return the list of assets without balance
	return operations.NewGetAllAssetsOK().WithPayload(assets)
}

// getBalance is a function to get the balance for an asset.
// For now, it returns 10.0 as the balance.
func getBalance(assetAddress string) float64 {
	// For now, we are returning a mocked balance of 10.0.
	return 10.0
}
