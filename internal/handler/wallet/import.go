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

func NewImport(prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface) operations.ImportAccountHandler {
	return &wImport{prompterApp: prompterApp, massaClient: massaClient}
}

type wImport struct {
	prompterApp prompt.WalletPrompterInterface
	massaClient network.NodeFetcherInterface
}

func (h *wImport) Handle(_ operations.ImportAccountParams) middleware.Responder {
	wlt, err := prompt.PromptImport(h.prompterApp)
	if err != nil {
		return operations.NewImportAccountUnauthorized().WithPayload(
			&models.Error{
				Code:    errorImportWallet,
				Message: "Unable to import account",
			})
	}

	err = wlt.Persist()
	if err != nil {
		errStr := "Unable to persist imported account: " + err.Error()
		h.prompterApp.EmitEvent(walletapp.PasswordResultEvent,
			walletapp.EventData{Success: false, Data: errStr})
		return operations.NewImportAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorImportWallet,
				Message: errStr,
			})
	}

	h.prompterApp.EmitEvent(walletapp.PasswordResultEvent,
		walletapp.EventData{Success: true, Data: "Import Success"})

	infos, err := h.massaClient.GetAccountsInfos([]wallet.Wallet{*wlt})
	if err != nil {
		return operations.NewImportAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: "Unable to retrieve accounts infos",
			})
	}

	return operations.NewImportAccountOK().WithPayload(
		&models.Account{
			Nickname:         models.Nickname(wlt.Nickname),
			CandidateBalance: models.Amount(fmt.Sprint(infos[0].CandidateBalance)),
			Balance:          models.Amount(fmt.Sprint(infos[0].Balance)),
			Address:          wlt.Address,
			KeyPair: models.AccountKeyPair{
				PrivateKey: "",
				PublicKey:  wlt.GetPupKey(),
				Salt:       "",
				Nonce:      "",
			},
		})
}
