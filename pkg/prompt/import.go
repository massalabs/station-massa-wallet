package prompt

import (
	"context"
	"fmt"
	"strings"
	"time"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

const TIMEOUT = 5 * time.Minute

func PromptImport(
	prompterApp WalletPrompterInterface,
) (*wallet.Wallet, error) {
	data := &PromptRequestData{
		Msg:  "Import",
		Data: nil,
	}
	prompterApp.PromptRequest(walletapp.Import, data.Msg, data.Data)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	for {
		select {

		case filePath := <-prompterApp.App().WalletFileChan:
			fmt.Println("filePath received: ", filePath)
			if !strings.HasSuffix(filePath, ".yml") {
				errStr := "invalid account file"
				fmt.Println(errStr)
				prompterApp.EmitEvent(walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Data: errStr})
				continue
			}
			wallet, err := wallet.LoadFile(filePath)
			if err != nil {
				errStr := "unable to load account file: " + err.Error()
				fmt.Println(errStr)
				prompterApp.EmitEvent(walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Data: errStr})
				continue
			}
			return &wallet, nil

		case walletInfo := <-prompterApp.App().PrivateKeyChan:
			wallet, err := wallet.Import(walletInfo.Nickname, walletInfo.PrivateKey, walletInfo.Password)
			if err != nil {
				errStr := "unable to import private key: " + err.Error()
				fmt.Println(errStr)
				prompterApp.EmitEvent(walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Data: errStr})
				return nil, fmt.Errorf(errStr)
			}
			return wallet, nil

		case <-prompterApp.App().CtrlChan:
			msg := "Action canceled by user"
			fmt.Println(msg)
			return nil, fmt.Errorf(msg)
		case <-ctxTimeout.Done():
			errStr := "Password prompt reached timeout"
			fmt.Println(errStr)
			prompterApp.EmitEvent(walletapp.PasswordResultEvent,
				walletapp.EventData{Success: false, Data: errStr, Error: "timeoutError"})
			return nil, fmt.Errorf(errStr)
		}
	}
}
