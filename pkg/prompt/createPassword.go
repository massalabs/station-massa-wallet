package prompt

import (
	"fmt"
	"strings"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
)

func handleNewPasswordPrompt(prompterApp WalletPrompterInterface, input interface{}) (*string, bool, error) {
	password, ok := input.(string)
	if !ok {
		return nil, false, InputTypeError(prompterApp)
	}

	trimmedPassword := strings.TrimSpace(password)
	if len(trimmedPassword) < PASSWORD_MIN_LENGTH {
		// TODO implement password strength check
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, Data: passwordLengthErr})
		return nil, true, fmt.Errorf(passwordLengthErr)
	}

	return &trimmedPassword, false, nil
}
