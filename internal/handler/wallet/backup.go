package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
	"github.com/massalabs/thyra-plugin-wallet/pkg/utils"
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

	promptRequest := prompt.PromptRequest{
		Action: walletapp.Backup,
		Msg:    fmt.Sprintf("Backup wallet %s:", wlt.Nickname),
		Data:   nil,
	}
	promptOutput, err := prompt.WakeUpPrompt(w.prompterApp, promptRequest, wlt)
	if err != nil {
		return operations.NewBackupAccountUnauthorized().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: "Unable to backup wallet",
			})
	}

	// If the user choose to backup the wallet using the yml file, promptOutput will be a BackupMethod
	// Else, it will be the password
	_, ok := promptOutput.(*prompt.BackupMethod)
	if !ok {
		// for private key backup, send the private key to the wails frontend
		w.prompterApp.EmitEvent(walletapp.PromptDataEvent,
			walletapp.EventData{Data: wlt.GetPrivKey()})
	}

	w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, CodeMessage: utils.MsgBackupSuccess})
	return operations.NewBackupAccountNoContent()
}
