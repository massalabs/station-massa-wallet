package wallet

import (
	"fmt"
	"io"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
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

	promptRequest := prompt.PromptRequest{
		Action: walletapp.Backup,
		Msg:    fmt.Sprintf("Backup wallet %s:", wlt.Nickname),
	}
	promptOutput, err := prompt.WakeUpPrompt(w.prompterApp, promptRequest, wlt)
	if err != nil {
		return operations.NewBackupAccountUnauthorized().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: "Unable to backup account",
			})
	}

	// If the user choose to backup the wallet using the yml file, promptOutput will be a BackupMethod
	// Else, it will be the password
	_, ok := promptOutput.(*prompt.BackupMethod)
	isYmlBackup := ok
	var privateKey string = ""

	if isYmlBackup {
		err = w.saveAccountFile(params.Nickname)
		if err != nil {
			return operations.NewBackupAccountBadRequest().WithPayload(
				&models.Error{
					Code:    errorSaveAccount,
					Message: "Unable to backup account file",
				})
		}
	} else {
		// for private key backup, send the private key to the wails frontend
		privateKey = wlt.GetPrivKey()
	}

	w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, CodeMessage: utils.MsgBackupSuccess, Data: privateKey})

	return operations.NewBackupAccountNoContent()
}

func (w *walletBackupAccount) saveAccountFile(nickname string) error {
	dstFile, err := w.prompterApp.SelectBackupFilepath(nickname)
	if err != nil {
		return err
	}

	// Create the destination file
	destination, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer destination.Close()

	srcFile, err := wallet.FilePath(nickname)
	if err != nil {
		return err
	}
	source, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}
	return nil
}
