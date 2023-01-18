package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/password"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/privateKey"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/wallet"
)

func NewImport(pkPrompt privateKey.PrivateKeyAsker, pwdPrompt password.PasswordAsker) operations.RestWalletImportHandler {
	return &wImport{pwdPrompt: pwdPrompt, pkPrompt: pkPrompt}
}

type wImport struct {
	pwdPrompt password.PasswordAsker
	pkPrompt  privateKey.PrivateKeyAsker
}

func (c *wImport) Handle(params operations.RestWalletImportParams) middleware.Responder {

	walletName := params.Nickname

	password, err := c.pwdPrompt.Ask(walletName)
	if err != nil {
		return ImportWalletError(errorImportWalletCanceled, errorImportWalletCanceled)
	}

	privateKey, err := c.pkPrompt.Ask()
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
