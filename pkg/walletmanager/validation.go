package walletmanager

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/massalabs/station-massa-wallet/pkg/types"
)

var (
	ErrAddressNotUnique  = errors.New("this account address already exists")
	ErrNicknameNotUnique = errors.New("this account nickname already exists")
)

func (w *Wallet) NicknameIsUnique(nickname string) error {
	for _, account := range w.accounts {
		if strings.EqualFold(account.Nickname, nickname) {
			return fmt.Errorf("%w: %s", ErrNicknameNotUnique, nickname)
		}
	}

	return nil
}

func (w *Wallet) AddressIsUnique(address types.Address) error {
	for _, account := range w.accounts {
		if bytes.Equal(account.Address.Data, address.Data) {
			return ErrAddressNotUnique
		}
	}

	return nil
}
