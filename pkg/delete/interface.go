package delete

// Confirmer is the delete account confirm interface
type Confirmer interface {
	Confirm(walletName string) (clearPassword string, err error)
}
