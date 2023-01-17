package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/guiModal"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/wallet"
)

func NewImport(walletInfoModal guiModal.WalletInfoAsker) operations.RestWalletImportHandler {
	return &wImport{walletInfoModal: walletInfoModal}
}

type wImport struct {
	walletInfoModal guiModal.WalletInfoAsker
}

func (c *wImport) Handle(operations.RestWalletImportParams) middleware.Responder {

	password, walletName, privateKey, err := c.walletInfoModal.WalletInfo()
	if err != nil {
		return ImportWalletError(errorImportWalletCanceled, errorImportWalletCanceled)
	}

	_, err = wallet.Load(walletName)
	if err == nil {
		return ImportWalletError(err.Error(), "Error: a wallet with the same nickname already exists.")
	}

	newWallet, err := wallet.Imported(walletName, privateKey, password)
	if err != nil {
		return ImportWalletError(err.Error(), err.Error())
	}

	return New(newWallet)
}

func ImportWalletError(code string, message string) middleware.Responder {
	return operations.NewRestWalletCreateInternalServerError().WithPayload(
		&models.Error{
			Code:    code,
			Message: message,
		})
}
