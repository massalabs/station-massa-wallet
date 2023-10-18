package prompt

import (
	"fmt"

	"github.com/awnumar/memguard"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

// handlePasswordPrompt returns the password as a LockedBuffer, or an error if the input is not a string.
func handlePasswordPrompt(prompterApp WalletPrompterInterface, input interface{}, acc *account.Account) (*memguard.LockedBuffer, bool, error) {
	inputString, ok := input.(string)
	if !ok {
		return nil, false, InputTypeError(prompterApp)
	}

	password := memguard.NewBufferFromBytes([]byte(inputString))

	// password will be destroy in acc.HasAccess, so we need to create a new one.
	passwordReturned := memguard.NewBufferFromBytes([]byte(inputString))

	if acc != nil && !acc.HasAccess(password) {
		msg := fmt.Sprintf("Invalid password for account %s", acc.Nickname)

		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.WrongPassword})

		return nil, true, fmt.Errorf("%w: %s", utils.ErrWrongPassword, msg)
	}

	return passwordReturned, false, nil
}
