package privateKey

type PrivateKeyAsker interface {
	Ask() (privateKey string, err error)
}
