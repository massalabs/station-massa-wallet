package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
)

func NewImport() operations.RestWalletImportHandler {
	return &wImport{}
}

type wImport struct{}

func (c *wImport) Handle(params operations.RestWalletImportParams) middleware.Responder {
	return nil
}

func errorBadRequest(code string, message string) middleware.Responder {
	return operations.NewRestWalletImportBadRequest().WithPayload(
		&models.Error{
			Code:    code,
			Message: message,
		})
}
