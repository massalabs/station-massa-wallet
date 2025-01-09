package wallet

import (
	"fmt"

	"github.com/awnumar/memguard"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

func PromptPassword(prompterApp prompt.WalletPrompterInterface, acc *account.Account, promptRequest *prompt.PromptRequest) (*memguard.LockedBuffer, error) {
	promptOutput, err := prompt.WakeUpPrompt(prompterApp, *promptRequest, acc)
	if err != nil {
		return nil, fmt.Errorf("prompting password: %w", err)
	}

	output, ok := promptOutput.(*memguard.LockedBuffer)
	if !ok {
		return nil, fmt.Errorf("prompting password: %s", utils.ErrInvalidInputType.Error())
	}

	prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true})

	return output, nil
}

func PromptForOperation(prompterApp prompt.WalletPrompterInterface, acc *account.Account, promptRequest *prompt.PromptRequest) (*walletapp.SignPromptOutput, error) {
	promptOutput, err := prompt.WakeUpPrompt(prompterApp, *promptRequest, acc)
	if err != nil {
		return nil, fmt.Errorf("prompting password: %w", err)
	}

	output, ok := promptOutput.(*walletapp.SignPromptOutput)
	if !ok {
		return nil, fmt.Errorf("prompting password: %s", utils.ErrInvalidInputType.Error())
	}

	prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true})

	return output, nil
}
