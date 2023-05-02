package wallet

import (
	"github.com/bluele/gcache"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
)

// AppendEndpoints appends wallet endpoints to the API
// Note: the password prompter is mandatory for sign endpoint
func AppendEndpoints(api *operations.MassaWalletAPI, prompterApp prompt.WalletPrompterInterface, gc gcache.Cache) {
	api.RestCreateAccountHandler = NewCreateAccount(prompterApp)
	api.RestAccountDeleteHandler = NewDelete(prompterApp)
	api.RestAccountImportHandler = NewImport()
	api.RestAccountListHandler = operations.RestAccountListHandlerFunc(HandleList)
	api.RestAccountSignOperationHandler = NewSign(prompterApp, gc)
	api.RestAccountGetHandler = NewGet(prompterApp)
	api.RestAccountExportFileHandler = operations.RestAccountExportFileHandlerFunc(HandleExportFile)
}
