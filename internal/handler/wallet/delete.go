package wallet

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletApp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func NewDelete(walletApp *walletApp.WalletApp) operations.RestWalletDeleteHandler {
	return &walletDelete{walletApp: walletApp}
}

type walletDelete struct {
	walletApp *walletApp.WalletApp
}

// HandleDelete handles a delete request
func (w *walletDelete) Handle(params operations.RestWalletDeleteParams) middleware.Responder {
	wallet, err := wallet.Load(params.Nickname)
	if err != nil {
		return operations.NewRestWalletDeleteBadRequest().WithPayload(
			&models.Error{
				Code:    errorGetWallet,
				Message: "Error cannot load wallet: " + err.Error(),
			})
	}

	go wallet.Delete(w.walletApp)

	return operations.NewRestWalletDeleteNoContent()
}
