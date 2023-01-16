package wallet

import (
	"sync"

	"fyne.io/fyne/v2"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/password"
)

// AppendEndpoints appends wallet endpoints to the API
// Note: the password prompter is mandatory for sign endpoint
func AppendEndpoints(api *operations.MassaWalletAPI, app *fyne.App, passwordPrompter password.Asker) {
	var walletStorage sync.Map
	api.RestWalletCreateHandler = operations.RestWalletCreateHandlerFunc(HandleCreate)
	api.RestWalletDeleteHandler = operations.RestWalletDeleteHandlerFunc(HandleDelete)
	api.RestWalletImportHandler = NewImport(&walletStorage, app)
	api.RestWalletListHandler = operations.RestWalletListHandlerFunc(HandleList)
	api.RestWalletSignOperationHandler = NewSign(passwordPrompter)
}
