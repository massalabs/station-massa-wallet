package privateKey

type Asker interface {
	Ask() (privateKey string, err error)
}
