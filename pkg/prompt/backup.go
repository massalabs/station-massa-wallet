package prompt

const (
	YmlFileBackup    = "yml"
	PrivateKeyBackup = "privateKey"
)

func handleBackupMethod(prompterApp WalletPrompterInterface, input interface{}) (bool, error) {
	method, ok := input.(string)
	if !ok {
		return false, InputTypeError(prompterApp)
	}
	if method == YmlFileBackup {
		return false, nil
	}
	if method == PrivateKeyBackup {
		return true, nil
	}

	return false, UserChoiceError(prompterApp)
}
