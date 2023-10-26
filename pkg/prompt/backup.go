package prompt

import (
	"fmt"

	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
)

type BackupMethod string

const (
	YamlFileBackup   BackupMethod = "yaml"
	PrivateKeyBackup BackupMethod = "privateKey"
)

// Returns output, keepListening, error
func handleBackupMethod(prompterApp WalletPrompterInterface, input interface{}) (*BackupMethod, bool, error) {
	method, ok := input.(string)
	if !ok {
		return nil, false, InputTypeError(prompterApp)
	}

	switch method {
	case string(YamlFileBackup):
		res := BackupMethod(method)
		return &res, false, nil

	case string(PrivateKeyBackup):
		res := BackupMethod(method)
		return &res, true, nil

	default:
		prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.ErrPromptInputType})
		return nil, false, fmt.Errorf("invalid backup method: %s", method)
	}
}
