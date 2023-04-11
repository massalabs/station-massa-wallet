package wallet

import (
	"github.com/bluele/gcache"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
)

// AppendEndpoints appends wallet endpoints to the API
// Note: the password prompter is mandatory for sign endpoint
func AppendEndpoints(api *operations.MassaWalletAPI, gc gcache.Cache) {
	api.RestWalletCreateHandler = operations.RestWalletCreateHandlerFunc(HandleCreate)
	api.RestWalletDeleteHandler = NewDelete()
	api.RestWalletImportHandler = NewImport()
	api.RestWalletListHandler = operations.RestWalletListHandlerFunc(HandleList)
	api.RestWalletSignOperationHandler = NewSign(gc)
	api.RestWalletGetHandler = operations.RestWalletGetHandlerFunc(HandleGet)
}
