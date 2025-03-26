package wallet

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/openapi"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	walletpkg "github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station/pkg/logger"
)

// AppendEndpoints appends wallet endpoints to the API
// Note: the password prompter is mandatory for sign endpoint
func AppendEndpoints(api *operations.MassaWalletAPI, prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface) {
	api.CreateAccountHandler = NewCreateAccount(prompterApp, massaClient)
	api.DeleteAccountHandler = NewDelete(prompterApp, massaClient)
	api.ImportAccountHandler = NewImport(prompterApp, massaClient)
	api.AccountListHandler = NewGetAll(prompterApp.App().Wallet, massaClient)
	api.SignHandler = NewSign(prompterApp)
	api.SignMessageHandler = NewSignMessage(prompterApp)
	api.GetAccountHandler = NewGet(prompterApp, massaClient)
	api.ExportAccountFileHandler = NewWalletExportFile(prompterApp.App().Wallet)
	api.TransferCoinHandler = NewTransferCoin(prompterApp, massaClient)
	api.TradeRollsHandler = NewTradeRolls(prompterApp, massaClient)
	api.BackupAccountHandler = NewBackupAccount(prompterApp)
	api.UpdateAccountHandler = NewUpdateAccount(prompterApp, massaClient)
	api.AddAssetHandler = NewAddAsset(massaClient)
	api.GetAllAssetsHandler = NewGetAllAssets(prompterApp.App().Wallet, massaClient)
	api.DeleteAssetHandler = NewDeleteAsset()
	api.GetConfigHandler = NewGetConfig()
	api.AddSignRuleHandler = NewAddSignRuleHandler(prompterApp)
	api.DeleteSignRuleHandler = NewDeleteSignRuleHandler(prompterApp)
	api.UpdateSignRuleHandler = NewUpdateSignRuleHandler(prompterApp)
}

// loadAccount loads a wallet from the file system or returns an error.
// Here it is acceptable to return a middleware.Responder to simplify the code.
func loadAccount(wallet *walletpkg.Wallet, nickname string) (*account.Account, middleware.Responder) {
	acc, err := wallet.GetAccount(nickname)
	if err == nil {
		return acc, nil
	}

	if errors.Is(err, walletpkg.AccountNotFoundError) {
		return nil, newErrorResponse(fmt.Sprintf("%v: %s", err.Error(), nickname), errorGetAccount, http.StatusNotFound)
	} else {
		return nil, newErrorResponse(err.Error(), errorGetAccount, http.StatusBadRequest)
	}
}

func newErrorResponse(message, code string, statusCode int) middleware.Responder {
	logger.Error(message)

	payload := &models.Error{
		Code:    code,
		Message: message,
	}

	return openapi.NewPayloadResponder(statusCode, payload)
}
