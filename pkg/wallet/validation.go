package wallet

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/massalabs/station-massa-wallet/pkg/types"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

var (
	ErrAddressNotUnique  = errors.New("this account address already exists")
	ErrNicknameNotUnique = errors.New("this account nickname already exists")
)

func (w *Wallet) NicknameIsUnique(nickname string) error {
	var duplicateNickname string

	w.accounts.Range(func(_, acc interface{}) bool {
		account, ok := acc.(*account.Account)
		if !ok {
			return true
		}

		if strings.EqualFold(account.Nickname, nickname) {
			duplicateNickname = nickname
			return false
		}

		return true
	})

	if duplicateNickname != "" {
		return fmt.Errorf("%w: %s", ErrNicknameNotUnique, duplicateNickname)
	}

	return nil
}

func (w *Wallet) AddressIsUnique(address *types.Address) error {
	duplicateAddress := false

	w.accounts.Range(func(_, acc interface{}) bool {
		account, ok := acc.(*account.Account)
		if !ok {
			return true
		}

		if bytes.Equal(account.Address.Data, address.Data) {
			duplicateAddress = true
			return false
		}

		return true
	})

	if duplicateAddress {
		return ErrAddressNotUnique
	}

	return nil
}
