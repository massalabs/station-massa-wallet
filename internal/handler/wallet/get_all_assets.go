package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
)

func NewGetAllAssets(accountsStore *assets.AccountsStore) operations.GetAllAssetsHandler {
	return &getAllAssets{
		accountsStore: accountsStore,
	}
}

type getAllAssets struct {
	accountsStore *assets.AccountsStore
}

func (h *getAllAssets) Handle(params operations.GetAllAssetsParams) middleware.Responder {
	wlt, resp := loadWallet(params.WalletNickname)
	if resp != nil {
		return resp
	}

	AssetsWithBalance := make([]*models.AssetInfoWithBalance, 0)

	// Retrieve all assets from the selected account
	for assetAddress, assetInfo := range h.accountsStore.Accounts[params.WalletNickname].ContractAssets {
		balance, err := assets.Balance(assetAddress, wlt.Address)
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
		AssetsWithBalance = append(AssetsWithBalance, assetWithBalance)
	}

	// Return the list of assets with balance
	return operations.NewGetAllAssetsOK().WithPayload(AssetsWithBalance)
}
