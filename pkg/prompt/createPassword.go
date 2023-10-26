package prompt

import (
	"strings"

	"github.com/awnumar/memguard"
)

func handleNewPasswordPrompt(prompterApp WalletPrompterInterface, input interface{}) (*memguard.LockedBuffer, bool, error) {
	password, ok := input.(string)
	if !ok {
		return nil, false, InputTypeError(prompterApp)
	}

	trimmedPassword := strings.TrimSpace(password)

	return memguard.NewBufferFromBytes([]byte(trimmedPassword)), false, nil
}
