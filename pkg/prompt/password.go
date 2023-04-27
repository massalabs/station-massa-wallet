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
				errStr := "error unprotecting wallet:" + err.Error()
				fmt.Println(errStr)
				prompterApp.EmitEvent(walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Data: errStr})
				continue
			}

			return password, nil
		case <-prompterApp.App().CtrlChan:
			msg := "Action canceled by user"
			fmt.Println(msg)
			return "", fmt.Errorf(msg)
		case <-ctxTimeout.Done():
			errStr := "Password prompt reached timeout"
			fmt.Println(errStr)
			prompterApp.EmitEvent(walletapp.PasswordResultEvent,
				walletapp.EventData{Success: false, Data: errStr, Error: "timeoutError"})
			return "", fmt.Errorf(errStr)
		}
	}
}
