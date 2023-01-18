package wallet

import (
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/password"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/privateKey"
)

// AppendEndpoints appends wallet endpoints to the API
// Note: the password prompter is mandatory for sign endpoint
func AppendEndpoints(api *operations.MassaWalletAPI, privateKeyPrompter privateKey.PrivateKeyAsker, passwordPrompter password.PasswordAsker) {
	api.RestWalletCreateHandler = operations.RestWalletCreateHandlerFunc(HandleCreate)
	api.RestWalletDeleteHandler = operations.RestWalletDeleteHandlerFunc(HandleDelete)
	api.RestWalletImportHandler = NewImport(privateKeyPrompter, passwordPrompter)
	api.RestWalletListHandler = operations.RestWalletListHandlerFunc(HandleList)
	api.RestWalletSignOperationHandler = NewSign(passwordPrompter)
}
