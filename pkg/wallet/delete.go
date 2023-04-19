package wallet

import (
	"fmt"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
)

func (w *Wallet) Delete(prompterApp WalletPrompterInterface) {
	msg := fmt.Sprintf("Deleting wallet %s:", w.Nickname)
	prompterApp.PromptRequest(walletapp.Password, msg, interface{}(nil))

	for {
		select {
		case password := <-prompterApp.App().PasswordChan:
			err := w.Unprotect(password)
			if err != nil {
				errStr := "error unprotecting wallet:" + err.Error()
				fmt.Println(errStr)
				prompterApp.EmitEvent(walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Data: errStr})
				continue
			}

			err = w.DeleteFile()
			if err != nil {
				errStr := "error deleting wallet:" + err.Error()
				fmt.Println(errStr)
				prompterApp.EmitEvent(walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Data: errStr})
				continue
			}

			prompterApp.EmitEvent(walletapp.PasswordResultEvent,
				walletapp.EventData{Success: true, Data: "Delete Success"})
			return
		case <-prompterApp.App().CtrlChan:
			fmt.Println("Action canceled by user")
			return
		}
	}
}
