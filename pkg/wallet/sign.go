package wallet

import (
	"fmt"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (w *Wallet) UnprotectWalletAskingPassword(walletApp *walletapp.WalletApp) bool {
	msg := fmt.Sprintf("Unprotect wallet %s", w.Nickname)
	walletApp.PromptRequest(walletapp.Sign, msg, interface{}(nil))

	for {
		select {
		case password := <-walletApp.PasswordChan:
			err := w.Unprotect(password)
			if err != nil {
				errStr := "error unprotecting wallet:" + err.Error()
				fmt.Println(errStr)
				runtime.EventsEmit(walletApp.Ctx, walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Data: errStr})
				continue
			}

			runtime.EventsEmit(walletApp.Ctx, walletapp.PasswordResultEvent,
				walletapp.EventData{Success: true, Data: "Unprotect Success"})
			return true
		case <-walletApp.CtrlChan:
			fmt.Println("Action canceled by user")
			return false
		}
	}
}
