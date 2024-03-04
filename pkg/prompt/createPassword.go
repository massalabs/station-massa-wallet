package prompt

import (
	"strings"

	"github.com/awnumar/memguard"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
)

func handleNewPasswordPrompt(prompterApp WalletPrompterInterface, input interface{}) (*memguard.LockedBuffer, bool, error) {
	inputObject, ok := input.(*walletapp.StringPromptInput)
	if !ok {
		return nil, false, InputTypeError(prompterApp)
	}

	trimmedPassword := strings.TrimSpace(inputObject.Message)

	return memguard.NewBufferFromBytes([]byte(trimmedPassword)), false, nil
}
