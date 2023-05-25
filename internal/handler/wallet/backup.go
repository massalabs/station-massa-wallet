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

func NewBackupAccount(prompterApp prompt.WalletPrompterInterface) operations.BackupAccountHandler {
	return &walletBackupAccount{prompterApp: prompterApp}
}

type walletBackupAccount struct {
	prompterApp prompt.WalletPrompterInterface
}

func (w *walletBackupAccount) Handle(params operations.BackupAccountParams) middleware.Responder {
	wlt, resp := loadWallet(params.Nickname)
	if resp != nil {
		return resp
	}

	go handleBackup(wlt, w.prompterApp)

	return operations.NewDeleteAccountNoContent()
}

func handleBackup(wlt *wallet.Wallet, prompterApp prompt.WalletPrompterInterface) {
	promptRequest := prompt.PromptRequest{
		Action: walletapp.Backup,
		Msg:    fmt.Sprintf("Backup wallet %s:", wlt.Nickname),
		Data:   nil,
	}

	_, err := prompt.WakeUpPrompt(prompterApp, promptRequest, wlt)
	if err != nil {
		return
	}

	prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, CodeMessage: utils.MsgBackupSuccess})
}
