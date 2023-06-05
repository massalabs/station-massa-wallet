package wallet

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/network"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
	"github.com/massalabs/thyra-plugin-wallet/pkg/utils"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
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
	wlt, resp := loadWallet(params.Nickname)
	if resp != nil {
		return resp
	}

	infos, err := w.massaClient.GetAccountsInfos([]wallet.Wallet{*wlt})
	if err != nil {
		return operations.NewDeleteAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallet,
				Message: "Unable to retrieve account infos",
			})
	}

	promptRequest := prompt.PromptRequest{
		Action: walletapp.Delete,
		Msg:    "Delete an account",
		Data: PromptRequestDeleteData{
			Nickname: wlt.Nickname,
			Balance:  fmt.Sprint(infos[0].CandidateBalance),
		},
	}

	_, err = prompt.WakeUpPrompt(w.prompterApp, promptRequest, wlt)
	if err != nil {
		return operations.NewDeleteAccountUnauthorized().WithPayload(
			&models.Error{
				Code:    fmt.Sprint(http.StatusUnauthorized),
				Message: "Unable to unprotect wallet",
			})
	}

	if wlt.DeleteFile() != nil {
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
		walletapp.EventData{Success: true, CodeMessage: utils.MsgAccountDeleted})

	return operations.NewDeleteAccountNoContent()
}
