package wallet

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station/pkg/logger"
)

func NewGetAllAssets(wallet *wallet.Wallet, AssetsStore *assets.AssetsStore, massaClient network.NodeFetcherInterface) operations.GetAllAssetsHandler {
	return &getAllAssets{
		wallet:      wallet,
		AssetsStore: AssetsStore,
		massaClient: massaClient,
	}
}

type getAllAssets struct {
	wallet      *wallet.Wallet
	AssetsStore *assets.AssetsStore
	massaClient network.NodeFetcherInterface
}

func (g *getAllAssets) Handle(params operations.GetAllAssetsParams) middleware.Responder {
	// Load the wallet based on the provided Nickname
	acc, errResp := loadAccount(g.wallet, params.Nickname)
	if errResp != nil {
		return errResp
	}

	// Create a slice to store the assets with their balances
	assetsWithBalance := make([]*models.AssetInfoWithBalance, 0)

	massaAsset, resp := g.getMASAsset(acc)
	if resp != nil {
		return resp
	}

	assetsWithBalance = append(assetsWithBalance, massaAsset)

	userAssetData, resp := g.getAssetsData(acc)
	if resp != nil {
		return resp
	}

	assetsWithBalance = append(assetsWithBalance, userAssetData...)

	// sort AssetsWithBalance by name
	sort.Slice(assetsWithBalance, func(i, j int) bool {
		return assetsWithBalance[i].Name < assetsWithBalance[j].Name
	})

	// Return the list of assets with balance
	return operations.NewGetAllAssetsOK().WithPayload(assetsWithBalance)
}

func (g *getAllAssets) getMASAsset(acc *account.Account) (*models.AssetInfoWithBalance, middleware.Responder) {
	// Fetch the account information for the wallet using the massaClient
	infos, err := g.massaClient.GetAccountsInfos([]*account.Account{acc})
	if err != nil {
		// Handle the error and return an internal server error response
		errorMsg := fmt.Sprintf("Failed to fetch balance for asset %s: %s", "MASSA", err.Error())

		return nil, operations.NewGetAllAssetsInternalServerError().WithPayload(&models.Error{
			Code:    errorFetchAssetBalance,
			Message: errorMsg,
		})
	}

	// Create the asset info for the Massa token and append it to the result slice
	massaAsset := &models.AssetInfoWithBalance{
		AssetInfo: assets.MASInfo(),
		Balance:   fmt.Sprint(infos[0].CandidateBalance),
	}

	return massaAsset, nil
}

func (g *getAllAssets) getAssetsData(acc *account.Account) ([]*models.AssetInfoWithBalance, middleware.Responder) {
	assetsWithBalance := make([]*models.AssetInfoWithBalance, 0)

	// Retrieve all assets from the selected nickname
	for assetAddress, assetInfo := range g.AssetsStore.Assets[acc.Nickname].ContractAssets {
		// Fetch the balance for the current asset
		balance, resp := g.fetchAssetData(assetAddress, acc)
		if resp != nil {
			return nil, resp
		}

		// If the asset does not exist in the network, skip it and go to the next one
		if balance == nil {
			continue
		}

		// Create the asset info with balance and append it to the result slice
		assetWithBalance := &models.AssetInfoWithBalance{
			AssetInfo: assetInfo,
			Balance:   *balance,
		}

		assetsWithBalance = append(assetsWithBalance, assetWithBalance)
	}

	return assetsWithBalance, nil
}

func (g *getAllAssets) fetchAssetData(assetAddress string, acc *account.Account) (*string, middleware.Responder) {
	// First, check if the asset exists in the network
	if !g.massaClient.AssetExistInNetwork(assetAddress) {
		logger.Infof("Asset %s does not exist in the network", assetAddress)
		return nil, nil
	}

	address, err := acc.Address.MarshalText()
	if err != nil {
		return nil, newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
	}

	balance, err := g.massaClient.DatastoreAssetBalance(assetAddress, string(address))
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to fetch balance for asset %s: %s", assetAddress, err.Error())

		// Handle the error and return an internal server error response
		return nil, operations.NewGetAllAssetsInternalServerError().WithPayload(&models.Error{
			Code:    errorFetchAssetBalance,
			Message: errorMsg,
		})
	}

	return &balance, nil
}
