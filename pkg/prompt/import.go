package prompt

import (
	"fmt"
	"strings"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/utils"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
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
	if !strings.HasSuffix(filePath, ".yml") {
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, Error: utils.ErrAccountFile})
		return nil, true, fmt.Errorf(InvalidAccountFileErr)
	}
	wallet, loadErr := wallet.LoadFile(filePath)
	if loadErr != nil {
		errStr := fmt.Sprintf("%v: %v", AccountLoadErr, loadErr.Err.Error())
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, Error: loadErr.CodeErr})
		return nil, true, fmt.Errorf(errStr)

	}

	return &wallet, false, nil
}

func handleImportPrivateKey(prompterApp WalletPrompterInterface, walletInfo walletapp.ImportFromPKey) (*wallet.Wallet, bool, error) {
	wallet, importErr := wallet.Import(walletInfo.Nickname, walletInfo.PrivateKey, walletInfo.Password)
	if importErr != nil {
		errStr := fmt.Sprintf("%v: %v", ImportPrivateKeyErr, importErr.Err)
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, Error: importErr.CodeErr})
		return nil, false, fmt.Errorf(errStr)
	}

	return wallet, false, nil
}
