package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

// NewGet instantiates a Get Handler
// The "classical" way is not possible because we need to pass to the handler a wallet.WalletPrompterInterface.
func NewGet(prompterApp wallet.WalletPrompterInterface) operations.RestWalletGetHandler {
	return &walletGet{prompterApp: prompterApp}
}

type walletGet struct {
	prompterApp wallet.WalletPrompterInterface
}

func (g *walletGet) Handle(params operations.RestWalletGetParams) middleware.Responder {
	wlt, err := wallet.Load(params.Nickname)
	if err == wallet.ErrorAccountNotFound {
		return operations.NewRestWalletGetNotFound().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: err.Error(),
			})
	} else if err != nil {
		return operations.NewRestWalletGetBadRequest().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: err.Error(),
			})
	}

	modelWallet := createModelWallet(*wlt)

	// if request not ciphered data, ask for password and unprotect the wallet
	if params.Ciphered != nil && !*params.Ciphered {
		promptData := &wallet.PromptRequestData{
			Msg:  fmt.Sprintf("Unprotect wallet %s", wlt.Nickname),
			Data: nil,
		}
		_, err := wlt.PromptPassword(g.prompterApp, walletapp.Export, promptData)
		if err != nil {
			return operations.NewRestWalletGetLocked().WithPayload(
				&models.Error{
					Code:    errorGetWallets,
					Message: "Unable to unprotect wallet",
				})
		}

		g.prompterApp.EmitEvent(walletapp.PasswordResultEvent,
			walletapp.EventData{Success: true, Data: "Unprotect Success"})

		modelWallet.KeyPair = models.WalletKeyPair{
			PrivateKey: wlt.GetPrivKey(),
			PublicKey:  wlt.GetPupKey(),
			Salt:       wlt.GetSalt(),
			Nonce:      wlt.GetNonce(),
		}
	}

	return operations.NewRestWalletGetOK().WithPayload(&modelWallet)
}

// HandleList handles a list request
func HandleList(params operations.RestWalletListParams) middleware.Responder {
	wallets, err := wallet.LoadAll()
	if err != nil {
		return operations.NewRestWalletListBadRequest().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: err.Error(),
			})
	}

	var wlts []*models.Wallet

	for i := 0; i < len(wallets); i++ {
		modelWallet := createModelWallet(wallets[i])
		wlts = append(wlts, &modelWallet)
	}

	return operations.NewRestWalletListOK().WithPayload(wlts)
}

func createModelWallet(wlt wallet.Wallet) models.Wallet {
	return models.Wallet{
		Nickname: models.Nickname(wlt.Nickname),
		Address:  wlt.Address,
		KeyPair:  models.WalletKeyPair{},
	}
}
