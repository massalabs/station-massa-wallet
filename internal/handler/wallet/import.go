package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/network"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
	"github.com/massalabs/thyra-plugin-wallet/pkg/utils"
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
	promptRequest := prompt.PromptRequest{
		Action: walletapp.Import,
		Msg:    "Import",
		Data:   nil,
	}

	promptOutput, err := prompt.WakeUpPrompt(h.prompterApp, promptRequest, nil)
	if err != nil {
		// an event has been emitted during WakeUpPrompt
		errStr := fmt.Sprintf("Unable to import account: %v", err)
		return operations.NewImportAccountUnauthorized().WithPayload(
			&models.Error{
				Code:    errorImportWallet,
				Message: errStr,
			})
	}

	wlt, _ := promptOutput.(*wallet.Wallet)

	err = wlt.Persist()
	if err != nil {
		errStr := fmt.Sprintf("Unable to persist imported account: %v", err)
		h.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.ErrAccountFile})
		return operations.NewImportAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorImportWallet,
				Message: errStr,
			})
	}

	h.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, CodeMessage: utils.MsgAccountImported})

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
