package wallet

import (
	"fmt"
	"sort"
	"sync"

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
	assetsWithBalance := make([]*assets.AssetInfoWithBalances, 0)

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
		if assetsWithBalance[i].AssetInfo.Symbol == "MAS" {
			return true
		}
		if assetsWithBalance[j].AssetInfo.Symbol == "MAS" {
			return false
		}
		if assetsWithBalance[i].DollarValue == nil {
			return false
		}
		if assetsWithBalance[j].DollarValue == nil {
			return false
		}
		return *assetsWithBalance[i].DollarValue > *assetsWithBalance[j].DollarValue
	})

	// Return the list of assets with balance
	return operations.NewGetAllAssetsOK().WithPayload(convertToModel(assetsWithBalance))
}

func (g *getAllAssets) getMASAsset(acc *account.Account) (*assets.AssetInfoWithBalances, middleware.Responder) {
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
	massaAsset := &assets.AssetInfoWithBalances{
		AssetInfo:  &asset,
		Balance:    fmt.Sprint(infos[0].CandidateBalance),
		IsDefault:  true,
		MEXCSymbol: "MASUSDT",
	}

	dollarValue, err := assets.DollarValue(massaAsset.Balance, massaAsset.MEXCSymbol, asset.Symbol, *asset.Decimals)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to fetch dollar value for asset %s: %s", asset.Address, err.Error())
		logger.Errorf(errorMsg)
	}

	massaAsset.DollarValue = dollarValue

	return massaAsset, nil
}

// getAssetsData fetches the balance and dollar value for each asset in the account.
// If user has asset that are deployed on another network, it will not be included.
func (g *getAllAssets) getAssetsData(acc *account.Account) ([]*assets.AssetInfoWithBalances, middleware.Responder) {
	assetsInfo := g.AssetsStore.All(acc.Nickname)

	assetsWithBalance := make([]*assets.AssetInfoWithBalances, 0)
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	resultsCh := make(chan *assets.AssetInfoWithBalances)

	// Retrieve all assets from the selected nickname
	for _, asset := range assetsInfo {
		wg.Add(1)

		go func(asset *assets.AssetInfoWithBalances) {
			defer wg.Done()

			// First, check if the asset exists in the network
			if !g.massaClient.AssetExistInNetwork(asset.AssetInfo.Address) {
				logger.Infof("Asset %s does not exist in the network", asset.AssetInfo.Address)
				return
			}

			// Fetch the balance for the current asset
			balance, dollarValue := g.fetchAssetData(asset, acc)

			asset.Balance = balance
			asset.DollarValue = dollarValue

			resultsCh <- asset
		}(asset)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	for result := range resultsCh {
		mu.Lock()
		assetsWithBalance = append(assetsWithBalance, result)
		mu.Unlock()
	}

	return assetsWithBalance, nil
}

func (g *getAllAssets) fetchAssetData(asset *assets.AssetInfoWithBalances, acc *account.Account) (string, *float64) {
	assetAddress := asset.AssetInfo.Address

	// Balance
	address, err := acc.Address.MarshalText()
	if err != nil {
		logger.Errorf("Failed to marshal address: %s", err.Error())

		return "", nil
	}

	balance, err := g.massaClient.DatastoreAssetBalance(assetAddress, string(address))
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to fetch balance for asset %s: %s", assetAddress, err.Error())
		logger.Errorf(errorMsg)

		return "", nil
	}

	// Dollar value
	dollarValue, err := assets.DollarValue(balance, asset.MEXCSymbol, asset.AssetInfo.Symbol, *asset.AssetInfo.Decimals)
	if err != nil {
		logger.Warnf(fmt.Sprintf("Failed to fetch dollar value for asset %s: %s", assetAddress, err.Error()))

		return balance, nil
	}

	return balance, dollarValue
}

func convertToModel(assetsWithBalance []*assets.AssetInfoWithBalances) []*models.AssetInfoWithBalance {
	result := make([]*models.AssetInfoWithBalance, 0)

	for _, asset := range assetsWithBalance {
		assetInfo := models.AssetInfo{
			Address:  asset.AssetInfo.Address,
			Decimals: asset.AssetInfo.Decimals,
			Name:     asset.AssetInfo.Name,
			Symbol:   asset.AssetInfo.Symbol,
		}

		if asset.DollarValue == nil {
			result = append(result, &models.AssetInfoWithBalance{
				AssetInfo: assetInfo,
				Balance:   asset.Balance,
				IsDefault: asset.IsDefault,
			})
		} else {
			result = append(result, &models.AssetInfoWithBalance{
				AssetInfo:   assetInfo,
				Balance:     asset.Balance,
				DollarValue: fmt.Sprintf("%.2f", *asset.DollarValue),
				IsDefault:   asset.IsDefault,
			})
		}
	}

	return result
}
