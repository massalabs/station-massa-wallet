package walletmanager

import (
	"fmt"

	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

type Wallet struct {
	Accounts map[string]*account.Account // Mapping from nickname to account
}

func New() *Wallet {
	return &Wallet{
		Accounts: make(map[string]*account.Account),
	}
}

// Add an account into the wallet
func (w *Wallet) AddAccount(acc *account.Account) error {
	if acc == nil {
		return fmt.Errorf("account is nil")
	}

	err := w.Persist(*acc)
	if err != nil {
		return fmt.Errorf("persisting account: %w", err)
	}

	if w.Accounts[acc.Nickname] == nil {
		w.Accounts[acc.Nickname] = acc
	}

	return nil
}

// Get an account from the wallet by nickname
func (w *Wallet) GetAccount(nickname string) (*account.Account, error) {
	if w.Accounts[nickname] != nil {
		return w.Accounts[nickname], nil
	}

	accountPath, err := w.accountPath(nickname)
	if err != nil {
		return nil, fmt.Errorf("getting account path: %w", err)
	}
	acc, err := w.Load(accountPath)
	if err != nil {
		return nil, fmt.Errorf("loading account: %w", err)
	}

	w.Accounts[nickname] = acc

	return acc, nil
}

// Delete an account from the wallet
func (w *Wallet) DeleteAccount(nickname string) error {
	if w.Accounts[nickname] == nil {
		return fmt.Errorf("account not found")
	}

	accountPath, err := w.accountPath(nickname)
	if err != nil {
		return fmt.Errorf("getting account path: %w", err)
	}

	err = w.deleteFile(accountPath)
	if err != nil {
		return fmt.Errorf("deleting account file: %w", err)
	}

	delete(w.Accounts, nickname)

	return nil
}

// Get the number of accounts in the wallet
func (w *Wallet) GetAccountCount() int {
	return len(w.Accounts)
}
