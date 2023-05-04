package wallet

import (
	"github.com/bluele/gcache"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-wallet/pkg/network"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
)

// AppendEndpoints appends wallet endpoints to the API
// Note: the password prompter is mandatory for sign endpoint
func AppendEndpoints(api *operations.MassaWalletAPI, prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface, gc gcache.Cache) {
	api.CreateAccountHandler = NewCreateAccount(prompterApp, massaClient)
	api.DeleteAccountHandler = NewDelete(prompterApp)
	api.ImportAccountHandler = NewImport(prompterApp, massaClient)
	api.AccountListHandler = NewGetAll(massaClient)
	api.SignHandler = NewSign(prompterApp, gc)
	api.GetAccountHandler = NewGet(prompterApp, massaClient)
	api.SignHandler = NewSign(prompterApp, gc)
	api.ExportAccountFileHandler = operations.ExportAccountFileHandlerFunc(HandleExportFile)
}
