package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/network"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
	"github.com/massalabs/thyra-plugin-wallet/pkg/utils"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

type PromptRequestDeleteDate struct {
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

	go handleDelete(wlt, w)

	return operations.NewDeleteAccountNoContent()
}

func handleDelete(wlt *wallet.Wallet, w *walletDelete) {
	infos, err := w.massaClient.GetAccountsInfos([]wallet.Wallet{*wlt})
	if err != nil {
		w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.ErrNetwork})
		return
	}

	promptRequest := prompt.PromptRequest{
		Action: walletapp.Delete,
		Msg:    "Delete an account",
		Data: PromptRequestDeleteDate{
			Nickname: wlt.Nickname,
			Balance:  fmt.Sprint(infos[0].CandidateBalance),
		},
	}

	_, err = prompt.WakeUpPrompt(w.prompterApp, promptRequest, wlt)
	if err != nil {
		return
	}

	if wlt.DeleteFile() != nil {
		errStr := fmt.Sprintf("error deleting wallet: %v", err.Error())
		fmt.Println(errStr)
		w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.ErrAccountFile})
	}

	w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, CodeMessage: utils.MsgAccountDeleted})
}
