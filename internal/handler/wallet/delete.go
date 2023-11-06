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
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

type PromptRequestDeleteData struct {
	Nickname string
	Balance  string
}

func NewDelete(prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface) operations.DeleteAccountHandler {
	return &walletDelete{prompterApp: prompterApp, massaClient: massaClient}
}

type walletDelete struct {
	prompterApp prompt.WalletPrompterInterface
	massaClient network.NodeFetcherInterface
}

// HandleDelete handles a delete request
func (w *walletDelete) Handle(params operations.DeleteAccountParams) middleware.Responder {
	acc, errResp := loadAccount(w.prompterApp.App().Wallet, params.Nickname)
	if errResp != nil {
		return errResp
	}

	infos, err := w.massaClient.GetAccountsInfos([]*account.Account{acc})
	if err != nil {
		return operations.NewDeleteAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetAccount,
				Message: "Unable to retrieve account infos",
			})
	}

	promptRequest := prompt.PromptRequest{
		Action: walletapp.Delete,
		Msg:    "Delete an account",
		Data: PromptRequestDeleteData{
			Nickname: acc.Nickname,
			Balance:  fmt.Sprint(infos[0].CandidateBalance),
		},
	}

	// Ask for password, validate password.
	_, err = prompt.WakeUpPrompt(w.prompterApp, promptRequest, acc)
	if err != nil {
		return operations.NewDeleteAccountUnauthorized().WithPayload(
			&models.Error{
				Code:    fmt.Sprint(http.StatusUnauthorized),
				Message: fmt.Sprintf("error validating password: %v", err.Error()),
			})
	}

	err = w.prompterApp.App().Wallet.DeleteAccount(acc.Nickname)
	if err != nil {
		errStr := fmt.Sprintf("error deleting wallet: %v", err.Error())
		fmt.Println(errStr)
		w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.ErrAccountFile})

		return operations.NewDeleteAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    utils.ErrAccountFile,
				Message: "Unable to delete account file",
			})
	}

	w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true})

	return operations.NewDeleteAccountNoContent()
}
