package wallet

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func NewCreateWallet(prompterApp wallet.WalletPrompterInterface) operations.CreateWalletHandler {
	return &walletCreate{prompterApp: prompterApp}
}

type walletCreate struct {
	prompterApp wallet.WalletPrompterInterface
}

func (w *walletCreate) Handle(params operations.CreateWalletParams) middleware.Responder {

	nickname := string(params.Body.Nickname)
	//nolint:gosimple
	password, err := wallet.PromptCreatePassword(w.prompterApp, nickname)
	if err != nil {
		return operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: "Unable to unprotect wallet",
			})
	}

	w.prompterApp.EmitEvent(walletapp.PasswordResultEvent,
		walletapp.EventData{Success: true, Data: "New password created"})

	newWallet, err := wallet.Generate(nickname, password)
	if err != nil {
		return operations.NewRestWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCreateNew,
				Message: err.Error(),
			})
	}

	return New(newWallet)
}

func New(newWallet *wallet.Wallet) middleware.Responder {
	return operations.NewCreateWalletOK().WithPayload(
		&models.Wallet{
			Nickname: models.Nickname(newWallet.Nickname),
			Address:  newWallet.Address,
			KeyPair: models.WalletKeyPair{
				PrivateKey: "",
				PublicKey:  newWallet.GetPupKey(),
				Salt:       "",
				Nonce:      "",
			},
		})
}
