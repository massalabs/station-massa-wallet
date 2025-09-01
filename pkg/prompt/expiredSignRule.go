package prompt

import (
	"fmt"

	"github.com/awnumar/memguard"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

func handleExpiredSignRulePrompt(prompterApp WalletPrompterInterface, input interface{}, acc account.Account) (*walletapp.ExpiredSignRulePromptOutput, bool, error) {
	expiredSignRuleObject, ok := input.(*walletapp.ExpiredSignRulePromptInput)
	if !ok {
		return nil, true, InputTypeError(prompterApp)
	}

	password := memguard.NewBufferFromBytes([]byte(expiredSignRuleObject.Password))

	// password will be destroy in acc.HasAccess, so we need to create a new one.
	returnedPassword := memguard.NewBufferFromBytes([]byte(expiredSignRuleObject.Password))

	if !acc.HasAccess(password) {
		msg := fmt.Sprintf("Invalid password for account %s", acc.Nickname)

		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.WrongPassword})

		return nil, false, fmt.Errorf("%w: %s", utils.ErrWrongPassword, msg)
	}

	return &walletapp.ExpiredSignRulePromptOutput{
		PasswordPromptOutput: walletapp.PasswordPromptOutput{
			Password: returnedPassword,
		},
		ToDelete: expiredSignRuleObject.ToDelete,
	}, false, nil
}
