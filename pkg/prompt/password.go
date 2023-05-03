package prompt

import (
	"context"
	"fmt"
	"time"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func PromptPassword(
	prompterApp WalletPrompterInterface,
	wallet *wallet.Wallet,
	action walletapp.PromptRequest,
	data *PromptRequestData,
) (string, error) {
	prompterApp.PromptRequest(action, data.Msg, data.Data)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	for {
		select {
		case password := <-prompterApp.App().PasswordChan:
			err := wallet.Unprotect(password)
			if err != nil {
				errStr := UnprotectErr + ": " + err.Error()
				fmt.Println(errStr)
				prompterApp.EmitEvent(walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Data: errStr})
				continue
			}

			return password, nil
		case <-prompterApp.App().CtrlChan:
			fmt.Println(ActionCanceledErr)
			return "", fmt.Errorf(ActionCanceledErr)
		case <-ctxTimeout.Done():
			fmt.Println(TimeoutErr)
			prompterApp.EmitEvent(walletapp.PasswordResultEvent,
				walletapp.EventData{Success: false, Data: TimeoutErr, Error: "timeoutError"})
			return "", fmt.Errorf(TimeoutErr)
		}
	}
}
