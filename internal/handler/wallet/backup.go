package wallet

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/awnumar/memguard"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
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
	acc, errResp := loadAccount(w.prompterApp.App().Wallet, params.Nickname)
	if errResp != nil {
		return errResp
	}

	promptRequest := prompt.PromptRequest{
		Action: walletapp.Backup,
		Msg:    acc.Nickname,
	}

	promptOutput, err := prompt.WakeUpPrompt(w.prompterApp, promptRequest, acc)
	if err != nil {
		return newErrorResponse("Unable to backup account", errorCanceledAction, http.StatusUnauthorized)
	}

	// If the user choose to backup the wallet using the yaml file, promptOutput will be a BackupMethod
	// Else, it will be the password
	_, ok := promptOutput.(*prompt.BackupMethod)
	isYmlBackup := ok
	var privateKey string = ""
	var publicKey string = ""

	if isYmlBackup {
		err = w.saveAccountFile(params.Nickname)
		if err != nil {
			w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
				walletapp.EventData{Success: false, CodeMessage: utils.WailsErrorCode(err)})

			return newErrorResponse(utils.ErrAccountFile, errorSaveAccount, http.StatusBadRequest)
		}
	} else {
		password, _ := promptOutput.(*memguard.LockedBuffer)
		guardedPrivateKey, err := acc.PrivateKeyTextInClear(password)
		//nolint:wsl
		if err != nil {
			return newErrorResponse(err.Error(), errorGetWallets, http.StatusInternalServerError)
		}

		defer guardedPrivateKey.Destroy()
		defer password.Destroy()

		privateKey = string(guardedPrivateKey.Bytes())

		publicKeyBytes, err := acc.PublicKey.MarshalText()
		if err != nil {
			return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
		}

		publicKey = string(publicKeyBytes)
	}

	data := KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, Data: data})

	return operations.NewBackupAccountNoContent()
}

func (w *walletBackupAccount) saveAccountFile(nickname string) error {
	dstFile, err := w.prompterApp.SelectBackupFilepath(nickname)
	if err != nil {
		return err
	}

	if dstFile == "" {
		return fmt.Errorf("no file selected: %w", utils.ErrActionCanceled)
	}

	// Create the destination file
	destination, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer destination.Close()

	srcFile, err := w.prompterApp.App().Wallet.AccountPath(nickname)
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
