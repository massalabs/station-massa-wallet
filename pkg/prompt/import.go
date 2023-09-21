package prompt

import (
	"fmt"
	"strings"

	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

func handleImportPrompt(prompterApp WalletPrompterInterface, input interface{}) (*wallet.Wallet, bool, error) {
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

func handleImportFile(prompterApp WalletPrompterInterface, filePath string) (*wallet.Wallet, bool, error) {
	if !strings.HasSuffix(filePath, ".yaml") && !strings.HasSuffix(filePath, ".yml") {
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.ErrInvalidFileExtension})
		return nil, false, fmt.Errorf(InvalidAccountFileErr)
	}
	acc, loadErr := wallet.LoadFile(filePath)
	if loadErr != nil {
		errStr := fmt.Sprintf("%v: %v", AccountLoadErr, loadErr.Err)
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: loadErr.CodeErr})
		return nil, false, fmt.Errorf(errStr)
	}

	// Validate nickname
	if !account.NicknameIsValid(acc.Nickname) {
		errorCode := utils.ErrInvalidNickname
		fmt.Printf("error while importing: invalid nickname: '%s'\n", acc.Nickname)
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: errorCode})
		return nil, false, fmt.Errorf(errorCode)
	}

	// Validate nickname uniqueness
	err := wallet.NicknameIsUnique(acc.Nickname)
	if err != nil {
		errorCode := utils.ErrDuplicateNickname
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: errorCode})
		return nil, false, fmt.Errorf(errorCode)
	}

	// Validate unique private key
	err = wallet.AddressIsUnique(acc.Address)
	if err != nil {
		errorCode := utils.ErrDuplicateKey
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: errorCode})
		return nil, false, fmt.Errorf(errorCode)
	}

	return &acc, false, nil
}

func handleImportPrivateKey(prompterApp WalletPrompterInterface, walletInfo walletapp.ImportFromPKey) (*wallet.Wallet, bool, error) {
	wallet, importErr := wallet.Import(walletInfo.Nickname, walletInfo.PrivateKey, walletInfo.Password)
	if importErr != nil {
		errStr := fmt.Sprintf("%v: %v", ImportPrivateKeyErr, importErr.Err)
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: importErr.CodeErr})
		return nil, false, fmt.Errorf(errStr)
	}

	return wallet, false, nil
}
