package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/network"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

// NewGet instantiates a Get Handler
// The "classical" way is not possible because we need to pass to the handler a wallet.WalletPrompterInterface.
func NewGet(prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface) operations.GetAccountHandler {
	return &walletGet{prompterApp: prompterApp, massaClient: massaClient}
}

type walletGet struct {
	prompterApp prompt.WalletPrompterInterface
	massaClient network.NodeFetcherInterface
}

func (g *walletGet) Handle(params operations.GetAccountParams) middleware.Responder {
	wlt, resp := loadWallet(params.Nickname)
	if resp != nil {
		return resp
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
			return operations.NewGetAccountUnauthorized().WithPayload(
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

	infos, err := g.massaClient.GetAccountsInfos([]wallet.Wallet{*wlt})
	if err != nil {
		return operations.NewGetAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: "Unable to retrieve account infos",
			})
	}

	modelWallet.CandidateBalance = models.Amount(fmt.Sprint(infos[0].CandidateBalance))
	modelWallet.Balance = models.Amount(fmt.Sprint(infos[0].Balance))

	return operations.NewGetAccountOK().WithPayload(&modelWallet)
}

func createModelWallet(wlt wallet.Wallet) models.Account {
	return models.Account{
		Nickname: models.Nickname(wlt.Nickname),
		Address:  wlt.Address,
		KeyPair:  models.AccountKeyPair{},
	}
}
