package walletmanager

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/awnumar/memguard"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station/pkg/logger"
)

var AccountNotFoundError = errors.New("account not found")

func ErrorAccountNotFound(nickname string) error {
	return fmt.Errorf("account '%s' not found", nickname)
}

type WalletError struct {
	Err     error
	CodeErr string // Sentinel error code from utils package, can be used as a translation key.
}

type Wallet struct {
	Accounts                map[string]*account.Account // Mapping from nickname to account
	InvalidAccountNicknames []string                    // List of invalid account nicknames
}

func New() (*Wallet, error) {
	wallet := &Wallet{
		Accounts: make(map[string]*account.Account),
	}

	err := MigrateWallet()
	if err != nil {
		logger.Errorf("migrating wallet: %s", err)
	}

	err = wallet.discover()
	if err != nil {
		return nil, fmt.Errorf("discovering accounts: %s\n", err)
	}

	return wallet, nil
}

func (w *Wallet) discover() error {
	accountsPath, err := Path()
	if err != nil {
		return fmt.Errorf("getting accounts path: %w", err)
	}

	files, err := os.ReadDir(accountsPath)
	if err != nil {
		return fmt.Errorf("reading accounts path: %w", err)
	}

	for _, f := range files {
		fileName := f.Name()
		filePath := path.Join(accountsPath, fileName)

		if strings.HasPrefix(fileName, "wallet_") && strings.HasSuffix(fileName, ".yaml") {
			acc, err := w.Load(filePath)
			if err != nil {
				nickname := w.nicknameFromFilePath(filePath)
				logger.Infof("invalid account found: %s", nickname)
				w.InvalidAccountNicknames = append(w.InvalidAccountNicknames, nickname)

				continue
			}

			err = w.AddAccount(acc, false)
			if err != nil {
				return fmt.Errorf("adding account: %w", err)
			}
		}
	}

	return nil
}

// Add an account into the wallet
func (w *Wallet) AddAccount(acc *account.Account, persist bool) error {
	if acc == nil {
		return fmt.Errorf("account is nil")
	}

	// Validate nickname uniqueness
	err := w.NicknameIsUnique(acc.Nickname)
	if err != nil {
		return fmt.Errorf("nickname is not unique: %w", err)
	}

	// Validate unique private key
	err = w.AddressIsUnique(acc.Address)
	if err != nil {
		return fmt.Errorf("address is not unique: %w", err)
	}

	if persist {
		err = w.Persist(*acc)
		if err != nil {
			return fmt.Errorf("persisting account: %w", err)
		}
	}

	if w.Accounts[acc.Nickname] == nil {
		w.Accounts[acc.Nickname] = acc
	}

	return nil
}

func (w *Wallet) GenerateAccount(password *memguard.LockedBuffer, nickname string) (*account.Account, error) {
	acc, err := account.NewGenerated(password, nickname)
	if err != nil {
		return nil, fmt.Errorf("generating account: %w", err)
	}

	err = w.AddAccount(acc, true)
	if err != nil {
		return nil, fmt.Errorf("adding account: %w", err)
	}

	return acc, nil
}

// Get an account from the wallet by nickname
func (w *Wallet) GetAccount(nickname string) (*account.Account, error) {
	if w.Accounts[nickname] != nil {
		return w.Accounts[nickname], nil
	}

	accountPath, err := w.AccountPath(nickname)
	if err != nil {
		return nil, fmt.Errorf("getting account path: %w", err)
	}

	acc, err := w.Load(accountPath)
	if err != nil {
		logger.Errorf("loading account: %s", err)
		return nil, AccountNotFoundError
	}

	err = w.AddAccount(acc, false)
	if err != nil {
		return nil, fmt.Errorf("adding account: %w", err)
	}

	return acc, nil
}

// Delete an account from the wallet
func (w *Wallet) DeleteAccount(nickname string) error {
	if w.Accounts[nickname] == nil {
		return fmt.Errorf("account not found")
	}

	accountPath, err := w.AccountPath(nickname)
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
