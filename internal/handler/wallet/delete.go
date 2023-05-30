package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
	"github.com/massalabs/thyra-plugin-wallet/pkg/utils"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func NewDelete(prompterApp prompt.WalletPrompterInterface) operations.DeleteAccountHandler {
	return &walletDelete{prompterApp: prompterApp}
}

type walletDelete struct {
	prompterApp prompt.WalletPrompterInterface
}

// HandleDelete handles a delete request
func (w *walletDelete) Handle(params operations.DeleteAccountParams) middleware.Responder {
	wlt, resp := loadWallet(params.Nickname)
	if resp != nil {
		return resp
	}

	go handleDelete(wlt, w.prompterApp)

	return operations.NewDeleteAccountNoContent()
}

func handleDelete(wlt *wallet.Wallet, prompterApp prompt.WalletPrompterInterface) {
	promptRequest := prompt.PromptRequest{
		Action: walletapp.Delete,
		Msg:    "Delete an account",
	}

	_, err := prompt.WakeUpPrompt(prompterApp, promptRequest, wlt)
	if err != nil {
		return
	}

	if wlt.DeleteFile() != nil {
		errStr := fmt.Sprintf("error deleting wallet: %v", err.Error())
		fmt.Println(errStr)
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.ErrAccountFile})
	}

	prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, CodeMessage: utils.MsgAccountDeleted})
}
