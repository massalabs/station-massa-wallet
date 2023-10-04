package walletmanager

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/massalabs/station-massa-wallet/pkg/types"
)

func (w *Wallet) NicknameIsUnique(nickname string) error {
	for _, account := range w.Accounts {
		if strings.EqualFold(account.Nickname, nickname) {
			return fmt.Errorf("this account name already exists: %s", nickname)
		}
	}

	return nil
}

func (w *Wallet) AddressIsUnique(address types.Address) error {
	for _, account := range w.Accounts {
		if bytes.Equal(account.Address.Data, address.Data) {
			return fmt.Errorf("this account address already exists")
		}
	}

	return nil
}
