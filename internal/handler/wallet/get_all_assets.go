package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station-massa-wallet/pkg/walletmanager"
)

func NewGetAllAssets(wallet *walletmanager.Wallet, AssetsStore *assets.AssetsStore, massaClient network.NodeFetcherInterface) operations.GetAllAssetsHandler {
	return &getAllAssets{
		wallet:      wallet,
		AssetsStore: AssetsStore,
		massaClient: massaClient,
	}
}

type getAllAssets struct {
	wallet      *walletmanager.Wallet
	AssetsStore *assets.AssetsStore
	massaClient network.NodeFetcherInterface
}

func (g *getAllAssets) Handle(params operations.GetAllAssetsParams) middleware.Responder {
	// Load the wallet based on the provided Nickname
	acc, resp := loadAccount(g.wallet, params.Nickname)
	if resp != nil {
		return resp
	}

	// Create a slice to store the assets with their balances
	AssetsWithBalance := make([]*models.AssetInfoWithBalance, 0)

	// Fetch the account information for the wallet using the massaClient
	infos, err := g.massaClient.GetAccountsInfos([]account.Account{*acc})
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
		AssetInfo: assets.MASInfo(),
		Balance:   fmt.Sprint(infos[0].CandidateBalance),
	}
	AssetsWithBalance = append(AssetsWithBalance, MassaAsset)

	// Retrieve all assets from the selected WalletNickname
	for assetAddress, assetInfo := range g.AssetsStore.Assets[params.Nickname].ContractAssets {
		// First, check if the asset exists in the network
		if !g.massaClient.AssetExistInNetwork(assetAddress) {
			// If the asset does not exist in the network, skip it and go to the next one
			continue
		}
		// Fetch the balance for the current asset
		address, err := acc.Address.MarshalText()
		if err != nil {
			return operations.NewGetAllAssetsInternalServerError().WithPayload(
				&models.Error{
					Code:    errorGetWallet,
					Message: ErrorAddressInvalid.Error(),
				})
		}

		balance, err := g.massaClient.DatastoreAssetBalance(assetAddress, string(address))
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
