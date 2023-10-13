package wallet

import (
	"testing"
)

func BenchmarkProtectWallet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iteration()
	}
}

func iteration() {
	password := "password"

	wlt, errGenerate := Generate("nickname", password)
	if errGenerate != nil {
		panic(errGenerate.Err)
	}

	err := wlt.Protect(password)
	if err != nil {
		panic(err)
	}
}
