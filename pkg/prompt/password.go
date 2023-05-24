package prompt

import (
	"fmt"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func handlePasswordPrompt(prompterApp WalletPrompterInterface, input interface{}, wallet *wallet.Wallet) (*string, bool, error) {
	password, ok := input.(string)
	if !ok {
		return nil, false, InputTypeError(prompterApp)
	}

	err := wallet.Unprotect(password)
	if err != nil {
		errStr := fmt.Sprintf("%v: %v", UnprotectErr, err.Error())
		fmt.Println(errStr)
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, Data: errStr})
		return nil, true, err
	}

	return &password, false, nil
}
