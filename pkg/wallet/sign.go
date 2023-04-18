package wallet

import (
	"fmt"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
)

func (w *Wallet) UnprotectWalletAskingPassword(prompterApp WalletPrompterInterface) bool {
	msg := fmt.Sprintf("Unprotect wallet %s", w.Nickname)
	prompterApp.PromptRequest(walletapp.Sign, msg, interface{}(nil))

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

			prompterApp.EmitEvent(walletapp.PasswordResultEvent,
				walletapp.EventData{Success: true, Data: "Unprotect Success"})

			return true
		case <-prompterApp.App().CtrlChan:
			fmt.Println("Action canceled by user")
			return false
		}
	}
}
