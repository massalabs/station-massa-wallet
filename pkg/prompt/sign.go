package prompt

import (
	"fmt"
	"strconv"

	"github.com/awnumar/memguard"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

// handleSignPrompt returns the password as a LockedBuffer, or an error if the input is not a string.
func handleSignPrompt(prompterApp WalletPrompterInterface, input interface{}, acc *account.Account) (*walletapp.SignPromptOutput, bool, error) {
	inputObject, ok := input.(walletapp.SignPromptInput)
	if !ok {
		return nil, false, InputTypeError(prompterApp)
	}

	fees, err := strconv.ParseUint(inputObject.Fees, 10, 64)
	if err != nil {
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.InvalidFees})

		return nil, true, fmt.Errorf("failed to parse fees: %w", err)
	}

	inputString := inputObject.Password

	password := memguard.NewBufferFromBytes([]byte(inputString))

	// password will be destroy in acc.HasAccess, so we need to create a new one.
	passwordReturned := memguard.NewBufferFromBytes([]byte(inputString))

	if acc != nil && !acc.HasAccess(password) {
		msg := fmt.Sprintf("Invalid password for account %s", acc.Nickname)

		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.WrongPassword})

		return nil, true, fmt.Errorf("%w: %s", utils.ErrWrongPassword, msg)
	}

	output := &walletapp.SignPromptOutput{
		Password: passwordReturned,
		Fees:     fees,
	}

	return output, false, nil
}
