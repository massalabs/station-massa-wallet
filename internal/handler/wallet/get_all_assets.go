package wallet

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

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
		dollarValueFloatI, err := strconv.ParseFloat(assetsWithBalance[i].DollarValue, 64)
		if err != nil {
			logger.Errorf("Failed to parse dollar value for asset %s: %s", assetsWithBalance[i].Name, err.Error())
			return false
		}
		dollarValueFloatJ, err := strconv.ParseFloat(assetsWithBalance[j].DollarValue, 64)
		if err != nil {
			logger.Errorf("Failed to parse dollar value for asset %s: %s", assetsWithBalance[j].Name, err.Error())
			return false
		}
		return dollarValueFloatI > dollarValueFloatJ
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
	asset := assets.MASInfo()
	massaAsset := &models.AssetInfoWithBalance{
		AssetInfo: asset,
		Balance:   fmt.Sprint(infos[0].CandidateBalance),
		IsDefault: true,
	}

	dollarValue, err := assets.DollarValue(massaAsset.Balance, asset.Symbol, *asset.Decimals)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to fetch dollar value for asset %s: %s", asset.Address, err.Error())
		logger.Errorf(errorMsg)
	}

	massaAsset.DollarValue = dollarValue

	return massaAsset, nil
}

func (g *getAllAssets) getAssetsData(acc *account.Account) ([]*models.AssetInfoWithBalance, middleware.Responder) {
	defaultAssets, err := assets.GetDefaultAssets()
	if err != nil {
		logger.Errorf("Failed to get default assets: %s", err.Error())
	}

	assetsInfo := make([]*models.AssetInfoWithBalance, 0)

	// Initialize map to track addressed already added
	includedAddresses := map[string]bool{}

	for _, asset := range defaultAssets {
		completeAsset := &models.AssetInfoWithBalance{
			AssetInfo: asset,
			Balance:   "",
			IsDefault: true,
		}
		assetsInfo = append(assetsInfo, completeAsset)
		includedAddresses[asset.Address] = true
	}

	// Append default assets ensuring no duplication
	for _, asset := range g.AssetsStore.Assets[acc.Nickname].ContractAssets {
		// Append the asset info to the result slice if it is not already in the list
		if _, exists := includedAddresses[asset.Address]; !exists {
			completeAsset := &models.AssetInfoWithBalance{
				AssetInfo: asset,
				Balance:   "",
				IsDefault: false,
			}
			assetsInfo = append(assetsInfo, completeAsset)
			includedAddresses[asset.Address] = true
		}
	}

	assetsWithBalance := make([]*models.AssetInfoWithBalance, 0)

	// Retrieve all assets from the selected nickname
	for _, asset := range assetsInfo {
		// First, check if the asset exists in the network
		if !g.massaClient.AssetExistInNetwork(asset.Address) {
			logger.Infof("Asset %s does not exist in the network", asset.Address)
			continue
		}

		// Fetch the balance for the current asset
		balance, dollarValue, resp := g.fetchAssetData(asset, acc)
		if resp != nil {
			return nil, resp
		}

		asset.Balance = balance
		asset.DollarValue = dollarValue
		assetsWithBalance = append(assetsWithBalance, asset)
	}

	return assetsWithBalance, nil
}

func (g *getAllAssets) fetchAssetData(asset *models.AssetInfoWithBalance, acc *account.Account) (string, string, middleware.Responder) {
	assetAddress := asset.Address

	// Balance
	address, err := acc.Address.MarshalText()
	if err != nil {
		return "", "", newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
	}

	balance, err := g.massaClient.DatastoreAssetBalance(assetAddress, string(address))
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to fetch balance for asset %s: %s", assetAddress, err.Error())
		return "", "", newErrorResponse(errorMsg, errorFetchAssetBalance, http.StatusInternalServerError)
	}

	// Dollar value
	dollarValue, err := assets.DollarValue(balance, asset.Symbol, *asset.Decimals)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to fetch dollar value for asset %s: %s", assetAddress, err.Error())
		logger.Errorf(errorMsg)

		return balance, "", nil
	}

	return balance, dollarValue, nil
}
