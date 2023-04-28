package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
)

func NewImport() operations.RestAccountImportHandler {
	return &wImport{}
}

type wImport struct{}

func (c *wImport) Handle(params operations.RestAccountImportParams) middleware.Responder {
	return nil
}
