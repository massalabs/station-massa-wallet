package wallet

import (
	"errors"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/gui"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/wallet"
)

const fileModeUserRW = 0o600

//nolint:nolintlint,ireturn
func NewImport(walletStorage *sync.Map, app *fyne.App) operations.RestWalletImportHandler {
	return &wImport{walletStorage: walletStorage, app: app}
}

type wImport struct {
	walletStorage *sync.Map
	app           *fyne.App
}

//nolint:nolintlint,ireturn,funlen
func (c *wImport) Handle(params operations.RestWalletImportParams) middleware.Responder {
	password, walletName, privateKey, err := gui.AskWalletInfo(c.app)
	if err != nil {
		return NewWalletError(errorCreateNew, err.Error())
	}

	if len(walletName) == 0 {
		return operations.NewRestWalletCreateBadRequest().WithPayload(
			&models.Error{
				Code:    errorCreateNoNickname,
				Message: "Error: nickname field is mandatory.",
			})
	}

	_, inStore := c.walletStorage.Load(walletName)
	if inStore {
		return NewWalletError(errorAlreadyExists, "Error: a wallet with the same nickname already exists.")
	}

	if len(password) == 0 {
		return NewWalletError(errorCreateNoPassword, "Error: password field is mandatory.")
	}

	newWallet, err := wallet.Imported(walletName, privateKey)
	if err != nil {
		if errors.Is(err, wallet.ErrWalletAlreadyImported) {
			return NewWalletError(errorAlreadyImported, err.Error())
		}

		return NewWalletError(errorCreateNew, err.Error())
	}

	return CreateNewWallet(&walletName, &password, c.walletStorage, newWallet)
}

//nolint:nolintlint,ireturn
func NewWalletError(code string, message string) middleware.Responder {
	return operations.NewRestWalletCreateInternalServerError().WithPayload(
		&models.Error{
			Code:    code,
			Message: message,
		})
}
