package wallet

import (
	"fmt"
	"testing"
)

// func BenchmarkProtectWallet(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		iteration()
// 	}
// }

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

func CreateWallet(b *testing.B) {
	password := "ze password"
	nickname := "ze nickname"
	wlt, err := Generate(nickname, password)
	if err != nil {
		panic(err)
	}

	// test wallet is well formed by importing it
	res, err := ImportFromMnemonic(nickname, wlt.Mnemonic, password)
	if err != nil {
		panic(err)
	}
	if res.Nickname != nickname {
		panic("nickname should be the same")
	}
	if res.Mnemonic != wlt.Mnemonic {
		panic("mnemonic should be the same")
	}
	if res.Address != wlt.Address {
		panic("address should be the same")
	}

	// test wallet is protected by trying to import it without password
	// _, err = Import(wlt.Bytes())
	// if err == nil {
	// 	panic("wallet should be protected")
	// }
	// // test wallet is protected by trying to import it with wrong password
	// _, err = Import(wlt.Bytes())
	// if err == nil {
	// 	panic("wallet should be protected")
	// }

	fmt.Println("WEEEESSSS")
	fmt.Println(wlt)
}
