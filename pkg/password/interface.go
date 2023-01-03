package password

// Asker is the password ask interface
type Asker interface {
	Ask(walletName string) (clearPassword string, err error)
}
