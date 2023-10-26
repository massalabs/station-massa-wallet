package prompt

import (
	"fmt"

	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

func handleImportPrompt(prompterApp WalletPrompterInterface, input interface{}) (*account.Account, bool, error) {
	filePath, ok := input.(string)
	if ok {
		return handleImportFile(prompterApp, filePath)
	}

	walletInfo, ok := input.(walletapp.ImportFromPKey)
	if ok {
		return handleImportPrivateKey(prompterApp, walletInfo)
	}

	return nil, false, InputTypeError(prompterApp)
}

func handleImportFile(prompterApp WalletPrompterInterface, filePath string) (*account.Account, bool, error) {
	wallet := prompterApp.App().Wallet

	acc, err := wallet.Load(filePath)
	if err != nil {
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.WailsErrorCode(err)})

		return nil, false, fmt.Errorf("unable to load account file: %w", err)
	}

	err = wallet.AddAccount(acc, true)
	if err != nil {
		msg := fmt.Sprintf("failed to add account: %v", err)
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.WailsErrorCode(err)})

		return nil, false, fmt.Errorf(msg)
	}

	return acc, false, nil
}

func handleImportPrivateKey(prompterApp WalletPrompterInterface, walletInfo walletapp.ImportFromPKey) (*account.Account, bool, error) {
	acc, err := account.NewFromPrivateKey(walletInfo.Password, walletInfo.Nickname, walletInfo.PrivateKey)
	if err != nil {
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.WailsErrorCode(err)})

		return nil, false, fmt.Errorf("unable to import private key: %w", err)
	}

	wallet := prompterApp.App().Wallet

	err = wallet.AddAccount(acc, true)
	if err != nil {
		msg := fmt.Sprintf("failed to add account: %v", err)
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.WailsErrorCode(err)})

		return nil, false, fmt.Errorf(msg)
	}

	return acc, false, nil
}
