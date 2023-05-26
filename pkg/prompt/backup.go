package prompt

type BackupMethod string

const (
	YmlFileBackup    BackupMethod = "yml"
	PrivateKeyBackup BackupMethod = "privateKey"
)

func handleBackupMethod(prompterApp WalletPrompterInterface, input interface{}) (*BackupMethod, bool, error) {
	method, ok := input.(BackupMethod)
	if !ok {
		return nil, false, InputTypeError(prompterApp)
	}

	if method == YmlFileBackup {
		return &method, false, nil
	} else {
		return &method, true, nil
	}
}
