package wallet

import (
	"github.com/bluele/gcache"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletApp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
)

// AppendEndpoints appends wallet endpoints to the API
// Note: the password prompter is mandatory for sign endpoint
func AppendEndpoints(api *operations.MassaWalletAPI, walletApp *walletApp.WalletApp, gc gcache.Cache) {
	api.RestWalletCreateHandler = operations.RestWalletCreateHandlerFunc(HandleCreate)
	api.RestWalletDeleteHandler = NewDelete(walletApp)
	api.RestWalletImportHandler = NewImport()
	api.RestWalletListHandler = operations.RestWalletListHandlerFunc(HandleList)
	api.RestWalletSignOperationHandler = NewSign(gc)
	api.RestWalletGetHandler = operations.RestWalletGetHandlerFunc(HandleGet)
}
