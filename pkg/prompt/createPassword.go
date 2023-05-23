package prompt

import (
	"fmt"
	"strings"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/utils"
)

func handleNewPasswordPrompt(prompterApp WalletPrompterInterface, input interface{}) (*string, bool, error) {
	password, ok := input.(string)
	if !ok {
		return nil, false, InputTypeError(prompterApp)
	}

	trimmedPassword := strings.TrimSpace(password)
	if len(trimmedPassword) < PASSWORD_MIN_LENGTH {
		// TODO implement/refactor password strength check
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.ErrInvalidPassword})
		return nil, true, fmt.Errorf(passwordLengthErr)
	}

	return &trimmedPassword, false, nil
}
