package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
)

func NewImport(prompterApp prompt.WalletPrompterInterface) operations.ImportAccountHandler {
	return &wImport{prompterApp: prompterApp}
}

type wImport struct {
	prompterApp prompt.WalletPrompterInterface
}

func (h *wImport) Handle(_ operations.ImportAccountParams) middleware.Responder {
	wallet, err := prompt.PromptImport(h.prompterApp)
	if err != nil {
		return operations.NewImportAccountUnauthorized().WithPayload(
			&models.Error{
				Code:    errorImportWallet,
				Message: "Unable to import account",
			})
	}

	err = wallet.Persist()
	if err != nil {
		return operations.NewImportAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorImportWallet,
				Message: "Unable to persist imported account",
			})
	}

	h.prompterApp.EmitEvent(walletapp.PasswordResultEvent,
		walletapp.EventData{Success: true, Data: "Import Success"})
	return operations.NewImportAccountOK().WithPayload(
		&models.Account{
			Nickname: models.Nickname(wallet.Nickname),
			Address:  wallet.Address,
			KeyPair: models.AccountKeyPair{
				PrivateKey: "",
				PublicKey:  wallet.GetPupKey(),
				Salt:       "",
				Nonce:      "",
			},
		})
}
