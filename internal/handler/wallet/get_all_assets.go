package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
)

func NewGetAllAssets(AssetsStore *assets.AssetsStore, massaClient network.NodeFetcherInterface) operations.GetAllAssetsHandler {
	return &getAllAssets{
		AssetsStore: AssetsStore,
		massaClient: massaClient,
	}
}

type getAllAssets struct {
	AssetsStore *assets.AssetsStore
	massaClient network.NodeFetcherInterface
}

func (h *getAllAssets) Handle(params operations.GetAllAssetsParams) middleware.Responder {
	// Load the wallet based on the provided Nickname
	wlt, resp := loadWallet(params.Nickname)
	if resp != nil {
		return resp
	}

	// Create a slice to store the assets with their balances
	AssetsWithBalance := make([]*models.AssetInfoWithBalance, 0)

	// Fetch the account information for the wallet using the massaClient
	infos, err := h.massaClient.GetAccountsInfos([]wallet.Wallet{*wlt})
	if err != nil {
		// Handle the error and return an internal server error response
		errorMsg := fmt.Sprintf("Failed to fetch balance for asset %s: %s", "MASSA", err.Error())
		return operations.NewGetAllAssetsInternalServerError().WithPayload(&models.Error{
			Code:    errorFetchAssetBalance,
			Message: errorMsg,
		})
	}

	// Create the asset info for the Massa token and append it to the result slice
	MassaAsset := &models.AssetInfoWithBalance{
		AssetInfo: assets.XMAInfo(),
		Balance:   fmt.Sprint(infos[0].CandidateBalance),
	}
	AssetsWithBalance = append(AssetsWithBalance, MassaAsset)

	// Retrieve all assets from the selected WalletNickname
	for assetAddress, assetInfo := range h.AssetsStore.Assets[params.Nickname].ContractAssets {
		// Fetch the balance for the current asset
		balance, err := assets.Balance(assetAddress, wlt.Address)
		if err != nil {
			// Handle the error and return an internal server error response
			errorMsg := fmt.Sprintf("Failed to fetch balance for asset %s: %s", assetAddress, err.Error())
			return operations.NewGetAllAssetsInternalServerError().WithPayload(&models.Error{
				Code:    errorFetchAssetBalance,
				Message: errorMsg,
			})
		}

		// Create the asset info with balance and append it to the result slice
		assetWithBalance := &models.AssetInfoWithBalance{
			AssetInfo: assetInfo,
			Balance:   balance,
		}
		AssetsWithBalance = append(AssetsWithBalance, assetWithBalance)
	}

	// Return the list of assets with balance
	return operations.NewGetAllAssetsOK().WithPayload(AssetsWithBalance)
}
