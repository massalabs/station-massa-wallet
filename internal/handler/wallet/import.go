package wallet

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

func NewImport(prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface) operations.ImportAccountHandler {
	return &walletImport{prompterApp: prompterApp, massaClient: massaClient}
}

type walletImport struct {
	prompterApp prompt.WalletPrompterInterface
	massaClient network.NodeFetcherInterface
}

func (w *walletImport) Handle(_ operations.ImportAccountParams) middleware.Responder {
	promptRequest := prompt.PromptRequest{
		Action: walletapp.Import,
		Msg:    "Import",
	}

	promptOutput, err := prompt.WakeUpPrompt(w.prompterApp, promptRequest, nil)
	if err != nil {
		// an event has been emitted during WakeUpPrompt
		errStr := fmt.Sprintf("Unable to import account: %v", err)

		return operations.NewImportAccountUnauthorized().WithPayload(
			&models.Error{
				Code:    errorImportWallet,
				Message: errStr,
			})
	}

	acc, _ := promptOutput.(*account.Account)

	w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true})

	infos, err := w.massaClient.GetAccountsInfos([]*account.Account{acc})
	if err != nil {
		return operations.NewImportAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: "Unable to retrieve accounts infos",
			})
	}

	modelWallet, err := newAccountModel(acc)
	if err != nil {
		return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
	}

	modelWallet.CandidateBalance = models.Amount(fmt.Sprint(infos[0].CandidateBalance))
	modelWallet.Balance = models.Amount(fmt.Sprint(infos[0].Balance))

	return operations.NewImportAccountOK().WithPayload(modelWallet)
}
