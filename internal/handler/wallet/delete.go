package wallet

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"

	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/wallet"
)

// HandleDelete handles a delete request
func HandleDelete(params operations.RestWalletDeleteParams) middleware.Responder {
	if len(params.Nickname) == 0 {
		return operations.NewRestWalletDeleteBadRequest().WithPayload(
			&models.Error{
				Code:    errorDeleteNoNickname,
				Message: "Error: nickname field is mandatory.",
			})
	}

	err := wallet.Delete(params.Nickname)
	if err != nil {
		return operations.NewRestWalletDeleteInternalServerError().WithPayload(
			&models.Error{
				Code:    errorDeleteFile,
				Message: err.Error(),
			})
	}

	return operations.NewRestWalletDeleteNoContent()
}
