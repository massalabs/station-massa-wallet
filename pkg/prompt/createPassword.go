package prompt

import (
	"fmt"
	"strings"

	"github.com/awnumar/memguard"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
)

func handleNewPasswordPrompt(prompterApp WalletPrompterInterface, input interface{}) (*memguard.LockedBuffer, bool, error) {
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

	guardedPassword := memguard.NewBufferFromBytes([]byte(trimmedPassword))

	return guardedPassword, false, nil
}
