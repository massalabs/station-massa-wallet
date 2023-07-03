package wallet

import (
	"fmt"
	"io"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
)

// KeyPair represents a pair of private and public keys.
type KeyPair struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

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
		Msg:    wlt.Nickname,
	}
	promptOutput, err := prompt.WakeUpPrompt(w.prompterApp, promptRequest, wlt)
	if err != nil {
		return operations.NewBackupAccountUnauthorized().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: "Unable to backup account",
			})
	}

	// If the user choose to backup the wallet using the yaml file, promptOutput will be a BackupMethod
	// Else, it will be the password
	_, ok := promptOutput.(*prompt.BackupMethod)
	isYmlBackup := ok
	var privateKey string = ""
	var publicKey string = ""

	if isYmlBackup {
		walletErr := w.saveAccountFile(params.Nickname)
		if walletErr != nil {
			w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
				walletapp.EventData{Success: false, CodeMessage: walletErr.CodeErr})
			return operations.NewBackupAccountBadRequest().WithPayload(
				&models.Error{
					Code:    errorSaveAccount,
					Message: walletErr.CodeErr,
				})
		}
	} else {

		privateKey = wlt.GetPrivKey()
		publicKey = wlt.GetPupKey()
	}

	data := KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, CodeMessage: utils.MsgBackupSuccess, Data: data})

	return operations.NewBackupAccountNoContent()
}

func (w *walletBackupAccount) saveAccountFile(nickname string) *wallet.WalletError {
	dstFile, err := w.prompterApp.SelectBackupFilepath(nickname)
	if err != nil {
		return &wallet.WalletError{
			Err:     err,
			CodeErr: utils.ErrAccountFile,
		}
	}

	if dstFile == "" {
		return &wallet.WalletError{
			Err:     fmt.Errorf("no file selected"),
			CodeErr: utils.ActionCanceled,
		}
	}

	// Create the destination file
	destination, err := os.Create(dstFile)
	if err != nil {
		return &wallet.WalletError{
			Err:     err,
			CodeErr: utils.ErrAccountFile,
		}
	}
	defer destination.Close()

	srcFile, err := wallet.FilePath(nickname)
	if err != nil {
		return &wallet.WalletError{
			Err:     err,
			CodeErr: utils.ErrAccountFile,
		}
	}
	source, err := os.Open(srcFile)
	if err != nil {
		return &wallet.WalletError{
			Err:     err,
			CodeErr: utils.ErrAccountFile,
		}
	}
	defer source.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return &wallet.WalletError{
			Err:     err,
			CodeErr: utils.ErrAccountFile,
		}
	}

	return nil
}
