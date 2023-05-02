package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

// NewGet instantiates a Get Handler
// The "classical" way is not possible because we need to pass to the handler a wallet.WalletPrompterInterface.
func NewGet(prompterApp prompt.WalletPrompterInterface) operations.RestAccountGetHandler {
	return &walletGet{prompterApp: prompterApp}
}

type walletGet struct {
	prompterApp prompt.WalletPrompterInterface
}

func (g *walletGet) Handle(params operations.RestAccountGetParams) middleware.Responder {
	// params.Nickname length is already checked by go swagger
	wlt, err := wallet.Load(params.Nickname)
	if err != nil {
		if err.Error() == wallet.ErrorAccountNotFound(params.Nickname).Error() {
			return operations.NewRestAccountGetNotFound().WithPayload(
				&models.Error{
					Code:    errorGetWallets,
					Message: err.Error(),
				})
		} else {
			return operations.NewRestAccountGetBadRequest().WithPayload(
				&models.Error{
					Code:    errorGetWallets,
					Message: err.Error(),
				})
		}
	}

	modelWallet := createModelWallet(*wlt)

	// if request not ciphered data, ask for password and unprotect the wallet
	if params.Ciphered != nil && !*params.Ciphered {
		promptData := &prompt.PromptRequestData{
			Msg:  fmt.Sprintf("Unprotect wallet %s", wlt.Nickname),
			Data: nil,
		}
		_, err := prompt.PromptPassword(g.prompterApp, wlt, walletapp.Export, promptData)
		if err != nil {
			return operations.NewRestAccountGetUnauthorized().WithPayload(
				&models.Error{
					Code:    errorGetWallets,
					Message: "Unable to unprotect wallet",
				})
		}

		g.prompterApp.EmitEvent(walletapp.PasswordResultEvent,
			walletapp.EventData{Success: true, Data: "Unprotect Success"})

		modelWallet.KeyPair = models.AccountKeyPair{
			PrivateKey: wlt.GetPrivKey(),
			PublicKey:  wlt.GetPupKey(),
			Salt:       wlt.GetSalt(),
			Nonce:      wlt.GetNonce(),
		}
	}

	return operations.NewRestAccountGetOK().WithPayload(&modelWallet)
}

// HandleList handles a list request
func HandleList(params operations.RestAccountListParams) middleware.Responder {
	wallets, err := wallet.LoadAll()
	if err != nil {
		return operations.NewRestAccountListInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: err.Error(),
			})
	}

	var wlts []*models.Account

	for i := 0; i < len(wallets); i++ {
		modelWallet := createModelWallet(wallets[i])
		wlts = append(wlts, &modelWallet)
	}

	return operations.NewRestAccountListOK().WithPayload(wlts)
}

func createModelWallet(wlt wallet.Wallet) models.Account {
	return models.Account{
		Nickname: models.Nickname(wlt.Nickname),
		Address:  wlt.Address,
		KeyPair:  models.AccountKeyPair{},
	}
}
