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
				fmt.Println(InvalidAccountFileErr)
				prompterApp.EmitEvent(walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Error: InvalidAccountFileErr})
				continue
			}
			wallet, loadErr := wallet.LoadFile(filePath)
			if loadErr != nil {
				errStr := fmt.Sprintf("%v: %v", AccountLoadErr, loadErr.Err.Error())
				fmt.Println(errStr)
				prompterApp.EmitEvent(walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Error: loadErr.CodeErr})
				continue
			}
			return &wallet, nil

		case walletInfo := <-prompterApp.App().PrivateKeyChan:
			wallet, importErr := wallet.Import(walletInfo.Nickname, walletInfo.PrivateKey, walletInfo.Password)
			if importErr != nil {
				errStr := fmt.Sprintf("%v: %v", ImportPrivateKeyErr, importErr.Err)
				fmt.Println(errStr)
				prompterApp.EmitEvent(walletapp.PasswordResultEvent,
					walletapp.EventData{Success: false, Error: importErr.CodeErr})
				return nil, fmt.Errorf(errStr)
			}
			return wallet, nil

		case <-prompterApp.App().CtrlChan:
			fmt.Println(ActionCanceledErr)
			return nil, fmt.Errorf(ActionCanceledErr)
		case <-ctxTimeout.Done():
			fmt.Println(TimeoutErr)
			prompterApp.EmitEvent(walletapp.PasswordResultEvent,
				walletapp.EventData{Success: false, Data: TimeoutErr, Error: "timeoutError"})
			return nil, fmt.Errorf(TimeoutErr)
		}
	}
}
