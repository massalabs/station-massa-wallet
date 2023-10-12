package prompt

import (
	"errors"
	"fmt"

	"github.com/awnumar/memguard"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

// handlePasswordPrompt returns the password as a LockedBuffer, or an error if the input is not a string.
func handlePasswordPrompt(prompterApp WalletPrompterInterface, input interface{}, acc *account.Account) (*memguard.LockedBuffer, bool, error) {
	password, ok := input.(string)
	if !ok {
		return nil, false, InputTypeError(prompterApp)
	}

	guardedPassword := memguard.NewBufferFromBytes([]byte(password))

	// guardedPassword will be destroy in acc.PasswordIsValid, so we need to create a new one.
	guardedPasswordReturned := memguard.NewBufferFromBytes([]byte(password))

	if acc != nil && !acc.PasswordIsValid(guardedPassword) {
		msg := fmt.Sprintf("Invalid password for account %s", acc.Nickname)
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.WrongPassword})
		return nil, true, errors.New(msg)
	}

	return guardedPasswordReturned, false, nil
}
