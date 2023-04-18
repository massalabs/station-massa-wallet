package wallet

import (
	"fmt"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (w *Wallet) Delete(walletApp *walletapp.WalletApp) {
	msg := fmt.Sprintf("Deleting wallet %s:", w.Nickname)
	walletApp.PromptRequest(walletapp.Password, msg, interface{}(nil))

	for {
		select {
		case password := <-walletApp.PasswordChan:
			fmt.Println("password received" + password)
			err := w.Unprotect(password)
			if err != nil {
				errStr := "error unprotecting wallet:" + err.Error()
				fmt.Println(errStr)
				runtime.EventsEmit(walletApp.Ctx, walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Data: errStr})
				continue
			}

			err = w.DeleteFile()
			if err != nil {
				errStr := "error deleting wallet:" + err.Error()
				fmt.Println(errStr)
				runtime.EventsEmit(walletApp.Ctx, walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Data: errStr})
				continue
			}

			runtime.EventsEmit(walletApp.Ctx, walletapp.PasswordResultEvent,
				walletapp.EventData{Success: true, Data: "Delete Success"})
		case <-walletApp.CtrlChan:
			fmt.Println("Action canceled by user")
			return
		}
	}
}
