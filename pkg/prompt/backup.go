package prompt

type BackupMethod string

const (
	YmlFileBackup    BackupMethod = "yml"
	PrivateKeyBackup BackupMethod = "privateKey"
)

func handleBackupMethod(prompterApp WalletPrompterInterface, input interface{}) (*BackupMethod, bool, error) {
	method, ok := input.(string)
	if !ok {
		return nil, false, InputTypeError(prompterApp)
	}
	switch method {
	case string(YmlFileBackup):
		res := BackupMethod(method)
		return &res, false, nil
	case string(PrivateKeyBackup):
		res := BackupMethod(method)
		return &res, true, nil
	default:
		return nil, false, InputTypeError(prompterApp)
	}
}
