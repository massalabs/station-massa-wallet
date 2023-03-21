package wallet

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"

	"github.com/massalabs/thyra-plugin-wallet/pkg/delete"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func NewDelete(deleteConfirm delete.Confirmer) operations.RestWalletDeleteHandler {
	return &walletDelete{deleteConfirmPrompt: deleteConfirm}
}

type walletDelete struct {
	deleteConfirmPrompt delete.Confirmer
}

// HandleDelete handles a delete request
func (c *walletDelete) Handle(params operations.RestWalletDeleteParams) middleware.Responder {
	if len(params.Nickname) == 0 {
		return operations.NewRestWalletDeleteBadRequest().WithPayload(
			&models.Error{
				Code:    errorDeleteNoNickname,
				Message: "Error: nickname field is mandatory.",
			})
	}

	walletLoaded, err := wallet.Load(params.Nickname)
	if err != nil {
		return operations.NewRestWalletDeleteInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallet,
				Message: "Error cannot load wallet: " + err.Error(),
			})
	}

	password, err := c.deleteConfirmPrompt.Confirm(params.Nickname)
	if err != nil {
		return operations.NewRestWalletDeleteInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: errorCanceledAction,
			})
	}

	err = walletLoaded.Unprotect(password)
	if err != nil {
		return operations.NewRestWalletDeleteInternalServerError().WithPayload(
			&models.Error{
				Code:    errorDeleteFile,
				Message: err.Error(),
			})
	}

	err = wallet.Delete(params.Nickname)
	if err != nil {
		return operations.NewRestWalletDeleteInternalServerError().WithPayload(
			&models.Error{
				Code:    errorDeleteFile,
				Message: err.Error(),
			})
	}

	return operations.NewRestWalletDeleteNoContent()
}
