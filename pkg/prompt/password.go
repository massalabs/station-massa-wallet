package prompt

import (
	"fmt"

	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
)

func handlePasswordPrompt(prompterApp WalletPrompterInterface, input interface{}, wallet *wallet.Wallet) (*string, bool, error) {
	password, ok := input.(string)
	if !ok {
		return nil, false, InputTypeError(prompterApp)
	}

	errUprotect := wallet.Unprotect(password)
	if errUprotect != nil {
		errStr := fmt.Sprintf("%v: %v", UnprotectErr, errUprotect.Err.Error())
		fmt.Println(errStr)
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: errUprotect.CodeErr})
		return nil, true, errUprotect.Err
	}

	return &password, false, nil
}
