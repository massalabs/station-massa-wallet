package wallet

import (
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-wallet/pkg/delete"
	"github.com/massalabs/thyra-plugin-wallet/pkg/password"
	"github.com/massalabs/thyra-plugin-wallet/pkg/privateKey"
)

// AppendEndpoints appends wallet endpoints to the API
// Note: the password prompter is mandatory for sign endpoint
func AppendEndpoints(api *operations.MassaWalletAPI, passwordPrompter password.Asker, privateKeyPrompter privateKey.Asker, deletePrompter delete.Confirmer) {
	api.RestWalletCreateHandler = operations.RestWalletCreateHandlerFunc(HandleCreate)
	api.RestWalletDeleteHandler = NewDelete(deletePrompter)
	api.RestWalletImportHandler = NewImport(privateKeyPrompter, passwordPrompter)
	api.RestWalletListHandler = operations.RestWalletListHandlerFunc(HandleList)
	api.RestWalletSignOperationHandler = NewSign(passwordPrompter)
}
