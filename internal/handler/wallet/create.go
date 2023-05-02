package wallet

import (
	"strings"

	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func NewCreateAccount(prompterApp prompt.WalletPrompterInterface) operations.RestCreateAccountHandler {
	return &walletCreate{prompterApp: prompterApp}
}

type walletCreate struct {
	prompterApp prompt.WalletPrompterInterface
}

func (w *walletCreate) Handle(params operations.RestCreateAccountParams) middleware.Responder {
	nickname := strings.TrimSpace(string(params.Nickname))

	if len(nickname) == 0 {
		return operations.NewRestCreateAccountBadRequest().WithPayload(
			&models.Error{
				Code:    errorCreateNoNickname,
				Message: "Error: nickname field is mandatory.",
			})
	}

	//nolint:gosimple
	password, err := prompt.PromptCreatePassword(w.prompterApp, nickname)
	if err != nil {
		return operations.NewRestCreateAccountUnauthorized().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: "Unable to create wallet",
			})
	}

	w.prompterApp.EmitEvent(walletapp.PasswordResultEvent,
		walletapp.EventData{Success: true, Data: "New password created"})

	newWallet, err := wallet.Generate(nickname, password)
	if err != nil {
		return operations.NewRestCreateAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCreateNew,
				Message: err.Error(),
			})
	}

	return New(newWallet)
}

func New(newWallet *wallet.Wallet) middleware.Responder {
	return operations.NewRestCreateAccountOK().WithPayload(
		&models.Account{
			Nickname: models.Nickname(newWallet.Nickname),
			Address:  newWallet.Address,
			KeyPair: models.AccountKeyPair{
				PrivateKey: "",
				PublicKey:  newWallet.GetPupKey(),
				Salt:       "",
				Nonce:      "",
			},
		})
}
