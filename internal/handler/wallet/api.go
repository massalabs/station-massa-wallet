package wallet

import (
	"github.com/bluele/gcache"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

// AppendEndpoints appends wallet endpoints to the API
// Note: the password prompter is mandatory for sign endpoint
func AppendEndpoints(api *operations.MassaWalletAPI, prompterApp wallet.WalletPrompterInterface, gc gcache.Cache) {
	api.RestWalletCreateHandler = operations.RestWalletCreateHandlerFunc(HandleCreate)
	api.RestWalletDeleteHandler = NewDelete(prompterApp)
	api.RestWalletImportHandler = NewImport()
	api.RestWalletListHandler = operations.RestWalletListHandlerFunc(HandleList)
	api.RestWalletSignOperationHandler = NewSign(prompterApp, gc)
	api.RestWalletGetHandler = NewGet(prompterApp)
	api.RestWalletExportFileHandler = operations.RestWalletExportFileHandlerFunc(HandleExportFile)
}
