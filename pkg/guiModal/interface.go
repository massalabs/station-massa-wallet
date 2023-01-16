package guiModal

// Asker is the password ask interface
type PasswordAsker interface {
	Ask(walletName string) (clearPassword string, err error)
}

type WalletInfoAsker interface {
	WalletInfo() (walletName string, password string, privateKey string, err error)
}
