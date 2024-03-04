package prompt

import walletapp "github.com/massalabs/station-massa-wallet/pkg/app"

type BackupMethod string

const (
	YamlFileBackup   BackupMethod = "yaml"
	PrivateKeyBackup BackupMethod = "privateKey"
)

// Returns output, keepListening, error
func handleBackupMethod(prompterApp WalletPrompterInterface, input interface{}) (*BackupMethod, bool, error) {
	inputObject, ok := input.(*walletapp.StringPromptInput)
	if !ok {
		return nil, false, InputTypeError(prompterApp)
	}

	method := inputObject.Message

	switch method {
	case string(YamlFileBackup):
		res := BackupMethod(method)
		return &res, false, nil

	case string(PrivateKeyBackup):
		res := BackupMethod(method)
		return &res, true, nil

	default:
		return nil, false, InputTypeError(prompterApp)
	}
}
