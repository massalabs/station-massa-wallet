package password

// Asker is the password ask interface
type PasswordAsker interface {
	Ask(walletName string) (clearPassword string, err error)
}
