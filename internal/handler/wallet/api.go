package wallet

import (
	"github.com/bluele/gcache"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
)

// AppendEndpoints appends wallet endpoints to the API
// Note: the password prompter is mandatory for sign endpoint
func AppendEndpoints(api *operations.MassaWalletAPI, prompterApp prompt.WalletPrompterInterface, gc gcache.Cache) {
	api.CreateAccountHandler = NewCreateAccount(prompterApp)
	api.DeleteAccountHandler = NewDelete(prompterApp)
	api.ImportAccountHandler = NewImport(prompterApp)
	api.AccountListHandler = operations.AccountListHandlerFunc(HandleList)
	api.SignOperationHandler = NewSign(prompterApp, gc)
	api.GetAccountHandler = NewGet(prompterApp)
	api.ExportAccountFileHandler = operations.ExportAccountFileHandlerFunc(HandleExportFile)
}
