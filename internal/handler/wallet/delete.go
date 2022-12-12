package wallet

import (
	"sync"

	"fyne.io/fyne/v2"

	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"

	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/gui"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/wallet"
)

//nolint:nolintlint,ireturn
func NewDelete(walletStorage *sync.Map, app *fyne.App) operations.RestWalletDeleteHandler {
	return &walletDelete{walletStorage: walletStorage, app: app}
}

type walletDelete struct {
	walletStorage *sync.Map
	app           *fyne.App
}

//nolint:nolintlint,ireturn
func (c *walletDelete) Handle(params operations.RestWalletDeleteParams) middleware.Responder {
	walletLoaded, err := wallet.Load(params.Nickname)
	if err != nil {
		return createInternalServerError(errorGetWallets, err.Error())
	}
	if len(params.Nickname) == 0 {
		return createInternalServerError(errorDeleteNoNickname, "Error: nickname field is mandatory.")
	}
	password := gui.AskPasswordDeleteWallet(params.Nickname, c.app)

	err = walletLoaded.Unprotect(password, 0)

	if err != nil {
		return createInternalServerError(errorWrongPassword, err.Error())
	}

	c.walletStorage.Delete(params.Nickname)

	err = wallet.Delete(params.Nickname)
	if err != nil {
		return createInternalServerError(errorDeleteFile, err.Error())
	}

	return operations.NewRestWalletDeleteNoContent()
}

//nolint:nolintlint,ireturn
func createInternalServerError(errorCode string, errorMessage string) middleware.Responder {
	return operations.NewRestWalletDeleteInternalServerError().
		WithPayload(
			&models.Error{
				Code:    errorCode,
				Message: errorMessage,
			})
}
