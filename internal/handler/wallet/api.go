package wallet

import (
	"github.com/bluele/gcache"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/assets"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
)

// AppendEndpoints appends wallet endpoints to the API
// Note: the password prompter is mandatory for sign endpoint
func AppendEndpoints(api *operations.MassaWalletAPI, prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface, assetsStore *assets.AssetsStore, gc gcache.Cache) {
	api.CreateAccountHandler = NewCreateAccount(prompterApp, massaClient)
	api.DeleteAccountHandler = NewDelete(prompterApp, massaClient)
	api.ImportAccountHandler = NewImport(prompterApp, massaClient)
	api.AccountListHandler = NewGetAll(massaClient)
	api.SignHandler = NewSign(prompterApp, gc)
	api.GetAccountHandler = NewGet(prompterApp, massaClient)
	api.ExportAccountFileHandler = operations.ExportAccountFileHandlerFunc(HandleExportFile)
	api.TransferCoinHandler = NewTransferCoin(prompterApp, massaClient)
	api.TradeRollsHandler = NewTradeRolls(prompterApp, massaClient)
	api.BackupAccountHandler = NewBackupAccount(prompterApp)
	api.UpdateAccountHandler = NewUpdateAccount(prompterApp, massaClient)
	api.AddAssetHandler = NewAddAsset(assetsStore)
	api.GetAllAssetsHandler = NewGetAllAssets(assetsStore)
}

// loadWallet loads a wallet from the file system or returns an error.
func loadWallet(nickname string) (*wallet.Wallet, middleware.Responder) {
	w, err := wallet.Load(nickname)
	if err == nil {
		return w, nil
	}

	errorObj := models.Error{
		Code:    errorGetWallets,
		Message: err.Error(),
	}

	if err.Error() == wallet.ErrorAccountNotFound(nickname).Error() {
		return nil, operations.NewGetAccountNotFound().WithPayload(&errorObj)
	} else {
		return nil, operations.NewGetAccountBadRequest().WithPayload(&errorObj)
	}
}
