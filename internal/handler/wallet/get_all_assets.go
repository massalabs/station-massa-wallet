package wallet

import (
	"fmt"

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
	wlt, resp := loadWallet(params.WalletNickname)
	if resp != nil {
		return resp
	}

	assets := make([]*models.AssetInfoWithBalance, 0)

	// Retrieve all assets from the asset store
	for assetAddress, assetInfo := range h.assetsStore.ContractAssets {
		balance, err := h.assetsStore.Balance(assetAddress, wlt.Address)
		if err != nil {
			// Handle the error and return an internal server error response
			errorMsg := fmt.Sprintf("Failed to fetch balance for asset %s: %s", assetAddress, err.Error())
			return operations.NewGetAllAssetsInternalServerError().WithPayload(&models.Error{
				Code:    errorFetchAssetBalance,
				Message: errorMsg,
			})
		}

		assetWithBalance := &models.AssetInfoWithBalance{
			AssetInfo: assetInfo,
			Balance:   balance,
		}
		assets = append(assets, assetWithBalance)
	}

	// Return the list of assets with balance
	return operations.NewGetAllAssetsOK().WithPayload(assets)
}
