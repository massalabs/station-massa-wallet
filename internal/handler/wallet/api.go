package wallet

import (
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/guiModal"
)

// AppendEndpoints appends wallet endpoints to the API
// Note: the password prompter is mandatory for sign endpoint
func AppendEndpoints(api *operations.MassaWalletAPI, walletInfoModal guiModal.WalletInfoAsker, passwordPrompter guiModal.PasswordAsker) {
	api.RestWalletCreateHandler = operations.RestWalletCreateHandlerFunc(HandleCreate)
	api.RestWalletDeleteHandler = operations.RestWalletDeleteHandlerFunc(HandleDelete)
	api.RestWalletImportHandler = NewImport(walletInfoModal)
	api.RestWalletListHandler = operations.RestWalletListHandlerFunc(HandleList)
	api.RestWalletSignOperationHandler = NewSign(passwordPrompter)
}
