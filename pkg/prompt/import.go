package prompt

import (
	"fmt"
	"strings"

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
	if !strings.HasSuffix(filePath, ".yaml") && !strings.HasSuffix(filePath, ".yml") {
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.ErrInvalidFileExtension})
		return nil, false, fmt.Errorf(InvalidAccountFileErr)
	}

	wallet := prompterApp.App().WalletManager

	acc, err := wallet.Load(filePath)
	if err != nil {
		msg := fmt.Sprintf("%v: %v", AccountLoadErr, err)
		code := utils.WailsErrorCode(err)
		fmt.Println("code is: ", code)
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: code})
		return nil, false, fmt.Errorf(msg)
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
		errStr := fmt.Sprintf("%v: %v", ImportPrivateKeyErr, err)
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.WailsErrorCode(err)})
		return nil, false, fmt.Errorf(errStr)
	}

	wallet := prompterApp.App().WalletManager

	err = wallet.AddAccount(acc, true)
	if err != nil {
		msg := fmt.Sprintf("failed to add account: %v", err)
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.WailsErrorCode(err)})
		return nil, false, fmt.Errorf(msg)
	}

	return acc, false, nil
}
