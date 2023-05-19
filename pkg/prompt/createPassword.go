package prompt

import (
	"context"
	"fmt"
	"strings"
	"time"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
)

func PromptCreatePassword(
	prompterApp WalletPrompterInterface,
	nickname string,
) (string, error) {
	data := &PromptRequestData{
		Msg:  "Define a password",
		Data: nil,
	}
	prompterApp.PromptRequest(walletapp.NewPassword, data.Msg, data.Data)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	for {
		select {
		case password := <-prompterApp.App().PasswordChan:
			trimmedPassword := strings.TrimSpace(password)
			if len(trimmedPassword) < PASSWORD_MIN_LENGTH {
				// TODO implement password strength check
				errStr := fmt.Sprintf(PasswordLengthErr, PASSWORD_MIN_LENGTH)
				fmt.Println(errStr)
				prompterApp.EmitEvent(walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Data: errStr})
				continue
			}

			return trimmedPassword, nil

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
