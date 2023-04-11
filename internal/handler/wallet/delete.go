package wallet

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
)

func NewDelete() operations.RestWalletDeleteHandler {
	return &walletDelete{}
}

type walletDelete struct {
}

// HandleDelete handles a delete request
func (c *walletDelete) Handle(params operations.RestWalletDeleteParams) middleware.Responder {

	return operations.NewRestWalletDeleteNoContent()
}
