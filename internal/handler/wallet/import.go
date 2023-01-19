package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/password"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/privateKey"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/wallet"
)

func NewImport(pkPrompt privateKey.Asker, pwdPrompt password.Asker) operations.RestWalletImportHandler {
	return &wImport{pwdPrompt: pwdPrompt, pkPrompt: pkPrompt}
}

type wImport struct {
	pwdPrompt password.Asker
	pkPrompt  privateKey.Asker
}

func (c *wImport) Handle(params operations.RestWalletImportParams) middleware.Responder {
	walletName := params.Nickname

	password, err := c.pwdPrompt.Ask(walletName)
	if err != nil {
		return ImportWalletErrorBadRequest(errorImportWalletCanceled, errorImportWalletCanceled)

	}

	privateKey, err := c.pkPrompt.Ask()
	if err != nil {
		return ImportWalletErrorBadRequest(errorImportWalletCanceled, errorImportWalletCanceled)
	}

	_, err = wallet.Load(walletName)
	if err == nil {
		return operations.NewRestWalletImportInternalServerError().WithPayload(
			&models.Error{
				Code:    errorImportNickNameAlreadyTaken,
				Message: errorImportNickNameAlreadyTaken,
			})
	}

	newWallet, err := wallet.Import(walletName, privateKey, password)
	if err != nil {
		return ImportWalletErrorBadRequest(err.Error(), err.Error())
	}

	return New(newWallet)
}

func ImportWalletErrorBadRequest(code string, message string) middleware.Responder {
	return operations.NewRestWalletImportBadRequest().WithPayload(
		&models.Error{
			Code:    code,
			Message: message,
		})
}
