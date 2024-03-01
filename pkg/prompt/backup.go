package prompt

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
		return nil, false, InputTypeError(prompterApp)
	}
}
