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
	wlt, err := Generate("nickname", password)
	if err != nil {
		panic(err)
	}
	err = wlt.Protect(password)
	if err != nil {
		panic(err)
	}
}
